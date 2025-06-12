# Структура Баз Данных Проекта

## 1. Введение

Данный документ предназначен для описания общей структуры баз данных, используемых в проекте "Российский Аналог Платформы Steam". Он должен предоставлять высокоуровневый обзор моделей данных каждого микросервиса и, по возможности, общую схему взаимодействия данных на уровне платформы.

**Цели документа:**

*   Дать общее представление о данных, которыми оперирует каждый микросервис.
*   Описать ключевые сущности и их связи на уровне отдельных сервисов.
*   Если существуют общие модели данных или разделяемые базы данных (не рекомендуется в чистой микросервисной архитектуре, но возможно для некоторых случаев), описать их здесь.
*   Служить отправной точкой для разработчиков при изучении моделей данных конкретных сервисов.

## 2. Общие Принципы Проектирования Баз Данных

*   **Децентрализация данных:** Каждый микросервис владеет своими данными и своей схемой базы данных. Прямой доступ к базе данных другого сервиса запрещен; взаимодействие осуществляется только через API.
*   **Выбор СУБД:** Тип СУБД (SQL, NoSQL, графовая и т.д.) выбирается исходя из специфических требований каждого микросервиса (например, PostgreSQL для транзакционных данных, Elasticsearch для поиска, Cassandra для сценариев с высокой нагрузкой на запись).
*   **Миграции схемы:** Для реляционных баз данных используются инструменты управления миграциями схемы (например, golang-migrate, Flyway, Alembic).
*   **Резервное копирование и восстановление:** Для каждой базы данных должны быть настроены процедуры регулярного резервного копирования и протестированы сценарии восстановления.
*   **Безопасность:** Доступ к базам данных должен быть строго ограничен. Учетные данные для доступа должны храниться в секретах.

## 3. Структура Баз Данных по Микросервисам

Этот раздел должен содержать ссылки на документацию по структуре БД каждого микросервиса или краткое описание их моделей данных.

### 3.1. Auth Service
*   **СУБД:** PostgreSQL, Redis.
*   **Основные сущности:** Пользователи (учетные данные), Роли, Разрешения, Сессии, Refresh-токены, API-ключи, Коды верификации, Аудит лог.
*   **Ссылка на детальную документацию:** `services/auth-service/docs/README.md#4-модели-данных-data-models`

### 3.2. Account Service
*   **СУБД:** PostgreSQL, Redis.
*   **Основные сущности:** Аккаунты пользователей (профильная информация), Контактная информация, Настройки пользователя, Аватары.
*   **Ссылка на детальную документацию:** `services/account-service/docs/README.md#4-модели-данных-data-models`

### 3.3. Catalog Service
*   **СУБД:** PostgreSQL, Elasticsearch, Redis.
*   **Основные сущности:** Продукты (Игры, DLC, ПО, Комплекты), Жанры, Теги, Категории, Цены, Медиа-контент, Метаданные достижений.
*   **Ссылка на детальную документацию:** `services/catalog-service/docs/README.md#4-модели-данных-data-models`

### 3.4. Library Service
*   **СУБД:** PostgreSQL, Redis, S3-совместимое хранилище.
*   **Основные сущности:** Записи библиотеки пользователя, Игровое время, Достижения пользователя, Списки желаемого, Метаданные сохранений игр.
*   **Ссылка на детальную документацию:** `services/library-service/docs/README.md#4-модели-данных-data-models`

### 3.5. Payment Service
*   **СУБД:** PostgreSQL, Redis.
*   **Основные сущности:** Транзакции, Платежные методы, Фискальные чеки, Балансы разработчиков, Промокоды, Подарочные сертификаты, Выплаты.
*   **Ссылка на детальную документацию:** `services/payment-service/docs/README.md#4-модели-данных-data-models`

### 3.6. Download Service
*   **СУБД:** PostgreSQL, Redis.
*   **Основные сущности:** Задачи на загрузку, Элементы загрузки, Информация об обновлениях, Метаданные файлов.
*   **Ссылка на детальную документацию:** `services/download-service/docs/README.md#4-модели-данных-data-models`

### 3.7. Social Service
*   **СУБД:** PostgreSQL, Cassandra, Neo4j, Redis.
*   **Основные сущности:** Профили пользователей (социальная часть), Друзья, Группы, Сообщения чатов, Записи ленты активности, Отзывы, Комментарии.
*   **Ссылка на детальную документацию:** `services/social-service/docs/README.md#4-модели-данных-data-models`

### 3.8. Notification Service
*   **СУБД:** PostgreSQL или ClickHouse (для статистики), Redis.
*   **Основные сущности:** Шаблоны уведомлений, Отправленные уведомления, Пользовательские предпочтения, Маркетинговые кампании, Токены устройств.
*   **Ссылка на детальную документацию:** `services/notification-service/docs/README.md#4-модели-данных-data-models`

### 3.9. Developer Service
*   **СУБД:** PostgreSQL, S3-совместимое хранилище.
*   **Основные сущности:** Аккаунты разработчиков, Игры разработчика, Версии игр (билды), Метаданные игр, Ценовая информация, Запросы на выплаты, API-ключи.
*   **Ссылка на детальную документацию:** `services/developer-service/docs/README.md#4-модели-данных-data-models`

### 3.10. Admin Service
*   **СУБД:** PostgreSQL, MongoDB.
*   **Основные сущности:** Административные пользователи, Очереди модерации, Решения по модерации, Тикеты поддержки, Статьи базы знаний, Системные настройки, Маркетинговые кампании.
*   **Ссылка на детальную документацию:** `services/admin-service/docs/README.md#4-модели-данных-data-models`

### 3.11. Analytics Service
*   **СУБД:** ClickHouse, PostgreSQL (метаданные).
*   **Основные сущности:** События, Метрики, Отчеты, Сегменты, ML Модели, Прогнозы.
*   **Ссылка на детальную документацию:** `services/analytics-service/docs/README.md#4-модели-данных-data-models`

### 3.12. API Gateway
*   **СУБД:** Может использовать собственную БД (например, PostgreSQL или Redis) для хранения конфигурации маршрутов, плагинов, потребителей, если это не управляется исключительно через CRD Kubernetes.
*   **Основные сущности:** Маршруты, Сервисы, Потребители, Плагины.
*   **Ссылка на детальную документацию:** `services/api-gateway/docs/README.md#4-модели-данных-data-models`

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

---
*Этот документ должен регулярно обновляться по мере развития моделей данных микросервисов.*
