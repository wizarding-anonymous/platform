# Спецификация Микросервиса: Notification Service

**Версия:** 1.0
**Дата последнего обновления:** 2025-05-25

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Notification Service.
*   **Роль в общей архитектуре платформы:** Notification Service отвечает за централизованную отправку и управление всеми видами уведомлений внутри платформы (транзакционные, системные, маркетинговые, от других сервисов). Абстрагирует сложность взаимодействия с различными провайдерами и каналами доставки.
*   **Основные бизнес-задачи:** Обеспечение надежной, своевременной и релевантной доставки уведомлений, управление пользовательскими подписками, поддержка маркетинговых кампаний.
*   (Источник: Спецификация Notification Service, разделы 1.1, 1.2, 2.1)

### 1.2. Ключевые Функциональности
*   Поддержка различных каналов доставки: Email, Push (FCM/APNS), SMS, In-App.
*   Управление шаблонами уведомлений: создание, редактирование, хранение, локализация, поддержка переменных.
*   Отправка уведомлений: одиночных, массовых, отложенная отправка, приоритизация.
*   Управление пользовательскими предпочтениями: настройка типов уведомлений и каналов, глобальные отписки.
*   Маркетинговые кампании: создание, управление, сегментация аудитории, планирование, A/B тестирование.
*   Отслеживание и статистика: статусы доставки (отправлено, доставлено, ошибка, открыто, кликнуто), статистика по кампаниям.
*   Интеграция с провайдерами: гибкая конфигурация, поддержка нескольких провайдеров, обработка обратной связи.
*   (Источник: Спецификация Notification Service, раздел 2.2)

### 1.3. Основные Технологии
*   **Язык программирования:** Go (предпочтительно) или Java/Kotlin.
*   **Фреймворки:** Gin/Echo (Go), Spring Boot (Java/Kotlin).
*   **Очереди сообщений:** Kafka (предпочтительно) или RabbitMQ.
*   **Шаблонизатор:** text/template, html/template (Go); Thymeleaf, Freemarker (Java).
*   **Базы данных:** PostgreSQL или ClickHouse (для статистики), Redis или Memcached (кэш).
*   **Инфраструктура:** Docker, Kubernetes, Prometheus, Grafana, ELK/Loki.
*   (Источник: Спецификация Notification Service, раздел 3.2)

### 1.4. Термины и Определения (Glossary)
*   **Уведомление (Notification):** Сообщение пользователю или системе.
*   **Канал доставки (Delivery Channel):** Способ доставки (Email, Push, SMS, In-App).
*   **Шаблон (Template):** Предопределенная структура сообщения.
*   **Провайдер (Provider):** Внешний сервис для отправки (Email-провайдер, SMS-шлюз, FCM/APNS).
*   (Полный глоссарий см. в Спецификации Notification Service, раздел 1.3)

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Notification Service построен с использованием событийно-ориентированного подхода.
*   Компоненты включают: API Gateway/Ingress, Message Consumer (Kafka), Notification Orchestrator, Channel Dispatchers (Email, Push, SMS, In-App), Template Engine, Preference Manager, Campaign Manager, Stats Collector, Repository Layer.
*   Диаграмма компонентов:
    ```mermaid
    graph TD
        subgraph Notification Service
            API[API Gateway / Ingress] --> PM(Preference Manager)
            API --> CM(Campaign Manager)
            Consumer[Message Consumer] --> Orch(Notification Orchestrator)
            Orch --> PM; Orch --> TE(Template Engine)
            Orch --> EmailQueue; Orch --> PushQueue; Orch --> SMSQueue; Orch --> InAppQueue
            EmailQueue --> EmailDispatcher; PushQueue --> PushDispatcher; SMSQueue --> SMSDispatcher; InAppQueue --> InAppDispatcher
            EmailDispatcher --> EmailProvider; PushDispatcher --> PushProvider; SMSDispatcher --> SMSProvider; InAppDispatcher --> WebSocketGateway
            EmailDispatcher --> SC(Stats Collector); PushDispatcher --> SC; SMSDispatcher --> SC; InAppDispatcher --> SC
            PM --> DB_Prefs[(Database: Preferences)]; CM --> DB_Camp[(Database: Campaigns)]; TE --> DB_Tmpl[(Database: Templates)]; SC --> DB_Stats[(Database: Stats)]
        end
        OtherServices --> Kafka(Message Broker: Kafka); Kafka --> Consumer
        Users --> API; Admin --> API
        EmailProvider --|> ExternalSystems; PushProvider --|> ExternalSystems; SMSProvider --|> ExternalSystems
        WebSocketGateway --|> ClientApps
    ```
*   (Источник: Спецификация Notification Service, раздел 3.1)

### 2.2. Слои Сервиса
(На основе компонентов, описанных в исходной спецификации)

#### 2.2.1. Presentation Layer (Слой Представления / API Layer & Message Consumer)
*   Ответственность: Прием HTTP-запросов (REST/gRPC) для управления шаблонами, кампаниями, предпочтениями; прием событий из Kafka.
*   Ключевые компоненты/модули: API Gateway/Ingress, Message Consumer (Kafka).

#### 2.2.2. Application Layer (Прикладной Слой / Orchestration & Management)
*   Ответственность: Ядро сервиса. Определение типа уведомления, проверка предпочтений, выбор шаблона и канала, обогащение данными, постановка задач на отправку. Управление кампаниями, предпочтениями, шаблонами.
*   Ключевые компоненты/модули: Notification Orchestrator, Preference Manager, Campaign Manager, Template Engine.

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности и их правила.
*   Ключевые компоненты/модули: `NotificationTemplate`, `NotificationMessage`, `UserPreferences`, `NotificationCampaign`, `DeviceToken`.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой / Dispatching & Persistence)
*   Ответственность: Взаимодействие с провайдерами доставки (Email, Push, SMS), WebSocket Gateway. Взаимодействие с хранилищами данных (PostgreSQL/ClickHouse, Redis). Сбор статистики.
*   Ключевые компоненты/модули: Channel Dispatchers (Email, Push, SMS, In-App), Stats Collector, Repository Layer (PostgreSQL, Redis), клиенты для внешних провайдеров.

## 3. API Endpoints

### 3.1. REST API
*   **Префикс:** `/api/v1/notifications`
*   **Аутентификация:** JWT, API-ключи для сервисов.
*   **Авторизация:** На основе ролей.
*   **Основные группы эндпоинтов:**
    *   Управление шаблонами: `POST /templates`, `GET /templates`, `GET /templates/{id}`, `PUT /templates/{id}`, `DELETE /templates/{id}`.
    *   Управление пользовательскими предпочтениями: `GET /preferences/{user_id}`, `PUT /preferences/{user_id}`.
    *   Управление маркетинговыми кампаниями: `POST /campaigns`, `GET /campaigns`, `GET /campaigns/{id}`, `POST /campaigns/{id}/start`.
    *   Управление устройствами: `POST /devices`, `PUT /devices/{id}`.
    *   Отправка уведомлений (для сервисов/админов): `POST /send`, `POST /send/batch`.
    *   Статистика: `GET /stats/delivery`, `GET /stats/messages/{message_id}`.
    *   Проверка состояния: `GET /health`.
*   (Детали см. в Спецификации Notification Service, раздел 5.1).

### 3.2. gRPC API
*   Опционально, для высокопроизводительного взаимодействия.
*   Пример сервисов: `NotificationService` (SendNotification, GetNotificationStatus).
*   (Детали см. в Спецификации Notification Service, раздел 5.3).

### 3.3. WebSocket API (для In-App уведомлений)
*   Взаимодействие с WebSocket Gateway для доставки In-App уведомлений.
*   (Источник: Спецификация Notification Service, разделы 3.1, 6.2.4)

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности
*   **NotificationTemplate**: Шаблон уведомления (ID, канал, язык, тема, тело, переменные).
*   **NotificationMessage**: Экземпляр уведомления (ID, user_id, template_id, канал, статус, данные).
*   **UserPreferences**: Настройки уведомлений пользователя.
*   **NotificationCampaign**: Маркетинговая кампания.
*   **DeviceToken**: Токены устройств для Push.
*   **ProviderConfig**: Конфигурация провайдеров.
*   (Детали см. в Спецификации Notification Service, раздел 3.3.1).

### 4.2. Схема Базы Данных
*   **PostgreSQL/ClickHouse**: Хранение шаблонов, кампаний, статистики, логов доставки, токенов устройств, предпочтений.
    ```sql
    -- Пример таблицы notification_templates (сокращенно)
    CREATE TABLE notification_templates (id UUID PRIMARY KEY, name VARCHAR(255) NOT NULL, channel_type VARCHAR(20) NOT NULL ...);
    -- Пример таблицы notification_messages (сокращенно)
    CREATE TABLE notification_messages (id BIGSERIAL PRIMARY KEY, user_id UUID NOT NULL, channel_type VARCHAR(20) NOT NULL, status VARCHAR(50) ...);
    ```
*   **Redis/Memcached**: Кэш статусов доставки, пользовательских предпочтений.
*   (Полную схему см. в Спецификации Notification Service, раздел 3.3.2).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Kafka.
*   **Формат событий:** CloudEvents JSON.
*   **Основные топики:** `notification.status.events` (для статусов доставки, открытий, ошибок), `notification.campaign.stats` (для агрегированной статистики по кампаниям).

#### Примеры структур Payload для публикуемых событий:

*   **`notification.message.status.updated.v1`**
    *   `Структура Payload (пример):`
        ```json
        {
          "event_id": "uuid_event",
          "event_type": "notification.message.status.updated.v1",
          "timestamp": "ISO8601_timestamp",
          "source_service": "notification-service",
          "message_id": "uuid_notification_message",
          "user_id": "uuid_user",
          "channel_type": "email" | "push" | "sms" | "in_app",
          "new_status": "sent" | "delivered" | "failed" | "opened" | "clicked",
          "provider_name": "sendgrid" | "fcm" | "smsc",
          "status_details": "Текст ошибки от провайдера или дополнительная информация",
          "campaign_id": "uuid_campaign"
        }
        ```
*   **`notification.campaign.stats.aggregated.v1`**
    *   `Структура Payload (пример):`
        ```json
        {
          "event_id": "uuid_event",
          "event_type": "notification.campaign.stats.aggregated.v1",
          "timestamp": "ISO8601_timestamp",
          "source_service": "notification-service",
          "campaign_id": "uuid_campaign",
          "aggregation_period_hours": 1,
          "metrics": {
            "sent_count": 10000,
            "delivered_count": 9500,
            "failed_count": 500,
            "opened_count": 3000,
            "clicked_count": 500
          }
        }
        ```
*   *(Примечание: Исходная спецификация упоминает топики `notification.events` и `notification.stats`. Предложенные выше события и их структура детализируют это.)*

### 5.2. Потребляемые События (Consumed Events)
*   **Основные топики и типы событий:**
    *   **`notification.send.request.v1`** (из любого сервиса, через Kafka)
        *   **Источник:** Любой микросервис платформы.
        *   **Назначение:** Запрос на отправку уведомления пользователю или группе пользователей.
        *   **Структура Payload (пример):**
            ```json
            {
              "request_id": "uuid_unique_request_id",
              "user_id": "uuid_user",
              "notification_type": "order_confirmation",
              "preferred_channels": ["email", "push"],
              "data": {
                "order_id": "ORD12345",
                "total_amount": "1999.00 RUB",
                "game_name": "Супер Игра"
              },
              "priority": "high",
              "send_at": "ISO8601_timestamp"
            }
            ```
        *   **Логика обработки:** Notification Orchestrator получает событие. Проверяет пользовательские предпочтения (`UserPreferences`). Выбирает соответствующий шаблон (`NotificationTemplate`) для языка пользователя и доступных/предпочтительных каналов. Рендерит шаблон с предоставленными `data`. Ставит задачи в очереди для соответствующих `ChannelDispatchers`.
    *   **`catalog.game.version.deleted.v1`**
        *   **Источник:** Catalog Service
        *   **Назначение:** Уведомление пользователей, которые могли быть заинтересованы в удаленной версии игры.
        *   **Структура Payload (пример):**
            ```json
            {
              "game_id": "uuid_game",
              "version_id": "uuid_game_version",
              "game_title": "Удаленная Игра",
              "reason_for_deletion": "Юридические ограничения"
            }
            ```
        *   **Логика обработки:** Определить затронутых пользователей. Создать и отправить им уведомление (тип `system.game.version.removed`).
    *   **`user.account.security.alert.v1`**
        *   **Источник:** Account Service или Auth Service
        *   **Назначение:** Уведомление пользователя о важном событии безопасности.
        *   **Структура Payload (пример):**
            ```json
            {
              "user_id": "uuid_user",
              "alert_type": "suspicious_login_attempt",
              "ip_address": "192.168.1.100",
              "device_info": "Chrome on Windows",
              "timestamp": "ISO8601_timestamp"
            }
            ```
        *   **Логика обработки:** Создать уведомление высокого приоритета (тип `account.security.alert`) и немедленно отправить его пользователю.
*   Также потребляются события из топиков: `user.events`, `payment.events`, `social.events`, `library.events` и др. для инициирования соответствующих уведомлений.
*   **Логика обработки:** Получение запроса/события, определение типа уведомления, проверка предпочтений, выбор шаблона и канала, отправка.
*   (Примеры см. в Спецификации Notification Service, раздел 5.2.1).

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Account Service**: Получение контактных данных, языка пользователя.
*   **Auth Service**: Получение событий (подозрительный вход, 2FA).
*   **Analytics Service**: Получение сегментов для кампаний, отправка статистики.
*   **Social, Library, Payment, etc.**: Получение событий, требующих уведомлений.
*   (Детали см. в Спецификации Notification Service, разделы 3.4, 6.1).

### 6.2. Внешние Системы
*   **Email-провайдеры** (SendGrid, Mailgun, Yandex.Mail): REST API.
*   **Push-провайдеры** (FCM, APNS): SDK/REST API.
*   **SMS-шлюзы** (SMSC.ru, Twilio): REST API.
*   **WebSocket-сервер**: Для In-App уведомлений.
*   (Детали см. в Спецификации Notification Service, раздел 6.2).

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `NOTIFICATION_SERVICE_HTTP_PORT`, `NOTIFICATION_SERVICE_GRPC_PORT`.
*   `POSTGRES_DSN` / `CLICKHOUSE_URL`.
*   `REDIS_ADDR`.
*   `KAFKA_BROKERS` (comma-separated, e.g., `kafka1:9092,kafka2:9092`)
*   `KAFKA_CONSUMER_GROUP_ID` (e.g., `notification-service-group`)
*   `KAFKA_TOPIC_SEND_REQUESTS` (e.g., `notification.send.request`)
*   `KAFKA_TOPIC_STATUS_EVENTS` (e.g., `notification.status.events`)
*   `EMAIL_PROVIDER_PRIMARY_API_KEY` (из Secrets)
*   `EMAIL_PROVIDER_PRIMARY_TYPE` (e.g., `sendgrid`, `mailgun`)
*   `EMAIL_SENDER_DEFAULT_FROM` (e.g., `noreply@gameplatform.ru`)
*   `SMS_PROVIDER_PRIMARY_API_KEY` (из Secrets)
*   `SMS_PROVIDER_PRIMARY_SENDER_ID` (e.g., `GamePlatform`)
*   `FCM_SERVER_KEY` (из Secrets)
*   `APNS_KEY_ID` (из Secrets)
*   `APNS_TEAM_ID` (из Secrets)
*   `APNS_PRIVATE_KEY_PATH` (путь к .p8 файлу, смонтированному из Secrets)
*   `APNS_BUNDLE_ID` (e.g., `com.gameplatform.client`)
*   `WEBSOCKET_GATEWAY_URL` (для In-App уведомлений, e.g., `ws://websocket-gateway/notifications`)
*   `LOG_LEVEL` (e.g., `info`, `debug`)
*   `DEFAULT_RETRY_ATTEMPTS` (e.g., `3`)
*   `DEFAULT_RETRY_DELAY_SECONDS` (e.g., `60`)
*   `ACCOUNT_SERVICE_GRPC_ADDR` (e.g., `account-service:9090`)
*   `ANALYTICS_SERVICE_GRPC_ADDR` (e.g., `analytics-service:9090`)
*   `OTEL_EXPORTER_JAEGER_ENDPOINT` (e.g., `http://jaeger-collector:14268/api/traces`)

### 7.2. Файлы Конфигурации (если применимо)
*   Конфигурация сервиса осуществляется преимущественно через переменные окружения. Если потребуются файлы конфигурации для сложных настроек (например, для детализированных правил выбора провайдеров уведомлений для разных сценариев или сложных структур шаблонов, не хранящихся в БД), их структура будет определена здесь.
*   (Источник: Спецификация Notification Service, раздел 8.2.4).

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   Обработка ошибок от провайдеров доставки.
*   Механизмы повторной отправки (retry) с экспоненциальной задержкой.
*   Использование DLQ для сообщений, которые не удалось обработать.
*   Подробное логирование ошибок.

### 8.2. Распространенные Коды Ошибок
*   `INVALID_RECIPIENT`
*   `TEMPLATE_NOT_FOUND`
*   `PROVIDER_ERROR`
*   `CHANNEL_DISABLED_FOR_USER`
*   `RATE_LIMIT_EXCEEDED`

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   API: JWT для пользовательских запросов, API-ключи для межсервисных.
*   Доступ к провайдерам: Защищенные API ключи.

### 9.2. Авторизация
*   RBAC для доступа к API управления шаблонами, кампаниями.

### 9.3. Защита Данных
*   ФЗ-152: Хранение и обработка ПД.
*   Шифрование: API ключи провайдеров, токены устройств.
*   Защита от спама и злоупотреблений: контроль частоты, верификация отправителей.

### 9.4. Управление Секретами
*   Kubernetes Secrets или HashiCorp Vault.
*   (Детали см. в Спецификации Notification Service, раздел 7.1).

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile.**
*   **Kubernetes манифесты/Helm-чарты.**

### 10.2. Зависимости при Развертывании
*   PostgreSQL/ClickHouse, Redis, Kafka.
*   Зависит от интегрируемых микросервисов для получения событий.

### 10.3. CI/CD
*   Автоматическая сборка, тестирование, развертывание.
*   (Детали см. в Спецификации Notification Service, раздел 8.2).

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   Структурированные JSON логи (Zap/Logback).
*   Централизованный сбор (ELK/Loki).

### 11.2. Мониторинг
*   Prometheus, Grafana.
*   Метрики: количество отправленных/доставленных/ошибочных уведомлений по каналам и типам, задержка доставки, размеры очередей, производительность провайдеров.

### 11.3. Трассировка
*   OpenTelemetry, Jaeger/Zipkin.
*   Отслеживание жизненного цикла уведомления от запроса до доставки.
*   (Детали см. в Спецификации Notification Service, раздел 8.2.3).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:** API P95 < 50 мс, RPS (прием) >= 5000, Задержка доставки P99 < 1 сек (high priority).
*   **Надежность:** Доступность 99.95%, Delivery Rate > 90-98% (зависит от канала).
*   **Масштабируемость:** Горизонтальное масштабирование для пиковых нагрузок.
*   (Детали см. в Спецификации Notification Service, раздел 2.3).

## 13. Приложения (Appendices) (Опционально)
*   Детальные схемы DDL для PostgreSQL/ClickHouse (для хранения шаблонов, статистики, логов доставки, токенов устройств, предпочтений), полные примеры REST API запросов/ответов (включая ответы на ошибки), а также форматы событий Kafka (для запросов на отправку и статусов доставки) будут добавлены по мере финализации дизайна и реализации соответствующих модулей.

---
*Этот шаблон является отправной точкой и может быть адаптирован под конкретные нужды проекта и сервиса.*
