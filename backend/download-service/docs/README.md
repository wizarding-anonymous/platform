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
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Управление загрузками игр:** Инициация, приостановка, возобновление и отмена загрузок. Управление очередью загрузок. Поддержка параллельных загрузок нескольких файлов/частей файлов.
*   **Управление обновлениями игр:** Автоматическая и ручная проверка наличия обновлений. Поддержка дельта-обновлений (патчей) для минимизации объема скачиваемых данных. Возможность отката к предыдущей версии (если поддерживается).
*   **Управление обновлениями клиентского приложения:** Доставка обновлений для основного клиентского приложения платформы, включая поддержку различных каналов (stable, beta).
*   **Проверка целостности файлов:** Верификация загруженных файлов по хеш-суммам (например, SHA256, MD5). Автоматическое исправление поврежденных или отсутствующих файлов путем их повторной загрузки.
*   **Взаимодействие с CDN:** Генерация безопасных ссылок на загрузку с CDN. Выбор оптимального CDN-сервера на основе геолокации пользователя или нагрузки.
*   **Мониторинг и статистика загрузок:** Отслеживание прогресса, текущей скорости загрузки, предполагаемого времени завершения. Сбор статистики для анализа производительности CDN и выявления проблем.
*   **Управление настройками загрузки:** Предоставление пользователям возможности настраивать параметры загрузки (например, ограничение скорости, расписание загрузок).

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+, согласно `../../../../project_technology_stack.md`).
*   **API:**
    *   REST API: Echo (`github.com/labstack/echo/v4`) или Gin (`github.com/gin-gonic/gin`) (для клиентского приложения, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
    *   gRPC: `google.golang.org/grpc` (для межсервисного взаимодействия, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
    *   WebSocket: (например, `github.com/gorilla/websocket`) для real-time обновлений прогресса загрузки на клиенте.
*   **База данных:** PostgreSQL (версия 15+) для хранения метаданных о загрузках, файлах, версиях, статусах. Драйвер: GORM (`gorm.io/gorm`) с `gorm.io/driver/postgres` или `pgx` (`github.com/jackc/pgx/v5`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Кэширование/Очереди:** Redis (версия 7.0+) для хранения временных данных сессий загрузки, состояния активных загрузок, кэширования токенов CDN, управления небольшими очередями задач. Клиент: `go-redis/redis` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Хранилище файлов (временное/staging):** S3-совместимое объектное хранилище (например, MinIO, Yandex Object Storage) может использоваться для временного хранения файлов перед их передачей в CDN или для сборки дельта-патчей. (согласно `../../../../project_technology_stack.md`).
*   **Брокер сообщений:** Apache Kafka (клиент `github.com/confluentinc/confluent-kafka-go` или `github.com/segmentio/kafka-go`, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Инфраструктура:** Docker, Kubernetes.
*   **Мониторинг/Трассировка:** OpenTelemetry SDK, Prometheus client (`github.com/prometheus/client_golang`). (согласно `../../../../project_observability_standards.md`).
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Логирование:** Zap (`go.uber.org/zap`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   **CDN (Content Delivery Network):** Сеть доставки (и дистрибуции) контента, используемая для ускорения загрузки файлов пользователями.
*   **Дельта-обновление (Delta Update/Patch):** Обновление, содержащее только измененные части файлов по сравнению с предыдущей версией, что позволяет уменьшить объем загрузки.
*   **Манифест Файлов (File Manifest):** Список файлов, входящих в состав продукта (игры/версии), с их размерами, хеш-суммами и путями.
*   **Chunk:** Часть файла, загружаемая отдельно при параллельной или возобновляемой загрузке.
*   Для других общих терминов см. `../../../../project_glossary.md`.

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
        GRPC_API --> DownloadUseCaseSvc
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
              "version_id": "version-uuid-abc",
              "installation_path": "/opt/games/MySuperGame"
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
              "status": "queued",
              "total_size_bytes": 10737418240,
              "downloaded_bytes": 0,
              "progress_percentage": 0.0,
              "estimated_time_left_seconds": null
            }
          }
        }
        ```
    *   Пример ответа (Ошибка 403 Forbidden - нет прав на продукт):
        ```json
        {
          "errors": [
            {
              "code": "PRODUCT_ACCESS_DENIED",
              "title": "Доступ к продукту запрещен",
              "detail": "У вас нет прав на загрузку продукта с ID 'game-uuid-123'."
            }
          ]
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
              "bandwidth_limit_kbps": 5000
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): (Возвращает обновленные настройки)
    *   Требуемые права доступа: `user`.

### 3.2. gRPC API (для межсервисного взаимодействия)
(Содержимое существующего раздела актуально).

### 3.3. WebSocket API
(Содержимое существующего раздела актуально).

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`DownloadSession` (Сессия Загрузки)** (Как в существующем документе)
*   **`DownloadItem` (Элемент Загрузки - файл в рамках сессии)** (Как в существующем документе)
*   **`FileMetadata` (Метаданные Файла Продукта)** (Как в существующем документе)
*   **`Update` (Обновление)** (Как в существующем документе)
*   **`UpdateHistory` (История Обновлений Пользователя)**
    *   `id` (UUID): Уникальный идентификатор записи.
    *   `user_id` (UUID): ID пользователя.
    *   `product_id` (UUID): ID продукта.
    *   `applied_version_id` (UUID): ID установленной версии.
    *   `previous_version_id` (UUID, Nullable): ID предыдущей установленной версии (если это было обновление).
    *   `applied_at` (TIMESTAMPTZ): Время применения обновления.
    *   `status` (ENUM: `success`, `failed`): Статус применения обновления.
*   **`VerificationResult` (Результат Верификации Файлов)**
    *   `id` (UUID): Уникальный идентификатор записи о верификации.
    *   `user_id` (UUID): ID пользователя.
    *   `product_id` (UUID): ID продукта.
    *   `version_id` (UUID): ID проверяемой версии.
    *   `status` (ENUM: `pending`, `in_progress`, `completed_ok`, `completed_errors`, `failed_to_repair`): Статус верификации.
    *   `corrupted_files_count` (INTEGER): Количество обнаруженных поврежденных файлов.
    *   `repaired_files_count` (INTEGER): Количество автоматически восстановленных файлов.
    *   `files_requiring_redownload` (JSONB, Nullable): Список файлов (их `file_metadata_id` или пути), которые не удалось восстановить и требуется их полная перезагрузка.
    *   `started_at` (TIMESTAMPTZ): Время начала верификации.
    *   `completed_at` (TIMESTAMPTZ, Nullable): Время завершения верификации.

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL
**ERD Диаграмма (дополненная):**
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
    UPDATES_HISTORY {
        UUID id PK
        UUID user_id "FK (User)"
        UUID product_id "FK (Product)"
        UUID applied_version_id "FK (ProductVersion)"
        UUID previous_version_id "nullable, FK (ProductVersion)"
        VARCHAR status "ENUM('success', 'failed')"
        TIMESTAMPTZ applied_at
    }
    VERIFICATION_RESULTS {
        UUID id PK
        UUID user_id "FK (User)"
        UUID product_id "FK (Product)"
        UUID version_id "FK (ProductVersion)"
        VARCHAR status "ENUM('pending', 'in_progress', 'completed_ok', 'completed_errors', 'failed_to_repair')"
        INTEGER corrupted_files_count
        INTEGER repaired_files_count
        JSONB files_requiring_redownload
        TIMESTAMPTZ started_at
        TIMESTAMPTZ completed_at
    }

    DOWNLOAD_SESSIONS ||--|{ DOWNLOAD_ITEMS : "contains"
    DOWNLOAD_ITEMS }o--|| FILE_METADATA : "references"
    FILE_METADATA }o--|| PRODUCTS_VERSIONS : "belongs_to_version"
    USERS ||--o{ DOWNLOAD_SESSIONS : "owns"
    USERS ||--o{ UPDATES_HISTORY : "has_applied"
    USERS ||--o{ VERIFICATION_RESULTS : "initiated_by"
    PRODUCTS_VERSIONS ||--o{ VERIFICATION_RESULTS : "verifies"
    PRODUCTS_VERSIONS ||--o{ UPDATES_HISTORY : "applied_version"
    PRODUCTS_VERSIONS ||--o{ UPDATES_HISTORY : "previous_version"

    PRODUCTS_VERSIONS {
        note "Managed by Catalog Service"
    }
    USERS {
        note "Managed by Auth/Account Service"
    }
```

**DDL (PostgreSQL - дополнения для `updates_history`, `verification_results`):**
```sql
CREATE TABLE updates_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    applied_version_id UUID NOT NULL,
    previous_version_id UUID,
    status VARCHAR(50) NOT NULL CHECK (status IN ('success', 'failed')),
    applied_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    error_message TEXT
);
CREATE INDEX idx_updates_history_user_product ON updates_history(user_id, product_id);
COMMENT ON TABLE updates_history IS 'История установки обновлений продуктов пользователями.';

CREATE TABLE verification_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    version_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed_ok', 'completed_errors', 'failed_to_repair')),
    corrupted_files_count INTEGER NOT NULL DEFAULT 0,
    repaired_files_count INTEGER NOT NULL DEFAULT 0,
    files_requiring_redownload JSONB,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ
);
CREATE INDEX idx_verification_results_user_product ON verification_results(user_id, product_id);
COMMENT ON TABLE verification_results IS 'Результаты проверки целостности файлов установленных продуктов.';
```

#### 4.2.2. Redis
*   **Очереди загрузок:** Используются Redis Lists или Sorted Sets (`download_queue:user:<user_id>`) для управления очередью загрузок конкретного пользователя, включая приоритеты.
*   **Состояние активных сессий загрузки:** Redis Hashes (`download_session_state:<session_id>`) для хранения часто обновляемой информации о прогрессе, скорости, статусе активных загрузок. Это позволяет быстро отдавать данные через WebSocket и снижает нагрузку на PostgreSQL.
*   **Кэш токенов/временных ссылок CDN:** Redis Strings (`cdn_link:<file_path_hash_or_id>`) с TTL.
*   **Блокировки:** Для предотвращения гонок состояний при обновлении сессий или файлов (`lock:download_session:<session_id>`).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Формат событий:** CloudEvents v1.0 JSON (согласно `../../../../project_api_standards.md`).
*   **Основной топик Kafka:** `com.platform.download.events.v1`.

*   **`com.platform.download.session.status.changed.v1`**
    *   `data` Payload: (Как в существующем документе, с корректным `type`).
*   **`com.platform.download.update.available.v1`**
    *   `data` Payload: (Как в существующем документе, с корректным `type`).
*   **`com.platform.download.file.verification.completed.v1`**
    *   `data` Payload: (Как в существующем документе, с корректным `type`).
*   **`com.platform.download.item.status.changed.v1`**
    *   Описание: Статус загрузки отдельного файла (элемента) в рамках сессии изменился.
    *   `data` Payload:
        ```json
        {
          "downloadSessionId": "download-session-uuid-xyz",
          "downloadItemId": "item-uuid-123",
          "fileMetadataId": "filemeta-uuid-abc",
          "relativePath": "bin/game.exe",
          "newStatus": "completed",
          "previousStatus": "downloading",
          "downloadedBytes": 104857600,
          "totalFileSizeBytes": 104857600,
          "timestamp": "2024-03-18T15:00:00Z",
          "errorDetails": null
        }
        ```
    *   Потребители: Library Service (для обновления статуса установки), Analytics Service.

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
*   **Безопасная генерация ссылок CDN:** Использование time-limited, подписанных URL (если поддерживается CDN) для предотвращения неавторизованного доступа и распространения.
*   **Проверка целостности файлов:** Обязательная проверка хеш-сумм после загрузки для гарантии отсутствия повреждений или модификаций.
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

### 14.1. PostgreSQL (Метаданные сессий загрузок, истории обновлений, результатов верификации)
*   **Процедура резервного копирования:**
    *   Ежедневный логический бэкап (`pg_dump`).
    *   Настроена непрерывная архивация WAL-сегментов (PITR), базовый бэкап еженедельно.
    *   **Хранение:** Бэкапы в S3, шифрование, версионирование, другой регион. Срок хранения: полные - 30 дней, WAL - 14 дней.
*   **Процедура восстановления:** Тестируется ежеквартально.
*   **RTO:** < 2 часов.
*   **RPO:** < 15 минут.
*   (Общие принципы см. `../../../../project_database_structure.md`).

### 14.2. Redis (Очереди загрузок, состояние активных сессий, кэш токенов CDN)
*   **Стратегия персистентности:**
    *   **AOF (Append Only File):** Включен с fsync `everysec` для очередей и состояний активных сессий, если их потеря критична для пользовательского опыта (например, чтобы не прерывать все активные загрузки при перезапуске Redis).
    *   **RDB Snapshots:** Регулярное создание снапшотов (например, каждые 1-6 часов).
*   **Резервное копирование (снапшотов):** RDB-снапшоты могут копироваться в S3 ежедневно. Срок хранения - 7 дней.
*   **Восстановление:** Из последнего RDB-снапшота и/или AOF. Большинство данных (кроме, возможно, долгоживущих очередей) могут быть перестроены или сессии будут переинициализированы клиентами.
*   **RTO:** < 30 минут.
*   **RPO:** < 1 минуты (для данных с AOF `everysec`). Для кэша токенов CDN RPO менее критичен.

### 14.3. S3-совместимое хранилище (если используется для временного staging)
*   **Стратегия:** Данные во временном хранилище обычно имеют короткий срок жизни. Основные файлы игр управляются Developer/Catalog Service и их S3 бакетами. Download Service не отвечает за бэкап этих основных файлов.
*   **Резервное копирование:** Обычно не требуется для временных файлов, специфичных для Download Service. Если есть критичные staging-файлы, создаваемые самим Download Service, можно настроить версионирование или репликацию в S3.
*   **RTO/RPO:** Неприменимо для временных данных; для критичных staging-данных зависит от настроек S3.

### 14.4. Общая стратегия
*   Восстановление PostgreSQL является приоритетным.
*   Redis восстанавливается для минимизации прерываний текущих операций.
*   Процедуры документированы и тестируются.

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [Процесс покупки игры и обновления библиотеки пользователя](../../../../project_workflows/game_purchase_flow.md)
*   [Процесс обновления клиентского приложения](../../../../project_workflows/client_update_flow.md) (Примечание: Создание документа `client_update_flow.md` является частью общей задачи по документированию проекта и выходит за рамки обновления документации данного микросервиса.)

---
*Этот документ является основной спецификацией для Download Service и должен поддерживаться в актуальном состоянии.*
