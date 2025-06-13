# Спецификация Микросервиса: Developer Service (Сервис для Разработчиков)

**Версия:** 1.1 (адаптировано из предыдущей версии)
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Developer Service, предназначенного для разработчиков и издателей игр на платформе "Российский Аналог Steam".
*   **Роль в общей архитектуре платформы:** Developer Service предоставляет интерфейс (Портал Разработчика) и API для управления аккаунтами разработчиков, загрузки и обновления игр и их контента, управления метаданными и маркетинговыми материалами, доступа к аналитике продаж и использования продуктов, а также взаимодействия с финансовыми аспектами платформы (управление выплатами).
*   **Основные бизнес-задачи:** Обеспечение разработчиков инструментами для самостоятельной публикации и поддержки их продуктов (игр, DLC, ПО) на платформе, управление жизненным циклом этих продуктов, предоставление релевантной аналитики и финансовых отчетов.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   Регистрация и управление аккаунтом разработчика/издателя (профиль компании, юридическая информация, управление командой и ролями доступа к порталу).
*   Управление проектами (играми/продуктами): создание карточки продукта, загрузка и управление билдами и версиями, управление метаданными (локализованные названия, описания, теги, жанры, системные требования, возрастные рейтинги), управление медиа-контентом (скриншоты, трейлеры, арты).
*   Управление ценообразованием: установка базовой цены, региональных цен, создание и управление скидками и промо-периодами для своих продуктов (в координации с Catalog Service).
*   Процесс публикации: подача продукта на модерацию, отслеживание статуса модерации, публикация одобренных продуктов, управление видимостью продуктов в каталоге.
*   Панель аналитики: доступ к дашбордам с метриками (продажи, доход, количество установок, DAU/MAU по своим продуктам), возможность генерации стандартных отчетов.
*   Финансовый раздел: просмотр баланса, истории транзакций (продажи, возвраты, комиссии), управление реквизитами для выплат, формирование запросов на выплаты.
*   Управление SDK и API ключами: доступ к SDK платформы, документации, управление API ключами для автоматизации процессов CI/CD и взаимодействия с API платформы.
*   Система уведомлений: получение уведомлений о важных событиях (статус модерации, обновления платформы, финансовые операции, сообщения от поддержки).

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+, согласно `../../../../project_technology_stack.md`).
*   **Веб-фреймворк (REST API):** Echo (`github.com/labstack/echo/v4`) или Gin (`github.com/gin-gonic/gin`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **База данных:** PostgreSQL (версия 15+) для хранения структурированных данных: аккаунты разработчиков, метаданные игр (черновики, специфичные для разработчика данные), информация о финансах. Драйвер: GORM (`gorm.io/gorm`) с `gorm.io/driver/postgres` или `pgx` (`github.com/jackc/pgx/v5`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Кэширование:** Redis (версия 7.0+) для кэширования часто запрашиваемых данных, сессий портала. Клиент: `go-redis/redis` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Очереди сообщений/События:** Apache Kafka (клиент `github.com/confluentinc/confluent-kafka-go` или `github.com/segmentio/kafka-go`, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Хранилище файлов (билды, медиа):** S3-совместимое объектное хранилище (например, MinIO, Yandex Object Storage) (согласно `../../../../project_technology_stack.md`).
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Логирование:** Zap (`go.uber.org/zap`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Мониторинг/Трассировка:** OpenTelemetry SDK, Prometheus client (`github.com/prometheus/client_golang`). (согласно `../../../../project_observability_standards.md`).
*   **Инфраструктура:** Docker, Kubernetes.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   **Разработчик (Developer Account):** Учетная запись компании или индивидуального разработчика/издателя.
*   **Продукт (Product/Game):** Игра, DLC или другой цифровой товар, управляемый разработчиком.
*   **Билд (Build):** Конкретная сборка исполняемых файлов продукта.
*   **Версия (Version):** Версия продукта, связанная с конкретным билдом и метаданными.
*   **Портал Разработчика (Developer Portal):** Веб-интерфейс, предоставляемый Developer Service для разработчиков.
*   Для других общих терминов см. `../../../../project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Developer Service будет реализован как модульный монолит или набор тесно связанных микросервисов, придерживаясь принципов Чистой Архитектуры (Clean Architecture) для разделения ответственностей.
*   Ключевые модули включают: Управление Аккаунтами Разработчиков, Управление Продуктами (Игры/DLC), Управление Загрузками (билды, медиа), Модуль Аналитики (отображение данных), Финансовый Модуль (выплаты), Управление API Ключами.

**Диаграмма Архитектуры (Clean Architecture):**
```mermaid
graph TD
    subgraph Developer Portal & External API Clients
        DevPortal[Портал Разработчика (Веб-интерфейс)]
        DevAPIClient[Клиенты API Разработчика (CI/CD, утилиты)]
    end

    subgraph Developer Service
        direction TB

        subgraph PresentationLayer [Presentation Layer (Адаптеры Транспорта)]
            REST_API[REST API (Echo/Gin)]
            GRPC_API[gRPC API (для внутренних нужд, если потребуется)]
        end

        subgraph ApplicationLayer [Application Layer (Сценарии Использования)]
            DevAccountSvc[Управление Аккаунтом Разработчика]
            GameManagementSvc[Управление Продуктами (Игры)]
            FileUploadSvc[Управление Загрузками]
            AnalyticsViewSvc[Просмотр Аналитики]
            FinanceSvc[Финансовые Операции (Выплаты)]
            ApiKeyMgmtSvc[Управление API Ключами]
        end

        subgraph DomainLayer [Domain Layer (Бизнес-логика и Сущности)]
            Entities[Сущности (Developer, Game, GameVersion, Payout)]
            Aggregates[Агрегаты (DeveloperProfile, GameProduct)]
            DomainEvents[Доменные События (GameSubmitted, PayoutRequested)]
            RepositoryIntf[Интерфейсы Репозиториев]
        end

        subgraph InfrastructureLayer [Infrastructure Layer (Внешние Зависимости)]
            PostgresAdapter[Адаптер PostgreSQL]
            S3Adapter[Адаптер S3-хранилища]
            RedisAdapter[Адаптер Redis (Кэш)]
            KafkaProducer[Продюсер Kafka (События)]
            ServiceClients[Клиенты других микросервисов (Auth, Catalog, Payment)]
            Config[Конфигурация (Viper)]
            Logging[Логирование (Zap)]
        end

        REST_API --> DevAccountSvc
        REST_API --> GameManagementSvc
        REST_API --> FileUploadSvc
        REST_API --> AnalyticsViewSvc
        REST_API --> FinanceSvc
        REST_API --> ApiKeyMgmtSvc

        ApplicationLayer --> DomainLayer
        ApplicationLayer --> InfrastructureLayer
        DomainLayer ----> RepositoryIntf
        InfrastructureLayer -- Implements --> RepositoryIntf
    end

    DevPortal --> REST_API
    DevAPIClient --> REST_API

    PostgresAdapter --> DB[(PostgreSQL)]
    S3Adapter --> S3[(S3 Хранилище)]
    RedisAdapter --> Cache[(Redis)]
    KafkaProducer --> Kafka[Kafka Broker]
    ServiceClients --> OtherServices[Другие Микросервисы]

    classDef layer_boundary fill:#f9f9f9,stroke:#333,stroke-width:2px,color:#333
    classDef component_major fill:#e6f0ff,stroke:#007bff,color:#000
    classDef component_minor fill:#d4edda,stroke:#28a745,color:#000
    classDef datastore fill:#f8d7da,stroke:#dc3545,color:#000

    class PresentationLayer,ApplicationLayer,DomainLayer,InfrastructureLayer layer_boundary
    class REST_API,GRPC_API,DevAccountSvc,GameManagementSvc,FileUploadSvc,AnalyticsViewSvc,FinanceSvc,ApiKeyMgmtSvc,Entities,Aggregates,DomainEvents,RepositoryIntf component_major
    class PostgresAdapter,S3Adapter,RedisAdapter,KafkaProducer,ServiceClients,Config,Logging component_minor
    class DB,S3,Cache,Kafka,OtherServices datastore
```

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления)
*   **Ответственность:** Обработка входящих HTTP REST запросов от Портала Разработчика и публичного API для разработчиков. Валидация данных запроса (DTO), аутентификация (через Auth Service), авторизация (на основе роли в команде разработчика), вызов соответствующей бизнес-логики в Application Layer.
*   **Ключевые компоненты/модули:** HTTP хендлеры (контроллеры) на базе Echo/Gin, DTO для запросов и ответов API, middleware для аутентификации и авторизации.

#### 2.2.2. Application Layer (Прикладной Слой)
*   **Ответственность:** Реализация сценариев использования (use cases), связанных с управлением аккаунтами разработчиков, их продуктами, финансами и т.д. Координирует взаимодействие между Domain Layer и Infrastructure Layer.
*   **Ключевые компоненты/модули:** Сервисы сценариев использования (например, `DeveloperAccountService`, `GameProductService`, `PayoutService`), обработчики команд и запросов (если используется CQRS).

#### 2.2.3. Domain Layer (Доменный Слой)
*   **Ответственность:** Содержит бизнес-сущности, агрегаты, доменные события и бизнес-правила, специфичные для Developer Service.
*   **Ключевые компоненты/модули:** Сущности (`Developer`, `DeveloperTeamMember`, `Game`, `GameVersion`, `GameMetadata`, `GamePricing`, `DeveloperPayout`, `DeveloperAPIKey`), объекты-значения, доменные сервисы, интерфейсы репозиториев.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   **Ответственность:** Реализация интерфейсов репозиториев для работы с PostgreSQL. Взаимодействие с S3-совместимым хранилищем для загрузки и управления файлами. Отправка событий в Kafka. Взаимодействие с другими микросервисами (Auth, Catalog, Payment, Admin, Analytics, Notification) через их gRPC/REST API.
*   **Ключевые компоненты/модули:** Реализации репозиториев для PostgreSQL, S3 клиент, Kafka продюсер, Redis клиент, gRPC/HTTP клиенты для других сервисов.

## 3. API Endpoints

### 3.1. REST API
*   **Префикс:** `/api/v1/developer` (маршрутизируется через API Gateway).
*   **Аутентификация:** JWT Bearer Token, полученный от Auth Service. Проверяется на API Gateway или в middleware Developer Service.
*   **Авторизация:** На основе `developer_id` (извлекается из JWT или связи `user_id` с `developer_id`) и роли пользователя в команде разработчика (например, `owner`, `admin`, `editor`, `viewer`).
*   **Формат ответа об ошибке (согласно `../../../../project_api_standards.md`):**
    ```json
    {
      "errors": [
        {
          "code": "ERROR_CODE_UPPER_SNAKE_CASE",
          "title": "Краткое описание ошибки на русском",
          "detail": "Полное описание ошибки с контекстом.",
          "source": { "pointer": "/data/attributes/field_name", "parameter": "query_param_name" }
        }
      ]
    }
    ```

#### 3.1.1. Аккаунты разработчиков (Developer Accounts)
*   **`POST /accounts`**
    *   Описание: Регистрация нового аккаунта разработчика/издателя.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "developerAccountCreation",
            "attributes": {
              "company_name": "Моя Игровая Студия",
              "legal_entity_type": "ООО",
              "tax_id": "1234567890",
              "contact_email": "dev@example.com",
              "country_code": "RU"
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created):
        ```json
        {
          "data": {
            "type": "developerAccount",
            "id": "dev-uuid-abc",
            "attributes": { /* ... поля созданного аккаунта ... */ }
          }
        }
        ```
    *    Пример ответа (Ошибка 400 Validation Error - стандартизированный):
        ```json
        {
          "errors": [
            {
              "code": "VALIDATION_ERROR",
              "title": "Ошибка валидации",
              "detail": "Поле 'tax_id' должно быть валидным ИНН.",
              "source": { "pointer": "/data/attributes/tax_id" }
            }
          ]
        }
        ```
    *   Требуемые права доступа: Аутентифицированный пользователь (для привязки к его User ID).
*   **`GET /accounts/me`**
    *   Описание: Получение информации о текущем аккаунте разработчика, к которому привязан пользователь.
    *   Пример ответа (Успех 200 OK): (Аналогично ответу POST /accounts)
    *   Требуемые права доступа: Участник команды разработчика.

#### 3.1.2. Управление Продуктами (Игры/Games)
*   **`POST /games`**
    *   Описание: Создание нового продукта (игры, DLC, ПО) в каталоге разработчика (начальная регистрация).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "productCreation",
            "attributes": {
              "title": "Моя Новая Супер Игра",
              "product_type": "game"
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created):
        ```json
        {
          "data": {
            "type": "game",
            "id": "game-uuid-xyz",
            "attributes": {
              "title": "Моя Новая Супер Игра",
              "product_type": "game",
              "status": "draft",
              "developer_id": "dev-uuid-abc",
              "created_at": "2024-03-15T12:00:00Z"
            }
          }
        }
        ```
    *   Требуемые права доступа: `developer_admin` или `game_manager` в команде.
*   **`GET /games/{game_id}`**
    *   Описание: Получение детальной информации о продукте разработчика.
    *   Пример ответа (Успех 200 OK): (Полная информация о Game, включая GameMetadata, GamePricing и т.д.)
    *   Требуемые права доступа: Участник команды разработчика (с правами на просмотр игры).

#### 3.1.3. Версии Игр и Загрузка Билдов
*   **`POST /games/{game_id}/versions`**
    *   Описание: Создание новой версии для игры (перед загрузкой билдов).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "gameVersionCreation",
            "attributes": {
              "version_name": "1.0.1",
              "changelog": { "ru-RU": "Исправлены ошибки, улучшена производительность." }
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Возвращает ID созданной версии)
    *   Требуемые права доступа: `developer_admin` или `build_manager`.
*   **`POST /games/{game_id}/versions/{version_id}/builds/upload-url`**
    *   Описание: Получение pre-signed URL для загрузки файла билда в S3.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "buildUploadRequest",
            "attributes": {
              "file_name": "mygame_v1.0.1_windows_x64.zip",
              "file_size_bytes": 1073741824,
              "content_type": "application/zip",
              "platform": "windows_x64"
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "presignedUploadUrl",
            "attributes": {
              "upload_url": "https://s3.example.com/bucket/path?AWSAccessKeyId=...",
              "method": "PUT",
              "expires_in_seconds": 3600,
              "internal_build_id": "build-uuid-temp"
            }
          }
        }
        ```
    *   Требуемые права доступа: `developer_admin` или `build_manager`.
*   **`POST /games/{game_id}/versions/{version_id}/builds/upload-complete`**
    *   Описание: Уведомление сервиса об успешной загрузке билда в S3.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "buildUploadCompletion",
            "attributes": {
              "internal_build_id": "build-uuid-temp",
              "s3_path": "path/to/mygame_v1.0.1_windows_x64.zip",
              "file_hash_sha256": "abcdef123..."
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): (Обновленный статус версии или билда)
    *   Требуемые права доступа: `developer_admin` или `build_manager`.

#### 3.1.4. Публикация
*   **`POST /games/{game_id}/versions/{version_id}/submit-for-review`**
    *   Описание: Отправка конкретной версии игры на модерацию.
    *   Пример ответа (Успех 200 OK): `{ "data": { "status": "in_review" } }`
    *   Требуемые права доступа: `developer_admin` или `release_manager`.

#### 3.1.5. Финансы
*   **`POST /finance/payouts/requests`**
    *   Описание: Создание запроса на выплату средств.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "payoutRequest",
            "attributes": {
              "amount_minor_units": 5000000,
              "currency_code": "RUB",
              "payment_method_id": "bank-account-uuid-123"
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Информация о созданном запросе на выплату)
    *   Требуемые права доступа: `developer_admin` или `finance_manager`.

### 3.2. gRPC API
*   Developer Service в основном **потребляет** gRPC API других сервисов (например, Auth Service для валидации токенов разработчиков, Catalog Service для получения информации о статусе их продуктов, Payment Service для финансовых данных).
*   На данный момент Developer Service **не предоставляет** публично доступных gRPC методов для других сервисов. Если в будущем потребуется специфичное межсервисное взаимодействие, где Developer Service будет выступать сервером, оно будет спроектировано и задокументировано отдельно.

### 3.3. WebSocket API
*   Не планируется для Developer Service на данном этапе.

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`Developer` (Разработчик/Издатель)** (Как в существующем документе)
*   **`DeveloperTeamMember` (Член Команды Разработчика)** (Как в существующем документе)
*   **`Game` (Продукт/Игра)**
    *   (Как в существующем документе, с добавлением)
    *   `draft_metadata` (JSONB): Временное хранилище метаданных продукта, редактируемых разработчиком перед отправкой в Catalog Service. Структура соответствует модели метаданных Catalog Service.
    *   `draft_pricing` (JSONB): Временное хранилище предложений по ценам, редактируемых разработчиком.
*   **`GameVersion` (Версия Продукта)** (Как в существующем документе)
*   **`GameMetadataLocalized`**: Эта сущность управляется и хранится преимущественно в **Catalog Service**. Developer Service взаимодействует с API Catalog Service для предложения изменений или просмотра этих данных. В Developer Service поле `Game.draft_metadata` используется для подготовки этих данных.
    *   Поля: `game_id`, `language_code`, `title`, `description_short`, `description_full`, `system_requirements`, `media_references`, `tags`, `genres`.
*   **`GamePricing`**: Эта сущность управляется и хранится преимущественно в **Catalog Service**. Developer Service предоставляет интерфейс для разработчиков, чтобы предлагать базовые цены и участвовать в промо-акциях, которые затем применяются и управляются через Catalog Service. В Developer Service поле `Game.draft_pricing` используется для подготовки этих данных.
    *   Поля: `product_id`, `region_code`, `currency_code`, `base_amount`, `discount_rules`.

*   **`DeveloperPayout` (Выплата Разработчику)**
    *   `id` (UUID): Уникальный идентификатор запроса на выплату.
    *   `developer_id` (UUID, FK to Developer): ID разработчика. Обязательность: Required.
    *   `amount_requested_minor` (BIGINT): Сумма запроса в минимальных единицах валюты. Обязательность: Required.
    *   `currency_code` (VARCHAR(3)): Код валюты. Обязательность: Required.
    *   `status` (VARCHAR(50)): Статус (`pending_request`, `processing`, `completed`, `failed`, `cancelled`). Обязательность: Required.
    *   `payment_method_details_snapshot` (JSONB): Снимок реквизитов на момент запроса (для истории и аудита). Обязательность: Required.
    *   `requested_at` (TIMESTAMPTZ), `processed_at` (TIMESTAMPTZ).
    *   `transaction_id_payment_service` (UUID): ID транзакции в Payment Service. Обязательность: Optional.
    *   `comment_developer` (TEXT): Комментарий от разработчика.
    *   `comment_admin` (TEXT): Комментарий от администратора/финансового отдела.

*   **`DeveloperAPIKey` (API Ключ Разработчика)**
    *   `id` (UUID): Уникальный идентификатор ключа.
    *   `developer_id` (UUID, FK to Developer): ID разработчика, которому принадлежит ключ. Обязательность: Required.
    *   `name` (VARCHAR(100)): Имя ключа, задаваемое разработчиком для идентификации. Пример: `ci_cd_pipeline_key`. Валидация: not null. Обязательность: Required.
    *   `prefix` (VARCHAR(8)): Первые несколько символов ключа (для отображения и быстрой идентификации). Пример: `dvk_`. Валидация: not null, unique. Обязательность: Required.
    *   `key_hash` (VARCHAR(255)): Хеш API ключа (сам ключ не хранится). Валидация: not null, unique. Обязательность: Required.
    *   `permissions` (JSONB): Список разрешений, связанных с этим API ключом (например, `upload_build`, `manage_metadata`). Обязательность: Required, default `[]`.
    *   `last_used_at` (TIMESTAMPTZ): Время последнего использования ключа. Обязательность: Optional.
    *   `expires_at` (TIMESTAMPTZ): Время истечения срока действия ключа. Обязательность: Optional (может быть бессрочным).
    *   `created_at` (TIMESTAMPTZ): Время создания. Обязательность: Required.
    *   `is_active` (BOOLEAN): Активен ли ключ. Обязательность: Required, default `true`.

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL
**ERD Диаграмма (обновленная):**
```mermaid
erDiagram
    DEVELOPERS {
        UUID id PK
        UUID owner_user_id FK "User(Auth)"
        VARCHAR company_name
        VARCHAR legal_entity_type
        VARCHAR tax_id
        VARCHAR contact_email
        VARCHAR status
        TIMESTAMPTZ created_at
    }
    DEVELOPER_TEAM_MEMBERS {
        UUID developer_id PK FK
        UUID user_id PK FK "User(Auth)"
        VARCHAR role_in_team
        TIMESTAMPTZ joined_at
    }
    GAMES {
        UUID id PK
        UUID developer_id FK
        VARCHAR title
        VARCHAR product_type
        VARCHAR status
        JSONB draft_metadata
        JSONB draft_pricing
        TIMESTAMPTZ created_at
        TIMESTAMPTZ published_at
    }
    GAME_VERSIONS {
        UUID id PK
        UUID game_id FK
        VARCHAR version_name
        VARCHAR status
        JSONB changelog
        VARCHAR build_s3_path
        VARCHAR build_platform
        TIMESTAMPTZ uploaded_at
    }
    DEVELOPER_PAYOUTS {
        UUID id PK
        UUID developer_id FK
        BIGINT amount_requested_minor
        VARCHAR currency_code
        VARCHAR status
        JSONB payment_method_snapshot
        TEXT comment_developer
        TEXT comment_admin
        UUID transaction_id_payment_service "nullable"
        TIMESTAMPTZ requested_at
        TIMESTAMPTZ processed_at
    }
    DEVELOPER_API_KEYS {
        UUID id PK
        UUID developer_id FK
        VARCHAR name
        VARCHAR prefix UK
        VARCHAR key_hash UK
        JSONB permissions
        BOOLEAN is_active
        TIMESTAMPTZ expires_at "nullable"
        TIMESTAMPTZ last_used_at "nullable"
        TIMESTAMPTZ created_at
    }

    DEVELOPERS ||--o{ DEVELOPER_TEAM_MEMBERS : "has"
    DEVELOPERS ||--o{ GAMES : "develops"
    DEVELOPERS ||--o{ DEVELOPER_PAYOUTS : "requests"
    DEVELOPERS ||--o{ DEVELOPER_API_KEYS : "owns"
    GAMES ||--o{ GAME_VERSIONS : "has"

    GAMES ..|> CATALOG_SERVICE_PRODUCTS : "Proposes/Views Metadata & Pricing via API"
    CATALOG_SERVICE_PRODUCTS {
        note "Managed by Catalog Service"
    }
```

**DDL (PostgreSQL - дополнения для `developer_payouts`, `developer_api_keys` и уточнение `games`):**
```sql
-- Добавление полей draft_metadata и draft_pricing в таблицу games, если она уже существует
-- ALTER TABLE games ADD COLUMN draft_metadata JSONB;
-- ALTER TABLE games ADD COLUMN draft_pricing JSONB;
-- COMMENT ON COLUMN games.draft_metadata IS 'Черновик метаданных продукта, редактируемый разработчиком (структура соответствует Catalog Service)';
-- COMMENT ON COLUMN games.draft_pricing IS 'Черновик предложений по ценам от разработчика (структура соответствует Catalog Service)';

CREATE TABLE developer_payouts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    developer_id UUID NOT NULL REFERENCES developers(id) ON DELETE CASCADE,
    amount_requested_minor BIGINT NOT NULL CHECK (amount_requested_minor > 0),
    currency_code VARCHAR(3) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_request' CHECK (status IN ('pending_request', 'processing', 'completed', 'failed', 'cancelled')),
    payment_method_details_snapshot JSONB NOT NULL,
    comment_developer TEXT,
    comment_admin TEXT,
    transaction_id_payment_service UUID,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ
);
CREATE INDEX idx_developer_payouts_developer_id ON developer_payouts(developer_id);
CREATE INDEX idx_developer_payouts_status ON developer_payouts(status);

CREATE TABLE developer_api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    developer_id UUID NOT NULL REFERENCES developers(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    prefix VARCHAR(8) NOT NULL UNIQUE,
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    permissions JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    expires_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_developer_api_keys_developer_id ON developer_api_keys(developer_id);
```

#### 4.2.2. S3-совместимое хранилище
*   **Роль:** Хранение бинарных файлов.
*   **Структура:**
    *   Билды игр: `s3://<bucket-name-game-builds>/<developer_id>/<game_id>/<version_id>/<platform>/<filename.zip>`
    *   Медиа-файлы: `s3://<bucket-name-game-media>/<developer_id>/<game_id>/media/<media_type>/<timestamp_filename.ext>`
    *   Документы для верификации разработчика: `s3://<bucket-name-developer-documents>/<developer_id>/verification/<document_type>/<filename.pdf>` (с строгими правами доступа).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Формат событий:** CloudEvents v1.0 JSON (согласно `../../../../project_api_standards.md`).
*   **Основной топик:** `com.platform.developer.events.v1`.

*   **`com.platform.developer.account.created.v1`**
    *   Описание: Аккаунт разработчика/издателя был успешно создан и ожидает верификации.
    *   `data` Payload:
        ```json
        {
          "developerId": "dev-uuid-abc",
          "ownerUserId": "user-uuid-owner",
          "companyName": "Моя Игровая Студия",
          "contactEmail": "dev@example.com",
          "status": "pending_verification",
          "creationTimestamp": "2024-03-18T10:00:00Z"
        }
        ```
    *   Потребители: Admin Service (для начала процесса верификации), Notification Service.
*   **`com.platform.developer.game.submitted.v1`**
    *   Описание: Разработчик отправил игру/версию на модерацию.
    *   `data` Payload:
        ```json
        {
          "developerId": "dev-uuid-abc",
          "gameId": "game-uuid-xyz",
          "versionId": "version-uuid-123",
          "versionName": "1.0.1",
          "submittedAt": "2024-03-15T14:30:00Z",
          "submittedByUserId": "user-uuid-dev-member"
        }
        ```
    *   Потребители: Admin Service (для начала процесса модерации).
*   **`com.platform.developer.game.published.v1`**
     *   Описание: Разработчик опубликовал одобренную игру/версию (или она была опубликована автоматически после одобрения).
    *   `data` Payload:
        ```json
        {
          "developerId": "dev-uuid-abc",
          "gameId": "game-uuid-xyz",
          "versionId": "version-uuid-123",
          "publishedAt": "2024-03-16T10:00:00Z",
          "publishedByUserId": "user-uuid-dev-member"
        }
        ```
    *   Потребители: Catalog Service, Download Service, Notification Service.
*   **`com.platform.developer.game.metadata.updated.v1`**
    *   Описание: Разработчик обновил метаданные для своего продукта (черновик сохранен в Developer Service, событие информирует о возможном последующем запросе на модерацию/публикацию в Catalog Service).
    *   `data` Payload:
        ```json
        {
          "developerId": "dev-uuid-abc",
          "gameId": "game-uuid-xyz",
          "updatedFields": ["descriptions", "system_requirements", "tags"],
          "updateTimestamp": "2024-03-18T11:00:00Z",
          "submittedByUserId": "user-uuid-dev-member"
        }
        ```
    *   Потребители: Catalog Service (для обновления своего представления, если это прямой пуш, или для информации), Admin Service (если требуется модерация изменений).
*   **`com.platform.developer.payout.requested.v1`**
    *   Описание: Разработчик запросил выплату средств.
    *   `data` Payload:
        ```json
        {
          "payoutRequestId": "payout-uuid-789",
          "developerId": "dev-uuid-abc",
          "amountMinorUnits": 5000000,
          "currencyCode": "RUB",
          "requestedAt": "2024-03-17T09:00:00Z",
          "paymentMethodId": "bank-account-uuid-123"
        }
        ```
    *   Потребители: Payment Service.

### 5.2. Потребляемые События (Consumed Events)
(Содержимое существующего раздела актуально, с коррекцией имен событий на формат `com.platform.*`).

## 6. Интеграции (Integrations)
(Содержимое существующего раздела актуально).

## 7. Конфигурация (Configuration)
(Содержимое существующего раздела YAML и описание переменных окружения в целом актуальны).

## 8. Обработка Ошибок (Error Handling)
(Содержимое существующего раздела актуально, форматы ошибок исправлены в разделе API).

## 9. Безопасность (Security)
(Содержимое существующего раздела в целом актуально).
*   **ФЗ-152 "О персональных данных":** Developer Service обрабатывает ПДн разработчиков (ФИО контактных лиц, email, телефон, юридические и банковские реквизиты ИП/самозанятых). Необходимо обеспечить шифрование этих данных при хранении и передаче, строгое управление доступом, получение согласий на обработку ПДн.
*   Ссылки на `../../../../project_security_standards.md` и `../../../../project_roles_and_permissions.md` актуальны.

## 10. Развертывание (Deployment)
(Содержимое существующего раздела актуально).

## 11. Мониторинг и Логирование (Logging and Monitoring)
(Содержимое существующего раздела актуально).

## 12. Нефункциональные Требования (NFRs)
(Содержимое существующего раздела актуально).

## 13. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API и Protobuf определения для gRPC API (если будут использоваться) поддерживаются в соответствующих репозиториях исходного кода сервиса и/или в централизованном репозитории `platform-protos`.
*   DDL схемы базы данных управляются через систему миграций (например, `golang-migrate/migrate`) и хранятся в репозитории исходного кода сервиса. Актуальные версии доступны во внутренней документации команды разработки и в GitOps репозитории.

## 14. Резервное Копирование и Восстановление (Backup and Recovery)

### 14.1. PostgreSQL (Данные аккаунтов разработчиков, черновики продуктов, запросы на выплаты и т.д.)
*   **Процедура резервного копирования:**
    *   **Логические бэкапы:** Ежедневный `pg_dump` для базы данных Developer Service.
    *   **Физические бэкапы (PITR):** Настроена непрерывная архивация WAL-сегментов. Базовый бэкап создается еженедельно.
    *   **Хранение:** Бэкапы и WAL-архивы хранятся в S3-совместимом хранилище с шифрованием и версионированием, в другом регионе. Срок хранения: полные логические бэкапы - 30 дней, WAL - 14 дней.
*   **Процедура восстановления:** Тестируется ежеквартально.
*   **RTO (Recovery Time Objective):** < 2 часов.
*   **RPO (Recovery Point Objective):** < 15 минут.
*   (Общие принципы см. `../../../../project_database_structure.md`).

### 14.2. S3-совместимое хранилище (Билды игр, медиа-активы, документы верификации)
*   **Процедура резервного копирования:**
    *   **Версионирование объектов:** Включено для всех бакетов.
    *   **Политики жизненного цикла (Lifecycle Policies):** Настроены для управления старыми версиями и, возможно, для перемещения неактивных билдов в более холодные классы хранения (если применимо).
    *   **Cross-Region Replication (CRR):** Настроена для бакетов с билдами игр и критически важными медиа-активами для обеспечения гео-резервирования. Для документов верификации также рекомендуется CRR.
*   **Процедура восстановления:**
    *   Восстановление отдельных объектов или версий из S3.
    *   В случае регионального сбоя – переключение на реплицированный бакет в другом регионе.
*   **RTO:** Зависит от объема данных и скорости S3, но обычно быстро для отдельных файлов. Полное восстановление всех данных может занять часы.
*   **RPO:** Близко к нулю при использовании версионирования и CRR (ограничено временем репликации S3).

### 14.3. Redis (Кэш, сессии Портала Разработчика)
*   **Стратегия:** Данные в Redis в основном являются кэшем или временными сессионными данными.
*   **Персистентность (опционально):** Может быть включена RDB-снапшотирование и/или AOF для ускорения восстановления после перезапуска Redis.
*   **Резервное копирование:** Не является критичным для большинства данных, так как они могут быть перестроены из PostgreSQL или пересозданы пользователем.
*   **RTO/RPO:** Неприменимо в контексте долгосрочного хранения данных. Восстановление функциональности кэша происходит по мере его заполнения.

### 14.4. Общая стратегия
*   Приоритет отдается восстановлению данных из PostgreSQL и S3.
*   Процедуры восстановления тестируются и документируются.
*   Мониторинг процессов резервного копирования.

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [Подача разработчиком новой игры на модерацию](../../../../project_workflows/game_submission_flow.md)
*   [Процесс выплат разработчикам] <!-- Workflow будет создан и описан в project_workflows/developer_payout_flow.md -->

---
*Этот документ является основной спецификацией для Developer Service и должен поддерживаться в актуальном состоянии.*
