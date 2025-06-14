# Структура Баз Данных Проекта

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16

## 1. Введение

Данный документ предназначен для описания общей структуры баз данных, используемых в проекте "Российский Аналог Платформы Steam". Он должен предоставлять высокоуровневый обзор моделей данных каждого микросервиса и, по возможности, общую схему взаимодействия данных на уровне платформы.

**Цели документа:**

*   Дать общее представление о данных, которыми оперирует каждый микросервис.
*   Описать ключевые сущности и их связи на уровне отдельных сервисов.
*   Если существуют общие модели данных или разделяемые базы данных (не рекомендуется в чистой микросервисной архитектуре, но возможно для некоторых случаев), описать их здесь.
*   Служить отправной точкой для разработчиков при изучении моделей данных конкретных сервисов.

## 2. Общие Принципы Проектирования Баз Данных

*   **Децентрализация данных (Database per Service):** Каждый микросервис владеет своими данными и своей схемой базы данных. Прямой доступ к базе данных другого сервиса строго запрещен; взаимодействие между сервисами должно осуществляться исключительно через их определенные API (REST, gRPC) или через асинхронный обмен событиями (Kafka).
*   **Полиглотное хранение (Polyglot Persistence):** Выбор конкретной технологии СУБД (SQL, NoSQL, графовая, колоночная, in-memory и т.д.) для каждого микросервиса диктуется его специфическими требованиями к данным, характером нагрузки (чтение/запись), требованиями к консистентности, масштабируемости и сложности запросов. Не существует единой СУБД для всех сервисов. Например, PostgreSQL может использоваться для транзакционных данных с сильной консистентностью, Elasticsearch для полнотекстового поиска, Cassandra для сценариев с высокой нагрузкой на запись и большими объемами временных рядов или сообщений, Neo4j для графовых связей, ClickHouse для аналитических запросов (OLAP), а Redis для кэширования и хранения эфемерных данных.
*   **Миграции схемы:** Для реляционных баз данных (например, PostgreSQL) и некоторых NoSQL баз данных, поддерживающих схемы (например, определения таблиц в Cassandra или ClickHouse), должны использоваться инструменты управления миграциями схемы (например, golang-migrate, Flyway, Alembic, или специфичные для СУБД инструменты). Миграции должны быть версионируемыми и интегрированы в процесс CI/CD.
*   **Резервное копирование и восстановление:** Для каждой базы данных должны быть настроены процедуры регулярного резервного копирования и протестированы сценарии восстановления.
*   **Безопасность:** Доступ к базам данных должен быть строго ограничен. Учетные данные для доступа должны храниться в секретах.

## 3. Структура Баз Данных по Микросервисам

Этот раздел должен содержать ссылки на документацию по структуре БД каждого микросервиса или краткое описание их моделей данных.

### 3.1. Auth Service
*   **СУБД:** PostgreSQL (основное хранилище), Redis (кэш сессий, временные токены, JTI blacklist, rate limiting).
*   **Основные сущности:** Пользователи (`users` - учетные данные, статус, 2FA), Роли (`roles`), Разрешения (`permissions`), Сессии (`sessions`), Refresh-токены (`refresh_tokens`), API-ключи (`api_keys`), Коды верификации (`verification_codes`), Логи аудита (`audit_logs`).
*   **Ссылка на детальную документацию:** `backend/auth-service/docs/README.md#4-модели-данных-data-models`

### 3.2. Account Service
*   **СУБД:** PostgreSQL (основное хранилище), Redis (кэширование).
*   **Основные сущности:** Аккаунты пользователей (`accounts` - связь с UserID из Auth), Профили (`profiles` - никнейм, аватар, bio), Контактная информация (`contact_infos`), Настройки пользователя (`user_settings`).
*   **Ссылка на детальную документацию:** `backend/account-service/docs/README.md#4-модели-данных-data-models`

### 3.3. Catalog Service
*   **СУБД:** PostgreSQL (основное хранилище), Elasticsearch (полнотекстовый поиск, фильтрация), Redis (кэширование).
*   **Основные сущности:** Продукты (`products` - Игры, DLC, ПО, Комплекты), Жанры (`genres`), Теги (`tags`), Категории (`categories`), Цены продуктов (`product_prices`), Медиа-контент (`media_items`), Метаданные достижений (`achievement_metadata`).
*   **Ссылка на детальную документацию:** `backend/catalog-service/docs/README.md#4-модели-данных-data-models`

### 3.4. Library Service
*   **СУБД:** PostgreSQL (основное хранилище), Redis (кэширование, текущие игровые сессии), S3-совместимое хранилище (для игровых сохранений).
*   **Основные сущности:** Элементы библиотеки пользователя (`user_library_items`), Игровые сессии (`playtime_sessions`), Достижения пользователя (`user_achievements`), Элементы списка желаемого (`wishlist_items`), Метаданные сохранений игр (`savegame_metadata`).
*   **Ссылка на детальную документацию:** `backend/library-service/docs/README.md#4-модели-данных-data-models`

### 3.5. Payment Service
*   **СУБД:** PostgreSQL (основное хранилище), Redis (кэширование сессий платежей, ключи идемпотентности).
*   **Основные сущности:** Транзакции (`transactions`), Элементы транзакций (`transaction_items`), Платежные методы пользователей (`user_payment_methods`), Фискальные чеки (`fiscal_receipts`), Балансы разработчиков (`developer_balances`), Выплаты разработчикам (`developer_payouts`), Промокоды (`promo_codes`).
*   **Ссылка на детальную документацию:** `backend/payment-service/docs/README.md#4-модели-данных-data-models`

### 3.6. Download Service
*   **СУБД:** PostgreSQL (метаданные загрузок, файлов), Redis (очереди загрузок, текущие статусы, кэш CDN токенов).
*   **Основные сущности:** Сессии загрузки (`download_sessions`), Элементы загрузки (`download_items`), Метаданные файлов продуктов (`file_metadata`), Информация об обновлениях (`updates_history`).
*   **Ссылка на детальную документацию:** `backend/download-service/docs/README.md#4-модели-данных-data-models`

### 3.7. Social Service
*   **СУБД:** PostgreSQL (профили, группы, форумы, отзывы), Apache Cassandra (чаты, ленты активности), Neo4j (социальный граф друзей), Redis (кэш, онлайн-статусы).
*   **Основные сущности:** Расширенные профили пользователей (`user_social_profiles`), Дружеские связи (`friendships` - в Neo4j и частично в PG), Группы (`groups`), Сообщения чатов (`chat_messages` - в Cassandra), Элементы ленты активности (`feed_items` - в Cassandra), Отзывы (`reviews`), Комментарии (`comments`).
*   **Ссылка на детальную документацию:** `backend/social-service/docs/README.md#4-модели-данных-data-models`

### 3.8. Notification Service
*   **СУБД:** PostgreSQL (шаблоны, предпочтения, недавние логи), ClickHouse (долгосрочная статистика и логи доставки), Redis (кэширование, счетчики).
*   **Основные сущности:** Шаблоны уведомлений (`notification_templates`), Логи отправленных уведомлений (`notification_message_log_recent` в PG, `notification_delivery_stats` в ClickHouse), Пользовательские предпочтения (`user_notification_preferences`), Маркетинговые кампании (`notification_campaigns`), Токены устройств (`device_tokens`).
*   **Ссылка на детальную документацию:** `backend/notification-service/docs/README.md#4-модели-данных-data-models`

### 3.9. Developer Service
*   **СУБД:** PostgreSQL (основное хранилище), S3-совместимое хранилище (билды игр, медиа-файлы).
*   **Основные сущности:** Аккаунты разработчиков (`developers`), Члены команд разработчиков (`developer_team_members`), Продукты/Игры (`games`), Версии игр (`game_versions`), Локализованные метаданные игр (`game_metadata_localized`), Запросы на выплаты (`developer_payouts`), API-ключи разработчиков (`developer_api_keys`).
*   **Ссылка на детальную документацию:** `backend/developer-service/docs/README.md#4-модели-данных-data-models`

### 3.10. Admin Service
*   **СУБД:** PostgreSQL (структурированные данные), MongoDB (логи аудита, сложные данные модерации).
*   **Основные сущности:** Административные пользователи (`admin_users`), Элементы модерации (`moderation_items`), Правила модерации (`moderation_rules`), Тикеты поддержки (`support_tickets`), Ответы на тикеты (`support_ticket_responses`), Статьи базы знаний (`knowledge_base_articles`), Настройки платформы (`platform_settings`), Логи аудита администраторов (`audit_log_admin` - в MongoDB или PG).
*   **Ссылка на детальную документацию:** `backend/admin-service/docs/README.md#4-модели-данных-data-models`

### 3.11. Analytics Service
*   **СУБД:** ClickHouse (DWH - хранилище данных для аналитики), PostgreSQL (хранение метаданных: определения метрик, отчетов, конфигурации пайплайнов). S3 Data Lake для сырых событий.
*   **Основные сущности (в метаданных PG):** Определения метрик (`metric_definitions`), Определения отчетов (`report_definitions`), Определения сегментов (`segment_definitions`), Метаданные ML моделей (`ml_model_metadata`), Запуски пайплайнов (`data_pipeline_runs`). В ClickHouse хранятся таблицы фактов и витрины данных.
*   **Ссылка на детальную документацию:** `backend/analytics-service/docs/README.md#4-модели-данных-data-models`

### 3.12. API Gateway
*   **СУБД:** Обычно не использует собственную выделенную базу данных для бизнес-сущностей. Конфигурация маршрутов, плагинов и потребителей управляется через CRD Kubernetes или, в некоторых решениях (например, Kong), может использовать PostgreSQL или Cassandra для хранения своей конфигурации. Redis может использоваться для кэширования и rate limiting.
*   **Основные сущности (конфигурационные):** Маршруты, Сервисы, Потребители, Плагины.
*   **Ссылка на детальную документацию:** `backend/api-gateway/docs/README.md#4-модели-данных-data-models`

## 4. Общая Схема Данных Платформы (ERD)

```mermaid
erDiagram
    USER ||--o{ ACCOUNT_PROFILE : "has (Account Svc)"
    USER ||--o{ USER_ROLE : "has (Auth Svc)"
    USER_ROLE ||--|{ ROLE : "is (Auth Svc)"
    USER ||--o{ SESSION : "has active (Auth Svc)"
    USER ||--o{ LIBRARY_ITEM : "owns (Library Svc)"
    USER ||--o{ ORDER : "places (Payment Svc)"
    USER ||--o{ PAYMENT_METHOD : "uses (Payment Svc)"
    USER ||--o{ FRIENDSHIP : "initiates/receives (Social Svc)"
    USER ||--o{ CHAT_MESSAGE : "sends/receives (Social Svc)"
    USER ||--o{ REVIEW : "writes (Social Svc)"
    USER ||--o{ WISHLIST_ITEM : "adds to (Library Svc)"
    USER ||--o{ SUPPORT_TICKET : "creates (Admin Svc)"
    USER ||--o{ NOTIFICATION_PREFERENCE : "sets (Notification Svc)"
    USER {
        string user_id PK "UUID (Auth Svc)"
        string email "varchar (Auth Svc)"
        string username "varchar (Auth Svc)"
        string password_hash "varchar (Auth Svc)"
        datetime registration_date "timestamp (Auth Svc)"
    }

    ACCOUNT_PROFILE {
        string profile_id PK "UUID (Account Svc)"
        string user_id FK "UUID (Account Svc)"
        string display_name "varchar (Account Svc)"
        string avatar_url "varchar (Account Svc)"
    }

    DEVELOPER_ACCOUNT ||--o{ GAME : "develops (Developer Svc)"
    PUBLISHER_ACCOUNT ||--o{ GAME : "publishes (Developer Svc)"
    DEVELOPER_ACCOUNT {
        string developer_id PK "UUID (Developer Svc)"
        string user_id FK "UUID (links to USER)"
        string company_name "varchar (Developer Svc)"
    }
    PUBLISHER_ACCOUNT {
        string publisher_id PK "UUID (Developer Svc)"
        string company_name "varchar (Developer Svc)"
    }

    GAME ||--|{ GAME_METADATA : "has (Catalog Svc)"
    GAME ||--o{ GAME_VERSION : "has (Developer/Catalog Svc)"
    GAME ||--o{ DLC : "can have (Catalog Svc)"
    GAME ||--o{ GAME_PRICE : "has (Catalog Svc)"
    GAME ||--o{ REVIEW : "is for (Social Svc)"
    GAME ||--o{ ACHIEVEMENT_METADATA : "defines (Catalog Svc)"
    GAME ||--o{ GAME_GENRE_ASSOC : "belongs to (Catalog Svc)"
    GAME ||--o{ GAME_TAG_ASSOC : "has (Catalog Svc)"
    GAME {
        string game_id PK "UUID (Catalog Svc)"
        string title "varchar (Catalog Svc)"
        string developer_id FK "UUID (Catalog Svc)"
        string publisher_id FK "UUID (Catalog Svc)"
        datetime release_date "timestamp (Catalog Svc)"
    }
    GAME_METADATA {
        string game_id FK "UUID (Catalog Svc)"
        string description "text (Catalog Svc)"
        string short_description "varchar (Catalog Svc)"
        json system_requirements "jsonb (Catalog Svc)"
    }
    DLC {
        string dlc_id PK "UUID (Catalog Svc)"
        string game_id FK "UUID (Catalog Svc)"
        string title "varchar (Catalog Svc)"
    }
    BUNDLE ||--o{ GAME_BUNDLE_ITEM : "contains (Catalog Svc)"
    GAME_BUNDLE_ITEM }|--|| GAME : "item (Catalog Svc)"
    GAME_BUNDLE_ITEM }|--|| DLC : "item (Catalog Svc)"

    LIBRARY_ITEM {
        string library_item_id PK "UUID (Library Svc)"
        string user_id FK "UUID (Library Svc)"
        string game_id FK "UUID (Library Svc)"
        datetime purchase_date "timestamp (Library Svc)"
        int total_playtime_minutes "integer (Library Svc)"
    }
    WISHLIST_ITEM {
        string wishlist_item_id PK "UUID (Library Svc)"
        string user_id FK "UUID (Library Svc)"
        string game_id FK "UUID (Library Svc)"
        datetime added_date "timestamp (Library Svc)"
    }

    ORDER ||--o{ ORDER_ITEM : "contains (Payment Svc)"
    ORDER ||--|{ PAYMENT_TRANSACTION : "results in (Payment Svc)"
    ORDER {
        string order_id PK "UUID (Payment Svc)"
        string user_id FK "UUID (Payment Svc)"
        datetime order_date "timestamp (Payment Svc)"
        decimal total_amount "decimal (Payment Svc)"
        string status "varchar (Payment Svc)"
    }
    ORDER_ITEM {
        string order_item_id PK "UUID (Payment Svc)"
        string order_id FK "UUID (Payment Svc)"
        string product_id FK "UUID (game_id or dlc_id, Catalog Svc)"
        string product_type "varchar (Payment Svc)"
        decimal price "decimal (Payment Svc)"
    }
    PAYMENT_TRANSACTION ||--o{ FISCAL_RECEIPT : "generates (Payment Svc)"
    PAYMENT_TRANSACTION {
        string transaction_id PK "UUID (Payment Svc)"
        string order_id FK "UUID (Payment Svc)"
        string payment_method_id FK "UUID (Payment Svc)"
        string status "varchar (Payment Svc)"
        datetime transaction_date "timestamp (Payment Svc)"
    }

    GAME_VERSION ||--o{ DOWNLOAD_FILE : "consists of (Download Svc)"
    GAME_VERSION {
        string version_id PK "UUID (Developer/Catalog Svc)"
        string game_id FK "UUID (Developer/Catalog Svc)"
        string version_number "varchar (Developer/Catalog Svc)"
    }

    GENRE {
        string genre_id PK "UUID (Catalog Svc)"
        string name "varchar (Catalog Svc)"
    }
    GAME_GENRE_ASSOC {
        string game_id FK "UUID (Catalog Svc)"
        string genre_id FK "UUID (Catalog Svc)"
    }
    GAME }|--o| GAME_GENRE_ASSOC : "association"
    GENRE ||--o| GAME_GENRE_ASSOC : "association"

    TAG {
        string tag_id PK "UUID (Catalog Svc)"
        string name "varchar (Catalog Svc)"
    }
    GAME_TAG_ASSOC {
        string game_id FK "UUID (Catalog Svc)"
        string tag_id FK "UUID (Catalog Svc)"
    }
    GAME }|--o| GAME_TAG_ASSOC : "association"
    TAG ||--o| GAME_TAG_ASSOC : "association"

    NOTIFICATION ||--|{ USER : "targets (Notification Svc)"
    NOTIFICATION {
        string notification_id PK "UUID (Notification Svc)"
        string user_id FK "UUID (Notification Svc)"
        string content "text (Notification Svc)"
        datetime sent_at "timestamp (Notification Svc)"
    }

    ADMIN_USER ||--o{ MODERATION_ACTION : "performs (Admin Svc)"
    MODERATION_ACTION ||--|{ REVIEW : "applies to (Admin Svc)"
    MODERATION_ACTION ||--|{ GAME : "applies to (Admin Svc)"

    ANALYTICS_EVENT ||--|{ USER : "relates to (Analytics Svc)"
    ANALYTICS_EVENT ||--|{ GAME : "relates to (Analytics Svc)"

    API_GATEWAY_ROUTE --> MICROSERVICE_API : "routes to"
```
*(Примечание: Это логическая диаграмма, показывающая ключевые сущности и их предполагаемое размещение по сервисам. Реальные схемы БД каждого сервиса будут детализированы в их собственной документации и могут содержать дополнительные таблицы и связи для оптимизации и нормализации).*

## 5. Вопросы Синхронизации и Консистентности Данных

В микросервисной архитектуре, где каждый сервис владеет своими данными, обеспечение консистентности данных между сервисами является ключевой задачей. Вместо традиционных ACID-транзакций, охватывающих несколько баз данных (что противоречит принципам слабой связанности микросервисов), применяются другие подходы.

### 5.1. Eventual Consistency (Согласованность в конечном счете)

**Eventual Consistency** является преобладающим подходом в нашей платформе. Этот подход означает, что система достигнет согласованного состояния через некоторое время после выполнения операции, но не обязательно немедленно.

*   **Принцип работы:**
    *   Когда сервис изменяет свои данные (например, `Auth Service` регистрирует нового пользователя), он публикует **доменное событие** (например, `UserRegisteredEvent`) в централизованную очередь сообщений (например, Apache Kafka).
    *   Другие микросервисы, заинтересованные в этом событии (например, `Account Service` для создания профиля, `Notification Service` для отправки приветственного письма), подписываются на эти события.
    *   При получении события каждый сервис-подписчик асинхронно обновляет свои локальные данные или выполняет необходимые действия. Например, `Account Service` может создать локальную копию некоторых данных пользователя или просто ссылку на `user_id`.

*   **Преимущества:**
    *   **Низкая связанность (Decoupling):** Сервисы не зависят напрямую от доступности других сервисов для обновления данных.
    *   **Отказоустойчивость (Resilience):** Если один сервис временно недоступен, другие сервисы могут продолжать обрабатывать события после его восстановления.
    *   **Масштабируемость (Scalability):** Асинхронная обработка событий позволяет лучше распределять нагрузку.
    *   **Производительность:** Операции в исходном сервисе завершаются быстрее, так как не ждут подтверждения от всех зависимых сервисов.

*   **Недостатки:**
    *   **Временная рассинхронизация данных:** Существует задержка между изменением данных в одном сервисе и их отражением в других. Пользовательский интерфейс должен быть спроектирован с учетом этого (например, не отображать данные, которые могут быть еще не консистентны, или указывать на их возможное обновление).
    *   **Сложность отладки и тестирования:** Отслеживание потока данных через асинхронные события может быть сложнее.
    *   **Обработка дубликатов и порядка событий:** Требуются механизмы для идемпотентной обработки событий и, в некоторых случаях, для сохранения порядка событий.

*   **Пример:**
    1.  `Payment Service` успешно обрабатывает платеж и публикует событие `OrderPaidEvent`, содержащее `order_id` и `user_id`.
    2.  `Library Service` подписан на `OrderPaidEvent`. При получении события он добавляет соответствующие игры из заказа в библиотеку пользователя.
    3.  `Notification Service` также подписан на `OrderPaidEvent` и отправляет пользователю уведомление об успешной покупке.
    Между шагом 1 и шагами 2-3 существует небольшая задержка, в течение которой пользователь может не сразу увидеть игру в библиотеке.

### 5.2. Saga Pattern (Паттерн Сага)

Для бизнес-процессов, которые требуют координации нескольких сервисов и должны либо полностью завершиться успешно, либо полностью отмениться (имитируя атомарность распределенной транзакции), используется паттерн **Сага**. Сага представляет собой последовательность локальных транзакций в каждом участвующем сервисе. Если какая-либо локальная транзакция завершается неудачей, сага запускает компенсирующие транзакции для отмены предыдущих успешно выполненных локальных транзакций.

*   **Типы реализации Саги:**
    *   **Хореография (Choreography-based Saga):**
        *   Каждый сервис, выполнив свою локальную транзакцию, публикует событие.
        *   Другие сервисы слушают эти события и выполняют свои локальные транзакции.
        *   **Преимущества:** Простота реализации для небольшого числа участников, отсутствие единой точки отказа.
        *   **Недостатки:** Сложность отслеживания состояния саги, риск циклических зависимостей, затрудненная отладка при большом количестве сервисов.
    *   **Оркестрация (Orchestration-based Saga):**
        *   Центральный координатор (оркестратор саги) управляет всем процессом.
        *   Оркестратор отправляет команды каждому сервису для выполнения локальной транзакции и ожидает ответа.
        *   В случае сбоя оркестратор отвечает за запуск компенсирующих транзакций в правильном порядке.
        *   **Преимущества:** Централизованное управление логикой саги, явное определение состояния, упрощенное добавление новых шагов, лучшая наблюдаемость.
        *   **Недостатки:** Риск превращения оркестратора в "божественный объект", дополнительный компонент для разработки и поддержки.

*   **Компенсирующие транзакции:** Для каждой транзакции в саге, которая может завершиться неудачей, должна быть предусмотрена компенсирующая транзакция. Компенсирующая транзакция должна быть идемпотентной и гарантированно выполнимой (насколько это возможно).

*   **Пример (Покупка игры - Оркестрованная Сага):**
    Оркестратор: `OrderSagaOrchestrator`
    1.  **Клиент -> `Order Service`:** Создать заказ (статус "Pending").
        *   `Order Service` -> Оркестратор: `OrderCreated` (сообщает об успешном создании заказа).
    2.  **Оркестратор -> `Payment Service`:** Обработать платеж для `order_id`.
        *   `Payment Service`: Выполняет локальную транзакцию по списанию средств.
        *   `Payment Service` -> Оркестратор: `PaymentProcessed` (успех) или `PaymentFailed` (неудача).
    3.  **Если `PaymentFailed`:**
        *   Оркестратор -> `Order Service`: Отменить заказ `order_id` (компенсирующая транзакция).
        *   Сага завершается неудачей.
    4.  **Если `PaymentProcessed`:**
        *   Оркестратор -> `Library Service`: Добавить игры из `order_id` в библиотеку `user_id`.
        *   `Library Service`: Выполняет локальную транзакцию по добавлению игр.
        *   `Library Service` -> Оркестратор: `GamesAddedToLibrary` (успех) или `AddToLibraryFailed` (неудача).
    5.  **Если `AddToLibraryFailed`:**
        *   Оркестратор -> `Payment Service`: Вернуть платеж для `order_id` (компенсирующая транзакция).
        *   Оркестратор -> `Order Service`: Отменить заказ `order_id` (компенсирующая транзакция).
        *   Сага завершается неудачей.
    6.  **Если `GamesAddedToLibrary`:**
        *   Оркестратор -> `Notification Service`: Отправить уведомление об успешной покупке.
        *   Оркестратор -> `Order Service`: Установить статус заказа `order_id` как "Completed".
        *   Сага успешно завершена.

Выбор между хореографией и оркестрацией зависит от сложности процесса. Для простых саг с 2-3 участниками может подойти хореография, для более сложных и длительных процессов предпочтительнее оркестрация.

Эти подходы помогают поддерживать целостность данных на уровне бизнес-транзакций в распределенной системе, где традиционные ACID-транзакции неприменимы или нежелательны.

## 6. Стратегия Резервного Копирования и Восстановления (Backup and Recovery Strategy)

### 6.1. Общие Принципы Резервного Копирования
*   Все критически важные данные, хранящиеся в базах данных микросервисов (PostgreSQL, MongoDB, Cassandra, Neo4j, ClickHouse, Elasticsearch), подлежат регулярному резервному копированию.
*   Резервные копии должны храниться в безопасном, географически удаленном месте от основных серверов данных для предотвращения потери данных в случае катастрофы на основной площадке.
*   Процедуры восстановления из резервных копий должны быть документированы и регулярно тестироваться (не реже одного раза в квартал) для каждого типа СУБД.
*   Необходимо обеспечить мониторинг процессов резервного копирования и алертинг в случае сбоев.

### 6.2. Стратегии по Типам СУБД и Сервисам
*   **PostgreSQL (для Auth, Account, Catalog, Library, Payment, Developer, Admin Services):**
    *   **Частота:** Ежедневные полные бэкапы. Инкрементальные бэкапы или WAL (Write-Ahead Logging) archiving для возможности восстановления на определенный момент времени (Point-in-Time Recovery - PITR).
    *   **RPO (Recovery Point Objective):** Не более 15-30 минут (достигается за счет WAL).
    *   **RTO (Recovery Time Objective):** 2-4 часа для полного восстановления сервиса из бэкапа.
    *   **Хранение:** Полные бэкапы хранятся 30 дней, WAL-архивы - 7 дней. Еженедельные бэкапы архивируются на 3 месяца, ежемесячные - на 1 год.
    *   **Инструменты:** `pg_dumpall` для полных логических бэкапов, `pg_basebackup` для физических, инструменты для управления WAL-архивированием (например, `pgBackRest`, `wal-g`).
*   **Redis (для кэшей, сессий):**
    *   **Частота:** Ежедневные снимки (RDB snapshots). AOF (Append Only File) для повышения долговечности (если используется для данных, которые не могут быть легко восстановлены из основного хранилища).
    *   **RPO:** Для кэша - потеря данных за последние 24 часа допустима. Если Redis используется как основное хранилище для некоторых данных, RPO должен быть ниже.
    *   **RTO:** 1-2 часа.
    *   **Хранение:** Снимки хранятся 7 дней.
    *   **Примечание:** "Для большинства сценариев использования Redis (кэширование, временные данные), данные могут быть восстановлены из основных источников (БД) или пересозданы приложением. Бэкапы важны для быстрого восстановления состояния или для случаев, когда Redis содержит уникальные данные."
*   **Elasticsearch (для Catalog, Admin Services):**
    *   **Частота:** Ежедневные снимки (snapshots) индексов.
    *   **RPO:** 24 часа.
    *   **RTO:** 4-8 часов (в зависимости от объема данных).
    *   **Хранение:** Снимки хранятся 14 дней.
    *   **Инструменты:** Встроенный механизм Elasticsearch Snapshot and Restore.
*   **ClickHouse (для Analytics, Notification Services):**
    *   **Частота:** Ежедневные бэкапы.
    *   **RPO:** 24 часа.
    *   **RTO:** 4-12 часов (в зависимости от объема данных).
    *   **Хранение:** Бэкапы хранятся 30 дней.
    *   **Инструменты:** `clickhouse-backup` или аналогичные.
*   **Cassandra (для Social Service - чаты, ленты):**
    *   **Частота:** Ежедневные снимки (snapshots) на каждом узле. Регулярное инкрементальное резервное копирование.
    *   **RPO:** 24 часа (для снимков), меньше для инкрементальных.
    *   **RTO:** 6-12 часов.
    *   **Хранение:** Снимки хранятся 7-14 дней.
    *   **Инструменты:** `nodetool snapshot`, кастомные скрипты или специализированные решения для бэкапа Cassandra.
*   **Neo4j (для Social Service - граф):**
    *   **Частота:** Ежедневные полные бэкапы.
    *   **RPO:** 24 часа.
    *   **RTO:** 2-4 часа.
    *   **Хранение:** Бэкапы хранятся 14-30 дней.
    *   **Инструменты:** `neo4j-admin backup`.
*   **MongoDB (для Admin Service - логи аудита):**
    *   **Частота:** Ежедневные полные бэкапы.
    *   **RPO:** 24 часа.
    *   **RTO:** 2-4 часа.
    *   **Хранение:** Бэкапы хранятся 30 дней, с возможностью долгосрочного архивирования для логов аудита согласно политикам.
    *   **Инструменты:** `mongodump`.
*   **S3-совместимое хранилище (для игровых файлов, медиа, сохранений):**
    *   **Стратегия:** "Данные в S3 обычно имеют высокую встроенную долговечность за счет репликации на стороне провайдера. Дополнительно может быть настроено версионирование объектов и межрегиональная репликация для критически важных бакетов."
    *   **RPO/RTO:** Определяются возможностями S3-провайдера и настройками версионирования/репликации.

### 6.3. Резервное Копирование Конфигураций
*   Конфигурации Kubernetes (манифесты, Helm-чарты), конфигурации CI/CD пайплайнов и другие важные конфигурационные файлы должны храниться в системе контроля версий (Git) и регулярно бэкапироваться вместе с репозиториями.
*   Секреты, управляемые через Vault или Kubernetes Secrets, должны иметь свою процедуру резервного копирования и восстановления, специфичную для этих систем.

### 6.4. Ответственность
*   Команда DevOps/SRE отвечает за настройку, мониторинг и тестирование процедур резервного копирования и восстановления.
*   Владельцы микросервисов (Tech Leads) отвечают за определение критичности данных своих сервисов и требований к RPO/RTO.

---
*Этот документ должен регулярно обновляться по мере развития моделей данных микросервисов.*
