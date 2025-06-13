# Спецификация Микросервиса: Library Service (Сервис Библиотеки Пользователя)

**Версия:** 1.2
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение сервиса:** Library Service предназначен для управления личными библиотеками игр и приложений пользователей платформы "Российский Аналог Steam". Это включает в себя отслеживание принадлежащих пользователю продуктов, их статуса установки, игрового времени, достижений, управление списком желаемого и синхронизацию игровых сохранений в облаке.
*   **Роль в общей архитектуре платформы:** Является центральным компонентом для хранения и управления данными, связанными с игровой активностью и владением цифровыми продуктами конкретного пользователя. Предоставляет эту информацию как самому пользователю через клиентские приложения, так и другим микросервисам.
*   **Основные бизнес-задачи:**
    *   Обеспечение пользователям доступа к приобретенным играм и приложениям.
    *   Отслеживание и отображение игрового времени и прогресса достижений.
    *   Предоставление функционала списка желаемого.
    *   Обеспечение надежной синхронизации игровых сохранений между устройствами пользователя.
    *   Управление пользовательскими настройками для игр.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Управление библиотекой пользователя:** Добавление игр в библиотеку после покупки/активации, скрытие игр, фильтрация и сортировка списка игр, создание пользовательских категорий/коллекций, отслеживание статуса установки игры (интеграция с Download Service). Поддержка концепции "семейного доступа" (Family Sharing) для предоставления временного доступа к играм другим пользователям.
*   **Отслеживание игрового времени:** Автоматическая регистрация начала и окончания игровых сессий, подсчет общего и сессионного игрового времени для каждого продукта, отображение статистики последней активности, синхронизация игрового времени между устройствами.
*   **Управление достижениями:** Регистрация полученных пользователем достижений (анлоков), хранение прогресса по частично выполненным достижениям, отображение списка всех доступных и полученных достижений для игры, глобальная статистика по достижениям (редкость).
*   **Управление списком желаемого (Wishlist):** Добавление и удаление продуктов из списка желаемого, возможность приоритизации элементов, получение уведомлений о скидках на игры в списке (через Notification Service).
*   **Синхронизация игровых сохранений (Cloud Saves):** Загрузка файлов сохранений в облачное S3-хранилище, скачивание сохранений на другие устройства пользователя, разрешение конфликтов версий сохранений, версионирование (опционально), автоматическая и ручная синхронизация, управление квотами на облачное хранилище для сохранений.
*   **Настройки игр:** Хранение и синхронизация пользовательских настроек для конкретных игр (например, настройки графики, управления), если игра поддерживает такую интеграцию.

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+, согласно `../../../../project_technology_stack.md`).
*   **Веб-фреймворк (REST API):** Echo (`github.com/labstack/echo/v4`) или Gin (`github.com/gin-gonic/gin`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **RPC фреймворк (gRPC):** `google.golang.org/grpc` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **WebSocket:** (например, `github.com/gorilla/websocket`) для real-time уведомлений.
*   **База данных (основная):** PostgreSQL (версия 15+) для хранения информации о библиотеках, достижениях, списках желаемого, метаданных сохранений. Драйвер: GORM (`gorm.io/gorm`) с `gorm.io/driver/postgres` или `pgx` (`github.com/jackc/pgx/v5`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Кэширование/Оперативные данные:** Redis (версия 7.0+) для кэширования часто запрашиваемых данных и оперативных данных (например, текущие игровые сессии). Клиент: `go-redis/redis` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Облачное хранилище (для сохранений):** S3-совместимое объектное хранилище (например, MinIO, Yandex Object Storage) (согласно `../../../../project_technology_stack.md`).
*   **Брокер сообщений:** Apache Kafka (клиент `github.com/confluentinc/confluent-kafka-go` или `github.com/segmentio/kafka-go`, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Логирование:** Zap (`go.uber.org/zap`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Мониторинг/Трассировка:** OpenTelemetry SDK, Prometheus client (`github.com/prometheus/client_golang`). (согласно `../../../../project_observability_standards.md`).
*   **Инфраструктура:** Docker, Kubernetes, Helm.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   **Библиотека (Library):** Персональная коллекция игр и других цифровых продуктов, принадлежащих пользователю или доступных ему (например, через Family Sharing).
*   **Игровое время (Playtime):** Общее время, проведенное пользователем в конкретной игре.
*   **Достижение (Achievement):** Виртуальная награда, выдаваемая пользователю за выполнение определенных условий или задач в игре.
*   **Список желаемого (Wishlist):** Персональный список продуктов, которые пользователь хочет приобрести в будущем.
*   **Игровое сохранение (Savegame/Cloud Save):** Файл или набор файлов, содержащих прогресс пользователя в игре, который может быть синхронизирован с облачным хранилищем.
*   **Семейный доступ (Family Sharing):** Функция, позволяющая делиться играми из своей библиотеки с членами семьи или близкими друзьями.
*   Для других общих терминов см. `../../../../project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Library Service будет реализован с использованием принципов Чистой Архитектуры (Clean Architecture) для достижения слабой связанности, высокой тестируемости и гибкости системы.
*   Сервис будет состоять из нескольких ключевых модулей, отвечающих за отдельные аспекты функциональности: управление библиотекой, отслеживание игрового времени, управление достижениями, список желаемого и синхронизация сохранений.

**Диаграмма Архитектуры (Clean Architecture):**
```mermaid
graph TD
    subgraph User Clients & Other Services
        UserClient[Клиент Пользователя (Веб, Десктоп)]
        APIGateway[API Gateway]
        OtherInternalServices[Другие Микросервисы (Catalog, Download, Payment)]
    end

    subgraph Library Service
        direction TB

        subgraph PresentationLayer [Presentation Layer (Адаптеры Транспорта)]
            REST_API[REST API (Echo/Gin) - для клиента]
            GRPC_API[gRPC API - для других сервисов]
            WebSocket_Handler[WebSocket Handler - для real-time уведомлений]
            KafkaConsumers[Kafka Consumers - для входящих событий]
        end

        subgraph ApplicationLayer [Application Layer (Сценарии Использования)]
            LibraryAppSvc[Управление Библиотекой]
            PlaytimeAppSvc[Учет Игрового Времени]
            AchievementAppSvc[Управление Достижениями]
            WishlistAppSvc[Управление Списком Желаемого]
            SavegameAppSvc[Синхронизация Сохранений]
            GameSettingsAppSvc[Настройки Игр Пользователя]
        end

        subgraph DomainLayer [Domain Layer (Бизнес-логика и Сущности)]
            Entities[Сущности (UserLibraryItem, PlaytimeSession, UserAchievement, WishlistItem, Savegame)]
            Aggregates[Агрегаты (UserLibrary, UserGameProfile)]
            DomainEvents[Доменные События (GameAddedToLibrary, AchievementUnlocked)]
            RepositoryIntf[Интерфейсы Репозиториев]
            DomainServices[Доменные Сервисы (PlaytimeCalculator, AchievementUnlocker)]
        end

        subgraph InfrastructureLayer [Infrastructure Layer (Внешние Зависимости)]
            PostgresAdapter[Адаптер PostgreSQL]
            RedisAdapter[Адаптер Redis (Кэш, Сессии)]
            S3Adapter[Адаптер S3-хранилища (Сохранения)]
            KafkaProducer[Продюсер Kafka (События)]
            InternalServiceClients[Клиенты других микросервисов (Catalog, Download, Auth)]
            Config[Конфигурация (Viper)]
            Logging[Логирование (Zap)]
        end

        PresentationLayer --> ApplicationLayer
        ApplicationLayer --> DomainLayer
        ApplicationLayer --> InfrastructureLayer
        DomainLayer ----> RepositoryIntf
        InfrastructureLayer -- Implements --> RepositoryIntf
    end

    UserClient -- HTTP/WebSocket --> APIGateway
    APIGateway -- HTTP/WebSocket --> REST_API
    APIGateway -- HTTP/WebSocket --> WebSocket_Handler
    OtherInternalServices -- gRPC --> GRPC_API
    OtherInternalServices -- Kafka Events --> KafkaConsumers

    PostgresAdapter --> DB[(PostgreSQL)]
    RedisAdapter --> Cache[(Redis)]
    S3Adapter --> S3[(S3 Хранилище)]
    KafkaProducer --> KafkaBroker[Kafka Broker]
    InternalServiceClients --> ExtServices[Внешние gRPC Сервисы]


    classDef layer_boundary fill:#f9f9f9,stroke:#333,stroke-width:2px,color:#333
    classDef component_major fill:#e6f0ff,stroke:#007bff,color:#000
    classDef component_minor fill:#d4edda,stroke:#28a745,color:#000
    classDef datastore fill:#f8d7da,stroke:#dc3545,color:#000

    class PresentationLayer,ApplicationLayer,DomainLayer,InfrastructureLayer layer_boundary
    class REST_API,GRPC_API,WebSocket_Handler,KafkaConsumers,LibraryAppSvc,PlaytimeAppSvc,AchievementAppSvc,WishlistAppSvc,SavegameAppSvc,GameSettingsAppSvc,Entities,Aggregates,DomainEvents,RepositoryIntf,DomainServices component_major
    class PostgresAdapter,RedisAdapter,S3Adapter,KafkaProducer,InternalServiceClients,Config,Logging component_minor
    class DB,Cache,S3,KafkaBroker,ExtServices datastore
```

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления / Транспортный слой)
*   **Ответственность:** Обработка входящих REST API запросов от клиентских приложений (через API Gateway), gRPC запросов от других микросервисов, а также асинхронных сообщений из Kafka. Управление WebSocket соединениями для real-time уведомлений клиентам. Валидация данных запроса (DTO), вызов соответствующей логики в Application Layer.
*   **Ключевые компоненты/модули:** HTTP хендлеры (Echo/Gin), gRPC серверные реализации, WebSocket хендлеры, обработчики Kafka сообщений, DTO для запросов/ответов.

#### 2.2.2. Application Layer (Прикладной Слой / Сервисный слой / Use Case Layer)
*   **Ответственность:** Реализация сценариев использования системы, связанных с управлением библиотекой, игровым временем, достижениями, списком желаемого и сохранениями. Координирует взаимодействие между Domain Layer и Infrastructure Layer. Не содержит бизнес-правил напрямую, а делегирует их Domain Layer.
*   **Ключевые компоненты/модули:** Сервисы сценариев использования (например, `UserLibraryApplicationService`, `PlaytimeTrackingService`, `UserAchievementService`, `WishlistManagementService`, `SavegameSyncService`, `GameSettingsApplicationService`).

#### 2.2.3. Domain Layer (Доменный Слой)
*   **Ответственность:** Содержит бизнес-сущности, агрегаты, доменные события и бизнес-правила, специфичные для управления библиотекой пользователя и связанной с ней информацией. Этот слой не зависит от деталей реализации других слоев.
*   **Ключевые компоненты/модули:**
    *   **Entities (Сущности):** `UserLibraryItem`, `PlaytimeSession`, `UserAchievement`, `WishlistItem`, `SavegameMetadata`, `GameSpecificSettings`, `UserGameCategory`, `FamilySharingLink`.
    *   **Aggregates (Агрегаты):** Например, `UserLibrary` (включающий все `UserLibraryItem` пользователя), `UserGameProfile` (включающий `PlaytimeSession` и `UserAchievement` для конкретной игры пользователя).
    *   **Value Objects (Объекты-значения):** `GameDuration`, `AchievementProgress`.
    *   **Domain Services:** Сервисы, инкапсулирующие доменную логику, не принадлежащую одной сущности (например, сервис для проверки условий разблокировки достижений, если они зависят от нескольких сущностей).
    *   **Domain Events:** `GameAddedToLibraryEvent`, `PlaytimeSessionStartedEvent`, `AchievementUnlockedEvent`, `SavegameUploadedEvent`.
    *   **Repository Interfaces:** Интерфейсы, определяющие контракты для сохранения и извлечения агрегатов и сущностей.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   **Ответственность:** Реализация интерфейсов репозиториев для работы с PostgreSQL и Redis. Взаимодействие с S3-совместимым хранилищем для файлов сохранений. Публикация доменных событий в Kafka. Клиенты для взаимодействия с другими микросервисами (Catalog, Auth, Download, Notification, Payment, Account).
*   **Ключевые компоненты/модули:** Реализации репозиториев для PostgreSQL, Redis клиент, S3 клиент, Kafka продюсеры и консьюмеры (если сервис также потребляет события напрямую для каких-то нужд), gRPC/HTTP клиенты к другим сервисам.

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL:** `/api/v1/library` (маршрутизируется через API Gateway).
*   **Аутентификация:** JWT Bearer Token (`X-User-Id` извлекается).
*   **Авторизация:** Все эндпоинты требуют аутентификации и проверяют, что пользователь имеет доступ только к своим данным (например, `user_self_only`).
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

#### 3.1.1. Библиотека Игр (User Library)
*   **`GET /items`**
    *   Описание: Получение списка игр в библиотеке текущего пользователя.
    *   Query параметры: `status` (installed, not_installed), `category_id`, `search` (по названию), `sort_by` (name_asc, last_played_desc), `page`, `limit`.
    *   Пример ответа (Успех 200 OK): (Как в существующем документе)
    *   Пример ответа (Ошибка 401 Unauthorized - стандартизированный):
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
    *   Требуемые права доступа: `user_self_only`.
*   **`PATCH /items/{library_item_id}`**
    *   Описание: Изменение атрибутов элемента библиотеки (например, добавление в категорию, скрытие).
    *   Тело запроса: (Как в существующем документе)
    *   Пример ответа (Успех 200 OK): (Обновленный `libraryItem`)
    *   Требуемые права доступа: `user_self_only`.

#### 3.1.2. Игровое Время (Playtime)
(Описания эндпоинтов `/playtime/sessions/start`, `/playtime/sessions/{session_id}/heartbeat`, `/playtime/sessions/{session_id}/end` как в существующем документе).

#### 3.1.3. Достижения (Achievements)
(Описание эндпоинта `/achievements/status/{product_id}` как в существующем документе).

#### 3.1.4. Список Желаемого (Wishlist)
(Описание эндпоинта `/wishlist/items` как в существующем документе).

#### 3.1.5. Игровые Сохранения (Savegames)
(Описания эндпоинтов `/savegames/upload-url`, `/savegames/confirm-upload` как в существующем документе).

### 3.2. gRPC API
(Содержимое существующего раздела актуально).

### 3.3. WebSocket API
(Содержимое существующего раздела актуально).

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`UserLibraryItem` (Элемент Библиотеки Пользователя)** (Как в существующем документе)
*   **`PlaytimeSession` (Игровая Сессия)** (Как в существующем документе)
*   **`UserAchievement` (Достижение Пользователя)** (Как в существующем документе)
*   **`WishlistItem` (Элемент Списка Желаемого)** (Как в существующем документе)
*   **`SavegameMetadata` (Метаданные Игрового Сохранения)** (Как в существующем документе)
*   **`UserGameCategory` (Пользовательская Категория Игр)**
    *   `id` (UUID): Уникальный идентификатор категории.
    *   `user_id` (UUID): ID пользователя-владельца категории.
    *   `name` (VARCHAR(100)): Название категории (например, "Любимые", "Пройденные", "Для стрима").
    *   `display_order` (INTEGER): Порядок отображения категории.
    *   `created_at` (TIMESTAMPTZ), `updated_at` (TIMESTAMPTZ).
*   **`FamilySharingLink` (Связь Семейного Доступа)**
    *   `id` (UUID): Уникальный идентификатор связи.
    *   `owner_user_id` (UUID): ID пользователя, который делится библиотекой.
    *   `shared_with_user_id` (UUID): ID пользователя, которому предоставлен доступ.
    *   `status` (ENUM: `pending_approval`, `active`, `revoked`, `declined`): Статус связи.
    *   `shared_at` (TIMESTAMPTZ): Время создания запроса/активации.
    *   `last_access_at` (TIMESTAMPTZ, Nullable): Время последнего доступа к библиотеке по этой связи.
*   **`GameSpecificSettings` (Пользовательские Настройки Игры)**
    *   `user_id` (UUID, PK): ID пользователя.
    *   `product_id` (UUID, PK): ID продукта (игры).
    *   `settings_payload` (JSONB): Произвольные настройки игры в формате JSON (например, `{"graphics": "high", "sound_volume": 0.8, "key_bindings": {...}}`).
    *   `last_updated_at` (TIMESTAMPTZ): Время последнего обновления настроек.
    *   `last_synced_with_cloud_at` (TIMESTAMPTZ, Nullable): Время последней успешной синхронизации с облаком.

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL
**ERD Диаграмма (дополненная):**
```mermaid
erDiagram
    USERS {
        UUID id PK "User ID"
    }
    USER_LIBRARY_ITEMS {
        UUID id PK
        UUID user_id FK
        UUID product_id "FK (Product)"
        TIMESTAMPTZ added_at
        VARCHAR acquisition_type
        UUID shared_from_user_id "FK (User, nullable)"
        VARCHAR installation_status
        TIMESTAMPTZ last_played_at
        BIGINT total_playtime_seconds
        BOOLEAN is_hidden
    }
    PLAYTIME_SESSIONS {
        UUID id PK
        UUID user_id FK
        UUID product_id "FK (Product)"
        TIMESTAMPTZ start_time
        TIMESTAMPTZ end_time "nullable"
        INTEGER duration_seconds "nullable"
        TIMESTAMPTZ last_heartbeat_at
    }
    USER_ACHIEVEMENTS {
        UUID user_id PK FK
        UUID achievement_meta_id PK "FK (AchievementMeta)"
        UUID product_id "FK (Product)"
        BOOLEAN is_unlocked
        TIMESTAMPTZ unlocked_at "nullable"
        INTEGER current_progress "nullable"
    }
    WISHLIST_ITEMS {
        UUID id PK
        UUID user_id FK
        UUID product_id "FK (Product)"
        TIMESTAMPTZ added_at
        INTEGER priority "nullable"
    }
    SAVEGAME_METADATA {
        UUID id PK
        UUID user_id FK
        UUID product_id "FK (Product)"
        VARCHAR slot_name
        VARCHAR s3_object_key
        BIGINT file_size_bytes
        VARCHAR file_hash_sha256
        TIMESTAMPTZ client_modified_at
        TIMESTAMPTZ server_uploaded_at
        INTEGER version
    }
    USER_GAME_CATEGORIES {
        UUID id PK
        UUID user_id FK
        VARCHAR name
        INTEGER display_order
    }
    USER_LIBRARY_ITEM_TO_CATEGORIES { /* Исправлено имя таблицы */
        UUID library_item_id PK FK
        UUID category_id PK FK
    }
    FAMILY_SHARING_LINKS {
        UUID id PK
        UUID owner_user_id FK
        UUID shared_with_user_id FK
        VARCHAR status
        TIMESTAMPTZ shared_at
    }
    GAME_SPECIFIC_SETTINGS {
        UUID user_id PK FK
        UUID product_id PK "FK (Product)"
        JSONB settings_payload
        TIMESTAMPTZ last_updated_at
    }

    USERS ||--o{ USER_LIBRARY_ITEMS : "owns"
    USERS ||--o{ PLAYTIME_SESSIONS : "has"
    USERS ||--o{ USER_ACHIEVEMENTS : "earns"
    USERS ||--o{ WISHLIST_ITEMS : "wishes_for"
    USERS ||--o{ SAVEGAME_METADATA : "has_saves_for"
    USERS ||--o{ USER_GAME_CATEGORIES : "creates"
    USERS ||--o{ FAMILY_SHARING_LINKS : "owner_of"
    USERS ||--o{ FAMILY_SHARING_LINKS : "shared_with"
    USERS ||--o{ GAME_SPECIFIC_SETTINGS : "configures_for_game"

    USER_LIBRARY_ITEMS }o--|| PRODUCTS : "references"
    USER_LIBRARY_ITEMS }o--|{ USER_LIBRARY_ITEM_TO_CATEGORIES : "assigned_to"
    USER_GAME_CATEGORIES ||--o{ USER_LIBRARY_ITEM_TO_CATEGORIES : "has_items"

    PLAYTIME_SESSIONS }o--|| PRODUCTS : "references"
    USER_ACHIEVEMENTS }o--|| PRODUCTS : "references"
    USER_ACHIEVEMENTS }o--|| ACHIEVEMENT_METADATA : "references"
    WISHLIST_ITEMS }o--|| PRODUCTS : "references"
    SAVEGAME_METADATA }o--|| PRODUCTS : "references"
    GAME_SPECIFIC_SETTINGS }o--|| PRODUCTS : "references"

    USER_LIBRARY_ITEMS }o--|| USERS : "shared_by (family)"


    entity PRODUCTS { note "From Catalog Service" }
    entity ACHIEVEMENT_METADATA { note "From Catalog Service" }
```

**DDL (PostgreSQL - дополнения для TODO таблиц):**
```sql
CREATE TABLE user_game_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, name),
    CONSTRAINT fk_user_game_categories_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE -- Предполагается наличие таблицы users
);
COMMENT ON TABLE user_game_categories IS 'Пользовательские категории для организации игр в библиотеке.';

CREATE TABLE user_library_item_to_categories (
    library_item_id UUID NOT NULL REFERENCES user_library_items(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES user_game_categories(id) ON DELETE CASCADE,
    PRIMARY KEY (library_item_id, category_id)
);
COMMENT ON TABLE user_library_item_to_categories IS 'Связь элементов библиотеки с пользовательскими категориями.';

CREATE TABLE family_sharing_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_user_id UUID NOT NULL,
    shared_with_user_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_approval' CHECK (status IN ('pending_approval', 'active', 'revoked_by_owner', 'declined_by_user', 'revoked_by_admin')),
    shared_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_access_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (owner_user_id, shared_with_user_id),
    CONSTRAINT check_different_users CHECK (owner_user_id <> shared_with_user_id),
    CONSTRAINT fk_family_sharing_owner FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_family_sharing_shared_with FOREIGN KEY (shared_with_user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_family_sharing_links_owner ON family_sharing_links(owner_user_id);
CREATE INDEX idx_family_sharing_links_shared_with ON family_sharing_links(shared_with_user_id);
COMMENT ON TABLE family_sharing_links IS 'Связи для функции семейного доступа к библиотекам.';

CREATE TABLE game_specific_settings (
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    settings_payload JSONB NOT NULL,
    last_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_synced_with_cloud_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, product_id),
    CONSTRAINT fk_game_settings_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    -- CONSTRAINT fk_game_settings_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE -- Предполагается наличие products из Catalog
);
COMMENT ON TABLE game_specific_settings IS 'Пользовательские настройки для конкретных игр (синхронизируемые).';
```

#### 4.2.2. Redis
(Описание структуры данных в Redis остается как в существующем документе).

#### 4.2.3. S3-совместимое хранилище
*   **Роль:** Хранение бинарных файлов игровых сохранений.
*   **Структура (примерная):**
    *   `s3://<bucket-name-savegames>/<user_id_hash_prefix>/<user_id>/<product_id>/<savegame_slot_name_or_hash>/<timestamp_version_filename.sav>`
    *   `<user_id_hash_prefix>`: Первые несколько символов хеша от `user_id` для лучшего распределения объектов в S3 (если пользователей очень много).
    *   Файлы сохранений могут быть сжаты перед загрузкой.

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Формат событий:** CloudEvents v1.0 JSON (согласно `../../../../project_api_standards.md`).
*   **Основной топик Kafka:** `com.platform.library.events.v1`.

*   **`com.platform.library.game.added.v1`**
    *   `data` Payload: (Как в существующем документе, с корректным `type`).
*   **`com.platform.library.achievement.unlocked.v1`**
    *   `data` Payload: (Как в существующем документе, с корректным `type`).
*   **`com.platform.library.playtime.session.ended.v1`**
    *   `data` Payload: (Как в существующем документе, с корректным `type`).
*   **`com.platform.library.wishlist.item.added.v1`**
    *   Описание: Продукт добавлен в список желаемого пользователя.
    *   `data` Payload:
        ```json
        {
          "userId": "user-uuid-123",
          "productId": "game-uuid-789",
          "wishlistItemId": "wishlist-item-uuid-def",
          "addedAt": "2024-03-18T16:00:00Z"
        }
        ```
    *   Потребители: Notification Service (для информирования о скидках), Analytics Service.
*   **`com.platform.library.savegame.uploaded.v1`**
    *   Описание: Файл сохранения игры был успешно загружен в облако.
    *   `data` Payload:
        ```json
        {
          "userId": "user-uuid-123",
          "productId": "game-uuid-abc",
          "savegameId": "savegame-meta-uuid-1",
          "slotName": "slot1_manual_save",
          "s3Path": "s3://bucket/user_id/product_id/...",
          "fileSizeBytes": 1048576,
          "clientTimestamp": "2024-03-18T16:30:00Z",
          "serverTimestamp": "2024-03-18T16:30:05Z"
        }
        ```
    *   Потребители: Analytics Service (для статистики использования облачных сохранений).

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
*   **Безопасность облачных сохранений:**
    *   Использование pre-signed URLs для S3 с коротким TTL для загрузки и скачивания файлов сохранений.
    *   Проверка прав доступа пользователя к `product_id` перед генерацией URL.
    *   Валидация хеш-сумм файлов на стороне сервера после загрузки.
    *   Рассмотреть возможность клиентского шифрования сохранений перед загрузкой в S3, если требуется максимальная приватность (ключи шифрования управляются клиентом, Library Service не имеет к ним доступа). Это усложнит некоторые сценарии (например, просмотр сохранений на разных устройствах без экспорта ключа).
*   Ссылки на `../../../../project_security_standards.md`.

## 10. Развертывание (Deployment)
(Содержимое существующего раздела актуально).

## 11. Мониторинг и Логирование (Logging and Monitoring)
(Содержимое существующего раздела актуально).

## 12. Нефункциональные Требования (NFRs)
(Содержимое существующего раздела актуально).

## 13. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API, Protobuf определения для gRPC API, и форматы сообщений WebSocket поддерживаются в соответствующих репозиториях исходного кода сервиса и/или в централизованном репозитории `platform-protos`.
*   DDL схемы базы данных управляются через систему миграций (например, `golang-migrate/migrate`) и хранятся в репозитории исходного кода сервиса. Актуальные версии доступны во внутренней документации команды разработки и в GitOps репозитории.

## 14. Резервное Копирование и Восстановление (Backup and Recovery)

### 14.1. PostgreSQL (Данные библиотеки, достижений, игрового времени, списка желаемого, метаданные сохранений)
*   **Процедура резервного копирования:**
    *   Ежедневный логический бэкап (`pg_dump`).
    *   Настроена непрерывная архивация WAL-сегментов (PITR), базовый бэкап еженедельно.
    *   **Хранение:** Бэкапы в S3, шифрование, версионирование, другой регион. Срок хранения: полные - 30 дней, WAL - 14 дней.
*   **Процедура восстановления:** Тестируется ежеквартально.
*   **RTO:** < 1 часа.
*   **RPO:** < 5 минут.
*   (Общие принципы см. `../../../../project_database_structure.md`).

### 14.2. Redis (Кэш, текущие игровые сессии)
*   **Стратегия персистентности:**
    *   **AOF (Append Only File):** Включен с fsync `everysec` для данных текущих игровых сессий, если их сохранение при перезапуске Redis критично.
    *   **RDB Snapshots:** Регулярное создание снапшотов (например, каждые 1-6 часов).
*   **Резервное копирование (снапшотов):** RDB-снапшоты могут копироваться в S3 ежедневно. Срок хранения - 3-7 дней.
*   **Восстановление:** Из последнего RDB-снапшота и/или AOF. Кэшированные данные будут перестроены.
*   **RTO:** < 30 минут.
*   **RPO:** < 1 минуты (для данных с AOF `everysec`). Для кэша RPO менее критичен.

### 14.3. S3-совместимое хранилище (Игровые сохранения)
*   **Процедура резервного копирования:**
    *   **Версионирование объектов:** Включено для бакета с игровыми сохранениями.
    *   **Политики жизненного цикла (Lifecycle Policies):** Для управления старыми версиями сохранений (например, удаление версий старше X дней/месяцев, если не помечены как "значимые").
    *   **Cross-Region Replication (CRR):** Настроена для бакета с сохранениями для обеспечения гео-резервирования.
*   **Процедура восстановления:** Восстановление отдельных объектов/версий из S3 или переключение на реплицированный бакет.
*   **RTO:** Зависит от объема восстанавливаемых данных, но обычно быстро для отдельных сохранений.
*   **RPO:** Близко к нулю (ограничено временем репликации S3).

### 14.4. Общая стратегия
*   Восстановление PostgreSQL и метаданных сохранений в S3 является приоритетным.
*   Процедуры восстановления тестируются и документируются.
*   Мониторинг процессов резервного копирования.

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [Процесс покупки игры и обновления библиотеки пользователя](../../../../project_workflows/game_purchase_flow.md)
*   [Синхронизация игровых сохранений с облаком](../../../../project_workflows/cloud_save_sync_flow.md) (Примечание: Создание документа `cloud_save_sync_flow.md` является частью общей задачи по документированию проекта и выходит за рамки обновления документации данного микросервиса.)
*   [Процесс разблокировки достижений](../../../../project_workflows/achievement_unlocking_flow.md) (Примечание: Создание документа `achievement_unlocking_flow.md` является частью общей задачи по документированию проекта и выходит за рамки обновления документации данного микросервиса.)

---
*Этот документ является основной спецификацией для Library Service и должен поддерживаться в актуальном состоянии.*
