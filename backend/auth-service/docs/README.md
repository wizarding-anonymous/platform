# Спецификация Микросервиса: Auth Service (Сервис Аутентификации и Авторизации)

**Версия:** 2.0
**Дата последнего обновления:** {{YYYY-MM-DD}}

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение и Роль:** Auth Service является центральным компонентом платформы "Российский Аналог Steam", отвечающим за управление идентификацией пользователей, аутентификацию и авторизацию. Он обеспечивает безопасный, масштабируемый и надежный фундамент для модели безопасности всей платформы. Его основная цель - проверка подлинности пользователей, контроль доступа к ресурсам на основе ролей и разрешений, управление сессиями и токенами (JWT), а также предоставление доверенной информации об аутентифицированных пользователях другим микросервисам.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   Регистрация пользователей (в координации с Account Service).
*   Аутентификация по логину/паролю.
*   Двухфакторная аутентификация (2FA): TOTP, SMS/Email коды.
*   Внешняя аутентификация: OAuth 2.0 / OIDC (Telegram, VK, Odnoklassniki).
*   Управление JSON Web Token (JWT): Генерация, валидация, ротация, отзыв.
*   Управление сессиями.
*   Управление паролями: Хранение (хеширование Argon2id), сброс, изменение.
*   Подтверждение Email.
*   Role-Based Access Control (RBAC).
*   Управление API ключами.
*   Аудит событий безопасности.
*   Обнаружение подозрительной активности.
*   Предоставление информации о пользователе и правах другим сервисам.
*   Административные функции (через Admin Service).

### 1.3. Основные Технологии
*   **Язык программирования:** Go (1.21+).
*   **API Фреймворки:** REST (Gin/Echo), gRPC.
*   **Базы данных:** PostgreSQL (GORM/pgx), Redis.
*   **Брокер сообщений/События:** Apache Kafka.
*   **Безопасность (криптография):** JWT (RS256), Argon2id, TOTP.
*   **Логирование:** Zap.
*   **Мониторинг:** Prometheus.
*   **Трассировка:** OpenTelemetry.
*   **Управление конфигурацией:** Viper.
*   **Контейнеризация/Оркестрация:** Docker, Kubernetes.
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `../../../../PACKAGE_STANDARDIZATION.md`.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`, `../../../../project_observability_standards.md`.

### 1.4. Термины и Определения (Glossary)
*   См. `../../../../project_glossary.md` для терминов: Access Token, Refresh Token, JWT, RBAC, 2FA, TOTP, Argon2id, OAuth, OIDC, JWKS, JTI, User, Role, Permission, Session, API Key.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Auth Service разработан как stateless микросервис на Go. Состояние сессий и временная информация хранится в Redis. PostgreSQL используется для постоянного хранения. Kafka для асинхронного обмена событиями. Архитектура многослойная, ориентирована на безопасность и масштабируемость.
*   Диаграмма верхнеуровневой архитектуры:
```mermaid
graph TD
    subgraph AuthService [Auth Service]
        direction LR
        subgraph PresentationLayer [API Layer (Presentation)]
            direction TB
            API_REST[REST API (Gin/Echo)]
            API_GRPC[gRPC API]
        end

        subgraph ApplicationLayer [Business Logic Layer (Application)]
            direction TB
            RegSvc[Registration Service]
            LoginSvc[Login Service]
            TokenSvc[Token Management (JWT, Refresh)]
            TwoFASvc[2FA Service (TOTP, SMS/Email)]
            RBACSvc[RBAC Service (Roles, Permissions)]
            ApiKeySvc[API Key Management]
            ExternalAuthSvc[External Auth (OAuth)]
            PasswordSvc[Password Management (Reset, Change)]
            EmailVerifySvc[Email Verification Service]
            SessionSvc[Session Management]
        end

        subgraph DomainLayer [Domain Layer]
            direction TB
            Entities[Entities (User, Role, Session, Token)]
            ValueObjects[Value Objects (HashedPassword, Email)]
            DomainEvents[Domain Events (UserRegistered, SessionCreated)]
            RepositoriesIntf[Repository Interfaces]
        end

        subgraph InfrastructureLayer [Data Access & Infrastructure Layer]
            direction TB
            RepoPostgreSQL[PostgreSQL Repositories Impl]
            CacheRedis[Redis Client (Session Cache, 2FA codes, Rate Limits, JTI Blacklist)]
            ProducerKafka[Kafka Producer (Domain Events)]
            CryptoUtils[Cryptography (Argon2id, JWT Signing/Parsing)]
            ExtSystemClients[External System Clients (e.g., Notification Service client)]
        end

        PresentationLayer --> ApplicationLayer
        ApplicationLayer --> DomainLayer
        ApplicationLayer --> InfrastructureLayer
        DomainLayer --> RepositoriesIntf
        InfrastructureLayer -- Implements --> RepositoriesIntf
    end

    Clients[Clients (Web, Mobile, Desktop Apps)] --> APIGateway[API Gateway]
    APIGateway -- REST/gRPC Requests --> PresentationLayer

    InfrastructureLayer -- CRUD Operations --> DB[(Database: PostgreSQL)]
    InfrastructureLayer -- Cache Operations --> Cache[(Cache: Redis)]

    ApplicationLayer -- Publishes Domain Events --> KafkaBroker[Kafka Broker]
    KafkaBroker -- com.platform.auth.user.registered.v1 etc. --> AccountSvc[Account Service]
    KafkaBroker -- com.platform.auth.user.verification_code_sent.v1 etc. --> NotificationSvc[Notification Service]

    InternalMS[Other Microservices] -- gRPC: ValidateToken, CheckPermission --> API_GRPC
```

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления)
*   Ответственность: Обработка HTTP (REST) и gRPC запросов, валидация DTO, вызов Application Layer.
*   Ключевые компоненты/модули: HTTP Handlers (Gin/Echo), gRPC Service Implementations, DTOs, парсеры/сериализаторы.

#### 2.2.2. Application Layer (Прикладной Слой / Слой Сценариев Использования)
*   Ответственность: Координация бизнес-логики, реализация use cases.
*   Ключевые компоненты/модули: `UserRegistrationService`, `AuthenticationService`, `TokenService`, `TwoFactorAuthService`, `RBACService`, `SessionManagementService`. Интерфейсы для репозиториев и внешних сервисов.

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности, value objects, доменные события, бизнес-правила.
*   Ключевые компоненты/модули: Entities (`User`, `Role`, `Session`, `RefreshToken`, `APIKey`, `VerificationCode`, `MFASecret`, `ExternalAccount`), Value Objects (`Email`, `HashedPassword`), Domain Events (`UserRegisteredEvent`), интерфейсы репозиториев.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Реализация интерфейсов для взаимодействия с PostgreSQL, Redis, Kafka, крипто-утилитами, внешними сервисами.
*   Ключевые компоненты/модули: PostgreSQL Repositories, Redis Cache Client, Kafka Producers, Cryptography Utilities, клиенты для внешних сервисов (Notification Service).

## 3. API Endpoints
Общие принципы и форматы см. `../../../../project_api_standards.md`.

### 3.1. REST API
*   **Базовый URL:** `/api/v1/auth` (маршрутизируется через API Gateway).
*   **Формат данных:** JSON.
*   **Аутентификация:** Большинство эндпоинтов требуют `Authorization: Bearer <access_token>`.

#### 3.1.1. Регистрация и Вход
*   **`POST /register`**
    *   Описание: Регистрация нового пользователя.
    *   Тело запроса: (см. существующий документ).
    *   Пример ответа (Успех 201 Created): (см. существующий документ).
    *   Пример ответа (Ошибка 400 Validation Error): (см. существующий документ).
    *   Требуемые права доступа: Публичный.
*   **`POST /login`**
    *   Описание: Аутентификация пользователя и получение токенов.
    *   Тело запроса: (см. существующий документ).
    *   Пример ответа (Успех 200 OK): (см. существующий документ, refresh token в HttpOnly cookie).
    *   Пример ответа (Ошибка 401 Invalid Credentials): (см. существующий документ).
    *   Требуемые права доступа: Публичный.
*   **`POST /login/2fa`**
    *   Описание: Завершение логина с использованием кода 2FA (TOTP, SMS/Email код).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "2faLoginCompletion",
            "attributes": {
              "mfa_session_token": "temp_token_after_password_ok", // Токен, полученный после успешного ввода пароля, если требуется 2FA
              "mfa_code": "123456" // Код TOTP или SMS/Email
            }
          }
        }
        ```
    *   Пример ответа: Аналогичен `/login`.
    *   Требуемые права доступа: Публичный (но требует `mfa_session_token`).
*   **`POST /oauth/{provider}/callback`**
    *   Описание: Обработка callback от OAuth 2.0 / OIDC провайдера (например, `/oauth/vk/callback`).
    *   Query параметры: `code`, `state` (от провайдера).
    *   Пример ответа: Аналогичен `/login` или редирект с сессией.
    *   Требуемые права доступа: Публичный.
    *   `[NEEDS DEVELOPER INPUT: Detailed flow for each OAuth provider (VK, Telegram, etc.) including initial redirect and specific parameters, or link to a general OAuth workflow document if it exists.]`

#### 3.1.2. Управление токенами
*   **`POST /refresh-token`**
    *   Описание: Обновление access token.
    *   Тело запроса (если refresh token в теле): (см. существующий документ).
    *   Пример ответа (Успех 200 OK): (см. существующий документ).
    *   Требуемые права доступа: Валидный refresh token.
*   **`POST /logout`**
    *   Описание: Отзыв текущего refresh token.
    *   Пример ответа (Успех 204 No Content).
    *   Требуемые права доступа: Валидный refresh token.
*   **`GET /.well-known/jwks.json`**
    *   Описание: Публичный эндпоинт для получения JWKS (JSON Web Key Set) для валидации подписи JWT.
    *   Пример ответа: `[NEEDS DEVELOPER INPUT: Example JWKS response]`
    *   Требуемые права доступа: Публичный.

#### 3.1.3. Верификация Email
*   **`POST /verify-email`**
    *   Описание: Подтверждение email с кодом верификации.
    *   Тело запроса: (см. существующий документ).
    *   Пример ответа (Успех 200 OK): (см. существующий документ).
    *   Требуемые права доступа: Публичный или аутентифицированный пользователь.
*   **`POST /resend-verification-email`**
    *   Описание: Повторная отправка письма для верификации email.
    *   Тело запроса: `{"data": {"type": "resendVerification", "attributes": {"email": "user@example.com"}}}` (если не аутентифицирован) или пустое тело (если аутентифицирован).
    *   Требуемые права доступа: Публичный или аутентифицированный пользователь.

#### 3.1.4. Управление Паролем
*   **`POST /forgot-password`**
    *   Описание: Запрос на сброс пароля.
    *   Тело запроса: (см. существующий документ).
    *   Пример ответа (Успех 200 OK): (см. существующий документ).
    *   Требуемые права доступа: Публичный.
*   **`POST /reset-password`**
    *   Описание: Установка нового пароля с использованием токена сброса.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "passwordReset",
            "attributes": {
              "reset_token": "long_secure_reset_token_string",
              "new_password": "NewSecurePassword456!",
              "confirm_new_password": "NewSecurePassword456!"
            }
          }
        }
        ```
    *   Требуемые права доступа: Публичный.
*   **`PUT /me/password`**
    *   Описание: Изменение пароля аутентифицированным пользователем.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "passwordChange",
            "attributes": {
              "current_password": "OldSecurePassword123!",
              "new_password": "NewSecurePassword789!",
              "confirm_new_password": "NewSecurePassword789!"
            }
          }
        }
        ```
    *   Требуемые права доступа: Аутентифицированный пользователь.

#### 3.1.5. Управление Двухфакторной Аутентификацией (2FA)
*   **`POST /me/2fa/totp/enable`**: Инициирует включение TOTP 2FA. Ответ: (см. существующий документ).
*   **`POST /me/2fa/totp/verify`**: Подтверждает включение TOTP 2FA. Ответ: (см. существующий документ, включает резервные коды).
*   **`POST /me/2fa/totp/disable`**
    *   Описание: Отключение TOTP 2FA.
    *   Тело запроса: `{"data": {"type": "totpDisable", "attributes": {"totp_code": "123456"}}}` (или пароль в зависимости от политики).
    *   Требуемые права доступа: Аутентифицированный пользователь с активным TOTP.
*   `[NEEDS DEVELOPER INPUT: Endpoints for SMS/Email based 2FA setup and verification if implemented, including request/response examples]`
*   **`GET /me/2fa/backup-codes`**: Получение новых резервных кодов (старые инвалидируются).
*   **`GET /me/2fa/status`**: Получение текущего статуса 2FA для пользователя.

#### 3.1.6. Управление API Ключами
*   **`GET /me/api-keys`**: Получение списка API ключей пользователя (без самих ключей, только метаданные).
*   **`POST /me/api-keys`**: Создание нового API ключа. Ответ: (см. существующий документ, ключ возвращается один раз).
*   **`DELETE /me/api-keys/{key_id_or_prefix}`**: Отзыв API ключа.
*   Требуемые права доступа: Аутентифицированный пользователь (например, `developer`).

#### 3.1.7. Управление Сессиями
*   **`GET /me/sessions`**
    *   Описание: Получение списка активных сессий пользователя.
    *   Пример ответа: `[NEEDS DEVELOPER INPUT: Example response for GET /me/sessions]`
    *   Требуемые права доступа: Аутентифицированный пользователь.
*   **`DELETE /me/sessions/{session_id}`**
    *   Описание: Отзыв конкретной сессии пользователя.
    *   Требуемые права доступа: Аутентифицированный пользователь.
*   **`DELETE /me/sessions/all-others`**
    *   Описание: Отзыв всех сессий, кроме текущей.
    *   Требуемые права доступа: Аутентифицированный пользователь.

### 3.2. gRPC API
*   **Пакет:** `auth.v1`
*   **Файл .proto:** `proto/auth/v1/auth_service.proto` (или в `platform-protos`).
*   **Аутентификация:** mTLS для межсервисного взаимодействия.

#### 3.2.1. Сервис: AuthService
*   **`rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse)`**
    *   Описание: Валидация access token и возврат информации о пользователе и его правах.
    *   `message ValidateTokenRequest { string access_token = 1; }`
    *   `message ValidateTokenResponse { bool is_valid = 1; string user_id = 2; string username = 3; repeated string roles = 4; repeated string permissions = 5; google.protobuf.Struct claims = 6; }`
*   **`rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse)`**
    *   Описание: Проверка, имеет ли пользователь указанное разрешение.
    *   `message CheckPermissionRequest { string user_id = 1; string permission_code = 2; map<string, string> context = 3; }`
    *   `message CheckPermissionResponse { bool is_granted = 1; }`
*   `[NEEDS DEVELOPER INPUT: Add other relevant gRPC methods if any, e.g., for internal service-to-service API key validation or user info retrieval by ID for trusted services]`

### 3.3. WebSocket API
*   Не применимо для данного сервиса.

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`User`**: Пользователь (логин, хеш пароля, email, статус, флаг 2FA).
*   **`Role`**: Роль (имя, описание, список разрешений).
*   **`Permission`**: Разрешение (код, описание) - `[NEEDS DEVELOPER INPUT: Clarify if permissions are predefined strings or stored in DB. If in DB, add Permission entity and User/Role to Permission mapping table if many-to-many.]`
*   **`Session`**: Сессия пользователя (ID, UserID, IP, User-Agent, время активности).
*   **`RefreshToken`**: Токен обновления (ID, UserID, хеш токена, SessionID, статус отзыва).
*   **`APIKey`**: API Ключ (ID, UserID, имя, префикс, хеш ключа, разрешения).
*   **`MFASecret`**: Секрет MFA (UserID, тип, зашифрованный секрет, статус верификации).
*   **`MFABackupCode`**: Резервный код MFA (ID, UserID, хеш кода, статус использования).
*   **`VerificationCode`**: Код верификации (ID, UserID, тип, хеш кода, цель, срок действия).
*   **`AuditLogAuth`**: Лог аудита (ID, UserID, ActorID, тип действия, IP, User-Agent, детали, время).
*   **`ExternalAccount`**: Внешний аккаунт OAuth/OIDC (ID, UserID, провайдер, ID пользователя у провайдера, детали, токены).

### 4.2. Схема Базы Данных (PostgreSQL)
*   ERD Диаграмма: (см. существующий документ, дополнена новыми сущностями).
*   DDL: (см. существующий документ, дополнен DDL для новых таблиц: `mfa_secrets`, `mfa_backup_codes`, `verification_codes`, `audit_logs_auth`, `external_accounts`).
*   `[NEEDS DEVELOPER INPUT: Review and confirm DDL for all tables, especially constraints, indexes, and relationships for auth-service. Clarify Permission entity storage.]`

#### Redis Data Structures
*   **JTI Blacklist:** `blacklist:jti:<JTI_value>` -> `expiry_timestamp` (SETEX)
*   **Session Data:** `session:<session_id>` -> `user_id, ip_address, user_agent, last_activity` (HASH/JSON) (SETEX with session TTL)
*   **2FA Codes (SMS/Email):** `2fa_code:<user_id>:<type>` -> `hashed_code` (SETEX with short TTL)
*   **Rate Limiting:** `rate_limit:<action>:<user_id_or_ip>` -> `count` (INCR with EXPIRE)
*   **Password Reset Tokens:** `pwd_reset_token:<token_hash>` -> `user_id` (SETEX with TTL)
*   **Email Verification Tokens:** `email_verify_token:<token_hash>` -> `user_id` (SETEX with TTL)

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   Система: Apache Kafka. Формат: CloudEvents JSON. Топик: `com.platform.auth.events.v1`.
*   **`com.platform.auth.user.registered.v1`**: (см. существующий документ).
*   **`com.platform.auth.user.login_success.v1`**: (см. существующий документ).
*   **`com.platform.auth.user.email_verified.v1`**: (см. существующий документ).
*   **`com.platform.auth.user.password_changed.v1`**: (см. существующий документ).
*   **`com.platform.auth.user.mfa_status.updated.v1`**: (см. существующий документ).
*   **`com.platform.auth.session.revoked.v1`**: (см. существующий документ).
*   `[NEEDS DEVELOPER INPUT: List or confirm other key published events like password_reset_requested, login_failed, account_locked, roles_changed, api_key_created/revoked, etc., with their payloads.]`

### 5.2. Потребляемые События (Consumed Events)
*   **`com.platform.account.user.profile_updated.v1`**: (см. существующий документ).
*   **`com.platform.admin.user.force_logout.v1`**: (см. существующий документ).
*   **`com.platform.admin.user.status.updated.v1`**: (см. существующий документ).
*   `[NEEDS DEVELOPER INPUT: List or confirm other key consumed events and their processing logic.]`

## 6. Интеграции (Integrations)
См. `../../../../project_integrations.md`.
[NEEDS DEVELOPER INPUT: Review and confirm service boundaries and reliability strategies for each integration for auth-service, especially Notification Service and third-party OAuth providers.]

### 6.1. Внутренние Микросервисы
*   **API Gateway**: Валидация токенов.
*   **Account Service**: Координация при регистрации.
*   **Notification Service**: Отправка email/SMS. Надежность: `[NEEDS DEVELOPER INPUT: Strategy for Notification Service unavailability, e.g., retry, DLQ for notifications]`
*   **Admin Service**: Управление пользователями/ролями.
*   **Другие микросервисы**: Проверка токенов/прав через gRPC.

### 6.2. Внешние Системы
*   **OAuth 2.0/OIDC Провайдеры (VK, Telegram, etc.)**: Для внешней аутентификации.
    *   Контракт: Стандартные потоки OAuth 2.0/OIDC.
    *   Надежность: `[NEEDS DEVELOPER INPUT: Handling of external provider unavailability or errors.]`
*   **Системы доставки SMS/Email (через Notification Service)**.

## 7. Конфигурация (Configuration)
Общие стандарты: `../../../../project_api_standards.md` и `../../../../DOCUMENTATION_GUIDELINES.md`.

### 7.1. Переменные Окружения
*   `AUTH_HTTP_PORT`, `AUTH_GRPC_PORT`
*   `POSTGRES_DSN`, `REDIS_ADDR`, `KAFKA_BROKERS`
*   `AUTH_JWT_PRIVATE_KEY_PATH`, `AUTH_JWT_PUBLIC_KEY_PATH` (монтируются из Secrets)
*   `AUTH_JWT_ACCESS_TOKEN_TTL_SECONDS`, `AUTH_JWT_REFRESH_TOKEN_TTL_SECONDS`
*   `AUTH_SESSION_TTL_SECONDS`
*   `AUTH_ARGON2ID_MEMORY_KB`, `AUTH_ARGON2ID_ITERATIONS`, `AUTH_ARGON2ID_PARALLELISM`
*   OAuth provider client IDs/secrets (из Secrets, например `AUTH_PROVIDER_VK_CLIENT_ID_SECRET_PATH`)
*   `AUTH_MFA_ENCRYPTION_KEY_SECRET_PATH` (ключ для шифрования MFA секретов, из Secrets)
*   `LOG_LEVEL`, `OTEL_EXPORTER_JAEGER_ENDPOINT`
*   `[NEEDS DEVELOPER INPUT: Review and add any other critical environment variables, especially for security parameters and rate limits.]`

### 7.2. Файлы Конфигурации (`config/config.yaml`)
*   Расположение: `backend/auth-service/config/config.yaml`.
*   Структура: (см. существующий документ, с акцентом на то, что секреты загружаются из файлов, указанных переменными окружения).
    ```yaml
    # ...
    jwt:
      private_key_file: ${AUTH_JWT_PRIVATE_KEY_PATH} # e.g., /run/secrets/jwt_private.pem
      public_key_file: ${AUTH_JWT_PUBLIC_KEY_PATH}  # e.g., /app/config/jwt_public.pem
      access_token_ttl: ${AUTH_JWT_ACCESS_TOKEN_TTL_SECONDS:900}
      refresh_token_ttl: ${AUTH_JWT_REFRESH_TOKEN_TTL_SECONDS:2592000} # 30 days
    # ...
    oauth_providers:
      vk:
        client_id_secret_path: ${AUTH_PROVIDER_VK_CLIENT_ID_SECRET_PATH} # File containing client_id
        # ...
    mfa:
      encryption_key_secret_path: ${AUTH_MFA_ENCRYPTION_KEY_SECRET_PATH}
    # ...
    ```

## 8. Обработка Ошибок (Error Handling)
*   Стандартные HTTP коды и форматы ошибок REST согласно `../../../../project_api_standards.md`.
*   gRPC ошибки согласно стандартам gRPC.
*   Логирование с `trace_id`.

### 8.1. Общие Принципы
*   Стандартный JSON формат ответа об ошибке для REST.

### 8.2. Распространенные Коды Ошибок
*   `VALIDATION_ERROR`, `INVALID_CREDENTIALS`, `TOKEN_EXPIRED`, `TOKEN_INVALID`, `SESSION_EXPIRED`, `USER_NOT_FOUND`, `EMAIL_ALREADY_EXISTS`, `USERNAME_ALREADY_EXISTS`, `MFA_REQUIRED`, `MFA_CODE_INVALID`, `ACCOUNT_LOCKED`, `PERMISSION_DENIED`.
*   `[NEEDS DEVELOPER INPUT: Review and add other specific error codes for auth-service.]`

## 9. Безопасность (Security)
См. `../../../../project_security_standards.md`.

### 9.1. Аутентификация
*   JWT (RS256), OAuth 2.0/OIDC, 2FA (TOTP, SMS/Email).

### 9.2. Авторизация
*   RBAC. Проверка прав доступа на уровне gRPC и API Gateway.

### 9.3. Защита Данных
*   Хеширование паролей (Argon2id).
*   Шифрование MFA секретов (AES-GCM, ключ из Secrets).
*   TLS 1.3.
*   Защита от CSRF, XSS (для любых веб-интерфейсов, управляемых Auth Service, если есть).
*   Rate limiting, защита от брутфорса.
*   ФЗ-152: Обработка ПДн (email, username, IP, User-Agent, телефон для SMS 2FA, данные OAuth). Хранение и обработка в РФ. Согласие на обработку ПДн. Логирование операций с ПДн.

### 9.4. Управление Секретами
*   Приватные ключи JWT, ключи шифрования MFA, OAuth секреты - через Kubernetes Secrets / HashiCorp Vault.

## 10. Развертывание (Deployment)
См. `../../../../project_deployment_standards.md`.

### 10.1. Инфраструктурные Файлы
*   Dockerfile: `backend/auth-service/Dockerfile`.
*   Helm-чарты: `deploy/charts/auth-service/`.

### 10.2. Зависимости при Развертывании
*   PostgreSQL, Redis, Kafka.
*   Notification Service (для 2FA, верификации).
*   API Gateway.

### 10.3. CI/CD
*   Стандартный пайплайн. Автоматические миграции БД. Тестирование безопасности (SAST, DAST).

## 11. Мониторинг и Логирование (Logging and Monitoring)
См. `../../../../project_observability_standards.md`.

### 11.1. Логирование
*   Формат: JSON (Zap).
*   Ключевые события: Попытки входа (успешные/неудачные), регистрация, смена пароля, операции 2FA, генерация/отзыв токенов, ошибки.
*   Интеграция: Loki/ELK.

### 11.2. Мониторинг
*   Метрики (Prometheus): Количество регистраций, логинов, ошибок аутентификации, активных сессий, сгенерированных токенов, задержки API.
*   Дашборды (Grafana): Обзор состояния, производительности, безопасности.
*   Алерты (AlertManager): Аномальное количество ошибок входа, проблемы с доступностью зависимостей, истечение срока действия сертификатов/ключей.

### 11.3. Трассировка
*   Интеграция: OpenTelemetry. Экспорт: Jaeger.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**: P99 логин/регистрация < 300мс. P99 валидация токена < 20мс.
*   **Масштабируемость**: Поддержка >10 млн пользователей, >5000 запросов/сек на аутентификацию.
*   **Надежность**: Доступность 99.99%. RTO/RPO см. раздел 14.
*   **Безопасность**: Соответствие OWASP ASVS Level 2 (целевой).
*   [NEEDS DEVELOPER INPUT: Confirm or update specific NFR values for auth-service.]

## 13. Резервное Копирование и Восстановление (Backup and Recovery)
Общие принципы см. в `../../../../project_database_structure.md`.

### 13.1. PostgreSQL
*   Процедура: Ежедневный `pg_dumpall`/`pg_dump`. Непрерывная архивация WAL (PITR).
*   Хранение: S3. RTO: < 1ч. RPO: < 5мин.

### 13.2. Redis
*   Стратегия: AOF с fsync `everysec` для JTI blacklist. RDB snapshots (каждые 1-6ч).
*   Хранение: RDB в S3. RTO: < 30мин. RPO: < 1мин (для JTI).

### 13.3. Общая стратегия
*   Критически важно. Автоматизировано, регулярно тестируется. Мониторинг бэкапов.
*   Ключи шифрования (MFA, JWT) должны быть частью стратегии бэкапа системы управления секретами.

## 14. Приложения (Appendices) (Опционально)
*   OpenAPI спецификация: `[NEEDS DEVELOPER INPUT: Link to OpenAPI spec or state if generated from code for auth-service]`
*   Protobuf схемы: `platform-protos/auth/v1/auth_service.proto` (предположительно).
*   `[NEEDS DEVELOPER INPUT: Add any other appendices if necessary for auth-service, e.g., detailed OAuth 2.0 flow diagrams if not covered by workflow docs.]`

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [Регистрация пользователя и начальная настройка профиля](../../../../project_workflows/user_registration_flow.md)
*   [Аутентификация пользователя (логин, 2FA, OAuth)](../../../../project_workflows/user_authentication_flow.md) `[NEEDS DEVELOPER INPUT: Create and link this workflow document]`
*   [Сброс и восстановление пароля](../../../../project_workflows/password_recovery_flow.md) `[NEEDS DEVELOPER INPUT: Create and link this workflow document]`
*   [Процесс управления сессиями пользователей] `[NEEDS DEVELOPER INPUT: Create and link this workflow document if specific flows are defined, e.g., project_workflows/user_session_management_flow.md]`
*   [Процесс работы с API ключами] `[NEEDS DEVELOPER INPUT: Create and link this workflow document if specific flows are defined, e.g., project_workflows/api_key_management_flow.md]`

---
*Этот документ является основной спецификацией для Auth Service и должен поддерживаться в актуальном состоянии.*
