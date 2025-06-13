# Спецификация Микросервиса: Admin Service

**Версия:** 1.0
**Дата последнего обновления:** {{YYYY-MM-DD}}

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
*   **REST Framework:** Echo (`github.com/labstack/echo/v4`)
*   **gRPC Framework:** `google.golang.org/grpc` (для потенциальных внутренних API)
*   **Базы данных:**
    *   PostgreSQL (версия 15+): для структурированных данных.
    *   MongoDB (версия 6.x+): для логов аудита, деталей модерации.
*   **Поисковый движок:** Elasticsearch (версия 8.x+): для поиска по тикетам, базе знаний, логам.
*   **Кэширование:** Redis (версия 7.0+)
*   **Брокер сообщений:** Kafka (Apache Kafka версии 3.x+)
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`)
*   **Логирование:** Zap (`go.uber.org/zap`)
*   **Валидация:** `go-playground/validator/v10`
*   **Трассировка и метрики:** OpenTelemetry (`go.opentelemetry.io/otel`), Prometheus client (`github.com/prometheus/client_golang`)
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `../../../../PACKAGE_STANDARDIZATION.md`.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`, `../../../../project_observability_standards.md`.

### 1.4. Термины и Определения (Glossary)
*   См. `../../../../project_glossary.md`.
*   **Модерация (Moderation):** Процесс проверки и управления пользовательским и разработческим контентом.
*   **Тикет поддержки (Support Ticket):** Формализованный запрос в службу поддержки.
*   **База знаний (Knowledge Base):** Коллекция статей и инструкций.
*   **Лог Аудита (Audit Log):** Запись действий администраторов.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Сервис Admin Service спроектирован как модульная система, придерживающаяся принципов чистой архитектуры (Clean Architecture).
*   Основные модули: Модерация контента, Управление пользователями платформы, Управление административными пользователями, Техническая поддержка, Мониторинг безопасности, Управление настройками платформы, Административная аналитика, Управление маркетинговыми кампаниями.
*   Диаграмма верхнеуровневой архитектуры сервиса:
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

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления)
*   Ответственность: Обработка входящих HTTP (REST) запросов, валидация DTO, вызов Application Layer.
*   Ключевые компоненты/модули:
    *   HTTP Handlers (Echo): Эндпоинты для административных функций.
    *   DTOs: Структуры для запросов и ответов API. Валидация с `go-playground/validator`.
    *   Middleware: Аутентификация, авторизация, логирование.

#### 2.2.2. Application Layer (Прикладной Слой / Слой Сценариев Использования)
*   Ответственность: Координация бизнес-логики, реализация use cases.
*   Ключевые компоненты/модули:
    *   Use Case Services: `PlatformUserManagementService`, `ContentModerationService`, etc.
    *   Интерфейсы для репозиториев и внешних сервисов.

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности, агрегаты, доменные события, бизнес-правила.
*   Ключевые компоненты/модули:
    *   Entities: `AdminUser`, `ModerationItem`, `SupportTicket`, `PlatformSetting`, etc.
    *   Value Objects: `ModerationReason`, `TicketStatus`, `AdminRole`.
    *   Domain Services: `ModerationRuleEngine`.
    *   Domain Events: `ContentModeratedEvent`, `UserStatusChangedByAdminEvent`.
    *   Интерфейсы репозиториев.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Реализация интерфейсов для взаимодействия с PostgreSQL, MongoDB, Elasticsearch, Redis, Kafka. Клиенты для других микросервисов.
*   Ключевые компоненты/модули:
    *   PostgreSQL Repositories (GORM).
    *   MongoDB Repositories (Official Go Driver).
    *   Elasticsearch Client.
    *   Redis Cache.
    *   Kafka Producers/Consumers.
    *   gRPC/REST клиенты для других сервисов.

## 3. API Endpoints
Общие принципы и форматы см. `../../../../project_api_standards.md`.

### 3.1. REST API
*   **Базовый URL (через API Gateway):** `/api/v1/admin`
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway, `X-Admin-UserID`, `X-Admin-Roles` передаются).
*   **Авторизация:** На основе ролей администратора (см. `../../../../project_roles_and_permissions.md`).

#### 3.1.1. Ресурс: Управление Пользователями Платформы
*   **`GET /platform-users`**
    *   Описание: Поиск и получение списка пользователей платформы.
    *   Query параметры: `search_query`, `status`, `role`, `page`, `per_page`, `sort_by`.
    *   Пример ответа (Успех 200 OK): `[NEEDS DEVELOPER INPUT: Example success response for GET /platform-users]`
    *   Пример ответа (Ошибка 4xx/5xx): `[NEEDS DEVELOPER INPUT: Example error response for GET /platform-users]`
    *   Требуемые права доступа: `admin:platform_users:read`, `support:platform_users:read_basic`.
*   **`GET /platform-users/{user_id}`**
    *   Описание: Получение детальной информации о пользователе платформы.
    *   Требуемые права доступа: `admin:platform_users:read_detailed`, `support:platform_users:read_detailed`.
*   **`PUT /platform-users/{user_id}/status`**
    *   Описание: Изменение статуса пользователя.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for PUT /platform-users/{user_id}/status]`
    *   Требуемые права доступа: `admin:platform_users:update_status`.
*   **`PUT /platform-users/{user_id}/roles`**
    *   Описание: Изменение ролей пользователя.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for PUT /platform-users/{user_id}/roles]`
    *   Требуемые права доступа: `admin:platform_users:update_roles`.
*   **`GET /platform-users/{user_id}/activity-log`**
    *   Описание: Получение лога активности пользователя.
    *   Требуемые права доступа: `admin:platform_users:read_activity_log`.
*   **`GET /platform-users/{user_id}/moderation-history`**
    *   Описание: Получение истории модерации, связанной с пользователем.
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
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for POST /moderation/items/{item_id}/decisions]`
    *   Требуемые права доступа: `moderator:items:decide`, `admin:items:decide`.
*   **`POST /moderation/rules`**
    *   Описание: Создание нового правила автоматической модерации.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for POST /moderation/rules]`
    *   Требуемые права доступа: `admin:moderation_rules:create`.
*   **`GET /moderation/items/history`**
    *   Описание: Получение истории решений по модерации.
    *   Query параметры: `item_reference_id`, `item_type`, `moderator_id`, `decision_type`, `start_date`, `end_date`, `page`, `per_page`.
    *   Требуемые права доступа: `admin:moderation_history:read`, `moderator_lead:moderation_history:read`.

#### 3.1.3. Ресурс: Техническая Поддержка
*   **`GET /support/tickets`**
    *   Описание: Получение списка тикетов поддержки.
    *   Требуемые права доступа: `support:tickets:read`, `admin:tickets:read`.
*   **`POST /support/tickets`**
    *   Описание: Создание тикета от имени пользователя.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for POST /support/tickets]`
    *   Требуемые права доступа: `support:tickets:create`, `admin:tickets:create`.
*   **`GET /support/tickets/{ticket_id}`**
    *   Описание: Получение информации о конкретном тикете.
    *   Требуемые права доступа: `support:tickets:read_detailed`, `admin:tickets:read_detailed`.
*   **`POST /support/tickets/{ticket_id}/responses`**
    *   Описание: Добавление ответа в тикет.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for POST /support/tickets/{ticket_id}/responses]`
    *   Требуемые права доступа: `support:tickets:respond`, `admin:tickets:respond`.
*   **`PUT /support/tickets/{ticket_id}/status`**
    *   Описание: Изменение статуса тикета.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for PUT /support/tickets/{ticket_id}/status]`
    *   Требуемые права доступа: `support:tickets:update_status`, `admin:tickets:update_status`.
*   **`PUT /support/tickets/{ticket_id}/assignee`**
    *   Описание: Назначение ответственного агента.
    *   Тело запроса: `{"data": {"type": "ticketAssignee", "attributes": {"assignee_admin_id": "support-agent-uuid-xyz"}}}`
    *   Требуемые права доступа: `support_lead:tickets:assign`, `admin:tickets:assign`.
*   **`GET /support/knowledge-base`**
    *   Описание: Поиск и получение статей из базы знаний.
    *   Требуемые права доступа: `support:kb:read`, `admin:kb:read`.
*   **`POST /support/knowledge-base`**
    *   Описание: Создание новой статьи в базе знаний.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for POST /support/knowledge-base]`
    *   Требуемые права доступа: `admin:kb:create`, `support_lead:kb:create`.
*   **`GET /support/reports/performance`**
    *   Описание: Получение отчета по эффективности службы поддержки.
    *   Требуемые права доступа: `admin:support_reports:read`, `support_lead:support_reports:read`.
*   **`GET /support/categories`**
    *   Описание: Получение списка категорий тикетов поддержки.
    *   Требуемые права доступа: `support:ticket_categories:read`, `admin:ticket_categories:read`.
*   **`POST /support/ticket-templates`**
    *   Описание: Создание шаблона ответа для тикетов.
    *   Тело запроса: `{"data": {"type": "ticketTemplate", "attributes": {"name": "Шаблон для сброса пароля", "text_ru": "...", "category_id": "account_issues", "language_code": "ru"}}}`
    *   Требуемые права доступа: `admin:ticket_templates:create`, `support_lead:ticket_templates:create`.

#### 3.1.4. Ресурс: Настройки Платформы
*   **`GET /platform-settings`**
    *   Описание: Получение текущих настроек платформы.
    *   Требуемые права доступа: `admin:platform_settings:read`.
*   **`PUT /platform-settings`**
    *   Описание: Обновление настроек платформы.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for PUT /platform-settings, should match GET response structure]`
    *   Требуемые права доступа: `admin:platform_settings:update`.
*   **`POST /platform-settings/maintenance-schedule`**
    *   Описание: Планирование окна обслуживания.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for POST /platform-settings/maintenance-schedule]`
    *   Требуемые права доступа: `admin:platform_settings:schedule_maintenance`.

#### 3.1.5. Ресурс: Управление Администраторами
*   **`GET /admin-users`**
    *   Описание: Список административных пользователей.
    *   Требуемые права доступа: `admin:admin_users:read`.
*   **`POST /admin-users`**
    *   Описание: Создание нового административного пользователя.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for POST /admin-users]`
    *   Требуемые права доступа: `admin:admin_users:create`.
*   **`GET /admin-users/{admin_user_id}`**
    *   Описание: Получение информации о конкретном административном пользователе.
    *   Требуемые права доступа: `admin:admin_users:read_detailed`.
*   **`PUT /admin-users/{admin_user_id}/permissions`**
    *   Описание: Изменение ролей и разрешений администратора.
    *   Тело запроса: `[NEEDS DEVELOPER INPUT: Example request body for PUT /admin-users/{admin_user_id}/permissions]`
    *   Требуемые права доступа: `admin:admin_users:update_permissions`.
*   **`GET /admin-users/{admin_user_id}/audit-log`**
    *   Описание: Получение лога действий конкретного администратора.
    *   Требуемые права доступа: `admin:admin_users:read_audit_log`.

### 3.2. gRPC API
*   В настоящее время gRPC API для внешнего использования Admin Service не планируется.
*   Ссылка на `.proto` файл: `[NEEDS DEVELOPER INPUT: Path to .proto if internal gRPC API exists and needs to be documented, otherwise state "Not Applicable"]`

### 3.3. WebSocket API (если применимо)
*   Не планируется для Admin Service на данном этапе.

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`AdminUser`**: Учетная запись администратора (PostgreSQL).
*   **`ModerationItem`**: Объект модерации (PostgreSQL/MongoDB).
*   **`ModerationRule`**: Правило модерации (PostgreSQL).
*   **`SupportTicket`**: Тикет поддержки (PostgreSQL).
*   **`SupportTicketResponse`**: Ответ в тикете (PostgreSQL).
*   **`SupportTicketCategory`**: Категория тикета (PostgreSQL).
*   **`KnowledgeBaseArticle`**: Статья базы знаний (PostgreSQL, Elasticsearch).
*   **`KnowledgeBaseCategory`**: Категория статьи БЗ (PostgreSQL).
*   **`PlatformSetting`**: Настройка платформы (PostgreSQL, Redis).
*   **`AuditLogAdmin`**: Лог действий администраторов (MongoDB/PostgreSQL).
*   **`ModerationItemDetail`**: Детали элемента модерации (MongoDB).
*   [NEEDS DEVELOPER INPUT: Review and add any other key entities for admin-service]

### 4.2. Схема Базы Данных (если применимо)
*   Диаграмма ERD: `[NEEDS DEVELOPER INPUT: Mermaid ERD diagram for PostgreSQL tables specific to admin-service, or confirm if existing diagrams cover it]`
*   Описание DDL для PostgreSQL: `[NEEDS DEVELOPER INPUT: DDL for admin-service specific PostgreSQL tables if not covered elsewhere]`
*   **MongoDB Коллекции:**
    *   `audit_logs_admin`: поля: `admin_user_id`, `timestamp`, `action_type`, `target_entity_type`, `target_entity_id`, `details`, `ip_address`. Индексы: `admin_user_id`, `timestamp`, `action_type`, `target_entity_id`.
    *   `moderation_item_details`: поля: `_id` (соотв. `ModerationItem.id`), `full_content_snapshot`, `external_analysis_reports`. Индексы: `_id`.
*   **Elasticsearch Индексы:**
    *   `support_tickets_idx`: поля: `id`, `subject`, `description`, `responses.text`, `user_id`, `status`, `priority`, `tags`, `created_at`.
    *   `knowledge_base_articles_idx`: поля: `id`, `title`, `content_markdown`, `tags`, `category_name`, `is_published`.
    *   `admin_audit_log_idx`: (опционально) поля: `id`, `admin_username`, `action_type`, `target_entity_type`, `target_entity_id`, `timestamp`, `details`.
    *   `moderation_items_idx`: поля: `id`, `item_type`, `status`, `priority`, `content_snapshot`, `user_id`.

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   Используемая система сообщений: Kafka.
*   Формат событий: CloudEvents JSON (согласно `../../../../project_api_standards.md`).
*   Основные топики Kafka: `com.platform.admin.events.v1`.

*   **`com.platform.admin.content.moderated.v1`**
    *   Описание: Контент прошел модерацию.
    *   Топик: `com.platform.admin.events.v1`
    *   Структура Payload: (см. существующий документ)
*   **`com.platform.admin.user.status.updated.v1`**
    *   Описание: Статус пользователя платформы изменен администратором.
    *   Топик: `com.platform.admin.events.v1`
    *   Структура Payload: (см. существующий документ)
*   **`com.platform.admin.platform.setting.updated.v1`**
    *   Описание: Изменена глобальная настройка платформы.
    *   Топик: `com.platform.admin.events.v1`
    *   Структура Payload: (см. существующий документ)
*   **`com.platform.admin.support.ticket.status.updated.v1`**
    *   Описание: Статус тикета поддержки изменен.
    *   Топик: `com.platform.admin.events.v1`
    *   Структура Payload: (см. существующий документ)
*   **`com.platform.admin.marketing.campaign.status.updated.v1`**
    *   Описание: Статус маркетинговой кампании изменен.
    *   Топик: `com.platform.admin.events.v1`
    *   Структура Payload: `{"campaignId": "uuid", "newStatus": "active", "adminId": "uuid", "timestamp": "ISO8601"}`

### 5.2. Потребляемые События (Consumed Events)
*   **`com.platform.user.complaint.submitted.v1`** (от Social Service, Catalog Service)
    *   Описание: Поступила жалоба от пользователя.
    *   Топик: `[NEEDS DEVELOPER INPUT: Specific topic name for user complaints]`
    *   Логика обработки: Создать `ModerationItem`.
*   **`com.platform.developer.game.submitted.v1`** (от Developer Service)
    *   Описание: Разработчик отправил игру/обновление на модерацию.
    *   Топик: `[NEEDS DEVELOPER INPUT: Specific topic name for game submissions]`
    *   Логика обработки: Создать `ModerationItem`.
*   **`com.platform.payment.transaction.suspicious.v1`** (от Payment Service, Analytics Service)
    *   Описание: Обнаружена подозрительная финансовая активность.
    *   Топик: `[NEEDS DEVELOPER INPUT: Specific topic name for suspicious transactions]`
    *   Логика обработки: Создать инцидент безопасности.
*   **`com.platform.system.health.degraded.v1`** (от системы мониторинга)
    *   Описание: Обнаружена деградация производительности сервиса.
    *   Топик: `[NEEDS DEVELOPER INPUT: Specific topic name for system health events]`
    *   Логика обработки: Создать алерт/тикет.

## 6. Интеграции (Integrations)
См. `../../../../project_integrations.md` для общей карты и деталей.
Admin Service активно взаимодействует с большинством других микросервисов.
[NEEDS DEVELOPER INPUT: Review and confirm service boundaries and reliability strategies for each integration for admin-service]

### 6.1. Внутренние Микросервисы
*   **Account Service**: Тип: gRPC/REST. Назначение: Получение данных пользователей. Контракт: `GetAccountInfo`, `GetUserProfile`. Надежность: `[NEEDS DEVELOPER INPUT]`
*   **Auth Service**: Тип: gRPC/REST. Назначение: Управление сессиями администраторов, проверка токенов. Контракт: `ValidateToken`, `GetAdminUserRoles`. Надежность: `[NEEDS DEVELOPER INPUT]`
*   **Catalog Service**: Тип: gRPC/REST, Kafka Events. Назначение: Модерация контента каталога, управление видимостью. Контракт: `GetGameDetails`, `UpdateGameStatus`, `com.platform.admin.content.moderated.v1`. Надежность: `[NEEDS DEVELOPER INPUT]`
*   **Developer Service**: Тип: Kafka Events. Назначение: Получение контента на модерацию. Контракт: `com.platform.developer.game.submitted.v1`. Надежность: `[NEEDS DEVELOPER INPUT]`
*   **Notification Service**: Тип: gRPC/Kafka. Назначение: Отправка уведомлений по результатам модерации, тикетам. Контракт: `SendEmailNotification`, `com.platform.admin.support.ticket.status.updated.v1`. Надежность: `[NEEDS DEVELOPER INPUT]`
*   **Payment Service**: Тип: gRPC/REST, Kafka Events. Назначение: Просмотр транзакций, управление возвратами (если применимо), получение событий о подозрительных транзакциях. Контракт: `GetTransactionDetails`, `com.platform.payment.transaction.suspicious.v1`. Надежность: `[NEEDS DEVELOPER INPUT]`
*   **Analytics Service**: Тип: REST/Kafka Events. Назначение: Доступ к административной аналитике, получение событий. Контракт: `GetPlatformReport`, `com.platform.payment.transaction.suspicious.v1`. Надежность: `[NEEDS DEVELOPER INPUT]`
*   **Social Service**: Тип: Kafka Events. Назначение: Получение жалоб на пользовательский контент. Контракт: `com.platform.user.complaint.submitted.v1`. Надежность: `[NEEDS DEVELOPER INPUT]`
*   ... [NEEDS DEVELOPER INPUT: Add other relevant microservice integrations]

### 6.2. Внешние Системы
*   На данном этапе не предполагается прямых интеграций с внешними системами, кроме тех, что обеспечиваются другими сервисами.
*   [NEEDS DEVELOPER INPUT: Confirm if any direct external system integrations are planned for admin-service]

## 7. Конфигурация (Configuration)
Общие стандарты: `../../../../project_api_standards.md` (раздел 7) и `../../../../DOCUMENTATION_GUIDELINES.md` (раздел 6).

### 7.1. Переменные Окружения
*   `ADMIN_HTTP_PORT`: Порт REST API (e.g., `8080`)
*   `ADMIN_GRPC_PORT`: Порт gRPC API (если используется) (e.g., `9090`)
*   `POSTGRES_DSN`: DSN для PostgreSQL.
*   `MONGO_DSN`: DSN для MongoDB.
*   `ELASTICSEARCH_URLS`: URL(ы) Elasticsearch.
*   `REDIS_ADDR`: Адрес Redis.
*   `KAFKA_BROKERS`: Брокеры Kafka.
*   `JWT_PUBLIC_KEY_PATH`: Путь к публичному ключу для JWT.
*   `LOG_LEVEL`: Уровень логирования.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`: Jaeger endpoint.
*   [NEEDS DEVELOPER INPUT: Add any other critical environment variables for admin-service]

### 7.2. Файлы Конфигурации (`configs/config.yaml`)
*   Расположение: `backend/admin-service/configs/config.yaml`
*   Структура:
    ```yaml
    http_server:
      port: ${ADMIN_HTTP_PORT:"8080"}
      # ...
    postgres:
      dsn: ${POSTGRES_DSN}
      # ...
    mongodb:
      dsn: ${MONGO_DSN}
      # ...
    elasticsearch:
      urls: ${ELASTICSEARCH_URLS}
      # ...
    redis:
      address: ${REDIS_ADDR}
      # ...
    kafka:
      brokers: ${KAFKA_BROKERS}
      # ...
    logging:
      level: ${LOG_LEVEL:"info"}
    security:
      jwt_public_key_path: ${JWT_PUBLIC_KEY_PATH}
    # ... [NEEDS DEVELOPER INPUT: Add other specific config sections for admin-service]
    ```

## 8. Обработка Ошибок (Error Handling)
*   Стандартные HTTP коды и форматы ошибок согласно `../../../../project_api_standards.md`.
*   Логирование с `trace_id`.

### 8.1. Общие Принципы
*   Стандартный формат JSON ответа об ошибке.
*   Использование стандартных HTTP кодов.

### 8.2. Распространенные Коды Ошибок
*   **`ADMIN_USER_NOT_FOUND`**
*   **`PLATFORM_USER_NOT_FOUND`**
*   **`MODERATION_ITEM_NOT_FOUND`**
*   **`SUPPORT_TICKET_NOT_FOUND`**
*   **`PLATFORM_SETTING_VALIDATION_ERROR`**
*   **`ACTION_NOT_PERMITTED_FOR_ROLE`**
*   **`TARGET_SERVICE_UNAVAILABLE`**
*   [NEEDS DEVELOPER INPUT: Review and add other specific error codes for admin-service]

## 9. Безопасность (Security)
См. `../../../../project_security_standards.md`.

### 9.1. Аутентификация
*   JWT для администраторов через Auth Service. API Gateway валидирует токен.

### 9.2. Авторизация
*   RBAC для администраторов. Роли в `../../../../project_roles_and_permissions.md`.

### 9.3. Защита Данных
*   **ФЗ-152 "О персональных данных"**: Обработка ПДн пользователей и администраторов. Логирование доступа.
*   Шифрование при передаче (TLS).
*   Управление секретами.

### 9.4. Управление Секретами
*   Kubernetes Secrets или HashiCorp Vault.

### 9.5. Аудит Действий
*   Все значимые действия администраторов логируются в `AuditLogAdmin`.

## 10. Развертывание (Deployment)
См. `../../../../project_deployment_standards.md`.

### 10.1. Инфраструктурные Файлы
*   Dockerfile: `backend/admin-service/Dockerfile`
*   Helm-чарты: `deploy/charts/admin-service/`

### 10.2. Зависимости при Развертывании
*   PostgreSQL, MongoDB, Elasticsearch, Redis, Kafka.
*   Доступ к API других микросервисов. API Gateway.

### 10.3. CI/CD
*   Стандартный пайплайн CI/CD. Автоматическое применение миграций БД.

## 11. Мониторинг и Логирование (Logging and Monitoring)
См. `../../../../project_observability_standards.md`.

### 11.1. Логирование
*   Формат: JSON (Zap). Ключевые поля: `timestamp`, `level`, `service`, `trace_id`, `admin_user_id`.
*   Интеграция: Fluent Bit в Elasticsearch/Loki.

### 11.2. Мониторинг
*   Метрики (Prometheus): HTTP запросы, длительность, размеры очередей модерации, статусы тикетов, производительность БД, внешних вызовов, Kafka.
*   Дашборды (Grafana): Обзор состояния, API, очередей, тикетов, активности администраторов.
*   Алерты (AlertManager): Ошибки API, задержки, проблемы с зависимостями, рост очередей.

### 11.3. Трассировка
*   Интеграция: OpenTelemetry SDK. Экспорт: Jaeger.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**: API интерактивных операций P95 < 800мс; API операций изменения P95 < 500мс. Загрузка главной страницы админ-панели P95 < 2с.
*   **Надежность**: Доступность 99.9%. RTO < 1ч (критичные функции). RPO < 15мин (конфигурации, тикеты), < 1ч (логи аудита).
*   **Безопасность**: Соответствие `../../../../project_security_standards.md` и ФЗ-152.
*   **Масштабируемость**: До 1000 одновременных сессий администраторов.
*   **Сопровождаемость**: Покрытие тестами > 80%.
*   [NEEDS DEVELOPER INPUT: Confirm or update NFRs for admin-service]

## 13. Резервное Копирование и Восстановление (Backup and Recovery)
Общие принципы см. в `../../../../project_database_structure.md`.

### 13.1. PostgreSQL
*   Процедура: Ежедневный `pg_dump`, непрерывная архивация WAL (PITR).
*   Хранение: S3. RTO: < 2ч. RPO: < 5мин.

### 13.2. MongoDB (Логи аудита, детали модерации)
*   Процедура: Ежедневный `mongodump`.
*   Хранение: S3. RTO: < 4ч. RPO: 24ч.

### 13.3. Elasticsearch (Индексы для поиска)
*   Процедура: Ежедневные Elasticsearch Snapshots.
*   Хранение: S3. RTO: < 3ч (из снапшота). RPO: 24ч.

### 13.4. Redis (Кэш)
*   Не критично для бэкапа, данные восстановимы из основных хранилищ.

## 14. Приложения (Appendices) (Опционально)
*   OpenAPI спецификация: `[NEEDS DEVELOPER INPUT: Link to OpenAPI spec or state if generated from code for admin-service]`
*   [NEEDS DEVELOPER INPUT: Add any other appendices if necessary for admin-service]

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [Подача разработчиком новой игры на модерацию](../../../../project_workflows/game_submission_flow.md)
*   [NEEDS DEVELOPER INPUT: Add link to user_complaint_handling_flow.md when created]
*   [NEEDS DEVELOPER INPUT: Add link to support_ticket_resolution_flow.md when created]
*   [NEEDS DEVELOPER INPUT: Add links to other relevant high-level workflow documents for admin-service]

---
*Этот документ является отправной точкой и должен регулярно обновляться по мере развития сервиса.*
