<!-- backend\social-service\docs\README.md -->
# Спецификация Микросервиса: Social Service (Социальный Сервис)

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение:** Social Service отвечает за реализацию и поддержку всех социальных взаимодействий между пользователями платформы "Российский Аналог Steam". Это включает управление расширенными профилями пользователей, систему дружбы, создание и управление группами по интересам, обмен личными и групповыми сообщениями (чаты), формирование и отображение лент активности, систему отзывов и комментариев к играм и другому контенту, а также форумы для обсуждений.
*   **Роль в общей архитектуре платформы:** Является ключевым компонентом для построения сообщества вокруг платформы и игр, повышения вовлеченности пользователей и предоставления им инструментов для общения и самовыражения. Тесно интегрируется с другими сервисами для получения информации (например, о играх, пользователях, достижениях) и для инициирования уведомлений.
*   **Основные бизнес-задачи:**
    *   Стимулирование общения и взаимодействия между пользователями платформы.
    *   Создание и поддержка активного игрового сообщества.
    *   Предоставление пользователям инструментов для самовыражения и поиска единомышленников.
    *   Повышение вовлеченности пользователей за счет социальных механик и предоставления возможности публиковать текстовые материалы (например, отзывы, комментарии, посты на форумах).
    *   Обеспечение механизмов модерации для пользовательских текстовых материалов.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Управление Расширенными Профилями Пользователей:**
    *   Создание и редактирование профилей (никнейм, аватар, фон профиля, статусное сообщение, раздел "Обо мне").
    *   Управление пользовательскими витринами (например, любимые игры, достижения, скриншоты).
    *   Настройки приватности профиля (кто может просматривать профиль, список друзей, активность и т.д.).
    *   Управление черным списком пользователей (блокировка взаимодействий).
*   **Управление Друзьями (Социальный Граф):**
    *   Отправка, принятие, отклонение и отмена запросов на добавление в друзья.
    *   Просмотр списка друзей, их онлайн-статусов и текущей активности (игры, в которые играют).
    *   Удаление из списка друзей.
    *   (Опционально, [TODO: Реализовать предложения друзей]) Предложения друзей на основе общих игр, групп или друзей второго порядка.
*   **Управление Группами Пользователей:**
    *   Создание и настройка публичных и приватных групп (с различными моделями вступления: открытое, по заявкам, по приглашениям).
    *   Управление членством в группах (вступление, выход, приглашение, исключение, назначение ролей – администратор, модератор группы).
    *   Публикация объявлений и новостей внутри группы.
    *   Создание и модерирование обсуждений/подфорумов внутри группы.
*   **Обмен Сообщениями (Чаты):**
    *   Личные (1-на-1) чаты между друзьями.
    *   Групповые чаты для участников групп.
    *   Отправка текстовых сообщений, эмодзи. ([Поддержка изображений/файлов в чатах: будет рассмотрена на будущих этапах]).
    *   Отображение статусов доставки и прочтения сообщений.
    *   Доступ к истории сообщений с возможностью поиска ([История чатов: глубина поиска и хранения будет уточнена]).
    *   Уведомления о новых сообщениях (через Notification Service).
*   **Лента Активности:**
    *   Формирование персонализированной ленты событий от друзей (например, "друг X начал играть в Y", "друг Z получил достижение A", "друг W оставил отзыв об игре Q") и из групп, на которые подписан пользователь.
    *   Возможность лайкать и комментировать элементы ленты.
    *   Фильтрация и настройка отображения ленты.
*   **Отзывы и Комментарии:**
    *   Возможность оставлять текстовые отзывы и оценки (например, по 5-звездочной шкале или "рекомендую/не рекомендую") для игр и других продуктов платформы.
    *   Написание комментариев к отзывам, новостям, элементам ленты, постам на форумах.
    *   Редактирование и удаление собственных отзывов/комментариев (с соблюдением правил).
    *   Система голосования за полезность отзывов ("полезный" / "бесполезный").
*   **Форумы и Обсуждения:**
    *   Создание и управление общими форумами по интересам или по конкретным играм (может быть инициировано администраторами или разработчиками через Admin/Developer Service).
    *   Создание тем (тредов) на форумах.
    *   Публикация постов (ответов) в темах.
    *   Базовое модерирование форумов (прикрепление тем, закрытие тем, удаление постов).
    *   Подписка на темы для получения уведомлений о новых постах.
*   **Модерация Пользовательских Текстовых Материалов:**
    *   Интеграция с Admin Service для обработки жалоб на пользовательские текстовые материалы (например, тексты профилей, сообщения, отзывы, посты на форумах). [Стратегия модерации контента (например, автоматическая с использованием AI + ручная)].
    *   API для применения решений модерации (скрытие, удаление контента, применение санкций к пользователям).

### 1.3. Основные Технологии
*   **Язык программирования:** Go (основной язык для API и бизнес-логики, в соответствии с `../../../../project_technology_stack.md`). Python может использоваться для специфичных задач (например, ML-компоненты для рекомендаций, если будут).
*   **API:**
    *   REST API: Echo (`github.com/labstack/echo/v4`) (для клиентских приложений, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
    *   gRPC: `google.golang.org/grpc` (для внутреннего межсервисного взаимодействия, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
    *   WebSocket: `github.com/gorilla/websocket` (или аналогичная Go-библиотека) для обеспечения real-time функциональности чатов, обновлений ленты и онлайн-статусов. [Детали компонента реального времени (например, WebSockets, технология, масштабирование)].
*   **Базы данных (Polyglot Persistence):**
    *   PostgreSQL (версия 15+): Для структурированных данных (профили пользователей, информация о группах, структура форумов, отзывы, комментарии, метаданные для управления связями в Neo4j, если требуется). Драйвер: GORM (`gorm.io/gorm`) с `gorm.io/driver/postgres` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
    *   Apache Cassandra (версия 4.x+): Для хранения данных с высокой интенсивностью записи и чтения, таких как сообщения чатов и элементы ленты активности. Клиент: `github.com/gocql/gocql`. (согласно `../../../../project_technology_stack.md`).
    *   Neo4j (версия 5.x+): Для хранения и обработки социального графа (связи дружбы, членство в группах, рекомендации). Клиент: официальный Go драйвер `github.com/neo4j/neo4j-go-driver`. (согласно `../../../../project_technology_stack.md`).
    *   Redis (версия 7.0+): Применяется для кэширования часто запрашиваемых данных (профили, списки друзей), хранения онлайн-статусов пользователей, управления сессиями WebSocket, временных данных (например, счетчики непрочитанных сообщений). Клиент: `go-redis/redis` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Брокер сообщений:** Apache Kafka (клиент `github.com/confluentinc/confluent-kafka-go`, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Поисковый движок (опционально, [Elasticsearch для поиска: необходимость будет уточнена]):** Elasticsearch (версия 8.x+) может использоваться для поиска по профилям, группам, форумам, если возможности PostgreSQL окажутся недостаточными.
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`).
*   **Логирование:** Zap (`go.uber.org/zap`).
*   **Мониторинг/Трассировка:** OpenTelemetry SDK, Prometheus client.
*   **Инфраструктура:** Docker, Kubernetes.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   **Профиль Пользователя (UserProfile):** Публичная или частично публичная страница пользователя с информацией о нем, его активности, друзьях и т.д.
*   **Социальный Граф (Social Graph):** Структура, представляющая пользователей как узлы и их социальные связи (например, дружба) как ребра.
*   **Группа (Group):** Сообщество пользователей, объединенных общими интересами или принадлежностью к чему-либо.
*   **Лента Активности (Activity Feed):** Хронологический список социальных событий, релевантных для пользователя.
*   **Отзыв (Review):** Мнение пользователя о продукте (игре, DLC), обычно с оценкой.
*   **Пользовательские Текстовые Материалы (User-Submitted Text Content):** Текстовый контент, создаваемый пользователями в рамках функционала платформы, такой как данные профиля, сообщения, отзывы, комментарии, посты на форумах. Не включает игровые моды, предметы или другой сложный контент типа "Мастерской" (Workshop), который находится вне рамок проекта.
*   Для других общих терминов см. `../../../../project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Social Service использует многослойную архитектуру (Clean Architecture) с разделением на компоненты, отвечающие за конкретные социальные функции.
*   Для различных типов данных используются специализированные хранилища: PostgreSQL для структурированных реляционных данных, Cassandra для данных с высокой интенсивностью записи/чтения (чаты, ленты), Neo4j для графовых данных (друзья), Redis для кэширования и эфемерных данных.
*   CQRS может применяться для некоторых компонентов (например, лента активности, где модель чтения оптимизирована). Event Sourcing может использоваться для чатов для обеспечения полной истории и восстановления состояния, если это будет признано необходимым.

### 2.2. Диаграмма Архитектуры
```mermaid
graph TD
    subgraph UserClientsAndAPIGateway ["Клиенты Пользователя и API Gateway"]
        UserClient["Клиент Пользователя (Веб, Десктоп, Мобильный)"] -- HTTP/WebSocket --> APIGateway["API Gateway"]
    end

    subgraph SocialService ["Social Service (Чистая Архитектура)"]
        direction TB

        subgraph PresentationLayer [Presentation Layer (Адаптеры Транспорта)]
            REST_API[REST API (Echo) - для клиентских приложений]
            GRPC_API[gRPC API (для межсервисного взаимодействия)]
            WebSocket_Hub[WebSocket Hub (Управление соединениями для чатов, статусов, лент)]
            KafkaConsumers[Kafka Consumers (Входящие события от других сервисов)]
        end

        subgraph ApplicationLayer [Application Layer (Сценарии Использования)]
            ProfileAppSvc["Управление Профилями (вкл. расширенные данные)"]
            FriendshipAppSvc["Управление Друзьями и Социальным Графом"]
            GroupAppSvc["Управление Группами и Членством"]
            ChatAppSvc["Обработка Сообщений Чатов (личные, групповые)"]
            ActivityFeedAppSvc["Формирование и Управление Лентой Активности"]
            ReviewCommentAppSvc["Управление Отзывами и Комментариями"]
            ForumAppSvc["Управление Форумами, Темами и Постами"]
            ModerationIntegrationAppSvc["Интеграция с Сервисом Модерации"]
        end

        subgraph DomainLayer [Domain Layer (Бизнес-логика и Сущности)]
            Entities["Сущности (UserProfile, Friendship, Group, ChatMessage, FeedItem, Review, ForumTopic, etc.)"]
            Aggregates["Агрегаты (например, UserSocialProfile, GroupCommunity)"]
            DomainEvents["Доменные События (FriendRequestSent, NewChatMessage, ReviewPosted, etc.)"]
            RepositoryIntf["Интерфейсы Репозиториев (PostgreSQL, Cassandra, Neo4j, Redis)"]
            DomainServices["Доменные Сервисы (например, FeedGenerationService, FriendRecommendationService)"]
        end

        subgraph InfrastructureLayer [Infrastructure Layer (Внешние Зависимости и Реализации)"]
            direction LR
            subgraph DataStoreAdapters [Адаптеры Хранилищ Данных]
                PostgresAdapter[Адаптер PostgreSQL] --> DB_PG[(PostgreSQL)]
                CassandraAdapter[Адаптер Cassandra] --> DB_Cass[(Cassandra)]
                Neo4jAdapter[Адаптер Neo4j] --> DB_Neo4j[(Neo4j)]
                RedisAdapter[Адаптер Redis] --> Cache[(Redis)]
            end
            KafkaIO [Продюсеры/Консьюмеры Kafka] --> KafkaBroker[Kafka Message Bus]
            InternalServiceClients[Клиенты других микросервисов (Auth, Catalog, Account, Notification, Admin)]
            ConfigLoggingMonitoring[Конфигурация, Логирование, Мониторинг]
        end

        APIGateway -- HTTP/WebSocket/gRPC --> PresentationLayer

        PresentationLayer --> ApplicationLayer
        ApplicationLayer --> DomainLayer
        ApplicationLayer --> InfrastructureLayer
        DomainLayer ----> RepositoryIntf
        InfrastructureLayer -- Implements --> RepositoryIntf
    end

    InternalServiceClients --> OtherServices[Другие Микросервисы]


    classDef layer_boundary fill:#f9f9f9,stroke:#333,stroke-width:2px,color:#333
    classDef component_major fill:#e6f0ff,stroke:#007bff,color:#000
    classDef datastore fill:#f8d7da,stroke:#dc3545,color:#000
    classDef external_actor fill:#FEF9E7,stroke:#F1C40F,color:#000

    class PresentationLayer,ApplicationLayer,DomainLayer,InfrastructureLayer layer_boundary
    class REST_API,GRPC_API,WebSocket_Hub,KafkaConsumers,ProfileAppSvc,FriendshipAppSvc,GroupAppSvc,ChatAppSvc,ActivityFeedAppSvc,ReviewCommentAppSvc,ForumAppSvc,ModerationIntegrationAppSvc,Entities,Aggregates,DomainEvents,RepositoryIntf,DomainServices component_major
    class DB_PG,DB_Cass,DB_Neo4j,Cache,KafkaBroker datastore
    class UserClient,APIGateway,OtherServices external_actor
```

### 2.3. Слои Сервиса
(Описания слоев аналогичны предыдущим сервисам, с акцентом на специфику Social Service: управление профилями, друзьями, группами, чатами, лентами, отзывами, форумами и взаимодействие с разнородными БД.)

## 3. API Endpoints
(Раздел API как в существующем документе, с уточнением payload и добавлением эндпоинтов для форумов и расширенного управления группами).

### 3.1. REST API
*   **Базовый URL (через API Gateway):** `/api/v1/social`
*   (Примеры эндпоинтов REST API для управления профилями, друзьями, группами, отзывами, комментариями, форумами будут детализированы здесь, следуя `project_api_standards.md`. Например: `GET /users/{userId}/profile`, `POST /users/me/friends`, `GET /groups`, `POST /groups/{groupId}/posts`, `GET /products/{productId}/reviews`)

### 3.2. gRPC API
*   **Пакет:** `social.v1`
*   **Файл .proto:** `proto/social/v1/social_service.proto` (или в общем репозитории `platform-protos`)

*   **`rpc GetUserProfile(GetUserProfileRequest) returns (GetUserProfileResponse)`**
    *   Описание: Получение расширенного профиля пользователя.
    *   Пример запроса: `message GetUserProfileRequest { string user_id = 1; }`
    *   Пример ответа: `message UserProfileResponse { string user_id = 1; string nickname = 2; string avatar_url = 3; string status_message = 4; string about_me = 5; ... }`
    *   Ошибки: `NOT_FOUND`, `PERMISSION_DENIED`.
*   **`rpc UpdateUserProfile(UpdateUserProfileRequest) returns (UserProfileResponse)`**
    *   Описание: Обновление профиля текущего пользователя.
    *   Пример запроса: `message UpdateUserProfileRequest { string user_id = 1; optional string nickname = 2; optional string avatar_url = 3; ... }`
    *   Ошибки: `INVALID_ARGUMENT`, `UNAUTHENTICATED`, `PERMISSION_DENIED`.
*   **`rpc SendFriendRequest(SendFriendRequestRequest) returns (FriendRequestResponse)`**
    *   Описание: Отправка запроса на добавление в друзья.
    *   Пример запроса: `message SendFriendRequestRequest { string target_user_id = 1; }` // user_id отправителя из контекста
    *   Пример ответа: `message FriendRequestResponse { string request_id = 1; string status = 2; /* pending, already_friends, etc. */ }`
    *   Ошибки: `NOT_FOUND` (target_user_id), `ALREADY_EXISTS` (уже друзья или запрос отправлен), `FAILED_PRECONDITION` (нельзя отправить себе).
*   **`rpc AcceptFriendRequest(AcceptFriendRequestRequest) returns (AcceptFriendRequestResponse)`**
    *   Описание: Принятие запроса в друзья.
    *   Пример запроса: `message AcceptFriendRequestRequest { string requester_user_id = 1; }`
    *   Пример ответа: `message AcceptFriendRequestResponse { bool success = 1; }`
    *   Ошибки: `NOT_FOUND` (запрос не найден), `PERMISSION_DENIED`.
*   **`rpc GetFriendsList(GetFriendsListRequest) returns (GetFriendsListResponse)`**
    *   Описание: Получение списка друзей пользователя.
    *   Пример запроса: `message GetFriendsListRequest { string user_id = 1; int32 page_size = 2; string page_token = 3; }`
    *   Пример ответа: `message Friend { string user_id = 1; string nickname = 2; string avatar_url = 3; string online_status = 4; } message GetFriendsListResponse { repeated Friend friends = 1; string next_page_token = 2; }`
    *   Ошибки: `NOT_FOUND` (пользователь не найден).
*   **`rpc CreatePost(CreatePostRequest) returns (PostResponse)`**
    *   Описание: Создание поста (например, в группе или на форуме).
    *   Пример запроса: `message CreatePostRequest { string target_id = 1; /* groupId или forumTopicId */ string content = 2; }`
    *   Пример ответа: `message PostResponse { string post_id = 1; string author_id = 2; string content = 3; google.protobuf.Timestamp created_at = 4; }`
    *   Ошибки: `INVALID_ARGUMENT`, `PERMISSION_DENIED`.
*   **`rpc GetPosts(GetPostsRequest) returns (GetPostsResponse)`**
    *   Описание: Получение постов (из группы, темы форума, ленты).
    *   Пример запроса: `message GetPostsRequest { string target_id = 1; int32 page_size = 2; string page_token = 3; }`
    *   Пример ответа: `message GetPostsResponse { repeated PostResponse posts = 1; string next_page_token = 2; }`
*   (Аналогичные gRPC методы для других функций: управление группами, отзывами, комментариями, форумами).

### 3.3. WebSocket API
*   **Эндпоинт:** `/ws/social/updates` (требует аутентификации).
*   **Назначение:** Real-time доставка сообщений чата, уведомлений о новой активности в ленте, изменений онлайн-статусов друзей.
*   **Формат сообщений:** JSON, согласно `project_api_standards.md`.
    *   Пример нового сообщения в чате: `{"type": "social.chat.message.new", "payload": {"chat_room_id": "...", "message_id": "...", "sender_id": "...", "text": "...", "timestamp": "..."}}`
    *   Пример обновления статуса друга: `{"type": "social.friend.status.updated", "payload": {"user_id": "...", "online_status": "in_game", "current_game_id": "..."}}`

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`UserProfile` (Расширенный Профиль Пользователя - PostgreSQL)**
    *   `user_id` (UUID, PK, FK на User в Auth/Account Service). **Обязательность: Да.**
    *   `custom_status_message` (VARCHAR, Nullable). **Обязательность: Нет.**
    *   `about_me_text` (TEXT, Nullable). **Обязательность: Нет.**
    *   `profile_background_url` (VARCHAR, Nullable). **Обязательность: Нет.**
    *   `showcase_items` (JSONB, Nullable): Витрина (любимые игры, достижения). **Обязательность: Нет.**
    *   `privacy_settings` (JSONB): Настройки приватности (кто видит профиль, активность, список друзей). **Обязательность: Да (DEFAULT '{}').**
    *   `last_online_timestamp` (TIMESTAMPTZ, Nullable). **Обязательность: Нет (обновляется через Redis).**
    *   `created_at`, `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируются БД).**
*   **`Friendship` (Дружба - Neo4j)**
    *   Отношение `[:FRIENDS_WITH]` между двумя узлами `:User`.
    *   Свойства: `since` (TIMESTAMPTZ). **Обязательность: Да.**
    *   `status` (ENUM: `pending_request`, `accepted`, `blocked`). **Обязательность: Да.**
*   **`Group` (Группа - PostgreSQL)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `name` (VARCHAR, UK). **Обязательность: Да.**
    *   `description` (TEXT, Nullable). **Обязательность: Нет.**
    *   `avatar_url` (VARCHAR, Nullable). **Обязательность: Нет.**
    *   `owner_user_id` (UUID, FK). **Обязательность: Да.**
    *   `group_type` (ENUM: `public`, `private`, `invite_only`). **Обязательность: Да.**
    *   `member_count` (INTEGER, Default: 1). **Обязательность: Да.**
    *   `created_at`, `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируются БД).**
*   **`GroupMember` (Участник Группы - PostgreSQL, и отношение в Neo4j)**
    *   `group_id` (UUID, PK, FK). **Обязательность: Да.**
    *   `user_id` (UUID, PK, FK). **Обязательность: Да.**
    *   `role` (ENUM: `member`, `moderator`, `admin`). **Обязательность: Да (DEFAULT 'member').**
    *   `joined_at` (TIMESTAMPTZ). **Обязательность: Да (DEFAULT now()).**
*   **`ChatMessage` (Сообщение Чата - Cassandra)**
    *   `chat_room_id` (TEXT, PK - партишн ключ, может быть `user1_id::user2_id` или `group_id`). **Обязательность: Да.**
    *   `message_id` (TIMEUUID, PK - кластерный ключ). **Обязательность: Да.**
    *   `sender_id` (UUID). **Обязательность: Да.**
    *   `sender_nickname` (TEXT). **Обязательность: Да.**
    *   `content_text` (TEXT). **Обязательность: Да.**
    *   `attachments` (LIST<TEXT>, Nullable): Ссылки на медиа. **Обязательность: Нет.**
    *   `created_at` (TIMESTAMP). **Обязательность: Да.**
    *   `is_edited` (BOOLEAN, Default: false). **Обязательность: Да.**
*   **`ActivityFeedItem` (Элемент Ленты Активности - Cassandra)**
    *   `user_id` (UUID, PK - партишн ключ, для чьей ленты это событие). **Обязательность: Да.**
    *   `event_time` (TIMEUUID, PK - кластерный ключ). **Обязательность: Да.**
    *   `actor_id` (UUID): Кто совершил действие. **Обязательность: Да.**
    *   `actor_nickname` (TEXT). **Обязательность: Да.**
    *   `verb` (TEXT): Тип действия (e.g., "unlocked_achievement", "became_friends_with", "posted_review"). **Обязательность: Да.**
    *   `object_id` (UUID, Nullable): ID объекта действия. **Обязательность: Нет.**
    *   `object_type` (TEXT, Nullable): Тип объекта (game, achievement, review). **Обязательность: Нет.**
    *   `object_preview` (TEXT, Nullable): Краткое описание/ссылка на объект. **Обязательность: Нет.**
    *   `target_id` (UUID, Nullable): ID цели действия (если есть). **Обязательность: Нет.**
    *   `target_type` (TEXT, Nullable): Тип цели. **Обязательность: Нет.**
    *   `target_preview` (TEXT, Nullable). **Обязательность: Нет.**
*   **`Review` (Отзыв - PostgreSQL)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `product_id` (UUID, FK). **Обязательность: Да.**
    *   `user_id` (UUID, FK). **Обязательность: Да.**
    *   `rating_score` (SMALLINT, Nullable): Оценка (1-5 или 0/1). **Обязательность: Нет.**
    *   `review_text` (TEXT). **Обязательность: Да.**
    *   `status` (ENUM: `published`, `pending_moderation`, `rejected`, `hidden`). **Обязательность: Да.**
    *   `upvotes` (INTEGER, Default: 0). **Обязательность: Да.**
    *   `downvotes` (INTEGER, Default: 0). **Обязательность: Да.**
    *   `created_at`, `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируются БД).**
*   **`Comment` (Комментарий - PostgreSQL)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `parent_entity_type` (ENUM: `review`, `post`, `feed_item`). **Обязательность: Да.**
    *   `parent_entity_id` (UUID). **Обязательность: Да.**
    *   `user_id` (UUID, FK). **Обязательность: Да.**
    *   `comment_text` (TEXT). **Обязательность: Да.**
    *   `status` (ENUM: `published`, `pending_moderation`, `rejected`, `hidden`). **Обязательность: Да.**
    *   `created_at`, `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируются БД).**
*   **`Like` (Лайк - PostgreSQL или Redis для счетчиков)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `target_entity_type` (ENUM: `review`, `comment`, `post`, `feed_item`). **Обязательность: Да.**
    *   `target_entity_id` (UUID). **Обязательность: Да.**
    *   `user_id` (UUID, FK). **Обязательность: Да.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   UNIQUE (`target_entity_type`, `target_entity_id`, `user_id`).
*   **`Forum` (Форум - PostgreSQL)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `name` (VARCHAR, UK). **Обязательность: Да.**
    *   `description` (TEXT, Nullable). **Обязательность: Нет.**
    *   `is_locked` (BOOLEAN, Default: false). **Обязательность: Да.**
*   **`ForumTopic` (Тема Форума - PostgreSQL)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `forum_id` (UUID, FK). **Обязательность: Да.**
    *   `title` (VARCHAR). **Обязательность: Да.**
    *   `author_user_id` (UUID, FK). **Обязательность: Да.**
    *   `is_locked` (BOOLEAN, Default: false). **Обязательность: Да.**
    *   `is_pinned` (BOOLEAN, Default: false). **Обязательность: Да.**
    *   `post_count` (INTEGER, Default: 1). **Обязательность: Да.**
    *   `last_post_timestamp` (TIMESTAMPTZ). **Обязательность: Да.**
    *   `created_at`, `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируются БД).**
*   **`ForumPost` (Пост на Форуме - PostgreSQL)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `topic_id` (UUID, FK). **Обязательность: Да.**
    *   `author_user_id` (UUID, FK). **Обязательность: Да.**
    *   `content_markdown` (TEXT). **Обязательность: Да.**
    *   `status` (ENUM: `published`, `pending_moderation`, `rejected`, `hidden`). **Обязательность: Да.**
    *   `created_at`, `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируются БД).**

#### 4.2.2. Apache Cassandra (Чаты, Ленты Активности)
*   **Концептуальная диаграмма `chat_messages`:**
    ```mermaid
    graph TD
        subgraph ChatMessagesTable [Table: chat_messages (Cassandra)]
            direction LR
            PKey["Partition Key: (chat_room_id)"]
            CKey["Clustering Key: (message_id DESC)"]

            PKey --> CKey;
            CKey --> sender_id["sender_id (uuid)"]
            CKey --> sender_nickname["sender_nickname (text)"]
            CKey --> content_text["content_text (text)"]
            CKey --> attachments["attachments (list<text>)"]
            CKey --> created_at["created_at (timestamp)"]
            CKey --> is_edited["is_edited (boolean)"]
        end
        note["Optimized for querying messages per chat room, sorted by time."]
    ```
*   **Концептуальная диаграмма `user_activity_feed`:**
    ```mermaid
    graph TD
        subgraph UserActivityFeedTable [Table: user_activity_feed (Cassandra)]
            direction LR
            PKey["Partition Key: (user_id)"]
            CKey["Clustering Key: (event_time DESC)"]

            PKey --> CKey;
            CKey --> actor_id["actor_id (uuid)"]
            CKey --> actor_nickname["actor_nickname (text)"]
            CKey --> verb["verb (text)"]
            CKey --> object_id["object_id (uuid, optional)"]
            CKey --> object_type["object_type (text, optional)"]
            CKey --> object_preview["object_preview (text, optional)"]
            CKey --> target_id["target_id (uuid, optional)"]
            CKey --> target_type["target_type (text, optional)"]
            CKey --> target_preview["target_preview (text, optional)"]
        end
        note["Fan-out-on-write: User's feed is pre-computed and stored per user."]
    ```

#### 4.2.3. Neo4j (Социальный Граф)
*   **Концептуальная диаграмма графа:**
    ```mermaid
    graph LR
        User1["(:User {userId, nickname})"]
        User2["(:User {userId, nickname})"]
        User3["(:User {userId, nickname})"]
        Game1["(:Game {productId, title})"]
        Group1["(:Group {groupId, name})"]

        User1 -- FRIENDS_WITH <br/> "{since, status}" --> User2
        User1 -- SENT_FRIEND_REQUEST_TO <br/> "{requested_at}" --> User3
        User1 -- MEMBER_OF <br/> "{role, joined_at}" --> Group1
        User1 -- PLAYED <br/> "{last_played, total_playtime}" --> Game1
        User1 -- REVIEWED <br/> "{rating, review_id}" --> Game1
        User1 -- HAS_INTEREST_IN <br/> "{interest_level}" --> Game1
        User2 -- MEMBER_OF <br/> "{role, joined_at}" --> Group1

        style User1 fill:#blue,color:white
        style User2 fill:#blue,color:white
        style User3 fill:#blue,color:white
        style Game1 fill:#green,color:white
        style Group1 fill:#orange,color:white
    ```
*   **Примеры запросов Cypher:**
    *   Найти друзей пользователя: `MATCH (u:User {userId: $userId})-[:FRIENDS_WITH]-(friend:User) RETURN friend.userId, friend.nickname`
    *   Найти друзей друзей (рекомендации): `MATCH (u:User {userId: $userId})-[:FRIENDS_WITH]-(friend:User)-[:FRIENDS_WITH]-(fof:User) WHERE NOT (u)-[:FRIENDS_WITH]-(fof) AND u <> fof RETURN DISTINCT fof.userId, fof.nickname LIMIT 20`
    *   Рекомендовать игры, в которые играют друзья: `MATCH (u:User {userId: $userId})-[:FRIENDS_WITH]-(friend:User)-[:PLAYED]->(game:Game) WHERE NOT (u)-[:PLAYED]->(game) RETURN DISTINCT game.productId, game.title, count(friend) as friends_who_played ORDER BY friends_who_played DESC LIMIT 10`
    *   Найти группы, в которых состоят друзья пользователя: `MATCH (u:User {userId: $userId})-[:FRIENDS_WITH]-(friend:User)-[:MEMBER_OF]->(group:Group) WHERE NOT (u)-[:MEMBER_OF]->(group) RETURN DISTINCT group.groupId, group.name, count(friend) as friends_in_group ORDER BY friends_in_group DESC LIMIT 10`

## 5. Потоковая Обработка Событий (Event Streaming)
## 5. Потоковая Обработка Событий (Event Streaming)
*   **Формат событий:** CloudEvents v1.0 JSON (согласно `../../../../project_api_standards.md`).
*   **Основной топик для публикуемых событий:** `com.platform.social.events.v1`.

### 5.1. Публикуемые События (Produced Events)
*   **`com.platform.social.friend.request.sent.v1`**
    *   Описание: Отправлен запрос на добавление в друзья.
    *   `data` Payload: `{"requesterUserId": "uuid-user-A", "targetUserId": "uuid-user-B", "requestTimestamp": "ISO8601"}`
*   **`com.platform.social.friend.request.accepted.v1`**
    *   Описание: Запрос на добавление в друзья принят.
    *   `data` Payload: `{"accepterUserId": "uuid-user-B", "requesterUserId": "uuid-user-A", "acceptedTimestamp": "ISO8601"}`
*   **`com.platform.social.chat.message.sent.v1`**
    *   Описание: Отправлено новое сообщение в чате (личном или групповом).
    *   `data` Payload: `{"messageId": "uuid-msg", "chatRoomId": "uuid-room", "senderUserId": "uuid-user-A", "receiverUserId": "uuid-user-B", "groupId": null, "text": "Привет!", "sentTimestamp": "ISO8601"}`
*   **`com.platform.social.review.submitted.v1`**
    *   Описание: Пользователь оставил отзыв об игре.
    *   `data` Payload: `{"reviewId": "uuid-review", "userId": "uuid-user", "productId": "uuid-game", "rating": 5, "submissionTimestamp": "ISO8601"}`
*   **`com.platform.social.comment.posted.v1`**
    *   Описание: Оставлен комментарий к сущности (отзыв, пост, и т.д.).
    *   `data` Payload: `{"commentId": "uuid-comment", "userId": "uuid-user", "parentEntityType": "review", "parentEntityId": "uuid-review", "text": "Согласен!", "postedTimestamp": "ISO8601"}`
*   **`com.platform.social.user_activity.event.v1`**
    *   Описание: Общее событие активности пользователя для формирования ленты.
    *   `data` Payload: `{"userId": "uuid-user-B", "activityType": "NEW_ACHIEVEMENT", "details": {"achievementName": "Мастер Клинка", "gameName": "Супер РПГ"}, "timestamp": "ISO8601"}`
    *   (Примеры `activityType`: `NEW_POST_IN_GROUP`, `JOINED_GROUP`, `NEW_REVIEW`, `NEW_FRIEND`)

### 5.2. Потребляемые События (Consumed Events)
*   **`com.platform.account.created.v1`** (от Account Service)
    *   Описание: Создан новый аккаунт пользователя.
    *   Логика обработки: Создать `UserProfile` для нового пользователя.
*   **`com.platform.account.profile.updated.v1`** (от Account Service)
    *   Описание: Базовый профиль пользователя обновлен (например, никнейм, аватар).
    *   Логика обработки: Обновить соответствующие поля в `UserProfile` Social Service.
*   **`com.platform.library.achievement.unlocked.v1`** (от Library Service)
    *   Описание: Пользователь разблокировал достижение.
    *   Логика обработки: Создать событие для ленты активности.
*   **`com.platform.catalog.game.published.v1`** (от Catalog Service)
    *   Описание: Новая игра опубликована.
    *   Логика обработки: Может использоваться для создания автоматических постов на форуме игры или в группах, связанных с игрой.

## 6. Интеграции (Integrations)
(Как в существующем документе).

## 7. Конфигурация (Configuration)
(Как в существующем документе).

## 8. Обработка Ошибок (Error Handling)
(Как в существующем документе, с добавлением специфичных кодов).
*   **`PROFILE_UPDATE_FORBIDDEN`**: Попытка обновить чужой профиль.
*   **`FRIEND_REQUEST_INVALID_TARGET`**: Нельзя отправить запрос самому себе.
*   **`CHAT_ROOM_ACCESS_DENIED`**: Пользователь не является участником чата.
*   **`GROUP_JOIN_POLICY_VIOLATION`**: Нарушение правил вступления в группу (например, закрытая группа).
*   **`FORUM_TOPIC_LOCKED`**: Попытка ответа в закрытой теме форума.

## 9. Безопасность (Security)
(Как в существующем документе, с акцентом на модерацию пользовательских текстовых материалов, приватность, защиту от спама/харассмента).
*   **ФЗ-152 "О персональных данных":** Social Service обрабатывает значительные объемы ПДн (профили, сообщения, списки друзей, контент). Все данные российских пользователей должны храниться и обрабатываться на территории РФ. Настройки приватности должны строго соблюдаться. Содержимое чатов и другой пользовательский контент должен быть защищен от несанкционированного доступа.
*   **Модерация:** Интеграция с Admin Service для модерации пользовательского контента. [Стратегия модерации контента (например, автоматическая с использованием AI + ручная)].

## 10. Развертывание (Deployment)
(Как в существующем документе).

## 11. Мониторинг и Логирование (Logging and Monitoring)
(Как в существующем документе).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность API:** P95 < 200 мс для большинства запросов. P99 < 500 мс.
*   **Лента активности:** P95 < 300 мс на генерацию/отображение.
*   **Чаты:** Доставка сообщений P99 < 1 секунды (в идеале < 500 мс).
*   **Масштабируемость:** Горизонтальное масштабирование для поддержки >1 млн активных пользователей в социальных функциях. Поддержка до [Уточнить значение, например, 1000] друзей на пользователя (`{{MAX_FRIENDS_PER_USER}}` заменено).
*   **Надежность:** Доступность > 99.9%.
*   **Согласованность данных:** Eventual consistency для лент активности и счетчиков. Strong consistency для операций с дружбой и членством в группах (в рамках Neo4j/PostgreSQL).
*   **Scalability (Стратегии масштабирования):**
    *   Использование специализированных БД (Cassandra для чатов/лент, Neo4j для графа).
    *   Кэширование (Redis) для профилей, онлайн-статусов, списков друзей.
    *   Асинхронная обработка событий (Kafka) для генерации лент, уведомлений.
    *   Горизонтальное масштабирование stateless компонентов сервиса.
    *   Потенциальное использование CQRS для разделения моделей чтения и записи для высоконагруженных компонентов (например, лента).

## 13. Приложения (Appendices)

## 13. Приложения (Appendices)
(Как в существующем документе).

## 14. Пользовательские Сценарии (User Flows)

В этом разделе описаны ключевые пользовательские сценарии, связанные с Social Service.

### 14.1. Отправка и Принятие Запроса на Дружбу
*   **Описание:** Пользователь А отправляет запрос на добавление в друзья пользователю Б. Пользователь Б принимает запрос.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor UserA
        participant ClientA as Клиент UserA
        participant APIGW as API Gateway
        participant SocialSvc as Social Service
        participant Neo4jDB as Neo4j
        participant Kafka as Kafka
        participant NotificationSvc as Notification Service
        actor UserB
        participant ClientB as Клиент UserB

        UserA->>ClientA: Найти UserB, нажать "Добавить в друзья"
        ClientA->>APIGW: POST /api/v1/social/users/me/friends/requests (target_user_id: userB_id)
        APIGW->>SocialSvc: Forward request (UserA_id, target_user_id: userB_id)
        SocialSvc->>Neo4jDB: Создать отношение (UserA)-[:SENT_FRIEND_REQUEST_TO]->(UserB)
        SocialSvc->>PostgresDB: (Опционально) Запись в FriendshipRequests для истории/уведомлений
        SocialSvc->>Kafka: Publish `com.platform.social.friend.request.sent.v1` (requester:UserA, target:UserB)
        SocialSvc-->>APIGW: HTTP 202 Accepted
        APIGW-->>ClientA: Запрос отправлен

        Kafka-->>NotificationSvc: Consume `friend.request.sent`
        NotificationSvc->>UserB_Client: (Push/WebSocket) Уведомление "UserA хочет добавить вас в друзья"

        UserB->>ClientB: Открывает уведомление/список запросов
        ClientB->>APIGW: PUT /api/v1/social/users/me/friends/requests/{request_id_or_UserA_id} (action: "accept")
        APIGW->>SocialSvc: Forward request (UserB_id, action: "accept", requester_id: UserA_id)
        SocialSvc->>Neo4jDB: Удалить :SENT_FRIEND_REQUEST_TO, Создать (UserA)-[:FRIENDS_WITH]->(UserB) и (UserB)-[:FRIENDS_WITH]->(UserA)
        SocialSvc->>PostgresDB: (Опционально) Обновить статус FriendshipRequest
        SocialSvc->>Kafka: Publish `com.platform.social.friend.request.accepted.v1` (accepter:UserB, requester:UserA)
        SocialSvc-->>APIGW: HTTP 200 OK
        APIGW-->>ClientB: Друг добавлен

        Kafka-->>NotificationSvc: Consume `friend.request.accepted`
        NotificationSvc->>ClientA: (Push/WebSocket) Уведомление "UserB принял ваш запрос в друзья"
        NotificationSvc->>ClientB: (WebSocket, если онлайн) Обновление списка друзей
    ```

### 14.2. Пользователь Создает Группу, Другие Пользователи Вступают
*   **Описание:** Пользователь А создает новую публичную группу. Пользователь Б находит и вступает в эту группу.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor UserA
        participant ClientA as Клиент UserA
        actor UserB
        participant ClientB as Клиент UserB
        participant APIGW as API Gateway
        participant SocialSvc as Social Service
        participant PostgresDB as PostgreSQL
        participant Neo4jDB as Neo4j (опционально для групп)
        participant Kafka as Kafka

        UserA->>ClientA: Заполняет форму создания группы (название, описание, тип: public)
        ClientA->>APIGW: POST /api/v1/social/groups (данные группы)
        APIGW->>SocialSvc: Forward request (UserA_id)
        SocialSvc->>PostgresDB: Создание записи Group (owner_user_id: UserA_id)
        SocialSvc->>PostgresDB: Создание записи GroupMember (group_id, user_id: UserA_id, role: 'owner')
        SocialSvc->>Neo4jDB: (Опционально) Создание узла (:Group) и отношения (UserA)-[:OWNS_GROUP]->(:Group), (UserA)-[:MEMBER_OF]->(:Group)
        SocialSvc->>Kafka: Publish `com.platform.social.group.created.v1`
        SocialSvc-->>APIGW: HTTP 201 Created (данные группы)
        APIGW-->>ClientA: Группа создана

        UserB->>ClientB: Ищет группу "Фанаты Супер Игры"
        ClientB->>APIGW: GET /api/v1/social/groups?search="Фанаты Супер Игры"
        APIGW->>SocialSvc: Forward request
        SocialSvc->>PostgresDB: Поиск групп (или через Elasticsearch если интегрирован)
        SocialSvc-->>APIGW: Список групп
        APIGW-->>ClientB: Результаты поиска

        UserB->>ClientB: Нажимает "Вступить" для группы X
        ClientB->>APIGW: POST /api/v1/social/groups/{groupX_id}/members (user_id: UserB_id)
        APIGW->>SocialSvc: Forward request (UserB_id)
        SocialSvc->>PostgresDB: Проверка типа группы (public)
        SocialSvc->>PostgresDB: Создание записи GroupMember (group_id, user_id: UserB_id, role: 'member')
        SocialSvc->>PostgresDB: Обновление счетчика участников в Group
        SocialSvc->>Neo4jDB: (Опционально) Создание отношения (UserB)-[:MEMBER_OF]->(:Group {groupId: groupX_id})
        SocialSvc->>Kafka: Publish `com.platform.social.group.member.joined.v1`
        SocialSvc-->>APIGW: HTTP 200 OK
        APIGW-->>ClientB: Вы вступили в группу
    ```

### 14.3. Пользователь Отправляет Сообщение в Групповой Чат
*   **Описание:** Пользователь А отправляет сообщение в чат группы, участником которой он является. Другие участники группы (например, Пользователь Б, если он онлайн) получают это сообщение в реальном времени.
*   **Диаграмма:** (См. диаграмму "Sending/Receiving a Chat Message" в разделе 3.3, адаптировать для группового чата)
    ```mermaid
    sequenceDiagram
        actor UserA
        participant ClientA as Клиент UserA
        participant WebSocket_GW as WebSocket Gateway / SocialSvc Hub
        participant ChatAppSvc as Social Service (ChatAppService)
        participant CassandraDB as Cassandra (chat_messages)
        participant Kafka as Kafka
        participant UserB_Client as Клиент UserB (участник группы, онлайн)

        UserA->>ClientA: Вводит сообщение в чат группы G
        ClientA->>WebSocket_GW: WebSocket: {action_type: "chat.message.send", payload: {chat_room_id: "groupG_id", text: "Всем привет в группе G!"}}
        WebSocket_GW->>ChatAppSvc: Обработка сообщения от UserA
        ChatAppSvc->>ChatAppSvc: Проверка членства UserA в группе G
        ChatAppSvc->>CassandraDB: Сохранение сообщения (chat_room_id: groupG_id, sender_id: UserA_id, ...)
        CassandraDB-->>ChatAppSvc: Сообщение сохранено

        ChatAppSvc->>Kafka: Publish `com.platform.social.chat.message.sent.v1` (chat_type: "group", group_id: groupG_id)

        ChatAppSvc->>WebSocket_GW: Разослать сообщение всем участникам группы G (кроме UserA)
        WebSocket_GW-->>UserB_Client: WebSocket: {event_type: "chat.message.new", payload: {message_details}}
        WebSocket_GW-->>UserA_Client: WebSocket: {event_type: "chat.message.sent.ack", payload: {client_message_id, server_message_id, status: "delivered_to_server"}}
    ```

### 14.4. Обновление Ленты Активности Пользователя
*   **Описание:** Пользователь Б (друг Пользователя А) разблокирует достижение. Это событие попадает в ленту Пользователя А.
*   **Диаграмма:** (См. диаграмму "Generating User's Activity Feed" в разделе 2, фокусируясь на асинхронной части)
    ```mermaid
    sequenceDiagram
        participant LibrarySvc as Library Service
        participant Kafka as Kafka Message Bus
        participant SocialSvc_FeedPopulator as Social Service (Feed Populator/Aggregator)
        participant Neo4jDB as Neo4j (для поиска друзей)
        participant CassandraDB as Cassandra (user_activity_feed)
        participant UserA_Client as Клиент UserA (потребитель ленты)

        LibrarySvc->>Kafka: Publish `com.platform.library.achievement.unlocked.v1` (userId: UserB_id, achievementDetails, gameDetails)

        SocialSvc_FeedPopulator->>Kafka: Consume `achievement.unlocked.v1` for UserB
        SocialSvc_FeedPopulator->>Neo4jDB: Найти друзей UserB (например, UserA)
        Neo4jDB-->>SocialSvc_FeedPopulator: Список друзей [UserA_id, UserC_id, ...]

        loop Для каждого друга (например, UserA)
            SocialSvc_FeedPopulator->>SocialSvc_FeedPopulator: Создать FeedItem (actor: UserB, verb: unlocked_achievement, object: achievement, target: game)
            SocialSvc_FeedPopulator->>CassandraDB: Записать FeedItem в user_activity_feed для UserA_id
        end

        UserA_Client->>SocialSvc_FeedPopulator: (Позже, через WebSocket или REST API) Запрос на обновление ленты
        SocialSvc_FeedPopulator->>CassandraDB: Чтение user_activity_feed для UserA_id
        CassandraDB-->>SocialSvc_FeedPopulator: Список FeedItems
        SocialSvc_FeedPopulator-->>UserA_Client: Обновленная лента
    ```

### 14.5. Пользователь Публикует Отзыв, Другой Пользователь Комментирует
*   **Описание:** Пользователь А пишет отзыв на игру. Пользователь Б видит этот отзыв (возможно, в ленте или на странице игры) и оставляет комментарий.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor UserA
        participant ClientA as Клиент UserA
        actor UserB
        participant ClientB as Клиент UserB
        participant APIGW as API Gateway
        participant SocialSvc as Social Service
        participant PostgresDB as PostgreSQL
        participant Kafka as Kafka Message Bus

        UserA->>ClientA: Пишет отзыв на игру X (рейтинг, текст)
        ClientA->>APIGW: POST /api/v1/social/products/{gameX_id}/reviews (payload)
        APIGW->>SocialSvc: Forward request (UserA_id)
        SocialSvc->>PostgresDB: Сохранение Review (status: 'pending_moderation' или 'published' если премодерация не строгая)
        SocialSvc->>Kafka: Publish `com.platform.social.review.submitted.v1`
        SocialSvc-->>APIGW: HTTP 201 Created (данные отзыва)
        APIGW-->>ClientA: Отзыв опубликован/отправлен на модерацию

        UserB->>ClientB: Читает отзыв UserA (на странице игры или в ленте)
        UserB->>ClientB: Пишет комментарий к отзыву
        ClientB->>APIGW: POST /api/v1/social/reviews/{reviewA_id}/comments (text)
        APIGW->>SocialSvc: Forward request (UserB_id)
        SocialSvc->>PostgresDB: Сохранение Comment, привязанного к Review
        SocialSvc->>Kafka: Publish `com.platform.social.comment.posted.v1` (parent_entity_type: 'review', parent_entity_id: reviewA_id)
        SocialSvc-->>APIGW: HTTP 201 Created (данные комментария)
        APIGW-->>ClientB: Комментарий добавлен

        Note over Kafka: Notification Service может уведомить UserA о новом комментарии.
    ```

### 14.6. Пользователь Участвует в Обсуждении на Форуме
*   **Описание:** Пользователь А создает новую тему на форуме. Пользователь Б отвечает в этой теме.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor UserA
        participant ClientA as Клиент UserA
        actor UserB
        participant ClientB as Клиент UserB
        participant APIGW as API Gateway
        participant SocialSvc as Social Service
        participant PostgresDB as PostgreSQL
        participant Kafka as Kafka Message Bus

        UserA->>ClientA: Открывает Форум X, нажимает "Создать тему"
        ClientA->>APIGW: POST /api/v1/social/forums/{forumX_id}/topics (title, initial_post_content)
        APIGW->>SocialSvc: Forward request (UserA_id)
        SocialSvc->>PostgresDB: Создание ForumTopic и первого ForumPost
        SocialSvc->>Kafka: Publish `com.platform.social.forum.topic.created.v1`
        SocialSvc->>Kafka: Publish `com.platform.social.forum.post.created.v1`
        SocialSvc-->>APIGW: HTTP 201 Created (данные темы и первого поста)
        APIGW-->>ClientA: Тема создана

        UserB->>ClientB: Открывает тему UserA
        UserB->>ClientB: Пишет ответ
        ClientB->>APIGW: POST /api/v1/social/topics/{topicA_id}/posts (content)
        APIGW->>SocialSvc: Forward request (UserB_id)
        SocialSvc->>PostgresDB: Создание ForumPost, обновление счетчиков в ForumTopic
        SocialSvc->>Kafka: Publish `com.platform.social.forum.post.created.v1`
        SocialSvc-->>APIGW: HTTP 201 Created (данные поста)
        APIGW-->>ClientB: Ответ опубликован

        Note over Kafka: Notification Service может уведомить UserA и других подписчиков темы о новом ответе.
    ```

## 15. Резервное Копирование и Восстановление (Backup and Recovery)
(Как в существующем документе, с уточнением RPO/RTO для каждой БД).

## 16. Приложения (Appendices)
(Как в существующем документе).

## 17. Связанные Рабочие Процессы (Related Workflows)
(Как в существующем документе, с добавлением плейсхолдеров для новых воркфлоу).
*   [Процесс Модерации Пользовательских Текстовых Материалов](../../../../project_workflows/user_text_content_moderation_flow.md) <!-- [Ссылка на user_text_content_moderation_flow.md - документ в разработке] -->
*   [Процесс Формирования Рекомендаций Друзей и Контента](../../../../project_workflows/social_recommendation_flow.md) <!-- [Ссылка на social_recommendation_flow.md - документ в разработке] -->
*   [Процесс Обработки Жалоб Пользователей](../../../../project_workflows/user_complaint_handling_flow.md) <!-- [Ссылка на user_complaint_handling_flow.md - документ в разработке] -->

---
*Этот документ является основной спецификацией для Social Service и должен поддерживаться в актуальном состоянии.*
