**Дата последнего обновления:** 2024-07-16

# User Registration and Initial Profile Setup Workflow

This diagram illustrates user registration, including direct email/password signup and registration via external OAuth providers like VKontakte and Telegram.

```mermaid
sequenceDiagram
    actor User
    participant ClientApp as Client Application
    participant APIGW as API Gateway
    participant AuthSvc as Auth Service
    participant KafkaBus as Kafka Message Bus
    participant AccountSvc as Account Service
    participant NotificationSvc as Notification Service
    participant SocialLoginProvider as External OAuth Provider (VK/Telegram)

    alt Direct Email/Password Registration
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
            AuthSvc-->>KafkaBus: Publish event `auth.user.registered.v1` (user_id, email, username, verification_code_or_token, source: 'direct', email_verified: false)
            AuthSvc-->>APIGW: HTTP 201 Created (user_id, status: "pending_verification")
            APIGW-->>ClientApp: HTTP 201 Created
            ClientApp-->>User: Display success message (e.g., "Регистрация почти завершена. Проверьте ваш email для подтверждения.")
        end
    else Registration via VK/Telegram
        User->>ClientApp: Clicks "Register with VK/Telegram"
        ClientApp->>APIGW: GET /api/v1/auth/oauth/url?provider=vk (or telegram)
        APIGW->>AuthSvc: Forward GET /oauth/url?provider=vk
        AuthSvc->>AuthSvc: Generate OAuth state, build redirect URL to SocialLoginProvider
        AuthSvc-->>APIGW: HTTP 200 OK (redirect_url)
        APIGW-->>ClientApp: HTTP 200 OK (redirect_url)
        ClientApp->>User: Redirect to SocialLoginProvider
        User->>SocialLoginProvider: Authenticates and grants permissions
        SocialLoginProvider-->>AuthSvc: Callback to redirect_uri (e.g., /api/v1/auth/oauth/callback?provider=vk&code=...&state=...)
        Note over APIGW, AuthSvc: Callback is routed through APIGW to AuthSvc
        AuthSvc->>AuthSvc: Validate state, exchange code for token with SocialLoginProvider
        AuthSvc->>SocialLoginProvider: Request user info (email, name, etc.)
        SocialLoginProvider-->>AuthSvc: User info response
        AuthSvc->>AuthSvc: Check if user exists by provider_id or email. If not, create User record (status: 'active', email_verified: true if email provided and trusted by provider)
        AuthSvc->>AuthSvc: Generate internal JWTs (access & refresh tokens)
        AuthSvc-->>KafkaBus: Publish event `auth.user.registered.v1` (user_id, email, username_from_provider, source: 'vk'/'telegram', email_verified: true/false based on provider trust)
        AuthSvc-->>APIGW: HTTP 200 OK (tokens, user_info)
        APIGW-->>ClientApp: HTTP 200 OK (tokens, user_info)
        ClientApp-->>User: Registration/Login successful, redirect to dashboard
    end

    subgraph Asynchronous Processing Post-Registration
        KafkaBus-->>AccountSvc: Consume event `auth.user.registered.v1` (user_id, email, username, source)
        AccountSvc->>AccountSvc: Create Account record (linked to user_id)
        AccountSvc->>AccountSvc: Create default Profile record (e.g., with username as nickname, or from provider data if source is 'vk'/'telegram')
        AccountSvc->>AccountSvc: Create default UserSetting record
        AccountSvc->>AccountSvc: Add email to ContactInfo (status: if source=='direct' then 'pending_verification' else if provider email is trusted then 'verified' else 'pending_verification')
        AccountSvc-->>KafkaBus: Publish event `account.created.v1` (account_id, user_id)
        AccountSvc-->>KafkaBus: Publish event `account.contact.added.v1` (user_id, email, type: email, status: determined by source)

        KafkaBus-->>NotificationSvc: Consume event `auth.user.registered.v1` (user_id, email, verification_code_or_token, source, email_verified)
        alt source == 'direct' AND email_verified == false
            NotificationSvc->>NotificationSvc: Get user contact preferences (default to email)
            NotificationSvc->>NotificationSvc: Prepare email verification message (using template `email_verification` and verification_code_or_token)
            NotificationSvc->>NotificationSvc: Send email via Email Provider (External)
            NotificationSvc-->>KafkaBus: Publish event `notification.message.status.updated.v1` (status: sent/failed, type: 'email_verification')
        else if source != 'direct' AND email_verified == true
            NotificationSvc->>NotificationSvc: Prepare welcome message (template `welcome_oauth_user`)
            NotificationSvc->>NotificationSvc: Send email via Email Provider (External)
            NotificationSvc-->>KafkaBus: Publish event `notification.message.status.updated.v1` (status: sent/failed, type: 'welcome_oauth')
        end

        KafkaBus-->>SocialSvc: Consume event `account.created.v1` (user_id, username from AccountSvc based on auth.user.registered event)
        SocialSvc->>SocialSvc: Create UserSocialProfile (with default privacy)
        SocialSvc->>SocialSvc: Create User node in Neo4j
        SocialSvc-->>KafkaBus: Publish event `social.user.profile.created.v1`
    end

    opt Direct Email Registration Verification
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

        subgraph Asynchronous Processing Post-Email-Verification (for direct registration)
            KafkaBus-->>AccountSvc: Consume event `auth.user.email_verified.v1` (user_id, email)
            AccountSvc->>AccountSvc: Update ContactInfo record for the email to `is_verified = true`, `is_primary = true` (if no other primary)
            AccountSvc-->>KafkaBus: Publish event `account.contact.verified.v1` (user_id, email)
            AccountSvc->>AccountSvc: Update Account status to `active` (if it was `pending_verification_email`)
            AccountSvc-->>KafkaBus: Publish event `account.status.updated.v1` (user_id, status: `active`)
        end
    end
```

This diagram outlines the primary registration flows, including direct email/password and OAuth via external providers. Error handling within each service (e.g., database connection errors) is omitted for brevity but assumed to be handled according to service specifications.
The email verification flow, involving a "verification_code_or_token" and the user clicking a link, is primarily for direct email/password registration. For OAuth registrations, if an email is provided by the OAuth provider and is considered trusted, it may be marked as verified automatically by `AuthSvc`. If an email from an OAuth provider is not available, not trusted, or if the platform policy requires explicit verification for all emails, the email verification flow might still be initiated.
The Social Service integration creates a basic social profile upon account creation for both registration types.
The `auth.user.registered.v1` event now includes `source` (e.g., 'direct', 'vk', 'telegram') and `email_verified` fields to distinguish registration methods and reflect pre-verified email status.
```
