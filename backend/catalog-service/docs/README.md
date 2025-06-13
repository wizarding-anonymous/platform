# Спецификация Микросервиса: Catalog Service

**Версия:** 1.1 (адаптировано из v3.0 внешней спецификации)
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой спецификацию микросервиса Catalog Service, являющегося ядром платформы "Российский Аналог Steam" и отвечающего за всю информацию о цифровых продуктах.
*   **Роль в общей архитектуре платформы:** Централизованное управление каталогом продуктов (игры, DLC, программное обеспечение, комплекты), ценообразованием, акциями, таксономией (жанры, теги, категории), медиа-контентом и метаданными достижений. Предоставляет данные другим микросервисам (например, Поисковому Сервису, Сервису Рекомендаций, Сервису Библиотек Пользователей) и обеспечивает поиск и обнаружение продуктов пользователями через клиентские приложения (веб, мобильные, десктопные).
*   **Основные бизнес-задачи:**
    *   Управление жизненным циклом продуктов (от создания до снятия с продажи).
    *   Управление метаданными продуктов, включая локализованные названия, описания, системные требования.
    *   Управление ценами (базовые, региональные), скидками и промо-акциями.
    *   Структурирование каталога через жанры, теги, категории, коллекции и франшизы.
    *   Обеспечение возможностей поиска, фильтрации и навигации по каталогу для пользователей и других сервисов.
    *   Управление медиа-контентом (скриншоты, трейлеры, арты, обложки).
    *   Управление метаданными достижений для игр.
    *   Поддержка процесса модерации контента продуктов (интеграция с Admin Service).
    *   Предоставление данных для системы рекомендаций.
*   Разработка сервиса должна вестись в соответствии с `CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   CRUD операции для продуктов, жанров, тегов, категорий, цен, медиа-элементов, метаданных достижений.
*   Управление локализованными метаданными продуктов.
*   Управление ценообразованием: установка базовых цен, региональных цен, управление скидками и промо-акциями.
*   Управление таксономией: присвоение продуктам жанров, тегов, категорий; управление франшизами и коллекциями.
*   Предоставление API для полнотекстового поиска продуктов с возможностью фильтрации (по жанрам, тегам, ценам, платформам, языкам и т.д.) и сортировки.
*   Пагинированный вывод списков продуктов и других сущностей каталога.
*   Управление медиа-контентом: загрузка, хранение ссылок, определение типов медиа.
*   Управление метаданными достижений: названия, описания, иконки, условия разблокировки.
*   Отслеживание и предоставление статусов модерации продуктов.
*   Предоставление API для систем рекомендаций (например, "похожие товары", "новые релизы").
*   Публикация событий об изменениях в каталоге (например, создание продукта, обновление цены).

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+).
*   **Веб-фреймворк (REST):** Echo (v4+).
*   **RPC фреймворк (gRPC):** `google.golang.org/grpc`.
*   **База данных (основная):** PostgreSQL (версия 15+) для хранения структурированных данных каталога.
*   **Поисковый движок:** Elasticsearch (версия 8.x+) для полнотекстового поиска и сложной фильтрации.
*   **Кэширование:** Redis (версия 7.0+) для кэширования часто запрашиваемых данных (детали продуктов, списки).
*   **Брокер сообщений:** Apache Kafka (версия 3.x+) для асинхронной публикации событий об изменениях в каталоге.
*   **Конфигурация:** Viper (для управления конфигурацией).
*   **Логирование:** Zap (для структурированного логирования).
*   **Мониторинг/Трассировка:** OpenTelemetry SDK, Prometheus client library, Grafana, Jaeger/Tempo.
*   **ORM/Query Builder (PostgreSQL):** GORM или Squirrel (будет уточнено).
*   **Клиент Elasticsearch:** Официальный Go клиент для Elasticsearch.
*   **Клиент Kafka:** `segmentio/kafka-go` или `confluent-kafka-go`.
*   **Инфраструктура:** Docker, Kubernetes.
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `PACKAGE_STANDARDIZATION.md`.

### 1.4. Термины и Определения (Glossary)
*   **Продукт (Product):** Любой цифровой товар, доступный на платформе (например, Игра, Дополнение (DLC), Программное обеспечение, Комплект продуктов).
*   **Метаданные (Metadata):** Описательная информация о Продукте, такая как название, описание, разработчик, издатель, системные требования, дата выпуска и т.д.
*   **Таксономия (Taxonomy):** Система классификации Продуктов, включающая Жанры, Теги, Категории.
*   **Цена (Price):** Стоимость Продукта. Может включать базовую цену, региональные цены и временные скидки.
*   **Медиа-контент (Media Content):** Графические и видео материалы, связанные с продуктом (скриншоты, трейлеры, арты).
*   **Достижение (Achievement):** Внутриигровое достижение, метаданные которого хранятся в каталоге.
*   Для других общих терминов см. `project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Catalog Service реализуется как независимый микросервис, следуя принципам **Чистой Архитектуры (Clean Architecture)** для разделения ответственностей и улучшения тестируемости.
*   Для оптимизации операций чтения и записи используется подход **CQRS (Command Query Responsibility Segregation)** на логическом уровне. Команды (операции записи) работают с основной доменной моделью и персистентным хранилищем (PostgreSQL). Запросы (операции чтения) могут использовать оптимизированные для чтения реплики, денормализованные данные или специализированные индексы (Elasticsearch, Redis кэш) для повышения производительности и гибкости.

**Диаграмма Архитектуры (Clean Architecture + CQRS):**
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
*   **Ответственность:** Обработка входящих REST (через Echo) и gRPC запросов от внешних клиентов или других микросервисов. Валидация данных запроса (DTO), аутентификация/авторизация (частично, основная проверка на API Gateway), вызов соответствующих сервисов в Application Layer. Преобразование результатов из Application Layer в формат ответа (JSON, Protobuf).
*   **Ключевые компоненты/модули:** HTTP хендлеры (контроллеры), gRPC серверы и их реализации, DTO для запросов и ответов, парсеры запросов, форматеры ответов.

#### 2.2.2. Application Layer (Прикладной Слой)
*   **Ответственность:** Реализация сценариев использования (use cases) системы. Координирует взаимодействие между Domain Layer и Infrastructure Layer. Содержит логику команд (изменение состояния системы) и запросов (получение данных). Не содержит бизнес-правил напрямую, а делегирует их Domain Layer.
*   **Ключевые компоненты/модули:** Сервисы сценариев использования (например, `ProductApplicationService`, `PriceApplicationService`), обработчики команд (`CreateProductCommandHandler`), обработчики запросов (`GetProductQueryHandler`), DTO, используемые для передачи данных между слоями.

#### 2.2.3. Domain Layer (Доменный Слой)
*   **Ответственность:** Содержит всю бизнес-логику, бизнес-правила, сущности и агрегаты. Этот слой не зависит от деталей реализации других слоев (например, от конкретной СУБД или фреймворка).
*   **Ключевые компоненты/модули:**
    *   **Entities (Сущности):** Основные объекты домена, обладающие идентичностью (например, `Product`, `Genre`, `Tag`, `Category`, `ProductPrice`, `MediaItem`, `AchievementMeta`).
    *   **Aggregates (Агрегаты):** Кластеры сущностей и объектов-значений, рассматриваемые как единое целое для операций изменения данных (например, `ProductAggregate` может включать сам продукт, его цены, медиа-файлы). Корень агрегата обеспечивает консистентность.
    *   **Value Objects (Объекты-значения):** Неизменяемые объекты, характеризующиеся своими атрибутами (например, `LocalizedString` для локализованных строк, `SystemRequirements`).
    *   **Domain Services:** Сервисы, инкапсулирующие доменную логику, которая не принадлежит естественным образом ни одной сущности.
    *   **Domain Events:** События, отражающие значимые изменения в домене (например, `ProductCreated`, `PriceUpdated`).
    *   **Repository Interfaces:** Интерфейсы, определяющие контракты для сохранения и извлечения агрегатов из хранилища.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   **Ответственность:** Реализация интерфейсов, определенных в Domain Layer (например, репозиториев) и Application Layer (например, для внешних сервисов). Взаимодействие с внешними системами: базы данных (PostgreSQL), поисковые движки (Elasticsearch), кэши (Redis), брокеры сообщений (Kafka). Также включает утилиты для логирования, конфигурации, метрик, трассировки.
*   **Ключевые компоненты/модули:** Адаптеры для PostgreSQL (например, с использованием GORM или Squirrel), адаптеры для Elasticsearch, адаптеры для Redis, продюсеры/консьюмеры Kafka, клиенты для других микросервисов (если Catalog Service их вызывает), утилиты для работы с конфигурацией (Viper), логированием (Zap), метриками (Prometheus), трассировкой (OpenTelemetry).

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL:** `/api/v1/catalog` (маршрутизируется через API Gateway).
*   **Формат данных:** JSON.
*   **Аутентификация:** Большинство публичных эндпоинтов для чтения данных доступны анонимно или требуют JWT с правами пользователя. Управляющие эндпоинты (например, с префиксом `/manage`) требуют JWT с правами администратора или менеджера каталога. Информация о пользователе (ID, роли, регион) передается API Gateway в заголовках (например, `X-User-Id`, `X-User-Roles`, `X-User-Region`).
*   **Авторизация:** На основе ролей (RBAC).
*   **Пагинация:** Используются query-параметры `page` (номер страницы, по умолчанию 1) и `limit` (количество элементов на странице, по умолчанию 20, максимум 100). Ответ содержит метаданные пагинации.
*   **Локализация:** Для локализованных полей (названия, описания) используется заголовок `Accept-Language` (например, `ru-RU, ru;q=0.9, en-US;q=0.8, en;q=0.7`). Если заголовок не указан, используется язык по умолчанию (`DEFAULT_LANGUAGE` из конфигурации).
*   **Стандартный формат ответа об ошибке:**
    ```json
    {
      "errors": [
        {
          "status": "4XX/5XX",
          "code": "ERROR_CODE_UPPER_SNAKE_CASE",
          "title": "Краткое описание ошибки на русском",
          "detail": "Полное описание ошибки с контекстом, если применимо.",
          "source": { "pointer": "/data/attributes/field_name", "parameter": "query_param_name" } // Опционально
        }
      ]
    }
    ```

#### 3.1.1. Ресурс: Продукты (Products)
*   **`GET /products`**
    *   Описание: Получение списка продуктов с возможностью фильтрации, сортировки и пагинации.
    *   Query параметры:
        *   `search` (string, опционально): Поисковый запрос по названию, описанию.
        *   `genre_ids` (string, опционально): ID жанров через запятую (например, `uuid1,uuid2`).
        *   `tag_ids` (string, опционально): ID тегов через запятую.
        *   `category_ids` (string, опционально): ID категорий через запятую.
        *   `platform` (string, опционально): Фильтр по платформе (`windows`, `linux`, `macos`, `steam_deck_verified`).
        *   `min_price` (integer, опционально): Минимальная цена.
        *   `max_price` (integer, опционально): Максимальная цена.
        *   `is_on_sale` (boolean, опционально): Только товары со скидкой.
        *   `sort_by` (string, опционально): Поле для сортировки (например, `name_asc`, `price_desc`, `release_date_desc`). По умолчанию `name_asc`.
        *   `page` (integer, опционально, default: 1).
        *   `limit` (integer, опционально, default: 20).
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "product",
              "id": "game-uuid-123",
              "attributes": {
                "title": "Супер Игра X",
                "short_description": "Захватывающее приключение в мире Y.",
                "cover_image_url": "https://cdn.example.com/covers/game-uuid-123.jpg",
                "current_price": { "amount": 1999, "currency": "RUB" },
                "release_date": "2023-10-26",
                "developer_name": "Крутые Разрабы"
              },
              "relationships": {
                "genres": { "data": [{ "type": "genre", "id": "genre-uuid-rpg" }] }
              },
              "links": { "self": "/api/v1/catalog/products/game-uuid-123" }
            }
            // ... другие продукты
          ],
          "meta": {
            "total_items": 150,
            "total_pages": 8,
            "current_page": 1,
            "per_page": 20
          },
          "links": {
            "self": "/api/v1/catalog/products?page=1&limit=20",
            "next": "/api/v1/catalog/products?page=2&limit=20"
          }
        }
        ```
    *   Требуемые права доступа: Публичный.
*   **`GET /products/{product_id}`**
    *   Описание: Получение детальной информации о продукте.
    *   Пример ответа (Успех 200 OK): (Более полная структура Product, см. раздел 4.1)
        ```json
        {
          "data": {
            "type": "product",
            "id": "game-uuid-123",
            "attributes": {
              // ... полные атрибуты продукта ...
              "title": "Супер Игра X",
              "description": "Полное описание игры...",
              "release_date": "2023-10-26",
              "developer_name": "Крутые Разрабы",
              "publisher_name": "Известный Издатель",
              "current_price": { "amount": 199900, "currency": "RUB", "formatted": "1999 руб." }, // Сумма в копейках/центах
              "base_price": { "amount": 199900, "currency": "RUB" },
              "discount_percentage": 0,
              "system_requirements": { /* ... */ },
              "media": [ /* ... массив MediaItem ... */ ],
              "achievements_count": 50
            },
            "relationships": { /* ... жанры, теги, категории ... */ },
            "links": { "self": "/api/v1/catalog/products/game-uuid-123" }
          }
        }
        ```
    *   Требуемые права доступа: Публичный.

#### 3.1.2. Ресурс: Жанры (Genres)
*   **`GET /genres`**
    *   Описание: Получение списка всех жанров.
    *   Query параметры: `page`, `limit`.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": [
            {
              "type": "genre",
              "id": "genre-uuid-rpg",
              "attributes": { "name": "Ролевые игры", "slug": "rpg" }
            },
            {
              "type": "genre",
              "id": "genre-uuid-strategy",
              "attributes": { "name": "Стратегии", "slug": "strategy" }
            }
          ],
          "meta": { "total_items": 25, "current_page": 1, "per_page": 20 }
        }
        ```
    *   Требуемые права доступа: Публичный.

#### 3.1.3. Управление Каталогом (Management API - пример)
*   **`POST /manage/products`**
    *   Описание: Создание нового продукта (административная функция).
    *   Тело запроса: (JSON с полными данными нового продукта, см. `Product` в разделе 4.1)
        ```json
        {
          "data": {
            "type": "productCreationRequest",
            "attributes": {
              "product_type": "game", // game, dlc, software, bundle
              "titles": { "ru-RU": "Новая Игра", "en-US": "New Game" },
              "descriptions": { "ru-RU": "Описание новой игры.", "en-US": "Description of the new game." },
              "developer_id": "dev-uuid-abc", // ID разработчика из Developer Service
              "publisher_id": "pub-uuid-xyz", // ID издателя
              // ... другие необходимые поля ...
              "initial_price": { "amount": 299900, "currency": "RUB" },
              "genre_ids": ["genre-uuid-rpg"]
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Возвращает созданный продукт)
    *   Требуемые права доступа: `catalog_admin`, `product_manager`.

### 3.2. gRPC API
*   Используется для внутреннего межсервисного взаимодействия.
*   Пакет: `catalog.v1`.
*   Определение Protobuf: `catalog/v1/catalog_internal_service.proto`.

#### 3.2.1. Сервис: `CatalogInternalService`
*   **`rpc GetProductInternal (GetProductInternalRequest) returns (GetProductInternalResponse)`**
    *   Описание: Получение детальной информации о продукте по его ID для внутреннего использования (может содержать больше полей, чем публичный REST API).
    *   `message GetProductInternalRequest { string product_id = 1; string language_code = 2; /* например, "ru-RU" */ string region_code = 3; /* например, "RU" */ }`
    *   `message ProductInternal { string id = 1; string type = 2; map<string, string> titles = 3; map<string, string> descriptions = 4; /* ... другие поля ... */ ProductPriceInternal current_price = 10; }`
    *   `message ProductPriceInternal { int64 amount = 1; string currency_code = 2; }`
    *   `message GetProductInternalResponse { ProductInternal product = 1; }`
*   **`rpc GetProductsInternal (GetProductsInternalRequest) returns (GetProductsInternalResponse)`**
    *   Описание: Получение списка продуктов по их ID.
    *   `message GetProductsInternalRequest { repeated string product_ids = 1; string language_code = 2; string region_code = 3; }`
    *   `message GetProductsInternalResponse { repeated ProductInternal products = 1; }`
*   **`rpc GetProductPrice (GetProductPriceRequest) returns (GetProductPriceResponse)`**
    *   Описание: Получение актуальной цены продукта для указанного региона и пользователя (если есть персональные скидки).
    *   `message GetProductPriceRequest { string product_id = 1; string user_id = 2; /* опционально */ string region_code = 3; string currency_code_override = 4; /* опционально */ }`
    *   `message GetProductPriceResponse { ProductPriceInternal price = 1; ProductPriceInternal base_price = 2; int32 discount_percentage = 3; }`

### 3.3. WebSocket API
*   Не применимо для данного сервиса.

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности

*   **`Product` (Продукт)**
    *   `id` (UUID): Уникальный идентификатор. Обязательность: Required.
    *   `product_type` (ENUM: `game`, `dlc`, `software`, `bundle`): Тип продукта. Обязательность: Required.
    *   `status` (ENUM: `draft`, `in_review`, `approved`, `rejected`, `published`, `unpublished`, `archived`): Статус продукта. Обязательность: Required.
    *   `titles` (JSONB/Map<String, String>): Локализованные названия (ключ - locale, например, `ru-RU`). Пример: `{"ru-RU": "Супер Игра X", "en-US": "Super Game X"}`. Обязательность: Required.
    *   `descriptions` (JSONB/Map<String, String>): Локализованные описания. Обязательность: Required.
    *   `release_date` (TIMESTAMPTZ): Дата релиза. Обязательность: Optional.
    *   `developer_ids` (ARRAY of UUID): ID разработчиков (ссылка на Developer Service или его сущность). Обязательность: Optional.
    *   `publisher_ids` (ARRAY of UUID): ID издателей. Обязательность: Optional.
    *   `system_requirements` (JSONB): Системные требования по платформам. Пример: `{"windows": {"minimum": "...", "recommended": "..."}, "linux": {"minimum": "..."}}`. Обязательность: Optional.
    *   `age_rating` (VARCHAR(10)): Возрастной рейтинг (например, `PEGI_18`, `ESRB_M`). Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ): Время создания. Обязательность: Required.
    *   `updated_at` (TIMESTAMPTZ): Время последнего обновления. Обязательность: Required.
    *   `average_rating` (FLOAT, 0-5): Средний рейтинг (обновляется из Social Service). Обязательность: Optional.
    *   `review_count` (INTEGER): Количество отзывов. Обязательность: Optional.
    *   `tags` (ARRAY of UUID, FK to Tags): Теги продукта.
    *   `genres` (ARRAY of UUID, FK to Genres): Жанры продукта.
    *   `categories` (ARRAY of UUID, FK to Categories): Категории продукта.
    *   `default_price_id` (UUID, FK to ProductPrices): Ссылка на текущую активную цену. Обязательность: Optional.

*   **`Genre` (Жанр)**
    *   `id` (UUID): Уникальный идентификатор. Обязательность: Required.
    *   `name` (JSONB/Map<String, String>): Локализованное имя жанра. Пример: `{"ru-RU": "Ролевая игра", "en-US": "Role-playing game"}`. Обязательность: Required.
    *   `slug` (VARCHAR(100)): Уникальный текстовый идентификатор (для URL). Пример: `rpg`. Валидация: unique, slug format. Обязательность: Required.
    *   `description` (JSONB/Map<String, String>): Локализованное описание. Обязательность: Optional.

*   **`ProductPrice` (Цена Продукта)**
    *   `id` (UUID): Уникальный идентификатор цены. Обязательность: Required.
    *   `product_id` (UUID, FK to Product): ID продукта. Обязательность: Required.
    *   `region_code` (VARCHAR(10)): Код региона (например, `RU`, `US`, `EU`, `GLOBAL`). `GLOBAL` для цены по умолчанию. Обязательность: Required.
    *   `currency_code` (VARCHAR(3)): Код валюты (например, `RUB`, `USD`, `EUR`). Обязательность: Required.
    *   `base_amount` (BIGINT): Базовая цена в минимальных единицах валюты (копейки, центы). Пример: `199900` (для 1999.00). Обязательность: Required.
    *   `discount_amount` (BIGINT): Сумма скидки (если есть). Обязательность: Optional, default 0.
    *   `effective_from` (TIMESTAMPTZ): Дата начала действия цены/скидки. Обязательность: Required.
    *   `effective_to` (TIMESTAMPTZ): Дата окончания действия цены/скидки. Обязательность: Optional.
    *   `is_active` (BOOLEAN): Является ли эта запись о цене текущей активной для продукта/региона. Обязательность: Required.

*   **`MediaItem` (Медиа-элемент)**
    *   `id` (UUID): Уникальный идентификатор. Обязательность: Required.
    *   `product_id` (UUID, FK to Product): ID продукта, к которому относится медиа. Обязательность: Required.
    *   `media_type` (ENUM: `screenshot`, `trailer_video`, `cover_art_small`, `cover_art_large`, `background_image`): Тип медиа. Обязательность: Required.
    *   `url` (VARCHAR(2048)): URL медиа-файла (на CDN). Валидация: valid URL. Обязательность: Required.
    *   `thumbnail_url` (VARCHAR(2048)): URL превью (для видео). Обязательность: Optional.
    *   `sort_order` (INTEGER): Порядок сортировки медиа-элементов. Обязательность: Optional.
    *   `metadata` (JSONB): Дополнительные метаданные (разрешение, длительность для видео). Обязательность: Optional.

*   **`AchievementMeta` (Метаданные Достижения)**
    *   `id` (UUID): Уникальный идентификатор. Обязательность: Required.
    *   `product_id` (UUID, FK to Product): ID игры, к которой относится достижение. Обязательность: Required.
    *   `achievement_api_name` (VARCHAR(100)): Уникальный API-идентификатор достижения (для интеграции с игрой). Валидация: unique per product. Обязательность: Required.
    *   `name` (JSONB/Map<String, String>): Локализованное название достижения. Обязательность: Required.
    *   `description` (JSONB/Map<String, String>): Локализованное описание. Обязательность: Required.
    *   `icon_url_unlocked` (VARCHAR(2048)): URL иконки разблокированного достижения. Обязательность: Required.
    *   `icon_url_locked` (VARCHAR(2048)): URL иконки заблокированного достижения. Обязательность: Required.
    *   `is_hidden` (BOOLEAN): Скрытое ли достижение до разблокировки. Обязательность: Required.
    *   `sort_order` (INTEGER): Порядок отображения. Обязательность: Optional.

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL

**ERD Диаграмма (ключевые таблицы):**
```mermaid
erDiagram
    PRODUCTS {
        UUID id PK
        VARCHAR product_type
        VARCHAR status
        JSONB titles
        JSONB descriptions
        TIMESTAMPTZ release_date
        ARRAY_UUID developer_ids
        ARRAY_UUID publisher_ids
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    GENRES {
        UUID id PK
        JSONB name UK
        VARCHAR slug UK
    }
    PRODUCT_GENRES {
        UUID product_id PK FK
        UUID genre_id PK FK
    }
    PRODUCT_PRICES {
        UUID id PK
        UUID product_id FK
        VARCHAR region_code
        VARCHAR currency_code
        BIGINT base_amount
        BIGINT discount_amount
        TIMESTAMPTZ effective_from
        TIMESTAMPTZ effective_to
        BOOLEAN is_active
    }
    MEDIA_ITEMS {
        UUID id PK
        UUID product_id FK
        VARCHAR media_type
        VARCHAR url
        VARCHAR thumbnail_url
        INTEGER sort_order
    }
    ACHIEVEMENT_METADATA {
        UUID id PK
        UUID product_id FK
        VARCHAR achievement_api_name UK_PerProduct
        JSONB name
        JSONB description
        VARCHAR icon_url_unlocked
        VARCHAR icon_url_locked
        BOOLEAN is_hidden
    }

    PRODUCTS ||--o{ PRODUCT_GENRES : "has"
    GENRES ||--o{ PRODUCT_GENRES : "belongs_to"
    PRODUCTS ||--o{ PRODUCT_PRICES : "has_prices"
    PRODUCTS ||--o{ MEDIA_ITEMS : "has_media"
    PRODUCTS ||--o{ ACHIEVEMENT_METADATA : "has_achievements"
```

**DDL (PostgreSQL - примеры ключевых таблиц):**
```sql
-- Расширение для UUID, если не создано
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Типы ENUM
CREATE TYPE product_type_enum AS ENUM ('game', 'dlc', 'software', 'bundle');
CREATE TYPE product_status_enum AS ENUM ('draft', 'in_review', 'approved', 'rejected', 'published', 'unpublished', 'archived');
CREATE TYPE media_type_enum AS ENUM ('screenshot', 'trailer_video', 'cover_art_small', 'cover_art_large', 'background_image');

-- Таблица продуктов
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_type product_type_enum NOT NULL,
    status product_status_enum NOT NULL DEFAULT 'draft',
    titles JSONB NOT NULL DEFAULT '{}'::jsonb, -- {"ru-RU": "Название", "en-US": "Title"}
    descriptions JSONB NOT NULL DEFAULT '{}'::jsonb,
    release_date TIMESTAMPTZ,
    -- developer_ids и publisher_ids могут быть массивами UUID или ссылками на отдельные таблицы/сервисы
    developer_ids UUID[],
    publisher_ids UUID[],
    system_requirements JSONB, -- {"windows": {"minimum": "...", "recommended": "..."}}
    age_rating VARCHAR(10),
    average_rating FLOAT,
    review_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_products_titles_gin ON products USING GIN (titles); -- Для поиска по названиям
CREATE INDEX idx_products_status ON products(status);

-- Таблица жанров
CREATE TABLE genres (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name JSONB NOT NULL DEFAULT '{}'::jsonb, -- {"ru-RU": "Название", "en-US": "Title"}
    slug VARCHAR(100) NOT NULL UNIQUE,
    description JSONB
);
CREATE INDEX idx_genres_name_gin ON genres USING GIN (name);

-- Связь продуктов и жанров (многие-ко-многим)
CREATE TABLE product_genres (
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    genre_id UUID NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, genre_id)
);

-- Таблица цен продуктов
CREATE TABLE product_prices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    region_code VARCHAR(10) NOT NULL DEFAULT 'GLOBAL', -- 'GLOBAL', 'RU', 'US', etc.
    currency_code VARCHAR(3) NOT NULL, -- 'RUB', 'USD', 'EUR'
    base_amount BIGINT NOT NULL, -- Цена в минимальных единицах валюты (копейки/центы)
    discount_amount BIGINT DEFAULT 0,
    effective_from TIMESTAMPTZ NOT NULL DEFAULT now(),
    effective_to TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- Для простой выборки текущей цены
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_product_prices_product_region_active ON product_prices(product_id, region_code, is_active);

-- Таблица медиа-элементов
CREATE TABLE media_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    media_type media_type_enum NOT NULL,
    url VARCHAR(2048) NOT NULL,
    thumbnail_url VARCHAR(2048),
    sort_order INTEGER DEFAULT 0,
    metadata JSONB, -- Например, разрешение для изображений, длительность для видео
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_media_items_product_id_type ON media_items(product_id, media_type);

-- TODO: Добавить DDL для tags, categories, product_tags, product_categories, achievement_metadata.
```

#### 4.2.2. Elasticsearch
*   **Роль:** Используется для обеспечения продвинутого полнотекстового поиска по продуктам, а также для сложной фильтрации и агрегации данных каталога.
*   **Индексы:** Основной индекс `products` содержит денормализованные данные о продуктах, включая названия, описания (на разных языках), жанры, теги, цены, рейтинг и т.д.
*   **Пример структуры документа в индексе `products`:**
    ```json
    {
      "id": "game-uuid-123",
      "product_type": "game",
      "status": "published",
      "titles": { "ru_RU": "Супер Игра X", "en_US": "Super Game X" }, // анализаторы для каждого языка
      "descriptions": { "ru_RU": "...", "en_US": "..." },
      "release_date": "2023-10-26T00:00:00Z",
      "developer_names": ["Крутые Разрабы"],
      "publisher_names": ["Известный Издатель"],
      "genres": ["rpg", "action"], // slug'и или ID
      "tags": ["open_world", "fantasy"],
      "platforms": ["windows", "steam_deck_verified"],
      "price_rub": 1999.00, // Актуальная цена в основной валюте для сортировки/фильтрации
      "average_rating": 4.75,
      "review_count": 125
      // ... другие поля для поиска и фильтрации
    }
    ```

#### 4.2.3. Redis
*   **Роль:** Кэширование часто запрашиваемых данных для снижения нагрузки на PostgreSQL и Elasticsearch и ускорения ответов API.
*   **Типы данных и ключи (примеры):**
    *   Детали продукта: `product:<product_id>:<lang>:<region>` (JSON строка или HASH). TTL: минуты/часы.
    *   Списки продуктов (например, главная страница, списки по жанрам): `products_list:<query_hash_or_page_params>`. TTL: минуты.
    *   Списки жанров, тегов, категорий: `genres_list:<lang>`, `tags_list:<lang>`. TTL: часы/день.
    *   Цены: `price:<product_id>:<region_code>`. TTL: минуты/часы.

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Apache Kafka.
*   **Формат событий:** CloudEvents v1.0 (JSON encoding), согласно `project_api_standards.md`.
*   **Топики (примеры):** `catalog.product.v1`, `catalog.price.v1`, `catalog.taxonomy.v1`. (Могут быть более гранулярными).
*   **Ключ сообщения Kafka:** Обычно ID сущности (например, `product_id`).

*   **`catalog.product.created.v1`**
    *   Описание: Новый продукт был создан.
    *   Пример Payload (`data` секция CloudEvent):
        ```json
        {
          "product_id": "game-uuid-abc",
          "product_type": "game",
          "titles": { "ru-RU": "Новая Игра", "en-US": "New Game" },
          "status": "draft",
          "created_by": "user-uuid-admin" // ID пользователя или сервиса, создавшего продукт
        }
        ```
*   **`catalog.product.updated.v1`**
    *   Описание: Метаданные продукта были обновлены.
    *   Пример Payload:
        ```json
        {
          "product_id": "game-uuid-abc",
          "updated_fields": ["titles", "descriptions", "status"], // Список измененных полей
          "new_status": "in_review",
          "updated_by": "user-uuid-developer"
        }
        ```
*   **`catalog.price.updated.v1`**
    *   Описание: Цена на продукт была обновлена (добавлена новая цена, изменена существующая, установлена скидка).
    *   Пример Payload:
        ```json
        {
          "product_id": "game-uuid-abc",
          "price_id": "price-uuid-xyz",
          "region_code": "RU",
          "currency_code": "RUB",
          "base_amount": 299900,
          "discount_amount": 50000, // 0 если скидки нет
          "effective_from": "2024-04-01T00:00:00Z",
          "effective_to": "2024-04-15T23:59:59Z" // null если бессрочно
        }
        ```
*   **Другие события:** `catalog.product.status.changed.v1`, `catalog.discount.started.v1` / `catalog.discount.ended.v1`, `catalog.genre.created.v1`, `catalog.achievement_meta.updated.v1`.
    *   TODO: Детализировать структуру Payload для этих событий по мере необходимости.

### 5.2. Потребляемые События (Consumed Events)

*   **`social.review.stats.updated.v1`** (от Social Service)
    *   Описание: Обновлена агрегированная статистика по отзывам для продукта (средний рейтинг, количество отзывов).
    *   Ожидаемый Payload (пример):
        ```json
        {
          "product_id": "game-uuid-abc",
          "average_rating": 4.78,
          "total_reviews_count": 152,
          "ratings_breakdown": { "1": 5, "2": 10, "3": 25, "4": 50, "5": 62 }
        }
        ```
    *   Логика обработки: Обновить поля `average_rating` и `review_count` для продукта в PostgreSQL. Обновить соответствующий документ в Elasticsearch. Инвалидировать кэш для данного продукта.
*   **`user.preference.changed.v1`** (гипотетическое, от User Profile Service или Personalization Service)
    *   Описание: Предпочтения пользователя (например, любимые жанры) изменились.
    *   Ожидаемый Payload (пример):
        ```json
        {
          "user_id": "user-uuid-xyz",
          "updated_preferences": {
            "favorite_genres": ["genre-uuid-rpg", "genre-uuid-strategy"],
            "ignored_tags": ["tag-uuid-horror"]
          }
        }
        ```
    *   Логика обработки: Это событие может быть не напрямую потреблено Catalog Service, а скорее системой рекомендаций, которая использует данные каталога. Если же Catalog Service предоставляет API персонализированных подборок, он может использовать это для обновления кэшей таких подборок.
*   **`moderation.content.approved.v1`** / **`moderation.content.rejected.v1`** (от Admin Service)
    *   Описание: Контент продукта прошел модерацию.
    *   Ожидаемый Payload (пример):
        ```json
        {
          "content_type": "product", // "product_description", "media_item", etc.
          "content_id": "game-uuid-abc", // ID продукта или медиа-элемента
          "moderator_id": "admin-uuid-123",
          "decision": "approved", // "rejected"
          "reason": "Соответствует правилам" // или причина отклонения
        }
        ```
    *   Логика обработки: Обновить поле `status` у продукта или соответствующего элемента. Если `approved`, продукт может стать видимым для пользователей. Опубликовать событие `catalog.product.status.changed.v1`.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **API Gateway**: Проксирование всех входящих REST API запросов, первичная аутентификация/авторизация.
*   **Developer Service**: Получение информации о разработчиках/издателях продуктов. Может вызывать API Catalog Service для создания/обновления продуктов от имени разработчика.
*   **Admin Service**: Управление статусами модерации продуктов. Может вызывать API Catalog Service для редактирования любых данных каталога.
*   **Payment Service**: Получение актуальных цен на продукты для формирования заказов. Catalog Service может уведомлять Payment Service об изменениях цен через Kafka.
*   **Library Service**: Получение метаданных продуктов и информации о достижениях для отображения в библиотеках пользователей.
*   **Download Service**: Получение информации о продуктах для предоставления файлов для скачивания.
*   **Analytics Service**: Catalog Service публикует события об изменениях (создание, обновление продуктов, цен), которые потребляются Analytics Service. Может получать данные (например, популярность товаров) от Analytics Service для сортировок или рекомендаций.
*   **Notification Service**: Инициирование уведомлений пользователям (например, о выходе игры из списка желаемого, о старте скидки на отслеживаемый продукт) через публикацию событий в Kafka, которые потребляет Notification Service.
*   **Auth Service**: Валидация JWT токенов и проверка прав доступа (обычно выполняется на уровне API Gateway, но может быть дополнительная проверка в Catalog Service для критичных операций).
*   **Search Service (если выделен)**: Catalog Service предоставляет данные для индексации или Search Service потребляет события Kafka от Catalog Service для обновления своих индексов. Если Elasticsearch используется напрямую, то это часть Catalog Service.
*   **Recommendation Service (если выделен)**: Предоставление данных о продуктах, их связях и пользовательских взаимодействиях (косвенно, через Analytics) для генерации рекомендаций.

### 6.2. Внешние Системы
*   **CDN (Content Delivery Network)**: Для хранения и быстрой доставки медиа-контента (изображения, видео). Catalog Service хранит URL-ы на ресурсы в CDN.
*   **S3-совместимое хранилище**: Может использоваться как первичное хранилище для загружаемых медиа-файлов перед их обработкой и передачей в CDN.

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `CATALOG_HTTP_PORT`: Порт для REST API (например, `8080`).
*   `CATALOG_GRPC_PORT`: Порт для gRPC API (например, `9090`).
*   `POSTGRES_DSN`: Строка подключения к PostgreSQL.
*   `ELASTICSEARCH_URLS`: URL(ы) Elasticsearch (через запятую).
*   `REDIS_ADDR`: Адрес Redis (например, `redis-master:6379`).
*   `REDIS_PASSWORD`: Пароль для Redis (если есть).
*   `REDIS_DB_CATALOG`: Номер базы Redis для Catalog Service (например, `0`).
*   `KAFKA_BROKERS`: Список брокеров Kafka.
*   `KAFKA_TOPIC_PRODUCT_EVENTS`: Топик для событий продуктов.
*   `KAFKA_TOPIC_PRICE_EVENTS`: Топик для событий цен.
*   `CDN_BASE_URL`: Базовый URL для медиа-контента на CDN.
*   `S3_ENDPOINT`, `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY`, `S3_BUCKET_MEDIA_RAW`, `S3_USE_SSL`: Параметры S3.
*   `LOG_LEVEL`: Уровень логирования (`debug`, `info`, `warn`, `error`).
*   `AUTH_SERVICE_GRPC_ADDR`: Адрес gRPC Auth Service.
*   `DEFAULT_LANGUAGE`: Язык по умолчанию (например, `ru-RU`).
*   `DEFAULT_REGION_CODE`: Код региона по умолчанию (например, `RU`).
*   `DEFAULT_CURRENCY_CODE`: Код валюты по умолчанию (например, `RUB`).
*   `CACHE_PRODUCT_DETAILS_TTL_SECONDS`: TTL для кэша деталей продукта.
*   `CACHE_GENRE_LIST_TTL_SECONDS`: TTL для кэша списка жанров.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`: Эндпоинт Jaeger для экспорта трейсов.

### 7.2. Файлы Конфигурации (если применимо)
*   **`configs/catalog_config.yaml`**: Может использоваться для задания структуры настроек, которые затем переопределяются переменными окружения.
    ```yaml
    server:
      http_port: ${CATALOG_HTTP_PORT:-8080}
      grpc_port: ${CATALOG_GRPC_PORT:-9090}
      read_timeout_seconds: 15
      write_timeout_seconds: 15

    database:
      postgres_dsn: ${POSTGRES_DSN}
      max_open_conns: 100
      max_idle_conns: 25
      conn_max_lifetime_seconds: 3600

    elasticsearch:
      urls: ${ELASTICSEARCH_URLS} # "http://es1:9200,http://es2:9200"
      username: ${ELASTICSEARCH_USERNAME:-}
      password: ${ELASTICSEARCH_PASSWORD:-}
      default_index_prefix: "catalog_"

    redis:
      address: ${REDIS_ADDR}
      password: ${REDIS_PASSWORD:-""}
      db: ${REDIS_DB_CATALOG:-0}

    kafka:
      brokers: ${KAFKA_BROKERS} # "kafka1:9092,kafka2:9092"
      topics:
        product_events: ${KAFKA_TOPIC_PRODUCT_EVENTS:-catalog.product.v1}
        price_events: ${KAFKA_TOPIC_PRICE_EVENTS:-catalog.price.v1}

    localization:
      default_language: ${DEFAULT_LANGUAGE:-ru-RU}
      supported_languages: ["ru-RU", "en-US"]

    pricing:
      default_region_code: ${DEFAULT_REGION_CODE:-RU}
      default_currency_code: ${DEFAULT_CURRENCY_CODE:-RUB}

    cache_settings:
      product_details_ttl_seconds: ${CACHE_PRODUCT_DETAILS_TTL_SECONDS:-300}
      genre_list_ttl_seconds: ${CACHE_GENRE_LIST_TTL_SECONDS:-3600}
      # ... другие TTL для кэшей

    # Настройки для полнотекстового поиска (могут быть специфичны для языка)
    search_settings:
      default_fuzziness: "AUTO"
      min_match_percentage: "75%"
      language_analyzers:
        ru: "russian_custom_analyzer" # Имя анализатора в Elasticsearch
        en: "english_custom_analyzer"
    ```

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   REST API: Используются стандартные коды состояния HTTP. Тело ответа об ошибке соответствует формату, определенному в `project_api_standards.md` (см. секцию 3.1).
*   gRPC API: Используются стандартные коды состояния gRPC. Дополнительная информация об ошибке передается через `google.rpc.Status` и `google.rpc.ErrorInfo`.
*   Все ошибки логируются с `trace_id` (из OpenTelemetry) и, если применимо, `request_id`.

### 8.2. Распространенные Коды Ошибок (для REST API)
*   **`400 Bad Request` (`INVALID_ARGUMENT`)**: Некорректные входные данные (например, неверный формат ID, отсутствуют обязательные поля в теле запроса, ошибка валидации).
*   **`401 Unauthorized` (`UNAUTHENTICATED`)**: Ошибка аутентификации (например, недействительный или отсутствующий JWT токен). Обычно обрабатывается на уровне API Gateway.
*   **`403 Forbidden` (`PERMISSION_DENIED`)**: Недостаточно прав для выполнения операции (например, попытка редактировать каталог без роли `catalog_admin`).
*   **`404 Not Found` (`RESOURCE_NOT_FOUND`)**: Запрашиваемый ресурс не найден (продукт, жанр, цена и т.д.).
*   **`409 Conflict` (`ALREADY_EXISTS`)**: Попытка создания ресурса, который уже существует с конфликтующими уникальными полями (например, продукт с тем же `developer_product_id` для данного разработчика).
*   **`422 Unprocessable Entity` (`VALIDATION_ERROR_BUSINESS_LOGIC`)**: Запрос корректен синтаксически, но нарушает бизнес-правила (например, установка скидки больше базовой цены).
*   **`500 Internal Server Error` (`INTERNAL_ERROR`)**: Внутренняя ошибка сервера (например, ошибка при работе с БД, непредвиденное исключение).
*   **`503 Service Unavailable` (`SERVICE_UNAVAILABLE`)**: Сервис временно недоступен или одна из его критических зависимостей (БД, Kafka) недоступна.

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   **Внешние запросы (REST API):** Проверка JWT токенов, полученных от Auth Service, выполняется на уровне API Gateway. Catalog Service доверяет информации о пользователе, переданной из API Gateway в заголовках (`X-User-Id`, `X-User-Roles`).
*   **Межсервисные запросы (gRPC):** Используется mTLS для установления защищенного соединения. Дополнительно могут использоваться API-ключи или JWT токены с сервисными ролями, передаваемые в метаданных gRPC запроса и валидируемые через Auth Service.

### 9.2. Авторизация
*   **Модель:** Role-Based Access Control (RBAC).
*   **Проверка прав:** На основе ролей пользователя, полученных из заголовка `X-User-Roles` (для REST API) или из контекста gRPC запроса. Для управляющих эндпоинтов (`/manage/*`) требуются специфичные роли (например, `catalog_admin`, `product_manager`, `price_manager`). Для публичных эндпоинтов проверка прав может не требоваться или ограничиваться видимостью определенных полей.
*   В некоторых случаях может применяться проверка владения ресурсом (например, разработчик может редактировать только свои продукты).

### 9.3. Защита Данных
*   **Шифрование:**
    *   TLS для всех внешних и внутренних коммуникаций (HTTPS, gRPCs).
    *   Шифрование дисков (at-rest encryption) для баз данных (PostgreSQL, Elasticsearch, Redis) и хранилища Kafka.
*   **Обработка персональных данных:** Catalog Service может хранить ID разработчиков/издателей, которые могут быть связаны с персональными данными в других сервисах. Прямых ПДн пользователей обычно не хранит, кроме, возможно, региональных предпочтений для цен, если это не получается из других источников. Необходимо соблюдать ФЗ-152.
*   **Защита от уязвимостей:**
    *   Предотвращение SQL-инъекций (использование ORM/параметризованных запросов).
    *   Валидация и санитизация всех входных данных для предотвращения XSS (хотя сервис в основном API, но данные могут отображаться в админ-панелях) и других атак на ввод данных.
    *   Ограничение скорости запросов (Rate Limiting) на уровне API Gateway.

### 9.4. Управление Секретами
*   Пароли к базам данных, ключи для Kafka, секреты для Elasticsearch и другие чувствительные данные конфигурации должны храниться в безопасном хранилище секретов (например, HashiCorp Vault или зашифрованные Kubernetes Secrets) и доставляться в сервис во время выполнения.

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Многоэтапная сборка (multi-stage build) на основе официального образа Go (для сборки) и легковесного образа (например, Alpine Linux) для runtime.
*   **Kubernetes манифесты/Helm-чарты:** Включают Deployment, Service, ConfigMap, Secret, HorizontalPodAutoscaler (HPA), PodDisruptionBudget (PDB), NetworkPolicy для управления развертыванием и сетевым доступом.
*   (Ссылка на `project_deployment_standards.md` и репозиторий GitOps).

### 10.2. Зависимости при Развертывании
*   Кластер Kubernetes.
*   PostgreSQL (может быть развернут как StatefulSet или использоваться управляемый сервис).
*   Elasticsearch (развертывается как кластер).
*   Redis (может быть кластером или отдельным экземпляром).
*   Apache Kafka.
*   Доступность API Gateway для маршрутизации запросов.
*   Доступность Auth Service для валидации токенов (если не полностью делегировано API Gateway).

### 10.3. CI/CD
*   Автоматизированные пайплайны (например, GitLab CI, Jenkins, GitHub Actions) для:
    *   Сборки бинарного файла сервиса.
    *   Запуска юнит-тестов и интеграционных тестов (с использованием тестовых БД/Elasticsearch/Redis).
    *   Статического анализа кода (SAST).
    *   Сборки Docker-образа и его публикации в приватный registry.
    *   Развертывания в различные окружения (dev, staging, production) с использованием Helm-чартов и GitOps-подхода (ArgoCD/Flux).
    *   Запуска E2E тестов после развертывания.
*   **Операционные процедуры:** Документированные процедуры для масштабирования (HPA настроен, но ручное вмешательство может понадобиться для БД), обновления/отката версий, миграции схемы БД (например, с использованием Alembic, Flyway или встроенных средств GORM), переиндексации Elasticsearch, диагностики проблем (траблшутинг), резервного копирования и восстановления данных.

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат:** Структурированные логи в формате JSON (с использованием библиотеки Zap). Обязательные поля: `timestamp`, `level` (DEBUG, INFO, WARN, ERROR), `service_name` ("catalog-service"), `version` (версия сервиса), `trace_id`, `span_id` (из OpenTelemetry), `message`, `caller` (файл и строка кода). Дополнительные контекстные поля (например, `product_id`, `user_id`, `error_details`).
*   **Интеграция:** Сбор логов через FluentBit/Vector и отправка в централизованную систему логирования (например, Loki, ELK Stack).
*   (Ссылка на `project_observability_standards.md`).

### 11.2. Мониторинг
*   **Метрики (Prometheus):** Экспортируются через эндпоинт `/metrics`.
    *   Количество запросов (gRPC/HTTP) с разделением по методу/пути/сервису и коду ответа: `catalog_http_requests_total`, `catalog_grpc_requests_total`.
    *   Длительность обработки запросов: `catalog_http_request_duration_seconds`, `catalog_grpc_request_duration_seconds` (гистограммы).
    *   Производительность Go runtime: стандартные метрики (goroutines, GC, heap size).
    *   Производительность и ошибки при работе с PostgreSQL, Elasticsearch, Redis, Kafka.
    *   Попадание в кэш (cache hit/miss rate) для Redis.
    *   Количество опубликованных/потребленных сообщений Kafka, задержки.
    *   Бизнес-метрики: количество продуктов в каталоге, количество поисковых запросов.
*   **Дашборды (Grafana):** Настроенные дашборды для визуализации состояния сервиса, его производительности, использования ресурсов, состояния зависимостей, ключевых бизнес-метрик.
*   **Алертинг (AlertManager):** Настроены алерты для критических ситуаций: высокий процент ошибок (>5% за 5 мин), значительное увеличение времени ответа (P99 > 500ms), ошибки подключения к БД/Elasticsearch/Redis/Kafka, переполнение очередей Kafka, низкий cache hit rate (<80%).
*   (Ссылка на `project_observability_standards.md`).

### 11.3. Трассировка
*   **Инструментация:** Используется OpenTelemetry SDK для Go.
*   **Создание спанов:** Для всех входящих REST и gRPC запросов, для вызовов репозиториев PostgreSQL, запросов к Elasticsearch, Redis, публикации сообщений в Kafka, и для исходящих gRPC/HTTP вызовов к другим сервисам.
*   **Контекст трассировки:** Автоматическая и ручная пропагация W3C Trace Context.
*   **Экспорт:** Трейсы экспортируются в Jaeger или Tempo.
*   (Ссылка на `project_observability_standards.md`).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность (Latency):**
    *   API чтения (например, `GET /products/{id}`, `GET /products` с фильтрами): P95 < 150 мс, P99 < 300 мс.
    *   API поиска (Elasticsearch): P95 < 200 мс для типичных запросов.
    *   API записи (например, `POST /manage/products`): P95 < 500 мс (без учета асинхронных операций, таких как индексация).
*   **Производительность (Throughput):**
    *   Способность обрабатывать не менее 1000 запросов в секунду (RPS) на чтение (публичные API).
    *   Способность обрабатывать не менее 100 RPS на запись (управляющие API).
*   **Масштабируемость:** Горизонтальное масштабирование для увеличения пропускной способности. Способность управлять каталогом до 100,000 продуктов с миллионами связанных сущностей (цены, медиа).
*   **Доступность:** >= 99.95% (время простоя не более ~22 минут в месяц).
*   **Консистентность данных:** Strong consistency для операций записи в PostgreSQL. Eventual consistency для данных в Elasticsearch и Redis кэше (задержка репликации/обновления кэша < 1 минуты).
*   **Надежность:** Отсутствие потери данных при сбоях. Устойчивость к сбоям зависимостей (деградация функциональности, но не полный отказ, где возможно).
*   **Безопасность:** Соответствие требованиям `project_security_standards.md`.
*   **Сопровождаемость:** Покрытие кода тестами > 80%. Четкая структура кода. Актуальная документация.

## 13. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API и Protobuf определения для gRPC API будут храниться в соответствующих репозиториях или артефактах CI/CD.
*   Полные DDL схемы базы данных и конфигурации Elasticsearch будут поддерживаться в актуальном состоянии в системе миграций и в GitOps репозитории конфигураций.
*   TODO: Добавить ссылки на репозиторий с OpenAPI/Protobuf и на систему управления миграциями БД, когда они будут определены.

---
*Этот документ является основной спецификацией для Catalog Service и должен поддерживаться в актуальном состоянии.*

## 14. Связанные Рабочие Процессы (Related Workflows)
*   [Game Purchase and Library Update](../../../project_workflows/game_purchase_flow.md)
*   [Developer Submits a New Game for Moderation](../../../project_workflows/game_submission_flow.md)
