# Спецификация Микросервиса: Download Service (Сервис Загрузок)

**Версия:** 1.1 (адаптировано из предыдущей версии)
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Download Service платформы "Российский Аналог Steam".
*   **Роль в общей архитектуре платформы:** Download Service является критически важным компонентом, отвечающим за эффективную и надежную загрузку игрового контента (игр, обновлений, DLC) и самого клиентского приложения платформы на устройства пользователей. Он управляет процессом доставки файлов, взаимодействуя с CDN, и обеспечивает целостность загружаемых данных.
*   **Основные бизнес-задачи:**
    *   Предоставление пользователям возможности скачивать приобретенные игры и их обновления.
    *   Обеспечение высокой скорости и надежности загрузок, включая возможность приостановки и возобновления.
    *   Оптимизация использования сетевых ресурсов и мощностей CDN.
    *   Гарантия целостности и безопасности загружаемых файлов.
    *   Предоставление информации о статусе и прогрессе загрузок клиентскому приложению.
    *   Управление обновлениями клиентского приложения платформы.
*   Разработка сервиса должна вестись в соответствии с `CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Управление загрузками игр:** Инициация, приостановка, возобновление и отмена загрузок. Управление очередью загрузок. Поддержка параллельных загрузок нескольких файлов/частей файлов.
*   **Управление обновлениями игр:** Автоматическая и ручная проверка наличия обновлений. Поддержка дельта-обновлений (патчей) для минимизации объема скачиваемых данных. Возможность отката к предыдущей версии (если поддерживается).
*   **Управление обновлениями клиентского приложения:** Доставка обновлений для основного клиентского приложения платформы, включая поддержку различных каналов (stable, beta).
*   **Проверка целостности файлов:** Верификация загруженных файлов по хеш-суммам (например, SHA256, MD5). Автоматическое исправление поврежденных или отсутствующих файлов путем их повторной загрузки.
*   **Взаимодействие с CDN:** Генерация безопасных ссылок на загрузку с CDN. Выбор оптимального CDN-сервера на основе геолокации пользователя или нагрузки.
*   **Мониторинг и статистика загрузок:** Отслеживание прогресса, текущей скорости загрузки, предполагаемого времени завершения. Сбор статистики для анализа производительности CDN и выявления проблем.
*   **Управление настройками загрузки:** Предоставление пользователям возможности настраивать параметры загрузки (например, ограничение скорости, расписание загрузок).

### 1.3. Основные Технологии
*   **Язык программирования:** Go (предпочтительно).
*   **API:**
    *   REST API (для клиентского приложения, например, через Echo/Gin).
    *   gRPC (для межсервисного взаимодействия).
    *   WebSocket (для real-time обновлений прогресса загрузки на клиенте).
*   **База данных:** PostgreSQL (для хранения метаданных о загрузках, файлах, версиях, статусах).
*   **Кэширование/Очереди:** Redis (для хранения временных данных сессий загрузки, состояния активных загрузок, кэширования токенов CDN, управления небольшими очередями задач).
*   **Хранилище файлов (временное/staging):** S3-совместимое объектное хранилище (например, MinIO) может использоваться для временного хранения файлов перед их передачей в CDN или для сборки дельта-патчей. Основные игровые файлы предполагаются доступными через Catalog/Developer service и физически находятся на CDN или в S3, управляемом ими.
*   **Брокер сообщений:** Apache Kafka (для асинхронного обмена событиями с другими сервисами).
*   **Инфраструктура:** Docker, Kubernetes.
*   **Мониторинг/Трассировка:** OpenTelemetry, Prometheus, Grafana, Jaeger/Tempo.
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `PACKAGE_STANDARDIZATION.md`.
*   *Примечание: Приведенные в последующих разделах примеры API и конфигураций основаны на предположении использования Go, PostgreSQL, Redis и Kafka.*

### 1.4. Термины и Определения (Glossary)
*   **CDN (Content Delivery Network):** Сеть доставки (и дистрибуции) контента, используемая для ускорения загрузки файлов пользователями.
*   **Дельта-обновление (Delta Update/Patch):** Обновление, содержащее только измененные части файлов по сравнению с предыдущей версией, что позволяет уменьшить объем загрузки.
*   **Манифест Файлов (File Manifest):** Список файлов, входящих в состав продукта (игры/версии), с их размерами, хеш-суммами и путями.
*   **Chunk:** Часть файла, загружаемая отдельно при параллельной или возобновляемой загрузке.
*   Для других общих терминов см. `project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Download Service будет реализован с использованием многослойной архитектуры (Clean Architecture) для обеспечения четкого разделения ответственностей, улучшения тестируемости и гибкости.
*   Ключевые компоненты включают управление сессиями загрузок, взаимодействие с CDN для получения ссылок, обработку запросов на скачивание и обновление, проверку целостности файлов, а также управление очередями и приоритетами загрузок.

**Диаграмма Архитектуры:**
```mermaid
graph TD
    subgraph User Client App
        ClientApp[Клиентское приложение платформы]
    end

    subgraph Download Service
        direction TB

        subgraph PresentationLayer [Presentation Layer (Адаптеры Транспорта)]
            REST_API[REST API (Echo/Gin) - для клиента]
            GRPC_API[gRPC API - для других сервисов]
            WebSocket_API[WebSocket API - для real-time прогресса]
        end

        subgraph ApplicationLayer [Application Layer (Сценарии Использования)]
            DownloadUseCaseSvc[Управление Загрузками (Start, Pause, Resume)]
            UpdateUseCaseSvc[Управление Обновлениями (Check, Apply)]
            IntegrityCheckSvc[Проверка Целостности Файлов]
            SettingsUseCaseSvc[Управление Настройками Загрузки]
            CDNTokenSvc[Сервис Токенов/Ссылок CDN]
        end

        subgraph DomainLayer [Domain Layer (Бизнес-логика и Сущности)]
            Entities[Сущности (DownloadSession, DownloadItem, FileManifest, UpdateInfo)]
            ValueObjects[Объекты-Значения (FileSize, Hash, DownloadSpeed)]
            DomainEvents[Доменные События (DownloadStarted, FileChunkCompleted)]
            RepositoryIntf[Интерфейсы Репозиториев (Download, FileMetadata)]
        end

        subgraph InfrastructureLayer [Infrastructure Layer (Внешние Зависимости)]
            PostgresAdapter[Адаптер PostgreSQL (Метаданные загрузок)]
            RedisAdapter[Адаптер Redis (Кэш сессий, Очереди)]
            S3Adapter[Адаптер S3 (Временное хранилище, если нужно)]
            CDNClient[Клиент CDN (Получение/валидация ссылок)]
            KafkaProducer[Продюсер Kafka (События)]
            ServiceClients[Клиенты др. сервисов (Catalog, Library, Auth)]
            Hasher[Утилита Хеширования]
        end

        REST_API & WebSocket_API --> DownloadUseCaseSvc
        REST_API --> UpdateUseCaseSvc
        REST_API --> IntegrityCheckSvc
        REST_API --> SettingsUseCaseSvc
        GRPC_API --> DownloadUseCaseSvc # например, для инициирования загрузки из LibraryService
        GRPC_API --> UpdateUseCaseSvc

        DownloadUseCaseSvc & UpdateUseCaseSvc --> CDNTokenSvc
        ApplicationLayer --> DomainLayer
        ApplicationLayer --> InfrastructureLayer
        DomainLayer ----> RepositoryIntf
        InfrastructureLayer -- Implements --> RepositoryIntf
    end

    ClientApp -- REST/WebSocket --> PresentationLayer

    PostgresAdapter --> DB[(PostgreSQL)]
    RedisAdapter --> Cache[(Redis)]
    S3Adapter --> S3[(S3 Хранилище)]
    CDNClient -.-> ExternalCDN[CDN Провайдеры]
    KafkaProducer --> Kafka[Kafka Broker]
    ServiceClients --> OtherServices[Другие Микросервисы]

    classDef layer_boundary fill:#f9f9f9,stroke:#333,stroke-width:2px,color:#333
    classDef component_major fill:#e6f0ff,stroke:#007bff,color:#000
    classDef component_minor fill:#d4edda,stroke:#28a745,color:#000
    classDef datastore fill:#f8d7da,stroke:#dc3545,color:#000

    class PresentationLayer,ApplicationLayer,DomainLayer,InfrastructureLayer layer_boundary
    class REST_API,GRPC_API,WebSocket_API,DownloadUseCaseSvc,UpdateUseCaseSvc,IntegrityCheckSvc,SettingsUseCaseSvc,CDNTokenSvc,Entities,Aggregates,DomainEvents,RepositoryIntf component_major
    class PostgresAdapter,RedisAdapter,S3Adapter,CDNClient,KafkaProducer,ServiceClients,Hasher component_minor
    class DB,Cache,S3,Kafka,OtherServices,ExternalCDN datastore
```

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления / Транспортный слой)
*   **Ответственность:** Обработка входящих запросов от клиентского приложения (REST API, WebSocket для прогресса) и других микросервисов (gRPC). Валидация DTO, вызов соответствующей логики в Application Layer.
*   **Ключевые компоненты/модули:** HTTP хендлеры (Echo/Gin), gRPC серверы, WebSocket хендлеры, DTO для запросов/ответов.

#### 2.2.2. Application Layer (Прикладной Слой / Сервисный слой)
*   **Ответственность:** Реализация бизнес-логики управления загрузками, обновлениями, проверкой целостности. Координирует взаимодействие между Domain Layer и Infrastructure Layer. Управляет потоком операций.
*   **Ключевые компоненты/модули:** Сервисы сценариев использования (`DownloadApplicationService`, `UpdateApplicationService`, `FileIntegrityService`, `DownloadSettingsService`), сервисы для работы с CDN (например, генерация защищенных ссылок, выбор оптимального CDN).

#### 2.2.3. Domain Layer (Доменный Слой)
*   **Ответственность:** Содержит бизнес-сущности, агрегаты, доменные события и бизнес-правила, связанные с процессом загрузки и обновления.
*   **Ключевые компоненты/модули:** Сущности (`DownloadSession`, `DownloadItem`, `UpdatePackage`, `FileMetadata`, `VerificationAttempt`), объекты-значения (`FileSize`, `HashChecksum`, `DownloadSpeed`, `URL`), доменные события (`DownloadInitiated`, `DownloadProgressUpdated`, `FileVerified`), интерфейсы репозиториев.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   **Ответственность:** Реализация интерфейсов репозиториев для работы с PostgreSQL (метаданные загрузок) и Redis (временные данные, очереди, кэш). Взаимодействие с CDN провайдерами. Взаимодействие с S3 (если используется для временного хранения). Публикация событий в Kafka. Клиенты для взаимодействия с другими микросервисами (Catalog, Library, Auth, Account). Утилиты для хеширования файлов и обработки дельта-патчей.
*   **Ключевые компоненты/модули:** Реализации репозиториев, CDN клиенты, S3 клиенты, Kafka продюсеры, gRPC/HTTP клиенты к другим сервисам, модуль для работы с файловой системой (для временных файлов при обработке патчей).

## 3. API Endpoints

### 3.1. REST API (для клиентского приложения)
*   **Базовый URL:** `/api/v1/download-client` (маршрутизируется через API Gateway).
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway или в middleware сервиса).
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

#### 3.1.1. Управление Загрузками
*   **`POST /downloads`**
    *   Описание: Инициировать новую загрузку продукта (игры/версии). Сервис проверяет права доступа (через Library Service) и доступность файлов (через Catalog Service).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "downloadRequest",
            "attributes": {
              "product_id": "game-uuid-123",
              "version_id": "version-uuid-abc", // Опционально, если не указано - последняя доступная
              "installation_path": "/opt/games/MySuperGame" // Предложенный пользователем путь
            }
          }
        }
        ```
    *   Пример ответа (Успех 202 Accepted - загрузка поставлена в очередь):
        ```json
        {
          "data": {
            "type": "downloadSession",
            "id": "download-session-uuid-xyz",
            "attributes": {
              "product_id": "game-uuid-123",
              "version_id": "version-uuid-abc",
              "status": "queued", // queued, downloading, paused, completed, error
              "total_size_bytes": 10737418240, // 10 GB
              "downloaded_bytes": 0,
              "progress_percentage": 0.0,
              "estimated_time_left_seconds": null
            }
          }
        }
        ```
    *   Требуемые права доступа: `user` (владелец игры).
*   **`GET /downloads`**
    *   Описание: Получение списка текущих и недавних загрузок пользователя.
    *   Query параметры: `status` (queued, downloading, paused, completed, error), `page`, `limit`.
    *   Пример ответа (Успех 200 OK): (Массив объектов `downloadSession`)
    *   Требуемые права доступа: `user`.

#### 3.1.2. Управление Обновлениями
*   **`POST /updates/check`**
    *   Описание: Проверка наличия обновлений для установленных продуктов.
    *   Тело запроса:
        ```json
        {
          "data": [
            { "type": "installedProduct", "id": "game-uuid-123", "attributes": { "current_version_id": "version-uuid-abc" } },
            { "type": "installedProduct", "id": "game-uuid-456", "attributes": { "current_version_id": "version-uuid-def" } }
          ]
        }
        ```
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "updateAvailability", "id": "game-uuid-123",
              "attributes": { "is_update_available": true, "latest_version_id": "version-uuid-new", "update_size_bytes": 536870912 }
            }
            // ...
          ]
        }
        ```
    *   Требуемые права доступа: `user`.

#### 3.1.3. Настройки Загрузки
*   **`PATCH /settings`**
    *   Описание: Обновление пользовательских настроек загрузки.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "downloadSettings",
            "attributes": {
              "max_concurrent_downloads": 3,
              "bandwidth_limit_kbps": 5000 // 0 для безлимитного
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): (Возвращает обновленные настройки)
    *   Требуемые права доступа: `user`.

### 3.2. gRPC API (для межсервисного взаимодействия)
*   Пакет: `download.v1`.
*   Определение Protobuf: `download/v1/download_service.proto`.

#### 3.2.1. Сервис: `DownloadInternalService`
*   **`rpc GetFileMetadataAndUrl(GetFileMetadataAndUrlRequest) returns (GetFileMetadataAndUrlResponse)`**
    *   Описание: Получение метаданных файла (размер, хеш) и защищенной ссылки на загрузку с CDN. Используется клиентским приложением перед началом скачивания каждого файла или чанка.
    *   `message GetFileMetadataAndUrlRequest { string product_id = 1; string version_id = 2; string file_path_in_manifest = 3; string client_ip_address = 4; /* для выбора CDN */ }`
    *   `message FileLocation { string url = 1; string cdn_provider = 2; int64 expires_at_unix = 3; }`
    *   `message GetFileMetadataAndUrlResponse { string file_id_internal = 1; int64 file_size_bytes = 2; string file_hash_sha256 = 3; repeated FileLocation locations = 4; }`
    *   Требуемые права доступа: Аутентифицированный пользователь, имеющий права на продукт (проверяется через Library Service).
*   **`rpc ReportDownloadStatus(ReportDownloadStatusRequest) returns (ReportDownloadStatusResponse)`**
    *   Описание: Клиент сообщает о статусе загрузки файла/чанка (успех, ошибка, прогресс).
    *   `message ReportDownloadStatusRequest { string download_session_id = 1; string file_id_internal = 2; enum Status { STARTED = 0; DOWNLOADING = 1; COMPLETED = 2; FAILED = 3; } status = 3; int64 downloaded_bytes = 4; /* для DOWNLOADING */ string error_message = 5; /* для FAILED */ }`
    *   `message ReportDownloadStatusResponse { bool acknowledged = 1; }`

### 3.3. WebSocket API
*   **Эндпоинт:** `/api/v1/ws/download-progress`
*   **Аутентификация:** JWT токен передается при установлении соединения (например, в query параметре).
*   **Сообщения от сервера к клиенту:**
    *   Тип: `downloadProgressUpdate`
    *   Пример:
        ```json
        {
          "event_type": "downloadProgressUpdate",
          "payload": {
            "download_session_id": "download-session-uuid-xyz",
            "product_id": "game-uuid-123",
            "status": "downloading",
            "progress_percentage": 45.5,
            "downloaded_bytes": 4865392640,
            "total_size_bytes": 10737418240,
            "current_speed_bps": 15728640, // 15 Mbps
            "estimated_time_left_seconds": 3720
          }
        }
        ```
    *   Тип: `downloadStatusChanged` (например, `paused`, `completed`, `error`)
    *   Пример:
        ```json
        {
          "event_type": "downloadStatusChanged",
          "payload": {
            "download_session_id": "download-session-uuid-xyz",
            "product_id": "game-uuid-123",
            "new_status": "completed",
            "error_details": null // или объект ошибки
          }
        }
        ```

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности

*   **`DownloadSession` (Сессия Загрузки)**
    *   `id` (UUID): Уникальный идентификатор сессии загрузки. Обязательность: Required.
    *   `user_id` (UUID): ID пользователя. Обязательность: Required.
    *   `product_id` (UUID): ID продукта (игры/DLC). Обязательность: Required.
    *   `version_id` (UUID): ID версии продукта. Обязательность: Required.
    *   `status` (ENUM: `queued`, `initializing`, `downloading`, `paused`, `verifying`, `completed`, `error`, `cancelled`): Текущий статус загрузки. Обязательность: Required.
    *   `total_size_bytes` (BIGINT): Общий размер всех файлов для загрузки. Обязательность: Required.
    *   `downloaded_bytes` (BIGINT): Общий объем уже загруженных данных. Обязательность: Required.
    *   `progress_percentage` (FLOAT, 0-100): Процент выполнения загрузки. Обязательность: Required.
    *   `current_speed_bps` (BIGINT): Текущая скорость загрузки в битах/сек. Обязательность: Optional.
    *   `estimated_time_left_seconds` (INTEGER): Примерное оставшееся время. Обязательность: Optional.
    *   `installation_path` (TEXT): Путь установки на диске пользователя. Обязательность: Required.
    *   `priority` (INTEGER): Приоритет загрузки в очереди. Обязательность: Required.
    *   `created_at` (TIMESTAMPTZ), `updated_at` (TIMESTAMPTZ).

*   **`DownloadItem` (Элемент Загрузки - файл в рамках сессии)**
    *   `id` (UUID): Уникальный идентификатор элемента. Обязательность: Required.
    *   `download_session_id` (UUID, FK to DownloadSession): ID сессии. Обязательность: Required.
    *   `file_metadata_id` (UUID, FK to FileMetadata): Ссылка на метаданные файла из каталога. Обязательность: Required.
    *   `relative_path` (TEXT): Относительный путь файла в структуре продукта. Пример: `bin/game.exe`. Обязательность: Required.
    *   `status` (ENUM: `pending`, `downloading`, `completed`, `verifying`, `error`): Статус загрузки конкретного файла. Обязательность: Required.
    *   `downloaded_bytes` (BIGINT): Объем загруженных данных для этого файла. Обязательность: Required.
    *   `cdn_url_current` (TEXT): Текущая ссылка на CDN, с которой идет загрузка. Обязательность: Optional.
    *   `retry_count` (INTEGER): Количество попыток загрузки. Обязательность: Required.

*   **`FileMetadata` (Метаданные Файла Продукта)**
    *   `id` (UUID): Уникальный идентификатор метаданных файла. Обязательность: Required.
    *   `product_id` (UUID): ID продукта. Обязательность: Required.
    *   `version_id` (UUID): ID версии продукта. Обязательность: Required.
    *   `relative_path` (TEXT): Относительный путь файла. Валидация: not null. Обязательность: Required.
    *   `size_bytes` (BIGINT): Размер файла в байтах. Валидация: not null, >0. Обязательность: Required.
    *   `hash_sha256` (VARCHAR(64)): SHA256 хеш файла. Валидация: not null, hex format. Обязательность: Required.
    *   `is_delta_patch` (BOOLEAN): Является ли файл дельта-патчем. Обязательность: Required.
    *   `base_version_id_for_patch` (UUID): Если `is_delta_patch`=true, ID базовой версии. Обязательность: Optional.
    *   `cdn_references` (JSONB): Информация о доступности на различных CDN или регионах CDN. Пример: `[{"provider": "akamai", "region": "eu-central", "priority": 1}, ...]`. Обязательность: Required.

*   **`Update` (Обновление)**: Представляет доступное или примененное обновление для продукта.
    *   `id` (UUID): Уникальный идентификатор записи об обновлении.
    *   `product_id` (UUID): ID продукта.
    *   `from_version_id` (UUID): С какой версии обновление.
    *   `to_version_id` (UUID): До какой версии обновление.
    *   `update_type` (ENUM: `full`, `delta`): Тип обновления.
    *   `size_bytes` (BIGINT): Размер обновления.
    *   `release_notes_url` (TEXT): Ссылка на заметки к выпуску.
    *   `status` (ENUM: `available`, `downloading`, `applied`, `failed`): Статус для конкретного пользователя. (Может храниться в связке с user_id).

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL

**ERD Диаграмма (ключевые таблицы):**
```mermaid
erDiagram
    DOWNLOAD_SESSIONS {
        UUID id PK
        UUID user_id "FK (User)"
        UUID product_id "FK (Product)"
        UUID version_id "FK (ProductVersion)"
        VARCHAR status
        BIGINT total_size_bytes
        BIGINT downloaded_bytes
        FLOAT progress_percentage
        TEXT installation_path
        INTEGER priority
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    DOWNLOAD_ITEMS {
        UUID id PK
        UUID download_session_id FK
        UUID file_metadata_id FK
        TEXT relative_path
        VARCHAR status
        BIGINT downloaded_bytes
        INTEGER retry_count
    }
    FILE_METADATA {
        UUID id PK
        UUID product_id "FK (Product)"
        UUID version_id "FK (ProductVersion)"
        TEXT relative_path
        BIGINT size_bytes
        VARCHAR hash_sha256
        BOOLEAN is_delta_patch
        UUID base_version_id_for_patch "nullable"
        JSONB cdn_references
        TIMESTAMPTZ created_at
    }
    UPDATES_HISTORY { # Для отслеживания установленных обновлений пользователем
        UUID id PK
        UUID user_id "FK (User)"
        UUID product_id "FK (Product)"
        UUID applied_version_id "FK (ProductVersion)"
        TIMESTAMPTZ applied_at
    }
    VERIFICATION_RESULTS {
        UUID id PK
        UUID user_id "FK (User)"
        UUID product_id "FK (Product)"
        UUID version_id "FK (ProductVersion)"
        VARCHAR status -- pending, in_progress, completed_ok, completed_errors
        TIMESTAMPTZ started_at
        TIMESTAMPTZ completed_at
    }

    DOWNLOAD_SESSIONS ||--|{ DOWNLOAD_ITEMS : "contains"
    DOWNLOAD_ITEMS }o--|| FILE_METADATA : "references"
    FILE_METADATA }o--|| PRODUCTS_VERSIONS : "belongs_to_version" # Предполагаемая таблица из Catalog
    USERS ||--o{ DOWNLOAD_SESSIONS : "owns" # Предполагаемая таблица пользователей
    USERS ||--o{ UPDATES_HISTORY : "has_applied"
    USERS ||--o{ VERIFICATION_RESULTS : "initiated_by"
    PRODUCTS_VERSIONS ||--o{ VERIFICATION_RESULTS : "verifies"
```

**DDL (PostgreSQL - примеры ключевых таблиц):**
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE download_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL, -- FK to users table in Account/Auth Service
    product_id UUID NOT NULL, -- FK to products table in Catalog Service
    version_id UUID NOT NULL, -- FK to product_versions table in Catalog Service
    status VARCHAR(50) NOT NULL DEFAULT 'queued' CHECK (status IN ('queued', 'initializing', 'downloading', 'paused', 'verifying', 'completed', 'error', 'cancelled')),
    total_size_bytes BIGINT NOT NULL DEFAULT 0,
    downloaded_bytes BIGINT NOT NULL DEFAULT 0,
    progress_percentage REAL NOT NULL DEFAULT 0.0,
    installation_path TEXT NOT NULL,
    priority INTEGER NOT NULL DEFAULT 5,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_download_sessions_user_status ON download_sessions(user_id, status);

CREATE TABLE file_metadata (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL,
    version_id UUID NOT NULL,
    relative_path TEXT NOT NULL, -- e.g., "bin/game.exe", "data/level1.pak"
    size_bytes BIGINT NOT NULL,
    hash_sha256 VARCHAR(64) NOT NULL,
    is_delta_patch BOOLEAN NOT NULL DEFAULT FALSE,
    base_version_id_for_patch UUID, -- if is_delta_patch is true
    cdn_references JSONB NOT NULL, -- [{"provider": "akamai", "region": "eu", "url_template": "https://cdn1.example.com/{path}"}, ...]
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (product_id, version_id, relative_path)
);
CREATE INDEX idx_file_metadata_product_version ON file_metadata(product_id, version_id);

CREATE TABLE download_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    download_session_id UUID NOT NULL REFERENCES download_sessions(id) ON DELETE CASCADE,
    file_metadata_id UUID NOT NULL REFERENCES file_metadata(id) ON DELETE RESTRICT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'downloading', 'completed', 'verifying', 'error')),
    downloaded_bytes BIGINT NOT NULL DEFAULT 0,
    current_cdn_url TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    error_message TEXT,
    UNIQUE (download_session_id, file_metadata_id)
);
CREATE INDEX idx_download_items_session_status ON download_items(download_session_id, status);

-- TODO: Добавить DDL для updates_history, verification_results.
```

#### 4.2.2. Redis
*   **Очереди загрузок (Download Queues):**
    *   Ключ: `download_queue:<user_id>` (LIST или ZSET для приоритетов). Значение: `download_session_id`.
*   **Текущие статусы/прогресс активных загрузок (для быстрого доступа и WebSocket):**
    *   Ключ: `download_status:<download_session_id>` (HASH). Поля: `product_id`, `status`, `progress_percentage`, `downloaded_bytes`, `current_speed_bps`, `eta_seconds`. TTL короткий, обновляется часто.
*   **Кэш токенов/временных ссылок CDN:**
    *   Ключ: `cdn_token:<file_id_or_path_hash>` (STRING). Значение: валидная ссылка на CDN. TTL по времени жизни ссылки.
*   **Счетчики для Rate Limiting (если применяются к API Download Service):**
    *   Ключ: `rl:download_api:<user_id_or_ip>:<action>` (COUNTER).
*   **Блокировки для конкурентного доступа (например, при обновлении статуса сессии):**
    *   Ключ: `lock:download_session:<download_session_id>` (STRING с NX/EX).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Apache Kafka.
*   **Формат событий:** CloudEvents v1.0 (JSON encoding), согласно `project_api_standards.md`.
*   **Основной топик:** `download.events.v1`.

*   **`download.session.status.changed.v1`**
    *   Описание: Статус сессии загрузки изменился (например, started, completed, failed, paused, resumed).
    *   Пример Payload (`data` секция CloudEvent):
        ```json
        {
          "download_session_id": "download-session-uuid-xyz",
          "user_id": "user-uuid-123",
          "product_id": "game-uuid-abc",
          "version_id": "version-uuid-def",
          "new_status": "completed", // "started", "paused", "failed", etc.
          "previous_status": "downloading",
          "timestamp": "2024-03-15T10:30:00Z",
          "details": { // Опционально, в зависимости от статуса
            "error_code": null, // если new_status = "failed"
            "error_message": null
          }
        }
        ```
*   **`download.update.available.v1`** (Может публиковаться при проверке обновлений)
    *   Описание: Для продукта пользователя доступно обновление.
    *   Пример Payload:
        ```json
        {
          "user_id": "user-uuid-123",
          "product_id": "game-uuid-abc",
          "current_version_id": "version-uuid-old",
          "new_version_id": "version-uuid-new",
          "update_size_bytes": 536870912,
          "release_notes_summary": "Исправлены критические ошибки, добавлена новая карта.",
          "timestamp": "2024-03-15T11:00:00Z"
        }
        ```
*   **`download.file.verification.completed.v1`**
    *   Описание: Завершена проверка целостности файлов для продукта.
    *   Пример Payload:
        ```json
        {
          "download_session_id": "download-session-uuid-xyz", // Если в рамках сессии загрузки
          "user_id": "user-uuid-123",
          "product_id": "game-uuid-abc",
          "version_id": "version-uuid-def",
          "verification_status": "success", // "failed_repair_needed", "failed_cannot_repair"
          "corrupted_files_count": 0,
          "repaired_files_count": 0,
          "timestamp": "2024-03-15T12:00:00Z"
        }
        ```
*   TODO: Детализировать другие события, если необходимо (например, `download.item.status.changed`).

### 5.2. Потребляемые События (Consumed Events)

*   **`catalog.product.version.published.v1`** (от Catalog Service)
    *   Описание: Новая версия продукта (игры/DLC/ПО) была опубликована в каталоге и стала доступна для загрузки.
    *   Ожидаемый Payload:
        ```json
        {
          "product_id": "game-uuid-abc",
          "version_id": "version-uuid-new",
          "version_name": "1.1.0",
          "release_date": "2024-03-15T00:00:00Z",
          "file_manifest_url": "s3://catalog-data/manifests/game-uuid-abc/version-uuid-new_manifest.json", // URL к манифесту файлов версии
          "total_size_bytes": 12884901888 // 12 GB
        }
        ```
    *   Логика обработки: Download Service должен получить манифест файлов (список файлов, их размеры, хеши, пути в CDN или S3) для новой версии из Catalog Service (или по указанному URL). Сохранить эти метаданные в свою таблицу `file_metadata`. Если пользователи имеют право на это обновление (например, автоматическое обновление включено), инициировать или запланировать проверку обновлений для них.
*   **`library.user.game.added.v1`** (от Library Service)
    *   Описание: Игра была добавлена в библиотеку пользователя (например, после покупки или активации ключа).
    *   Ожидаемый Payload:
        ```json
        {
          "user_id": "user-uuid-123",
          "product_id": "game-uuid-abc",
          "added_at": "2024-03-15T09:00:00Z",
          "entitlement_source": "purchase" // "gift", "key_activation"
        }
        ```
    *   Логика обработки: Если пользователь настроил автоматическую загрузку приобретенных игр, Download Service может инициировать загрузку последней доступной версии игры для этого пользователя. В противном случае, просто регистрирует право на загрузку.
*   **`user.settings.download.preferences.updated.v1`** (от Account Service или User Settings Service)
    *   Описание: Пользователь обновил свои предпочтения по загрузкам (например, ограничение скорости, расписание).
    *   Ожидаемый Payload:
        ```json
        {
          "user_id": "user-uuid-123",
          "preferences": {
            "max_concurrent_downloads": 2,
            "bandwidth_limit_kbps": 10000, // 10 Mbps
            "enable_auto_updates_for_all": false,
            "specific_game_auto_update_settings": { "game-uuid-xyz": true }
          }
        }
        ```
    *   Логика обработки: Обновить сохраненные настройки загрузки для пользователя. Применить новые настройки к текущим и будущим загрузкам.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Catalog Service:** Получение метаданных игр, информации о версиях, файловых манифестов, хешей файлов, ссылок на CDN/S3 (через gRPC).
*   **Library Service:** Проверка прав пользователя на загрузку/обновление игры (через gRPC). Уведомление Library Service о статусе установки/обновления игры (через Kafka или gRPC).
*   **Auth Service:** Аутентификация пользователя для API запросов, валидация токенов (обычно через API Gateway, но может быть прямой вызов gRPC).
*   **Account Service (или User Settings Service):** Получение и обновление настроек загрузки пользователя (например, лимиты скорости, авто-обновления) (через gRPC).
*   **API Gateway:** Проксирование REST API и WebSocket соединений от клиентского приложения.

### 6.2. Внешние Системы
*   **CDN (Content Delivery Network):** Основной источник для загрузки файлов игр и обновлений клиентами. Download Service генерирует и предоставляет клиентам безопасные, возможно, временные ссылки на ресурсы в CDN. Может взаимодействовать с API CDN для управления кэшем или получения статистики.
*   **S3-совместимое хранилище (опционально):** Может использоваться для временного хранения файлов перед их загрузкой в CDN, для хранения дельта-патчей или оригинальных билдов, если CDN используется только как кэширующий слой.

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `DOWNLOAD_SERVICE_HTTP_PORT`: Порт для REST API (например, `8082`).
*   `DOWNLOAD_SERVICE_GRPC_PORT`: Порт для gRPC API (например, `9092`).
*   `DOWNLOAD_SERVICE_WS_PORT`: Порт для WebSocket API (например, `8083`).
*   `POSTGRES_DSN`: Строка подключения к PostgreSQL.
*   `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB_DOWNLOAD`: Параметры Redis.
*   `KAFKA_BROKERS`: Список брокеров Kafka.
*   `KAFKA_TOPIC_DOWNLOAD_EVENTS`: Топик для публикуемых событий Download Service.
*   `CDN_PRIMARY_BASE_URL`: Базовый URL основного CDN. Пример: `https://cdn-provider1.example.com/content`.
*   `CDN_SECONDARY_BASE_URL` (опционально): URL резервного CDN.
*   `CDN_LINK_SECRET_KEY`: Секретный ключ для генерации подписанных/защищенных ссылок CDN (если используется).
*   `CDN_LINK_TTL_SECONDS`: Время жизни генерируемых ссылок CDN.
*   `S3_TEMP_STORAGE_ENDPOINT`, `S3_TEMP_STORAGE_ACCESS_KEY`, `S3_TEMP_STORAGE_SECRET_KEY`, `S3_TEMP_STORAGE_BUCKET_NAME`: Параметры S3, если используется для временного хранения.
*   `LOG_LEVEL`: Уровень логирования.
*   `AUTH_SERVICE_GRPC_ADDR`: Адрес gRPC Auth Service.
*   `CATALOG_SERVICE_GRPC_ADDR`: Адрес gRPC Catalog Service.
*   `LIBRARY_SERVICE_GRPC_ADDR`: Адрес gRPC Library Service.
*   `ACCOUNT_SERVICE_GRPC_ADDR`: Адрес gRPC Account Service.
*   `DOWNLOAD_CHUNK_SIZE_MB`: Размер чанка для параллельной загрузки (например, `16`).
*   `MAX_CONCURRENT_DOWNLOADS_PER_USER`: Максимальное количество одновременных загрузок для одного пользователя.
*   `DEFAULT_DOWNLOAD_PRIORITY`: Приоритет загрузки по умолчанию.
*   `HASH_VERIFICATION_ALGORITHM`: Алгоритм хеширования для проверки файлов (например, `SHA256`).
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`: Эндпоинт Jaeger.

### 7.2. Файлы Конфигурации (если применимо)
*   **`configs/download_config.yaml`**: Может использоваться для более сложных или менее часто изменяемых настроек.
    ```yaml
    server:
      http_port: ${DOWNLOAD_SERVICE_HTTP_PORT:-8082}
      grpc_port: ${DOWNLOAD_SERVICE_GRPC_PORT:-9092}
      ws_port: ${DOWNLOAD_SERVICE_WS_PORT:-8083}

    download_manager:
      max_concurrent_tasks_global: 1000 # Общий лимит одновременных задач обработки в сервисе
      default_chunk_size_bytes: ${DOWNLOAD_CHUNK_SIZE_MB:-16} * 1024 * 1024
      max_retries_per_item: 5
      retry_delay_seconds_base: 10
      verification_on_complete: true # Включить проверку целостности после каждой загрузки файла

    cdn_strategy:
      default_provider: "primary_cdn" # Имя провайдера по умолчанию
      providers:
        primary_cdn:
          base_url: ${CDN_PRIMARY_BASE_URL}
          api_key_env: "CDN_API_KEY_PRIMARY" # Имя переменной окружения с ключом
          # Другие специфичные настройки для этого CDN
        secondary_cdn:
          base_url: ${CDN_SECONDARY_BASE_URL}
          api_key_env: "CDN_API_KEY_SECONDARY"
      selection_policy: "geo_proximity_latency" # round_robin, failover_only, geo_proximity_latency

    bandwidth_management:
      global_limit_mbps: ${GLOBAL_BANDWIDTH_LIMIT_MBPS:-0} # 0 - без ограничений
      user_default_limit_kbps: 0 # Пользовательские настройки из Account Service будут иметь приоритет

    # Настройки для дельта-обновлений
    delta_updates:
      enabled: true
      patch_tool_path: "/usr/bin/bspatch" # Пример пути к утилите для применения патчей
      max_patch_size_gb: 5
    ```

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   Использование паттерна Circuit Breaker при взаимодействии с внешними CDN и другими микросервисами.
*   Механизмы Retry с экспоненциальной задержкой для временных сетевых ошибок или ошибок CDN.
*   Fallback на альтернативные CDN или источники файлов, если это возможно.
*   Подробное логирование всех ошибок с `trace_id`.

### 8.2. Распространенные Коды Ошибок (для REST API)
*   **`400 Bad Request` (`INVALID_ARGUMENT`)**: Некорректные параметры запроса (например, неверный `product_id`).
*   **`401 Unauthorized` (`UNAUTHENTICATED`)**: Ошибка аутентификации пользователя.
*   **`403 Forbidden` (`PERMISSION_DENIED`)**: У пользователя нет прав на загрузку данного продукта (проверка через Library Service).
*   **`404 Not Found` (`PRODUCT_NOT_FOUND`, `VERSION_NOT_FOUND`, `FILE_NOT_FOUND`)**: Запрошенный продукт, версия или файл не найдены в каталоге или системе.
*   **`429 Too Many Requests` (`RATE_LIMIT_EXCEEDED`)**: Превышен лимит запросов на создание загрузок или другие операции.
*   **`502 Bad Gateway` (`CDN_ERROR`)**: Ошибка при взаимодействии с CDN (не удалось получить валидную ссылку, CDN вернул ошибку).
*   **`503 Service Unavailable` (`SERVICE_UNAVAILABLE`, `QUEUE_OVERLOADED`)**: Сервис временно недоступен, перегружен или его зависимости (БД, Redis, Kafka) недоступны.
*   **`507 Insufficient Storage` (`INSUFFICIENT_CLIENT_DISK_SPACE`)**: (Может сообщаться клиентом, сервис логирует) У клиента недостаточно места.

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   Все запросы к API (REST, WebSocket, gRPC для клиентов) требуют JWT аутентификации пользователя. Валидация токенов через Auth Service (или локально с использованием JWKS).
*   Межсервисное gRPC взаимодействие защищено с помощью mTLS.

### 9.2. Авторизация
*   Проверка прав пользователя на загрузку конкретного продукта/версии осуществляется через Library Service перед инициацией загрузки.
*   Доступ к административным функциям сервиса (если таковые будут) защищен ролями.

### 9.3. Защита Данных и Контента
*   Использование HTTPS для всех клиентских коммуникаций.
*   Генерация безопасных, временных (time-limited) и одноразовых (если возможно CDN) ссылок на загрузку с CDN. Это предотвращает несанкционированное распространение прямых ссылок.
*   Проверка хеш-сумм (SHA256 или аналогичный) всех загруженных файлов и чанков для гарантии целостности данных.
*   Цифровая подпись критически важных файлов (например, исполняемых файлов установщика/обновления клиента платформы, основных игровых исполняемых файлов) для защиты от подмены.
*   Защита от несанкционированного доступа к файлам в S3 (если используется) через строгие политики доступа.

### 9.4. Управление Секретами
*   Секретные ключи для генерации подписанных URL для CDN, пароли к базам данных, ключи для Kafka должны храниться в Kubernetes Secrets или HashiCorp Vault.

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Стандартный многоэтапный Dockerfile для Go-приложения.
*   **Kubernetes манифесты/Helm-чарты:** Для управления развертыванием (Deployment/StatefulSet), сервисами (Service), конфигурациями (ConfigMap), секретами (Secret), автомасштабированием (HPA).
*   (Ссылка на `project_deployment_standards.md`).

### 10.2. Зависимости при Развертывании
*   Кластер Kubernetes.
*   PostgreSQL, Redis, Kafka.
*   Доступность Auth Service, Catalog Service, Library Service, Account Service.
*   Настроенное взаимодействие с CDN провайдерами.
*   S3-совместимое хранилище (если используется).

### 10.3. CI/CD
*   Автоматизированная сборка, юнит- и интеграционное тестирование.
*   Сборка Docker-образа, публикация в registry.
*   Развертывание в окружения с использованием GitOps (ArgoCD/Flux).
*   Тестирование производительности и нагрузочное тестирование эндпоинтов загрузки.

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат:** Структурированные логи в формате JSON (Zap).
*   **Ключевые события:** Начало/завершение/ошибка загрузки сессии и отдельных файлов, проверка целостности, генерация ссылок CDN, ошибки взаимодействия с другими сервисами, изменения статусов.
*   **Интеграция:** С централизованной системой логирования (Loki, ELK Stack).
*   (Ссылка на `project_observability_standards.md`).

### 11.2. Мониторинг
*   **Метрики (Prometheus):**
    *   Количество активных загрузок, загрузок в очереди.
    *   Средняя и пиковая скорость загрузки (глобально, по CDN, по регионам).
    *   Количество ошибок загрузки (по типам ошибок, по CDN).
    *   Процент успешных загрузок / проверок целостности.
    *   Время ответа API эндпоинтов.
    *   Задержки в очередях Redis/Kafka (если используются для задач).
    *   Производительность и ошибки при работе с PostgreSQL, Redis.
    *   Статистика использования CDN (количество запросов к CDN, объем трафика).
*   **Дашборды (Grafana):** Визуализация состояния системы загрузок, производительности, ошибок, использования CDN.
*   **Алертинг (AlertManager):** Срабатывание при высоком проценте ошибок загрузки, недоступности CDN, переполнении очередей, высокой задержке API.
*   (Ссылка на `project_observability_standards.md`).

### 11.3. Трассировка
*   **Интеграция:** OpenTelemetry SDK для Go. Экспорт трейсов в Jaeger или Tempo.
*   Трассировка полного жизненного цикла запроса на загрузку: от инициации до завершения, включая вызовы к другим сервисам и CDN.
*   (Ссылка на `project_observability_standards.md`).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:**
    *   API для инициации загрузки и получения ссылок: P95 < 200 мс.
    *   Скорость загрузки файлов: должна быть максимально приближена к пропускной способности канала пользователя и возможностям CDN (цель > 100 Mbps для большинства пользователей с хорошим каналом).
    *   WebSocket обновления прогресса: задержка < 1 секунды.
*   **Масштабируемость:**
    *   Поддержка до 1 миллиона одновременных загрузок файлов (чанков).
    *   Способность утилизировать пропускную способность CDN до нескольких Терабит/сек.
    *   Горизонтальное масштабирование инстансов Download Service.
*   **Надежность:**
    *   Доступность сервиса: 99.95%.
    *   Успешное завершение загрузок (с учетом retry): > 99.5%.
    *   Поддержка возобновления загрузок после обрыва связи или перезапуска клиента.
    *   Возможность использования резервных CDN при сбое основного.
*   **Безопасность:**
    *   Гарантия целостности файлов (проверка хеш-сумм).
    *   Безопасная передача файлов (HTTPS).
    *   Защита от несанкционированного доступа к файлам (подписанные URL).

## 13. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API, Protobuf определения для gRPC API, и форматы сообщений WebSocket будут храниться в соответствующих репозиториях или артефактах CI/CD.
*   Полные DDL схемы базы данных будут поддерживаться в актуальном состоянии в системе миграций.
*   TODO: Добавить ссылки на репозиторий с OpenAPI/Protobuf и на систему управления миграциями БД, когда они будут определены.

---
*Этот документ является основной спецификацией для Download Service и должен поддерживаться в актуальном состоянии.*

## 14. Связанные Рабочие Процессы (Related Workflows)
*   [Game Purchase and Library Update](../../../project_workflows/game_purchase_flow.md)
*   [Client Application Update Flow](../../../project_workflows/client_update_flow.md) (TODO: Создать этот документ, если его нет)
