# Спецификация Микросервиса: Social Service (Социальный Сервис)

**Версия:** 1.1 (адаптировано из предыдущей версии)
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Social Service платформы "Российский Аналог Steam".
*   **Роль в общей архитектуре платформы:** Social Service отвечает за реализацию и поддержку всех социальных взаимодействий между пользователями платформы. Это включает управление профилями пользователей, систему дружбы, создание и управление группами по интересам, обмен личными и групповыми сообщениями (чаты), формирование и отображение лент активности, систему отзывов и комментариев к играм и другому контенту, а также форумы для обсуждений.
*   **Основные бизнес-задачи:**
    *   Стимулирование общения и взаимодействия между пользователями платформы.
    *   Создание и поддержка активного игрового сообщества.
    *   Предоставление пользователям инструментов для самовыражения и поиска единомышленников.
    *   Повышение вовлеченности пользователей за счет социальных механик.
*   Разработка сервиса должна вестись в соответствии с `CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Управление профилями пользователей:** Создание и редактирование расширенных профилей (аватар, фон, описание, интересы, ссылки на соцсети), настройки приватности профиля, управление пользовательским контентом (например, скриншоты, арты), формирование персональной ленты активности, управление черным списком пользователей.
*   **Управление друзьями:** Отправка и принятие/отклонение запросов на дружбу, просмотр списка друзей, отображение онлайн-статусов друзей, поиск пользователей, рекомендации друзей (на основе общих игр, друзей и т.д.).
*   **Управление группами:** Создание публичных и приватных групп, управление членством (вступление, приглашение, исключение, роли в группе), публикация объявлений и новостей группы, создание обсуждений/подфорумов внутри группы.
*   **Обмен сообщениями (Чаты):** Личные (1-на-1) чаты, групповые чаты. Отправка текстовых сообщений, эмодзи, изображений (ссылками). Статусы доставки и прочтения сообщений. История сообщений. Уведомления о новых сообщениях.
*   **Лента активности:** Формирование персонализированной ленты событий от друзей (например, "друг X начал играть в Y", "друг Z получил достижение A") и из групп, на которые подписан пользователь. Возможность лайкать и комментировать элементы ленты. Фильтрация ленты.
*   **Отзывы и комментарии:** Возможность оставлять отзывы и оценки для игр и другого контента на платформе. Написание комментариев к отзывам, новостям, элементам ленты. Редактирование и удаление своих отзывов/комментариев. Система голосования за полезность отзывов.
*   **Форумы и обсуждения:** Создание и управление общими форумами по интересам или по конкретным играм. Создание тем, написание постов, модерирование. Подписка на темы.
*   **Модерация пользовательского контента:** Инструменты для модераторов (просмотр жалоб, удаление/скрытие контента, блокировка пользователей). Система жалоб на контент и поведение пользователей. Автоматическая фильтрация нежелательного контента (базовая).

### 1.3. Основные Технологии
*   **Язык программирования:** Go (основной язык для API и бизнес-логики, в соответствии с `project_technology_stack.md`). Python может использоваться для специфичных задач (например, ML-компоненты для рекомендаций, если будут).
*   **Базы данных (Polyglot Persistence):**
    *   PostgreSQL: Основная реляционная СУБД для структурированных данных (профили пользователей, информация о группах, структура форумов, отзывы, комментарии, данные модерации). Соответствует `project_technology_stack.md`.
    *   Apache Cassandra: Используется для хранения данных с высокой интенсивностью записи и чтения, таких как сообщения чатов и элементы ленты активности. Соответствует `project_technology_stack.md`.
    *   Neo4j: Используется для хранения и обработки социального графа (связи дружбы, рекомендации друзей). Соответствует `project_technology_stack.md`.
    *   Redis: Применяется для кэширования часто запрашиваемых данных, хранения онлайн-статусов пользователей, управления сессиями WebSocket и других эфемерных данных. Соответствует `project_technology_stack.md`.
*   **Брокер сообщений:** Apache Kafka (для асинхронного обмена событиями между Social Service и другими сервисами, а также для внутренних задач). Соответствует `project_technology_stack.md`.
*   **WebSocket:** Для обеспечения real-time функциональности чатов и обновлений ленты/статусов (например, стандартные библиотеки Go для WebSocket или Socket.IO, если фронтенд использует его).
*   **Поисковый движок (опционально):** Elasticsearch может использоваться для поиска по профилям, группам, форумам, если возможности PostgreSQL окажутся недостаточными.
*   **Инфраструктура:** Docker, Kubernetes.
*   **Мониторинг/Трассировка:** OpenTelemetry, Prometheus, Grafana, Jaeger/Tempo.
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `PACKAGE_STANDARDIZATION.md`.
*   *Примечание: Приведенные ниже примеры API и конфигураций основаны на предположении использования Go (Echo/Gin для REST), PostgreSQL, Cassandra, Neo4j, Redis и Kafka.*

### 1.4. Термины и Определения (Glossary)
*   **Профиль Пользователя (UserProfile):** Публичная или частично публичная страница пользователя с информацией о нем, его активности, друзьях и т.д.
*   **Социальный Граф (Social Graph):** Структура, представляющая пользователей как узлы и их социальные связи (например, дружба) как ребра.
*   **Группа (Group):** Сообщество пользователей, объединенных общими интересами или принадлежностью к чему-либо.
*   **Лента Активности (Activity Feed):** Хронологический список социальных событий, релевантных для пользователя.
*   **Отзыв (Review):** Мнение пользователя о продукте (игре, DLC), обычно с оценкой.
*   **UGC (User-Generated Content):** Контент, создаваемый пользователями (профили, сообщения, отзывы, посты на форумах).
*   Для других общих терминов см. `project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Social Service использует многослойную архитектуру (Clean Architecture) с разделением на компоненты, отвечающие за конкретные социальные функции.
*   Для различных типов данных используются специализированные хранилища: PostgreSQL для структурированных реляционных данных, Cassandra для данных с высокой интенсивностью записи/чтения (чаты, ленты), Neo4j для графовых данных (друзья), Redis для кэширования и эфемерных данных.
*   CQRS может применяться для некоторых компонентов (например, лента активности, где модель чтения оптимизирована). Event Sourcing может использоваться для чатов для обеспечения полной истории и восстановления состояния.

**Диаграмма Архитектуры:**
```mermaid
graph TD
    subgraph User Clients & API Gateway
        UserClient[Клиент Пользователя (Веб, Десктоп, Мобильный)] -- HTTP/WebSocket --> APIGW[API Gateway]
    end

    subgraph Social Service
        direction TB

        subgraph PresentationLayer [Presentation Layer (Адаптеры Транспорта)]
            REST_API[REST API (Echo/Gin)]
            GRPC_API[gRPC API (для межсервисного)]
            WebSocket_Hub[WebSocket Hub (Управление соединениями)]
            KafkaConsumers[Kafka Consumers (Входящие события)]
        end

        subgraph ApplicationLayer [Application Layer (Сценарии Использования)]
            ProfileAppSvc[Управление Профилями]
            FriendsAppSvc[Управление Друзьями]
            GroupsAppSvc[Управление Группами]
            ChatAppSvc[Чаты и Сообщения]
            FeedAppSvc[Лента Активности]
            ReviewAppSvc[Отзывы и Комментарии]
            ForumAppSvc[Форумы]
            ModerationAppSvc[Инструменты Модерации]
        end

        subgraph DomainLayer [Domain Layer (Бизнес-логика и Сущности)]
            Entities[Сущности (UserProfile, Friendship, Group, ChatMessage, FeedItem, Review)]
            Aggregates[Агрегаты (UserSocialGraph, GroupCommunity)]
            DomainEvents[Доменные События (FriendRequestSent, NewChatMessagePosted)]
            RepositoryIntf[Интерфейсы Репозиториев (PostgreSQL, Cassandra, Neo4j, Redis)]
        end

        subgraph InfrastructureLayer [Infrastructure Layer (Внешние Зависимости)]
            direction LR
            subgraph DataStores [Хранилища Данных]
                PostgresAdapter[Адаптер PostgreSQL] --> DB_PG[(PostgreSQL)]
                CassandraAdapter[Адаптер Cassandra] --> DB_Cass[(Cassandra)]
                Neo4jAdapter[Адаптер Neo4j] --> DB_Neo4j[(Neo4j)]
                RedisAdapter[Адаптер Redis] --> Cache[(Redis)]
            end
            KafkaProducer[Продюсер Kafka (Исходящие события)] --> KafkaBroker[Kafka Broker]
            InternalServiceClients[Клиенты др. микросервисов (Auth, Catalog, Account, Notification)]
        end

        APIGW -- HTTP/WebSocket --> PresentationLayer

        PresentationLayer --> ApplicationLayer
        ApplicationLayer --> DomainLayer
        ApplicationLayer --> InfrastructureLayer
        DomainLayer ----> RepositoryIntf
        InfrastructureLayer -- Implements --> RepositoryIntf
    end

    classDef layer_boundary fill:#f9f9f9,stroke:#333,stroke-width:2px,color:#333
    classDef component_major fill:#e6f0ff,stroke:#007bff,color:#000
    classDef datastore fill:#f8d7da,stroke:#dc3545,color:#000

    class PresentationLayer,ApplicationLayer,DomainLayer,InfrastructureLayer layer_boundary
    class REST_API,GRPC_API,WebSocket_Hub,KafkaConsumers,ProfileAppSvc,FriendsAppSvc,GroupsAppSvc,ChatAppSvc,FeedAppSvc,ReviewAppSvc,ForumAppSvc,ModerationAppSvc,Entities,Aggregates,DomainEvents,RepositoryIntf component_major
    class DB_PG,DB_Cass,DB_Neo4j,Cache,KafkaBroker datastore
```

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления)
*   **Ответственность:** Обработка всех входящих взаимодействий: REST API запросы от клиентов (через API Gateway), gRPC запросы от других микросервисов, управление WebSocket соединениями для real-time функций (чаты, уведомления о статусах), потребление событий из Kafka. Валидация данных, DTO преобразования.
*   **Ключевые компоненты/модули:** REST контроллеры (Echo/Gin), gRPC серверные реализации, WebSocket хаб/менеджер соединений, Kafka консьюмеры.

#### 2.2.2. Application Layer (Прикладной Слой)
*   **Ответственность:** Реализация сценариев использования для каждой социальной функции (профили, друзья, группы, чаты, ленты, отзывы, форумы, модерация). Координация операций между Domain Layer и Infrastructure Layer. Управление транзакциями на уровне приложения (если не делегировано ниже).
*   **Ключевые компоненты/модули:** Сервисы для каждой функциональной области (`ProfileService`, `FriendshipService`, `GroupService`, `ChatService`, `FeedService`, `ReviewService`, `ForumService`, `ModerationService`), обработчики команд и запросов (CQRS).

#### 2.2.3. Domain Layer (Доменный Слой)
*   **Ответственность:** Содержит бизнес-сущности, агрегаты, доменные события и бизнес-правила, специфичные для социальных взаимодействий. Определяет контракты репозиториев.
*   **Ключевые компоненты/модули:** Сущности (`UserProfile`, `FriendshipLink`, `Group`, `GroupMembership`, `ChatMessage`, `FeedItem`, `Review`, `Comment`, `ForumTopic`, `ForumPost`), объекты-значения (`UserID`, `GroupID`, `ContentText`), доменные сервисы (например, для логики рекомендаций друзей, если она сложная), доменные события.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   **Ответственность:** Реализация интерфейсов репозиториев для взаимодействия с различными базами данных (PostgreSQL, Cassandra, Neo4j, Redis). Взаимодействие с Kafka (продюсеры). Клиенты для других микросервисов (Auth, Account, Catalog, Notification, Admin). Управление конфигурацией, логированием.
*   **Ключевые компоненты/модули:** Реализации репозиториев для каждой БД, Kafka продюсер, gRPC/HTTP клиенты, утилиты для работы с WebSocket соединениями (если управление ими частично здесь).

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL:** `/api/v1/social` (маршрутизируется через API Gateway).
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway, `X-User-Id` передается в заголовках).
*   **Авторизация:** На основе `X-User-Id` и ролей (`X-User-Roles`), проверка приватности и владения ресурсами.
*   **Стандартный формат ответа об ошибке:**
    ```json
    {
      "errors": [
        {
          "id": "unique-error-instance-uuid-optional",
          "status": "4XX/5XX",
          "code": "ERROR_CODE_UPPER_SNAKE_CASE",
          "title": "Краткое описание ошибки на русском",
          "detail": "Полное описание ошибки с контекстом и, возможно, информацией о некорректных значениях.",
          "source": {
            "pointer": "/data/attributes/field_name",
            "parameter": "query_param_name"
          }
        }
      ]
    }
    ```
Пример `source` при ошибке валидации поля в теле запроса:
```json
    {
      "errors": [
        {
          "id": "c7a8b9f0-e1d2-4c3b-a890-1234567890ab",
          "status": "400",
          "code": "VALIDATION_ERROR",
          "title": "Ошибка валидации",
          "detail": "Поле 'nickname' не может быть пустым.",
          "source": {
            "pointer": "/data/attributes/nickname"
          }
        }
      ]
    }
```

#### 3.1.1. Профили Пользователей (User Profiles)
*   **`GET /users/{user_id}/profile`**
    *   Описание: Получение публичного или частично видимого профиля пользователя.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "userProfile",
            "id": "user-uuid-target",
            "attributes": {
              "nickname": "TargetUser",
              "avatar_url": "https://example.com/avatars/target.png",
              "status_message": "Играю в Супер Игру!",
              "profile_visibility": "public", // public, friends_only, private
              "last_online_at": "2024-03-15T10:00:00Z", // может быть приблизительным или скрытым
              "common_friends_count": 5, // если запрос от другого пользователя
              "custom_sections": [ { "title": "Мои любимые игры", "content_markdown": "- Супер Игра X\n- Приключение Z" } ]
            },
            "relationships": { "friends": { "links": { "related": "/api/v1/social/users/user-uuid-target/friends" } } }
          }
        }
        ```
    *   Требуемые права доступа: `public` (с учетом настроек приватности профиля) или `friend` (для доступа к доп. информации).
*   **`PUT /users/me/profile`**
    *   Описание: Обновление профиля текущего аутентифицированного пользователя.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "userProfileUpdate",
            "attributes": {
              "status_message": "Новый статус!",
              "profile_visibility": "friends_only",
              "custom_sections": [ { "title": "Обо мне", "content_markdown": "Всем привет!" } ]
            }
          }
        }
        ```
    *   Пример ответа (Успех 200 OK): (Обновленный профиль)
    *   Требуемые права доступа: `user_self_only`.

#### 3.1.2. Друзья (Friends)
*   **`POST /users/me/friends/requests`**
    *   Описание: Отправка запроса на добавление в друзья другому пользователю.
    *   Тело запроса: `{"data": {"type": "friendRequest", "attributes": {"target_user_id": "user-uuid-friend"}}} `
    *   Пример ответа (Успех 202 Accepted): `{ "meta": { "message": "Запрос на добавление в друзья отправлен." } }`
    *   Требуемые права доступа: `user_self_only`.
*   **`PUT /users/me/friends/requests/{request_id}`**
    *   Описание: Принятие или отклонение входящего запроса на дружбу.
    *   Тело запроса: `{"data": {"type": "friendRequestResponse", "attributes": {"action": "accept"}}} ` (action: `accept` или `reject`)
    *   Пример ответа (Успех 200 OK): `{ "data": { "type": "friendship", "id": "friendship-uuid-123", "attributes": {"status": "accepted"} } }`
    *   Требуемые права доступа: `user_self_only` (получатель запроса).

#### 3.1.3. Группы (Groups)
*   **`POST /groups`**
    *   Описание: Создание новой группы.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "groupCreation",
            "attributes": {
              "name": "Фанаты Супер Игры X",
              "description": "Обсуждаем все о Супер Игре X!",
              "group_type": "public" // public, private_invite_only, private_approval_required
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Данные созданной группы)
    *   Требуемые права доступа: `user` (любой аутентифицированный пользователь).

#### 3.1.4. Лента Активности (Feed)
*   **`GET /users/me/feed`**
    *   Описание: Получение персонализированной ленты активности для текущего пользователя.
    *   Query параметры: `before_event_id` (для пагинации "в прошлое"), `limit`.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "feedItem",
              "id": "feed-item-uuid-1",
              "attributes": {
                "actor": { "type": "user", "id": "friend-uuid-1", "nickname": "Друг1" },
                "verb": "unlocked_achievement", // posted_review, added_friend, joined_group
                "object": { "type": "achievement", "id": "ach-uuid-5", "name": "Мастер Исследователь" },
                "target": { "type": "game", "id": "game-uuid-123", "name": "Супер Игра X" }, // Опционально
                "timestamp": "2024-03-15T09:30:00Z",
                "likes_count": 5,
                "comments_count": 2
              }
            }
          ],
          "meta": { "next_cursor": "cursor_for_next_page" }
        }
        ```
    *   Требуемые права доступа: `user_self_only`.

#### 3.1.5. Отзывы (Reviews)
*   **`POST /products/{product_id}/reviews`**
    *   Описание: Публикация отзыва на продукт.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "reviewCreation",
            "attributes": {
              "rating": 5, // 1-5
              "title": "Отличная игра!",
              "content_text": "Мне очень понравилось, рекомендую всем!",
              "is_anonymous": false // Опционально
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Данные созданного отзыва)
    *   Требуемые права доступа: `user` (купивший игру).

### 3.2. gRPC API
*   Предназначен для внутреннего межсервисного взаимодействия.
*   Пакет: `social.v1`.
*   Определение Protobuf: `social/v1/social_internal_service.proto`.

#### 3.2.1. Сервис: `SocialInternalService`
*   **`rpc CheckFriendship(CheckFriendshipRequest) returns (CheckFriendshipResponse)`**
    *   Описание: Проверка статуса дружбы между двумя пользователями.
    *   `message CheckFriendshipRequest { string user_id_1 = 1; string user_id_2 = 2; }`
    *   `message CheckFriendshipResponse { enum Status { NO_RELATIONSHIP = 0; FRIENDS = 1; REQUEST_SENT_BY_1 = 2; REQUEST_SENT_BY_2 = 3; } Status status = 1; }`
*   **`rpc GetUserProfileSummary(GetUserProfileSummaryRequest) returns (GetUserProfileSummaryResponse)`**
    *   Описание: Получение краткой сводки по профилю пользователя для отображения в других сервисах.
    *   `message GetUserProfileSummaryRequest { string user_id = 1; }`
    *   `message UserProfileSummary { string user_id = 1; string nickname = 2; string avatar_url = 3; string online_status = 4; /* online, offline, busy */ }`
    *   `message GetUserProfileSummaryResponse { UserProfileSummary profile_summary = 1; }`
*   **`rpc BatchGetUsersProfileSummary(BatchGetUsersProfileSummaryRequest) returns (BatchGetUsersProfileSummaryResponse)`**
    *   Описание: Получение кратких сводок по профилям нескольких пользователей.
    *   `message BatchGetUsersProfileSummaryRequest { repeated string user_ids = 1; }`
    *   `message BatchGetUsersProfileSummaryResponse { repeated UserProfileSummary profile_summaries = 1; }`
*   *Примечание: Другие специфичные gRPC методы будут добавлены по мере финализации потребностей межсервисной интеграции.*

### 3.3. WebSocket API
*   **Эндпоинт:** `/api/v1/social/ws` (устанавливается через API Gateway).
*   **Аутентификация:** JWT токен при установлении WebSocket соединения.
*   **Сообщения от сервера к клиенту (примеры):**
    *   Новое сообщение в чате:
        ```json
        {
          "event_type": "chat.message.new",
          "payload": {
            "chat_room_id": "chat-room-uuid-personal-123", // или group_id
            "message_id": "msg-uuid-abc",
            "sender_id": "user-uuid-friend",
            "sender_nickname": "Друг1",
            "content_text": "Привет! Как дела?",
            "timestamp": "2024-03-15T10:00:00Z"
          }
        }
        ```
    *   Обновление статуса пользователя (онлайн/офлайн):
        ```json
        {
          "event_type": "user.status.updated",
          "payload": {
            "user_id": "user-uuid-friend",
            "new_status": "online", // "offline", "ingame"
            "last_seen_at": "2024-03-15T10:00:00Z" // если offline
          }
        }
        ```
    *   Новый элемент в ленте активности:
        ```json
        {
          "event_type": "feed.item.new",
          "payload": { /* ... структура FeedItem как в REST API ... */ }
        }
        ```
*   **Сообщения от клиента к серверу (примеры):**
    *   Отправка сообщения в чат:
        ```json
        {
          "action_type": "chat.message.send",
          "payload": {
            "chat_room_id": "chat-room-uuid-personal-123",
            "client_message_id": "local-temp-id-001", // для отслеживания доставки
            "content_text": "Все отлично, спасибо!"
          }
        }
        ```
    *   Уведомление о прочтении сообщений:
        ```json
        {
          "action_type": "chat.message.mark_read",
          "payload": {
            "chat_room_id": "chat-room-uuid-personal-123",
            "last_read_message_id": "msg-uuid-abc"
          }
        }
        ```

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности

*   **`UserProfile` (Профиль Пользователя - расширенный, хранится в PostgreSQL)**
    *   `user_id` (UUID, PK): ID пользователя (из Auth Service). Обязательность: Required.
    *   `nickname` (VARCHAR(100)): Никнейм (может отличаться от логина). Валидация: unique (опционально), profanity filter. Обязательность: Required.
    *   `avatar_url` (TEXT): URL аватара. Валидация: URL format. Обязательность: Optional.
    *   `profile_background_url` (TEXT): URL фона профиля. Обязательность: Optional.
    *   `status_message` (VARCHAR(255)): Короткий статус пользователя. Обязательность: Optional.
    *   `about_me_markdown` (TEXT): Раздел "Обо мне" в Markdown. Обязательность: Optional.
    *   `privacy_settings` (JSONB): Настройки приватности профиля. Пример: `{"profile_visibility": "public", "friend_list_visibility": "friends_only", "show_real_name": false}`. Обязательность: Required.
    *   `custom_sections` (JSONB): Пользовательские разделы профиля. Пример: `[{"title": "Мои любимые жанры", "content_markdown": "- RPG\n- Strategy"}]`. Обязательность: Optional.
    *   `last_online_at` (TIMESTAMPTZ): Время последнего захода в онлайн (приблизительное, из Redis). Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ), `updated_at` (TIMESTAMPTZ).

*   **`Friendship` (Дружба - хранится в Neo4j как отношение, метаданные в PostgreSQL)**
    *   `user_id1` (UUID, PK): ID первого пользователя.
    *   `user_id2` (UUID, PK): ID второго пользователя.
    *   `status` (ENUM: `pending_user1_to_user2`, `pending_user2_to_user1`, `friends`, `blocked_by_user1`, `blocked_by_user2`): Статус отношений. Обязательность: Required.
    *   `requested_at` (TIMESTAMPTZ): Время отправки запроса на дружбу. Обязательность: Optional.
    *   `accepted_at` (TIMESTAMPTZ): Время принятия дружбы. Обязательность: Optional.
    *   `blocked_at` (TIMESTAMPTZ): Время блокировки. Обязательность: Optional.
    *   *Neo4j Relationship:* `(:User {userId: "uuid1"})-[:FRIENDS_WITH {since: timestamp, status: "friends"}]->(:User {userId: "uuid2"})`

*   **`Group` (Группа - PostgreSQL)**
    *   `id` (UUID, PK): Уникальный идентификатор группы. Обязательность: Required.
    *   `name` (VARCHAR(255)): Название группы. Валидация: not null, unique. Обязательность: Required.
    *   `description` (TEXT): Описание группы. Обязательность: Optional.
    *   `avatar_url` (TEXT): URL аватара группы. Обязательность: Optional.
    *   `group_type` (ENUM: `public`, `private_invite_only`, `private_approval_required`): Тип группы. Обязательность: Required.
    *   `owner_user_id` (UUID, FK to User): ID создателя/владельца группы. Обязательность: Required.
    *   `created_at` (TIMESTAMPTZ).

*   **`ChatMessage` (Сообщение в Чате - Cassandra)**
    *   `chat_room_id` (UUID, PK (Partition Key)): ID комнаты чата (может быть ID пользователя для личного чата 1-на-1, или ID группы для группового чата). Обязательность: Required.
    *   `message_id` (TIMEUUID, PK (Clustering Key)): Уникальный, сортируемый по времени ID сообщения. Обязательность: Required.
    *   `sender_id` (UUID): ID отправителя. Обязательность: Required.
    *   `sender_nickname` (TEXT): Никнейм отправителя (денормализация для быстрого отображения). Обязательность: Required.
    *   `content_text` (TEXT): Текст сообщения. Валидация: max_length. Обязательность: Required.
    *   `attachments` (LIST<TEXT>): Список URL вложений (изображения, файлы). Обязательность: Optional.
    *   `created_at` (TIMESTAMP): Время создания сообщения (дублирует часть TimeUUID для запросов по диапазону). Обязательность: Required.
    *   `is_edited` (BOOLEAN), `is_deleted` (BOOLEAN).

*   **`FeedItem` (Элемент Ленты Активности - Cassandra)**
    *   `user_id` (UUID, PK (Partition Key)): ID пользователя, для которого предназначен этот элемент ленты. Обязательность: Required.
    *   `event_time` (TIMEUUID, PK (Clustering Key)): Время события, вызвавшего появление элемента в ленте. Обязательность: Required.
    *   `event_id` (UUID): Уникальный ID самого события-источника. Обязательность: Required.
    *   `actor_id` (UUID): ID пользователя, совершившего действие. Обязательность: Required.
    *   `actor_nickname` (TEXT): Никнейм актора. Обязательность: Required.
    *   `verb` (VARCHAR(50)): Тип действия (например, `posted_review`, `unlocked_achievement`, `added_friend`). Обязательность: Required.
    *   `object_id` (UUID): ID объекта действия (например, ID отзыва, ID достижения, ID друга). Обязательность: Optional.
    *   `object_type` (VARCHAR(50)): Тип объекта (например, `review`, `achievement`, `user`). Обязательность: Optional.
    *   `object_name_or_preview` (TEXT): Название или превью объекта (например, название игры для достижения). Обязательность: Optional.
    *   `target_id` (UUID): ID цели действия (например, ID игры, для которой оставлен отзыв). Обязательность: Optional.
    *   `target_type` (VARCHAR(50)): Тип цели. Обязательность: Optional.
    *   `target_name_or_preview` (TEXT): Название или превью цели. Обязательность: Optional.
    *   `likes_count` (COUNTER): Счетчик лайков.
    *   `comments_count` (COUNTER): Счетчик комментариев.

*   **`Review` (Отзыв - PostgreSQL)**
    *   `id` (UUID, PK): Уникальный идентификатор отзыва. Обязательность: Required.
    *   `product_id` (UUID): ID продукта (из Catalog Service). Обязательность: Required.
    *   `user_id` (UUID, FK to User): ID автора отзыва. Обязательность: Required.
    *   `rating` (SMALLINT): Оценка (например, 1-5 звезд). Валидация: 1-5. Обязательность: Required.
    *   `title` (VARCHAR(255)): Заголовок отзыва. Обязательность: Optional.
    *   `content_text` (TEXT): Текст отзыва. Валидация: not null, min/max length. Обязательность: Required.
    *   `is_anonymous` (BOOLEAN): Анонимный ли отзыв. Обязательность: Required, default false.
    *   `status` (ENUM: `pending_moderation`, `published`, `rejected`, `edited_after_publication`): Статус отзыва. Обязательность: Required.
    *   `positive_votes_count` (INTEGER), `negative_votes_count` (INTEGER).
    *   `created_at` (TIMESTAMPTZ), `updated_at` (TIMESTAMPTZ).

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL (Профили, Группы, Форумы, Отзывы, Связи Дружбы)

**ERD Диаграмма (ключевые таблицы PostgreSQL):**
```mermaid
erDiagram
    USER_SOCIAL_PROFILES {
        UUID user_id PK "FK to Auth.User"
        VARCHAR nickname UK
        TEXT avatar_url
        TEXT profile_background_url
        VARCHAR status_message
        TEXT about_me_markdown
        JSONB privacy_settings
        JSONB custom_sections
        TIMESTAMPTZ last_online_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    FRIENDSHIP_REQUESTS { # Для отслеживания запросов, основной граф в Neo4j
        UUID requester_user_id PK "FK to Auth.User"
        UUID target_user_id PK "FK to Auth.User"
        VARCHAR status -- pending, ignored
        TIMESTAMPTZ requested_at
    }
    GROUPS {
        UUID id PK
        VARCHAR name UK
        TEXT description
        TEXT avatar_url
        VARCHAR group_type -- public, private_invite, private_approval
        UUID owner_user_id FK "FK to Auth.User"
        TIMESTAMPTZ created_at
    }
    GROUP_MEMBERS {
        UUID group_id PK FK
        UUID user_id PK FK "FK to Auth.User"
        VARCHAR role_in_group -- admin, moderator, member
        TIMESTAMPTZ joined_at
    }
    REVIEWS {
        UUID id PK
        UUID product_id "FK to Catalog.Product"
        UUID user_id "FK to Auth.User"
        SMALLINT rating
        VARCHAR title
        TEXT content_text
        BOOLEAN is_anonymous
        VARCHAR status
        INTEGER positive_votes_count
        INTEGER negative_votes_count
        TIMESTAMPTZ created_at
    }
    COMMENTS { # Комментарии к отзывам, постам в ленте, статьям форума
        UUID id PK
        UUID parent_entity_id -- ID отзыва, поста ленты, поста форума
        VARCHAR parent_entity_type -- review, feed_item, forum_post
        UUID user_id "FK to Auth.User"
        TEXT content_text
        TIMESTAMPTZ created_at
    }

    USER_SOCIAL_PROFILES ||--o{ FRIENDSHIP_REQUESTS : "sends/receives"
    USER_SOCIAL_PROFILES ||--o{ GROUPS : "owns"
    USER_SOCIAL_PROFILES ||--o{ GROUP_MEMBERS : "is_member_of"
    GROUPS ||--o{ GROUP_MEMBERS : "has_members"
    USER_SOCIAL_PROFILES ||--o{ REVIEWS : "writes"
    USER_SOCIAL_PROFILES ||--o{ COMMENTS : "writes"
    REVIEWS ||--o{ COMMENTS : "has"
    PRODUCTS ||--o{ REVIEWS : "has" # PRODUCTS из Catalog Service
    USER_SOCIAL_PROFILES ||--o{ FORUM_POSTS : "creates"
    FORUM_TOPICS ||--o{ FORUM_POSTS : "has"
    FORUMS ||--o{ FORUM_TOPICS : "has"
    USER_SOCIAL_PROFILES ||--o{ FORUM_TOPICS : "creates" # Как автор темы
    USER_SOCIAL_PROFILES }o--|| FRIENDSHIP_REQUESTS : "sends/receives"


    entity PRODUCTS { # Предполагается из Catalog Service
        UUID id PK
        VARCHAR title
    }
```

**DDL (PostgreSQL - примеры):**
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_social_profiles (
    user_id UUID PRIMARY KEY, -- Ссылается на ID пользователя в Auth Service
    nickname VARCHAR(100) UNIQUE,
    avatar_url TEXT,
    profile_background_url TEXT,
    status_message VARCHAR(255),
    about_me_markdown TEXT,
    privacy_settings JSONB NOT NULL DEFAULT '{"profile_visibility": "public", "friend_list_visibility": "friends_only"}',
    custom_sections JSONB,
    last_online_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    avatar_url TEXT,
    group_type VARCHAR(50) NOT NULL CHECK (group_type IN ('public', 'private_invite_only', 'private_approval_required')),
    owner_user_id UUID NOT NULL, -- FK to users table in Auth Service
    members_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_groups_name_search ON groups USING GIN (to_tsvector('russian', name));

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL, -- FK to products table in Catalog Service
    user_id UUID NOT NULL, -- FK to users table in Auth Service
    rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(255),
    content_text TEXT NOT NULL,
    is_anonymous BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_moderation' CHECK (status IN ('pending_moderation', 'published', 'rejected', 'hidden_by_user', 'edited_after_publication')),
    positive_votes_count INTEGER NOT NULL DEFAULT 0,
    negative_votes_count INTEGER NOT NULL DEFAULT 0,
    language_code VARCHAR(10) DEFAULT 'ru-RU',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (product_id, user_id) -- Один пользователь - один отзыв на продукт
);
CREATE INDEX idx_reviews_product_id ON reviews(product_id);
CREATE INDEX idx_reviews_user_id ON reviews(user_id);

CREATE TABLE friendship_requests (
    requester_user_id UUID NOT NULL REFERENCES user_social_profiles(user_id) ON DELETE CASCADE,
    target_user_id UUID NOT NULL REFERENCES user_social_profiles(user_id) ON DELETE CASCADE,
    status VARCHAR(30) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected', 'ignored_by_target', 'cancelled_by_requester')),
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    responded_at TIMESTAMPTZ,
    PRIMARY KEY (requester_user_id, target_user_id)
);
CREATE INDEX idx_friendship_requests_target_id_status ON friendship_requests(target_user_id, status);

CREATE TABLE group_members (
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES user_social_profiles(user_id) ON DELETE CASCADE,
    role_in_group VARCHAR(50) NOT NULL DEFAULT 'member' CHECK (role_in_group IN ('admin', 'moderator', 'member')),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (group_id, user_id)
);
CREATE INDEX idx_group_members_user_id ON group_members(user_id);

CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_entity_id UUID NOT NULL, -- ID отзыва, поста ленты, поста форума и т.д.
    parent_entity_type VARCHAR(50) NOT NULL CHECK (parent_entity_type IN ('review', 'feed_item', 'forum_post', 'group_announcement')),
    user_id UUID NOT NULL REFERENCES user_social_profiles(user_id) ON DELETE CASCADE,
    content_text TEXT NOT NULL,
    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE, -- Мягкое удаление
    deleted_by_moderator BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_comments_parent_entity ON comments(parent_entity_id, parent_entity_type);
CREATE INDEX idx_comments_user_id ON comments(user_id);

CREATE TABLE forums (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    category VARCHAR(100), -- Общий, Игровой, Технический и т.д.
    is_locked BOOLEAN NOT NULL DEFAULT FALSE, -- Закрыт для новых тем/постов
    topics_count INTEGER NOT NULL DEFAULT 0,
    posts_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE forum_topics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    forum_id UUID NOT NULL REFERENCES forums(id) ON DELETE CASCADE,
    author_user_id UUID NOT NULL REFERENCES user_social_profiles(user_id),
    title VARCHAR(255) NOT NULL,
    is_sticky BOOLEAN NOT NULL DEFAULT FALSE, -- Закрепленная тема
    is_locked BOOLEAN NOT NULL DEFAULT FALSE, -- Закрыта для новых ответов
    views_count INTEGER NOT NULL DEFAULT 0,
    replies_count INTEGER NOT NULL DEFAULT 0,
    last_post_at TIMESTAMPTZ,
    last_post_user_id UUID REFERENCES user_social_profiles(user_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_forum_topics_forum_id_last_post_at ON forum_topics(forum_id, last_post_at DESC);
CREATE INDEX idx_forum_topics_author_id ON forum_topics(author_user_id);

CREATE TABLE forum_posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    topic_id UUID NOT NULL REFERENCES forum_topics(id) ON DELETE CASCADE,
    author_user_id UUID NOT NULL REFERENCES user_social_profiles(user_id),
    content_markdown TEXT NOT NULL,
    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_by_moderator BOOLEAN NOT NULL DEFAULT FALSE,
    position_in_topic SERIAL, -- Для примерной сортировки постов в теме, если не по created_at
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_forum_posts_topic_id_created_at ON forum_posts(topic_id, created_at ASC);
CREATE INDEX idx_forum_posts_author_id ON forum_posts(author_user_id);

```

#### 4.2.2. Apache Cassandra (Чаты, Ленты Активности)
*   **Роль:** Хранение данных с высокой интенсивностью записи и чтения, таких как сообщения в чатах и элементы ленты активности. Модель данных оптимизирована под конкретные запросы чтения.
*   **Пример DDL (CQL) для `chat_messages`:**
    ```cql
    CREATE TABLE IF NOT EXISTS social_service.chat_messages (
        chat_room_id uuid, // ID комнаты (может быть составным из user1_id + user2_id для личных, или group_id)
        message_id timeuuid,
        sender_id uuid,
        sender_nickname text,
        content_text text,
        attachments list<text>, // Список URL вложений
        created_at timestamp, // Обычно совпадает с временем из message_id
        is_edited boolean,
        PRIMARY KEY ((chat_room_id), message_id)
    ) WITH CLUSTERING ORDER BY (message_id DESC)
       AND compaction = {'class': 'TimeWindowCompactionStrategy', 'compaction_window_size': '7', 'compaction_window_unit': 'DAYS'};
    ```
*   **Пример DDL (CQL) для `user_activity_feed`:**
    ```cql
    CREATE TABLE IF NOT EXISTS social_service.user_activity_feed (
        user_id uuid,         // ID пользователя, чья это лента
        event_time timeuuid,   // Время события, для сортировки ленты
        actor_id uuid,        // ID пользователя, совершившего действие
        actor_nickname text,
        verb text,            // Тип действия (например, "posted_review", "added_friend")
        object_id uuid,       // ID объекта действия (ID отзыва, ID друга)
        object_type text,     // Тип объекта ("review", "user")
        object_preview text,  // Краткое описание/название объекта
        target_id uuid,       // ID цели действия (например, ID игры, для которой оставлен отзыв) (опционально)
        target_type text,     // Тип цели (опционально)
        target_preview text,  // Краткое описание/название цели (опционально)
        PRIMARY KEY (user_id, event_time)
    ) WITH CLUSTERING ORDER BY (event_time DESC);
    ```

#### 4.2.3. Neo4j (Социальный Граф)
*   **Роль:** Хранение и обработка связей дружбы между пользователями, поиск общих друзей, рекомендации друзей.
*   **Модель данных (узлы и отношения):**
    *   Узлы: `(:User {userId: "uuid", nickname: "Display Name"})`
    *   Отношения:
        *   `[:FRIENDS_WITH {since: datetime, status: "accepted"}]`
        *   `[:SENT_FRIEND_REQUEST_TO {requested_at: datetime}]`
        *   `[:BLOCKED {by_user_id: "uuid", at: datetime}]` (может быть и в PostgreSQL)
*   **Пример создания (Cypher):**
    ```cypher
    MERGE (u1:User {userId: $userId1}) ON CREATE SET u1.nickname = $nickname1;
    MERGE (u2:User {userId: $userId2}) ON CREATE SET u2.nickname = $nickname2;
    MERGE (u1)-[r:SENT_FRIEND_REQUEST_TO]->(u2) SET r.requested_at = datetime();
    ```

#### 4.2.4. Redis
*   **Роль:**
    *   Кэширование часто запрашиваемых данных: профили пользователей, списки друзей (ID), списки членов групп.
    *   Хранение онлайн-статусов пользователей: `user_status:<user_id>` (STRING: "online", "offline", "ingame:<game_id>"). TTL для статуса "online".
    *   Управление сессиями WebSocket (если сам сервис управляет, а не API Gateway).
    *   Кэширование сгенерированных лент активности (частично или полностью).
    *   Счетчики непрочитанных сообщений, уведомлений.

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Apache Kafka.
*   **Формат событий:** CloudEvents v1.0 (JSON encoding), согласно `project_api_standards.md`.
*   **Основной топик:** `social.events.v1` (может быть разделен на более гранулярные топики, например, `platform.social.profiles.v1`, `platform.social.friends.v1` и т.д., если потребуется).

*   **`com.platform.social.user.profile.updated.v1`** (Ранее `social.user.profile.updated.v1`)
    *   Описание: Профиль пользователя был обновлен (никнейм, аватар, статус и т.д.).
    *   Пример Payload (`data` секция CloudEvent):
        ```json
        {
          "user_id": "user-uuid-123",
          "updated_fields": ["nickname", "avatar_url"],
          "new_nickname": "SuperUser123",
          "new_avatar_url": "https://example.com/avatars/new.png",
          "updated_at": "2024-03-15T10:00:00Z"
        }
        ```
*   **`com.platform.social.friend.request.sent.v1`** (Ранее `social.friend.request.sent.v1`)
    *   Описание: Пользователь отправил запрос на добавление в друзья.
    *   Пример Payload:
        ```json
        {
          "requester_user_id": "user-uuid-sender",
          "target_user_id": "user-uuid-receiver",
          "request_id": "freq-uuid-abc",
          "sent_at": "2024-03-15T10:05:00Z"
        }
        ```
*   **`com.platform.social.friend.request.accepted.v1`**
    *   Описание: Пользователь принял запрос на добавление в друзья.
    *   Пример Payload:
        ```json
        {
          "accepter_user_id": "user-uuid-receiver", // Кто принял
          "requester_user_id": "user-uuid-sender", // Кто отправлял запрос
          "friendship_id": "fs-uuid-123", // Опционально, ID записи о дружбе
          "accepted_at": "2024-03-16T11:00:00Z"
        }
        ```
*   **`com.platform.social.review.submitted.v1`**
    *   Описание: Пользователь оставил отзыв о продукте.
    *   Пример Payload:
        ```json
        {
          "review_id": "review-uuid-abc",
          "user_id": "user-uuid-author",
          "product_id": "game-uuid-xyz",
          "product_type": "game", // game, dlc
          "rating": 5,
          "title_preview": "Отличная игра!", // Первые N символов заголовка
          "content_preview": "Мне очень понравилось...", // Первые N символов отзыва
          "submitted_at": "2024-03-16T12:00:00Z"
        }
        ```
*   **`com.platform.social.chat.message.sent.v1`** (Ранее `social.chat.message.sent.v1`)
    *   Описание: Новое сообщение отправлено в чат (личное или групповое).
    *   Пример Payload:
        ```json
        {
          "message_id": "msg-uuid-xyz",
          "chat_room_id": "chatroom-uuid-personal-or-group",
          "chat_type": "personal" , // "group"
          "sender_id": "user-uuid-sender",
          "recipient_id": "user-uuid-receiver", // Для личного чата
          "group_id": null, // Для группового чата
          "content_preview": "Привет! Как дела?...", // Первые N символов
          "sent_at": "2024-03-15T10:10:00Z"
        }
        ```
*   *Примечание: Другие события, такие как `com.platform.social.group.created.v1`, `com.platform.social.comment.posted.v1`, `com.platform.social.content.reported.v1` будут детализированы по аналогии по мере необходимости.*

### 5.2. Потребляемые События (Consumed Events)

*   **`account.user.created.v1`** (от Account Service)
    *   Описание: Новый пользователь зарегистрирован на платформе.
    *   Ожидаемый Payload: `{"user_id": "user-uuid-new", "username": "NewUser", "email": "new@example.com", "registration_timestamp": "..."}`
    *   Логика обработки: Создать начальный `UserProfile` в PostgreSQL. Создать узел `:User` в Neo4j.
*   **`library.achievement.unlocked.v1`** (от Library Service)
    *   Описание: Пользователь разблокировал достижение в игре.
    *   Ожидаемый Payload: `{"user_id": "user-uuid-123", "product_id": "game-uuid-abc", "achievement_id": "ach-meta-uuid-001", "achievement_name": "Первый шаг", "achievement_icon_url": "...", "unlocked_at": "..."}`
    *   Логика обработки: Создать элемент `FeedItem` типа `achievement_unlocked` для пользователя и, возможно, для его друзей (в зависимости от настроек приватности и логики формирования ленты).
*   **`catalog.game.review_period_opened.v1`** (от Catalog Service, гипотетическое)
    *   Описание: Продукт стал доступен для написания отзывов.
    *   Ожидаемый Payload: `{"product_id": "game-uuid-abc"}`
    *   Логика обработки: Разрешить создание отзывов для `product_id`.
*   **`moderation.user.sanction.applied.v1`** (от Admin Service)
    *   Описание: К пользователю применены санкции (например, бан чата, временная блокировка профиля).
    *   Ожидаемый Payload: `{"user_id": "user-uuid-123", "sanction_type": "chat_ban", "expires_at": "...", "reason": "..."}`
    *   Логика обработки: Обновить статус пользователя в Social Service, ограничить его возможности (например, отправку сообщений).

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Account Service:** Получение и обновление базовой информации о профиле пользователя (никнейм, аватар, если они управляются централизованно).
*   **Auth Service:** Аутентификация пользователей для доступа к API Social Service. Валидация токенов.
*   **Library Service:** Получение информации об играх пользователя, полученных достижениях для отображения в профиле и формирования событий ленты.
*   **Catalog Service:** Получение метаданных игр и другого контента для отображения в отзывах, обсуждениях, ленте.
*   **Notification Service:** Отправка уведомлений пользователям о новых сообщениях, запросах в друзья, ответах на комментарии, событиях в группах и т.д. (через Kafka).
*   **Admin Service:** Получение информации о статусах модерации контента, применение санкций к пользователям. Social Service предоставляет API для инструментов модерации.
*   **Analytics Service:** Social Service публикует события о социальной активности (лайки, комментарии, создание групп, дружба и т.д.), которые потребляются Analytics Service для анализа.
*   **API Gateway:** Маршрутизация REST API запросов и управление WebSocket соединениями.

### 6.2. Внешние Системы
*   В текущей версии не предполагается прямых интеграций с внешними социальными сетями для кросс-постинга или импорта друзей, кроме как через OAuth на стороне Auth Service для регистрации/входа.

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `SOCIAL_SERVICE_HTTP_PORT`: Порт для REST API (например, `8083`).
*   `SOCIAL_SERVICE_GRPC_PORT`: Порт для gRPC API (например, `9093`).
*   `SOCIAL_SERVICE_WS_PORT`: Порт для WebSocket (если обрабатывается напрямую).
*   `POSTGRES_DSN_SOCIAL`: DSN для PostgreSQL.
*   `CASSANDRA_HOSTS`: Хосты Cassandra (через запятую).
*   `CASSANDRA_KEYSPACE_SOCIAL`: Кейсспейс для Social Service в Cassandra.
*   `NEO4J_URI`, `NEO4J_USERNAME`, `NEO4J_PASSWORD`: Параметры подключения к Neo4j. Пароль (`NEO4J_PASSWORD`) должен загружаться из системы управления секретами.
*   `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB_SOCIAL`: Параметры Redis. Пароль (`REDIS_PASSWORD`) должен загружаться из системы управления секретами, если установлен.
*   `KAFKA_BROKERS`: Список брокеров Kafka.
*   `KAFKA_TOPIC_SOCIAL_EVENTS`: Топик для публикуемых событий.
*   `KAFKA_CONSUMER_GROUP_ID_SOCIAL`: ID группы консьюмеров.
*   `LOG_LEVEL`.
*   `AUTH_SERVICE_GRPC_ADDR`, `ACCOUNT_SERVICE_GRPC_ADDR`, `CATALOG_SERVICE_GRPC_ADDR`, `NOTIFICATION_SERVICE_KAFKA_TOPIC_SEND_REQUESTS`.
*   `MAX_CHAT_HISTORY_DAYS_DEFAULT`: Глубина хранения истории чатов по умолчанию.
*   `FEED_ITEMS_PAGE_SIZE_DEFAULT`: Количество элементов ленты на странице.
*   `MAX_FRIENDS_LIMIT_PER_USER`, `MAX_GROUPS_JOINED_PER_USER`.
*   `WEBSOCKET_MAX_CONNECTIONS`, `WEBSOCKET_WRITE_BUFFER_SIZE_BYTES`, `WEBSOCKET_READ_BUFFER_SIZE_BYTES`, `WEBSOCKET_PING_INTERVAL_SECONDS`.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`.

### 7.2. Файлы Конфигурации (если применимо)
*   **`configs/social_config.yaml`**: Может использоваться для более сложных или менее часто изменяемых настроек.
    ```yaml
    server:
      http_port: ${SOCIAL_SERVICE_HTTP_PORT:-8083}
      # ... другие серверные настройки

    feed_generation:
      default_page_size: ${FEED_ITEMS_PAGE_SIZE_DEFAULT:-25}
      # Параметры для алгоритмов формирования ленты (веса событий и т.д.)
      # event_weights: { "new_friend": 1.0, "game_achievement": 0.8, "review_posted": 0.7 }

    chat:
      max_message_length: 2000
      # Настройки хранения истории, если отличаются от глобальных
      # default_history_retention_days: ${MAX_CHAT_HISTORY_DAYS_DEFAULT:-90}

    moderation: # Базовые правила, если не вынесены в Admin Service или отдельную систему правил
      profanity_filter_enabled: true
      # profanity_dictionary_path: "/app/config/profanity_ru.txt"
      image_moderation_enabled: false # Если планируется UGC с картинками

    websocket:
      max_connections: ${WEBSOCKET_MAX_CONNECTIONS:-10000}
      write_buffer_size_bytes: ${WEBSOCKET_WRITE_BUFFER_SIZE_BYTES:-2048}
      read_buffer_size_bytes: ${WEBSOCKET_READ_BUFFER_SIZE_BYTES:-2048}
      ping_interval_seconds: ${WEBSOCKET_PING_INTERVAL_SECONDS:-30}

    # Настройки приватности по умолчанию для новых профилей
    default_privacy_settings:
      profile_visibility: "public" # friends_only, private
      friend_list_visibility: "friends_only"
      activity_feed_visibility: "friends_only"
    ```

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   REST API: Используются стандартные коды состояния HTTP. Тело ответа об ошибке соответствует формату, определенному в `project_api_standards.md`.
*   gRPC API: Используются стандартные коды состояния gRPC.
*   WebSocket: Ошибки передаются в виде специальных сообщений.
*   Использование механизмов retry и Circuit Breaker для межсервисных вызовов. Идемпотентность для операций, которые могут быть повторены.

### 8.2. Распространенные Коды Ошибок (для REST API)
*   **`400 Bad Request` (`VALIDATION_ERROR`)**: Некорректные входные данные.
*   **`401 Unauthorized` (`UNAUTHENTICATED`)**: Ошибка аутентификации.
*   **`403 Forbidden` (`PERMISSION_DENIED`, `USER_BLOCKED_TARGET`, `PROFILE_ACCESS_DENIED`)**: Недостаточно прав или доступ ограничен.
*   **`404 Not Found` (`PROFILE_NOT_FOUND`, `GROUP_NOT_FOUND`, `FRIEND_REQUEST_NOT_FOUND`)**: Ресурс не найден.
*   **`409 Conflict` (`FRIEND_REQUEST_ALREADY_EXISTS`, `ALREADY_FRIENDS`, `ALREADY_IN_GROUP`)**: Конфликт состояния.
*   **`422 Unprocessable Entity` (`MAX_FRIENDS_LIMIT_REACHED`, `MESSAGE_TOO_LONG`)**: Нарушение бизнес-правил.
*   **`500 Internal Server Error` (`INTERNAL_ERROR`)**: Внутренняя ошибка сервера.
*   **`503 Service Unavailable` (`DATABASE_UNAVAILABLE`)**: Сервис или его зависимости временно недоступны.

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   Все API (REST, gRPC, WebSocket) требуют JWT аутентификации через Auth Service.

### 9.2. Авторизация
*   RBAC для доступа к административным функциям (если есть).
*   Проверка владения ресурсами (например, пользователь может редактировать только свой профиль, отправлять сообщения от своего имени).
*   Учет настроек приватности профилей и контента.
*   Проверка статуса дружбы или членства в группе для доступа к соответствующим ресурсам.

### 9.3. Защита Данных
*   Соблюдение настроек приватности пользователей при отображении информации.
*   TLS для всех коммуникаций.
*   Шифрование чувствительных данных в покое (если таковые будут определены, например, приватные заметки в профиле).
*   **Чаты:** В первой версии End-to-End шифрование (E2EE) для личных сообщений не планируется, будет использоваться шифрование на транспортном уровне (TLS) и шифрование в покое для данных чатов в Cassandra. Вопрос внедрения E2EE может быть рассмотрен в будущем как отдельная задача, требующая значительных усилий по управлению ключами на клиентах.
*   Защита от спама и вредоносного контента через механизмы жалоб, модерации и автоматической фильтрации.
*   Валидация и санитизация всего пользовательского ввода для предотвращения XSS и других атак.
*   **ФЗ-152 "О персональных данных":** Все данные пользователей, включая профили, переписку в чатах и любой User Generated Content (UGC), обрабатываются в строгом соответствии с ФЗ-152 'О персональных данных', включая сбор согласий, хранение и защиту. См. также `project_security_standards.md`.
*   **Модерация UGC:** Механизмы для отправки жалоб на UGC (отзывы, комментарии, посты на форумах, сообщения в чатах) и их последующей модерации предоставляются и управляются через Admin Service. Social Service интегрируется с Admin Service для применения решений по модерации (например, скрытие или удаление контента, применение санкций к пользователям). См. также `backend/admin-service/docs/README.md`.

### 9.4. Управление Секретами
*   Пароли к базам данных (PostgreSQL, Neo4j, Redis, Cassandra, если используются), ключи для Kafka и другие секреты должны храниться в Kubernetes Secrets или HashiCorp Vault и безопасно внедряться в приложение во время выполнения. См. `project_security_standards.md`.

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Многоэтапная сборка для Go.
*   **Kubernetes манифесты/Helm-чарты.**
*   (Ссылка на `project_deployment_standards.md`).

### 10.2. Зависимости при Развертывании
*   Кластер Kubernetes.
*   PostgreSQL, Cassandra, Neo4j, Redis, Kafka.
*   Доступность Auth Service, Account Service, Catalog Service, Notification Service, Admin Service, API Gateway.

### 10.3. CI/CD
*   Автоматизированная сборка, юнит- и интеграционное тестирование (с использованием тестовых контейнеров для БД).
*   Развертывание в окружения с использованием GitOps.
*   Процедуры миграции схемы для PostgreSQL. Схемы для Cassandra и Neo4j также должны версионироваться и применяться управляемо.

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат:** Структурированные JSON логи (Zap).
*   **Ключевые события:** Все API запросы, операции с друзьями/группами, отправка/получение сообщений в чате, создание элементов ленты, операции модерации, ошибки.
*   **Интеграция:** С централизованной системой логирования (Loki, ELK Stack).
*   (Ссылка на `project_observability_standards.md`).

### 11.2. Мониторинг
*   **Метрики (Prometheus):**
    *   Количество API запросов (по типам, статусам). Длительность обработки.
    *   Количество активных WebSocket соединений.
    *   Количество отправленных/полученных сообщений в чатах.
    *   Скорость генерации элементов ленты.
    *   Количество операций с графом друзей.
    *   Ошибки при работе с PostgreSQL, Cassandra, Neo4j, Redis, Kafka.
*   **Дашборды (Grafana):** Для визуализации метрик.
*   **Алертинг (AlertManager):** Для критических ошибок, проблем с зависимостями, аномалий.
*   (Ссылка на `project_observability_standards.md`).

### 11.3. Трассировка
*   **Интеграция:** OpenTelemetry SDK для Go. Экспорт трейсов в Jaeger или Tempo.
*   Трассировка всех входящих запросов и исходящих вызовов.
*   (Ссылка на `project_observability_standards.md`).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:**
    *   REST API для чтения (профили, списки друзей, лента): P95 < 150 мс.
    *   REST API для записи (создание поста, отправка запроса в друзья): P95 < 100 мс.
    *   Доставка сообщений в чатах (WebSocket): P99 < 1 секунды (в пределах одного дата-центра).
    *   Генерация ленты активности для пользователя: P95 < 500 мс.
*   **Масштабируемость:**
    *   Поддержка до 10 миллионов зарегистрированных пользователей с профилями.
    *   Поддержка до 1 миллиона активных пользователей онлайн (статусы).
    *   Обработка до 10,000 сообщений в секунду в системе чатов.
    *   Горизонтальное масштабирование всех компонентов.
*   **Надежность:**
    *   Доступность сервиса: 99.95%.
    *   Гарантированная доставка сообщений в чатах (с учетом возможных сетевых проблем на клиенте).
    *   RPO < 5 минут для данных в PostgreSQL/Neo4j, RPO < 1 час для данных в Cassandra (допустима потеря недавних сообщений/ленты при катастрофе).
    *   RTO < 1 час.
*   **Согласованность данных:**
    *   Strong consistency для операций с профилями, группами, форумами в PostgreSQL.
    *   Eventual consistency для данных в Cassandra (ленты, чаты) и Neo4j (граф друзей может обновляться асинхронно). Статусы онлайн в Redis также eventually consistent.

## 13. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API, Protobuf определения для gRPC API, форматы сообщений WebSocket, а также полные DDL/CQL/Cypher схемы баз данных и их миграции поддерживаются в актуальном состоянии в соответствующих репозиториях исходного кода сервиса и во внутренней документации для разработчиков.
*   Примеры фискальных чеков и форматы взаимодействия с конкретными ОФД. (Это предложение здесь неуместно, удалено)

---
*Этот документ является основной спецификацией для Social Service и должен поддерживаться в актуальном состоянии.*

## 14. Резервное Копирование и Восстановление (Backup and Recovery)

### 14.1. Общая Стратегия
Social Service использует полиглотное хранилище, поэтому стратегия резервного копирования и восстановления должна учитывать особенности каждой используемой СУБД. Цель - минимизировать потерю данных (RPO) и время восстановления (RTO) для каждой категории данных в соответствии с их критичностью. Ссылки на общие принципы см. в `project_database_structure.md`.

### 14.2. PostgreSQL (Профили, Группы, Форумы, Отзывы)
*   **Метод:**
    *   Регулярные полные бэкапы (например, ежедневно).
    *   Непрерывное архивирование WAL (Write-Ahead Logging) для возможности Point-In-Time Recovery (PITR).
*   **Целевые показатели:**
    *   RPO: < 15 минут.
    *   RTO: < 2 часов.
*   **Хранение:** Географически распределенное, защищенное хранилище.
*   **Тестирование:** Регулярные учения по восстановлению.

### 14.3. Apache Cassandra (Чаты, Ленты Активности)
*   **Метод:**
    *   `nodetool snapshot` для создания снэпшотов на каждой ноде (например, ежедневно или еженедельно, в зависимости от объема данных и требований к RPO).
    *   Рассмотреть использование инкрементальных бэкапов (`nodetool backup` или сторонние инструменты, совместимые с Cassandra) для уменьшения объема хранимых бэкапов и потенциального улучшения RPO.
    *   Резервное копирование схемы (`cqlsh -e "DESC KEYSPACE social_service" > schema.cql`).
*   **Целевые показатели:**
    *   RPO: < 1-24 часов (в зависимости от частоты снэпшотов и использования инкрементальных бэкапов). Потеря самых последних сообщений в чате или элементов ленты может быть приемлема в случае катастрофического сбоя, так как эти данные часто имеют характер "горячих", но не критически важных для долгосрочного хранения в полном объеме.
    *   RTO: < 4-8 часов (сильно зависит от объема данных и количества нод).
*   **Хранение:** Снэпшоты и инкрементальные бэкапы должны копироваться с нод Cassandra в централизованное защищенное хранилище.
*   **Тестирование:** Регулярное тестирование процедур восстановления на тестовом кластере.

### 14.4. Neo4j (Социальный Граф)
*   **Метод:**
    *   Использование `neo4j-admin backup` для создания полных (full) или инкрементальных (incremental) бэкапов. Рекомендуется ежедневный полный бэкап или комбинация еженедельного полного и ежедневных инкрементальных.
*   **Целевые показатели:**
    *   RPO: < 30 минут (при использовании инкрементальных бэкапов чаще) или < 24 часа (при ежедневных полных).
    *   RTO: < 2 часов.
*   **Хранение:** Файлы бэкапов должны храниться в отдельном, защищенном хранилище.
*   **Тестирование:** Регулярное тестирование восстановления.

### 14.5. Redis (Кэш, Онлайн-статусы)
*   **Метод:**
    *   RDB снэпшоты (например, каждые 1-6 часов).
    *   AOF (Append-Only File) может быть включен, если данные в Redis (например, счетчики, которые сложно восстановить) считаются критичными.
*   **Целевые показатели:**
    *   RPO: < 1 часа (для RDB). С AOF `everysec` - до 1-2 секунд.
    *   RTO: < 30 минут.
*   **Примечание:** Большая часть данных в Redis для Social Service является кэшем (профили, списки друзей) или эфемерными данными (онлайн-статусы, сессии WebSocket). Эти данные часто могут быть перестроены из основных СУБД или потеряны без значительного ущерба. Стратегия бэкапирования Redis должна фокусироваться на тех данных, которые действительно сложно или невозможно восстановить другим способом.

### 14.6. Конфигурации и Секреты
*   Конфигурационные файлы сервиса версионируются в Git.
*   Секреты управляются через HashiCorp Vault или Kubernetes Secrets и бэкапируются в рамках процедур этих систем.

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [Процесс Модерации Пользовательского Контента (UGC)](../../../project_workflows/ugc_moderation_flow.md) (Описание будет добавлено в указанный документ).
*   [Процесс Формирования Рекомендаций Друзей и Контента](../../../project_workflows/social_recommendation_flow.md) (Описание будет добавлено в указанный документ).
*   [Процесс Обработки Жалоб Пользователей](../../../project_workflows/user_complaint_handling_flow.md) (Описание будет добавлено в указанный документ).
