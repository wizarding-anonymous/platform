# Спецификация Микросервиса: Account Service

**Версия:** 1.0
**Дата последнего обновления:** {{YYYY-MM-DD}} <!-- Placeholder for date -->

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   Account Service является основным микросервисом платформы "Российский Аналог Steam", отвечающим за управление учетными записями пользователей, их профилями и связанными данными, такими как настройки, информация для верификации и контактная информация.
*   Он служит авторитетным источником данных, связанных с пользователями (кроме учетных данных аутентификации, управляемых Auth Service), предоставляя эту информацию другим сервисам и обеспечивая целостность и безопасность данных.
*   Основные бизнес-задачи: управление жизненным циклом аккаунта пользователя, управление профильной информацией, управление контактными данными и их верификацией, управление пользовательскими настройками.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Управление учетными записями:** Создание базовой записи аккаунта после регистрации пользователя (координируется с Auth Service), хранение и управление основной информацией об аккаунте (ID, статус), поддержка изменения статуса аккаунта (активация, блокировка, удаление).
*   **Управление профилями пользователей:** Создание и редактирование деталей профиля (никнейм, биография, страна, аватар, кастомный URL), управление настройками видимости профиля.
*   **Управление контактной информацией и верификацией:** Хранение и управление контактными данными пользователя (email, телефон), обработка процесса верификации.
*   **Управление настройками:** Хранение и обновление пользовательских настроек, разделенных по категориям (например, приватность, уведомления, интерфейс).
*   **Генерация событий:** Публикация событий, связанных с изменениями в аккаунтах и профилях (например, создание аккаунта, обновление профиля) в брокер сообщений Kafka.

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+, согласно `../../../../project_technology_stack.md`)
*   **REST Framework:** Echo (`github.com/labstack/echo/v4`)
*   **gRPC Framework:** `google.golang.org/grpc`
*   **База данных:** PostgreSQL (версия 15+)
*   **Кэширование:** Redis (версия 7.0+)
*   **Брокер сообщений:** Kafka (`github.com/confluentinc/confluent-kafka-go` или `github.com/segmentio/kafka-go`)
*   **Валидация:** `go-playground/validator/v10`
*   **ORM/DB Driver:** GORM (`gorm.io/gorm`) с драйвером `gorm.io/driver/postgres`
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`)
*   **Логирование:** Zap (`go.uber.org/zap`)
*   **Трассировка и метрики:** OpenTelemetry (`go.opentelemetry.io/otel`), Prometheus client (`github.com/prometheus/client_golang`)
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `../../../../PACKAGE_STANDARDIZATION.md`.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`, `../../../../project_observability_standards.md`.

### 1.4. Термины и Определения (Glossary)
*   См. `../../../../project_glossary.md`.
*   **Аккаунт (Account):** Основная запись пользователя в системе, связанная с его UserID из Auth Service. Содержит статус и базовую информацию.
*   **Профиль (Profile):** Публичная и приватная информация пользователя, отображаемая другим пользователям (никнейм, аватар, биография и т.д.).
*   **Контактная информация (ContactInfo):** Email, телефон пользователя, используемые для связи и верификации.
*   **Настройки (Settings):** Пользовательские конфигурации для управления поведением платформы и видимостью данных.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Сервис построен с использованием принципов чистой архитектуры (Clean Architecture) для разделения ответственностей на слои: Presentation, Application, Domain, Infrastructure.
*   Основные модули: управление аккаунтами, управление профилями, управление контактами, управление настройками.
*   Диаграмма верхнеуровневой архитектуры сервиса:
```mermaid
graph TD
    A[Presentation Layer (HTTP/gRPC Handlers)] --> B(Application Layer (Use Cases/Services))
    B --> C{Domain Layer (Entities, Repos Interfaces)}
    B --> D{Infrastructure Layer (DB Repos, Kafka Client)}
    D -- Implements --> C
```

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления)
*   Ответственность: Обработка входящих HTTP (REST) и gRPC запросов, валидация DTO, вызов Application Layer.
*   Ключевые компоненты/модули:
    *   HTTP Handlers (Echo): Эндпоинты для CRUD операций над аккаунтами, профилями, контактами, настройками.
    *   gRPC Service Implementations: Методы для межсервисного взаимодействия (например, `GetAccountInfo`, `GetUserProfile`).
    *   DTOs: Структуры для передачи данных (например, `CreateAccountRequest`, `UserProfileResponse`), валидируемые с помощью `go-playground/validator`.
    *   Валидаторы: Кастомные правила валидации, если стандартных недостаточно.

#### 2.2.2. Application Layer (Прикладной Слой / Слой Сценариев Использования)
*   Ответственность: Координация бизнес-логики, реализация use cases. Не содержит бизнес-правил напрямую.
*   Ключевые компоненты/модули:
    *   Use Case Services: `AccountUseCaseService`, `ProfileUseCaseService`, `ContactUseCaseService`, `SettingsUseCaseService`.
    *   Интерфейсы для репозиториев (`AccountRepository`, `ProfileRepository`) и внешних сервисов (например, `NotificationServiceIntegration` - интерфейс для клиента Notification Service).

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности, агрегаты, доменные события, бизнес-правила. Независим от других слоев.
*   Ключевые компоненты/модули:
    *   Entities: `Account` (ID, UserID (из Auth Service), Status, Timestamps), `Profile` (Nickname, Bio, AvatarURL, CountryCode), `ContactInfo` (Email, Phone, VerifiedStatus), `UserSetting` (SettingsData JSONB).
    *   Value Objects: [NEEDS DEVELOPER INPUT: Specific Value Objects for account-service, if any, e.g., EmailAddress, PhoneNumber]
    *   Domain Services: [NEEDS DEVELOPER INPUT: Specific Domain Services for account-service, if any]
    *   Domain Events: `AccountCreatedEvent`, `ProfileUpdatedEvent`, `ContactVerifiedEvent`.
    *   Интерфейсы репозиториев (определяются здесь, реализуются в Infrastructure Layer).

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Реализация интерфейсов, определенных в Application и Domain Layers, для взаимодействия с PostgreSQL, Redis, Kafka.
*   Ключевые компоненты/модули:
    *   PostgreSQL Repositories (GORM): Реализации `AccountRepository`, `ProfileRepository` и т.д.
    *   Redis Cache: Для кэширования профилей, настроек.
    *   Kafka Producer: Для отправки доменных событий в формате CloudEvents.
    *   External Service Clients: Клиенты для других сервисов (например, gRPC клиент для Auth Service).
    *   Адаптеры для файловых хранилищ: [NEEDS DEVELOPER INPUT: Specific S3/Object Storage adapter details if Account Service handles uploads directly]

## 3. API Endpoints
Общие принципы и форматы см. `../../../../project_api_standards.md`.

### 3.1. REST API
*   **Базовый URL (через API Gateway):** `/api/v1/account` (уточненный базовый URL для ресурсов Account Service).
*   **Версионирование:** `/v1/`.
*   **Формат данных:** `application/json` (согласно `../../../../project_api_standards.md`).
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway, UserID и роли передаются в заголовках, например, `X-User-ID`, `X-User-Roles`).
*   **Стандартные заголовки:** `X-Request-ID`.

#### 3.1.1. Ресурс: Мой Аккаунт и Профиль (My Account & Profile)
*   **`GET /me`**
    *   Описание: Получение информации о текущем аутентифицированном пользователе (его аккаунте, профиле, основных настройках).
    *   Query параметры: `include=contact_info,full_settings` (опционально, для включения дополнительных связанных данных).
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "id": "uuid-account-id",
            "type": "userAccount",
            "attributes": {
              "userId": "uuid-auth-user-id",
              "status": "active",
              "profile": {
                "nickname": "User123",
                "bio": "Hello world!",
                "avatarUrl": "https://example.com/avatar.jpg",
                "countryCode": "RU",
                "customUrlSlug": "user123_profile"
              },
              "createdAt": "2023-10-28T10:00:00Z",
              "updatedAt": "2023-10-28T10:00:00Z"
            },
            "relationships": {
              "contactInfo": { "links": { "related": "/api/v1/account/me/contact-info" } },
              "settings": { "links": { "related": "/api/v1/account/me/settings" } }
            }
          }
        }
        ```
    *   Пример ответа (Ошибка 4xx/5xx): `[NEEDS DEVELOPER INPUT: Example error response for GET /me]`
    *   Требуемые права доступа: Аутентифицированный пользователь (`user`).
*   **`PUT /me/profile`**
    *   Описание: Обновление профиля текущего пользователя.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "profileUpdate",
            "attributes": {
              "nickname": "NewUser123",
              "bio": "Updated bio.",
              "countryCode": "RU",
              "customUrlSlug": "new_user123_profile",
              "privacySettings": { "showRealName": false, "inventoryPublic": true }
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): `[NEEDS DEVELOPER INPUT: Example success response for PUT /me/profile]`
    *   Требуемые права доступа: Владелец профиля (`user`).
*   **`POST /me/profile/avatar`**
    *   Описание: Загрузка нового аватара для текущего пользователя (multipart/form-data).
    *   Тело запроса: `file: (binary_avatar_data)`
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "avatarUpdate",
            "attributes": {
              "avatarUrl": "https://example.com/new_avatar.jpg"
            }
          }
        }
        ```
    *   Требуемые права доступа: Владелец профиля (`user`).
*   **`DELETE /me`**
    *   Описание: Запрос на удаление аккаунта текущего пользователя (инициирует процесс деактивации/удаления согласно политикам платформы).
    *   Пример ответа (Успех 202 Accepted): `[NEEDS DEVELOPER INPUT: Example success response for DELETE /me]`
    *   Требуемые права доступа: Владелец аккаунта (`user`).

#### 3.1.2. Ресурс: Публичные Профили (Public Profiles)
*   **`GET /profiles/{profile_id_or_custom_url}`**
    *   Описание: Получение публичного профиля пользователя по его ID профиля или кастомному URL. Объем возвращаемых данных зависит от настроек приватности профиля.
    *   Пример ответа (Успех 200 OK): `[NEEDS DEVELOPER INPUT: Example success response for GET /profiles/{id}]`
    *   Требуемые права доступа: `anonymous` (для публичных данных), `user` (для данных, видимых другим пользователям).

#### 3.1.3. Ресурс: Контактная Информация (Contact Info) - `/me/contact-info`
*   **`GET /me/contact-info`**
    *   Описание: Получение контактной информации текущего пользователя.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {"type": "contactMethod", "id": "uuid-email-id", "attributes": {"contactType": "email", "value": "user@example.com", "isVerified": true, "isPrimary": true}},
            {"type": "contactMethod", "id": "uuid-phone-id", "attributes": {"contactType": "phone", "value": "+79001234567", "isVerified": false, "isPrimary": false}}
          ]
        }
        ```
    *   Требуемые права доступа: Владелец (`user`).
*   **`POST /me/contact-info`**
    *   Описание: Добавление нового контактного метода (email/телефон). Инициирует процесс верификации (отправка кода через Notification Service).
    *   Тело запроса: `{"data": {"type": "contactMethodCreation", "attributes": {"contactType": "email", "value": "new@example.com"}}}`.
    *   Пример ответа (Успех 201 Created): `[NEEDS DEVELOPER INPUT: Example success response for POST /me/contact-info]`
    *   Требуемые права доступа: Владелец (`user`).
*   **`DELETE /me/contact-info/{contact_id}`**
    *   Описание: Удаление контактного метода. Нельзя удалить основной метод, если он единственный верифицированный.
    *   Пример ответа (Успех 204 No Content): `[NEEDS DEVELOPER INPUT: Example success response for DELETE /me/contact-info/{id}]`
    *   Требуемые права доступа: Владелец (`user`).
*   **`POST /me/contact-info/{contact_id}/set-primary`**
    *   Описание: Установка контактного метода как основного (если он верифицирован).
    *   Пример ответа (Успех 200 OK): `[NEEDS DEVELOPER INPUT: Example success response for POST /me/contact-info/{id}/set-primary]`
    *   Требуемые права доступа: Владелец (`user`).
*   **`POST /me/contact-info/{contact_id}/request-verification`**
    *   Описание: Повторный запрос кода верификации.
    *   Пример ответа (Успех 200 OK): `[NEEDS DEVELOPER INPUT: Example success response for POST /me/contact-info/{id}/request-verification]`
    *   Требуемые права доступа: Владелец (`user`).
*   **`POST /me/contact-info/{contact_id}/verify`**
    *   Описание: Подтверждение контактного метода с помощью кода.
    *   Тело запроса: `{"data": {"type": "verificationConfirmation", "attributes": {"code": "123456"}}}`.
    *   Пример ответа (Успех 200 OK): `[NEEDS DEVELOPER INPUT: Example success response for POST /me/contact-info/{id}/verify]`
    *   Требуемые права доступа: Владелец (`user`).

#### 3.1.4. Ресурс: Настройки (Settings) - `/me/settings`
*   **`GET /me/settings`**
    *   Описание: Получение всех настроек текущего пользователя.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "userSettings",
            "id": "uuid-user-settings-id",
            "attributes": {
              "privacy": {"inventoryVisibility": "friends_only", "activityFeedVisibility": "public"},
              "notifications": {"newFriendRequestEmail": true, "gameUpdatePush": false},
              "interface": {"language": "ru", "theme": "dark"}
            }
          }
        }
        ```
    *   Требуемые права доступа: Владелец (`user`).
*   **`PUT /me/settings`**
    *   Описание: Обновление всех настроек текущего пользователя (полное обновление).
    *   Тело запроса: (Аналогично структуре ответа).
    *   Пример ответа (Успех 200 OK): `[NEEDS DEVELOPER INPUT: Example success response for PUT /me/settings]`
    *   Требуемые права доступа: Владелец (`user`).
*   **`PATCH /me/settings`**
    *   Описание: Частичное обновление настроек текущего пользователя.
    *   Тело запроса: (Только изменяемые категории/поля).
        ```json
        {
          "data": {
            "type": "userSettingsUpdate",
            "attributes": {
              "notifications": {"newFriendRequestEmail": false}
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): `[NEEDS DEVELOPER INPUT: Example success response for PATCH /me/settings]`
    *   Требуемые права доступа: Владелец (`user`).

#### 3.1.5. Административные эндпоинты
*   Префикс: `/admin/accounts` (доступ через API Gateway с проверкой роли `admin`).
*   **`GET /admin/accounts/{user_id}`**: Получение аккаунта и профиля пользователя по UserID (из Auth Service).
*   **`PUT /admin/accounts/{user_id}/status`**: Обновление статуса аккаунта (например, `active`, `blocked`, `pending_deletion`).
    *   Тело запроса: `{"data": {"type": "accountStatusUpdate", "attributes": {"status": "blocked", "reason": "Violation of terms"}}}`.
*   **`PUT /admin/accounts/{user_id}/profile`**: Редактирование профиля пользователя администратором.
*   Требуемые права доступа: `admin` (согласно `../../../../project_roles_and_permissions.md`).

### 3.2. gRPC API
*   **Пакет:** `account.v1`
*   **Аутентификация:** mTLS для межсервисного взаимодействия, опционально передача UserID/токена в метаданных для служебных запросов.

#### 3.2.1. Сервис: AccountService
*   Ссылка на `.proto` файл: `proto/account/v1/account_service.proto` (предполагаемое расположение в репозитории `platform-protos` или локально в `backend/account-service/proto/v1/account_service.proto`).
*   **`rpc GetAccountInfo(GetAccountInfoRequest) returns (GetAccountInfoResponse)`**
    *   Описание: Получение основной информации об аккаунте пользователя (статус, UserID).
    *   Сообщение запроса `GetAccountInfoRequest`: `message GetAccountInfoRequest { string user_id = 1; }`
    *   Сообщение ответа `GetAccountInfoResponse`: `message GetAccountInfoResponse { AccountInfo account_info = 1; } message AccountInfo { string account_id = 1; string user_id = 2; string status = 3; google.protobuf.Timestamp created_at = 4; }`
    *   Требуемые права доступа: [NEEDS DEVELOPER INPUT: Permissions for gRPC GetAccountInfo, e.g., internal service]
*   **`rpc GetUserProfile(GetUserProfileRequest) returns (GetUserProfileResponse)`**
    *   Описание: Получение профиля пользователя.
    *   Сообщение запроса `GetUserProfileRequest`: `message GetUserProfileRequest { string user_id = 1; }`
    *   Сообщение ответа `GetUserProfileResponse`: `message GetUserProfileResponse { Profile profile = 1; } message Profile { string profile_id = 1; string account_id = 2; string nickname = 3; string bio = 4; string avatar_url = 5; string country_code = 6; string custom_url_slug = 7; google.protobuf.Struct privacy_settings = 8; }`
    *   Требуемые права доступа: [NEEDS DEVELOPER INPUT: Permissions for gRPC GetUserProfile]
*   **`rpc GetUserProfilesBatch(GetUserProfilesBatchRequest) returns (GetUserProfilesBatchResponse)`**
    *   Описание: Получение профилей нескольких пользователей.
    *   Сообщение запроса `GetUserProfilesBatchRequest`: `message GetUserProfilesBatchRequest { repeated string user_ids = 1; }`
    *   Сообщение ответа `GetUserProfilesBatchResponse`: `message GetUserProfilesBatchResponse { repeated Profile profiles = 1; }`
    *   Требуемые права доступа: [NEEDS DEVELOPER INPUT: Permissions for gRPC GetUserProfilesBatch]
*   **`rpc CheckUsernameAvailability(CheckUsernameAvailabilityRequest) returns (CheckUsernameAvailabilityResponse)`**
    *   Описание: Проверка, доступен ли никнейм или кастомный URL профиля.
    *   Сообщение запроса `CheckUsernameAvailabilityRequest`: `message CheckUsernameAvailabilityRequest { string username = 1; }`
    *   Сообщение ответа `CheckUsernameAvailabilityResponse`: `message CheckUsernameAvailabilityResponse { bool is_available = 1; }`
    *   Требуемые права доступа: [NEEDS DEVELOPER INPUT: Permissions for gRPC CheckUsernameAvailability]
*   **`rpc GetUserSettings(GetUserSettingsRequest) returns (GetUserSettingsResponse)`**
    *   Описание: Получение настроек пользователя для использования другими сервисами.
    *   Сообщение запроса `GetUserSettingsRequest`: `message GetUserSettingsRequest { string user_id = 1; repeated string categories = 2; }`
    *   Сообщение ответа `GetUserSettingsResponse`: `message GetUserSettingsResponse { google.protobuf.Struct settings_data = 1; }`
    *   Требуемые права доступа: [NEEDS DEVELOPER INPUT: Permissions for gRPC GetUserSettings]

### 3.3. WebSocket API (если применимо)
*   Не предполагается для основного функционала Account Service. Обновления статуса пользователя или профиля могут передаваться через события Kafka другим сервисам, которые поддерживают WebSocket.

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **Account**:
    *   `id` (UUID, PK): Уникальный идентификатор аккаунта в Account Service.
    *   `user_id` (UUID, FK, Unique): Идентификатор пользователя из Auth Service. **Обязательное поле.**
    *   `status` (ENUM: `active`, `inactive`, `blocked`, `pending_deletion`, `deleted`): Статус аккаунта. **Обязательное поле.**
    *   `created_at` (TIMESTAMPTZ).
    *   `updated_at` (TIMESTAMPTZ).
*   **Profile**:
    *   `id` (UUID, PK).
    *   `account_id` (UUID, FK, Unique): Ссылка на Account. **Обязательное поле.**
    *   `nickname` (VARCHAR, Nullable, Unique): Отображаемое имя.
    *   `bio` (TEXT, Nullable).
    *   `avatar_url` (VARCHAR, Nullable).
    *   `country_code` (CHAR(2), Nullable): ISO 3166-1 alpha-2.
    *   `custom_url_slug` (VARCHAR, Nullable, Unique): Кастомный URL для профиля.
    *   `privacy_settings` (JSONB): Настройки приватности профиля.
    *   `created_at` (TIMESTAMPTZ).
    *   `updated_at` (TIMESTAMPTZ).
*   **ContactInfo**:
    *   `id` (UUID, PK).
    *   `account_id` (UUID, FK): Ссылка на Account. **Обязательное поле.**
    *   `type` (ENUM: `email`, `phone`). **Обязательное поле.**
    *   `value` (VARCHAR, Not Null). **Обязательное поле.**
    *   `is_verified` (BOOLEAN, Default: false).
    *   `is_primary` (BOOLEAN, Default: false).
    *   `verification_code_hash` (VARCHAR, Nullable): Хэш кода верификации.
    *   `verification_code_expires_at` (TIMESTAMPTZ, Nullable).
    *   `created_at` (TIMESTAMPTZ).
    *   `updated_at` (TIMESTAMPTZ).
*   **UserSetting**:
    *   `id` (UUID, PK).
    *   `account_id` (UUID, FK, Unique): Ссылка на Account. **Обязательное поле.**
    *   `settings_data` (JSONB): Все настройки пользователя, сгруппированные по категориям.
    *   `created_at` (TIMESTAMPTZ).
    *   `updated_at` (TIMESTAMPTZ).

### 4.2. Схема Базы Данных (PostgreSQL)
*   ERD Диаграмма:
    ```mermaid
    erDiagram
        ACCOUNT ||--o{ PROFILE : "has one"
        ACCOUNT ||--o{ CONTACT_INFO : "has many"
        ACCOUNT ||--o{ USER_SETTING : "has one"

        ACCOUNT {
            UUID id PK
            UUID user_id FK "Refers to Auth Service User ID, Unique"
            VARCHAR status "ENUM('active', 'inactive', 'blocked', 'pending_deletion', 'deleted')"
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }
        PROFILE {
            UUID id PK
            UUID account_id FK "Refers to ACCOUNT(id), Unique"
            VARCHAR nickname "Nullable, Unique"
            TEXT bio "Nullable"
            VARCHAR avatar_url "Nullable"
            CHAR(2) country_code "Nullable, ISO 3166-1 alpha-2"
            VARCHAR custom_url_slug "Nullable, Unique"
            JSONB privacy_settings
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }
        CONTACT_INFO {
            UUID id PK
            UUID account_id FK "Refers to ACCOUNT(id)"
            VARCHAR type "ENUM('email', 'phone')"
            VARCHAR value
            BOOLEAN is_verified "Default: false"
            BOOLEAN is_primary "Default: false"
            VARCHAR verification_code_hash "Nullable"
            TIMESTAMPTZ verification_code_expires_at "Nullable"
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }
        USER_SETTING {
            UUID id PK
            UUID account_id FK "Refers to ACCOUNT(id), Unique"
            JSONB settings_data
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }
    ```
*   Описание основных таблиц и индексов:
    ```sql
    -- Таблица: accounts
    CREATE TABLE accounts (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        user_id UUID NOT NULL UNIQUE, -- Внешний ключ на Auth Service (логическая связь, не физическая)
        status VARCHAR(50) NOT NULL DEFAULT 'inactive',
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
    CREATE INDEX idx_accounts_user_id ON accounts(user_id);
    CREATE INDEX idx_accounts_status ON accounts(status);

    -- Таблица: profiles
    CREATE TABLE profiles (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        account_id UUID NOT NULL UNIQUE REFERENCES accounts(id) ON DELETE CASCADE,
        nickname VARCHAR(255) UNIQUE,
        bio TEXT,
        avatar_url VARCHAR(2048),
        country_code CHAR(2),
        custom_url_slug VARCHAR(100) UNIQUE,
        privacy_settings JSONB DEFAULT '{}'::jsonb,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
    CREATE INDEX idx_profiles_nickname ON profiles(nickname text_pattern_ops);
    CREATE INDEX idx_profiles_custom_url_slug ON profiles(custom_url_slug text_pattern_ops);

    -- Таблица: contact_infos
    CREATE TABLE contact_infos (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
        type VARCHAR(20) NOT NULL, -- 'email', 'phone'
        value VARCHAR(255) NOT NULL,
        is_verified BOOLEAN NOT NULL DEFAULT FALSE,
        is_primary BOOLEAN NOT NULL DEFAULT FALSE,
        verification_code_hash VARCHAR(255),
        verification_code_expires_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        CONSTRAINT unique_contact_value_per_account_type UNIQUE (account_id, type, value),
        CONSTRAINT unique_primary_contact_per_account_type UNIQUE (account_id, type, is_primary) WHERE is_primary = TRUE
    );
    CREATE INDEX idx_contact_infos_account_id_type ON contact_infos(account_id, type);
    CREATE INDEX idx_contact_infos_value ON contact_infos(value);

    -- Таблица: user_settings
    CREATE TABLE user_settings (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        account_id UUID NOT NULL UNIQUE REFERENCES accounts(id) ON DELETE CASCADE,
        settings_data JSONB DEFAULT '{}'::jsonb,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
    CREATE INDEX idx_user_settings_account_id ON user_settings(account_id);
    ```

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   Используемая система сообщений: Kafka (согласно `../../../../project_technology_stack.md`).
*   Формат событий: CloudEvents JSON (согласно `../../../../project_api_standards.md`).
*   Основные топики Kafka: `ru.steam.account.events.v1` (или более гранулярные, если потребуется).

*   **`ru.steam.account.created.v1`**
    *   Описание: Публикуется после успешного создания базовой записи аккаунта.
    *   Топик: `ru.steam.account.events.v1`
    *   Структура Payload (`data` в CloudEvent):
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user-id",
          "status": "inactive",
          "createdAt": "ISO8601_timestamp"
        }
        ```
*   **`ru.steam.account.status.updated.v1`**
    *   Описание: Статус аккаунта изменен.
    *   Топик: `ru.steam.account.events.v1`
    *   Структура Payload:
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user-id",
          "oldStatus": "active",
          "newStatus": "blocked",
          "reason": "Violation of terms",
          "updatedAt": "ISO8601_timestamp"
        }
        ```
*   **`ru.steam.account.profile.updated.v1`**
    *   Описание: Профиль пользователя обновлен.
    *   Топик: `ru.steam.account.events.v1`
    *   Структура Payload:
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user_id",
          "profileId": "uuid-profile-id",
          "updatedFields": ["nickname", "avatarUrl"],
          "updatedAt": "ISO8601_timestamp"
        }
        ```
*   **`ru.steam.account.contact.added.v1`**
    *   Описание: Добавлен новый контактный метод.
    *   Топик: `ru.steam.account.events.v1`
    *   Структура Payload:
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user_id",
          "contactId": "uuid-contact-id",
          "type": "email",
          "value": "new_contact@example.com",
          "addedAt": "ISO8601_timestamp"
        }
        ```
*   **`ru.steam.account.contact.verified.v1`**
    *   Описание: Контактный метод успешно верифицирован.
    *   Топик: `ru.steam.account.events.v1`
    *   Структура Payload:
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user_id",
          "contactId": "uuid-contact-id",
          "type": "email",
          "value": "verified_contact@example.com",
          "verifiedAt": "ISO8601_timestamp"
        }
        ```
*   **`ru.steam.account.settings.updated.v1`**
    *   Описание: Настройки пользователя обновлены.
    *   Топик: `ru.steam.account.events.v1`
    *   Структура Payload:
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user_id",
          "updatedCategories": ["privacy", "notifications"],
          "updatedAt": "ISO8601_timestamp"
        }
        ```

### 5.2. Потребляемые События (Consumed Events)
*   **`ru.steam.auth.user.registered.v1`**
    *   Описание: От Auth Service. Инициирует создание аккаунта и профиля по умолчанию.
    *   Топик: `ru.steam.auth.events.v1` (согласно спецификации Auth Service)
    *   Структура Payload (`data` в CloudEvent):
        ```json
        {
          "userId": "uuid-auth-user-id",
          "email": "user@example.com",
          "username": "User123",
          "registrationTimestamp": "ISO8601_timestamp"
        }
        ```
    *   Логика обработки: Создать `Account`, `Profile`, `ContactInfo`. Опубликовать `ru.steam.account.created.v1`.
*   **`ru.steam.auth.user.email.updated.v1`**
    *   Описание: От Auth Service. Обновление основного email пользователя.
    *   Топик: `ru.steam.auth.events.v1`
    *   Структура Payload:
        ```json
        {
          "userId": "uuid-auth-user-id",
          "newEmail": "new_primary@example.com",
          "isVerified": true
        }
        ```
    *   Логика обработки: Обновить/создать `ContactInfo`, пометить как основной и верифицированный.
*   **`ru.steam.admin.user.action.v1`**
    *   Описание: От Admin Service для административных действий.
    *   Топик: `ru.steam.admin.events.v1`
    *   Структура Payload:
        ```json
        {
          "userId": "uuid-auth-user-id",
          "actionType": "BLOCK_ACCOUNT",
          "payload": { "reason": "Terms of Service violation" }
        }
        ```
    *   Логика обработки: Выполнить административное действие.

## 6. Интеграции (Integrations)
См. `../../../../project_integrations.md` для общей карты и деталей.
[NEEDS DEVELOPER INPUT: Review and confirm service boundaries and reliability strategies for each integration for account-service]

### 6.1. Внутренние Микросервисы
*   **Auth Service:**
    *   Тип интеграции: Потребление событий Kafka, gRPC вызовы (опционально).
    *   Назначение: Получение информации о зарегистрированных пользователях.
    *   Контракт: `ru.steam.auth.user.registered.v1`, `GetUserByID` (gRPC).
    *   Надежность: [NEEDS DEVELOPER INPUT: Reliability strategy for Auth Service integration]
*   **Notification Service:**
    *   Тип интеграции: События Kafka или gRPC вызовы.
    *   Назначение: Отправка пользовательских уведомлений.
    *   Контракт: `ru.steam.account.contact.verification.requested.v1` (Kafka), `SendEmail` (gRPC).
    *   Надежность: [NEEDS DEVELOPER INPUT: Reliability strategy for Notification Service integration]
*   **Social Service:**
    *   Тип интеграции: Потребление событий Kafka, gRPC вызовы.
    *   Назначение: Обмен информацией о профиле.
    *   Контракт: `ru.steam.account.profile.updated.v1` (Kafka), `CheckNicknameUniqueness` (gRPC).
    *   Надежность: [NEEDS DEVELOPER INPUT: Reliability strategy for Social Service integration]
*   **Payment Service:**
    *   Тип интеграции: gRPC вызовы.
    *   Назначение: Предоставление информации о пользователе.
    *   Контракт: `GetAccountInfo`, `GetUserProfile` (gRPC).
    *   Надежность: [NEEDS DEVELOPER INPUT: Reliability strategy for Payment Service integration]
*   **Admin Service:**
    *   Тип интеграции: REST/gRPC эндпоинты.
    *   Назначение: Администрирование пользователей.
    *   Контракт: См. раздел 3.1.5.
    *   Надежность: [NEEDS DEVELOPER INPUT: Reliability strategy for Admin Service integration]
*   **API Gateway:**
    *   Тип интеграции: Проксирование REST запросов.
    *   Назначение: Единая точка входа.
    *   Контракт: REST API Account Service.
    *   Надежность: [NEEDS DEVELOPER INPUT: Reliability strategy for API Gateway integration]

### 6.2. Внешние Системы
*   **S3-совместимое хранилище (Yandex Object Storage, VK Cloud, etc.):**
    *   Тип интеграции: SDK.
    *   Назначение: Хранение аватаров пользователей (если управляется сервисом).
    *   Контракт: S3 API.
*   **PostgreSQL:**
    *   Тип интеграции: Прямое подключение (DB Driver).
    *   Назначение: Основное хранилище данных.
*   **Redis:**
    *   Тип интеграции: Прямое подключение (Redis Client).
    *   Назначение: Кэширование.
*   **Kafka:**
    *   Тип интеграции: Прямое подключение (Kafka Client).
    *   Назначение: Брокер сообщений.

## 7. Конфигурация (Configuration)
Общие стандарты конфигурационных файлов определены в `../../../../project_api_standards.md` (раздел 7) и `../../../../DOCUMENTATION_GUIDELINES.md` (раздел 6).

### 7.1. Переменные Окружения
*   `ACCOUNT_HTTP_PORT`: Порт для REST API (например, `8081`)
*   `ACCOUNT_GRPC_PORT`: Порт для gRPC API (например, `9091`)
*   `POSTGRES_DSN`: Строка подключения к PostgreSQL.
*   `REDIS_ADDR`: Адрес Redis.
*   `REDIS_PASSWORD`: Пароль для Redis.
*   `REDIS_DB_ACCOUNT`: Номер базы данных Redis.
*   `KAFKA_BROKERS`: Список брокеров Kafka.
*   `KAFKA_TOPIC_ACCOUNT_EVENTS`: Имя топика для публикуемых событий.
*   `KAFKA_CONSUMER_GROUP_ACCOUNT`: Имя группы потребителей Kafka.
*   `LOG_LEVEL`: Уровень логирования (например, `info`).
*   `AVATAR_MAX_SIZE_BYTES`: Максимальный размер аватара.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`: Endpoint для экспорта трейсов.
*   `DEFAULT_USER_SETTINGS_JSON`: JSON с настройками по умолчанию.
*   `JWT_PUBLIC_KEY_PATH`: [NEEDS DEVELOPER INPUT: Clarify if needed, or remove if API Gateway handles all JWT validation]

### 7.2. Файлы Конфигурации (`configs/config.yaml`)
*   Расположение: `backend/account-service/configs/config.yaml`
*   Структура:
    ```yaml
    http_server:
      port: ${ACCOUNT_HTTP_PORT:"8081"}
      timeout_seconds: 30
    grpc_server:
      port: ${ACCOUNT_GRPC_PORT:"9091"}
      timeout_seconds: 30
    postgres:
      dsn: ${POSTGRES_DSN}
      pool_max_conns: 10
    redis:
      address: ${REDIS_ADDR}
      password: ${REDIS_PASSWORD:""}
      db: ${REDIS_DB_ACCOUNT:0}
    kafka:
      brokers: ${KAFKA_BROKERS}
      producer_topics:
        account_events: ${KAFKA_TOPIC_ACCOUNT_EVENTS:"ru.steam.account.events.v1"}
      consumer_topics:
        auth_events: ${KAFKA_TOPIC_AUTH_EVENTS:"ru.steam.auth.events.v1"}
      consumer_group: ${KAFKA_CONSUMER_GROUP_ACCOUNT:"account-service-consumer-group"}
    logging:
      level: ${LOG_LEVEL:"info"}
      format: "json"
    security:
      jwt_public_key_path: ${JWT_PUBLIC_KEY_PATH:""} # [NEEDS DEVELOPER INPUT: Confirm or remove]
    avatar:
      max_size_bytes: ${AVATAR_MAX_SIZE_BYTES:5242880} # 5MB
      # s3: # [NEEDS DEVELOPER INPUT: Uncomment and fill if service handles S3 uploads directly]
      #   bucket: ${AVATAR_S3_BUCKET:""}
      #   region: ${AVATAR_S3_REGION:""}
      #   endpoint: ${AVATAR_S3_ENDPOINT:""}
      #   access_key_id: ${AVATAR_S3_ACCESS_KEY_ID:""}
      #   secret_access_key: ${AVATAR_S3_SECRET_ACCESS_KEY:""}
    otel:
      exporter_jaeger_endpoint: ${OTEL_EXPORTER_JAEGER_ENDPOINT:""}
      service_name: "account-service"
    default_settings:
      user_settings_json: ${DEFAULT_USER_SETTINGS_JSON:'{"privacy": {"inventoryVisibility": "private"}, "notifications": {"newFriendRequestEmail": true}}'}
    ```

## 8. Обработка Ошибок (Error Handling)
Следование `../../../../project_api_standards.md` для форматов ошибок REST и gRPC.

### 8.1. Общие Принципы
*   Стандартный формат ответа об ошибке (JSON).
*   Использование стандартных HTTP кодов состояния для REST API.
*   Использование стандартных кодов gRPC для gRPC API.
*   Подробное логирование ошибок с `trace_id` и контекстом.

### 8.2. Распространенные Коды Ошибок
*   **`ACCOUNT_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`): Аккаунт не найден.
*   **`PROFILE_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`): Профиль не найден.
*   **`CONTACT_INFO_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`): Контактная информация не найдена.
*   **`SETTINGS_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`): Настройки не найдены.
*   **`NICKNAME_TAKEN`** (HTTP 409, gRPC `ALREADY_EXISTS`): Никнейм уже занят.
*   **`CUSTOM_URL_TAKEN`** (HTTP 409, gRPC `ALREADY_EXISTS`): Кастомный URL уже занят.
*   **`EMAIL_ALREADY_EXISTS_FOR_USER`** (HTTP 409, gRPC `ALREADY_EXISTS`): Email уже используется этим пользователем.
*   **`PHONE_ALREADY_EXISTS_FOR_USER`** (HTTP 409, gRPC `ALREADY_EXISTS`): Телефон уже используется этим пользователем.
*   **`VERIFICATION_CODE_INVALID`** (HTTP 400, gRPC `INVALID_ARGUMENT`): Неверный код верификации.
*   **`VERIFICATION_CODE_EXPIRED`** (HTTP 400, gRPC `INVALID_ARGUMENT`): Срок действия кода верификации истек.
*   **`CANNOT_DELETE_PRIMARY_CONTACT`** (HTTP 400, gRPC `FAILED_PRECONDITION`): Нельзя удалить основной контакт.
*   **`CONTACT_NOT_VERIFIED`** (HTTP 400, gRPC `FAILED_PRECONDITION`): Контакт не верифицирован.
*   **`INVALID_AVATAR_FILE`** (HTTP 400, gRPC `INVALID_ARGUMENT`): Некорректный файл аватара.
*   **`UPDATE_CONFLICT`** (HTTP 409, gRPC `ABORTED`): Конфликт одновременного обновления.
*   Пример ответа (для `NICKNAME_TAKEN`):
    ```json
    {
      "error": {
        "code": "NICKNAME_TAKEN",
        "message": "Выбранный никнейм уже занят.",
        "details": { "nickname": "User123" }
      }
    }
    ```

## 9. Безопасность (Security)
См. `../../../../project_security_standards.md` для общих стандартов.

### 9.1. Аутентификация
*   REST API: JWT аутентификация (проверяется API Gateway).
*   gRPC API: mTLS для межсервисного взаимодействия. UserID передается в метаданных.

### 9.2. Авторизация
*   Проверка владения ресурсом (на основе `X-User-ID`).
*   Административные эндпоинты требуют роль `admin` (на основе `X-User-Roles`).
*   Детализация ролей: см. `../../../../project_roles_and_permissions.md`.

### 9.3. Защита Данных
*   **ФЗ-152 "О персональных данных":**
    *   Хранение ПДн в РФ.
    *   Согласие на обработку ПДн.
    *   Логирование доступа к ПДн.
    *   Шифрование при передаче (TLS).
    *   Хеширование кодов верификации.
    *   Анонимизация/удаление данных.
*   Валидация входных данных.
*   Ограничение на размер и тип загружаемых файлов.

### 9.4. Управление Секретами
*   Использование Kubernetes Secrets или HashiCorp Vault (согласно `../../../../project_security_standards.md`).
*   Доступ через переменные окружения или API Vault.

## 10. Развертывание (Deployment)
См. `../../../../project_deployment_standards.md` для общих стандартов.

### 10.1. Инфраструктурные Файлы
*   Dockerfile: `backend/account-service/Dockerfile`
*   Helm-чарты/Kubernetes манифесты: `deploy/charts/account-service/` (предполагаемое расположение).

### 10.2. Зависимости при Развертывании
*   PostgreSQL.
*   Redis.
*   Kafka.
*   Auth Service.
*   API Gateway.
*   Jaeger/OpenTelemetry Collector (опционально).

### 10.3. CI/CD
*   Стандартный пайплайн CI/CD (сборка, тесты, SAST/DAST, сборка образа, публикация, деплой).
*   Автоматическое применение миграций БД.
*   Файлы CI/CD: `.github/workflows/account-service.yml` или аналогичный.

## 11. Мониторинг и Логирование (Logging and Monitoring)
См. `../../../../project_observability_standards.md` для общих стандартов.

### 11.1. Логирование
*   Формат: JSON (Zap).
*   Уровни: DEBUG, INFO, WARN, ERROR, FATAL.
*   Ключевые поля: `timestamp`, `level`, `service`, `instance`, `message`, `trace_id`, `span_id`, `user_id`, `error_details`.
*   Интеграция: Fluent Bit для сбора логов в ELK/Loki.

### 11.2. Мониторинг
*   Метрики (Prometheus): HTTP, gRPC, База данных, Kafka, Кэш, Бизнес-метрики.
*   Эндпоинт для метрик: `/metrics`
*   Дашборды (Grafana): Стандартный дашборд для Account Service.
*   Алерты (AlertManager): Высокий % ошибок, большая задержка, проблемы с подключением к зависимостям, ошибки Kafka.

### 11.3. Трассировка
*   Интеграция: OpenTelemetry SDK для Go.
*   Экспорт: Jaeger.
*   Контекст трассировки: передается через HTTP/gRPC запросы и Kafka события.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**:
    *   API чтения профиля (`GET /me`): P95 < 100мс.
    *   API обновления профиля (`PUT /me/profile`): P95 < 200мс.
    *   gRPC `GetAccountInfo`: P99 < 50мс.
    *   [NEEDS DEVELOPER INPUT: Confirm or update performance NFRs for account-service]
*   **Надежность**:
    *   Доступность > 99.95%.
    *   RTO/RPO: См. раздел "Резервное Копирование и Восстановление".
    *   [NEEDS DEVELOPER INPUT: Confirm or update reliability NFRs for account-service]
*   **Масштабируемость**:
    *   Горизонтальное масштабирование для поддержки >10 млн активных пользователей.
    *   Способность обрабатывать до 1000 RPS на чтение и 200 RPS на запись.
    *   [NEEDS DEVELOPER INPUT: Confirm or update scalability NFRs for account-service]
*   **Сопровождаемость**:
    *   Покрытие кода unit-тестами > 80%.
    *   Покрытие интеграционными тестами > 70%.
    *   [NEEDS DEVELOPER INPUT: Confirm or update maintainability NFRs for account-service]
*   **Безопасность**: Соответствие `../../../../project_security_standards.md` и ФЗ-152.

## 13. Резервное Копирование и Восстановление (Backup and Recovery)
Общие принципы см. в `../../../../project_database_structure.md`.

### 13.1. PostgreSQL
*   **Процедура резервного копирования:** Ежедневный `pg_dumpall`, непрерывная архивация WAL (PITR).
*   **Частота:** Базовый бэкап еженедельно, WAL непрерывно.
*   **Хранение:** S3, шифрование, версионирование, другой регион. Срок хранения: полные - 30 дней, WAL - 14 дней.
*   **Процедура восстановления:** Тестируется ежеквартально.
*   **RTO (Recovery Time Objective):** < 4 часов.
*   **RPO (Recovery Point Objective):** < 5 минут.

### 13.2. Redis
*   **Процедура резервного копирования:** RDB Snapshots каждые 6 часов, AOF с fsync `everysec`.
*   **Хранение:** S3, ежедневно. Срок хранения - 7 дней.
*   **Процедура восстановления:** Из RDB + AOF. Тестируется ежеквартально.
*   **RTO:** < 1 часа.
*   **RPO:** < 1 минуты. (Может быть менее строгим, т.к. в основном кэш).

### 13.3. Общая стратегия
*   Часть общей стратегии BCDR. Документировано, регулярно пересматривается. Мониторинг бэкапов.

## 14. Приложения (Appendices) (Опционально)
*   Детальные схемы Protobuf: репозиторий `platform-protos` или локально `backend/account-service/proto/v1/`.
*   [NEEDS DEVELOPER INPUT: Add full JSON examples for DTOs if required, or link to OpenAPI spec for account-service]

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [Регистрация пользователя и начальная настройка профиля](../../../../project_workflows/user_registration_flow.md)
*   [NEEDS DEVELOPER INPUT: Add links to other relevant high-level workflow documents for account-service]

---
*Этот документ является отправной точкой и должен регулярно обновляться по мере развития сервиса.*
