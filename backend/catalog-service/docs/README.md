# Спецификация Микросервиса: Catalog Service

**Версия:** 1.1
**Дата последнего обновления:** {{YYYY-MM-DD}}

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой спецификацию микросервиса Catalog Service, являющегося ядром платформы "Российский Аналог Steam" и отвечающего за всю информацию о цифровых продуктах.
*   **Роль в общей архитектуре платформы:** Централизованное управление каталогом продуктов (игры, DLC, программное обеспечение, комплекты), ценообразованием, акциями, таксономией (жанры, теги, категории), медиа-контентом и метаданными достижений. Предоставляет данные другим микросервисам и обеспечивает поиск и обнаружение продуктов пользователями.
*   **Основные бизнес-задачи:** Управление жизненным циклом продуктов, метаданными, ценами, скидками, таксономией, медиа-контентом, метаданными достижений; обеспечение поиска и навигации; поддержка модерации.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   CRUD операции для продуктов, жанров, тегов, категорий, цен, медиа, достижений.
*   Управление локализованными метаданными.
*   Управление ценообразованием (базовые/региональные цены, скидки, акции).
*   Управление таксономией (жанры, теги, категории, франшизы, коллекции).
*   API для полнотекстового поиска, фильтрации, сортировки. Пагинация.
*   Управление медиа-контентом.
*   Управление метаданными достижений.
*   Отслеживание и предоставление статусов модерации.
*   API для систем рекомендаций.
*   Публикация событий об изменениях в каталоге.

### 1.3. Основные Технологии
*   **Язык программирования:** Go (1.21+).
*   **Веб-фреймворк (REST):** Echo.
*   **RPC фреймворк (gRPC):** `google.golang.org/grpc`.
*   **База данных (основная):** PostgreSQL (GORM/pgx).
*   **Поисковый движок:** Elasticsearch.
*   **Кэширование:** Redis.
*   **Брокер сообщений:** Apache Kafka.
*   **Управление конфигурацией:** Viper.
*   **Логирование:** Zap.
*   **Мониторинг/Трассировка:** OpenTelemetry SDK, Prometheus client.
*   **Инфраструктура:** Docker, Kubernetes.
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `../../../../PACKAGE_STANDARDIZATION.md`.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`, `../../../../project_observability_standards.md`.

### 1.4. Термины и Определения (Glossary)
*   См. `../../../../project_glossary.md`.
*   **Продукт (Product):** Цифровой товар (Игра, DLC, ПО, Комплект).
*   **Метаданные (Metadata):** Описательная информация о Продукте.
*   **Таксономия (Taxonomy):** Система классификации (Жанры, Теги, Категории).
*   **Цена (Price):** Стоимость Продукта.
*   **Медиа-контент (Media Content):** Графические/видео материалы.
*   **Достижение (Achievement):** Внутриигровое достижение.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Catalog Service реализуется как независимый микросервис, следуя принципам **Чистой Архитектуры (Clean Architecture)** и **CQRS (Command Query Responsibility Segregation)** на логическом уровне.
*   Диаграмма Архитектуры (Clean Architecture + CQRS):
```mermaid
graph TD
    subgraph User/Client Interaction
        Clients[Клиенты (Web, Mobile, API Gateway)]
    end

    subgraph Catalog Service
        direction TB

        subgraph PresentationLayer [Presentation Layer (Адаптеры Транспорта)]
            APIs[REST API (Echo) / gRPC API]
        end

        subgraph ApplicationLayer [Application Layer (Сценарии Использования)]
            CommandHandlers[Обработчики Команд (CreateProduct, UpdatePrice)]
            QueryServices[Сервисы Запросов (GetProductDetails, SearchProducts)]
        end

        subgraph DomainLayer [Domain Layer (Бизнес-логика и Сущности)]
            Entities[Сущности (Product, Genre, Price)]
            Aggregates[Агрегаты (ProductAggregate)]
            DomainEvents[Доменные События (ProductCreated, PriceUpdated)]
            RepositoryIntf[Интерфейсы Репозиториев]
        end

        subgraph InfrastructureLayer [Infrastructure Layer (Внешние Зависимости)]
            PostgresAdapter[Адаптер PostgreSQL (GORM/Squirrel)]
            ElasticAdapter[Адаптер Elasticsearch]
            RedisAdapter[Адаптер Redis (Кэш)]
            KafkaProducer[Продюсер Kafka (События)]
            Config[Конфигурация (Viper)]
            Logging[Логирование (Zap)]
            Metrics[Метрики (Prometheus)]
        end

        APIs --> CommandHandlers
        APIs --> QueryServices
        CommandHandlers --> RepositoryIntf
        CommandHandlers --> DomainEvents
        QueryServices --> RepositoryIntf
        QueryServices --> ElasticAdapter  # Для поисковых запросов
        QueryServices --> RedisAdapter    # Для кэшированных запросов

        RepositoryIntf -- Реализуются --> PostgresAdapter
        DomainEvents -- Публикуются через --> KafkaProducer
    end

    Clients --> APIs

    PostgresAdapter --> DB[(PostgreSQL)]
    ElasticAdapter --> ES[(Elasticsearch)]
    RedisAdapter --> Cache[(Redis)]
    KafkaProducer --> Kafka[Kafka Broker]

    classDef layer_boundary fill:#f9f9f9,stroke:#333,stroke-width:2px,color:#333
    classDef component_major fill:#e6f0ff,stroke:#007bff,color:#000
    classDef component_minor fill:#d4edda,stroke:#28a745,color:#000
    classDef datastore fill:#f8d7da,stroke:#dc3545,color:#000

    class PresentationLayer,ApplicationLayer,DomainLayer,InfrastructureLayer layer_boundary
    class APIs,CommandHandlers,QueryServices,Entities,Aggregates,DomainEvents,RepositoryIntf component_major
    class PostgresAdapter,ElasticAdapter,RedisAdapter,KafkaProducer,Config,Logging,Metrics component_minor
    class DB,ES,Cache,Kafka datastore
```

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления / Адаптеры Транспорта)
*   Ответственность: Обработка REST и gRPC запросов, валидация DTO, вызов Application Layer, преобразование результатов.
*   Ключевые компоненты/модули: HTTP хендлеры, gRPC серверы, DTO.

#### 2.2.2. Application Layer (Прикладной Слой)
*   Ответственность: Реализация use cases, координация Domain и Infrastructure Layers.
*   Ключевые компоненты/модули: Сервисы сценариев использования (`ProductApplicationService`), обработчики команд/запросов, DTO.

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-логика, правила, сущности, агрегаты.
*   Ключевые компоненты/модули: Entities (`Product`, `Genre`, `Tag`, `Category`), Aggregates (`ProductAggregate`), Value Objects (`LocalizedString`), Domain Events (`ProductCreated`), Repository Interfaces.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Реализация интерфейсов, взаимодействие с PostgreSQL, Elasticsearch, Redis, Kafka.
*   Ключевые компоненты/модули: Адаптеры БД/поиска/кэша, Kafka продюсеры/консьюмеры, утилиты.

## 3. API Endpoints
Общие принципы и форматы см. `../../../../project_api_standards.md`.

### 3.1. REST API
*   **Базовый URL:** `/api/v1/catalog`.
*   **Аутентификация:** Публичные эндпоинты для чтения. Управляющие эндпоинты (`/manage`) требуют JWT (admin, manager).
*   **Авторизация:** RBAC.
*   **Пагинация:** `page`, `limit`.
*   **Локализация:** Заголовок `Accept-Language`.

#### 3.1.1. Ресурс: Продукты (Products)
*   **`GET /products`**
    *   Описание: Список продуктов с фильтрацией, сортировкой, пагинацией.
    *   Query параметры: `search`, `genre_ids`, `tag_ids`, `category_ids`, `platform`, `min_price`, `max_price`, `is_on_sale`, `sort_by`, `page`, `limit`.
    *   Пример ответа (Успех 200 OK): `[NEEDS DEVELOPER INPUT: Example JSON response for GET /products, including localized fields and price info]`
    *   Пример ответа (Ошибка 400 Validation Error): (см. существующий документ).
    *   Требуемые права доступа: Публичный.
*   **`GET /products/{product_id}`**
    *   Описание: Детальная информация о продукте.
    *   Пример ответа (Успех 200 OK): `[NEEDS DEVELOPER INPUT: Example JSON response for GET /products/{product_id}, including all relevant nested structures like media, achievements, system requirements, etc.]`
    *   Требуемые права доступа: Публичный.
*   **`GET /products/slug/{product_slug}`**
    *   Описание: Получение продукта по его уникальному текстовому идентификатору (slug).
    *   Требуемые права доступа: Публичный.

#### 3.1.2. Ресурс: Жанры (Genres)
*   **`GET /genres`**
    *   Описание: Список всех жанров.
    *   Пример ответа (Успех 200 OK): `[NEEDS DEVELOPER INPUT: Example JSON response for GET /genres]`
    *   Требуемые права доступа: Публичный.
*   **`GET /genres/{genre_id_or_slug}/products`**
    *   Описание: Список продуктов, принадлежащих указанному жанру.
    *   Требуемые права доступа: Публичный.

#### 3.1.3. Ресурс: Теги (Tags)
*   **`GET /tags`**
    *   Описание: Список всех тегов.
    *   Требуемые права доступа: Публичный.
*   **`GET /tags/{tag_id_or_slug}/products`**
    *   Описание: Список продуктов с указанным тегом.
    *   Требуемые права доступа: Публичный.

#### 3.1.4. Ресурс: Категории (Categories)
*   **`GET /categories`**
    *   Описание: Список всех категорий (может быть иерархическим).
    *   Требуемые права доступа: Публичный.
*   **`GET /categories/{category_id_or_slug}/products`**
    *   Описание: Список продуктов из указанной категории.
    *   Требуемые права доступа: Публичный.

#### 3.1.5. Ресурс: Поиск (Search)
*   **`GET /search`**
    *   Описание: Полнотекстовый поиск по продуктам с расширенными возможностями фильтрации. Является основным способом поиска для клиентов.
    *   Query параметры: `query` (текстовый запрос), `filters` (JSON объект или набор параметров для фильтрации по жанрам, тегам, цене, платформе, языку и т.д.), `sort_by`, `page`, `limit`.
    *   Пример ответа: `[NEEDS DEVELOPER INPUT: Example JSON response for GET /search, showing search results structure and facets/aggregations if supported]`
    *   Требуемые права доступа: Публичный.

#### 3.1.6. Управление Каталогом (Management API - префикс `/manage`)
*   **`POST /manage/products`**: Создание продукта. Тело запроса: (см. существующий документ). Требуемые права: `catalog_admin`, `product_manager`.
*   **`PUT /manage/products/{product_id}`**: Полное обновление продукта.
*   **`PATCH /manage/products/{product_id}`**: Частичное обновление продукта.
*   **`DELETE /manage/products/{product_id}`**: Удаление продукта (логическое или физическое).
*   **`POST /manage/genres`**: Создание жанра.
*   **`PUT /manage/genres/{genre_id}`**: Обновление жанра.
*   `[NEEDS DEVELOPER INPUT: List other key management endpoints for products, genres, tags, categories, prices, media, achievements, including request/response examples and permissions.]`

### 3.2. gRPC API
*   **Пакет:** `catalog.v1`
*   **Файл .proto:** `proto/catalog/v1/catalog_service.proto` (или в `platform-protos`).
*   **Аутентификация:** mTLS.
#### 3.2.1. Сервис: CatalogService
*   **`rpc GetProductInternal(GetProductInternalRequest) returns (ProductInternalResponse)`**: Для внутреннего использования (например, Payment Service для получения цены).
    *   `message GetProductInternalRequest { string product_id = 1; string region_code = 2; string currency_code = 3; string language_code = 4; }`
    *   `message ProductInternalResponse { Product product = 1; ProductPrice current_price = 2; }` (`Product` и `ProductPrice` - proto-версии сущностей).
*   **`rpc GetProductsBatchInternal(GetProductsBatchInternalRequest) returns (GetProductsBatchInternalResponse)`**: Получение информации по нескольким продуктам.
*   `[NEEDS DEVELOPER INPUT: Add other relevant gRPC methods for internal service communication, e.g., for price validation, stock checking if applicable (though less for digital goods).]`

### 3.3. WebSocket API
*   Не применимо для данного сервиса.

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`Product`**: Продукт (ID, тип, статус, локализованные названия/описания, дата релиза, разработчики, издатели).
*   **`Genre`**: Жанр (ID, локализованное имя, slug).
*   **`Tag`**: Тег (ID, локализованное имя, slug).
*   **`Category`**: Категория (ID, локализованное имя, slug, родительская категория).
*   **`ProductPrice`**: Цена Продукта (ID, ProductID, регион, валюта, базовая сумма, сумма скидки, даты действия).
*   **`MediaItem`**: Медиа-элемент (ID, ProductID, тип, URL, URL превью, порядок).
*   **`AchievementMeta`**: Метаданные Достижения (ID, ProductID, API имя, локализованные имя/описание, URL иконок, скрытость).
*   **`SystemRequirements`**: Системные требования (тип: min/rec, платформа: pc/mac, ОС, процессор, память, графика, диск). Внедряется в `Product`.
*   **`LocalizedString`**: Объект для локализованных строк (например, `{"ru-RU": "текст", "en-US": "text"}`).
*   `[NEEDS DEVELOPER INPUT: Review and confirm if any other key entities are missing, e.g., Bundles, Editions, Franchises, Collections, and their detailed attributes and relationships.]`

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL
*   ERD Диаграмма: (см. существующий документ).
*   DDL: (см. существующий документ, включая таблицы `tags`, `product_tags`, `categories`, `product_categories`, `achievement_metadata`).
*   `[NEEDS DEVELOPER INPUT: Ensure DDL includes tables for any missing entities like Bundles, Editions, Franchises, Collections and their linking tables if they are part of Catalog Service's direct responsibility.]`

#### 4.2.2. Elasticsearch
*   **Индексы:** `products_catalog_vX`.
*   **Структура документа:** Денормализованный документ продукта, включающий ключевые поля для поиска и фильтрации (названия, описания, жанры, теги, категории, цены, платформы, языки, рейтинг).
*   **Анализаторы:** Специфичные для языка анализаторы (russian, english) для текстовых полей.
*   `[NEEDS DEVELOPER INPUT: Provide a more detailed Elasticsearch mapping example for the 'products_catalog_vX' index, showing key fields, types, and analyzer configurations.]`

#### 4.2.3. Redis (Кэширование)
*   **Стратегия:** Cache-Aside. Инвалидация при обновлении данных в PostgreSQL или через события Kafka.
*   **Ключи:** `product:<product_id>:<lang>:<region>`, `products_list:<hash_params>`, `genres_all:<lang>`, `price:<product_id>:<region>`.
*   TTL: Варьируется (1 мин - 24 часа).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   Формат: CloudEvents JSON. Топик: `com.platform.catalog.events.v1`.
*   **`com.platform.catalog.product.created.v1`**: `data`: `{ "productId": "uuid", "type": "game", ... }`
*   **`com.platform.catalog.product.updated.v1`**: `data`: `{ "productId": "uuid", "updatedFields": ["titles", "descriptions"], ... }`
*   **`com.platform.catalog.product.status.changed.v1`**: `data`: `{ "productId": "uuid", "oldStatus": "in_review", "newStatus": "published", ... }`
*   **`com.platform.catalog.price.updated.v1`**: `data`: `{ "productId": "uuid", "priceId": "uuid", ... }`
*   **`com.platform.catalog.genre.created.v1`**: `data`: `{ "genreId": "uuid", "names": {...}, ... }`
*   **`com.platform.catalog.tag.assigned.v1`**: `data`: `{ "productId": "uuid", "tagId": "uuid", ... }`
*   `[NEEDS DEVELOPER INPUT: List or confirm other key published events for catalog changes (e.g., product_deleted, category_created, media_added, achievement_updated) including their payload structures.]`

### 5.2. Потребляемые События (Consumed Events)
*   **`com.platform.developer.product.submitted.v1`** (от Developer Service)
    *   Описание: Разработчик отправил новый продукт или обновление на рассмотрение/модерацию.
    *   Логика обработки: Создание/обновление продукта в статусе "на рассмотрении", запуск процесса модерации.
*   **`com.platform.admin.product.moderation.approved.v1`** (от Admin Service)
    *   Описание: Продукт одобрен модератором.
    *   Логика обработки: Изменение статуса продукта на "опубликован" или "готов к публикации".
*   **`com.platform.admin.product.moderation.rejected.v1`** (от Admin Service)
    *   Описание: Продукт отклонен модератором.
    *   Логика обработки: Изменение статуса продукта на "отклонен", уведомление разработчика.
*   `[NEEDS DEVELOPER INPUT: List or confirm other key consumed events (e.g., from Payment Service for sales events if they influence 'top sellers' categories directly, or from Social Service for review counts/ratings if denormalized here).]`

## 6. Интеграции (Integrations)
См. `../../../../project_integrations.md`.
[NEEDS DEVELOPER INPUT: Review and confirm service boundaries and reliability strategies for each integration for catalog-service.]

### 6.1. Внутренние Микросервисы
*   **API Gateway**: Проксирование запросов.
*   **Developer Service**: Управление продуктами со стороны разработчиков (через события или прямое API Catalog Service для создания/обновления). `[NEEDS DEVELOPER INPUT: Clarify if Developer Service calls Catalog's management API or if it's event-driven]`
*   **Admin Service**: Модерация контента, управление глобальными настройками каталога.
*   **Payment Service**: Получение актуальных цен (через gRPC).
*   **Library Service, Download Service**: Предоставление метаданных продуктов.
*   **Analytics Service**: Обмен данными о продуктах.
*   **Auth Service**: Валидация токенов.
*   **Social Service**: Потенциально для получения информации о пользовательских отзывах/рейтингах для агрегации в каталоге. `[NEEDS DEVELOPER INPUT: Clarify integration with Social Service for reviews/ratings.]`

### 6.2. Внешние Системы
*   Не предполагается прямых интеграций.

## 7. Конфигурация (Configuration)
Общие стандарты: `../../../../project_api_standards.md` и `../../../../DOCUMENTATION_GUIDELINES.md`.

### 7.1. Переменные Окружения
*   `CATALOG_HTTP_PORT`, `CATALOG_GRPC_PORT`
*   `POSTGRES_DSN`, `ELASTICSEARCH_URLS`, `REDIS_ADDR`, `KAFKA_BROKERS`
*   `DEFAULT_LANGUAGE`, `DEFAULT_REGION`, `DEFAULT_CURRENCY`
*   `JWT_PUBLIC_KEY_PATH` (для валидации токенов админов/менеджеров)
*   `LOG_LEVEL`, `OTEL_EXPORTER_JAEGER_ENDPOINT`
*   `[NEEDS DEVELOPER INPUT: Add specific env vars for search indexing parameters, cache TTLs, or any recommendation algorithm parameters if managed here.]`

### 7.2. Файлы Конфигурации (`configs/config.yaml`)
*   Расположение: `backend/catalog-service/configs/config.yaml`.
*   Структура: (см. существующий документ, проверено).

## 8. Обработка Ошибок (Error Handling)
*   Стандартные HTTP коды и форматы ошибок REST согласно `../../../../project_api_standards.md`.
*   gRPC ошибки согласно стандартам gRPC.

### 8.1. Общие Принципы
*   Стандартный JSON формат ответа об ошибке для REST.

### 8.2. Распространенные Коды Ошибок
*   `PRODUCT_NOT_FOUND`, `GENRE_NOT_FOUND`, `PRICE_NOT_FOUND`, `VALIDATION_ERROR`, `SEARCH_ENGINE_UNAVAILABLE`.
*   `[NEEDS DEVELOPER INPUT: Review and add other specific error codes for catalog-service.]`

## 9. Безопасность (Security)
См. `../../../../project_security_standards.md`.

### 9.1. Аутентификация
*   Публичные API для чтения. JWT для управляющих API.

### 9.2. Авторизация
*   RBAC для управляющих API (роли `catalog_admin`, `product_manager`).

### 9.3. Защита Данных
*   ФЗ-152: Локализованные метаданные могут содержать ПДн (если созданы пользователями). ID разработчиков.
*   Защита от SQL-инъекций (через ORM/драйверы), XSS (при отображении данных).

### 9.4. Управление Секретами
*   Пароли БД, ключи доступа к Kafka/Elasticsearch через Kubernetes Secrets / Vault.

## 10. Развертывание (Deployment)
См. `../../../../project_deployment_standards.md`.

### 10.1. Инфраструктурные Файлы
*   Dockerfile: `backend/catalog-service/Dockerfile`.
*   Helm-чарты: `deploy/charts/catalog-service/`.

### 10.2. Зависимости при Развертывании
*   PostgreSQL, Elasticsearch, Redis, Kafka.
*   API Gateway, Auth Service.

### 10.3. CI/CD
*   Стандартный пайплайн. Автоматические миграции БД. Индексация Elasticsearch при необходимости.

## 11. Мониторинг и Логирование (Logging and Monitoring)
См. `../../../../project_observability_standards.md`.

### 11.1. Логирование
*   Формат: JSON (Zap).
*   Ключевые события: CRUD операции, поиск, ошибки, обновления цен/статусов.
*   Интеграция: Loki/ELK.

### 11.2. Мониторинг
*   Метрики (Prometheus): Запросы API, ошибки, задержки, производительность Elasticsearch, Kafka лаги, размеры кэша.
*   Дашборды (Grafana): Обзор состояния, API, Elasticsearch, Kafka, PostgreSQL.
*   Алерты (AlertManager): Ошибки API, недоступность зависимостей, проблемы с индексацией/кэшем.

### 11.3. Трассировка
*   Интеграция: OpenTelemetry. Экспорт: Jaeger.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**: P99 поиск < 200мс. P99 детали продукта < 100мс. P99 обновление < 300мс.
*   **Масштабируемость**: >1 млн продуктов. >1000 запросов/сек на чтение.
*   **Надежность**: Доступность 99.95%.
*   **Консистентность данных**: Strong consistency для основного хранилища (PostgreSQL), eventual consistency для Elasticsearch и кэшей (задержка < 1 мин).
*   [NEEDS DEVELOPER INPUT: Confirm or update specific NFR values for catalog-service.]

## 13. Резервное Копирование и Восстановление (Backup and Recovery)
Общие принципы см. в `../../../../project_database_structure.md`.

### 13.1. PostgreSQL
*   Процедура: Ежедневный `pg_dump`. Непрерывная архивация WAL (PITR).
*   Хранение: S3. RTO: < 2ч. RPO: < 5мин.

### 13.2. Elasticsearch
*   Стратегия: Переиндексация из PostgreSQL (источник правды). Снапшоты для ускорения.
*   RTO: < 1ч (снапшот), < 6-12ч (переиндексация). RPO: 24ч (снапшот), 0 (переиндексация).

### 13.3. Redis
*   Стратегия: Данные - кэш, восстановимы. Персистентность для "прогрева".
*   RTO/RPO: Неприменимо для потери данных.

## 14. Приложения (Appendices) (Опционально)
*   OpenAPI спецификация: `[NEEDS DEVELOPER INPUT: Link to OpenAPI spec or state if generated from code for catalog-service]`
*   Protobuf схемы: `platform-protos/catalog/v1/catalog_service.proto` (предположительно).
*   `[NEEDS DEVELOPER INPUT: Add any other appendices if necessary for catalog-service.]`

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [Процесс покупки игры и обновления библиотеки пользователя](../../../../project_workflows/game_purchase_flow.md)
*   [Процесс подачи разработчиком новой игры на модерацию](../../../../project_workflows/game_submission_flow.md)
*   [Процесс управления промо-акциями и скидками] `[NEEDS DEVELOPER INPUT: Create and link workflow document, e.g., project_workflows/promo_management_flow.md]`

---
*Этот документ является основной спецификацией для Catalog Service и должен поддерживаться в актуальном состоянии.*
