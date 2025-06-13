# Спецификация Микросервиса: Notification Service (Сервис Уведомлений)

**Версия:** 1.1 (адаптировано из предыдущей версии)
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Notification Service платформы "Российский Аналог Steam".
*   **Роль в общей архитектуре платформы:** Notification Service отвечает за централизованную отправку и управление всеми видами уведомлений пользователям платформы. Это включает транзакционные уведомления (например, подтверждение регистрации, информация о заказе), системные оповещения (например, о технических работах), маркетинговые рассылки и уведомления, инициированные другими микросервисами (например, о выходе новой игры из списка желаемого, о разблокировке достижения). Сервис абстрагирует сложность взаимодействия с различными провайдерами доставки (Email, Push, SMS) и каналами (In-App).
*   **Основные бизнес-задачи:**
    *   Обеспечение надежной и своевременной доставки уведомлений пользователям по предпочтительным для них каналам.
    *   Управление пользовательскими подписками и предпочтениями по типам уведомлений.
    *   Предоставление инструментов для создания и управления шаблонами уведомлений.
    *   Поддержка маркетинговых кампаний и массовых рассылок.
    *   Сбор статистики по доставке и взаимодействию с уведомлениями.
*   Разработка сервиса должна вестись в соответствии с `CODING_STANDARDS.md`.

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
*   **Язык программирования:** Go (предпочтительно).
*   **Веб-фреймворк (REST API):** Echo или Gin (для Go).
*   **RPC фреймворк (gRPC):** Опционально, для внутреннего высокопроизводительного взаимодействия.
*   **Очереди сообщений (основной способ получения запросов на отправку):** Apache Kafka.
*   **Шаблонизатор:** Стандартные пакеты Go `text/template` и `html/template`, или сторонние библиотеки при необходимости.
*   **Базы данных:**
    *   PostgreSQL: Для хранения метаданных (шаблоны, кампании, пользовательские предпочтения, токены устройств, логи некоторых сообщений).
    *   ClickHouse (опционально): Для хранения и агрегации больших объемов статистики по доставке и взаимодействию с уведомлениями.
*   **Кэширование:** Redis (для кэширования пользовательских предпочтений, шаблонов, счетчиков для rate limiting, временных данных кампаний).
*   **Инфраструктура:** Docker, Kubernetes.
*   **Мониторинг/Трассировка:** OpenTelemetry, Prometheus, Grafana, Jaeger/Tempo.
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `PACKAGE_STANDARDIZATION.md`.
*   *Примечание: Приведенные в последующих разделах примеры API и конфигураций основаны на предположении использования Go (Echo/Gin для REST), PostgreSQL, ClickHouse, Redis и Kafka.*

### 1.4. Термины и Определения (Glossary)
*   **Уведомление (Notification):** Информационное сообщение, отправляемое пользователю или системе.
*   **Канал Доставки (Delivery Channel):** Способ, которым уведомление достигает пользователя (Email, Push, SMS, In-App).
*   **Шаблон (Template):** Предопределенная структура сообщения с плейсхолдерами для персонализации.
*   **Провайдер Доставки (Provider):** Внешний сервис, используемый для фактической отправки уведомлений по определенному каналу (например, SendGrid для Email, Firebase Cloud Messaging (FCM) для Push).
*   **Кампания (Notification Campaign):** Организованная рассылка одного или нескольких уведомлений целевой аудитории, обычно в маркетинговых целях.
*   **Токен Устройства (Device Token):** Уникальный идентификатор, используемый для отправки Push-уведомлений на конкретное мобильное устройство.
*   **In-App Уведомление:** Уведомление, отображаемое внутри клиентского приложения платформы.
*   Для других общих терминов см. `project_glossary.md`.

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

#### 2.2.1. Presentation Layer (Слой Представления / API и Потребители Событий)
*   **Ответственность:** Прием входящих запросов на управление и отправку уведомлений.
*   **Ключевые компоненты/модули:**
    *   **REST API (Gin/Echo):** Предоставляет эндпоинты для управления шаблонами, кампаниями, пользовательскими предпочтениями, устройствами, а также для прямой отправки уведомлений (например, администраторами).
    *   **Kafka Consumers:** Основной канал для получения запросов на отправку транзакционных и системных уведомлений от других микросервисов. Каждый консьюмер обрабатывает сообщения из определенных топиков.

#### 2.2.2. Application Layer (Прикладной Слой / Оркестрация и Управление)
*   **Ответственность:** Центральная логика сервиса. Обработка запросов, принятых из Presentation Layer, управление жизненным циклом уведомления.
*   **Ключевые компоненты/модули:**
    *   **Notification Orchestrator:** Получает запрос на уведомление (из Kafka или API), определяет тип уведомления, проверяет пользовательские предпочтения (через Preference Manager), выбирает подходящий шаблон (через Template Manager), обогащает шаблон данными, определяет каналы доставки и ставит задачи в очереди для соответствующих диспетчеров каналов.
    *   **Template Manager:** Управляет CRUD операциями для шаблонов уведомлений, включая их локализацию и рендеринг с переменными.
    *   **Preference Manager:** Управляет настройками уведомлений пользователей (подписки, отписки, предпочитаемые каналы).
    *   **Campaign Manager:** Управляет созданием, планированием и выполнением маркетинговых кампаний. Взаимодействует с Analytics Service для получения сегментов аудитории.
    *   **Device Token Manager:** Управляет регистрацией и обновлением токенов устройств для Push-уведомлений.

#### 2.2.3. Domain Layer (Доменный Слой)
*   **Ответственность:** Определение бизнес-сущностей, их состояний и правил валидации.
*   **Ключевые компоненты/модули:**
    *   **Entities:** `NotificationTemplate` (шаблон), `NotificationMessage` (экземпляр отправляемого или отправленного сообщения), `UserPreferences` (предпочтения пользователя), `NotificationCampaign` (маркетинговая кампания), `DeviceToken` (токен устройства), `ProviderConfig` (конфигурация внешнего провайдера).
    *   **Value Objects:** `ChannelType` (Email, Push, SMS, InApp), `NotificationStatus` (pending, sent, delivered, failed, opened, clicked), `MessagePriority`.
    *   **Domain Events:** `NotificationRequestedEvent`, `NotificationSentToProviderEvent`, `NotificationDeliveryStatusUpdatedEvent`.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой / Диспетчеризация и Взаимодействие с Хранилищами)
*   **Ответственность:** Непосредственная отправка уведомлений через внешних провайдеров. Взаимодействие с базами данных и кэшем. Публикация событий о статусах.
*   **Ключевые компоненты/модули:**
    *   **Channel Dispatchers:** Специализированные компоненты для каждого канала доставки (EmailDispatcher, PushDispatcher, SMSDispatcher, InAppDispatcher). Получают готовые к отправке сообщения от Orchestrator и взаимодействуют с соответствующими внешними провайдерами. Обрабатывают ответы от провайдеров.
    *   **Stats Collector:** Агрегирует информацию о статусах доставки и взаимодействии с уведомлениями, периодически сохраняя ее в ClickHouse и публикуя агрегированные события.
    *   **Repository Layer:** Реализации интерфейсов репозиториев для PostgreSQL (хранение шаблонов, предпочтений, логов сообщений, токенов устройств) и ClickHouse (хранение статистики).
    *   **Redis Client:** Для кэширования шаблонов, предпочтений, управления очередями задач для диспетчеров (если не используется Kafka для этого), счетчиков rate limiting.
    *   **Kafka Producer:** Для публикации событий о статусах доставки.
    *   **Клиенты внешних провайдеров:** HTTP/SMTP/APNS/FCM клиенты для взаимодействия с сервисами отправки.

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL:** `/api/v1/notifications` (маршрутизируется через API Gateway).
*   **Аутентификация:** JWT Bearer Token для запросов от имени пользователей или администраторов. API-ключи для межсервисных запросов (например, для прямой отправки уведомления от другого сервиса).
*   **Авторизация:** На основе ролей (`notification_admin`, `marketing_manager`, `user_self_only`, `service_internal`).
*   **Стандартный формат ответа об ошибке:**
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

#### 3.1.1. Управление Шаблонами Уведомлений
*   **`POST /templates`**
    *   Описание: Создание нового шаблона уведомления.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "notificationTemplate",
            "attributes": {
              "name": "order_confirmation_email",
              "channel_type": "email", // email, sms, push, in_app
              "default_language_code": "ru-RU",
              "subject_template": "Заказ {{order_id}} подтвержден!", // Для email/push
              "body_template_text": "Уважаемый {{user_name}}, ваш заказ {{order_id}} на сумму {{total_amount}} успешно оформлен.", // Для sms, простого push
              "body_template_html": "<p>Уважаемый {{user_name}},</p><p>Ваш заказ <b>{{order_id}}</b> на сумму {{total_amount}} успешно оформлен.</p>", // Для email, in_app
              "example_variables": ["order_id", "user_name", "total_amount"]
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Возвращает созданный шаблон)
    *   Требуемые права доступа: `notification_admin`.
*   **`GET /templates/{template_id}`**
    *   Описание: Получение информации о шаблоне.
    *   Пример ответа (Успех 200 OK): (Структура аналогична запросу на создание, но с `id` и метаданными)
    *   Требуемые права доступа: `notification_admin`, `marketing_manager`.

#### 3.1.2. Управление Пользовательскими Предпочтениями
*   **`GET /preferences/users/{user_id}`**
    *   Описание: Получение текущих настроек уведомлений для пользователя.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "userNotificationPreferences",
            "id": "user-uuid-123",
            "attributes": {
              "global_unsubscribe_all": false,
              "channel_preferences": {
                "email": { "enabled": true, "preferred_email": "user@example.com" },
                "push": { "enabled": true },
                "sms": { "enabled": false, "phone_number": "+79XXXXXXXXX" }
              },
              "notification_type_preferences": { // Ключи - типы уведомлений
                "order_updates": { "email": true, "push": true, "sms": false },
                "new_game_in_wishlist_on_sale": { "email": true, "push": false }
              }
            }
          }
        }
        ```
    *   Требуемые права доступа: `user_self_only` или `notification_admin`.
*   **`PUT /preferences/users/{user_id}`**
    *   Описание: Обновление настроек уведомлений для пользователя.
    *   Тело запроса: (Аналогично ответу GET)
    *   Пример ответа (Успех 200 OK): (Обновленные предпочтения)
    *   Требуемые права доступа: `user_self_only` или `notification_admin`.

#### 3.1.3. Отправка Уведомлений (Прямая, для сервисов/админов)
*   **`POST /send/direct`**
    *   Описание: Прямая отправка одиночного уведомления (используется, когда нет подходящего события Kafka или для ручной отправки).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "directNotificationRequest",
            "attributes": {
              "user_id": "user-uuid-123", // или recipient_details если пользователь не идентифицирован
              // "recipient_details": { "email_to": "test@example.com" },
              "template_name": "custom_admin_message", // Имя ранее созданного шаблона
              // или "content_override": { "channel": "email", "subject": "Важное сообщение", "body_text": "..." }
              "data_payload": { "admin_message_text": "Проверка системы." },
              "priority": "normal" // high, normal, low
            }
          }
        }
        ```
    *   Пример ответа (Успех 202 Accepted): `{ "data": { "type": "messageStatus", "id": "message-uuid-xyz", "status": "queued" } }`
    *   Требуемые права доступа: `service_internal` (с API ключом) или `notification_admin`.

### 3.2. gRPC API
*   В настоящее время основной упор делается на REST API для управления и Kafka для получения запросов на отправку. gRPC может быть рассмотрен для высокопроизводительных внутренних вызовов, например, `rpc SendNotificationInternal(SendNotificationInternalRequest) returns (SendNotificationInternalResponse)` для прямого вызова от другого сервиса с минимальными накладными расходами, но это требует дополнительного проектирования. Пока не является приоритетом.
*   Если будет реализован, то `message SendNotificationInternalRequest { string user_id = 1; string template_name = 2; map<string, string> data_payload = 3; ... }`

### 3.3. WebSocket API (для In-App уведомлений)
*   Notification Service напрямую не предоставляет WebSocket API и не управляет WebSocket соединениями с клиентами.
*   Для доставки In-App уведомлений он отправляет сообщения в специализированный **WebSocket Gateway** (предположительно, отдельный сервис или компонент API Gateway), который уже поддерживает постоянные соединения с клиентскими приложениями.
*   Сообщение, отправляемое Notification Service в WebSocket Gateway (например, через Kafka или внутренний gRPC вызов), может иметь следующий концептуальный формат:
    ```json
    {
      "target_user_id": "user-uuid-123", // Идентификатор пользователя для WebSocket Gateway
      "event_type": "in_app_notification", // Тип события для клиента
      "payload": { // Полезная нагрузка, которую клиентское приложение сможет отобразить
        "notification_id": "inapp-uuid-abc",
        "title": "Новое сообщение от поддержки!",
        "short_text": "У вас новое сообщение в тикете #T12345.",
        "icon_url": "https://example.com/icons/support.png",
        "deep_link": "mygameplatform://support/ticket/T12345", // Для перехода в приложении
        "urgency": "high",
        "display_type": "toast" // "modal", "banner", "toast", "feed_item"
      },
      "send_options": {
        "require_ack": false // Требуется ли подтверждение от WebSocket Gateway
      }
    }
    ```

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности

*   **`NotificationTemplate` (Шаблон Уведомления)**
    *   `id` (UUID): Уникальный идентификатор. Обязательность: Required.
    *   `name` (VARCHAR(255)): Уникальное имя шаблона (для использования в API/событиях). Пример: `order_confirmation_email`. Валидация: not null, unique. Обязательность: Required.
    *   `description` (TEXT): Описание шаблона. Обязательность: Optional.
    *   `channel_type` (ENUM: `email`, `sms`, `push_fcm`, `push_apns`, `in_app`): Канал доставки. Обязательность: Required.
    *   `default_language_code` (VARCHAR(10)): Код языка по умолчанию (например, `ru-RU`). Обязательность: Required.
    *   `versions` (JSONB): Массив версий шаблона, каждая со своим телом, темой и языком. Пример: `[{"version": 1, "lang": "ru-RU", "subject": "Заказ {{order_id}}", "body_text": "...", "body_html": "...", "is_active": true}]`. Обязательность: Required.
    *   `example_variables` (ARRAY of TEXT): Список переменных, используемых в шаблоне. Пример: `["order_id", "user_name"]`. Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ), `updated_at` (TIMESTAMPTZ).

*   **`NotificationMessageLog` (Запись об Отправленном Уведомлении)** (может храниться в PostgreSQL для недавних и в ClickHouse для долгосрочных)
    *   `id` (UUID или BIGSERIAL): Уникальный идентификатор сообщения. Обязательность: Required.
    *   `user_id` (UUID): ID пользователя-получателя. Обязательность: Required (если не системное широковещательное).
    *   `template_id` (UUID, FK to NotificationTemplate): ID использованного шаблона. Обязательность: Optional (если контент был задан напрямую).
    *   `campaign_id` (UUID, FK to NotificationCampaign): ID кампании, если сообщение является ее частью. Обязательность: Optional.
    *   `channel_type` (ENUM): Канал доставки. Обязательность: Required.
    *   `provider_name` (VARCHAR(100)): Имя провайдера, через которого была попытка отправки. Обязательность: Optional.
    *   `recipient_address` (TEXT): Адрес получателя (email, номер телефона, токен устройства). Обязательность: Required.
    *   `status` (ENUM: `queued`, `sent`, `delivered`, `failed`, `opened`, `clicked`, `undeliverable_address`, `user_unsubscribed`): Статус доставки. Обязательность: Required.
    *   `status_details` (TEXT): Дополнительная информация о статусе (например, сообщение об ошибке от провайдера). Обязательность: Optional.
    *   `requested_at` (TIMESTAMPTZ): Время запроса на отправку. Обязательность: Required.
    *   `sent_at` (TIMESTAMPTZ): Время отправки провайдеру. Обязательность: Optional.
    *   `delivered_at` (TIMESTAMPTZ): Время подтверждения доставки. Обязательность: Optional.
    *   `opened_at` (TIMESTAMPTZ): Время открытия (для email/push с трекингом). Обязательность: Optional.
    *   `clicked_at` (TIMESTAMPTZ): Время клика по ссылке (для email/push с трекингом). Обязательность: Optional.
    *   `rendered_subject` (TEXT): Отрендеренная тема (для email/push). Обязательность: Optional.
    *   `rendered_body_preview` (TEXT): Начало отрендеренного тела (для логов). Обязательность: Optional.

*   **`UserNotificationPreferences` (Предпочтения Пользователя по Уведомлениям)**
    *   `user_id` (UUID, PK): ID пользователя. Обязательность: Required.
    *   `global_unsubscribe_all` (BOOLEAN): Глобальная отписка от всех уведомлений (кроме критически важных). Обязательность: Required, default false.
    *   `channel_settings` (JSONB): Настройки для каждого канала. Пример: `{"email": {"enabled": true, "address": "user@example.com"}, "push_fcm": {"enabled": false}}`. Обязательность: Optional.
    *   `type_preferences` (JSONB): Настройки для каждого типа уведомлений. Пример: `{"order_updates": {"email": true, "push": false}, "marketing_promo": {"email": false}}`. Обязательность: Optional.
    *   `updated_at` (TIMESTAMPTZ).

*   **`DeviceToken` (Токен Устройства для Push)**
    *   `id` (UUID): Уникальный идентификатор.
    *   `user_id` (UUID, FK to User): ID пользователя. Обязательность: Required.
    *   `device_platform` (ENUM: `android_fcm`, `ios_apns`): Платформа устройства. Обязательность: Required.
    *   `token` (TEXT): Сам токен устройства. Валидация: not null, unique. Обязательность: Required.
    *   `last_seen_at` (TIMESTAMPTZ): Время последнего использования токена. Обязательность: Required.
    *   `is_active` (BOOLEAN): Активен ли токен. Обязательность: Required, default true.
    *   `app_version` (VARCHAR(50)): Версия клиентского приложения. Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ).

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL (Метаданные, Шаблоны, Предпочтения, Логи недавних сообщений)

**ERD Диаграмма (ключевые таблицы PostgreSQL):**
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
        VARCHAR target_segment_id "FK (Segment from Analytics)"
        UUID template_id FK
        VARCHAR status -- draft, scheduled, active, completed, archived
        TIMESTAMPTZ scheduled_at
        TIMESTAMPTZ started_at
        TIMESTAMPTZ completed_at
        JSONB stats_summary
    }
    NOTIFICATION_MESSAGE_LOG { # Для недавних сообщений, основные логи в ClickHouse
        UUID id PK
        UUID user_id "FK (User ID, nullable for system)"
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

    USERS { # Предполагается из Auth/Account Service
        UUID id PK
        VARCHAR email
        VARCHAR phone_number
    }

    NOTIFICATION_TEMPLATES ||--o{ NOTIFICATION_CAMPAIGNS : "uses"
    USERS ||--|| USER_NOTIFICATION_PREFERENCES : "has"
    USERS ||--o{ DEVICE_TOKENS : "has_devices"
    USERS ||--o{ NOTIFICATION_MESSAGE_LOG : "receives"
    NOTIFICATION_TEMPLATES ||--o{ NOTIFICATION_MESSAGE_LOG : "based_on_template"
    NOTIFICATION_CAMPAIGNS ||--o{ NOTIFICATION_MESSAGE_LOG : "part_of_campaign"
```

**DDL (PostgreSQL - примеры ключевых таблиц):**
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE notification_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    channel_type VARCHAR(20) NOT NULL CHECK (channel_type IN ('email', 'sms', 'push_fcm', 'push_apns', 'in_app')),
    default_language_code VARCHAR(10) NOT NULL DEFAULT 'ru-RU',
    versions JSONB NOT NULL DEFAULT '[]'::jsonb, -- [{"version": 1, "lang": "ru-RU", "subject_template": "...", "body_template_text": "...", "body_template_html": "...", "is_active": true}]
    example_variables TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_notification_preferences (
    user_id UUID PRIMARY KEY, -- FK to users table in Account/Auth Service
    global_unsubscribe_all BOOLEAN NOT NULL DEFAULT FALSE,
    channel_settings JSONB, -- {"email": {"enabled": true, "address": "user@example.com"}, "push_fcm": {"enabled": false, "device_tokens": ["token1"]}}
    type_preferences JSONB, -- {"order_updates": {"email": true, "push": false}, "marketing_promo": {"email": false}}
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE device_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL, -- FK to users table
    device_platform VARCHAR(20) NOT NULL CHECK (device_platform IN ('android_fcm', 'ios_apns')),
    token TEXT NOT NULL,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    app_version VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, token, device_platform) -- Один токен на юзера и платформу
);
CREATE INDEX idx_device_tokens_user_id ON device_tokens(user_id);
CREATE INDEX idx_device_tokens_token ON device_tokens(token); -- Для поиска по токену (например, при обратной связи от FCM/APNS)

-- Таблица для недавних логов, основные данные идут в ClickHouse
CREATE TABLE notification_message_log_recent (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID, -- Nullable for system-wide or non-user-specific notifications
    template_id UUID REFERENCES notification_templates(id) ON DELETE SET NULL,
    campaign_id UUID, -- FK to notification_campaigns if created
    channel_type VARCHAR(20) NOT NULL,
    provider_name VARCHAR(100),
    recipient_address TEXT NOT NULL, -- Email, phone, device token, or user_id for in-app
    status VARCHAR(50) NOT NULL DEFAULT 'queued',
    status_details TEXT,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    sent_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    opened_at TIMESTAMPTZ,
    clicked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_notification_log_user_id ON notification_message_log_recent(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_notification_log_status_channel ON notification_message_log_recent(status, channel_type);
CREATE INDEX idx_notification_log_created_at ON notification_message_log_recent(created_at DESC);

-- TODO: Добавить DDL для notification_campaigns, provider_configs.
```

#### 4.2.2. ClickHouse (Статистика и Долгосрочные Логи)
*   **Роль:** Хранение и агрегация больших объемов данных о доставке уведомлений и взаимодействии с ними для аналитики.
*   **Пример таблицы `notification_delivery_stats`:**
    ```sql
    CREATE TABLE notification_delivery_stats (
        event_date Date,
        event_datetime DateTime,
        message_id String, -- UUID в строковом представлении
        user_id String,    -- UUID в строковом представлении
        campaign_id Nullable(String),
        template_name String,
        channel_type Enum8('email'=1, 'sms'=2, 'push_fcm'=3, 'push_apns'=4, 'in_app'=5),
        provider_name String,
        status Enum8('queued'=1, 'sent'=2, 'delivered'=3, 'failed'=4, 'opened'=5, 'clicked'=6, 'user_unsubscribed'=7, 'invalid_address'=8),
        country_code FixedString(2), -- Из профиля пользователя или IP
        app_version String,          -- Версия клиентского приложения
        error_code Nullable(String), -- Код ошибки от провайдера
        processing_time_ms UInt32   -- Время от запроса до отправки/ошибки
    )
    ENGINE = MergeTree()
    PARTITION BY toYYYYMM(event_date)
    ORDER BY (event_date, channel_type, status, template_name, user_id);
    ```

#### 4.2.3. Redis
*   **Кэширование:**
    *   Пользовательские предпочтения: `user_prefs:<user_id>` (JSON). TTL: часы/дни.
    *   Активные шаблоны уведомлений: `template:<template_name_or_id>:<lang>` (скомпилированный текст/html). TTL: часы/дни.
*   **Счетчики для Rate Limiting:**
    *   `rl:notification:<user_id>:<notification_type>` (COUNTER). TTL: минуты/часы.
    *   `rl:provider:<provider_name>:api_calls` (COUNTER).
*   **Прогресс/статус маркетинговых кампаний (для быстрых обновлений):**
    *   `campaign_stats:<campaign_id>` (HASH). Поля: `sent_count`, `delivered_count`, `opened_count`.
*   **Очереди задач (если Kafka используется только для входа, а внутренняя диспетчеризация через Redis):**
    *   `queue:email_notifications`, `queue:push_notifications` (LIST или ZSET).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Apache Kafka.
*   **Формат событий:** CloudEvents v1.0 (JSON encoding), согласно `project_api_standards.md`.
*   **Основной топик для статусов:** `notification.status.events.v1`.

*   **`notification.message.status.updated.v1`**
    *   Описание: Статус доставки конкретного сообщения был обновлен (например, отправлено, доставлено, ошибка, открыто).
    *   Пример Payload (`data` секция CloudEvent):
        ```json
        {
          "message_id": "msg-uuid-12345",
          "user_id": "user-uuid-123",
          "channel_type": "email",
          "new_status": "delivered", // sent, failed, opened, clicked
          "provider_name": "sendgrid",
          "timestamp": "2024-03-15T10:05:00Z",
          "campaign_id": "campaign-uuid-abc", // Опционально
          "template_name": "order_confirmation_email", // Опционально
          "status_details": "Successfully delivered to recipient's mail server." // Опционально
        }
        ```
*   **`notification.campaign.stats.aggregated.v1`** (менее частое событие, для аналитики)
    *   Описание: Периодическая агрегированная статистика по маркетинговой кампании.
    *   Пример Payload:
        ```json
        {
          "campaign_id": "campaign-uuid-abc",
          "aggregation_period_start": "2024-03-15T10:00:00Z",
          "aggregation_period_end": "2024-03-15T11:00:00Z",
          "metrics": {
            "sent_count": 500,
            "delivered_count": 480,
            "opened_count": 150,
            "clicked_count": 30,
            "failure_count_by_reason": {
              "invalid_email": 5,
              "provider_bounce": 15
            }
          }
        }
        ```

### 5.2. Потребляемые События (Consumed Events)
*   **Основной топик для запросов на отправку:** `notification.send.request.v1` (из любого сервиса).
*   **`notification.send.request.v1`**
    *   Описание: Запрос на отправку уведомления пользователю или группе пользователей.
    *   Ожидаемый Payload:
        ```json
        {
          "request_id": "req-uuid-unique", // Для трассировки и идемпотентности
          "recipients": [ // Один или несколько получателей
            { "user_id": "user-uuid-123" }
            // или { "email_to": "external@example.com" } для неидентифицированных
            // или { "device_token": "fcm_token_xyz", "platform": "android_fcm" }
          ],
          "template_name": "new_feature_announcement", // Имя шаблона
          "language_override": "en-US", // Опционально, иначе используется язык пользователя
          "data_payload": { // Переменные для шаблона
            "feature_name": "Cloud Saves",
            "launch_date": "2024-04-01"
          },
          "delivery_preferences": { // Опционально
            "channels": ["email", "push_fcm"], // Предпочтительные каналы (сервис выберет доступные)
            "priority": "normal", // high, normal, low
            "send_at_utc": null // null для немедленной отправки, или ISO8601 для отложенной
          },
          "correlation_id": "triggering_event_or_campaign_id" // Для связи с источником
        }
        ```
    *   Логика обработки: Notification Orchestrator получает событие. Для каждого получателя: определяет его `user_id` (если не указан напрямую), проверяет его предпочтения и подписки, выбирает подходящий шаблон и язык, рендерит шаблон с `data_payload`, определяет доступные и разрешенные каналы, ставит задачи в очереди для соответствующих диспетчеров каналов.
*   **Другие события от различных сервисов (примеры):**
    *   `user.account.registered.v1` (от Auth/Account Service) -> Отправка "Добро пожаловать" и "Подтвердите Email".
    *   `payment.order.completed.v1` (от Payment Service) -> Отправка "Чек по заказу".
    *   `social.new_reply_to_comment.v1` (от Social Service) -> Отправка "Вам ответили на комментарий".
    *   `library.game.added_to_wishlist_on_sale.v1` (от Library/Catalog Service) -> Отправка "Игра из вашего списка желаемого теперь со скидкой!".

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Account Service:** Получение контактных данных пользователя (email, телефон), языка, глобальных настроек уведомлений (через gRPC).
*   **Auth Service:** Получение информации о валидности сессий/устройств, возможно, событий безопасности, требующих уведомления (через Kafka или gRPC).
*   **Analytics Service:** Получение сегментов пользователей для маркетинговых кампаний (через gRPC/REST). Отправка статистики по доставке и взаимодействию с уведомлениями (через Kafka в ClickHouse, управляемый Analytics Service, или напрямую).
*   **Любой другой микросервис:** Может публиковать события в Kafka (например, `payment.order.completed.v1`), которые Notification Service потребляет для инициирования отправки соответствующих уведомлений.
*   **API Gateway:** Маршрутизация REST API запросов. Может управлять WebSocket соединениями для In-App уведомлений, получая сообщения от Notification Service через внутреннюю шину (например, Kafka или NATS).

### 6.2. Внешние Системы
*   **Email Провайдеры (например, SendGrid, Mailgun, Amazon SES):** Интеграция через их REST API или SMTP для отправки email.
*   **Push Уведомлений Провайдеры (Firebase Cloud Messaging - FCM, Apple Push Notification service - APNS):** Интеграция через их HTTP API для отправки push-уведомлений на мобильные устройства.
*   **SMS Шлюзы (например, Twilio, SMSC.ru):** Интеграция через их REST API для отправки SMS.
*   **WebSocket Gateway (если это отдельный компонент):** Внутренняя интеграция для доставки In-App уведомлений.

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `NOTIFICATION_SERVICE_HTTP_PORT`: Порт для REST API.
*   `NOTIFICATION_SERVICE_GRPC_PORT` (если используется).
*   `POSTGRES_DSN`: DSN для PostgreSQL.
*   `CLICKHOUSE_DSN` (если используется): DSN для ClickHouse.
*   `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB_NOTIFICATION`: Параметры Redis.
*   `KAFKA_BROKERS`: Список брокеров Kafka.
*   `KAFKA_CONSUMER_GROUP_ID_NOTIFICATIONS`: ID группы консьюмеров для топиков запросов на отправку.
*   `KAFKA_TOPIC_SEND_REQUESTS`: Основной топик для получения запросов на отправку.
*   `KAFKA_TOPIC_NOTIFICATION_STATUS_EVENTS`: Топик для публикации событий о статусах.
*   `EMAIL_PROVIDER_PRIMARY_TYPE` (sendgrid, mailgun, etc.), `EMAIL_PROVIDER_PRIMARY_API_KEY` (из Secrets), `EMAIL_SENDER_DEFAULT_FROM`, `EMAIL_SENDER_DEFAULT_NAME`.
*   `SMS_PROVIDER_PRIMARY_TYPE`, `SMS_PROVIDER_PRIMARY_API_KEY` (из Secrets), `SMS_PROVIDER_PRIMARY_SENDER_ID`.
*   `FCM_SERVER_KEY` (из Secrets).
*   `APNS_KEY_ID`, `APNS_TEAM_ID`, `APNS_PRIVATE_KEY_PATH`, `APNS_BUNDLE_ID` (секреты и путь к .p8 файлу).
*   `WEBSOCKET_GATEWAY_INTERNAL_API_URL`: URL для отправки сообщений в WebSocket Gateway.
*   `LOG_LEVEL`.
*   `DEFAULT_MESSAGE_TTL_SECONDS_EMAIL`, `DEFAULT_MESSAGE_TTL_SECONDS_PUSH`: Время жизни сообщения в очереди на отправку.
*   `MAX_RETRY_ATTEMPTS_PER_MESSAGE`: Максимальное количество попыток повторной отправки.
*   `RETRY_DELAY_SECONDS_BASE`: Базовая задержка перед повторной попыткой.
*   `ACCOUNT_SERVICE_GRPC_ADDR`, `ANALYTICS_SERVICE_GRPC_ADDR`.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`.

### 7.2. Файлы Конфигурации (если применимо)
*   **`configs/notification_config.yaml`**: Может использоваться для более сложных или менее часто изменяемых настроек.
    ```yaml
    server:
      http_port: ${NOTIFICATION_SERVICE_HTTP_PORT:-8085}

    kafka:
      brokers: ${KAFKA_BROKERS}
      topics:
        send_requests: ${KAFKA_TOPIC_SEND_REQUESTS:-notification.send.request.v1}
        status_events: ${KAFKA_TOPIC_NOTIFICATION_STATUS_EVENTS:-notification.status.events.v1}
      consumer_groups:
        send_request_group: ${KAFKA_CONSUMER_GROUP_ID_NOTIFICATIONS:-notification-service-main}

    providers:
      email:
        default_provider: "sendgrid_primary" # Имя конфигурации провайдера
        sendgrid_primary:
          api_key_env_var: "EMAIL_PROVIDER_PRIMARY_API_KEY" # Имя переменной окружения с ключом
          sender_email: ${EMAIL_SENDER_DEFAULT_FROM:-"noreply@example.com"}
          sender_name: "Наша Платформа"
          retry_policy: { "max_attempts": 3, "initial_backoff_ms": 1000, "max_backoff_ms": 30000 }
      sms:
        default_provider: "smsc_ru"
        smsc_ru:
          api_key_env_var: "SMS_PROVIDER_PRIMARY_API_KEY"
          sender_id: ${SMS_PROVIDER_PRIMARY_SENDER_ID:-"Platform"}
      push_fcm:
        server_key_env_var: "FCM_SERVER_KEY"
      push_apns:
        key_id_env_var: "APNS_KEY_ID"
        team_id_env_var: "APNS_TEAM_ID"
        private_key_path: ${APNS_PRIVATE_KEY_PATH} # Путь к .p8 файлу в контейнере
        default_bundle_id: ${APNS_BUNDLE_ID}

    template_engine:
      cache_enabled: true
      cache_ttl_seconds: 3600
      # Настройки специфичные для выбранного шаблонизатора

    delivery_retry_policy:
      default_max_attempts: ${DEFAULT_RETRY_ATTEMPTS:-3}
      default_initial_delay_ms: ${DEFAULT_RETRY_DELAY_SECONDS:-60} * 1000
      per_channel_overrides:
        sms: { "max_attempts": 2 }

    # Настройки для In-App уведомлений через WebSocket Gateway
    in_app_notifications:
      websocket_gateway_url: ${WEBSOCKET_GATEWAY_URL}
      timeout_seconds: 5
    ```

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   Обработка ошибок от внешних провайдеров доставки (API ошибки, неверные адреса, отписки).
*   Механизмы повторной отправки (retry) с экспоненциальной задержкой и ограниченным количеством попыток.
*   Использование DLQ (Dead Letter Queue) в Kafka для сообщений, которые не удалось обработать после нескольких попыток, для последующего анализа.
*   Подробное логирование всех ошибок с `trace_id`.

### 8.2. Распространенные Коды Ошибок (для REST API)
*   **`400 Bad Request` (`VALIDATION_ERROR`, `TEMPLATE_VARIABLES_MISMATCH`)**: Некорректные входные данные, отсутствуют обязательные переменные для шаблона.
*   **`401 Unauthorized` (`UNAUTHENTICATED`)**: Ошибка аутентификации.
*   **`403 Forbidden` (`PERMISSION_DENIED`)**: Недостаточно прав для выполнения операции (например, управление шаблонами).
*   **`404 Not Found` (`TEMPLATE_NOT_FOUND`, `USER_PREFERENCES_NOT_FOUND`)**: Запрашиваемый ресурс не найден.
*   **`429 Too Many Requests` (`RATE_LIMIT_EXCEEDED`)**: Превышен лимит на отправку уведомлений (например, для пользователя или глобально).
*   **`500 Internal Server Error` (`INTERNAL_ERROR`)**: Внутренняя ошибка сервера.
*   **`503 Service Unavailable` (`PROVIDER_UNAVAILABLE`)**: Внешний провайдер доставки временно недоступен.

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   REST API: JWT для пользовательских запросов и запросов от администраторов. Защищенные API-ключи для межсервисных запросов на отправку уведомлений.
*   Доступ к API внешних провайдеров (Email, SMS, Push) осуществляется с использованием API ключей или других механизмов аутентификации, предоставляемых провайдерами. Эти ключи должны безопасно храниться.

### 9.2. Авторизация
*   RBAC для доступа к административным функциям API (управление шаблонами, кампаниями, просмотр глобальной статистики).
*   Пользователи могут управлять только своими предпочтениями и просматривать свои уведомления.

### 9.3. Защита Данных
*   Соблюдение ФЗ-152 "О персональных данных": хранение и обработка контактных данных пользователей (email, телефон), токенов устройств, текстов уведомлений. Получение согласия на получение уведомлений.
*   Шифрование при передаче (TLS) для всех API и при взаимодействии с внешними провайдерами.
*   Шифрование чувствительных конфигурационных данных и API ключей провайдеров при хранении (at-rest).
*   Защита от спама и злоупотреблений: контроль частоты отправки уведомлений, проверка валидности адресов/токенов, обработка отписок. Недопущение утечки персональных данных через переменные в шаблонах.

### 9.4. Управление Секретами
*   API ключи для доступа к внешним провайдерам, секреты для баз данных и Kafka должны храниться в Kubernetes Secrets или HashiCorp Vault и безопасно передаваться в приложение.

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Многоэтапная сборка для Go-приложения.
*   **Kubernetes манифесты/Helm-чарты:** Для управления развертыванием, сервисами, конфигурациями.
*   (Ссылка на `project_deployment_standards.md`).

### 10.2. Зависимости при Развертывании
*   Кластер Kubernetes.
*   PostgreSQL (и/или ClickHouse), Redis, Kafka.
*   Доступность Account Service (для получения контактных данных и предпочтений пользователя).
*   Доступность внешних провайдеров уведомлений.
*   WebSocket Gateway (для In-App уведомлений).

### 10.3. CI/CD
*   Автоматизированная сборка, юнит- и интеграционное тестирование (с моками внешних провайдеров).
*   Развертывание в окружения с использованием GitOps.
*   Процедуры миграции схемы БД.

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат:** Структурированные JSON логи (Zap).
*   **Ключевые события:** Запросы на отправку, попытки отправки через провайдеров, статусы доставки, ошибки, изменения предпочтений, операции с шаблонами и кампаниями.
*   **Интеграция:** С централизованной системой логирования (Loki, ELK Stack).
*   (Ссылка на `project_observability_standards.md`).

### 11.2. Мониторинг
*   **Метрики (Prometheus):**
    *   Количество полученных запросов на отправку (по типам, каналам).
    *   Количество успешно отправленных уведомлений (по каналам, провайдерам, типам).
    *   Количество ошибок доставки (по каналам, провайдерам, типам ошибок).
    *   Количество открытий/кликов (если отслеживается).
    *   Задержка обработки и доставки уведомлений.
    *   Размеры очередей на отправку (Kafka, Redis).
    *   Производительность и ошибки при работе с PostgreSQL, ClickHouse, Redis.
*   **Дашборды (Grafana):** Для визуализации метрик.
*   **Алертинг (AlertManager):** Для критических ошибок, большого количества недоставленных сообщений, проблем с провайдерами, переполнения очередей.
*   (Ссылка на `project_observability_standards.md`).

### 11.3. Трассировка
*   **Интеграция:** OpenTelemetry SDK для Go. Экспорт трейсов в Jaeger или Tempo.
*   Трассировка полного жизненного цикла уведомления: от получения запроса в Kafka/API до попытки отправки через провайдера и обновления статуса.
*   (Ссылка на `project_observability_standards.md`).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:**
    *   API управления (шаблоны, предпочтения): P95 < 200 мс.
    *   Пропускная способность приема запросов на отправку из Kafka: > 10,000 сообщений/сек.
    *   Задержка от получения запроса до постановки в очередь на отправку провайдеру: P99 < 100 мс (для высокоприоритетных).
    *   Скорость отправки через провайдеров: зависит от лимитов провайдеров, но сервис должен быть способен утилизировать их по максимуму.
*   **Надежность:**
    *   Доступность сервиса: 99.95%.
    *   Гарантированная доставка критически важных уведомлений (с учетом retry и failover на резервных провайдеров, если настроено): > 99.9%.
    *   Delivery Rate (доставлено/отправлено) для Email: > 98%. Для Push: > 90% (зависит от валидности токенов). Для SMS: > 95%.
*   **Масштабируемость:** Горизонтальное масштабирование компонентов для обработки пиковых нагрузок (например, во время крупных маркетинговых кампаний или системных событий). Способность обрабатывать миллионы уведомлений в сутки.
*   **Безопасность:** Соответствие требованиям `project_security_standards.md`.

## 13. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API и Protobuf определения для gRPC API (если будут использоваться) будут храниться в соответствующих репозиториях или артефактах CI/CD.
*   Полные DDL схемы баз данных будут поддерживаться в актуальном состоянии в системе миграций.
*   TODO: Добавить ссылки на репозиторий с OpenAPI/Protobuf и на систему управления миграциями БД, когда они будут определены.

---
*Этот документ является основной спецификацией для Notification Service и должен поддерживаться в актуальном состоянии.*

## 14. Связанные Рабочие Процессы (Related Workflows)
*   [User Registration and Initial Profile Setup](../../../project_workflows/user_registration_flow.md)
*   [Game Purchase and Library Update](../../../project_workflows/game_purchase_flow.md)
*   [Developer Submits a New Game for Moderation](../../../project_workflows/game_submission_flow.md)
*   [Password Reset Flow](../../../project_workflows/password_reset_flow.md) (TODO: Создать этот документ)
