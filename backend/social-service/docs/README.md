# Спецификация Микросервиса: Social Service

**Версия:** 1.0
**Дата последнего обновления:** 2025-05-25

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Social Service.
*   **Роль в общей архитектуре платформы:** Social Service отвечает за все социальные взаимодействия на платформе, включая управление профилями пользователей, списки друзей, группы, чаты, ленты активности, отзывы, комментарии, форумы и другие социальные функции.
*   **Основные бизнес-задачи:** Способствование общению и взаимодействию между пользователями, создание активного сообщества.
*   (Источник: Спецификация Social Service, разделы 1.1, 1.2, 2.1)

### 1.2. Ключевые Функциональности
*   Управление профилями пользователей (расширенная информация, приватность, кастомизация, лента активности, черный список).
*   Управление друзьями (запросы, списки, статусы, поиск, рекомендации).
*   Управление группами (создание, членство, приватность, объявления, обсуждения).
*   Обмен сообщениями (личные и групповые чаты, статусы доставки/прочтения, история, уведомления).
*   Лента активности (события друзей и групп, лайки, комментарии, фильтрация).
*   Отзывы и комментарии (к играм и другому контенту, оценки, модерация).
*   Форумы и обсуждения (создание, управление, подписки, модерация).
*   Модерация пользовательского контента (инструменты, жалобы, автоматическая фильтрация).
*   (Источник: Спецификация Social Service, раздел 2.3)

### 1.3. Основные Технологии
*   **Язык программирования:** Go (предпочтительно) или Java/Kotlin, Python.
*   **Базы данных:** PostgreSQL (профили, группы, форумы, отзывы), Cassandra (чаты, ленты), Neo4j (граф друзей), Redis (кэш, статусы онлайн).
*   **Брокер сообщений:** Kafka или RabbitMQ.
*   **WebSocket:** Для чатов и обновлений в реальном времени (например, Gorilla WebSocket, Socket.IO).
*   **Инфраструктура:** Docker, Kubernetes, Prometheus, Grafana, ELK/Loki, OpenTelemetry/Jaeger.
*   (Источник: Спецификация Social Service, разделы 3.1.2, 3.3, 8)

### 1.4. Термины и Определения (Glossary)
*   **Лента активности:** Хронологический список событий друзей и групп.
*   **Социальный граф:** Структура связей между пользователями (дружба).
*   **Группа:** Объединение пользователей по интересам.
*   **Отзыв:** Мнение пользователя о продукте.
*   (Полный глоссарий см. в Спецификации Social Service, раздел 9)

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Social Service использует многослойную архитектуру, ориентированную на обработку большого количества запросов на чтение и запись, поддержку реального времени.
*   Применяются подходы CQRS и Event Sourcing (для чатов).
*   Компоненты включают управление профилями, друзьями, группами, чатами, лентой, отзывами/комментариями, форумами, модерацией.
*   (Источник: Спецификация Social Service, разделы 3.1, 3.2)

### 2.2. Слои Сервиса
(На основе Спецификации Social Service, раздел 3.1.1)

#### 2.2.1. Presentation Layer (Слой Представления / Транспортный слой)
*   Ответственность: Обработка входящих REST, gRPC и WebSocket запросов.
*   Ключевые компоненты/модули: REST контроллеры, gRPC серверы, WebSocket хендлеры.

#### 2.2.2. Application Layer (Прикладной Слой / Сервисный слой)
*   Ответственность: Реализация бизнес-логики для каждой социальной функции. Обработка команд (запись) и запросов (чтение).
*   Ключевые компоненты/модули: Сервисы для профилей, друзей, групп, чатов, ленты, отзывов, форумов (с разделением на Command/Query handlers).

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности (UserProfile, Friendship, Group, ChatMessage, FeedItem, Review, ForumPost) и их правила.
*   Ключевые компоненты/модули: Entities, Value Objects, Domain Services, Repository Interfaces.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Взаимодействие с PostgreSQL, Cassandra, Neo4j, Redis, Kafka/RabbitMQ. Управление WebSocket соединениями. Интеграция с Notification Service.
*   Ключевые компоненты/модули: Repositories (PostgreSQL, Cassandra, Neo4j, Redis), Event producers/consumers, WebSocket management, клиенты для других сервисов.

*   (Подробную структуру проекта см. в Спецификации Social Service, раздел 3.1.2)

## 3. API Endpoints

### 3.1. REST API
*   **Аутентификация:** JWT (через Auth Service).
*   **Основные эндпоинты:**
    *   Профили: `GET /users/{userId}/profile`, `PUT /users/me/profile`.
    *   Друзья: `GET /users/me/friends`, `POST /users/me/friends/requests`, `PUT /users/me/friends/requests/{requestId}`.
    *   Группы: `GET /groups`, `POST /groups`, `GET /groups/{groupId}`.
    *   Лента: `GET /feed`, `POST /feed/{itemId}/like`.
    *   Отзывы: `GET /games/{gameId}/reviews`, `POST /games/{gameId}/reviews`.
*   (Более полный список см. в Спецификации Social Service, раздел 5.2).

### 3.2. gRPC API
*   Используется для внутреннего взаимодействия.
*   **Примеры методов:** `CheckFriendship(userId1, userId2)`, `GetUserProfileSummary(userId)`.
*   (Детали см. в Спецификации Social Service, раздел 5.3).

### 3.3. WebSocket API
*   Используется для чатов и обновлений в реальном времени.
*   **События от сервера:** `new_message`, `message_read`, `user_status_update`, `feed_update`.
*   **Сообщения от клиента:** `send_message`, `mark_message_read`.
*   (Детали см. в Спецификации Social Service, раздел 5.4).

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности
*   **UserProfile**: Расширенный профиль пользователя.
*   **Friendship**: Связь дружбы между пользователями.
*   **Group**: Группа пользователей.
*   **GroupMember**: Членство в группе.
*   **ChatMessage**: Сообщение в чате.
*   **FeedItem**: Элемент ленты активности.
*   **Review**: Отзыв на игру.
*   **Comment**: Комментарий.
*   **Forum**: Форум.
*   **ForumTopic**: Тема на форуме.
*   **ForumPost**: Сообщение на форуме.
*   (Детали см. в Спецификации Social Service, раздел 5.1).

### 4.2. Схема Базы Данных
*   **PostgreSQL**: Профили, группы, форумы, отзывы, данные модерации.
    *   *Примечание: Пример DDL для PostgreSQL (профили, группы, форумы, отзывы) будет включать таблицы `user_social_profiles` (для расширенных данных профиля, если они не полностью в Account Service), `friendships`, `groups`, `group_members`, `forums`, `forum_topics`, `forum_posts`, `reviews`, `comments`, `user_blocks`, `content_reports`. Детальная DDL будет представлена в отдельных файлах миграций или в приложении.*
*   **Cassandra**: Сообщения чатов, ленты активности (для высокой нагрузки на запись/чтение).
    ```cql
    // Пример таблицы для сообщений чата в Cassandra (концептуальный)
    // CREATE TABLE IF NOT EXISTS chat_messages (
    //   chat_room_id uuid, // Может быть ID пользователя для личного чата или ID группы
    //   message_id timeuuid,
    //   sender_id uuid,
    //   sender_nickname text, // Денормализация для отображения
    //   content text,
    //   attachments list<text>, // Ссылки на медиа
    //   created_at timestamp,
    //   PRIMARY KEY ((chat_room_id), message_id)
    // ) WITH CLUSTERING ORDER BY (message_id DESC);

    // Пример таблицы для ленты активности в Cassandra (концептуальный)
    // CREATE TABLE IF NOT EXISTS user_activity_feed (
    //   user_id uuid, // ID пользователя, чья это лента
    //   event_time timeuuid, // Время события, для сортировки
    //   event_id uuid,
    //   actor_id uuid,
    //   actor_nickname text, // Денормализация
    //   action_type text, // e.g., 'posted_review', 'added_friend', 'achieved_goal'
    //   object_id uuid,
    //   object_type text, // e.g., 'game', 'review', 'user'
    //   object_name text, // Денормализация (название игры, ник друга)
    //   content_preview text,
    //   PRIMARY KEY (user_id, event_time)
    // ) WITH CLUSTERING ORDER BY (event_time DESC);
    // Примечание: Детальные схемы CQL для Cassandra будут определены в соответствующей технической документации.
    ```
*   **Neo4j**: Социальный граф (друзья).
    ```cypher
    // Пример создания узлов и связей в Neo4j (концептуальный)
    // CREATE CONSTRAINT user_unique_id IF NOT EXISTS FOR (u:User) REQUIRE u.userId IS UNIQUE;
    // CREATE (u1:User {userId: 'user-uuid-1', nickname: 'UserOne'}),
    //        (u2:User {userId: 'user-uuid-2', nickname: 'UserTwo'}),
    //        (u3:User {userId: 'user-uuid-3', nickname: 'UserThree'})
    // CREATE (u1)-[:FRIENDS_WITH {since: timestamp(), status: 'accepted'}]->(u2)
    // CREATE (u1)-[:SENT_FRIEND_REQUEST_TO {requested_at: timestamp()}]->(u3)
    // CREATE (g1:Game {gameId: 'game-uuid-1', title: 'Awesome Game'})
    // CREATE (u1)-[:PLAYED {last_played: timestamp(), playtime_hours: 120}]->(g1)
    // Примечание: Детальная модель данных и индексы для Neo4j будут определены в соответствующей технической документации.
    ```
*   **Redis**: Кэш, статусы онлайн, WebSocket соединения.

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Kafka или RabbitMQ.
*   **Формат событий:** CloudEvents JSON (предположительно).
*   **Основные публикуемые события:** `user.profile.updated`, `friend.request.sent`, `friend.request.accepted`, `group.created`, `chat.message.sent`, `review.submitted`, `comment.posted`, `moderation.required`.
*   (Детали см. в Спецификации Social Service, раздел 5.5).

### 5.2. Потребляемые События (Consumed Events)
*   `account.user.created.v1`:
    *   **Источник:** Account Service
    *   **Назначение:** Создание начального социального профиля пользователя.
    *   **Структура Payload (пример):**
        ```json
        {
          "user_id": "uuid_user",
          "username": "user123",
          "email": "user@example.com",
          "registration_date": "ISO8601_timestamp"
        }
        ```
    *   **Логика обработки:** Создать запись `UserProfile` с базовыми значениями.
*   `account.user.profile.updated.v1`:
    *   **Источник:** Account Service
    *   **Назначение:** Обновление данных в Social Service, которые дублируются или зависят от Account Service (например, никнейм).
    *   **Структура Payload (пример):**
        ```json
        {
          "user_id": "uuid_user",
          "updated_fields": {
            "nickname": "NewUserNick",
            "avatar_url": "new_avatar_link"
          },
          "updated_at": "ISO8601_timestamp"
        }
        ```
    *   **Логика обработки:** Обновить соответствующие поля в `UserProfile`.
*   `library.achievement.unlocked.v1`:
    *   **Источник:** Library Service
    *   **Назначение:** Добавление события о получении достижения в ленту активности.
    *   **Структура Payload (пример):**
        ```json
        {
          "user_id": "uuid_user",
          "game_id": "uuid_game",
          "achievement_id": "uuid_achievement",
          "achievement_name": "Мастер Клинка",
          "achievement_icon_url": "url_to_icon",
          "unlocked_at": "ISO8601_timestamp"
        }
        ```
    *   **Логика обработки:** Создать `FeedItem` типа `achievement_unlocked`.
*   `catalog.game.published.v1`:
    *   **Источник:** Catalog Service
    *   **Назначение:** Активация возможности оставлять отзывы и создавать обсуждения для новой игры.
    *   **Структура Payload (пример):**
        ```json
        {
          "game_id": "uuid_game",
          "title": "Новая Опубликованная Игра",
          "release_date": "ISO8601_date"
        }
        ```
    *   **Логика обработки:** Разрешить создание отзывов и форумов для `game_id`.
*   *(Примечание: События типа `library.game.purchased` также могут использоваться для генерации записей в ленте активности).*

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Account Service**: Получение/обновление основной информации профиля.
*   **Auth Service**: Аутентификация/авторизация.
*   **Library Service**: Данные об играх, достижениях для профилей и лент.
*   **Catalog Service**: Метаданные игр для отзывов, обсуждений.
*   **Notification Service**: Отправка уведомлений (новые сообщения, запросы в друзья).
*   **Admin Service**: Модерация контента.
*   **Analytics Service**: Сбор данных о социальной активности.
*   (Детали см. в Спецификации Social Service, разделы 1.3 и 6).

### 6.2. Внешние Системы
*   Information not found in existing documentation.

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `SOCIAL_SERVICE_HTTP_PORT`, `SOCIAL_SERVICE_GRPC_PORT`, `SOCIAL_SERVICE_WS_PORT`.
*   DSN для PostgreSQL, Cassandra, Neo4j.
*   `REDIS_ADDR`.
*   `KAFKA_BROKERS` (comma-separated) / `RABBITMQ_URL`
*   `KAFKA_TOPIC_SOCIAL_EVENTS` (e.g., `social.events`)
*   `ACCOUNT_SERVICE_GRPC_ADDR`
*   `AUTH_SERVICE_GRPC_ADDR`
*   `LIBRARY_SERVICE_GRPC_ADDR`
*   `CATALOG_SERVICE_GRPC_ADDR`
*   `NOTIFICATION_SERVICE_KAFKA_TOPIC` (для отправки уведомлений)
*   `LOG_LEVEL` (e.g., `info`, `debug`)
*   `MAX_CHAT_HISTORY_DAYS` (e.g., `90`)
*   `FEED_ITEMS_PER_PAGE` (e.g., `20`)
*   `MAX_FRIENDS_LIMIT` (e.g., `1000`)
*   `MAX_GROUPS_JOIN_LIMIT` (e.g., `100`)
*   `WEBSOCKET_MAX_CONNECTIONS` (e.g., `10000`)
*   `WEBSOCKET_WRITE_BUFFER_SIZE` (e.g., `1024`)
*   `WEBSOCKET_READ_BUFFER_SIZE` (e.g., `1024`)
*   `WEBSOCKET_PING_INTERVAL_SECONDS` (e.g., `30`)
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`

### 7.2. Файлы Конфигурации (если применимо)
*   Конфигурация сервиса осуществляется преимущественно через переменные окружения. Если потребуются файлы конфигурации для сложных настроек (например, для правил модерации, параметров алгоритмов рекомендаций друзей, или настроек WebSocket), их структура будет определена здесь.

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   Стандартные коды HTTP/gRPC/WebSocket.
*   Подробное логирование.
*   Механизмы retry для межсервисных вызовов.

### 8.2. Распространенные Коды Ошибок
*   `PROFILE_NOT_FOUND`
*   `FRIEND_REQUEST_ALREADY_SENT`
*   `USER_ALREADY_IN_GROUP`
*   `MESSAGE_VALIDATION_FAILED`
*   `CONTENT_MODERATION_FAILED`

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   JWT (через Auth Service) для всех API.

### 9.2. Авторизация
*   RBAC. Проверка прав на создание/редактирование/удаление контента, управление группами, приватность профиля.

### 9.3. Защита Данных
*   Соблюдение настроек приватности пользователя.
*   На данный момент end-to-end шифрование для личных сообщений не планируется в первой версии, будет использоваться шифрование на транспортном уровне (TLS) и шифрование в покое для данных чатов. Вопрос внедрения E2EE может быть рассмотрен в будущем.
*   Защита от спама и вредоносного контента (автоматическая и ручная модерация).

### 9.4. Управление Секретами
*   Kubernetes Secrets или Vault.
*   (Детали см. в Спецификации Social Service, раздел 7 - Безопасность).

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile.**
*   **Kubernetes манифесты/Helm-чарты.**

### 10.2. Зависимости при Развертывании
*   PostgreSQL, Cassandra, Neo4j, Redis, Kafka/RabbitMQ.
*   Account, Auth, Library, Catalog, Notification, Admin Services.

### 10.3. CI/CD
*   Автоматическая сборка, тестирование, развертывание.
*   (Детали см. в Спецификации Social Service, раздел 8).

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   Структурированные JSON логи.
*   Интеграция с ELK/Loki.

### 11.2. Мониторинг
*   Prometheus, Grafana.
*   Метрики: количество активных пользователей, сообщений в чатах, постов в ленте, созданных групп, ошибок API, задержки.

### 11.3. Трассировка
*   OpenTelemetry/Jaeger для отслеживания запросов.
*   (Детали см. в Спецификации Social Service, раздел 8).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:** API чтения P95 < 150мс; API записи P95 < 100мс; доставка сообщений в чатах < 1 сек.
*   **Масштабируемость:** Горизонтальная. Поддержка миллионов активных пользователей.
*   **Надежность:** Доступность 99.95%. Гарантированная доставка сообщений.
*   **Безопасность:** Защита ПД, приватность, защита от спама.
*   **Согласованность данных:** Eventual consistency для лент и счетчиков, strong для остального.
*   (Детали см. в Спецификации Social Service, раздел 2.4).

## 13. Приложения (Appendices) (Опционально)
*   Детальные схемы DDL для PostgreSQL, CQL для Cassandra, и модель данных для Neo4j, а также полные примеры REST API, WebSocket сообщений, и форматы событий Kafka/RabbitMQ будут добавлены по мере финализации дизайна и реализации соответствующих модулей.

---
*Этот шаблон является отправной точкой и может быть адаптирован под конкретные нужды проекта и сервиса.*
