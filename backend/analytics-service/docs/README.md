<!-- backend\analytics-service\docs\README.md -->
# Спецификация Микросервиса: Analytics Service

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   Analytics Service является центральным компонентом платформы "Российский Аналог Steam", ответственным за сбор, обработку, анализ и предоставление данных и инсайтов, генерируемых на платформе.
*   Его основная роль - поддержка принятия решений на основе данных для бизнес-стратегии, операционных улучшений, улучшения пользовательского опыта, а также предоставление разработчикам релевантной статистики по их продуктам.
*   Основные бизнес-задачи: сбор и обработка данных, расчет метрик и KPI, генерация отчетов, анализ поведения пользователей, сегментация аудитории, поддержка предиктивной аналитики (прогнозирование, рекомендации, обнаружение аномалий).
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Сбор данных:** Прием событий и данных от всех микросервисов платформы (действия пользователей, транзакции, системные логи, маркетинговые взаимодействия). Поддержка потоковой и пакетной загрузки данных.
*   **Обработка данных:** Потоковая и пакетная обработка сырых данных: трансформация, агрегация, обогащение, очистка. Строгое соблюдение процедур анонимизации и псевдонимизации персональных данных (ПДн) в соответствии с ФЗ-152 для аналитических целей.
*   **Метрики и Отчетность:** Расчет ключевых показателей эффективности (KPI) и продуктовых метрик (DAU, MAU, ARPU, конверсии, удержание и т.д.). Генерация стандартных и настраиваемых отчетов для различных групп пользователей (администраторы, разработчики, маркетологи).
*   **Анализ поведения пользователей:** Предоставление инструментов и данных для анализа путей пользователей, когортного анализа, результатов A/B тестов, воронок конверсий.
*   **Сегментация аудитории:** Создание и управление статическими и динамическими сегментами пользователей на основе их атрибутов и поведения для целевого маркетинга и персонализации.
*   **Предиктивная аналитика:** Разработка, тренировка, развертывание и мониторинг моделей машинного обучения (ML) для прогнозирования продаж, оттока пользователей, генерации персонализированных рекомендаций, обнаружения мошенничества и аномалий.
*   **Мониторинг производительности системы:** Агрегация и анализ технических метрик здоровья и производительности микросервисов платформы для выявления узких мест и оптимизации.

### 1.3. Основные Технологии
*   **API Layer (Слой API и управления):**
    *   Язык программирования: Go (версия 1.21+) (согласно `../../../../project_technology_stack.md`, основной для API Layer). Java (Spring Boot) может рассматриваться для отдельных компонентов API или управления, если это будет признано целесообразным командой.
    *   REST Framework: Echo (`github.com/labstack/echo/v4`) для Go, Spring Boot для Java. (согласно `../../../../PACKAGE_STANDARDIZATION.md`)
    *   GraphQL (потенциально): Apollo Server (Node.js) или Hasura (Go/GraphQL engine) если потребуется.
    *   WebSocket (потенциально): Для real-time дашбордов.
*   **Data Processing Layer (Слой Обработки Данных):**
    *   Языки программирования: Python (версия 3.10+), Scala, Java (согласно `../../../../project_technology_stack.md`).
    *   Фреймворки обработки: Apache Spark (для пакетной обработки), Apache Flink или Kafka Streams (для потоковой обработки). (согласно `../../../../project_technology_stack.md`)
*   **Data Storage Layer (Слой Хранения Данных):**
    *   Аналитическая СУБД (DWH): ClickHouse (версия 23.x+). (согласно `../../../../project_technology_stack.md`)
    *   Data Lake (Хранилище сырых данных): S3-совместимое хранилище (например, MinIO, Yandex Object Storage). (согласно `../../../../project_technology_stack.md`)
    *   Хранилище метаданных и конфигураций: PostgreSQL (версия 15+). (согласно `../../../../project_technology_stack.md`)
*   **Messaging Layer (Слой Обмена Сообщениями):**
    *   Брокер сообщений: Apache Kafka (версия 3.x+). (согласно `../../../../project_technology_stack.md` и `../../../../PACKAGE_STANDARDIZATION.md`)
*   **Machine Learning Layer (Слой Машинного Обучения):**
    *   Библиотеки: Python с TensorFlow, PyTorch, Scikit-learn.
    *   Управление моделями: MLflow.
*   **Общие для API и Processing (где применимо):**
    *   Управление конфигурацией: Viper (`github.com/spf13/viper`) для Go.
    *   Логирование: Zap (`go.uber.org/zap`) для Go.
    *   Трассировка и метрики: OpenTelemetry, Prometheus client.
*   **Визуализация данных (внешние инструменты):** Grafana (для операционных метрик и некоторых аналитических дашбордов), Apache Superset или аналогичный BI-инструмент (для продвинутой бизнес-аналитики).
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   См. `../../../../project_glossary.md`.
*   **Событие (Event):** Атомарная запись о действии пользователя или системы, являющаяся основным источником данных для аналитики.
*   **KPI (Key Performance Indicator):** Ключевой показатель эффективности, используемый для оценки достижения стратегических и операционных целей.
*   **Сегмент (Segment):** Группа пользователей или других сущностей, объединенных общими характеристиками или поведением, для таргетированного анализа или воздействия.
*   **ML Модель (Machine Learning Model):** Алгоритм, обученный на исторических данных для выполнения прогнозов или классификаций.
*   **Data Lake:** Централизованное хранилище для сырых данных в их исходном формате.
*   **DWH (Data Warehouse):** Централизованное хранилище структурированных и агрегированных данных, оптимизированное для аналитических запросов.
*   **ETL/ELT (Extract, Transform, Load / Extract, Load, Transform):** Процессы извлечения, преобразования и загрузки данных.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Analytics Service построен на основе событийно-ориентированной архитектуры, оптимизированной для обработки больших объемов данных (Big Data) и поддержки жизненного цикла ML-моделей. Он сочетает элементы Lambda и Kappa архитектур для обеспечения как пакетной (batch), так и потоковой (stream) обработки данных.
*   **Для API и управляющих компонентов** (например, управление определениями метрик, отчетов), сервис может придерживаться стандартной слоистой архитектуры (Presentation, Application, Domain, Infrastructure), где это применимо. Однако основная часть сервиса — это конвейеры данных (data pipelines).
*   Ключевые компоненты: Сбор данных (Data Ingestion), Обработка данных (Data Processing - batch/stream), Хранение данных (Data Storage - DWH, Data Lake, Metadata Store), Слой доступа к данным (Data Access API) и Слой ML моделей (ML Model Management & Serving).

### 2.2. Диаграмма Архитектуры
Ниже представлена диаграмма верхнеуровневой архитектуры Analytics Service:
```mermaid
graph TD
    subgraph "Источники Данных (Все Микросервисы Платформы)"
        direction LR
        KafkaTopics[Kafka: Topics (account.events, catalog.events, etc.)]
    end

    subgraph "Analytics Service"
        direction TB

        IngestionLayer["Data Ingestion Layer (Kafka Consumers, API Endpoints for Batch Push)"]
        KafkaTopics --> IngestionLayer

        subgraph "Data Processing & Storage"
            direction TB
            RawDataLake["Data Lake (S3: сырые и очищенные события - Parquet/Avro)"]
            IngestionLayer -- Сырые события --> RawDataLake

            StreamProcessing["Stream Processing (Flink/Kafka Streams: обогащение, real-time агрегаты, feature engineering)"]
            IngestionLayer -- Поток событий --> StreamProcessing

            RealtimeDWH["Real-time/Speed Layer (ClickHouse/Redis: витрины реального времени, быстрые счетчики)"]
            StreamProcessing -- Real-time витрины --> RealtimeDWH

            BatchProcessing["Batch Processing (Spark ETL/ELT: сложные агрегаты, исторические витрины, feature engineering)"]
            RawDataLake -- Очищенные данные --> BatchProcessing
            StreamProcessing -- Обработанные события (опционально) --> BatchProcessingViaKafka["Kafka (для триггеров батчей)"] -- опционально --> BatchProcessing

            DWH["Core DWH (ClickHouse: исторические агрегаты, детальные витрины данных)"]
            BatchProcessing -- Агрегированные витрины --> DWH

            MLDataStore["ML Feature Store & Model Artifacts (S3/ClickHouse/PostgreSQL - [Конкретная реализация ML Feature Store, например, Feast/Hopsworks или кастомное решение, подлежит уточнению])"]
            BatchProcessing -- Признаки для ML --> MLDataStore
            StreamProcessing -- Real-time признаки для ML --> MLDataStore
        end

        subgraph "API & Management Layer (Go / Java - Clean Architecture)"
            direction TB
            MetadataDB["Metadata Store (PostgreSQL: определения метрик, отчетов, ML-моделей, пайплайнов, прав доступа к данным)"]

            APIService["API Service (REST/GraphQL - Presentation Layer)"]
            APIServiceLogic["Application & Domain Logic for API (Use Cases, Entities)"]
            APIService --> APIServiceLogic
            APIServiceLogic -- CRUD метаданных --> MetadataDB
            APIServiceLogic -- Запросы данных --> DWH
            APIServiceLogic -- Запросы real-time данных --> RealtimeDWH
            APIServiceLogic -- Запросы прогнозов --> MLServingAPI["ML Model Serving API"]

            PipelineOrchestrator["Data Pipeline Orchestrator (Apache Airflow, Prefect, или Argo Workflows)"]
            PipelineOrchestrator -- Управление --> BatchProcessing
            PipelineOrchestrator -- Управление --> StreamProcessing
            PipelineOrchestrator -- Чтение/Запись метаданных --> MetadataDB
        end

        subgraph "ML Layer (MLOps)"
            direction TB
            MLTraining["ML Model Training & Experimentation (Jupyter, Python libs, SparkML)"]
            MLDataStore -- Данные для тренировки --> MLTraining
            MLTraining --> MLModelRegistry["ML Model Registry (MLflow / DVC + PostgreSQL/S3)"]
            MetadataDB <--> MLModelRegistry

            MLServingAPI -- Загрузка моделей --> MLModelRegistry
            MLDataStore -- Артефакты моделей --> MLServingAPI
        end

        APIService --> ExternalConsumers["Потребители API (Admin Panel, Developer Portal, BI Tools, Маркетинговые системы)"]
        MLServingAPI --> ExternalConsumersML["Потребители ML API (Recommendation Svc, Fraud Detection Svc, Personalization Svc)"]
    end

    classDef dataFlow fill:#e6f3ff,stroke:#007bff,stroke-width:2px;
    classDef component fill:#d4edda,stroke:#28a745,stroke-width:2px;
    classDef datastore fill:#f8d7da,stroke:#dc3545,stroke-width:2px;
    classDef api fill:#fff3cd,stroke:#ffc107,stroke-width:2px;
    classDef external fill:#e2e3e5,stroke:#6c757d,stroke-width:2px;

    class KafkaTopics,BatchProcessingViaKafka dataFlow;
    class IngestionLayer,StreamProcessing,BatchProcessing,MLTraining,MLModelRegistry,PipelineOrchestrator,APIServiceLogic component;
    class RawDataLake,DWH,RealtimeDWH,MetadataDB,MLDataStore datastore;
    class APIService,MLServingAPI api;
    class ExternalConsumers,ExternalConsumersML external;
```

### 2.3. Слои и Компоненты (детальнее)

#### 2.3.1. Data Ingestion Layer (Слой Приема Данных)
*   Ответственность: Подписка на топики Kafka от всех микросервисов платформы. Сериализация/десериализация событий. Базовая валидация схем. Сохранение сырых данных в Data Lake (S3) для долговременного хранения и пакетной обработки. Передача данных в слой потоковой обработки.
*   Ключевые технологии: Kafka Consumers (Java/Scala/Python), коннекторы S3, Avro/Protobuf для схем событий.

#### 2.3.2. Data Processing Layer (Слой Обработки Данных)
*   Ответственность: Очистка, трансформация, обогащение, агрегация данных.
*   **Stream Processing (Потоковая обработка):**
    *   Технологии: Apache Flink или Kafka Streams.
    *   Задачи: Расчет метрик в реальном времени (например, количество активных пользователей), обогащение событий (например, геолокация), обнаружение простых паттернов/аномалий, обновление витрин данных реального времени в ClickHouse.
*   **Batch Processing (Пакетная обработка):**
    *   Технологии: Apache Spark.
    *   Задачи: Сложные ETL/ELT процессы, пересчет исторических данных, построение агрегированных витрин в DWH (ClickHouse), подготовка данных для обучения ML моделей (feature engineering). Запускается по расписанию или триггерам.

#### 2.3.3. Data Storage Layer (Слой Хранения Данных)
*   Ответственность: Обеспечение надежного и эффективного хранения данных на разных уровнях обработки.
*   **Data Lake (S3-совместимое хранилище):** Хранение всех сырых событий в их исходном или минимально обработанном виде (например, в формате Parquet или ORC). Структура обычно партиционирована по дате и типу события (например, `s3://bucket/raw_events/event_type=user_login/year=2024/month=03/day=15/`).
*   **Data Warehouse (DWH - ClickHouse):** Хранение структурированных, агрегированных данных и витрин данных (data marts), оптимизированных для быстрых аналитических запросов.
*   **Real-time Data Marts (ClickHouse/Redis):** Хранение часто обновляемых метрик или данных для дашбордов реального времени.
*   **Metadata Store (PostgreSQL):** Хранение метаданных: определения метрик, отчетов, сегментов, конфигурации пайплайнов, схемы данных, метаданные ML моделей (если не используется специализированный реестр типа MLflow эксклюзивно).

#### 2.3.4. API & Management Layer (Слой API и Управления) - Принципы Clean Architecture
Для компонентов этого слоя (API Service, управление метаданными) применяется стандартная слоистая архитектура:
*   **Presentation Layer (API Service):**
    *   Ответственность: Обработка HTTP REST (и потенциально GraphQL/WebSocket) запросов от внешних потребителей (админ-панель, BI-инструменты, другие сервисы). Валидация запросов, аутентификация, авторизация.
    *   Технологии: Go (Echo/Gin) или Java (Spring Boot).
*   **Application Layer:**
    *   Ответственность: Координация выполнения запросов, вызов соответствующих сервисов домена или инфраструктуры для получения/обработки данных. Формирование ответов.
    *   Компоненты: `MetricService`, `ReportService`, `SegmentService`.
*   **Domain Layer:**
    *   Ответственность: Бизнес-логика управления метаданными (метрики, отчеты, сегменты), правила валидации.
    *   Компоненты: Сущности `MetricDefinition`, `ReportDefinition`, `SegmentDefinition`.
*   **Infrastructure Layer:**
    *   Ответственность: Взаимодействие с Metadata Store (PostgreSQL) для CRUD операций с метаданными. Формирование и выполнение запросов к DWH (ClickHouse) и Real-time Data Marts для получения данных. Взаимодействие с ML Model Serving API.

#### 2.3.5. Machine Learning Layer (Слой Машинного Обучения)
*   Ответственность: Полный жизненный цикл ML моделей: от подготовки данных и тренировки до развертывания и мониторинга.
*   **ML Model Training:**
    *   Технологии: Python, SparkML, TensorFlow, PyTorch, Scikit-learn.
    *   Задачи: Feature engineering, обучение моделей, оценка качества, версионирование.
*   **ML Model Registry (MLflow или аналог):**
    *   Задачи: Хранение артефактов моделей, версий, метрик производительности, параметров обучения.
*   **ML Model Serving:**
    *   Технологии: Python (Flask/FastAPI), Java (Spring Boot), или специализированные решения (KFServing, Seldon Core).
    *   Задачи: Предоставление REST/gRPC API для получения прогнозов от развернутых моделей.

#### 2.3.6. Data Pipeline Orchestration
*   Ответственность: Управление и мониторинг выполнения пайплайнов обработки данных (batch и stream).
*   Технологии: Apache Airflow, Prefect, или кастомные решения на базе Kubernetes CronJobs/Argo Workflows.

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL (через API Gateway):** `/api/v1/analytics`
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway). Роли из токена используются для авторизации. (см. `../../../../project_security_standards.md`)
*   **Авторизация:** На основе ролей (см. `../../../../project_roles_and_permissions.md`). Например:
    *   `platform_admin`: доступ ко всем данным.
    *   `developer`: доступ к агрегированной статистике по своим играм.
    *   `marketing_manager`: доступ к отчетам по кампаниям, сегментам.
*   **Формат ответа об ошибке (согласно `../../../../project_api_standards.md`):**
    ```json
    {
      "errors": [
        {
          "code": "ERROR_CODE_STRING",
          "title": "Человекочитаемый заголовок ошибки",
          "detail": "Детальное описание проблемы.",
          "source": { "pointer": "/data/attributes/field_name", "parameter": "query_param_name" } // Опционально
        }
      ]
    }
    ```
*   **Примечание по схемам запросов/ответов:** Детальные JSON схемы для всех API запросов и ответов будут доступны через публикуемую OpenAPI спецификацию сервиса.

#### 3.1.1. Ресурс: Метрики (Metrics)
*   Эндпоинты `GET /metrics/definitions`, `GET /metrics/definitions/{metric_name}`, `GET /metrics/values/{metric_name}` (как в существующем документе, с уточнением прав доступа).
    *   Пример ответа (Ошибка 404 Not Found для `GET /metrics/values/{metric_name}`):
        ```json
        // Example for GET /metrics/values/{metric_name} - Error 404
        {
          "errors": [
            {
              "code": "METRIC_NOT_FOUND",
              "title": "Метрика не найдена",
              "detail": "Метрика с именем 'non_existent_metric' не найдена или для нее нет данных."
            }
          ]
        }
        ```

#### 3.1.2. Ресурс: Отчеты (Reports)
*   Эндпоинты `GET /reports/definitions`, `POST /reports/instances`, `GET /reports/instances/{instance_id}`, `GET /reports/instances/{instance_id}/download` (как в существующем документе, с уточнением прав доступа).

#### 3.1.3. Ресурс: Сегменты (Segments)
*    Эндпоинты `GET /segments/definitions`, `POST /segments/definitions`, `GET /segments/{segment_id}/users-count` (как в существующем документе, с уточнением прав доступа).

#### 3.1.4. Ресурс: Предиктивная Аналитика (Predictions)
*   Эндпоинты `GET /predictions/models`, `POST /predictions/{model_name_or_id}/predict` (как в существующем документе, с уточнением прав доступа).

### 3.2. GraphQL API
*   **Эндпоинт:** `/api/v1/analytics/graphql` (потенциальный)
*   **Статус:** Не реализован. Рассматривается для будущих версий для предоставления гибкого доступа к данным для продвинутых пользователей или BI-инструментов.

### 3.3. WebSocket API
*   **Эндпоинт:** `/api/v1/analytics/ws/streaming` (потенциальный)
*   **Статус:** Не реализован. Рассматривается для будущих версий для стриминга real-time метрик на дашборды.

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности/Структуры Данных
*   **`Event` (Событие):** Базовая структура события, соответствующая CloudEvents. Содержит поля: `id`, `type`, `source`, `time`, `dataContentType`, `data`, `subject`, `correlationId`. Детали поля `data` зависят от конкретного типа события. Хранится в S3 (сырые, очищенные), используется в Kafka, обрабатывается в Spark/Flink, агрегированные данные попадают в ClickHouse.
*   **`MetricDefinition` (Определение Метрики):** Описание метрики, как она рассчитывается, ее источники. Хранится в PostgreSQL.
    *   Поля: `id`, `name`, `display_name_ru`, `display_name_en`, `description_ru`, `description_en`, `metric_type` (counter, gauge, histogram), `calculation_method` (sum, avg, count_distinct), `source_event_type`, `value_field`, `filters` (JSONB), `dimensions` (JSONB), `granularity` (realtime, hourly, daily), `unit`, `is_realtime`, `owner_service`, `tags`, `created_at`, `updated_at`.
*   **`ReportDefinition` (Определение Отчета):** Шаблон отчета, включая запросы к данным, параметры. Хранится в PostgreSQL.
    *   Поля: `id`, `name`, `display_name_ru`, `display_name_en`, `description_ru`, `description_en`, `source_type` (sql_query_clickhouse, pre_aggregated_metrics), `source_query_or_config` (TEXT), `parameters` (JSONB), `default_schedule`, `output_formats` (JSONB), `owner_admin_id`, `created_at`, `updated_at`.
*   **`ReportInstance` (Экземпляр Отчета):** Сгенерированный отчет. Метаданные хранятся в PostgreSQL, сам файл отчета - в S3.
    *   Поля: `id`, `report_definition_id`, `generation_requested_at`, `generation_started_at`, `generation_completed_at`, `status`, `parameters_used` (JSONB), `output_format`, `file_path_s3`, `file_size_bytes`, `error_message`, `requested_by_user_id`, `expires_at`.
*   **`SegmentDefinition` (Определение Сегмента):** Критерии для формирования пользовательского сегмента. Хранится в PostgreSQL.
    *   Поля: `id`, `name`, `display_name_ru`, `display_name_en`, `description_ru`, `description_en`, `criteria` (JSONB), `segment_type` (dynamic, static), `refresh_schedule`, `last_calculated_at`, `current_user_count`, `created_by_admin_id`, `created_at`, `updated_at`.
*   **`MLModelMetadata` (Метаданные ML Модели):** Информация о ML модели. Хранится в PostgreSQL или специализированном реестре (MLflow).
    *   Поля: `id`, `model_name`, `model_version`, `description_ru`, `description_en`, `algorithm`, `hyperparameters` (JSONB), `training_dataset_ref`, `performance_metrics` (JSONB), `artifact_path`, `deployment_status`, `input_schema` (JSONB), `output_schema` (JSONB), `registered_by_user_id`, `trained_at`, `created_at`, `updated_at`.
*   **`DataPipelineRun` (Запуск Пайплайна Данных):** Информация о запусках ETL/ELT пайплайнов. Хранится в PostgreSQL.
    *   Поля: `id`, `pipeline_name`, `start_time`, `end_time`, `status`, `parameters` (JSONB), `logs_summary`, `processed_records_count`, `source_data_start_time`, `source_data_end_time`.

### 4.2. Схема Базы Данных (Концептуальные диаграммы и DDL)

#### 4.2.1. PostgreSQL (Metadata Store)
*   ERD Диаграмма (PostgreSQL):
    ```mermaid
    erDiagram
        METRIC_DEFINITION {
            UUID id PK
            VARCHAR name UK
            VARCHAR display_name_ru
            TEXT description_ru
            VARCHAR metric_type
            VARCHAR calculation_method
            VARCHAR source_event_type
            JSONB filters
            JSONB dimensions
            VARCHAR granularity
            BOOLEAN is_realtime
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }

        REPORT_DEFINITION {
            UUID id PK
            VARCHAR name UK
            VARCHAR display_name_ru
            TEXT description_ru
            VARCHAR source_type
            TEXT source_query_or_config
            JSONB parameters
            VARCHAR default_schedule
            JSONB output_formats
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }
        REPORT_DEFINITION ||--o{ REPORT_INSTANCE : "has many"

        REPORT_INSTANCE {
            UUID id PK
            UUID report_definition_id FK
            TIMESTAMPTZ generation_requested_at
            VARCHAR status
            JSONB parameters_used
            VARCHAR file_path_s3
            TIMESTAMPTZ expires_at
        }

        SEGMENT_DEFINITION {
            UUID id PK
            VARCHAR name UK
            VARCHAR display_name_ru
            JSONB criteria
            VARCHAR segment_type
            VARCHAR refresh_schedule
            BIGINT current_user_count
            TIMESTAMPTZ last_calculated_at
            TIMESTAMPTZ created_at
        }

        ML_MODEL_METADATA {
            UUID id PK
            VARCHAR model_name
            VARCHAR model_version
            TEXT description_ru
            JSONB hyperparameters
            JSONB performance_metrics
            VARCHAR artifact_path
            VARCHAR deployment_status
            TIMESTAMPTZ trained_at
            UNIQUE (model_name, model_version)
        }

        DATA_PIPELINE_RUN {
            UUID id PK
            VARCHAR pipeline_name
            TIMESTAMPTZ start_time
            TIMESTAMPTZ end_time
            VARCHAR status
            JSONB parameters
            BIGINT processed_records_count
        }
    ```
*   DDL (PostgreSQL - основные таблицы, как в существующем документе, с дополнениями):
```sql
-- Расширение для UUID, если не создано
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Определения метрик
CREATE TABLE metric_definitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    display_name_ru VARCHAR(255) NOT NULL,
    display_name_en VARCHAR(255),
    description_ru TEXT,
    description_en TEXT,
    metric_type VARCHAR(50) NOT NULL CHECK (metric_type IN ('counter', 'gauge', 'histogram', 'timer')),
    calculation_method VARCHAR(50) NOT NULL, -- sum, average, count_distinct, percentile_95, etc.
    source_event_type VARCHAR(255), -- Например, com.platform.payment.transaction.completed.v1
    value_field VARCHAR(255), -- Путь к значению в JSON событии, e.g., data.amount_rub
    filters JSONB, -- Фильтры для событий, e.g., {"data.status": "success"}
    dimensions JSONB, -- Измерения для группировки, e.g., ["data.game_id", "data.payment_method"]
    granularity VARCHAR(50) NOT NULL CHECK (granularity IN ('realtime', 'hourly', 'daily', 'monthly', 'yearly', 'raw')),
    unit VARCHAR(50), -- Единица измерения: users, seconds, RUB, events
    is_realtime BOOLEAN NOT NULL DEFAULT FALSE,
    owner_service VARCHAR(100), -- Сервис-владелец или основной потребитель метрики
    tags JSONB DEFAULT '[]'::jsonb, -- Теги для категоризации и поиска (например, "sales", "user_activity")
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE metric_definitions IS 'Определения метрик, используемых в системе аналитики.';
CREATE INDEX idx_metric_definitions_name ON metric_definitions(name);
CREATE INDEX idx_metric_definitions_tags ON metric_definitions USING gin(tags);

-- Определения отчетов
CREATE TABLE report_definitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE, -- Уникальное системное имя отчета
    display_name_ru VARCHAR(255) NOT NULL,
    display_name_en VARCHAR(255),
    description_ru TEXT,
    description_en TEXT,
    source_type VARCHAR(50) NOT NULL CHECK (source_type IN ('sql_query_clickhouse', 'pre_aggregated_metrics', 'api_external', 'python_script_placeholder')), -- Тип источника данных для отчета
    source_query_or_config TEXT NOT NULL, -- SQL запрос для ClickHouse, конфигурация для других источников
    parameters JSONB, -- Описание параметров отчета, e.g., [{"name": "period_start", "display_name_ru": "Начало периода", "type": "date", "required": true, "default_value": "yesterday"}]
    default_schedule VARCHAR(50), -- Расписание генерации в формате cron, e.g., "0 3 * * *" (ежедневно в 3 AM UTC)
    output_formats JSONB DEFAULT '["csv", "json", "pdf_placeholder"]'::jsonb, -- Доступные форматы выгрузки
    owner_admin_id UUID, -- REFERENCES admin_users(id) ON DELETE SET NULL, -- Ссылка на администратора-владельца
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE report_definitions IS 'Определения отчетов, которые могут быть сгенерированы.';
CREATE INDEX idx_report_definitions_name ON report_definitions(name);

-- Экземпляры сгенерированных отчетов
CREATE TABLE report_instances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_definition_id UUID NOT NULL REFERENCES report_definitions(id) ON DELETE CASCADE,
    generation_requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    generation_started_at TIMESTAMPTZ,
    generation_completed_at TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL CHECK (status IN ('requested', 'generating', 'completed', 'failed', 'cancelled')),
    parameters_used JSONB, -- Использованные параметры при генерации
    output_format VARCHAR(20),
    file_path_s3 VARCHAR(1024), -- Ссылка на файл отчета в S3
    file_size_bytes BIGINT,
    error_message TEXT, -- Сообщение об ошибке, если генерация не удалась
    requested_by_user_id UUID, -- REFERENCES users(id) ON DELETE SET NULL, -- Пользователь или админ, запросивший отчет
    expires_at TIMESTAMPTZ -- Время, когда файл отчета может быть удален из S3 (политика очистки)
);
COMMENT ON TABLE report_instances IS 'Экземпляры сгенерированных отчетов.';
CREATE INDEX idx_report_instances_status ON report_instances(status);
CREATE INDEX idx_report_instances_definition_id ON report_instances(report_definition_id);
CREATE INDEX idx_report_instances_requested_by ON report_instances(requested_by_user_id);

-- Определения пользовательских сегментов
CREATE TABLE segment_definitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE, -- Системное имя сегмента
    display_name_ru VARCHAR(255) NOT NULL,
    display_name_en VARCHAR(255),
    description_ru TEXT,
    description_en TEXT,
    criteria JSONB NOT NULL, -- Описание критериев сегментации, e.g., {"type": "AND", "conditions": [{"field": "user.total_payments_sum_rub", "operator": ">=", "value": 10000}, {"field": "user.registration_date", "operator": "<", "value": "2023-01-01"}]}
    segment_type VARCHAR(50) NOT NULL DEFAULT 'dynamic' CHECK (segment_type IN ('dynamic', 'static')), -- dynamic (пересчитываемый), static (фиксированный список user_id)
    refresh_schedule VARCHAR(50), -- Для dynamic сегментов, расписание пересчета в формате cron
    last_calculated_at TIMESTAMPTZ, -- Время последнего пересчета
    current_user_count BIGINT, -- Количество пользователей в сегменте после последнего пересчета
    created_by_admin_id UUID, -- REFERENCES admin_users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE segment_definitions IS 'Определения пользовательских сегментов.';
CREATE INDEX idx_segment_definitions_name ON segment_definitions(name);

-- Метаданные ML моделей
CREATE TABLE ml_model_metadata (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    model_name VARCHAR(255) NOT NULL, -- Имя модели, e.g., "game_recommendation_als"
    model_version VARCHAR(50) NOT NULL, -- Версия модели, e.g., "v1.2.3"
    description_ru TEXT,
    description_en TEXT,
    algorithm VARCHAR(100), -- Алгоритм, e.g., "ALS", "LightGBM"
    hyperparameters JSONB, -- Гиперпараметры, использованные при тренировке
    training_dataset_ref VARCHAR(1024), -- Ссылка на датасет для тренировки (e.g., S3 path, DVC reference)
    performance_metrics JSONB, -- Метрики качества модели, e.g., {"auc": 0.85, "f1_score": 0.78, "training_duration_sec": 3600}
    artifact_path VARCHAR(1024), -- Путь к артефактам модели (файл модели, веса) в MLflow или S3
    deployment_status VARCHAR(50) NOT NULL DEFAULT 'development' CHECK (deployment_status IN ('development', 'staging', 'production', 'archived', 'failed')),
    input_schema JSONB, -- Описание ожидаемой структуры входных данных для модели
    output_schema JSONB, -- Описание структуры выходных данных (прогнозов)
    registered_by_user_id UUID, -- REFERENCES users(id) ON DELETE SET NULL, -- Пользователь, зарегистрировавший модель
    trained_at TIMESTAMPTZ, -- Время завершения тренировки
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (model_name, model_version)
);
COMMENT ON TABLE ml_model_metadata IS 'Метаданные моделей машинного обучения.';
CREATE INDEX idx_ml_model_metadata_name_version ON ml_model_metadata(model_name, model_version);

-- Запуски пайплайнов обработки данных
CREATE TABLE data_pipeline_runs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pipeline_name VARCHAR(255) NOT NULL, -- Имя пайплайна (e.g., "daily_user_activity_etl", "ml_recommendation_feature_engineering")
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL CHECK (status IN ('running', 'completed', 'failed', 'skipped', 'manual_run')),
    parameters JSONB, -- Параметры запуска пайплайна
    logs_summary TEXT, -- Краткое содержание логов или ссылка на полные логи (S3, ELK/Loki)
    processed_records_count BIGINT, -- Количество обработанных записей
    source_data_start_time TIMESTAMPTZ, -- Для инкрементальных загрузок: начало временного окна исходных данных
    source_data_end_time TIMESTAMPTZ   -- Для инкрементальных загрузок: конец временного окна исходных данных
);
COMMENT ON TABLE data_pipeline_runs IS 'Информация о запусках пайплайнов обработки данных.';
CREATE INDEX idx_data_pipeline_runs_name_status_start_time ON data_pipeline_runs(pipeline_name, status, start_time DESC);

-- [Спецификации таблиц для конфигураций A/B тестов будут добавлены здесь, если эта функциональность будет реализована в Analytics Service.]
-- [Спецификации таблиц для определений алертов будут добавлены здесь, если потребуется их хранение в Analytics Service, а не только в Alertmanager.]
```
#### 4.2.2. ClickHouse (DWH) - DDL Примеры
*   **Таблица фактов: Игровые сессии (пример)**
    ```sql
    CREATE TABLE analytics_db.fact_game_sessions (
        session_id String DEFAULT generateUUIDv4(),  -- Уникальный ID сессии
        user_id String,                             -- ID пользователя (UUID в строковом представлении)
        game_id String,                             -- ID игры (UUID в строковом представлении)
        start_timestamp DateTime64(3, 'UTC'),       -- Время начала сессии
        end_timestamp DateTime64(3, 'UTC'),         -- Время окончания сессии
        duration_seconds UInt32,                    -- Длительность сессии в секундах
        platform Enum8('pc' = 1, 'mobile_android' = 2, 'mobile_ios' = 3, 'web' = 4, 'unknown' = 0),
        country_code FixedString(2),                -- Код страны пользователя (ISO 3166-1 alpha-2)
        client_version String,                      -- Версия игрового клиента
        ip_address IPv6,                            -- IP адрес (для геолокации, если разрешено)

        -- Поля для агрегации и анализа, могут добавляться по мере необходимости
        first_session_for_user_game Bool,           -- Первая ли это сессия для данной пары пользователь-игра
        total_ingame_purchases_count UInt16 DEFAULT 0,
        total_ingame_purchases_amount Decimal64(2) DEFAULT 0.00, -- Сумма внутриигровых покупок за сессию

        event_date Date ALIAS toDate(start_timestamp), -- Для партиционирования
        processed_at DateTime DEFAULT now()            -- Время обработки записи в аналитике
    )
    ENGINE = MergeTree()
    PARTITION BY toYYYYMM(event_date)
    ORDER BY (game_id, event_date, user_id, start_timestamp)
    SETTINGS index_granularity = 8192, allow_nullable_key = 1; -- allow_nullable_key для Nullable полей в ORDER BY, если есть
    ```
    *Комментарий: Эта таблица агрегирует данные по игровым сессиям. Источником могут быть события начала и конца сессии, либо события активности внутри сессии.*

*   **Витрина данных: Ежедневная активность пользователей по играм (пример)**
    ```sql
    CREATE MATERIALIZED VIEW analytics_db.mart_daily_game_activity_mv
    ENGINE = SummingMergeTree() -- или AggregatingMergeTree для более сложных агрегаций
    PARTITION BY toYYYYMM(activity_date)
    ORDER BY (activity_date, game_id, country_code, platform)
    POPULATE -- Заполнить данными при создании
    AS SELECT
        toDate(start_timestamp) AS activity_date,
        game_id,
        country_code,
        platform,
        count(DISTINCT user_id) AS dau_count, -- Количество уникальных активных пользователей
        sum(duration_seconds) AS total_playtime_seconds,
        avg(duration_seconds) AS avg_playtime_seconds,
        sum(total_ingame_purchases_amount) AS total_revenue_ingame
        -- Для AggregatingMergeTree можно использовать:
        -- uniqExactState(user_id) AS unique_users_state -- для последующего использования с uniqExactMerge
    FROM analytics_db.fact_game_sessions
    GROUP BY
        activity_date,
        game_id,
        country_code,
        platform;
    ```
    *Комментарий: Эта материализованная витрина автоматически агрегирует данные из `fact_game_sessions` для быстрого получения DAU, общего и среднего времени игры, и дохода по играм, странам и платформам.*

#### 4.2.3. S3 Data Lake Структура и Схемы
*   Сырые и обработанные события, а также артефакты ML и экспорты отчетов хранятся в S3-совместимом хранилище.
*   **Диаграмма структуры S3 (концептуальная):**
    ```mermaid
    graph LR
        S3Bucket["S3 Bucket (Data Lake)"]
        S3Bucket --> RawEvents["/raw_events/"]
        RawEvents --> ServiceName["service_name=&lt;service&gt;/"]
        ServiceName --> EventType["event_type=&lt;event_type_fqdn&gt;/"]
        EventType --> Year["year=&lt;YYYY&gt;/"]
        Year --> Month["month=&lt;MM&gt;/"]
        Month --> Day["day=&lt;DD&gt;/"]
        Day --> Hour["hour=&lt;HH&gt;/"]
        Hour --> EventFile["&lt;uuid&gt;.parquet (или .json.gz)"]

        S3Bucket --> ProcessedEvents["/processed_events/"]
        ProcessedEvents --> ProcessedEventType["event_type=&lt;processed_event_type&gt;/"]
        ProcessedEventType --> ProcessedYear["year=&lt;YYYY&gt;/"]
        ProcessedYear --> ProcessedMonth["month=&lt;MM&gt;/"]
        ProcessedMonth --> ProcessedFile["&lt;uuid_or_batch_id&gt;.parquet"]

        S3Bucket --> MLEvents["/ml_artifacts/"]
        MLEvents --> ModelName["model_name=&lt;model_name&gt;/"]
        ModelName --> Version["version=&lt;version&gt;/"]
        Version --> ArtifactFiles["&lt;artifact_files&gt; (модели, словари)"]

        S3Bucket --> Reports["/reports/"]
        Reports --> ReportInstance["report_instance_id=&lt;instance_uuid&gt;/"]
        ReportInstance --> ReportFile["&lt;filename&gt;.&lt;format&gt; (csv, json)"]

        style S3Bucket fill:#f9f,stroke:#333,stroke-width:2px
        style RawEvents fill:#ccf,stroke:#333,stroke-width:2px
        style ProcessedEvents fill:#cfc,stroke:#333,stroke-width:2px
        style MLEvents fill:#fec,stroke:#333,stroke-width:2px
        style Reports fill:#ffc,stroke:#333,stroke-width:2px
    ```
*   **Пример схемы события в Parquet (для `com.platform.library.session.started.v1`):**
    ```
    message GameSessionStartedEvent {
      required binary id (STRING);             // ID события (UUID)
      required binary type (STRING);           // Тип события (com.platform.library.session.started.v1)
      required binary source (STRING);         // Источник (library-service)
      required int64 time (TIMESTAMP_MICROS);  // Время события UTC
      required binary subject (STRING);        // ID сессии (UUID)
      optional binary correlationId (STRING);  // ID корреляции (UUID)
      message Data {
        required binary userId (STRING);       // ID пользователя (UUID)
        required binary gameId (STRING);       // ID игры (UUID)
        required int64 sessionStartTime (TIMESTAMP_MICROS); // Время начала сессии UTC
        required binary platform (ENUM);     // pc, mobile_android, mobile_ios, web, unknown
        optional binary clientVersion (STRING);
        optional binary countryCode (STRING);  // ISO 3166-1 alpha-2
        // Другие релевантные поля...
      }
      required Data data;
    }
    ```
    *Комментарий: Использование Parquet или Avro предпочтительно для типизированного, сжатого и колоночного хранения, что ускоряет обработку в Spark/Flink и запросы из DWH (например, через внешние таблицы ClickHouse).*

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   Analytics Service в основном является потребителем событий. Однако он может публиковать:
*   **Формат событий:** CloudEvents JSON (согласно `../../../../project_api_standards.md`).
*   **Основные топики Kafka для публикуемых событий:** `com.platform.analytics.events.v1`.

*   **`com.platform.analytics.report.generated.v1`**
    *   Описание: Отчет сгенерирован и доступен для скачивания или просмотра.
    *   `data` Payload:
        ```json
        {
          "reportInstanceId": "instance-uuid-xyz",
          "reportDefinitionName": "monthly_sales_summary",
          "status": "completed", // "completed", "failed"
          "downloadUrl": "s3://bucket/reports/instance-uuid-xyz/monthly_sales.csv", // Опционально, если доступ прямой
          "generationCompletedAt": "ISO8601_timestamp",
          "requestedByUserId": "user-uuid-admin" // Опционально
        }
        ```
    *   Потребители: Notification Service (для уведомления пользователя/администратора), Admin Service (для отображения статуса).
*   **`com.platform.analytics.segment.updated.v1`**
    *   Описание: Пользовательский сегмент был обновлен (например, пересчитан состав участников).
    *   `data` Payload:
        ```json
        {
          "segmentId": "segment-uuid-rpg-active",
          "segmentName": "active_rpg_players_last_month",
          "userCount": 15230,
          "calculationTimestamp": "ISO8601_timestamp"
        }
        ```
    *   Потребители: Marketing Service, Notification Service (для целевых кампаний), Personalization Service.
*   **`com.platform.analytics.alert.triggered.v1`**
    *   Описание: Сработал аналитический алерт, требующий внимания (например, резкое падение DAU, рост ошибок транзакций).
    *   `data` Payload:
        ```json
        {
          "alertName": "CriticalDAUDrop",
          "severity": "critical", // "critical", "warning", "info"
          "description": "DAU упал на 30% по сравнению со средним за последние 7 дней.",
          "metricName": "daily_active_users",
          "currentValue": 7000,
          "thresholdValue": 10000,
          "dimensions": {"platform": "pc"}, // Опционально, для уточнения
          "triggerTimestamp": "ISO8601_timestamp"
        }
        ```
    *   Потребители: Admin Service (для отображения на дашборде инцидентов), Notification Service (для оповещения администраторов/ответственных лиц).

### 5.2. Потребляемые События (Consumed Events)
*   Analytics Service является основным потребителем событий от **всех** других микросервисов платформы. Это включает, но не ограничивается:
    *   `com.platform.auth.user.registered.v1`, `com.platform.auth.user.loggedin.v1`
    *   `com.platform.account.profile.updated.v1`, `com.platform.account.status.updated.v1`
    *   `com.platform.catalog.game.viewed.v1`, `com.platform.catalog.game.searched.v1`
    *   `com.platform.library.game.added.v1`, `com.platform.library.session.started.v1`, `com.platform.library.session.ended.v1`
    *   `com.platform.payment.transaction.completed.v1`, `com.platform.payment.refund.processed.v1`
    *   `com.platform.download.started.v1`, `com.platform.download.completed.v1`
    *   `com.platform.social.review.created.v1`, `com.platform.social.friend.added.v1`
    *   `com.platform.notification.sent.v1`, `com.platform.notification.failed.v1`
    *   `com.platform.developer.game.published.v1`
    *   `com.platform.admin.user.banned.v1`
*   **Топики Kafka:** Подписка на все релевантные топики событий других сервисов.
*   **Формат событий:** Ожидается CloudEvents JSON (согласно `../../../../project_api_standards.md`).
*   **Логика обработки:** Описана в разделе "Внутренняя Архитектура".

## 6. Интеграции (Integrations)
См. `../../../../project_integrations.md` для общей карты и деталей.

### 6.1. Внутренние Микросервисы
*   **Потребление данных:** Analytics Service интегрируется со всеми микросервисами платформы через Kafka для сбора событий. Это основной механизм получения данных.
*   **Предоставление данных (через API Analytics Service):**
    *   **Admin Service:** Для отображения дашбордов, отчетов, статистики платформы.
    *   **Developer Service:** Для предоставления разработчикам статистики по их играм (продажи, активность игроков, отзывы).
    *   **Marketing Service (гипотетический):** Для получения сегментов пользователей, результатов A/B тестов, анализа эффективности кампаний.
    *   **Recommendation Service (гипотетический):** Может использовать данные о поведении пользователей, популярности игр, ML-прогнозы от Analytics Service для генерации рекомендаций.
    *   **Notification Service:** Получает события о готовности отчетов или срабатывании алертов для уведомления соответствующих пользователей.
*   **Auth Service:** Для аутентификации и авторизации запросов к API Analytics Service (проверка JWT токенов).

### 6.2. Внешние Системы
*   **S3-совместимое хранилище:** Критически важная интеграция для Data Lake и хранения артефактов.
*   **Внешние BI-инструменты (например, Apache Superset, Metabase, Tableau):** Могут подключаться к DWH (ClickHouse) Analytics Service для построения кастомных отчетов и визуализаций данных.
*   **Системы визуализации (Grafana):** Для операционных метрик и некоторых аналитических дашбордов.

## 7. Конфигурация (Configuration)
Общие стандарты конфигурационных файлов (формат YAML, структура, управление переменными окружения и секретами) определены в `../../../../project_api_standards.md` (раздел 7) и `../../../../DOCUMENTATION_GUIDELINES.md` (раздел 6).

### 7.1. Переменные Окружения (Примеры)
*   `ANALYTICS_API_HTTP_PORT`: Порт для REST API Analytics Service (например, `8090`).
*   `POSTGRES_DSN_ANALYTICS_META`: DSN для подключения к PostgreSQL (Metadata Store).
*   `CLICKHOUSE_HOSTS_ANALYTICS_DWH`: Список хостов ClickHouse (например, `ch1:8123,ch2:8123`).
*   `CLICKHOUSE_USER_ANALYTICS_DWH`, `CLICKHOUSE_PASSWORD_ANALYTICS_DWH`.
*   `S3_ENDPOINT_DATALAKE`, `S3_ACCESS_KEY_DATALAKE`, `S3_SECRET_KEY_DATALAKE`, `S3_BUCKET_DATALAKE`.
*   `KAFKA_BROKERS_ANALYTICS`: Список брокеров Kafka для Analytics Service.
*   `SPARK_MASTER_URL`: URL Spark Master (если используется).
*   `FLINK_MASTER_URL`: URL Flink JobManager (если используется).
*   `MLFLOW_TRACKING_URI`: URI для MLflow Tracking Server.
*   `LOG_LEVEL_ANALYTICS`: Уровень логирования для компонентов Analytics Service.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT_ANALYTICS`: Endpoint для Jaeger.

### 7.2. Файлы Конфигурации (`configs/analytics_service_config.yaml` - для API Layer)
*   Структура (для API Layer, компоненты обработки данных могут иметь свои конфигурационные файлы):
    ```yaml
    http_server:
      port: ${ANALYTICS_API_HTTP_PORT:"8090"}
      timeout_seconds: 60
    postgres_metadata:
      dsn: ${POSTGRES_DSN_ANALYTICS_META}
      pool_max_conns: 10
    clickhouse_dwh:
      hosts: ${CLICKHOUSE_HOSTS_ANALYTICS_DWH} # "host1:port,host2:port"
      database: "analytics_db"
      username: ${CLICKHOUSE_USER_ANALYTICS_DWH}
      password: ${CLICKHOUSE_PASSWORD_ANALYTICS_DWH}
      # Параметры подключения, таймауты
    s3_datalake:
      endpoint: ${S3_ENDPOINT_DATALAKE}
      access_key: ${S3_ACCESS_KEY_DATALAKE}
      secret_key: ${S3_SECRET_KEY_DATALAKE}
      bucket: ${S3_BUCKET_DATALAKE}
      region: "ru-central1" # Пример
    kafka:
      brokers: ${KAFKA_BROKERS_ANALYTICS}
      # Конфигурации Producer/Consumer для API Layer (если он публикует/читает напрямую)
      producer_topics:
        analytics_events: ${KAFKA_TOPIC_ANALYTICS_EVENTS:"com.platform.analytics.events.v1"}
      consumer_group_api: "analytics-api-group"
    logging:
      level: ${LOG_LEVEL_ANALYTICS:"info"}
      format: "json"
    security:
      jwt_public_key_path: ${JWT_PUBLIC_KEY_PATH} # Общий ключ для токенов платформы
    otel:
      exporter_jaeger_endpoint: ${OTEL_EXPORTER_JAEGER_ENDPOINT_ANALYTICS}
      service_name: "analytics-service-api"
    # Конфигурации для доступа к ML Serving API, если он внешний
    # ml_serving_api:
    #   base_url: ${ML_SERVING_API_URL}
    #   timeout_seconds: 5
    ```
*   **Конфигурация для Spark/Flink джобов:** Обычно передается через параметры командной строки при запуске джобы или через файлы конфигурации, специфичные для этих фреймворков (например, `spark-submit --properties-file <path>`).

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   **Data Pipelines:**
    *   Использование Dead Letter Queues (DLQ) в Kafka для событий, которые не удалось обработать.
    *   Механизмы retry с экспоненциальной задержкой для временных ошибок при обработке или записи данных.
    *   Мониторинг и алертинг по количеству ошибок в пайплайнах и DLQ.
    *   Подробное логирование ошибок на каждом этапе обработки.
*   **API Layer:**
    *   Стандартизированные коды и форматы ошибок для REST API (согласно `../../../../project_api_standards.md`).
    *   Информативные сообщения об ошибках для клиентов API.

### 8.2. Распространенные Коды Ошибок (для API)
*   **`METRIC_DEFINITION_NOT_FOUND`**: Запрошенное определение метрики не найдено.
*   **`REPORT_DEFINITION_NOT_FOUND`**: Запрошенное определение отчета не найдено.
*   **`REPORT_INSTANCE_NOT_FOUND`**: Экземпляр отчета не найден или еще не готов.
*   **`REPORT_GENERATION_FAILED`**: Ошибка при генерации отчета.
*   **`SEGMENT_DEFINITION_NOT_FOUND`**: Определение сегмента не найдено.
*   **`ML_MODEL_NOT_FOUND`**: ML модель не найдена или недоступна.
*   **`PREDICTION_FAILED`**: Ошибка при выполнении прогноза ML моделью.
*   **`INVALID_QUERY_PARAMETERS`**: Некорректные параметры запроса (например, неверный диапазон дат, неверный формат фильтра).
*   **`DATA_NOT_YET_AVAILABLE`**: Данные для запрошенного периода или метрики еще не рассчитаны или недоступны.
*   **`ACCESS_DENIED_TO_RESOURCE`**: У пользователя нет прав на доступ к запрошенным данным/метрикам.

## 9. Безопасность (Security)
См. `../../../../project_security_standards.md` для общих стандартов.

### 9.1. Аутентификация
*   Доступ к API Analytics Service защищен JWT токенами, валидируемыми через Auth Service или API Gateway.
*   Для внутренних компонентов (Spark/Flink джобы, доступ к Kafka, S3, ClickHouse, PostgreSQL) используются механизмы аутентификации, специфичные для этих систем (например, SASL для Kafka, IAM роли или ключи для S3, логины/пароли для БД).

### 9.2. Авторизация
*   Применяется RBAC (Role-Based Access Control) для доступа к API и данным. Роли определяются в `../../../../project_roles_and_permissions.md`.
*   Возможна реализация Attribute-Based Access Control (ABAC) для более гранулярного контроля доступа к данным на уровне строк или колонок в DWH, если это потребуется.

### 9.3. Защита Данных
*   **ФЗ-152 "О персональных данных":**
    *   Все персональные данные (ПДн), поступающие в Analytics Service, должны обрабатываться с соблюдением требований ФЗ-152.
    *   **Анонимизация/Псевдонимизация:** Перед сохранением в Data Lake и DWH для общего аналитического использования, ПДн должны проходить процедуру псевдонимизации (например, замена реальных UserID на псевдонимы) или анонимизации (если данные используются для публичных отчетов или исследований без возможности деанонимизации). Методы и степень анонимизации/псевдонимизации должны быть документированы.
    *   Доступ к сырым данным, содержащим ПДн (даже псевдонимизированным), должен быть строго ограничен.
    *   Логирование доступа к данным, содержащим ПДн.
*   Шифрование данных при передаче (TLS/SSL) и в состоянии покоя (шифрование на стороне S3, шифрование дисков для БД).
*   Контроль доступа к Data Lake, DWH, и Metadata Store на уровне пользователей и сервисных аккаунтов.

### 9.4. Управление Секретами
*   Все секреты (пароли к БД, ключи доступа к S3, Kafka, MLflow) хранятся в Kubernetes Secrets или HashiCorp Vault. Доступ к ним осуществляется через переменные окружения или API Vault.

## 10. Развертывание (Deployment)
См. `../../../../project_deployment_standards.md` для общих стандартов.

### 10.1. Инфраструктурные Файлы
*   **Dockerfiles:** Отдельные Dockerfile для API сервиса, для Spark/Flink приложений, для ML моделей.
*   **Helm-чарты/Kubernetes манифесты:** `deploy/charts/analytics-service/` (для API сервиса), отдельные конфигурации для запуска Spark/Flink джобов на Kubernetes.

### 10.2. Зависимости при Развертывании
*   ClickHouse, PostgreSQL (для метаданных), S3-совместимое хранилище, Kafka.
*   MLflow (если используется для управления моделями).
*   Доступ к Kafka топикам всех других микросервисов платформы.
*   Кластер Kubernetes с поддержкой запуска Spark/Flink приложений (если используется такой подход).

### 10.3. CI/CD
*   Пайплайны для сборки и тестирования API сервисов.
*   Пайплайны для сборки, тестирования и развертывания ETL/ELT джобов (Spark, Flink).
*   Пайплайны для тренировки, валидации и развертывания ML моделей (MLOps).
*   Автоматическое применение миграций схемы для PostgreSQL (метаданные) и ClickHouse (если применимо).

## 11. Мониторинг и Логирование (Logging and Monitoring)
См. `../../../../project_observability_standards.md` для общих стандартов.

### 11.1. Логирование
*   **Формат:** JSON (Zap для Go, стандартные логгеры для Python/Scala/Java с JSON-форматтерами).
*   **Ключевые события:**
    *   API Layer: Запросы к API, ошибки, авторизация.
    *   Data Ingestion: Статус приема данных из Kafka, количество полученных/отброшенных событий, ошибки парсинга.
    *   Data Processing (Stream/Batch): Статус запуска и завершения джобов, количество обработанных/ошибочных записей, длительность выполнения, использование ресурсов.
    *   ML Layer: Статус тренировки моделей, метрики качества, ошибки при прогнозировании.
*   **Интеграция:** Fluent Bit для сбора логов и отправки в Elasticsearch/Loki/ClickHouse.

### 11.2. Мониторинг
*   **Метрики (Prometheus):**
    *   API Layer: `http_requests_total`, `http_request_duration_seconds`, `api_errors_total`.
    *   Data Ingestion (Kafka Consumers): `kafka_consumer_lag_seconds`, `kafka_messages_consumed_total`, `kafka_consumer_errors_total`.
    *   Stream Processing (Flink/Kafka Streams): Метрики специфичные для фреймворка (например, `flink_job_uptime`, `flink_records_processed_per_second`, `flink_job_last_checkpoint_duration_ms`).
    *   Batch Processing (Spark): Метрики Spark джобов (длительность выполнения, использование ресурсов, количество прочитанных/записанных данных).
    *   DWH (ClickHouse): Производительность запросов, использование диска, количество активных соединений, ошибки.
    *   ML Model Serving: `ml_prediction_requests_total`, `ml_prediction_duration_seconds`, `ml_prediction_errors_total`.
    *   Актуальность данных в витринах (Data Freshness): `data_mart_last_update_timestamp_gauge{mart_name="daily_user_activity"}`.
*   **Дашборды (Grafana):**
    *   Обзор состояния пайплайнов данных (ingestion, stream, batch).
    *   Производительность и здоровье DWH (ClickHouse).
    *   Использование ресурсов компонентами Analytics Service.
    *   Производительность API Analytics Service.
    *   Ключевые бизнес-метрики в реальном времени (если применимо).
*   **Алерты (AlertManager):**
    *   Сбои в пайплайнах обработки данных (ETL/ELT джобы).
    *   Значительная задержка обработки данных в Kafka (высокий consumer lag).
    *   Недоступность или проблемы с производительностью ClickHouse/PostgreSQL/S3.
    *   Высокий процент ошибок API Analytics Service.
    *   Превышение порогов по актуальности данных (data freshness SLOs).

### 11.3. Трассировка
*   **Интеграция:** OpenTelemetry для API Layer и, где возможно, для отслеживания прохождения данных через критичные этапы обработки.
*   **Экспорт:** Jaeger или другой совместимый коллектор.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**: См. детализированные NFRs в существующем документе. Ключевые аспекты: задержка сбора и обработки данных, скорость ответа API, время генерации отчетов.
*   **Масштабируемость**: Горизонтальное масштабирование всех компонентов. Способность обрабатывать рост объема данных и количества запросов.
*   **Надежность**: Гарантированная доставка событий, отказоустойчивость пайплайнов, доступность API. RTO/RPO для хранилищ данных.
*   **Актуальность данных (Data Freshness)**: Определенные SLO для обновления real-time метрик, ежедневных и еженедельных отчетов.
*   **Точность данных (Data Accuracy)**: Минимальное расхождение с операционными системами, механизмы валидации.
*   **Сопровождаемость**: Покрытие кода тестами, актуальная документация, легкость развертывания и мониторинга.

## 13. Приложения (Appendices)
*   Детальные JSON схемы для API запросов/ответов будут доступны через публикуемую OpenAPI спецификацию сервиса.
*   Схемы событий (Avro, Protobuf, или JSON Schema), используемые в Kafka, будут храниться в централизованном реестре схем (Schema Registry) или в общем репозитории артефактов (`platform-protos` или аналогичном).

## 14. Пользовательские Сценарии (User Flows)

В этом разделе описаны ключевые сценарии использования Analytics Service, демонстрирующие его роль в экосистеме платформы.

### 14.1. Поток Приема Данных для Нового Типа События

*   **Описание:** Разработчик нового микросервиса (или существующего) внедряет отправку нового типа события в Kafka. Analytics Service должен быть сконфигурирован для приема, обработки и сохранения этого нового типа события.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        participant Developer as Разработчик
        participant NewService as Новый/Обновленный Сервис
        participant KafkaBus as Kafka Message Bus
        participant AnalyticsIngestion as Analytics Service (Ingestion Layer)
        participant DataLakeS3 as S3 Data Lake
        participant StreamProc as Analytics Service (Stream Processing)
        participant MetadataDB as Analytics Service (Metadata DB - PostgreSQL)
        participant Alerting as Система Алертинга

        Developer->>NewService: Внедряет генерацию нового события (например, `com.platform.newfeature.interaction.v1`)
        Developer->>MetadataDB: (Через UI или API Admin/Analytics) Регистрирует схему нового события, если используется Schema Registry
        NewService->>KafkaBus: Publish `com.platform.newfeature.interaction.v1` (сообщение)

        AnalyticsIngestion->>KafkaBus: Consume событие (подписка на топик или wildcard)
        alt Схема события известна и валидна
            AnalyticsIngestion->>DataLakeS3: Сохранение сырого события (например, в Parquet/JSON.gz)
            AnalyticsIngestion->>StreamProc: Передача события для потоковой обработки (если применимо)
            StreamProc->>StreamProc: Обработка события (например, обогащение, агрегация real-time метрик)
        else Схема неизвестна или невалидна
            AnalyticsIngestion->>KafkaBus: Отправка события в Dead Letter Queue (DLQ)
            AnalyticsIngestion->>Alerting: Уведомление о неизвестном/невалидном событии
        end

        Note over Developer, MetadataDB: Администратор аналитики может настроить правила обработки и агрегации для нового события в MetadataDB.
    ```

### 14.2. Генерация Отчета "Daily Active Users (DAU)" Аналитиком

*   **Описание:** Аналитик или менеджер продукта запрашивает отчет по DAU за определенный период через API Analytics Service (например, из админ-панели или BI-инструмента).
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor Analyst as Аналитик
        participant ClientTool as BI Tool / Admin Panel
        participant AnalyticsAPI as Analytics Service (API Layer)
        participant DWHClickHouse as Data Warehouse (ClickHouse)
        participant MetadataDB as Analytics Service (Metadata DB - PostgreSQL)

        Analyst->>ClientTool: Запрос отчета DAU (период: последние 7 дней, группировка: по платформе)
        ClientTool->>AnalyticsAPI: GET /api/v1/analytics/reports/instances (definition_name="dau_report", params={"period_start":"YYYY-MM-DD", "period_end":"YYYY-MM-DD", "group_by":"platform"})

        AnalyticsAPI->>MetadataDB: Загрузка определения отчета "dau_report" (SQL-шаблон, параметры)
        AnalyticsAPI->>AnalyticsAPI: Формирование SQL-запроса к DWH на основе шаблона и параметров
        AnalyticsAPI->>DWHClickHouse: Выполнение SQL-запроса (например, SELECT activity_date, platform, dau_count FROM analytics_db.mart_daily_platform_activity_mv WHERE ...)
        DWHClickHouse-->>AnalyticsAPI: Результаты запроса (агрегированные данные DAU)
        AnalyticsAPI->>AnalyticsAPI: Форматирование данных в запрошенный формат (JSON/CSV)
        opt Сохранение экземпляра отчета
            AnalyticsAPI->>MetadataDB: Сохранение информации об экземпляре отчета
            AnalyticsAPI->>DataLakeS3: (Если отчет большой) Сохранение файла отчета в S3
        end
        AnalyticsAPI-->>ClientTool: HTTP 200 OK (данные отчета или ссылка на файл)
        ClientTool-->>Analyst: Отображение отчета (графики, таблицы)
    ```

### 14.3. Тренировка Новой Версии Модели Рекомендаций Игр
*   **Описание:** Data Scientist запускает пайплайн для тренировки новой версии ML-модели рекомендаций игр, используя свежие данные о поведении пользователей.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor DataScientist as Data Scientist
        participant Orchestrator as Data Pipeline Orchestrator (Airflow/Prefect)
        participant BatchProcSpark as Analytics Service (Batch Processing - Spark)
        participant DataLakeS3 as S3 Data Lake (Features, Model Artifacts)
        participant DWHClickHouse as DWH (Aggregated Features)
        participant MLTrainingEnv as ML Training Environment (Python, SparkML)
        participant MLModelRegistry as ML Model Registry (MLflow)
        participant MetadataDB as Analytics Service (Metadata DB)

        DataScientist->>Orchestrator: Запуск пайплайна "train_recommendation_model_v2"
        Orchestrator->>BatchProcSpark: Запуск Spark-джобы для подготовки признаков (feature engineering)
        BatchProcSpark->>DataLakeS3: Чтение сырых/обработанных данных о поведении пользователей
        BatchProcSpark->>DWHClickHouse: Чтение агрегированных данных
        BatchProcSpark->>DataLakeS3: Сохранение набора данных для тренировки (feature store)

        Orchestrator->>MLTrainingEnv: Запуск скрипта тренировки модели (параметры, путь к данным)
        MLTrainingEnv->>DataLakeS3: Загрузка тренировочных данных
        MLTrainingEnv->>MLTrainingEnv: Тренировка модели (например, ALS, LightFM)
        MLTrainingEnv->>MLTrainingEnv: Оценка качества модели (метрики: precision@k, recall@k)
        MLTrainingEnv->>MLModelRegistry: Регистрация новой версии модели (артефакты, параметры, метрики)
        MLModelRegistry->>DataLakeS3: Сохранение артефактов модели
        MLModelRegistry->>MetadataDB: Обновление метаданных о модели (статус, путь к артефактам)

        Orchestrator-->>DataScientist: Уведомление о завершении тренировки (статус, метрики)
    ```

### 14.4. Получение Системой Алерта на Основе Обнаружения Аномалий
*   **Описание:** Analytics Service (его компонент потоковой обработки или специальный ML-сервис) обнаруживает аномалию в данных (например, резкое падение количества транзакций) и генерирует алерт.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        participant StreamProcFlink as Analytics Service (Stream Processing - Flink/Kafka Streams)
        participant AnomalyDetectionSvc as Analytics Service (Anomaly Detection Model/Rules)
        participant KafkaBus as Kafka Message Bus
        participant AlertingSystem as Система Алертинга (внешняя или часть Admin Service)
        participant NotificationSvc as Notification Service

        StreamProcFlink->>StreamProcFlink: Обработка потока событий транзакций
        StreamProcFlink->>AnomalyDetectionSvc: Передача данных/метрик для анализа (например, количество транзакций в минуту)
        AnomalyDetectionSvc->>AnomalyDetectionSvc: Обнаружение аномалии (например, значение ниже порога N сигм)
        AnomalyDetectionSvc->>KafkaBus: Publish `com.platform.analytics.alert.triggered.v1` (alertName="TransactionRateDrop", severity="critical", details)

        KafkaBus->>AlertingSystem: Consume событие алерта
        AlertingSystem->>AlertingSystem: Регистрация алерта, определение правил эскалации
        AlertingSystem->>NotificationSvc: Запрос на отправку уведомления (ответственным лицам)
        NotificationSvc->>NotificationSvc: Отправка уведомлений (Email, SMS, Push)
    ```

## 15. Резервное Копирование и Восстановление (Backup and Recovery)

### 14.1. PostgreSQL (Metadata Store)
*   **Процедура резервного копирования:**
    *   Ежедневный логический бэкап (`pg_dump`) базы метаданных.
    *   Настроена непрерывная архивация WAL-сегментов (PITR) для возможности восстановления на любой момент времени.
    *   **Хранение:** Бэкапы и WAL-архивы хранятся в S3-совместимом хранилище с шифрованием и версионированием, в другом регионе. Срок хранения: полные бэкапы - 30 дней, WAL - 14 дней.
*   **Процедура восстановления:** Тестируется ежеквартально.
*   **RTO:** < 2 часов.
*   **RPO:** < 10 минут.

### 14.2. ClickHouse (DWH Data)
*   **Процедура резервного копирования:**
    *   Использование инструмента `clickhouse-backup` (или аналогичного) для создания инкрементальных или полных бэкапов таблиц.
    *   Для критически важных таблиц или кластеров может быть настроена репликация ClickHouse.
    *   **Частота:** Ежедневно для основных витрин данных и агрегатов. Менее часто для исторических партиций, которые не изменяются.
    *   **Хранение:** Бэкапы хранятся в S3-совместимом хранилище. Срок хранения зависит от критичности данных и требований к хранению (например, 30-90 дней).
*   **Процедура восстановления:**
    *   Восстановление из бэкапов. При наличии репликации – переключение на реплику.
    *   Тестируется ежеквартально для ключевых наборов данных.
*   **RTO:** < 8 часов (для больших объемов данных).
*   **RPO:** 24 часа (для ежедневных бэкапов). Если используется асинхронная репликация, RPO может быть значительно меньше для реплицированных данных.

### 14.3. S3 Data Lake (Raw Events, ML Artifacts)
*   **Процедура резервного копирования:**
    *   **Версионирование объектов:** Включено для всех бакетов Data Lake для защиты от случайного удаления или перезаписи.
    *   **Политики жизненного цикла (Lifecycle Policies):** Для управления хранением старых версий и перемещения редко используемых данных в более холодные классы хранения.
    *   **Cross-Region Replication (CRR):** Может быть настроена для критически важных бакетов для обеспечения гео-резервирования.
*   **Процедура восстановления:** Восстановление предыдущих версий объектов или восстановление из реплики в другом регионе.
*   **RTO:** Зависит от объема данных и скорости S3, но обычно быстро для отдельных объектов/партиций.
*   **RPO:** Близко к нулю при использовании версионирования и CRR.

### 14.4. Kafka (События)
*   **Стратегия:** Основной упор на хранение сырых данных в S3 Data Lake. Kafka используется как буфер и для потоковой обработки.
*   **Retention в Kafka:** Топики с сырыми событиями могут иметь ограниченный срок хранения (например, 7-14 дней), достаточный для обработки и выгрузки в S3.
*   **Резервное копирование топиков Kafka:** Обычно не производится, если данные надежно архивируются в S3. Однако, для некоторых критичных обработанных/промежуточных топиков, если их восстановление из S3 трудоемко, может рассматриваться использование инструментов типа Kafka MirrorMaker для репликации в другой кластер или бэкап с помощью специализированных инструментов (редко).

### 14.5. MLflow Artifacts & Metadata (если MLflow используется и его БД отдельна)
*   **Артефакты (модели, данные):** Хранятся в S3 и подпадают под стратегию бэкапа S3 Data Lake.
*   **Метаданные MLflow (если используется внешняя БД, например, PostgreSQL):** Бэкапируются аналогично другим PostgreSQL базам данных (см. п. 14.1).

### 14.6. Общая стратегия
*   Резервное копирование и восстановление компонентов Analytics Service являются частью общей стратегии обеспечения непрерывности бизнеса платформы.
*   Процедуры документированы и регулярно пересматриваются.
*   Мониторинг процессов резервного копирования.
*   Общие принципы резервного копирования для различных СУБД описаны в `../../../../project_database_structure.md`.

## 16. Связанные Рабочие Процессы (Related Workflows)
*   [Генерация аналитического отчета по запросу администратора] Подробное описание этого рабочего процесса будет добавлено в [admin_report_generation_flow.md](../../../../project_workflows/admin_report_generation_flow.md) (документ в разработке).
*   [Процесс тренировки и развертывания ML модели для рекомендаций] Подробное описание этого рабочего процесса будет добавлено в [ml_model_training_deployment_flow.md](../../../../project_workflows/ml_model_training_deployment_flow.md) (документ в разработке).

---
*Этот документ является отправной точкой и должен регулярно обновляться по мере развития сервиса.*
