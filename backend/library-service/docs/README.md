# Спецификация Микросервиса: Library Service (Сервис Библиотеки Пользователя)

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение сервиса:** Library Service предназначен для управления личными библиотеками игр и приложений пользователей платформы "Российский Аналог Steam". Это включает в себя отслеживание принадлежащих пользователю продуктов, их статуса установки, игрового времени, достижений, управление списком желаемого, синхронизацию игровых сохранений в облаке и пользовательских настроек игр.
*   **Роль в общей архитектуре платформы:** Является центральным компонентом для хранения и управления данными, связанными с игровой активностью и владением цифровыми продуктами конкретного пользователя. Предоставляет эту информацию как самому пользователю через клиентские приложения, так и другим микросервисам (например, Social Service для отображения достижений, Download Service для управления загрузками).
*   **Основные бизнес-задачи:**
    *   Обеспечение пользователям доступа к приобретенным играм и приложениям.
    *   Отслеживание и отображение игрового времени и прогресса достижений.
    *   Предоставление функционала списка желаемого.
    *   Обеспечение надежной синхронизации игровых сохранений и пользовательских настроек игр между устройствами пользователя.
    *   Управление пользовательскими коллекциями игр и функцией "Семейный доступ".
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Управление библиотекой игр (Entitlement Management):**
    *   Добавление продуктов (игры, DLC) в библиотеку пользователя после покупки (через событие от Payment Service) или активации ключа.
    *   Отображение списка принадлежащих пользователю продуктов с их метаданными (получаемыми из Catalog Service).
    *   Фильтрация и сортировка игр в библиотеке (по названию, дате добавления, последнему запуску, жанру и т.д.).
    *   Создание и управление пользовательскими категориями/коллекциями для организации игр в библиотеке.
    *   Скрытие игр из основного вида библиотеки (без удаления права владения).
    *   Отслеживание статуса установки игры на текущем устройстве (интеграция с клиентским приложением и Download Service).
    *   (Опционально, [Family Sharing: функционал будет уточнен]) Реализация функции "Семейный доступ" (Family Sharing) для предоставления временного доступа к играм другим пользователям согласно правилам платформы.
*   **Отслеживание игрового времени (Playtime Tracking):**
    *   Автоматическая регистрация начала и окончания игровых сессий через API, вызываемое игровым клиентом или лаунчером.
    *   Периодические "heartbeat" сигналы от клиента для подтверждения активной игровой сессии.
    *   Подсчет общего игрового времени для каждого продукта и отображение статистики (например, "последний запуск", "часов за последние 2 недели").
    *   Синхронизация игрового времени между различными устройствами пользователя.
*   **Управление достижениями (Achievement Progress):**
    *   Регистрация полученных пользователем достижений (анлоков) и времени их получения.
    *   Хранение прогресса по частично выполненным достижениям (если игра предоставляет такую информацию).
    *   Отображение списка всех доступных и полученных достижений для каждой игры (метаданные достижений из Catalog Service).
    *   Публикация событий о разблокировке достижений для интеграции с Social Service (лента активности, уведомления друзьям).
*   **Управление списком желаемого (Wishlist Management):**
    *   Добавление и удаление продуктов из списка желаемого пользователя.
    *   Возможность установки приоритета или заметок для элементов списка желаемого.
    *   Интеграция с Notification Service для информирования пользователя о скидках на игры из его списка желаемого.
*   **Синхронизация игровых сохранений (Cloud Saves):**
    *   Предоставление API для игровых клиентов для загрузки файлов сохранений в облачное S3-хранилище.
    *   Скачивание последних (или выбранных пользователем) сохранений на другие устройства пользователя.
    *   Реализация стратегии разрешения конфликтов версий сохранений (например, "выбрать локальное", "выбрать облачное", "сохранить оба с разными именами", или автоматическое на основе временных меток).
    *   Версионирование файлов сохранений (опционально, [Cloud Save Versioning: глубина версионирования будет уточнена]).
    *   Управление квотами на облачное хранилище для сохранений на пользователя или на игру.
*   **Управление Пользовательскими Настройками Игр (User-Specific Game Settings Storage):**
    *   Хранение и синхронизация пользовательских настроек для конкретных игр (например, настройки графики, управления, звука), если игра поддерживает такую интеграцию.
    *   API для игр для сохранения и загрузки этих настроек.

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+, согласно `../../../../project_technology_stack.md`).
*   **Веб-фреймворк (REST API):** Echo (`github.com/labstack/echo/v4`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **RPC фреймворк (gRPC):** `google.golang.org/grpc` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **WebSocket:** `github.com/gorilla/websocket` (или аналогичная библиотека Go) для real-time уведомлений (например, о разблокировке достижений).
*   **База данных (основная):** PostgreSQL (версия 15+) для хранения информации о библиотеках, игровом времени, достижениях, списках желаемого, метаданных сохранений и пользовательских настройках игр. Драйвер: GORM (`gorm.io/gorm`) с `gorm.io/driver/postgres` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Кэширование/Оперативные данные:** Redis (версия 7.0+) для кэширования часто запрашиваемых данных (например, содержимое библиотеки пользователя, списки желаемого) и оперативных данных (например, текущие активные игровые сессии, прогресс достижений перед записью в БД). Клиент: `go-redis/redis` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Облачное хранилище (для сохранений):** S3-совместимое объектное хранилище (например, MinIO, Yandex Object Storage) (согласно `../../../../project_technology_stack.md`).
*   **Брокер сообщений:** Apache Kafka (клиент `github.com/confluentinc/confluent-kafka-go`, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Логирование:** Zap (`go.uber.org/zap`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Мониторинг/Трассировка:** OpenTelemetry SDK, Prometheus client (`github.com/prometheus/client_golang`). (согласно `../../../../project_observability_standards.md`).
*   **Инфраструктура:** Docker, Kubernetes, Helm.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   **Библиотека (Library):** Персональная коллекция игр и других цифровых продуктов, принадлежащих пользователю или доступных ему (например, через Family Sharing).
*   **Право Владения (Entitlement):** Запись, подтверждающая право пользователя на доступ к определенному продукту.
*   **Игровое время (Playtime):** Общее время, проведенное пользователем в конкретной игре.
*   **Достижение (Achievement):** Виртуальная награда, выдаваемая пользователю за выполнение определенных условий или задач в игре.
*   **Список желаемого (Wishlist):** Персональный список продуктов, которые пользователь хочет приобрести в будущем.
*   **Игровое сохранение (Savegame/Cloud Save):** Файл или набор файлов, содержащих прогресс пользователя в игре, который может быть синхронизирован с облачным хранилищем.
*   **Семейный доступ (Family Sharing):** Функция, позволяющая делиться играми из своей библиотеки с членами семьи или близкими друзьями.
*   **Слот сохранения (Save Slot):** Именованная ячейка для игрового сохранения, позволяющая иметь несколько параллельных прогрессов.
*   Для других общих терминов см. `../../../../project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Library Service будет реализован с использованием принципов Чистой Архитектуры (Clean Architecture) для достижения слабой связанности, высокой тестируемости и гибкости системы.
*   Сервис будет состоять из нескольких ключевых модулей, отвечающих за отдельные аспекты функциональности: управление библиотекой, отслеживание игрового времени, управление достижениями, список желаемого, синхронизация сохранений и управление настройками игр.

### 2.2. Диаграмма Архитектуры (Clean Architecture)
```mermaid
graph TD
    subgraph UserClientsAndOtherServices ["Клиенты Пользователя и Другие Сервисы"]
        UserClient["Клиент Пользователя (Веб, Десктоп, Мобильный)"]
        GameClient["Игровой Клиент (запущенная игра)"]
        APIGateway["API Gateway"]
        OtherInternalServices["Другие Микросервисы (Catalog, Payment, Social, Auth, Download)"]
    end

    subgraph LibraryService ["Library Service (Чистая Архитектура)"]
        direction TB

        subgraph PresentationLayer [Presentation Layer (Адаптеры Транспорта)]
            REST_API[REST API (Echo) - для клиента платформы]
            GRPC_API[gRPC API - для игрового клиента и других сервисов]
            WebSocket_Handler[WebSocket Handler - для real-time уведомлений (например, ачивки)]
            KafkaConsumers[Kafka Consumers - для входящих событий (например, покупка игры)]
        end

        subgraph ApplicationLayer [Application Layer (Сценарии Использования)]
            LibraryAppSvc["Управление Библиотекой (добавление, просмотр, категории)"]
            PlaytimeAppSvc["Учет Игрового Времени (старт/стоп сессии, heartbeat)"]
            AchievementAppSvc["Управление Достижениями (анлок, прогресс)"]
            WishlistAppSvc["Управление Списком Желаемого"]
            CloudSaveAppSvc["Синхронизация Игровых Сохранений (загрузка, скачивание, конфликты)"]
            UserGameSettingsAppSvc["Управление Настройками Игр Пользователя"]
            FamilySharingAppSvc["Управление Семейным Доступом"]
        end

        subgraph DomainLayer [Domain Layer (Бизнес-логика и Сущности)]
            Entities["Сущности (UserLibraryItem, PlaytimeSession, UserAchievement, WishlistItem, SavegameMetadata, UserGameSetting, FamilyLink)"]
            Aggregates["Агрегаты (UserLibrary, UserGameProfile)"]
            DomainEvents["Доменные События (GameAddedToLibrary, AchievementUnlocked, PlaytimeUpdated, SavegameSynced)"]
            RepositoryIntf["Интерфейсы Репозиториев (PostgreSQL, Redis)"]
            DomainServices["Доменные Сервисы (PlaytimeCalculator, AchievementUnlocker, CloudSaveConflictResolver)"]
        end

        subgraph InfrastructureLayer [Infrastructure Layer (Внешние Зависимости и Реализации)"]
            PostgresAdapter["Адаптер PostgreSQL (реализация репозиториев)"]
            RedisAdapter["Адаптер Redis (кэш, активные сессии, очереди задач)"]
            S3Adapter["Адаптер S3-хранилища (для облачных сохранений)"]
            KafkaProducer["Продюсер Kafka (исходящие события)"]
            InternalServiceClients["Клиенты других микросервисов (Catalog, Auth, Payment, Download, Social, Notification)"]
            Config["Конфигурация (Viper)"]
            Logging["Логирование (Zap)"]
        end

        PresentationLayer --> ApplicationLayer
        ApplicationLayer --> DomainLayer
        ApplicationLayer --> InfrastructureLayer
        DomainLayer ----> RepositoryIntf
        InfrastructureLayer -- Implements --> RepositoryIntf
    end

    UserClient -- HTTP/WebSocket --> APIGateway
    GameClient -- gRPC/HTTP --> APIGateway
    APIGateway -- HTTP/WebSocket/gRPC --> PresentationLayer

    OtherInternalServices -- gRPC / Kafka --> PresentationLayer

    PostgresAdapter --> DB[(PostgreSQL)]
    RedisAdapter --> Cache[(Redis)]
    S3Adapter --> S3Storage[("S3 Cloud Storage")]
    KafkaProducer --> KafkaBroker[Kafka Message Bus]
    InternalServiceClients --> ExtServices[("Внешние gRPC/REST Сервисы")]


    classDef layer_boundary fill:#f9f9f9,stroke:#333,stroke-width:2px,color:#333
    classDef component_major fill:#e6f0ff,stroke:#007bff,color:#000
    classDef component_minor fill:#d4edda,stroke:#28a745,color:#000
    classDef datastore fill:#f8d7da,stroke:#dc3545,color:#000
    classDef external_actor fill:#FEF9E7,stroke:#F1C40F,color:#000

    class PresentationLayer,ApplicationLayer,DomainLayer,InfrastructureLayer layer_boundary
    class REST_API,GRPC_API,WebSocket_Handler,KafkaConsumers,LibraryAppSvc,PlaytimeAppSvc,AchievementAppSvc,WishlistAppSvc,CloudSaveAppSvc,UserGameSettingsAppSvc,FamilySharingAppSvc,Entities,Aggregates,DomainEvents,RepositoryIntf,DomainServices component_major
    class PostgresAdapter,RedisAdapter,S3Adapter,KafkaProducer,InternalServiceClients,Config,Logging component_minor
    class DB,Cache,S3Storage,KafkaBroker,ExtServices datastore
    class UserClient,GameClient,APIGateway,OtherInternalServices external_actor
```

### 2.3. Слои Сервиса
(Описания слоев аналогичны предыдущим сервисам, с акцентом на специфику Library Service: управление библиотекой, игровым временем, достижениями, списком желаемого, облачными сохранениями и настройками игр.)

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL:** `/api/v1/library` (маршрутизируется через API Gateway).
*   **Аутентификация:** JWT Bearer Token (`X-User-Id` извлекается).
*   **Авторизация:** Все эндпоинты требуют аутентификации и проверяют, что пользователь имеет доступ только к своим данным (например, `user_self_only`).
*   **Формат ответа об ошибке:** Согласно `../../../../project_api_standards.md`.

#### 3.1.1. Библиотека Игр (User Library)
*   **`GET /me/items`**
    *   Описание: Получение списка игр в библиотеке текущего пользователя.
    *   Query параметры: `status` (installed, not_installed), `category_id`, `search` (по названию), `sort_by` (name_asc, last_played_desc, added_at_desc), `page`, `limit`.
    *   Ответ: (Массив `UserLibraryItem` с основной информацией о продукте из Catalog Service и пользовательскими данными).
*   **`PATCH /me/items/{library_item_id}`**
    *   Описание: Изменение атрибутов элемента библиотеки (например, добавление в категорию, скрытие, статус установки).
    *   Тело запроса: `{"data": {"type": "libraryItemUpdate", "attributes": {"is_hidden": true, "category_ids": ["uuid1", "uuid2"], "installation_status": "installed"}}}`
    *   Ответ: (Обновленный `UserLibraryItem`).
*   **`GET /me/categories`**
    *   Описание: Получение списка пользовательских категорий для игр.
    *   Ответ: (Массив `UserGameCategory`).
*   **`POST /me/categories`**
    *   Описание: Создание новой пользовательской категории.
    *   Тело запроса: `{"data": {"type": "userGameCategoryCreation", "attributes": {"name": "Избранное"}}}`
    *   Ответ: (Созданный `UserGameCategory`).

#### 3.1.2. Игровое Время (Playtime)
*   **`GET /me/playtime`**
    *   Описание: Получение агрегированной статистики игрового времени по всем или конкретным играм.
    *   Query параметры: `product_ids` (comma-separated), `period` (all_time, last_2_weeks).
    *   Ответ: (Статистика игрового времени).
*   **`POST /me/playtime/sessions/heartbeat`**
    *   Описание: Игровой клиент отправляет "heartbeat" для активной игровой сессии, чтобы подтвердить присутствие пользователя и обновить игровое время. Auth Service может использовать это для проверки активности сессии.
    *   Тело запроса: `{"data": {"type": "playtimeHeartbeat", "attributes": {"productId": "game-uuid-123", "sessionId": "session-uuid-abc", "currentTimestamp": "ISO8601"}}}`
    *   Ответ: (200 OK, опционально обновленное состояние сессии).

#### 3.1.3. Достижения (Achievements)
*   **`GET /me/achievements/status`**
    *   Описание: Получение статуса всех достижений для указанных продуктов.
    *   Query параметры: `product_ids` (comma-separated).
    *   Ответ: (Массив статусов достижений).

#### 3.1.4. Список Желаемого (Wishlist)
*   **`GET /me/wishlist`**
    *   Описание: Получение списка желаемого текущего пользователя.
    *   Ответ: (Массив `WishlistItem` с информацией о продуктах из Catalog Service).
*   **`POST /me/wishlist`**
    *   Описание: Добавление продукта в список желаемого.
    *   Тело запроса: `{"data": {"type": "wishlistItemCreation", "attributes": {"productId": "game-uuid-789"}}}`
    *   Ответ: (Созданный `WishlistItem`).
*   **`DELETE /me/wishlist/{product_id}`**
    *   Описание: Удаление продукта из списка желаемого.
    *   Ответ: (204 No Content).

#### 3.1.5. Игровые Сохранения (Cloud Saves)
*   **`GET /me/savegames/{product_id}`**
    *   Описание: Получение списка метаданных доступных облачных сохранений для игры.
    *   Ответ: (Массив `SavegameMetadata`).
*   **`POST /me/savegames/{product_id}/slots/{slot_name}/upload-url`**
    *   Описание: Запрос pre-signed URL для загрузки файла сохранения в S3.
    *   Тело запроса: `{"data": {"type": "savegameUploadRequest", "attributes": {"fileSizeBytes": 1048576, "fileHashSha256": "hash_value"}}}`
    *   Ответ: (Pre-signed URL и `savegameId`).
*   **`POST /me/savegames/confirm-upload`**
    *   Описание: Подтверждение успешной загрузки файла сохранения в S3.
    *   Тело запроса: `{"data": {"type": "savegameUploadConfirmation", "attributes": {"savegameId": "savegame-uuid-1", "s3ObjectKey": "path/to/save.dat", "clientModifiedAt": "ISO8601"}}}`
    *   Ответ: (Обновленный `SavegameMetadata`).
*   **`GET /me/savegames/{savegame_id}/download-url`**
    *   Описание: Запрос pre-signed URL для скачивания файла сохранения из S3.
    *   Ответ: (Pre-signed URL).

#### 3.1.6. Пользовательские Настройки Игр
*   **`GET /me/game-settings/{product_id}`**
    *   Описание: Получение пользовательских настроек для конкретной игры.
    *   Ответ: (Объект `GameSpecificSettings`).
*   **`PUT /me/game-settings/{product_id}`**
    *   Описание: Обновление пользовательских настроек для игры.
    *   Тело запроса: `{"data": {"type": "gameSettingsUpdate", "attributes": {"settingsPayload": {"graphics": "ultra", "sound": 0.7}}}}`
    *   Ответ: (Обновленный `GameSpecificSettings`).

### 3.2. gRPC API (для игрового клиента и межсервисного взаимодействия)
*   **Пакет:** `library.v1`
*   **Сервис:** `PlaytimeTrackerService`
    *   `rpc StartPlaytimeSession(StartPlaytimeSessionRequest) returns (StartPlaytimeSessionResponse)`
    *   `rpc RecordPlaytimeHeartbeat(RecordPlaytimeHeartbeatRequest) returns (google.protobuf.Empty)`
    *   `rpc EndPlaytimeSession(EndPlaytimeSessionRequest) returns (EndPlaytimeSessionResponse)`
*   **Сервис:** `AchievementService`
    *   `rpc UpdateAchievementProgress(UpdateAchievementProgressRequest) returns (UpdateAchievementProgressResponse)` (включая полный анлок)
*   **Сервис:** `CloudSaveService`
    *   `rpc RequestSaveGameUploadURL(RequestSaveGameUploadURLRequest) returns (RequestSaveGameUploadURLResponse)`
    *   `rpc ConfirmSaveGameUpload(ConfirmSaveGameUploadRequest) returns (ConfirmSaveGameUploadResponse)`
    *   `rpc ListSaveGames(ListSaveGamesRequest) returns (ListSaveGamesResponse)`
    *   `rpc RequestSaveGameDownloadURL(RequestSaveGameDownloadURLRequest) returns (RequestSaveGameDownloadURLResponse)`
*   **Сервис:** `GameSettingsSyncService`
    *   `rpc GetGameSettings(GetGameSettingsRequest) returns (GetGameSettingsResponse)`
    *   `rpc UpdateGameSettings(UpdateGameSettingsRequest) returns (UpdateGameSettingsResponse)`

### 3.3. WebSocket API
*   **Эндпоинт:** `/ws/library/notifications` (требует аутентификации).
*   **Сообщения от сервера к клиенту:**
    *   **`achievementUnlocked`**: `{"type": "achievementUnlocked", "payload": {"productId": "...", "achievementId": "...", "achievementName": "...", "timestamp": "..."}}`
    *   **`cloudSaveConflict`**: `{"type": "cloudSaveConflict", "payload": {"productId": "...", "message": "Обнаружен конфликт сохранений. Выберите версию для использования."}}`
    *   **`cloudSaveSynced`**: `{"type": "cloudSaveSynced", "payload": {"productId": "...", "slotName":"...", "timestamp":"..."}}`
    *   **`wishlistGameOnSale`**: `{"type": "wishlistGameOnSale", "payload": {"productId": "...", "productName":"...", "discountPercent":"30%"}}`

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`UserLibraryItem` (Элемент Библиотеки Пользователя)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK). **Обязательность: Да.**
    *   `product_id` (UUID, FK). **Обязательность: Да.**
    *   `acquisition_type` (ENUM: `purchase`, `gift`, `free`, `family_shared`). **Обязательность: Да.**
    *   `acquisition_timestamp` (TIMESTAMPTZ). **Обязательность: Да (DEFAULT now()).**
    *   `last_played_timestamp` (TIMESTAMPTZ, Nullable). **Обязательность: Нет.**
    *   `total_playtime_seconds` (BIGINT, Default: 0). **Обязательность: Да (DEFAULT 0).**
    *   `installation_status` (ENUM: `not_installed`, `installing`, `installed`, `update_required`, `error`). **Обязательность: Да (DEFAULT 'not_installed').**
    *   `installed_version_id` (UUID, Nullable). **Обязательность: Нет.**
    *   `is_hidden` (BOOLEAN, Default: false). **Обязательность: Да (DEFAULT FALSE).**
    *   `custom_categories` (ARRAY_TEXT, Nullable). **Обязательность: Нет.**
    *   `added_to_library_at` (TIMESTAMPTZ). **Обязательность: Да (DEFAULT now()).**
*   **`PlaytimeSession` (Игровая Сессия)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `user_library_item_id` (UUID, FK). **Обязательность: Да.**
    *   `user_id` (UUID). **Обязательность: Да.**
    *   `product_id` (UUID). **Обязательность: Да.**
    *   `start_timestamp` (TIMESTAMPTZ). **Обязательность: Да.**
    *   `end_timestamp` (TIMESTAMPTZ, Nullable). **Обязательность: Нет (пока сессия активна).**
    *   `duration_seconds` (INTEGER, Nullable). **Обязательность: Нет (рассчитывается при завершении).**
    *   `last_heartbeat_timestamp` (TIMESTAMPTZ). **Обязательность: Да (для активных сессий).**
*   **`UserAchievement` (Достижение Пользователя)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK). **Обязательность: Да.**
    *   `product_id` (UUID, FK). **Обязательность: Да.**
    *   `achievement_api_name` (VARCHAR, FK на AchievementMeta в Catalog Service). **Обязательность: Да.**
    *   `is_unlocked` (BOOLEAN, Default: false). **Обязательность: Да (DEFAULT FALSE).**
    *   `unlocked_timestamp` (TIMESTAMPTZ, Nullable). **Обязательность: Нет.**
    *   `progress_percentage` (INTEGER, Nullable, 0-100). **Обязательность: Нет.**
*   **`WishlistItem` (Элемент Списка Желаемого)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK). **Обязательность: Да.**
    *   `product_id` (UUID, FK). **Обязательность: Да.**
    *   `added_timestamp` (TIMESTAMPTZ). **Обязательность: Да (DEFAULT now()).**
    *   `priority` (INTEGER, Nullable). **Обязательность: Нет.**
    *   `notes` (TEXT, Nullable). **Обязательность: Нет.**
*   **`SavegameMetadata` (Метаданные Игрового Сохранения)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK). **Обязательность: Да.**
    *   `product_id` (UUID, FK). **Обязательность: Да.**
    *   `slot_name` (VARCHAR). **Обязательность: Да.**
    *   `s3_object_key` (VARCHAR, UK). **Обязательность: Да.**
    *   `file_size_bytes` (BIGINT). **Обязательность: Да.**
    *   `file_hash_sha256` (VARCHAR). **Обязательность: Да.**
    *   `client_modified_timestamp` (TIMESTAMPTZ). **Обязательность: Да.**
    *   `server_uploaded_timestamp` (TIMESTAMPTZ). **Обязательность: Да (DEFAULT now()).**
    *   `version_number` (INTEGER, Default: 1). **Обязательность: Да (для версионирования).**
    *   `tags` (JSONB, Nullable). **Обязательность: Нет.**
*   **`UserGameSetting` (Пользовательские Настройки Игры)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID, FK). **Обязательность: Да.**
    *   `product_id` (UUID, FK). **Обязательность: Да.**
    *   `settings_payload` (JSONB). **Обязательность: Да.**
    *   `last_synced_timestamp` (TIMESTAMPTZ). **Обязательность: Да (DEFAULT now()).**
    *   UNIQUE (`user_id`, `product_id`).
*   **`FamilyLink` (Связь Семейного Доступа)**
    *   `link_id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `owner_user_id` (UUID, FK). **Обязательность: Да.**
    *   `shared_user_id` (UUID, FK). **Обязательность: Да.**
    *   `status` (ENUM: `pending_approval`, `active`, `revoked`). **Обязательность: Да.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (DEFAULT now()).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (DEFAULT now()).**
    *   UNIQUE (`owner_user_id`, `shared_user_id`).

### 4.2. Схема Базы Данных (PostgreSQL)
(ERD и DDL для PostgreSQL, описание Redis и S3 структуры как в предыдущем моем ответе).
*Примечание: DDL для всех перечисленных выше сущностей должен быть создан или проверен на соответствие. Примеры DDL для ключевых таблиц, таких как `user_library_items`, `playtime_sessions`, `user_achievements`, `wishlist_items`, `savegame_metadata`, `user_game_settings`, `family_links`, должны быть детализированы в этом разделе, следуя структуре других сервисных документов.*

## 5. Потоковая Обработка Событий (Event Streaming)
*   **Формат событий:** CloudEvents v1.0 JSON (согласно `../../../../project_api_standards.md`).
*   **Основной топик для публикуемых событий:** `com.platform.library.events.v1`.

### 5.1. Публикуемые События (Produced Events)
*   **`com.platform.library.item.added.v1`**
    *   Описание: Продукт добавлен в библиотеку пользователя.
    *   `data` Payload: `{"userId": "user-uuid", "productId": "prod-uuid", "acquisitionType": "purchase", "timestamp": "ISO8601"}`
*   **`com.platform.library.playtime.updated.v1`**
    *   Описание: Обновлено игровое время для продукта.
    *   `data` Payload: `{"userId": "user-uuid", "productId": "prod-uuid", "totalPlaytimeSeconds": 7200, "lastPlayedTimestamp": "ISO8601"}`
*   **`com.platform.library.achievement.unlocked.v1`**
    *   Описание: Пользователь разблокировал достижение.
    *   `data` Payload: `{"userId": "user-uuid", "productId": "prod-uuid", "achievementApiName": "FIRST_BLOOD", "unlockedTimestamp": "ISO8601"}`
*   **`com.platform.library.savegame.uploaded.v1`**
    *   Описание: Новое игровое сохранение загружено в облако.
    *   `data` Payload: `{"userId": "user-uuid", "productId": "prod-uuid", "savegameId": "save-uuid", "slotName": "slot1", "timestamp": "ISO8601"}`
*   **`com.platform.library.wishlist.item.added.v1`**
    *   Описание: Продукт добавлен в список желаемого.
    *   `data` Payload: `{"userId": "user-uuid", "productId": "prod-uuid", "timestamp": "ISO8601"}`
*   **`com.platform.library.wishlist.item.removed.v1`**
    *   Описание: Продукт удален из списка желаемого.
    *   `data` Payload: `{"userId": "user-uuid", "productId": "prod-uuid", "timestamp": "ISO8601"}`

### 5.2. Потребляемые События (Consumed Events)
*   **`com.platform.payment.transaction.completed.v1`** (от Payment Service)
    *   Описание: Успешная транзакция покупки.
    *   Ожидаемый `data` Payload: `{"userId": "user-uuid", "orderId": "order-uuid", "items": [{"productId": "prod-uuid", "productType": "game"}], "timestamp": "ISO8601"}`
    *   Логика обработки: Добавить купленные продукты в библиотеку пользователя.
*   **`com.platform.catalog.product.updated.v1`** (от Catalog Service)
    *   Описание: Обновлены метаданные продукта (например, название, иконки достижений).
    *   Ожидаемый `data` Payload: `{"productId": "prod-uuid", "updatedFields": ["name", "achievement_definitions"], ...}`
    *   Логика обработки: Обновить кэшированную информацию о продукте в Library Service, если это необходимо для отображения в библиотеке или для достижений.
*   **`com.platform.social.review.submitted.v1`** (от Social Service)
    *   Описание: Пользователь оставил отзыв на игру.
    *   Ожидаемый `data` Payload: `{"userId": "user-uuid", "productId": "prod-uuid", "rating": 5, "reviewId": "review-uuid"}`
    *   Логика обработки: Может использоваться для отображения информации о том, оставлял ли пользователь отзыв на игру в своей библиотеке (опционально).

## 6. Интеграции (Integrations)
(Описание интеграций как в предыдущем моем ответе).

## 7. Конфигурация (Configuration)
(Описание файла конфигурации и переменных окружения как в предыдущем моем ответе).

## 8. Обработка Ошибок (Error Handling)
*   Стандартные ошибки API согласно `../../../../project_api_standards.md`.
*   **Специфичные коды ошибок:**
    *   `LIBRARY_ITEM_NOT_FOUND`: Элемент библиотеки не найден.
    *   `ACHIEVEMENT_NOT_FOUND`: Метаданные достижения не найдены.
    *   `WISHLIST_ITEM_ALREADY_EXISTS`: Продукт уже в списке желаемого.
    *   `SAVEGAME_SLOT_NOT_FOUND`: Указанный слот сохранения не найден.
    *   `SAVEGAME_UPLOAD_FAILED`: Ошибка при загрузке сохранения в S3.
    *   `SAVEGAME_DOWNLOAD_FAILED`: Ошибка при скачивании сохранения из S3.
    *   `SAVEGAME_CONFLICT_RESOLUTION_REQUIRED`: Требуется разрешение конфликта версий сохранений.
    *   `MAX_SAVEGAME_QUOTA_EXCEEDED`: Превышена квота на облачное хранилище.
    *   `FAMILY_SHARING_LINK_INVALID_STATE`: Недопустимая операция для текущего статуса семейного доступа.

## 9. Безопасность (Security)
(Описание безопасности как в предыдущем моем ответе, с акцентом на Cloud Saves).

## 10. Развертывание (Deployment)
(Как в предыдущем моем ответе).

## 11. Мониторинг и Логирование (Logging and Monitoring)
(Как в предыдущем моем ответе).

## 12. Нефункциональные Требования (NFRs)
(Как в предыдущем моем ответе).

## 13. Приложения (Appendices)
(Как в предыдущем моем ответе).

## 14. Пользовательские Сценарии (User Flows)

В этом разделе описаны ключевые пользовательские сценарии, связанные с Library Service.

### 14.1. Пользователь Просматривает Свою Библиотеку Игр и Запускает Игру
*   **Описание:** Пользователь открывает клиентское приложение, просматривает свою библиотеку игр и запускает одну из них.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Клиентское Приложение
        participant APIGW as API Gateway
        participant LibrarySvc as Library Service
        participant CatalogSvc as Catalog Service
        participant DownloadSvc as Download Service (для проверки установки/запуска)

        User->>ClientApp: Открывает раздел "Библиотека"
        ClientApp->>APIGW: GET /api/v1/library/me/items?sort_by=last_played_desc
        APIGW->>LibrarySvc: Forward request (с User JWT)
        LibrarySvc->>LibrarySvc: Получение списка UserLibraryItem из PostgreSQL
        LibrarySvc->>CatalogSvc: (gRPC) Запрос метаданных для списка ProductID
        CatalogSvc-->>LibrarySvc: Метаданные продуктов
        LibrarySvc-->>APIGW: HTTP 200 OK (список игр с метаданными)
        APIGW-->>ClientApp: HTTP 200 OK
        ClientApp-->>User: Отображение списка игр

        User->>ClientApp: Нажимает "Играть" на игре X
        ClientApp->>DownloadSvc: (через API или локальный IPC) Проверка статуса установки/обновления игры X
        alt Игра установлена и обновлена
            ClientApp->>LibrarySvc: (gRPC) POST /api/v1/library/playtime/sessions/start (productId)
            LibrarySvc-->>ClientApp: Ответ (sessionId)
            ClientApp->>ClientApp: Запуск игрового процесса (локально)
        else Игра не установлена или требует обновления
            ClientApp-->>User: Предложение установить/обновить игру
            Note over ClientApp, DownloadSvc: Запускается сценарий загрузки/обновления (см. User Flows Download Service)
        end
    ```

### 14.2. Игровой Клиент Сообщает об Игровом Времени
*   **Описание:** Запущенный игровой клиент периодически отправляет "heartbeat" и информацию о завершении сессии для учета игрового времени.
*   **Диаграмма:** (См. диаграмму "Playtime Tracking & Achievement Unlocking" в разделе 2 или адаптировать сюда)
    ```mermaid
    sequenceDiagram
        participant GameClient as Игровой Клиент
        participant LibrarySvc as Library Service (gRPC API)
        participant DB_PostgreSQL as PostgreSQL
        participant Cache_Redis as Redis (Active Sessions)
        participant Kafka as Kafka Message Bus

        GameClient->>LibrarySvc: RecordPlaytimeHeartbeatRequest (userId, productId, sessionId, timestamp)
        LibrarySvc->>Cache_Redis: Обновление времени последнего heartbeat для активной сессии
        LibrarySvc-->>GameClient: RecordPlaytimeHeartbeatResponse (OK)

        Note over GameClient, LibrarySvc: При завершении игры:
        GameClient->>LibrarySvc: EndPlaytimeSessionRequest (userId, productId, sessionId, endTime, duration)
        LibrarySvc->>DB_PostgreSQL: Обновление UserLibraryItem (total_playtime, last_played_at)
        LibrarySvc->>DB_PostgreSQL: Сохранение PlaytimeSession
        LibrarySvc->>Cache_Redis: Удаление активной сессии из кэша
        LibrarySvc->>Kafka: Publish `com.platform.library.playtime.session.ended.v1`
        LibrarySvc-->>GameClient: EndPlaytimeSessionResponse (OK)
    ```

### 14.3. Пользователь Разблокирует Достижение
*   **Описание:** Игровой клиент сообщает о разблокировке достижения. Library Service сохраняет это и уведомляет другие системы.
*   **Диаграмма:** (Часть диаграммы "Playtime Tracking & Achievement Unlocking" или адаптировать сюда)
    ```mermaid
    sequenceDiagram
        participant GameClient as Игровой Клиент
        participant LibrarySvc as Library Service (gRPC API)
        participant DB_PostgreSQL as PostgreSQL
        participant Kafka as Kafka Message Bus
        participant SocialSvc as Social Service (через Kafka)
        participant NotificationSvcWS as Notification Service (WebSocket)

        GameClient->>LibrarySvc: UpdateAchievementProgressRequest (userId, productId, achievementApiName, progress?, unlockedAt?)
        LibrarySvc->>DB_PostgreSQL: Сохранение/Обновление UserAchievement (is_unlocked=true, unlocked_at)
        alt Достижение действительно разблокировано
            LibrarySvc->>Kafka: Publish `com.platform.library.achievement.unlocked.v1` (userId, productId, achievementId, achievementName)
            LibrarySvc-->>NotificationSvcWS: (через WebSocket) Отправка уведомления клиенту о разблокировке
            Kafka-->>SocialSvc: Consume `achievement.unlocked` -> Публикация в ленте активности
        end
        LibrarySvc-->>GameClient: UpdateAchievementProgressResponse (OK, текущий статус ачивки)
    ```

### 14.4. Игровой Клиент Загружает Новое Сохранение в Облако
*   **Описание:** Игра автоматически или по команде пользователя загружает файл сохранения в облако.
*   **Диаграмма:** (Часть диаграммы "Cloud Save Synchronization" или адаптировать сюда)
    ```mermaid
    sequenceDiagram
        actor User
        participant GameClient as Игровой Клиент
        participant LibrarySvc as Library Service (REST/gRPC)
        participant S3Store as S3 Cloud Storage

        User->>GameClient: Сохраняет игру (например, F5)
        GameClient->>GameClient: Подготовка файла сохранения (save_slot_1.dat, hash)
        GameClient->>LibrarySvc: RequestSaveGameUploadURL(productId, slotName="slot_1", fileHash, fileSize)
        LibrarySvc->>LibrarySvc: Проверка квот, генерация pre-signed S3 URL
        LibrarySvc-->>GameClient: Pre-signed S3 URL, savegameId (для подтверждения)
        GameClient->>S3Store: PUT <pre-signed_url> (файл сохранения)
        S3Store-->>GameClient: HTTP 200 OK
        GameClient->>LibrarySvc: ConfirmSaveGameUpload(savegameId, s3Path, clientModifiedAt)
        LibrarySvc->>LibrarySvc: Обновление SavegameMetadata в PostgreSQL
        LibrarySvc->>Kafka: Publish `com.platform.library.savegame.uploaded.v1`
        LibrarySvc-->>GameClient: UploadConfirmed (OK)
    ```

### 14.5. Пользователь Устанавливает Игру на Новом Устройстве и Загружает Облачное Сохранение
*   **Описание:** После установки игры на новом устройстве, клиент игры запрашивает и скачивает последнее облачное сохранение.
*   **Диаграмма:** (Часть диаграммы "Cloud Save Synchronization" или адаптировать сюда)
    ```mermaid
    sequenceDiagram
        actor User
        participant GameClient as Игровой Клиент (на новом устройстве)
        participant LibrarySvc as Library Service (REST/gRPC)
        participant S3Store as S3 Cloud Storage

        User->>GameClient: Первый запуск игры X
        GameClient->>LibrarySvc: ListSaveGames(productId)
        LibrarySvc-->>GameClient: Список SavegameMetadata (включая самое последнее)
        alt Есть сохранения в облаке
            GameClient->>LibrarySvc: RequestSaveGameDownloadURL(latest_savegame_id)
            LibrarySvc-->>GameClient: Pre-signed S3 URL
            GameClient->>S3Store: GET <pre-signed_url>
            S3Store-->>GameClient: Файл сохранения
            GameClient->>GameClient: Загрузка сохранения в игру
            GameClient-->>User: Предложение загрузить сохранение / Автоматическая загрузка
        else Нет сохранений в облаке
            GameClient-->>User: Начало новой игры
        end
    ```
    *   **Разрешение конфликтов:** Если клиент обнаруживает локальное сохранение, которое новее облачного, или несинхронизированное локальное сохранение при наличии облачного, он должен предложить пользователю выбор:
        1.  Загрузить облачное сохранение (перезаписав локальное).
        2.  Загрузить локальное сохранение в облако (перезаписав облачное).
        3.  (Опционально) Сохранить обе версии или отменить синхронизацию.
        Выбор пользователя инициирует соответствующий поток загрузки или скачивания.

### 14.6. Пользователь Добавляет/Удаляет Игру из Списка Желаемого
*   **Описание:** Пользователь добавляет игру в свой список желаемого или удаляет ее оттуда.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Клиентское Приложение
        participant APIGW as API Gateway
        participant LibrarySvc as Library Service
        participant Kafka as Kafka Message Bus

        User->>ClientApp: Нажимает "Добавить в желаемое" на странице игры X
        ClientApp->>APIGW: POST /api/v1/library/me/wishlist (productId)
        APIGW->>LibrarySvc: Forward request
        LibrarySvc->>LibrarySvc: Добавление WishlistItem в PostgreSQL
        LibrarySvc->>Kafka: Publish `com.platform.library.wishlist.item.added.v1`
        LibrarySvc-->>APIGW: HTTP 201 Created (созданный WishlistItem)
        APIGW-->>ClientApp: HTTP 201 Created
        ClientApp-->>User: Игра добавлена в желаемое

        User->>ClientApp: Нажимает "Удалить из желаемого"
        ClientApp->>APIGW: DELETE /api/v1/library/me/wishlist/{productId}
        APIGW->>LibrarySvc: Forward request
        LibrarySvc->>LibrarySvc: Удаление WishlistItem из PostgreSQL
        LibrarySvc->>Kafka: Publish `com.platform.library.wishlist.item.removed.v1`
        LibrarySvc-->>APIGW: HTTP 204 No Content
        APIGW-->>ClientApp: HTTP 204 No Content
        ClientApp-->>User: Игра удалена из желаемого
    ```

## 15. Резервное Копирование и Восстановление (Backup and Recovery)
(Как в предыдущем моем ответе).

## 16. Приложения (Appendices)
(Как в предыдущем моем ответе).

## 17. Связанные Рабочие Процессы (Related Workflows)
*   [Процесс покупки игры и обновления библиотеки пользователя](../../../../project_workflows/game_purchase_flow.md)
*   [Синхронизация игровых сохранений с облаком](../../../../project_workflows/cloud_save_sync_flow.md) <!-- [Ссылка на cloud_save_sync_flow.md - документ в разработке] -->
*   [Процесс разблокировки достижений](../../../../project_workflows/achievement_unlocking_flow.md) <!-- [Ссылка на achievement_unlocking_flow.md - документ в разработке] -->

---
*Этот документ является основной спецификацией для Library Service и должен поддерживаться в актуальном состоянии.*
