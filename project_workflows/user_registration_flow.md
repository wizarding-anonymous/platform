# User Registration and Initial Profile Setup Workflow

This diagram illustrates the sequence of interactions between services when a new user registers on the platform.

```mermaid
sequenceDiagram
    actor User
    participant ClientApp as Client Application
    participant APIGW as API Gateway
    participant AuthSvc as Auth Service
    participant KafkaBus as Kafka Message Bus
    participant AccountSvc as Account Service
    participant NotificationSvc as Notification Service

    User->>ClientApp: Submits registration form (email, password, username)
    ClientApp->>APIGW: POST /api/v1/auth/register (payload)
    APIGW->>AuthSvc: Forward POST /register (payload: username, email, password)

    AuthSvc->>AuthSvc: Validate input (email format, password strength, username uniqueness)
    alt Username/Email already exists
        AuthSvc-->>APIGW: HTTP 409 Conflict (Error: USERNAME_TAKEN or EMAIL_TAKEN)
        APIGW-->>ClientApp: HTTP 409 Conflict
        ClientApp-->>User: Display error message
    else Input is valid
        AuthSvc->>AuthSvc: Hash password (Argon2id)
        AuthSvc->>AuthSvc: Create User record in DB (status: pending_verification, user_id generated)
        AuthSvc->>AuthSvc: Generate verification code/token
        AuthSvc->>AuthSvc: Store verification code (e.g., in Redis with TTL, associated with user_id)
        AuthSvc-->>KafkaBus: Publish event `auth.user.registered.v1` (user_id, email, username, verification_code_or_token)
        AuthSvc-->>APIGW: HTTP 201 Created (user_id, status: "pending_verification")
        APIGW-->>ClientApp: HTTP 201 Created
        ClientApp-->>User: Display success message (e.g., "Регистрация почти завершена. Проверьте ваш email для подтверждения.")
    end

    subgraph Asynchronous Processing Post-Registration
        KafkaBus-->>AccountSvc: Consume event `auth.user.registered.v1` (user_id, email, username)
        AccountSvc->>AccountSvc: Create Account record (linked to user_id)
        AccountSvc->>AccountSvc: Create default Profile record (e.g., with username as nickname)
        AccountSvc->>AccountSvc: Create default UserSetting record
        AccountSvc->>AccountSvc: Add email to ContactInfo (status: pending_verification)
        AccountSvc-->>KafkaBus: Publish event `account.created.v1` (account_id, user_id)
        AccountSvc-->>KafkaBus: Publish event `account.contact.added.v1` (user_id, email, type: email, status: pending_verification)

        KafkaBus-->>NotificationSvc: Consume event `auth.user.registered.v1` (user_id, email, verification_code_or_token)
        NotificationSvc->>NotificationSvc: Get user contact preferences (default to email)
        NotificationSvc->>NotificationSvc: Prepare email verification message (using template `email_verification` and verification_code_or_token)
        NotificationSvc->>NotificationSvc: Send email via Email Provider (External)
        NotificationSvc-->>KafkaBus: Publish event `notification.message.status.updated.v1` (status: sent/failed)

        KafkaBus-->>SocialSvc: Consume event `account.created.v1` (user_id, username from AccountSvc based on auth.user.registered)
        SocialSvc->>SocialSvc: Create UserSocialProfile (with default privacy)
        SocialSvc->>SocialSvc: Create User node in Neo4j
        SocialSvc-->>KafkaBus: Publish event `social.user.profile.created.v1`
    end

    User->>ClientApp: Clicks verification link in email (e.g., https://client.app/verify-email?token=<verification_token>)
    ClientApp->>APIGW: POST /api/v1/auth/verify-email (payload: { "verification_token": "<token>" })
    APIGW->>AuthSvc: Forward POST /verify-email

    AuthSvc->>AuthSvc: Validate verification token (check existence, expiry, user_id match from token)
    alt Token is valid
        AuthSvc->>AuthSvc: Update User record status to 'active' in DB
        AuthSvc->>AuthSvc: Delete/Invalidate verification token
        AuthSvc-->>KafkaBus: Publish event `auth.user.email_verified.v1` (user_id, email)
        AuthSvc-->>APIGW: HTTP 200 OK (message: "Email verified successfully")
        APIGW-->>ClientApp: HTTP 200 OK
        ClientApp-->>User: Display success (e.g., "Email подтвержден. Теперь вы можете войти.")
    else Token is invalid/expired
        AuthSvc-->>APIGW: HTTP 400 Bad Request (Error: INVALID_VERIFICATION_TOKEN or TOKEN_EXPIRED)
        APIGW-->>ClientApp: HTTP 400 Bad Request
        ClientApp-->>User: Display error message
    end

    subgraph Asynchronous Processing Post-Email-Verification
        KafkaBus-->>AccountSvc: Consume event `auth.user.email_verified.v1` (user_id, email)
        AccountSvc->>AccountSvc: Update ContactInfo record for the email to `is_verified = true`, `is_primary = true` (if no other primary)
        AccountSvc-->>KafkaBus: Publish event `account.contact.verified.v1` (user_id, email)
        AccountSvc->>AccountSvc: Update Account status to `active` (if it was `pending_verification_email`)
        AccountSvc-->>KafkaBus: Publish event `account.status.updated.v1` (user_id, status: `active`)
    end
```

This diagram outlines the primary flow. Error handling within each service (e.g., database connection errors) is omitted for brevity but assumed to be handled according to service specifications.
The "verification_code_or_token" can be a simple code sent in the email body or a more complex JWT used in a link.
The Social Service integration is added to create a basic social profile upon account creation.
```
