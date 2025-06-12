# Спецификация Микросервиса: Catalog Service

**Версия:** 1.0 (на основе исходной v3.0)
**Дата последнего обновления:** 2025-05-25

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой спецификацию микросервиса Catalog Service, являющегося ядром платформы и отвечающего за всю информацию о цифровых продуктах.
*   **Роль в общей архитектуре платформы:** Централизованное управление каталогом продуктов (игры, DLC, ПО, комплекты), ценообразованием, акциями, таксономией, медиа-контентом и метаданными достижений. Предоставляет данные другим микросервисам и обеспечивает поиск и обнаружение продуктов пользователями.
*   **Основные бизнес-задачи:**
    *   Управление жизненным циклом и метаданными продуктов.
    *   Управление ценами, скидками, акциями.
    *   Структурирование каталога (жанры, теги, категории, коллекции, франшизы).
    *   Обеспечение поиска, фильтрации и навигации по каталогу.
    *   Управление медиа-контентом и метаданными достижений.
    *   Поддержка процесса модерации контента.
    *   Предоставление данных для системы рекомендаций.
*   (Источник: Каталог Спец. v3.0, разделы 1, 2.1, 2.2)

### 1.2. Ключевые Функциональности
*   CRUD операции для продуктов.
*   Управление метаданными продуктов (названия, описания, системные требования, локализация).
*   Управление ценообразованием (базовые/региональные цены, скидки, акции).
*   Управление таксономией (жанры, теги, категории, франшизы, коллекции).
*   Полнотекстовый поиск, фильтрация, сортировка, пагинация.
*   Управление медиа-контентом (скриншоты, трейлеры, арты, обложки).
*   Управление метаданными достижений.
*   Хранение статусов модерации.
*   Предоставление API для рекомендаций.
*   (Источник: Каталог Спец. v3.0, раздел 3.3)

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+).
*   **Веб-фреймворк (REST):** Echo (v4+).
*   **RPC фреймворк (gRPC):** `google.golang.org/grpc`.
*   **База данных:** PostgreSQL (версия 15+).
*   **Поисковый движок:** Elasticsearch (версия 8.x+).
*   **Кэширование:** Redis (версия 7.0+).
*   **Брокер сообщений:** Apache Kafka (версия 3.x+).
*   **Конфигурация:** Viper.
*   **Логирование:** Zap.
*   **Мониторинг/Трассировка:** OpenTelemetry, Prometheus, Grafana, Jaeger.
*   **Инфраструктура:** Docker, Kubernetes.
*   (Источник: Каталог Спец. v3.0, раздел 4.4)

### 1.4. Термины и Определения (Glossary)
*   **Продукт (Product):** Любой цифровой товар (Игра, DLC, ПО, Комплект).
*   **Метаданные (Metadata):** Описательная информация о Продукте.
*   **Таксономия (Taxonomy):** Система классификации (Жанры, Теги, Категории).
*   **Цена (Price):** Стоимость Продукта, включая региональные цены и скидки.
*   (Полный глоссарий см. в Каталог Спец. v3.0, раздел 1, или в "Едином глоссарии терминов...")

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Catalog Service реализуется как независимый микросервис, следуя принципам **Чистой Архитектуры (Clean Architecture)** и элементам **CQRS (Command Query Responsibility Segregation)**.
*   Это обеспечивает разделение ответственностей, тестируемость и оптимизацию операций чтения и записи. Запросы на чтение могут использовать денормализованные данные из Elasticsearch/кэша, команды работают с основной моделью в PostgreSQL.
*   (Источник: Каталог Спец. v3.0, раздел 4.1)

### 2.2. Слои Сервиса

#### 2.2.1. Presentation Layer (Слой Представления / Transport Adapters)
*   Ответственность: Обработка входящих REST (Echo) и gRPC запросов, валидация, вызов Application Layer.
*   Ключевые компоненты/модули: HTTP хендлеры, gRPC серверы, обработчики сообщений Kafka (для входящих команд, если есть).
*   (Источник: Каталог Спец. v3.0, раздел 4.1.1 - Infrastructure Layer / Entrypoints)

#### 2.2.2. Application Layer (Прикладной Слой)
*   Ответственность: Реализация сценариев использования (use cases), координация взаимодействия. Содержит команды, запросы и их обработчики.
*   Ключевые компоненты/модули: Use Case Services (например, `CatalogProductService`), Command Handlers (`CreateProductCommandHandler`), Query Handlers (`GetProductQueryHandler`), DTOs.
*   (Источник: Каталог Спец. v3.0, раздел 4.1.1 - Application Layer)

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности (Product, Genre, Tag, PriceRule), агрегаты, доменные события, бизнес-правила.
*   Ключевые компоненты/модули: Entities (Product, Genre, Tag, Price, MediaItem, AchievementMeta), Value Objects (LocalizedString, SystemRequirements), Domain Services, Repository Interfaces.
*   (Источник: Каталог Спец. v3.0, раздел 4.1.1 - Domain Layer)

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Реализация интерфейсов для взаимодействия с PostgreSQL (GORM/Squirrel), Elasticsearch, Redis, Kafka (segmentio/kafka-go или confluent-kafka-go), клиенты для других сервисов. Управление конфигурацией (Viper), логирование (Zap), метрики (Prometheus), трассировка (OpenTelemetry).
*   Ключевые компоненты/модули: Data Persistence Adapters (PostgreSQL, Elasticsearch), Cache Adapters (Redis), Event Bus Adapters (Kafka), External Service Clients, Frameworks & Libraries.
*   (Источник: Каталог Спец. v3.0, раздел 4.1.1 - Infrastructure Layer)

*   (Детальную структуру проекта см. в Каталог Спец. v3.0, раздел 4.1.2)

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL:** `/api/v1/catalog` (устанавливается API Gateway).
*   **Формат:** JSON.
*   **Аутентификация:** JWT (проверяется API Gateway, информация о пользователе в заголовках `X-User-Id`, `X-User-Roles`, `X-User-Region`).
*   **Авторизация:** RBAC.
*   **Пагинация:** `page`, `limit`.
*   **Локализация:** Заголовок `Accept-Language`.
*   **Основные эндпоинты:**
    *   `GET /products`: Получение списка продуктов с фильтрацией и сортировкой.
    *   `GET /products/{product_id}`: Получение деталей продукта.
    *   `GET /genres`: Список жанров.
    *   `GET /tags`: Список тегов.
    *   `GET /categories`: Список категорий.
    *   `GET /products/{product_id}/achievements`: Метаданные достижений.
    *   `GET /recommendations`: Получение рекомендаций.
    *   Управляющие эндпоинты (префикс `/manage`): `POST /manage/products`, `PUT /manage/products/{product_id}` и т.д. (требуют соответствующих прав).
*   (Полный список и детали см. в Каталог Спец. v3.0, раздел 6.3)

### 3.2. gRPC API
*   Используется для внутреннего взаимодействия. Пакет `catalog.v1`.
*   **Основные сервисы и методы:**
    *   `CatalogInternalService`:
        *   `GetProductInternal`: Получить детали продукта по ID.
        *   `GetProductsInternal`: Получить список продуктов по ID.
        *   `GetProductPrice`: Получить актуальную цену продукта.
        *   `CreateProduct`, `UpdateProduct`: Управление продуктами.
        *   `GetAchievementMetadata`.
*   (Полный Protobuf и детали см. в Каталог Спец. v3.0, раздел 6.4)

### 3.3. WebSocket API (если применимо)
*   Information not found in existing documentation. (Не ожидается для Catalog Service).

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности
*   **Product**: Продукт (Игра, DLC, ПО, Комплект) со всеми метаданными.
*   **Genre**: Жанр.
*   **Tag**: Тег.
*   **Category**: Категория.
*   **MediaItem**: Элемент медиа-контента (скриншот, трейлер, обложка).
*   **ProductPrice**: Цена продукта (региональная, со скидками).
*   **AchievementMeta**: Метаданные достижения.
*   (Go структуры см. в Каталог Спец. v3.0, раздел 6.1)

### 4.2. Схема Базы Данных (PostgreSQL)
*   Основные таблицы: `products`, `genres`, `tags`, `categories`, `product_genres`, `product_tags`, `product_categories`, `media_items`, `product_prices`, `achievement_metadata`.
    ```sql
    -- Пример таблицы products (сокращенно)
    CREATE TABLE IF NOT EXISTS products (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        type product_type NOT NULL, -- ENUM('game', 'dlc', 'software', 'bundle')
        status product_status NOT NULL DEFAULT 'draft', -- ENUM('draft', 'in_review', ...)
        titles JSONB NOT NULL DEFAULT '{}'::jsonb,
        -- ... другие поля ...
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
*   Используются индексы, включая GIN для JSONB и полнотекстового поиска.
*   (Полную DDL см. в Каталог Спец. v3.0, раздел 6.2)

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Kafka.
*   **Формат событий:** CloudEvents v1.0 (JSON encoding).
*   **Топики:** `catalog.<entity>.<verb>` (например, `catalog.product.created`).
*   **Ключ сообщения:** ID сущности (например, `product_id`).
*   **Основные публикуемые события:**
    *   `catalog.product.created`
    *   `catalog.product.updated`
    *   `catalog.product.status.changed`
    *   `catalog.price.updated`
    *   `catalog.discount.started` / `ended`
    *   `catalog.genre.created` / `updated` / `deleted`
    *   `catalog.tag.created` / `updated` / `deleted`
    *   `catalog.achievement.metadata.updated`
*   (Пример события и детали см. в Каталог Спец. v3.0, раздел 6.5)

### 5.2. Потребляемые События (Consumed Events)
Взаимодействие с такими сервисами, как Developer Service (при подаче нового продукта или обновлении) и Admin Service (при изменении статуса модерации), как правило, осуществляется через прямые API вызовы (gRPC/REST) к Catalog Service, а не через потребление им событий от этих сервисов. Однако, Catalog Service может потреблять следующие типы событий для обновления агрегированных или связанных данных:

*   **`social.review.stats.updated.v1`** (название предложено, исходя из логики)
    *   **Источник:** Social Service
    *   **Назначение:** Обновление агрегированной информации о рейтингах и количестве отзывов для продукта в каталоге.
    *   **Структура Payload (пример):**
        ```json
        {
          "product_id": "uuid_game_or_dlc",
          "average_rating": 4.75,
          "review_count": 125,
          "updated_at": "ISO8601_timestamp"
        }
        ```
    *   **Логика обработки:** Обновить поля `average_rating` и `review_count` для соответствующего продукта в базе данных Catalog Service. Инвалидировать кэш для данного продукта.

*   **`user.preference.changed.v1`** (гипотетическое событие для системы рекомендаций)
    *   **Источник:** User Profile Service / Personalization Service
    *   **Назначение:** Уведомление об изменении предпочтений пользователя, что может повлиять на данные, используемые для формирования рекомендаций в каталоге.
    *   **Структура Payload (пример):**
        ```json
        {
          "user_id": "uuid_user",
          "changed_preferences": {
            "preferred_genres": ["rpg", "strategy"],
            "ignored_tags": ["horror"]
          },
          "timestamp": "ISO8601_timestamp"
        }
        ```
    *   **Логика обработки:** Передать информацию в модуль рекомендаций для возможного обновления персонализированных предложений или пересчета весов для пользователя.

*Примечание: Специфичные события, потребляемые сервисом Каталога, будут окончательно детализированы по мере проработки интеграционных сценариев и дизайна системы рекомендаций. Основной поток создания и модерации контента инициируется через API Catalog Service.*

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **API Gateway**: Проксирование запросов, аутентификация/авторизация.
*   **Developer Service**: Управление продуктами (CRUD).
*   **Admin Service**: Модерация, управление глобальными настройками каталога.
*   **Payment Service**: Получение актуальных цен, уведомления об изменении цен.
*   **Library Service**: Получение метаданных продуктов и достижений.
*   **Download Service**: Получение метаданных для загрузок.
*   **Analytics Service**: Отправка событий об изменениях в каталоге; получение данных для рекомендаций.
*   **Notification Service**: Инициирование уведомлений (выход игры из wish-list, старт скидки).
*   **Auth Service**: Валидация JWT, проверка прав (если не делегировано API Gateway).
*   (Детали по протоколам и API см. в Каталог Спец. v3.0, разделы 2.3 и 7)

### 6.2. Внешние Системы
*   **CDN**: Хранение и доставка медиа-контента. Catalog Service хранит URL.
*   **S3-совместимое хранилище**: Первичное хранение медиа-файлов.
*   (Источник: Каталог Спец. v3.0, раздел 7)

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `CATALOG_HTTP_PORT`, `CATALOG_GRPC_PORT`
*   `POSTGRES_DSN`
*   `ELASTICSEARCH_URLS`
*   `REDIS_ADDR`
*   `KAFKA_BROKERS`
*   `CDN_BASE_URL`
*   `S3_ENDPOINT`
*   `S3_ACCESS_KEY_ID`
*   `S3_SECRET_ACCESS_KEY`
*   `S3_BUCKET_MEDIA_RAW`
*   `S3_USE_SSL` (e.g., `true`)
*   `CDN_BASE_URL`
*   `LOG_LEVEL` (e.g., `info`, `debug`)
*   `AUTH_SERVICE_GRPC_ADDR` (e.g., `auth-service:9090`)
*   `PAYMENT_SERVICE_GRPC_ADDR` (e.g., `payment-service:9090`)
*   `DEVELOPER_SERVICE_GRPC_ADDR` (e.g., `developer-service:9090`)
*   `ADMIN_SERVICE_GRPC_ADDR` (e.g., `admin-service:9090`)
*   `DEFAULT_LANGUAGE` (e.g., `ru`)
*   `DEFAULT_REGION_CODE` (e.g., `RU`)
*   `DEFAULT_CURRENCY_CODE` (e.g., `RUB`)
*   `CACHE_PRODUCT_DETAILS_TTL_SECONDS` (e.g., `300`)
*   `CACHE_GENRE_LIST_TTL_SECONDS` (e.g., `3600`)
*   `CACHE_TAG_LIST_TTL_SECONDS` (e.g., `3600`)
*   `CACHE_CATEGORY_LIST_TTL_SECONDS` (e.g., `3600`)
*   `OTEL_EXPORTER_JAEGER_ENDPOINT` (e.g., `http://jaeger-collector:14268/api/traces`)

### 7.2. Файлы Конфигурации (если применимо)
*   `config.yaml` (пример структуры в Каталог Спец. v3.0, раздел 4.1.1 - Infrastructure/config). Файлы конфигурации используются для задания структуры, которая может быть переопределена переменными окружения.

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   REST API: Стандартные коды HTTP, JSON-формат ошибки.
*   gRPC API: Стандартные коды gRPC (`google.rpc.Status`).
*   Подробное логирование ошибок с `trace_id`.

### 8.2. Распространенные Коды Ошибок
*   `INVALID_ARGUMENT` / HTTP 400: Некорректные входные данные.
*   `NOT_FOUND` / HTTP 404: Ресурс не найден (продукт, жанр и т.д.).
*   `UNAUTHENTICATED` / HTTP 401: Ошибка аутентификации.
*   `PERMISSION_DENIED` / HTTP 403: Недостаточно прав.
*   `ALREADY_EXISTS` / HTTP 409: Попытка создания уже существующего ресурса (например, продукт с тем же developer_product_id).
*   `INTERNAL` / HTTP 500: Внутренняя ошибка сервера.
*   `UNAVAILABLE` / HTTP 503: Сервис временно недоступен или зависимость недоступна.

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   Внешние запросы: JWT через API Gateway.
*   Межсервисные gRPC: mTLS + API-ключи в метаданных.

### 9.2. Авторизация
*   RBAC. Проверка ролей (`X-User-Roles`) и владения ресурсом.
*   (Матрица доступа см. в Каталог Спец. v3.0, раздел 11.2.3).

### 9.3. Защита Данных
*   Шифрование дисков для БД, Elasticsearch, Redis. TLS для соединений.
*   Не хранение особо чувствительных данных (например, полных данных кредитных карт).

### 9.4. Управление Секретами
*   Kubernetes Secrets или HashiCorp Vault.
*   **Защита от уязвимостей:** Параметризованные запросы (SQL-инъекции), валидация и санитизация ввода (XSS, NoSQL-инъекции), ограничение скорости (DoS).
*   (Детали см. в Каталог Спец. v3.0, раздел 11).

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Многостадийный (golang:alpine -> alpine:latest). (см. Каталог Спец. v3.0, раздел 10.1.1).
*   **Kubernetes манифесты/Helm-чарты:** Deployment, Service, ConfigMap, Secret, HPA. (см. Каталог Спец. v3.0, раздел 10.2).

### 10.2. Зависимости при Развертывании
*   PostgreSQL, Elasticsearch, Redis, Kafka.
*   Auth Service, API Gateway.

### 10.3. CI/CD
*   Сборка, Unit/Integration тесты, статический анализ, сборка образа, развертывание (Staging, Production), E2E тесты.
*   (Детали см. в Каталог Спец. v3.0, раздел 10.3).
*   **Операционные процедуры:** Масштабирование (HPA), обновление/откат, миграции БД, переиндексация Elasticsearch, траблшутинг, бэкапы. (см. Каталог Спец. v3.0, раздел 10.4).

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат:** Структурированные JSON логи (Zap). Поля: `timestamp`, `level`, `service`, `version`, `trace_id`, `message`, `caller`, контекстные поля.
*   **Интеграция:** Сбор через Fluent Bit в ELK/Loki.
*   (Детали см. в Каталог Спец. v3.0, раздел 8.2).

### 11.2. Мониторинг
*   **Метрики (Prometheus):** Запросы (gRPC/HTTP), производительность Go, БД, кэш, Kafka, Elasticsearch, бизнес-метрики.
*   **Дашборды (Grafana):** Обзор сервиса, производительность, ресурсы, зависимости, бизнес-метрики.
*   **Алертинг (AlertManager):** Высокая задержка, ошибки, недоступность сервиса/зависимостей, низкий cache hit rate.
*   (Детали см. в Каталог Спец. v3.0, раздел 8.1).
*   **SLO/SLI:** Доступность 99.95%; Задержка чтения P95 < 150ms; Задержка записи P95 < 500ms. (см. Каталог Спец. v3.0, раздел 9).

### 11.3. Трассировка
*   **Инструментация:** OpenTelemetry.
*   **Создание спанов:** Входящие запросы, вызовы репозиториев, кэша, Kafka, Elasticsearch, исходящие вызовы.
*   **Контекст:** Пропагация W3C Trace Context.
*   **Экспорт:** Jaeger/Tempo.
*   (Детали см. в Каталог Спец. v3.0, раздел 8.3).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность (Latency):** P95 API чтения < 150ms; P95 API записи < 500ms.
*   **Производительность (Throughput):** >= 1000 RPS на чтение.
*   **Масштабируемость:** Горизонтальная.
*   **Доступность:** >= 99.95%.
*   **Консистентность:** Eventual consistency для реплик/кэшей, strong для основных записей.
*   (Полный список см. в Каталог Спец. v3.0, раздел 3.4).

## 13. Приложения (Appendices) (Опционально)
*   Полная спецификация OpenAPI.
*   Полные Protobuf определения.
*   Детальная схема DDL.
*   (Многие из этих деталей содержатся в Каталог Спец. v3.0).

---
*Этот шаблон является отправной точкой и может быть адаптирован под конкретные нужды проекта и сервиса.*

## 14. Связанные Рабочие Процессы (Related Workflows)
*   [Game Purchase and Library Update](../../../project_workflows/game_purchase_flow.md)
*   [Developer Submits a New Game for Moderation](../../../project_workflows/game_submission_flow.md)
