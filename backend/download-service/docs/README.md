<!-- backend\download-service\docs\README.md -->
# Спецификация Микросервиса: Download Service (Сервис Загрузок)

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение:** Download Service является критически важным компонентом платформы "Российский Аналог Steam", отвечающим за эффективную, надежную и безопасную доставку цифрового контента (игр, DLC, обновлений, клиентского приложения платформы) на устройства пользователей.
*   **Роль в общей архитектуре платформы:** Сервис управляет всем процессом загрузки, начиная от авторизации запроса на скачивание до фактической передачи файлов через сети доставки контента (CDN). Он обеспечивает возможность приостановки/возобновления загрузок, проверку целостности файлов, применение дельта-обновлений и информирование пользователя о прогрессе.
*   **Основные бизнес-задачи:**
    *   Предоставление пользователям быстрого и надежного способа скачивания приобретенных продуктов и их обновлений.
    *   Оптимизация использования сетевых ресурсов платформы и пропускной способности CDN.
    *   Обеспечение целостности и безопасности загружаемого контента.
    *   Улучшение пользовательского опыта за счет управления очередями загрузок, приоритетов и предоставления обратной связи о процессе.
    *   Управление доставкой обновлений для основного клиентского приложения платформы.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Управление Загрузками:**
    *   Инициация загрузки по запросу от клиентского приложения (после проверки прав через Library Service).
    *   Постановка загрузок в очередь (пользовательскую и, возможно, глобальную с приоритетами).
    *   Поддержка приостановки (pause) и возобновления (resume) загрузок.
    *   Отмена загрузок.
    *   Поддержка параллельной загрузки нескольких файлов или частей одного файла (чанкинга) для ускорения.
*   **Управление Обновлениями:**
    *   Автоматическая и ручная проверка наличия обновлений для установленных продуктов.
    *   Расчет и применение дельта-обновлений (патчей) для минимизации объема скачиваемых данных.
    *   Обработка сценариев, когда дельта-обновление невозможно и требуется полная загрузка новой версии.
    *   (Опционально) Возможность отката к предыдущей стабильной версии продукта [Детали и необходимость данной функции будут уточнены. См. также GAME_UPDATE_MECHANISM.md].
*   **Доставка Клиентского Приложения Платформы:**
    *   Управление загрузкой и обновлением основного клиентского приложения платформы.
    *   Поддержка различных каналов обновлений (например, `stable`, `beta`).
*   **Обеспечение Целостности и Безопасности:**
    *   Предоставление файловых манифестов (список файлов, их размеры, хеш-суммы).
    *   Проверка хеш-сумм (например, SHA256, MD5) загруженных файлов и их частей.
    *   Автоматическое исправление поврежденных или отсутствующих файлов путем их повторной загрузки.
*   **Взаимодействие с CDN:**
    *   Генерация безопасных, временно ограниченных ссылок на загрузку файлов с CDN (например, с использованием токенизации URL или подписанных URL).
    *   Стратегии выбора оптимального CDN-сервера/региона на основе геолокации пользователя, нагрузки на CDN или других критериев.
    *   Мониторинг доступности и производительности CDN.
    *   Приоритет будет отдаваться российским CDN-провайдерам или CDN с точками присутствия в России для обеспечения быстрой и надежной доставки контента российским пользователям.
*   **Информирование о Прогрессе:**
    *   Предоставление клиентскому приложению информации о статусе и прогрессе загрузки в реальном времени (через WebSocket).
    *   Расчет текущей скорости загрузки и предполагаемого времени до завершения.
*   **Управление Настройками Загрузки:**
    *   Предоставление пользователям возможности настраивать параметры загрузки через клиентское приложение (например, ограничение максимальной скорости загрузки, настройка расписания для загрузок).

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+, согласно `../../../../project_technology_stack.md`).
*   **API:**
    *   REST API: Echo (`github.com/labstack/echo/v4`) (для взаимодействия с клиентским приложением, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
    *   gRPC: `google.golang.org/grpc` (для межсервисного взаимодействия, например, с Library Service, Catalog Service, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
    *   WebSocket: `github.com/gorilla/websocket` или аналогичная библиотека Go (для real-time обновлений прогресса загрузки на клиенте).
*   **База данных:** PostgreSQL (версия 15+) для хранения метаданных о сессиях загрузок, файловых манифестах, истории загрузок и обновлений, конфигурациях CDN. Драйвер: GORM (`gorm.io/gorm`) с `gorm.io/driver/postgres` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Кэширование/Очереди/Состояние:** Redis (версия 7.0+) для хранения состояния активных загрузок, пользовательских очередей, кэширования временных ссылок CDN, счетчиков для rate limiting. Клиент: `go-redis/redis` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Хранилище файлов (временное/staging):** S3-совместимое объектное хранилище (например, MinIO) может использоваться для временного хранения файлов при сборке дельта-патчей или если файлы сначала загружаются на сервер платформы перед распространением через CDN. (согласно `../../../../project_technology_stack.md`).
*   **Брокер сообщений:** Apache Kafka (клиент `github.com/confluentinc/confluent-kafka-go`, согласно `../../../../PACKAGE_STANDARDIZATION.md`) для асинхронной обработки задач и публикации событий.
*   **Инфраструктура:** Docker, Kubernetes.
*   **Мониторинг/Трассировка:** OpenTelemetry SDK, Prometheus client (`github.com/prometheus/client_golang`). (согласно `../../../../project_observability_standards.md`).
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Логирование:** Zap (`go.uber.org/zap`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Алгоритмы дельта-патчей:** [Алгоритм дельта-патчей: подлежит выбору/разработке].
    *   **Сети Доставки Контента (CDN):** Интеграция с CDN-провайдерами. **Приоритет отдается российским CDN (например, NGENIX, G-Core Labs, CDNvideo) или глобальным CDN с обширными точками присутствия в РФ.**
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   **CDN (Content Delivery Network):** Сеть доставки (и дистрибуции) контента, используемая для ускорения загрузки файлов пользователями. **Приоритет отдается российским CDN (например, NGENIX, G-Core Labs, CDNvideo) или глобальным CDN с обширными точками присутствия в РФ.**
*   **Дельта-обновление (Delta Update/Patch):** Обновление, содержащее только измененные части файлов по сравнению с предыдущей версией, что позволяет уменьшить объем загрузки.
*   **Манифест Файлов (File Manifest):** Список файлов, входящих в состав продукта (игры/версии), с их размерами, хеш-суммами, путями и информацией о чанках (если файлы разбиваются на части).
*   **Чанк (Chunk):** Часть файла, загружаемая отдельно при параллельной или возобновляемой загрузке.
*   **Токенизация URL (URL Tokenization):** Метод защиты ссылок CDN, при котором к URL добавляется временный токен, разрешающий доступ на ограниченное время или для конкретного пользователя/IP.
*   Для других общих терминов см. `../../../../project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Download Service спроектирован с использованием принципов Чистой Архитектуры (Clean Architecture) для обеспечения модульности, тестируемости и независимости от конкретных технологий инфраструктуры.
*   Сервис обрабатывает запросы на загрузку, управляет взаимодействием с CDN, следит за состоянием загрузок и обеспечивает целостность данных. Он также отвечает за логику применения дельта-обновлений.

### 2.2. Диаграмма Архитектуры
```mermaid
graph TD
    subgraph UserClientApp ["Клиентское Приложение"]
        ClientAppUI["UI Загрузок"]
    end

    subgraph DownloadService ["Download Service (Clean Architecture)"]
        direction TB

        subgraph PresentationLayer ["Presentation Layer (Адаптеры Транспорта)"]
            REST_API["REST API (Echo)"]
            GRPC_API["gRPC API (межсервисный)"]
            WebSocket_API["WebSocket API (прогресс)"]
        end

        subgraph ApplicationLayer ["Application Layer (Сценарии Использования)"]
            DownloadManagerSvc["Менеджер Загрузок (старт, пауза, возобновление, очередь)"]
            UpdateManagerSvc["Менеджер Обновлений (проверка, применение патчей)"]
            FileIntegritySvc["Сервис Целостности Файлов (проверка хешей)"]
            CDNOrchestratorSvc["Оркестратор CDN (генерация ссылок, выбор CDN)"]
            ProgressReporterSvc["Сервис Отчетности о Прогрессе"]
        end

        subgraph DomainLayer ["Domain Layer (Бизнес-логика и Сущности)"]
            Entities["Сущности (DownloadTask, FileChunk, ProductManifest, DeltaPatch)"]
            ValueObjects["Объекты-Значения (FileSize, Hash, DownloadURL)"]
            DomainEvents["Доменные События (DownloadInitiated, ChunkDownloaded, UpdateApplied)"]
            RepositoryIntf["Интерфейсы Репозиториев (DownloadTaskRepo, ManifestRepo)"]
            DeltaLogic["Логика Применения Дельта-Патчей"]
        end

        subgraph InfrastructureLayer ["Infrastructure Layer (Внешние Зависимости и Реализации)"]
            PostgresAdapter["Адаптер PostgreSQL (хранение задач, манифестов)"]
            RedisAdapter["Адаптер Redis (сессии, очереди, прогресс)"]
            CDNClient["Клиент CDN (взаимодействие с API CDN)"]
            S3Client["Клиент S3 (для патчей или временных файлов)"]
            KafkaProducer["Продюсер Kafka (публикация событий)"]
            KafkaConsumer["Консьюмер Kafka (получение обновлений манифестов)"]
            ServiceClients["Клиенты других сервисов (Catalog, Library, Auth)"]
            FileHasher["Утилита Хеширования/Патчинга"]
        end

        PresentationLayer --> ApplicationLayer
        ApplicationLayer --> DomainLayer
        ApplicationLayer --> InfrastructureLayer
        DomainLayer ----> RepositoryIntf
        InfrastructureLayer -- Implements --> RepositoryIntf
    end

    ClientAppUI -- REST/WebSocket --> PresentationLayer

    PostgresAdapter --> DB[(PostgreSQL)]
    RedisAdapter --> Cache[(Redis)]
    CDNClient -.-> ExternalCDN[("CDN Провайдеры")]
    S3Client -.-> ExternalS3[(S3 Хранилище)]
    KafkaProducer --> KafkaBroker[Kafka Message Bus]
    KafkaConsumer --> KafkaBroker
    ServiceClients --> InternalServices[Другие Микросервисы (Catalog, Library, Auth)]

    classDef layer_boundary fill:#f9f9f9,stroke:#333,stroke-width:2px,color:#333
    classDef component_major fill:#e6f0ff,stroke:#007bff,color:#000
    classDef component_minor fill:#d4edda,stroke:#28a745,color:#000
    classDef datastore fill:#f8d7da,stroke:#dc3545,color:#000
    classDef external_actor fill:#FEF9E7,stroke:#F1C40F,color:#000

    class PresentationLayer,ApplicationLayer,DomainLayer,InfrastructureLayer layer_boundary
    class REST_API,GRPC_API,WebSocket_API,DownloadManagerSvc,UpdateManagerSvc,FileIntegritySvc,CDNOrchestratorSvc,ProgressReporterSvc,Entities,ValueObjects,DomainEvents,RepositoryIntf,DeltaLogic component_major
    class PostgresAdapter,RedisAdapter,CDNClient,S3Client,KafkaProducer,KafkaConsumer,ServiceClients,FileHasher component_minor
    class DB,Cache,KafkaBroker datastore
    class ExternalCDN,ExternalS3,InternalServices external_actor
```

### 2.3. Слои Сервиса

#### 2.3.1. Presentation Layer (Слой Представления)
*   **Ответственность:** Обработка входящих запросов от клиентского приложения (REST API для управления загрузками, WebSocket для обновлений прогресса) и от других микросервисов (gRPC для авторизации загрузок, получения информации о файлах). Валидация DTO, вызов Application Layer.
*   **Ключевые компоненты:** HTTP хендлеры, gRPC серверы, WebSocket хендлеры.

#### 2.3.2. Application Layer (Прикладной Слой)
*   **Ответственность:** Реализация сценариев использования: инициация и управление загрузками/обновлениями, проверка наличия обновлений, применение патчей, проверка целостности. Координация между Domain Layer и Infrastructure Layer.
*   **Ключевые компоненты:** `DownloadManagerService`, `UpdateManagerService`, `CDNOrchestrationService`, `ProgressReportingService`.

#### 2.3.3. Domain Layer (Доменный Слой)
*   **Ответственность:** Бизнес-логика и правила, связанные с процессом загрузки. Сущности (`DownloadTask`, `FileChunk`, `ProductManifest`, `DeltaPatch`), объекты-значения (`FileSize`, `Hash`), доменные события. Логика расчета дельта-патчей и их применения.
*   **Ключевые компоненты:** Сущности, репозитории (интерфейсы), сервисы домена (например, `PatchingService`).

#### 2.3.4. Infrastructure Layer (Инфраструктурный Слой)
*   **Ответственность:** Реализация интерфейсов репозиториев (PostgreSQL, Redis). Взаимодействие с CDN (генерация и валидация ссылок). Взаимодействие с S3 для временного хранения или работы с патчами. Публикация и потребление событий Kafka. Клиенты для Auth, Catalog, Library Service. Утилиты для хеширования и применения патчей.
*   **Ключевые компоненты:** Реализации репозиториев, клиенты CDN/S3/Kafka, gRPC клиенты.

## 3. API Endpoints

### 3.1. REST API (для клиентского приложения)
*   **Базовый URL:** `/api/v1/download` (маршрутизируется через API Gateway).
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway или в middleware сервиса).
*   **Формат ответа об ошибке:** Согласно `../../../../project_api_standards.md`.

#### 3.1.1. Управление Загрузками
*   **`POST /tasks`**
    *   Описание: Инициировать новую загрузку продукта (игры/версии) или группы файлов. Сервис проверяет права доступа (через Library Service) и доступность файлов (через Catalog Service).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "downloadTaskRequest",
            "attributes": {
              "productId": "game-uuid-123",
              "versionId": "version-uuid-abc", // Опционально, если запрашивается последняя версия
              "targetPath": "/user/games/MySuperGame", // Предлагаемый пользователем путь
              "priority": 1 // 0 - normal, 1 - high etc.
            }
          }
        }
        ```
    *   Ответ (202 Accepted): (Объект `DownloadTask` в статусе `queued` или `preparing`)
*   **`GET /tasks`**
    *   Описание: Получение списка текущих и недавних задач на загрузку пользователя.
    *   Query параметры: `status` (queued, downloading, paused, completed, error), `page`, `limit`.
    *   Ответ (200 OK): (Массив объектов `DownloadTask`)
*   **`GET /tasks/{taskId}`**
    *   Описание: Получение детальной информации о конкретной задаче на загрузку.
    *   Ответ (200 OK): (Объект `DownloadTask` с детальным списком файлов и их статусом)
*   **`POST /tasks/{taskId}/pause`**
    *   Описание: Приостановить активную загрузку.
    *   Ответ (200 OK): (Обновленный объект `DownloadTask`)
*   **`POST /tasks/{taskId}/resume`**
    *   Описание: Возобновить приостановленную загрузку.
    *   Ответ (200 OK): (Обновленный объект `DownloadTask`)
*   **`DELETE /tasks/{taskId}`**
    *   Описание: Отменить загрузку (и удалить частично загруженные файлы).
    *   Ответ (204 No Content)
*   **`GET /tasks/queue`**
    *   Описание: Получение текущей очереди загрузок пользователя.
    *   Ответ (200 OK): `{ "data": [ /* ordered list of DownloadTask summaries */ ] }`
*   **`POST /tasks/queue/reorder`**
    *   Описание: Изменение порядка загрузок в очереди.
    *   Тело запроса: `{"data": [{"taskId": "uuid1", "order": 0}, {"taskId": "uuid2", "order": 1}]}`
    *   Ответ (200 OK)

#### 3.1.2. Управление Обновлениями (для установленных продуктов)
*   **`POST /updates/check`**
    *   Описание: Проверка наличия обновлений для списка установленных продуктов.
    *   Тело запроса:
        ```json
        {
          "data": [
            {"type": "installedProduct", "attributes": {"productId": "game-uuid-123", "currentVersionId": "version-uuid-abc"}},
            {"type": "installedProduct", "attributes": {"productId": "game-uuid-456", "currentVersionId": "version-uuid-def"}}
          ]
        }
        ```
    *   Ответ (200 OK):
        ```json
        {
          "data": [
            {"type": "updateInfo", "attributes": {"productId": "game-uuid-123", "isUpdateAvailable": true, "latestVersionId": "version-uuid-new", "updateType": "delta", "downloadSizeBytes": 536870912, "requiredDiskSpaceBytes": 1073741824}},
            {"type": "updateInfo", "attributes": {"productId": "game-uuid-456", "isUpdateAvailable": false}}
          ]
        }
        ```
*   **`POST /updates/initiate`**
    *   Описание: Инициировать загрузку обновления для продукта (добавляет задачу в общую очередь загрузок).
    *   Тело запроса: `{"data": {"type": "updateInitiateRequest", "attributes": {"productId": "game-uuid-123", "targetVersionId": "version-uuid-new"}}}`
    *   Ответ (202 Accepted): (Объект `DownloadTask` для обновления)

#### 3.1.3. Проверка Целостности Файлов
*   **`POST /integrity/verify`**
    *   Описание: Инициировать проверку целостности файлов для установленного продукта.
    *   Тело запроса: `{"data": {"type": "integrityVerificationRequest", "attributes": {"productId": "game-uuid-123", "installedVersionId": "version-uuid-abc"}}}`
    *   Ответ (202 Accepted): (Объект `VerificationTask` со статусом `pending` или `in_progress`)
*   **`GET /integrity/verify/{taskId}`**
    *   Описание: Получить статус и результат задачи проверки целостности.
    *   Ответ (200 OK): (Объект `VerificationTask` с деталями)

### 3.2. gRPC API (для межсервисного взаимодействия)
*   **Пакет:** `download.v1`
*   **Сервис:** `DownloadInternalService`
    *   **`rpc AuthorizeDownload(AuthorizeDownloadRequest) returns (AuthorizeDownloadResponse)`**
        *   Описание: Используется Library Service для проверки прав пользователя и получения метаданных для начала загрузки.
        *   `message AuthorizeDownloadRequest { string user_id = 1; string product_id = 2; string version_id = 3; }`
        *   `message AuthorizeDownloadResponse { bool authorized = 1; string download_session_id = 2; FileManifest manifest = 3; repeated CDNSource sources = 4; string error_message = 5; }`
    *   **`rpc NotifyNewVersionAvailable(NotifyNewVersionRequest) returns (google.protobuf.Empty)`**
        *   Описание: Используется Catalog Service для уведомления Download Service о появлении новой версии продукта или манифеста.
        *   `message NotifyNewVersionRequest { string product_id = 1; string version_id = 2; FileManifest manifest = 3; }`

### 3.3. WebSocket API
*   **Эндпоинт:** `/ws/download/progress` (требует аутентификации, например, через токен в query-параметре при установке соединения).
*   **Сообщения от сервера к клиенту:**
    *   **`downloadTaskStatusUpdate`**: Обновление общего статуса задачи на загрузку.
        ```json
        {
          "type": "downloadTaskStatusUpdate",
          "payload": {
            "taskId": "session-uuid-xyz",
            "status": "downloading", // queued, downloading, paused, completed, error, verifying_files
            "totalProgressPercentage": 25.5,
            "downloadSpeedBps": 5242880, // Байт в секунду
            "estimatedTimeLeftSeconds": 1800,
            "downloadedBytes": 2684354560,
            "totalSizeBytes": 10737418240,
            "errorDetails": null // if status is 'error'
          }
        }
        ```
    *   **`downloadItemProgressUpdate`**: Обновление прогресса по конкретному файлу/чанку в задаче.
        ```json
        {
          "type": "downloadItemProgressUpdate",
          "payload": {
            "taskId": "session-uuid-xyz",
            "itemId": "item-uuid-123", // ID файла/чанка
            "relativePath": "data/level1.pak",
            "status": "downloading",
            "downloadedBytes": 52428800,
            "totalFileSizeBytes": 104857600
          }
        }
        ```
    *   **`downloadError`**: Сообщение об ошибке загрузки.
        ```json
        {
          "type": "downloadError",
          "payload": {
            "taskId": "session-uuid-xyz",
            "itemId": "item-uuid-123", // Опционально, если ошибка по конкретному файлу
            "errorCode": "CDN_UNREACHABLE",
            "errorMessage": "Не удалось подключиться к серверу CDN."
          }
        }
        ```
*   **Сообщения от клиента к серверу:** (Обычно не требуются, управление через REST API)

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`DownloadTask` (Задача на Загрузку)**: Представляет собой задачу пользователя на загрузку продукта (игры, обновления). Включает несколько `DownloadItem`.
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID). **Обязательность: Да.**
    *   `product_id` (UUID). **Обязательность: Да.**
    *   `version_id` (UUID). **Обязательность: Да.**
    *   `target_path` (TEXT). **Обязательность: Да.**
    *   `status` (ENUM: `queued`, `preparing`, `downloading`, `paused`, `verifying`, `completed`, `error`, `cancelled`). **Обязательность: Да (DEFAULT 'queued').**
    *   `priority` (INTEGER). **Обязательность: Да (DEFAULT 0).**
    *   `total_size_bytes` (BIGINT). **Обязательность: Да (DEFAULT 0).**
    *   `downloaded_bytes` (BIGINT). **Обязательность: Да (DEFAULT 0).**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `error_message` (TEXT). **Обязательность: Нет.**
*   **`DownloadItem` (Элемент Загрузки)**: Конкретный файл или чанк файла в рамках `DownloadTask`.
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `download_task_id` (UUID, FK). **Обязательность: Да.**
    *   `file_manifest_id` (UUID, FK на `ProductFileManifest`). **Обязательность: Да (или ссылка на конкретную запись в JSONB манифеста).**
    *   `relative_path` (TEXT). **Обязательность: Да.**
    *   `status` (ENUM: `pending`, `downloading`, `paused`, `completed`, `error`, `verifying`). **Обязательность: Да (DEFAULT 'pending').**
    *   `downloaded_bytes` (BIGINT). **Обязательность: Да (DEFAULT 0).**
    *   `total_size_bytes` (BIGINT). **Обязательность: Да.**
    *   `retry_count` (INTEGER). **Обязательность: Да (DEFAULT 0).**
    *   `cdn_url_current` (TEXT, nullable). **Обязательность: Нет.**
    *   `expected_hash_sha256` (VARCHAR). **Обязательность: Нет (если нет в манифесте).**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
*   **`ProductFileManifest` (Манифест Файлов Продукта)**: Описывает структуру файлов для конкретной версии продукта. Получается из Catalog Service. Может кэшироваться в Download Service.
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `product_id` (UUID). **Обязательность: Да.**
    *   `version_id` (UUID). **Обязательность: Да.**
    *   `file_entries` (JSONB: `[{"path": "bin/game.exe", "size": 12345, "hash_sha256": "...", "chunks": [{"offset":0, "size":6000, "hash_sha256":"..."}]}]`). **Обязательность: Да.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
*   **`CDNConfig` (Конфигурация CDN)**: Настройки для работы с различными CDN провайдерами.
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `provider_name` (VARCHAR, UK). **Обязательность: Да.**
    *   `api_endpoint` (TEXT). **Обязательность: Нет.**
    *   `auth_type` (VARCHAR). **Обязательность: Да (DEFAULT 'none').**
    *   `auth_credentials_encrypted` (TEXT). **Обязательность: Нет (если auth_type != 'none').**
    *   `url_generation_template` (TEXT). **Обязательность: Да.**
    *   `priority` (INTEGER). **Обязательность: Да (DEFAULT 0).**
    *   `is_active` (BOOLEAN). **Обязательность: Да (DEFAULT TRUE).**
*   **`UserDownloadHistory` (История Загрузок Пользователя)**: Запись о завершенных или отмененных загрузках.
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `user_id` (UUID). **Обязательность: Да.**
    *   `product_id` (UUID). **Обязательность: Да.**
    *   `version_id` (UUID). **Обязательность: Да.**
    *   `status` (ENUM: `completed`, `cancelled`, `failed_permanently`). **Обязательность: Да.**
    *   `completed_at` (TIMESTAMPTZ). **Обязательность: Да (если статус completed).**
    *   `total_download_time_seconds` (INTEGER). **Обязательность: Нет.**
    *   `total_bytes_downloaded` (BIGINT). **Обязательность: Да.**
*   **`DeltaPatchInfo` (Информация о Дельта-Патче)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `product_id` (UUID). **Обязательность: Да.**
    *   `from_version_id` (UUID). **Обязательность: Да.**
    *   `to_version_id` (UUID). **Обязательность: Да.**
    *   `patch_s3_path` (TEXT). **Обязательность: Да.**
    *   `patch_size_bytes` (BIGINT). **Обязательность: Да.**
    *   `patch_hash_sha256` (TEXT). **Обязательность: Да.**
    *   `instructions_s3_path` (TEXT, опционально, для сложных патчей). **Обязательность: Нет.**

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL
**ERD Диаграмма:**
```mermaid
erDiagram
    USERS {
        UUID id PK
        note "From Auth/Account Service"
    }
    PRODUCTS {
        UUID id PK
        note "From Catalog Service"
    }
    PRODUCT_VERSIONS {
        UUID id PK
        UUID product_id FK
        note "From Catalog Service"
    }
    DOWNLOAD_TASKS {
        UUID id PK
        UUID user_id FK
        UUID product_id FK
        UUID version_id FK
        VARCHAR target_path
        VARCHAR status
        INTEGER priority
        BIGINT total_size_bytes
        BIGINT downloaded_bytes
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    DOWNLOAD_ITEMS {
        UUID id PK
        UUID download_task_id FK
        UUID product_file_manifest_entry_id "Refers to an entry in ProductFileManifest (conceptual)"
        TEXT relative_path
        VARCHAR status
        BIGINT downloaded_bytes
        BIGINT total_size_bytes
        INTEGER retry_count
    }
    PRODUCT_FILE_MANIFESTS {
        UUID id PK
        UUID product_id FK
        UUID version_id FK
        JSONB file_entries "List of files, sizes, hashes, chunks"
        TIMESTAMPTZ created_at
    }
    CDN_CONFIGS {
        UUID id PK
        VARCHAR provider_name UK
        TEXT api_endpoint
        TEXT url_generation_template
        INTEGER priority
        BOOLEAN is_active
    }
    USER_DOWNLOAD_HISTORY {
        UUID id PK
        UUID user_id FK
        UUID product_id FK
        UUID version_id FK
        VARCHAR status
        TIMESTAMPTZ completed_at
        BIGINT total_bytes_downloaded
    }
    DELTA_PATCH_INFO {
        UUID id PK
        UUID product_id FK
        UUID from_version_id FK
        UUID to_version_id FK
        TEXT patch_s3_path
        BIGINT patch_size_bytes
        VARCHAR patch_hash_sha256
    }

    USERS ||--o{ DOWNLOAD_TASKS : "initiates"
    PRODUCTS ||--o{ DOWNLOAD_TASKS : "target_product"
    PRODUCT_VERSIONS ||--o{ DOWNLOAD_TASKS : "target_version"
    DOWNLOAD_TASKS ||--|{ DOWNLOAD_ITEMS : "contains"
    PRODUCT_FILE_MANIFESTS ||--|| PRODUCT_VERSIONS : "describes"
    PRODUCTS ||--|| PRODUCT_FILE_MANIFESTS : "for_product"
    USERS ||--o{ USER_DOWNLOAD_HISTORY : "has"
    PRODUCTS ||--o{ USER_DOWNLOAD_HISTORY : "of_product"
    PRODUCT_VERSIONS ||--o{ USER_DOWNLOAD_HISTORY : "of_version"
    PRODUCTS ||--o{ DELTA_PATCH_INFO : "for_product"
    PRODUCT_VERSIONS ||--o{ DELTA_PATCH_INFO : "from_version"
    PRODUCT_VERSIONS ||--o{ DELTA_PATCH_INFO : "to_version"

    DOWNLOAD_ITEMS ..> PRODUCT_FILE_MANIFESTS : "references_entry_in"
```

**DDL (PostgreSQL - ключевые таблицы):**
```sql
CREATE TABLE download_tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL, -- FK to Users table (conceptual)
    product_id UUID NOT NULL, -- FK to Products table in Catalog Service (conceptual)
    version_id UUID NOT NULL, -- FK to ProductVersions table in Catalog Service (conceptual)
    target_path TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'queued' CHECK (status IN ('queued', 'preparing', 'downloading', 'paused', 'verifying', 'completed', 'error', 'cancelled')),
    priority INTEGER NOT NULL DEFAULT 0,
    total_size_bytes BIGINT NOT NULL DEFAULT 0,
    downloaded_bytes BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    error_message TEXT
);
CREATE INDEX idx_download_tasks_user_status ON download_tasks(user_id, status);

CREATE TABLE download_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    download_task_id UUID NOT NULL REFERENCES download_tasks(id) ON DELETE CASCADE,
    -- file_manifest_id UUID NOT NULL, -- conceptually links to an entry in ProductFileManifest
    relative_path TEXT NOT NULL, -- Путь файла относительно корня продукта
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'downloading', 'paused', 'completed', 'error', 'verifying')),
    downloaded_bytes BIGINT NOT NULL DEFAULT 0,
    total_size_bytes BIGINT NOT NULL,
    expected_hash_sha256 VARCHAR(64),
    current_cdn_url TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_download_items_task_id ON download_items(download_task_id);

CREATE TABLE product_file_manifests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL,
    version_id UUID NOT NULL,
    file_entries JSONB NOT NULL, -- [{"path": "...", "size": ..., "hash": "...", "chunks": [...]}]
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (product_id, version_id)
);
COMMENT ON TABLE product_file_manifests IS 'Хранит манифесты файлов для версий продуктов, полученные из Catalog Service.';

CREATE TABLE cdn_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider_name VARCHAR(100) NOT NULL UNIQUE,
    api_endpoint TEXT,
    url_generation_template TEXT NOT NULL, -- e.g., "https://{host}/{path}?token={token}"
    auth_type VARCHAR(50) DEFAULT 'none', -- none, token_param, signed_url
    secret_key_id_for_signing VARCHAR(255), -- Если используется подпись URL
    priority INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT
);
COMMENT ON TABLE cdn_configs IS 'Конфигурации для различных CDN провайдеров.';

-- (UserDownloadHistory и DeltaPatchInfo DDLs как в существующем документе, если они там есть, или создать аналогично)
```

#### 4.2.2. Redis
*   **Активные сессии/задачи загрузки:** `active_download_task:<taskId>` (HASH) - хранит часто обновляемые поля: `status`, `downloaded_bytes`, `current_speed_bps`, `active_item_id`, `active_item_progress_bytes`. TTL на случай зависания.
*   **Очередь загрузок пользователя:** `user_download_queue:<userId>` (LIST или SORTED SET) - ID задач в порядке их выполнения.
*   **Прогресс по файлам/чанкам:** `download_item_progress:<itemId>` (HASH) - `downloaded_bytes`, `status`.
*   **Кэш ссылок CDN:** `cdn_url:<file_path_or_chunk_id>` (STRING) - кэшированная подписанная ссылка CDN с TTL.
*   **Счетчики для Rate Limiting:** (Если Download Service реализует свой rate limiting для запросов на генерацию ссылок) `rate_limit:user:<userId>:cdn_url_requests` (COUNTER с TTL).

## 5. Потоковая Обработка Событий (Event Streaming)
*   **Формат событий:** CloudEvents v1.0 JSON (согласно `../../../../project_api_standards.md`).
*   **Основной топик для публикуемых событий:** `com.platform.download.events.v1`.

*   **`com.platform.download.task.status.changed.v1`** (Ранее `download.session.status.changed.v1`)
    *   Описание: Статус задачи на загрузку изменился (например, начата, приостановлена, завершена, ошибка).
    *   `data` Payload:
        ```json
        {
          "taskId": "session-uuid-xyz",
          "userId": "user-uuid-123",
          "productId": "game-uuid-abc",
          "versionId": "version-uuid-def",
          "newStatus": "completed", // "downloading", "paused", "error", "verifying"
          "oldStatus": "downloading",
          "totalProgressPercentage": 100.0,
          "downloadedBytes": 10737418240,
          "totalSizeBytes": 10737418240,
          "changeTimestamp": "2024-07-12T10:00:00Z",
          "errorDetails": null
        }
        ```
*   **`com.platform.download.item.status.changed.v1`**
    *   Описание: Статус отдельного файла/чанка в задаче на загрузку изменился.
    *   `data` Payload:
        ```json
        {
          "taskId": "session-uuid-xyz",
          "itemId": "item-uuid-123",
          "relativePath": "data/level1.pak",
          "newStatus": "completed", // "downloading", "paused", "error", "verifying"
          "downloadedBytes": 104857600,
          "totalFileSizeBytes": 104857600,
          "changeTimestamp": "2024-07-12T09:55:00Z",
          "errorDetails": null
        }
        ```
*   **`com.platform.download.update.available.v1`**
    *   Описание: Для установленного продукта пользователя доступно обновление.
    *   `data` Payload:
        ```json
        {
          "userId": "user-uuid-123",
          "productId": "game-uuid-abc",
          "latestVersionId": "version-uuid-new",
          "updateType": "delta", // "delta" or "full"
          "downloadSizeBytes": 536870912,
          "releaseNotesUrl": "https://example.com/updates/game-abc/v2.1.0"
        }
        ```
    *   Потребители: Notification Service (для уведомления пользователя), Клиентское приложение (для отображения).
*   **`com.platform.download.manifest.updated.v1`**
    *   Описание: Манифест файлов для продукта/версии был обновлен (например, после публикации новой версии в Catalog Service).
    *   `data` Payload:
        ```json
        {
          "productId": "game-uuid-abc",
          "versionId": "version-uuid-new",
          "manifestUrl": "s3://bucket/manifests/product_abc_v_new.json", // или сам манифест
          "updateTimestamp": "2024-07-12T08:00:00Z"
        }
        ```
    *   Потребители: Download Service (для обновления своего кэша манифестов).

### 5.2. Потребляемые События (Consumed Events)
*   **`com.platform.catalog.product.version.published.v1`** (от Catalog Service)
    *   Описание: Новая версия продукта была опубликована в каталоге.
    *   Ожидаемый `data` Payload: `{"productId": "game-uuid-abc", "versionId": "version-uuid-new", "manifestS3Path": "path/to/manifest.json", ...}`
    *   Логика обработки: Загрузить/обновить `ProductFileManifest`. Если есть пользователи с установленной предыдущей версией, инициировать проверку обновлений и возможно опубликовать `com.platform.download.update.available.v1`.
*   **`com.platform.library.download.request.authorized.v1`** (от Library Service, альтернатива gRPC вызову)
    *   Описание: Пользователю разрешена загрузка продукта/версии.
    *   Ожидаемый `data` Payload: `{"userId": "user-uuid-123", "productId": "game-uuid-abc", "versionId": "version-uuid-def"}`
    *   Логика обработки: Инициировать создание `DownloadTask`.

## 6. Интеграции (Integrations)
(Содержимое существующего раздела актуально, с уточнением взаимодействия с CDN и S3).
*   **CDN Провайдеры:** Ключевая интеграция. Download Service запрашивает у CDN (или генерирует для CDN) безопасные ссылки для скачивания файлов.
*   **S3-совместимое хранилище:** Для хранения оригиналов билдов (если загружаются через Developer Service -> Download Service -> CDN), дельта-патчей, временных файлов.

## 7. Конфигурация (Configuration)
Общие стандарты конфигурационных файлов (формат YAML, структура, управление переменными окружения и секретами) определены в `../../../../project_api_standards.md` (раздел 7) и `../../../../DOCUMENTATION_GUIDELINES.md` (раздел 6).

### 7.1. Переменные Окружения (Примеры)
*   `DOWNLOAD_HTTP_PORT`: Порт для REST API (например, `8083`)
*   `DOWNLOAD_GRPC_PORT`: Порт для gRPC API (например, `9093`)
*   `DOWNLOAD_WEBSOCKET_PORT`: Порт для WebSocket API (например, `8084`)
*   `POSTGRES_DSN_DOWNLOAD`: Строка подключения к PostgreSQL.
*   `REDIS_ADDR_DOWNLOAD`: Адрес Redis.
*   `KAFKA_BROKERS_DOWNLOAD`: Список брокеров Kafka.
*   `LOG_LEVEL_DOWNLOAD`: Уровень логирования.
*   `CDN_PROVIDER_DEFAULT_API_KEY`: API ключ для CDN по умолчанию.
*   `CDN_PROVIDER_DEFAULT_SECRET`: Секрет для CDN по умолчанию.
*   `S3_STAGING_BUCKET_PATCHES`: S3 бакет для временного хранения патчей.
*   `DELTA_PATCH_MIN_FILE_SIZE_MB`: Минимальный размер файла для применения дельта-патча.
*   `JWT_PUBLIC_KEY_PATH`: Путь к публичному ключу для валидации JWT.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT_DOWNLOAD`: Endpoint для Jaeger.

### 7.2. Файлы Конфигурации (`configs/download_service_config.yaml`)
*   Структура:
    ```yaml
    http_server:
      port: ${DOWNLOAD_HTTP_PORT:"8083"}
      timeout_seconds: 30
    grpc_server:
      port: ${DOWNLOAD_GRPC_PORT:"9093"}
      timeout_seconds: 30
    websocket_server:
      port: ${DOWNLOAD_WEBSOCKET_PORT:"8084"}
      # ...другие настройки WebSocket
    postgres:
      dsn: ${POSTGRES_DSN_DOWNLOAD}
      pool_max_conns: 10
    redis:
      address: ${REDIS_ADDR_DOWNLOAD}
      password: ${REDIS_PASSWORD_DOWNLOAD:""}
      db: ${REDIS_DB_DOWNLOAD:3} # Пример
    kafka:
      brokers: ${KAFKA_BROKERS_DOWNLOAD}
      producer_topics:
        download_events: "com.platform.download.events.v1"
      consumer_topics:
        catalog_events: "com.platform.catalog.product.version.published.v1"
        # ...другие потребляемые топики
      consumer_group: "download-service-group"
    logging:
      level: ${LOG_LEVEL_DOWNLOAD:"info"}
    cdn:
      default_provider: "my_cdn_provider_key" # Key to specific provider config
      providers:
        my_cdn_provider_key:
          api_key: ${CDN_PROVIDER_DEFAULT_API_KEY}
          secret: ${CDN_PROVIDER_DEFAULT_SECRET} # Store actual secret in k8s secrets
          url_template: "https://cdn.example.com/{path}?token={token}"
        # ...другие провайдеры
    s3_staging_patches:
      bucket: ${S3_STAGING_BUCKET_PATCHES}
      endpoint: ${S3_ENDPOINT_STAGING} # Пример, если отличается от основного S3
      access_key: ${S3_ACCESS_KEY_STAGING}
      secret_key: ${S3_SECRET_KEY_STAGING}
      region: "ru-central1"
    delta_updates:
      min_file_size_mb: ${DELTA_PATCH_MIN_FILE_SIZE_MB:10}
      patch_algorithm: "bsdiff" # [bsdiff, xdelta, custom_rsync_like - подлежит выбору]
    security:
      jwt_public_key_path: ${JWT_PUBLIC_KEY_PATH}
    otel:
      exporter_jaeger_endpoint: ${OTEL_EXPORTER_JAEGER_ENDPOINT_DOWNLOAD}
      service_name: "download-service"
    ```

## 8. Обработка Ошибок (Error Handling)
(Содержимое существующего раздела актуально).
*   **`CDN_LINK_GENERATION_FAILED`**: Ошибка при генерации ссылки на CDN.
*   **`FILE_HASH_MISMATCH`**: Хеш-сумма загруженного файла не совпадает с ожидаемой.
*   **`DELTA_PATCH_APPLICATION_FAILED`**: Ошибка применения дельта-патча.
*   **`DOWNLOAD_QUEUE_FULL`**: Очередь загрузок пользователя заполнена.

## 9. Безопасность (Security)
(Содержимое существующего раздела актуально).
*   **Защита ссылок CDN:** Использование короткоживущих, подписанных URL или токенизированных ссылок для предотвращения неавторизованного доступа и распространения контента.
*   **Проверка целостности файлов:** Обязательная проверка хеш-сумм после загрузки для гарантии отсутствия повреждений или модификаций во время передачи.
*   **Авторизация загрузок:** Тесная интеграция с Library Service и Auth Service для проверки, что пользователь имеет право на загрузку контента.
*   **Защита от злоупотреблений:** Rate limiting на запросы генерации ссылок CDN.

## 10. Развертывание (Deployment)
(Содержимое существующего раздела актуально).

## 11. Мониторинг и Логирование (Logging and Monitoring)
См. `../../../../project_observability_standards.md` для общих стандартов.
*   **Ключевые метрики для мониторинга включают:**
    *   Количество активных загрузок.
    *   Общая пропускная способность (байт/сек).
    *   Средняя скорость загрузки на пользователя.
    *   Количество ошибок при генерации ссылок CDN.
    *   Количество ошибок при скачивании с CDN (по кодам ошибок).
    *   Процент успешно завершенных загрузок.
    *   Процент успешно примененных дельта-патчей.
    *   Размер очередей задач на загрузку/обработку.
    *   Задержки при взаимодействии с внешними сервисами (Catalog, Library, CDN API).
    *   Скорость загрузки (средняя, по регионам, по CDN), количество активных загрузок, количество ошибок CDN, процент успешных/неуспешных загрузок, время применения патчей.
*   **Логирование:** Детальное логирование всех этапов процесса загрузки, включая запросы на авторизацию, получение манифестов, генерацию ссылок, применение патчей, ошибки.

## 12. Нефункциональные Требования (NFRs)
(Содержимое существующего раздела актуально).
*   **Пропускная способность:** Способность обслуживать X Гбит/с суммарного трафика загрузок (зависит от CDN).
*   **Задержка генерации ссылок CDN:** P99 < 200 мс.
*   **Надежность применения патчей:** > 99.9% успешных применений.

## 13. Приложения (Appendices)
(Содержимое существующего раздела актуально).

## 14. Пользовательские Сценарии (User Flows)

В этом разделе описаны ключевые пользовательские сценарии, связанные с Download Service.

### 14.1. Пользователь Начинает Загрузку Новой Игры
*   **Описание:** Пользователь, имеющий права на игру, инициирует ее загрузку через клиентское приложение.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Клиентское Приложение
        participant APIGW as API Gateway
        participant LibrarySvc as Library Service
        participant AuthSvc as Auth Service
        participant DownloadSvc as Download Service
        participant CatalogSvc as Catalog Service
        participant CDN as CDN
        participant Kafka as Kafka Message Bus

        User->>ClientApp: Нажимает "Скачать игру X"
        ClientApp->>APIGW: POST /api/v1/library/me/downloads (productId, versionId)
        APIGW->>LibrarySvc: Forward request (с User JWT)
        LibrarySvc->>AuthSvc: (gRPC) ValidateToken(JWT) -> user_id
        AuthSvc-->>LibrarySvc: Token valid, user_id
        LibrarySvc->>LibrarySvc: Проверка прав пользователя на продукт (entitlement)
        alt Права есть
            LibrarySvc->>DownloadSvc: (gRPC) RequestDownload(user_id, productId, versionId)
            DownloadSvc->>CatalogSvc: (gRPC) GetFileManifest(productId, versionId)
            CatalogSvc-->>DownloadSvc: FileManifest (список файлов, хеши, размеры, чанки)
            DownloadSvc->>DownloadSvc: Создание DownloadTask, DownloadItems (status: 'queued')
            DownloadSvc->>DownloadSvc: Генерация безопасных ссылок на CDN для первых чанков/файлов
            DownloadSvc-->>LibrarySvc: DownloadQueuedResponse (taskId, queue_position)
            LibrarySvc-->>APIGW: HTTP 202 Accepted
            APIGW-->>ClientApp: HTTP 202 Accepted (taskId)
            ClientApp-->>User: Загрузка добавлена в очередь

            ClientApp->>DownloadSvc: (WebSocket) Connect /ws/download/progress?taskId={taskId}
            DownloadSvc-->>ClientApp: (WebSocket) Connection established; TaskStatusUpdate (status: 'downloading')

            loop Для каждого файла/чанка в задаче
                DownloadSvc->>ClientApp: (WebSocket) ItemProgressUpdate (CDN_URL_for_chunk, item_id, path)
                ClientApp->>CDN: GET <CDN_URL_for_chunk>
                CDN-->>ClientApp: Данные чанка
                ClientApp->>ClientApp: Сохранение чанка
                ClientApp->>DownloadSvc: (REST или WebSocket) Сообщение о завершении чанка (опционально, или DownloadSvc сам отслеживает)
                DownloadSvc->>DownloadSvc: Проверка хеша чанка (если есть)
            end
            DownloadSvc->>DownloadSvc: Проверка целостности всех файлов после завершения
            DownloadSvc->>Kafka: Publish `com.platform.download.session.status.changed.v1` (taskId, status: 'completed')
            DownloadSvc-->>ClientApp: (WebSocket) DownloadTaskStatusUpdate (status: 'completed')
            ClientApp-->>User: Загрузка завершена
        else Права отсутствуют
            LibrarySvc-->>APIGW: HTTP 403 Forbidden
            APIGW-->>ClientApp: HTTP 403 Forbidden
            ClientApp-->>User: Ошибка: нет прав на загрузку
        end
    ```

### 14.2. Пользователь Приостанавливает и Возобновляет Загрузку
*   **Описание:** Пользователь приостанавливает активную загрузку, а затем возобновляет ее.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor User
        participant ClientApp as Клиентское Приложение
        participant APIGW as API Gateway
        participant DownloadSvc as Download Service
        participant CDN as CDN

        User->>ClientApp: Нажимает "Пауза" для загрузки {taskId}
        ClientApp->>APIGW: POST /api/v1/download/tasks/{taskId}/pause
        APIGW->>DownloadSvc: Forward request
        DownloadSvc->>DownloadSvc: Обновляет статус DownloadTask на 'paused'. Отменяет текущие HTTP запросы к CDN для этой задачи.
        DownloadSvc-->>APIGW: HTTP 200 OK (обновленный DownloadTask)
        APIGW-->>ClientApp: HTTP 200 OK
        ClientApp->>DownloadSvc: (WebSocket) Получает DownloadTaskStatusUpdate (status: 'paused')
        ClientApp-->>User: Загрузка приостановлена

        User->>ClientApp: Нажимает "Возобновить" для загрузки {taskId}
        ClientApp->>APIGW: POST /api/v1/download/tasks/{taskId}/resume
        APIGW->>DownloadSvc: Forward request
        DownloadSvc->>DownloadSvc: Обновляет статус DownloadTask на 'downloading'.
        DownloadSvc->>DownloadSvc: Определяет недокачанные файлы/чанки. Генерирует новые ссылки CDN.
        DownloadSvc-->>APIGW: HTTP 200 OK (обновленный DownloadTask)
        APIGW-->>ClientApp: HTTP 200 OK
        ClientApp->>DownloadSvc: (WebSocket) Получает DownloadTaskStatusUpdate (status: 'downloading') и ItemProgressUpdate для возобновления
        ClientApp-->>User: Загрузка возобновлена
    ```

### 14.3. Применение Дельта-Обновления для Установленной Игры
*   **Описание:** Клиентское приложение проверяет наличие обновлений, обнаруживает дельта-патч, скачивает его и применяет.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        participant ClientApp as Клиентское Приложение
        participant APIGW as API Gateway
        participant DownloadSvc as Download Service
        participant CatalogSvc as Catalog Service
        participant CDN as CDN
        participant LocalFileSystem as Локальная ФС Клиента

        ClientApp->>APIGW: POST /api/v1/download/updates/check (productId, currentVersionId)
        APIGW->>DownloadSvc: Forward request
        DownloadSvc->>CatalogSvc: (gRPC) GetLatestVersionInfo(productId)
        CatalogSvc-->>DownloadSvc: LatestVersionInfo
        DownloadSvc->>DownloadSvc: Сравнение версий. Определение наличия дельта-патча.
        alt Дельта-патч доступен
            DownloadSvc->>CatalogSvc: (gRPC) GetDeltaPatchInfo(productId, currentVersionId, latestVersionId)
            CatalogSvc-->>DownloadSvc: DeltaPatchInfo (S3 путь к патчу, хеш, размер)
            DownloadSvc-->>APIGW: HTTP 200 OK (updateAvailable=true, type='delta', patchInfo)
        else Дельта-патч недоступен, доступно полное обновление
            DownloadSvc-->>APIGW: HTTP 200 OK (updateAvailable=true, type='full', fullDownloadInfo)
        end
        APIGW-->>ClientApp: Ответ о наличии обновления

        ClientApp->>APIGW: POST /api/v1/download/updates/initiate (productId, targetVersionId, type='delta')
        APIGW->>DownloadSvc: Forward request
        DownloadSvc->>DownloadSvc: Создание DownloadTask для патча. Генерация CDN ссылок для патча.
        DownloadSvc-->>APIGW: HTTP 202 Accepted (DownloadTask для патча)
        APIGW-->>ClientApp: HTTP 202 Accepted

        ClientApp->>CDN: Скачивание файла(ов) патча
        CDN-->>ClientApp: Данные патча
        ClientApp->>LocalFileSystem: Применение патча к локальным файлам игры (используя логику из DeltaLogic)
        alt Патч применен успешно
            ClientApp->>LocalFileSystem: Проверка целостности обновленных файлов
            alt Целостность подтверждена
                ClientApp->>APIGW: POST /api/v1/download/updates/applied (productId, newVersionId, status='success')
                APIGW->>DownloadSvc: Forward request
                DownloadSvc->>DownloadSvc: Обновление истории обновлений пользователя
                DownloadSvc->>LibrarySvc: (через Kafka/gRPC) Уведомление об успешном обновлении
            else Ошибка целостности
                ClientApp->>APIGW: POST /api/v1/download/updates/applied (productId, newVersionId, status='verification_failed')
                Note over ClientApp: Запрос полной переустановки или повторной проверки.
            end
        else Ошибка применения патча
            ClientApp->>APIGW: POST /api/v1/download/updates/applied (productId, currentVersionId, status='patch_failed')
            Note over ClientApp: Запрос полной переустановки.
        end
    ```

### 14.4. Управление Очередью Загрузок Пользователя
*   **Описание:** Download Service управляет несколькими задачами на загрузку для одного пользователя, учитывая их приоритеты и глобальные лимиты.
*   **Диаграмма:**
    ```mermaid
    graph TD
        U[Пользователь] --> Q1{Запрос Загрузки 1 (Игра A, Приоритет Normal)}
        U --> Q2{Запрос Загрузки 2 (Патч B, Приоритет High)}
        U --> Q3{Запрос Загрузки 3 (Игра C, Приоритет Normal)}

        subgraph DownloadService
            PackageManager[Менеджер Пакетов/Очереди]
            ActiveDownloader1[Активный Загрузчик 1]
            ActiveDownloader2[Активный Загрузчик 2 (если max_concurrent > 1)]
        end

        Q1 --> PackageManager
        Q2 --> PackageManager
        Q3 --> PackageManager

        PackageManager -- Учитывает приоритет и лимит --> ActiveDownloader1
        PackageManager -- Учитывает приоритет и лимит --> ActiveDownloader2

        ActiveDownloader1 --> CDN1[CDN]
        ActiveDownloader2 --> CDN2[CDN]

        UserClient[Клиент Пользователя] <-->|WebSocket| PackageManager
        PackageManager --> UserClient: Обновления статуса для всех задач в очереди

        note right of PackageManager
         - Обработка очереди (например, Redis Sorted Set по приоритету и времени добавления).
         - Запуск N параллельных загрузок согласно настройкам пользователя/системы.
         - Приостановка низкоприоритетных задач при добавлении высокоприоритетной, если лимит исчерпан.
        end
    ```

### 14.5. Обработка Недоступности CDN или Переключение на Альтернативный CDN
*   **Описание:** В процессе загрузки файла CDN становится недоступен. Download Service пытается переключиться на другой CDN или повторить запрос через некоторое время.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        participant ClientApp as Клиентское Приложение
        participant DownloadSvc as Download Service
        participant PrimaryCDN as Основной CDN
        participant SecondaryCDN as Резервный CDN (если есть)
        participant CatalogSvc as Catalog Service (для списка зеркал/CDN)

        ClientApp->>PrimaryCDN: GET /path/to/file_chunk_1 (по ссылке от DownloadSvc)
        alt PrimaryCDN недоступен или ошибка
            PrimaryCDN-->>ClientApp: Ошибка (Timeout, 5xx)
            ClientApp->>DownloadSvc: (REST/WebSocket) Сообщение об ошибке загрузки чанка (itemId, errorCode)
            DownloadSvc->>DownloadSvc: Логирование ошибки. Инкремент счетчика ошибок для PrimaryCDN / файла.
            DownloadSvc->>CatalogSvc: (gRPC, если информация о CDN там) Запрос альтернативных CDN для файла/продукта
            CatalogSvc-->>DownloadSvc: Список альтернативных CDN или стратегия отката
            alt Есть резервный CDN
                DownloadSvc->>DownloadSvc: Генерация новой ссылки для SecondaryCDN
                DownloadSvc->>ClientApp: (WebSocket) Новая ссылка для загрузки чанка с SecondaryCDN
                ClientApp->>SecondaryCDN: GET /path/to/file_chunk_1 (новая ссылка)
                SecondaryCDN-->>ClientApp: Данные чанка
            else Нет резервного CDN или он тоже недоступен
                DownloadSvc->>DownloadSvc: Планирование повторной попытки для PrimaryCDN через N секунд (backoff)
                DownloadSvc->>ClientApp: (WebSocket) Уведомление о временной проблеме и планируемой повторной попытке
            end
        else PrimaryCDN отвечает успешно
            PrimaryCDN-->>ClientApp: Данные чанка
        end
    ```

## 15. Резервное Копирование и Восстановление (Backup and Recovery)

### 15.1. PostgreSQL (Метаданные сессий загрузок, истории обновлений, результатов верификации)
*   **Процедура резервного копирования:**
    *   Ежедневный логический бэкап (`pg_dump`).
    *   Настроена непрерывная архивация WAL-сегментов (PITR), базовый бэкап еженедельно.
    *   **Хранение:** Бэкапы в S3, шифрование, версионирование, другой регион. Срок хранения: полные - 30 дней, WAL - 14 дней.
*   **Процедура восстановления:** Тестируется ежеквартально.
*   **RTO:** < 2 часов.
*   **RPO:** < 15 минут.
*   (Общие принципы см. `../../../../project_database_structure.md`).

### 15.2. Redis (Очереди загрузок, состояние активных сессий, кэш токенов CDN)
*   **Стратегия персистентности:**
    *   **AOF (Append Only File):** Включен с fsync `everysec` для очередей и состояний активных сессий, если их потеря критична для пользовательского опыта (например, чтобы не прерывать все активные загрузки при перезапуске Redis).
    *   **RDB Snapshots:** Регулярное создание снапшотов (например, каждые 1-6 часов).
*   **Резервное копирование (снапшотов):** RDB-снапшоты могут копироваться в S3 ежедневно. Срок хранения - 7 дней.
*   **Восстановление:** Из последнего RDB-снапшота и/или AOF. Большинство данных (кроме, возможно, долгоживущих очередей) могут быть перестроены или сессии будут переинициализированы клиентами.
*   **RTO:** < 30 минут.
*   **RPO:** < 1 минуты (для данных с AOF `everysec`). Для кэша токенов CDN RPO менее критичен.

### 15.3. S3-совместимое хранилище (если используется для временного staging)
*   **Стратегия:** Данные во временном хранилище обычно имеют короткий срок жизни. Основные файлы игр управляются Developer/Catalog Service и их S3 бакетами. Download Service не отвечает за бэкап этих основных файлов.
*   **Резервное копирование:** Обычно не требуется для временных файлов, специфичных для Download Service. Если есть критичные staging-файлы, создаваемые самим Download Service, можно настроить версионирование или репликацию в S3.
*   **RTO/RPO:** Неприменимо для временных данных; для критичных staging-данных зависит от настроек S3.

### 15.4. Общая стратегия
*   Восстановление PostgreSQL является приоритетным.
*   Redis восстанавливается для минимизации прерываний текущих операций.
*   Процедуры документированы и тестируются.

## 16. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API, Protobuf определения для gRPC API, и форматы сообщений WebSocket поддерживаются в соответствующих репозиториях исходного кода сервиса и/или в централизованном репозитории `platform-protos`.
*   DDL схемы базы данных управляются через систему миграций (например, `golang-migrate/migrate`) и хранятся в репозитории исходного кода сервиса. Актуальные версии доступны во внутренней документации команды разработки и в GitOps репозитории.

## 17. Связанные Рабочие Процессы (Related Workflows)
*   [Процесс покупки игры и обновления библиотеки пользователя](../../../../project_workflows/game_purchase_flow.md)
*   [Процесс обновления клиентского приложения](../../../../project_workflows/client_update_flow.md) Подробное описание этого рабочего процесса будет добавлено в [client_update_flow.md](../../../../project_workflows/client_update_flow.md) (документ в разработке).

## Game Update Mechanism Integration
Download Service играет ключевую роль в процессе обновления игр. Подробное описание общего механизма обновления игр и их компонентов на платформе представлено в документе [GAME_UPDATE_MECHANISM.md](../../../../GAME_UPDATE_MECHANISM.md). Download Service отвечает за доставку билдов, обновлений и патчей, как описано в указанном документе, взаимодействуя с Catalog Service и Developer Service для получения информации о версиях и манифестах.

---
*Этот документ является основной спецификацией для Download Service и должен поддерживаться в актуальном состоянии.*
