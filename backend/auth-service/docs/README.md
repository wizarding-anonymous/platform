# Спецификация Микросервиса: Auth Service (Сервис Аутентификации и Авторизации)

**Версия:** 2.0
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение и Роль:** Auth Service является центральным компонентом платформы "Российский Аналог Steam", отвечающим за управление идентификацией пользователей, аутентификацию и авторизацию. Он обеспечивает безопасный, масштабируемый и надежный фундамент для модели безопасности всей платформы. Его основная цель - проверка подлинности пользователей, контроль доступа к ресурсам на основе ролей и разрешений, управление сессиями и токенами (JWT), а также предоставление доверенной информации об аутентифицированных пользователях другим микросервисам.
*   Разработка сервиса должна вестись в соответствии с `CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   Регистрация пользователей (в координации с Account Service).
*   Аутентификация по логину/паролю.
*   Двухфакторная аутентификация (2FA): TOTP, SMS/Email коды (через Notification Service), резервные коды.
*   Внешняя аутентификация: OAuth 2.0 / OIDC (например, Telegram, VK, Odnoklassniki).
*   Управление JSON Web Token (JWT): Генерация (Access Token RS256, Refresh Token), валидация, ротация, отзыв (через JTI blacklist).
*   Управление сессиями: Отслеживание и отзыв активных сессий пользователей.
*   Управление паролями: Безопасное хранение (хеширование Argon2id), сброс и изменение пароля.
*   Подтверждение Email.
*   Role-Based Access Control (RBAC): Управление ролями и разрешениями пользователей.
*   Управление API ключами для внешних и внутренних сервисов.
*   Аудит событий безопасности, связанных с аутентификацией и авторизацией.
*   Обнаружение подозрительной активности (например, частые неудачные попытки входа).
*   Предоставление информации о пользователе и его правах другим сервисам.
*   Административные функции для управления пользователями и ролями (через Admin Service).

### 1.3. Основные Технологии
*   **Язык программирования:** Go.
*   **API Фреймворки:** Gin (REST), standard `net/http` и `grpc-go` (gRPC).
*   **Базы данных:** PostgreSQL (пользовательские данные, роли, разрешения, токены), Redis (кэш сессий, временные коды 2FA, JTI blacklist, счетчики rate limiting).
*   **Сообщения/События:** Apache Kafka (`confluent-kafka-go`) для асинхронного обмена событиями.
*   **Безопасность (криптография):** `golang-jwt/jwt/v5` (JWT RS256), `golang.org/x/crypto/argon2` (хеширование паролей Argon2id), стандартные библиотеки для TOTP.
*   **Логирование:** Zap.
*   **Мониторинг:** Prometheus (`prometheus-client_golang`).
*   **Трассировка:** OpenTelemetry.
*   **Контейнеризация:** Docker.
*   **Оркестрация:** Kubernetes.
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `PACKAGE_STANDARDIZATION.md`.

### 1.4. Термины и Определения (Glossary)
*   Для используемых терминов (Access Token, Refresh Token, JWT, RBAC, 2FA, TOTP, Argon2id, OAuth, OIDC, JWKS, JTI, User, Role, Permission, Session, API Key и др.) см. `project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Auth Service разработан как stateless (не хранящий состояние между запросами) микросервис на языке Go. Состояние сессий и другая временная информация хранится в Redis.
*   Он взаимодействует с PostgreSQL для постоянного хранения данных (учетные записи, роли, разрешения, токены обновления и т.д.), Redis для кэширования и временных данных, и Kafka для асинхронного обмена доменными событиями.
*   Архитектура ориентирована на безопасность, масштабируемость, отказоустойчивость и сопровождаемость, следуя принципам многослойной архитектуры.

**Диаграмма верхнеуровневой архитектуры:**
```mermaid
graph TD
    subgraph Auth Service
        direction LR
        subgraph PresentationLayer [API Layer (Presentation)]
            direction TB
            API_REST[REST API (Gin)]
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
    KafkaBroker -- auth.user.registered etc. --> AccountSvc[Account Service]
    KafkaBroker -- auth.user.verification_code_sent etc. --> NotificationSvc[Notification Service]

    InternalMS[Other Microservices] -- gRPC: ValidateToken, CheckPermission --> API_GRPC
```

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления)
*   **Ответственность:** Обработка входящих HTTP (RESTful через Gin) и gRPC запросов, валидация входных данных (DTO), преобразование данных и вызов соответствующей бизнес-логики в Application Layer.
*   **Ключевые компоненты/модули:**
    *   **HTTP Handlers (Gin):** Контроллеры для каждого REST эндпоинта (например, `POST /register`, `POST /login`).
    *   **gRPC Service Implementations:** Реализации gRPC сервисов (например, `AuthServiceV1` с методами `ValidateToken`, `CheckPermission`).
    *   **DTOs (Data Transfer Objects):** Структуры для запросов и ответов API, включая валидацию данных.
    *   **Парсеры и сериализаторы:** Для JSON и Protobuf.

#### 2.2.2. Application Layer (Прикладной Слой / Слой Сценариев Использования)
*   **Ответственность:** Координация бизнес-логики, реализация сценариев использования (use cases). Не содержит бизнес-правил напрямую, а оркестрирует взаимодействие между Domain Layer и Infrastructure Layer.
*   **Ключевые компоненты/модули:**
    *   **Use Case Services / Application Services:** Сервисы для каждого основного бизнес-процесса (например, `UserRegistrationService`, `AuthenticationService`, `TokenService`, `TwoFactorAuthService`, `RBACService`, `SessionManagementService`).
    *   **Интерфейсы для репозиториев и внешних сервисов:** Определяются здесь, а их реализации находятся в Infrastructure Layer.

#### 2.2.3. Domain Layer (Доменный Слой)
*   **Ответственность:** Содержит бизнес-сущности (entities), объекты-значения (value objects), доменные события и бизнес-правила. Независим от других слоев.
*   **Ключевые компоненты/модули:**
    *   **Entities (Сущности):** `User` (с методами для управления паролем, статусом), `Role`, `Permission`, `Session`, `RefreshToken`, `APIKey`, `VerificationCode`, `MFASecret`, `ExternalAccount`.
    *   **Value Objects (Объекты-значения):** Например, `Email`, `HashedPassword`, `PhoneNumber`, `UserID`.
    *   **Domain Services:** Для сложной доменной логики, не принадлежащей одной сущности (например, сервис для проверки уникальности username/email, если это не просто ограничение БД).
    *   **Domain Events (События домена):** Например, `UserRegisteredEvent`, `PasswordChangedEvent`, `SessionCreatedEvent`.
    *   **Интерфейсы репозиториев:** Определяют контракты для сохранения и извлечения сущностей (например, `UserRepository`, `SessionRepository`).

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   **Ответственность:** Реализация интерфейсов, определенных в Application и Domain Layers, для взаимодействия с внешними системами (базы данных, брокеры сообщений, другие микросервисы, криптографические утилиты).
*   **Ключевые компоненты/модули:**
    *   **Database Repositories (PostgreSQL):** Реализации интерфейсов репозиториев для PostgreSQL.
    *   **Cache Implementations (Redis):** Реализация для кэширования сессий, хранения JTI, кодов 2FA, счетчиков rate limiting.
    *   **Message Queue Producers (Kafka):** Продюсеры для отправки доменных событий в Kafka.
    *   **Cryptography Utilities:** Модули для хеширования паролей (Argon2id), генерации и валидации JWT (RS256), генерации TOTP.
    *   **Клиенты для внешних сервисов:** Например, HTTP-клиент для взаимодействия с Notification Service (если отправка кодов 2FA идет через него синхронно, а не через события).

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL:** `/api/v1/auth` (маршрутизируется через API Gateway).
*   **Формат данных:** JSON. Стандартный формат ответа об ошибке соответствует `project_api_standards.md`:
    ```json
    {
      "errors": [
        {
          "status": "4XX/5XX",
          "code": "ERROR_CODE_UPPER_SNAKE_CASE",
          "title": "Краткое описание ошибки на русском",
          "detail": "Полное описание ошибки с контекстом, если применимо.",
          "source": { "pointer": "/data/attributes/field_name", "parameter": "query_param_name" } // Опционально
        }
      ]
    }
    ```
*   **Аутентификация:** Большинство эндпоинтов требуют `Authorization: Bearer <access_token>` в заголовке. Публичные эндпоинты (регистрация, логин, восстановление пароля и т.д.) не требуют аутентификации.

#### 3.1.1. Регистрация и Вход
*   **`POST /register`**
    *   Описание: Регистрация нового пользователя.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "userRegistration",
            "attributes": {
              "username": "new_user",
              "email": "user@example.com",
              "password": "SecurePassword123!",
              "confirm_password": "SecurePassword123!"
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created):
        ```json
        {
          "data": {
            "type": "user",
            "id": "user-uuid-123",
            "attributes": {
              "username": "new_user",
              "email": "user@example.com",
              "status": "pending_verification",
              "created_at": "2024-03-15T10:00:00Z"
            }
          },
          "meta": {
            "message": "Регистрация прошла успешно. Пожалуйста, проверьте ваш email для верификации."
          }
        }
        ```
    *   Пример ответа (Ошибка 400 Validation Error):
        ```json
        {
          "errors": [
            {
              "status": "400",
              "code": "VALIDATION_ERROR",
              "title": "Ошибка валидации",
              "detail": "Пароли не совпадают.",
              "source": { "pointer": "/data/attributes/confirm_password" }
            }
          ]
        }
        ```
    *   Требуемые права доступа: Публичный.
*   **`POST /login`**
    *   Описание: Аутентификация пользователя и получение токенов.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "userLogin",
            "attributes": {
              "login": "user@example.com", // Может быть username или email
              "password": "SecurePassword123!"
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): (Refresh token обычно устанавливается в HttpOnly cookie)
        ```json
        {
          "data": {
            "type": "tokens",
            "attributes": {
              "access_token": "eyJhbGciOiJSUzI1NiI...",
              "token_type": "Bearer",
              "expires_in": 900, // секунды
              "user_id": "user-uuid-123",
              "username": "new_user",
              "roles": ["user"]
            }
          }
        }
        ```
    *   Пример ответа (Ошибка 401 Invalid Credentials):
        ```json
        {
          "errors": [
            {
              "status": "401",
              "code": "INVALID_CREDENTIALS",
              "title": "Неверные учетные данные",
              "detail": "Предоставлены неверный логин или пароль."
            }
          ]
        }
        ```
    *   Требуемые права доступа: Публичный.

#### 3.1.2. Управление токенами
*   **`POST /refresh-token`**
    *   Описание: Обновление access token с использованием refresh token (refresh token передается в HttpOnly cookie или в теле запроса, если это безопасно для типа клиента).
    *   Тело запроса (если refresh token в теле):
        ```json
        {
          "data": {
            "type": "refreshToken",
            "attributes": {
              "refresh_token": "your-long-refresh-token-string"
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): (Новый refresh token также может быть установлен в HttpOnly cookie)
        ```json
        {
          "data": {
            "type": "tokens",
            "attributes": {
              "access_token": "eyJhbGciOiJSUzI1NiI...", // Новый access token
              "token_type": "Bearer",
              "expires_in": 900
            }
          }
        }
        ```
    *   Требуемые права доступа: Требуется валидный refresh token.
*   **`POST /logout`**
    *   Описание: Отзыв текущего refresh token (и access token через JTI blacklist). Refresh token передается в HttpOnly cookie или в теле.
    *   Пример ответа (Успех 204 No Content): (Тело ответа отсутствует)
    *   Требуемые права доступа: Требуется валидный refresh token.

#### 3.1.3. Верификация Email
*   **`POST /verify-email`**
    *   Описание: Подтверждение адреса электронной почты с использованием кода верификации.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "emailVerification",
            "attributes": {
              "verification_code": "123456abcdef"
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "meta": {
            "message": "Email успешно подтвержден."
          }
        }
        ```
    *   Требуемые права доступа: Публичный или аутентифицированный пользователь (в зависимости от сценария).

#### 3.1.4. Управление Паролем
*   **`POST /forgot-password`**
    *   Описание: Запрос на сброс пароля. Отправляет код сброса на email пользователя.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "passwordForgot",
            "attributes": {
              "email": "user@example.com"
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "meta": {
            "message": "Если пользователь с таким email существует, на него будет отправлена инструкция по сбросу пароля."
          }
        }
        ```
    *   Требуемые права доступа: Публичный.

#### 3.1.5. Управление Двухфакторной Аутентификацией (2FA)
*   **`POST /me/2fa/totp/enable`**
    *   Описание: Инициирует включение TOTP 2FA, возвращает QR-код (или секрет).
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "totpEnableDetails",
            "attributes": {
              "secret": "BASE32ENCODEDSECRET",
              "qr_code_image_uri": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..." // Data URI QR-кода
            }
          }
        }
        ```
    *   Требуемые права доступа: Аутентифицированный пользователь.
*   **`POST /me/2fa/totp/verify`**
    *   Описание: Подтверждает включение TOTP 2FA с помощью кода из аутентификатора.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "totpVerification",
            "attributes": {
              "totp_code": "123456"
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "meta": {
            "message": "Двухфакторная аутентификация (TOTP) успешно включена."
          },
          "data": { // Возвращаем резервные коды
            "type": "backupCodes",
            "attributes": {
              "codes": ["abcdef01", "uvwxyz02", "..."]
            }
          }
        }
        ```
    *   Требуемые права доступа: Аутентифицированный пользователь (в процессе включения 2FA).

#### 3.1.6. Управление API Ключами (Пример)
*   **`POST /me/api-keys`**
    *   Описание: Создание нового API ключа для текущего пользователя.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "apiKeyCreation",
            "attributes": {
              "name": "My Application Key",
              "expires_at": "2025-12-31T23:59:59Z" // Опционально
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created):
        ```json
        {
          "data": {
            "type": "apiKey",
            "id": "apikey-uuid-789",
            "attributes": {
              "name": "My Application Key",
              "key_value": "secret_api_key_value_returned_only_once", // Это значение показывается только один раз
              "prefix": "pak_", // Префикс для идентификации ключа
              "created_at": "2024-03-15T10:00:00Z",
              "expires_at": "2025-12-31T23:59:59Z"
            }
          }
        }
        ```
    *   Требуемые права доступа: Аутентифицированный пользователь (например, `developer`).

### 3.2. gRPC API
*   Предназначен для внутреннего взаимодействия между микросервисами.
*   Определение Protobuf: `auth/v1/auth_service.proto`.
*   **Основные RPC Методы:**
    *   **`ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse)`**
        *   Описание: Проверяет валидность предоставленного access token.
        *   `message ValidateTokenRequest { string access_token = 1; }`
        *   `message ValidateTokenResponse { bool is_valid = 1; string user_id = 2; repeated string roles = 3; map<string, string> claims = 4; int64 expires_at_unix = 5; }`
    *   **`CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse)`**
        *   Описание: Проверяет, имеет ли пользователь (идентифицированный по user_id или токену) указанное разрешение.
        *   `message CheckPermissionRequest { string user_id = 1; string permission_code = 2; string access_token = 3; /* Либо user_id, либо access_token */ }`
        *   `message CheckPermissionResponse { bool granted = 1; }`
    *   **`GetUserInfo(GetUserInfoRequest) returns (GetUserInfoResponse)`**
        *   Описание: Получает основную информацию о пользователе по его ID или токену.
        *   `message GetUserInfoRequest { string user_id = 1; string access_token = 2; /* Либо user_id, либо access_token */ }`
        *   `message UserInfo { string user_id = 1; string username = 2; string email = 3; bool is_email_verified = 4; string status = 5; repeated string roles = 6; }`
        *   `message GetUserInfoResponse { UserInfo user_info = 1; }`
    *   **`GetJWKS(GetJWKSRequest) returns (GetJWKSResponse)`**
        *   Описание: Предоставляет набор публичных ключей JSON Web Key Set (JWKS) для верификации JWT подписей на стороне клиентов или других сервисов.
        *   `message GetJWKSRequest {}`
        *   `message JWK { string kty = 1; string use = 2; string kid = 3; string alg = 4; string n = 5; string e = 6; /* ... другие поля для EC и т.д. ... */ }`
        *   `message GetJWKSResponse { repeated JWK keys = 1; }`
    *   **`HealthCheck(HealthCheckRequest) returns (HealthCheckResponse)`**
        *   Описание: Стандартная проверка работоспособности сервиса.
        *   `message HealthCheckRequest {}`
        *   `message HealthCheckResponse { enum ServingStatus { UNKNOWN = 0; SERVING = 1; NOT_SERVING = 2; } ServingStatus status = 1; }`

### 3.3. WebSocket API
*   Не применимо для данного сервиса. Auth Service не предоставляет WebSocket API напрямую.

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности

*   **`User` (Пользователь)**
    *   `id` (UUID): Уникальный идентификатор. Пример: `a1b2c3d4-e5f6-7890-1234-567890abcdef`. Валидация: not null, unique. Обязательность: Required.
    *   `username` (VARCHAR(100)): Имя пользователя. Пример: `john_doe`. Валидация: not null, unique, min_length=3, max_length=100, alphanumeric. Обязательность: Required.
    *   `email` (VARCHAR(255)): Адрес электронной почты. Пример: `john.doe@example.com`. Валидация: not null, unique, valid email format. Обязательность: Required.
    *   `password_hash` (VARCHAR(255)): Хеш пароля (Argon2id). Пример: `$argon2id$v=19$m=65536,t=3,p=4$...`. Валидация: not null. Обязательность: Required.
    *   `status` (VARCHAR(50)): Статус пользователя (`pending_verification`, `active`, `blocked`, `deleted`). Пример: `active`. Валидация: not null, enum. Обязательность: Required.
    *   `is_2fa_enabled` (BOOLEAN): Включена ли двухфакторная аутентификация. Пример: `false`. Валидация: not null. Обязательность: Required.
    *   `email_verified_at` (TIMESTAMPTZ): Дата и время подтверждения email. Пример: `2024-01-15T10:30:00Z`. Обязательность: Optional.
    *   `blocked_at` (TIMESTAMPTZ): Дата и время блокировки. Обязательность: Optional.
    *   `last_login_at` (TIMESTAMPTZ): Дата и время последнего входа. Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ): Дата и время создания. Обязательность: Required.
    *   `updated_at` (TIMESTAMPTZ): Дата и время последнего обновления. Обязательность: Required.

*   **`Role` (Роль)**
    *   `id` (UUID): Уникальный идентификатор роли. Пример: `r1o2l3e4-a5b6-7890-1234-567890abcdef`. Валидация: not null, unique. Обязательность: Required.
    *   `name` (VARCHAR(50)): Имя роли (например, `user`, `admin`, `moderator`, `developer`). Пример: `admin`. Валидация: not null, unique. Обязательность: Required.
    *   `description` (TEXT): Описание роли. Пример: `Администратор платформы с полными правами`. Обязательность: Optional.
    *   `permissions` (JSONB): Список кодов разрешений, связанных с этой ролью. Пример: `["manage_users", "view_reports"]`. Валидация: array of strings. Обязательность: Optional.

*   **`Session` (Сессия пользователя)**
    *   `id` (UUID): Уникальный идентификатор сессии. Пример: `s1e2s3s4-i5o6-7890-1234-567890abcdef`. Валидация: not null, unique. Обязательность: Required.
    *   `user_id` (UUID, foreign key to User): ID пользователя. Валидация: not null. Обязательность: Required.
    *   `ip_address` (VARCHAR(45)): IP-адрес, с которого была создана сессия. Пример: `192.168.1.100`. Обязательность: Required.
    *   `user_agent` (TEXT): User-Agent клиента. Пример: `Mozilla/5.0 (...) Chrome/90.0`. Обязательность: Optional.
    *   `last_activity_at` (TIMESTAMPTZ): Время последней активности в сессии. Обязательность: Required.
    *   `created_at` (TIMESTAMPTZ): Время создания сессии. Обязательность: Required.
    *   `expires_at` (TIMESTAMPTZ): Время истечения сессии (если применимо, обычно для Redis сессий). Обязательность: Optional.

*   **`RefreshToken` (Токен обновления)**
    *   `id` (UUID): Уникальный идентификатор токена.
    *   `user_id` (UUID, foreign key to User): ID пользователя. Валидация: not null. Обязательность: Required.
    *   `token_hash` (VARCHAR(255)): Хеш самого refresh token (для безопасного хранения). Валидация: not null, unique. Обязательность: Required.
    *   `session_id` (UUID, foreign key to Session): ID сессии, к которой привязан токен. Обязательность: Optional.
    *   `user_agent` (TEXT): User-Agent на момент выпуска токена. Обязательность: Optional.
    *   `ip_address` (VARCHAR(45)): IP-адрес на момент выпуска токена. Обязательность: Optional.
    *   `is_revoked` (BOOLEAN): Отозван ли токен. Пример: `false`. Обязательность: Required.
    *   `expires_at` (TIMESTAMPTZ): Время истечения срока действия токена. Обязательность: Required.
    *   `created_at` (TIMESTAMPTZ): Время создания. Обязательность: Required.

*   **`APIKey` (API Ключ)**
    *   `id` (UUID): Уникальный идентификатор ключа.
    *   `user_id` (UUID, foreign key to User): ID пользователя-владельца ключа. Обязательность: Required.
    *   `name` (VARCHAR(100)): Имя ключа, задаваемое пользователем. Пример: `my_external_app_key`. Валидация: not null. Обязательность: Required.
    *   `prefix` (VARCHAR(8)): Префикс ключа (первые символы для идентификации). Пример: `pak_`. Валидация: not null, unique. Обязательность: Required.
    *   `key_hash` (VARCHAR(255)): Хеш API ключа. Валидация: not null, unique. Обязательность: Required.
    *   `permissions` (JSONB): Список разрешений, связанных с этим API ключом. Обязательность: Optional.
    *   `last_used_at` (TIMESTAMPTZ): Время последнего использования. Обязательность: Optional.
    *   `expires_at` (TIMESTAMPTZ): Время истечения срока действия ключа. Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ): Время создания. Обязательность: Required.

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL

**ERD Диаграмма (ключевые таблицы):**
```mermaid
erDiagram
    USERS {
        UUID id PK
        VARCHAR username UK
        VARCHAR email UK
        VARCHAR password_hash
        VARCHAR status
        BOOLEAN is_2fa_enabled
        TIMESTAMPTZ email_verified_at
        TIMESTAMPTZ last_login_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    SESSIONS {
        UUID id PK
        UUID user_id FK
        VARCHAR ip_address
        TEXT user_agent
        TIMESTAMPTZ last_activity_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ expires_at
    }
    REFRESH_TOKENS {
        UUID id PK
        UUID user_id FK
        VARCHAR token_hash UK
        UUID session_id FK "nullable"
        BOOLEAN is_revoked
        TIMESTAMPTZ expires_at
        TIMESTAMPTZ created_at
    }
    ROLES {
        UUID id PK
        VARCHAR name UK
        TEXT description
        JSONB permissions_list
    }
    USER_ROLES {
        UUID user_id PK FK
        UUID role_id PK FK
    }
    API_KEYS {
        UUID id PK
        UUID user_id FK
        VARCHAR name
        VARCHAR prefix UK
        VARCHAR key_hash UK
        JSONB permissions
        TIMESTAMPTZ last_used_at
        TIMESTAMPTZ expires_at
        TIMESTAMPTZ created_at
    }
    MFA_SECRETS {
        UUID user_id PK FK
        VARCHAR type -- 'totp', 'sms'
        VARCHAR secret_encrypted
        BOOLEAN is_verified
        TIMESTAMPTZ created_at
    }
    MFA_BACKUP_CODES {
        UUID id PK
        UUID user_id FK
        VARCHAR code_hash UK
        BOOLEAN is_used
        TIMESTAMPTZ created_at
    }
    VERIFICATION_CODES {
        UUID id PK
        UUID user_id FK
        VARCHAR type -- 'email_verification', 'password_reset'
        VARCHAR code_hash UK
        TIMESTAMPTZ expires_at
        TIMESTAMPTZ created_at
    }

    USERS ||--o{ SESSIONS : "has"
    USERS ||--o{ REFRESH_TOKENS : "has"
    USERS ||--o{ USER_ROLES : "has"
    ROLES ||--o{ USER_ROLES : "is_part_of"
    USERS ||--o{ API_KEYS : "owns"
    USERS ||--o{ MFA_SECRETS : "has"
    USERS ||--o{ MFA_BACKUP_CODES : "has"
    USERS ||--o{ VERIFICATION_CODES : "has"
    SESSIONS ||--o{ REFRESH_TOKENS : "can_be_linked_to"
```

**DDL (PostgreSQL - примеры ключевых таблиц):**
```sql
-- Расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица пользователей
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    -- соль уже включена в формат Argon2id, отдельное поле не нужно
    status VARCHAR(50) NOT NULL DEFAULT 'pending_verification' CHECK (status IN ('pending_verification', 'active', 'blocked', 'deleted')),
    is_2fa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    email_verified_at TIMESTAMPTZ,
    blocked_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- Таблица ролей
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL UNIQUE, -- e.g., 'user', 'admin', 'developer'
    description TEXT,
    permissions_list JSONB DEFAULT '[]'::jsonb, -- Список кодов разрешений
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Таблица связи пользователей и ролей (многие-ко-многим)
CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Таблица сессий (может также храниться в Redis для высокой производительности)
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ -- Для долгоживущих сессий, которые могут быть отозваны
);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Таблица refresh-токенов
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE, -- Хеш от самого refresh token
    session_id UUID REFERENCES sessions(id) ON DELETE SET NULL, -- Опциональная связь с сессией
    user_agent TEXT,
    ip_address VARCHAR(45),
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- Таблица API ключей
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    prefix VARCHAR(8) NOT NULL UNIQUE, -- Первые несколько символов ключа для быстрой идентификации
    key_hash VARCHAR(255) NOT NULL UNIQUE, -- Хеш от API ключа
    permissions JSONB DEFAULT '[]'::jsonb,
    last_used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);

-- TODO: Добавить DDL для mfa_secrets, mfa_backup_codes, verification_codes, audit_logs, external_accounts.
```

#### 4.2.2. Redis
*   **Сессии пользователей (опционально, если не используется PostgreSQL для сессий):**
    *   Ключ: `session:<session_id>`
    *   Значение: JSON строка или HASH с данными сессии (user_id, ip_address, user_agent, last_activity_at).
    *   TTL: Устанавливается в соответствии со сроком жизни сессии.
*   **Коды 2FA (TOTP, SMS/Email):**
    *   Ключ: `2fa_code:<user_id>:<type>` (type: `totp_setup`, `sms_login`, `email_login`)
    *   Значение: Код или секрет (для TOTP setup), количество попыток.
    *   TTL: Короткий (1-5 минут).
*   **JTI Blacklist (для отзыва Access Token):**
    *   Ключ: `jti_bl:<jti_access_token>`
    *   Значение: `true` или временная метка истечения access token.
    *   TTL: Равен времени жизни access token.
*   **Счетчики для Rate Limiting:**
    *   Ключ: `rl:<endpoint_or_action>:<user_id_or_ip>`
    *   Значение: Счетчик (integer).
    *   TTL: Зависит от окна ограничения (например, 1 минута, 1 час).
*   **Временные токены (сброс пароля, верификация email):**
    *   Ключ: `verification_token:<hashed_token>`
    *   Значение: `user_id:<user_id>,type:<password_reset|email_verify>`
    *   TTL: 1-24 часа.

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Apache Kafka.
*   **Формат событий:** CloudEvents v1.0 JSON (согласно `project_api_standards.md`).
*   **Основной топик для публикуемых событий:** `auth.events.v1`. (Может быть разделен на более гранулярные топики при необходимости).

*   **`auth.user.registered.v1`**
    *   Описание: Пользователь успешно зарегистрировался (но может требовать верификации email).
    *   Пример Payload (`data` секция CloudEvent):
        ```json
        {
          "user_id": "user-uuid-123",
          "username": "new_user",
          "email": "user@example.com",
          "status": "pending_verification",
          "registration_timestamp": "2024-03-15T10:00:00Z"
        }
        ```
*   **`auth.user.login_success.v1`**
    *   Описание: Пользователь успешно аутентифицировался.
    *   Пример Payload:
        ```json
        {
          "user_id": "user-uuid-123",
          "session_id": "session-uuid-abc",
          "ip_address": "192.168.1.100",
          "user_agent": "Mozilla/5.0...",
          "login_timestamp": "2024-03-15T10:05:00Z",
          "mfa_method_used": "none" // "totp", "sms", "backup_code"
        }
        ```
*   **`auth.session.revoked.v1`**
    *   Описание: Сессия пользователя была отозвана (logout, logout all, admin action).
    *   Пример Payload:
        ```json
        {
          "user_id": "user-uuid-123",
          "session_id": "session-uuid-abc", // или "all" для logout-all
          "revocation_reason": "user_logout", // "admin_action", "token_compromised"
          "revoked_at": "2024-03-15T11:00:00Z"
        }
        ```
*   **Другие события:** `auth.user.email_verified.v1`, `auth.user.password_reset_requested.v1`, `auth.user.password_changed.v1`, `auth.user.login_failed.v1`, `auth.user.account_locked.v1`, `auth.user.roles_changed.v1`, `auth.2fa.enabled.v1`, `auth.2fa.disabled.v1`, `auth.api_key.created.v1`, `auth.api_key.revoked.v1`.
    *   TODO: Детализировать структуру Payload для каждого из этих событий.

### 5.2. Потребляемые События (Consumed Events)
*   **`account.user.profile_updated.v1`** (из Account Service)
    *   Описание: Профиль пользователя был обновлен (например, изменен email, username, статус).
    *   Ожидаемый Payload (пример): `{"user_id": "user-uuid-123", "updated_fields": ["email", "status"], "new_email": "new_user@example.com", "new_status": "active"}`
    *   Логика обработки: Обновить кэшированную информацию о пользователе в Auth Service, если такая есть. Если изменен email, может потребоваться повторная верификация. Если статус изменен на `deleted` или `blocked` не через Auth Service, отозвать все сессии и токены.
*   **`admin.user.force_logout.v1`** (из Admin Service)
    *   Описание: Администратор принудительно завершил все сессии пользователя.
    *   Ожидаемый Payload (пример): `{"user_id": "user-uuid-123", "reason": "Suspicious activity reported"}`
    *   Логика обработки: Отозвать все активные сессии и refresh-токены для указанного `user_id`.
*   **`admin.user.block.v1`** (из Admin Service)
    *   Описание: Администратор заблокировал пользователя.
    *   Ожидаемый Payload (пример): `{"user_id": "user-uuid-123", "reason": "Violation of terms"}`
    *   Логика обработки: Установить статус пользователя `blocked`, отозвать все активные сессии и refresh-токены.
*   TODO: Детализировать структуру Payload и логику обработки для других потребляемых событий.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **API Gateway**:
    *   Тип интеграции: gRPC вызовы от API Gateway к Auth Service (`ValidateToken`, `CheckPermission`, `GetJWKS`).
    *   Назначение: Валидация JWT токенов и API ключей на уровне шлюза, получение публичных ключей для самостоятельной валидации токенов шлюзом. API Gateway также проксирует публичные REST эндпоинты Auth Service.
*   **Account Service**:
    *   Тип интеграции: Асинхронная через Kafka (Auth Service публикует `auth.user.registered`, Account Service подписывается). Может быть прямой gRPC вызов от Auth Service к Account Service для проверки существования пользователя перед регистрацией или для получения дополнительных данных профиля, если Auth Service не хранит их.
    *   Назначение: Создание профиля пользователя в Account Service после успешной регистрации в Auth Service. Реакция на изменения в профиле пользователя.
*   **Notification Service**:
    *   Тип интеграции: Асинхронная через Kafka (Auth Service публикует события типа `auth.user.verification_code_sent`, `auth.user.password_reset_code_sent`, `auth.2fa.code_sent`). Notification Service подписывается и отправляет email/SMS.
    *   Назначение: Отправка email и SMS сообщений для верификации, сброса пароля, кодов 2FA.
*   **Admin Service**:
    *   Тип интеграции: REST/gRPC API, предоставляемый Auth Service для Admin Service. Потребление событий Kafka от Admin Service.
    *   Назначение: Управление пользователями, ролями, сессиями, просмотр аудита через Admin Service. Реакция на административные действия (force_logout, block/unblock).
*   **Другие микросервисы платформы**:
    *   Тип интеграции: gRPC вызовы к Auth Service (`ValidateToken`, `CheckPermission`, `GetUserInfo`).
    *   Назначение: Аутентификация и авторизация межсервисных запросов и запросов от имени пользователей.

### 6.2. Внешние Системы
*   **OAuth Провайдеры (например, Telegram, VK, Odnoklassniki)**:
    *   Тип интеграции: HTTP редиректы и API вызовы по протоколам OAuth 2.0 / OIDC.
    *   Назначение: Внешняя аутентификация пользователей через сторонние сервисы.
*   **Платежные системы (косвенно)**: Могут влиять на статус подписки пользователя, что может отражаться на его ролях/разрешениях, управляемых Auth Service. Интеграция через Payment Service и Account Service.

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `AUTH_SERVICE_HTTP_PORT`: Порт для REST API (например, `8080`).
*   `AUTH_SERVICE_GRPC_PORT`: Порт для gRPC API (например, `9090`).
*   `AUTH_DB_HOST`, `AUTH_DB_PORT`, `AUTH_DB_USER`, `AUTH_DB_PASSWORD`, `AUTH_DB_NAME`, `AUTH_DB_SSL_MODE`: Параметры подключения к PostgreSQL.
*   `AUTH_REDIS_ADDR`, `AUTH_REDIS_PASSWORD`, `AUTH_REDIS_DB`: Параметры подключения к Redis.
*   `AUTH_KAFKA_BROKERS`: Список брокеров Kafka (например, `kafka1:9092,kafka2:9092`).
*   `AUTH_KAFKA_TOPIC_AUTH_EVENTS`: Имя топика для публикуемых событий (например, `auth.events.v1`).
*   `AUTH_JWT_PRIVATE_KEY_PATH`: Путь к файлу с приватным RSA ключом для подписи JWT.
*   `AUTH_JWT_PUBLIC_KEY_PATH`: Путь к файлу с публичным RSA ключом для проверки JWT (также используется для JWKS эндпоинта).
*   `AUTH_JWT_ACCESS_TOKEN_TTL_SECONDS`: Срок жизни access token в секундах (например, `900` для 15 минут).
*   `AUTH_JWT_REFRESH_TOKEN_TTL_SECONDS`: Срок жизни refresh token в секундах (например, `2592000` для 30 дней).
*   `AUTH_ARGON2ID_MEMORY_KB`, `AUTH_ARGON2ID_ITERATIONS`, `AUTH_ARGON2ID_PARALLELISM`, `AUTH_ARGON2ID_KEY_LENGTH`, `AUTH_ARGON2ID_SALT_LENGTH`: Параметры для хеширования паролей Argon2id.
*   `AUTH_EMAIL_VERIFICATION_CODE_TTL_MINUTES`: Срок жизни кода верификации email.
*   `AUTH_PASSWORD_RESET_CODE_TTL_MINUTES`: Срок жизни кода сброса пароля.
*   `AUTH_2FA_TOTP_ISSUER_NAME`: Имя издателя для TOTP (например, "Моя Платформа").
*   `AUTH_2FA_CODE_TTL_SECONDS`: Срок жизни кодов 2FA (для SMS/Email).
*   `LOG_LEVEL`: Уровень логирования (`debug`, `info`, `warn`, `error`).
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`: Эндпоинт Jaeger для экспорта трейсов.
*   OAuth provider client IDs and secrets (для VK, Telegram и т.д., например `AUTH_VK_CLIENT_ID`, `AUTH_VK_CLIENT_SECRET`).
*   `CORS_ALLOWED_ORIGINS`: Список разрешенных источников для CORS.

### 7.2. Файлы Конфигурации (если применимо)
*   **`config/config.yaml` (опционально):** Может использоваться для задания структуры настроек по умолчанию, которые затем переопределяются переменными окружения. Например, параметры Argon2id, таймауты, настройки CORS.
    ```yaml
    server:
      http_port: ${AUTH_SERVICE_HTTP_PORT:-8080}
      grpc_port: ${AUTH_SERVICE_GRPC_PORT:-9090}
    jwt:
      access_token_ttl_seconds: ${AUTH_JWT_ACCESS_TOKEN_TTL_SECONDS:-900}
      refresh_token_ttl_seconds: ${AUTH_JWT_REFRESH_TOKEN_TTL_SECONDS:-2592000}
      private_key_path: ${AUTH_JWT_PRIVATE_KEY_PATH:-/run/secrets/jwt_private.pem}
      public_key_path: ${AUTH_JWT_PUBLIC_KEY_PATH:-/app/jwt_public.pem} # Публичный ключ может быть вшит в образ
    password_hashing:
      argon2id:
        memory_kb: ${AUTH_ARGON2ID_MEMORY_KB:-65536}
        iterations: ${AUTH_ARGON2ID_ITERATIONS:-3}
        parallelism: ${AUTH_ARGON2ID_PARALLELISM:-4}
        salt_length: ${AUTH_ARGON2ID_SALT_LENGTH:-16}
        key_length: ${AUTH_ARGON2ID_KEY_LENGTH:-32}
    # ... другие настройки
    ```
*   **JWKS (JSON Web Key Set):** Публичные ключи для проверки JWT подписей могут предоставляться через эндпоинт `/jwks.json` (генерируется на основе `AUTH_JWT_PUBLIC_KEY_PATH`) и/или храниться в файле, доступном сервису.

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   REST API: Используются стандартные коды состояния HTTP. Тело ответа об ошибке соответствует формату, определенному в `project_api_standards.md` (см. секцию 3.1).
*   gRPC API: Используются стандартные коды состояния gRPC. Дополнительная информация об ошибке передается через `google.rpc.Status` и `google.rpc.ErrorInfo`.
*   Все ошибки логируются с `trace_id` (из OpenTelemetry) и, если применимо, `request_id`.
*   Не раскрывать чувствительную информацию в сообщениях об ошибках для конечных пользователей.

### 8.2. Распространенные Коды Ошибок (для REST API)
*   **`400 Bad Request` (`VALIDATION_ERROR`)**: Некорректные входные данные (например, неверный формат email, короткий пароль).
*   **`401 Unauthorized` (`INVALID_CREDENTIALS`, `TOKEN_EXPIRED`, `INVALID_TOKEN`, `INVALID_2FA_CODE`)**: Ошибка аутентификации.
*   **`403 Forbidden` (`EMAIL_NOT_VERIFIED`, `USER_BLOCKED`, `INSUFFICIENT_PERMISSIONS`)**: Доступ запрещен.
*   **`404 Not Found` (`USER_NOT_FOUND`, `RESOURCE_NOT_FOUND`)**: Запрашиваемый ресурс не найден.
*   **`409 Conflict` (`USERNAME_ALREADY_EXISTS`, `EMAIL_ALREADY_EXISTS`)**: Конфликт данных.
*   **`429 Too Many Requests` (`RATE_LIMIT_EXCEEDED`)**: Превышен лимит запросов.
*   **`500 Internal Server Error` (`INTERNAL_ERROR`)**: Внутренняя ошибка сервера.

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   **Механизмы:** Логин/пароль, JWT (Access Token + Refresh Token), OAuth 2.0/OIDC (например, Telegram, VK), API ключи.
*   **JWT:** Алгоритм RS256. Access Token имеет короткий срок жизни (например, 15 минут). Refresh Token имеет длительный срок жизни (например, 30 дней), хранится в HttpOnly cookie для веб-клиентов, подвержен ротации и может быть отозван.
*   **Двухфакторная Аутентификация (2FA):** TOTP (RFC 6238) как основной метод. SMS/Email коды как альтернатива. Поддержка резервных кодов.

### 9.2. Авторизация
*   **Модель:** Role-Based Access Control (RBAC). Роли назначаются пользователям. Разрешения привязаны к ролям.
*   **Проверка прав:** Может выполняться как на API Gateway (грубая проверка по наличию роли), так и в каждом микросервисе (точная проверка конкретного разрешения через gRPC вызов к Auth Service - метод `CheckPermission`). JWT содержит информацию о ролях и ключевых разрешениях пользователя.

### 9.3. Защита Данных
*   **Хеширование паролей:** Используется стойкий алгоритм Argon2id с настраиваемыми параметрами.
*   **Шифрование:** Секреты для TOTP и другие чувствительные данные конфигурации шифруются при хранении (например, с использованием Vault или шифрованных Kubernetes Secrets). Мастер-ключ для шифрования управляется отдельно.
*   **Транспортная безопасность:** TLS 1.2+ для всех внешних и внутренних коммуникаций (REST, gRPC, PostgreSQL, Redis, Kafka).
*   **Защита от атак:**
    *   Rate limiting для предотвращения brute-force атак на эндпоинты логина, сброса пароля, проверки кодов.
    *   Блокировка аккаунтов после нескольких неудачных попыток входа.
    *   Использование CAPTCHA (планируется) для публичных эндпоинтов.
    *   Защита от CSRF/XSS для любых веб-интерфейсов, связанных с Auth Service (например, страницы логина, если они являются частью сервиса).
    *   Предотвращение user enumeration.
*   **Соответствие ФЗ-152 "О персональных данных":** Локализация баз данных на территории РФ, получение согласия на обработку ПДн, минимизация собираемых данных, обеспечение прав субъектов ПДн.

### 9.4. Управление Секретами
*   Приватные RSA ключи для подписи JWT, секреты OAuth клиентов, пароли к базам данных, мастер-ключ для шифрования конфигураций должны храниться в безопасном хранилище секретов (например, HashiCorp Vault или зашифрованные Kubernetes Secrets).
*   Плановая ротация ключей и секретов.

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Многоэтапная сборка на основе официального образа Go и Alpine Linux для минимизации размера и уязвимостей.
*   **Kubernetes манифесты/Helm-чарты:** Включают Deployment, Service, ConfigMap, Secret, HorizontalPodAutoscaler (HPA), PodDisruptionBudget (PDB), NetworkPolicy.
*   (Ссылка на `project_deployment_standards.md` и репозиторий GitOps).

### 10.2. Зависимости при Развертывании
*   Кластер Kubernetes.
*   PostgreSQL (может быть развернут как StatefulSet или использоваться управляемый сервис).
*   Redis (может быть развернут как StatefulSet или использоваться управляемый сервис).
*   Apache Kafka.
*   Доступность API Gateway для маршрутизации запросов.
*   Доступность Notification Service (для отправки кодов) и Account Service (для координации при регистрации).
*   Настроенное хранилище секретов.

### 10.3. CI/CD
*   Автоматизированные пайплайны (например, GitLab CI, Jenkins, GitHub Actions) для:
    *   Сборки Docker-образа.
    *   Запуска юнит-тестов и интеграционных тестов.
    *   Сканирования на уязвимости (SAST, DAST, сканирование образов).
    *   Развертывания в различные окружения (dev, staging, production) с использованием Helm-чартов и GitOps-подхода (ArgoCD/Flux).
*   Стратегии развертывания: RollingUpdate, Canary (опционально).

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат:** Структурированные логи в формате JSON (с использованием библиотеки Zap).
*   **Ключевые события для логирования:** Входящие запросы (с `request_id`), результаты аутентификации и авторизации, ошибки, события жизненного цикла токенов, административные действия, события безопасности (например, неудачные попытки входа, блокировки).
*   **Уровни логирования:** `DEBUG`, `INFO`, `WARN`, `ERROR`. Уровень настраивается через переменные окружения.
*   **Интеграция:** Сбор логов через FluentBit/Vector и отправка в централизованную систему логирования (например, Loki, ELK Stack).
*   (Ссылка на `project_observability_standards.md`).

### 11.2. Мониторинг
*   **Метрики (Prometheus):** Экспортируются через эндпоинт `/metrics`.
    *   `auth_requests_total{method, path, status_code}`: Количество HTTP запросов.
    *   `auth_request_duration_seconds{method, path}`: Гистограмма или summary длительности HTTP запросов.
    *   `auth_grpc_requests_total{service, method, status_code}`: Количество gRPC запросов.
    *   `auth_grpc_request_duration_seconds{service, method}`: Длительность gRPC запросов.
    *   `auth_token_validation_total{status}`: Результаты валидации токенов (success, failure).
    *   `auth_active_sessions_count`: Количество активных сессий (если отслеживается).
    *   `auth_db_connection_pool_status`: Состояние пула соединений к БД.
    *   Стандартные метрики Go-приложения (goroutines, GC, etc.).
*   **Дашборды (Grafana):** Визуализация ключевых метрик производительности, ошибок, событий безопасности.
*   **Алертинг (AlertManager):** Настроены алерты для критических ситуаций: высокий процент ошибок (>5% за 5 мин), значительное увеличение времени ответа, ошибки валидации токенов, проблемы с подключением к БД/Redis/Kafka.
*   (Ссылка на `project_observability_standards.md`).

### 11.3. Трассировка
*   **Интеграция:** OpenTelemetry SDK для Go. Экспорт трейсов в Jaeger или Tempo.
*   **Контекст трассировки:** Распространение `trace_id` через все входящие и исходящие запросы (HTTP, gRPC, Kafka).
*   Трассируются ключевые операции: аутентификация, генерация токенов, взаимодействие с БД и другими сервисами.
*   (Ссылка на `project_observability_standards.md`).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:**
    *   Время ответа для эндпоинта логина (включая проверку пароля и генерацию токенов): P95 < 150 мс, P99 < 300 мс при нагрузке 1000 RPS.
    *   Время ответа для валидации access token (gRPC `ValidateToken`): P95 < 20 мс, P99 < 50 мс при нагрузке 5000 RPS.
    *   Время генерации нового access token по refresh token: P95 < 50 мс.
*   **Масштабируемость:** Горизонтальное масштабирование для обработки до 5000 RPS на эндпоинтах аутентификации и валидации токенов. Поддержка до 10 миллионов зарегистрированных пользователей и 1 миллиона активных сессий.
*   **Надежность:** Доступность сервиса: 99.98% (допустимое время простоя ~1.75 часа в год). RTO (Recovery Time Objective) < 5 минут. RPO (Recovery Point Objective) < 1 минуты для критичных данных (транзакции по токенам, сессии).
*   **Безопасность:** Соответствие OWASP Top 10, использование Argon2id для хеширования паролей, RS256 для JWT. Регулярные аудиты безопасности. Соответствие требованиям ФЗ-152.
*   **Сопровождаемость:** Покрытие кода юнит-тестами > 85%. Наличие интеграционных тестов для ключевых сценариев. Логирование и метрики должны обеспечивать полную наблюдаемость.

## 13. Приложения (Appendices)
*   Все детальные примеры API запросов/ответов, схемы Protobuf, JSON Schemas для событий и моделей данных, а также детальные DDL теперь интегрированы в соответствующие разделы этого документа.

---
*Этот документ является основной спецификацией для Auth Service и должен поддерживаться в актуальном состоянии.*

## 14. Связанные Рабочие Процессы (Related Workflows)
*   [User Registration and Initial Profile Setup](../../../project_workflows/user_registration_flow.md)
