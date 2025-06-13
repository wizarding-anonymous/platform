# Спецификация Микросервиса: Analytics Service

**Версия:** 1.0
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   Analytics Service является центральным компонентом платформы "Российский Аналог Steam", ответственным за сбор, обработку, анализ и предоставление данных и инсайтов, генерируемых на платформе.
*   Его основная роль - поддержка принятия решений на основе данных для бизнес-стратегии, операционных улучшений, улучшения пользовательского опыта, а также предоставление разработчикам релевантной статистики по их продуктам.
*   Основные бизнес-задачи: сбор и обработка данных, расчет метрик и KPI, генерация отчетов, анализ поведения пользователей, сегментация аудитории, поддержка предиктивной аналитики.

### 1.2. Ключевые Функциональности
*   **Сбор данных:** Прием событий и данных от всех микросервисов (действия пользователей, транзакции, системные логи, маркетинговые взаимодействия). Поддержка потоковой и пакетной загрузки.
*   **Обработка данных:** Потоковая и пакетная обработка сырых данных (трансформация, агрегация, очистка, анонимизация PII в соответствии с 152-ФЗ).
*   **Метрики и Отчетность:** Расчет KPI и метрик (DAU, MAU, ARPU, конверсии, удержание). Генерация стандартных и настраиваемых отчетов.
*   **Анализ поведения пользователей:** Инструменты для анализа путей пользователей, когортного анализа, результатов A/B тестов.
*   **Сегментация аудитории:** Создание статических и динамических сегментов пользователей.
*   **Предиктивная аналитика:** Разработка и развертывание ML моделей (прогнозирование продаж, оттока, рекомендации, обнаружение мошенничества).
*   **Мониторинг производительности системы:** Агрегация и анализ метрик здоровья и производительности микросервисов платформы.

### 1.3. Основные Технологии
*   **Языки обработки данных:** Scala, Python, Java (согласно `project_technology_stack.md`).
*   **Языки API слоя:** Go, Java (согласно `project_technology_stack.md`).
*   **Фреймворки обработки данных:** Apache Spark (batch), Kafka Streams или Apache Flink (stream).
*   **Хранилища данных:**
    *   Аналитическая СУБД: ClickHouse.
    *   Data Lake / Сырые данные: S3-совместимое хранилище (например, MinIO).
    *   Метаданные и конфигурация: PostgreSQL.
*   **Брокер сообщений / Потоковая обработка событий:** Apache Kafka.
*   **Машинное обучение:** TensorFlow, PyTorch, Scikit-learn, MLflow.
*   **Типы API:** REST (Spring Boot для Java, Echo/Gin для Go), GraphQL (Apollo Server, Hasura).
*   **Визуализация данных:** Grafana, Apache Superset.
*   (Ссылки на `project_technology_stack.md`, `PACKAGE_STANDARDIZATION.md`, `project_glossary.md`)

### 1.4. Термины и Определения (Glossary)
*   См. `project_glossary.md`.
*   **Событие (Event):** Запись о действии пользователя или системы.
*   **KPI (Key Performance Indicator):** Ключевой показатель эффективности.
*   **Сегмент (Segment):** Группа пользователей с общими характеристиками.
*   **ML Модель (Machine Learning Model):** Модель для предиктивной аналитики.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Analytics Service использует архитектуру, оптимизированную для обработки больших данных (Big Data) и машинного обучения, сочетая элементы Lambda и Kappa архитектур для обеспечения как пакетной, так и потоковой обработки.
*   Ключевые компоненты: Сбор данных (Data Ingestion), Обработка данных (Data Processing - batch/stream), Хранение данных (Data Storage - DWH, Data Lake, Metadata Store), API доступа к данным (Data Access API) и Слой ML моделей (ML Layer).
*   Ниже представлена диаграмма верхнеуровневой архитектуры Analytics Service:

```mermaid
graph TD
    subgraph "Источники Данных (Другие Микросервисы)"
        direction LR
        MS1[Сервис Аккаунтов] --> KI{Kafka Input}
        MS2[Сервис Каталога] --> KI
        MS3[Сервис Платежей] --> KI
        MSn[Другие Сервисы] --> KI
    end

    subgraph "Analytics Service"
        direction TB
        KI --> Ingestion[Data Ingestion Layer]

        Ingestion --> RawDataLake[Data Lake (S3 - сырые данные)]
        Ingestion --> StreamProc[Stream Processing (Kafka Streams/Flink)]

        StreamProc --> RealtimeMetrics[Real-time Metrics (ClickHouse/Redis)]
        StreamProc --> ProcessedEventsTopic{Kafka Processed Events}

        RawDataLake --> BatchProc[Batch Processing (Spark ETL)]
        ProcessedEventsTopic --> BatchProc

        BatchProc --> DWH[Data Warehouse (ClickHouse - агрегаты, витрины)]
        BatchProc --> MLDataPrep[ML Data Preparation]

        MLDataPrep --> MLLayer[ML Layer (Model Training & Serving)]
        MLLayer --> PredictionsAPI[Predictions API]
        MLLayer --> ModelMetadata[ML Model Metadata (PostgreSQL/MLflow)]


        DWH --> DataAPI[Data Access API (REST/GraphQL)]
        RealtimeMetrics --> DataAPI
        PredictionsAPI --> DataAPI

        MetadataDB[Metadata Store (PostgreSQL)] <--> DataAPI
        MetadataDB <--> BatchProc
        MetadataDB <--> StreamProc
        MetadataDB <--> MLLayer


        DataAPI --> Consumers[Потребители API (Админ-панель, Разработчики, BI-инструменты)]

        subgraph "Хранилища Данных"
            direction LR
            RawDataLake
            DWH
            RealtimeMetrics
            MetadataDB
            ModelMetadata
        end
    end

    classDef dataFlow fill:#e6f3ff,stroke:#007bff,stroke-width:2px;
    classDef component fill:#d4edda,stroke:#28a745,stroke-width:2px;
    classDef datastore fill:#f8d7da,stroke:#dc3545,stroke-width:2px;
    classDef api fill:#fff3cd,stroke:#ffc107,stroke-width:2px;
    classDef external fill:#e2e3e5,stroke:#6c757d,stroke-width:2px;

    class KI,ProcessedEventsTopic dataFlow;
    class Ingestion,StreamProc,BatchProc,MLDataPrep,MLLayer component;
    class RawDataLake,DWH,RealtimeMetrics,MetadataDB,ModelMetadata datastore;
    class DataAPI,PredictionsAPI api;
    class MS1,MS2,MS3,MSn,Consumers external;
```

### 2.2. Слои Сервиса / Компоненты

#### 2.2.1. Data Ingestion Layer (Слой Приема Данных)
*   Ответственность: Прием событий из Kafka от всех микросервисов платформы. Валидация схем событий (базовая). Сохранение сырых данных в Data Lake (S3).
*   Ключевые компоненты/модули: Kafka Consumers, коннекторы к S3.

#### 2.2.2. Data Processing Layer (Слой Обработки Данных)
*   Ответственность: Трансформация, агрегация, обогащение, очистка данных. Расчет метрик.
*   Ключевые компоненты/модули:
    *   Stream Processing Engine (Kafka Streams/Flink): Для обработки данных в реальном времени и расчета real-time метрик.
    *   Batch Processing Engine (Spark): Для сложных пакетных расчетов, ETL процессов, подготовки данных для ML.
    *   Загрузка обработанных данных в DWH (ClickHouse).

#### 2.2.3. Data Storage Layer (Слой Хранения Данных)
*   Ответственность: Хранение сырых, обработанных данных и метаданных.
*   Ключевые компоненты/модули:
    *   Data Lake (S3): Хранение всех сырых событий.
    *   Data Warehouse (ClickHouse): Хранение агрегированных данных, витрин данных, метрик для быстрого анализа.
    *   Metadata Store (PostgreSQL): Хранение схем данных, конфигураций пайплайнов, определений отчетов, метаданных ML моделей.

#### 2.2.4. Data Access & API Layer (Слой Доступа к Данным и API)
*   Ответственность: Предоставление доступа к данным и результатам анализа через API.
*   Ключевые компоненты/модули:
    *   REST API (Go/Java): Для запроса метрик, отчетов, сегментов.
    *   GraphQL API: Для гибких запросов к данным (будет рассмотрено в будущем).
    *   WebSocket API: Для стриминга real-time метрик (будет рассмотрено в будущем).

#### 2.2.5. Machine Learning Layer (Слой Машинного Обучения)
*   Ответственность: Тренировка, развертывание и обслуживание ML моделей. Предоставление API для получения прогнозов.
*   Ключевые компоненты/модули: MLflow (управление моделями), TensorFlow/PyTorch/Scikit-learn (библиотеки), API для прогнозов (может быть частью Data Access API или отдельным).

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL (через API Gateway):** `/api/v1/analytics`
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway или самим сервисом). Используется для доступа администраторов, разработчиков (к статистике по их продуктам), маркетологов и других внутренних пользователей. Возможно использование API ключей для межсервисного взаимодействия.
*   **Авторизация:** На основе ролей и, возможно, прав доступа к конкретным наборам данных/метрикам (например, `platform_admin` для общей статистики, `game_developer` для статистики по своим играм, `marketing_manager` для доступа к сегментам).
*   **Формат ответа об ошибке (стандартный):**
    ```json
    {
      "errors": [
        {
          "status": "4XX/5XX",
          "code": "ERROR_CODE",
          "title": "Краткое описание ошибки",
          "detail": "Полное описание ошибки с контекстом.",
          "source": { "pointer": "/data/attributes/field_name", "parameter": "query_param_name" } // Опционально
        }
      ]
    }
    ```
*   (Общие принципы см. `project_api_standards.md`)

#### 3.1.1. Ресурс: Метрики (Metrics)
*   **`GET /metrics`**
    *   Описание: Получение списка доступных определений метрик.
    *   Query параметры:
        *   `tag` (string, опционально): Фильтр по тегу метрики.
        *   `owner_service` (string, опционально): Фильтр по сервису-владельцу.
        *   `is_realtime` (boolean, опционально): Фильтр по поддержке real-time.
        *   `page` (integer, опционально, default: 1): Номер страницы.
        *   `per_page` (integer, опционально, default: 20): Количество элементов на странице.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "metricDefinition",
              "id": "metric-def-uuid-dau",
              "attributes": {
                "name": "daily_active_users",
                "display_name": "Дневная Активная Аудитория",
                "description": "Количество уникальных пользователей, активных за сутки.",
                "metric_type": "counter",
                "granularity": "daily",
                "dimensions": ["country", "platform_type"],
                "tags": ["core_kpi", "user_activity"]
              },
              "links": {
                "self": "/api/v1/analytics/metrics/definitions/daily_active_users"
              }
            }
          ],
          "meta": { "total_items": 50, "current_page": 1, "per_page": 20 },
          "links": { "next": "/api/v1/analytics/metrics?page=2" }
        }
        ```
    *   Требуемые права доступа: `platform_admin`, `game_developer`, `marketing_manager`.
*   **`GET /metrics/definitions/{metric_name}`**
    *   Описание: Получение детального определения конкретной метрики.
    *   Пример ответа (Успех 200 OK): (Структура аналогична элементу из `GET /metrics`)
    *   Требуемые права доступа: `platform_admin`, `game_developer`, `marketing_manager`.
*   **`GET /metrics/values/{metric_name}`**
    *   Описание: Получение рассчитанных значений для конкретной метрики.
    *   Query параметры:
        *   `start_date` (date, YYYY-MM-DD, обязательно): Начало периода.
        *   `end_date` (date, YYYY-MM-DD, обязательно): Конец периода.
        *   `dimensions` (string, опционально): Список измерений для группировки через запятую (например, `country,platform_type`).
        *   `filters` (string, опционально): Фильтры в формате `dimension_name=value` (например, `game_id=game-uuid-123&country=RU`). URL-кодировать значения.
        *   `granularity` (enum: `daily`, `hourly`, `monthly`, `raw`, обязательно): Гранулярность данных.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "metricValues",
            "id": "daily_active_users-20240101-20240131",
            "attributes": {
              "metric_name": "daily_active_users",
              "granularity": "daily",
              "start_date": "2024-01-01",
              "end_date": "2024-01-31",
              "values": [
                { "date": "2024-01-01", "value": 10500, "dimensions": {"country": "RU", "platform_type": "pc"} },
                { "date": "2024-01-01", "value": 5200, "dimensions": {"country": "US", "platform_type": "mobile"} },
                { "date": "2024-01-02", "value": 10650, "dimensions": {"country": "RU", "platform_type": "pc"} }
                // ...
              ]
            }
          },
          "meta": { "total_points": 60 }
        }
        ```
    *   Пример ответа (Ошибка 400 Bad Request - неверный параметр):
        ```json
        {
          "errors": [
            {
              "status": "400",
              "code": "INVALID_QUERY_PARAMETER",
              "title": "Неверный параметр запроса",
              "detail": "Параметр 'granularity' имеет недопустимое значение 'weekly'. Допустимые значения: daily, hourly, monthly, raw.",
              "source": { "parameter": "granularity" }
            }
          ]
        }
        ```
    *   Требуемые права доступа: `platform_admin`, `game_developer` (с ограничениями по доступу к данным, например, только по своим `game_id`).

#### 3.1.2. Ресурс: Отчеты (Reports)
*   **`GET /reports/definitions`**
    *   Описание: Получение списка доступных определений отчетов.
    *   Query параметры: `tag`, `owner_admin_id`, `page`, `per_page`.
    *   Пример ответа (Успех 200 OK): (Структура аналогична `GET /metrics/definitions`, но с полями `ReportDefinition`)
    *   Требуемые права доступа: `platform_admin`, `game_developer`, `marketing_manager`.
*   **`POST /reports/instances`**
    *   Описание: Запрос на генерацию нового экземпляра отчета (асинхронная операция).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "reportInstanceRequest",
            "attributes": {
              "report_definition_id": "report-def-uuid-sales",
              "parameters": {
                "month": "2024-02",
                "game_genre": "RPG"
              },
              "output_format": "csv"
            }
          }
        }
        ```
    *   Пример ответа (Успех 202 Accepted):
        ```json
        {
          "data": {
            "type": "reportInstance",
            "id": "instance-uuid-xyz",
            "attributes": {
              "report_definition_id": "report-def-uuid-sales",
              "status": "requested",
              "requested_at": "2024-03-15T12:30:00Z"
            },
            "links": {
              "status": "/api/v1/analytics/reports/instances/instance-uuid-xyz"
            }
          }
        }
        ```
    *   Требуемые права доступа: `platform_admin`, `game_developer`.
*   **`GET /reports/instances/{instance_id}`**
    *   Описание: Получение статуса генерации экземпляра отчета и ссылки на скачивание, если готов.
    *   Пример ответа (Успех 200 OK - отчет готов):
        ```json
        {
          "data": {
            "type": "reportInstance",
            "id": "instance-uuid-xyz",
            "attributes": {
              "report_definition_id": "report-def-uuid-sales",
              "status": "completed",
              "requested_at": "2024-03-15T12:30:00Z",
              "completed_at": "2024-03-15T12:35:00Z",
              "output_format": "csv",
              "file_size_bytes": 102400
            },
            "links": {
              "download": "/api/v1/analytics/reports/instances/instance-uuid-xyz/download"
            }
          }
        }
        ```
    *   Требуемые права доступа: `platform_admin`, `game_developer`.
*   **`GET /reports/instances/{instance_id}/download`**
    *   Описание: Загрузка сгенерированного отчета. Ответ будет файлом.
    *   Требуемые права доступа: `platform_admin`, `game_developer`.

#### 3.1.3. Ресурс: Сегменты (Segments)
*   **`GET /segments/definitions`**
    *   Описание: Получение списка определений пользовательских сегментов.
    *   Query параметры: `tag`, `segment_type` (`dynamic`/`static`), `page`, `per_page`.
    *   Требуемые права доступа: `platform_admin`, `marketing_manager`.
*   **`POST /segments/definitions`**
    *   Описание: Создание нового определения пользовательского сегмента.
    *   Тело запроса: (Содержит поля `SegmentDefinition`, особенно `criteria`)
        ```json
        {
          "data": {
            "type": "segmentDefinition",
            "attributes": {
              "name": "active_rpg_players_last_month",
              "display_name": "Активные RPG игроки за последний месяц",
              "segment_type": "dynamic",
              "refresh_schedule": "daily",
              "criteria": {
                "type": "AND",
                "conditions": [
                  {"field": "user.genres_played", "operator": "contains", "value": "RPG"},
                  {"field": "user.last_activity_date", "operator": ">=", "value": "NOW() - INTERVAL '30 days'"}
                ]
              }
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Возвращает созданное определение сегмента)
    *   Требуемые права доступа: `platform_admin`, `marketing_manager`.
*   **`GET /segments/{segment_id}/users-count`**
    *   Описание: Получение текущего количества пользователей в сегменте (для динамических сегментов может инициировать пересчет, если данные устарели).
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "segmentUserCount",
            "id": "segment-uuid-rpg-active",
            "attributes": {
              "user_count": 15230,
              "last_calculated_at": "2024-03-15T10:00:00Z"
            }
          }
        }
        ```
    *   Требуемые права доступа: `platform_admin`, `marketing_manager`.

#### 3.1.4. Ресурс: Предиктивная Аналитика (Predictions)
*   **`GET /predictions/models`**
    *   Описание: Получение списка доступных ML моделей и их метаданных.
    *   Query параметры: `status` (`development`, `staging`, `production`), `tag`.
    *   Требуемые права доступа: `platform_admin`, `data_scientist`.
*   **`POST /predictions/{model_name_or_id}/predict`**
    *   Описание: Запрос прогноза от конкретной ML модели.
    *   Тело запроса: Входные данные для модели в формате JSON. Структура зависит от модели.
        ```json
        // Пример для модели предсказания оттока
        {
          "data": {
            "type": "predictionRequest",
            "attributes": {
              "model_version": "v1.2", // Опционально, если не указано - используется последняя production версия
              "features": {
                "user_id": "user-uuid-abc",
                "days_since_last_login": 15,
                "total_sessions_last_30d": 5,
                "avg_session_duration_minutes": 25
                // ... другие фичи
              }
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "predictionResult",
            "id": "pred-uuid-xyz",
            "attributes": {
              "model_name": "churn_prediction",
              "model_version": "v1.2",
              "prediction": {
                "user_id": "user-uuid-abc",
                "churn_probability": 0.78,
                "predicted_label": true // true - склонен к оттоку
              },
              "timestamp": "2024-03-15T14:00:00Z"
            }
          }
        }
        ```
    *   Требуемые права доступа: `platform_admin` или специфические сервисные роли (например, `recommendation_service_role`).

### 3.2. GraphQL API
*   **Эндпоинт:** `/api/v1/analytics/graphql` (гипотетический)
*   **Описание:** В настоящее время GraphQL API не планируется. Приоритет отдан REST API для предоставления данных. Возможность добавления GraphQL будет рассмотрена в будущем при наличии соответствующих запросов от потребителей на более гибкие и кастомные запросы к данным.
*   **Пример возможного запроса (иллюстративный):**
    ```graphql
    # query GetDauAndRevenue {
    #   metric(name: "daily_active_users", startDate: "2024-01-01", endDate: "2024-01-07") {
    #     date
    #     value
    #     dimensions { country }
    #   }
    #   revenue(period: MONTHLY, year: 2024, month: 1) {
    #     totalAmount
    #     averagePerUser
    #   }
    # }
    ```

### 3.3. WebSocket API (для стриминга метрик)
*   **Эндпоинт:** `/api/v1/analytics/ws/streaming` (гипотетический)
*   **Описание:** В настоящее время WebSocket API для потоковой передачи данных не планируется. Потребности в real-time стриминге метрик будут оцениваться отдельно. Альтернативой может быть использование Server-Sent Events (SSE) через REST API или периодический опрос для real-time дашбордов, если задержка в несколько секунд приемлема.
*   **Пример возможного сообщения (иллюстративный, если бы использовался WebSocket):**
    ```json
    // Клиент подписывается на: { "action": "subscribe", "metric": "realtime_active_users_pc" }
    // Сервер отправляет обновления:
    // {
    //   "metric": "realtime_active_users_pc",
    //   "timestamp": "2024-03-15T14:05:30Z",
    //   "value": 12345,
    //   "dimensions": { "platform": "pc" }
    // }
    ```

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности/Структуры Данных

*   **`Event` (Событие)**: Базовая единица информации, отражающая действие пользователя или системы.
    *   **Соответствие стандарту**: Рекомендуется придерживаться спецификации [CloudEvents](https://cloudevents.io/) для унификации структуры событий по всей платформе (см. `project_api_standards.md`).
    *   **Ключевые атрибуты (примеры из CloudEvents)**:
        *   `id` (String, UUID): Уникальный идентификатор события. Обязательный.
        *   `source` (URI): Источник события (например, имя микросервиса). Обязательный.
        *   `specversion` (String): Версия спецификации CloudEvents. Обязательный (например, "1.0").
        *   `type` (String): Тип события в формате reverse-DNS (например, `com.example-platform.user.created`). Обязательный.
        *   `datacontenttype` (String): MIME-тип данных в `data` (например, `application/json`). Опциональный.
        *   `dataschema` (URI): Ссылка на схему данных события. Опциональный.
        *   `subject` (String): Идентификатор субъекта события (например, ID пользователя или игры). Опциональный.
        *   `time` (Timestamp): Время возникновения события. Опциональный, но крайне рекомендуемый.
        *   `data` (Object/Binary): Полезная нагрузка события, специфичная для `type`. Опциональный.
    *   **Хранение**:
        *   Сырые события: Data Lake (S3) в формате JSON или Avro.
        *   Обработанные/агрегированные события: Таблицы фактов в DWH (ClickHouse).

*   **`MetricDefinition` (Определение Метрики)**: Описывает, как рассчитывается и интерпретируется метрика. Хранится в PostgreSQL.
    *   `id` (UUID): Уникальный идентификатор определения метрики.
    *   `name` (VARCHAR(255)): Уникальное имя метрики (например, `daily_active_users`, `average_session_duration`).
    *   `display_name` (VARCHAR(255)): Человекочитаемое имя метрики (например, "Дневная Активная Аудитория").
    *   `description` (TEXT): Подробное описание метрики, включая бизнес-смысл.
    *   `metric_type` (VARCHAR(50)): Тип метрики (`counter`, `gauge`, `histogram`, `timer`).
    *   `calculation_method` (VARCHAR(50)): Метод расчета (`sum`, `average`, `count_distinct`, `percentile_95`).
    *   `source_event_type` (VARCHAR(255)): Тип события, на основе которого рассчитывается метрика (если применимо).
    *   `value_field` (VARCHAR(255)): Поле в событии/таблице, содержащее значение для расчета (если применимо).
    *   `filters` (JSONB): Предопределенные фильтры, применяемые при расчете.
    *   `dimensions` (JSONB): Список доступных измерений/срезов для этой метрики (например, `["country", "game_id", "platform_type"]`).
    *   `granularity` (VARCHAR(50)): Типичная гранулярность агрегации (`daily`, `hourly`, `monthly`).
    *   `unit` (VARCHAR(50)): Единица измерения (например, `users`, `seconds`, `RUB`).
    *   `is_realtime` (BOOLEAN): Поддерживается ли расчет в реальном времени.
    *   `owner_service` (VARCHAR(100)): Сервис-владелец или основной потребитель метрики.
    *   `tags` (JSONB): Теги для категоризации и поиска.

*   **`ReportDefinition` (Определение Отчета)**: Описывает структуру, параметры и метод генерации отчета. Хранится в PostgreSQL.
    *   `id` (UUID): Уникальный идентификатор определения отчета.
    *   `name` (VARCHAR(255)): Уникальное имя отчета (например, `monthly_sales_summary`, `user_retention_by_cohort`).
    *   `display_name` (VARCHAR(255)): Человекочитаемое имя отчета.
    *   `description` (TEXT): Описание отчета.
    *   `source_type` (VARCHAR(50)): Тип источника данных (`sql_query`, `pre_aggregated_metrics`, `external_api`).
    *   `source_query_or_config` (TEXT): SQL-запрос к DWH или конфигурация для другого источника.
    *   `parameters` (JSONB): Определение параметров отчета (например, `{"name": "period_start", "type": "date", "required": true}`).
    *   `default_schedule` (VARCHAR(50)): Расписание генерации по умолчанию (например, `daily_03_00_utc`, `monday_weekly`).
    *   `output_formats` (JSONB): Доступные форматы вывода (например, `["csv", "pdf", "json"]`).
    *   `owner_admin_id` (UUID): Администратор, ответственный за отчет.

*   **`SegmentDefinition` (Определение Сегмента)**: Описывает критерии для группировки пользователей. Хранится в PostgreSQL.
    *   `id` (UUID): Уникальный идентификатор определения сегмента.
    *   `name` (VARCHAR(255)): Уникальное имя сегмента (например, `active_players_last_7_days`, `whales_paid_over_10000`).
    *   `display_name` (VARCHAR(255)): Человекочитаемое имя сегмента.
    *   `description` (TEXT): Описание сегмента.
    *   `criteria` (JSONB): Критерии для включения пользователя в сегмент (например, основанные на атрибутах пользователя, его действиях, покупках). Пример: `{"type": "AND", "conditions": [{"field": "total_payments_sum", "operator": ">=", "value": 10000}, {"field": "last_activity_date", "operator": ">=", "value": "NOW() - INTERVAL '30 days'"}]}`.
    *   `segment_type` (VARCHAR(50)): Тип сегмента (`dynamic` - пересчитывается регулярно, `static` - разовый снимок).
    *   `refresh_schedule` (VARCHAR(50)): Расписание обновления для динамических сегментов.

*   **`MLModelMetadata` (Метаданные ML Модели)**: Информация о модели машинного обучения. Может храниться в PostgreSQL или специализированном инструменте типа MLflow.
    *   `id` (UUID): Уникальный идентификатор метаданных модели.
    *   `model_name` (VARCHAR(255)): Имя модели (например, `churn_prediction_v1`, `game_recommendation_user_cf`).
    *   `model_version` (VARCHAR(50)): Версия модели.
    *   `description` (TEXT): Описание модели, ее назначение.
    *   `algorithm` (VARCHAR(100)): Используемый алгоритм (например, `RandomForestClassifier`, `CollaborativeFiltering`).
    *   `hyperparameters` (JSONB): Гиперпараметры, использованные при обучении.
    *   `training_dataset_ref` (VARCHAR(255)): Ссылка на датасет, использованный для обучения (например, путь в S3 или ID в системе версионирования данных).
    *   `performance_metrics` (JSONB): Метрики качества модели (например, `{"auc": 0.85, "f1_score": 0.78}`).
    *   `artifact_path` (VARCHAR(255)): Путь к артефакту модели (например, в MLflow или S3).
    *   `deployment_status` (VARCHAR(50)): Статус развертывания (`development`, `staging`, `production`, `archived`).
    *   `created_at` (TIMESTAMPTZ): Дата создания записи о модели.
    *   `trained_at` (TIMESTAMPTZ): Дата обучения модели.

*   **`DataPipelineRun` (Запуск Пайплайна Данных)**: Информация о выполнении ETL/ELT пайплайна. Хранится в PostgreSQL.
    *   `id` (UUID): Уникальный идентификатор запуска.
    *   `pipeline_name` (VARCHAR(255)): Имя пайплайна (например, `daily_user_aggregation`, `events_to_clickhouse_stream`).
    *   `start_time` (TIMESTAMPTZ): Время начала выполнения.
    *   `end_time` (TIMESTAMPTZ): Время завершения выполнения.
    *   `status` (VARCHAR(50)): Статус (`running`, `completed`, `failed`, `skipped`).
    *   `parameters` (JSONB): Параметры, с которыми был запущен пайплайн.
    *   `logs_summary` (TEXT): Краткая сводка логов или ссылка на полные логи.
    *   `processed_records_count` (BIGINT): Количество обработанных записей.

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL (Metadata Store) - ERD Диаграмма

```mermaid
erDiagram
    METRIC_DEFINITIONS {
        UUID id PK
        VARCHAR name UK
        VARCHAR display_name
        TEXT description
        VARCHAR metric_type
        VARCHAR calculation_method
        VARCHAR source_event_type
        VARCHAR value_field
        JSONB filters
        JSONB dimensions
        VARCHAR granularity
        VARCHAR unit
        BOOLEAN is_realtime
        VARCHAR owner_service
        JSONB tags
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    REPORT_DEFINITIONS {
        UUID id PK
        VARCHAR name UK
        VARCHAR display_name
        TEXT description
        VARCHAR source_type
        TEXT source_query_or_config
        JSONB parameters
        VARCHAR default_schedule
        JSONB output_formats
        UUID owner_admin_id FK "nullable"
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    REPORT_INSTANCES {
        UUID id PK
        UUID report_definition_id FK
        TIMESTAMPTZ generation_requested_at
        TIMESTAMPTZ generation_started_at
        TIMESTAMPTZ generation_completed_at
        VARCHAR status
        JSONB parameters_used
        VARCHAR output_format
        VARCHAR file_path_s3 "nullable"
        TEXT error_message "nullable"
        UUID requested_by_user_id FK "nullable"
    }

    SEGMENT_DEFINITIONS {
        UUID id PK
        VARCHAR name UK
        VARCHAR display_name
        TEXT description
        JSONB criteria
        VARCHAR segment_type
        VARCHAR refresh_schedule "nullable"
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
        UUID created_by_admin_id FK "nullable"
    }

    ML_MODEL_METADATA {
        UUID id PK
        VARCHAR model_name
        VARCHAR model_version
        TEXT description
        VARCHAR algorithm
        JSONB hyperparameters
        VARCHAR training_dataset_ref
        JSONB performance_metrics
        VARCHAR artifact_path
        VARCHAR deployment_status
        TIMESTAMPTZ created_at
        TIMESTAMPTZ trained_at
        UUID registered_by_user_id FK "nullable"
    }

    DATA_PIPELINE_RUNS {
        UUID id PK
        VARCHAR pipeline_name
        TIMESTAMPTZ start_time
        TIMESTAMPTZ end_time
        VARCHAR status
        JSONB parameters
        TEXT logs_summary
        BIGINT processed_records_count
    }

    ADMIN_USERS { # Предполагается наличие таблицы из Admin Service или аналогичной
        UUID id PK
        VARCHAR username
    }

    USERS { # Предполагается наличие таблицы пользователей платформы
        UUID id PK
        VARCHAR username
    }

    REPORT_DEFINITIONS ||--o{ REPORT_INSTANCES : "defines"
    REPORT_DEFINITIONS }o--|| ADMIN_USERS : "owned_by"
    REPORT_INSTANCES }o--|| USERS : "requested_by"
    SEGMENT_DEFINITIONS }o--|| ADMIN_USERS : "created_by"
    ML_MODEL_METADATA }o--|| USERS : "registered_by (user/admin)"


    ADMIN_USERS ||--|{ REPORT_DEFINITIONS : "can own"
    USERS ||--|{ REPORT_INSTANCES : "can request"
    ADMIN_USERS ||--|{ SEGMENT_DEFINITIONS : "can create"
    USERS ||--|{ ML_MODEL_METADATA : "can register"
```

#### 4.2.2. PostgreSQL (Metadata Store) - DDL

```sql
-- Расширение для UUID, если не создано
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE metric_definitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    metric_type VARCHAR(50) NOT NULL, -- counter, gauge, histogram, timer
    calculation_method VARCHAR(50) NOT NULL, -- sum, average, count_distinct, percentile_95
    source_event_type VARCHAR(255), -- e.g., com.example.user.played_game
    value_field VARCHAR(255), -- Path to value in event data, e.g., data.duration_seconds
    filters JSONB, -- e.g., {"data.game_genre": "action"}
    dimensions JSONB, -- e.g., ["country", "game_id"]
    granularity VARCHAR(50) NOT NULL, -- daily, hourly, monthly
    unit VARCHAR(50), -- users, seconds, RUB
    is_realtime BOOLEAN NOT NULL DEFAULT FALSE,
    owner_service VARCHAR(100),
    tags JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE metric_definitions IS 'Определения метрик, используемых в системе аналитики.';

CREATE TABLE report_definitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    source_type VARCHAR(50) NOT NULL, -- sql_query, pre_aggregated_metrics, external_api
    source_query_or_config TEXT NOT NULL, -- SQL query or JSON config
    parameters JSONB, -- [{"name": "period_start", "type": "date", "required": true, "default_value": "yesterday"}]
    default_schedule VARCHAR(50), -- daily_03_00_utc, monday_weekly
    output_formats JSONB DEFAULT '["csv", "json"]'::jsonb,
    owner_admin_id UUID REFERENCES admin_users(id) ON DELETE SET NULL, -- Ссылка на таблицу admin_users
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE report_definitions IS 'Определения отчетов, которые могут быть сгенерированы.';

CREATE TABLE report_instances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_definition_id UUID NOT NULL REFERENCES report_definitions(id) ON DELETE CASCADE,
    generation_requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    generation_started_at TIMESTAMPTZ,
    generation_completed_at TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL, -- requested, generating, completed, failed
    parameters_used JSONB,
    output_format VARCHAR(20),
    file_path_s3 VARCHAR(1024),
    error_message TEXT,
    requested_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL -- Ссылка на таблицу users
);
COMMENT ON TABLE report_instances IS 'Экземпляры сгенерированных отчетов.';
CREATE INDEX idx_report_instances_status ON report_instances(status);
CREATE INDEX idx_report_instances_definition_id ON report_instances(report_definition_id);


CREATE TABLE segment_definitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    criteria JSONB NOT NULL, -- {"type": "AND", "conditions": [{"field": "total_payments_sum", "operator": ">=", "value": 10000}]}
    segment_type VARCHAR(50) NOT NULL DEFAULT 'dynamic', -- dynamic, static
    refresh_schedule VARCHAR(50), -- Для dynamic сегментов, например, 'daily'
    created_by_admin_id UUID REFERENCES admin_users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE segment_definitions IS 'Определения пользовательских сегментов.';

CREATE TABLE ml_model_metadata (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    model_name VARCHAR(255) NOT NULL,
    model_version VARCHAR(50) NOT NULL,
    description TEXT,
    algorithm VARCHAR(100),
    hyperparameters JSONB,
    training_dataset_ref VARCHAR(255), -- e.g., S3 path or DVC reference
    performance_metrics JSONB, -- {"auc": 0.85, "f1_score": 0.78}
    artifact_path VARCHAR(1024), -- Path to model file in MLflow or S3
    deployment_status VARCHAR(50) NOT NULL DEFAULT 'development', -- development, staging, production, archived
    registered_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Пользователь или админ, зарегистрировавший модель
    trained_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (model_name, model_version)
);
COMMENT ON TABLE ml_model_metadata IS 'Метаданные моделей машинного обучения.';

CREATE TABLE data_pipeline_runs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pipeline_name VARCHAR(255) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL, -- running, completed, failed, skipped
    parameters JSONB,
    logs_summary TEXT, -- Could be a link to full logs in S3 or ELK
    processed_records_count BIGINT
);
COMMENT ON TABLE data_pipeline_runs IS 'Информация о запусках пайплайнов обработки данных.';
CREATE INDEX idx_data_pipeline_runs_name_status ON data_pipeline_runs(pipeline_name, status);

-- ПРЕДПОЛАГАЕТСЯ НАЛИЧИЕ ТАБЛИЦ admin_users и users ИЗ ДРУГИХ СЕРВИСОВ ИЛИ ОБЩЕЙ БД
-- CREATE TABLE admin_users (id UUID PRIMARY KEY, username VARCHAR(100) ...);
-- CREATE TABLE users (id UUID PRIMARY KEY, username VARCHAR(100) ...);
```

#### 4.2.3. ClickHouse - Примеры DDL

*   **Таблица событий (примерная, для игровых событий):**
    ```sql
    CREATE TABLE default.game_events (
        event_id String DEFAULT generateUUIDv4(), -- Уникальный ID события
        event_type String,                       -- Тип события (например, 'game_start', 'level_complete', 'item_purchase')
        event_timestamp DateTime64(3, 'UTC'),    -- Время события с миллисекундами
        user_id String,                          -- ID пользователя
        session_id String,                       -- ID игровой сессии
        game_id String,                          -- ID игры
        platform Enum8('pc' = 1, 'mobile_android' = 2, 'mobile_ios' = 3, 'web' = 4),
        country_code FixedString(2),             -- Код страны пользователя

        -- Поля, специфичные для разных event_type (могут быть в JSON/Map или вынесены)
        -- Для 'level_complete':
        level_name Nullable(String),
        time_spent_seconds Nullable(UInt32),
        score Nullable(Int32),
        -- Для 'item_purchase':
        item_id Nullable(String),
        item_price Nullable(Decimal64(2)),
        currency_code Nullable(FixedString(3)),
        -- Общие/дополнительные параметры
        event_data Map(String, String),          -- Дополнительные неструктурированные параметры события

        received_at DateTime DEFAULT now()       -- Время получения события сервером аналитики
    )
    ENGINE = MergeTree()
    PARTITION BY toYYYYMM(event_timestamp)
    ORDER BY (game_id, event_type, event_timestamp, user_id)
    SETTINGS index_granularity = 8192;
    ```
    *Комментарий: Эта таблица предназначена для хранения "сырых" или минимально обработанных событий. Для реального использования схема будет значительно шире и сложнее, возможно, с использованием `Nested` структур для повторяющихся данных или `JSON` для гибкости.*

*   **Таблица агрегированных метрик (уточненная `metrics_daily`):**
    ```sql
    CREATE TABLE default.metrics_daily (
        metric_date Date,                         -- Дата, за которую рассчитана метрика
        metric_name String,                       -- Имя метрики (e.g., 'dau', 'registrations', 'total_revenue')

        -- Измерения (dimensions) - могут быть разными для разных метрик.
        -- Здесь примеры, лучше иметь отдельные таблицы или более гибкую структуру (например, Materialized Views).
        game_id String,                           -- ID игры (0 или '' для общеплатформенных метрик)
        country_code FixedString(2),              -- Код страны
        platform Enum8('pc' = 1, 'mobile' = 2, 'web' = 3, 'total' = 0), -- Платформа (total для всех)
        segment_id String,                        -- ID сегмента пользователей (0 или '' для всех)

        value Float64,                            -- Значение метрики

        -- Дополнительная информация
        calculation_time DateTime DEFAULT now()   -- Время расчета этой записи
    )
    ENGINE = SummingMergeTree(value) -- Или AggregatingMergeTree для более сложных агрегаций
    PARTITION BY toYYYYMM(metric_date)
    ORDER BY (metric_name, metric_date, game_id, country_code, platform, segment_id);
    -- При SummingMergeTree все столбцы кроме value и тех, что в ORDER BY, должны быть в нем.
    -- Если нужны другие агрегации (avg, uniq), используется AggregatingMergeTree с -State/-Merge функциями.
    ```

*   **Таблица измерений: Профили пользователей (пример):**
    ```sql
    CREATE TABLE default.user_profiles_dim (
        user_id String,                          -- ID пользователя
        registration_date Date,                  -- Дата регистрации
        first_payment_date Nullable(Date),       -- Дата первой покупки
        total_payments_amount Decimal64(2) DEFAULT 0, -- Общая сумма покупок
        country_code FixedString(2),             -- Код страны
        platform_last_used Enum8('pc' = 1, 'mobile' = 2, 'web' = 3),
        is_active Boolean,                       -- Активен ли пользователь (например, заходил в посл. 30 дней)
        -- Другие атрибуты пользователя, используемые для сегментации и анализа

        updated_at DateTime DEFAULT now()       -- Время последнего обновления записи в этой таблице
    )
    ENGINE = ReplacingMergeTree(updated_at) -- Для обновления атрибутов пользователя
    PRIMARY KEY user_id
    ORDER BY user_id;
    ```

*   **S3 (Data Lake)**: Хранение сырых событий (например, в формате Parquet или ORC, партицированных по дате и типу события), артефактов ML моделей, больших временных файлов для Spark/Flink джобов, экспорт больших отчетов. Структура директорий в S3 будет зависеть от конкретных пайплайнов.

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   Analytics Service в основном является потребителем событий. Однако он может публиковать:
*   **`analytics.report.generated.v1`**
    *   Описание: Отчет сгенерирован и доступен.
    *   Топик: `analytics.events`
    *   Структура Payload: `{ "report_instance_id": "...", "report_name": "...", "status": "success", "download_url": "...", "generated_at": "..." }`
    *   Потребители: Notification Service (уведомить пользователя), Developer Service (если отчет для разработчика).
*   **`analytics.segment.updated.v1`**
    *   Описание: Пользовательский сегмент обновлен (например, пересчитан состав).
    *   Топик: `analytics.events`
    *   Структура Payload: `{ "segment_id": "...", "user_count": 12345, "updated_at": "..." }`
    *   Потребители: Marketing tools, Notification Service (для кампаний).
*   **`analytics.alert.triggered.v1`**
    *   Описание: Сработал аналитический алерт (например, резкое падение DAU).
    *   Топик: `analytics.alerts`
    *   Структура Payload: `{ "alert_name": "DAU_drop", "severity": "critical", "details": {...}, "timestamp": "..." }`
    *   Потребители: Admin Service, Notification Service (для оповещения администраторов).

### 5.2. Потребляемые События (Consumed Events)
*   Analytics Service является основным потребителем событий от **всех** других микросервисов платформы.
*   **Топики:** `account.events`, `auth.events`, `catalog.events`, `library.events`, `payment.events`, `download.events`, `social.events`, `notification.events`, `developer.events`, `admin.events`.
*   **Формат событий:** CloudEvents JSON (согласно `project_api_standards.md`).
*   **Логика обработки:**
    *   Валидация схемы события.
    *   Сохранение сырого события в Data Lake (S3).
    *   Потоковая обработка (Kafka Streams/Flink) для:
        *   Обогащения данных (например, добавление геолокации по IP).
        *   Расчета real-time метрик.
        *   Обновления пользовательских профилей для аналитики.
        *   Обнаружения паттернов / аномалий.
    *   Пакетная обработка (Spark) для:
        *   ETL в DWH (ClickHouse).
        *   Пересчета сложных метрик и KPI.
        *   Тренировки ML моделей.
        *   Генерации отчетов.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Все микросервисы (через Kafka):** Основной источник данных для Analytics Service.
*   **Developer Service:** Предоставление аналитики и отчетов по играм для разработчиков (через API Analytics Service).
*   **Admin Service:** Предоставление общеплатформенных дашбордов и отчетов для администраторов (через API Analytics Service).
*   **Catalog Service / Recommendation Service (потенциально):** Может получать обработанные данные (например, метрики популярности, пользовательские предпочтения) для улучшения рекомендаций или ранжирования контента.
*   **Notification Service:** Для отправки уведомлений о готовности отчетов или срабатывании алертов.
*   **Auth Service:** Для аутентификации и авторизации запросов к API Analytics Service.

### 6.2. Внешние Системы
*   **S3-совместимое хранилище:** Для Data Lake и хранения артефактов.
*   **Внешние маркетинговые инструменты (потенциально):** Экспорт сегментов пользователей.
*   **Системы BI (потенциально):** Подключение к ClickHouse для построения кастомных отчетов (например, Tableau, PowerBI).

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `ANALYTICS_API_HTTP_PORT`, `ANALYTICS_API_GRPC_PORT` (если есть gRPC API).
*   `CLICKHOUSE_HOST`, `CLICKHOUSE_PORT`, `CLICKHOUSE_USER`, `CLICKHOUSE_PASSWORD`, `CLICKHOUSE_DATABASE`.
*   `POSTGRES_DSN_ANALYTICS_META`.
*   `S3_ENDPOINT_DATALAKE`, `S3_ACCESS_KEY_DATALAKE`, `S3_SECRET_KEY_DATALAKE`, `S3_BUCKET_DATALAKE`.
*   `KAFKA_BROKERS`.
*   `KAFKA_CONSUMER_GROUP_ANALYTICS`.
*   `SPARK_MASTER_URL` (если используется Spark).
*   `FLINK_JOBMANAGER_RPC_ADDRESS` (если используется Flink).
*   `MLFLOW_TRACKING_URI`.
*   `LOG_LEVEL`.
*   `AUTH_SERVICE_GRPC_ADDR` (для валидации токенов доступа к API).
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`: Эндпоинт Jaeger для экспорта трейсов OpenTelemetry. Пример: `http://jaeger-collector:14268/api/traces`.
*   `CORS_ALLOWED_ORIGINS_ANALYTICS_API`: Список разрешенных источников для CORS API аналитики. Пример: `https://bi.myplatform.com,http://localhost:3002`.

### 7.2. Файлы Конфигурации (если применимо)
*   Расположение: `configs/analytics_config.yaml` (путь внутри контейнера или при монтировании).
*   Этот файл используется для конфигураций, которые менее вероятно изменятся между окружениями или имеют сложную структуру. Переменные окружения могут переопределять значения из файла.
*   Пример `configs/analytics_config.yaml`:
    ```yaml
    server:
      http_port: ${ANALYTICS_API_HTTP_PORT:-8090}
      grpc_port: ${ANALYTICS_API_GRPC_PORT:-9091} # Если используется gRPC
      read_timeout_seconds: 60
      write_timeout_seconds: 60

    logging:
      level: ${LOG_LEVEL:-info} # error, warn, info, debug, trace
      format: "json"

    # Настройки подключения к основным хранилищам (часть может быть только в ENV VARS из-за секретности)
    # DSN/URI обычно лучше держать в ENV VARS для безопасности и гибкости.
    # Здесь могут быть нечувствительные параметры пулов соединений и т.д.
    clickhouse:
      # connection_string: ${CLICKHOUSE_DSN} # Предпочтительно из ENV
      default_database: ${CLICKHOUSE_DATABASE:-analytics_db}
      request_timeout_seconds: 120
      max_open_connections: 50

    postgresql_metadata:
      # dsn: ${POSTGRES_DSN_ANALYTICS_META} # Предпочтительно из ENV
      max_open_connections: 20

    s3_datalake:
      # endpoint: ${S3_ENDPOINT_DATALAKE} # Предпочтительно из ENV
      default_bucket: ${S3_BUCKET_DATALAKE:-platform-events-lake}
      # access_key: ${S3_ACCESS_KEY_DATALAKE} # Секрет, только ENV
      # secret_key: ${S3_SECRET_KEY_DATALAKE} # Секрет, только ENV

    kafka:
      # brokers: ${KAFKA_BROKERS} # Предпочтительно из ENV
      consumer_groups:
        events_processor_main: "analytics_events_main_consumer_group"
        realtime_metrics_updater: "analytics_realtime_metrics_consumer_group"
      producer_defaults:
        ack_mode: "all"
        compression_type: "snappy" # gzip, lz4

    # Конфигурация для Spark джобов (пример)
    spark_processing:
      default_master_url: ${SPARK_MASTER_URL:-local[*]} # Для локального запуска или указания мастера
      default_executor_memory: "2g"
      default_driver_memory: "1g"
      default_parallelism: 200
      jobs:
        daily_user_aggregation:
          app_name: "DailyUserAggregation"
          main_class: "com.example.analytics.spark.DailyUserAggregationJob"
          schedule: "0 2 * * *" # cron-like
          input_path_template: "s3a://${S3_BUCKET_DATALAKE:-platform-events-lake}/raw_events/date={yyyy-MM-dd}/"
          output_table_dwh: "user_daily_summary"
          extra_configs:
            spark.sql.shuffle.partitions: "100"
        # ... другие Spark джобы

    # Конфигурация для Flink джобов (пример)
    flink_processing:
      default_jobmanager_rpc_address: ${FLINK_JOBMANAGER_RPC_ADDRESS:-localhost:6123}
      default_parallelism: 4
      jobs:
        realtime_session_analyzer:
          app_name: "RealtimeSessionAnalyzer"
          main_class: "com.example.analytics.flink.RealtimeSessionAnalyzerJob"
          input_topic: "platform.user.activity.v1" # Kafka топик
          output_metrics_table: "realtime_session_metrics" # ClickHouse таблица
          checkpoint_interval_ms: 60000
          # ... другие Flink джобы

    # Конфигурация для MLflow (если используется для трекинга)
    mlflow:
      tracking_uri: ${MLFLOW_TRACKING_URI:-http://mlflow-server:5000}
      default_experiment_name: "Platform_Analytics_ML"

    api_settings:
      default_page_size: 20
      max_page_size: 100
      realtime_metrics_cache_ttl_seconds: 30
    ```
*   SQL файлы для создания витрин в ClickHouse и инициализации метаданных в PostgreSQL также являются частью конфигурации сервиса и должны версионироваться вместе с кодом.

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   Надежная обработка ошибок в пайплайнах приема и обработки данных (DLQ для Kafka, retry механизмы).
*   Мониторинг состояния джобов обработки данных.
*   Информативные ошибки для API запросов.

### 8.2. Распространенные Коды Ошибок (для API)
*   **`METRIC_NOT_FOUND`**
*   **`REPORT_GENERATION_FAILED`**
*   **`INVALID_QUERY_PARAMETERS`**
*   **`DATA_NOT_YET_AVAILABLE`** (для real-time или недавно загруженных данных)
*   (В дополнение к стандартным HTTP ошибкам)

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   JWT/API ключи для доступа к API.
*   Защищенный доступ к хранилищам данных и брокерам сообщений.

### 9.2. Авторизация
*   RBAC для доступа к API и данным (например, разработчики видят только агрегированную статистику по своим играм, администраторы видят все).
*   (Ссылка на `project_roles_and_permissions.md`)

### 9.3. Защита Данных
*   Анонимизация и псевдонимизация PII при обработке и хранении.
*   Соблюдение ФЗ-152.
*   Контроль доступа к Data Lake и DWH.
*   Шифрование чувствительных данных в покое и при передаче.
*   (Ссылка на `project_security_standards.md`)

### 9.4. Управление Секретами
*   Использование Kubernetes Secrets или Vault.

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   Dockerfiles для различных компонентов (API сервис, Spark/Flink джобы).
*   Helm-чарты/Kubernetes манифесты.
*   (Ссылка на `project_deployment_standards.md`)

### 10.2. Зависимости при Развертывании
*   ClickHouse, PostgreSQL, S3, Kafka.
*   MLflow (если используется).
*   Доступ к Kafka топикам от всех других сервисов.

### 10.3. CI/CD
*   Пайплайны для сборки API сервисов, джобов обработки данных, ML моделей.
*   Тестирование ETL пайплайнов, валидация данных.
*   (Ссылка на `project_deployment_standards.md`)

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   Формат: JSON.
*   Ключевые события: Статус приема данных, ошибки обработки, выполнение ETL джобов, запросы к API.
*   Интеграция: ELK/Loki.
*   (Ссылка на `project_observability_standards.md`)

### 11.2. Мониторинг
*   Метрики (Prometheus):
    *   Объем входящих/обработанных событий Kafka.
    *   Задержка обработки событий.
    *   Состояние и производительность Spark/Flink джобов.
    *   Производительность запросов к ClickHouse.
    *   Ошибки API Analytics Service.
    *   Актуальность данных в витринах.
*   Дашборды (Grafana): Мониторинг состояния пайплайнов данных, производительности DWH, использования ресурсов.
*   Алерты (AlertManager): Сбои в обработке данных, большая задержка данных, недоступность хранилищ.
*   (Ссылка на `project_observability_standards.md`)

### 11.3. Трассировка
*   Интеграция: OpenTelemetry, Jaeger (для API слоя и некоторых критичных частей обработки).
*   (Ссылка на `project_observability_standards.md`)

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**:
    *   **Сбор данных (Data Ingestion)**:
        *   Способность обрабатывать пиковую нагрузку до 50,000 событий/сек.
        *   Средняя задержка от получения события от Kafka до сохранения в Data Lake (S3): P95 < 500 мс.
    *   **Обработка данных (Data Processing)**:
        *   Пакетные ETL-процессы для формирования ежедневных агрегатов должны завершаться в течение 3 часов.
        *   Задержка потоковой обработки для ключевых real-time метрик (например, количество активных пользователей онлайн): P99 < 5 секунд от поступления события в Kafka до обновления метрики.
        *   Время пересчета основных пользовательских сегментов (динамических): не более 1 часа (для ежедневного пересчета).
    *   **API Доступа к данным**:
        *   Запросы к часто используемым агрегированным метрикам (например, DAU, MAU за период): P95 < 1000 мс.
        *   Запросы к сложным аналитическим отчетам (требующие агрегации на лету): P95 < 10 секунд.
        *   Запросы к API предиктивных моделей (например, вероятность оттока): P99 < 200 мс.
*   **Масштабируемость**:
    *   Горизонтальное масштабирование всех компонентов (слой приема, потоковая/пакетная обработка, API).
    *   Способность хранить и обрабатывать данные объемом до 1 Петабайта в DWH (ClickHouse) и до 5 Петабайт в Data Lake (S3) с прогнозируемым ростом на 50% в год.
    *   Поддержка до 100 одновременных аналитических запросов к API.
    *   Увеличение количества источников данных (микросервисов) на 20% не должно требовать кардинального изменения архитектуры.
*   **Надежность**:
    *   Гарантированная доставка событий в Data Lake: не менее 99.99% (без потерь).
    *   Устойчивость ETL-пайплайнов к сбоям отдельных узлов; наличие механизмов retry и восстановления.
    *   Доступность API Сервиса Аналитики: 99.9%.
    *   RTO (Recovery Time Objective) для DWH после сбоя: < 4 часов.
    *   RPO (Recovery Point Objective) для DWH: < 24 часов (для пакетных данных), < 1 часа (для real-time витрин, если есть).
*   **Актуальность данных (Data Freshness)**:
    *   Real-time метрики для дашбордов: обновляются с задержкой не более 1 минуты от события.
    *   Ежедневные агрегированные отчеты: доступны к 06:00 UTC за предыдущие сутки.
    *   Еженедельные агрегированные отчеты: доступны к понедельнику 09:00 UTC за предыдущую неделю.
    *   Данные для ML моделей: обновляются по расписанию, соответствующему циклу переобучения моделей (например, ежедневно или еженедельно).
*   **Точность данных (Data Accuracy)**:
    *   Расхождение между данными в операционных системах и DWH после ETL: < 0.1% для ключевых финансовых и пользовательских метрик.
    *   Механизмы валидации данных на различных этапах обработки.
*   **Сопровождаемость**:
    *   Покрытие кода тестами (unit, integration, e2e для пайплайнов): > 75%.
    *   Время развертывания новой версии API или пайплайна в production: < 30 минут (после прохождения всех тестов).
    *   Наличие документации по всем метрикам, отчетам и пайплайнам данных.

## 13. Приложения (Appendices) (Опционально)
*   TODO: Детальные схемы событий (если отличаются от CloudEvents), полные DDL для всех таблиц ClickHouse/PostgreSQL, развернутые примеры API запросов/ответов, схемы GraphQL (если будет реализовано) и другие приложения могут быть добавлены сюда или в отдельные специализированные документы.

---
*Этот документ является отправной точкой и должен регулярно обновляться по мере развития сервиса.*
