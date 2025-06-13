# Спецификация Микросервиса: Account Service

**Версия:** 1.0
**Дата последнего обновления:** 2023-10-28

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   Account Service является основным микросервисом платформы "Российский Аналог Steam", отвечающим за управление учетными записями пользователей, их профилями и связанными данными, такими как настройки, информация для верификации и контактная информация.
*   Он служит авторитетным источником данных, связанных с пользователями, предоставляя эту информацию другим сервисам и обеспечивая целостность и безопасность данных.
*   Основные бизнес-задачи: управление жизненным циклом аккаунта пользователя, управление профильной информацией, управление контактными данными и их верификацией, управление пользовательскими настройками.

### 1.2. Ключевые Функциональности
*   **Управление учетными записями:** Регистрация новых пользователей (координируется с Auth Service), хранение и управление основной информацией об аккаунте (ID, статус), поддержка изменения статуса аккаунта (активация, блокировка, удаление).
*   **Управление профилями пользователей:** Создание и редактирование деталей профиля (никнейм, биография, страна, аватар), управление настройками видимости профиля.
*   **Управление контактной информацией и верификацией:** Хранение и управление контактными данными пользователя (email, телефон), обработка процесса верификации.
*   **Управление настройками:** Хранение и обновление пользовательских настроек, разделенных по категориям (например, приватность, уведомления, интерфейс).
*   **Генерация событий:** Публикация событий, связанных с изменениями в аккаунтах и профилях (например, создание аккаунта, обновление профиля) в брокер сообщений Kafka.

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+)
*   **REST Framework:** Gin или Echo (согласно `project_technology_stack.md`, предпочтительно Echo)
*   **gRPC Framework:** google.golang.org/grpc
*   **База данных:** PostgreSQL (версия 15+)
*   **Кэширование/Управление сессиями:** Redis
*   **Брокер сообщений:** Kafka (например, confluent-kafka-go или sarama)
*   **Валидация:** go-playground/validator
*   **ORM/DB Driver:** GORM или sqlx (согласно `project_technology_stack.md`, предпочтительно GORM или pgx)
*   (Ссылки на `project_technology_stack.md`, `project_glossary.md`)

### 1.4. Термины и Определения (Glossary)
*   См. `project_glossary.md`.
*   **Аккаунт (Account):** Основная учетная запись пользователя.
*   **Профиль (Profile):** Публичная и приватная информация пользователя.
*   **Контактная информация (ContactInfo):** Email, телефон пользователя.
*   **Настройки (Settings):** Пользовательские конфигурации.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Сервис построен с использованием принципов чистой архитектуры (Clean Architecture) или слоистой архитектуры для разделения ответственностей.
*   Основные модули: управление аккаунтами, управление профилями, управление контактами, управление настройками.
*   [Диаграмма верхнеуровневой архитектуры сервиса будет добавлена в будущих версиях документации.]

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления)
*   Ответственность: Обработка входящих HTTP (REST) и gRPC запросов, валидация DTO, вызов Application Layer.
*   Ключевые компоненты/модули:
    *   HTTP Handlers (Echo/Gin): Эндпоинты для CRUD операций над аккаунтами, профилями, контактами, настройками.
    *   gRPC Service Implementations: Методы для межсервисного взаимодействия (например, `GetAccount`, `GetProfile`).
    *   DTOs: Структуры для передачи данных (например, `CreateAccountRequest`, `UserProfileResponse`).
    *   Валидаторы: Для проверки корректности входных данных.

#### 2.2.2. Application Layer (Прикладной Слой / Слой Сценариев Использования)
*   Ответственность: Координация бизнес-логики, реализация use cases.
*   Ключевые компоненты/модули:
    *   Use Case Services: `AccountUseCaseService`, `ProfileUseCaseService`, `ContactUseCaseService`, `SettingsUseCaseService`.
    *   Интерфейсы для репозиториев (`AccountRepository`, `ProfileRepository`) и внешних сервисов (например, `NotificationService`).

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности, агрегаты, доменные события.
*   Ключевые компоненты/модули:
    *   Entities: `Account` (ID, UserID (из Auth Service), Status, Timestamps), `Profile` (Nickname, Bio, AvatarURL, Country), `ContactInfo` (Email, Phone, VerifiedStatus), `Setting` (Category, ValuesJSONB).
    *   Domain Events: `AccountCreatedEvent`, `ProfileUpdatedEvent`, `ContactVerifiedEvent`.
    *   Интерфейсы репозиториев.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Взаимодействие с PostgreSQL, Redis, Kafka.
*   Ключевые компоненты/модули:
    *   PostgreSQL Repositories: Реализации `AccountRepository`, `ProfileRepository` etc.
    *   Redis Cache: Для кэширования профилей, настроек.
    *   Kafka Producer: Для отправки `AccountCreatedEvent`, `ProfileUpdatedEvent`.
    *   Клиенты для других сервисов (например, gRPC клиент для Auth Service для получения UserID по необходимости).

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL (через API Gateway):** `/api/v1/users` или `/api/v1/accounts` (будет уточнено в `project_api_standards.md` или конфигурации API Gateway).
*   **Версионирование:** `/v1/`.
*   **Формат данных:** `application/json`.
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway, UserID передается в заголовке, например, `X-User-ID`).
*   **Стандартные заголовки:** `X-Request-ID`.
*   (Общие принципы см. `project_api_standards.md`)

#### 3.1.1. Ресурс: Аккаунты (Accounts)
*   **`GET /me`**
    *   Описание: Получение информации о текущем аутентифицированном пользователе (его аккаунте).
    *   Query параметры: `include=profile,settings` (опционально, для включения связанных данных).
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "id": "uuid-account-id",
            "type": "account",
            "attributes": {
              "user_id": "uuid-auth-user-id",
              "status": "active",
              "created_at": "2023-10-28T10:00:00Z",
              "updated_at": "2023-10-28T10:00:00Z"
            },
            "relationships": { // если запрошено через include
              "profile": { "data": { "type": "profile", "id": "uuid-profile-id" } },
              "settings": { "data": { "type": "settings", "id": "uuid-settings-id" } }
            }
          }
        }
        ```
    *   Требуемые права доступа: Аутентифицированный пользователь.
*   **`PUT /me`**
    *   Описание: Обновление основной информации аккаунта текущего пользователя (например, смена основного email после верификации нового, если это разрешено здесь, а не в Auth). [Подлежит уточнению, какие поля аккаунта изменяемы].
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "account",
            "attributes": {
              // Поля для обновления
            }
          }
        }
        ```
    *   Требуемые права доступа: Владелец аккаунта.
*   **`DELETE /me`**
    *   Описание: Запрос на удаление аккаунта текущего пользователя (может инициировать процесс удаления или деактивации).
    *   Требуемые права доступа: Владелец аккаунта.
*   **Административные эндпоинты (префикс `/admin/accounts`):**
    *   `GET /admin/accounts/{user_id}`: Получение аккаунта пользователя по UserID (из Auth Service).
    *   `PUT /admin/accounts/{user_id}/status`: Обновление статуса аккаунта (например, `active`, `blocked`, `pending_deletion`).
        *   Тело запроса: `{"data": {"type": "accountStatus", "attributes": {"status": "blocked", "reason": "Violation of terms"}}}`.
    *   Требуемые права доступа: `admin`.

#### 3.1.2. Ресурс: Профили (Profiles)
*   **`GET /me/profile`**
    *   Описание: Получение профиля текущего пользователя.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "id": "uuid-profile-id",
            "type": "profile",
            "attributes": {
              "nickname": "User123",
              "bio": "Hello world!",
              "avatar_url": "https://example.com/avatar.jpg",
              "country": "RU",
              "custom_url": "user123_profile",
              "privacy_settings": { "show_real_name": false, "inventory_public": true }
            }
          }
        }
        ```
    *   Требуемые права доступа: Владелец профиля.
*   **`PUT /me/profile`**
    *   Описание: Обновление профиля текущего пользователя.
    *   Тело запроса: (Аналогично структуре ответа, но только изменяемые поля).
    *   Требуемые права доступа: Владелец профиля.
*   **`POST /me/profile/avatar`**
    *   Описание: Загрузка нового аватара для текущего пользователя. (multipart/form-data).
    *   Требуемые права доступа: Владелец профиля.
*   **`GET /profiles/{profile_id_or_custom_url}`**
    *   Описание: Получение публичного профиля пользователя по его ID профиля или кастомному URL.
    *   Требуемые права доступа: `anonymous` (для публичных данных), `user` (для данных, видимых другим пользователям).

#### 3.1.3. Ресурс: Контактная Информация (Contact Info)
*   **`GET /me/contact-info`**
    *   Описание: Получение контактной информации текущего пользователя.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {"type": "email", "id": "uuid-email-id", "attributes": {"value": "user@example.com", "verified": true, "is_primary": true}},
            {"type": "phone", "id": "uuid-phone-id", "attributes": {"value": "+79001234567", "verified": false, "is_primary": false}}
          ]
        }
        ```
    *   Требуемые права доступа: Владелец.
*   **`POST /me/contact-info`**
    *   Описание: Добавление нового контактного метода (email/телефон). Инициирует процесс верификации.
    *   Тело запроса: `{"data": {"type": "email/phone", "attributes": {"value": "new@example.com"}}}`.
    *   Требуемые права доступа: Владелец.
*   **`DELETE /me/contact-info/{contact_id}`**
    *   Описание: Удаление контактного метода.
    *   Требуемые права доступа: Владелец.
*   **`POST /me/contact-info/{contact_id}/set-primary`**
    *   Описание: Установка контактного метода как основного (если он верифицирован).
    *   Требуемые права доступа: Владелец.
*   **`POST /me/contact-info/{contact_id}/request-verification`**
    *   Описание: Повторный запрос кода верификации.
    *   Требуемые права доступа: Владелец.
*   **`POST /me/contact-info/{contact_id}/verify`**
    *   Описание: Подтверждение контактного метода с помощью кода.
    *   Тело запроса: `{"data": {"type": "verification", "attributes": {"code": "123456"}}}`.
    *   Требуемые права доступа: Владелец.

#### 3.1.4. Ресурс: Настройки (Settings)
*   **`GET /me/settings`**
    *   Описание: Получение всех настроек текущего пользователя.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "settings",
            "id": "uuid-user-settings-id",
            "attributes": {
              "privacy": {"inventory_visibility": "friends_only", "activity_feed_visibility": "public"},
              "notifications": {"new_friend_request_email": true, "game_update_push": false},
              "interface": {"language": "ru", "theme": "dark"}
            }
          }
        }
        ```
    *   Требуемые права доступа: Владелец.
*   **`PUT /me/settings`**
    *   Описание: Обновление настроек текущего пользователя.
    *   Тело запроса: (Аналогично структуре ответа, но только изменяемые поля).
    *   Требуемые права доступа: Владелец.
*   **`GET /me/settings/{category}`** (например, `/me/settings/privacy`)
    *   Описание: Получение настроек для конкретной категории.
    *   Требуемые права доступа: Владелец.
*   **`PUT /me/settings/{category}`**
    *   Описание: Обновление настроек для конкретной категории.
    *   Требуемые права доступа: Владелец.

### 3.2. gRPC API
*   **Пакет:** `account.v1`
*   **Файл .proto:** `proto/account/v1/account_service.proto` (предполагаемое расположение)
*   **Аутентификация:** mTLS для межсервисного взаимодействия.

#### 3.2.1. Сервис: AccountService
*   **`rpc GetAccountInfo(GetAccountInfoRequest) returns (GetAccountInfoResponse)`**
    *   Описание: Получение основной информации об аккаунте пользователя.
    *   `GetAccountInfoRequest`: `{ string user_id }` (UserID из Auth Service)
    *   `GetAccountInfoResponse`: `{ Account account }` (содержит ID аккаунта, статус и т.д.)
*   **`rpc GetUserProfile(GetUserProfileRequest) returns (GetUserProfileResponse)`**
    *   Описание: Получение профиля пользователя.
    *   `GetUserProfileRequest`: `{ string user_id }`
    *   `GetUserProfileResponse`: `{ Profile profile }`
*   **`rpc GetUserProfiles(GetUserProfilesRequest) returns (GetUserProfilesResponse)`**
    *   Описание: Получение профилей нескольких пользователей.
    *   `GetUserProfilesRequest`: `{ repeated string user_ids }`
    *   `GetUserProfilesResponse`: `{ repeated Profile profiles }`
*   **`rpc CheckUsernameAvailability(CheckUsernameAvailabilityRequest) returns (CheckUsernameAvailabilityResponse)`**
    *   Описание: Проверка, доступен ли никнейм (кастомный URL профиля).
    *   `CheckUsernameAvailabilityRequest`: `{ string username }`
    *   `CheckUsernameAvailabilityResponse`: `{ bool is_available }`
*   **`rpc GetUserSettings(GetUserSettingsRequest) returns (GetUserSettingsResponse)`**
    *   Описание: Получение настроек пользователя для использования другими сервисами (например, Notification Service для проверки предпочтений).
    *   `GetUserSettingsRequest`: `{ string user_id, repeated string categories }`
    *   `GetUserSettingsResponse`: `{ map<string, google.protobuf.Struct> settings_by_category }`

### 3.3. WebSocket API (если применимо)
*   Не предполагается для основного функционала Account Service. Обновления статуса пользователя или профиля могут передаваться через события Kafka другим сервисам, которые поддерживают WebSocket (например, Social Service).

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности
*   **Account**:
    *   `id` (UUID, PK): Уникальный идентификатор аккаунта в Account Service.
    *   `user_id` (UUID, FK, Unique): Идентификатор пользователя из Auth Service.
    *   `status` (ENUM: `active`, `inactive`, `blocked`, `pending_deletion`, `deleted`): Статус аккаунта.
    *   `created_at` (TIMESTAMPTZ).
    *   `updated_at` (TIMESTAMPTZ).
*   **Profile**:
    *   `id` (UUID, PK).
    *   `account_id` (UUID, FK, Unique): Ссылка на Account.
    *   `nickname` (VARCHAR, Nullable, Unique): Отображаемое имя.
    *   `bio` (TEXT, Nullable).
    *   `avatar_url` (VARCHAR, Nullable).
    *   `country_code` (CHAR(2), Nullable): ISO 3166-1 alpha-2.
    *   `custom_url_slug` (VARCHAR, Nullable, Unique): Кастомный URL для профиля.
    *   `privacy_settings` (JSONB): Настройки приватности профиля (например, `{"inventory_visibility": "public", "friends_list_visibility": "private"}`).
    *   `created_at` (TIMESTAMPTZ).
    *   `updated_at` (TIMESTAMPTZ).
*   **ContactInfo**:
    *   `id` (UUID, PK).
    *   `account_id` (UUID, FK): Ссылка на Account.
    *   `type` (ENUM: `email`, `phone`).
    *   `value` (VARCHAR, Not Null).
    *   `is_verified` (BOOLEAN, Default: false).
    *   `is_primary` (BOOLEAN, Default: false).
    *   `verification_code` (VARCHAR, Nullable): Хэш кода верификации.
    *   `verification_code_expires_at` (TIMESTAMPTZ, Nullable).
    *   `created_at` (TIMESTAMPTZ).
    *   `updated_at` (TIMESTAMPTZ).
    *   (Constraint: Unique(account_id, type, value))
*   **UserSetting**:
    *   `id` (UUID, PK).
    *   `account_id` (UUID, FK, Unique): Ссылка на Account.
    *   `settings_data` (JSONB): Все настройки пользователя, сгруппированные по категориям (например, `{"notifications": {"email_updates": true}, "interface": {"language": "ru"}}`).
    *   `created_at` (TIMESTAMPTZ).
    *   `updated_at` (TIMESTAMPTZ).
*   **Avatar**: (Может быть частью Profile или отдельной сущностью если требуется история/управление)
    *   `id` (UUID, PK).
    *   `account_id` (UUID, FK).
    *   `file_path_s3` (VARCHAR).
    *   `uploaded_at` (TIMESTAMPTZ).

### 4.2. Схема Базы Данных (PostgreSQL)
*   Диаграмма ERD: [Mermaid ERD диаграмма будет добавлена в будущих версиях документации.]
    ```mermaid
    erDiagram
        ACCOUNT ||--o{ PROFILE : "has one"
        ACCOUNT ||--o{ CONTACT_INFO : "has many"
        ACCOUNT ||--o{ USER_SETTING : "has one"

        ACCOUNT {
            UUID id PK
            UUID user_id FK "Refers to Auth Service User"
            VARCHAR status
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }
        PROFILE {
            UUID id PK
            UUID account_id FK
            VARCHAR nickname
            TEXT bio
            VARCHAR avatar_url
            CHAR(2) country_code
            VARCHAR custom_url_slug
            JSONB privacy_settings
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }
        CONTACT_INFO {
            UUID id PK
            UUID account_id FK
            VARCHAR type "ENUM('email', 'phone')"
            VARCHAR value
            BOOLEAN is_verified
            BOOLEAN is_primary
            VARCHAR verification_code_hash
            TIMESTAMPTZ verification_code_expires_at
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }
        USER_SETTING {
            UUID id PK
            UUID account_id FK
            JSONB settings_data
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }
    ```
*   Описание основных таблиц:
    ```sql
    -- Таблица: accounts
    CREATE TABLE accounts (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        user_id UUID NOT NULL UNIQUE, -- Внешний ключ на Auth Service (логическая связь)
        status VARCHAR(50) NOT NULL DEFAULT 'pending_verification_email', -- e.g., pending_verification_email, active, blocked, deleted
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
    CREATE INDEX idx_accounts_user_id ON accounts(user_id);

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
    CREATE INDEX idx_profiles_nickname ON profiles(nickname);
    CREATE INDEX idx_profiles_custom_url_slug ON profiles(custom_url_slug);

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
        UNIQUE (account_id, type, value),
        UNIQUE (account_id, type, is_primary) WHERE is_primary = TRUE -- Только один основной email/телефон
    );
    CREATE INDEX idx_contact_infos_account_id_type ON contact_infos(account_id, type);

    -- Таблица: user_settings
    CREATE TABLE user_settings (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        account_id UUID NOT NULL UNIQUE REFERENCES accounts(id) ON DELETE CASCADE,
        settings_data JSONB DEFAULT '{}'::jsonb,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
    ```

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Kafka (согласно `project_technology_stack.md`).
*   **Формат событий:** CloudEvents JSON (согласно `project_api_standards.md`).
*   **Основные топики:** `account.events` (или более гранулярные, например, `account.profile.events`, `account.settings.events`).

*   **`account.created.v1`**
    *   Описание: Публикуется после успешного создания базовой записи аккаунта (обычно после события `auth.user.registered` от Auth Service).
    *   Топик: `account.events`
    *   Структура Payload:
        ```json
        {
          "account_id": "uuid-account-id",
          "user_id": "uuid-auth-user-id",
          "status": "pending_verification_email", // Начальный статус
          "created_at": "ISO8601_timestamp"
        }
        ```
*   **`account.status.updated.v1`**
    *   Описание: Статус аккаунта изменен (например, активирован, заблокирован).
    *   Топик: `account.events`
    *   Структура Payload:
        ```json
        {
          "account_id": "uuid-account-id",
          "user_id": "uuid-auth-user-id",
          "old_status": "active",
          "new_status": "blocked",
          "reason": "Violation of terms", // Опционально
          "updated_at": "ISO8601_timestamp"
        }
        ```
*   **`account.profile.updated.v1`**
    *   Описание: Профиль пользователя обновлен.
    *   Топик: `account.events` (или `account.profile.events`)
    *   Структура Payload:
        ```json
        {
          "account_id": "uuid-account-id",
          "user_id": "uuid-auth-user_id",
          "profile_id": "uuid-profile-id",
          "updated_fields": ["nickname", "avatar_url"], // Список измененных полей
          "updated_at": "ISO8601_timestamp"
        }
        ```
*   **`account.contact.added.v1`**
    *   Описание: Добавлен новый контактный метод.
    *   Топик: `account.events`
    *   Структура Payload:
        ```json
        {
          "account_id": "uuid-account-id",
          "user_id": "uuid-auth-user_id",
          "contact_id": "uuid-contact-id",
          "type": "email", // "email" or "phone"
          "value": "new_contact@example.com",
          "added_at": "ISO8601_timestamp"
        }
        ```
*   **`account.contact.verified.v1`**
    *   Описание: Контактный метод успешно верифицирован.
    *   Топик: `account.events`
    *   Структура Payload:
        ```json
        {
          "account_id": "uuid-account-id",
          "user_id": "uuid-auth-user_id",
          "contact_id": "uuid-contact-id",
          "type": "email",
          "value": "verified_contact@example.com",
          "verified_at": "ISO8601_timestamp"
        }
        ```
*   **`account.settings.updated.v1`**
    *   Описание: Настройки пользователя обновлены.
    *   Топик: `account.events` (или `account.settings.events`)
    *   Структура Payload:
        ```json
        {
          "account_id": "uuid-account-id",
          "user_id": "uuid-auth-user_id",
          "updated_categories": ["privacy", "notifications"], // Список обновленных категорий настроек
          "updated_at": "ISO8601_timestamp"
        }
        ```

### 5.2. Потребляемые События (Consumed Events)
*   **`auth.user.registered.v1`**
    *   Описание: От Auth Service. Инициирует создание аккаунта и профиля по умолчанию.
    *   Топик: `auth.events` (или аналогичный)
    *   Структура Payload (ожидаемая):
        ```json
        {
          "user_id": "uuid-auth-user-id",
          "email": "user@example.com", // Опционально, если email используется как username
          "username": "User123", // Опционально
          "registration_timestamp": "ISO8601_timestamp"
        }
        ```
    *   Логика обработки: Создать новую запись `Account` и связанную `Profile` (с дефолтными значениями). Опубликовать событие `account.created.v1`.
*   **`auth.user.email.updated.v1`** (Если email управляется в Auth Service и может меняться)
    *   Описание: От Auth Service. Обновление основного email пользователя.
    *   Топик: `auth.events`
    *   Структура Payload (ожидаемая):
        ```json
        {
          "user_id": "uuid-auth-user-id",
          "new_email": "new_primary@example.com",
          "old_email": "old_primary@example.com"
        }
        ```
    *   Логика обработки: Обновить соответствующий `ContactInfo` типа `email`, пометить его как основной и верифицированный.
*   **`admin.user.action.v1`** (Общее событие от Admin Service)
    *   Описание: От Admin Service, например, для принудительного обновления профиля или блокировки.
    *   Топик: `admin.events`
    *   Структура Payload: [Подлежит уточнению в зависимости от действий Admin Service]
    *   Логика обработки: Выполнить соответствующее действие над аккаунтом/профилем.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Auth Service:**
    *   Тип интеграции: Потребление событий Kafka (`auth.user.registered`), gRPC вызовы (опционально, для получения деталей пользователя по UserID, если они не приходят в событии).
    *   Назначение: Получение информации о зарегистрированных пользователях для создания аккаунтов. Account Service предоставляет UserID, полученный от Auth Service, во всех своих событиях и API для связи.
    *   Контракт: Событие `auth.user.registered.v1`.
*   **Notification Service:**
    *   Тип интеграции: Публикация событий Kafka (например, `account.contact.added` для отправки верификационного письма/SMS) или прямые gRPC вызовы к Notification Service.
    *   Назначение: Отправка уведомлений пользователю (например, о необходимости верификации email/телефона, об изменении статуса аккаунта).
    *   Контракт: События, публикуемые Account Service (см. 5.1), или gRPC методы Notification Service (например, `SendVerificationEmail`).
*   **Social Service:**
    *   Тип интеграции: Social Service потребляет события `account.profile.updated.v1` от Account Service. Account Service может вызывать gRPC Social Service для получения расширенных социальных данных профиля, если это необходимо для каких-то агрегированных представлений (маловероятно, обычно Social Service сам обогащает данные).
    *   Назначение: Обмен информацией о профиле пользователя.
    *   Контракт: Событие `account.profile.updated.v1`.
*   **Payment Service:**
    *   Тип интеграции: Payment Service может вызывать gRPC Account Service для получения базовой информации о пользователе (например, страна для определения региональных цен или налогов, если это не решается на уровне Auth или Catalog).
    *   Назначение: Предоставление информации о пользователе для финансовых операций.
    *   Контракт: gRPC `GetAccountInfo` или `GetUserProfile`.
*   **Admin Service:**
    *   Тип интеграции: Admin Service вызывает REST/gRPC эндпоинты Account Service для управления аккаунтами и профилями. Account Service может публиковать события об административных действиях для аудита.
    *   Назначение: Администрирование пользователей.
    *   Контракт: Административные REST/gRPC эндпоинты Account Service.
*   **API Gateway:**
    *   Тип интеграции: Проксирование REST запросов.
    *   Назначение: Единая точка входа для клиентских приложений.
    *   Контракт: Публичные REST API эндпоинты Account Service.

### 6.2. Внешние Системы
*   Не предполагается прямых интеграций с внешними системами, кроме тех, что опосредованы другими микросервисами (например, S3 для аватаров через API Gateway или клиентское приложение).

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `ACCOUNT_HTTP_PORT`: (например, `8081`)
*   `ACCOUNT_GRPC_PORT`: (например, `9091`)
*   `POSTGRES_DSN`: (например, `postgres://user:pass@host:port/dbname?sslmode=disable`)
*   `REDIS_ADDR`: (например, `redis-host:6379`)
*   `REDIS_PASSWORD`: (если есть)
*   `REDIS_DB_ACCOUNT`: (например, `0`)
*   `KAFKA_BROKERS`: (например, `kafka1:9092,kafka2:9092`)
*   `KAFKA_TOPIC_ACCOUNT_EVENTS`: (например, `account.events`)
*   `KAFKA_CONSUMER_GROUP_ACCOUNT`: (например, `account-service-group`)
*   `JWT_PUBLIC_KEY_PATH`: (Путь к публичному ключу для валидации токенов от API Gateway, если валидация дублируется. Обычно не нужно, если API Gateway полностью берет на себя проверку JWT).
*   `LOG_LEVEL`: (например, `info`, `debug`)
*   `AVATAR_MAX_SIZE_MB`: (например, `5`)
*   `AVATAR_S3_BUCKET`: (Если аватары хранятся в S3 и Account Service управляет загрузкой)
*   `AVATAR_S3_REGION`
*   `AVATAR_S3_ENDPOINT`
*   `AVATAR_S3_ACCESS_KEY_ID`
*   `AVATAR_S3_SECRET_ACCESS_KEY`
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`: (например, `http://jaeger-collector:14268/api/traces`)

### 7.2. Файлы Конфигурации (если применимо)
*   Расположение: `configs/config.yaml` (стандартное для проекта).
*   Структура:
    ```yaml
    http_server:
      port: ${ACCOUNT_HTTP_PORT}
    grpc_server:
      port: ${ACCOUNT_GRPC_PORT}
    postgres:
      dsn: ${POSTGRES_DSN}
    redis:
      address: ${REDIS_ADDR}
      password: ${REDIS_PASSWORD}
      db: ${REDIS_DB_ACCOUNT}
    kafka:
      brokers: ${KAFKA_BROKERS}
      topics:
        account_events: ${KAFKA_TOPIC_ACCOUNT_EVENTS}
      consumer_groups:
         account_group: ${KAFKA_CONSUMER_GROUP_ACCOUNT}
    logging:
      level: ${LOG_LEVEL}
    # ... другие настройки
    ```

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   Следование `project_api_standards.md` для форматов ошибок REST и gRPC.
*   Подробное логирование ошибок с `trace_id`.

### 8.2. Распространенные Коды Ошибок (в дополнение к стандартным)
*   **`ACCOUNT_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`)
*   **`PROFILE_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`)
*   **`CONTACT_INFO_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`)
*   **`SETTINGS_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`)
*   **`NICKNAME_TAKEN`** (HTTP 409, gRPC `ALREADY_EXISTS`)
*   **`CUSTOM_URL_TAKEN`** (HTTP 409, gRPC `ALREADY_EXISTS`)
*   **`EMAIL_ALREADY_EXISTS_FOR_USER`** (HTTP 409, gRPC `ALREADY_EXISTS`)
*   **`PHONE_ALREADY_EXISTS_FOR_USER`** (HTTP 409, gRPC `ALREADY_EXISTS`)
*   **`VERIFICATION_CODE_INVALID`** (HTTP 400, gRPC `INVALID_ARGUMENT`)
*   **`VERIFICATION_CODE_EXPIRED`** (HTTP 400, gRPC `INVALID_ARGUMENT`)
*   **`CANNOT_DELETE_PRIMARY_CONTACT`** (HTTP 400, gRPC `FAILED_PRECONDITION`)
*   **`CONTACT_NOT_VERIFIED`** (HTTP 400, gRPC `FAILED_PRECONDITION`): Например, при попытке сделать не верифицированный контакт основным.

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   Все запросы к API (кроме, возможно, публичного GET profile) требуют JWT аутентификации, проверяемой API Gateway. UserID извлекается из токена.
*   Межсервисное gRPC взаимодействие защищено mTLS.

### 9.2. Авторизация
*   Проверка владения ресурсом (пользователь может редактировать только свой профиль, свои настройки и т.д.).
*   Административные эндпоинты требуют роль `admin` (проверяется на основе информации из JWT).
*   (Ссылка на `project_roles_and_permissions.md`)

### 9.3. Защита Данных
*   Соблюдение ФЗ-152 "О персональных данных".
*   Шифрование при передаче (TLS).
*   Хеширование кодов верификации.
*   Валидация всех входных данных для предотвращения XSS, SQL Injection (через ORM).
*   Ограничение на размер загружаемых аватаров.
*   (Ссылка на `project_security_standards.md`)

### 9.4. Управление Секретами
*   Пароли к БД, Redis, ключи Kafka хранятся в Kubernetes Secrets или Vault.
*   (Ссылка на `project_security_standards.md`)

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   Dockerfile: `backend/account-service/Dockerfile` (стандартный многоэтапный для Go).
*   Helm-чарты/Kubernetes манифесты: `deploy/charts/account-service/` (предполагаемое расположение).
*   (Ссылка на `project_deployment_standards.md`)

### 10.2. Зависимости при Развертывании
*   PostgreSQL, Redis, Kafka.
*   Auth Service (для получения UserID и событий о регистрации).
*   API Gateway.

### 10.3. CI/CD
*   Стандартный пайплайн: сборка, unit-тесты, интеграционные тесты, SAST/DAST сканирование, сборка Docker-образа, деплой на окружения.
*   (Ссылка на `project_deployment_standards.md` и `.github/workflows/`)

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   Формат: JSON (Zap).
*   Ключевые события: CRUD операции, ошибки, важные изменения статуса.
*   Интеграция: ELK/Loki.
*   (Ссылка на `project_observability_standards.md`)

### 11.2. Мониторинг
*   Метрики (Prometheus):
    *   `http_requests_total{handler="/me/profile", method="GET", status="200"}`
    *   `http_request_duration_seconds{handler="/me/profile", method="GET"}`
    *   `grpc_requests_total{service="AccountService", method="GetAccountInfo", status="OK"}`
    *   `db_query_duration_seconds{query="select_profile"}`
    *   `kafka_messages_produced_total{topic="account.events"}`
    *   `cache_hits_total{cache_name="profile"}`
    *   `cache_misses_total{cache_name="profile"}`
*   Дашборды (Grafana): Обзор состояния сервиса, производительность API, использование БД/кэша.
*   Алерты (AlertManager): Высокий % ошибок, большая задержка ответов, проблемы с подключением к БД/Kafka.
*   (Ссылка на `project_observability_standards.md`)

### 11.3. Трассировка
*   Интеграция: OpenTelemetry, Jaeger.
*   Контекст трассировки передается через все запросы и события.
*   (Ссылка на `project_observability_standards.md`)

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**: API чтения (GET /me/profile) P95 < 100мс; API записи (PUT /me/profile) P95 < 200мс.
*   **Масштабируемость**: Горизонтальное масштабирование для поддержки >1 млн активных пользователей.
*   **Надежность**: Доступность > 99.95%.
*   **Сопровождаемость**: Покрытие тестами > 80%.
*   [Конкретные цифры NFR для Account Service будут определены и добавлены на последующих этапах проектирования.]

## 13. Приложения (Appendices) (Опционально)
*   [Полные примеры JSON для всех DTO и детальные схемы Protobuf будут предоставлены в соответствующих спецификациях API или в будущих обновлениях этого документа.]
*   [Полные примеры JSON для всех DTO и детальные схемы Protobuf будут предоставлены в соответствующих спецификациях API или в будущих обновлениях этого документа.]

---
*Этот документ является отправной точкой и должен регулярно обновляться по мере развития сервиса.*

## 14. Связанные Рабочие Процессы (Related Workflows)
*   [User Registration and Initial Profile Setup](../../../project_workflows/user_registration_flow.md)
