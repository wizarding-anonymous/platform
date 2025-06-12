# Спецификация Микросервиса: Account Service

**Версия:** 1.0
**Дата последнего обновления:** 2025-05-24

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса "Account Service". Он описывает назначение, архитектуру, функциональные возможности, структуру данных, API, интеграции и нефункциональные требования, необходимые для реализации и эксплуатации сервиса управления Аккаунтами (Accounts) и Профилями Пользователей (User Profiles).
*   **Роль в общей архитектуре платформы:** Account Service является корневым и фундаментальным микросервисом платформы, отвечающим за управление Аккаунтами Пользователей, Профилями Пользователей и связанными с ними данными (настройки, верификации, контактная информация). Он выступает в качестве единого источника достоверной информации о Пользователях.
*   **Основные бизнес-задачи:**
    *   Регистрация и управление жизненным циклом Аккаунтов.
    *   Хранение и управление данными Профилей Пользователей.
    *   Управление контактной информацией и ее верификация.
    *   Управление настройками Пользователей.
    *   Предоставление API для доступа к данным Аккаунтов и Профилей.
    *   Генерация Событий об изменениях состояния Аккаунтов и Профилей.

### 1.2. Ключевые Функциональности
*   **Управление Аккаунтами:** Регистрация, хранение и управление базовой информацией (ID, Username, статус), обновление статуса, блокировка/разблокировка, удаление (мягкое), поиск.
*   **Управление Профилями Пользователей:** Создание, редактирование (Nickname, bio, страна, город и т.д.), управление изображениями (Avatar, баннеры), настройка видимости, история изменений.
*   **Верификация Пользователей:** Управление контактной информацией (email, телефон), запрос и подтверждение верификации, управление статусами верификации.
*   **Управление Настройками:** Хранение и обновление пользовательских настроек по категориям (приватность, уведомления, интерфейс, безопасность).

### 1.3. Основные Технологии
*   **Backend:** Go (версия 1.21+), Gin или Echo (REST), google.golang.org/grpc (gRPC), Gorilla WebSocket.
*   **Базы данных:** PostgreSQL (версия 15+), Redis (кэш/сессии).
*   **Сообщения:** Kafka.
*   **Frontend (Клиент):** Flutter, Dart.
*   **Инфраструктура:** Docker, Kubernetes, Prometheus, Grafana, ELK Stack/Loki, Jaeger.

### 1.4. Термины и Определения (Glossary)
*   **Аккаунт (Account)**: Учетная запись Пользователя.
*   **Профиль пользователя (User Profile)**: Публичная или частично публичная информация о Пользователе.
*   **Событие (Event)**: Сообщение о произошедшем действии (CloudEvents формат).
*   (Полный глоссарий см. в "Едином глоссарии терминов и определений для российского аналога Steam.txt" и разделе 1.3 исходной спецификации Account Service).

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Account Service построен на принципах чистой архитектуры (Clean Architecture) и микросервисного подхода.
*   Диаграмма слоев (из исходной спецификации):
    ```mermaid
    graph TD
        subgraph "API Layer"
            A[REST API]
            B[gRPC API]
            C[WebSocket API]
        end
        subgraph "Use Case Layer"
            D[RegisterAccountUseCase]
            E[UpdateProfileUseCase]
            F[VerifyEmailUseCase]
            G[GetSettingsUseCase]
        end
        subgraph "Domain Layer"
            H[Account Entity]
            I[Profile Entity]
            J[Verification Entity]
            K[Settings Entity]
            L[Repository Interfaces]
        end
        subgraph "Infrastructure Layer"
            M[PostgreSQL Repositories]
            N[Redis Cache]
            O[Kafka Producer/Consumer]
            P[gRPC Clients]
            Q[External Service Clients]
        end

        A --> D; B --> E; C --> F
        D --> L; E --> L; F --> L; G --> L
        L -- Implemented by --> M
        L -- Implemented by --> N
        D --> O; E --> O; D --> P
        M --> R[(PostgreSQL)]; N --> S[(Redis)]; O --> T[(Kafka)]
        P --> U[Other Microservices]; Q --> V[External Systems]
    ```

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления / API Layer)
*   Ответственность: Обработка входящих запросов (REST, gRPC, WebSocket), валидация данных, преобразование DTO, вызов Use Case Layer.
*   Ключевые компоненты/модули:
    *   HTTP Handlers/Controllers: Gin или Echo для REST.
    *   gRPC Service Implementations: Стандартные библиотеки Go для gRPC.
    *   WebSocket Handlers: Gorilla WebSocket.
    *   DTOs: Для преобразования данных между API и Use Case слоями.
    *   Валидаторы: go-playground/validator.

#### 2.2.2. Application Layer (Прикладной Слой / Use Case Layer)
*   Ответственность: Содержит основную бизнес-логику сервиса, оркестрируя взаимодействие между различными компонентами и репозиториями для выполнения конкретных сценариев (например, регистрация Аккаунта, обновление Профиля).
*   Ключевые компоненты/модули:
    *   Use Case Services: `RegisterAccountUseCase`, `UpdateProfileUseCase`, `VerifyEmailUseCase`, `GetSettingsUseCase`.
    *   Интерфейсы для репозиториев и внешних сервисов.

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Содержит основные доменные сущности (Аккаунт, Профиль, Верификация, Настройки), их бизнес-правила и интерфейсы репозиториев.
*   Ключевые компоненты/модули:
    *   Entities: `Account`, `AuthMethod`, `Profile`, `ContactInfo`, `Setting`, `Avatar`, `ProfileHistory`.
    *   Repository Interfaces.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Реализация интерфейсов, определенных в Domain и Application Layers, для взаимодействия с PostgreSQL, Redis, Kafka и другими микросервисами.
*   Ключевые компоненты/модули:
    *   PostgreSQL Repositories: GORM или sqlx.
    *   Redis Cache: go-redis/redis.
    *   Kafka Producer/Consumer: confluent-kafka-go или sarama.
    *   gRPC Clients: Для взаимодействия с другими микросервисами.
    *   External Service Clients: (например, для SMS-шлюза, S3).

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL**: `/api/v1`
*   **Формат данных**: JSON
*   **Аутентификация**: JWT Bearer Token (валидируется через API Gateway / Auth Service).
*   **Стандарт ответа (Успех)**: (см. раздел 5.1 исходной спецификации)
*   **Стандарт ответа (Ошибка)**: (см. раздел 5.1 исходной спецификации)

#### 3.1.1. Ресурс: Аккаунты (`/accounts`)
*   **`POST /accounts`**: Регистрация нового Аккаунта.
    *   Запрос: `{ "username": "...", "email": "...", "password": "..." }` (пароль передается в Auth Service) или `{ "provider": "google", "token": "..." }`.
    *   Ответ: `201 Created`, `data: { "id": "...", "username": "...", "status": "pending" }`.
*   **`GET /accounts/{id}`**: Получение информации об Аккаунте.
*   **`GET /accounts/me`**: Получение информации о текущем Аккаунте.
*   **`PUT /accounts/{id}/status`**: Обновление статуса Аккаунта (только Admin).
*   **`DELETE /accounts/{id}`**: Удаление Аккаунта (мягкое).
*   **`GET /accounts`**: Поиск Аккаунтов (только Admin).

#### 3.1.2. Ресурс: Профили (`/accounts/{id}/profile`)
*   **`GET /accounts/{id}/profile`**: Получение Профиля Пользователя.
*   **`GET /accounts/me/profile`**: Получение Профиля текущего Пользователя.
*   **`PUT /accounts/{id}/profile`**: Обновление Профиля.
*   **`POST /accounts/{id}/avatar`**: Загрузка Аватара.
*   **`GET /accounts/{id}/profile/history`**: Получение истории изменений Профиля.

#### 3.1.3. Ресурс: Контактная информация (`/accounts/{id}/contact-info`)
*   **`GET /accounts/{id}/contact-info`**: Получение списка контактной информации.
*   **`POST /accounts/{id}/contact-info`**: Добавление новой контактной информации.
*   **`PUT /accounts/{id}/contact-info/{contact_id}`**: Обновление контактной информации.
*   **`DELETE /accounts/{id}/contact-info/{contact_id}`**: Удаление контактной информации.
*   **`POST /accounts/{id}/contact-info/{type}/verification-request`**: Запрос кода верификации.
*   **`POST /accounts/{id}/contact-info/{type}/verify`**: Подтверждение кода верификации.

#### 3.1.4. Ресурс: Настройки (`/accounts/{id}/settings`)
*   **`GET /accounts/{id}/settings`**: Получение всех категорий настроек.
*   **`GET /accounts/{id}/settings/{category}`**: Получение настроек конкретной категории.
*   **`PUT /accounts/{id}/settings/{category}`**: Обновление настроек категории.
*   (Более детальное описание см. в разделе 5.1 исходной спецификации)

### 3.2. gRPC API
*   Используется для синхронного взаимодействия между микросервисами.
*   Пример `.proto` файла (`account.proto` из раздела 5.2 исходной спецификации):
    ```protobuf
    syntax = "proto3";
    package account.v1;
    // ... (содержимое proto файла) ...
    service AccountService {
      rpc GetAccount(GetAccountRequest) returns (AccountResponse);
      rpc GetAccounts(GetAccountsRequest) returns (GetAccountsResponse);
      rpc GetProfile(GetProfileRequest) returns (ProfileResponse);
      // ... и другие методы ...
    }
    // ... (определения сообщений) ...
    ```
*   (Полные определения см. в разделе 5.2 исходной спецификации).

### 3.3. WebSocket API (если применимо)
*   Используется для отправки уведомлений об изменениях состояния в реальном времени клиентам (Flutter).
*   Эндпоинт: `/ws/v1/notifications` (может быть предоставлен отдельным сервисом или API Gateway).
*   Аутентификация: JWT токен.
*   Примеры типов уведомлений: `profile.updated`, `account.status.changed`, `contact.verified`.
*   (Детали см. в разделе 5.3 исходной спецификации).

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности
*   **Account**: Базовая информация об Аккаунте (ID, Username, Status, Timestamps).
*   **AuthMethod**: Способ аутентификации (ID, AccountID, Type, Identifier, IsVerified).
*   **Profile**: Расширенная информация Профиля (ID, AccountID, Nickname, Bio, AvatarURL).
*   **ContactInfo**: Контактная информация (ID, AccountID, Type, Value, IsVerified).
*   **Setting**: Пользовательские настройки (ID, AccountID, Category, SettingsJSON).
*   **Avatar**: Загруженные изображения профиля (ID, AccountID, URL, IsCurrent).
*   **ProfileHistory**: История изменений Профиля.
*   (Полный список полей см. в разделе 3.3.1 исходной спецификации).

### 4.2. Схема Базы Данных (PostgreSQL)
*   Диаграмма ERD:
    ```mermaid
    erDiagram
        ACCOUNTS {
            UUID id PK
            VARCHAR(64) username
            VARCHAR(20) status
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
            TIMESTAMPTZ deleted_at
        }

        AUTH_METHODS {
            UUID id PK
            UUID account_id FK
            VARCHAR(20) type
            VARCHAR(255) identifier
            TEXT secret
            BOOLEAN is_verified
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }

        PROFILES {
            UUID id PK
            UUID account_id FK
            VARCHAR(64) nickname
            TEXT bio
            CHAR(2) country
            VARCHAR(100) city
            DATE birth_date
            VARCHAR(10) gender
            VARCHAR(20) visibility
            TEXT avatar_url
            TEXT banner_url
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }

        CONTACT_INFO {
            UUID id PK
            UUID account_id FK
            VARCHAR(10) type
            VARCHAR(255) value
            BOOLEAN is_primary
            BOOLEAN is_verified
            VARCHAR(10) verification_code
            INT verification_attempts
            TIMESTAMPTZ verification_expires_at
            VARCHAR(20) visibility
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }

        USER_SETTINGS {
            UUID id PK
            UUID account_id FK
            VARCHAR(50) category
            JSONB settings
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }

        AVATARS {
            UUID id PK
            UUID account_id FK
            TEXT url
            BOOLEAN is_current
            TIMESTAMPTZ created_at
        }

        PROFILE_HISTORY {
            BIGSERIAL id PK
            UUID profile_id FK
            VARCHAR(50) field_name
            TEXT old_value
            TEXT new_value
            UUID changed_by_account_id FK
            TIMESTAMPTZ changed_at
        }

        ACCOUNTS ||--o{ AUTH_METHODS : "has"
        ACCOUNTS ||--|{ PROFILES : "has one"
        ACCOUNTS ||--o{ CONTACT_INFO : "has"
        ACCOUNTS ||--o{ USER_SETTINGS : "has"
        ACCOUNTS ||--o{ AVATARS : "has"
        PROFILES ||--o{ PROFILE_HISTORY : "has history"
        ACCOUNTS ||--o{ PROFILE_HISTORY : "changed by (optional)"


    ```
*   DDL (из раздела 3.3.2 исходной спецификации):
    ```sql
    -- Аккаунты
    CREATE TABLE accounts ( id UUID PRIMARY KEY ..., username VARCHAR(64) NOT NULL UNIQUE, status VARCHAR(20) ...);
    -- Методы аутентификации
    CREATE TABLE auth_methods ( id UUID PRIMARY KEY ..., account_id UUID NOT NULL REFERENCES accounts(id) ...);
    -- Профили пользователей
    CREATE TABLE profiles ( id UUID PRIMARY KEY ..., account_id UUID NOT NULL UNIQUE REFERENCES accounts(id) ...);
    -- Контактная информация
    CREATE TABLE contact_info ( id UUID PRIMARY KEY ..., account_id UUID NOT NULL REFERENCES accounts(id) ...);
    -- Настройки пользователей
    CREATE TABLE user_settings ( id UUID PRIMARY KEY ..., account_id UUID NOT NULL REFERENCES accounts(id) ...);
    -- Аватары
    CREATE TABLE avatars ( id UUID PRIMARY KEY ..., account_id UUID NOT NULL REFERENCES accounts(id) ...);
    -- История изменений профилей
    CREATE TABLE profile_history ( id BIGSERIAL PRIMARY KEY ..., profile_id UUID NOT NULL REFERENCES profiles(id) ...);
    ```
*   (Полную схему см. в разделе 3.3.2 исходной спецификации).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений**: Kafka.
*   **Формат событий**: CloudEvents JSON.
*   **Топики Kafka**: `account-events`, `profile-events`.
*   **Основные типы событий**:
    *   `account.created`
    *   `account.status.updated` (включая blocked, deleted, activated)
    *   `account.contact.added`
    *   `account.contact.updated`
    *   `account.contact.verified`
    *   `account.contact.verification.requested`
    *   `profile.updated`
    *   `profile.avatar.updated`
    *   `account.settings.updated`
*   (Пример события `account.created` и детали см. в разделе 5.4 исходной спецификации).

### 5.2. Потребляемые События (Consumed Events)

Account Service подписывается на следующие события от других микросервисов для поддержания консистентности данных и выполнения зависимых бизнес-процессов:

*   **`auth.user.credentials.created.v1`** (ранее мог называться `auth.user.registered.v1`)
    *   **Источник:** Auth Service
    *   **Назначение:** Завершение создания Аккаунта после того, как Auth Service успешно создал основные учетные данные (например, логин/пароль или привязал внешний ID провайдера). Account Service может на основе этого события активировать аккаунт, создать начальный профиль или выполнить другие пост-регистрационные задачи.
    *   **Ожидаемая структура `data` (пример):**
        ```json
        {
          "auth_method_id": "uuid", // ID метода аутентификации, созданного в Auth Service
          "account_id": "uuid",     // ID аккаунта, который был создан (или должен быть создан) в Account Service
          "username": "user123",
          "email": "user@example.com", // Если email является частью учетных данных
          "auth_provider": "password", // или "google", "telegram"
          "registration_timestamp": "ISO8601_timestamp"
        }
        ```
    *   **Действия Account Service:** Создание или обновление записи `accounts` и `auth_methods`, создание `profiles`, установка начального статуса.

*   **`admin.user.status.updated.v1`**
    *   **Источник:** Admin Service
    *   **Назначение:** Обновление статуса Аккаунта Пользователя (например, `blocked`, `unblocked`, `deleted_by_admin`) по инициативе администратора.
    *   **Ожидаемая структура `data` (пример):**
        ```json
        {
          "account_id": "uuid",
          "new_status": "blocked", // "active", "deleted"
          "reason": "Нарушение правил платформы, пункт 5.2.", // Опционально
          "changed_by_admin_id": "uuid",
          "timestamp": "ISO8601_timestamp"
        }
        ```
    *   **Действия Account Service:** Обновление поля `status` в таблице `accounts`. Может также инициировать анонимизацию данных, если статус `deleted`.

*   **`payment.subscription.status.changed.v1` (Условно/На будущее)**
    *   **Источник:** Payment Service
    *   **Назначение:** Обновление статуса Аккаунта или его функций на основе изменения статуса подписки Пользователя (если такая логика будет реализована, например, премиум-аккаунты).
    *   **Ожидаемая структура `data` (пример):**
        ```json
        {
          "account_id": "uuid",
          "subscription_id": "uuid",
          "new_status": "active", // "expired", "cancelled"
          "plan_type": "premium",
          "valid_until": "ISO8601_timestamp", // Опционально
          "timestamp": "ISO8601_timestamp"
        }
        ```
    *   **Действия Account Service:** Потенциально, обновление специальных полей в `accounts` или `user_settings`, влияющих на доступные функции. *На данный момент, детальная логика обработки этого события не определена и требует дальнейшей проработки при внедрении соответствующего функционала.*

*Примечание: Точные имена событий, их версии и структура полезной нагрузки (`data`) должны быть согласованы и задокументированы в рамках общих "Стандартов событий" платформы и спецификаций взаимодействующих сервисов.*

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Auth Service**: Делегирование создания/проверки учетных данных (gRPC/REST).
*   **Notification Service**: Отправка email/SMS уведомлений (Kafka события или API).
*   **Social Service**: Предоставление данных профилей, получение событий об изменениях соц. связей (gRPC, Kafka).
*   **Payment Service**: Предоставление информации об аккаунте для транзакций (gRPC).
*   **Library Service / Catalog Service / etc.**: Предоставление информации об аккаунте/профиле (gRPC).
*   **Admin Service**: Предоставление API для административных действий (REST/gRPC).
*   **API Gateway**: Проксирование REST-запросов, обогащение запросов.
*   (Детали см. в разделе 3.4 и 6.1 исходной спецификации).

### 6.2. Внешние Системы
*   **SMS-шлюз**: Отправка SMS с кодами верификации (через Notification Service или напрямую).
*   **Email-провайдер**: Отправка email с кодами верификации (через Notification Service).
*   **Хранилище файлов (S3-совместимое)**: Загрузка/удаление Аватаров и баннеров.
*   **Системы аутентификации (Google, Telegram, etc.)**: Через Auth Service.
*   (Детали см. в разделе 6.2 исходной спецификации).

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `SERVICE_PORT_HTTP`: Порт для HTTP REST API.
*   `SERVICE_PORT_GRPC`: Порт для gRPC API.
*   `DATABASE_URL` / `DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME`: Параметры подключения к PostgreSQL.
*   `REDIS_ADDR, REDIS_PASSWORD`: Параметры подключения к Redis.
*   `KAFKA_BROKERS`: Список брокеров Kafka.
*   `AUTH_SERVICE_GRPC_ADDR`: Адрес gRPC Auth Service.
*   `NOTIFICATION_SERVICE_ADDR`: Адрес Notification Service (API/gRPC).
*   `LOG_LEVEL`: Уровень логирования.
*   `JAEGER_AGENT_HOST`: Адрес Jaeger агента.
*   (Полный список см. в разделе 9.2 исходной спецификации).

### 7.2. Файлы Конфигурации (если применимо)
*   Обычно конфигурация через переменные окружения, но может использоваться YAML файл, как указано в разделе 9.2 исходной спецификации.

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   REST API: Стандартные коды HTTP, JSON-формат ошибки (см. раздел 5.1 исходной спецификации).
*   gRPC API: Стандартные коды gRPC.
*   Логирование всех ошибок с контекстом.

### 8.2. Распространенные Коды Ошибок
*   `VALIDATION_ERROR` (HTTP 400): Невалидные данные.
*   `RESOURCE_NOT_FOUND` (HTTP 404): Ресурс не найден.
*   `UNAUTHENTICATED` (HTTP 401): Ошибка аутентификации (делегировано Auth Service / API Gateway).
*   `PERMISSION_DENIED` (HTTP 403): Недостаточно прав.
*   `CONFLICT` (HTTP 409): Конфликт данных (например, Username/Email уже занят).
*   `INTERNAL_SERVER_ERROR` (HTTP 500): Внутренняя ошибка сервера.
*   (Детали и примеры см. в разделе 4 и 5.6 исходной спецификации).

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   JWT (Access Token + Refresh Token), выдаваемые Auth Service. Валидация через API Gateway / Auth Service.
*   (Детали см. в разделе 7.1 исходной спецификации).

### 9.2. Авторизация
*   RBAC. Роли из JWT. Проверка разрешений через middleware или Casbin.
*   Матрица доступа приведена в Приложении 10.1 исходной спецификации.
*   (Детали см. в разделе 7.2 исходной спецификации).

### 9.3. Защита Данных
*   Шифрование TLS для внешних и внутренних коммуникаций.
*   Шифрование чувствительных данных в покое (если применимо).
*   Соответствие ФЗ-152, анонимизация при удалении.
*   Валидация ввода.
*   (Детали см. в разделе 7.3 исходной спецификации).

### 9.4. Управление Секретами
*   Kubernetes Secrets.
*   (Детали см. в разделе 7.4 исходной спецификации).
*   **Защита от атак**: Rate limiting (API Gateway), защита от брутфорса (Auth Service), CSRF, XSS. (см. раздел 7.5 исходной спецификации).

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile**: Многоэтапная сборка (golang -> alpine/distroless), non-root пользователь. (см. раздел 9.1 исходной спецификации).
*   **Kubernetes Manifests (Helm Chart)**: Deployment, Service, ConfigMap, Secret, HPA, PDB, Probes. (см. раздел 9.1 исходной спецификации).

### 10.2. Зависимости при Развертывании
*   PostgreSQL, Redis, Kafka, Auth Service, Notification Service.
*   (Детали см. в разделе 3.5 исходной спецификации).

### 10.3. CI/CD
*   GitLab CI. Этапы: Build, Test, Lint, Security Scan, Push, Deploy. Семантическое версионирование.
*   (Детали см. в разделе 9.3 исходной спецификации).

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат**: Структурированные логи JSON (Zap/Logrus).
*   **Содержание**: Timestamp, level, message, trace ID, span ID, контекст.
*   **Интеграция**: ELK/EFK.
*   (Детали см. в разделе 8.2 исходной спецификации).

### 11.2. Мониторинг
*   **Метрики**: Prometheus (`/metrics`). Запросы API/gRPC (количество, задержка, ошибки), активные WebSocket, события Kafka, ресурсы (CPU, RAM), зависимости.
*   **Дашборды**: Grafana (RED).
*   **Алертинг**: Alertmanager для критических ситуаций.
*   (Детали см. в разделе 8.1 исходной спецификации).

### 11.3. Трассировка
*   **Инструментация**: OpenTelemetry.
*   **Контекст**: Передача Trace ID, Span ID.
*   **Экспорт**: Jaeger.
*   (Детали см. в разделе 8.3 исходной спецификации).
*   **Аудит**: Логирование важных событий безопасности и бизнеса. (см. раздел 8.4 исходной спецификации).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**: P95 < 100мс (чтение), < 200мс (запись); P99 < 200мс (чтение), < 400мс (запись); RPS: 2000+ (чтение), 500+ (запись).
*   **Надежность**: Доступность 99.99%; RTO < 5 мин; RPO < 1 мин.
*   **Масштабируемость**: 10-50 млн. аккаунтов, горизонтальное масштабирование.
*   **Безопасность**: ФЗ-152, шифрование, OWASP Top 10.
*   **Сопровождаемость**: Чистый код, тесты > 80%, CI/CD.
*   **Совместимость**: Flutter клиент, другие микросервисы.
*   (Полный список см. в разделе 2.3 исходной спецификации).

## 13. Приложения (Appendices) (Опционально)
*   Матрица доступа (см. Приложение 10.1 в исходной спецификации Account Service).

---
*Этот шаблон является отправной точкой и может быть адаптирован под конкретные нужды проекта и сервиса.*
