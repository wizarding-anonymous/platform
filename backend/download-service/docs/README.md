# Спецификация Микросервиса: Download Service

**Версия:** 1.0
**Дата последнего обновления:** 2025-05-25

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Download Service.
*   **Роль в общей архитектуре платформы:** Download Service является критически важным компонентом, отвечающим за загрузку, обновление и установку игр на устройства пользователей. Обеспечивает надежную и эффективную доставку игрового контента.
*   **Основные бизнес-задачи:** Управление загрузкой, обновлением и установкой игр, оптимизация использования сетевых ресурсов, обеспечение целостности и безопасности файлов, предоставление информации о статусе загрузок.

### 1.2. Ключевые Функциональности
*   **Управление загрузками игр:** Инициация, приостановка, возобновление, отмена, управление очередью, параллельные загрузки.
*   **Управление обновлениями игр:** Проверка наличия, автоматическое/ручное обновление, дельта-обновления, откат.
*   **Управление клиентским приложением:** Обновление клиента платформы, поддержка различных каналов обновлений.
*   **Проверка целостности:** Верификация файлов по хеш-суммам, автоматическое исправление.
*   **Мониторинг и статистика:** Отслеживание прогресса, скорости загрузки, производительности CDN, использования ресурсов.
*   (Источник: Спецификация Download Service, раздел 2.3)

### 1.3. Основные Технологии
*   **Язык программирования:** Go.
*   **API:** gRPC (межсервисное), REST (клиентское), WebSocket (real-time прогресс).
*   **Базы данных:** PostgreSQL (метаданные загрузок), Redis (кэширование, очереди).
*   **Хранилище файлов:** MinIO или S3-совместимое (для временных файлов, если требуется).
*   **Инфраструктура:** Docker, Kubernetes, Prometheus, Jaeger, ELK/Loki.
*   (Источник: Спецификация Download Service, раздел 8.1)

### 1.4. Термины и Определения (Glossary)
*   **CDN (Content Delivery Network):** Сеть доставки контента.
*   **Дельта-обновление:** Обновление, содержащее только изменения между версиями.
*   **Хеш-сумма:** Уникальная строка для проверки целостности файла.
*   **Circuit Breaker:** Паттерн для предотвращения каскадных отказов.
*   **mTLS (mutual TLS):** Взаимная TLS-аутентификация.
*   (Полный глоссарий см. в Спецификации Download Service, раздел 9)

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Download Service построен на многослойной архитектуре, обеспечивающей разделение ответственности и масштабируемость.
*   Ключевые компоненты включают управление загрузками, обновлениями, проверку целостности, взаимодействие с CDN и управление очередями.
*   (Источник: Спецификация Download Service, раздел 3.1)

### 2.2. Слои Сервиса
(На основе раздела 3.1.1 исходной спецификации)

#### 2.2.1. Presentation Layer (Слой Представления / Транспортный слой)
*   Ответственность: Обработка входящих запросов от клиентского приложения (REST, WebSocket) и других микросервисов (gRPC).
*   Ключевые компоненты/модули:
    *   REST API Handlers (Gin или Echo).
    *   gRPC Service Implementations.
    *   WebSocket Handlers (для отслеживания прогресса).

#### 2.2.2. Application Layer (Прикладной Слой / Сервисный слой)
*   Ответственность: Реализация бизнес-логики управления загрузками, обновлениями, проверкой целостности. Оркестрация взаимодействия компонентов.
*   Ключевые компоненты/модули:
    *   Сервисы: `DownloadServiceLogic`, `UpdateServiceLogic`, `VerificationServiceLogic`, `ClientUpdateServiceLogic`.
    *   (Компоненты из раздела 3.2 исходной спецификации: Download Manager, Update Manager, Integrity Checker, Delta Processor).

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности (Download, Update, FileMetadata), правила валидации.
*   Ключевые компоненты/модули:
    *   Entities: `Download`, `DownloadItem`, `Update`, `UpdateItem`, `FileMetadata`, `VerificationResult`, `VerificationIssue`.
    *   Enums: `DownloadStatus`, `ItemStatus`, `UpdateType`, `VerificationStatus`.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Взаимодействие с PostgreSQL, Redis, CDN. Управление очередями. Мониторинг, логирование, трассировка.
*   Ключевые компоненты/модули:
    *   PostgreSQL Repositories, Redis Repositories.
    *   CDN Client.
    *   Queue Manager (может использовать Redis или другую систему очередей).
    *   Модули для хеширования, обработки дельта-патчей.
    *   Клиенты для взаимодействия с Catalog, Library, Auth, Account Services.

*   (Структуру проекта см. в Спецификации Download Service, раздел 3.1.2)

## 3. API Endpoints

### 3.1. REST API (для клиентского приложения)
*   **Префикс:** `/api/v1/downloads` (пример)
*   **Аутентификация:** JWT (через Auth Service).
*   **Основные группы эндпоинтов:**
    *   Управление загрузками: `POST /`, `GET /`, `GET /{download_id}`, `PATCH /{download_id}` (статус), `PATCH /{download_id}/priority`.
    *   Управление обновлениями: `GET /updates/check`, `POST /updates/apply`, `GET /updates/history`.
    *   Проверка целостности: `POST /verification`, `GET /verification/{verification_id}`, `POST /verification/{verification_id}/fix`.
    *   Настройки загрузки: `GET /settings`, `PATCH /settings`, `POST /settings/bandwidth-schedule`.
*   (Детали см. в Спецификации Download Service, раздел 5.3).

### 3.2. gRPC API (для межсервисного взаимодействия)
*   **Сервисы:** `DownloadService`, `UpdateService`, `VerificationService`, `SettingsService`.
*   **Основные методы:** `InitiateDownload`, `GetDownloadStatus`, `CheckForUpdates`, `VerifyGameFiles`, `SyncGameInstallStatus`.
*   (Детали см. в Спецификации Download Service, раздел 5.4).

### 3.3. WebSocket API
*   **Эндпоинт:** `/api/v1/ws/downloads`
*   **Назначение:** Получение обновлений о прогрессе загрузок в реальном времени.
*   (Детали см. в Спецификации Download Service, раздел 5.3.5).

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности
*   **Download**: Информация о загрузке (ID, UserID, GameID, Status, Progress, Speed, Size).
*   **DownloadItem**: Элемент загрузки/файл (ID, FileID, FileName, Path, Size, Hash, Status).
*   **Update**: Информация об обновлении (ID, GameID, FromVersion, ToVersion, Type, Size).
*   **FileMetadata**: Метаданные файла (ID, GameID, FileName, Path, Size, Hash, Version, CDNLocations).
*   **VerificationResult**: Результат проверки целостности.
*   (Go структуры см. в Спецификации Download Service, раздел 5.1.1).

### 4.2. Схема Базы Данных
*   **PostgreSQL**: Хранит метаданные загрузок, обновлений, файлов, историю, статистику CDN, результаты верификации, настройки.
    ```sql
    -- Пример таблицы downloads (сокращенно)
    CREATE TABLE downloads (id VARCHAR PRIMARY KEY, user_id VARCHAR, game_id VARCHAR, status VARCHAR, progress REAL ...);
    -- Пример таблицы file_metadata (сокращенно)
    CREATE TABLE file_metadata (id VARCHAR PRIMARY KEY, game_id VARCHAR, file_name VARCHAR, file_path VARCHAR, file_size BIGINT, file_hash VARCHAR ...);
    ```
*   **Redis**: Хранит очереди загрузок, статусы текущих загрузок, кэш метаданных, временные токены CDN.
*   (Основные таблицы и структуры Redis см. в Спецификации Download Service, разделы 3.3.1, 3.3.2, 5.2).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   `download.started`: Начало загрузки.
*   `download.progress`: Обновление прогресса загрузки (может идти через WebSocket).
*   `download.completed`: Завершение загрузки.
*   `download.failed`: Ошибка загрузки.
*   `update.available`: Доступно обновление.
*   `update.completed`: Обновление завершено.
    *   `Структура Payload (пример):`
        ```json
        {
          "event_id": "uuid_event",
          "event_type": "update.completed.v1",
          "timestamp": "ISO8601_timestamp",
          "source_service": "download-service",
          "user_id": "uuid_user",
          "game_id": "uuid_game",
          "version_id": "uuid_game_version_updated_to",
          "update_duration_seconds": 600
        }
        ```
*   `verification.completed`: Проверка целостности завершена.
    *   `Структура Payload (пример):`
        ```json
        {
          "event_id": "uuid_event",
          "event_type": "verification.completed.v1",
          "timestamp": "ISO8601_timestamp",
          "source_service": "download-service",
          "user_id": "uuid_user",
          "game_id": "uuid_game",
          "version_id": "uuid_game_version_verified",
          "status": "success" | "failed_needs_repair",
          "corrupted_files_count": 0,
          "repaired_files_count": 0
        }
        ```
*   **Система сообщений:** Apache Kafka
*   **Основные топики:** `download.lifecycle.events` (для событий начала, завершения, отмены загрузки), `download.progress.events` (для обновлений прогресса, если не используется исключительно WebSocket).

### 5.2. Потребляемые События (Consumed Events)
*   `catalog.game.version.published`:
    *   **Источник:** Catalog Service
    *   **Назначение:** Уведомление Download Service о том, что новая версия игры опубликована и готова для загрузки пользователями. Download Service должен обновить свои метаданные, подготовить ссылки на CDN.
    *   **Структура Payload (пример):**
        ```json
        {
          "game_id": "uuid_game",
          "version_id": "uuid_game_version",
          "version_string": "1.1.0",
          "manifest_url": "url_to_download_manifest_from_catalog_or_developer_service",
          "file_hashes": [ { "file_path": "/bin/game.exe", "hash": "sha256_hash_value" } ],
          "total_size_bytes": 10737418240,
          "published_at": "ISO8601_timestamp"
        }
        ```
    *   **Логика обработки:** Добавить/обновить метаданные файлов версии в своей базе данных (`file_metadata`). Подготовить информацию для предоставления клиентам ссылок на загрузку через CDN.
*   `library.game.install.requested`:
    *   **Источник:** Library Service
    *   **Назначение:** Уведомление Download Service о том, что пользователь запросил установку игры, и права доступа подтверждены.
    *   **Структура Payload (пример):**
        ```json
        {
          "user_id": "uuid_user",
          "game_id": "uuid_game",
          "version_id": "uuid_game_version_to_install",
          "entitlement_id": "uuid_library_item",
          "requested_at": "ISO8601_timestamp"
        }
        ```
    *   **Логика обработки:** Инициировать процесс загрузки для пользователя, если он еще не начат. Добавить в очередь загрузок.
*   `catalog.game.version.deleted`:
    *   **Источник:** Catalog Service
    *   **Назначение:** Уведомление о том, что определенная версия игры была удалена из каталога. Download Service должен прекратить предоставление этой версии для загрузки.
    *   **Структура Payload (пример):**
        ```json
        {
          "game_id": "uuid_game",
          "version_id": "uuid_game_version",
          "deleted_at": "ISO8601_timestamp",
          "reason": "Версия содержит критическую уязвимость"
        }
        ```
    *   **Логика обработки:** Пометить метаданные файлов для данной версии как неактивные. Запретить новые загрузки этой версии.
*   `user.account.deleted`:
    *   **Источник:** Account Service
    *   **Назначение:** Уведомление об удалении аккаунта пользователя. Download Service может потребоваться анонимизировать или удалить историю загрузок этого пользователя.
    *   **Структура Payload (пример):**
        ```json
        {
          "user_id": "uuid_user",
          "deleted_at": "ISO8601_timestamp",
          "anonymization_required": true
        }
        ```
    *   **Логика обработки:** Анонимизировать или удалить историю загрузок пользователя. Отменить активные загрузки.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Catalog Service**: Получение метаданных игр, информации о версиях, обновлениях, хешах файлов (gRPC).
*   **Library Service**: Проверка прав доступа, обновление статуса установки игры (gRPC).
*   **Auth Service**: Аутентификация/авторизация пользователей (gRPC).
*   **Account Service**: Получение/обновление настроек загрузки пользователя (gRPC).
*   (Детали и обработка ошибок см. в Спецификации Download Service, раздел 6).

### 6.2. Внешние Системы
*   **CDN (Content Delivery Network)**: Загрузка файлов игр и обновлений (HTTPS).
*   (Детали и обработка ошибок см. в Спецификации Download Service, раздел 6.5).

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `DOWNLOAD_SERVICE_HTTP_PORT`, `DOWNLOAD_SERVICE_GRPC_PORT`, `DOWNLOAD_SERVICE_WS_PORT`.
*   `POSTGRES_DSN`.
*   `REDIS_ADDR`.
*   `CDN_BASE_URL_PRIMARY`
*   `CDN_SECONDARY_BASE_URL` (optional, for fallback)
*   `CDN_API_KEY_PRIMARY` (if CDN requires API key for URL signing/management)
*   `S3_TEMP_STORAGE_ENDPOINT` (if used for temporary storage before CDN distribution)
*   `S3_TEMP_STORAGE_ACCESS_KEY_ID`
*   `S3_TEMP_STORAGE_SECRET_ACCESS_KEY`
*   `S3_TEMP_STORAGE_BUCKET`
*   `KAFKA_TOPIC_DOWNLOAD_LIFECYCLE_EVENTS` (e.g., `download.lifecycle.events`)
*   `KAFKA_TOPIC_DOWNLOAD_PROGRESS_EVENTS` (e.g., `download.progress.events`)
*   `AUTH_SERVICE_GRPC_ADDR`
*   `CATALOG_SERVICE_GRPC_ADDR`
*   `LIBRARY_SERVICE_GRPC_ADDR`
*   `ACCOUNT_SERVICE_GRPC_ADDR` (for user settings)
*   `LOG_LEVEL` (e.g., `info`, `debug`)
*   `DOWNLOAD_CHUNK_SIZE_MB` (e.g., `10`)
*   `MAX_CONCURRENT_USER_DOWNLOADS` (e.g., `3`)
*   `GLOBAL_BANDWIDTH_LIMIT_MBPS` (e.g., `1000`)
*   `DEFAULT_DOWNLOAD_PRIORITY` (e.g., `5`)
*   `HASH_VERIFICATION_ALGORITHM` (e.g., `SHA256`)
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`

### 7.2. Файлы Конфигурации (если применимо)
*   Конфигурация сервиса осуществляется преимущественно через переменные окружения. Если потребуются файлы конфигурации для сложных настроек (например, для детализированных правил выбора CDN или стратегий обработки ошибок для разных типов файлов), их структура будет определена здесь.

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   Использование Circuit Breaker и Retry с экспоненциальной задержкой при взаимодействии с другими сервисами.
*   Fallback на кэшированные данные или альтернативные CDN.
*   Подробное логирование ошибок.

### 8.2. Распространенные Коды Ошибок
*   Для REST API: Стандартные HTTP коды (400, 401, 403, 404, 500, 503).
*   Для gRPC: Стандартные коды gRPC.
*   Внутренние ошибки: `CDN_UNAVAILABLE`, `FILE_HASH_MISMATCH`, `INSUFFICIENT_DISK_SPACE`.

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   Все запросы к API требуют JWT аутентификации (через Auth Service).
*   mTLS для межсервисного gRPC взаимодействия.

### 9.2. Авторизация
*   Проверка прав доступа к игре через Library Service перед началом загрузки.

### 9.3. Защита Данных
*   TLS 1.3 для всех внешних коммуникаций.
*   Проверка хеш-сумм всех загруженных файлов.
*   Цифровая подпись критических файлов (например, исполняемых файлов обновлений).

### 9.4. Управление Секретами
*   Ключи доступа к CDN, пароли к БД через Kubernetes Secrets или Vault.
*   **Защита от атак:** Защита CDN от DDoS, Rate limiting на API, защита от подмены контента.
*   (Детали см. в Спецификации Download Service, раздел 7.1).

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Стандартный Dockerfile для Go приложения.
*   **Kubernetes манифесты/Helm-чарты.**

### 10.2. Зависимости при Развертывании
*   PostgreSQL, Redis.
*   Catalog Service, Library Service, Auth Service, Account Service.
*   Доступ к CDN.

### 10.3. CI/CD
*   Автоматическая сборка, тестирование, развертывание.
*   (Детали см. в Спецификации Download Service, раздел 8).

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   Структурированные логи (JSON).
*   Централизованный сбор (ELK/Loki).
*   Логирование прогресса загрузок, ошибок, событий безопасности.

### 11.2. Мониторинг
*   Prometheus, Grafana.
*   Метрики: скорость загрузки, количество активных загрузок, ошибки, использование CDN, производительность проверки целостности.

### 11.3. Трассировка
*   Интеграция с Jaeger/OpenTelemetry для отслеживания запросов через CDN и внутренние компоненты.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:** Высокая пропускная способность, минимизация задержек.
*   **Масштабируемость:** Горизонтальное масштабирование, балансировка нагрузки CDN.
*   **Надежность:** Устойчивость к сетевым сбоям, автоматическое восстановление, отказоустойчивость CDN.
*   **Безопасность:** Защита от несанкционированного доступа и подмены контента.
*   (Детали см. в Спецификации Download Service, раздел 2.4).

## 13. Приложения (Appendices) (Опционально)
*   Детальные схемы DDL для PostgreSQL, полные примеры REST API и WebSocket сообщений (включая ответы на ошибки), а также форматы событий Kafka будут добавлены по мере финализации дизайна и реализации соответствующих модулей.

---
*Этот шаблон является отправной точкой и может быть адаптирован под конкретные нужды проекта и сервиса.*
