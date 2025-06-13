# Спецификация Микросервиса: Library Service (Сервис Библиотеки Пользователя)

**Версия:** 1.2 (адаптировано из v1.1 внешней спецификации)
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
*   Разработка сервиса должна вестись в соответствии с `CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Управление библиотекой пользователя:** Добавление игр в библиотеку после покупки/активации, скрытие игр, фильтрация и сортировка списка игр, создание пользовательских категорий/коллекций, отслеживание статуса установки игры (интеграция с Download Service). Поддержка концепции "семейного доступа" (Family Sharing) для предоставления временного доступа к играм другим пользователям.
*   **Отслеживание игрового времени:** Автоматическая регистрация начала и окончания игровых сессий, подсчет общего и сессионного игрового времени для каждого продукта, отображение статистики последней активности, синхронизация игрового времени между устройствами.
*   **Управление достижениями:** Регистрация полученных пользователем достижений (анлоков), хранение прогресса по частично выполненным достижениям, отображение списка всех доступных и полученных достижений для игры, глобальная статистика по достижениям (редкость).
*   **Управление списком желаемого (Wishlist):** Добавление и удаление продуктов из списка желаемого, возможность приоритизации элементов, получение уведомлений о скидках на игры в списке (через Notification Service).
*   **Синхронизация игровых сохранений (Cloud Saves):** Загрузка файлов сохранений в облачное S3-хранилище, скачивание сохранений на другие устройства пользователя, разрешение конфликтов версий сохранений, версионирование (опционально), автоматическая и ручная синхронизация, управление квотами на облачное хранилище для сохранений.
*   **Настройки игр:** Хранение и синхронизация пользовательских настроек для конкретных игр (например, настройки графики, управления), если игра поддерживает такую интеграцию.

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+).
*   **Веб-фреймворк (REST API):** Echo (v4+) или Gin.
*   **RPC фреймворк (gRPC):** `google.golang.org/grpc` для межсервисного взаимодействия.
*   **База данных (основная):** PostgreSQL (версия 15+) для хранения информации о библиотеках, достижениях, списках желаемого, метаданных сохранений.
*   **Кэширование/Оперативные данные:** Redis (версия 7+) для кэширования часто запрашиваемых данных (например, содержимое библиотеки пользователя, текущие игровые сессии), хранения временных данных (например, токены для загрузки/скачивания сохранений из S3).
*   **Облачное хранилище (для сохранений):** S3-совместимое объектное хранилище (например, MinIO).
*   **Брокер сообщений:** Apache Kafka (версия 3.x+) для асинхронного обмена событиями с другими сервисами.
*   **Инфраструктура:** Docker, Kubernetes, Helm.
*   **Мониторинг/Трассировка:** OpenTelemetry, Prometheus, Grafana, Jaeger/Tempo.
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `PACKAGE_STANDARDIZATION.md`.
*   *Примечание: Приведенные в последующих разделах примеры API и конфигураций основаны на предположении использования Go (Echo/Gin для REST), PostgreSQL, Redis, Kafka и S3-совместимого хранилища.*

### 1.4. Термины и Определения (Glossary)
*   **Библиотека (Library):** Персональная коллекция игр и других цифровых продуктов, принадлежащих пользователю или доступных ему (например, через Family Sharing).
*   **Игровое время (Playtime):** Общее время, проведенное пользователем в конкретной игре.
*   **Достижение (Achievement):** Виртуальная награда, выдаваемая пользователю за выполнение определенных условий или задач в игре.
*   **Список желаемого (Wishlist):** Персональный список продуктов, которые пользователь хочет приобрести в будущем.
*   **Игровое сохранение (Savegame/Cloud Save):** Файл или набор файлов, содержащих прогресс пользователя в игре, который может быть синхронизирован с облачным хранилищем.
*   **Семейный доступ (Family Sharing):** Функция, позволяющая делиться играми из своей библиотеки с членами семьи или близкими друзьями.
*   Для других общих терминов см. `project_glossary.md`.

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
    *   **Entities (Сущности):** `UserLibraryItem`, `PlaytimeSession`, `UserAchievement`, `WishlistItem`, `SavegameMetadata`, `GameSpecificSettings`.
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
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway, `X-User-Id` передается в заголовках).
*   **Авторизация:** Все эндпоинты требуют аутентификации и проверяют, что пользователь имеет доступ только к своим данным (например, `user_self_only`).
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

#### 3.1.1. Библиотека Игр (User Library)
*   **`GET /items`**
    *   Описание: Получение списка игр в библиотеке текущего пользователя.
    *   Query параметры: `status` (installed, not_installed), `category_id`, `search` (по названию), `sort_by` (name_asc, last_played_desc), `page`, `limit`.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "libraryItem",
              "id": "lib-item-uuid-1",
              "attributes": {
                "product_id": "game-uuid-123",
                "title": "Супер Игра X", // Получено из Catalog Service или закэшировано
                "cover_image_url": "https://cdn.example.com/covers/game-uuid-123.jpg",
                "added_at": "2023-01-15T10:00:00Z",
                "last_played_at": "2024-03-10T18:30:00Z",
                "total_playtime_seconds": 72000,
                "installation_status": "installed", // "not_installed", "updating"
                "user_rating": 5 // Опционально, если пользователь ставил оценку
              },
              "relationships": { "categories": { "data": [{"type": "userCategory", "id": "cat-uuid-fav"}] } },
              "links": { "self": "/api/v1/library/items/lib-item-uuid-1" }
            }
          ],
          "meta": { "total_items": 50, "current_page": 1, "per_page": 20 },
          "links": { "next": "/api/v1/library/items?page=2" }
        }
        ```
    *   Требуемые права доступа: `user_self_only`.
*   **`PATCH /items/{library_item_id}`**
    *   Описание: Изменение атрибутов элемента библиотеки (например, добавление в категорию, скрытие).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "libraryItemUpdate",
            "id": "lib-item-uuid-1",
            "attributes": {
              "is_hidden": true,
              "category_ids": ["cat-uuid-fav", "cat-uuid-completed"]
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): (Обновленный `libraryItem`)
    *   Требуемые права доступа: `user_self_only`.

#### 3.1.2. Игровое Время (Playtime)
*   **`POST /playtime/sessions/start`**
    *   Описание: Уведомление о начале игровой сессии.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "playtimeSessionStart",
            "attributes": { "product_id": "game-uuid-123" }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created):
        ```json
        {
          "data": {
            "type": "playtimeSession",
            "id": "session-uuid-abc",
            "attributes": {
              "product_id": "game-uuid-123",
              "start_time": "2024-03-15T18:00:00Z",
              "status": "active"
            }
          }
        }
        ```
    *   Требуемые права доступа: `user_self_only`.
*   **`POST /playtime/sessions/{session_id}/heartbeat`**
    *   Описание: Периодическое уведомление о продолжении активной игровой сессии.
    *   Пример ответа (Успех 204 No Content)
    *   Требуемые права доступа: `user_self_only`.
*   **`POST /playtime/sessions/{session_id}/end`**
    *   Описание: Уведомление о завершении игровой сессии.
    *   Пример ответа (Успех 200 OK): (Обновленная сессия со статусом `completed` и `duration_seconds`)
    *   Требуемые права доступа: `user_self_only`.

#### 3.1.3. Достижения (Achievements)
*   **`GET /achievements/status/{product_id}`**
    *   Описание: Получение статуса всех достижений для указанной игры для текущего пользователя.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "userAchievementStatus",
              "id": "ach-meta-uuid-001", // ID метаданных достижения из Catalog Service
              "attributes": {
                "name": "Первый шаг",
                "description": "Завершить обучение.",
                "icon_url_unlocked": "...",
                "is_unlocked": true,
                "unlocked_at": "2024-03-01T10:00:00Z",
                "current_progress": null, // или { "current": 5, "total": 10 } для прогрессивных
                "rarity_percentage": 75.5 // Глобальная редкость
              }
            }
          ]
        }
        ```
    *   Требуемые права доступа: `user_self_only`.

#### 3.1.4. Список Желаемого (Wishlist)
*   **`POST /wishlist/items`**
    *   Описание: Добавление продукта в список желаемого.
    *   Тело запроса:
        ```json
        {
          "data": { "type": "wishlistItemCreation", "attributes": { "product_id": "game-uuid-789" } }
        }
        ```
    *   Пример ответа (Успех 201 Created):
        ```json
        {
          "data": {
            "type": "wishlistItem",
            "id": "wishlist-item-uuid-def",
            "attributes": { "product_id": "game-uuid-789", "added_at": "2024-03-15T15:00:00Z" }
          }
        }
        ```
    *   Требуемые права доступа: `user_self_only`.

#### 3.1.5. Игровые Сохранения (Savegames)
*   **`POST /savegames/upload-url`**
    *   Описание: Запрос pre-signed URL для загрузки файла сохранения в S3.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "savegameUploadRequest",
            "attributes": {
              "product_id": "game-uuid-123",
              "save_slot_name": "slot1_manual_save", // Имя слота или файла сохранения
              "file_size_bytes": 1048576, // 1MB
              "file_hash_sha256": "abcdef123..." // Хеш файла для проверки целостности
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
              "expires_in_seconds": 600,
              "internal_save_id": "save-uuid-temp" // Временный ID для подтверждения
            }
          }
        }
        ```
    *   Требуемые права доступа: `user_self_only`.
*   **`POST /savegames/confirm-upload`**
    *   Описание: Подтверждение успешной загрузки файла сохранения в S3.
    *   Тело запроса: `{ "data": { "type": "savegameUploadConfirmation", "attributes": { "internal_save_id": "save-uuid-temp", "s3_path": "actual/path/in/s3/slot1.sav" } } }`
    *   Пример ответа (Успех 200 OK): (Метаданные созданного/обновленного сохранения)
    *   Требуемые права доступа: `user_self_only`.

### 3.2. gRPC API
*   Пакет: `library.v1`.
*   Определение Protobuf: `library/v1/library_query_service.proto`.

#### 3.2.1. Сервис: `LibraryQueryService`
*   **`rpc GetLibraryItems(GetLibraryItemsRequest) returns (GetLibraryItemsResponse)`**
    *   Описание: Получение списка элементов библиотеки пользователя для внутреннего использования другими сервисами.
    *   `message GetLibraryItemsRequest { string user_id = 1; repeated string product_ids_filter = 2; /* опционально */ }`
    *   `message LibraryItemInternal { string product_id = 1; google.protobuf.Timestamp added_at = 2; string installation_status = 3; int64 total_playtime_seconds = 4; }`
    *   `message GetLibraryItemsResponse { repeated LibraryItemInternal items = 1; }`
*   **`rpc CheckGameAccess(CheckGameAccessRequest) returns (CheckGameAccessResponse)`**
    *   Описание: Проверка, имеет ли пользователь доступ (владеет ли) к указанной игре.
    *   `message CheckGameAccessRequest { string user_id = 1; string product_id = 2; }`
    *   `message CheckGameAccessResponse { bool has_access = 1; google.protobuf.Timestamp acquisition_date = 2; /* опционально */ }`

### 3.3. WebSocket API
*   **Эндпоинт:** `/api/v1/ws/library-updates` (маршрутизируется через API Gateway).
*   **Аутентификация:** JWT токен при установлении соединения.
*   **Сообщения от сервера к клиенту (примеры):**
    *   **Событие: `library.item.updated`** (например, изменился статус установки, добавлена категория)
        ```json
        {
          "event_type": "library.item.updated",
          "payload": {
            "library_item_id": "lib-item-uuid-1",
            "product_id": "game-uuid-123",
            "changes": { "installation_status": "installed" }
          }
        }
        ```
    *   **Событие: `achievement.unlocked`**
        ```json
        {
          "event_type": "achievement.unlocked",
          "payload": {
            "product_id": "game-uuid-123",
            "achievement_id": "ach-meta-uuid-001", // ID метаданных достижения
            "achievement_name": "Первый шаг",
            "unlocked_at": "2024-03-15T19:00:00Z"
          }
        }
        ```
    *   **Событие: `savegame.sync.status`**
        ```json
        {
          "event_type": "savegame.sync.status",
          "payload": {
            "product_id": "game-uuid-123",
            "status": "sync_completed", // "sync_started", "sync_failed"
            "last_synced_at": "2024-03-15T20:00:00Z",
            "error_message": null // если failed
          }
        }
        ```

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности

*   **`UserLibraryItem` (Элемент Библиотеки Пользователя)**
    *   `id` (UUID): Уникальный идентификатор записи в библиотеке. Обязательность: Required.
    *   `user_id` (UUID): ID пользователя. Обязательность: Required.
    *   `product_id` (UUID): ID продукта (игры, DLC и т.д.) из Catalog Service. Обязательность: Required.
    *   `added_at` (TIMESTAMPTZ): Дата и время добавления в библиотеку. Обязательность: Required.
    *   `acquisition_type` (ENUM: `purchase`, `gift`, `key_activation`, `family_shared_access`, `free_to_play`): Тип приобретения. Обязательность: Required.
    *   `shared_from_user_id` (UUID): Если `family_shared_access`, ID пользователя, предоставившего доступ. Обязательность: Optional.
    *   `installation_status` (ENUM: `not_installed`, `installing`, `installed`, `updating`, `repair_needed`, `error`): Статус установки. Обязательность: Required.
    *   `last_played_at` (TIMESTAMPTZ): Дата и время последнего запуска. Обязательность: Optional.
    *   `total_playtime_seconds` (BIGINT): Общее наигранное время в секундах. Обязательность: Required, default 0.
    *   `is_hidden` (BOOLEAN): Скрыта ли игра в библиотеке пользователя. Обязательность: Required, default false.
    *   `user_categories` (ARRAY of UUID): ID пользовательских категорий, к которым отнесена игра. Обязательность: Optional.

*   **`PlaytimeSession` (Игровая Сессия)**
    *   `id` (UUID): Уникальный идентификатор сессии. Обязательность: Required.
    *   `user_id` (UUID): ID пользователя. Обязательность: Required.
    *   `product_id` (UUID): ID продукта. Обязательность: Required.
    *   `start_time` (TIMESTAMPTZ): Время начала сессии. Обязательность: Required.
    *   `end_time` (TIMESTAMPTZ): Время окончания сессии. Обязательность: Optional (null для активных сессий).
    *   `duration_seconds` (INTEGER): Длительность сессии в секундах (рассчитывается при завершении). Обязательность: Optional.
    *   `last_heartbeat_at` (TIMESTAMPTZ): Время последнего "пульса" от клиента. Обязательность: Required для активных сессий.

*   **`UserAchievement` (Достижение Пользователя)**
    *   `user_id` (UUID, PK): ID пользователя. Обязательность: Required.
    *   `achievement_meta_id` (UUID, PK, FK to AchievementMeta in Catalog Service): ID метаданных достижения. Обязательность: Required.
    *   `product_id` (UUID, FK to Product in Catalog Service): ID продукта, к которому относится достижение. Обязательность: Required.
    *   `is_unlocked` (BOOLEAN): Разблокировано ли достижение. Обязательность: Required, default false.
    *   `unlocked_at` (TIMESTAMPTZ): Дата и время разблокировки. Обязательность: Optional.
    *   `current_progress` (INTEGER): Текущий прогресс для прогрессивных достижений. Обязательность: Optional.
    *   `total_progress_needed` (INTEGER): Общий прогресс, необходимый для разблокировки (из метаданных). Обязательность: Optional.

*   **`WishlistItem` (Элемент Списка Желаемого)**
    *   `id` (UUID): Уникальный идентификатор элемента. Обязательность: Required.
    *   `user_id` (UUID): ID пользователя. Обязательность: Required.
    *   `product_id` (UUID): ID продукта. Обязательность: Required.
    *   `added_at` (TIMESTAMPTZ): Дата добавления. Обязательность: Required.
    *   `priority` (INTEGER): Приоритет (если поддерживается). Обязательность: Optional.
    *   `notified_on_sale_id` (UUID): ID последней скидки, о которой было уведомление. Обязательность: Optional.

*   **`SavegameMetadata` (Метаданные Игрового Сохранения)**
    *   `id` (UUID): Уникальный идентификатор метаданных сохранения. Обязательность: Required.
    *   `user_id` (UUID): ID пользователя. Обязательность: Required.
    *   `product_id` (UUID): ID продукта. Обязательность: Required.
    *   `slot_name` (VARCHAR(255)): Имя слота сохранения или идентифицирующее имя файла. Пример: `quicksave_20240315_2000.sav`. Валидация: not null. Обязательность: Required.
    *   `s3_path` (VARCHAR(1024)): Путь к файлу сохранения в S3. Валидация: not null. Обязательность: Required.
    *   `file_size_bytes` (BIGINT): Размер файла сохранения. Обязательность: Required.
    *   `file_hash_sha256` (VARCHAR(64)): SHA256 хеш файла. Обязательность: Required.
    *   `client_timestamp` (TIMESTAMPTZ): Временная метка сохранения, установленная клиентом игры. Обязательность: Required.
    *   `server_timestamp` (TIMESTAMPTZ): Временная метка загрузки на сервер. Обязательность: Required.
    *   `version` (INTEGER): Версия сохранения (для разрешения конфликтов). Обязательность: Optional, default 1.
    *   `comment` (TEXT): Комментарий пользователя к сохранению. Обязательность: Optional.

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL

**ERD Диаграмма (ключевые таблицы):**
```mermaid
erDiagram
    USERS {
        UUID id PK "User ID (from Auth Service)"
        -- Other user-related info not stored here
    }
    USER_LIBRARY_ITEMS {
        UUID id PK
        UUID user_id FK
        UUID product_id "FK (Product from Catalog)"
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
        UUID product_id "FK (Product from Catalog)"
        TIMESTAMPTZ start_time
        TIMESTAMPTZ end_time "nullable"
        INTEGER duration_seconds "nullable"
        TIMESTAMPTZ last_heartbeat_at
    }
    USER_ACHIEVEMENTS {
        UUID user_id PK FK
        UUID achievement_meta_id PK "FK (AchievementMeta from Catalog)"
        UUID product_id "FK (Product from Catalog)"
        BOOLEAN is_unlocked
        TIMESTAMPTZ unlocked_at "nullable"
        INTEGER current_progress "nullable"
    }
    WISHLIST_ITEMS {
        UUID id PK
        UUID user_id FK
        UUID product_id "FK (Product from Catalog)"
        TIMESTAMPTZ added_at
        INTEGER priority "nullable"
    }
    SAVEGAME_METADATA {
        UUID id PK
        UUID user_id FK
        UUID product_id "FK (Product from Catalog)"
        VARCHAR slot_name
        VARCHAR s3_path
        BIGINT file_size_bytes
        VARCHAR file_hash_sha256
        TIMESTAMPTZ client_timestamp
        TIMESTAMPTZ server_timestamp
        INTEGER version
    }

    USERS ||--o{ USER_LIBRARY_ITEMS : "owns"
    USERS ||--o{ PLAYTIME_SESSIONS : "has"
    USERS ||--o{ USER_ACHIEVEMENTS : "earns"
    USERS ||--o{ WISHLIST_ITEMS : "wishes_for"
    USERS ||--o{ SAVEGAME_METADATA : "has_saves_for"

    USER_LIBRARY_ITEMS }o--|| PRODUCTS : "references_product"
    PLAYTIME_SESSIONS }o--|| PRODUCTS : "references_product"
    USER_ACHIEVEMENTS }o--|| PRODUCTS : "references_product"
    USER_ACHIEVEMENTS }o--|| ACHIEVEMENT_METADATA : "references_achievement_meta"
    WISHLIST_ITEMS }o--|| PRODUCTS : "references_product"
    SAVEGAME_METADATA }o--|| PRODUCTS : "references_product"

    entity PRODUCTS { # Предполагается из Catalog Service
        UUID id PK
        VARCHAR title
    }
    entity ACHIEVEMENT_METADATA { # Предполагается из Catalog Service
        UUID id PK
        VARCHAR name
    }

```
*Примечание: `USERS`, `PRODUCTS`, `ACHIEVEMENT_METADATA` - это концептуальные ссылки на сущности из других сервисов (Auth, Catalog).*

**DDL (PostgreSQL - примеры ключевых таблиц):**
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_library_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL, -- Внешний ключ на таблицу пользователей в Auth Service (неявный)
    product_id UUID NOT NULL, -- Внешний ключ на таблицу продуктов в Catalog Service (неявный)
    added_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    acquisition_type VARCHAR(50) NOT NULL CHECK (acquisition_type IN ('purchase', 'gift', 'key_activation', 'family_shared_access', 'free_to_play')),
    shared_from_user_id UUID, -- Для family_shared_access
    installation_status VARCHAR(50) NOT NULL DEFAULT 'not_installed' CHECK (installation_status IN ('not_installed', 'installing', 'installed', 'updating', 'repair_needed', 'error')),
    last_played_at TIMESTAMPTZ,
    total_playtime_seconds BIGINT NOT NULL DEFAULT 0,
    is_hidden BOOLEAN NOT NULL DEFAULT FALSE,
    user_defined_categories UUID[], -- Массив ID пользовательских категорий
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, product_id)
);
CREATE INDEX idx_user_library_items_user_id ON user_library_items(user_id);
CREATE INDEX idx_user_library_items_last_played ON user_library_items(user_id, last_played_at DESC NULLS LAST);

CREATE TABLE playtime_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    duration_seconds INTEGER,
    last_heartbeat_at TIMESTAMPTZ NOT NULL,
    platform VARCHAR(50), -- pc, mobile, web
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_playtime_sessions_user_product ON playtime_sessions(user_id, product_id);
CREATE INDEX idx_playtime_sessions_active ON playtime_sessions(user_id, product_id, end_time) WHERE end_time IS NULL;

CREATE TABLE user_achievements (
    user_id UUID NOT NULL,
    achievement_meta_id UUID NOT NULL, -- FK to AchievementMeta in Catalog Service
    product_id UUID NOT NULL, -- FK to Product in Catalog Service
    is_unlocked BOOLEAN NOT NULL DEFAULT FALSE,
    unlocked_at TIMESTAMPTZ,
    current_progress INTEGER,
    PRIMARY KEY (user_id, achievement_meta_id)
);
CREATE INDEX idx_user_achievements_user_product ON user_achievements(user_id, product_id);

CREATE TABLE wishlist_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    priority INTEGER DEFAULT 0,
    UNIQUE (user_id, product_id)
);

CREATE TABLE savegame_metadata (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    slot_name VARCHAR(255) NOT NULL,
    s3_object_key VARCHAR(1024) NOT NULL UNIQUE, -- Полный путь в S3 бакете
    file_size_bytes BIGINT NOT NULL,
    file_hash_sha256 VARCHAR(64) NOT NULL,
    client_modified_at TIMESTAMPTZ NOT NULL, -- Время изменения файла на клиенте
    server_uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now(), -- Время загрузки на сервер
    version INTEGER NOT NULL DEFAULT 1,
    comment TEXT,
    UNIQUE (user_id, product_id, slot_name)
);
CREATE INDEX idx_savegame_metadata_user_product ON savegame_metadata(user_id, product_id, server_uploaded_at DESC);

-- TODO: Добавить DDL для user_game_categories, family_sharing_links, game_settings.
```

#### 4.2.2. Redis
*   **Кэширование:**
    *   Содержимое библиотеки пользователя: `library:<user_id>:items` (JSON или HASH). TTL: минуты.
    *   Статус достижений пользователя для игры: `achievements:<user_id>:<product_id>` (JSON). TTL: минуты/часы.
    *   Список желаемого: `wishlist:<user_id>:items` (JSON). TTL: минуты.
*   **Текущие игровые сессии (для быстрого доступа и heartbeat):**
    *   Ключ: `playtime_session:<user_id>:<product_id>` (HASH). Поля: `session_id`, `start_time`, `last_heartbeat_at`. TTL: часы (обновляется при heartbeat).
*   **Временные токены/URL для S3 (если генерируются сервисом):**
    *   Ключ: `s3_presigned_url_token:<user_id>:<savegame_id>` (STRING). Значение: токен или URL. TTL: короткий (минуты).
*   **Очереди задач (если используются для фоновой обработки):**
    *   Например, очередь для обработки событий синхронизации сохранений.

#### 4.2.3. S3-совместимое хранилище
*   **Роль:** Хранение бинарных файлов игровых сохранений.
*   **Структура (примерная):**
    *   `s3://<bucket-savegames>/<user_id>/<product_id>/<savegame_id_or_slot_name_versioned>`
    *   Каждый объект представляет собой файл сохранения или архив с файлами сохранения.

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Apache Kafka.
*   **Формат событий:** CloudEvents v1.0 (JSON encoding), согласно `project_api_standards.md`.
*   **Основной топик:** `library.events.v1`. (Может быть разделен на более гранулярные).

*   **`library.game.added.v1`**
    *   Описание: Игра добавлена в библиотеку пользователя.
    *   Пример Payload (`data` секция CloudEvent):
        ```json
        {
          "user_id": "user-uuid-123",
          "product_id": "game-uuid-abc",
          "acquisition_type": "purchase",
          "added_at": "2024-03-15T10:00:00Z"
        }
        ```
*   **`library.achievement.unlocked.v1`**
    *   Описание: Пользователь разблокировал достижение.
    *   Пример Payload:
        ```json
        {
          "user_id": "user-uuid-123",
          "product_id": "game-uuid-abc",
          "achievement_meta_id": "ach-meta-uuid-001",
          "unlocked_at": "2024-03-15T12:30:00Z"
        }
        ```
*   **`library.playtime.session.ended.v1`**
    *   Описание: Игровая сессия пользователя завершилась, обновлено общее игровое время.
    *   Пример Payload:
        ```json
        {
          "user_id": "user-uuid-123",
          "product_id": "game-uuid-abc",
          "session_id": "session-uuid-xyz",
          "session_duration_seconds": 3600,
          "total_playtime_seconds": 75600,
          "ended_at": "2024-03-15T19:00:00Z"
        }
        ```
*   TODO: Детализировать другие события, например, `library.wishlist.item.added.v1`, `library.savegame.uploaded.v1`, `library.game.settings.updated.v1`.

### 5.2. Потребляемые События (Consumed Events)

*   **`payment.purchase.completed.v1`** (от Payment Service)
    *   Описание: Успешная покупка продукта пользователем.
    *   Ожидаемый Payload:
        ```json
        {
          "user_id": "user-uuid-123",
          "order_id": "order-uuid-xyz",
          "items": [
            { "product_id": "game-uuid-abc", "product_type": "game", "quantity": 1 },
            { "product_id": "dlc-uuid-def", "product_type": "dlc", "quantity": 1 }
          ],
          "purchase_timestamp": "2024-03-15T09:55:00Z"
        }
        ```
    *   Логика обработки: Добавить каждый купленный продукт (`product_id`) в таблицу `user_library_items` для данного `user_id` с типом приобретения `purchase`. Опубликовать событие `library.game.added.v1` для каждого добавленного элемента.
*   **`catalog.product.updated.v1`** (от Catalog Service)
    *   Описание: Обновлена информация о продукте в каталоге (например, название, иконки).
    *   Ожидаемый Payload:
        ```json
        {
          "product_id": "game-uuid-abc",
          "updated_fields": ["title", "cover_image_url"],
          "new_title": "Супер Игра X - Издание Года",
          "new_cover_image_url": "https://cdn.example.com/covers/new_cover.jpg"
        }
        ```
    *   Логика обработки: Обновить кэшированные данные о продукте в Redis, если такие есть. При следующем запросе библиотеки пользователя данные будут подтянуты из Catalog Service или обновленного кэша.
*   **`download.game.installed.v1`** / **`download.game.uninstalled.v1`** (от Download Service)
    *   Описание: Статус установки игры изменился.
    *   Ожидаемый Payload:
        ```json
        {
          "user_id": "user-uuid-123",
          "product_id": "game-uuid-abc",
          "version_id": "version-uuid-123",
          "event_type": "installed", // "uninstalled"
          "timestamp": "2024-03-15T20:00:00Z"
        }
        ```
    *   Логика обработки: Обновить поле `installation_status` в `user_library_items` для данного пользователя и продукта. Опубликовать событие `library.item.updated.v1` через WebSocket.
*   **`account.user.deleted.v1`** (от Account Service)
    *   Описание: Аккаунт пользователя был удален.
    *   Ожидаемый Payload: `{"user_id": "user-uuid-123", "deleted_at": "..."}`
    *   Логика обработки: Анонимизировать или удалить все данные пользователя из Library Service (библиотека, игровое время, достижения, сохранения, список желаемого) в соответствии с политиками GDPR/ФЗ-152.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Account Service:** Получение базовой информации о пользователе (ID, статус) для привязки данных библиотеки (через gRPC).
*   **Catalog Service:** Получение метаданных игр, DLC, приложений, а также метаданных достижений (через gRPC).
*   **Payment Service:** Получение событий об успешных покупках для добавления продуктов в библиотеку (через Kafka).
*   **Download Service:** Получение событий об установке/удалении игр для обновления статуса в библиотеке (через Kafka). Инициирование загрузки игры из библиотеки (через gRPC к Download Service).
*   **Social Service:** Публикация событий о разблокировке достижений или других игровых событиях для отображения в ленте активности или для друзей (через Kafka).
*   **Analytics Service:** Отправка анонимизированных данных об игровой активности (время игры, популярные игры, разблокированные достижения) для анализа (через Kafka).
*   **Notification Service:** Отправка событий для уведомления пользователей (например, о скидках на игры из списка желаемого, о завершении синхронизации сохранений) (через Kafka).
*   **API Gateway:** Обработка и маршрутизация клиентских запросов (REST API, WebSocket) к Library Service.
*   **Auth Service:** Валидация JWT токенов (обычно на уровне API Gateway, но может быть дополнительная проверка).

### 6.2. Внешние Системы
*   **S3-совместимое облачное хранилище:** Для хранения файлов игровых сохранений пользователей. Library Service управляет метаданными и генерирует pre-signed URL для безопасной загрузки/скачивания файлов клиентом напрямую из/в S3.

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `LIBRARY_SERVICE_HTTP_PORT`: Порт для REST API (например, `8084`).
*   `LIBRARY_SERVICE_GRPC_PORT`: Порт для gRPC API (например, `9094`).
*   `LIBRARY_SERVICE_WS_PORT`: Порт для WebSocket (если обрабатывается напрямую, а не через API Gateway)
*   `POSTGRES_DSN`: Строка подключения к PostgreSQL.
*   `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB_LIBRARY`: Параметры Redis.
*   `KAFKA_BROKERS`: Список брокеров Kafka.
*   `KAFKA_TOPIC_LIBRARY_EVENTS`: Топик для публикуемых событий.
*   `KAFKA_CONSUMER_GROUP_LIBRARY`: Группа консьюмеров для входящих событий.
*   `S3_ENDPOINT`, `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY`, `S3_BUCKET_SAVEGAMES`, `S3_REGION`, `S3_USE_SSL`: Параметры S3 для хранилища сохранений.
*   `S3_PRESIGNED_URL_TTL_SECONDS`: Время жизни pre-signed URL для S3.
*   `LOG_LEVEL`: Уровень логирования.
*   `AUTH_SERVICE_GRPC_ADDR`: Адрес gRPC Auth Service.
*   `CATALOG_SERVICE_GRPC_ADDR`: Адрес gRPC Catalog Service.
*   `DOWNLOAD_SERVICE_GRPC_ADDR`: Адрес gRPC Download Service.
*   `ACCOUNT_SERVICE_GRPC_ADDR`: Адрес gRPC Account Service.
*   `PLAYTIME_HEARTBEAT_INTERVAL_SECONDS`: Ожидаемый интервал heartbeat от клиента.
*   `MAX_SAVEGAME_SIZE_BYTES`: Максимальный размер одного файла сохранения.
*   `USER_SAVEGAME_QUOTA_BYTES`: Общая квота на сохранения для одного пользователя.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`: Эндпоинт Jaeger.

### 7.2. Файлы Конфигурации (если применимо)
*   **`configs/library_config.yaml`**: Может использоваться для более сложных или менее часто изменяемых настроек.
    ```yaml
    server:
      http_port: ${LIBRARY_SERVICE_HTTP_PORT:-8084}
      grpc_port: ${LIBRARY_SERVICE_GRPC_PORT:-9094}

    database:
      postgres_dsn: ${POSTGRES_DSN}
      # Параметры пула соединений

    redis:
      address: ${REDIS_ADDR}
      # ...

    s3_storage:
      bucket_savegames: ${S3_BUCKET_SAVEGAMES}
      presigned_url_ttl_seconds: ${S3_PRESIGNED_URL_TTL_SECONDS:-300}
      max_save_file_size_bytes: ${MAX_SAVEGAME_SIZE_BYTES:-52428800} # 50MB
      user_total_quota_bytes: ${USER_SAVEGAME_QUOTA_BYTES:-1073741824} # 1GB

    playtime_tracking:
      heartbeat_interval_seconds: ${PLAYTIME_HEARTBEAT_INTERVAL_SECONDS:-60}
      session_timeout_inactive_seconds: 300 # Если нет heartbeat в течение 5 минут

    # Настройки для достижений (если есть специфичные)
    achievements:
      default_unlock_validation_rules: {} # TODO: Определить, если есть серверная валидация

    # Настройки для списка желаемого
    wishlist:
      max_items_per_user: 100
    ```

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   REST API: Используются стандартные коды состояния HTTP. Тело ответа об ошибке соответствует формату, определенному в `project_api_standards.md` (см. секцию 3.1).
*   gRPC API: Используются стандартные коды состояния gRPC.
*   Использование таймаутов, retry с экспоненциальной задержкой, Circuit Breaker при взаимодействии с другими сервисами.
*   Идемпотентность обработчиков событий Kafka, использование DLQ (Dead Letter Queue) для неразрешимых ошибок.

### 8.2. Распространенные Коды Ошибок (для REST API)
*   **`400 Bad Request` (`VALIDATION_ERROR`)**: Некорректные входные данные.
*   **`401 Unauthorized` (`UNAUTHENTICATED`)**: Ошибка аутентификации.
*   **`403 Forbidden` (`PERMISSION_DENIED`, `FAMILY_ACCESS_RESTRICTED`)**: Недостаточно прав или доступ ограничен (например, игра скрыта настройками семейного доступа).
*   **`404 Not Found` (`LIBRARY_ITEM_NOT_FOUND`, `GAME_NOT_FOUND_IN_WISHLIST`, `SAVEGAME_SLOT_NOT_FOUND`)**: Запрашиваемый ресурс не найден.
*   **`409 Conflict` (`ACHIEVEMENT_ALREADY_UNLOCKED`, `SAVEGAME_VERSION_CONFLICT`)**: Конфликт состояния (например, попытка разблокировать уже разблокированное достижение, конфликт версий сохранения).
*   **`413 Payload Too Large` (`SAVEGAME_FILE_TOO_LARGE`)**: Файл сохранения превышает допустимый размер.
*   **`422 Unprocessable Entity` (`USER_QUOTA_EXCEEDED`)**: Превышена квота на облачные сохранения.
*   **`500 Internal Server Error` (`INTERNAL_ERROR`)**: Внутренняя ошибка сервера.
*   **`503 Service Unavailable` (`SERVICE_UNAVAILABLE`)**: Сервис временно недоступен или зависимость недоступна (например, S3).

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   Все запросы к API (REST, WebSocket, gRPC для клиентов) требуют JWT аутентификации пользователя.
*   Межсервисное gRPC взаимодействие защищено с помощью mTLS.

### 9.2. Авторизация
*   Пользователь имеет доступ только к своим данным библиотеки, игрового времени, достижений, списка желаемого и сохранений.
*   Для функций семейного доступа реализуются специфичные проверки прав на основе связей между аккаунтами.

### 9.3. Защита Данных
*   TLS 1.2+ для всех внешних и внутренних коммуникаций.
*   Шифрование секретов и конфигураций.
*   Для игровых сохранений в S3:
    *   Использование pre-signed URL для загрузки/скачивания, ограничивающих доступ по времени и операции.
    *   Возможность шифрования на стороне сервера (SSE-S3) или шифрования на стороне клиента перед загрузкой (если требуется повышенная приватность, но усложняет управление).
*   Валидация входных данных для предотвращения инъекций и других атак.
*   Защита от подделки информации о достижениях или игровом времени (например, серверная валидация критических достижений, анализ аномального игрового времени).

### 9.4. Управление Секретами
*   Пароли к БД, ключи для S3, секреты для Kafka должны храниться в Kubernetes Secrets или HashiCorp Vault.

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Многоэтапная сборка для Go.
*   **Kubernetes манифесты/Helm-чарты.**
*   (Ссылка на `project_deployment_standards.md`).

### 10.2. Зависимости при Развертывании
*   PostgreSQL, Redis, Kafka, S3-совместимое хранилище.
*   Доступность Auth Service, Catalog Service, Download Service, Notification Service, Account Service, API Gateway.

### 10.3. CI/CD
*   Автоматизированная сборка, юнит- и интеграционное тестирование.
*   Развертывание в окружения с использованием GitOps.
*   Процедуры миграции схемы БД (например, golang-migrate).
*   Процедуры резервного копирования и восстановления для PostgreSQL, Redis и S3.

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат:** Структурированные JSON логи (Zap).
*   **Ключевые события:** Операции с библиотекой, начало/конец игровых сессий, разблокировка достижений, операции со списком желаемого, загрузка/скачивание сохранений, ошибки.
*   **Интеграция:** С централизованной системой логирования (Loki, ELK Stack).
*   (Ссылка на `project_observability_standards.md`).

### 11.2. Мониторинг
*   **Метрики (Prometheus):**
    *   Количество запросов к API (по типам, статусам).
    *   Длительность обработки запросов.
    *   Количество активных игровых сессий.
    *   Количество и объем синхронизаций сохранений.
    *   Частота разблокировки достижений.
    *   Ошибки при работе с БД, Redis, S3, Kafka.
*   **Дашборды (Grafana):** Для визуализации метрик.
*   **Алертинг (AlertManager):** Для критических ошибок, проблем с зависимостями, аномалий (например, резкое падение числа активных сессий).
*   (Ссылка на `project_observability_standards.md`).

### 11.3. Трассировка
*   **Интеграция:** OpenTelemetry SDK для Go. Экспорт трейсов в Jaeger или Tempo.
*   Трассировка всех входящих запросов и исходящих вызовов.
*   (Ссылка на `project_observability_standards.md`).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:**
    *   REST API для получения библиотеки пользователя (первая страница): P95 < 200 мс.
    *   Запись игровой сессии (heartbeat): P99 < 50 мс.
    *   Разблокировка достижения: P99 < 100 мс.
    *   Загрузка/скачивание файла сохранения (генерация URL): P95 < 150 мс (без учета времени передачи файла в/из S3).
*   **Масштабируемость:**
    *   Поддержка до 10 миллионов пользователей с активными библиотеками.
    *   Поддержка до 100 миллионов записей об играх в библиотеках.
    *   Поддержка до 1 миллиона одновременных игровых сессий.
    *   Объем хранилища для сохранений: до 1 ПБ с возможностью дальнейшего роста.
*   **Надежность:**
    *   Доступность сервиса: 99.99%.
    *   Сохранность данных библиотеки, достижений, игрового времени: RPO = 0 (без потерь при сбое).
    *   Сохранность игровых сохранений в S3: обеспечивается S3 (например, 99.999999999% durability). RPO для метаданных сохранений < 5 минут.
    *   RTO < 30 минут.
*   **Безопасность:** Защита пользовательских данных, включая приватность игровых сессий и сохранений. Предотвращение читерства с достижениями (базовая серверная валидация, где возможно).

## 13. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API, Protobuf определения для gRPC API, и форматы сообщений WebSocket будут храниться в соответствующих репозиториях или артефактах CI/CD.
*   Полные DDL схемы базы данных будут поддерживаться в актуальном состоянии в системе миграций.
*   TODO: Добавить ссылки на репозиторий с OpenAPI/Protobuf и на систему управления миграциями БД.

---
*Этот документ является основной спецификацией для Library Service и должен поддерживаться в актуальном состоянии.*

## 14. Связанные Рабочие Процессы (Related Workflows)
*   [Game Purchase and Library Update](../../../project_workflows/game_purchase_flow.md)
*   [Cloud Save Synchronization](../../../project_workflows/cloud_save_sync_flow.md) (TODO: Создать этот документ)
*   [Achievement Unlocking Flow](../../../project_workflows/achievement_unlocking_flow.md) (TODO: Создать этот документ)
