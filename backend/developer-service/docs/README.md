# Спецификация Микросервиса: Developer Service (Сервис для Разработчиков)

**Версия:** 1.1 (адаптировано из предыдущей версии)
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Developer Service, предназначенного для разработчиков и издателей игр на платформе "Российский Аналог Steam".
*   **Роль в общей архитектуре платформы:** Developer Service предоставляет интерфейс (Портал Разработчика) и API для управления аккаунтами разработчиков, загрузки и обновления игр и их контента, управления метаданными и маркетинговыми материалами, доступа к аналитике продаж и использования продуктов, а также взаимодействия с финансовыми аспектами платформы (управление выплатами).
*   **Основные бизнес-задачи:** Обеспечение разработчиков инструментами для самостоятельной публикации и поддержки их продуктов (игр, DLC, ПО) на платформе, управление жизненным циклом этих продуктов, предоставление релевантной аналитики и финансовых отчетов.
*   Разработка сервиса должна вестись в соответствии с `CODING_STANDARDS.md`.

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
*   **Язык программирования:** Go (предпочтительно).
*   **Веб-фреймворк (REST API):** Echo или Gin (для Go).
*   **База данных:** PostgreSQL (для хранения структурированных данных: аккаунты разработчиков, метаданные игр, информация о финансах).
*   **Кэширование:** Redis (для кэширования часто запрашиваемых данных, сессий портала).
*   **Очереди сообщений:** Apache Kafka (для асинхронной обработки задач и публикации событий).
*   **Хранилище файлов (билды, медиа):** S3-совместимое объектное хранилище (например, MinIO).
*   **Инфраструктура:** Docker, Kubernetes.
*   **Фронтенд Портала Разработчика:** (Технологии фронтенда здесь не описываются, но Developer Service предоставляет для него API).
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `PACKAGE_STANDARDIZATION.md`.
*   *Примечание: Приведенные в последующих разделах примеры API и конфигураций основаны на предположении использования Go (Echo/Gin для REST), PostgreSQL, Redis и Kafka.*

### 1.4. Термины и Определения (Glossary)
*   **Разработчик (Developer Account):** Учетная запись компании или индивидуального разработчика/издателя.
*   **Продукт (Product/Game):** Игра, DLC или другой цифровой товар, управляемый разработчиком.
*   **Билд (Build):** Конкретная сборка исполняемых файлов продукта.
*   **Версия (Version):** Версия продукта, связанная с конкретным билдом и метаданными.
*   **Портал Разработчика (Developer Portal):** Веб-интерфейс, предоставляемый Developer Service для разработчиков.
*   Для других общих терминов см. `project_glossary.md`.

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
*   **Формат ответа об ошибке (стандартный):**
    ```json
    {
      "errors": [
        {
          "status": "4XX/5XX",
          "code": "ERROR_CODE_UPPER_SNAKE_CASE",
          "title": "Краткое описание ошибки на русском",
          "detail": "Полное описание ошибки с контекстом.",
          "source": { "pointer": "/data/attributes/field_name", "parameter": "query_param_name" } // Опционально
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
              "legal_entity_type": "ООО", // "ИП", "Самозанятый"
              "tax_id": "1234567890",
              "contact_email": "dev@example.com",
              "country_code": "RU"
              // ... другие юридические и контактные данные
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
              "product_type": "game" // game, dlc, software
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
              "status": "draft", // Начальный статус
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
              "file_size_bytes": 1073741824, // 1GB
              "content_type": "application/zip",
              "platform": "windows_x64" // или другая платформа
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
              "internal_build_id": "build-uuid-temp" // ID для отслеживания после загрузки
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
              "s3_path": "path/to/mygame_v1.0.1_windows_x64.zip", // Путь в S3 после загрузки
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
              "amount_minor_units": 5000000, // 50,000.00 RUB (в копейках)
              "currency_code": "RUB",
              "payment_method_id": "bank-account-uuid-123" // ID сохраненного способа выплаты
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Информация о созданном запросе на выплату)
    *   Требуемые права доступа: `developer_admin` или `finance_manager`.

### 3.2. gRPC API
*   Для межсервисного взаимодействия могут быть определены специфичные gRPC методы (например, для получения информации о разработчике или игре по запросу от Admin Service или Catalog Service).
*   Пример:
    *   `rpc GetDeveloperInfo(GetDeveloperInfoRequest) returns (GetDeveloperInfoResponse)`
    *   `rpc GetGamePublicationStatus(GetGamePublicationStatusRequest) returns (GetGamePublicationStatusResponse)`
*   TODO: Детализировать gRPC API по мере необходимости проектирования межсервисных взаимодействий.

### 3.3. WebSocket API
*   В настоящее время WebSocket API не планируется. Возможным вариантом использования может быть real-time оповещение разработчиков в их портале о важных событиях (например, завершение модерации, получение выплаты, завершение сборки билда), но это требует дополнительного проектирования и может быть реализовано через другие механизмы (например, Server-Sent Events или периодический опрос).

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности

*   **`Developer` (Разработчик/Издатель)**
    *   `id` (UUID): Уникальный идентификатор. Обязательность: Required.
    *   `user_id_owner` (UUID): ID пользователя-владельца из Auth Service. Обязательность: Required.
    *   `company_name` (VARCHAR(255)): Название компании или имя разработчика. Пример: "Моя Игровая Студия". Валидация: not null. Обязательность: Required.
    *   `legal_entity_type` (VARCHAR(50)): Тип юр. лица (ООО, ИП, Самозанятый). Обязательность: Required.
    *   `tax_id` (VARCHAR(50)): ИНН/аналог. Обязательность: Required.
    *   `contact_email` (VARCHAR(255)): Контактный email. Валидация: valid email. Обязательность: Required.
    *   `country_code` (VARCHAR(2)): Код страны регистрации. Пример: `RU`. Обязательность: Required.
    *   `status` (VARCHAR(50)): Статус аккаунта (`pending_verification`, `active`, `suspended`). Обязательность: Required.
    *   `created_at` (TIMESTAMPTZ), `updated_at` (TIMESTAMPTZ).

*   **`DeveloperTeamMember` (Член Команды Разработчика)**
    *   `developer_id` (UUID, FK to Developer): ID аккаунта разработчика. Обязательность: Required.
    *   `user_id` (UUID, FK to User in Auth Service): ID пользователя. Обязательность: Required.
    *   `role_in_team` (VARCHAR(50)): Роль в команде (`owner`, `admin`, `game_editor`, `finance_viewer`). Обязательность: Required.
    *   `joined_at` (TIMESTAMPTZ).

*   **`Game` (Продукт/Игра)**
    *   `id` (UUID): Уникальный идентификатор продукта. Обязательность: Required.
    *   `developer_id` (UUID, FK to Developer): ID разработчика. Обязательность: Required.
    *   `title` (VARCHAR(255)): Основное название продукта (может быть заменено локализованным из GameMetadata). Пример: "Моя Новая Супер Игра". Валидация: not null. Обязательность: Required.
    *   `product_type` (VARCHAR(50)): Тип (`game`, `dlc`, `software`). Обязательность: Required.
    *   `status` (VARCHAR(50)): Общий статус продукта (`draft`, `in_review`, `published`, `rejected`, `archived`). Обязательность: Required.
    *   `created_at` (TIMESTAMPTZ), `updated_at` (TIMESTAMPTZ).
    *   `published_at` (TIMESTAMPTZ): Дата публикации. Обязательность: Optional.

*   **`GameVersion` (Версия Продукта)**
    *   `id` (UUID): Уникальный идентификатор версии. Обязательность: Required.
    *   `game_id` (UUID, FK to Game): ID продукта. Обязательность: Required.
    *   `version_name` (VARCHAR(50)): Имя версии (например, `1.0.0`, `beta-2`). Валидация: not null. Обязательность: Required.
    *   `status` (VARCHAR(50)): Статус версии (`uploading`, `processing`, `failed_processing`, `ready_for_review`, `in_review`, `approved`, `rejected`, `live`, `deprecated`). Обязательность: Required.
    *   `changelog` (JSONB): Локализованный список изменений. Пример: `{"ru-RU": "Исправлены баги...", "en-US": "Bugfixes..."}`. Обязательность: Optional.
    *   `build_s3_path` (VARCHAR(1024)): Путь к файлу билда в S3. Обязательность: Optional (появляется после загрузки).
    *   `build_platform` (VARCHAR(50)): Платформа билда (`windows_x64`, `linux_x64`, `macos_arm64`). Обязательность: Optional.
    *   `uploaded_at` (TIMESTAMPTZ), `processed_at` (TIMESTAMPTZ).

*   **`GameMetadata` (Метаданные Продукта)**: Связана с `Game`. Содержит локализуемые поля, которые передаются в Catalog Service.
    *   `game_id` (UUID, PK, FK to Game).
    *   `language_code` (VARCHAR(10), PK): Код языка (например, `ru-RU`, `en-US`).
    *   `title` (VARCHAR(255)): Локализованное название.
    *   `description_short` (TEXT): Краткое описание.
    *   `description_full` (TEXT): Полное описание.
    *   `system_requirements` (JSONB): Системные требования для данной локали/платформы.
    *   `media_references` (JSONB): Ссылки на медиа-файлы (скриншоты, трейлеры) в S3 или CDN, специфичные для локали.
    *   `tags` (ARRAY of VARCHAR): Локализованные теги или ссылки на глобальные теги.
    *   `genres` (ARRAY of VARCHAR): Локализованные жанры или ссылки на глобальные жанры.

*   **`DeveloperPayout` (Выплата Разработчику)**
    *   `id` (UUID): Уникальный идентификатор запроса на выплату.
    *   `developer_id` (UUID, FK to Developer): ID разработчика. Обязательность: Required.
    *   `amount_requested_minor` (BIGINT): Сумма запроса в минимальных единицах валюты. Обязательность: Required.
    *   `currency_code` (VARCHAR(3)): Код валюты. Обязательность: Required.
    *   `status` (VARCHAR(50)): Статус (`pending_request`, `processing`, `completed`, `failed`, `cancelled`). Обязательность: Required.
    *   `payment_method_details_snapshot` (JSONB): Снимок реквизитов на момент запроса. Обязательность: Required.
    *   `requested_at` (TIMESTAMPTZ), `processed_at` (TIMESTAMPTZ).
    *   `transaction_id_payment_service` (UUID): ID транзакции в Payment Service. Обязательность: Optional.

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL

**ERD Диаграмма (ключевые таблицы):**
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
    GAME_METADATA_LOCALIZED { # Пример локализованных метаданных
        UUID game_id PK FK
        VARCHAR language_code PK
        VARCHAR title
        TEXT description_short
        TEXT description_full
    }
    DEVELOPER_PAYOUTS {
        UUID id PK
        UUID developer_id FK
        BIGINT amount_requested_minor
        VARCHAR currency_code
        VARCHAR status
        JSONB payment_method_snapshot
        TIMESTAMPTZ requested_at
        TIMESTAMPTZ processed_at
    }
    DEVELOPER_API_KEYS {
        UUID id PK
        UUID developer_id FK
        VARCHAR name
        VARCHAR prefix UK
        VARCHAR key_hash UK
        TIMESTAMPTZ expires_at
        TIMESTAMPTZ last_used_at
    }

    DEVELOPERS ||--o{ DEVELOPER_TEAM_MEMBERS : "has"
    DEVELOPERS ||--o{ GAMES : "develops"
    DEVELOPERS ||--o{ DEVELOPER_PAYOUTS : "requests"
    DEVELOPERS ||--o{ DEVELOPER_API_KEYS : "owns"
    GAMES ||--o{ GAME_VERSIONS : "has"
    GAMES ||--o{ GAME_METADATA_LOCALIZED : "has_localized"
```

**DDL (PostgreSQL - примеры ключевых таблиц):**
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE developers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_user_id UUID NOT NULL, -- Ссылка на User ID из Auth Service
    company_name VARCHAR(255) NOT NULL,
    legal_entity_type VARCHAR(50) NOT NULL,
    tax_id VARCHAR(50) NOT NULL UNIQUE,
    contact_email VARCHAR(255) NOT NULL UNIQUE,
    country_code VARCHAR(2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_verification', -- pending_verification, active, suspended
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_developer_owner FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE RESTRICT -- Предполагается наличие таблицы users из Auth
);

CREATE TABLE developer_team_members (
    developer_id UUID NOT NULL REFERENCES developers(id) ON DELETE CASCADE,
    user_id UUID NOT NULL, -- User ID из Auth Service
    role_in_team VARCHAR(50) NOT NULL, -- owner, admin, editor, viewer, finance_manager
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (developer_id, user_id),
    CONSTRAINT fk_team_member_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    developer_id UUID NOT NULL REFERENCES developers(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL, -- Основное нелокализованное название или ключ для локализации
    product_type VARCHAR(50) NOT NULL DEFAULT 'game', -- game, dlc, software
    status VARCHAR(50) NOT NULL DEFAULT 'draft', -- draft, in_review, published, rejected, archived
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ
);
CREATE INDEX idx_games_developer_id ON games(developer_id);

CREATE TABLE game_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    version_name VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'uploading', -- uploading, processing, ready_for_review, in_review, approved, rejected, live, deprecated
    changelog JSONB, -- {"ru-RU": "...", "en-US": "..."}
    build_s3_path VARCHAR(1024),
    build_platform VARCHAR(50), -- windows_x64, linux_x64, etc.
    file_size_bytes BIGINT,
    file_hash_sha256 VARCHAR(64),
    uploaded_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (game_id, version_name)
);
CREATE INDEX idx_game_versions_game_id ON game_versions(game_id);

-- TODO: Добавить DDL для game_metadata_localized, game_pricing (может быть ссылкой на Catalog), developer_payouts, developer_api_keys.
```

#### 4.2.2. S3-совместимое хранилище
*   **Роль:** Хранение бинарных файлов:
    *   Билды игр и ПО.
    *   Медиа-файлы (скриншоты, трейлеры, обложки, маркетинговые материалы), загружаемые разработчиками.
    *   Возможно, SDK и документация для разработчиков.
*   **Структура (примерная):**
    *   `s3://<bucket-game-builds>/<developer_id>/<game_id>/<version_id>/<platform>/<filename>`
    *   `s3://<bucket-game-media>/<developer_id>/<game_id>/media/<media_type>/<filename>`

### 4.2.3. Redis
*   **Роль:**
    *   Кэширование часто запрашиваемых данных из PostgreSQL (например, профиль разработчика, список игр).
    *   Хранение временных данных (например, pre-signed URL для загрузки в S3, токены сессий Портала Разработчика, если он stateful).
    *   Счетчики для rate limiting API разработчика.

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Apache Kafka.
*   **Формат событий:** CloudEvents v1.0 (JSON encoding), согласно `project_api_standards.md`.
*   **Топики (примеры):** `developer.events.v1`. (Может быть разделен).

*   **`developer.game.submitted.v1`**
    *   Описание: Разработчик отправил игру/версию на модерацию.
    *   Пример Payload (`data` секция CloudEvent):
        ```json
        {
          "developer_id": "dev-uuid-abc",
          "game_id": "game-uuid-xyz",
          "version_id": "version-uuid-123",
          "version_name": "1.0.1",
          "submitted_at": "2024-03-15T14:30:00Z",
          "submitted_by_user_id": "user-uuid-dev-member"
        }
        ```
    *   Потребители: Admin Service (для начала процесса модерации).
*   **`developer.game.published.v1`**
    *   Описание: Разработчик опубликовал одобренную игру/версию (или она была опубликована автоматически после одобрения).
    *   Пример Payload:
        ```json
        {
          "developer_id": "dev-uuid-abc",
          "game_id": "game-uuid-xyz",
          "version_id": "version-uuid-123",
          "published_at": "2024-03-16T10:00:00Z",
          "published_by_user_id": "user-uuid-dev-member"
        }
        ```
    *   Потребители: Catalog Service, Download Service, Notification Service.
*   **`developer.payout.requested.v1`**
    *   Описание: Разработчик запросил выплату средств.
    *   Пример Payload:
        ```json
        {
          "payout_request_id": "payout-uuid-789",
          "developer_id": "dev-uuid-abc",
          "amount_minor_units": 5000000,
          "currency_code": "RUB",
          "requested_at": "2024-03-17T09:00:00Z",
          "payment_method_id": "bank-account-uuid-123"
        }
        ```
    *   Потребители: Payment Service.
*   TODO: Добавить другие события, например, `developer.account.created.v1`, `developer.game.metadata.updated.v1`.

### 5.2. Потребляемые События (Consumed Events)

*   **`admin.game.moderation.approved.v1`** (от Admin Service)
    *   Описание: Игра/версия одобрена модерацией.
    *   Ожидаемый Payload:
        ```json
        {
          "game_id": "game-uuid-xyz",
          "version_id": "version-uuid-123",
          "moderator_id": "admin-user-uuid-456",
          "decision_details_url": "https://admin.example.com/moderation/log/log-uuid-abc",
          "approved_at": "2024-03-15T18:00:00Z"
        }
        ```
    *   Логика обработки: Обновить статус `GameVersion` на `approved`. Уведомить разработчика через Notification Service (или опубликовать событие для Notification Service).
*   **`admin.game.moderation.rejected.v1`** (от Admin Service)
    *   Описание: Игра/версия отклонена модерацией.
    *   Ожидаемый Payload:
        ```json
        {
          "game_id": "game-uuid-xyz",
          "version_id": "version-uuid-123",
          "moderator_id": "admin-user-uuid-456",
          "reasons": [ {"code": "VIOLATION_CODE_1", "comment": "Недопустимый контент в описании."} ],
          "rejected_at": "2024-03-15T19:00:00Z"
        }
        ```
    *   Логика обработки: Обновить статус `GameVersion` на `rejected`. Сохранить причины отклонения. Уведомить разработчика.
*   **`payment.payout.status.changed.v1`** (от Payment Service)
    *   Описание: Статус выплаты разработчику изменился.
    *   Ожидаемый Payload:
        ```json
        {
          "payout_request_id": "payout-uuid-789", // ID из Developer Service
          "payment_service_transaction_id": "txn-payment-uuid-qwerty",
          "new_status": "completed", // "processing", "failed"
          "processed_at": "2024-03-18T10:00:00Z",
          "failure_reason_code": null, // или код ошибки, если failed
          "failure_reason_message": null
        }
        ```
    *   Логика обработки: Обновить статус `DeveloperPayout`. Уведомить разработчика.
*   **`analytics.report.ready.v1`** (от Analytics Service)
    *   Описание: Отчет по аналитике для разработчика готов.
    *   Ожидаемый Payload:
        ```json
        {
          "report_id": "report-instance-uuid-111",
          "developer_id": "dev-uuid-abc",
          "report_type": "monthly_sales", // "game_dau_mau", etc.
          "game_id": "game-uuid-xyz", // опционально, если отчет по конкретной игре
          "period_start": "2024-02-01",
          "period_end": "2024-02-29",
          "download_url_or_reference": "s3://analytics-reports/dev-uuid-abc/monthly_sales_202402.csv",
          "generated_at": "2024-03-10T08:00:00Z"
        }
        ```
    *   Логика обработки: Сохранить ссылку на отчет или метаданные. Уведомить разработчика о доступности нового отчета в Портале.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Auth Service:** Для аутентификации пользователей-разработчиков и управления API ключами разработчиков.
*   **Account Service:** Для получения базовой информации о пользователях, входящих в команду разработчика.
*   **Catalog Service:** Developer Service инициирует создание/обновление продуктов в Catalog Service через его API или путем публикации событий, которые Catalog Service потребляет.
*   **Download Service:** Developer Service передает информацию о загруженных билдах и их версиях в Download Service для организации их скачивания пользователями.
*   **Payment Service:** Для обработки запросов на выплаты разработчикам и получения статусов этих выплат.
*   **Admin Service:** Для передачи продуктов на модерацию и получения результатов модерации.
*   **Analytics Service:** Для отображения статистики по продажам, использованию продуктов и другим метрикам в Портале Разработчика. Developer Service запрашивает данные у Analytics Service.
*   **Notification Service:** Для отправки уведомлений разработчикам о различных событиях (статус модерации, финансовые операции, системные сообщения).

### 6.2. Внешние Системы
*   **S3-совместимое объектное хранилище:** Для хранения билдов игр, медиа-файлов (скриншоты, трейлеры, арты), SDK и другой документации, загружаемой разработчиками.

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `DEVELOPER_SERVICE_HTTP_PORT`: Порт для REST API (например, `8081`).
*   `DEVELOPER_SERVICE_GRPC_PORT`: Порт для gRPC API (например, `9091`), если используется.
*   `POSTGRES_DSN`: Строка подключения к PostgreSQL.
*   `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB_DEVELOPER`: Параметры Redis.
*   `KAFKA_BROKERS`: Список брокеров Kafka.
*   `KAFKA_TOPIC_DEVELOPER_EVENTS`: Топик для публикуемых событий.
*   `S3_ENDPOINT`, `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY`, `S3_BUCKET_GAME_BUILDS`, `S3_BUCKET_GAME_MEDIA`, `S3_REGION`, `S3_USE_SSL`: Параметры S3.
*   `LOG_LEVEL`: Уровень логирования.
*   `AUTH_SERVICE_GRPC_ADDR`: Адрес gRPC Auth Service.
*   `CATALOG_SERVICE_GRPC_ADDR` (или REST API URL): Адрес Catalog Service.
*   `PAYMENT_SERVICE_GRPC_ADDR`: Адрес Payment Service.
*   `ADMIN_SERVICE_GRPC_ADDR`: Адрес Admin Service.
*   `ANALYTICS_SERVICE_API_URL`: URL API Analytics Service.
*   `NOTIFICATION_SERVICE_GRPC_ADDR`: Адрес Notification Service.
*   `MAX_BUILD_FILE_SIZE_BYTES`: Максимальный размер файла билда.
*   `MAX_MEDIA_FILE_SIZE_BYTES`: Максимальный размер медиа-файла.
*   `TEMPORARY_UPLOAD_DIR`: Временная директория для загрузок.
*   `DEVELOPER_PORTAL_BASE_URL`: Базовый URL портала разработчика (для генерации ссылок в уведомлениях).
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`: Эндпоинт Jaeger.

### 7.2. Файлы Конфигурации (если применимо)
*   **`configs/developer_config.yaml`**: Может использоваться для более сложных или менее часто изменяемых настроек.
    ```yaml
    server:
      http_port: ${DEVELOPER_SERVICE_HTTP_PORT:-8081}
      read_timeout_seconds: 30
      write_timeout_seconds: 300 # Для загрузки файлов

    file_storage:
      s3:
        game_builds_bucket: ${S3_BUCKET_GAME_BUILDS}
        game_media_bucket: ${S3_BUCKET_GAME_MEDIA}
        max_build_size_gb: 100 # Переопределяет ENV VAR, если нужно
        allowed_build_mime_types: ["application/zip", "application/octet-stream"]
        allowed_media_mime_types: ["image/jpeg", "image/png", "video/mp4"]
      temporary_upload_dir: ${TEMPORARY_UPLOAD_DIR:-/tmp/dev_uploads}

    payouts:
      min_payout_amount_rub: 100000 # 1000.00 RUB
      default_payout_schedule: "monthly" # monthly, weekly, on_request
      # Ограничения или правила для выплат

    analytics_integration:
      default_report_period_days: 30

    # Настройки для различных этапов публикации продукта
    publication_workflow:
      require_all_media_types: true
      min_screenshots_count: 3
    ```

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   REST API: Используются стандартные коды состояния HTTP. Тело ответа об ошибке соответствует формату, определенному в `project_api_standards.md` (см. секцию 3.1).
*   gRPC API: Используются стандартные коды состояния gRPC.
*   Все ошибки логируются с `trace_id`.

### 8.2. Распространенные Коды Ошибок (для REST API)
*   **`400 Bad Request` (`VALIDATION_ERROR`)**: Некорректные входные данные (например, неверный формат, отсутствуют обязательные поля).
*   **`401 Unauthorized` (`UNAUTHENTICATED`)**: Ошибка аутентификации (требуется логин).
*   **`403 Forbidden` (`PERMISSION_DENIED`)**: Недостаточно прав для выполнения операции (например, попытка редактировать чужую игру, или роль в команде не позволяет).
*   **`404 Not Found` (`RESOURCE_NOT_FOUND`)**: Запрашиваемый ресурс не найден (игра, версия, аккаунт разработчика).
*   **`409 Conflict` (`RESOURCE_ALREADY_EXISTS`)**: Попытка создания ресурса, который уже существует (например, игра с таким же названием у этого разработчика).
*   **`413 Payload Too Large` (`FILE_TOO_LARGE`)**: Загружаемый файл превышает допустимый размер.
*   **`422 Unprocessable Entity` (`BUSINESS_LOGIC_ERROR`)**: Запрос корректен, но нарушает бизнес-правила (например, попытка опубликовать игру без билдов).
*   **`500 Internal Server Error` (`INTERNAL_ERROR`)**: Внутренняя ошибка сервера.
*   **`503 Service Unavailable` (`SERVICE_UNAVAILABLE`)**: Сервис временно недоступен или зависимость недоступна.

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   Все запросы к API Developer Service (кроме, возможно, некоторых публичных информационных эндпоинтов, если такие будут) требуют JWT аутентификации пользователя (разработчика). Аутентификация выполняется через Auth Service.

### 9.2. Авторизация
*   Используется модель RBAC внутри команды разработчика. Пользователь, привязанный к `Developer` аккаунту, может иметь разные роли (например, `owner`, `admin`, `editor`, `viewer`, `finance_manager`), которые определяют доступ к различным функциям и данным в рамках этого `Developer` аккаунта.
*   Проверка принадлежности ресурса (например, игры) к аутентифицированному `Developer` аккаунту обязательна для всех операций CRUD.

### 9.3. Защита Данных
*   Шифрование конфиденциальных данных, таких как финансовые реквизиты (если хранятся) и API ключи разработчиков (хранятся только хэши).
*   HTTPS для всех коммуникаций.
*   Безопасная загрузка файлов: проверка типов файлов, ограничение размера, возможно, антивирусная проверка загружаемых билдов на стороне S3 или после загрузки.
*   Защита от несанкционированного доступа к билдам и медиа-файлам в S3 (например, через pre-signed URLs с коротким сроком жизни для загрузки/скачивания SDK).

### 9.4. Управление Секретами
*   API ключи, используемые самим Developer Service для доступа к другим сервисам, должны храниться в Kubernetes Secrets или HashiCorp Vault.
*   API ключи, генерируемые для разработчиков, должны храниться в виде хэшей в базе данных. Само значение ключа показывается разработчику только один раз при создании.

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Стандартный многоэтапный Dockerfile для Go-приложения.
*   **Kubernetes манифесты/Helm-чарты:** Для управления развертыванием, сервисами, конфигурациями, секретами и т.д.
*   (Ссылка на `project_deployment_standards.md`).

### 10.2. Зависимости при Развертывании
*   PostgreSQL, S3-совместимое хранилище, Kafka (или RabbitMQ), Redis.
*   Доступность Auth Service, Account Service, Payment Service, Catalog Service, Download Service, Analytics Service, Admin Service, Notification Service.

### 10.3. CI/CD
*   Автоматизированная сборка, модульное и интеграционное тестирование.
*   Статический анализ кода, сканирование на уязвимости.
*   Сборка Docker-образа и публикация в приватный registry.
*   Развертывание в различные окружения (dev, staging, production) с использованием GitOps-подхода.

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат:** Структурированные логи в формате JSON (например, с использованием Zap).
*   **Ключевые события:** Все запросы к API, операции с файлами, изменения статусов продуктов, финансовые операции, ошибки, события безопасности.
*   **Интеграция:** С централизованной системой логирования (ELK Stack, Loki/Grafana).
*   (Ссылка на `project_observability_standards.md`).

### 11.2. Мониторинг
*   **Метрики (Prometheus):**
    *   Количество запросов к API (по эндпоинтам, методам, статусам ответа).
    *   Длительность обработки запросов.
    *   Количество и размер загруженных/скачанных файлов.
    *   Количество активных разработчиков/сессий в портале.
    *   Ошибки при взаимодействии с БД, S3, Kafka, другими сервисами.
    *   Стандартные метрики Go-приложения.
*   **Дашборды (Grafana):** Для визуализации метрик производительности, ошибок, использования ресурсов.
*   **Алертинг (AlertManager):** Для критических ошибок, недоступности зависимостей, аномального поведения (например, резкий рост ошибок загрузки файлов).
*   (Ссылка на `project_observability_standards.md`).

### 11.3. Трассировка
*   **Интеграция:** OpenTelemetry SDK для Go. Экспорт трейсов в Jaeger или Tempo.
*   Трассировка всех входящих API запросов и исходящих вызовов к другим микросервисам и базам данных.
*   (Ссылка на `project_observability_standards.md`).

## 12. Нефункциональные Требования (NFRs)
*   **Безопасность:**
    *   Надежная защита аккаунтов разработчиков и их интеллектуальной собственности (билды игр, исходный код, если применимо).
    *   Соответствие OWASP Top 10 для веб-портала.
    *   Безопасное хранение API ключей и финансовых данных.
*   **Производительность:**
    *   Время отклика API Портала Разработчика для типичных операций: P95 < 500 мс.
    *   Скорость загрузки билдов: должна быть ограничена только пропускной способностью сети и S3 (цель: поддержка загрузки файлов до 100GB+ без таймаутов на стороне сервиса).
    *   Время обработки уведомления о завершении загрузки билда: < 1 минуты.
*   **Масштабируемость:**
    *   Горизонтальное масштабирование для поддержки до 10,000 активных разработчиков и управления до 50,000 продуктов.
    *   Масштабируемость S3-хранилища для петабайтов данных (билды, медиа).
*   **Надежность:**
    *   Доступность Портала Разработчика и API: >= 99.9%.
    *   Сохранность загруженных билдов и метаданных: отсутствие потерь данных. RPO < 1 часа.
    *   RTO < 2 часов для критически важных функций (например, возможность загрузить билд или обновить статус игры).
*   **Удобство использования (Usability):** Интуитивно понятный и хорошо документированный Портал Разработчика и API.

## 13. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API и Protobuf определения для gRPC API (если будут использоваться) будут храниться в соответствующих репозиториях или артефактах CI/CD.
*   Полные DDL схемы базы данных будут поддерживаться в актуальном состоянии в системе миграций.
*   TODO: Добавить ссылки на репозиторий с OpenAPI/Protobuf и на систему управления миграциями БД, когда они будут определены.

---
*Этот документ является основной спецификацией для Developer Service и должен поддерживаться в актуальном состоянии.*

## 14. Связанные Рабочие Процессы (Related Workflows)
*   [Developer Submits a New Game for Moderation](../../../project_workflows/game_submission_flow.md)
