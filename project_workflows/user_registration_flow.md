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
    APIGW->>AuthSvc: Forward POST /register (payload)

    AuthSvc->>AuthSvc: Validate input (email format, password strength, username availability)
    alt Username/Email already exists
        AuthSvc-->>APIGW: HTTP 409 Conflict (Error: Username/Email taken)
        APIGW-->>ClientApp: HTTP 409 Conflict
        ClientApp-->>User: Display error message
    else Input is valid
        AuthSvc->>AuthSvc: Hash password (Argon2id)
        AuthSvc->>AuthSvc: Create User record (status: pending_verification)
        AuthSvc->>AuthSvc: Generate verification token/code
        AuthSvc-->>KafkaBus: Publish event `auth.user.registered.v1` (user_id, email, username, verification_code)
        AuthSvc-->>APIGW: HTTP 201 Created (User partially registered, needs verification)
        APIGW-->>ClientApp: HTTP 201 Created
        ClientApp-->>User: Display success message (e.g., "Check your email for verification")
    end

    subgraph Asynchronous Processing
        KafkaBus-->>AccountSvc: Consume event `auth.user.registered.v1`
        AccountSvc->>AccountSvc: Create Account record (linked to user_id from event)
        AccountSvc->>AccountSvc: Create default Profile record
        AccountSvc->>AccountSvc: Create default UserSetting record
        AccountSvc-->>KafkaBus: Publish event `account.created.v1` (account_id, user_id)
        AccountSvc-->>KafkaBus: Publish event `account.contact.added.v1` (for email, needs verification)

        KafkaBus-->>NotificationSvc: Consume event `auth.user.registered.v1` (or `account.contact.added.v1`)
        NotificationSvc->>NotificationSvc: Prepare email verification message (using template)
        NotificationSvc->>NotificationSvc: Send email via Email Provider (External)
        NotificationSvc-->>KafkaBus: Publish event `notification.sent.v1` (status: success/failure)
    end

    User->>ClientApp: Clicks verification link in email (contains token)
    ClientApp->>APIGW: GET /api/v1/auth/verify-email?token=<verification_token>
    APIGW->>AuthSvc: Forward GET /verify-email

    AuthSvc->>AuthSvc: Validate verification token
    alt Token is valid
        AuthSvc->>AuthSvc: Update User record status to 'active'
        AuthSvc-->>KafkaBus: Publish event `auth.user.email_verified.v1` (user_id)
        AuthSvc-->>APIGW: HTTP 200 OK (Email verified)
        APIGW-->>ClientApp: HTTP 200 OK
        ClientApp-->>User: Display success (e.g., "Email verified, you can now log in")
    else Token is invalid/expired
        AuthSvc-->>APIGW: HTTP 400 Bad Request (Error: Invalid token)
        APIGW-->>ClientApp: HTTP 400 Bad Request
        ClientApp-->>User: Display error message
    end

    KafkaBus-->>AccountSvc: Consume event `auth.user.email_verified.v1`
    AccountSvc->>AccountSvc: Update ContactInfo record for email to `is_verified = true`
    AccountSvc-->>KafkaBus: Publish event `account.contact.verified.v1`
```

This diagram outlines the primary flow. Error handling within each service (e.g., database connection errors) is omitted for brevity but assumed to be handled according to service specifications.
```
