# Спецификация Микросервиса: Admin Service

**Версия:** 1.0
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   Admin Service является центральным компонентом платформы "Российский Аналог Steam", предоставляющим инструменты и интерфейсы для администраторов платформы и сотрудников службы поддержки.
*   Его основная роль - управление и надзор за различными аспектами платформы, включая модерацию контента, управление пользователями (как платформенными, так и административными), операции технической поддержки, мониторинг безопасности, конфигурацию общесистемных настроек, административную аналитику и управление маркетинговыми кампаниями.
*   Основные бизнес-задачи: обеспечение операционного контроля, поддержка пользователей, модерация контента, управление системными параметрами, обеспечение безопасности и стабильности платформы.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Управление модерацией контента:** Создание и управление очередями модерации, автоматизированные правила модерации, ручная обработка элементов модерации (игры, DLC, отзывы, комментарии, профили пользователей), управление историей модерации.
*   **Управление пользователями платформы:** Поиск, просмотр детальной информации, изменение статусов (активация, блокировка, временное ограничение), управление ролями и разрешениями пользователей.
*   **Управление административными пользователями:** Создание и управление учетными записями администраторов, назначение административных ролей и разрешений, аудит действий администраторов.
*   **Техническая поддержка:** Система управления тикетами поддержки (создание, категоризация, назначение, отслеживание статуса, ответы), управление базой знаний (создание, редактирование, публикация статей).
*   **Мониторинг безопасности:** Просмотр логов безопасности, выявление подозрительной активности, управление IP-блокировками, анализ действий пользователей и администраторов.
*   **Управление настройками платформы:** Конфигурация глобальных и региональных параметров платформы, операционных лимитов, управление функциями (feature flags).
*   **Административная аналитика:** Доступ к дашбордам и отчетам по ключевым метрикам платформы (продажи, активность пользователей, эффективность модерации, производительность поддержки).
*   **Управление маркетинговыми кампаниями:** Создание и управление промо-акциями, скидками, специальными предложениями, баннерами.

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+, согласно `../../../../project_technology_stack.md`)
*   **REST Framework:** Echo (`github.com/labstack/echo/v4`) (согласно `../../../../project_technology_stack.md` и `../../../../PACKAGE_STANDARDIZATION.md`)
*   **gRPC Framework:** `google.golang.org/grpc` (может использоваться для внутренних API или API для специализированных админ-инструментов, если потребуется в будущем; основной интерфейс - REST). (согласно `../../../../project_technology_stack.md` и `../../../../PACKAGE_STANDARDIZATION.md`)
*   **Базы данных:**
    *   PostgreSQL (версия 15+): для структурированных данных (административные пользователи, тикеты поддержки, настройки платформы, правила модерации, категории базы знаний). (согласно `../../../../project_technology_stack.md`)
    *   MongoDB (версия 6.x+): для хранения логов аудита действий администраторов, деталей элементов модерации с гибкой схемой. (согласно `../../../../project_technology_stack.md`)
*   **Поисковый движок:** Elasticsearch (версия 8.x+): для индексации и поиска по тикетам поддержки, статьям базы знаний, логам аудита. (согласно `../../../../project_technology_stack.md`)
*   **Кэширование:** Redis (версия 7.0+): для кэширования часто запрашиваемых данных, сессий администраторов (если применимо). (согласно `../../../../project_technology_stack.md`)
*   **Брокер сообщений:** Kafka (Apache Kafka версии 3.x+): для асинхронных задач (например, генерация отчетов, массовые операции по модерации, уведомления о событиях). (согласно `../../../../project_technology_stack.md` и `../../../../PACKAGE_STANDARDIZATION.md`)
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`)
*   **Логирование:** Zap (`go.uber.org/zap`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`)
*   **Валидация:** `go-playground/validator/v10` (согласно `../../../../PACKAGE_STANDARDIZATION.md`)
*   **Трассировка и метрики:** OpenTelemetry (`go.opentelemetry.io/otel`), Prometheus client (`github.com/prometheus/client_golang`) (согласно `../../../../PACKAGE_STANDARDIZATION.md` и `../../../../project_observability_standards.md`)
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   См. `../../../../project_glossary.md`.
*   **Модерация (Moderation):** Процесс проверки и управления пользовательским и разработческим контентом на соответствие правилам платформы.
*   **Тикет поддержки (Support Ticket):** Формализованный запрос пользователя или системы в службу поддержки для решения проблемы или получения информации.
*   **База знаний (Knowledge Base):** Коллекция статей, инструкций и ответов на часто задаваемые вопросы для самопомощи пользователей и поддержки агентов.
*   **Лог Аудита (Audit Log):** Запись действий администраторов и системы, имеющих значение для безопасности и контроля.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Сервис Admin Service спроектирован как модульная система, придерживающаяся принципов чистой архитектуры (Clean Architecture) для разделения ответственностей. Он может рассматриваться как набор тесно связанных внутренних модулей, каждый из которых отвечает за определенную область администрирования.
*   Основные модули: Модерация контента, Управление пользователями платформы, Управление административными пользователями, Техническая поддержка, Мониторинг безопасности, Управление настройками платформы, Административная аналитика, Управление маркетинговыми кампаниями.
*   Взаимодействие между модулями осуществляется через внутренние сервисные интерфейсы. Внешние взаимодействия с другими микросервисами платформы происходят через их опубликованные API (REST/gRPC) или через общие шины событий (Kafka).

### 2.2. Диаграмма Архитектуры
Ниже представлена диаграмма верхнеуровневой архитектуры Admin Service, показывающая его основные компоненты и взаимодействия:
```mermaid
graph TD
    APIGateway[API Gateway] --> AdminPanelUI[Admin Panel UI (Frontend)]
    AdminPanelUI --> AdminServiceAPI[Admin Service REST API]

    subgraph AdminService [Admin Service]
        direction LR
        AdminServiceAPI

        subgraph Modules
            direction TB
            Mod[Модерация]
            UserMgmt[Управление пользователями платформы]
            AdminUserMgmt[Управление администраторами]
            Support[Техническая поддержка]
            SecurityMon[Мониторинг безопасности]
            SettingsMgmt[Управление настройками платформы]
            AnalyticsAccess[Доступ к Аналитике]
            MarketingMgmt[Управление маркетингом]
        end

        AdminServiceAPI --> Mod
        AdminServiceAPI --> UserMgmt
        AdminServiceAPI --> AdminUserMgmt
        AdminServiceAPI --> Support
        AdminServiceAPI --> SecurityMon
        AdminServiceAPI --> SettingsMgmt
        AdminServiceAPI --> AnalyticsAccess
        AdminServiceAPI --> MarketingMgmt

        CoreLogic[Application & Domain Layers]

        Mod --> CoreLogic
        UserMgmt --> CoreLogic
        AdminUserMgmt --> CoreLogic
        Support --> CoreLogic
        SecurityMon --> CoreLogic
        SettingsMgmt --> CoreLogic
        AnalyticsAccess --> CoreLogic
        MarketingMgmt --> CoreLogic

        subgraph DataStores [Хранилища Данных]
            direction TB
            PostgresDB[(PostgreSQL)]
            MongoDB[(MongoDB)]
            Elasticsearch[(Elasticsearch)]
            RedisCache[(Redis)]
        end

        CoreLogic --> PostgresDB
        CoreLogic --> MongoDB
        CoreLogic --> Elasticsearch
        CoreLogic --> RedisCache

        KafkaProducerConsumer[Kafka Producer/Consumer]
        CoreLogic --> KafkaProducerConsumer

        ExtServiceClients[Клиенты Внешних Сервисов]
        CoreLogic --> ExtServiceClients
    end

    ExtServiceClients --> AccountService[Account Service (gRPC/REST)]
    ExtServiceClients --> CatalogService[Catalog Service (gRPC/REST)]
    ExtServiceClients --> PaymentService[Payment Service (gRPC/REST)]
    ExtServiceClients --> AuthService[Auth Service (gRPC/REST)]
    ExtServiceClients --> NotificationService[Notification Service (gRPC/Kafka)]
    ExtServiceClients --> AnalyticsService[Analytics Service (REST/Kafka)]

    KafkaProducerConsumer --> KafkaBroker[Apache Kafka]
    KafkaBroker --> KafkaProducerConsumer

    classDef service fill:#D6EAF8,stroke:#3498DB,stroke-width:2px;
    classDef module fill:#E8F8F5,stroke:#1ABC9C,stroke-width:2px;
    classDef datastore fill:#FDEDEC,stroke:#E74C3C,stroke-width:2px;
    classDef external fill:#FEF9E7,stroke:#F1C40F,stroke-width:2px;
    classDef api fill:#FADBD8,stroke:#C0392B,stroke-width:2px;

    class AdminService,AccountService,CatalogService,PaymentService,AuthService,NotificationService,AnalyticsService service;
    class Modules,Mod,UserMgmt,AdminUserMgmt,Support,SecurityMon,SettingsMgmt,AnalyticsAccess,MarketingMgmt module;
    class DataStores,PostgresDB,MongoDB,Elasticsearch,RedisCache datastore;
    class APIGateway,AdminPanelUI,KafkaBroker external;
    class AdminServiceAPI,KafkaProducerConsumer,ExtServiceClients,CoreLogic api;
```
*Примечание: Диаграмма упрощена для наглядности. Реальные взаимодействия могут быть сложнее и включать больше компонентов.*

### 2.3. Слои Сервиса (детальнее)

#### 2.3.1. Presentation Layer (Слой Представления)
*   Ответственность: Обработка входящих HTTP (REST) запросов от административной панели. Валидация DTO. Преобразование запросов к внутренним командам и вызов Application Layer.
*   Ключевые компоненты/модули:
    *   HTTP Handlers (Echo): Эндпоинты для всех административных функций, сгруппированные по ресурсам (например, `/platform-users`, `/moderation/items`, `/support/tickets`).
    *   DTOs (Data Transfer Objects): Структуры для запросов и ответов API, включая параметры пагинации, фильтрации, сортировки. Валидация с использованием `go-playground/validator`.
    *   Middleware: Аутентификация администратора, проверка ролей и разрешений, логирование запросов.

#### 2.3.2. Application Layer (Прикладной Слой / Слой Сценариев Использования)
*   Ответственность: Координация бизнес-логики для каждого административного модуля. Реализация сценариев использования (use cases). Взаимодействие с Domain Layer и Infrastructure Layer.
*   Ключевые компоненты/модули:
    *   Use Case Services: Например, `PlatformUserManagementService`, `ContentModerationService`, `SupportTicketOrchestrationService`, `PlatformSettingsFacade`.
    *   Интерфейсы для репозиториев (определяются здесь, реализуются в Infrastructure Layer): `AdminUserRepository`, `ModerationItemRepository`, `SupportTicketRepository`, etc.
    *   Интерфейсы для клиентов других сервисов: `AccountServiceClient`, `CatalogServiceClient`, etc.

#### 2.3.3. Domain Layer (Доменный Слой)
*   Ответственность: Содержит бизнес-сущности (entities), агрегаты, доменные события и бизнес-правила, специфичные для администрирования.
*   Ключевые компоненты/модули:
    *   Entities: `AdminUser`, `ModerationItem`, `ModerationDecision`, `SupportTicket`, `SupportTicketResponse`, `KnowledgeBaseArticle`, `SystemSetting`, `PlatformUserView` (агрегированные данные о пользователе для отображения).
    *   Value Objects: `ModerationReason`, `TicketStatus`, `AdminRole`.
    *   Domain Services: Логика, не относящаяся к конкретной сущности (например, `ModerationRuleEngine`).
    *   Domain Events: `ContentModeratedEvent`, `UserStatusChangedByAdminEvent`, `TicketResolvedEvent`, `PlatformSettingUpdatedEvent`.
    *   Интерфейсы репозиториев (определяются здесь).

#### 2.3.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Реализация интерфейсов для взаимодействия с PostgreSQL, MongoDB, Elasticsearch, Redis, Kafka. Реализация клиентов для других микросервисов.
*   Ключевые компоненты/модули:
    *   PostgreSQL Repositories (GORM): Реализации интерфейсов для работы с `admin_users`, `support_tickets`, `platform_settings`, etc.
    *   MongoDB Repositories (Official Go Driver): Реализации для работы с `audit_log_admin`, `moderation_item_details`.
    *   Elasticsearch Client: Для индексации и поиска по данным.
    *   Redis Cache: Для кэширования настроек, сессий администраторов, часто запрашиваемых данных.
    *   Kafka Producers/Consumers: Для публикации и потребления событий.
    *   gRPC/REST клиенты для Account Service, Catalog Service, Payment Service и других сервисов платформы.

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL (через API Gateway):** `/api/v1/admin`
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway, UserID администратора и его роли передаются в заголовках, например, `X-Admin-UserID`, `X-Admin-Roles`).
*   **Авторизация:** На основе ролей администратора (см. `../../../../project_roles_and_permissions.md`).
*   **Формат ответа об ошибке (согласно `project_api_standards.md`):**
    ```json
    {
      "errors": [
        {
          "code": "ERROR_CODE_STRING", // Например, "RESOURCE_NOT_FOUND", "VALIDATION_ERROR"
          "title": "Человекочитаемый заголовок ошибки",
          "detail": "Детальное описание проблемы, возможно с указанием полей.",
          "source": { "pointer": "/data/attributes/field_name" } // Опционально, для ошибок валидации
        }
      ]
    }
    ```
*   (Общие принципы пагинации, сортировки, фильтрации см. `../../../../project_api_standards.md`)

#### 3.1.1. Ресурс: Управление Пользователями Платформы
*   **`GET /platform-users`**
    *   Описание: Поиск и получение списка пользователей платформы.
    *   Query параметры: `search_query` (string), `status` (enum), `role` (string), `page`, `per_page`, `sort_by`.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [ /* ...список пользователей... */ ],
          "meta": { "total_items": 100, "total_pages": 5, "current_page": 1, "per_page": 20 },
          "links": { "self": "...", "next": "..." }
        }
        ```
    *   Пример ответа (Ошибка 400 Bad Request - неверный параметр `status`):
        ```json
        {
          "errors": [
            {
              "code": "INVALID_QUERY_PARAMETER",
              "title": "Неверный параметр запроса",
              "detail": "Параметр 'status' имеет недопустимое значение 'unknown'.",
              "source": { "pointer": "/query/status" }
            }
          ]
        }
        ```
    *   Требуемые права доступа: `admin:platform_users:read`, `support:platform_users:read_basic`.
*   **`GET /platform-users/{user_id}`**
    *   Описание: Получение детальной информации о пользователе платформы.
    *   Пример ответа (Ошибка 404 Not Found):
        ```json
        {
          "errors": [
            {
              "code": "USER_NOT_FOUND",
              "title": "Пользователь не найден",
              "detail": "Пользователь с ID 'uuid-user-unknown' не найден."
            }
          ]
        }
        ```
    *   Требуемые права доступа: `admin:platform_users:read_detailed`, `support:platform_users:read_detailed`.
*   **`PUT /platform-users/{user_id}/status`**
    *   Описание: Изменение статуса пользователя.
    *   Требуемые права доступа: `admin:platform_users:update_status`.
*   **`PUT /platform-users/{user_id}/roles`**
    *   Описание: Изменение ролей пользователя.
    *   Требуемые права доступа: `admin:platform_users:update_roles`.
*   **`GET /platform-users/{user_id}/activity-log`**
    *   Описание: Получение лога активности пользователя.
    *   Требуемые права доступа: `admin:platform_users:read_activity_log`.
*   **`GET /platform-users/{user_id}/moderation-history`**
    *   Описание: Получение истории модерации, связанной с пользователем.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "moderationActionLog",
              "id": "modlog-uuid-1",
              "attributes": {
                "moderation_item_id": "item-uuid-abc",
                "item_type": "game_review",
                "decision": "rejected",
                "reason_code": "profanity",
                "moderator_comment": "Ненормативная лексика в тексте отзыва.",
                "action_taken_at": "2024-03-10T15:00:00Z"
              },
              "relationships": {
                "moderator": { "data": { "type": "adminUser", "id": "admin-uuid-moderator" } }
              }
            }
          ],
          "meta": { "total_items": 5, "current_page": 1, "per_page": 10 }
        }
        ```
    *   Требуемые права доступа: `admin:platform_users:read_moderation_history`, `moderator:platform_users:read_moderation_history`.

#### 3.1.2. Ресурс: Модерация Контента
*   **`GET /moderation/queues`**
    *   Описание: Получение списка очередей модерации.
    *   Требуемые права доступа: `moderator:queues:read`, `admin:queues:read`.
*   **`GET /moderation/queues/{queue_id}/items`**
    *   Описание: Получение элементов из конкретной очереди модерации.
    *   Требуемые права доступа: `moderator:items:read`, `admin:items:read`.
*   **`GET /moderation/items/{item_id}`**
    *   Описание: Получение конкретного элемента на модерацию.
    *   Требуемые права доступа: `moderator:items:read_detailed`, `admin:items:read_detailed`.
*   **`POST /moderation/items/{item_id}/decisions`**
    *   Описание: Принятие решения по элементу модерации.
    *   Требуемые права доступа: `moderator:items:decide`, `admin:items:decide`.
*   **`POST /moderation/rules`**
    *   Описание: Создание нового правила автоматической модерации.
    *   Требуемые права доступа: `admin:moderation_rules:create`.
*   **`GET /moderation/items/history`**
    *   Описание: Получение истории решений по модерации (например, по ID контента или ID пользователя).
    *   Query параметры: `item_reference_id` (string, опционально), `item_type` (string, опционально), `moderator_id` (uuid, опционально), `decision_type` (enum, опционально: `approved`, `rejected`), `start_date` (date, опционально), `end_date` (date, опционально), `page`, `per_page`.
    *   Пример ответа: (Аналогично `/platform-users/{user_id}/moderation-history`, но с фокусом на элементы контента).
        ```json
        {
          "data": [
            {
              "type": "moderationDecisionLog",
              "id": "decisionlog-uuid-1",
              "attributes": {
                "item_reference_id": "game-review-xyz",
                "item_type": "game_review",
                "decision": "approved",
                "decided_at": "2024-03-11T10:00:00Z"
              },
              "relationships": {
                "moderator": { "data": { "type": "adminUser", "id": "admin-uuid-moderator" } },
                "item_submitter": { "data": { "type": "platformUser", "id": "user-uuid-submitter" } }
              }
            }
          ],
          "meta": { "total_items": 15, "current_page": 1, "per_page": 10 }
        }
        ```
    *   Требуемые права доступа: `admin:moderation_history:read`, `moderator_lead:moderation_history:read`.

#### 3.1.3. Ресурс: Техническая Поддержка
*   **`GET /support/tickets`**
    *   Описание: Получение списка тикетов поддержки.
    *   Требуемые права доступа: `support:tickets:read`, `admin:tickets:read`.
*   **`POST /support/tickets`**
    *   Описание: Создание тикета от имени пользователя.
    *   Требуемые права доступа: `support:tickets:create`, `admin:tickets:create`.
*   **`GET /support/tickets/{ticket_id}`**
    *   Описание: Получение информации о конкретном тикете.
    *   Требуемые права доступа: `support:tickets:read_detailed`, `admin:tickets:read_detailed`.
*   **`POST /support/tickets/{ticket_id}/responses`**
    *   Описание: Добавление ответа в тикет.
    *   Требуемые права доступа: `support:tickets:respond`, `admin:tickets:respond`.
*   **`PUT /support/tickets/{ticket_id}/status`**
    *   Описание: Изменение статуса тикета.
    *   Требуемые права доступа: `support:tickets:update_status`, `admin:tickets:update_status`.
*   **`PUT /support/tickets/{ticket_id}/assignee`**
    *   Описание: Назначение или изменение ответственного агента для тикета.
    *   Тело запроса: `{"data": {"type": "ticketAssignee", "attributes": {"assignee_admin_id": "support-agent-uuid-xyz"}}}`
    *   Требуемые права доступа: `support_lead:tickets:assign`, `admin:tickets:assign`.
*   **`GET /support/knowledge-base`**
    *   Описание: Поиск и получение статей из базы знаний.
    *   Требуемые права доступа: `support:kb:read`, `admin:kb:read`.
*   **`POST /support/knowledge-base`**
    *   Описание: Создание новой статьи в базе знаний.
    *   Требуемые права доступа: `admin:kb:create`, `support_lead:kb:create`.
*   **`GET /support/reports/performance`**
    *   Описание: Получение отчета по эффективности службы поддержки.
    *   Требуемые права доступа: `admin:support_reports:read`, `support_lead:support_reports:read`.
*   **`GET /support/categories`**
    *   Описание: Получение списка категорий тикетов поддержки.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            { "type": "supportTicketCategory", "id": "billing", "attributes": { "name_ru": "Вопросы по оплате", "name_en": "Billing Issues", "is_active": true, "description_ru": "Проблемы с покупками, подписками, возвратами." } },
            { "type": "supportTicketCategory", "id": "technical", "attributes": { "name_ru": "Технические проблемы", "name_en": "Technical Problems", "is_active": true, "description_ru": "Ошибки запуска игр, проблемы с клиентом." } }
          ]
        }
        ```
    *   Требуемые права доступа: `support:ticket_categories:read`, `admin:ticket_categories:read`.
*   **`POST /support/ticket-templates`**
    *   Описание: Создание шаблона ответа для тикетов.
    *   Тело запроса: `{"data": {"type": "ticketTemplate", "attributes": {"name": "Шаблон для сброса пароля", "text_ru": "Здравствуйте! Для сброса пароля...", "category_id": "account_issues", "language_code": "ru"}}}`
    *   Пример ответа (Успех 201 Created): (Возвращает созданный шаблон)
    *   Требуемые права доступа: `admin:ticket_templates:create`, `support_lead:ticket_templates:create`.

#### 3.1.4. Ресурс: Настройки Платформы
*   **`GET /platform-settings`**
    *   Описание: Получение текущих настроек платформы (сгруппированных по категориям).
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "platformSettings",
            "attributes": {
              "general": {
                "platform_name": "Моя Игровая Платформа",
                "maintenance_mode": false,
                "global_announcement_ru": null,
                "global_announcement_en": null
              },
              "user_registration": {
                "allow_new_registrations": true,
                "default_user_roles": ["user"],
                "require_email_verification": true,
                "min_password_length": 8,
                "banned_usernames": ["admin", "root", "support"]
              },
              "content_moderation": {
                "default_review_queue": "pending_manual", // "pending_auto", "pending_manual"
                "profanity_filter_level": "medium", // "off", "low", "medium", "high"
                "image_moderation_enabled": true
              },
              "payment_gateway": {
                "default_provider_key": "yoomoney_prod",
                "commission_percent_platform": 5.0,
                "regional_settings": [
                  { "region_code": "RU", "currency_code": "RUB", "enabled_providers": ["yoomoney_prod", "sbp_prod"] }
                ]
              }
              // ... другие категории настроек
            }
          }
        }
        ```
    *   Требуемые права доступа: `admin:platform_settings:read`.
*   **`PUT /platform-settings`**
    *   Описание: Обновление настроек платформы (отправляются только изменяемые поля/группы).
    *   Требуемые права доступа: `admin:platform_settings:update`.
*   **`POST /platform-settings/maintenance-schedule`**
    *   Описание: Планирование окна обслуживания.
    *   Требуемые права доступа: `admin:platform_settings:schedule_maintenance`.

#### 3.1.5. Ресурс: Управление Администраторами
*   **`GET /admin-users`**
    *   Описание: Список административных пользователей.
    *   Требуемые права доступа: `admin:admin_users:read`.
*   **`POST /admin-users`**
    *   Описание: Создание нового административного пользователя.
    *   Требуемые права доступа: `admin:admin_users:create`.
*   **`GET /admin-users/{admin_user_id}`**
    *   Описание: Получение информации о конкретном административном пользователе.
    *   Требуемые права доступа: `admin:admin_users:read_detailed`.
*   **`PUT /admin-users/{admin_user_id}/permissions`**
    *   Описание: Изменение ролей и прямых разрешений администратора.
    *   Требуемые права доступа: `admin:admin_users:update_permissions`.
*   **`GET /admin-users/{admin_user_id}/audit-log`**
    *   Описание: Получение лога действий конкретного администратора.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "adminActionLogEntry",
              "id": "auditlog-uuid-1",
              "attributes": {
                "timestamp": "2024-03-15T14:30:00Z",
                "action_type": "update_user_status",
                "target_entity_type": "PlatformUser",
                "target_entity_id": "user-uuid-xyz",
                "details": { "old_status": "active", "new_status": "blocked", "reason": "spam" },
                "ip_address": "195.10.20.30"
              }
            }
          ],
          "meta": { "total_items": 25, "current_page": 1, "per_page": 10 }
        }
        ```
    *   Требуемые права доступа: `admin:admin_users:read_audit_log`.

### 3.2. gRPC API
*   В настоящее время gRPC API для внешнего использования Admin Service не планируется. Внутренние gRPC интерфейсы могут существовать для взаимодействия между модулями Admin Service, если он реализован как набор микро-модулей, но они не являются частью публичного контракта.

### 3.3. WebSocket API
*   Не планируется для Admin Service на данном этапе.

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`AdminUser` (Администратор)**: Учетная запись администратора/сотрудника. Хранится в PostgreSQL.
*   **`ModerationItem` (Элемент модерации)**: Объект, требующий модерации (жалоба, контент). Основная информация в PostgreSQL, детали могут быть в MongoDB.
*   **`ModerationRule` (Правило модерации)**: Автоматическое правило. Хранится в PostgreSQL.
*   **`SupportTicket` (Тикет поддержки)**: Запрос пользователя. Хранится в PostgreSQL.
*   **`SupportTicketResponse` (Ответ в тикете)**: Ответ/комментарий в тикете. Хранится в PostgreSQL.
*   **`SupportTicketCategory` (Категория тикета)**: Категория для классификации тикетов. Хранится в PostgreSQL.
*   **`KnowledgeBaseArticle` (Статья базы знаний)**: Статья для самопомощи. Хранится в PostgreSQL, индексируется в Elasticsearch.
*   **`KnowledgeBaseCategory` (Категория статьи БЗ)**: Категория для статей. Хранится в PostgreSQL.
*   **`PlatformSetting` (Настройка платформы)**: Глобальный параметр конфигурации. Хранится в PostgreSQL, кэшируется в Redis.
*   **`AuditLogAdmin` (Лог действий администраторов)**: Запись о действии администратора. Хранится в MongoDB или PostgreSQL.
*   **`ModerationItemDetail` (Детали элемента модерации)**: Расширенные данные для элемента модерации. Хранится в MongoDB.

### 4.2. Схема Базы Данных
*   ERD-диаграмма и DDL для PostgreSQL таблиц приведены выше в тексте документа.
*   **MongoDB Коллекции:**
    *   `audit_logs_admin`: Хранит документы `AuditLogAdmin`. Поля: `admin_user_id`, `timestamp`, `action_type`, `target_entity_type`, `target_entity_id`, `details` (JSONB), `ip_address`.
        *   Индексы: `admin_user_id`, `timestamp`, `action_type`, `target_entity_id`.
    *   `moderation_item_details`: Хранит расширенные данные для `ModerationItem` (например, полный текст переписки, если модерируется чат, или детализированный отчет системы анализа контента). Поля: `_id` (соответствует `ModerationItem.id`), `full_content_snapshot`, `external_analysis_reports`.
        *   Индексы: `_id`.
*   **Elasticsearch Индексы:**
    *   `support_tickets_idx`: Поля: `id`, `subject`, `description`, `responses.text`, `user_id`, `status`, `priority`, `tags`, `created_at`.
    *   `knowledge_base_articles_idx`: Поля: `id`, `title`, `content_markdown`, `tags`, `category_name`, `is_published`.
    *   `admin_audit_log_idx` (если логи аудита также индексируются для поиска): Поля `id`, `admin_username`, `action_type`, `target_entity_type`, `target_entity_id`, `timestamp`, текстовый поиск по `details`.
    *   `moderation_items_idx`: Поля `id`, `item_type`, `status`, `priority`, текстовый поиск по `content_snapshot`, `user_id`.

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Kafka.
*   **Формат событий:** CloudEvents JSON (согласно `../../../../project_api_standards.md`).
*   **Основные топики Kafka:** `com.platform.admin.events.v1`.

*   **`com.platform.admin.content.moderated.v1`**
    *   Описание: Контент прошел модерацию (одобрен/отклонен/другое решение).
    *   `data` Payload:
        ```json
        {
          "moderationItemId": "item-uuid-1",
          "itemReferenceId": "review-uuid-456",
          "itemType": "game_review",
          "decision": "rejected",
          "reasonCode": "hate_speech",
          "moderatorComment": "Комментарий содержит разжигание ненависти.",
          "moderatorId": "admin-uuid-moderator",
          "userId": "user-uuid-submitter",
          "decisionTimestamp": "ISO8601_timestamp"
        }
        ```
    *   Потребители: Catalog Service, Social Service, Developer Service, Notification Service.
*   **`com.platform.admin.user.status.updated.v1`**
    *   Описание: Статус пользователя платформы изменен администратором.
    *   `data` Payload:
        ```json
        {
          "userId": "uuid-user-1",
          "oldStatus": "active",
          "newStatus": "blocked",
          "reason": "Нарушение пункта 3.4 правил сообщества.",
          "expiresAt": "ISO8601_timestamp",
          "adminId": "admin-uuid-moderator",
          "updateTimestamp": "ISO8601_timestamp"
        }
        ```
    *   Потребители: Auth Service, Account Service, Notification Service.
*   **`com.platform.admin.platform.setting.updated.v1`**
    *   Описание: Изменена глобальная настройка платформы.
    *   `data` Payload:
        ```json
        {
          "settingKey": "user_registration.allow_new_registrations",
          "oldValue": true,
          "newValue": false,
          "adminId": "admin-uuid-root",
          "updateTimestamp": "ISO8601_timestamp"
        }
        ```
    *   Потребители: Все сервисы, которые зависят от данной настройки.
*   **`com.platform.admin.support.ticket.status.updated.v1`**
    *   Описание: Статус тикета поддержки изменен.
    *   `data` Payload:
        ```json
        {
          "ticketId": "ticket-uuid-1",
          "newStatus": "resolved",
          "oldStatus": "pending_agent_reply",
          "userId": "user-uuid-xyz",
          "assigneeAdminId": "support-agent-uuid-abc",
          "updateTimestamp": "ISO8601_timestamp"
        }
        ```
    *   Потребители: Notification Service.
*   **`com.platform.admin.marketing.campaign.status.updated.v1`** (Пример)
    *   Описание: Статус маркетинговой кампании изменен (например, запущена, остановлена).
    *   `data` Payload: `{"campaignId": "uuid", "newStatus": "active", "adminId": "uuid", "timestamp": "ISO8601"}`
    *   Потребители: Catalog Service (для применения скидок), Notification Service.

### 5.2. Потребляемые События (Consumed Events)
*   **`com.platform.user.complaint.submitted.v1`** (от Social Service, Catalog Service)
    *   Описание: Пользователь подал жалобу на контент или пользователя.
    *   Логика обработки: Создать новый `ModerationItem`.
*   **`com.platform.developer.game.submitted.v1`** (от Developer Service)
    *   Описание: Разработчик отправил игру/обновление на модерацию.
    *   Логика обработки: Создать новый `ModerationItem`.
*   **`com.platform.payment.transaction.suspicious.v1`** (от Payment Service, Analytics Service)
    *   Описание: Обнаружена подозрительная финансовая активность.
    *   Логика обработки: Создать инцидент безопасности или задачу для администратора.
*   **`com.platform.system.health.degraded.v1`** (от системы мониторинга, концептуально)
    *   Описание: Обнаружена деградация производительности или сбой в одном из сервисов.
    *   Логика обработки: Создать алерт или тикет для технической команды.

## 6. Интеграции (Integrations)
См. `../../../../project_integrations.md` для общей карты и деталей.

### 6.1. Внутренние Микросервисы
Admin Service активно взаимодействует с большинством других микросервисов платформы для выполнения своих функций (см. предыдущий раздел "Обзор Сервиса" и `project_integrations.md`).

### 6.2. Внешние Системы
*   На данном этапе не предполагается прямых интеграций Admin Service с внешними российскими сервисами, выходящих за рамки тех, что обеспечиваются другими микросервисами (например, Notification Service для отправки SMS/Email через российских провайдеров).
*   Если в будущем потребуется интеграция с внешними системами аналитики контента, анти-фрод системами или государственными информационными системами (ГИС) для отчетности, это будет задокументировано дополнительно.

## 7. Конфигурация (Configuration)
Общие стандарты конфигурационных файлов (формат YAML, структура, управление переменными окружения и секретами) определены в `../../../../project_api_standards.md` (раздел 7) и `../../../../DOCUMENTATION_GUIDELINES.md` (раздел 6). Специфичные для Admin Service переменные и структура файла `configs/admin_config.yaml` приведены в разделе 3.1.4 примера API.

## 8. Обработка Ошибок (Error Handling)
*   Используются стандартные HTTP коды состояния для REST API и форматы ошибок согласно `../../../../project_api_standards.md`.
*   Внутренние ошибки сервиса логируются с высоким уровнем детализации, включая `trace_id`.
*   Пользователям админ-панели отображаются понятные сообщения об ошибках.

### 8.1. Распространенные Коды Ошибок (специфичные для Admin Service)
*   **`ADMIN_USER_NOT_FOUND`**: Административный пользователь не найден.
*   **`PLATFORM_USER_NOT_FOUND`**: Пользователь платформы не найден.
*   **`MODERATION_ITEM_NOT_FOUND`**: Элемент модерации не найден.
*   **`MODERATION_RULE_INVALID`**: Ошибка в конфигурации правила модерации.
*   **`SUPPORT_TICKET_NOT_FOUND`**: Тикет поддержки не найден.
*   **`KNOWLEDGE_BASE_ARTICLE_NOT_FOUND`**: Статья базы знаний не найдена.
*   **`PLATFORM_SETTING_NOT_FOUND`**: Настройка платформы не найдена.
*   **`PLATFORM_SETTING_VALIDATION_ERROR`**: Ошибка валидации значения настройки.
*   **`ACTION_NOT_PERMITTED_FOR_ROLE`**: Действие не разрешено для текущей роли администратора.
*   **`TARGET_SERVICE_UNAVAILABLE`**: Зависимый микросервис недоступен для выполнения операции.

## 9. Безопасность (Security)
См. `../../../../project_security_standards.md` для общих стандартов.

### 9.1. Аутентификация
*   Доступ к Admin API осуществляется по JWT, полученному административным пользователем через Auth Service. API Gateway валидирует токен и передает информацию об администраторе (ID, роли) в Admin Service.

### 9.2. Авторизация
*   Используется детализированная модель RBAC (Role-Based Access Control) для администраторов. Роли (например, `superuser`, `content_moderator`, `support_agent_l1`, `finance_admin`) определены в `../../../../project_roles_and_permissions.md`.
*   Каждый эндпоинт Admin Service проверяет наличие необходимых ролей/разрешений у аутентифицированного администратора.

### 9.3. Защита Данных
*   **ФЗ-152 "О персональных данных":** Admin Service обрабатывает значительный объем персональных данных пользователей платформы (при просмотре профилей, тикетов поддержки, модерации контента) и ПДн административных пользователей.
    *   Все операции с ПДн логируются в `AuditLogAdmin`.
    *   Доступ к ПДн строго регламентирован ролями и разрешениями.
    *   Обеспечиваются меры по защите ПДн при их отображении и обработке в админ-панели (например, маскирование части данных для некоторых ролей поддержки).
*   Шифрование при передаче (TLS) для всех API.
*   Чувствительные данные в конфигурации (пароли, ключи) хранятся в секретах.
*   Регулярный аудит безопасности и проверка на уязвимости.

### 9.4. Управление Секретами
*   Секреты сервиса (пароли к БД, ключи API для связи с другими сервисами, если есть) управляются через Kubernetes Secrets или HashiCorp Vault, согласно `../../../../project_security_standards.md`.

### 9.5. Аудит Действий
*   Все значимые действия, выполняемые администраторами через Admin Service, логируются в коллекцию `AuditLogAdmin` (MongoDB или PostgreSQL). Лог содержит информацию об администраторе, времени действия, типе действия, целевой сущности и измененных данных.

## 10. Развертывание (Deployment)
См. `../../../../project_deployment_standards.md` для общих стандартов.

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** `backend/admin-service/Dockerfile` (стандартный многоэтапный для Go).
*   **Helm-чарты/Kubernetes манифесты:** `deploy/charts/admin-service/` (включая Deployment, Service, ConfigMap, Secret, HPA).

### 10.2. Зависимости при Развертывании
*   PostgreSQL, MongoDB, Elasticsearch, Redis, Kafka.
*   Доступ к API других микросервисов платформы (Account, Auth, Catalog, etc.).
*   API Gateway для проксирования REST API.

### 10.3. CI/CD
*   Стандартный пайплайн CI/CD, включающий сборку, тестирование (unit, integration), статический анализ, сканирование безопасности, сборку Docker-образа, публикацию в реестр и развертывание на окружения с использованием Helm.
*   Автоматическое применение миграций БД (PostgreSQL) при деплое.

## 11. Мониторинг и Логирование (Logging and Monitoring)
См. `../../../../project_observability_standards.md` для общих стандартов.

### 11.1. Логирование
*   **Формат:** JSON (с использованием Zap).
*   **Уровни:** DEBUG, INFO, WARN, ERROR, FATAL.
*   **Ключевые поля:** `timestamp`, `level`, `service`, `instance`, `message`, `trace_id`, `span_id`, `admin_user_id` (если применимо), `target_entity_id`, `error_details`.
*   **Интеграция:** Fluent Bit для сбора логов и отправки в Elasticsearch/Loki.

### 11.2. Мониторинг
*   **Метрики (Prometheus):**
    *   `http_requests_total{handler="/platform-users", method="GET", status="200"}`
    *   `http_request_duration_seconds{handler="/moderation/items", method="POST"}`
    *   `moderation_queue_size_gauge{queue_id="game_reviews"}`
    *   `support_tickets_open_gauge{priority="high"}`
    *   `db_query_duration_seconds{db_instance="admin_postgres", operation="select_admin_users"}`
    *   `external_service_request_duration_seconds{service_name="account_service", operation="get_user_details"}`
    *   `kafka_messages_produced_total{topic="com.platform.admin.events.v1"}`
*   **Дашборды (Grafana):** Обзор состояния Admin Service, производительность API, размеры очередей модерации, количество и статус тикетов поддержки, активность администраторов.
*   **Алерты (AlertManager):**
    *   Высокий % ошибок API Admin Service.
    *   Большая задержка ответов API.
    *   Проблемы с подключением к зависимым БД (PostgreSQL, MongoDB, Elasticsearch) или Kafka.
    *   Аномальный рост очередей модерации или тикетов поддержки.
    *   Ошибки при выполнении критически важных административных задач.

### 11.3. Трассировка
*   **Интеграция:** OpenTelemetry SDK для Go.
*   **Экспорт:** Jaeger или другой совместимый коллектор.
*   **Контекст трассировки:** Передается через все HTTP запросы и Kafka события.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**:
    *   Время отклика API для интерактивных операций (получение списка, открытие сущности): P95 < 800 мс, P99 < 1500 мс.
    *   Время отклика API для операций изменения (создание, обновление, удаление): P95 < 500 мс, P99 < 1000 мс.
    *   Загрузка главной страницы админ-панели со всеми виджетами: P95 < 2 секунд.
*   **Надежность**:
    *   Доступность сервиса: 99.9%.
    *   RTO: < 1 часа для критических функций.
    *   RPO: < 15 минут для данных конфигураций и тикетов; < 1 час для логов аудита.
*   **Безопасность**: Соответствие `../../../../project_security_standards.md` и ФЗ-152. Логирование всех действий администраторов.
*   **Масштабируемость**: Горизонтальное масштабирование для поддержки до 1000 одновременных сессий администраторов/модераторов.
*   **Сопровождаемость**: Покрытие кода тестами > 80%. Актуальная документация.

## 13. Приложения (Appendices)
*   Детальные JSON схемы для API могут быть предоставлены в виде OpenAPI спецификации, генерируемой из кода или аннотаций.

## 14. Резервное Копирование и Восстановление (Backup and Recovery)

### 14.1. PostgreSQL (Административные данные, тикеты, настройки)
*   **Процедура резервного копирования:**
    *   **Логические бэкапы:** Ежедневный `pg_dump` для базы данных Admin Service.
    *   **Физические бэкапы (PITR):** Настроена непрерывная архивация WAL-сегментов, аналогично другим критичным PostgreSQL базам данных платформы.
    *   **Частота:** Базовый физический бэкап еженедельно, WAL-файлы архивируются непрерывно, логический бэкап ежедневно.
    *   **Хранение:** Бэкапы хранятся в S3-совместимом хранилище (шифрование, версионирование, другой регион). Срок хранения: полные бэкапы - 30 дней, WAL - 14 дней.
*   **Процедура восстановления:** Тестируется ежеквартально.
*   **RTO (Recovery Time Objective):** < 2 часов.
*   **RPO (Recovery Point Objective):** < 5 минут.

### 14.2. MongoDB (Логи аудита, детали модерации)
*   **Процедура резервного копирования:**
    *   Использование `mongodump` для создания логических бэкапов.
    *   **Частота:** Ежедневно для коллекции `audit_logs_admin` и `moderation_item_details`.
    *   **Хранение:** Бэкапы сжимаются, шифруются и сохраняются в S3-совместимом хранилище. Срок хранения: `audit_logs_admin` - до 1 года (согласно требованиям безопасности/комплаенса), `moderation_item_details` - 90 дней.
*   **Процедура восстановления:** Тестируется ежеквартально.
*   **RTO:** < 4 часов.
*   **RPO:** 24 часа (допустима потеря данных за последние сутки для этих типов данных, если более частое резервное копирование нецелесообразно).

### 14.3. Elasticsearch (Индексы для поиска)
*   **Процедура резервного копирования:**
    *   Использование Elasticsearch Snapshots для создания резервных копий индексов.
    *   **Частота:** Ежедневно для всех поисковых индексов Admin Service.
    *   **Хранение:** Снапшоты хранятся в S3-совместимом репозитории. Срок хранения - 14 дней.
*   **Процедура восстановления:**
    *   Восстановление из снапшота. В случае полной потери данных, возможна переиндексация из PostgreSQL и MongoDB, но это займет значительно больше времени.
    *   Тестируется ежеквартально.
*   **RTO:** < 3 часов (из снапшота). > 12 часов (при переиндексации).
*   **RPO:** 24 часа.

### 14.4. Redis (Кэш)
*   Данные в Redis для Admin Service в основном являются кэшем или временными данными. Специализированное резервное копирование не является критичным, так как данные могут быть восстановлены из основных хранилищ (PostgreSQL, MongoDB) или пересозданы. Стандартные механизмы RDB snapshots и AOF Redis могут быть включены для ускорения восстановления после перезапуска.

### 14.5. Общая стратегия
*   Резервное копирование и восстановление Admin Service являются частью общей стратегии обеспечения непрерывности бизнеса платформы.
*   Процедуры документированы и регулярно пересматриваются.
*   Мониторинг процессов бэкапа.
*   Общие принципы резервного копирования для различных СУБД описаны в `../../../../project_database_structure.md`.

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [Подача разработчиком новой игры на модерацию](../../../../project_workflows/game_submission_flow.md)
*   [Обработка жалобы пользователя на контент] <!-- Workflow будет создан и описан в project_workflows/user_complaint_handling_flow.md -->
*   [Процесс решения тикета поддержки] <!-- Workflow будет создан и описан в project_workflows/support_ticket_resolution_flow.md -->

---
*Этот документ является отправной точкой и должен регулярно обновляться по мере развития сервиса.*
