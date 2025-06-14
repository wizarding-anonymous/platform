<!-- backend\account-service\docs\README.md -->
# Спецификация Микросервиса: Account Service

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16

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
*   **REST Framework:** Echo (`github.com/labstack/echo/v4`) (согласно `../../../../project_technology_stack.md` и `../../../../PACKAGE_STANDARDIZATION.md`)
*   **gRPC Framework:** `google.golang.org/grpc` (согласно `../../../../project_technology_stack.md` и `../../../../PACKAGE_STANDARDIZATION.md`)
*   **База данных:** PostgreSQL (версия 15+) (согласно `../../../../project_technology_stack.md`)
*   **Кэширование:** Redis (версия 7.0+) (согласно `../../../../project_technology_stack.md`)
*   **Брокер сообщений:** Kafka (`github.com/confluentinc/confluent-kafka-go` или `github.com/segmentio/kafka-go`, согласно `../../../../PACKAGE_STANDARDIZATION.md`)
*   **Валидация:** `go-playground/validator/v10` (согласно `../../../../PACKAGE_STANDARDIZATION.md`)
*   **ORM/DB Driver:** GORM (`gorm.io/gorm`) с драйвером `gorm.io/driver/postgres` (согласно `../../../../project_technology_stack.md` и `../../../../PACKAGE_STANDARDIZATION.md`)
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`)
*   **Логирование:** Zap (`go.uber.org/zap`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`)
*   **Трассировка и метрики:** OpenTelemetry (`go.opentelemetry.io/otel`), Prometheus client (`github.com/prometheus/client_golang`) (согласно `../../../../PACKAGE_STANDARDIZATION.md` и `../../../../project_observability_standards.md`)
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

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
*   Ниже представлена упрощенная диаграмма слоев сервиса. Детальная диаграмма компонентов и их взаимодействий будет добавлена в будущих версиях, если потребуется большая детализация. Текущая структура соответствует стандартному подходу Clean Architecture.
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
    *   Domain Events: `AccountCreatedEvent`, `ProfileUpdatedEvent`, `ContactVerifiedEvent`.
    *   Интерфейсы репозиториев (определяются здесь, реализуются в Infrastructure Layer).

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Реализация интерфейсов, определенных в Application и Domain Layers, для взаимодействия с PostgreSQL, Redis, Kafka.
*   Ключевые компоненты/модули:
    *   PostgreSQL Repositories (GORM): Реализации `AccountRepository`, `ProfileRepository` и т.д.
    *   Redis Cache: Для кэширования профилей, настроек.
    *   Kafka Producer: Для отправки доменных событий в формате CloudEvents.
    *   Клиенты для других сервисов (например, gRPC клиент для Auth Service для получения UserID по необходимости, если эта информация не приходит с событием).

## 3. API Endpoints
Общие принципы и форматы см. `../../../../project_api_standards.md`.

### 3.1. REST API
*   **Базовый URL (через API Gateway):** `/api/v1/account` (уточненный базовый URL для ресурсов Account Service).
*   **Версионирование:** `/v1/`.
*   **Формат данных:** `application/json` (согласно `project_api_standards.md`).
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
    *   Пример ответа (Ошибка 401 Unauthorized):
        ```json
        {
          "errors": [
            {
              "code": "UNAUTHENTICATED",
              "title": "Ошибка аутентификации",
              "detail": "Необходима аутентификация для доступа к этому ресурсу."
            }
          ]
        }
        ```
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
    *   Пример ответа (Ошибка 409 Conflict - Никнейм занят):
        ```json
        {
          "errors": [
            {
              "code": "NICKNAME_TAKEN",
              "title": "Никнейм уже занят",
              "detail": "Выбранный никнейм 'NewUser123' уже используется. Пожалуйста, выберите другой.",
              "source": { "pointer": "/data/attributes/nickname" }
            }
          ]
        }
        ```
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
    *   Требуемые права доступа: Владелец аккаунта (`user`).

#### 3.1.2. Ресурс: Публичные Профили (Public Profiles)
*   **`GET /profiles/{profile_id_or_custom_url}`**
    *   Описание: Получение публичного профиля пользователя по его ID профиля или кастомному URL. Объем возвращаемых данных зависит от настроек приватности профиля.
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
    *   Требуемые права доступа: Владелец (`user`).
*   **`DELETE /me/contact-info/{contact_id}`**
    *   Описание: Удаление контактного метода. Нельзя удалить основной метод, если он единственный верифицированный.
    *   Требуемые права доступа: Владелец (`user`).
*   **`POST /me/contact-info/{contact_id}/set-primary`**
    *   Описание: Установка контактного метода как основного (если он верифицирован).
    *   Требуемые права доступа: Владелец (`user`).
*   **`POST /me/contact-info/{contact_id}/request-verification`**
    *   Описание: Повторный запрос кода верификации.
    *   Требуемые права доступа: Владелец (`user`).
*   **`POST /me/contact-info/{contact_id}/verify`**
    *   Описание: Подтверждение контактного метода с помощью кода.
    *   Тело запроса: `{"data": {"type": "verificationConfirmation", "attributes": {"code": "123456"}}}`.
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
*   **Файл .proto:** `proto/account/v1/account_service.proto` (предполагаемое расположение в репозитории `platform-protos` или локально в директории `proto` сервиса. Необходимо уточнить и стандартизировать расположение proto-файлов в рамках всего проекта.)
*   **Аутентификация:** mTLS для межсервисного взаимодействия, опционально передача UserID/токена в метаданных для служебных запросов.

#### 3.2.1. Сервис: AccountService
*   **`rpc GetAccountInfo(GetAccountInfoRequest) returns (GetAccountInfoResponse)`**
    *   Описание: Получение основной информации об аккаунте пользователя (статус, UserID).
    *   `message GetAccountInfoRequest { string user_id = 1; }`
    *   `message AccountInfo { string account_id = 1; string user_id = 2; string status = 3; google.protobuf.Timestamp created_at = 4; }`
    *   `message GetAccountInfoResponse { AccountInfo account_info = 1; }`
*   **`rpc GetUserProfile(GetUserProfileRequest) returns (GetUserProfileResponse)`**
    *   Описание: Получение профиля пользователя.
    *   `message GetUserProfileRequest { string user_id = 1; }`
    *   `message Profile { string profile_id = 1; string account_id = 2; string nickname = 3; string bio = 4; string avatar_url = 5; string country_code = 6; string custom_url_slug = 7; google.protobuf.Struct privacy_settings = 8; }`
    *   `message GetUserProfileResponse { Profile profile = 1; }`
*   **`rpc GetUserProfilesBatch(GetUserProfilesBatchRequest) returns (GetUserProfilesBatchResponse)`**
    *   Описание: Получение профилей нескольких пользователей.
    *   `message GetUserProfilesBatchRequest { repeated string user_ids = 1; }`
    *   `message GetUserProfilesBatchResponse { repeated Profile profiles = 1; }`
*   **`rpc CheckUsernameAvailability(CheckUsernameAvailabilityRequest) returns (CheckUsernameAvailabilityResponse)`**
    *   Описание: Проверка, доступен ли никнейм или кастомный URL профиля.
    *   `message CheckUsernameAvailabilityRequest { string username = 1; }` // Может быть и custom_url_slug
    *   `message CheckUsernameAvailabilityResponse { bool is_available = 1; }`
*   **`rpc GetUserSettings(GetUserSettingsRequest) returns (GetUserSettingsResponse)`**
    *   Описание: Получение настроек пользователя для использования другими сервисами (например, Notification Service для проверки предпочтений).
    *   `message GetUserSettingsRequest { string user_id = 1; repeated string categories = 2; }` // categories - опционально, для фильтрации
    *   `message GetUserSettingsResponse { google.protobuf.Struct settings_data = 1; }` // JSONB как Struct

### 3.3. WebSocket API
*   Не предполагается для основного функционала Account Service. Обновления статуса пользователя или профиля могут передаваться через события Kafka другим сервисам, которые поддерживают WebSocket (например, Social Service, Notification Service).

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **Account**:
    *   `id` (UUID, PK): Уникальный идентификатор аккаунта в Account Service. **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK, Unique): Идентификатор пользователя из Auth Service. **Обязательность: Да.**
    *   `status` (ENUM: `active`, `inactive`, `blocked`, `pending_deletion`, `deleted`): Статус аккаунта. **Обязательность: Да (DEFAULT 'inactive').**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
*   **Profile**:
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `account_id` (UUID, FK, Unique): Ссылка на Account. **Обязательность: Да.**
    *   `nickname` (VARCHAR, Nullable, Unique): Отображаемое имя. **Обязательность: Нет (но может быть установлено при регистрации).**
    *   `bio` (TEXT, Nullable). **Обязательность: Нет.**
    *   `avatar_url` (VARCHAR, Nullable). **Обязательность: Нет.**
    *   `country_code` (CHAR(2), Nullable): ISO 3166-1 alpha-2. **Обязательность: Нет.**
    *   `custom_url_slug` (VARCHAR, Nullable, Unique): Кастомный URL для профиля. **Обязательность: Нет.**
    *   `privacy_settings` (JSONB): Настройки приватности профиля (например, `{"inventory_visibility": "public", "friends_list_visibility": "private"}`). **Обязательность: Да (DEFAULT '{}').**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
*   **ContactInfo**:
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `account_id` (UUID, FK): Ссылка на Account. **Обязательность: Да.**
    *   `type` (ENUM: `email`, `phone`). **Обязательность: Да.**
    *   `value` (VARCHAR, Not Null). **Обязательность: Да.**
    *   `is_verified` (BOOLEAN, Default: false). **Обязательность: Да (DEFAULT false).**
    *   `is_primary` (BOOLEAN, Default: false). **Обязательность: Да (DEFAULT false).**
    *   `verification_code_hash` (VARCHAR, Nullable): Хэш кода верификации. **Обязательность: Нет.**
    *   `verification_code_expires_at` (TIMESTAMPTZ, Nullable). **Обязательность: Нет.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   (Constraint: Unique(account_id, type, value), Unique(account_id, type, is_primary) WHERE is_primary = TRUE)
*   **UserSetting**:
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `account_id` (UUID, FK, Unique): Ссылка на Account. **Обязательность: Да.**
    *   `settings_data` (JSONB): Все настройки пользователя, сгруппированные по категориям (например, `{"notifications": {"email_updates": true}, "interface": {"language": "ru"}}`). **Обязательность: Да (DEFAULT '{}').**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**

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
            %% Unique(account_id, type, value)
            %% Unique(account_id, type, is_primary) WHERE is_primary = TRUE
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
        status VARCHAR(50) NOT NULL DEFAULT 'inactive', -- e.g., inactive, active, blocked, pending_deletion, deleted
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
    CREATE INDEX idx_profiles_nickname ON profiles(nickname text_pattern_ops); -- для LIKE запросов, если нужны
    CREATE INDEX idx_profiles_custom_url_slug ON profiles(custom_url_slug text_pattern_ops); -- для LIKE запросов

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
    CREATE INDEX idx_contact_infos_value ON contact_infos(value); -- для поиска по email/телефону

    -- Таблица: user_settings
    CREATE TABLE user_settings (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        account_id UUID NOT NULL UNIQUE REFERENCES accounts(id) ON DELETE CASCADE,
        settings_data JSONB DEFAULT '{}'::jsonb, -- {"privacy": {...}, "notifications": {...}}
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
    CREATE INDEX idx_user_settings_account_id ON user_settings(account_id);
    ```

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Kafka (согласно `../../../../project_technology_stack.md`).
*   **Формат событий:** CloudEvents JSON (согласно `../../../../project_api_standards.md`).
*   **Основные топики Kafka:** `ru.steam.account.events.v1` (или более гранулярные, если потребуется).

*   **`ru.steam.account.created.v1`**
    *   Описание: Публикуется после успешного создания базовой записи аккаунта (обычно после события `ru.steam.auth.user.registered.v1` от Auth Service).
    *   Структура Payload (`data` в CloudEvent):
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user-id",
          "status": "inactive", // Начальный статус
          "createdAt": "ISO8601_timestamp"
        }
        ```
*   **`ru.steam.account.status.updated.v1`**
    *   Описание: Статус аккаунта изменен (например, активирован, заблокирован).
    *   Структура Payload:
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user-id",
          "oldStatus": "active",
          "newStatus": "blocked",
          "reason": "Violation of terms", // Опционально
          "updatedAt": "ISO8601_timestamp"
        }
        ```
*   **`ru.steam.account.profile.updated.v1`**
    *   Описание: Профиль пользователя обновлен.
    *   Структура Payload:
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user_id",
          "profileId": "uuid-profile-id",
          "updatedFields": ["nickname", "avatarUrl"], // Список измененных полей
          "updatedAt": "ISO8601_timestamp"
        }
        ```
*   **`ru.steam.account.contact.added.v1`**
    *   Описание: Добавлен новый контактный метод.
    *   Структура Payload:
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user_id",
          "contactId": "uuid-contact-id",
          "type": "email", // "email" or "phone"
          "value": "new_contact@example.com", // Немаскированное значение для внутренних обработчиков, если необходимо (например, Notification Service)
          "addedAt": "ISO8601_timestamp"
        }
        ```
*   **`ru.steam.account.contact.verified.v1`**
    *   Описание: Контактный метод успешно верифицирован.
    *   Структура Payload:
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user_id",
          "contactId": "uuid-contact-id",
          "type": "email",
          "value": "verified_contact@example.com", // Немаскированное значение
          "verifiedAt": "ISO8601_timestamp"
        }
        ```
*   **`ru.steam.account.settings.updated.v1`**
    *   Описание: Настройки пользователя обновлены.
    *   Структура Payload:
        ```json
        {
          "accountId": "uuid-account-id",
          "userId": "uuid-auth-user_id",
          "updatedCategories": ["privacy", "notifications"], // Список обновленных категорий настроек
          "updatedAt": "ISO8601_timestamp"
        }
        ```

### 5.2. Потребляемые События (Consumed Events)
*   **`ru.steam.auth.user.registered.v1`**
    *   Описание: От Auth Service. Инициирует создание аккаунта и профиля по умолчанию.
    *   Топик: `ru.steam.auth.events.v1` (или аналогичный, согласно спецификации Auth Service)
    *   Ожидаемая структура Payload (`data` в CloudEvent):
        ```json
        {
          "userId": "uuid-auth-user-id",
          "email": "user@example.com", // Опционально, если email используется как username или начальный контакт
          "username": "User123", // Опционально, может быть использован для начального nickname
          "registrationTimestamp": "ISO8601_timestamp"
        }
        ```
    *   Логика обработки: Создать новую запись `Account` (статус `inactive` или `pending_email_verification`), связанный `Profile` (с дефолтными значениями или из данных события), и `ContactInfo` если email предоставлен. Опубликовать событие `ru.steam.account.created.v1`.
*   **`ru.steam.auth.user.email.updated.v1`** (Если основной email управляется в Auth Service и может меняться, и это изменение должно отражаться в Account Service как основной верифицированный email)
    *   Описание: От Auth Service. Обновление основного email пользователя.
    *   Топик: `ru.steam.auth.events.v1`
    *   Ожидаемая структура Payload:
        ```json
        {
          "userId": "uuid-auth-user-id",
          "newEmail": "new_primary@example.com",
          "isVerified": true // Предполагается, что Auth Service уже верифицировал email
        }
        ```
    *   Логика обработки: Найти или создать `ContactInfo` для `newEmail`, пометить его как основной и верифицированный. Если старый email был основным, он перестает им быть.
*   **`ru.steam.admin.user.action.v1`** (Пример общего события от Admin Service)
    *   Описание: От Admin Service, например, для принудительной блокировки аккаунта или изменения профиля.
    *   Топик: `ru.steam.admin.events.v1`
    *   Структура Payload: Зависит от конкретного действия, но должна включать `userId` и детали действия.
        ```json
        {
          "userId": "uuid-auth-user-id",
          "actionType": "BLOCK_ACCOUNT", // "UPDATE_PROFILE", "DELETE_ACCOUNT"
          "payload": { // Специфичные данные для действия
            "reason": "Terms of Service violation" // для BLOCK_ACCOUNT
          }
        }
        ```
    *   Логика обработки: Выполнить соответствующее административное действие над аккаунтом/профилем.

## 6. Интеграции (Integrations)
См. `../../../../project_integrations.md` для общей карты и деталей.

### 6.1. Внутренние Микросервисы
*   **Auth Service:**
    *   Тип интеграции: Потребление событий Kafka (`ru.steam.auth.user.registered.v1`), gRPC вызовы (опционально, для получения деталей пользователя по UserID, если они не приходят в событии и требуются синхронно).
    *   Назначение: Получение информации о зарегистрированных пользователях для создания аккаунтов. Account Service предоставляет UserID, полученный от Auth Service, во всех своих событиях и API для связи.
*   **Notification Service:**
    *   Тип интеграции: Account Service инициирует отправку уведомлений (например, верификационные письма/SMS, изменение статуса аккаунта) через Notification Service. Это может быть реализовано через публикацию специфичных событий в Kafka, на которые подписан Notification Service (например, `ru.steam.account.contact.verification.requested.v1`), либо через прямые gRPC вызовы к Notification Service.
    *   Назначение: Отправка пользовательских уведомлений.
*   **Social Service:**
    *   Тип интеграции: Social Service потребляет события `ru.steam.account.profile.updated.v1` от Account Service для обновления своей локальной копии профильных данных. Account Service может вызывать gRPC Social Service для проверки уникальности никнейма, если это не делается локально или если никнейм глобален для платформы.
    *   Назначение: Обмен информацией о профиле пользователя.
*   **Payment Service:**
    *   Тип интеграции: Payment Service может вызывать gRPC Account Service (`GetAccountInfo`, `GetUserProfile`) для получения базовой информации о пользователе (например, страна для определения региональных цен или налогов).
    *   Назначение: Предоставление информации о пользователе для финансовых операций.
*   **Admin Service:**
    *   Тип интеграции: Admin Service вызывает REST/gRPC эндпоинты Account Service для управления аккаунтами и профилями (см. раздел 3.1.5). Account Service может публиковать события об административных действиях для аудита.
    *   Назначение: Администрирование пользователей.
*   **API Gateway:**
    *   Тип интеграции: Проксирование REST запросов к Account Service.
    *   Назначение: Единая точка входа для клиентских приложений.

### 6.2. Внешние Системы
*   **S3-совместимое хранилище (Yandex Object Storage, VK Cloud, etc.):**
    *   Тип интеграции: Если аватары загружаются через Account Service, он будет взаимодействовать с S3 для сохранения файлов. Чаще загрузка идет через клиентское приложение или API Gateway напрямую в S3, а Account Service только сохраняет URL.
    *   Назначение: Хранение аватаров пользователей.
*   **Интеграция с системами аутентификации (VK, Telegram):**
    *   Account Service напрямую не интегрируется с VK/Telegram для аутентификации. Этим занимается Auth Service. Account Service потребляет событие `ru.steam.auth.user.registered.v1`, которое может быть результатом такой внешней аутентификации.

## 7. Конфигурация (Configuration)
Общие стандарты конфигурации см. в `../../../../project_api_standards.md` (раздел 7) и `../../../../DOCUMENTATION_GUIDELINES.md` (раздел 6).

### 7.1. Переменные Окружения (Примеры)
*   `ACCOUNT_HTTP_PORT`: Порт для REST API (например, `8081`)
*   `ACCOUNT_GRPC_PORT`: Порт для gRPC API (например, `9091`)
*   `POSTGRES_DSN`: Строка подключения к PostgreSQL (например, `postgres://user:pass@host:port/dbname?sslmode=disable`)
*   `REDIS_ADDR`: Адрес Redis (например, `redis-host:6379`)
*   `REDIS_PASSWORD`: Пароль для Redis (если используется)
*   `REDIS_DB_ACCOUNT`: Номер базы данных Redis для Account Service (например, `0`)
*   `KAFKA_BROKERS`: Список брокеров Kafka (например, `kafka1:9092,kafka2:9092`)
*   `KAFKA_TOPIC_ACCOUNT_EVENTS`: Имя топика для публикуемых событий (например, `ru.steam.account.events.v1`)
*   `KAFKA_CONSUMER_GROUP_ACCOUNT`: Имя группы потребителей Kafka (например, `account-service-consumer-group`)
*   `JWT_PUBLIC_KEY_PATH`: Путь к публичному ключу для валидации JWT токенов (если API Gateway не полностью берет на себя эту задачу, или для служебных нужд). Обычно не требуется, если валидация на API Gateway.
*   `LOG_LEVEL`: Уровень логирования (например, `info`, `debug`, `error`)
*   `AVATAR_MAX_SIZE_BYTES`: Максимальный размер аватара в байтах (например, `5242880` для 5MB).
*   `AVATAR_S3_BUCKET`: S3 бакет для аватаров (если Account Service управляет загрузкой).
*   `AVATAR_S3_REGION`: Регион S3 бакета.
*   `AVATAR_S3_ENDPOINT`: Endpoint S3-совместимого хранилища.
*   `AVATAR_S3_ACCESS_KEY_ID`: Ключ доступа S3.
*   `AVATAR_S3_SECRET_ACCESS_KEY`: Секретный ключ S3.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`: Endpoint для экспорта трейсов в Jaeger (например, `http://jaeger-collector:14268/api/traces`).
*   `DEFAULT_USER_SETTINGS_JSON`: Строка JSON с настройками по умолчанию для новых пользователей.

### 7.2. Файлы Конфигурации (`configs/config.yaml`)
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
      password: ${REDIS_PASSWORD:""} # Пустое значение по умолчанию, если пароль не указан
      db: ${REDIS_DB_ACCOUNT:0}
    kafka:
      brokers: ${KAFKA_BROKERS}
      producer_topics:
        account_events: ${KAFKA_TOPIC_ACCOUNT_EVENTS:"ru.steam.account.events.v1"}
      consumer_topics:
        auth_events: ${KAFKA_TOPIC_AUTH_EVENTS:"ru.steam.auth.events.v1"} # Пример потребляемого топика
      consumer_group: ${KAFKA_CONSUMER_GROUP_ACCOUNT:"account-service-consumer-group"}
    logging:
      level: ${LOG_LEVEL:"info"} # debug, info, warn, error, fatal
      format: "json" # text, json
    security:
      jwt_public_key_path: ${JWT_PUBLIC_KEY_PATH:""} # Если используется
      # ФЗ-152 специфичные настройки, если есть (например, пути к логам доступа к ПДн)
    avatar:
      max_size_bytes: ${AVATAR_MAX_SIZE_BYTES:5242880} # 5MB
      # S3 settings if service handles uploads directly
      # s3:
      #   bucket: ${AVATAR_S3_BUCKET:""}
      #   region: ${AVATAR_S3_REGION:""}
      #   endpoint: ${AVATAR_S3_ENDPOINT:""}
      #   access_key_id: ${AVATAR_S3_ACCESS_KEY_ID:""}
      #   secret_access_key: ${AVATAR_S3_SECRET_ACCESS_KEY:""}
    otel: # OpenTelemetry
      exporter_jaeger_endpoint: ${OTEL_EXPORTER_JAEGER_ENDPOINT:""}
      service_name: "account-service"
    default_settings:
      user_settings_json: ${DEFAULT_USER_SETTINGS_JSON:'{"privacy": {"inventoryVisibility": "private"}, "notifications": {"newFriendRequestEmail": true}}'}

    ```

## 8. Обработка Ошибок (Error Handling)
*   Следование `../../../../project_api_standards.md` для форматов ошибок REST и gRPC.
*   Подробное логирование ошибок с `trace_id` и контекстом.

### 8.1. Распространенные Коды Ошибок (в дополнение к стандартным HTTP/gRPC)
*   **`ACCOUNT_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`): Аккаунт не найден.
*   **`PROFILE_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`): Профиль не найден.
*   **`CONTACT_INFO_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`): Контактная информация не найдена.
*   **`SETTINGS_NOT_FOUND`** (HTTP 404, gRPC `NOT_FOUND`): Настройки не найдены (редко, обычно создаются по умолчанию).
*   **`NICKNAME_TAKEN`** (HTTP 409, gRPC `ALREADY_EXISTS`): Никнейм уже занят.
*   **`CUSTOM_URL_TAKEN`** (HTTP 409, gRPC `ALREADY_EXISTS`): Кастомный URL уже занят.
*   **`EMAIL_ALREADY_EXISTS_FOR_USER`** (HTTP 409, gRPC `ALREADY_EXISTS`): Email уже используется этим пользователем.
*   **`PHONE_ALREADY_EXISTS_FOR_USER`** (HTTP 409, gRPC `ALREADY_EXISTS`): Телефон уже используется этим пользователем.
*   **`VERIFICATION_CODE_INVALID`** (HTTP 400, gRPC `INVALID_ARGUMENT`): Неверный код верификации.
*   **`VERIFICATION_CODE_EXPIRED`** (HTTP 400, gRPC `INVALID_ARGUMENT`): Срок действия кода верификации истек.
*   **`CANNOT_DELETE_PRIMARY_CONTACT`** (HTTP 400, gRPC `FAILED_PRECONDITION`): Нельзя удалить основной контакт, если он единственный верифицированный.
*   **`CONTACT_NOT_VERIFIED`** (HTTP 400, gRPC `FAILED_PRECONDITION`): Например, при попытке сделать не верифицированный контакт основным.
*   **`INVALID_AVATAR_FILE`** (HTTP 400, gRPC `INVALID_ARGUMENT`): Некорректный файл аватара (формат, размер).
*   **`UPDATE_CONFLICT`** (HTTP 409, gRPC `ABORTED`): Конфликт одновременного обновления, требуется повторить операцию (с использованием ETag или версионирования).

## 9. Безопасность (Security)
См. `../../../../project_security_standards.md` для общих стандартов.

### 9.1. Аутентификация
*   REST API: Все эндпоинты (кроме, возможно, публичного GET `/profiles/{id_or_slug}`) требуют JWT аутентификации. Проверка JWT и извлечение UserID/ролей выполняется API Gateway. Account Service доверяет этим данным.
*   gRPC API: Межсервисное взаимодействие защищено mTLS. Для запросов, инициируемых пользователем и проходящих через другие сервисы, UserID передается в метаданных.

### 9.2. Авторизация
*   Проверка владения ресурсом: Пользователь может редактировать только свой профиль, свои настройки и т.д. (на основе `X-User-ID` из API Gateway).
*   Административные эндпоинты (префикс `/admin/`) требуют роль `admin` (на основе `X-User-Roles` из API Gateway).
*   Детализация ролей и разрешений: см. `../../../../project_roles_and_permissions.md`.

### 9.3. Защита Данных
*   **ФЗ-152 "О персональных данных":**
    *   Account Service обрабатывает значительный объем персональных данных российских граждан (никнейм, ФИО если указано, email, телефон, страна, IP-адреса из логов и т.д.).
    *   Все персональные данные российских граждан, собираемые и обрабатываемые Account Service, хранятся и обрабатываются на серверах, физически расположенных на территории Российской Федерации. Платформа использует инфраструктуру российского хостинг-провайдера Beget для размещения своих сервисов и баз данных (PostgreSQL, Redis), содержащих эти персональные данные.
    *   Согласие на обработку ПДн получается при регистрации (управляется Auth Service и клиентским приложением) и при предоставлении специфических данных в профиле.
    *   Логирование доступа к ПДн (кто, когда, какие данные запрашивал/изменял) обязательно.
    *   Шифрование при передаче (TLS).
    *   Хеширование кодов верификации (например, SHA256 + соль). Пароли пользователей не хранятся в Account Service.
    *   Анонимизация или удаление данных по запросу пользователя (в рамках политик платформы).
*   Валидация всех входных данных для предотвращения XSS, SQL Injection (через ORM и валидаторы).
*   Ограничение на размер загружаемых аватаров и их содержимого (проверка на стороне сервиса или через внешний сервис).

### 9.4. Управление Секретами
*   Пароли к БД, Redis, ключи для Kafka и другие секреты сервиса хранятся в Kubernetes Secrets или HashiCorp Vault (согласно `../../../../project_security_standards.md`).
*   Доступ к секретам осуществляется через переменные окружения, внедряемые Kubernetes, или через API Vault.

## 10. Развертывание (Deployment)
См. `../../../../project_deployment_standards.md` для общих стандартов.

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** `backend/account-service/Dockerfile` (стандартный многоэтапный для Go, см. `project_deployment_standards.md`).
*   **Helm-чарты/Kubernetes манифесты:** `deploy/charts/account-service/` (предполагаемое расположение, согласно `project_deployment_standards.md`).
    *   Включает Deployment, Service, ConfigMap, Secret, HorizontalPodAutoscaler.

### 10.2. Зависимости при Развертывании
*   PostgreSQL (доступен по DSN).
*   Redis (доступен по адресу).
*   Kafka (доступен список брокеров).
*   Auth Service (для координации при регистрации, хотя основное взаимодействие через события).
*   API Gateway (для проксирования REST API).
*   (Опционально) Jaeger/OpenTelemetry Collector (для экспорта трейсов).

### 10.3. CI/CD
*   Стандартный пайплайн CI/CD (определен в `project_deployment_standards.md`):
    1.  Сборка бинарного файла.
    2.  Unit-тесты.
    3.  Интеграционные тесты (с PostgreSQL, Redis).
    4.  SAST/DAST сканирование (например, `gosec`).
    5.  Сборка Docker-образа.
    6.  Публикация Docker-образа в приватный реестр.
    7.  Деплой на окружения (dev, staging, prod) с использованием Helm.
*   Автоматическое применение миграций БД при деплое (через init-контейнер или хуки Helm).
*   Файлы CI/CD: `.github/workflows/account-service.yml` или аналогичный для GitLab CI.

## 11. Мониторинг и Логирование (Logging and Monitoring)
См. `../../../../project_observability_standards.md` для общих стандартов.

### 11.1. Логирование
*   **Формат:** JSON (с использованием Zap).
*   **Уровни:** DEBUG, INFO, WARN, ERROR, FATAL.
*   **Ключевые поля:** `timestamp`, `level`, `service`, `instance`, `message`, `trace_id`, `span_id`, `user_id` (если применимо), `error_details` (для ошибок).
*   **Интеграция:** Fluent Bit для сбора логов и отправки в ELK/Loki.

### 11.2. Мониторинг
*   **Метрики (Prometheus):**
    *   **HTTP:** `http_requests_total{handler="/me/profile", method="GET", status="200"}`, `http_request_duration_seconds{handler="/me/profile", method="GET"}`.
    *   **gRPC:** `grpc_requests_total{service="AccountService", method="GetAccountInfo", status_code="OK"}`, `grpc_request_duration_seconds{...}`.
    *   **База данных (GORM/pgx):** `db_query_duration_seconds{query_type="select_profile"}`, `db_errors_total`.
    *   **Kafka:** `kafka_messages_produced_total{topic="ru.steam.account.events.v1"}`, `kafka_producer_errors_total`.
    *   **Кэш (Redis):** `cache_hits_total{cache_name="profile"}`, `cache_misses_total{cache_name="profile"}`.
    *   **Бизнес-метрики:** `accounts_created_total`, `profiles_updated_total`, `active_accounts_gauge`.
*   **Дашборды (Grafana):** Стандартный дашборд для Account Service, включающий обзор состояния, производительность API, использование БД/кэша, статистику по событиям Kafka.
*   **Алерты (AlertManager):**
    *   Высокий % ошибок REST/gRPC API.
    *   Большая задержка ответов API.
    *   Проблемы с подключением к PostgreSQL/Redis/Kafka.
    *   Ошибки публикации/потребления сообщений Kafka.
    *   Низкий уровень успешных верификаций контактов (аномалия).

### 11.3. Трассировка
*   **Интеграция:** OpenTelemetry SDK для Go.
*   **Экспорт:** Jaeger или другой совместимый коллектор.
*   **Контекст трассировки:** Передается через все HTTP/gRPC запросы и Kafka события (`trace_id`, `span_id`).
*   Автоматическая инструментация для HTTP, gRPC, SQL запросов. Ручная для специфичных бизнес-операций.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**:
    *   API чтения профиля (`GET /me`): P95 < 100мс.
    *   API обновления профиля (`PUT /me/profile`): P95 < 200мс.
    *   gRPC `GetAccountInfo`: P99 < 50мс.
*   **Масштабируемость**: Горизонтальное масштабирование для поддержки >10 млн активных пользователей (требуется уточнение профиля нагрузки). Способность обрабатывать до 1000 RPS на чтение и 200 RPS на запись профильных данных.
*   **Надежность**: Доступность > 99.95%. Отсутствие единой точки отказа (кроме БД, которая должна иметь свою стратегию HA).
*   **Сопровождаемость**: Покрытие кода unit-тестами > 80%. Покрытие интеграционными тестами ключевых сценариев > 70%.
*   **Безопасность**: Соответствие `project_security_standards.md` и ФЗ-152.
*   [Конкретные цифры NFR для Account Service могут быть доработаны на последующих этапах проектирования и нагрузочного тестирования.]

## 13. Приложения (Appendices)
*   Детальные схемы Protobuf для gRPC API находятся в репозитории `platform-protos` (или локально `proto/account/v1/`).
*   Примеры полных JSON для всех REST DTO могут быть добавлены при необходимости или генерироваться из OpenAPI спецификации.

## 14. Пользовательские Сценарии (User Flows)

В этом разделе описаны ключевые пользовательские сценарии, в которых участвует Account Service. Эти сценарии демонстрируют, как сервис взаимодействует с пользователями и другими микросервисами платформы для выполнения своих основных функций.

### 14.1. Регистрация Пользователя и Создание Аккаунта

Этот сценарий описывает процесс создания аккаунта пользователя после его первоначальной регистрации через Auth Service.

*   **Описание:** После того как пользователь успешно проходит регистрацию в Auth Service (например, предоставляет email, пароль, и возможно, проходит первичную верификацию email), Auth Service публикует событие `auth.user.registered.v1`. Account Service потребляет это событие, создает соответствующую запись `Account`, базовый `Profile`, `ContactInfo` (если email предоставлен и должен быть здесь сохранен) и `UserSetting`. Account Service затем публикует событие `account.created.v1`.
*   **Связанный документ:** Детальный воркфлоу регистрации описан в [Регистрация пользователя и начальная настройка профиля](../../../../project_workflows/user_registration_flow.md).
*   **Диаграмма (роль Account Service):**
    ```mermaid
    sequenceDiagram
        participant AuthSvc as Auth Service
        participant KafkaBus as Kafka Message Bus
        participant AccountSvc as Account Service

        AuthSvc->>KafkaBus: Publish `auth.user.registered.v1` (userId, email, username)
        KafkaBus->>AccountSvc: Consume `auth.user.registered.v1`
        AccountSvc->>AccountSvc: Create Account (userId, status='inactive')
        AccountSvc->>AccountSvc: Create Profile (nickname=username)
        AccountSvc->>AccountSvc: Create ContactInfo (email, is_verified=false)
        AccountSvc->>AccountSvc: Create UserSetting (default settings)
        AccountSvc->>KafkaBus: Publish `account.created.v1` (accountId, userId)
        AccountSvc->>KafkaBus: Publish `account.contact.added.v1` (accountId, userId, contactId, email)
        opt Email Verification by Auth Service
           AuthSvc->>KafkaBus: Publish `auth.user.email_verified.v1` (userId, email)
           KafkaBus->>AccountSvc: Consume `auth.user.email_verified.v1`
           AccountSvc->>AccountSvc: Update ContactInfo: email.is_verified = true
           AccountSvc->>AccountSvc: Update Account: status = 'active'
           AccountSvc->>KafkaBus: Publish `account.contact.verified.v1`
           AccountSvc->>KafkaBus: Publish `account.status.updated.v1`
        end
    ```

### 14.2. Обновление Профиля Пользователя

Этот сценарий описывает, как пользователь обновляет информацию своего профиля.

*   **Описание:** Аутентифицированный пользователь через клиентское приложение отправляет запрос на обновление своего профиля (например, никнейм, биография, аватар, страна). Запрос поступает в Account Service через API Gateway. Account Service валидирует данные, обновляет сущность `Profile` в своей базе данных и публикует событие `account.profile.updated.v1`.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Client Application
        participant APIGW as API Gateway
        participant AccountSvc as Account Service
        participant KafkaBus as Kafka Message Bus

        User->>ClientApp: Запрос на обновление профиля (новые данные)
        ClientApp->>APIGW: PUT /api/v1/account/me/profile (данные профиля)
        APIGW->>AccountSvc: Forward PUT /me/profile (X-User-ID, данные профиля)
        AccountSvc->>AccountSvc: Валидация данных (например, уникальность никнейма)
        AccountSvc->>AccountSvc: Обновление записи Profile в БД
        opt Загрузка аватара
            ClientApp->>APIGW: POST /api/v1/account/me/profile/avatar (файл)
            APIGW->>AccountSvc: Forward POST /me/profile/avatar
            AccountSvc->>AccountSvc: Сохранение аватара (например, в S3) и обновление avatar_url
        end
        AccountSvc->>KafkaBus: Publish `account.profile.updated.v1` (accountId, userId, updatedFields)
        AccountSvc-->>APIGW: HTTP 200 OK (обновленный профиль)
        APIGW-->>ClientApp: HTTP 200 OK
        ClientApp-->>User: Отображение обновленного профиля
    ```

### 14.3. Обновление и Верификация Контактной Информации

Этот сценарий описывает процесс добавления, обновления или верификации контактной информации пользователя (email, телефон).

*   **Описание:** Пользователь добавляет новый email или телефон. Account Service сохраняет эту информацию и инициирует процесс верификации. Для этого он может отправить запрос в Notification Service (напрямую или через Kafka), чтобы тот отправил код верификации пользователю. После получения кода пользователь вводит его, и Account Service проверяет код, обновляя статус контактной информации.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Client Application
        participant APIGW as API Gateway
        participant AccountSvc as Account Service
        participant NotificationSvc as Notification Service
        participant KafkaBus as Kafka Message Bus

        User->>ClientApp: Добавить email "new@example.com"
        ClientApp->>APIGW: POST /api/v1/account/me/contact-info (type: email, value: "new@example.com")
        APIGW->>AccountSvc: Forward POST /me/contact-info (X-User-ID, contact data)
        AccountSvc->>AccountSvc: Сохранить ContactInfo (is_verified=false, verification_code=generate())
        AccountSvc->>KafkaBus: Publish `account.contact.added.v1`
        AccountSvc->>NotificationSvc: (gRPC or Kafka) Request to send verification code (email, code)
        NotificationSvc->>User: Отправка email/SMS с кодом верификации
        AccountSvc-->>APIGW: HTTP 201 Created
        APIGW-->>ClientApp: HTTP 201 Created

        User->>ClientApp: Ввод кода верификации "123456"
        ClientApp->>APIGW: POST /api/v1/account/me/contact-info/{contact_id}/verify (code: "123456")
        APIGW->>AccountSvc: Forward POST /me/contact-info/{contact_id}/verify (X-User-ID, code)
        AccountSvc->>AccountSvc: Проверка кода и срока его действия
        alt Код верный
            AccountSvc->>AccountSvc: Обновить ContactInfo (is_verified=true)
            AccountSvc->>KafkaBus: Publish `account.contact.verified.v1`
            AccountSvc-->>APIGW: HTTP 200 OK
        else Код неверный или истек
            AccountSvc-->>APIGW: HTTP 400 Bad Request (VERIFICATION_CODE_INVALID/EXPIRED)
        end
        APIGW-->>ClientApp: Соответствующий ответ
        ClientApp-->>User: Отображение результата
    ```

### 14.4. Обновление Пользовательских Настроек

Этот сценарий описывает, как пользователь изменяет свои настройки.

*   **Описание:** Пользователь изменяет настройки приватности, уведомлений или интерфейса. Запрос на изменение настроек поступает в Account Service, который валидирует и сохраняет новые значения в сущности `UserSetting`. Публикуется событие `account.settings.updated.v1`.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Client Application
        participant APIGW as API Gateway
        participant AccountSvc as Account Service
        participant KafkaBus as Kafka Message Bus

        User->>ClientApp: Изменение настроек (например, язык интерфейса)
        ClientApp->>APIGW: PUT /api/v1/account/me/settings (новые настройки)
        APIGW->>AccountSvc: Forward PUT /me/settings (X-User-ID, настройки)
        AccountSvc->>AccountSvc: Валидация и сохранение UserSetting в БД
        AccountSvc->>KafkaBus: Publish `account.settings.updated.v1` (accountId, userId, updatedCategories)
        AccountSvc-->>APIGW: HTTP 200 OK (обновленные настройки)
        APIGW-->>ClientApp: HTTP 200 OK
        ClientApp-->>User: Отображение подтверждения
    ```

### 14.5. Запрос на Удаление Аккаунта

Этот сценарий описывает процесс, когда пользователь запрашивает удаление своего аккаунта.

*   **Описание:** Пользователь инициирует удаление своего аккаунта. Account Service получает запрос, может изменить статус аккаунта на `pending_deletion` и инициировать другие процессы (например, уведомление пользователя о периоде ожидания, запуск задач по анонимизации или удалению данных в других сервисах через события). Фактическое удаление или анонимизация может быть отложенным процессом.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Client Application
        participant APIGW as API Gateway
        participant AccountSvc as Account Service
        participant KafkaBus as Kafka Message Bus

        User->>ClientApp: Запрос на удаление аккаунта
        ClientApp->>APIGW: DELETE /api/v1/account/me
        APIGW->>AccountSvc: Forward DELETE /me (X-User-ID)
        AccountSvc->>AccountSvc: Проверка условий (например, нет активных подписок)
        AccountSvc->>AccountSvc: Обновление статуса Account на `pending_deletion` в БД
        AccountSvc->>KafkaBus: Publish `account.status.updated.v1` (accountId, userId, newStatus='pending_deletion')
        AccountSvc-->>APIGW: HTTP 202 Accepted (или HTTP 204 No Content)
        APIGW-->>ClientApp: Соответствующий ответ
        ClientApp-->>User: Отображение информации о процессе удаления

        Note over AccountSvc, KafkaBus: Дальнейшие шаги (анонимизация, удаление данных в других сервисах) могут быть инициированы событием `account.status.updated.v1` (newStatus='pending_deletion' или 'deleted') и выполняться асинхронно.
    ```

## 15. Резервное Копирование и Восстановление (Backup and Recovery)

### 14.1. PostgreSQL
*   **Процедура резервного копирования:**
    *   **Логические бэкапы:** Ежедневный `pg_dumpall` для полной копии схемы и данных. Хранение на отдельном защищенном хранилище.
    *   **Физические бэкапы (Point-in-Time Recovery - PITR):** Настроена непрерывная архивация WAL-сегментов (Write-Ahead Logging) с использованием `pg_basebackup` для создания базовой копии и `archive_command` для архивации WAL. Это позволяет восстановить состояние БД на любой момент времени.
    *   **Частота:** Базовый бэкап еженедельно, WAL-файлы архивируются непрерывно.
    *   **Хранение:** Бэкапы хранятся в S3-совместимом хранилище с шифрованием и версионированием, в другом регионе от основной БД. Срок хранения полных бэкапов - 30 дней, WAL-сегментов - 14 дней.
*   **Процедура восстановления:**
    *   Тестируется ежеквартально.
    *   Восстановление из `pg_dumpall` (для полного восстановления на новую систему).
    *   Восстановление PITR (для восстановления на конкретный момент времени).
*   **RTO (Recovery Time Objective):** < 4 часов для полного восстановления.
*   **RPO (Recovery Point Objective):** < 5 минут (максимальная потеря данных при сбое).

### 14.2. Redis
*   **Процедура резервного копирования:**
    *   **RDB Snapshots:** Автоматическое сохранение RDB-снапшотов каждые 6 часов.
    *   **AOF (Append Only File):** Включен режим AOF с fsync `everysec` для минимизации потерь данных при сбое сервера.
    *   **Хранение:** RDB-снапшоты и AOF-файлы (при необходимости) копируются в S3-совместимое хранилище ежедневно. Срок хранения - 7 дней.
*   **Процедура восстановления:**
    *   Восстановление из последнего RDB-снапшота. Если AOF включен, Redis также использует его для восстановления данных, записанных после последнего снапшота.
    *   Тестируется ежеквартально.
*   **RTO:** < 1 часа.
*   **RPO:** < 1 минуты (при fsync `everysec` для AOF). Redis используется в основном для кэширования и данных, которые могут быть пересозданы из PostgreSQL, поэтому требования к RPO могут быть менее строгими, чем для основной БД.

### 14.3. Общая стратегия
*   Резервное копирование и восстановление являются частью общей стратегии обеспечения непрерывности бизнеса и аварийного восстановления платформы.
*   Все процедуры документированы и регулярно пересматриваются.
*   Мониторинг процессов резервного копирования настроен для своевременного обнаружения сбоев.
*   Общие принципы резервного копирования для различных СУБД описаны в `../../../../project_database_structure.md`.

## 16. Связанные Рабочие Процессы (Related Workflows)
*   [Регистрация пользователя и начальная настройка профиля](../../../../project_workflows/user_registration_flow.md)

---
*Этот документ является отправной точкой и должен регулярно обновляться по мере развития сервиса.*
