<!-- backend\admin-service\docs\README.md -->
# Спецификация Микросервиса: Admin Service

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16

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
    AdminPanelUI --> AdminServiceAPI{Admin Service REST API}

    subgraph AdminService ["Admin Service (Чистая Архитектура)"]
        direction LR

        AdminServiceAPI --> PresentationLayer[Presentation Layer (HTTP Handlers)]
        PresentationLayer --> ApplicationLayer[Application Layer (Use Cases / Services)]
        ApplicationLayer --> DomainLayer[Domain Layer (Entities, Business Logic)]
        ApplicationLayer -- Интерфейсы репозиториев --> InfrastructureLayer[Infrastructure Layer]
        DomainLayer -- Интерфейсы репозиториев --> InfrastructureLayer

        subgraph Modules [Функциональные Модули (в Application/Domain Layers)]
            direction TB
            Mod[Модерация]
            UserMgmt[Управление Пользователями]
            AdminUserMgmt[Управление Администраторами]
            Support[Тех. Поддержка]
            SettingsMgmt[Настройки Платформы]
            AnalyticsMod[Аналитика (доступ)]
            MarketingMod[Маркетинг]
            SecurityMod[Безопасность (мониторинг)]
        end
        ApplicationLayer --> Modules

        subgraph DataStoresAdapter ["Адаптеры к Хранилищам Данных (в Infrastructure Layer)"]
            PostgresDB[(PostgreSQL)]
            MongoDB[(MongoDB)]
            Elasticsearch[(Elasticsearch)]
            RedisCache[(Redis)]
        end
        InfrastructureLayer --> DataStoresAdapter

        subgraph ExternalServicesAdapter ["Клиенты Внешних Сервисов (в Infrastructure Layer)"]
            AccountServiceExt[Account Service Client]
            CatalogServiceExt[Catalog Service Client]
            AuthServiceExt[Auth Service Client]
            NotificationServiceExt[Notification Service Client]
            AnalyticsServiceExt[Analytics Service Client]
            PaymentServiceExt[Payment Service Client]
        end
        InfrastructureLayer --> ExternalServicesAdapter

        subgraph MessageQueueAdapter ["Адаптеры к Брокеру Сообщений (в Infrastructure Layer)"]
            KafkaClient[Kafka Producer/Consumer]
        end
        InfrastructureLayer --> MessageQueueAdapter
    end

    ExternalServicesAdapter --> AccountService[Account Service (gRPC/REST)]
    ExternalServicesAdapter --> CatalogService[Catalog Service (gRPC/REST)]
    ExternalServicesAdapter --> AuthService[Auth Service (gRPC/REST)]
    ExternalServicesAdapter --> NotificationService[Notification Service (gRPC/Kafka)]
    ExternalServicesAdapter --> AnalyticsService[Analytics Service (REST/Kafka)]
    ExternalServicesAdapter --> PaymentService[Payment Service (gRPC/REST)]

    MessageQueueAdapter --> KafkaBroker[Apache Kafka]
    KafkaBroker --> MessageQueueAdapter

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

### 4.1. Основные Сущности (PostgreSQL)
*   **`AdminUser` (Администратор)**: Учетная запись администратора/сотрудника.
    *   `id` (UUID, PK): Уникальный идентификатор. **Обязательность: Да (генерируется БД).**
    *   `username` (VARCHAR, UK): Логин администратора. **Обязательность: Да.**
    *   `password_hash` (VARCHAR): Хеш пароля. **Обязательность: Да.**
    *   `email` (VARCHAR, UK): Email администратора. **Обязательность: Да.**
    *   `full_name` (VARCHAR): Полное имя. **Обязательность: Нет.**
    *   `roles` (ARRAY_TEXT): Список ролей. **Обязательность: Да.**
    *   `is_active` (BOOLEAN): Активен ли пользователь. **Обязательность: Да (DEFAULT TRUE).**
    *   `created_at` (TIMESTAMPTZ): Время создания. **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. **Обязательность: Да (генерируется БД).**
*   **`SupportTicket` (Тикет поддержки)**: Запрос пользователя.
    *   `id` (UUID, PK): Уникальный идентификатор. **Обязательность: Да (генерируется БД).**
    *   `subject` (VARCHAR): Тема тикета. **Обязательность: Да.**
    *   `description` (TEXT): Описание проблемы. **Обязательность: Да.**
    *   `platform_user_id` (UUID, FK): ID пользователя платформы. **Обязательность: Да.**
    *   `assignee_admin_id` (UUID, FK на `AdminUser`, Nullable): ID назначенного администратора. **Обязательность: Нет.**
    *   `category_id` (UUID, FK на `SupportTicketCategory`): ID категории тикета. **Обязательность: Да.**
    *   `status` (VARCHAR ENUM): Статус тикета (new, open, pending_user, resolved, closed). **Обязательность: Да (DEFAULT 'new').**
    *   `priority` (VARCHAR ENUM): Приоритет (low, medium, high, urgent). **Обязательность: Да (DEFAULT 'medium').**
    *   `custom_fields` (JSONB, Nullable): Дополнительные поля. **Обязательность: Нет.**
    *   `created_at` (TIMESTAMPTZ): Время создания. **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. **Обязательность: Да (генерируется БД).**
    *   `resolved_at` (TIMESTAMPTZ, Nullable): Время решения. **Обязательность: Нет.**
    *   `closed_at` (TIMESTAMPTZ, Nullable): Время закрытия. **Обязательность: Нет.**
*   **`SupportTicketResponse` (Ответ в тикете)**: Ответ/комментарий в тикете.
    *   `id` (UUID, PK): Уникальный идентификатор. **Обязательность: Да (генерируется БД).**
    *   `ticket_id` (UUID, FK): ID тикета. **Обязательность: Да.**
    *   `admin_user_id` (UUID, FK на `AdminUser`, Nullable): ID администратора, давшего ответ. **Обязательность: Нет (если ответ от пользователя).**
    *   `platform_user_id` (UUID, Nullable): ID пользователя платформы, давшего ответ. **Обязательность: Нет (если ответ от админа).**
    *   `body` (TEXT): Текст ответа. **Обязательность: Да.**
    *   `is_internal_note` (BOOLEAN): Является ли ответ внутренней заметкой. **Обязательность: Да (DEFAULT FALSE).**
    *   `created_at` (TIMESTAMPTZ): Время создания. **Обязательность: Да (генерируется БД).**
*   **`SupportTicketCategory` (Категория тикета)**: Категория для классификации тикетов.
    *   `id` (UUID, PK): Уникальный идентификатор. **Обязательность: Да (генерируется БД).**
    *   `name` (VARCHAR, UK): Название категории. **Обязательность: Да.**
    *   `description` (TEXT, Nullable): Описание. **Обязательность: Нет.**
    *   `is_active` (BOOLEAN): Активна ли категория. **Обязательность: Да (DEFAULT TRUE).**
*   **`KnowledgeBaseArticle` (Статья базы знаний)**: Статья для самопомощи.
    *   `id` (UUID, PK): Уникальный идентификатор. **Обязательность: Да (генерируется БД).**
    *   `title` (VARCHAR): Заголовок статьи. **Обязательность: Да.**
    *   `content_markdown` (TEXT): Содержимое статьи в Markdown. **Обязательность: Да.**
    *   `category_id` (UUID, FK на `KnowledgeBaseCategory`): ID категории. **Обязательность: Да.**
    *   `author_admin_id` (UUID, FK на `AdminUser`): ID автора-администратора. **Обязательность: Да.**
    *   `is_published` (BOOLEAN): Опубликована ли статья. **Обязательность: Да (DEFAULT FALSE).**
    *   `view_count` (INTEGER): Количество просмотров. **Обязательность: Да (DEFAULT 0).**
    *   `language_code` (VARCHAR): Код языка (ru, en). **Обязательность: Да.**
    *   `created_at` (TIMESTAMPTZ): Время создания. **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. **Обязательность: Да (генерируется БД).**
*   **`KnowledgeBaseCategory` (Категория статьи БЗ)**: Категория для статей.
    *   `id` (UUID, PK): Уникальный идентификатор. **Обязательность: Да (генерируется БД).**
    *   `name` (VARCHAR, UK): Название категории. **Обязательность: Да.**
    *   `description` (TEXT, Nullable): Описание. **Обязательность: Нет.**
    *   `parent_category_id` (UUID, FK, Nullable): ID родительской категории. **Обязательность: Нет.**
*   **`ModerationItem` (Элемент модерации - основная запись)**: Объект, требующий модерации.
    *   `id` (UUID, PK): Уникальный идентификатор. **Обязательность: Да (генерируется БД).**
    *   `item_reference_id` (VARCHAR): ID контента в другом сервисе. **Обязательность: Да.**
    *   `item_type` (VARCHAR ENUM): Тип контента. **Обязательность: Да.**
    *   `status` (VARCHAR ENUM): Статус модерации. **Обязательность: Да (DEFAULT 'pending_auto').**
    *   `reason_for_submission` (TEXT): Причина отправки на модерацию. **Обязательность: Нет.**
    *   `content_snapshot_summary` (TEXT): Краткий снимок контента. **Обязательность: Нет.**
    *   `submitter_user_id` (UUID, Nullable): ID пользователя, отправившего контент. **Обязательность: Нет.**
    *   `assigned_moderator_id` (UUID, FK на `AdminUser`, Nullable): ID назначенного модератора. **Обязательность: Нет.**
    *   `created_at` (TIMESTAMPTZ): Время создания. **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. **Обязательность: Да (генерируется БД).**
    *   `decision_at` (TIMESTAMPTZ, Nullable): Время принятия решения. **Обязательность: Нет.**
    *   `decision` (VARCHAR ENUM, Nullable): Решение (approved, rejected). **Обязательность: Нет.**
    *   `decision_reason_code` (VARCHAR, Nullable): Код причины решения. **Обязательность: Нет.**
    *   `moderator_comment` (TEXT, Nullable): Комментарий модератора. **Обязательность: Нет.**
*   **`ModerationRule` (Правило модерации)**: Автоматическое правило.
    *   `id` (UUID, PK): Уникальный идентификатор. **Обязательность: Да (генерируется БД).**
    *   `name` (VARCHAR, UK): Название правила. **Обязательность: Да.**
    *   `description` (TEXT): Описание правила. **Обязательность: Нет.**
    *   `item_type` (VARCHAR ENUM): Тип контента, к которому применяется правило. **Обязательность: Да.**
    *   `condition_script` (TEXT): Условие срабатывания правила. **Обязательность: Да.**
    *   `action_to_take` (VARCHAR ENUM): Действие при срабатывании. **Обязательность: Да.**
    *   `priority` (INTEGER): Приоритет правила. **Обязательность: Да (DEFAULT 0).**
    *   `is_active` (BOOLEAN): Активно ли правило. **Обязательность: Да (DEFAULT TRUE).**
    *   `created_at` (TIMESTAMPTZ): Время создания. **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. **Обязательность: Да (генерируется БД).**
*   **`PlatformSetting` (Настройка платформы)**: Глобальный параметр конфигурации.
    *   `key` (VARCHAR, PK): Ключ настройки. **Обязательность: Да.**
    *   `value` (TEXT): Значение настройки. **Обязательность: Да.**
    *   `description` (TEXT, Nullable): Описание. **Обязательность: Нет.**
    *   `value_type` (VARCHAR ENUM): Тип значения (string, integer, boolean, json). **Обязательность: Да.**
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. **Обязательность: Да (генерируется БД).**
    *   `updated_by_admin_id` (UUID, FK на `AdminUser`): ID администратора, обновившего настройку. **Обязательность: Да.**

### 4.2. Схема Базы Данных (PostgreSQL)
*   ERD Диаграмма (PostgreSQL):
    ```mermaid
    erDiagram
        ADMIN_USER ||--o{ SUPPORT_TICKET : "assigns/handles"
        ADMIN_USER ||--o{ MODERATION_ITEM : "moderates"
        ADMIN_USER ||--o{ KNOWLEDGE_BASE_ARTICLE : "authors"
        ADMIN_USER ||--o{ PLATFORM_SETTING : "updates"
        ADMIN_USER ||--o{ SUPPORT_TICKET_RESPONSE : "responds"

        SUPPORT_TICKET }o--|| SUPPORT_TICKET_CATEGORY : "belongs to"
        SUPPORT_TICKET ||--o{ SUPPORT_TICKET_RESPONSE : "has many"

        KNOWLEDGE_BASE_ARTICLE }o--|| KNOWLEDGE_BASE_CATEGORY : "belongs to"

        MODERATION_ITEM ||--o{ MODERATION_RULE : "can be affected by"

        ADMIN_USER {
            UUID id PK
            VARCHAR username UK
            VARCHAR password_hash
            VARCHAR email UK
            VARCHAR full_name
            ARRAY_TEXT roles "e.g. ['superuser', 'content_moderator']"
            BOOLEAN is_active
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }

        SUPPORT_TICKET {
            UUID id PK
            VARCHAR subject
            TEXT description
            UUID platform_user_id "FK to Platform User (Account Service)"
            UUID assignee_admin_id FK "Nullable, Refers to ADMIN_USER(id)"
            UUID category_id FK "Refers to SUPPORT_TICKET_CATEGORY(id)"
            VARCHAR status "ENUM('new', 'open', 'pending_user', 'pending_agent', 'resolved', 'closed')"
            VARCHAR priority "ENUM('low', 'medium', 'high', 'urgent')"
            JSONB custom_fields "Nullable"
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
            TIMESTAMPTZ resolved_at "Nullable"
            TIMESTAMPTZ closed_at "Nullable"
        }

        SUPPORT_TICKET_RESPONSE {
            UUID id PK
            UUID ticket_id FK "Refers to SUPPORT_TICKET(id)"
            UUID admin_user_id FK "Nullable, if response from admin, Refers to ADMIN_USER(id)"
            UUID platform_user_id "Nullable, if response from platform user"
            TEXT body
            BOOLEAN is_internal_note "Default false"
            TIMESTAMPTZ created_at
        }

        SUPPORT_TICKET_CATEGORY {
            UUID id PK
            VARCHAR name UK
            TEXT description "Nullable"
            BOOLEAN is_active "Default true"
        }

        KNOWLEDGE_BASE_ARTICLE {
            UUID id PK
            VARCHAR title
            TEXT content_markdown
            UUID category_id FK "Refers to KNOWLEDGE_BASE_CATEGORY(id)"
            UUID author_admin_id FK "Refers to ADMIN_USER(id)"
            BOOLEAN is_published "Default false"
            INTEGER view_count "Default 0"
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
            VARCHAR language_code "e.g. 'ru', 'en'"
        }

        KNOWLEDGE_BASE_CATEGORY {
            UUID id PK
            VARCHAR name UK
            TEXT description "Nullable"
            UUID parent_category_id FK "Nullable, self-referencing for subcategories"
        }

        MODERATION_ITEM {
            UUID id PK
            VARCHAR item_reference_id "ID of the content in another service, e.g. review_id, game_id"
            VARCHAR item_type "ENUM('game_review', 'user_comment', 'game_submission', 'user_profile_customization')"
            VARCHAR status "ENUM('pending_auto', 'pending_manual', 'approved', 'rejected', 'escalated')"
            TEXT reason_for_submission "Why this item is in moderation"
            TEXT content_snapshot_summary "Snapshot of content being moderated"
            UUID submitter_user_id "Nullable, FK to Platform User"
            UUID assigned_moderator_id FK "Nullable, Refers to ADMIN_USER(id)"
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
            TIMESTAMPTZ decision_at "Nullable"
            VARCHAR decision "Nullable, ENUM('approved', 'rejected')"
            VARCHAR decision_reason_code "Nullable"
            TEXT moderator_comment "Nullable"
        }

        MODERATION_RULE {
            UUID id PK
            VARCHAR name UK
            TEXT description
            VARCHAR item_type "ENUM matching MODERATION_ITEM.item_type or 'any'"
            TEXT condition_script "e.g., Groovy, Python, or specific DSL for rule engine"
            VARCHAR action_to_take "ENUM('auto_approve', 'auto_reject', 'escalate_to_human', 'flag_for_review')"
            INTEGER priority "Higher priority rules evaluated first"
            BOOLEAN is_active "Default true"
            TIMESTAMPTZ created_at
            TIMESTAMPTZ updated_at
        }

        PLATFORM_SETTING {
            VARCHAR key PK "e.g. 'user_registration.allow_new_registrations'"
            TEXT value
            TEXT description "Nullable"
            VARCHAR value_type "ENUM('string', 'integer', 'boolean', 'json')"
            TIMESTAMPTZ updated_at
            UUID updated_by_admin_id FK "Refers to ADMIN_USER(id)"
        }
    ```
*   Описание основных таблиц и индексов (DDL-подобное описание для PostgreSQL):
    ```sql
    -- Таблица: admin_users
    CREATE TABLE admin_users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        username VARCHAR(100) NOT NULL UNIQUE,
        password_hash VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE,
        full_name VARCHAR(255),
        roles TEXT[] NOT NULL, -- Массив ролей, например: '{superuser,content_moderator}'
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
    CREATE INDEX idx_admin_users_roles ON admin_users USING GIN (roles);

    -- Таблица: support_ticket_categories
    CREATE TABLE support_ticket_categories (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(255) NOT NULL UNIQUE,
        description TEXT,
        is_active BOOLEAN NOT NULL DEFAULT TRUE
    );

    -- Таблица: support_tickets
    CREATE TABLE support_tickets (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        subject VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        platform_user_id UUID, -- ID пользователя из Account Service
        assignee_admin_id UUID REFERENCES admin_users(id) ON DELETE SET NULL,
        category_id UUID REFERENCES support_ticket_categories(id) ON DELETE RESTRICT,
        status VARCHAR(50) NOT NULL DEFAULT 'new', -- new, open, pending_user, pending_agent, resolved, closed
        priority VARCHAR(50) NOT NULL DEFAULT 'medium', -- low, medium, high, urgent
        custom_fields JSONB,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        resolved_at TIMESTAMPTZ,
        closed_at TIMESTAMPTZ
    );
    CREATE INDEX idx_support_tickets_status ON support_tickets(status);
    CREATE INDEX idx_support_tickets_assignee_admin_id ON support_tickets(assignee_admin_id);
    CREATE INDEX idx_support_tickets_platform_user_id ON support_tickets(platform_user_id);

    -- (Другие таблицы: support_ticket_responses, knowledge_base_categories, knowledge_base_articles, moderation_items, moderation_rules, platform_settings - аналогично с полями, индексами и FK)
    ```

### 4.3. MongoDB Коллекции
*   **`audit_log_admin`**: Хранит документы `AuditLogAdmin`. Поля: `_id` (ObjectId), `admin_user_id` (UUID), `admin_username` (String), `timestamp` (Date), `action_type` (String, например, "update_user_status", "resolve_ticket"), `target_entity_type` (String, например, "PlatformUser", "SupportTicket"), `target_entity_id` (String), `details_before` (Object, опционально), `details_after` (Object, опционально), `ip_address` (String), `user_agent` (String).
    *   Индексы: `{admin_user_id: 1, timestamp: -1}`, `{timestamp: -1}`, `{action_type: 1}`, `{target_entity_id: 1}`.
*   **`moderation_item_details`**: Хранит расширенные данные для `ModerationItem`. Поля: `_id` (UUID, соответствует `ModerationItem.id` из PostgreSQL), `full_content_snapshot` (Object/String), `external_analysis_reports` (Array of Objects), `attachments` (Array of Objects).
    *   Индексы: `{_id: 1}`.

### 4.4. Elasticsearch Индексы
*   **`support_tickets_idx`**: Поля: `id` (keyword), `subject` (text), `description` (text), `responses_text` (text, агрегированные ответы), `platform_user_id` (keyword), `assignee_admin_username` (keyword), `category_name` (keyword), `status` (keyword), `priority` (keyword), `tags` (keyword[]), `created_at` (date), `updated_at` (date).
*   **`knowledge_base_articles_idx`**: Поля: `id` (keyword), `title` (text), `content_text` (text, из markdown), `tags` (keyword[]), `category_name` (keyword), `author_username` (keyword), `language_code` (keyword), `is_published` (boolean), `created_at` (date).
*   **`admin_audit_log_idx`** (если логи аудита также индексируются для детального поиска): Поля `id` (keyword, из MongoDB _id), `admin_username` (keyword), `action_type` (keyword), `target_entity_type` (keyword), `target_entity_id` (keyword), `timestamp` (date), `details_text` (text, для поиска по содержимому изменений).
*   **`moderation_items_idx`** (если требуется сложный поиск по элементам модерации): Поля `id` (keyword, из PostgreSQL ModerationItem.id), `item_type` (keyword), `status` (keyword), `content_summary_text` (text), `submitter_user_id` (keyword), `assigned_moderator_username` (keyword), `created_at` (date).

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
Общие стандарты конфигурационных файлов (формат YAML, структура, управление переменными окружения и секретами) определены в `../../../../project_api_standards.md` (раздел 7) и `../../../../DOCUMENTATION_GUIDELINES.md` (раздел 6).

### 7.1. Переменные Окружения (Примеры)
*   `ADMIN_HTTP_PORT`: Порт для REST API Admin Service (например, `8085`)
*   `ADMIN_GRPC_PORT`: Порт для внутреннего gRPC API, если используется (например, `9095`)
*   `POSTGRES_DSN_ADMIN`: Строка подключения к PostgreSQL для Admin Service.
*   `MONGODB_URI_ADMIN`: Строка подключения к MongoDB для Admin Service.
*   `ELASTICSEARCH_URLS_ADMIN`: URL(ы) Elasticsearch для Admin Service.
*   `REDIS_ADDR_ADMIN`: Адрес Redis для Admin Service.
*   `KAFKA_BROKERS_ADMIN`: Список брокеров Kafka.
*   `JWT_PUBLIC_KEY_PATH`: Путь к публичному ключу для валидации JWT токенов администраторов (обычно тот же, что и для пользователей, если токены выдаются одним Auth Service).
*   `LOG_LEVEL_ADMIN`: Уровень логирования для Admin Service (например, `info`, `debug`).
*   `OTEL_EXPORTER_JAEGER_ENDPOINT_ADMIN`: Endpoint для экспорта трейсов в Jaeger.
*   `ACCOUNT_SERVICE_GRPC_ADDR`: Адрес gRPC Account Service.
*   `CATALOG_SERVICE_GRPC_ADDR`: Адрес gRPC Catalog Service.
*   `AUTH_SERVICE_GRPC_ADDR`: Адрес gRPC Auth Service.
*   `DEFAULT_MODERATION_QUEUE_SIZE`: Максимальный размер очереди модерации по умолчанию.

### 7.2. Файлы Конфигурации (`configs/admin_service_config.yaml`)
*   Структура:
    ```yaml
    http_server:
      port: ${ADMIN_HTTP_PORT:"8085"}
      timeout_seconds: 60
    # grpc_server: # Если используется внутренний gRPC
    #   port: ${ADMIN_GRPC_PORT:"9095"}
    #   timeout_seconds: 60
    postgres:
      dsn: ${POSTGRES_DSN_ADMIN}
      pool_max_conns: 15
    mongodb:
      uri: ${MONGODB_URI_ADMIN}
      database_name: "admin_service_db"
      pool_max_size: 10
    elasticsearch:
      urls: ${ELASTICSEARCH_URLS_ADMIN} # "http://es1:9200,http://es2:9200"
      username: ${ELASTICSEARCH_USER:""}
      password: ${ELASTICSEARCH_PASSWORD:""}
    redis:
      address: ${REDIS_ADDR_ADMIN}
      password: ${REDIS_PASSWORD_ADMIN:""}
      db: ${REDIS_DB_ADMIN:2} # Отдельная база Redis для Admin Service
    kafka:
      brokers: ${KAFKA_BROKERS_ADMIN}
      producer_topics:
        admin_events: ${KAFKA_TOPIC_ADMIN_EVENTS:"com.platform.admin.events.v1"}
      consumer_topics:
        user_complaints: ${KAFKA_TOPIC_USER_COMPLAINTS:"com.platform.user.complaint.submitted.v1"}
        # ... другие потребляемые топики
      consumer_group: ${KAFKA_CONSUMER_GROUP_ADMIN:"admin-service-group"}
    logging:
      level: ${LOG_LEVEL_ADMIN:"info"}
      format: "json"
    security:
      jwt_public_key_path: ${JWT_PUBLIC_KEY_PATH}
      # Дополнительные параметры безопасности, например, лимиты на загрузку файлов
    otel:
      exporter_jaeger_endpoint: ${OTEL_EXPORTER_JAEGER_ENDPOINT_ADMIN}
      service_name: "admin-service"
    integrations:
      account_service_grpc_addr: ${ACCOUNT_SERVICE_GRPC_ADDR}
      catalog_service_grpc_addr: ${CATALOG_SERVICE_GRPC_ADDR}
      # ... адреса других сервисов
    default_limits:
      moderation_queue_size: ${DEFAULT_MODERATION_QUEUE_SIZE:1000}
      max_login_attempts_admin: 5
    ```

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

## 14. Пользовательские Сценарии (User Flows)

В этом разделе описаны ключевые пользовательские сценарии, в которых участвуют администраторы платформы и сотрудники поддержки, используя Admin Service.

### 14.1. Вход Администратора в Админ-Панель
*   **Описание:** Администратор входит в систему через специальный интерфейс админ-панели, используя свои учетные данные. Auth Service проверяет их и выдает JWT, который затем используется для доступа к Admin Service.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor Admin as Администратор
        participant AdminPanelUI as Админ-Панель (Frontend)
        participant APIGW as API Gateway
        participant AuthSvc as Auth Service
        participant AdminSvc as Admin Service

        Admin->>AdminPanelUI: Вводит логин/пароль
        AdminPanelUI->>APIGW: POST /api/v1/auth/admin/login (credentials)
        APIGW->>AuthSvc: Forward admin login request
        AuthSvc->>AuthSvc: Проверка учетных данных администратора
        alt Успешная аутентификация
            AuthSvc->>AuthSvc: Генерация JWT (с ролями администратора)
            AuthSvc-->>APIGW: HTTP 200 OK (JWT)
            APIGW-->>AdminPanelUI: HTTP 200 OK (JWT)
            AdminPanelUI->>AdminPanelUI: Сохранение JWT, загрузка начальной страницы
            AdminPanelUI->>APIGW: GET /api/v1/admin/dashboard-summary (Authorization: Bearer JWT)
            APIGW->>AdminSvc: Forward request (X-Admin-UserID, X-Admin-Roles)
            AdminSvc-->>APIGW: HTTP 200 OK (данные для дашборда)
            APIGW-->>AdminPanelUI: HTTP 200 OK
            AdminPanelUI-->>Admin: Отображение дашборда
        else Ошибка аутентификации
            AuthSvc-->>APIGW: HTTP 401 Unauthorized
            APIGW-->>AdminPanelUI: HTTP 401 Unauthorized
            AdminPanelUI-->>Admin: Сообщение об ошибке входа
        end
    ```

### 14.2. Управление Пользователями Платформы (Поиск, Просмотр, Блокировка)
*   **Описание:** Администратор ищет пользователя платформы, просматривает его детали и при необходимости блокирует его аккаунт. Это включает взаимодействие с Account Service.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor PlatformAdmin as Администратор Платформы
        participant AdminPanelUI as Админ-Панель
        participant AdminSvc as Admin Service (REST API)
        participant AccountSvc as Account Service (gRPC/REST)
        participant KafkaBus as Kafka

        PlatformAdmin->>AdminPanelUI: Поиск пользователя (например, по email)
        AdminPanelUI->>AdminSvc: GET /api/v1/admin/platform-users?search=user@example.com
        AdminSvc->>AccountSvc: (gRPC) GetUsers(filter_by_email="user@example.com")
        AccountSvc-->>AdminSvc: UserListResponse
        AdminSvc-->>AdminPanelUI: Список пользователей

        PlatformAdmin->>AdminPanelUI: Выбирает пользователя для просмотра деталей
        AdminPanelUI->>AdminSvc: GET /api/v1/admin/platform-users/{user_id}
        AdminSvc->>AccountSvc: (gRPC) GetUserProfile(user_id) & GetAccountInfo(user_id)
        AccountSvc-->>AdminSvc: UserProfileResponse & AccountInfoResponse
        AdminSvc-->>AdminPanelUI: Детальная информация о пользователе

        PlatformAdmin->>AdminPanelUI: Нажимает "Заблокировать пользователя" (указывает причину)
        AdminPanelUI->>AdminSvc: PUT /api/v1/admin/platform-users/{user_id}/status (payload: {status: "blocked", reason: "Spamming"})
        AdminSvc->>AccountSvc: (gRPC) UpdateAccountStatus(user_id, new_status="blocked", reason="Spamming")
        AccountSvc-->>AdminSvc: Success/Failure
        alt Успешная блокировка
            AdminSvc->>AdminSvc: Запись в AuditLogAdmin
            AdminSvc-->>KafkaBus: Publish `com.platform.admin.user.status.updated.v1` (userId, newStatus="blocked")
            AdminSvc-->>AdminPanelUI: HTTP 200 OK (статус обновлен)
            AdminPanelUI-->>PlatformAdmin: Уведомление об успешной блокировке
        else Ошибка блокировки
            AdminSvc-->>AdminPanelUI: HTTP Error (например, 404, 500)
            AdminPanelUI-->>PlatformAdmin: Сообщение об ошибке
        end
    ```

### 14.3. Модерация Пользовательского Контента
*   **Описание:** Модератор просматривает очередь контента, ожидающего модерации (например, отзывы об играх), принимает решение (одобрить/отклонить) и указывает причину.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor Moderator as Модератор
        participant AdminPanelUI as Админ-Панель
        participant AdminSvc as Admin Service
        participant KafkaBus as Kafka

        Moderator->>AdminPanelUI: Открывает очередь модерации "Отзывы об играх"
        AdminPanelUI->>AdminSvc: GET /api/v1/admin/moderation/queues/game_reviews/items
        AdminSvc->>AdminSvc: Запрос к БД (PostgreSQL/MongoDB) для получения элементов
        AdminSvc-->>AdminPanelUI: Список элементов для модерации

        Moderator->>AdminPanelUI: Выбирает элемент, просматривает контент
        AdminPanelUI->>AdminSvc: GET /api/v1/admin/moderation/items/{item_id}
        AdminSvc-->>AdminPanelUI: Детали элемента (снапшот контента)

        Moderator->>AdminPanelUI: Принимает решение "Отклонить" (причина: "Спам")
        AdminPanelUI->>AdminSvc: POST /api/v1/admin/moderation/items/{item_id}/decisions (payload: {decision: "rejected", reason_code: "spam", comment: "..."})
        AdminSvc->>AdminSvc: Обновление статуса ModerationItem в БД
        AdminSvc->>AdminSvc: Запись в AuditLogAdmin
        AdminSvc-->>KafkaBus: Publish `com.platform.admin.content.moderated.v1` (itemId, decision="rejected", reason="spam")
        AdminSvc-->>AdminPanelUI: HTTP 200 OK (решение принято)
        AdminPanelUI-->>Moderator: Элемент удален из очереди, уведомление об успехе
    ```

### 14.4. Управление Тикетом Поддержки
*   **Описание:** Агент поддержки просматривает назначенный ему тикет, отвечает пользователю и закрывает тикет.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor SupportAgent as Агент Поддержки
        participant AdminPanelUI as Админ-Панель
        participant AdminSvc as Admin Service
        participant NotificationSvc as Notification Service (через Kafka или gRPC)
        participant KafkaBus as Kafka

        SupportAgent->>AdminPanelUI: Открывает список "Мои открытые тикеты"
        AdminPanelUI->>AdminSvc: GET /api/v1/admin/support/tickets?assignee_id=current_agent&status=open
        AdminSvc->>AdminSvc: Запрос к БД (PostgreSQL)
        AdminSvc-->>AdminPanelUI: Список тикетов

        SupportAgent->>AdminPanelUI: Выбирает тикет, читает историю
        AdminPanelUI->>AdminSvc: GET /api/v1/admin/support/tickets/{ticket_id}
        AdminSvc-->>AdminPanelUI: Детали тикета

        SupportAgent->>AdminPanelUI: Пишет ответ пользователю
        AdminPanelUI->>AdminSvc: POST /api/v1/admin/support/tickets/{ticket_id}/responses (payload: {body: "...", is_internal_note: false})
        AdminSvc->>AdminSvc: Сохранение ответа в БД
        AdminSvc->>AdminSvc: Обновление статуса тикета (например, на 'pending_user')
        AdminSvc-->>KafkaBus: Publish `com.platform.admin.support.ticket.status.updated.v1` (ticketId, newStatus='pending_user')
        AdminSvc->>NotificationSvc: Запрос на уведомление пользователя об ответе
        AdminSvc-->>AdminPanelUI: HTTP 201 Created (ответ добавлен)

        SupportAgent->>AdminPanelUI: Решает закрыть тикет (после получения подтверждения от пользователя или по регламенту)
        AdminPanelUI->>AdminSvc: PUT /api/v1/admin/support/tickets/{ticket_id}/status (payload: {status: "resolved"})
        AdminSvc->>AdminSvc: Обновление статуса тикета в БД
        AdminSvc-->>KafkaBus: Publish `com.platform.admin.support.ticket.status.updated.v1` (ticketId, newStatus='resolved')
        AdminSvc-->>AdminPanelUI: HTTP 200 OK
        AdminPanelUI-->>SupportAgent: Тикет обновлен
    ```

### 14.5. Конфигурирование Настройки Платформы
*   **Описание:** Главный администратор изменяет глобальную настройку платформы, например, включает или выключает возможность регистрации новых пользователей.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor SuperAdmin as Главный Администратор
        participant AdminPanelUI as Админ-Панель
        participant AdminSvc as Admin Service
        participant KafkaBus as Kafka

        SuperAdmin->>AdminPanelUI: Переходит в раздел "Настройки Платформы" -> "Регистрация"
        AdminPanelUI->>AdminSvc: GET /api/v1/admin/platform-settings?group=user_registration
        AdminSvc->>AdminSvc: Запрос к БД (PostgreSQL/Redis)
        AdminSvc-->>AdminPanelUI: Текущие настройки регистрации

        SuperAdmin->>AdminPanelUI: Изменяет "Разрешить новые регистрации" на "false"
        AdminPanelUI->>AdminSvc: PUT /api/v1/admin/platform-settings (payload: {"user_registration": {"allow_new_registrations": false}})
        AdminSvc->>AdminSvc: Валидация и сохранение PlatformSetting в БД
        AdminSvc->>AdminSvc: Запись в AuditLogAdmin
        AdminSvc->>RedisCache: Очистка кэша для данной настройки (если кэшируется)
        AdminSvc-->>KafkaBus: Publish `com.platform.admin.platform.setting.updated.v1` (key="user_registration.allow_new_registrations", newValue=false)
        AdminSvc-->>AdminPanelUI: HTTP 200 OK (настройки обновлены)
        AdminPanelUI-->>SuperAdmin: Уведомление об успешном изменении
    ```

### 14.6. Просмотр Отчета в Админ-Аналитике
*   **Описание:** Администратор просматривает отчет по продажам игр за последний месяц. Admin Service проксирует или формирует запрос к Analytics Service.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor AnalyticsAdmin as Администратор Аналитики
        participant AdminPanelUI as Админ-Панель
        participant AdminSvc as Admin Service
        participant AnalyticsSvc as Analytics Service

        AnalyticsAdmin->>AdminPanelUI: Открывает "Аналитика" -> "Отчет по продажам"
        AdminPanelUI->>AdminSvc: GET /api/v1/admin/analytics/reports/sales?period=last_month&group_by=game
        alt AdminSvc проксирует запрос или использует данные из DWH AnalyticsSvc
            AdminSvc->>AnalyticsSvc: (REST/gRPC) GetSalesReport(period="last_month", group_by="game")
            AnalyticsSvc-->>AdminSvc: Данные отчета
        else AdminSvc сам имеет доступ к данным для некоторых отчетов
            AdminSvc->>AdminSvc: Запрос к своей БД или DWH (если часть данных агрегируется локально)
            AdminSvc-->>AdminPanelUI: Данные отчета
        end
        AdminPanelUI-->>AnalyticsAdmin: Отображение отчета с графиками и таблицами
    ```

## 15. Резервное Копирование и Восстановление (Backup and Recovery)

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

## 16. Связанные Рабочие Процессы (Related Workflows)
*   [Подача разработчиком новой игры на модерацию](../../../../project_workflows/game_submission_flow.md)
*   [Обработка жалобы пользователя на контент] Подробное описание этого рабочего процесса будет добавлено в [user_complaint_handling_flow.md](../../../../project_workflows/user_complaint_handling_flow.md) (документ в разработке).
*   [Процесс решения тикета поддержки] Подробное описание этого рабочего процесса будет добавлено в [support_ticket_resolution_flow.md](../../../../project_workflows/support_ticket_resolution_flow.md) (документ в разработке).

---
*Этот документ является отправной точкой и должен регулярно обновляться по мере развития сервиса.*
