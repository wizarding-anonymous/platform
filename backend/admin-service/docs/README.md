# Спецификация Микросервиса: Admin Service

**Версия:** 1.0
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   Admin Service является центральным компонентом платформы "Российский Аналог Steam", предоставляющим инструменты и интерфейсы для администраторов платформы и сотрудников службы поддержки.
*   Его основная роль - управление и надзор за различными аспектами платформы, включая модерацию контента, управление пользователями, операции технической поддержки, мониторинг безопасности, конфигурацию общесистемных настроек, административную аналитику и управление маркетинговыми кампаниями.
*   Основные бизнес-задачи: обеспечение операционного контроля, поддержка пользователей, модерация контента, управление системными параметрами.

### 1.2. Ключевые Функциональности
*   **Управление модерацией контента:** Управление рабочими процессами модерации пользовательского и разработческого контента (игры, отзывы, комментарии, профили). Поддержка автоматизированных и ручных процессов.
*   **Управление пользователями платформы:** Поиск, просмотр, управление статусами пользователей (активация, блокировка), управление ролями.
*   **Управление административными пользователями:** Управление учетными записями и разрешениями сотрудников.
*   **Техническая поддержка:** Система управления тикетами поддержки (категоризация, назначение, ответы), база знаний.
*   **Мониторинг безопасности:** Инструменты для выявления подозрительной активности, анализа платежей, управления IP-блокировками, аудита действий администраторов.
*   **Управление настройками платформы:** Конфигурация глобальных и региональных параметров, операционных лимитов.
*   **Административная аналитика:** Дашборды и отчеты по метрикам платформы (продажи, активность пользователей, эффективность модерации).
*   **Управление маркетинговыми кампаниями:** Создание и управление промо-акциями, скидками, специальными предложениями.

### 1.3. Основные Технологии
*   **Язык программирования:** Go (предпочтительно, в соответствии с `project_technology_stack.md`)
*   **REST Framework:** Echo (согласно `project_technology_stack.md`)
*   **gRPC Framework:** google.golang.org/grpc (для возможного внутреннего использования или API для других админ-тулов)
*   **База данных:** PostgreSQL (для структурированных данных: админ-пользователи, тикеты, настройки), MongoDB (для логов аудита, сложных данных модерации).
*   **Поисковый движок:** Elasticsearch (для поиска по логам, тикетам).
*   **Кэширование:** Redis.
*   **Брокер сообщений:** Kafka (для асинхронных задач, например, генерация отчетов, массовые операции).
*   (Ссылки на `project_technology_stack.md`, `project_glossary.md`)

### 1.4. Термины и Определения (Glossary)
*   См. `project_glossary.md`.
*   **Модерация (Moderation):** Процесс проверки и управления контентом.
*   **Тикет поддержки (Support Ticket):** Запрос пользователя в службу поддержки.
*   **База знаний (Knowledge Base):** Сборник статей для решения проблем.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Сервис может быть реализован как модульный монолит или набор тесно связанных микро-модулей, придерживаясь принципов чистой архитектуры.
*   Основные модули: Модерация, Управление пользователями платформы, Управление администраторами, Техническая поддержка, Мониторинг безопасности, Управление настройками платформы, Административная аналитика, Управление маркетинговыми кампаниями.
*   Ниже представлена диаграмма верхнеуровневой архитектуры Admin Service:

```mermaid
graph TD
    A[API Gateway] --> AdminService[Admin Service]

    subgraph Admin Service
        direction LR
        Mod[Модерация]
        UserMgmt[Управление пользователями платформы]
        AdminUserMgmt[Управление администраторами]
        Support[Техническая поддержка]
        Security[Мониторинг безопасности]
        Settings[Управление настройками платформы]
        Analytics[Административная аналитика]
        Marketing[Управление маркетинговыми кампаниями]
    end

    AdminService --> DB[(PostgreSQL/MongoDB)]
    AdminService --> Search[(Elasticsearch)]
    AdminService --> Cache[(Redis)]
    AdminService --> MQ[(Kafka)]

    Mod --- UserMgmt
    Mod --- Support
    UserMgmt --- Security
    Support --- UserMgmt
    Settings --- Mod
    Settings --- UserMgmt
    Analytics --- Mod
    Analytics --- UserMgmt
    Analytics --- Support
    Marketing --- Settings

    AdminService -->|gRPC/REST| AccountService[Account Service]
    AdminService -->|gRPC/REST| CatalogService[Catalog Service]
    AdminService -->|gRPC/REST| PaymentService[Payment Service]
    AdminService -->|gRPC/REST| DeveloperService[Developer Service]
    AdminService -->|gRPC/REST| SocialService[Social Service]
    AdminService -->|gRPC/REST| NotificationService[Notification Service]

    classDef service fill:#D6EAF8,stroke:#3498DB,stroke-width:2px;
    classDef module fill:#E8F8F5,stroke:#1ABC9C,stroke-width:2px;
    classDef datastore fill:#FDEDEC,stroke:#E74C3C,stroke-width:2px;
    classDef external fill:#FEF9E7,stroke:#F1C40F,stroke-width:2px;

    class AdminService,AccountService,CatalogService,PaymentService,DeveloperService,SocialService,NotificationService service;
    class Mod,UserMgmt,AdminUserMgmt,Support,Security,Settings,Analytics,Marketing module;
    class DB,Search,Cache,MQ datastore;
    class A external;
```

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления)
*   Ответственность: Обработка входящих HTTP (REST) запросов от административной панели. Валидация DTO.
*   Ключевые компоненты/модули:
    *   HTTP Handlers (Echo): Эндпоинты для всех административных функций.
    *   DTOs: Структуры для передачи данных.
    *   Валидаторы.

#### 2.2.2. Application Layer (Прикладной Слой)
*   Ответственность: Координация бизнес-логики для каждого административного модуля.
*   Ключевые компоненты/модули:
    *   Use Case Services: `ModerationService`, `UserManagementService`, `SupportTicketService`, `PlatformSettingsService`, etc.
    *   Интерфейсы для репозиториев и клиентов других сервисов.

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности и правила.
*   Ключевые компоненты/модули:
    *   Entities: `AdminUser`, `ModerationQueue`, `ModerationItem`, `ModerationDecision`, `SupportTicket`, `KnowledgeBaseArticle`, `SystemSetting`, `PlatformUserView` (агрегированные данные о пользователе).
    *   Domain Events: `ContentModeratedEvent`, `UserStatusChangedByAdminEvent`, `TicketResolvedEvent`.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Взаимодействие с PostgreSQL, MongoDB, Elasticsearch, Redis, Kafka. Клиенты для других микросервисов.
*   Ключевые компоненты/модули:
    *   PostgreSQL/MongoDB Repositories.
    *   Elasticsearch Client.
    *   Redis Cache.
    *   Kafka Producers/Consumers.
    *   gRPC/REST клиенты для Account Service, Catalog Service, Payment Service и др.

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL (через API Gateway):** `/api/v1/admin`
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway, UserID администратора и его роли передаются в заголовках).
*   **Авторизация:** На основе ролей администратора (см. `project_roles_and_permissions.md`).
*   **Формат ответа об ошибке (стандартный):**
    ```json
    {
      "errors": [
        {
          "status": "4XX/5XX",
          "code": "ERROR_CODE",
          "title": "Краткое описание ошибки",
          "detail": "Полное описание ошибки с контекстом",
          "source": { "pointer": "/data/attributes/field_name" } // Опционально, для ошибок валидации
        }
      ]
    }
    ```
*   (Общие принципы см. `project_api_standards.md`)

#### 3.1.1. Ресурс: Управление Пользователями Платформы
*   **`GET /platform-users`**
    *   Описание: Поиск и получение списка пользователей платформы.
    *   Query параметры:
        *   `search_query` (string, опционально): Поисковый запрос по имени пользователя, email.
        *   `status` (string, опционально, enum: `active`, `blocked`, `unverified`): Фильтр по статусу.
        *   `role` (string, опционально): Фильтр по роли пользователя.
        *   `page` (integer, опционально, default: 1): Номер страницы.
        *   `per_page` (integer, опционально, default: 20): Количество элементов на странице.
        *   `sort_by` (string, опционально, enum: `created_at_asc`, `created_at_desc`, `username_asc`): Поле и направление сортировки.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "platformUser",
              "id": "uuid-user-1",
              "attributes": {
                "username": "player123",
                "email": "player123@example.com",
                "status": "active",
                "roles": ["user"],
                "created_at": "2023-01-15T10:00:00Z",
                "last_login_at": "2024-03-10T12:00:00Z"
              },
              "links": {
                "self": "/api/v1/admin/platform-users/uuid-user-1"
              }
            }
          ],
          "meta": {
            "total_items": 100,
            "total_pages": 5,
            "current_page": 1,
            "per_page": 20
          },
          "links": {
            "self": "/api/v1/admin/platform-users?page=1&per_page=20",
            "next": "/api/v1/admin/platform-users?page=2&per_page=20"
          }
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
              "detail": "Параметр 'status' имеет недопустимое значение 'unknown'.",
              "source": { "parameter": "status" }
            }
          ]
        }
        ```
    *   Требуемые права доступа: `admin`, `support` (с ограничениями на фильтры и возвращаемые поля).
*   **`GET /platform-users/{user_id}`**
    *   Описание: Получение детальной информации о пользователе платформы (включая связанные данные из других сервисов, например, история заказов, количество сообщений).
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "platformUser",
            "id": "uuid-user-1",
            "attributes": {
              "username": "player123",
              "email": "player123@example.com",
              "status": "active",
              "roles": ["user", "beta_tester"],
              "profile": {
                "first_name": "Иван",
                "last_name": "Иванов",
                "bio": "Люблю играть в игры.",
                "avatar_url": "https://example.com/avatars/uuid-user-1.png"
              },
              "statistics": {
                "games_owned": 15,
                "reviews_written": 3,
                "support_tickets_opened": 1
              },
              "created_at": "2023-01-15T10:00:00Z",
              "updated_at": "2024-03-01T15:30:00Z",
              "last_login_at": "2024-03-10T12:00:00Z"
            },
            "relationships": {
              "moderation_history": {
                "links": {
                  "related": "/api/v1/admin/moderation/items?user_id=uuid-user-1"
                }
              }
            },
            "links": {
              "self": "/api/v1/admin/platform-users/uuid-user-1"
            }
          }
        }
        ```
    *   Пример ответа (Ошибка 404 Not Found):
        ```json
        {
          "errors": [
            {
              "status": "404",
              "code": "USER_NOT_FOUND",
              "title": "Пользователь не найден",
              "detail": "Пользователь с ID 'uuid-user-unknown' не найден."
            }
          ]
        }
        ```
    *   Требуемые права доступа: `admin`, `support`.
*   **`PUT /platform-users/{user_id}/status`**
    *   Описание: Изменение статуса пользователя (например, `block`, `unblock`, `verify_manually`, `mute`).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "userStatus",
            "attributes": {
              "status": "blocked",
              "reason": "Нарушение пункта 3.4 правил сообщества.",
              "comment": "Пользователь уведомлен.",
              "expires_at": "2024-12-31T23:59:59Z" // Опционально, для временных блокировок/мутов
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "platformUser",
            "id": "uuid-user-1",
            "attributes": {
              "status": "blocked",
              "reason_for_status_change": "Нарушение пункта 3.4 правил сообщества."
            }
          }
        }
        ```
    *   Требуемые права доступа: `admin`.
*   **`PUT /platform-users/{user_id}/roles`**
    *   Описание: Изменение ролей пользователя (полное замещение списка ролей).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "userRoles",
            "attributes": {
              "roles": ["user", "moderator_trainee", "trusted_reviewer"]
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "platformUser",
            "id": "uuid-user-1",
            "attributes": {
              "roles": ["user", "moderator_trainee", "trusted_reviewer"]
            }
          }
        }
        ```
    *   Требуемые права доступа: `admin`.
*   **`GET /platform-users/{user_id}/activity-log`**
    *   Описание: Получение лога активности пользователя на платформе (вход, ключевые действия, изменения профиля).
    *   Query параметры: `event_type`, `start_date`, `end_date`, `page`, `per_page`.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "userActivityLogEntry",
              "id": "log-entry-uuid-1",
              "attributes": {
                "event_type": "user_login",
                "timestamp": "2024-03-15T10:00:00Z",
                "ip_address": "192.168.1.100",
                "user_agent": "Chrome/100.0",
                "details": { "login_method": "password" }
              }
            },
            {
              "type": "userActivityLogEntry",
              "id": "log-entry-uuid-2",
              "attributes": {
                "event_type": "profile_update",
                "timestamp": "2024-03-15T10:05:00Z",
                "details": { "changed_fields": ["bio", "avatar_url"] }
              }
            }
          ],
          "meta": { "total_items": 50, "current_page": 1, "per_page": 10 },
          "links": { "next": "/api/v1/admin/platform-users/uuid-user-1/activity-log?page=2" }
        }
        ```
    *   Требуемые права доступа: `admin`, `security_officer`.

#### 3.1.2. Ресурс: Модерация Контента
*   **`GET /moderation/queues`**
    *   Описание: Получение списка очередей модерации и их текущего состояния (количество элементов).
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "moderationQueue",
              "id": "game_reviews",
              "attributes": {
                "display_name": "Отзывы на игры",
                "item_count": 125,
                "oldest_item_at": "2024-03-14T09:00:00Z"
              },
              "links": {
                "items": "/api/v1/admin/moderation/queues/game_reviews/items"
              }
            },
            {
              "type": "moderationQueue",
              "id": "user_avatars",
              "attributes": {
                "display_name": "Аватары пользователей",
                "item_count": 30,
                "oldest_item_at": "2024-03-15T11:00:00Z"
              },
              "links": {
                "items": "/api/v1/admin/moderation/queues/user_avatars/items"
              }
            }
          ]
        }
        ```
    *   Требуемые права доступа: `moderator`, `admin`.
*   **`GET /moderation/queues/{queue_id}/items`**
    *   Описание: Получение элементов из конкретной очереди модерации.
    *   Query параметры:
        *   `page` (integer, опционально, default: 1): Номер страницы.
        *   `per_page` (integer, опционально, default: 10): Количество элементов на странице.
        *   `status` (string, опционально, enum: `pending`, `in_review`): Фильтр по статусу элемента.
        *   `assigned_to` (string, опционально, `me` или `unassigned`): Фильтр по назначенному модератору.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "moderationItem",
              "id": "item-uuid-1",
              "attributes": {
                "item_type": "game_review",
                "submitted_at": "2024-03-14T09:00:00Z",
                "status": "pending",
                "content_preview": "Эта игра просто ужасна...",
                "priority": "normal"
              },
              "relationships": {
                "submitted_by": { "data": { "type": "platformUser", "id": "user-uuid-submitter" } }
              },
              "links": {
                "self": "/api/v1/admin/moderation/items/item-uuid-1"
              }
            }
          ],
          "meta": { "total_items": 125, "current_page": 1, "per_page": 10 },
          "links": { "next": "/api/v1/admin/moderation/queues/game_reviews/items?page=2" }
        }
        ```
    *   Требуемые права доступа: `moderator`, `admin`.
*   **`GET /moderation/items/{item_id}`**
    *   Описание: Получение конкретного элемента на модерацию, включая полное содержимое и историю.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "moderationItem",
            "id": "item-uuid-1",
            "attributes": {
              "item_type": "game_review",
              "item_reference_id": "review-uuid-456", // ID оригинального отзыва
              "submitted_at": "2024-03-14T09:00:00Z",
              "status": "in_review",
              "content": {
                "title": "Ужасная игра",
                "text": "Эта игра просто ужасна, не рекомендую никому.",
                "rating": 1
              },
              "history": [
                { "timestamp": "2024-03-14T09:00:00Z", "action": "submitted", "actor_id": "user-uuid-submitter" },
                { "timestamp": "2024-03-15T14:00:00Z", "action": "assigned", "actor_id": "admin-uuid-moderator", "details": { "assigned_to": "admin-uuid-moderator" } }
              ]
            },
            "relationships": {
              "submitted_by": { "data": { "type": "platformUser", "id": "user-uuid-submitter" } },
              "assigned_to": { "data": { "type": "adminUser", "id": "admin-uuid-moderator" } }
            },
            "links": {
              "self": "/api/v1/admin/moderation/items/item-uuid-1"
            }
          }
        }
        ```
    *   Требуемые права доступа: `moderator`, `admin`.
*   **`POST /moderation/items/{item_id}/decisions`**
    *   Описание: Принятие решения по элементу модерации (`approve`, `reject`, `escalate`, `request_changes`).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "moderationDecision",
            "attributes": {
              "decision": "rejected",
              "reason_code": "hate_speech", // Код из предопределенного списка причин
              "comment": "Комментарий содержит разжигание ненависти.",
              "send_notification_to_user": true,
              "action_details": { // Специфичные для решения действия
                "block_user_ لمدة_days": 7 // Например, если решение - отклонить и заблокировать
              }
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "moderationItem",
            "id": "item-uuid-1",
            "attributes": {
              "status": "rejected",
              "decision": {
                "type": "rejected",
                "reason_code": "hate_speech",
                "comment": "Комментарий содержит разжигание ненависти.",
                "decided_at": "2024-03-15T15:00:00Z",
                "decided_by": "admin-uuid-moderator"
              }
            }
          }
        }
        ```
    *   Требуемые права доступа: `moderator`, `admin`.
*   **`POST /moderation/rules`**
    *   Описание: Создание нового правила автоматической модерации.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "moderationRule",
            "attributes": {
              "name": "Блокировать обсценную лексику в отзывах",
              "item_type": "game_review",
              "condition_group": { // Группа условий (AND/OR)
                "operator": "AND",
                "conditions": [
                  { "field": "content.text", "operator": "contains_any_of", "values": ["плохое_слово1", "плохое_слово2"] }
                ],
                "groups": [] // Вложенные группы
              },
              "actions": [
                { "type": "auto_reject", "reason_code": "obscene_language" },
                { "type": "flag_for_review", "priority": "high" }
              ],
              "is_active": true,
              "priority_order": 10
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created):
        ```json
        {
          "data": {
            "type": "moderationRule",
            "id": "rule-uuid-123",
            "attributes": {
              "name": "Блокировать обсценную лексику в отзывах",
              "is_active": true,
              "priority_order": 10
            }
          }
        }
        ```
    *   Требуемые права доступа: `admin`, `lead_moderator`.
*   TODO: Эндпоинты для просмотра истории модерации конкретного пользователя или контента.

#### 3.1.3. Ресурс: Техническая Поддержка
*   **`GET /support/tickets`**
    *   Описание: Получение списка тикетов поддержки.
    *   Query параметры: `status` (enum: `open`, `pending_user`, `pending_agent`, `resolved`, `closed`), `assignee_id` (uuid или `me` или `unassigned`), `priority` (enum: `low`, `medium`, `high`, `urgent`), `user_id` (uuid), `search_query` (string), `page`, `per_page`, `sort_by`.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "supportTicket",
              "id": "ticket-uuid-1",
              "attributes": {
                "subject": "Не могу зайти в игру X",
                "status": "open",
                "priority": "high",
                "created_at": "2024-03-15T09:00:00Z",
                "updated_at": "2024-03-15T09:05:00Z",
                "last_reply_at": "2024-03-15T09:05:00Z"
              },
              "relationships": {
                "user": { "data": { "type": "platformUser", "id": "user-uuid-xyz" } },
                "assignee": { "data": { "type": "adminUser", "id": "support-agent-uuid-abc" } }
              },
              "links": { "self": "/api/v1/admin/support/tickets/ticket-uuid-1" }
            }
          ],
          "meta": { "total_items": 50, "current_page": 1, "per_page": 10 },
          "links": { "next": "/api/v1/admin/support/tickets?page=2" }
        }
        ```
    *   Требуемые права доступа: `support`, `admin`.
*   **`POST /support/tickets`**
    *   Описание: Создание тикета от имени пользователя (например, если обращение пришло по другому каналу, как телефон).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "supportTicket",
            "attributes": {
              "user_id": "user-uuid-xyz", // ID пользователя, если известен
              "user_email": "user@example.com", // Если ID не известен, но есть email
              "subject": "Проблема с оплатой",
              "description": "Пользователь сообщил, что не проходит оплата картой VISA.",
              "priority": "high",
              "category_id": "billing_issues",
              "tags": ["payment", "visa"]
            },
            "relationships": {
              "assignee": { "data": { "type": "adminUser", "id": "support-agent-uuid-abc" } } // Опционально
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Аналогично GET /support/tickets/{ticket_id})
    *   Требуемые права доступа: `support`, `admin`.
*   **`GET /support/tickets/{ticket_id}`**
    *   Описание: Получение информации о конкретном тикете, включая все ответы и историю изменений.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "supportTicket",
            "id": "ticket-uuid-1",
            "attributes": {
              "subject": "Не могу зайти в игру X",
              "description": "При попытке входа в игру X появляется ошибка 'Connection timed out'.",
              "status": "pending_agent",
              "priority": "high",
              "category_id": "game_issues",
              "tags": ["connection", "game_x"],
              "created_at": "2024-03-15T09:00:00Z",
              "updated_at": "2024-03-15T10:30:00Z"
            },
            "relationships": {
              "user": { "data": { "type": "platformUser", "id": "user-uuid-xyz" } },
              "assignee": { "data": { "type": "adminUser", "id": "support-agent-uuid-abc" } },
              "responses": {
                "data": [
                  { "type": "supportTicketResponse", "id": "response-uuid-1" },
                  { "type": "supportTicketResponse", "id": "response-uuid-2" }
                ],
                "links": { "related": "/api/v1/admin/support/tickets/ticket-uuid-1/responses" }
              }
            },
            "links": { "self": "/api/v1/admin/support/tickets/ticket-uuid-1" }
          },
          "included": [ // Опционально, если запрошено ?include=responses
            {
              "type": "supportTicketResponse",
              "id": "response-uuid-1",
              "attributes": {
                "text": "Здравствуйте! Пожалуйста, уточните, пробовали ли вы перезагрузить роутер?",
                "created_at": "2024-03-15T09:05:00Z",
                "is_internal_note": false
              },
              "relationships": {
                "author": { "data": { "type": "adminUser", "id": "support-agent-uuid-abc" } }
              }
            }
          ]
        }
        ```
    *   Требуемые права доступа: `support`, `admin`.
*   **`POST /support/tickets/{ticket_id}/responses`**
    *   Описание: Добавление ответа в тикет (от агента поддержки или как внутренняя заметка).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "supportTicketResponse",
            "attributes": {
              "text": "Мы проверили логи, проблема на стороне сервера игры. Передаем разработчикам.",
              "is_internal_note": false, // true для заметки, невидимой пользователю
              "update_ticket_status_to": "pending_third_party" // Опционально, изменить статус тикета
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Возвращает созданный ресурс ответа)
    *   Требуемые права доступа: `support`, `admin`.
*   **`PUT /support/tickets/{ticket_id}/status`**
    *   Описание: Изменение статуса тикета (open, pending_user, pending_agent, resolved, closed).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "supportTicketStatus",
            "attributes": {
              "status": "resolved",
              "resolution_comment": "Проблема решена перезапуском сервера игры."
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): (Возвращает обновленный тикет)
    *   Требуемые права доступа: `support`, `admin`.
*   **`GET /support/knowledge-base`**
    *   Описание: Поиск и получение статей из базы знаний.
    *   Query параметры: `search_query`, `category_id`, `tags`, `page`, `per_page`.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "knowledgeBaseArticle",
              "id": "kb-article-uuid-1",
              "attributes": {
                "title": "Как сбросить пароль?",
                "slug": "how-to-reset-password",
                "category_name": "Учетная запись",
                "tags": ["пароль", "сброс"],
                "last_updated_at": "2024-01-20T10:00:00Z"
              },
              "links": { "self": "/api/v1/admin/support/knowledge-base/kb-article-uuid-1" }
            }
          ],
          "meta": { "total_items": 15, "current_page": 1, "per_page": 5 }
        }
        ```
    *   Требуемые права доступа: `support`, `admin`, (может быть доступно и обычным пользователям через другой API).
*   **`POST /support/knowledge-base`**
    *   Описание: Создание новой статьи в базе знаний.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "knowledgeBaseArticle",
            "attributes": {
              "title": "Решение проблемы с черным экраном в игре Y",
              "content_markdown": "### Шаг 1: Обновите драйвера видеокарты...",
              "category_id": "game_troubleshooting",
              "tags": ["черный_экран", "игра_y", "драйвера"],
              "is_published": true
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Возвращает созданную статью)
    *   Требуемые права доступа: `admin`, `support_lead`, `knowledge_manager`.
*   **`GET /support/reports/performance`**
    *   Описание: Получение отчета по эффективности работы службы поддержки.
    *   Query параметры: `start_date`, `end_date`, `agent_id` (опционально), `group_by` (enum: `day`, `week`, `agent`).
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "supportPerformanceReport",
            "attributes": {
              "period_start": "2024-03-01T00:00:00Z",
              "period_end": "2024-03-15T23:59:59Z",
              "total_tickets_created": 250,
              "total_tickets_resolved": 230,
              "average_resolution_time_hours": 12.5,
              "first_response_time_median_minutes": 30,
              "satisfaction_rating_avg_percent": 85.2,
              "breakdown_by_agent": [ // если group_by=agent
                { "agent_id": "support-agent-uuid-abc", "resolved_tickets": 50, "avg_resolution_time_hours": 10.2 },
                { "agent_id": "support-agent-uuid-def", "resolved_tickets": 65, "avg_resolution_time_hours": 11.5 }
              ]
            }
          }
        }
        ```
    *   Требуемые права доступа: `admin`, `support_lead`.
*   TODO: Эндпоинты для управления категориями тикетов, назначения ответственных и шаблонов ответов.

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
                "global_announcement": null
              },
              "user_registration": {
                "allow_new_registrations": true,
                "default_user_roles": ["user"],
                "require_email_verification": true
              },
              "content_moderation": {
                "default_review_queue": "pending_manual",
                "profanity_filter_level": "medium"
              }
            }
          }
        }
        ```
    *   Требуемые права доступа: `admin`.
*   **`PUT /platform-settings`**
    *   Описание: Обновление настроек платформы (отправляются только изменяемые поля/группы).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "platformSettings",
            "attributes": {
              "general": {
                "maintenance_mode": true,
                "global_announcement": "Платформа будет недоступна с 02:00 до 03:00 МСК для проведения технических работ."
              }
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): (Возвращает все настройки после обновления)
    *   Требуемые права доступа: `admin`.
*   **`POST /platform-settings/maintenance-schedule`**
    *   Описание: Планирование окна обслуживания (включает/выключает maintenance_mode по расписанию).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "maintenanceSchedule",
            "attributes": {
              "start_time": "2024-03-20T02:00:00Z",
              "end_time": "2024-03-20T03:00:00Z",
              "message": "Плановые технические работы. Платформа будет временно недоступна."
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created):
        ```json
        {
          "data": {
            "type": "maintenanceSchedule",
            "id": "schedule-uuid-1",
            "attributes": {
              "start_time": "2024-03-20T02:00:00Z",
              "end_time": "2024-03-20T03:00:00Z",
              "message": "Плановые технические работы. Платформа будет временно недоступна.",
              "is_active": false, // станет true, когда наступит время
              "created_by": "admin-uuid-root"
            }
          }
        }
        ```
    *   Требуемые права доступа: `admin`.
*   TODO: Детальная структура настроек платформы и их возможные группы будут определены и описаны.

#### 3.1.5. Ресурс: Управление Администраторами
*   **`GET /admin-users`**
    *   Описание: Список административных пользователей.
    *   Query параметры: `role`, `is_active`, `page`, `per_page`.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "adminUser",
              "id": "admin-uuid-root",
              "attributes": {
                "username": "root_admin",
                "email": "root@example.com",
                "roles": ["superuser", "admin_management"],
                "is_active": true,
                "created_at": "2023-01-01T00:00:00Z"
              }
            }
          ]
        }
        ```
    *   Требуемые права доступа: `admin_management` (или `superuser`).
*   **`POST /admin-users`**
    *   Описание: Создание нового административного пользователя. Связывает существующего пользователя платформы (по `platform_user_id`) с административными ролями или создает нового пользователя в Auth Service, если необходимо.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "adminUser",
            "attributes": {
              "platform_user_id": "user-uuid-to-make-admin", // ID существующего пользователя платформы
              "username": "new_admin_user", // Используется, если platform_user_id не указан
              "email": "new_admin@example.com", // Используется, если platform_user_id не указан
              "password": "SecurePassword123!", // Используется, если platform_user_id не указан
              "roles": ["moderator", "support_agent_level1"],
              "is_active": true
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Возвращает созданного adminUser)
    *   Требуемые права доступа: `admin_management` (или `superuser`).
*   **`GET /admin-users/{admin_user_id}`**
    *   Описание: Получение информации о конкретном административном пользователе.
    *   Пример ответа (Успех 200 OK): (Аналогично элементу в GET /admin-users)
    *   Требуемые права доступа: `admin_management` (или `superuser`).
*   **`PUT /admin-users/{admin_user_id}/permissions`**
    *   Описание: Изменение ролей и прямых разрешений администратора.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "adminUserPermissions",
            "attributes": {
              "roles": ["lead_moderator", "support_agent_level2"],
              "direct_permissions": ["grant_temporary_platform_ban"], // Дополнительные, не входящие в роли
              "is_active": true // Можно также изменить статус активности
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): (Возвращает обновленного adminUser)
    *   Требуемые права доступа: `admin_management` (или `superuser`).
*   TODO: Эндпоинты для просмотра логов действий конкретного администратора.

### 3.2. gRPC API
*   В настоящее время gRPC API для Admin Service не планируется. Взаимодействие с другими внутренними сервисами осуществляется через их REST или gRPC API, а Admin Service предоставляет только REST API для административной панели.

### 3.3. WebSocket API (если применимо)
*   В настоящее время WebSocket API для Admin Service не планируется. Потребности в real-time дашбордах могут быть пересмотрены в будущем, но на данном этапе не являются приоритетом.

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности

*   **`AdminUser` (Администратор)**: Представляет учетную запись администратора или сотрудника поддержки.
    *   `id` (UUID, primary key): Уникальный идентификатор администратора. Пример: `a1b2c3d4-e5f6-7890-1234-567890abcdef`. Валидация: not null, unique. Обязательность: Required.
    *   `platform_user_id` (UUID, foreign key): Ссылка на ID пользователя в сервисе Account/Auth, если администратор также является пользователем платформы. Пример: `f0e1d2c3-b4a5-6789-0123-456789abcdef`. Валидация: unique. Обязательность: Optional (может быть чисто служебная учетная запись).
    *   `username` (VARCHAR(100)): Уникальное имя пользователя администратора. Пример: `admin_ivanov`. Валидация: not null, unique, min_length=3, max_length=100. Обязательность: Required.
    *   `email` (VARCHAR(255)): Email администратора. Пример: `ivanov.admin@example-platform.com`. Валидация: not null, unique, valid email format. Обязательность: Required.
    *   `password_hash` (VARCHAR(255)): Хеш пароля. Пример: `$2a$10$abcdef1234567890ABCD...`. Валидация: not null. Обязательность: Required (для администраторов без `platform_user_id`).
    *   `roles` (JSONB): Список ролей администратора (например, `["moderator", "support_level_2"]`). Пример: `["superuser"]`. Валидация: not null, valid array of strings. Обязательность: Required.
    *   `permissions` (JSONB): Список прямых разрешений (помимо тех, что дают роли). Пример: `["can_issue_permanent_ban"]`. Валидация: valid array of strings. Обязательность: Optional.
    *   `is_active` (BOOLEAN): Статус активности учетной записи. Пример: `true`. Валидация: not null. Обязательность: Required.
    *   `last_login_at` (TIMESTAMPTZ): Время последнего входа. Пример: `2024-03-15T10:00:00Z`. Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ): Время создания. Пример: `2023-01-01T00:00:00Z`. Валидация: not null. Обязательность: Required.
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. Пример: `2024-03-10T12:30:00Z`. Валидация: not null. Обязательность: Required.

*   **`ModerationItem` (Элемент модерации)**: Представляет объект (контент, жалоба), требующий модерации.
    *   `id` (UUID, primary key): Уникальный идентификатор элемента модерации. Пример: `m1n2o3p4-q5r6-7890-1234-567890abcdef`. Валидация: not null, unique. Обязательность: Required.
    *   `item_type` (VARCHAR(100)): Тип модерируемого контента (например, `game_review`, `user_profile`, `game_submission`, `user_complaint`). Пример: `game_review`. Валидация: not null. Обязательность: Required.
    *   `item_reference_id` (VARCHAR(255)): ID оригинального элемента в другом сервисе. Пример: `review-uuid-abc-123`. Валидация: not null. Обязательность: Required.
    *   `user_id` (UUID, foreign key): ID пользователя, создавшего контент (если применимо). Пример: `f0e1d2c3-b4a5-6789-0123-456789abcdef`. Обязательность: Optional.
    *   `status` (VARCHAR(50)): Текущий статус модерации (`pending`, `in_review`, `approved`, `rejected`, `escalated`, `action_taken`). Пример: `pending`. Валидация: not null, enum. Обязательность: Required.
    *   `priority` (VARCHAR(50)): Приоритет обработки (`low`, `normal`, `high`). Пример: `normal`. Валидация: not null, enum. Обязательность: Required.
    *   `content_snapshot` (JSONB): Снимок модерируемого контента на момент поступления в очередь. Пример: `{"text": "Эта игра ужасна!", "rating": 1}`. Обязательность: Required.
    *   `assigned_moderator_id` (UUID, foreign key, references AdminUser): ID администратора, которому назначен элемент. Пример: `a1b2c3d4-e5f6-7890-1234-567890abcdef`. Обязательность: Optional.
    *   `decision_made_at` (TIMESTAMPTZ): Время принятия решения. Обязательность: Optional.
    *   `decision_details` (JSONB): Детали решения (причина, комментарий модератора, принятые меры). Пример: `{"decision": "rejected", "reason_code": "profanity", "comment": "Ненормативная лексика."}`. Обязательность: Optional.
    *   `escalation_reason` (TEXT): Причина эскалации, если статус `escalated`. Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ): Время поступления в очередь. Обязательность: Required.
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. Обязательность: Required.

*   **`ModerationRule` (Правило модерации)**: Автоматическое правило для обработки элементов модерации.
    *   `id` (UUID, primary key): Уникальный идентификатор правила. Валидация: not null, unique. Обязательность: Required.
    *   `name` (VARCHAR(255)): Название правила. Пример: `Авто-отклонение спам-отзывов`. Валидация: not null. Обязательность: Required.
    *   `description` (TEXT): Описание правила. Обязательность: Optional.
    *   `item_type` (VARCHAR(100)): Тип контента, к которому применяется правило. Пример: `game_review`. Валидация: not null. Обязательность: Required.
    *   `conditions` (JSONB): Условия срабатывания правила (например, совпадение по ключевым словам, репутация пользователя). Пример: `{"field": "content_snapshot.text", "operator": "contains_any_of", "values": ["купить дешево", "продать аккаунт"]}`. Валидация: not null. Обязательность: Required.
    *   `actions` (JSONB): Действия, выполняемые при срабатывании правила (например, `auto_approve`, `auto_reject`, `flag_for_review`, `assign_to_group`). Пример: `[{"action_type": "auto_reject", "reason_code": "spam"}]`. Валидация: not null. Обязательность: Required.
    *   `priority_order` (INTEGER): Порядок применения правила. Пример: `10`. Валидация: not null. Обязательность: Required.
    *   `is_active` (BOOLEAN): Активно ли правило. Пример: `true`. Валидация: not null. Обязательность: Required.
    *   `created_by_admin_id` (UUID, foreign key, references AdminUser): Кто создал правило. Обязательность: Required.
    *   `updated_by_admin_id` (UUID, foreign key, references AdminUser): Кто последний раз обновил правило. Обязательность: Required.
    *   `created_at` (TIMESTAMPTZ): Время создания. Обязательность: Required.
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. Обязательность: Required.

*   **`SupportTicket` (Тикет поддержки)**: Запрос от пользователя в службу поддержки.
    *   `id` (UUID, primary key): Уникальный идентификатор тикета. Пример: `t1c2k3e4-t5f6-7890-1234-567890abcdef`. Валидация: not null, unique. Обязательность: Required.
    *   `user_id` (UUID, foreign key): ID пользователя платформы, создавшего тикет. Пример: `f0e1d2c3-b4a5-6789-0123-456789abcdef`. Обязательность: Required (если создан пользователем).
    *   `user_email_if_unregistered` (VARCHAR(255)): Email, если тикет создан от имени незарегистрированного пользователя. Обязательность: Optional.
    *   `subject` (VARCHAR(255)): Тема тикета. Пример: `Не могу запустить игру "Супер Гонки"`. Валидация: not null, max_length=255. Обязательность: Required.
    *   `description` (TEXT): Полное описание проблемы. Пример: `После обновления игра перестала запускаться...`. Валидация: not null. Обязательность: Required.
    *   `status` (VARCHAR(50)): Статус тикета (`open`, `pending_user_reply`, `pending_agent_reply`, `resolved`, `closed`, `on_hold`). Пример: `open`. Валидация: not null, enum. Обязательность: Required.
    *   `priority` (VARCHAR(50)): Приоритет (`low`, `medium`, `high`, `urgent`). Пример: `medium`. Валидация: not null, enum. Обязательность: Required.
    *   `category_id` (UUID, foreign key, references SupportTicketCategory): Категория тикета. Обязательность: Optional.
    *   `assignee_admin_id` (UUID, foreign key, references AdminUser): ID администратора (агента поддержки), которому назначен тикет. Обязательность: Optional.
    *   `tags` (JSONB): Теги для классификации. Пример: `["запуск_игры", "Супер Гонки"]`. Обязательность: Optional.
    *   `attachments` (JSONB): Ссылки на прикрепленные файлы. Пример: `[{"filename": "dxdiag.txt", "url": "..."}]`. Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ): Время создания. Обязательность: Required.
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. Обязательность: Required.
    *   `resolved_at` (TIMESTAMPTZ): Время решения. Обязательность: Optional.
    *   `closed_at` (TIMESTAMPTZ): Время закрытия. Обязательность: Optional.

*   **`SupportTicketResponse` (Ответ в тикете поддержки)**: Ответ агента или пользователя в рамках тикета.
    *   `id` (UUID, primary key): Уникальный идентификатор ответа. Валидация: not null, unique. Обязательность: Required.
    *   `ticket_id` (UUID, foreign key, references SupportTicket): ID тикета, к которому относится ответ. Валидация: not null. Обязательность: Required.
    *   `author_admin_id` (UUID, foreign key, references AdminUser): ID администратора-автора. Обязательность: Optional (если автор - пользователь).
    *   `author_user_id` (UUID): ID пользователя-автора. Обязательность: Optional (если автор - админ).
    *   `text` (TEXT): Текст ответа. Валидация: not null. Обязательность: Required.
    *   `is_internal_note` (BOOLEAN): Является ли ответ внутренней заметкой (невидимой пользователю). Пример: `false`. Валидация: not null. Обязательность: Required.
    *   `attachments` (JSONB): Ссылки на прикрепленные файлы в ответе. Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ): Время создания. Обязательность: Required.

*   **`KnowledgeBaseArticle` (Статья базы знаний)**: Статья для самопомощи пользователей.
    *   `id` (UUID, primary key): Уникальный идентификатор статьи. Валидация: not null, unique. Обязательность: Required.
    *   `title` (VARCHAR(255)): Заголовок статьи. Пример: `Как обновить драйверы видеокарты?`. Валидация: not null. Обязательность: Required.
    *   `slug` (VARCHAR(255)): URL-совместимый идентификатор статьи. Пример: `how-to-update-gpu-drivers`. Валидация: not null, unique. Обязательность: Required.
    *   `content_markdown` (TEXT): Содержимое статьи в формате Markdown. Валидация: not null. Обязательность: Required.
    *   `category_id` (UUID, foreign key, references KnowledgeBaseCategory): Категория статьи. Обязательность: Required.
    *   `tags` (JSONB): Теги для поиска. Пример: `["драйверы", "видеокарта", "обновление"]`. Обязательность: Optional.
    *   `is_published` (BOOLEAN): Опубликована ли статья. Пример: `true`. Валидация: not null. Обязательность: Required.
    *   `author_admin_id` (UUID, foreign key, references AdminUser): ID администратора, создавшего статью. Обязательность: Required.
    *   `created_at` (TIMESTAMPTZ): Время создания. Обязательность: Required.
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. Обязательность: Required.
    *   `view_count` (INTEGER, default: 0): Количество просмотров. Обязательность: Required.

*   **`PlatformSetting` (Настройка платформы)**: Глобальный или специфичный параметр конфигурации платформы.
    *   `id` (UUID, primary key): Уникальный идентификатор настройки. Валидация: not null, unique. Обязательность: Required.
    *   `key` (VARCHAR(255)): Уникальный ключ настройки (например, `general.maintenance_mode`, `user.registration.allow_new`). Валидация: not null, unique. Обязательность: Required.
    *   `value` (JSONB): Значение настройки (может быть строкой, числом, булевым, JSON объектом/массивом). Пример: `true` или `{"min_password_length": 8}`. Валидация: not null. Обязательность: Required.
    *   `description` (TEXT): Описание настройки. Пример: `Включает режим обслуживания платформы`. Обязательность: Optional.
    *   `category` (VARCHAR(100)): Категория настройки (например, `general`, `user_registration`, `moderation`). Пример: `general`. Обязательность: Optional.
    *   `is_editable_by_admin` (BOOLEAN, default: true): Можно ли редактировать через админ-панель. Обязательность: Required.
    *   `validation_rules` (JSONB): Правила валидации для значения (например, min/max, enum). Обязательность: Optional.
    *   `updated_by_admin_id` (UUID, foreign key, references AdminUser): Кто последний раз обновил настройку. Обязательность: Optional.
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. Обязательность: Required.

*   **`AuditLogAdmin` (Лог действий администраторов)**: Запись о действии, выполненном администратором.
    *   `id` (UUID, primary key): Уникальный идентификатор записи лога. Валидация: not null, unique. Обязательность: Required.
    *   `admin_user_id` (UUID, foreign key, references AdminUser): ID администратора, совершившего действие. Валидация: not null. Обязательность: Required.
    *   `timestamp` (TIMESTAMPTZ): Время совершения действия. Валидация: not null. Обязательность: Required.
    *   `action_type` (VARCHAR(100)): Тип действия (например, `user_ban`, `edit_platform_setting`, `resolve_ticket`). Пример: `user_ban`. Валидация: not null. Обязательность: Required.
    *   `target_entity_type` (VARCHAR(100)): Тип сущности, на которую было направлено действие (например, `PlatformUser`, `SupportTicket`). Обязательность: Optional.
    *   `target_entity_id` (VARCHAR(255)): ID сущности, на которую было направлено действие. Обязательность: Optional.
    *   `details` (JSONB): Детали действия (например, измененные поля, причина). Пример: `{"user_id": "uuid-user-1", "previous_status": "active", "new_status": "blocked", "reason": "spam"}`. Обязательность: Required.
    *   `ip_address` (VARCHAR(45)): IP-адрес администратора в момент совершения действия. Пример: `192.168.1.101`. Обязательность: Optional.

### 4.2. Схема Базы Данных
*   Ниже представлена ERD-диаграмма для основных таблиц PostgreSQL, используемых Admin Service.

```mermaid
erDiagram
    ADMIN_USERS {
        UUID id PK
        UUID platform_user_id FK "nullable"
        VARCHAR username UK
        VARCHAR email UK
        VARCHAR password_hash
        JSONB roles
        JSONB permissions
        BOOLEAN is_active
        TIMESTAMPTZ last_login_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    MODERATION_ITEMS {
        UUID id PK
        VARCHAR item_type
        VARCHAR item_reference_id
        UUID user_id FK "nullable"
        VARCHAR status
        VARCHAR priority
        JSONB content_snapshot
        UUID assigned_moderator_id FK "nullable"
        TIMESTAMPTZ decision_made_at
        JSONB decision_details
        TEXT escalation_reason
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    MODERATION_RULES {
        UUID id PK
        VARCHAR name
        TEXT description
        VARCHAR item_type
        JSONB conditions
        JSONB actions
        INTEGER priority_order
        BOOLEAN is_active
        UUID created_by_admin_id FK
        UUID updated_by_admin_id FK
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    SUPPORT_TICKETS {
        UUID id PK
        UUID user_id FK "nullable"
        VARCHAR user_email_if_unregistered
        VARCHAR subject
        TEXT description
        VARCHAR status
        VARCHAR priority
        UUID category_id FK "nullable"
        UUID assignee_admin_id FK "nullable"
        JSONB tags
        JSONB attachments
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
        TIMESTAMPTZ resolved_at
        TIMESTAMPTZ closed_at
    }

    SUPPORT_TICKET_RESPONSES {
        UUID id PK
        UUID ticket_id FK
        UUID author_admin_id FK "nullable"
        UUID author_user_id "nullable"
        TEXT text
        BOOLEAN is_internal_note
        JSONB attachments
        TIMESTAMPTZ created_at
    }

    SUPPORT_TICKET_CATEGORIES {
        UUID id PK
        VARCHAR name UK
        TEXT description
        BOOLEAN is_active
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    KNOWLEDGE_BASE_ARTICLES {
        UUID id PK
        VARCHAR title
        VARCHAR slug UK
        TEXT content_markdown
        UUID category_id FK
        JSONB tags
        BOOLEAN is_published
        UUID author_admin_id FK
        INTEGER view_count
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    KNOWLEDGE_BASE_CATEGORIES {
        UUID id PK
        VARCHAR name UK
        VARCHAR slug UK
        TEXT description
        INTEGER display_order
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    PLATFORM_SETTINGS {
        UUID id PK
        VARCHAR key UK
        JSONB value
        TEXT description
        VARCHAR category
        BOOLEAN is_editable_by_admin
        JSONB validation_rules
        UUID updated_by_admin_id FK "nullable"
        TIMESTAMPTZ updated_at
    }

    AUDIT_LOG_ADMIN {
        UUID id PK
        UUID admin_user_id FK
        TIMESTAMPTZ timestamp
        VARCHAR action_type
        VARCHAR target_entity_type
        VARCHAR target_entity_id
        JSONB details
        VARCHAR ip_address
    }

    ADMIN_USERS ||--o{ MODERATION_ITEMS : "assigned_to"
    ADMIN_USERS ||--o{ MODERATION_RULES : "created_by"
    ADMIN_USERS ||--o{ MODERATION_RULES : "updated_by"
    ADMIN_USERS ||--o{ SUPPORT_TICKETS : "assigned_to"
    ADMIN_USERS ||--o{ SUPPORT_TICKET_RESPONSES : "authored_by (admin)"
    ADMIN_USERS ||--o{ KNOWLEDGE_BASE_ARTICLES : "authored_by"
    ADMIN_USERS ||--o{ PLATFORM_SETTINGS : "updated_by"
    ADMIN_USERS ||--o{ AUDIT_LOG_ADMIN : "performed_by"
    SUPPORT_TICKETS ||--o{ SUPPORT_TICKET_RESPONSES : "belongs_to"
    SUPPORT_TICKET_CATEGORIES ||--o{ SUPPORT_TICKETS : "categorized_as"
    KNOWLEDGE_BASE_CATEGORIES ||--o{ KNOWLEDGE_BASE_ARTICLES : "categorized_as"

    MODERATION_ITEMS }o--|| ADMIN_USERS : "assigned_moderator_id"
    MODERATION_RULES }o--|| ADMIN_USERS : "created_by_admin_id"
    MODERATION_RULES }o--|| ADMIN_USERS : "updated_by_admin_id"
    SUPPORT_TICKETS }o--|| ADMIN_USERS : "assignee_admin_id"
    SUPPORT_TICKET_RESPONSES }o--|| ADMIN_USERS : "author_admin_id"
    KNOWLEDGE_BASE_ARTICLES }o--|| ADMIN_USERS : "author_admin_id"
    PLATFORM_SETTINGS }o--|| ADMIN_USERS : "updated_by_admin_id"
    AUDIT_LOG_ADMIN }o--|| ADMIN_USERS : "admin_user_id"
    SUPPORT_TICKET_RESPONSES }o--|| SUPPORT_TICKETS : "ticket_id"
    SUPPORT_TICKETS }o--|| SUPPORT_TICKET_CATEGORIES : "category_id"
    KNOWLEDGE_BASE_ARTICLES }o--|| KNOWLEDGE_BASE_CATEGORIES : "category_id"

```

*   **PostgreSQL DDL (детально):**
    ```sql
    -- Расширение для генерации UUID, если еще не создано
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    -- Таблица: admin_users (Администраторы)
    CREATE TABLE admin_users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        platform_user_id UUID UNIQUE, -- Ссылка на пользователя в Auth Service (может быть NULL)
        username VARCHAR(100) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password_hash VARCHAR(255), -- NULL если platform_user_id задан и используется SSO
        roles JSONB NOT NULL DEFAULT '[]'::jsonb, -- Список админских ролей
        permissions JSONB NOT NULL DEFAULT '[]'::jsonb, -- Список прямых разрешений
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        last_login_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
    COMMENT ON COLUMN admin_users.platform_user_id IS 'Ссылка на ID пользователя в сервисе Account/Auth';
    COMMENT ON COLUMN admin_users.password_hash IS 'Хеш пароля, если не используется platform_user_id для входа';

    -- Таблица: moderation_items (Элементы модерации)
    CREATE TABLE moderation_items (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        item_type VARCHAR(100) NOT NULL, -- 'game_review', 'user_profile', 'user_complaint', etc.
        item_reference_id VARCHAR(255) NOT NULL, -- ID сущности в другом сервисе
        user_id UUID, -- ID пользователя, создавшего контент (если применимо)
        status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'in_review', 'approved', 'rejected', 'escalated', 'action_taken'
        priority VARCHAR(50) NOT NULL DEFAULT 'normal', -- 'low', 'normal', 'high'
        content_snapshot JSONB NOT NULL, -- Снимок модерируемого контента
        assigned_moderator_id UUID REFERENCES admin_users(id) ON DELETE SET NULL,
        decision_made_at TIMESTAMPTZ,
        decision_details JSONB, -- {'decision': 'approved/rejected', 'reason_code': '...', 'comment': '...'}
        escalation_reason TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        CONSTRAINT fk_moderation_items_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL -- Предполагается наличие таблицы users
    );
    CREATE INDEX idx_moderation_items_status_type ON moderation_items(status, item_type);
    CREATE INDEX idx_moderation_items_assigned_moderator ON moderation_items(assigned_moderator_id);
    CREATE INDEX idx_moderation_items_item_reference ON moderation_items(item_type, item_reference_id);

    -- Таблица: moderation_rules (Правила модерации)
    CREATE TABLE moderation_rules (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(255) NOT NULL,
        description TEXT,
        item_type VARCHAR(100) NOT NULL, -- Тип контента, к которому применяется
        conditions JSONB NOT NULL, -- Условия срабатывания (например, regex, keyword match)
        actions JSONB NOT NULL, -- Действия (например, auto_approve, auto_reject, flag_for_review)
        priority_order INTEGER NOT NULL DEFAULT 0,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_by_admin_id UUID NOT NULL REFERENCES admin_users(id),
        updated_by_admin_id UUID NOT NULL REFERENCES admin_users(id),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

    -- Таблица: support_ticket_categories (Категории тикетов поддержки)
    CREATE TABLE support_ticket_categories (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(100) NOT NULL UNIQUE,
        description TEXT,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

    -- Таблица: support_tickets (Тикеты поддержки)
    CREATE TABLE support_tickets (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        user_id UUID, -- ID пользователя платформы (может быть NULL, если обращение от незарег.)
        user_email_if_unregistered VARCHAR(255),
        subject VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        status VARCHAR(50) NOT NULL DEFAULT 'open', -- 'open', 'pending_user_reply', 'pending_agent_reply', 'resolved', 'closed', 'on_hold'
        priority VARCHAR(50) NOT NULL DEFAULT 'medium', -- 'low', 'medium', 'high', 'urgent'
        category_id UUID REFERENCES support_ticket_categories(id) ON DELETE SET NULL,
        assignee_admin_id UUID REFERENCES admin_users(id) ON DELETE SET NULL,
        tags JSONB DEFAULT '[]'::jsonb,
        attachments JSONB DEFAULT '[]'::jsonb, -- Массив объектов {filename, url, size}
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        resolved_at TIMESTAMPTZ,
        closed_at TIMESTAMPTZ,
        CONSTRAINT fk_support_tickets_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL -- Предполагается наличие таблицы users
    );
    CREATE INDEX idx_support_tickets_status ON support_tickets(status);
    CREATE INDEX idx_support_tickets_assignee ON support_tickets(assignee_admin_id);
    CREATE INDEX idx_support_tickets_user ON support_tickets(user_id);

    -- Таблица: support_ticket_responses (Ответы в тикетах)
    CREATE TABLE support_ticket_responses (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        ticket_id UUID NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
        author_admin_id UUID REFERENCES admin_users(id) ON DELETE SET NULL, -- NULL если автор - пользователь
        author_user_id UUID, -- NULL если автор - админ (ссылка на users.id)
        text TEXT NOT NULL,
        is_internal_note BOOLEAN NOT NULL DEFAULT FALSE,
        attachments JSONB DEFAULT '[]'::jsonb,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        CONSTRAINT fk_support_responses_user FOREIGN KEY (author_user_id) REFERENCES users(id) ON DELETE SET NULL -- Предполагается наличие таблицы users
    );
    CREATE INDEX idx_support_ticket_responses_ticket_id ON support_ticket_responses(ticket_id);

    -- Таблица: knowledge_base_categories (Категории статей базы знаний)
    CREATE TABLE knowledge_base_categories (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(100) NOT NULL UNIQUE,
        slug VARCHAR(100) NOT NULL UNIQUE,
        description TEXT,
        display_order INTEGER DEFAULT 0,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

    -- Таблица: knowledge_base_articles (Статьи базы знаний)
    CREATE TABLE knowledge_base_articles (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        title VARCHAR(255) NOT NULL,
        slug VARCHAR(255) NOT NULL UNIQUE,
        content_markdown TEXT NOT NULL,
        category_id UUID NOT NULL REFERENCES knowledge_base_categories(id) ON DELETE RESTRICT,
        tags JSONB DEFAULT '[]'::jsonb,
        is_published BOOLEAN NOT NULL DEFAULT FALSE,
        author_admin_id UUID NOT NULL REFERENCES admin_users(id),
        view_count INTEGER NOT NULL DEFAULT 0,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
    CREATE INDEX idx_kb_articles_published_category ON knowledge_base_articles(is_published, category_id);

    -- Таблица: platform_settings (Настройки платформы)
    CREATE TABLE platform_settings (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        key VARCHAR(255) NOT NULL UNIQUE, -- e.g., 'general.maintenance_mode'
        value JSONB NOT NULL,
        description TEXT,
        category VARCHAR(100), -- e.g., 'general', 'user_registration'
        is_editable_by_admin BOOLEAN NOT NULL DEFAULT TRUE,
        validation_rules JSONB, -- e.g., {type: 'boolean'}, {type: 'integer', min: 0, max: 100}
        updated_by_admin_id UUID REFERENCES admin_users(id) ON DELETE SET NULL,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

    -- Таблица: audit_log_admin (Лог действий администраторов)
    CREATE TABLE audit_log_admin (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        admin_user_id UUID NOT NULL REFERENCES admin_users(id) ON DELETE RESTRICT,
        timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
        action_type VARCHAR(100) NOT NULL, -- e.g., 'user_ban', 'setting_change'
        target_entity_type VARCHAR(100), -- e.g., 'PlatformUser', 'SupportTicket'
        target_entity_id VARCHAR(255),
        details JSONB NOT NULL, -- { 'field_changed': 'status', 'old_value': 'active', 'new_value': 'banned', 'reason': '...' }
        ip_address VARCHAR(45)
    );
    CREATE INDEX idx_audit_log_admin_user_id ON audit_log_admin(admin_user_id);
    CREATE INDEX idx_audit_log_admin_action_type ON audit_log_admin(action_type);
    CREATE INDEX idx_audit_log_admin_target_entity ON audit_log_admin(target_entity_type, target_entity_id);
    ```
*   **MongoDB**: Используется для хранения сложных, неструктурированных или редко запрашиваемых деталей элементов модерации (например, полные логи чатов, большие медиа-файлы, если они не хранятся в S3). Также может использоваться для расширенных логов действий, требующих гибкой схемы.
    *   Пример коллекции `moderation_item_details`:
        ```json
        {
          "_id": "ObjectId(...)", // Соответствует moderation_items.id (UUID преобразованный)
          "original_content": { /* ... сложная структура ... */ },
          "violation_markers": [ { "type": "text_selection", "start": 10, "end": 25, "reason": "profanity" } ],
          "external_links_analysis": { /* ... */ }
        }
        ```
*   **Elasticsearch**: Индексы для полнотекстового поиска:
    *   `support_tickets_idx`: Индексируются поля `subject`, `description`, `responses.text`.
    *   `knowledge_base_articles_idx`: Индексируются `title`, `content_markdown`, `tags`.
    *   `admin_audit_log_idx`: Индексируются `action_type`, `target_entity_id`, `details` (ключевые поля из JSONB).
    *   `moderation_items_idx`: Индексируются ключевые поля из `content_snapshot` для поиска по содержимому.

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Kafka.
*   **Формат событий:** CloudEvents JSON.
*   **Основные топики:** `admin.events`.

*   **`admin.content.moderated.v1`**
    *   Описание: Контент прошел модерацию (одобрен/отклонен).
    *   Структура Payload: `{ "item_id": "...", "item_type": "...", "decision": "approved/rejected", "moderator_id": "...", "timestamp": "..." }`
    *   Потребители: Catalog Service, Social Service, Developer Service, Notification Service.
*   **`admin.user.status.changed.v1`**
    *   Описание: Статус пользователя изменен администратором.
    *   Структура Payload: `{ "user_id": "...", "new_status": "blocked/active", "admin_id": "...", "reason": "...", "timestamp": "..." }`
    *   Потребители: Auth Service, Account Service, Notification Service.
*   **`admin.platform.setting.updated.v1`**
    *   Описание: Изменена настройка платформы.
    *   Структура Payload: `{ "setting_key": "...", "old_value": "...", "new_value": "...", "admin_id": "...", "timestamp": "..." }`
    *   Потребители: Различные сервисы, которым важна эта настройка.
*   TODO: Другие релевантные публикуемые события будут определены и добавлены по мере развития функционала.

### 5.2. Потребляемые События (Consumed Events)
*   **`user.complaint.submitted.v1`** (от Social Service или других)
    *   Описание: Пользователь подал жалобу на контент или другого пользователя.
    *   Логика обработки: Создать новый `ModerationItem` в очереди.
*   **`developer.game.submitted.v1`** (от Developer Service)
    *   Описание: Разработчик отправил игру/обновление на модерацию.
    *   Логика обработки: Создать новый `ModerationItem` в очереди "Game Submissions".
*   **`payment.fraud.detected.v1`** (от Payment Service или Analytics Service)
    *   Описание: Обнаружена подозрительная финансовая активность.
    *   Логика обработки: Создать инцидент безопасности или задачу для администратора.
*   TODO: Другие релевантные потребляемые события будут определены и добавлены по мере развития интеграционных сценариев.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Account Service:** Для управления данными пользователей платформы (просмотр, изменение статуса, ролей).
*   **Auth Service:** Для управления административными пользователями и их правами; для получения информации о сессиях пользователей.
*   **Catalog Service:** Для модерации игрового контента, управления видимостью игр.
*   **Developer Service:** Для обработки заявок на публикацию игр от разработчиков.
*   **Social Service:** Для модерации пользовательского контента (отзывы, комментарии, сообщения).
*   **Payment Service:** Для обработки запросов на возврат средств, анализа транзакций, управления промо-акциями.
*   **Analytics Service:** Для получения данных для административных отчетов и дашбордов; для отправки событий о действиях админов.
*   **Notification Service:** Для отправки системных уведомлений, уведомлений пользователям по результатам модерации или решения тикетов.
*   **Download Service:** Для управления доступностью игровых файлов (например, после модерации).
*   **Library Service:** Для управления библиотеками пользователей в исключительных случаях.

### 6.2. Внешние Системы
*   TODO: Интеграции Admin Service с внешними системами (например, анти-чит, анализ контента) будут определены и описаны в последующих версиях, если потребуется.

## 7. Конфигурация (Configuration)

*   Общие стандарты конфигурационных файлов (формат YAML, структура, управление переменными окружения и секретами) определены в `project_api_standards.md` (раздел 7) и `DOCUMENTATION_GUIDELINES.md` (раздел 6). Данный раздел конкретизирует конфигурацию для Admin Service.

### 7.1. Переменные Окружения
*   `ADMIN_HTTP_PORT`: Порт для HTTP сервера Admin Service. Пример: `8080`. Значение по умолчанию: `8080`.
*   `ADMIN_GRPC_PORT`: Порт для gRPC сервера Admin Service (если используется). Пример: `9090`. Значение по умолчанию: `9090`.
*   `LOG_LEVEL`: Уровень логирования. Пример: `info`. Значение по умолчанию: `info`. (Допустимые: `error`, `warn`, `info`, `debug`, `trace`)
*   `POSTGRES_DSN_ADMIN`: DSN для подключения к основной базе данных PostgreSQL Admin Service. Пример: `postgresql://admin_user:password@postgres-admin:5432/admin_db?sslmode=disable`. Обязательная.
*   `MONGODB_URI_ADMIN`: URI для подключения к MongoDB (для деталей модерации, логов). Пример: `mongodb://mongo-admin:27017/admin_service_extra`. Обязательная, если используется MongoDB.
*   `ELASTICSEARCH_URLS_ADMIN`: URL(ы) Elasticsearch. Пример: `http://elasticsearch-admin:9200`. Обязательная, если используется Elasticsearch.
*   `REDIS_ADDR_ADMIN`: Адрес Redis (host:port). Пример: `redis-admin:6379`. Обязательная для кэширования.
*   `KAFKA_BROKERS`: Список брокеров Kafka. Пример: `kafka-main-1:9092,kafka-main-2:9092`. Обязательная для событийной интеграции.
*   `KAFKA_TOPIC_ADMIN_EVENTS`: Основной топик Kafka для публикуемых событий Admin Service. Пример: `admin.events.v1`. Значение по умолчанию: `admin.events.v1`.
*   `JWT_SECRET_KEY_ADMIN_API`: Секретный ключ для проверки JWT токенов, приходящих на Admin API (если проверка происходит на стороне сервиса, а не только Gateway). Хранится в Secrets.
*   `AUTH_SERVICE_GRPC_ADDR`: Адрес gRPC Auth Service. Пример: `auth-service:9000`. Обязательная.
*   `ACCOUNT_SERVICE_GRPC_ADDR`: Адрес gRPC Account Service. Пример: `account-service:9001`. Обязательная.
*   `CATALOG_SERVICE_API_URL`: URL REST API Catalog Service. Пример: `http://catalog-service:8000/api/v1`. Обязательная.
*   `PAYMENT_SERVICE_API_URL`: URL REST API Payment Service. Пример: `http://payment-service:8002/api/v1`. Обязательная.
*   `DEFAULT_MODERATION_QUEUE_SIZE`: Размер очереди модерации по умолчанию для некоторых операций. Пример: `100`. Значение по умолчанию: `50`.
*   `SUPPORT_TICKET_SLA_HOURS_HIGH_PRIORITY`: SLA в часах для ответа на тикеты с высоким приоритетом. Пример: `4`. Значение по умолчанию: `8`.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`: Эндпоинт Jaeger для экспорта трейсов OpenTelemetry. Пример: `http://jaeger-collector:14268/api/traces`.
*   `CORS_ALLOWED_ORIGINS_ADMIN`: Список разрешенных источников для CORS (для админ-панели). Пример: `https://admin.myplatform.com,http://localhost:3001`. Значение по умолчанию: `*` (не рекомендуется для production).

### 7.2. Файлы Конфигурации (если применимо)
*   Расположение: `configs/admin_config.yaml` (путь внутри контейнера или при монтировании).
*   Этот файл может использоваться для хранения конфигураций, которые не так часто меняются, как переменные окружения, или имеют более сложную структуру. Переменные окружения могут переопределять значения из файла.
*   Пример `configs/admin_config.yaml`:
    ```yaml
    server:
      http_port: ${ADMIN_HTTP_PORT:-8080} # Может быть переопределено переменной окружения
      grpc_port: ${ADMIN_GRPC_PORT:-9090}
      read_timeout_seconds: 30
      write_timeout_seconds: 30

    logging:
      level: ${LOG_LEVEL:-info}
      format: "json" # "text" or "json"

    database:
      postgres_dsn: ${POSTGRES_DSN_ADMIN} # Обязательно через env
      max_open_conns: 25
      max_idle_conns: 5
      conn_max_lifetime_minutes: 60

    redis:
      address: ${REDIS_ADDR_ADMIN} # Обязательно через env
      password: ${REDIS_PASSWORD_ADMIN:-} # Если Redis защищен паролем
      db: 0

    kafka:
      brokers: ${KAFKA_BROKERS} # Обязательно через env, список через запятую
      producer:
        retry_max: 5
        timeout_seconds: 10
      consumer_groups:
        user_complaints_group: "admin_user_complaints_consumer"

    moderation:
      default_items_per_page: 20
      rules_engine:
        # Пример встроенных правил, которые могут быть здесь, если не управляются через API/DB
        # - name: "Block specific keywords in reviews"
        #   item_type: "game_review"
        #   is_active: true
        #   conditions:
        #     field: "content_snapshot.text"
        #     operator: "contains_any_of_case_insensitive"
        #     values: ["badword1", "badword2"]
        #   actions:
        #     - type: "auto_reject"
        #       reason_code: "profanity"
        #     - type: "send_notification"
        #       template_id: "review_rejected_profanity"
        #   priority: 10
      queues:
        game_reviews:
          display_name: "Отзывы на игры"
          sla_hours: 24
        user_avatars:
          display_name: "Аватары пользователей"
          sla_hours: 12

    support:
      ticket_sla:
        high_priority_response_hours: ${SUPPORT_TICKET_SLA_HOURS_HIGH_PRIORITY:-8}
        medium_priority_response_hours: 24
        low_priority_response_hours: 72
      knowledge_base:
        recent_articles_count: 5

    security:
      audit_log_enabled: true
      # Настройки для rate limiting API эндпоинтов админки
      rate_limiter:
        enabled: true
        requests_per_second: 100
        burst: 200

    # Настройки для интеграции с другими сервисами (таймауты, ретраи)
    service_clients:
      auth_service:
        grpc_address: ${AUTH_SERVICE_GRPC_ADDR}
        timeout_seconds: 5
      account_service:
        grpc_address: ${ACCOUNT_SERVICE_GRPC_ADDR}
        timeout_seconds: 5
      # ... и так далее для других сервисов
    ```

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   Стандартные форматы ошибок REST API.
*   Подробное логирование для аудита и отладки.

### 8.2. Распространенные Коды Ошибок
*   **`USER_NOT_FOUND`**
*   **`ITEM_NOT_IN_MODERATION_QUEUE`**
*   **`INVALID_MODERATION_DECISION`**
*   **`TICKET_ALREADY_CLOSED`**
*   **`SETTING_VALIDATION_FAILED`**
*   (В дополнение к стандартным HTTP ошибкам)

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   JWT для доступа к Admin API, проверка ролей администратора.

### 9.2. Авторизация
*   Гранулярная RBAC для администраторов (супер-админ, модератор контента, специалист поддержки, финансовый администратор и т.д.).
*   (Ссылка на `project_roles_and_permissions.md`)

### 9.3. Защита Данных
*   Защита чувствительных данных пользователей и платформы, к которым имеют доступ администраторы.
*   Шифрование при передаче и хранении (где применимо).
*   (Ссылка на `project_security_standards.md`)

### 9.4. Управление Секретами
*   Использование Kubernetes Secrets или Vault.

### 9.5. Аудит Действий
*   Все действия администраторов должны подробно логироваться в специальный аудит-лог (`AuditLogAdmin`).

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   Dockerfile.
*   Helm-чарты/Kubernetes манифесты.
*   (Ссылка на `project_deployment_standards.md`)

### 10.2. Зависимости при Развертывании
*   PostgreSQL, MongoDB, Elasticsearch, Redis, Kafka.
*   Практически все остальные микросервисы платформы для выполнения своих функций.

### 10.3. CI/CD
*   Стандартный пайплайн.
*   (Ссылка на `project_deployment_standards.md`)

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   Формат: JSON.
*   Ключевые события: Действия администраторов, решения по модерации, изменения статусов тикетов, ошибки.
*   Интеграция: ELK/Loki.
*   (Ссылка на `project_observability_standards.md`)

### 11.2. Мониторинг
*   Метрики (Prometheus):
    *   Количество элементов в очередях модерации.
    *   Среднее время принятия решения по модерации.
    *   Количество открытых/закрытых тикетов поддержки.
    *   Активность администраторов.
    *   Ошибки API.
*   Дашборды (Grafana).
*   Алерты (AlertManager): Переполнение очередей модерации, большое количество нерешенных тикетов, ошибки сервиса.
*   (Ссылка на `project_observability_standards.md`)

### 11.3. Трассировка
*   Интеграция: OpenTelemetry, Jaeger.
*   (Ссылка на `project_observability_standards.md`)

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**:
    *   Время отклика API для интерактивных операций (получение списка, открытие сущности): P95 < 800 мс, P99 < 1500 мс.
    *   Время отклика API для операций изменения (создание, обновление, удаление): P95 < 500 мс, P99 < 1000 мс.
    *   Пропускная способность для пакетных операций (например, импорт правил модерации): будет определена по тестам нагрузки, но не менее 100 записей/сек.
    *   Обработка одного элемента в очереди ручной модерации (от открытия до принятия решения модератором): среднее время < 3 минут (без учета времени на само чтение/анализ контента).
    *   Загрузка главной страницы админ-панели со всеми виджетами: P95 < 2 секунд.
*   **Надежность**:
    *   Доступность сервиса: 99.9% (допустимое время простоя ~8.76 часов в год).
    *   RTO (Recovery Time Objective): < 1 часа для критических функций (управление пользователями, модерация).
    *   RPO (Recovery Point Objective): < 15 минут для данных конфигураций и тикетов поддержки; < 1 час для логов аудита.
    *   Ошибки на стороне сервера (5xx) для API: < 0.1% от всех запросов.
*   **Безопасность**:
    *   Соответствие `project_security_standards.md`.
    *   Все чувствительные данные (пароли, токены доступа к другим сервисам) должны храниться в зашифрованном виде или через систему управления секретами (Vault/Kubernetes Secrets).
    *   Регулярное сканирование на уязвимости.
    *   Защита от OWASP Top 10 уязвимостей.
    *   Обязательное логирование всех действий администраторов с достаточной детализацией для расследования инцидентов.
*   **Масштабируемость**:
    *   Горизонтальное масштабирование: Сервис должен поддерживать запуск нескольких экземпляров для распределения нагрузки.
    *   Способность обрабатывать до 1 миллиона активных пользователей платформы (с точки зрения запросов на данные пользователей, их жалоб и т.д.).
    *   Способность обрабатывать до 1000 одновременных сессий администраторов/модераторов.
    *   Очереди модерации должны справляться с пиковыми нагрузками до 1000 новых элементов в час без значительного увеличения времени ожидания.
    *   База данных должна быть рассчитана на рост объема данных (например, тикеты, логи аудита) на 50-100% в год без деградации производительности.
*   **Сопровождаемость**:
    *   Покрытие кода тестами (unit, integration): > 80%.
    *   Время развертывания новой версии в production: < 15 минут (после прохождения всех тестов).
    *   Логи должны быть структурированными и легко читаемыми.
    *   Метрики должны предоставлять ясное представление о состоянии сервиса.
    *   Документация (эта спецификация) должна поддерживаться в актуальном состоянии.

## 13. Приложения (Appendices) (Опционально)
*   TODO: Детальные JSON схемы для API, примеры конфигураций и другие приложения будут добавлены в последующих версиях документации или в отдельных специализированных документах.

---
*Этот документ является отправной точкой и должен регулярно обновляться по мере развития сервиса.*

## 14. Связанные Рабочие Процессы (Related Workflows)
*   [Developer Submits a New Game for Moderation](../../../project_workflows/game_submission_flow.md)
