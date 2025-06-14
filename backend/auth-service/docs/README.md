<!-- backend\auth-service\docs\README.md -->
# Спецификация Микросервиса: Auth Service

**Версия:** 2.0
**Дата последнего обновления:** 2024-07-16

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение и Роль:** Auth Service является центральным компонентом платформы "Российский Аналог Steam", отвечающим за управление идентификацией пользователей, аутентификацию и авторизацию. Он обеспечивает безопасный, масштабируемый и надежный фундамент для модели безопасности всей платформы. Его основная цель - проверка подлинности пользователей, контроль доступа к ресурсам на основе ролей и разрешений, управление сессиями и токенами (JWT), а также предоставление доверенной информации об аутентифицированных пользователях другим микросервисам.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   Регистрация пользователей (в координации с Account Service).
*   Аутентификация по логину/паролю.
*   Двухфакторная аутентификация (2FA): TOTP, SMS/Email коды (через Notification Service), резервные коды.
*   Внешняя аутентификация: OAuth 2.0 / OIDC (включая приоритетную поддержку российских провайдеров, таких как **ВКонтакте (VK ID)** и **Telegram Login**, а также других по мере необходимости).
*   Управление JSON Web Token (JWT): Генерация (Access Token RS256, Refresh Token), валидация, ротация, отзыв (через JTI blacklist).
*   Управление сессиями: Отслеживание и отзыв активных сессий пользователей.
*   Управление паролями: Безопасное хранение (хеширование Argon2id), сброс и изменение пароля.
*   Подтверждение Email.
*   Role-Based Access Control (RBAC): Управление ролями и разрешениями пользователей.
*   Управление API ключами для внешних и внутренних сервисов.
*   Аудит событий безопасности, связанных с аутентификацией и авторизацией.
*   Обнаружение подозрительной активности (например, частые неудачные попытки входа, вход с подозрительных IP).
*   Предоставление информации о пользователе и его правах другим сервисам.
*   Административные функции для управления пользователями и ролями (через Admin Service).

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+, согласно `../../../../project_technology_stack.md`).
*   **API Фреймворки:**
    *   REST: Gin (`github.com/gin-gonic/gin`) или Echo (`github.com/labstack/echo/v4`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
    *   gRPC: `google.golang.org/grpc` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Базы данных:**
    *   PostgreSQL (версия 15+): для хранения пользовательских данных, ролей, разрешений, refresh-токенов, API ключей, кодов верификации, секретов MFA. (согласно `../../../../project_technology_stack.md`). Драйвер: GORM (`gorm.io/gorm`) с `gorm.io/driver/postgres` или `pgx` (`github.com/jackc/pgx/v5`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
    *   Redis (версия 7.0+): для кэширования сессий, временных кодов 2FA, JTI blacklist для JWT, счетчиков rate limiting. Клиент: `go-redis/redis` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Брокер сообщений/События:** Apache Kafka (клиент `github.com/confluentinc/confluent-kafka-go` или `github.com/segmentio/kafka-go`, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Безопасность (криптография):**
    *   JWT: `github.com/golang-jwt/jwt/v5` (алгоритм RS256).
    *   Хеширование паролей: `golang.org/x/crypto/argon2` (Argon2id).
    *   TOTP: Стандартные библиотеки Go (например, `github.com/pquerna/otp`).
*   **Логирование:** Zap (`go.uber.org/zap`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Мониторинг:** Prometheus (`github.com/prometheus/client_golang`) (согласно `../../../../project_observability_standards.md`).
*   **Трассировка:** OpenTelemetry (`go.opentelemetry.io/otel`) (согласно `../../../../project_observability_standards.md`).
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Контейнеризация:** Docker.
*   **Оркестрация:** Kubernetes.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   Для используемых терминов (Access Token, Refresh Token, JWT, RBAC, 2FA, TOTP, Argon2id, OAuth, OIDC, JWKS, JTI, User, Role, Permission, Session, API Key и др.) см. `../../../../project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Auth Service разработан как stateless (не хранящий состояние между запросами) микросервис на языке Go. Состояние сессий и другая временная информация хранится в Redis.
*   Он взаимодействует с PostgreSQL для постоянного хранения данных (учетные записи, роли, разрешения, токены обновления и т.д.), Redis для кэширования и временных данных, и Kafka для асинхронного обмена доменными событиями.
*   Архитектура ориентирована на безопасность, масштабируемость, отказоустойчивость и сопровождаемость, следуя принципам многослойной архитектуры.

**Диаграмма верхнеуровневой архитектуры:**
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
*   **Ответственность:** Обработка входящих HTTP (RESTful через Gin/Echo) и gRPC запросов, валидация входных данных (DTO), преобразование данных и вызов соответствующей бизнес-логики в Application Layer.
*   **Ключевые компоненты/модули:**
    *   **HTTP Handlers (Gin/Echo):** Контроллеры для каждого REST эндпоинта (например, `POST /register`, `POST /login`).
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
*   **Формат данных:** JSON. Стандартный формат ответа об ошибке соответствует `../../../../project_api_standards.md`:
    ```json
    {
      "errors": [
        {
          "code": "ERROR_CODE_UPPER_SNAKE_CASE",
          "title": "Краткое описание ошибки на русском",
          "detail": "Полное описание ошибки с контекстом, если применимо.",
          "source": { "pointer": "/data/attributes/field_name" }
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
    *   Пример ответа (Ошибка 400 Validation Error - стандартизированный):
        ```json
        {
          "errors": [
            {
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
              "login": "user@example.com",
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
              "expires_in": 900,
              "user_id": "user-uuid-123",
              "username": "new_user",
              "roles": ["user"]
            }
          }
        }
        ```
    *   Пример ответа (Ошибка 401 Invalid Credentials - стандартизированный):
        ```json
        {
          "errors": [
            {
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
    *   Описание: Обновление access token с использованием refresh token.
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
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "tokens",
            "attributes": {
              "access_token": "eyJhbGciOiJSUzI1NiI...",
              "token_type": "Bearer",
              "expires_in": 900
            }
          }
        }
        ```
    *   Требуемые права доступа: Требуется валидный refresh token.
*   **`POST /logout`**
    *   Описание: Отзыв текущего refresh token.
    *   Пример ответа (Успех 204 No Content):
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
    *   Требуемые права доступа: Публичный или аутентифицированный пользователь.

#### 3.1.4. Управление Паролем
*   **`POST /forgot-password`**
    *   Описание: Запрос на сброс пароля.
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
    *   Описание: Инициирует включение TOTP 2FA.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "totpEnableDetails",
            "attributes": {
              "secret": "BASE32ENCODEDSECRET",
              "qr_code_image_uri": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
            }
          }
        }
        ```
    *   Требуемые права доступа: Аутентифицированный пользователь.
*   **`POST /me/2fa/totp/verify`**
    *   Описание: Подтверждает включение TOTP 2FA.
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
          "data": {
            "type": "backupCodes",
            "attributes": {
              "codes": ["abcdef01", "uvwxyz02", "..."] // Резервные коды генерируются и возвращаются ОДИН РАЗ
            }
          }
        }
        ```
    *   Требуемые права доступа: Аутентифицированный пользователь.
*   **`POST /me/2fa/disable`**
    *   Описание: Отключение 2FA для пользователя (требует текущего 2FA кода или пароля в зависимости от политики).
    *   Тело запроса: `{"data": {"type": "2faDisableRequest", "attributes": {"verification_code": "123456_or_password"}}}`
    *   Требуемые права доступа: Аутентифицированный пользователь.

#### 3.1.6. Управление Сессиями
*   **`GET /me/sessions`**
    *   Описание: Получение списка активных сессий пользователя.
    *   Требуемые права доступа: Аутентифицированный пользователь.
*   **`DELETE /me/sessions/{session_id}`**
    *   Описание: Отзыв конкретной сессии пользователя (кроме текущей, если не предусмотрен специальный механизм).
    *   Требуемые права доступа: Аутентифицированный пользователь.
*   **`DELETE /me/sessions/all-others`**
    *   Описание: Отзыв всех сессий пользователя, кроме текущей.
    *   Требуемые права доступа: Аутентифицированный пользователь.

#### 3.1.7. Управление API Ключами (Пример)
*   **`POST /me/api-keys`**
    *   Описание: Создание нового API ключа.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "apiKeyCreation",
            "attributes": {
              "name": "My Application Key",
              "expires_at": "2025-12-31T23:59:59Z"
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
              "key_value": "secret_api_key_value_returned_only_once",
              "prefix": "pak_",
              "created_at": "2024-03-15T10:00:00Z",
              "expires_at": "2025-12-31T23:59:59Z"
            }
          }
        }
        ```
    *   Требуемые права доступа: Аутентифицированный пользователь (например, `developer`).

### 3.2. gRPC API
(Описание остается как в существующем документе).

### 3.3. WebSocket API
*   Не применимо для данного сервиса.

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности
*   **`User` (Пользователь)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `username` (VARCHAR, UK). **Обязательность: Да.**
    *   `email` (VARCHAR, UK). **Обязательность: Да.**
    *   `password_hash` (VARCHAR). **Обязательность: Да.**
    *   `status` (VARCHAR, ENUM: `pending_verification`, `active`, `blocked`, `deleted`). **Обязательность: Да (DEFAULT 'pending_verification').**
    *   `is_2fa_enabled` (BOOLEAN). **Обязательность: Да (DEFAULT FALSE).**
    *   `email_verified_at` (TIMESTAMPTZ, Nullable). **Обязательность: Нет.**
    *   `last_login_at` (TIMESTAMPTZ, Nullable). **Обязательность: Нет.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
*   **`Role` (Роль)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `name` (VARCHAR, UK). **Обязательность: Да.**
    *   `description` (TEXT, Nullable). **Обязательность: Нет.**
    *   `permissions` (JSONB): Список разрешений для роли. **Обязательность: Да (DEFAULT '[]').**
*   **`Session` (Сессия пользователя)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK to User). **Обязательность: Да.**
    *   `ip_address` (VARCHAR(45), Nullable). **Обязательность: Нет.**
    *   `user_agent` (TEXT, Nullable). **Обязательность: Нет.**
    *   `last_activity_at` (TIMESTAMPTZ). **Обязательность: Да.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `expires_at` (TIMESTAMPTZ). **Обязательность: Да.**
*   **`RefreshToken` (Токен обновления)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK to User). **Обязательность: Да.**
    *   `token_hash` (VARCHAR, UK). **Обязательность: Да.**
    *   `session_id` (UUID, FK to Session, Nullable). **Обязательность: Нет.**
    *   `is_revoked` (BOOLEAN). **Обязательность: Да (DEFAULT FALSE).**
    *   `expires_at` (TIMESTAMPTZ). **Обязательность: Да.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
*   **`APIKey` (API Ключ)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK to User, Nullable): Пользователь, которому принадлежит ключ (если это пользовательский ключ). **Обязательность: Нет.**
    *   `service_name` (VARCHAR, Nullable): Имя сервиса, которому принадлежит ключ (если это ключ сервиса). **Обязательность: Нет.**
    *   `name` (VARCHAR): Название ключа для идентификации. **Обязательность: Да.**
    *   `prefix` (VARCHAR(8), UK): Префикс ключа (первые символы для быстрой идентификации). **Обязательность: Да.**
    *   `key_hash` (VARCHAR, UK): Хеш самого ключа. **Обязательность: Да.**
    *   `permissions` (JSONB): Разрешения, связанные с ключом. **Обязательность: Да (DEFAULT '[]').**
    *   `last_used_at` (TIMESTAMPTZ, Nullable). **Обязательность: Нет.**
    *   `expires_at` (TIMESTAMPTZ, Nullable). **Обязательность: Нет.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `is_active` (BOOLEAN). **Обязательность: Да (DEFAULT TRUE).**
*   **`MFASecret` (Секрет MFA)**
    *   `user_id` (UUID, PK, FK to User): ID пользователя. **Обязательность: Да.**
    *   `type` (VARCHAR(10)): Тип MFA ('totp', 'sms'). **Обязательность: Да.**
    *   `secret_encrypted` (TEXT): Зашифрованный секрет. **Обязательность: Да.**
    *   `is_verified` (BOOLEAN): Подтвержден ли данный метод MFA. **Обязательность: Да (DEFAULT FALSE).**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
*   **`MFABackupCode` (Резервный код MFA)**
    *   `id` (UUID, PK): Уникальный ID кода. **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK to User): ID пользователя. **Обязательность: Да.**
    *   `code_hash` (VARCHAR(255)): Хеш резервного кода. **Обязательность: Да.**
    *   `is_used` (BOOLEAN): Использован ли данный код. **Обязательность: Да (DEFAULT FALSE).**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `used_at` (TIMESTAMPTZ, Nullable): Время использования кода. **Обязательность: Нет.**
*   **`VerificationCode` (Код верификации)**
    *   `id` (UUID, PK): Уникальный ID кода. **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK to User): ID пользователя. **Обязательность: Да.**
    *   `type` (VARCHAR(50)): Тип кода. **Обязательность: Да.**
    *   `code_hash` (VARCHAR(255)): Хеш кода верификации. **Обязательность: Да.**
    *   `target` (VARCHAR(255)): Цель верификации. **Обязательность: Да.**
    *   `expires_at` (TIMESTAMPTZ): Время истечения срока действия кода. **Обязательность: Да.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
*   **`AuditLogAuth` (Лог аудита Auth Service)**
    *   `id` (UUID, PK): Уникальный ID записи лога. **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, Nullable, FK to User): ID пользователя. **Обязательность: Нет.**
    *   `actor_id` (UUID, Nullable, FK to User): ID субъекта. **Обязательность: Нет.**
    *   `action_type` (VARCHAR(100)): Тип действия. **Обязательность: Да.**
    *   `ip_address` (VARCHAR(45)): IP-адрес. **Обязательность: Нет.**
    *   `user_agent` (TEXT): User-Agent. **Обязательность: Нет.**
    *   `details` (JSONB): Дополнительные детали. **Обязательность: Нет.**
    *   `timestamp` (TIMESTAMPTZ): Время события. **Обязательность: Да (генерируется БД).**
*   **`ExternalAccount` (Внешний аккаунт OAuth/OIDC)**
    *   `id` (UUID, PK): Уникальный ID связи. **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK to User): ID пользователя в нашей системе. **Обязательность: Да.**
    *   `provider_name` (VARCHAR(50)): Имя провайдера. **Обязательность: Да.**
    *   `provider_user_id` (VARCHAR(255)): Уникальный ID пользователя у внешнего провайдера. **Обязательность: Да.**
    *   `provider_user_details` (JSONB, Nullable): Дополнительные детали от провайдера. **Обязательность: Нет.**
    *   `access_token_encrypted` (TEXT, Nullable): Зашифрованный access token от провайдера. **Обязательность: Нет.**
    *   `refresh_token_encrypted` (TEXT, Nullable): Зашифрованный refresh token от провайдера. **Обязательность: Нет.**
    *   `token_expires_at` (TIMESTAMPTZ, Nullable): Время истечения access token от провайдера. **Обязательность: Нет.**
    *   `linked_at` (TIMESTAMPTZ): Время привязки аккаунта. **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления данных от провайдера. **Обязательность: Да (генерируется БД).**

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL
**ERD Диаграмма (дополненная):**
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
        VARCHAR type "ENUM('totp', 'sms')"
        TEXT secret_encrypted
        BOOLEAN is_verified
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    MFA_BACKUP_CODES {
        UUID id PK
        UUID user_id FK
        VARCHAR code_hash UK
        BOOLEAN is_used
        TIMESTAMPTZ created_at
        TIMESTAMPTZ used_at "nullable"
    }
    VERIFICATION_CODES {
        UUID id PK
        UUID user_id FK
        VARCHAR type "ENUM('email_verification', 'password_reset', 'sms_verification', 'generic_otp')"
        VARCHAR code_hash UK
        VARCHAR target "email or phone"
        TIMESTAMPTZ expires_at
        TIMESTAMPTZ created_at
    }
    AUDIT_LOGS_AUTH {
        UUID id PK
        UUID user_id FK "nullable"
        UUID actor_id FK "nullable"
        VARCHAR action_type
        VARCHAR ip_address
        TEXT user_agent
        JSONB details
        TIMESTAMPTZ timestamp
    }
    EXTERNAL_ACCOUNTS {
        UUID id PK
        UUID user_id FK
        VARCHAR provider_name
        VARCHAR provider_user_id
        JSONB provider_user_details "nullable"
        TEXT access_token_encrypted "nullable"
        TEXT refresh_token_encrypted "nullable"
        TIMESTAMPTZ token_expires_at "nullable"
        TIMESTAMPTZ linked_at
        TIMESTAMPTZ updated_at
        UNIQUE (provider_name, provider_user_id)
    }

    USERS ||--o{ SESSIONS : "has"
    USERS ||--o{ REFRESH_TOKENS : "has"
    USERS ||--o{ USER_ROLES : "has"
    ROLES ||--o{ USER_ROLES : "is_part_of"
    USERS ||--o{ API_KEYS : "owns"
    USERS ||--o{ MFA_SECRETS : "has_one_per_type"
    USERS ||--o{ MFA_BACKUP_CODES : "has_many"
    USERS ||--o{ VERIFICATION_CODES : "has_many"
    USERS ||--o{ AUDIT_LOGS_AUTH : "subject_of (nullable)"
    USERS ||--o{ AUDIT_LOGS_AUTH : "actor_of (nullable)"
    USERS ||--o{ EXTERNAL_ACCOUNTS : "has_many"
    SESSIONS ||--o{ REFRESH_TOKENS : "can_be_linked_to"
```

**DDL (PostgreSQL - дополнения для таблиц из TODO):**
```sql
-- Таблица секретов MFA
CREATE TABLE mfa_secrets (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL CHECK (type IN ('totp', 'sms')),
    secret_encrypted TEXT NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    -- CONSTRAINT uq_user_mfa_type UNIQUE (user_id, type) -- Если разрешен только один метод каждого типа
);
COMMENT ON TABLE mfa_secrets IS 'Хранит зашифрованные секреты для двухфакторной аутентификации.';

-- Таблица резервных кодов MFA
CREATE TABLE mfa_backup_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash VARCHAR(255) NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    used_at TIMESTAMPTZ,
    UNIQUE (user_id, code_hash) -- Код должен быть уникален для пользователя
);
CREATE INDEX idx_mfa_backup_codes_user_id ON mfa_backup_codes(user_id);
COMMENT ON TABLE mfa_backup_codes IS 'Резервные коды для двухфакторной аутентификации.';

-- Таблица кодов верификации (email, сброс пароля, и т.д.)
CREATE TABLE verification_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('email_verification', 'password_reset', 'sms_verification', 'generic_otp')),
    code_hash VARCHAR(255) NOT NULL,
    target VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (type, target, code_hash) -- Комбинация должна быть уникальна, или просто code_hash если он глобально уникален
);
CREATE INDEX idx_verification_codes_user_id_type ON verification_codes(user_id, type);
CREATE INDEX idx_verification_codes_expires_at ON verification_codes(expires_at);
CREATE INDEX idx_verification_codes_target_type ON verification_codes(target, type);
COMMENT ON TABLE verification_codes IS 'Коды для различных операций верификации.';

-- Таблица логов аудита Auth Service
CREATE TABLE audit_logs_auth (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action_type VARCHAR(100) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    details JSONB,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_audit_logs_auth_user_id ON audit_logs_auth(user_id);
CREATE INDEX idx_audit_logs_auth_action_type ON audit_logs_auth(action_type);
CREATE INDEX idx_audit_logs_auth_timestamp ON audit_logs_auth(timestamp);
COMMENT ON TABLE audit_logs_auth IS 'Логи аудита событий аутентификации и авторизации.';

-- Таблица для связанных внешних аккаунтов (OAuth/OIDC)
CREATE TABLE external_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_name VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    provider_user_details JSONB,
    access_token_encrypted TEXT,
    refresh_token_encrypted TEXT,
    token_expires_at TIMESTAMPTZ,
    linked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (provider_name, provider_user_id),
    UNIQUE (user_id, provider_name) -- Пользователь может иметь только один аккаунт каждого провайдера
);
CREATE INDEX idx_external_accounts_user_id ON external_accounts(user_id);
CREATE INDEX idx_external_accounts_provider_user_id ON external_accounts(provider_name, provider_user_id);
COMMENT ON TABLE external_accounts IS 'Связанные внешние аккаунты для OAuth/OIDC аутентификации.';
```

#### 4.2.2. Redis
(Описание структуры данных в Redis остается как в существующем документе).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Apache Kafka.
*   **Формат событий:** CloudEvents v1.0 JSON (согласно `../../../../project_api_standards.md`).
*   **Основной топик для публикуемых событий:** `com.platform.auth.events.v1`.

*   **`com.platform.auth.user.registered.v1`**
    *   Описание: Пользователь успешно зарегистрировался (но может требовать верификации email).
    *   `data` Payload:
        ```json
        {
          "userId": "user-uuid-123",
          "username": "new_user",
          "email": "user@example.com",
          "status": "pending_verification",
          "registrationTimestamp": "2024-03-15T10:00:00Z"
        }
        ```
*   **`com.platform.auth.user.login_success.v1`**
    *   Описание: Пользователь успешно аутентифицировался.
    *   `data` Payload:
        ```json
        {
          "userId": "user-uuid-123",
          "sessionId": "session-uuid-abc",
          "ipAddress": "192.168.1.100",
          "userAgent": "Mozilla/5.0...",
          "loginTimestamp": "2024-03-15T10:05:00Z",
          "mfaMethodUsed": "none"
        }
        ```
*   **`com.platform.auth.user.email_verified.v1`**
    *   Описание: Email пользователя был успешно верифицирован.
    *   `data` Payload:
        ```json
        {
          "userId": "user-uuid-123",
          "email": "user@example.com",
          "verificationTimestamp": "2024-03-16T11:00:00Z"
        }
        ```
*   **`com.platform.auth.user.password_changed.v1`**
    *   Описание: Пароль пользователя был успешно изменен.
    *   `data` Payload:
        ```json
        {
          "userId": "user-uuid-123",
          "changeTimestamp": "2024-03-16T12:00:00Z",
          "method": "user_initiated"
        }
        ```
*   **`com.platform.auth.user.mfa_status.updated.v1`**
    *   Описание: Статус MFA для пользователя был изменен (включен/отключен).
    *   `data` Payload:
        ```json
        {
          "userId": "user-uuid-123",
          "mfaEnabled": true,
          "mfaType": "totp",
          "updateTimestamp": "2024-03-16T13:00:00Z"
        }
        ```
*   **`com.platform.auth.session.revoked.v1`**
    *   Описание: Сессия пользователя была отозвана (включая отзыв refresh token).
    *   `data` Payload:
        ```json
        {
          "userId": "user-uuid-123",
          "sessionId": "session-uuid-abc", // ID сессии, если применимо
          "refreshTokenId": "refresh-token-jti-xyz", // JTI отозванного refresh token
          "revocationReason": "user_logout", // "user_logout", "admin_action", "token_compromised", "idle_timeout"
          "revokedAt": "2024-03-15T11:00:00Z"
        }
        ```
*   **`com.platform.auth.user.password_reset_requested.v1`**
    * Описание: Пользователь запросил сброс пароля.
    * `data` Payload: `{"userId": "user-uuid-123", "email": "user@example.com", "resetTokenId": "opaque-token-id", "requestTimestamp": "ISO8601"}`
*   **`com.platform.auth.user.login_failed.v1`**
    * Описание: Неудачная попытка входа.
    * `data` Payload: `{"loginAttempt": "user@example.com", "reason": "invalid_password", "ipAddress": "1.2.3.4", "userAgent": "...", "failTimestamp": "ISO8601"}`
*   **`com.platform.auth.user.account_locked.v1`**
    * Описание: Аккаунт пользователя заблокирован из-за множества неудачных попыток входа.
    * `data` Payload: `{"userId": "user-uuid-123", "reason": "too_many_failed_login_attempts", "lockTimestamp": "ISO8601"}`
*   **`com.platform.auth.user.roles_changed.v1`**
    * Описание: Роли пользователя были изменены (например, администратором).
    * `data` Payload: `{"userId": "user-uuid-123", "oldRoles": ["user"], "newRoles": ["user", "developer"], "adminId": "admin-uuid-xyz", "changeTimestamp": "ISO8601"}`
*   **`com.platform.auth.api_key.created.v1`**
    * Описание: Создан новый API ключ для пользователя.
    * `data` Payload: `{"userId": "user-uuid-123", "apiKeyId": "apikey-uuid-789", "keyName": "My App Key", "prefix": "pak_", "permissions": ["read_catalog"], "creationTimestamp": "ISO8601"}`
*   **`com.platform.auth.api_key.revoked.v1`**
    * Описание: API ключ был отозван.
    * `data` Payload: `{"userId": "user-uuid-123", "apiKeyId": "apikey-uuid-789", "revocationTimestamp": "ISO8601"}`

### 5.2. Потребляемые События (Consumed Events)
*   **`com.platform.account.user.profile_updated.v1`** (от Account Service)
    *   Логика обработки: Обновить кэшированную информацию о пользователе. Если изменен email и он является основным логином - может потребовать дополнительных действий. Если статус изменен на `deleted` или `blocked` не через Auth Service, отозвать все сессии и токены.
*   **`com.platform.admin.user.force_logout.v1`** (из Admin Service)
    *   Ожидаемый `data` Payload:
        ```json
        {
          "userId": "user-uuid-123",
          "reason": "Suspicious activity reported by admin.",
          "adminId": "admin-uuid-xyz",
          "actionTimestamp": "2024-03-16T14:00:00Z"
        }
        ```
    *   Логика обработки: Отозвать все активные сессии и refresh-токены для указанного `userId`. Записать в аудит.
*   **`com.platform.admin.user.status.updated.v1`** (из Admin Service, если статус меняется там, например, блокировка)
     *   Ожидаемый `data` Payload:
        ```json
        {
          "userId": "user-uuid-123",
          "newStatus": "blocked",
          "reason": "Violation of terms by admin.",
          "adminId": "admin-uuid-xyz",
          "actionTimestamp": "2024-03-16T14:05:00Z"
        }
        ```
    *   Логика обработки: Обновить статус пользователя в Auth Service. Если статус `blocked` или `deleted`, отозвать все сессии и токены. Записать в аудит.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
(Как в существующем документе, например: Account Service, Notification Service, API Gateway, Admin Service)

### 6.2. Внешние Системы
*   **OAuth/OIDC Провайдеры:**
    *   Для аутентификации через социальные сети и внешние провайдеры, Auth Service обеспечивает поддержку **ВКонтакте (VK ID)** и **Telegram Login** как ключевых российских платформ. Интеграция с другими провайдерами OAuth 2.0/OIDC возможна и будет рассматриваться по мере необходимости.
    *   Взаимодействие происходит через стандартные протоколы OAuth 2.0 (Authorization Code Grant) и OpenID Connect.
*   **Сервисы отправки SMS/Email (для 2FA и верификации):** Интеграция с Notification Service, который, в свою очередь, взаимодействует с конкретными провайдерами.

(Описание интеграций остается как в существующем документе).

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
(Список переменных окружения в целом остается как в существующем документе. Ключевые моменты):
*   `AUTH_JWT_PRIVATE_KEY_PATH`: Путь к файлу с приватным RSA ключом (например, `/run/secrets/auth_jwt_private.pem`). **Этот файл должен монтироваться из Kubernetes Secret или аналогичного безопасного хранилища.**
*   `AUTH_JWT_PUBLIC_KEY_PATH`: Путь к файлу с публичным RSA ключом (например, `/app/config/auth_jwt_public.pem`). Может быть частью ConfigMap или образа.
*   OAuth provider client IDs and secrets (например `AUTH_VK_CLIENT_ID`, `AUTH_VK_CLIENT_SECRET`) **должны загружаться из Kubernetes Secrets, а не напрямую через переменные окружения, если это возможно, либо переменные окружения должны быть внедрены из секретов.**

### 7.2. Файлы Конфигурации (`config/config.yaml`)
(Структура YAML остается как в существующем документе, с акцентом на то, что пути к секретным ключам указывают на файлы, управляемые системой секретов Kubernetes/Vault).

## 8. Обработка Ошибок (Error Handling)
(Описание остается как в существующем документе, форматы ошибок исправлены в разделе API).

## 9. Безопасность (Security)

### 9.1. Аутентификация
(Описание остается как в существующем документе).

### 9.2. Авторизация
(Описание остается как в существующем документе, ссылка на `../../../../project_roles_and_permissions.md` корректна).

### 9.3. Защита Данных
*   **Хеширование паролей:** Используется стойкий алгоритм Argon2id с настраиваемыми параметрами (см. `../../../../project_security_standards.md` и конфигурацию сервиса).
*   **Шифрование:** Секреты для TOTP и другие чувствительные данные конфигурации (например, `secret_encrypted` в `mfa_secrets`) шифруются перед сохранением в БД с использованием симметричного шифрования (например, AES-GCM). Ключ шифрования этих данных управляется через систему управления секретами.
*   **Транспортная безопасность:** TLS 1.3 (рекомендуется, минимум 1.2) для всех внешних и внутренних коммуникаций.
*   **Защита от атак:** (Как в существующем документе).
*   **Соответствие ФЗ-152 "О персональных данных":** Auth Service обрабатывает ПДн (email, username, IP, User-Agent, потенциально телефон для SMS 2FA, данные из внешних OAuth-провайдеров). **Персональные данные российских граждан, собираемые и обрабатываемые Auth Service, хранятся и обрабатываются на серверах, физически расположенных на территории Российской Федерации, с использованием инфраструктуры российского хостинг-провайдера Beget.** Получение согласия на обработку ПДн является частью процесса регистрации. Все операции с ПДн логируются.
*   (Ссылка на `../../../../project_security_standards.md`).

### 9.4. Управление Секретами
(Как в существующем документе, ссылка на `../../../../project_security_standards.md`).

## 10. Развертывание (Deployment)
(Описание остается как в существующем документе, ссылки на `../../../../project_deployment_standards.md` корректны).

## 11. Мониторинг и Логирование (Logging and Monitoring)
(Описание остается как в существующем документе, ссылки на `../../../../project_observability_standards.md` корректны).

## 12. Нефункциональные Требования (NFRs)
(Описание NFRs остается как в существующем документе).

## 13. Приложения (Appendices)
*   **JWT Claims (Пример структуры Access Token Payload):**
    ```json
    {
      "iss": "https://auth.mygameplatform.ru", // Issuer
      "sub": "user-uuid-123",                 // Subject (User ID)
      "aud": ["https://api.mygameplatform.ru"], // Audience (Ресурсы, для которых предназначен токен)
      "exp": 1678886400,                      // Expiration Time (Unix timestamp)
      "nbf": 1678882800,                      // Not Before (Unix timestamp)
      "iat": 1678882800,                      // Issued At (Unix timestamp)
      "jti": "jwt-unique-id-abc",             // JWT ID (для отзыва)
      "username": "new_user",                 // Имя пользователя
      "email": "user@example.com",            // Email
      "roles": ["user", "beta_tester"],       // Роли пользователя
      "permissions": ["read:game_catalog", "write:game_review"], // Опционально, детальные разрешения
      "amr": ["pwd", "mfa_totp"],             // Authentication Methods References (как пользователь был аутентифицирован)
      "sid": "session-uuid-xyz"               // Опционально, ID сессии, если токены привязаны к сессиям
    }
    ```
*   Детальные схемы Protobuf для gRPC API находятся в репозитории `platform-protos` (или локально `proto/auth/v1/`).
*   Примеры полных JSON для всех REST DTO могут быть добавлены при необходимости или генерироваться из OpenAPI спецификации.

## 14. Пользовательские Сценарии (User Flows)

В этом разделе описаны ключевые пользовательские сценарии, связанные с Auth Service.

### 14.1. Регистрация Нового Пользователя (Email/Пароль, Верификация Email)
*   **Описание:** Пользователь регистрируется в системе, предоставляя email и пароль. На указанный email отправляется письмо с кодом/ссылкой для верификации.
*   **Связанный основной воркфлоу:** [Регистрация пользователя и начальная настройка профиля](../../../../project_workflows/user_registration_flow.md)
*   **Диаграмма (фокус на Auth Service):**
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Клиентское Приложение
        participant APIGW as API Gateway
        participant AuthSvc as Auth Service
        participant UserDB as База Данных (PostgreSQL)
        participant NotificationSvc as Notification Service (через Kafka)
        participant Kafka as Kafka

        User->>ClientApp: Заполняет форму регистрации (username, email, password)
        ClientApp->>APIGW: POST /api/v1/auth/register (payload)
        APIGW->>AuthSvc: Forward /register
        AuthSvc->>AuthSvc: Валидация данных (формат, уникальность username/email)
        alt Данные валидны
            AuthSvc->>UserDB: Создание пользователя (статус: pending_verification), хеширование пароля (Argon2id)
            AuthSvc->>UserDB: Генерация и сохранение кода верификации email (VerificationCode)
            AuthSvc-->>Kafka: Publish `com.platform.auth.user.registered.v1` (userId, email, username)
            AuthSvc-->>Kafka: Publish `com.platform.auth.email.verification_requested.v1` (userId, email, verification_code)
            Kafka-->>NotificationSvc: Consume `email.verification_requested`
            NotificationSvc->>User: Отправка email с кодом/ссылкой верификации
            AuthSvc-->>APIGW: HTTP 201 Created (userId, status)
            APIGW-->>ClientApp: HTTP 201 Created
            ClientApp-->>User: Сообщение об успешной регистрации и необходимости верификации
        else Ошибка валидации
            AuthSvc-->>APIGW: HTTP 400/409 (Ошибка валидации/Конфликт)
            APIGW-->>ClientApp: HTTP 400/409
            ClientApp-->>User: Отображение ошибки
        end

        User->>ClientApp: Переход по ссылке из email или ввод кода
        ClientApp->>APIGW: POST /api/v1/auth/verify-email (verification_code)
        APIGW->>AuthSvc: Forward /verify-email
        AuthSvc->>UserDB: Проверка кода верификации
        alt Код верен
            AuthSvc->>UserDB: Обновление статуса пользователя на 'active', email_verified_at
            AuthSvc->>UserDB: Удаление/пометка использованным кода верификации
            AuthSvc-->>Kafka: Publish `com.platform.auth.user.email_verified.v1` (userId, email)
            AuthSvc-->>APIGW: HTTP 200 OK
            APIGW-->>ClientApp: HTTP 200 OK
            ClientApp-->>User: Email успешно подтвержден
        else Код не верен или истек
            AuthSvc-->>APIGW: HTTP 400 Bad Request (INVALID_VERIFICATION_CODE)
            APIGW-->>ClientApp: HTTP 400
            ClientApp-->>User: Ошибка подтверждения email
        end
    ```

### 14.2. Вход Пользователя (Email/Пароль) и Выдача JWT
*   **Описание:** Пользователь входит в систему, используя свой email и пароль. Auth Service проверяет учетные данные и выдает пару JWT (Access Token и Refresh Token).
*   **Диаграмма:** (См. диаграмму "Login with Password & JWT Issuance" в разделе 2 или 3 данной документации, если она там размещена, или продублировать/адаптировать сюда).
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Клиентское Приложение
        participant APIGW as API Gateway
        participant AuthSvc as Auth Service
        participant UserDB as База Данных Пользователей (PostgreSQL)
        participant SessionCache as Кэш Сессий (Redis)

        User->>ClientApp: Ввод логина и пароля
        ClientApp->>APIGW: POST /api/v1/auth/login (login, password)
        APIGW->>AuthSvc: Forward login request
        AuthSvc->>UserDB: Поиск пользователя по логину
        UserDB-->>AuthSvc: Данные пользователя (включая хеш пароля, статус, 2FA статус)
        alt Пользователь найден и активен
            AuthSvc->>AuthSvc: Проверка пароля (Argon2id.Verify(password, hash))
            alt Пароль верен
                alt 2FA не включен
                    AuthSvc->>AuthSvc: Генерация Access Token (JWT RS256)
                    AuthSvc->>AuthSvc: Генерация Refresh Token
                    AuthSvc->>SessionCache: Сохранение сессии (опционально)
                    AuthSvc->>UserDB: Сохранение Refresh Token (хешированный, связан с сессией)
                    AuthSvc-->>Kafka: Publish `com.platform.auth.user.login_success.v1`
                    AuthSvc-->>APIGW: HTTP 200 OK (Access Token, Refresh Token в HttpOnly cookie)
                    APIGW-->>ClientApp: HTTP 200 OK (Access Token)
                    ClientApp-->>User: Успешный вход
                else 2FA включен
                    AuthSvc->>AuthSvc: Генерация временного токена/сессии для шага 2FA
                    AuthSvc-->>APIGW: HTTP 202 Accepted (требуется 2FA, временный токен)
                    APIGW-->>ClientApp: HTTP 202 Accepted
                    ClientApp-->>User: Запрос 2FA кода
                end
            else Пароль не верен
                AuthSvc-->>Kafka: Publish `com.platform.auth.user.login_failed.v1`
                AuthSvc-->>APIGW: HTTP 401 Unauthorized (INVALID_CREDENTIALS)
                APIGW-->>ClientApp: HTTP 401
                ClientApp-->>User: Ошибка: Неверный логин или пароль
            end
        else Пользователь не найден, не активен или заблокирован
            AuthSvc-->>Kafka: Publish `com.platform.auth.user.login_failed.v1`
            AuthSvc-->>APIGW: HTTP 401 Unauthorized (INVALID_CREDENTIALS или USER_INACTIVE/BLOCKED)
            APIGW-->>ClientApp: HTTP 401
            ClientApp-->>User: Ошибка входа
        end
    ```

### 14.3. Вход Пользователя с 2FA (TOTP)
*   **Описание:** После успешного ввода пароля, если у пользователя включена 2FA (TOTP), система запрашивает TOTP код.
*   **Диаграмма:** (См. диаграмму "Login with 2FA (TOTP)" в разделе 2 или 3, или продублировать/адаптировать сюда).
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Клиентское Приложение
        participant APIGW as API Gateway
        participant AuthSvc as Auth Service
        participant UserDB as База Данных Пользователей (PostgreSQL)

        User->>ClientApp: Ввод TOTP кода (после шага с паролем)
        ClientApp->>APIGW: POST /api/v1/auth/login/2fa-verify (временный_токен_сессии, totp_code)
        APIGW->>AuthSvc: Forward 2FA verification request
        AuthSvc->>UserDB: Получение секрета TOTP для пользователя (связанного с временным токеном)
        alt Секрет найден
            AuthSvc->>AuthSvc: Валидация TOTP кода
            alt TOTP код верен
                AuthSvc->>AuthSvc: Завершение процесса логина: генерация Access/Refresh токенов
                AuthSvc-->>Kafka: Publish `com.platform.auth.user.login_success.v1` (mfaMethodUsed="totp")
                AuthSvc-->>APIGW: HTTP 200 OK (Access Token, Refresh Token в HttpOnly cookie)
                APIGW-->>ClientApp: HTTP 200 OK (Access Token)
                ClientApp-->>User: Успешный вход
            else TOTP код не верен
                AuthSvc-->>Kafka: Publish `com.platform.auth.user.login_failed.v1` (reason="invalid_2fa_code")
                AuthSvc-->>APIGW: HTTP 401 Unauthorized (INVALID_2FA_CODE)
                APIGW-->>ClientApp: HTTP 401
                ClientApp-->>User: Ошибка: Неверный 2FA код
            end
        else Ошибка (например, временная сессия не найдена)
            AuthSvc-->>APIGW: HTTP 400 Bad Request
            APIGW-->>ClientApp: HTTP 400
            ClientApp-->>User: Ошибка
        end
    ```

### 14.4. Обновление Access Token с Использованием Refresh Token
*   **Описание:** Access token пользователя истек. Клиентское приложение использует Refresh Token для получения нового Access Token без повторного ввода пароля.
*   **Диаграмма:** (См. диаграмму "Access Token Refresh" в разделе 2 или 3, или продублировать/адаптировать сюда).
    ```mermaid
    sequenceDiagram
        actor ClientApp as Клиентское Приложение
        participant APIGW as API Gateway
        participant AuthSvc as Auth Service
        participant TokenDB as Хранилище Refresh Токенов (PostgreSQL/Redis JTI Blacklist)

        ClientApp->>APIGW: POST /api/v1/auth/refresh-token (с Refresh Token из HttpOnly cookie)
        APIGW->>AuthSvc: Forward refresh token request
        AuthSvc->>TokenDB: Поиск и валидация Refresh Token (проверка хеша, срока действия, не отозван ли JTI)
        alt Refresh Token валиден и не отозван
            AuthSvc->>AuthSvc: Генерация нового Access Token (JWT RS256)
            AuthSvc->>AuthSvc: (Опционально, если настроена ротация) Генерация нового Refresh Token, отзыв старого JTI и сохранение нового JTI.
            AuthSvc->>TokenDB: (Опционально) Обновление/сохранение нового Refresh Token, добавление старого JTI в blacklist.
            AuthSvc-->>APIGW: HTTP 200 OK (новый Access Token; новый Refresh Token в HttpOnly cookie если ротация)
            APIGW-->>ClientApp: HTTP 200 OK (новый Access Token)
        else Refresh Token невалиден или отозван
            AuthSvc-->>APIGW: HTTP 401 Unauthorized (INVALID_REFRESH_TOKEN)
            APIGW-->>ClientApp: HTTP 401 Unauthorized (требуется повторный логин)
        end
    ```

### 14.5. Сброс Пароля
*   **Описание:** Пользователь забыл пароль и инициирует процедуру сброса через email.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Клиентское Приложение
        participant APIGW as API Gateway
        participant AuthSvc as Auth Service
        participant UserDB as База Данных (PostgreSQL)
        participant NotificationSvc as Notification Service (через Kafka)
        participant Kafka as Kafka

        User->>ClientApp: Нажимает "Забыли пароль?"
        ClientApp->>APIGW: POST /api/v1/auth/forgot-password (email: "user@example.com")
        APIGW->>AuthSvc: Forward /forgot-password
        AuthSvc->>UserDB: Поиск пользователя по email
        alt Пользователь найден и email верифицирован
            AuthSvc->>UserDB: Генерация и сохранение токена/кода сброса пароля (VerificationCode type='password_reset')
            AuthSvc-->>Kafka: Publish `com.platform.auth.user.password_reset_requested.v1` (userId, email, reset_token_or_code)
            Kafka-->>NotificationSvc: Consume `password_reset_requested`
            NotificationSvc->>User: Отправка email с инструкцией и ссылкой/кодом для сброса
        end
        AuthSvc-->>APIGW: HTTP 200 OK (общее сообщение, не раскрывающее существование email)
        APIGW-->>ClientApp: HTTP 200 OK
        ClientApp-->>User: Сообщение об отправке инструкции (если email существует)

        User->>ClientApp: Переход по ссылке из email / Ввод кода и нового пароля
        ClientApp->>APIGW: POST /api/v1/auth/reset-password (reset_token_or_code, new_password, confirm_password)
        APIGW->>AuthSvc: Forward /reset-password
        AuthSvc->>UserDB: Проверка токена/кода сброса
        alt Токен/код верен и не истек
            AuthSvc->>UserDB: Обновление хеша пароля пользователя (Argon2id)
            AuthSvc->>UserDB: Удаление/пометка использованным токена/кода сброса
            AuthSvc-->>Kafka: Publish `com.platform.auth.user.password_changed.v1` (userId, method="password_reset")
            AuthSvc->>UserDB: (Опционально) Отзыв всех активных сессий пользователя
            AuthSvc-->>APIGW: HTTP 200 OK (Пароль успешно изменен)
            APIGW-->>ClientApp: HTTP 200 OK
            ClientApp-->>User: Пароль изменен, предложение войти
        else Токен/код не верен или истек
            AuthSvc-->>APIGW: HTTP 400 Bad Request (INVALID_RESET_TOKEN)
            APIGW-->>ClientApp: HTTP 400
            ClientApp-->>User: Ошибка сброса пароля
        end
    ```

### 14.6. Валидация Access Token Другим Микросервисом (gRPC)
*   **Описание:** Внутренний микросервис (например, Catalog Service) получает запрос от API Gateway с Access Token и обращается к Auth Service для его валидации и получения информации о пользователе.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        participant APIGW as API Gateway
        participant CatalogSvc as Catalog Service
        participant AuthSvcGRPC as Auth Service (gRPC API)

        APIGW->>CatalogSvc: Запрос к Catalog Service (с `Authorization: Bearer <access_token>`)
        CatalogSvc->>AuthSvcGRPC: ValidateTokenRequest (access_token)
        AuthSvcGRPC->>AuthSvcGRPC: Валидация подписи, срока действия, JTI (если используется blacklist в Redis)
        alt Токен валиден
            AuthSvcGRPC-->>CatalogSvc: ValidateTokenResponse (user_id, username, roles, permissions, is_valid=true)
            CatalogSvc->>CatalogSvc: Обработка запроса с учетом прав пользователя
            CatalogSvc-->>APIGW: Ответ Catalog Service
        else Токен невалиден
            AuthSvcGRPC-->>CatalogSvc: ValidateTokenResponse (is_valid=false, error_message)
            CatalogSvc-->>APIGW: HTTP 401/403 Unauthorized/Forbidden
        end
    ```

### 14.7. Вход Пользователя через Внешнего OAuth2/OIDC Провайдера (например, VK)
*   **Описание:** Пользователь выбирает вход через VK. Происходит редирект на VK, пользователь аутентифицируется там, затем редирект обратно на платформу с кодом авторизации. Auth Service обменивает код на токен VK, получает данные пользователя VK и создает/связывает аккаунт на платформе, выпуская JWT для сессии на платформе.
*   **Диаграмма:** (См. диаграмму "OAuth 2.0 Authorization Code Grant Flow" в разделе 2 или 3, или продублировать/адаптировать сюда).
     ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Клиентское Приложение
        participant APIGW as API Gateway
        participant AuthSvc as Auth Service
        participant VKAuth as VK Authorization Server
        participant UserDB as База Данных Пользователей (PostgreSQL)
        participant Kafka as Kafka

        User->>ClientApp: Нажимает "Войти через VK"
        ClientApp->>APIGW: GET /api/v1/auth/oauth/vk/login-url
        APIGW->>AuthSvc: Forward /oauth/vk/login-url
        AuthSvc-->>APIGW: HTTP 200 OK (redirect_url_to_vk)
        APIGW-->>ClientApp: HTTP 200 OK (redirect_url_to_vk)
        ClientApp->>User: Редирект на VK Authorization Server (redirect_url_to_vk)

        User->>VKAuth: Аутентификация на стороне VK, предоставление разрешений
        VKAuth-->>ClientApp: Редирект на callback URL платформы (с authorization_code)

        ClientApp->>APIGW: GET /api/v1/auth/oauth/vk/callback?code=<authorization_code>&state=<state_param_if_used>
        APIGW->>AuthSvc: Forward /oauth/vk/callback
        AuthSvc->>VKAuth: Обмен authorization_code на access_token VK (с использованием client_id, client_secret)
        VKAuth-->>AuthSvc: VK access_token, VK refresh_token (если есть), user_id_vk
        AuthSvc->>VKAuth: Запрос информации о пользователе VK (используя VK access_token)
        VKAuth-->>AuthSvc: Информация о пользователе VK (email, имя и т.д.)

        AuthSvc->>UserDB: Поиск пользователя по provider_name='vk' и provider_user_id
        alt Пользователь VK уже привязан
            AuthSvc->>UserDB: Обновление данных пользователя из VK (если необходимо)
            AuthSvc->>AuthSvc: Генерация JWT (Access/Refresh) для существующего пользователя
            AuthSvc-->>Kafka: Publish `com.platform.auth.user.login_success.v1` (method="oauth_vk")
        else Новый пользователь VK (или существующий пользователь по email, но без привязки VK)
            AuthSvc->>UserDB: Поиск пользователя по email из VK (если есть)
            alt Пользователь с таким email существует
                AuthSvc->>UserDB: Привязка VK аккаунта к существующему пользователю (создание ExternalAccount)
            else Новый пользователь платформы
                AuthSvc->>UserDB: Создание нового пользователя (статус 'active', email_verified_at из VK если есть)
                AuthSvc->>UserDB: Создание ExternalAccount
                AuthSvc-->>Kafka: Publish `com.platform.auth.user.registered.v1` (source="oauth_vk")
            end
            AuthSvc->>AuthSvc: Генерация JWT (Access/Refresh)
            AuthSvc-->>Kafka: Publish `com.platform.auth.user.login_success.v1` (method="oauth_vk")
        end
        AuthSvc-->>APIGW: HTTP 200 OK (Access Token, Refresh Token в HttpOnly cookie)
        APIGW-->>ClientApp: HTTP 200 OK (Access Token)
        ClientApp-->>User: Успешный вход / Регистрация через VK
    ```

## 15. Резервное Копирование и Восстановление (Backup and Recovery)

### 14.1. PostgreSQL (Данные пользователей, ролей, токенов и т.д.)
*   **Процедура резервного копирования:**
    *   **Логические бэкапы:** Ежедневный `pg_dumpall` или `pg_dump` для базы данных Auth Service.
    *   **Физические бэкапы (Point-in-Time Recovery - PITR):** Настроена непрерывная архивация WAL-сегментов. Базовый бэкап создается еженедельно.
    *   **Хранение:** Бэкапы и WAL-архивы хранятся в S3-совместимом хранилище с шифрованием на стороне сервера и версионированием, предпочтительно в другом географическом регионе от основной БД. Срок хранения: полные логические бэкапы - 60 дней, WAL-сегменты - 30 дней (для возможности восстановления на более длительный период).
*   **Процедура восстановления:**
    *   Тестируется ежеквартально. Восстановление из логического бэкапа или с использованием PITR.
*   **RTO (Recovery Time Objective):** < 1 часа.
*   **RPO (Recovery Point Objective):** < 5 минут.
*   (Общие принципы см. `../../../../project_database_structure.md`).

### 14.2. Redis (Кэш сессий, JTI blacklist, временные коды)
*   **Стратегия персистентности и бэкапа:**
    *   **AOF (Append Only File):** Включен с fsync `everysec` для минимизации потерь данных для критичных данных, таких как JTI blacklist.
    *   **RDB Snapshots:** Регулярное создание снапшотов (например, каждые 1-6 часов в зависимости от нагрузки и критичности данных в Redis).
    *   **Хранение:** RDB-снапшоты (и AOF при необходимости) могут копироваться в S3-совместимое хранилище ежедневно. Срок хранения - 7-14 дней.
    *   Большинство данных в Redis (сессии, временные коды) имеют TTL и могут быть некритичны для восстановления из бэкапа, так как пользователи могут просто перелогиниться или запросить код заново. JTI blacklist более важен для предотвращения использования отозванных токенов.
*   **Процедура восстановления:** Восстановление из последнего RDB-снапшота и/или AOF.
*   **RTO:** < 30 минут.
*   **RPO:** < 1 минуты (для данных с AOF `everysec`). Для менее критичных данных в Redis, RPO может быть равен интервалу RDB снапшотирования.

### 14.3. Общая стратегия
*   Резервное копирование и восстановление Auth Service являются критически важными для обеспечения непрерывности работы платформы.
*   Процедуры должны быть тщательно документированы, автоматизированы и регулярно тестироваться.
*   Мониторинг процессов резервного копирования обязателен.
*   Ключи шифрования, используемые для защиты данных в Auth Service (например, для `mfa_secrets`), должны быть частью стратегии резервного копирования и восстановления системы управления секретами (Vault/Kubernetes Secrets).

## 16. Приложения (Appendices)
*   **JWT Claims (Пример структуры Access Token Payload):**
    ```json
    {
      "iss": "https://auth.mygameplatform.ru", // Issuer
      "sub": "user-uuid-123",                 // Subject (User ID)
      "aud": ["https://api.mygameplatform.ru"], // Audience (Ресурсы, для которых предназначен токен)
      "exp": 1678886400,                      // Expiration Time (Unix timestamp)
      "nbf": 1678882800,                      // Not Before (Unix timestamp)
      "iat": 1678882800,                      // Issued At (Unix timestamp)
      "jti": "jwt-unique-id-abc",             // JWT ID (для отзыва)
      "username": "new_user",                 // Имя пользователя
      "email": "user@example.com",            // Email
      "roles": ["user", "beta_tester"],       // Роли пользователя
      "permissions": ["read:game_catalog", "write:game_review"], // Опционально, детальные разрешения
      "amr": ["pwd", "mfa_totp"],             // Authentication Methods References (как пользователь был аутентифицирован)
      "sid": "session-uuid-xyz"               // Опционально, ID сессии, если токены привязаны к сессиям
    }
    ```
*   Детальные схемы Protobuf для gRPC API находятся в репозитории `platform-protos` (или локально `proto/auth/v1/`).
*   Примеры полных JSON для всех REST DTO могут быть добавлены при необходимости или генерироваться из OpenAPI спецификации.


## 17. Связанные Рабочие Процессы (Related Workflows)
*   [Регистрация пользователя и начальная настройка профиля](../../../../project_workflows/user_registration_flow.md)
*   [Аутентификация пользователя (логин, 2FA, OAuth)] Подробное описание этого рабочего процесса будет добавлено в [user_authentication_flow.md](../../../../project_workflows/user_authentication_flow.md) (документ в разработке).
*   [Сброс и восстановление пароля] Подробное описание этого рабочего процесса будет добавлено в [password_recovery_flow.md](../../../../project_workflows/password_recovery_flow.md) (документ в разработке).

---
*Этот документ является основной спецификацией для Auth Service и должен поддерживаться в актуальном состоянии.*
