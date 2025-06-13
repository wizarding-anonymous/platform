# Спецификация Микросервиса: Analytics Service

**Версия:** 1.0
**Дата последнего обновления:** {{YYYY-MM-DD}}

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   Analytics Service является центральным компонентом платформы "Российский Аналог Steam", ответственным за сбор, обработку, анализ и предоставление данных и инсайтов, генерируемых на платформе.
*   Его основная роль - поддержка принятия решений на основе данных для бизнес-стратегии, операционных улучшений, улучшения пользовательского опыта, а также предоставление разработчикам релевантной статистики по их продуктам.
*   Основные бизнес-задачи: сбор и обработка данных, расчет метрик и KPI, генерация отчетов, анализ поведения пользователей, сегментация аудитории, поддержка предиктивной аналитики.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Сбор данных:** Прием событий и данных от всех микросервисов платформы.
*   **Обработка данных:** Потоковая и пакетная обработка; анонимизация и псевдонимизация ПДн.
*   **Метрики и Отчетность:** Расчет KPI и продуктовых метрик; генерация отчетов.
*   **Анализ поведения пользователей:** Анализ путей, когорт, A/B тестов, воронок.
*   **Сегментация аудитории:** Создание и управление статическими и динамическими сегментами.
*   **Предиктивная аналитика:** Разработка, тренировка, развертывание и мониторинг ML моделей.
*   **Мониторинг производительности системы:** Агрегация и анализ технических метрик.

### 1.3. Основные Технологии
*   **API Layer:** Go (Echo) или Java (Spring Boot); Потенциально GraphQL, WebSocket.
*   **Data Processing Layer:** Python, Scala, Java; Apache Spark, Apache Flink/Kafka Streams.
*   **Data Storage Layer:** ClickHouse (DWH), S3-совместимое хранилище (Data Lake), PostgreSQL (Метаданные).
*   **Messaging Layer:** Apache Kafka.
*   **Machine Learning Layer:** Python (TensorFlow, PyTorch, Scikit-learn), MLflow.
*   **Общие:** Viper, Zap (для Go), OpenTelemetry, Prometheus.
*   **Визуализация (внешние):** Grafana, Apache Superset.
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `../../../../PACKAGE_STANDARDIZATION.md`.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   См. `../../../../project_glossary.md`.
*   **Событие (Event):** Атомарная запись о действии пользователя или системы.
*   **KPI (Key Performance Indicator):** Ключевой показатель эффективности.
*   **Сегмент (Segment):** Группа пользователей с общими характеристиками.
*   **ML Модель (Machine Learning Model):** Алгоритм для прогнозов или классификаций.
*   **Data Lake:** Централизованное хранилище сырых данных.
*   **DWH (Data Warehouse):** Централизованное хранилище структурированных данных для аналитики.
*   **ETL/ELT:** Процессы извлечения, преобразования и загрузки данных.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Analytics Service построен на событийно-ориентированной архитектуре, сочетая элементы Lambda и Kappa архитектур для пакетной и потоковой обработки.
*   API и управляющие компоненты могут следовать стандартной слоистой архитектуре. Основная часть сервиса — конвейеры данных.
*   Ключевые компоненты: Сбор данных, Обработка данных, Хранение данных, Слой доступа к данным (API), Слой ML моделей.
*   Диаграмма верхнеуровневой архитектуры сервиса:
```mermaid
graph TD
    subgraph "Источники Данных (Все Микросервисы Платформы)"
        direction LR
        KafkaTopics[Kafka: Topics (account.events, catalog.events, etc.)]
    end

    subgraph "Analytics Service"
        direction TB

        IngestionLayer[Data Ingestion Layer (Kafka Consumers)]
        KafkaTopics --> IngestionLayer

        subgraph "Data Processing & Storage"
            direction TB
            RawDataLake[Data Lake (S3: сырые события - Parquet/Avro)]
            IngestionLayer --> RawDataLake

            StreamProcessing[Stream Processing (Flink/Kafka Streams)]
            IngestionLayer --> StreamProcessing

            RealtimeDWH[Real-time/Speed Layer (ClickHouse: витрины реального времени)]
            StreamProcessing --> RealtimeDWH

            BatchProcessing[Batch Processing (Spark ETL/ELT)]
            RawDataLake --> BatchProcessing
            StreamProcessing -- Processed/Enriched Events --> BatchProcessingViaKafka[Kafka (для батч-триггеров)] -- опционально --> BatchProcessing


            DWH[Core DWH (ClickHouse: агрегаты, исторические витрины)]
            BatchProcessing --> DWH

            MLDataStore[ML Feature Store & Model Artifacts (S3/ClickHouse/PostgreSQL)]
            BatchProcessing --> MLDataStore
            StreamProcessing --> MLDataStore
        end

        subgraph "API & Management Layer (Go / Java)"
            direction TB
            MetadataDB[Metadata Store (PostgreSQL: определения метрик, отчетов, ML-моделей, пайплайнов)]

            APIService[API Service (REST/GraphQL - Presentation)]
            APIServiceLogic[Application & Domain Logic for API]
            APIService --> APIServiceLogic
            APIServiceLogic --> MetadataDB
            APIServiceLogic -- Queries --> DWH
            APIServiceLogic -- Queries --> RealtimeDWH
            APIServiceLogic -- Predictions --> MLServingAPI[ML Model Serving API]

            PipelineOrchestrator[Data Pipeline Orchestrator (e.g., Airflow, Prefect, или кастомный)]
            PipelineOrchestrator --> BatchProcessing
            PipelineOrchestrator --> MetadataDB
        end

        subgraph "ML Layer"
            direction TB
            MLTraining[ML Model Training (SparkML, Python libs)]
            MLDataStore --> MLTraining
            MLTraining --> MLModelRegistry[ML Model Registry (MLflow / PostgreSQL)]
            MetadataDB <--> MLModelRegistry

            MLServingAPI[ML Model Serving API (Python/Java/Go)]
            MLModelRegistry --> MLServingAPI
            MLDataStore -- Model Artifacts --> MLServingAPI
        end

        APIService --> ExternalConsumers[Потребители API (Admin Panel, Developer Portal, BI Tools)]
        MLServingAPI --> ExternalConsumersML[Потребители ML API (Recommendation Svc, Fraud Detection)]

    end

    classDef dataFlow fill:#e6f3ff,stroke:#007bff,stroke-width:2px;
    classDef component fill:#d4edda,stroke:#28a745,stroke-width:2px;
    classDef datastore fill:#f8d7da,stroke:#dc3545,stroke-width:2px;
    classDef api fill:#fff3cd,stroke:#ffc107,stroke-width:2px;
    classDef external fill:#e2e3e5,stroke:#6c757d,stroke-width:2px;

    class KafkaTopics,BatchProcessingViaKafka dataFlow;
    class IngestionLayer,StreamProcessing,BatchProcessing,MLDataStore,MLTraining,MLModelRegistry,PipelineOrchestrator,APIServiceLogic component;
    class RawDataLake,DWH,RealtimeDWH,MetadataDB,MLDataStore datastore;
    class APIService,MLServingAPI api;
    class ExternalConsumers,ExternalConsumersML external;
```

### 2.2. Слои Сервиса (Компоненты)

#### 2.2.1. Data Ingestion Layer (Слой Приема Данных)
*   Ответственность: Прием событий из Kafka, базовая валидация, сохранение сырых данных в Data Lake (S3), передача в потоковую обработку.
*   Ключевые технологии: Kafka Consumers, S3 коннекторы, Avro/Protobuf.

#### 2.2.2. Data Processing Layer (Слой Обработки Данных)
*   Ответственность: Очистка, трансформация, обогащение, агрегация данных.
*   **Stream Processing:** Apache Flink или Kafka Streams. Расчет real-time метрик, обогащение, простые паттерны.
*   **Batch Processing:** Apache Spark. Сложные ETL/ELT, пересчет истории, построение витрин в DWH, подготовка данных для ML.

#### 2.2.3. Data Storage Layer (Слой Хранения Данных)
*   **Data Lake (S3):** Хранение сырых событий (Parquet/ORC), партиционирование.
*   **Data Warehouse (DWH - ClickHouse):** Структурированные, агрегированные данные, витрины данных.
*   **Real-time Data Marts (ClickHouse/Redis):** Часто обновляемые метрики.
*   **Metadata Store (PostgreSQL):** Определения метрик, отчетов, сегментов, конфигурации пайплайнов, схемы, метаданные ML моделей.

#### 2.2.4. API & Management Layer (Слой API и Управления)
Применяет принципы Clean Architecture для своих компонентов.
*   **Presentation Layer (API Service):** Go (Echo) или Java (Spring Boot). Обработка REST (GraphQL/WebSocket потенциально) запросов.
*   **Application Layer:** Координация запросов, формирование ответов. Компоненты: `MetricService`, `ReportService`, `SegmentService`.
*   **Domain Layer:** Бизнес-логика управления метаданными. Сущности: `MetricDefinition`, `ReportDefinition`.
*   **Infrastructure Layer:** Взаимодействие с Metadata Store (PostgreSQL), DWH (ClickHouse), ML Model Serving API.

#### 2.2.5. Machine Learning Layer (Слой Машинного Обучения)
*   **ML Model Training:** Python, SparkML, TensorFlow, PyTorch, Scikit-learn. Feature engineering, обучение, оценка, версионирование.
*   **ML Model Registry (MLflow или аналог):** Хранение артефактов моделей, версий, метрик.
*   **ML Model Serving:** Python (Flask/FastAPI), Java (Spring Boot) или спец. решения. REST/gRPC API для прогнозов.

#### 2.2.6. Data Pipeline Orchestration
*   Ответственность: Управление и мониторинг пайплайнов.
*   Технологии: Apache Airflow, Prefect, или кастомные решения.

## 3. API Endpoints
Общие принципы и форматы см. `../../../../project_api_standards.md`.
Детальные JSON схемы для API будут доступны через публикуемую OpenAPI спецификацию.

### 3.1. REST API
*   **Базовый URL (через API Gateway):** `/api/v1/analytics`
*   **Аутентификация:** JWT Bearer Token.
*   **Авторизация:** На основе ролей (см. `../../../../project_roles_and_permissions.md`).

#### 3.1.1. Ресурс: Метрики (Metrics)
*   **`GET /metrics/definitions`**
    *   Описание: Получение списка определений всех доступных метрик.
    *   Query параметры: `tag`, `owner_service`, `is_realtime`.
    *   Пример ответа: `[NEEDS DEVELOPER INPUT: Example success response for GET /metrics/definitions]`
    *   Требуемые права доступа: `platform_admin`, `developer`, `marketing_manager`.
*   **`GET /metrics/definitions/{metric_name}`**
    *   Описание: Получение детального определения конкретной метрики.
    *   Требуемые права доступа: `platform_admin`, `developer`, `marketing_manager`.
*   **`POST /metrics/definitions`**
    *   Описание: Создание нового определения метрики (для администраторов).
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for POST /metrics/definitions]`
    *   Требуемые права доступа: `platform_admin`.
*   **`GET /metrics/values/{metric_name}`**
    *   Описание: Получение рассчитанных значений метрики за период.
    *   Query параметры: `period_start`, `period_end`, `granularity` (day, hour), `dimensions` (json), `filters` (json).
    *   Пример ответа: `[NEEDS DEVELOPER INPUT: Example success response for GET /metrics/values/{metric_name}]`
    *   Требуемые права доступа: `platform_admin`, `developer` (с ограничениями по своим играм), `marketing_manager`.

#### 3.1.2. Ресурс: Отчеты (Reports)
*   **`GET /reports/definitions`**
    *   Описание: Получение списка доступных определений отчетов.
    *   Требуемые права доступа: `platform_admin`, `marketing_manager`.
*   **`POST /reports/instances`**
    *   Описание: Запрос на генерацию экземпляра отчета по его определению.
    *   Тело запроса: `{"report_definition_id": "uuid", "parameters": {"period_start": "YYYY-MM-DD", ...}, "output_format": "csv"}`
    *   Пример ответа: `[NEEDS DEVELOPER INPUT: Example success response for POST /reports/instances]`
    *   Требуемые права доступа: `platform_admin`, `marketing_manager`.
*   **`GET /reports/instances/{instance_id}`**
    *   Описание: Получение статуса и метаданных сгенерированного экземпляра отчета.
    *   Требуемые права доступа: Владелец запроса, `platform_admin`.
*   **`GET /reports/instances/{instance_id}/download`**
    *   Описание: Скачивание файла сгенерированного отчета.
    *   Требуемые права доступа: Владелец запроса, `platform_admin`.

#### 3.1.3. Ресурс: Сегменты (Segments)
*   **`GET /segments/definitions`**
    *   Описание: Получение списка определений пользовательских сегментов.
    *   Требуемые права доступа: `platform_admin`, `marketing_manager`.
*   **`POST /segments/definitions`**
    *   Описание: Создание нового определения сегмента.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for POST /segments/definitions]`
    *   Требуемые права доступа: `platform_admin`.
*   **`GET /segments/{segment_id}/users-count`**
    *   Описание: Получение текущего количества пользователей в сегменте.
    *   Требуемые права доступа: `platform_admin`, `marketing_manager`.
*   **`POST /segments/{segment_id}/refresh`**
    *   Описание: Принудительный пересчет сегмента.
    *   Требуемые права доступа: `platform_admin`.

#### 3.1.4. Ресурс: Предиктивная Аналитика (Predictions)
*   **`GET /predictions/models`**
    *   Описание: Получение списка доступных ML моделей и их метаданных.
    *   Требуемые права доступа: `platform_admin`, `developer` (для специфичных моделей).
*   **`POST /predictions/{model_name_or_id}/predict`**
    *   Описание: Запрос прогноза от указанной ML модели.
    *   Тело запроса: `{"features": {"feature1": "value1", ...}}`
    *   Пример ответа: `[NEEDS DEVELOPER INPUT: Example success response for POST /predictions/{model_name_or_id}/predict]`
    *   Требуемые права доступа: `platform_admin`, сервисные аккаунты (для Recommendation Svc и т.д.).

### 3.2. GraphQL API
*   **Эндпоинт:** `/api/v1/analytics/graphql` (потенциальный)
*   **Статус:** Не реализован. Рассматривается для будущего.
*   Ссылка на `.proto` файл: `[NEEDS DEVELOPER INPUT: Path to .proto if gRPC API exists and needs to be documented, otherwise state "Not Applicable"]`

### 3.3. WebSocket API (если применимо)
*   **Эндпоинт:** `/api/v1/analytics/ws/streaming` (потенциальный)
*   **Статус:** Не реализован. Рассматривается для будущего.

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности/Структуры Данных
*   **`Event` (Событие):** CloudEvents-совместимый формат. Хранится в S3 (сырые), ClickHouse (факты).
*   **`MetricDefinition`**: Определения метрик (PostgreSQL).
*   **`ReportDefinition`**: Определения отчетов (PostgreSQL).
*   **`ReportInstance`**: Экземпляры отчетов (PostgreSQL, файлы в S3).
*   **`SegmentDefinition`**: Определения сегментов (PostgreSQL).
*   **`MLModelMetadata`**: Метаданные ML моделей (PostgreSQL или MLflow).
*   **`DataPipelineRun`**: Информация о запусках пайплайнов (PostgreSQL).
*   [NEEDS DEVELOPER INPUT: Review and add any other key entities for analytics-service]

### 4.2. Схема Базы Данных (если применимо)
*   **PostgreSQL (Metadata Store):** DDL см. в существующем документе (секция 4.2.1). Диаграмма: `[NEEDS DEVELOPER INPUT: Mermaid ERD for PostgreSQL metadata tables if existing DDL is not sufficient or needs visualization]`
*   **ClickHouse (DWH):** Примеры DDL для таблиц фактов и витрин см. в существующем документе (секция 4.2.2). Диаграмма: `[NEEDS DEVELOPER INPUT: Mermaid ERD for key ClickHouse tables/views if needed]`
*   **S3 Data Lake Структура:** Описание см. в существующем документе (секция 4.2.3).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   Используемая система сообщений: Kafka. Формат: CloudEvents JSON.
*   Топик: `com.platform.analytics.events.v1`.
*   **`com.platform.analytics.report.generated.v1`**
    *   Описание: Отчет сгенерирован. Payload: (см. существующий документ).
*   **`com.platform.analytics.segment.updated.v1`**
    *   Описание: Пользовательский сегмент обновлен. Payload: (см. существующий документ).
*   **`com.platform.analytics.alert.triggered.v1`**
    *   Описание: Сработал аналитический алерт. Payload: (см. существующий документ).

### 5.2. Потребляемые События (Consumed Events)
*   Analytics Service потребляет события от **всех** других микросервисов (см. список в существующем документе).
*   Топики: Подписка на все релевантные топики событий.
*   Формат: CloudEvents JSON.

## 6. Интеграции (Integrations)
См. `../../../../project_integrations.md`.
[NEEDS DEVELOPER INPUT: Review and confirm service boundaries and reliability strategies for each integration for analytics-service]

### 6.1. Внутренние Микросервисы
*   **Потребление данных:** От всех микросервисов через Kafka.
*   **Предоставление данных (API):** Admin Service, Developer Service, Marketing Service (гипотетический), Recommendation Service (гипотетический).
*   **Notification Service:** Получение событий о готовности отчетов/алертов.
*   **Auth Service:** Аутентификация/авторизация запросов к API.

### 6.2. Внешние Системы
*   **S3-совместимое хранилище:** Для Data Lake и артефактов.
*   **Внешние BI-инструменты:** Подключение к DWH (ClickHouse).
*   **Grafana:** Для дашбордов.

## 7. Конфигурация (Configuration)
Общие стандарты: `../../../../project_api_standards.md` и `../../../../DOCUMENTATION_GUIDELINES.md`.

### 7.1. Переменные Окружения
*   `ANALYTICS_API_PORT`: Порт API.
*   `POSTGRES_DSN_METADATA`: DSN для PostgreSQL (метаданные).
*   `CLICKHOUSE_DSN`: DSN для ClickHouse.
*   `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`, `S3_BUCKET_LAKE`, `S3_BUCKET_REPORTS`, `S3_BUCKET_MLFLOW`: Настройки S3.
*   `KAFKA_BROKERS`: Брокеры Kafka.
*   `SPARK_MASTER_URL`: URL Spark Master (если используется).
*   `FLINK_MASTER_URL`: URL Flink Master (если используется).
*   `MLFLOW_TRACKING_URI`: URI для MLflow.
*   `LOG_LEVEL`.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`.
*   [NEEDS DEVELOPER INPUT: Add any other critical environment variables for analytics-service]

### 7.2. Файлы Конфигурации (`configs/config.yaml`)
*   Расположение: `backend/analytics-service/configs/config.yaml` (для API Layer).
*   Структура:
    ```yaml
    api_server:
      port: ${ANALYTICS_API_PORT:"8081"}
    postgres_metadata:
      dsn: ${POSTGRES_DSN_METADATA}
    clickhouse_dwh:
      dsn: ${CLICKHOUSE_DSN}
    s3_config:
      endpoint: ${S3_ENDPOINT}
      # ...
    kafka:
      brokers: ${KAFKA_BROKERS}
    # ... [NEEDS DEVELOPER INPUT: Add other specific config sections for analytics-service API layer, Spark/Flink jobs might have their own config files]
    ```

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   **Data Pipelines:** DLQ в Kafka, retry механизмы, мониторинг, логирование.
*   **API Layer:** Стандартные коды и форматы ошибок REST.

### 8.2. Распространенные Коды Ошибок (API)
*   **`METRIC_DEFINITION_NOT_FOUND`**
*   **`REPORT_DEFINITION_NOT_FOUND`**
*   **`REPORT_GENERATION_FAILED`**
*   **`ML_MODEL_NOT_FOUND`**
*   **`INVALID_QUERY_PARAMETERS`**
*   **`DATA_NOT_YET_AVAILABLE`**
*   **`ACCESS_DENIED_TO_RESOURCE`**
*   [NEEDS DEVELOPER INPUT: Review and add other specific error codes for analytics-service]

## 9. Безопасность (Security)
См. `../../../../project_security_standards.md`.

### 9.1. Аутентификация
*   API: JWT.
*   Внутренние компоненты: Специфичные для систем (SASL Kafka, IAM S3, etc.).

### 9.2. Авторизация
*   RBAC для API. Потенциально ABAC для данных.

### 9.3. Защита Данных
*   **ФЗ-152 "О персональных данных":** Анонимизация/псевдонимизация ПДн перед использованием в аналитике. Строгий контроль доступа к сырым данным.
*   Шифрование при передаче и в покое.

### 9.4. Управление Секретами
*   Kubernetes Secrets или HashiCorp Vault.

## 10. Развертывание (Deployment)
См. `../../../../project_deployment_standards.md`.

### 10.1. Инфраструктурные Файлы
*   Dockerfiles: для API, Spark/Flink приложений, ML моделей.
*   Helm-чарты: `deploy/charts/analytics-service/`. Конфигурации для Spark/Flink джобов.

### 10.2. Зависимости при Развертывании
*   ClickHouse, PostgreSQL, S3, Kafka, MLflow.
*   Доступ к Kafka топикам всех сервисов.
*   Kubernetes кластер (с поддержкой Spark/Flink если нужно).

### 10.3. CI/CD
*   Пайплайны для API, ETL/ELT джобов, MLOps.
*   Автоматические миграции схем БД.

## 11. Мониторинг и Логирование (Logging and Monitoring)
См. `../../../../project_observability_standards.md`.

### 11.1. Логирование
*   Формат: JSON.
*   Ключевые события: API, Data Ingestion, Data Processing, ML Layer.
*   Интеграция: Fluent Bit в Elasticsearch/Loki/ClickHouse.

### 11.2. Мониторинг
*   Метрики (Prometheus): API, Kafka Consumers, Stream Processing, Batch Processing, DWH, ML Model Serving, Data Freshness.
*   Дашборды (Grafana): Состояние пайплайнов, DWH, использование ресурсов, API, real-time метрики.
*   Алерты (AlertManager): Сбои пайплайнов, задержки Kafka, проблемы с БД/S3, ошибки API, нарушение SLO по Data Freshness.

### 11.3. Трассировка
*   Интеграция: OpenTelemetry для API и критичных этапов обработки. Экспорт: Jaeger.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**: Задержка сбора данных (stream: <5s, batch: <1h до Data Lake), API (P95 <500ms для метрик, <2s для отчетов), генерация отчетов (сложные: <15min).
*   **Масштабируемость**: Обработка до X TB данных в день, Y событий в секунду.
*   **Надежность**: Доступность API 99.9%. RTO/RPO для хранилищ.
*   **Актуальность данных**: Real-time метрики (<1min), ежедневные отчеты (<24h).
*   **Точность данных**: Расхождение < Z % с операционными системами.
*   [NEEDS DEVELOPER INPUT: Confirm or update specific NFR values for analytics-service, X, Y, Z etc.]

## 13. Резервное Копирование и Восстановление (Backup and Recovery)
Общие принципы см. в `../../../../project_database_structure.md`.

### 13.1. PostgreSQL (Metadata Store)
*   Процедура: Ежедневный `pg_dump`, непрерывная архивация WAL (PITR).
*   Хранение: S3. RTO: < 2ч. RPO: < 10мин.

### 13.2. ClickHouse (DWH Data)
*   Процедура: `clickhouse-backup` (инкрементальные/полные). Репликация.
*   Хранение: S3. RTO: < 8ч. RPO: 24ч.

### 13.3. S3 Data Lake (Raw Events, ML Artifacts)
*   Процедура: Версионирование, Lifecycle Policies, Cross-Region Replication (CRR).
*   RTO: Зависит от объема. RPO: Близко к нулю.

### 13.4. Kafka (События)
*   Стратегия: Основной упор на хранение в S3. Kafka - буфер с ограниченным retention (7-14 дней).

### 13.5. MLflow Artifacts & Metadata
*   Артефакты: S3 (см. 13.3). Метаданные (PostgreSQL): см. 13.1.

## 14. Приложения (Appendices) (Опционально)
*   OpenAPI спецификация: `[NEEDS DEVELOPER INPUT: Link to OpenAPI spec or state if generated from code for analytics-service]`
*   Схемы событий (Avro/Protobuf/JSON Schema): `[NEEDS DEVELOPER INPUT: Link to schema registry or artifact repository for analytics-service event schemas]`
*   [NEEDS DEVELOPER INPUT: Add any other appendices if necessary for analytics-service]

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [NEEDS DEVELOPER INPUT: Add link to admin_report_generation_flow.md when created]
*   [NEEDS DEVELOPER INPUT: Add link to ml_model_training_deployment_flow.md when created]
*   [NEEDS DEVELOPER INPUT: Add links to other relevant high-level workflow documents for analytics-service]

---
*Этот документ является отправной точкой и должен регулярно обновляться по мере развития сервиса.*
