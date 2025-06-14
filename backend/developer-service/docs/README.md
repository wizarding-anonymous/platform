<!-- backend\developer-service\docs\README.md -->
# Спецификация Микросервиса: Developer Service (Сервис для Разработчиков)

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение:** Developer Service предоставляет интерфейс (Портал Разработчика) и API для разработчиков и издателей игр на платформе "Российский Аналог Steam".
*   **Роль в общей архитектуре платформы:** Developer Service является ключевым компонентом для управления жизненным циклом продуктов (игры, DLC, ПО) со стороны их создателей. Он обеспечивает инструменты для загрузки контента, управления метаданными, ценообразованием, доступом к аналитике продаж и использования продуктов, а также взаимодействия с финансовыми аспектами платформы (управление выплатами). Сервис тесно интегрирован с Auth Service, Catalog Service, Admin Service (для модерации), Payment Service и Analytics Service.
*   **Основные бизнес-задачи:**
    *   Привлечение и поддержка разработчиков и издателей.
    *   Обеспечение разработчиков инструментами для самостоятельной публикации и поддержки их продуктов.
    *   Управление процессом подачи, модерации и публикации продуктов.
    *   Предоставление разработчикам релевантной аналитики и финансовых отчетов.
    *   Управление командами разработчиков и их доступом к порталу.
*   Разработка сервиса должна вестись в соответствии с `../../../../CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Управление Аккаунтом Разработчика/Издателя:**
    *   Регистрация и верификация аккаунтов разработчиков/издателей.
    *   Управление профилем компании/разработчика (название, юридическая информация, контактные данные).
    *   Управление командой: добавление/удаление участников, назначение ролей и разрешений внутри команды (например, `owner`, `admin`, `technical_contact`, `marketing_contact`, `finance_contact`).
    *   Управление юридическими соглашениями и документами.
*   **Управление Продуктами (Игры, DLC, ПО):**
    *   Создание и управление карточками продуктов (игры, DLC, ПО, комплекты).
    *   Загрузка и управление билдами продуктов через интеграцию с S3-совместимым хранилищем. Поддержка различных платформ (Windows, Linux, macOS, Android, iOS - [Платформы: Windows, Linux, macOS, Android, iOS - подлежит уточнению в соответствии с `project_cross_platform_support.md`]).
    *   Версионирование билдов и управление их статусами (например, `alpha`, `beta`, `release_candidate`, `live`).
    *   Управление метаданными продуктов: локализованные названия, описания, системные требования, возрастные рейтинги, информация о разработчиках и издателях.
    *   Настройка страницы продукта в магазине: кастомизация описания, загрузка медиа-контента (скриншоты, трейлеры, арты), управление промо-материалами.
    *   Управление метаданными достижений для игр.
*   **Управление Ценообразованием и Публикацией:**
    *   Установка базовой цены продукта.
    *   Предложение и управление региональными ценами (в координации с Catalog Service).
    *   Создание и управление скидками, участие в промо-акциях платформы.
    *   Подача продукта/обновления на модерацию.
    *   Отслеживание статуса модерации (через интеграцию с Admin Service).
    *   Публикация одобренных продуктов и обновлений, управление их видимостью в каталоге.
*   **Аналитика и Отчетность:**
    *   Доступ к дашбордам с аналитикой по продажам, доходам, количеству установок, DAU/MAU и другим метрикам для опубликованных продуктов (данные предоставляются Analytics Service).
    *   Возможность генерации стандартных отчетов по продажам и активности пользователей.
*   **Финансовый Менеджмент:**
    *   Просмотр финансового баланса разработчика.
    *   История транзакций (продажи, возвраты, комиссии платформы).
    *   Управление банковскими реквизитами и методами для получения выплат.
    *   Формирование и отслеживание запросов на выплату средств (в координации с Payment Service).
*   **Управление SDK и API Ключами:**
    *   Доступ к SDK платформы и документации по интеграции.
    *   Создание и управление API ключами для автоматизации процессов CI/CD (например, загрузка билдов) и взаимодействия с API платформы от имени разработчика.
*   **Система Уведомлений:**
    *   Получение уведомлений о важных событиях: изменение статуса модерации продукта, необходимость обновления информации, финансовые операции, сообщения от службы поддержки платформы.

### 1.3. Основные Технологии
*   **Язык программирования:** Go (версия 1.21+, согласно `../../../../project_technology_stack.md`).
*   **Веб-фреймворк (REST API):** Echo (`github.com/labstack/echo/v4`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **gRPC:** `google.golang.org/grpc` для внутреннего взаимодействия с другими сервисами (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **База данных:** PostgreSQL (версия 15+) для хранения структурированных данных: аккаунты разработчиков, метаданные игр (черновики, специфичные для разработчика данные), информация о финансах и выплатах, API ключи. Драйвер: GORM (`gorm.io/gorm`) с `gorm.io/driver/postgres` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Кэширование:** Redis (версия 7.0+) для кэширования часто запрашиваемых данных, сессий Портала Разработчика, временных данных форм. Клиент: `go-redis/redis` (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Очереди сообщений/События:** Apache Kafka (клиент `github.com/confluentinc/confluent-kafka-go` или `github.com/segmentio/kafka-go`, согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Хранилище файлов (билды, медиа):** S3-совместимое объектное хранилище (например, MinIO, Yandex Object Storage) (согласно `../../../../project_technology_stack.md`).
*   **Управление конфигурацией:** Viper (`github.com/spf13/viper`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Логирование:** Zap (`go.uber.org/zap`) (согласно `../../../../PACKAGE_STANDARDIZATION.md`).
*   **Мониторинг/Трассировка:** OpenTelemetry SDK, Prometheus client (`github.com/prometheus/client_golang`). (согласно `../../../../project_observability_standards.md`).
*   **Инфраструктура:** Docker, Kubernetes.
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../PACKAGE_STANDARDIZATION.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   **Разработчик (Developer Account):** Учетная запись компании или индивидуального разработчика/издателя, зарегистрированного на платформе.
*   **Издатель (Publisher):** Компания или лицо, ответственное за публикацию и маркетинг продукта. Может быть тем же, что и разработчик, или отдельной сущностью.
*   **Продукт (Product/Game):** Игра, DLC, программное обеспечение или другой цифровой товар, управляемый разработчиком через Developer Service.
*   **Билд (Build):** Конкретная сборка исполняемых файлов и связанных ассетов продукта для определенной платформы.
*   **Версия (Version):** Публикуемая или тестовая версия продукта, связанная с конкретным билдом и набором метаданных.
*   **Портал Разработчика (Developer Portal):** Веб-интерфейс, предоставляемый Developer Service для разработчиков и издателей.
*   **Метаданные Продукта:** Вся информация о продукте, видимая пользователям в магазине и в их библиотеке (названия, описания, скриншоты, системные требования и т.д.).
*   **Выплата (Payout):** Перечисление заработанных разработчиком средств от продаж его продуктов на платформе.
*   Для других общих терминов см. `../../../../project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Developer Service будет реализован как модульный сервис, придерживаясь принципов Чистой Архитектуры (Clean Architecture) для разделения ответственностей и обеспечения тестируемости.
*   Сервис предоставляет REST API для Портала Разработчика (фронтенд) и для внешних систем разработчиков (например, CI/CD для загрузки билдов).
*   Ключевые модули включают: Управление Аккаунтами Разработчиков, Управление Продуктами, Управление Загрузками (билды, медиа), Модуль Аналитики (отображение данных от Analytics Service), Финансовый Модуль (выплаты), Управление API Ключами.

### 2.2. Диаграмма Архитектуры (Clean Architecture)
```mermaid
graph TD
    subgraph UserInteraction ["Портал Разработчика / CI/CD Системы"]
        DevPortal[Портал Разработчика (Веб-интерфейс)]
        DevAPIClient[Клиенты API Разработчика (CI/CD, утилиты)]
    end

    subgraph DeveloperService ["Developer Service (Чистая Архитектура)"]
        direction TB

        subgraph PresentationLayer [Presentation Layer (Адаптеры Транспорта)]
            REST_API[REST API (Echo/Gin)]
            GRPC_Internal_API[gRPC API (для внутренних вызовов от других сервисов, если потребуется)]
        end

        subgraph ApplicationLayer [Application Layer (Сценарии Использования)]
            DevAccountSvc[Управление Аккаунтом Разработчика]
            ProductManagementSvc[Управление Продуктами (Игры, DLC)]
            BuildUploadSvc[Управление Загрузками Билдов]
            AnalyticsAccessSvc[Доступ к Аналитике Продуктов]
            PayoutManagementSvc[Управление Выплатами]
            ApiKeyManagementSvc[Управление API Ключами]
            ModerationCoordinationSvc[Координация с Модерацией]
        end

        subgraph DomainLayer [Domain Layer (Бизнес-логика и Сущности)]
            Entities[Сущности (DeveloperAccount, DeveloperTeamMember, ProductSubmission, BuildArtifact, PayoutRequest, DevApiKey)]
            Aggregates[Агрегаты (DeveloperProfile, Product)]
            DomainEvents[Доменные События (ProductSubmittedForReviewEvent, PayoutRequestedEvent)]
            RepositoryIntf[Интерфейсы Репозиториев]
        end

        subgraph InfrastructureLayer [Infrastructure Layer (Внешние Зависимости и Реализации)]
            PostgresAdapter[Адаптер PostgreSQL (Реализация Репозиториев)]
            S3Adapter[Адаптер S3-хранилища (Загрузка/Скачивание Файлов)]
            RedisAdapter[Адаптер Redis (Кэш Сессий, Черновиков)]
            KafkaProducer[Продюсер Kafka (Публикация Событий)]
            KafkaConsumer[Консьюмер Kafka (Потребление Событий от Admin/Analytics)]
            AuthSvcClient[Клиент Auth Service (gRPC)]
            CatalogSvcClient[Клиент Catalog Service (gRPC/REST)]
            PaymentSvcClient[Клиент Payment Service (gRPC/REST)]
            AnalyticsSvcClient[Клиент Analytics Service (gRPC/REST)]
            NotificationSvcClient[Клиент Notification Service (gRPC/Kafka)]
            Config[Конфигурация (Viper)]
            Logging[Логирование (Zap)]
        end

        REST_API --> ApplicationLayer
        GRPC_Internal_API --> ApplicationLayer

        ApplicationLayer --> DomainLayer
        ApplicationLayer --> InfrastructureLayer
        DomainLayer ----> RepositoryIntf
        InfrastructureLayer -- Implements --> RepositoryIntf
    end

    DevPortal --> REST_API
    DevAPIClient --> REST_API

    PostgresAdapter --> DB[(PostgreSQL)]
    S3Adapter --> S3[(S3 Хранилище Билдов и Медиа)]
    RedisAdapter --> Cache[(Redis)]
    KafkaProducer --> KafkaBroker[Kafka Message Bus]
    KafkaConsumer --> KafkaBroker

    AuthSvcClient --> AuthService[Auth Service]
    CatalogSvcClient --> CatalogService[Catalog Service]
    PaymentSvcClient --> PaymentService[Payment Service]
    AnalyticsSvcClient --> AnalyticsService[Analytics Service]
    NotificationSvcClient --> NotificationService[Notification Service]


    classDef layer_boundary fill:#f9f9f9,stroke:#333,stroke-width:2px,color:#333
    classDef component_major fill:#e6f0ff,stroke:#007bff,color:#000
    classDef component_minor fill:#d4edda,stroke:#28a745,color:#000
    classDef datastore fill:#f8d7da,stroke:#dc3545,color:#000
    classDef external_service fill:#FEF9E7,stroke:#F1C40F,color:#000

    class PresentationLayer,ApplicationLayer,DomainLayer,InfrastructureLayer layer_boundary
    class REST_API,GRPC_Internal_API,DevAccountSvc,ProductManagementSvc,BuildUploadSvc,AnalyticsAccessSvc,PayoutManagementSvc,ApiKeyManagementSvc,ModerationCoordinationSvc,Entities,Aggregates,DomainEvents,RepositoryIntf component_major
    class PostgresAdapter,S3Adapter,RedisAdapter,KafkaProducer,KafkaConsumer,AuthSvcClient,CatalogSvcClient,PaymentSvcClient,AnalyticsSvcClient,NotificationSvcClient,Config,Logging component_minor
    class DB,S3,Cache,KafkaBroker datastore
    class AuthService,CatalogService,PaymentService,AnalyticsService,NotificationService external_service
```

### 2.3. Слои Сервиса

#### 2.3.1. Presentation Layer (Слой Представления)
*   **Ответственность:** Обработка входящих HTTP REST запросов от Портала Разработчика и публичного API для разработчиков. Валидация данных запроса (DTO), аутентификация (через Auth Service, передача JWT), авторизация (на основе `developer_id` и роли пользователя в команде разработчика), вызов соответствующей бизнес-логики в Application Layer.
*   **Ключевые компоненты/модули:** HTTP хендлеры (контроллеры) на базе Echo, DTO для запросов и ответов API, middleware для аутентификации и авторизации.

#### 2.3.2. Application Layer (Прикладной Слой)
*   **Ответственность:** Реализация сценариев использования (use cases), связанных с управлением аккаунтами разработчиков, их продуктами, финансами, загрузкой контента и т.д. Координирует взаимодействие между Domain Layer и Infrastructure Layer.
*   **Ключевые компоненты/модули:** Сервисы сценариев использования (например, `DeveloperAccountService`, `ProductSubmissionService`, `BuildManagementService`, `PayoutOrchestrationService`, `DeveloperAnalyticsService`), обработчики команд и запросов (если используется CQRS).

#### 2.3.3. Domain Layer (Доменный Слой)
*   **Ответственность:** Содержит бизнес-сущности, агрегаты, доменные события и бизнес-правила, специфичные для Developer Service. Например, правила валидации данных продукта, логика смены статусов продукта, правила формирования запросов на выплаты.
*   **Ключевые компоненты/модули:** Сущности (`DeveloperAccount`, `DeveloperTeamMember`, `ProductSubmission` (черновик продукта), `BuildArtifact`, `PayoutRequest`, `DeveloperAPIKey`), объекты-значения (например, `LocalizedText`, `FinancialDetails`), доменные сервисы, интерфейсы репозиториев.

#### 2.3.4. Infrastructure Layer (Инфраструктурный Слой)
*   **Ответственность:** Реализация интерфейсов репозиториев для работы с PostgreSQL и Redis. Взаимодействие с S3-совместимым хранилищем для загрузки и управления файлами. Отправка и получение событий через Kafka. Взаимодействие с другими микросервисами (Auth, Catalog, Payment, Admin, Analytics, Notification) через их gRPC/REST API.
*   **Ключевые компоненты/модули:** Реализации репозиториев для PostgreSQL, S3 клиент, Kafka продюсер и консьюмер, Redis клиент, gRPC/HTTP клиенты для других сервисов.

## 3. API Endpoints

### 3.1. REST API
*   **Префикс:** `/api/v1/developer` (маршрутизируется через API Gateway).
*   **Аутентификация:** JWT Bearer Token, полученный от Auth Service для пользователя-разработчика. Проверяется на API Gateway или в middleware Developer Service. `developer_id` и `user_id` (члена команды) извлекаются из токена или контекста безопасности.
*   **Авторизация:** На основе `developer_id` и роли пользователя в команде разработчика (например, `owner`, `admin`, `editor`, `viewer`, `finance_manager`, `build_manager`). Разрешения проверяются для каждой операции.
*   **Формат ответа об ошибке:** Согласно `../../../../project_api_standards.md`.

#### 3.1.1. Аккаунты разработчиков (Developer Accounts)
*   **`POST /accounts`**
    *   Описание: Регистрация нового аккаунта разработчика/издателя. Связывает текущего аутентифицированного пользователя платформы как владельца.
    *   Тело запроса: (Как в существующем документе)
    *   Ответ: (Как в существующем документе)
    *   Требуемые права доступа: Аутентифицированный пользователь платформы.
*   **`GET /accounts/me`**
    *   Описание: Получение информации о текущем аккаунте разработчика, к которому привязан пользователь.
    *   Требуемые права доступа: Участник команды разработчика.
*   **`PUT /accounts/me`**
    *   Описание: Обновление профиля и юридической информации аккаунта разработчика.
    *   Требуемые права доступа: Роль `owner` или `admin` в команде разработчика.
*   **`GET /accounts/me/team`**
    *   Описание: Получение списка членов команды разработчика.
    *   Требуемые права доступа: Участник команды разработчика.
*   **`POST /accounts/me/team/members`**
    *   Описание: Приглашение нового участника в команду разработчика (по email).
    *   Требуемые права доступа: Роль `owner` или `admin` в команде разработчика.
*   **`PUT /accounts/me/team/members/{member_user_id}`**
    *   Описание: Изменение роли участника команды.
    *   Требуемые права доступа: Роль `owner` или `admin` в команде разработчика.
*   **`DELETE /accounts/me/team/members/{member_user_id}`**
    *   Описание: Удаление участника из команды.
    *   Требуемые права доступа: Роль `owner` или `admin` в команде разработчика.

#### 3.1.2. Управление Продуктами (Игры/DLC/ПО)
*   **`POST /products`**
    *   Описание: Создание нового продукта (игра, DLC, ПО) в Developer Service (создание черновика).
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "productDraftCreation",
            "attributes": {
              "title": {"ru-RU": "Моя Новая Игра", "en-US": "My New Game"},
              "product_type": "game" // game, dlc, software
            }
          }
        }
        ```
    *   Ответ: (Возвращает созданный черновик продукта с его ID)
    *   Требуемые права доступа: Роль `owner`, `admin`, `editor` в команде.
*   **`GET /products`**
    *   Описание: Получение списка продуктов, управляемых разработчиком.
    *   Query параметры: `status`, `product_type`, `page`, `limit`.
    *   Требуемые права доступа: Участник команды разработчика.
*   **`GET /products/{product_id}`**
    *   Описание: Получение детальной информации о продукте разработчика (черновик или опубликованная версия).
    *   Требуемые права доступа: Участник команды разработчика.
*   **`PUT /products/{product_id}`**
    *   Описание: Обновление метаданных черновика продукта.
    *   Тело запроса: (Полная структура метаданных продукта, включая локализованные поля, системные требования, медиа-ссылки, предлагаемые цены, теги, жанры и т.д.).
    *   Требуемые права доступа: Роль `owner`, `admin`, `editor` в команде.

#### 3.1.3. Версии Продуктов и Загрузка Билдов
*   **`POST /products/{product_id}/versions`**
    *   Описание: Создание новой версии для продукта.
    *   Тело запроса: (Как в существующем документе)
    *   Требуемые права доступа: Роль `owner`, `admin`, `build_manager`.
*   **`POST /products/{product_id}/versions/{version_id}/builds/upload-url`**
    *   Описание: Получение pre-signed URL для загрузки файла билда в S3.
    *   Тело запроса: (Как в существующем документе)
    *   Требуемые права доступа: Роль `owner`, `admin`, `build_manager`.
*   **`POST /products/{product_id}/versions/{version_id}/builds/upload-complete`**
    *   Описание: Уведомление сервиса об успешной загрузке билда в S3.
    *   Тело запроса: (Как в существующем документе)
    *   Требуемые права доступа: Роль `owner`, `admin`, `build_manager`.

#### 3.1.4. Публикация и Модерация
*   **`POST /products/{product_id}/submit-for-review`**
    *   Описание: Отправка продукта (и его текущей черновой версии метаданных/билдов) на модерацию.
    *   Ответ: `{ "data": { "status": "pending_moderation" } }`
    *   Требуемые права доступа: Роль `owner`, `admin`, `release_manager`.
*   **`GET /products/{product_id}/moderation-status`**
    *   Описание: Получение текущего статуса модерации для продукта.
    *   Требуемые права доступа: Участник команды разработчика.

#### 3.1.5. Финансы и Выплаты
*   **`GET /finance/balance`**
    *   Описание: Получение текущего финансового баланса разработчика.
    *   Требуемые права доступа: Роль `owner`, `admin`, `finance_manager`.
*   **`GET /finance/transactions`**
    *   Описание: Получение истории финансовых транзакций (продажи, возвраты, комиссии).
    *   Query параметры: `start_date`, `end_date`, `type`, `page`, `limit`.
    *   Требуемые права доступа: Роль `owner`, `admin`, `finance_manager`.
*   **`POST /finance/payouts/requests`**
    *   Описание: Создание запроса на выплату средств.
    *   Тело запроса: (Как в существующем документе)
    *   Требуемые права доступа: Роль `owner`, `admin`, `finance_manager`.
*   **`GET /finance/payouts/requests`**
    *   Описание: Получение истории запросов на выплаты.
    *   Требуемые права доступа: Роль `owner`, `admin`, `finance_manager`.

#### 3.1.6. Аналитика
*   **`GET /analytics/products/{product_id}/summary`**
    *   Описание: Получение сводной аналитики по продукту (продажи, DAU, MAU). Проксирует запрос к Analytics Service.
    *   Query параметры: `period` (например, `last_7_days`, `last_30_days`, `custom_range`).
    *   Требуемые права доступа: Участник команды разработчика.
*   **`GET /analytics/reports`**
    *   Описание: Запрос на генерацию или получение списка доступных отчетов от Analytics Service.
    *   Требуемые права доступа: Участник команды разработчика.

#### 3.1.7. API Ключи Разработчика
*   **`POST /api-keys`**
    *   Описание: Создание нового API ключа для разработчика.
    *   Тело запроса: `{"data": {"type": "apiKeyCreation", "attributes": {"name": "CI/CD Key", "permissions": ["upload_build:game_id_123"], "expires_at": "YYYY-MM-DDTHH:mm:ssZ"}}}`
    *   Ответ: (Возвращает созданный API ключ **один раз**)
    *   Требуемые права доступа: Роль `owner`, `admin`.
*   **`GET /api-keys`**
    *   Описание: Получение списка API ключей (без самих значений ключей, только метаданные).
    *   Требуемые права доступа: Роль `owner`, `admin`.
*   **`DELETE /api-keys/{api_key_id}`**
    *   Описание: Отзыв (удаление) API ключа.
    *   Требуемые права доступа: Роль `owner`, `admin`.

### 3.2. gRPC API
*   На данный момент Developer Service в основном **потребляет** gRPC API других сервисов.
*   Если в будущем возникнет необходимость в предоставлении собственных gRPC эндпоинтов (например, для CLI-утилиты разработчика или специфичных внутренних интеграций), они будут определены и задокументированы здесь в соответствии со стандартами `project_api_standards.md`.

### 3.3. WebSocket API
*   Не планируется для Developer Service на данном этапе. Может быть рассмотрено в будущем для real-time уведомлений на Портале Разработчика.

## 4. Модели Данных (Data Models)
См. также `../../../../project_database_structure.md`.

### 4.1. Основные Сущности
*   **`DeveloperAccount` (Аккаунт Разработчика/Издателя)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `owner_user_id` (UUID, FK на User в Auth Service): Пользователь-владелец аккаунта. **Обязательность: Да.**
    *   `company_name` (VARCHAR). **Обязательность: Да.**
    *   `legal_entity_type` (VARCHAR, например, "ООО", "ИП", "Самозанятый"). **Обязательность: Да (при регистрации, может быть уточнено при верификации).**
    *   `tax_id` (VARCHAR, например, ИНН для РФ). **Обязательность: Да (для юр.лиц и ИП РФ, может быть уточнено при верификации).**
    *   `country_code` (CHAR(2), ISO 3166-1 alpha-2). **Обязательность: Да.**
    *   `address_legal` (JSONB): Юридический адрес. Структура зависит от страны. **Обязательность: Да (для юр.лиц).**
    *   `address_postal` (JSONB): Почтовый адрес, если отличается от юридического. **Обязательность: Нет.**
    *   `contact_email` (VARCHAR, уникальный). **Обязательность: Да.**
    *   `contact_phone` (VARCHAR, nullable). **Обязательность: Нет.**
    *   `website_url` (VARCHAR, nullable). **Обязательность: Нет.**
    *   `status` (VARCHAR: `pending_verification`, `active`, `limited_access`, `suspended`, `rejected`, `closed`). **Обязательность: Да (DEFAULT 'pending_verification').**
    *   `verification_documents_s3_links` (JSONB, массив ссылок на файлы в S3). **Обязательность: Нет.**
    *   `default_payout_method_id` (UUID, FK на `DeveloperPaymentMethod`, nullable). **Обязательность: Нет.**
    *   `developer_agreement_version_accepted` (VARCHAR, nullable). **Обязательность: Нет.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `verified_at` (TIMESTAMPTZ, nullable). **Обязательность: Нет.**
    *   `closed_at` (TIMESTAMPTZ, nullable). **Обязательность: Нет.**
*   **`DeveloperTeamMember` (Член Команды Разработчика)**
    *   `developer_account_id` (UUID, PK, FK на DeveloperAccount). **Обязательность: Да.**
    *   `user_id` (UUID, PK, FK на User в Auth Service). **Обязательность: Да.**
    *   `role_in_team` (VARCHAR: `owner`, `admin`, `editor`, `viewer`, `finance_manager`, `build_manager`, `marketing_manager`, `support_agent`). **Обязательность: Да.**
    *   `permissions_override` (JSONB, nullable): Индивидуальные переопределения прав (например, доступ к конкретным продуктам). **Обязательность: Нет.**
    *   `invited_by_user_id` (UUID, FK на User в Auth Service, nullable). **Обязательность: Нет (если первый владелец).**
    *   `invitation_status` (VARCHAR: `pending`, `accepted`, `declined`, nullable). **Обязательность: Нет.**
    *   `joined_at` (TIMESTAMPTZ, nullable). **Обязательность: Нет (устанавливается после принятия приглашения).**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
*   **`ProductSubmission` (Представление Продукта/Черновик)**: Используется для хранения данных о продукте (игра, DLC, ПО), пока он находится в разработке или на модерации в Developer Service, перед его передачей/синхронизацией с Catalog Service.
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `developer_account_id` (UUID, FK на DeveloperAccount). **Обязательность: Да.**
    *   `catalog_product_id` (UUID, nullable, UK): ID продукта в Catalog Service. Заполняется после успешной первой синхронизации. **Обязательность: Нет.**
    *   `product_type` (VARCHAR: `game`, `dlc`, `software`, `bundle`, `soundtrack`, `artbook`, `demo`). **Обязательность: Да.**
    *   `title_main_language_code` (CHAR(5), например "ru-RU", "en-US"): Основной язык, на котором представлено название и описание. **Обязательность: Да.**
    *   `titles` (JSONB): Локализованные названия продукта. `{"ru-RU": "Моя Игра", "en-US": "My Game"}`. **Обязательность: Да (хотя бы для основного языка).**
    *   `status_developer` (VARCHAR: `draft`, `ready_for_review`, `in_review_internal`, `changes_requested_internal`, `approved_internal`). Статус со стороны разработчика и внутренней проверки Developer Service. **Обязательность: Да (DEFAULT 'draft').**
    *   `status_moderation_platform` (VARCHAR: `not_submitted`, `pending_moderation`, `in_moderation`, `changes_requested_by_moderator`, `approved_by_moderator`, `rejected_by_moderator`). Статус модерации со стороны платформы (Admin Service). **Обязательность: Да (DEFAULT 'not_submitted').**
    *   `current_draft_version_id` (UUID, FK на `ProductVersionDraft`, nullable): Ссылка на текущий активный черновик версии продукта. **Обязательность: Нет.**
    *   `live_version_id_in_catalog` (UUID, nullable): ID версии, которая сейчас "live" в Catalog Service (может быть не из этого сервиса, а из ProductVersionLive). **Обязательность: Нет.**
    *   `moderation_request_id_admin_service` (UUID, nullable): ID тикета/запроса на модерацию в Admin Service. **Обязательность: Нет.**
    *   `developer_notes_for_moderator` (TEXT, nullable). **Обязательность: Нет.**
    *   `moderator_feedback_to_developer` (TEXT, nullable). **Обязательность: Нет.**
    *   `tags_internal` (JSONB, массив строк): Внутренние теги для Developer Service. **Обязательность: Нет.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `submitted_for_moderation_at` (TIMESTAMPTZ, nullable). **Обязательность: Нет.**
    *   `last_moderation_decision_at` (TIMESTAMPTZ, nullable). **Обязательность: Нет.**
*   **`ProductVersionDraft` (Черновик Версии Продукта)**: Представляет собой конкретную версию продукта, над которой работает разработчик.
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `product_submission_id` (UUID, FK на ProductSubmission). **Обязательность: Да.**
    *   `version_name` (VARCHAR, например, "1.0.0", "Beta 2.1", "New Year Update"). **Обязательность: Да.**
    *   `version_number_internal` (INTEGER, автоинкремент в рамках product_submission_id, опционально). **Обязательность: Нет.**
    *   `status` (VARCHAR: `draft`, `builds_processing`, `ready_for_submission`, `submitted`, `approved`, `rejected`). **Обязательность: Да (DEFAULT 'draft').**
    *   `metadata_draft` (JSONB): Полная структура метаданных для этой версии (описания, системные требования, возрастные рейтинги, ссылки на медиа в S3 и т.д.). **Обязательность: Да (может быть пустым JSON '{}').**
    *   `pricing_draft` (JSONB): Предлагаемые цены для этой версии (базовая цена, региональные оверрайды, информация о скидках). **Обязательность: Да (может быть пустым JSON '{}').**
    *   `changelog` (JSONB, локализованный): Список изменений для этой версии. **Обязательность: Нет.**
    *   `release_date_planned` (TIMESTAMPTZ, nullable): Планируемая дата релиза этой версии. **Обязательность: Нет.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
*   **`BuildArtifact` (Артефакт Билда)**: Представляет собой конкретный файл билда для определенной платформы и версии продукта.
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `product_version_draft_id` (UUID, FK на ProductVersionDraft). **Обязательность: Да.**
    *   `platform_code` (VARCHAR: `windows_x64`, `linux_x86_64`, `macos_arm64`, `android_arm64v8a`, `ios_arm64`). См. `project_cross_platform_support.md`. **Обязательность: Да.**
    *   `s3_bucket_name` (VARCHAR). **Обязательность: Да.**
    *   `s3_object_key` (VARCHAR, UK в рамках product_version_draft_id + platform_code): Путь к файлу билда в S3. **Обязательность: Да.**
    *   `original_file_name` (VARCHAR). **Обязательность: Да.**
    *   `file_size_bytes` (BIGINT). **Обязательность: Да.**
    *   `file_hash_sha256` (VARCHAR(64), nullable). **Обязательность: Нет (вычисляется после загрузки).**
    *   `status` (VARCHAR: `pending_upload`, `uploading`, `uploaded`, `processing_requested`, `processing`, `ready`, `failed_processing`, `deprecated`). **Обязательность: Да (DEFAULT 'pending_upload').**
    *   `upload_session_id` (VARCHAR, nullable): Идентификатор сессии загрузки (для multipart uploads). **Обязательность: Нет.**
    *   `upload_expires_at` (TIMESTAMPTZ, для pre-signed URL). **Обязательность: Нет.**
    *   `processing_details` (JSONB, nullable): Информация о процессе обработки билда (например, результаты сканирования). **Обязательность: Нет.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
*   **`DeveloperPayoutRequest` (Запрос на Выплату)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `developer_account_id` (UUID, FK на DeveloperAccount). **Обязательность: Да.**
    *   `amount_requested_minor_units` (BIGINT, >0). **Обязательность: Да.**
    *   `currency_code` (VARCHAR(3), ISO 4217). **Обязательность: Да.**
    *   `status` (VARCHAR: `pending_review`, `approved_by_developer_service`, `rejected_by_developer_service`, `pending_processing_payment_service`, `processing_by_payment_service`, `completed`, `failed_payment_service`, `cancelled_by_developer`). **Обязательность: Да (DEFAULT 'pending_review').**
    *   `developer_payment_method_id` (UUID, FK на `DeveloperPaymentMethod`). **Обязательность: Да.**
    *   `payment_method_details_snapshot` (JSONB): Снимок реквизитов на момент запроса для аудита. **Обязательность: Да.**
    *   `requested_at` (TIMESTAMPTZ). **Обязательность: Да (DEFAULT now()).**
    *   `decision_by_developer_service_at` (TIMESTAMPTZ, nullable). **Обязательность: Нет.**
    *   `decision_notes_developer_service` (TEXT, nullable): Примечания от Developer Service (например, причина отклонения). **Обязательность: Нет.**
    *   `payment_service_transaction_id` (UUID, nullable): ID транзакции, присвоенный Payment Service. **Обязательность: Нет.**
    *   `payment_service_processing_started_at` (TIMESTAMPTZ, nullable). **Обязательность: Нет.**
    *   `payment_service_last_update_at` (TIMESTAMPTZ, nullable). **Обязательность: Нет.**
    *   `completed_at` (TIMESTAMPTZ, nullable). **Обязательность: Нет.**
    *   `developer_comment` (TEXT, nullable). **Обязательность: Нет.**
*   **`DeveloperPaymentMethod` (Платежный Метод Разработчика)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `developer_account_id` (UUID, FK на DeveloperAccount). **Обязательность: Да.**
    *   `method_type` (VARCHAR: `bank_transfer_ru_individual`, `bank_transfer_ru_entity`, `swift_transfer_usd`, `swift_transfer_eur`, `crypto_wallet_usdt_trc20_placeholder`). **Обязательность: Да.**
    *   `details_encrypted` (TEXT): Зашифрованные банковские реквизиты или адрес кошелька. Шифруется на уровне приложения перед сохранением в БД. **Обязательность: Да.**
    *   `is_default` (BOOLEAN). **Обязательность: Да (DEFAULT FALSE).** Должен быть только один метод по умолчанию для аккаунта.
    *   `is_verified_by_platform` (BOOLEAN). **Обязательность: Да (DEFAULT FALSE).**
    *   `verification_notes` (TEXT, nullable). **Обязательность: Нет.**
    *   `display_name` (VARCHAR, nullable, например "Мой счет в Альфа-Банке"). **Обязательность: Нет.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `verified_at` (TIMESTAMPTZ, nullable). **Обязательность: Нет.**
*   **`DeveloperAPIKey` (API Ключ Разработчика)**
    *   `id` (UUID, PK). **Обязательность: Да (генерируется БД).**
    *   `developer_account_id` (UUID, FK на DeveloperAccount). **Обызательность: Да.**
    *   `name` (VARCHAR). **Обязательность: Да.**
    *   `description` (TEXT, nullable). **Обязательность: Нет.**
    *   `prefix` (VARCHAR(8), UK): Короткий префикс ключа, видимый пользователю. **Обязательность: Да.**
    *   `key_hash` (VARCHAR(255), UK): Хеш самого API ключа (ключ целиком не хранится). **Обязательность: Да.**
    *   `permissions` (JSONB, массив строк, описывающих разрешения, например `["product:read:*", "product:write:prod_xyz", "build:upload:prod_xyz"]`). **Обязательность: Да (DEFAULT '[]').**
    *   `is_active` (BOOLEAN). **Обязательность: Да (DEFAULT TRUE).**
    *   `expires_at` (TIMESTAMPTZ, nullable): Дата и время истечения срока действия ключа. **Обязательность: Нет.**
    *   `last_used_at` (TIMESTAMPTZ, nullable). **Обязательность: Нет.**
    *   `last_used_from_ip` (VARCHAR, nullable). **Обязательность: Нет.**
    *   `created_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**
    *   `updated_at` (TIMESTAMPTZ). **Обязательность: Да (генерируется БД).**

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL
**ERD Диаграмма:**
```mermaid
erDiagram
    DEVELOPER_ACCOUNTS {
        UUID id PK
        UUID owner_user_id FK "USERS(id) from AuthSvc"
        VARCHAR company_name
        VARCHAR legal_entity_type
        VARCHAR tax_id
        VARCHAR contact_email UK
        VARCHAR status
        JSONB verification_documents
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    DEVELOPER_TEAM_MEMBERS {
        UUID developer_account_id PK FK
        UUID user_id PK FK "USERS(id) from AuthSvc"
        VARCHAR role_in_team
        JSONB permissions_override
        UUID invited_by_user_id FK "USERS(id) from AuthSvc"
        TIMESTAMPTZ joined_at
    }
    PRODUCT_SUBMISSIONS {
        UUID id PK
        UUID developer_account_id FK
        UUID catalog_product_id "nullable, FK to CatalogSvc.PRODUCTS(id)"
        VARCHAR product_type
        VARCHAR status
        UUID current_version_id "nullable, FK to PRODUCT_VERSIONS(id)"
        JSONB draft_metadata
        JSONB draft_pricing
        TEXT moderation_notes_to_admin
        TEXT moderation_feedback_from_admin
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
        TIMESTAMPTZ submitted_at
    }
    PRODUCT_VERSIONS {
        UUID id PK
        UUID product_submission_id FK
        VARCHAR version_name
        VARCHAR status
        JSONB changelog
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    BUILD_ARTIFACTS {
        UUID id PK
        UUID product_version_id FK
        VARCHAR platform
        VARCHAR s3_path
        VARCHAR file_name
        BIGINT file_size_bytes
        VARCHAR file_hash_sha256
        VARCHAR status
        TIMESTAMPTZ upload_expires_at
        TIMESTAMPTZ created_at
    }
    DEVELOPER_PAYMENT_METHODS {
        UUID id PK
        UUID developer_account_id FK
        VARCHAR method_type
        TEXT details_encrypted
        BOOLEAN is_default
        BOOLEAN is_verified
        TIMESTAMPTZ created_at
    }
    DEVELOPER_PAYOUT_REQUESTS {
        UUID id PK
        UUID developer_account_id FK
        BIGINT amount_requested_minor_units
        VARCHAR currency_code
        VARCHAR status
        UUID payment_method_id FK
        JSONB payment_method_details_snapshot
        TEXT comment_developer
        TEXT comment_admin
        UUID transaction_id_payment_service "nullable"
        TIMESTAMPTZ requested_at
        TIMESTAMPTZ processed_at
    }
    DEVELOPER_API_KEYS {
        UUID id PK
        UUID developer_account_id FK
        VARCHAR name
        VARCHAR prefix UK
        VARCHAR key_hash UK
        JSONB permissions
        BOOLEAN is_active
        TIMESTAMPTZ expires_at "nullable"
        TIMESTAMPTZ last_used_at "nullable"
        TIMESTAMPTZ created_at
    }

    DEVELOPER_ACCOUNTS ||--o{ DEVELOPER_TEAM_MEMBERS : "has team"
    DEVELOPER_ACCOUNTS ||--o{ PRODUCT_SUBMISSIONS : "submits"
    DEVELOPER_ACCOUNTS ||--o{ DEVELOPER_PAYMENT_METHODS : "has payment methods"
    DEVELOPER_ACCOUNTS ||--o{ DEVELOPER_PAYOUT_REQUESTS : "requests payouts"
    DEVELOPER_ACCOUNTS ||--o{ DEVELOPER_API_KEYS : "owns API keys"
    PRODUCT_SUBMISSIONS ||--o{ PRODUCT_VERSIONS : "has versions"
    PRODUCT_VERSIONS ||--o{ BUILD_ARTIFACTS : "has builds"
    DEVELOPER_PAYMENT_METHODS ||--o{ DEVELOPER_PAYOUT_REQUESTS : "used for"
```

**DDL (PostgreSQL - примеры для новых/ключевых таблиц):**
```sql
-- Аккаунты Разработчиков
CREATE TABLE developer_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_user_id UUID NOT NULL, -- Внешний ключ на Users из Auth Service
    company_name VARCHAR(255) NOT NULL,
    legal_entity_type VARCHAR(100),
    tax_id VARCHAR(50),
    country_code CHAR(2) NOT NULL,
    contact_email VARCHAR(255) NOT NULL UNIQUE,
    contact_phone VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'pending_verification', -- pending_verification, active, suspended, rejected
    verification_documents JSONB, -- Ссылки на документы в S3
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON COLUMN developer_accounts.owner_user_id IS 'ID пользователя-владельца из Auth Service';

-- Члены Команды Разработчика
CREATE TABLE developer_team_members (
    developer_account_id UUID NOT NULL REFERENCES developer_accounts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL, -- Внешний ключ на Users из Auth Service
    role_in_team VARCHAR(50) NOT NULL, -- owner, admin, editor, viewer, finance_manager, build_manager
    permissions_override JSONB,
    invited_by_user_id UUID, -- Внешний ключ на Users из Auth Service
    joined_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (developer_account_id, user_id)
);
COMMENT ON COLUMN developer_team_members.user_id IS 'ID пользователя из Auth Service';

-- Представления Продуктов (Черновики)
CREATE TABLE product_submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    developer_account_id UUID NOT NULL REFERENCES developer_accounts(id) ON DELETE CASCADE,
    catalog_product_id UUID, -- ID продукта в Catalog Service после первой синхронизации
    product_type VARCHAR(50) NOT NULL DEFAULT 'game', -- game, dlc, software, bundle
    status VARCHAR(50) NOT NULL DEFAULT 'draft', -- draft, in_review, changes_requested, approved_by_developer_service
    current_version_id UUID, -- FK на product_versions(id) - устанавливается позже
    draft_metadata JSONB,
    draft_pricing JSONB,
    moderation_notes_to_admin TEXT,
    moderation_feedback_from_admin TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    submitted_at TIMESTAMPTZ
);
CREATE INDEX idx_product_submissions_developer_id ON product_submissions(developer_account_id);
CREATE INDEX idx_product_submissions_status ON product_submissions(status);

-- Версии Продуктов
CREATE TABLE product_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_submission_id UUID NOT NULL REFERENCES product_submissions(id) ON DELETE CASCADE,
    version_name VARCHAR(100) NOT NULL, -- e.g., "1.0.0", "Beta 2.1"
    status VARCHAR(50) NOT NULL DEFAULT 'draft', -- draft, uploading_builds, ready_for_submission, submitted_for_review, live, deprecated
    changelog JSONB, -- Локализованный список изменений
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
ALTER TABLE product_submissions ADD CONSTRAINT fk_current_version FOREIGN KEY (current_version_id) REFERENCES product_versions(id) ON DELETE SET NULL;


-- Артефакты Билдов
CREATE TABLE build_artifacts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_version_id UUID NOT NULL REFERENCES product_versions(id) ON DELETE CASCADE,
    platform VARCHAR(100) NOT NULL, -- windows_x64, linux_x86_64, macos_arm64, etc.
    s3_path VARCHAR(1024) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    file_hash_sha256 VARCHAR(64),
    status VARCHAR(50) NOT NULL DEFAULT 'uploading', -- uploading, uploaded, processing, ready, failed_processing
    upload_expires_at TIMESTAMPTZ, -- Для pre-signed URL
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_build_artifacts_version_platform ON build_artifacts(product_version_id, platform);

-- Платежные Методы Разработчиков
CREATE TABLE developer_payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    developer_account_id UUID NOT NULL REFERENCES developer_accounts(id) ON DELETE CASCADE,
    method_type VARCHAR(50) NOT NULL, -- bank_transfer_ru, swift_transfer, etc.
    details_encrypted TEXT NOT NULL, -- Зашифрованные реквизиты
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Запросы на Выплаты
CREATE TABLE developer_payout_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    developer_account_id UUID NOT NULL REFERENCES developer_accounts(id) ON DELETE CASCADE,
    amount_requested_minor_units BIGINT NOT NULL CHECK (amount_requested_minor_units > 0),
    currency_code VARCHAR(3) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_request', -- pending_request, pending_approval, processing, completed, failed, cancelled
    payment_method_id UUID NOT NULL REFERENCES developer_payment_methods(id),
    payment_method_details_snapshot JSONB NOT NULL, -- Снимок реквизитов на момент запроса
    comment_developer TEXT,
    comment_admin TEXT,
    transaction_id_payment_service UUID, -- ID транзакции в Payment Service
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    approved_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ
);
CREATE INDEX idx_developer_payout_requests_developer_id ON developer_payout_requests(developer_account_id);
CREATE INDEX idx_developer_payout_requests_status ON developer_payout_requests(status);

-- API Ключи Разработчиков
CREATE TABLE developer_api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    developer_account_id UUID NOT NULL REFERENCES developer_accounts(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    prefix VARCHAR(8) NOT NULL UNIQUE,
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    permissions JSONB NOT NULL DEFAULT '[]'::jsonb, -- Список разрешений ключа
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    expires_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_developer_api_keys_developer_id ON developer_api_keys(developer_account_id);
```

#### 4.2.2. Redis
*   **Сессии Портала Разработчика:** `session:developer:<session_id>` - хранение данных сессии пользователя Портала Разработчика. TTL устанавливается.
*   **Черновики Форм (опционально):** `form_draft:developer:<user_id>:product_submission:<product_id>` - для сохранения промежуточного состояния больших форм (например, при создании/редактировании продукта). TTL короткий.
*   **Счетчики Rate Limiting:** Для API эндпоинтов Developer Service (если не управляется централизованно API Gateway).

#### 4.2.3. S3-совместимое хранилище
*   **Билды игр:** `s3://<bucket-name-game-builds>/<developer_account_id>/<product_submission_id>/<product_version_id>/<platform>/<filename.zip>`
*   **Медиа-файлы продуктов:** `s3://<bucket-name-product-media>/<developer_account_id>/<product_submission_id>/media/<media_type>/<timestamp_filename.ext>` (например, `screenshots`, `trailers`, `cover_art`)
*   **Документы для верификации разработчика:** `s3://<bucket-name-developer-documents>/<developer_account_id>/verification/<document_type>/<filename.pdf>` (с строгими правами доступа, возможно шифрование на стороне сервера S3).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Формат событий:** CloudEvents v1.0 JSON (согласно `../../../../project_api_standards.md`).
*   **Основной топик:** `com.platform.developer.events.v1`.

*   **`com.platform.developer.account.created.v1`**
    *   Описание: Аккаунт разработчика/издателя был успешно создан и ожидает верификации.
    *   `data` Payload:
        ```json
        {
          "developerAccountId": "dev-uuid-abc",
          "ownerUserId": "user-uuid-owner", // ID пользователя из Auth Service
          "companyName": "Моя Игровая Студия",
          "contactEmail": "dev@example.com",
          "status": "pending_verification",
          "creationTimestamp": "2024-07-11T10:00:00Z"
        }
        ```
    *   Потребители: Admin Service (для начала процесса верификации), Notification Service.
*   **`com.platform.developer.product.submitted.v1`**
    *   Описание: Разработчик отправил продукт (игру, DLC и т.д.) или его новую версию на модерацию.
    *   `data` Payload:
        ```json
        {
          "developerAccountId": "dev-uuid-abc",
          "productSubmissionId": "prodsub-uuid-xyz", // ID черновика/представления в Developer Service
          "catalogProductId": null, // Будет заполнен после создания в Catalog Service
          "productType": "game",
          "versionName": "1.0.0", // Если это подача конкретной версии
          "titles": {"ru-RU": "Моя Игра", "en-US": "My Game"}, // Ключевые метаданные
          "submissionTimestamp": "2024-07-11T14:30:00Z",
          "submittedByUserId": "user-uuid-dev-member"
        }
        ```
    *   Потребители: Admin Service (для начала процесса модерации), Catalog Service (для создания/обновления записи о продукте).
*   **`com.platform.developer.product.published.v1`**
     *   Описание: Разработчик опубликовал одобренный продукт/версию (или он был опубликован автоматически после одобрения).
    *   `data` Payload:
        ```json
        {
          "developerAccountId": "dev-uuid-abc",
          "productSubmissionId": "prodsub-uuid-xyz",
          "catalogProductId": "cat-prod-uuid-123",
          "versionName": "1.0.0",
          "publishedAt": "2024-07-11T16:00:00Z",
          "publishedByUserId": "user-uuid-dev-member"
        }
        ```
    *   Потребители: Catalog Service (для обновления статуса), Download Service (для подготовки билдов к скачиванию), Notification Service, Analytics Service.
*   **`com.platform.developer.build.uploaded.v1`**
    *   Описание: Разработчик успешно загрузил новый билд для версии продукта.
    *   `data` Payload:
        ```json
        {
            "developerAccountId": "dev-uuid-abc",
            "productSubmissionId": "prodsub-uuid-xyz",
            "productVersionId": "ver-uuid-456",
            "buildArtifactId": "build-uuid-789",
            "platform": "windows_x64",
            "fileName": "mygame_v1.0.0.zip",
            "s3Path": "path/to/build.zip",
            "uploadTimestamp": "2024-07-11T15:00:00Z"
        }
        ```
    *   Потребители: Download Service (для обработки и подготовки билда), Admin Service (информация для модерации).
*   **`com.platform.developer.payout.requested.v1`**
    *   Описание: Разработчик запросил выплату средств.
    *   `data` Payload:
        ```json
        {
          "payoutRequestId": "payout-uuid-789",
          "developerAccountId": "dev-uuid-abc",
          "amountMinorUnits": 5000000,
          "currencyCode": "RUB",
          "requestedAt": "2024-07-11T09:00:00Z",
          "paymentMethodId": "dev-pm-uuid-123"
        }
        ```
    *   Потребители: Payment Service (для обработки выплаты), Admin Service (для финансового контроля).
*   **`com.platform.developer.api_key.created.v1`**
    *   Описание: Разработчик создал новый API ключ.
    *   `data` Payload: `{"developerAccountId": "dev-uuid-abc", "apiKeyId": "devkey-uuid-123", "name": "CI/CD Key", "prefix": "dpk_", "permissions": ["upload_build"], "creationTimestamp": "ISO8601"}`
    *   Потребители: Auth Service (для регистрации ключа, если требуется централизованное управление), Admin Service (для аудита).

### 5.2. Потребляемые События (Consumed Events)
*   **`com.platform.admin.product.moderation.status.changed.v1`** (от Admin Service)
    *   Описание: Статус модерации продукта изменен администратором.
    *   Ожидаемый `data` Payload: `{"productId": "prodsub-uuid-xyz", "newStatus": "approved", "moderatorComment": "Все отлично!", "decisionTimestamp": "ISO8601"}`
    *   Логика обработки: Обновить статус `ProductSubmission`. Уведомить разработчика через Notification Service. Если статус `approved`, разрешить публикацию.
*   **`com.platform.analytics.developer_report.ready.v1`** (от Analytics Service)
    *   Описание: Ежемесячный или запрошенный отчет по аналитике для разработчика готов.
    *   Ожидаемый `data` Payload: `{"developerAccountId": "dev-uuid-abc", "reportId": "report-uuid-analytics", "reportType": "monthly_sales", "s3PathToReport": "path/to/report.pdf", "generationTimestamp": "ISO8601"}`
    *   Логика обработки: Сохранить ссылку на отчет, уведомить разработчика.
*   **`com.platform.payment.payout.status.changed.v1`** (от Payment Service)
    *   Описание: Статус запроса на выплату разработчику изменен.
    *   Ожидаемый `data` Payload: `{"payoutRequestId": "payout-uuid-789" /* ID из Developer Service */, "newStatus": "completed", "transactionId": "payment-txn-uuid", "processedAt": "ISO8601", "details": "Выплата успешно проведена."}`
    *   Логика обработки: Обновить статус `DeveloperPayoutRequest`. Уведомить разработчика.
*   **`com.platform.catalog.product.live_status.changed.v1`** (от Catalog Service)
    *   Описание: Продукт, поданный разработчиком, был опубликован в основном каталоге или снят с публикации.
    *   Ожидаемый `data` Payload: `{"catalogProductId": "cat-prod-uuid-123", "developerSuppliedId": "prodsub-uuid-xyz", "newLiveStatus": "published", "changeTimestamp": "ISO8601"}`
    *   Логика обработки: Обновить соответствующий статус в Developer Service, чтобы разработчик видел актуальное состояние своего продукта в магазине.

## 6. Интеграции (Integrations)
См. `../../../../project_integrations.md` для общей карты и деталей.
*   **Auth Service:** Для аутентификации пользователей Портала Разработчика и управления Developer API ключами. Developer Service вызывает gRPC методы Auth Service.
*   **Catalog Service:** Developer Service передает метаданные продуктов, цены, информацию о билдах в Catalog Service для публикации. Может запрашивать у Catalog Service текущий статус продуктов. Взаимодействие через Kafka и/или gRPC/REST.
*   **Admin Service:** Получает от Developer Service продукты на модерацию (через Kafka). Developer Service потребляет события об изменении статуса модерации от Admin Service.
*   **Payment Service:** Developer Service отправляет запросы на выплаты в Payment Service (через Kafka или API). Получает обновления статуса выплат.
*   **Analytics Service:** Developer Service запрашивает агрегированные данные и отчеты по продуктам разработчика у Analytics Service через его API для отображения в Портале Разработчика.
*   **Download Service:** Developer Service уведомляет Download Service о новых доступных билдах (через Kafka или API), передавая информацию о местоположении в S3.
*   **Notification Service:** Используется для отправки различных уведомлений разработчикам (статус модерации, финансовые операции, важные объявления платформы).
*   **S3-совместимое хранилище:** Для хранения билдов игр, медиа-контента, документов верификации.

## 7. Конфигурация (Configuration)
Общие стандарты конфигурационных файлов (формат YAML, структура, управление переменными окружения и секретами) определены в `../../../../project_api_standards.md` (раздел 7) и `../../../../DOCUMENTATION_GUIDELINES.md` (раздел 6).

### 7.1. Переменные Окружения (Примеры)
*   `DEVELOPER_HTTP_PORT`: Порт для REST API Developer Service (например, `8082`)
*   `POSTGRES_DSN_DEVELOPER`: Строка подключения к PostgreSQL.
*   `REDIS_ADDR_DEVELOPER`: Адрес Redis.
*   `KAFKA_BROKERS_DEVELOPER`: Список брокеров Kafka.
*   `S3_ENDPOINT_BUILDS`, `S3_ACCESS_KEY_BUILDS`, `S3_SECRET_KEY_BUILDS`, `S3_BUCKET_GAME_BUILDS`.
*   `S3_BUCKET_PRODUCT_MEDIA`, `S3_BUCKET_DEVELOPER_DOCS`.
*   `AUTH_SERVICE_GRPC_ADDR`, `CATALOG_SERVICE_API_ADDR`, `PAYMENT_SERVICE_API_ADDR`, `ANALYTICS_SERVICE_API_ADDR`.
*   `LOG_LEVEL_DEVELOPER`: Уровень логирования.
*   `JWT_PUBLIC_KEY_PATH`: Путь к публичному ключу для валидации JWT.

### 7.2. Файлы Конфигурации (`configs/developer_service_config.yaml`)
```yaml
http_server:
  port: ${DEVELOPER_HTTP_PORT:"8082"}
  timeout_seconds: 30
postgres:
  dsn: ${POSTGRES_DSN_DEVELOPER}
  pool_max_conns: 10
redis:
  address: ${REDIS_ADDR_DEVELOPER}
  password: ${REDIS_PASSWORD_DEVELOPER:""}
  db: ${REDIS_DB_DEVELOPER:1}
kafka:
  brokers: ${KAFKA_BROKERS_DEVELOPER}
  producer_topics:
    developer_events: ${KAFKA_TOPIC_DEVELOPER_EVENTS:"com.platform.developer.events.v1"}
  consumer_topics:
    admin_moderation_events: "com.platform.admin.product.moderation.status.changed.v1"
    analytics_reports_events: "com.platform.analytics.developer_report.ready.v1"
    payment_payout_events: "com.platform.payment.payout.status.changed.v1"
  consumer_group: "developer-service-group"
s3_storage:
  builds:
    endpoint: ${S3_ENDPOINT_BUILDS}
    access_key: ${S3_ACCESS_KEY_BUILDS}
    secret_key: ${S3_SECRET_KEY_BUILDS}
    bucket: ${S3_BUCKET_GAME_BUILDS}
    region: "ru-central1"
  media:
    bucket: ${S3_BUCKET_PRODUCT_MEDIA}
  documents:
    bucket: ${S3_BUCKET_DEVELOPER_DOCS}
integrations:
  auth_service_grpc_addr: ${AUTH_SERVICE_GRPC_ADDR}
  # ... другие адреса
logging:
  level: ${LOG_LEVEL_DEVELOPER:"info"}
security:
  jwt_public_key_path: ${JWT_PUBLIC_KEY_PATH}
  max_build_size_gb: 50
  max_media_file_size_mb: 100
```

## 8. Обработка Ошибок (Error Handling)
*   Используются стандартные HTTP коды состояния для REST API и форматы ошибок согласно `../../../../project_api_standards.md`.
*   Внутренние ошибки сервиса логируются с высоким уровнем детализации, включая `trace_id`.
*   Пользователям Портала Разработчика отображаются понятные сообщения об ошибках.
### 8.1. Распространенные Коды Ошибок (специфичные для Developer Service)
*   **`DEVELOPER_ACCOUNT_NOT_FOUND`**: Аккаунт разработчика не найден.
*   **`PRODUCT_SUBMISSION_NOT_FOUND`**: Черновик/заявка на продукт не найдена.
*   **`BUILD_UPLOAD_FAILED`**: Ошибка при загрузке билда.
*   **`PAYOUT_REQUEST_INVALID_AMOUNT`**: Некорректная сумма для запроса выплаты.
*   **`MAX_PRODUCTS_LIMIT_REACHED`**: Достигнут лимит на количество продуктов для разработчика.
*   **`UNSUPPORTED_BUILD_PLATFORM`**: Указанная платформа для билда не поддерживается.
*   **`LEGAL_AGREEMENT_NOT_ACCEPTED`**: Требуется принятие актуального юридического соглашения.

## 9. Безопасность (Security)
См. `../../../../project_security_standards.md` для общих стандартов.
*   **Аутентификация:** Пользователи Портала Разработчика аутентифицируются через Auth Service (JWT). API ключи используются для автоматизированных систем.
*   **Авторизация:** RBAC внутри команды разработчика. Developer Service проверяет права доступа для всех операций.
*   **Защита данных:**
    *   **ФЗ-152 "О персональных данных":** Developer Service обрабатывает ПДн разработчиков (ФИО контактных лиц, email, телефон, юридические и банковские реквизиты ИП/самозанятых). Необходимо обеспечить шифрование этих данных при хранении (например, `details_encrypted` в `DeveloperPaymentMethod`) и передаче, строгое управление доступом, получение согласий на обработку ПДн.
    *   Загружаемые билды и медиа-файлы должны сканироваться на вирусы и вредоносное ПО.
    *   Безопасное хранение API ключей (только хеши).
*   **Защита от атак:** Стандартные меры (валидация ввода, CSRF-токены для веб-форм, защита от XSS). Rate limiting на API.

## 10. Развертывание (Deployment)
(Содержимое существующего раздела актуально, ссылки на `../../../../project_deployment_standards.md` актуальны).

## 11. Мониторинг и Логирование (Logging and Monitoring)
(Содержимое существующего раздела актуально, ссылки на `../../../../project_observability_standards.md` актуальны).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность Портала Разработчика:** P95 < 500мс для большинства операций. Загрузка дашбордов с аналитикой P95 < 2с.
*   **Процесс загрузки билдов:** Должен поддерживать файлы до 50-100 ГБ ([Максимальный размер билда: будет определен на основе технических возможностей и политик платформы, целевой диапазон 50-100 ГБ].). Скорость загрузки зависит от S3 и сети клиента.
*   **Надежность:** Доступность сервиса > 99.9%.
*   **Масштабируемость:** Поддержка до 10,000+ активных разработчиков/издателей и 100,000+ продуктов.

## 13. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API и Protobuf определения для gRPC API (если будут использоваться) поддерживаются в соответствующих репозиториях исходного кода сервиса и/или в централизованном репозитории `platform-protos`.
*   DDL схемы базы данных управляются через систему миграций (например, `golang-migrate/migrate`) и хранятся в репозитории исходного кода сервиса. Актуальные версии доступны во внутренней документации команды разработки и в GitOps репозитории.

## 14. Пользовательские Сценарии (User Flows)

В этом разделе описаны ключевые пользовательские сценарии, связанные с Developer Service.

### 14.1. Регистрация Нового Разработчика/Издателя и Настройка Профиля
*   **Описание:** Новый пользователь платформы решает стать разработчиком/издателем, регистрирует аккаунт разработчика, заполняет юридическую информацию и настраивает профиль команды.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor User
        participant DevPortal as Developer Portal
        participant AuthSvc as Auth Service
        participant DevSvc as Developer Service
        participant Kafka as Kafka Message Bus
        participant AdminSvc as Admin Service
        participant NotificationSvc as Notification Service

        User->>DevPortal: Нажимает "Стать разработчиком"
        DevPortal->>AuthSvc: (Если не аутентифицирован) Процесс входа/регистрации пользователя платформы
        AuthSvc-->>DevPortal: JWT пользователя платформы
        DevPortal->>DevSvc: POST /api/v1/developer/accounts (данные компании/ИП)
        DevSvc->>DevSvc: Валидация данных, создание DeveloperAccount (status: 'pending_verification'), привязка UserID как owner
        DevSvc->>KafkaBus: Publish `com.platform.developer.account.created.v1`
        DevSvc-->>DevPortal: HTTP 201 Created (developer_account_id)
        DevPortal-->>User: Сообщение об успешной подаче заявки и необходимости верификации

        KafkaBus-->>AdminSvc: Consume `developer.account.created.v1` -> Задача на верификацию
        KafkaBus-->>NotificationSvc: Consume `developer.account.created.v1` -> Email разработчику о статусе
    ```

### 14.2. Разработчик Подает Новую Игру на Рассмотрение
*   **Описание:** Разработчик создает черновик игры, заполняет все метаданные, загружает билды и медиа, устанавливает цены, а затем отправляет игру на модерацию.
*   **Диаграмма:** (См. диаграмму "New Game Submission and Approval Process" в разделе 2 или 5, если она там размещена, или адаптировать сюда).
    ```mermaid
    sequenceDiagram
        actor Developer
        participant DevPortal as Developer Portal
        participant DevSvc as Developer Service
        participant S3Store as S3 Storage
        participant Kafka as Kafka Message Bus
        participant CatalogSvc as Catalog Service
        participant AdminSvc as Admin Service

        Developer->>DevPortal: Создает новую игру (POST /products) -> получает product_submission_id
        Developer->>DevPortal: Редактирует метаданные (PUT /products/{id}, draft_metadata)
        Developer->>DevPortal: Загружает билд (POST /products/{id}/versions/.../upload-url -> PUT to S3 -> POST .../upload-complete)
        Developer->>DevPortal: Устанавливает цены (PUT /products/{id}, draft_pricing)
        Developer->>DevPortal: Нажимает "Отправить на модерацию" (POST /products/{id}/submit-for-review)
        DevSvc->>DevSvc: Меняет статус ProductSubmission на 'in_review'
        DevSvc->>Kafka: Publish `com.platform.developer.game.submitted.v1` (product_submission_id)

        CatalogSvc->>Kafka: Consume `game.submitted.v1`
        CatalogSvc->>CatalogSvc: Создает/обновляет запись в своем каталоге со статусом 'in_review'
        CatalogSvc->>Kafka: Publish `com.platform.catalog.product.moderation.required.v1`

        AdminSvc->>Kafka: Consume `product.moderation.required.v1`
        AdminSvc->>AdminSvc: Создание задачи на модерацию
        Note over AdminSvc: Модератор проверяет игру.
        AdminSvc->>Kafka: Publish `com.platform.admin.product.moderation.decision.v1` (decision: 'approved'/'rejected')

        DevSvc->>Kafka: Consume `product.moderation.decision.v1`
        DevSvc->>DevSvc: Обновляет статус ProductSubmission
        DevSvc->>DevPortal: (Через WebSocket или SSE) Уведомление разработчика о решении
    ```

### 14.3. Разработчик Загружает Новую Версию/Билд для Опубликованной Игры
*   **Описание:** Для уже опубликованной игры разработчик создает новую версию, загружает для нее билды и отправляет на модерацию (если требуется для обновлений).
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor Developer
        participant DevPortal as Developer Portal
        participant DevSvc as Developer Service
        participant S3Store as S3 Storage
        participant Kafka as Kafka Message Bus

        Developer->>DevPortal: Выбирает опубликованную игру, нажимает "Создать новую версию"
        DevPortal->>DevSvc: POST /products/{product_id}/versions (version_name, changelog)
        DevSvc-->>DevPortal: Ответ (version_id)
        Developer->>DevPortal: Для новой версии загружает билды (аналогично п.14.2)
        Developer->>DevPortal: Нажимает "Отправить версию на модерацию/публикацию"
        DevSvc->>DevSvc: Обновляет статус ProductVersion
        DevSvc->>Kafka: Publish `com.platform.developer.game_version.submitted.v1` (или аналогичное событие)
        Note over DevSvc, Kafka: Дальнейший процесс модерации и публикации аналогичен п.14.2.
    ```

### 14.4. Разработчик Просматривает Аналитику Продаж и Использования
*   **Описание:** Разработчик заходит в Портал Разработчика, чтобы посмотреть статистику по своим продуктам.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor Developer
        participant DevPortal as Developer Portal
        participant DevSvc as Developer Service
        participant AnalyticsSvc as Analytics Service

        Developer->>DevPortal: Переходит в раздел "Аналитика" для игры X
        DevPortal->>DevSvc: GET /api/v1/developer/analytics/products/{game_id}/summary?period=last_30_days
        DevSvc->>DevSvc: Проверка прав доступа разработчика к игре X
        DevSvc->>AnalyticsSvc: (gRPC/REST) Запрос агрегированных данных (developer_id, game_id, period)
        AnalyticsSvc-->>DevSvc: Данные аналитики (продажи, DAU, MAU и т.д.)
        DevSvc-->>DevPortal: HTTP 200 OK (данные аналитики)
        DevPortal-->>Developer: Отображение дашбордов и графиков
    ```

### 14.5. Разработчик Инициирует Запрос на Выплату Средств
*   **Описание:** Разработчик, накопив достаточную сумму на балансе, формирует запрос на выплату на свои банковские реквизиты.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor Developer
        participant DevPortal as Developer Portal
        participant DevSvc as Developer Service
        participant Kafka as Kafka Message Bus
        participant PaymentSvc as Payment Service

        Developer->>DevPortal: Переходит в раздел "Финансы" -> "Выплаты"
        DevPortal->>DevSvc: GET /finance/balance (получение доступного баланса)
        DevSvc-->>DevPortal: Баланс
        Developer->>DevPortal: Заполняет форму запроса на выплату (сумма, платежный метод)
        DevPortal->>DevSvc: POST /api/v1/developer/finance/payouts/requests (сумма, payment_method_id)
        DevSvc->>DevSvc: Валидация запроса (баланс, лимиты, статус аккаунта разработчика)
        DevSvc->>DevSvc: Создание DeveloperPayoutRequest (status: 'pending_approval' или 'pending_request')
        DevSvc->>Kafka: Publish `com.platform.developer.payout.requested.v1`
        DevSvc-->>DevPortal: HTTP 201 Created (статус запроса)
        DevPortal-->>Developer: Уведомление о создании запроса на выплату

        PaymentSvc->>Kafka: Consume `payout.requested.v1`
        Note over PaymentSvc: Дальнейшая обработка выплаты Payment Service и Admin Service.
    ```

### 14.6. Разработчик Управляет Членами Команды и Разрешениями
*   **Описание:** Владелец или администратор аккаунта разработчика добавляет нового участника в команду и назначает ему роли/разрешения.
*   **Диаграмма:**
    ```mermaid
    sequenceDiagram
        actor DevAdmin as Администратор Команды Разработчика
        participant DevPortal as Developer Portal
        participant DevSvc as Developer Service
        participant UserDB_Auth as Auth Service (User DB - для поиска пользователя)
        participant NotificationSvc as Notification Service

        DevAdmin->>DevPortal: Переходит в "Управление командой" -> "Добавить участника"
        DevPortal->>DevSvc: POST /api/v1/developer/accounts/me/team/members (email_приглашаемого, роль_в_команде)
        DevSvc->>UserDB_Auth: (Через AuthSvc gRPC) Поиск пользователя по email
        alt Пользователь найден на платформе
            DevSvc->>DevSvc: Создание записи DeveloperTeamMember (статус 'pending_invitation' или сразу 'active')
            DevSvc->>NotificationSvc: (Через Kafka/API) Отправка приглашения пользователю на email
            DevSvc-->>DevPortal: HTTP 201 Created / HTTP 200 OK
            DevPortal-->>DevAdmin: Участник приглашен / добавлен
        else Пользователь не найден
            DevSvc-->>DevPortal: HTTP 404 Not Found (USER_NOT_FOUND_ON_PLATFORM)
            DevPortal-->>DevAdmin: Ошибка: пользователь не найден на платформе
        end
    ```

## 15. Резервное Копирование и Восстановление (Backup and Recovery)

### 15.1. PostgreSQL (Данные аккаунтов разработчиков, черновики продуктов, запросы на выплаты и т.д.)
*   **Процедура резервного копирования:**
    *   **Логические бэкапы:** Ежедневный `pg_dump` для базы данных Developer Service.
    *   **Физические бэкапы (PITR):** Настроена непрерывная архивация WAL-сегментов. Базовый бэкап создается еженедельно.
    *   **Хранение:** Бэкапы и WAL-архивы хранятся в S3-совместимом хранилище с шифрованием и версионированием, в другом регионе. Срок хранения: полные логические бэкапы - 30 дней, WAL - 14 дней.
*   **Процедура восстановления:** Тестируется ежеквартально.
*   **RTO (Recovery Time Objective):** < 2 часов.
*   **RPO (Recovery Point Objective):** < 15 минут.
*   (Общие принципы см. `../../../../project_database_structure.md`).

### 15.2. S3-совместимое хранилище (Билды игр, медиа-активы, документы верификации)
*   **Процедура резервного копирования:**
    *   **Версионирование объектов:** Включено для всех бакетов.
    *   **Политики жизненного цикла (Lifecycle Policies):** Настроены для управления старыми версиями и, возможно, для перемещения неактивных билдов в более холодные классы хранения (если применимо).
    *   **Cross-Region Replication (CRR):** Настроена для бакетов с билдами игр и критически важными медиа-активами для обеспечения гео-резервирования. Для документов верификации также рекомендуется CRR.
*   **Процедура восстановления:**
    *   Восстановление отдельных объектов или версий из S3.
    *   В случае регионального сбоя – переключение на реплицированный бакет в другом регионе.
*   **RTO:** Зависит от объема данных и скорости S3, но обычно быстро для отдельных файлов. Полное восстановление всех данных может занять часы.
*   **RPO:** Близко к нулю при использовании версионирования и CRR (ограничено временем репликации S3).

### 15.3. Redis (Кэш, сессии Портала Разработчика)
*   **Стратегия:** Данные в Redis в основном являются кэшем или временными сессионными данными.
*   **Персистентность (опционально):** Может быть включена RDB-снапшотирование и/или AOF для ускорения восстановления после перезапуска Redis.
*   **Резервное копирование:** Не является критичным для большинства данных, так как они могут быть перестроены из PostgreSQL или пересозданы пользователем.
*   **RTO/RPO:** Неприменимо в контексте долгосрочного хранения данных. Восстановление функциональности кэша происходит по мере его заполнения.

### 15.4. Общая стратегия
*   Приоритет отдается восстановлению данных из PostgreSQL и S3.
*   Процедуры восстановления тестируются и документируются.
*   Мониторинг процессов резервного копирования.

## 16. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API и Protobuf определения для gRPC API (если будут использоваться) поддерживаются в соответствующих репозиториях исходного кода сервиса и/или в централизованном репозитории `platform-protos`.
*   DDL схемы базы данных управляются через систему миграций (например, `golang-migrate/migrate`) и хранятся в репозитории исходного кода сервиса. Актуальные версии доступны во внутренней документации команды разработки и в GitOps репозитории.

## 17. Связанные Рабочие Процессы (Related Workflows)
*   [Подача разработчиком новой игры на модерацию](../../../../project_workflows/game_submission_flow.md)
*   [Процесс выплат разработчикам] Подробное описание этого рабочего процесса будет добавлено в [developer_payout_flow.md](../../../../project_workflows/developer_payout_flow.md) (документ в разработке).

---
*Этот документ является основной спецификацией для Developer Service и должен поддерживаться в актуальном состоянии.*
