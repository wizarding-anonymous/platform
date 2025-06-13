# Спецификация Микросервиса: Notification Service (Сервис Уведомлений)

**Версия:** 1.1
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Notification Service платформы "Российский Аналог Steam".
*   **Роль в общей архитектуре платформы:** Notification Service отвечает за централизованную отправку и управление всеми видами уведомлений пользователям платформы. Это включает транзакционные уведомления, системные оповещения, маркетинговые рассылки и уведомления, инициированные другими микросервисами. Сервис абстрагирует сложность взаимодействия с различными провайдерами доставки (Email, Push, SMS) и каналами (In-App).
*   **Основные бизнес-задачи:**
    *   Обеспечение надежной и своевременной доставки уведомлений.
    *   Управление пользовательскими подписками и предпочтениями.
    *   Предоставление инструментов для создания и управления шаблонами.
    *   Поддержка маркетинговых кампаний.
    *   Сбор статистики по доставке и взаимодействию с уведомлениями.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Многоканальная доставка:** Поддержка различных каналов доставки уведомлений: Email, Push-уведомления (FCM для Android, APNS для iOS), SMS, In-App уведомления (через WebSocket Gateway).
*   **Управление шаблонами уведомлений:** Создание, редактирование, версионирование и хранение шаблонов уведомлений. Поддержка локализации шаблонов и использования переменных (плейсхолдеров) для персонализации.
*   **Отправка уведомлений:** API для одиночной и массовой отправки уведомлений. Возможность отложенной отправки и установки приоритетов для разных типов уведомлений.
*   **Управление пользовательскими предпочтениями:** Предоставление пользователям возможности настраивать, какие типы уведомлений и по каким каналам они хотят получать. Управление глобальными отписками.
*   **Управление маркетинговыми кампаниями:** Создание и управление маркетинговыми рассылками, включая сегментацию аудитории (интеграция с Analytics Service), планирование времени отправки, A/B тестирование (базовое).
*   **Отслеживание статусов и сбор статистики:** Отслеживание статусов доставки уведомлений (отправлено, доставлено, ошибка доставки, открыто, переход по ссылке). Сбор и агрегация статистики для анализа эффективности уведомлений и кампаний.
*   **Интеграция с провайдерами:** Гибкая конфигурация для подключения к различным внешним провайдерам отправки Email, SMS и Push-уведомлений. Поддержка нескольких провайдеров для одного канала с возможностью выбора или failover. Обработка обратной связи от провайдеров (например, отписки от Email, недоставленные Push).
*   **Управление устройствами пользователя:** Регистрация и управление токенами устройств для Push-уведомлений.

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+, согласно `../../../../project_technology_stack.md`).
*   **Веб-фреймворк (REST API):** Echo (`github.com/labstack/echo/v4`) или Gin (`github.com/gin-gonic/gin`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **RPC фреймворк (gRPC):** Опционально, для внутреннего высокопроизводительного взаимодействия (`google.golang.org/grpc`, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Очереди сообщений (основной способ получения запросов на отправку):** Apache Kafka (клиент `github.com/confluentinc/confluent-kafka-go` или `github.com/segmentio/kafka-go`, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Шаблонизатор:** Стандартные пакеты Go `text/template` и `html/template`.
*   **Базы данных:**
    *   PostgreSQL (версия 15+): Для хранения метаданных (шаблоны, кампании, пользовательские предпочтения, токены устройств, логи некоторых сообщений). Драйвер: GORM (`gorm.io/gorm`) с `gorm.io/driver/postgres` или `pgx` (`github.com/jackc/pgx/v5`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
    *   ClickHouse (версия 23.x+): Для хранения и агрегации больших объемов статистики по доставке и взаимодействию с уведомлениями. (согласно `../../../../project_technology_stack.md`).
*   **Кэширование:** Redis (версия 7.0+) для кэширования пользовательских предпочтений, шаблонов, счетчиков для rate limiting, временных данных кампаний. Клиент: `go-redis/redis` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Интеграции с внешними провайдерами:** Через их HTTP API (Email, SMS, Push).
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Логирование:** Zap (`go.uber.org/zap`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Мониторинг/Трассировка:** OpenTelemetry SDK, Prometheus client (`github.com/prometheus/client_golang`). (согласно `../../../../project_observability_standards.md`).
*   **Инфраструктура:** Docker, Kubernetes.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   **Уведомление (Notification):** Информационное сообщение, отправляемое пользователю или системе.
*   **Канал Доставки (Delivery Channel):** Способ, которым уведомление достигает пользователя (Email, Push, SMS, In-App).
*   **Шаблон (Template):** Предопределенная структура сообщения с плейсхолдерами для персонализации.
*   **Провайдер Доставки (Provider):** Внешний сервис, используемый для фактической отправки уведомлений по определенному каналу (например, SendGrid для Email, Firebase Cloud Messaging (FCM) для Push).
*   **Кампания (Notification Campaign):** Организованная рассылка одного или нескольких уведомлений целевой аудитории, обычно в маркетинговых целях.
*   **Токен Устройства (Device Token):** Уникальный идентификатор, используемый для отправки Push-уведомлений на конкретное мобильное устройство.
*   **In-App Уведомление:** Уведомление, отображаемое внутри клиентского приложения платформы.
*   Для других общих терминов см. `../../../../project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Notification Service построен с использованием событийно-ориентированной архитектуры (EDA), где большинство запросов на отправку уведомлений поступает через Kafka. Сервис также предоставляет REST API для управления и некоторых операций.
*   Ключевые компоненты включают: потребители Kafka для входящих запросов, оркестратор уведомлений, менеджеры шаблонов и предпочтений, диспетчеры для каждого канала доставки, и сборщик статистики.

**Диаграмма Компонентов:**
```mermaid
graph TD
    subgraph Notification Service
        direction LR
        API_REST[REST API (Управление)]

        subgraph CoreLogic [Ядро Сервиса]
            direction TB
            KafkaConsumer[Kafka Consumers (Запросы на отправку)]
            Orchestrator[Notification Orchestrator]
            TemplateManager[Template Manager]
            PreferenceManager[Preference Manager]
            CampaignManager[Campaign Manager]
            DeviceManager[Device Token Manager]
        end

        subgraph Dispatchers [Диспетчеры Каналов]
            direction TB
            EmailDispatcher[Email Dispatcher]
            PushDispatcher[Push Notification Dispatcher]
            SMSDispatcher[SMS Dispatcher]
            InAppDispatcher[In-App Dispatcher (to WebSocket Gateway)]
        end

        StatsCollector[Stats Collector]

        subgraph Persistence [Хранилища Данных]
            direction TB
            PostgresDB[PostgreSQL (Шаблоны, Предпочтения, Кампании, Устройства, Логи)]
            ClickHouseDB[ClickHouse (Статистика Доставки)]
            RedisCache[Redis (Кэш, Очереди Задач)]
        end

        API_REST --> PreferenceManager
        API_REST --> CampaignManager
        API_REST --> TemplateManager
        API_REST --> DeviceManager
        API_REST --> Orchestrator # Для прямой отправки через API

        KafkaConsumer --> Orchestrator
        Orchestrator --> TemplateManager
        Orchestrator --> PreferenceManager
        Orchestrator --> DeviceManager
        Orchestrator --> Dispatchers

        Dispatchers --> ExternalProviders[Внешние Провайдеры (Email, SMS, Push)]
        InAppDispatcher --> WebSocketGateway[WebSocket Gateway]

        Dispatchers -- Статусы доставки --> StatsCollector
        StatsCollector --> ClickHouseDB
        StatsCollector --> KafkaProducerStatus[Kafka Producer (Статусы)]

        TemplateManager --> PostgresDB
        PreferenceManager --> PostgresDB
        CampaignManager --> PostgresDB
        DeviceManager --> PostgresDB
        Orchestrator -- Логирование сообщений --> PostgresDB

        CoreLogic --> RedisCache
    end

    ExternalServices[Другие Микросервисы] -- Запросы через Kafka --> KafkaConsumer
    AdminUsers[Администраторы] -- Управление через UI --> API_REST

    classDef component fill:#d4edda,stroke:#28a745,color:#000
    classDef datastore fill:#f8d7da,stroke:#dc3545,color:#000
    classDef external fill:#e2e3e5,stroke:#6c757d,stroke-width:2px;

    class API_REST,KafkaConsumer,Orchestrator,TemplateManager,PreferenceManager,CampaignManager,DeviceManager,EmailDispatcher,PushDispatcher,SMSDispatcher,InAppDispatcher,StatsCollector component;
    class PostgresDB,ClickHouseDB,RedisCache datastore;
    class ExternalServices,AdminUsers,ExternalProviders,WebSocketGateway external;
    class KafkaProducerStatus component_minor;
```

### 2.2. Слои Сервиса (и Компоненты)
(Содержимое существующего раздела актуально).

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL:** `/api/v1/notifications`.
*   **Аутентификация:** JWT Bearer Token или API-ключи.
*   **Авторизация:** На основе ролей.
*   **Стандартный формат ответа об ошибке (согласно `../../../../project_api_standards.md`):**
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

#### 3.1.1. Управление Шаблонами Уведомлений
*   **`POST /templates`**
    *   Пример ответа (Ошибка 400 Validation Error - стандартизированный):
        ```json
        {
          "errors": [
            {
              "code": "VALIDATION_ERROR",
              "title": "Ошибка валидации шаблона",
              "detail": "Поле 'channel_type' должно быть одним из допустимых значений.",
              "source": { "pointer": "/data/attributes/channel_type" }
            }
          ]
        }
        ```
    (Остальное описание эндпоинтов `/templates` и `/templates/{template_id}` как в существующем документе).

#### 3.1.2. Управление Пользовательскими Предпочтениями
(Описания эндпоинтов `/preferences/users/{user_id}` и `PUT /preferences/users/{user_id}` как в существующем документе).

#### 3.1.3. Отправка Уведомлений (Прямая, для сервисов/админов)
(Описание эндпоинта `POST /send/direct` как в существующем документе).

### 3.2. gRPC API
*   В настоящее время предоставление собственного gRPC API не является приоритетом. Сервис взаимодействует с другими микросервисами преимущественно через Kafka для получения запросов на отправку уведомлений и через REST API для функций управления. Если в будущем возникнет потребность в высокопроизводительных синхронных вызовах к Notification Service от других внутренних сервисов (например, для немедленного получения статуса критического уведомления), разработка gRPC API будет рассмотрена.

### 3.3. WebSocket API (для In-App уведомлений)
(Содержимое существующего раздела актуально).

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`NotificationTemplate` (Шаблон Уведомления)** (Как в существующем документе)
*   **`NotificationMessageLog` (Запись об Отправленном Уведомлении)** (Как в существующем документе)
*   **`UserNotificationPreferences` (Предпочтения Пользователя по Уведомлениям)** (Как в существующем документе)
*   **`DeviceToken` (Токен Устройства для Push)** (Как в существующем документе)
*   **`NotificationCampaign` (Маркетинговая Кампания)**
    *   `id` (UUID): Уникальный идентификатор кампании.
    *   `name` (VARCHAR(255)): Название кампании.
    *   `description` (TEXT): Описание целей и деталей кампании.
    *   `target_segment_id` (VARCHAR(255)): Идентификатор сегмента аудитории (из Analytics Service или определенный локально).
    *   `template_id` (UUID, FK to NotificationTemplate): ID используемого шаблона.
    *   `status` (ENUM: `draft`, `scheduled`, `active`, `paused`, `completed`, `archived`): Статус кампании.
    *   `scheduled_at` (TIMESTAMPTZ, Nullable): Планируемое время начала.
    *   `started_at` (TIMESTAMPTZ, Nullable): Фактическое время начала.
    *   `completed_at` (TIMESTAMPTZ, Nullable): Время завершения.
    *   `created_by` (UUID): ID администратора/маркетолога, создавшего кампанию.
    *   `stats_summary` (JSONB): Краткая статистика по кампании (отправлено, доставлено, открыто).
    *   `created_at` (TIMESTAMPTZ), `updated_at` (TIMESTAMPTZ).
*   **`ProviderConfig` (Конфигурация Провайдера Доставки):** Данные конфигурации для внешних провайдеров (API ключи, URL эндпоинтов, специфичные настройки) **не хранятся в базе данных**. Они управляются через переменные окружения, которые инжектируются из Kubernetes Secrets или Vault, как описано в разделе "Конфигурация".

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL
**ERD Диаграмма (дополненная):**
```mermaid
erDiagram
    NOTIFICATION_TEMPLATES {
        UUID id PK
        VARCHAR name UK
        TEXT description
        VARCHAR channel_type
        VARCHAR default_language_code
        JSONB versions
        ARRAY_TEXT example_variables
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    USER_NOTIFICATION_PREFERENCES {
        UUID user_id PK "FK (User ID)"
        BOOLEAN global_unsubscribe_all
        JSONB channel_settings
        JSONB type_preferences
        TIMESTAMPTZ updated_at
    }
    DEVICE_TOKENS {
        UUID id PK
        UUID user_id FK "FK (User ID)"
        VARCHAR device_platform
        TEXT token UK
        TIMESTAMPTZ last_seen_at
        BOOLEAN is_active
        VARCHAR app_version
        TIMESTAMPTZ created_at
    }
    NOTIFICATION_CAMPAIGNS {
        UUID id PK
        VARCHAR name
        TEXT description
        VARCHAR target_segment_id
        UUID template_id FK
        VARCHAR status
        TIMESTAMPTZ scheduled_at
        TIMESTAMPTZ started_at
        TIMESTAMPTZ completed_at
        JSONB stats_summary
        UUID created_by "FK (Admin User)"
    }
    NOTIFICATION_MESSAGE_LOG {
        UUID id PK
        UUID user_id "FK (User ID, nullable)"
        UUID template_id "FK (nullable)"
        UUID campaign_id "FK (nullable)"
        VARCHAR channel_type
        VARCHAR provider_name
        TEXT recipient_address
        VARCHAR status
        TEXT status_details
        TIMESTAMPTZ requested_at
        TIMESTAMPTZ sent_at
        TIMESTAMPTZ delivered_at
    }

    USERS {
        note "From Auth/Account Service"
    }
    ADMIN_USERS {
        note "From Admin Service"
    }


    NOTIFICATION_TEMPLATES ||--o{ NOTIFICATION_CAMPAIGNS : "uses"
    USERS ||--|| USER_NOTIFICATION_PREFERENCES : "has"
    USERS ||--o{ DEVICE_TOKENS : "has_devices"
    USERS ||--o{ NOTIFICATION_MESSAGE_LOG : "receives"
    ADMIN_USERS ||--o{ NOTIFICATION_CAMPAIGNS : "created_by"
    NOTIFICATION_TEMPLATES ||--o{ NOTIFICATION_MESSAGE_LOG : "based_on_template"
    NOTIFICATION_CAMPAIGNS ||--o{ NOTIFICATION_MESSAGE_LOG : "part_of_campaign"
```

**DDL (PostgreSQL - дополнение для `notification_campaigns`):**
```sql
CREATE TABLE notification_campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    target_segment_id VARCHAR(255),
    template_id UUID NOT NULL REFERENCES notification_templates(id) ON DELETE RESTRICT,
    status VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'scheduled', 'active', 'paused', 'completed', 'archived')),
    scheduled_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_by UUID,
    stats_summary JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_notification_campaigns_status ON notification_campaigns(status);
CREATE INDEX idx_notification_campaigns_scheduled_at ON notification_campaigns(scheduled_at) WHERE status = 'scheduled';
COMMENT ON TABLE notification_campaigns IS 'Маркетинговые и информационные кампании уведомлений.';
```

#### 4.2.2. ClickHouse (Статистика и Долгосрочные Логи)
(Схема таблицы `notification_delivery_stats` остается как в существующем документе).

#### 4.2.3. Redis
(Описание структуры данных в Redis остается как в существующем документе).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Формат событий:** CloudEvents v1.0 JSON (согласно `../../../../project_api_standards.md`).
*   **Основной топик Kafka:** `com.platform.notification.events.v1`.

*   **`com.platform.notification.message.status.updated.v1`**
    *   `data` Payload: (Как в существующем документе, с корректным `type`).
*   **`com.platform.notification.campaign.status.changed.v1`**
    *   Описание: Статус маркетинговой кампании изменился.
    *   `data` Payload:
        ```json
        {
          "campaignId": "campaign-uuid-abc",
          "campaignName": "Весенняя распродажа - Анонс",
          "newStatus": "active",
          "previousStatus": "scheduled",
          "changeTimestamp": "2024-03-20T08:00:00Z",
          "updatedBy": "system"
        }
        ```
    *   Потребители: Admin Service, Analytics Service.

### 5.2. Потребляемые События (Consumed Events)
*   **Основной топик для запросов на отправку:** `com.platform.notification.send.request.v1`.
*   **`com.platform.notification.send.request.v1`**
    *   Описание: Запрос на отправку уведомления.
    *   `data` Payload (уточненный):
        ```json
        {
          "requestId": "req-uuid-unique",
          "recipients": [
            {
              "userId": "user-uuid-123",
              "emailTo": "override_user@example.com",
              "phoneTo": "+79991234567",
              "deviceTokens": [
                {"token": "fcm_token_xyz", "platform": "android_fcm"},
                {"token": "apns_token_abc", "platform": "ios_apns"}
              ]
            }
          ],
          "templateName": "order_confirmation",
          "languageCode": "ru-RU",
          "payloadVariables": {
            "orderId": "ORD-2024-03-18-001",
            "userName": "Иван Петров",
            "totalAmount": "1500.00 RUB",
            "itemsListHtml": "<ul><li>Игра 'Супер Гонки' - 1000 RUB</li><li>DLC 'Новые Трассы' - 500 RUB</li></ul>",
            "orderDetailsUrl": "https://myplatform.com/orders/ORD-2024-03-18-001"
          },
          "deliveryOptions": {
            "channels": ["email", "push_fcm"],
            "priority": "high",
            "sendAtUtc": null,
            "ttlSeconds": 3600
          },
          "correlationId": "order-uuid-xyz"
        }
        ```
(Остальное описание потребляемых событий как в существующем документе, с коррекцией имен событий).

## 6. Интеграции (Integrations)
(Содержимое существующего раздела актуально).
*   **Внешние провайдеры:** При выборе конкретных внешних провайдеров (Email, SMS) следует отдавать предпочтение российским сервисам для обеспечения соответствия локальным требованиям и потенциально лучшей доставляемости внутри РФ, если это является целевым рынком.

## 7. Конфигурация (Configuration)
(Содержимое существующего раздела YAML и описание переменных окружения в целом актуальны).
*   **Управление API ключами провайдеров:** Все API ключи (`EMAIL_PROVIDER_PRIMARY_API_KEY`, `SMS_PROVIDER_PRIMARY_API_KEY`, `FCM_SERVER_KEY`, `APNS_KEY_ID`, `APNS_TEAM_ID`) и другие секреты должны загружаться исключительно из переменных окружения, которые, в свою очередь, инжектируются из Kubernetes Secrets или HashiCorp Vault.

## 8. Обработка Ошибок (Error Handling)
(Содержимое существующего раздела актуально, форматы ошибок исправлены в разделе API).

## 9. Безопасность (Security)
(Содержимое существующего раздела в целом актуально).
*   **ФЗ-152 "О персональных данных":** Notification Service обрабатывает значительный объем ПДн: контактные данные пользователей (email, номера телефонов), токены устройств, а также содержимое уведомлений, которое может косвенно раскрывать ПДн или другую чувствительную информацию. Необходимо обеспечить строгие меры по защите этих данных, включая шифрование, управление доступом, логирование операций и получение согласий пользователей на получение уведомлений различных типов.
*   Ссылки на `../../../../project_security_standards.md` и `../../../../project_roles_and_permissions.md`.

## 10. Развертывание (Deployment)
(Содержимое существующего раздела актуально).

## 11. Мониторинг и Логирование (Logging and Monitoring)
(Содержимое существующего раздела актуально).

## 12. Нефункциональные Требования (NFRs)
(Содержимое существующего раздела актуально).

## 13. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API и Protobuf определения для gRPC API (если будут использоваться) поддерживаются в соответствующих репозиториях исходного кода сервиса и/или в централизованном репозитории `platform-protos`.
*   DDL схемы базы данных управляются через систему миграций (например, `golang-migrate/migrate`) и хранятся в репозитории исходного кода сервиса. Актуальные версии доступны во внутренней документации команды разработки и в GitOps репозитории.

## 14. Резервное Копирование и Восстановление (Backup and Recovery)

### 14.1. PostgreSQL (Шаблоны, предпочтения, недавние логи, токены устройств, метаданные кампаний)
*   **Процедура резервного копирования:**
    *   Ежедневный логический бэкап (`pg_dump`).
    *   Непрерывная архивация WAL-сегментов (PITR), базовый бэкап еженедельно.
    *   **Хранение:** Бэкапы в S3, шифрование, версионирование, другой регион. Срок хранения: полные - 30 дней, WAL - 14 дней.
*   **Процедура восстановления:** Тестируется ежеквартально.
*   **RTO:** < 2 часов.
*   **RPO:** < 10 минут.
*   (Общие принципы см. `../../../../project_database_structure.md`).

### 14.2. ClickHouse (Статистика доставки и взаимодействия)
*   **Стратегия:**
    *   Использование инструмента `clickhouse-backup` для создания инкрементальных или полных бэкапов.
    *   Для кластерных инсталляций – использование встроенной репликации ClickHouse.
    *   **Частота:** Ежедневно для таблиц статистики.
    *   **Хранение:** Бэкапы в S3. Срок хранения 30-90 дней.
*   **Восстановление:** Из бэкапов `clickhouse-backup` или через восстановление реплики. Часть статистики может быть восстановлена путем пересчета из логов Kafka или логов PostgreSQL за определенный период, если это предусмотрено.
*   **RTO:** < 4-8 часов (зависит от объема данных).
*   **RPO:** 24 часа (для ежедневных бэкапов).

### 14.3. Redis (Кэш, очереди, счетчики)
*   **Стратегия:**
    *   **AOF (Append Only File):** Включен с fsync `everysec` для очередей и счетчиков, если их потеря критична.
    *   **RDB Snapshots:** Регулярное создание снапшотов.
*   **Резервное копирование (снапшотов):** RDB-снапшоты могут копироваться в S3.
*   **Восстановление:** Из RDB/AOF. Большинство данных кэша могут быть перестроены. Критичные очереди должны быть устойчивы к потере (например, через подтверждения в Kafka).
*   **RTO:** < 30 минут.
*   **RPO:** < 1 минуты (для AOF).

### 14.4. Общая стратегия
*   Восстановление PostgreSQL и ClickHouse является приоритетным.
*   Процедуры документированы и тестируются.

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [Регистрация пользователя и начальная настройка профиля](../../../../project_workflows/user_registration_flow.md)
*   [Процесс покупки игры и обновления библиотеки пользователя](../../../../project_workflows/game_purchase_flow.md)
*   [Процесс подачи разработчиком новой игры на модерацию](../../../../project_workflows/game_submission_flow.md)
*   [Процесс сброса и восстановления пароля](../../../../project_workflows/password_reset_flow.md) (Примечание: Создание документа `password_reset_flow.md` является частью общей задачи по документированию проекта и выходит за рамки обновления документации данного микросервиса.)

---
*Этот документ является основной спецификацией для Notification Service и должен поддерживаться в актуальном состоянии.*
