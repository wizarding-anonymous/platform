# Спецификация Микросервиса: Developer Service

**Версия:** 1.0
**Дата последнего обновления:** 2025-05-25

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Developer Service, предназначенного для разработчиков и издателей игр на платформе.
*   **Роль в общей архитектуре платформы:** Developer Service предоставляет интерфейс и функциональность для управления аккаунтами разработчиков, загрузки и обновления игр, управления метаданными и маркетинговыми материалами, доступа к аналитике продаж и использования продуктов, а также взаимодействия с финансовыми аспектами платформы.
*   **Основные бизнес-задачи:** Обеспечение инструментов для публикации и поддержки игр, управление жизненным циклом продуктов разработчиков, предоставление аналитики и финансовых отчетов.

### 1.2. Ключевые Функциональности
*   Регистрация и управление аккаунтом разработчика/издателя (профиль компании, управление командой и ролями).
*   Управление проектами (играми): создание, загрузка билдов, управление версиями, метаданными (описания, теги, жанры, системные требования, рейтинги), медиа-контентом.
*   Управление ценообразованием: установка базовой и региональных цен, создание и управление скидками.
*   Процесс публикации: управление стадиями жизненного цикла игры (черновик, на рассмотрении, опубликовано).
*   Панель аналитики: дашборды с метриками (продажи, доход, установки, DAU/MAU), генерация отчетов.
*   Финансовый раздел: баланс, история транзакций, управление реквизитами, запросы на выплаты.
*   Управление SDK и API ключами: доступ к SDK, документации, управление ключами API для CI/CD.
*   Система уведомлений: о статусе модерации, обновлениях, финансовых операциях.
*   (Источник: Спецификация Developer Service, раздел 2.3)

### 1.3. Основные Технологии
*   **Язык программирования:** Go / Python / Java / Node.js (будет согласован).
*   **Фреймворк:** Gin (Go), Django/Flask (Python), Spring Boot (Java), Express (Node.js).
*   **База данных:** PostgreSQL (реляционные данные), возможно NoSQL (MongoDB) для черновиков/сложных метаданных.
*   **Кэширование:** Redis.
*   **Очереди сообщений:** Kafka или RabbitMQ.
*   **Хранилище файлов:** S3-совместимое объектное хранилище.
*   **Инфраструктура:** Docker, Kubernetes.
*   (Источник: Спецификация Developer Service, раздел 3.3)

### 1.4. Термины и Определения (Glossary)
*   Для специфичных терминов см. "Единый глоссарий терминов и определений для российского аналога Steam.txt".
*   Ключевые понятия: "Аккаунт разработчика", "Билд игры", "Метаданные игры", "Портал разработчика".

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Developer Service может быть реализован как набор микросервисов или модульный монолит.
*   Архитектура будет следовать принципам чистой архитектуры или слоистой архитектуры.
*   Компоненты включают: API Gateway, Модуль управления аккаунтами разработчиков, Модуль управления играми, Модуль загрузки контента, Модуль аналитики, Модуль финансов, Модуль SDK и API, Веб-интерфейс (Портал разработчика).
*   (Источник: Спецификация Developer Service, разделы 3.1, 3.2)

### 2.2. Слои Сервиса
(Предполагаемая структура на основе общих принципов)

#### 2.2.1. Presentation Layer (Слой Представления)
*   Ответственность: Обработка входящих HTTP запросов от Портала разработчика и публичного API для разработчиков. Валидация, DTO.
*   Ключевые компоненты/модули:
    *   HTTP Handlers/Controllers: Для каждого функционального модуля (управление играми, финансы и т.д.).
    *   DTOs: Для запросов/ответов API.

#### 2.2.2. Application Layer (Прикладной Слой / Слой Сценариев Использования)
*   Ответственность: Координация бизнес-логики для управления продуктами разработчиков, аналитикой, финансами.
*   Ключевые компоненты/модули:
    *   Use Case Services: Например, `SubmitGameForModerationService`, `RequestPayoutService`, `GetGameAnalyticsService`.

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности (DeveloperAccount, GameProject, GameVersion, PayoutRequest), бизнес-правила.
*   Ключевые компоненты/модули:
    *   Entities: `Developer`, `DeveloperTeamMember`, `Game`, `GameVersion`, `GameMetadata`, `GamePricing`, `DeveloperPayout`, `DeveloperAPIKey`.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Взаимодействие с PostgreSQL, S3, Kafka/RabbitMQ, Redis, а также с другими микросервисами (Account, Payment, Catalog, Download, Analytics, Admin, Notification Services).
*   Ключевые компоненты/модули:
    *   Database Repositories.
    *   S3 Client: Для загрузки/управления билдами и медиа.
    *   Message Queue Producers/Consumers.
    *   Clients for other microservices.

## 3. API Endpoints

### 3.1. REST API
*   **Префикс:** `/api/v1/developer`
*   **Аутентификация:** JWT (через Auth Service).
*   **Авторизация:** На основе `developer_id` и роли пользователя в команде.
*   **Основные группы эндпоинтов:**
    *   Аккаунты разработчиков (`POST /accounts`, `GET /accounts/me`, `GET /accounts/me/team`).
    *   Игры (`POST /games`, `GET /games`, `GET /games/{game_id}`).
    *   Версии игр (`POST /games/{game_id}/versions`, `GET /games/{game_id}/versions`).
    *   Метаданные игр (`PUT /games/{game_id}/metadata`, `POST /games/{game_id}/media`).
    *   Ценообразование (`PUT /games/{game_id}/pricing`).
    *   Публикация (`POST /games/{game_id}/versions/{version_id}/submit`, `POST /games/{game_id}/publish`).
    *   Аналитика (`GET /analytics/summary`, `GET /analytics/games/{game_id}`).
    *   Финансы (`GET /finance/balance`, `POST /finance/payouts`).
    *   API Ключи (`POST /apikeys`, `GET /apikeys`).
*   (Более полный список см. в Спецификации Developer Service, раздел 5.2).

### 3.2. gRPC API
*   Information not found in existing documentation. (Может использоваться для внутренней коммуникации, если Developer Service реализован как несколько микросервисов).

### 3.3. WebSocket API (если применимо)
*   Information not found in existing documentation. (Может использоваться для real-time уведомлений в Портале разработчика).

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности
*   **Developer**: Компания-разработчик/издатель.
*   **DeveloperTeamMember**: Член команды разработчика (связь с UserID из Account Service).
*   **Game**: Проект игры.
*   **GameVersion**: Версия билда игры (путь к файлу в S3, хэш).
*   **GameMetadata**: Метаданные игры (описания, теги, жанры, системные требования, медиа).
*   **GamePricing**: Ценовая информация.
*   **DeveloperPayout**: Запрос на выплату.
*   **DeveloperAPIKey**: API ключ для доступа к публичному API разработчика.
*   (Основные таблицы PostgreSQL см. в Спецификации Developer Service, раздел 5.1).

### 4.2. Схема Базы Данных
*   **PostgreSQL**: Хранит структурированную информацию (профили разработчиков, метаданные игр, финансы).
    ```sql
    -- Пример таблицы developers (сокращенно)
    CREATE TABLE developers ( developer_id UUID PRIMARY KEY ..., company_name VARCHAR(255) NOT NULL ...);
    -- Пример таблицы games (сокращенно)
    CREATE TABLE games ( game_id UUID PRIMARY KEY ..., developer_id UUID REFERENCES developers(developer_id) ..., title VARCHAR(255) ...);
    -- ... и другие таблицы ...
    ```
*   **S3-совместимое хранилище**: Хранит бинарные данные (билды игр, медиа-файлы, SDK).
*   (Полный список таблиц и их полей см. в Спецификации Developer Service, раздел 5.1).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Kafka или RabbitMQ.
*   **События, которые Developer Service может публиковать:**
    *   `developer.game.submitted`: Новая игра/версия отправлена на модерацию -> Admin Service.
    *   `developer.game.published`: Игра опубликована -> Catalog Service, Download Service, Notification Service.
    *   `developer.game.updated`: Метаданные/цена игры обновлены -> Catalog Service.
        *   `Структура Payload (пример):`
            ```json
            {
              "developer_id": "uuid_developer_account",
              "game_id": "uuid_game",
              "updated_at": "ISO8601_timestamp",
              "updated_fields": ["metadata", "pricing"]
            }
            ```
    *   `developer.payout.requested`: Запрошена выплата -> Payment Service.
        *   `Структура Payload (пример):`
            ```json
            {
              "payout_request_id": "uuid_payout_request",
              "developer_id": "uuid_developer_account",
              "amount": 150000.75,
              "currency": "RUB",
              "requested_at": "ISO8601_timestamp",
              "payment_details_snapshot": {
                "account_type": "bank_transfer",
                "beneficiary_name": "ООО Разработчик Игр"
              }
            }
            ```

### 5.2. Потребляемые События (Consumed Events)
*   **События, которые Developer Service может потреблять:**
    *   `admin.game.moderation.approved`: Игра одобрена модерацией <- Admin Service.
        *   `Структура Payload (пример):`
            ```json
            {
              "game_id": "uuid_game",
              "version_id": "uuid_game_version",
              "moderator_id": "uuid_admin_user",
              "decision": "approved",
              "approved_at": "ISO8601_timestamp",
              "comments": "Игра соответствует всем правилам."
            }
            ```
        *   `Логика обработки:` Обновить статус игры/версии в Developer Service на "approved" или "ready_for_publish". Отправить уведомление разработчику.
    *   `admin.game.moderation.rejected`: Игра отклонена модерацией <- Admin Service.
        *   `Структура Payload (пример):`
            ```json
            {
              "game_id": "uuid_game",
              "version_id": "uuid_game_version",
              "moderator_id": "uuid_admin_user",
              "decision": "rejected",
              "rejected_at": "ISO8601_timestamp",
              "reason": "Ненормативный контент в описании.",
              "detailed_feedback_url": "link_to_admin_panel_feedback"
            }
            ```
        *   `Логика обработки:` Обновить статус игры/версии на "rejected". Сохранить причину отклонения. Отправить уведомление разработчику.
    *   `payment.payout.status.changed`: Статус выплаты изменен <- Payment Service.
        *   `Структура Payload (пример):`
            ```json
            {
              "payout_request_id": "uuid_payout_request_in_developer_service",
              "developer_id": "uuid_developer_account",
              "payment_service_transaction_id": "uuid_payment_service_transaction",
              "new_status": "completed" | "failed" | "processing",
              "processed_at": "ISO8601_timestamp",
              "failure_reason": "Неверные реквизиты"
            }
            ```
        *   `Логика обработки:` Обновить статус запроса на выплату. Уведомить разработчика.
    *   `analytics.report.ready`: Отчет по аналитике готов <- Analytics Service.
        *   `Структура Payload (пример):`
            ```json
            {
              "report_id": "uuid_report_instance",
              "developer_id": "uuid_developer_account",
              "report_type": "monthly_sales_summary",
              "game_id": "uuid_game",
              "period_start": "ISO8601_date",
              "period_end": "ISO8601_date",
              "report_url": "url_to_download_or_view_report",
              "generated_at": "ISO8601_timestamp"
            }
            ```
        *   `Логика обработки:` Сохранить ссылку на отчет. Уведомить разработчика о готовности отчета.
    *   `notification.message.received`: Ответ от техподдержки или важное системное уведомление <- Notification Service.
        *   `Структура Payload (пример):`
             ```json
            {
              "notification_id": "uuid",
              "recipient_type": "developer_team",
              "recipient_id": "uuid_developer_account_or_team_member_user_id",
              "category": "support_reply" | "platform_announcement" | "game_status_update",
              "title": "Получен ответ от службы поддержки по тикету #12345",
              "short_message": "Специалист поддержки ответил на ваш запрос...",
              "details_url": "/developer/support/tickets/12345",
              "sent_at": "ISO8601_timestamp"
            }
            ```
        *   `Логика обработки:` Отобразить уведомление в интерфейсе портала разработчика.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Account Service**: Управление базовой информацией об аккаунтах разработчиков.
*   **Payment Service**: Обработка выплат, финансовые отчеты.
*   **Catalog Service**: Получение данных об играх для публикации.
*   **Download Service**: Передача информации о версиях игр и файлах.
*   **Analytics Service**: Предоставление данных для панели аналитики разработчика.
*   **Admin Service**: Процессы модерации контента.
*   **Notification Service**: Отправка уведомлений разработчикам.
*   (Детали см. в Спецификации Developer Service, разделы 1.3 и 6).

### 6.2. Внешние Системы
*   **S3-совместимое хранилище**: Для билдов игр и медиа-файлов.

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `DEVELOPER_SERVICE_PORT`: Порт сервиса.
*   `POSTGRES_DSN`: Строка подключения к PostgreSQL.
*   `S3_ENDPOINT`
*   `S3_ACCESS_KEY_ID`
*   `S3_SECRET_ACCESS_KEY`
*   `S3_BUCKET_GAME_BUILDS`
*   `S3_BUCKET_GAME_MEDIA`
*   `S3_REGION` (optional)
*   `S3_USE_SSL` (e.g., `true`)
*   `KAFKA_BROKERS` (comma-separated) / `RABBITMQ_URL`
*   `LOG_LEVEL` (e.g., `info`, `debug`)
*   `ACCOUNT_SERVICE_GRPC_ADDR`
*   `PAYMENT_SERVICE_GRPC_ADDR`
*   `CATALOG_SERVICE_GRPC_ADDR`
*   `DOWNLOAD_SERVICE_GRPC_ADDR`
*   `ANALYTICS_SERVICE_GRPC_ADDR`
*   `ADMIN_SERVICE_GRPC_ADDR`
*   `NOTIFICATION_SERVICE_GRPC_ADDR`
*   `AUTH_SERVICE_GRPC_ADDR`
*   `MAX_BUILD_FILE_SIZE_BYTES` (e.g., `10737418240` for 10GB)
*   `MAX_MEDIA_FILE_SIZE_BYTES` (e.g., `104857600` for 100MB)
*   `TEMPORARY_UPLOAD_DIR` (e.g., `/tmp/uploads`)
*   `REDIS_ADDR` (e.g., `redis:6379`)
*   `REDIS_PASSWORD` (optional)
*   `REDIS_DB_DEVELOPER` (e.g., `1`)

### 7.2. Файлы Конфигурации (если применимо)
*   Конфигурация сервиса осуществляется преимущественно через переменные окружения. Если в будущем потребуются файлы конфигурации для сложных настроек (например, для определения специфических правил валидации для разных типов продуктов или детализированных настроек API), их структура будет определена здесь.

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   Единообразная и информативная обработка ошибок.
*   Возвращение корректных HTTP статусов.
*   Логирование всех ошибок.

### 8.2. Распространенные Коды Ошибок
*   **400 Bad Request**: Некорректные входные данные.
*   **401 Unauthorized**: Ошибка аутентификации.
*   **403 Forbidden**: Недостаточно прав (например, попытка редактировать чужую игру).
*   **404 Not Found**: Ресурс не найден (игра, версия).
*   **500 Internal Server Error**: Внутренняя ошибка.
*   **502 Bad Gateway / 503 Service Unavailable**: Ошибки взаимодействия с другими сервисами.

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   JWT через Auth Service для доступа к порталу и API.

### 9.2. Авторизация
*   RBAC внутри команды разработчика (администратор, разработчик, маркетолог, финансист).
*   Проверка принадлежности игры к аккаунту разработчика.

### 9.3. Защита Данных
*   Шифрование конфиденциальных данных (реквизиты, API ключи). HTTPS.
*   Безопасная загрузка билдов (антивирусная проверка).

### 9.4. Управление Секретами
*   API ключи разработчиков хранятся в виде хэшей.
*   Секреты сервиса через Kubernetes Secrets или Vault.
*   **Аудит**: Логирование важных действий.
*   (Детали см. в Спецификации Developer Service, раздел 7.1).

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Стандартный Dockerfile.
*   **Kubernetes манифесты/Helm-чарты.**

### 10.2. Зависимости при Развертывании
*   PostgreSQL, S3, Kafka/RabbitMQ, Redis.
*   Account, Auth, Payment, Catalog, Download, Analytics, Admin, Notification Services.

### 10.3. CI/CD
*   Автоматическая сборка, тестирование, развертывание.
*   (Детали см. в Спецификации Developer Service, раздел 8.2).

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   Структурированные логи (JSON).
*   Интеграция с ELK/Loki.

### 11.2. Мониторинг
*   Prometheus, Grafana.
*   Метрики: запросы, ошибки, производительность загрузок, использование ресурсов.

### 11.3. Трассировка
*   Интеграция с системой распределенной трассировки (например, Jaeger/OpenTelemetry) будет реализована согласно общепроектным стандартам. Контекст трассировки будет передаваться для всех входящих API запросов и исходящих вызовов к другим сервисам, а также включаться в логи.

## 12. Нефункциональные Требования (NFRs)
*   **Безопасность**: Защита билдов и данных разработчиков, RBAC.
*   **Производительность**: Обработка загрузки больших файлов, быстрый отклик портала.
*   **Масштабируемость**: Горизонтальное масштабирование, масштабируемость хранилища файлов.
*   **Надежность**: Доступность >= 99.8%, сохранность данных.
*   **Удобство использования**: Интуитивный портал разработчика.
*   (Детали см. в Спецификации Developer Service, раздел 2.4).

## 13. Приложения (Appendices) (Опционально)
*   Детальные схемы DDL для PostgreSQL, полные примеры API запросов/ответов (включая структуры для загрузки медиа-файлов и билдов), а также форматы событий Kafka/RabbitMQ будут добавлены по мере финализации дизайна и реализации соответствующих модулей.

---
*Этот шаблон является отправной точкой и может быть адаптирован под конкретные нужды проекта и сервиса.*
