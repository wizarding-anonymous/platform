# Спецификация Микросервиса: Payment Service

**Версия:** 1.0
**Дата последнего обновления:** 2025-05-25

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Payment Service.
*   **Роль в общей архитектуре платформы:** Payment Service является критически важным компонентом, отвечающим за обработку платежей, управление транзакциями и финансовыми операциями. Обеспечивает безопасное проведение платежей, фискализацию, обработку возвратов, управление балансами разработчиков, обработку промокодов и подарочных сертификатов, а также выплаты разработчикам.
*   **Основные бизнес-задачи:** Обработка платежей, управление транзакциями, фискализация (54-ФЗ), обработка возвратов, управление балансами разработчиков, обработка промокодов/подарков, выплаты разработчикам.
*   (Источник: Спецификация Payment Service, разделы 1.1, 1.2, 2.1)

### 1.2. Ключевые Функциональности
*   **Обработка платежей:** Интеграция с российскими платежными системами (СБП, МИР, ЮMoney), инициирование платежа, обработка колбэков, проверка статуса, сохранение платежных методов, поддержка валют.
*   **Управление транзакциями:** Создание, отслеживание статуса, история, детализация, группировка, экспорт.
*   **Фискализация:** Формирование чеков (54-ФЗ), интеграция с ОФД, хранение и доступ к чекам, коррекция чеков.
*   **Обработка возвратов:** Инициирование, проверка возможности, проведение возврата, фискализация возврата.
*   **Управление балансами разработчиков:** Учет доходов, расчет комиссий, отслеживание баланса, история операций, блокировка/корректировка баланса.
*   **Промокоды и подарочные сертификаты:** Создание, управление, применение промокодов; создание, активация подарочных сертификатов.
*   **Выплаты разработчикам:** Планирование, инициирование, поддержка методов выплат, отслеживание статуса, отчеты, налоговая информация.
*   (Источник: Спецификация Payment Service, раздел 2.3)

### 1.3. Основные Технологии
*   **Языки программирования:** Java/Kotlin (основные сервисы), Go (высоконагруженные компоненты).
*   **Фреймворки:** Spring Boot, Micronaut.
*   **Базы данных:** PostgreSQL (основное хранилище), Redis (кэширование, временные данные).
*   **Очереди сообщений:** Kafka.
*   **API:** RESTful API, gRPC (внутреннее).
*   **Инфраструктура:** Docker, Kubernetes, Prometheus, Grafana, ELK Stack, Vault.
*   (Источник: Спецификация Payment Service, раздел 3.4)

### 1.4. Термины и Определения (Glossary)
*   **Транзакция:** Любая финансовая операция на платформе.
*   **Фискализация:** Процесс формирования фискальных чеков согласно 54-ФЗ.
*   **ОФД (Оператор Фискальных Данных):** Посредник для передачи чеков в налоговую.
*   **Платежный шлюз:** Внешняя система для проведения платежей (СБП, МИР, ЮMoney).
*   (Для других терминов см. "Единый глоссарий терминов...")

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Payment Service построен на микросервисной архитектуре с четким разделением ответственности.
*   Ключевые компоненты включают: API Gateway, Transaction Service, Payment Processing Service, Refund Service, Fiscal Service, Balance Service, Promo Service, Payout Service, Notification Service, Reporting Service, Data Storage, Security Service, Monitoring Service.
*   Диаграмма взаимодействия компонентов:
    ```mermaid
    graph TD
        subgraph Payment Service
            API_GW[API Gateway] --> TS(Transaction Service)
            API_GW --> PPS(Payment Processing Service)
            API_GW --> FS(Fiscal Service)
            TS --> PPS; TS --> FS; TS --> BS(Balance Service); TS --> RS(Refund Service); TS --> PS(Promo Service); TS --> PayoutS(Payout Service)
            PPS --> ExtPS[External Payment Systems]
            FS --> OFD[Operator Fiscal Data]
            BS --> DS_DB[(Data Storage)]
            RS --> PPS
            PS --> DS_DB
            PayoutS --> PPS
            NS(Notification Service) -- Receives events from others
            ReportingS(Reporting Service) -- Reads from DS_DB
            DS_DB -- Security & Monitoring --> SecS(Security Service)
            DS_DB -- Security & Monitoring --> MonS(Monitoring Service)
        end
        OtherMicroservices --> API_GW
        Clients --> API_GW
        AdminInterfaces --> API_GW
    ```
*   (Источник: Спецификация Payment Service, разделы 3.1, 3.2, 3.3)

### 2.2. Слои Сервиса
(Предполагаемая структура на основе описанных компонентов)

#### 2.2.1. Presentation Layer (Слой Представления / API Gateway)
*   Ответственность: Прием и маршрутизация всех внешних запросов (от клиентов, других сервисов, админ-панели). Аутентификация и авторизация.
*   Ключевые компоненты/модули: REST API эндпоинты, Webhook-хендлеры.

#### 2.2.2. Application Layer (Прикладной Слой / Компоненты сервиса)
*   Ответственность: Оркестрация бизнес-логики для каждой основной функции (транзакции, обработка платежей, возвраты, фискализация, балансы, промо, выплаты).
*   Ключевые компоненты/модули: `TransactionService`, `PaymentProcessingService`, `RefundService`, `FiscalService`, `BalanceService`, `PromoService`, `PayoutService`. Каждый из этих компонентов реализует соответствующую бизнес-логику и сценарии использования.

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности (Транзакция, Платежный метод, Фискальный чек, Баланс разработчика, Промокод, Выплата), их состояния и правила валидации.
*   Ключевые компоненты/модули: Entities (`Transaction`, `TransactionItem`, `PaymentMethod`, `FiscalReceipt`, `DeveloperBalance`, `PromoCode`, `GiftCard`, `DeveloperPayoutMethod`, `DeveloperPayout`).

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Взаимодействие с PostgreSQL, Redis, Kafka. Интеграция с внешними платежными системами и ОФД. Отправка уведомлений через Notification Service. Формирование отчетов. Обеспечение безопасности и мониторинга.
*   Ключевые компоненты/модули: PostgreSQL repositories, Redis cache, Kafka producers/consumers, клиенты для платежных систем, ОФД, Notification Service, Reporting Service, Security Service, Monitoring Service.

## 3. API Endpoints

### 3.1. REST API
*   **Префикс:** `/api/v1/payments`
*   **Аутентификация:** JWT (через Auth Service).
*   **Авторизация:** На основе ролей.
*   **Основные группы эндпоинтов:**
    *   Транзакции: `POST /transactions`, `GET /transactions`, `GET /transactions/{transaction_id}`, `PATCH /transactions/{transaction_id}` (статус), `GET /transactions/{transaction_id}/receipt`.
    *   Платежные методы: `POST /payment-methods`, `GET /payment-methods`, `DELETE /payment-methods/{payment_method_id}`.
    *   Возвраты: `POST /refunds`, `GET /refunds`, `GET /refunds/{refund_id}`.
    *   Промокоды: `POST /promo-codes`, `GET /promo-codes`, `POST /promo-codes/validate`.
    *   Подарочные сертификаты: `POST /gift-cards`, `POST /gift-cards/activate`.
    *   Балансы разработчиков (админ): `GET /developer-balances`, `GET /developer-balances/{developer_id}`.
    *   Выплаты разработчикам (админ/разработчик): `POST /developer-payout-methods`, `GET /developer-payout-methods`, `POST /developer-payouts`, `GET /developer-payouts`.
    *   Фискальные данные (админ/пользователь): `GET /fiscal-receipts`, `GET /fiscal-receipts/{receipt_id}/download`.
*   (Детали см. в Спецификации Payment Service, раздел 5.2).

### 3.2. Webhook API
*   **Префикс:** `/api/v1/payments/webhooks`
*   `POST /webhooks/{provider}`: Эндпоинт для получения уведомлений от платежных систем (СБП, МИР, ЮMoney).
*   (Детали см. в Спецификации Payment Service, раздел 5.3).

### 3.3. Интеграционные API (внутренние)
*   **Префикс:** `/api/v1/payments/integration`
*   `POST /integration/purchase-complete`: Уведомление о завершении покупки.
*   `POST /integration/refund-complete`: Уведомление о завершении возврата.
*   (Детали см. в Спецификации Payment Service, раздел 5.4).

### 3.4. gRPC API
*   Для внутреннего взаимодействия, например, для проверки статуса транзакции или получения деталей платежного метода другими сервисами, может быть определен следующий gRPC сервис. Детальное определение `.proto` будет предоставлено по мере необходимости.
    ```protobuf
    // Примерный сервис PaymentInternalService
    service PaymentInternalService {
      // Проверить статус транзакции
      rpc GetTransactionStatus(GetTransactionStatusRequest) returns (GetTransactionStatusResponse);
      // Получить доступные платежные методы для пользователя
      rpc GetUserPaymentMethods(GetUserPaymentMethodsRequest) returns (GetUserPaymentMethodsResponse);
      // Инициировать удержание средств (холдирование)
      rpc HoldPayment(HoldPaymentRequest) returns (HoldPaymentResponse);
      // Подтвердить списание удержанных средств
      rpc ConfirmHold(ConfirmHoldRequest) returns (ConfirmHoldResponse);
      // Отменить удержание средств
      rpc CancelHold(CancelHoldRequest) returns (CancelHoldResponse);
    }

    message GetTransactionStatusRequest {
      string transaction_id = 1;
    }

    message GetTransactionStatusResponse {
      string transaction_id = 1;
      string status = 2; // e.g., "created", "processing", "completed", "failed"
      string payment_method_type = 3; // e.g., "card", "sbp"
      google.protobuf.Timestamp completed_at = 4;
    }

    message GetUserPaymentMethodsRequest {
      string user_id = 1;
    }

    message PaymentMethod {
      string payment_method_id = 1;
      string type = 2; // "card", "sbp", "yoomoney"
      string masked_identifier = 3; // "•••• 1234" for card, or phone for SBP
      bool is_default = 4;
    }

    message GetUserPaymentMethodsResponse {
      repeated PaymentMethod payment_methods = 1;
    }

    message HoldPaymentRequest {
      string user_id = 1;
      string order_id = 2; // ID заказа из другого сервиса
      string payment_method_id = 3; // опционально, если пользователь выбирает
      int64 amount_kop = 4; // сумма в копейках
      string currency_code = 5; // "RUB"
      string description = 6;
    }

    message HoldPaymentResponse {
      string transaction_id = 1;
      string status = 2; // "hold_pending", "hold_success", "hold_failed"
      string payment_gateway_reference_id = 3; // ID операции в платежном шлюзе
    }

    message ConfirmHoldRequest {
      string transaction_id = 1; // ID транзакции холдирования
      int64 final_amount_kop = 2; // может отличаться от суммы холда
    }

    message ConfirmHoldResponse {
      string transaction_id = 1;
      string status = 2; // "completed", "failed"
    }

    message CancelHoldRequest {
      string transaction_id = 1; // ID транзакции холдирования
      string reason = 2;
    }

    message CancelHoldResponse {
      string transaction_id = 1;
      string status = 2; // "cancelled", "failed"
    }
    // Другие необходимые сообщения должны быть определены дополнительно
    ```

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности
*   **Transaction**: Транзакция (покупка, возврат, выплата).
*   **TransactionItem**: Элемент транзакции (игра, подписка).
*   **PaymentMethod**: Сохраненный платежный метод пользователя.
*   **FiscalReceipt**: Фискальный чек.
*   **DeveloperBalance**: Баланс разработчика.
*   **DeveloperBalanceHistory**: История операций по балансу.
*   **PromoCode**: Промокод на скидку/товар.
*   **GiftCard**: Подарочный сертификат.
*   **DeveloperPayoutMethod**: Метод выплаты для разработчика.
*   **DeveloperPayout**: Операция выплаты разработчику.
*   (SQL DDL см. в Спецификации Payment Service, раздел 5.1).

### 4.2. Схема Базы Данных
*   **PostgreSQL**: Хранит все основные данные о транзакциях, платежных методах, чеках, балансах, промокодах, выплатах.
    ```sql
    -- Пример таблицы transactions (сокращенно)
    CREATE TABLE transactions (transaction_id UUID PRIMARY KEY, user_id UUID NOT NULL, type VARCHAR(50) NOT NULL, status VARCHAR(50) NOT NULL, amount DECIMAL(10, 2) NOT NULL ...);
    -- Пример таблицы payment_methods (сокращенно)
    CREATE TABLE payment_methods (payment_method_id UUID PRIMARY KEY, user_id UUID NOT NULL, type VARCHAR(50) NOT NULL, token VARCHAR(255) NOT NULL ...);
    ```
*   **Redis**: Используется для кэширования сессий платежей, временных токенов, лимитов скорости.
*   (Полные DDL см. в Спецификации Payment Service, раздел 5.1).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Kafka.
*   **Формат событий:** CloudEvents JSON (предположительно).
*   **Основные публикуемые события:**
    *   `payment.transaction.created`: Создана новая транзакция.
    *   `payment.transaction.completed`: Транзакция успешно завершена (оплачена). -> Library Service, Notification Service, Analytics Service.
    *   `payment.transaction.failed`: Ошибка транзакции. -> Notification Service.
    *   `payment.refund.processed`: Возврат обработан. -> Library Service, Notification Service, Analytics Service.
    *   `payment.payout.initiated`: Инициирована выплата разработчику.
        *   `Структура Payload (пример):`
            ```json
            {
              "event_id": "uuid_event",
              "event_type": "payment.payout.initiated.v1",
              "timestamp": "ISO8601_timestamp",
              "source_service": "payment-service",
              "payout_id": "uuid_payout_operation",
              "developer_id": "uuid_developer",
              "amount": 50000.00,
              "currency": "RUB",
              "payout_method_type": "bank_transfer"
            }
            ```
    *   `payment.payout.status.updated`: Статус выплаты разработчику изменен. -> Developer Service, Notification Service.
        *   `Структура Payload (пример):`
            ```json
            {
              "event_id": "uuid_event",
              "event_type": "payment.payout.status.updated.v1",
              "timestamp": "ISO8601_timestamp",
              "source_service": "payment-service",
              "payout_id": "uuid_payout_operation",
              "developer_id": "uuid_developer",
              "new_status": "completed" | "failed" | "processing",
              "processed_at": "ISO8601_timestamp",
              "external_transaction_id": "id_in_payment_system",
              "failure_reason": "Ошибка обработки банком"
            }
            ```
    *   `payment.fiscal.receipt.created`: Фискальный чек создан. -> Notification Service (для отправки пользователю).
        *   `Структура Payload (пример):`
            ```json
            {
              "event_id": "uuid_event",
              "event_type": "payment.fiscal.receipt.created.v1",
              "timestamp": "ISO8601_timestamp",
              "source_service": "payment-service",
              "transaction_id": "uuid_transaction",
              "receipt_id": "uuid_fiscal_receipt",
              "fiscal_document_number": "1234567890",
              "receipt_url_or_data": "url_to_ofd_receipt_or_base64_data"
            }
            ```

### 5.2. Потребляемые События (Consumed Events)
*   `user.account.created`:
    *   **Источник:** Account Service
    *   **Назначение:** Создание профиля платежных данных для нового пользователя, если это требуется (например, инициализация внутреннего баланса, если есть).
    *   **Структура Payload (пример):**
        ```json
        {
          "user_id": "uuid_user",
          "email": "user@example.com",
          "username": "user123",
          "registration_date": "ISO8601_timestamp"
        }
        ```
    *   **Логика обработки:** Payment Service может создать внутреннюю запись для пользователя, подготовить структуру для хранения его будущих платежных методов или транзакций.
*   `developer.payout.requested`:
    *   **Источник:** Developer Service
    *   **Назначение:** Инициирование процесса выплаты разработчику.
    *   **Структура Payload (пример):**
        ```json
        {
          "payout_request_id": "uuid_dev_service_payout_request",
          "developer_id": "uuid_developer",
          "amount": 75000.00,
          "currency": "RUB",
          "requested_payout_method_id": "uuid_developer_payout_method",
          "requested_at": "ISO8601_timestamp"
        }
        ```
    *   **Логика обработки:** Валидировать запрос, проверить баланс разработчика, создать транзакцию типа "payout" со статусом "processing", инициировать взаимодействие с платежной системой для осуществления выплаты.
*   `admin.refund.request.manual`:
    *   **Источник:** Admin Service
    *   **Назначение:** Ручное инициирование возврата средств пользователю.
    *   **Структура Payload (пример):**
        ```json
        {
          "original_transaction_id": "uuid_purchase_transaction",
          "user_id": "uuid_user",
          "amount_to_refund": 1000.00,
          "reason": "Решение администратора",
          "admin_id": "uuid_admin_user",
          "requested_at": "ISO8601_timestamp"
        }
        ```
    *   **Логика обработки:** Проверить исходную транзакцию, создать транзакцию возврата, инициировать возврат через платежную систему, создать фискальный чек возврата.
*   `order.created.v1`:
    *   **Источник:** Order Service (или Catalog/Library Service)
    *   **Назначение:** Уведомление о создании нового заказа, для которого необходимо инициировать платеж.
    *   **Структура Payload (пример):**
        ```json
        {
          "order_id": "uuid_order",
          "user_id": "uuid_user",
          "total_amount": 2499.00,
          "currency": "RUB",
          "items": [
            { "item_id": "uuid_game1", "item_type": "game", "price": 1999.00, "quantity": 1 },
            { "item_id": "uuid_dlc1", "item_type": "dlc", "price": 500.00, "quantity": 1 }
          ],
          "created_at": "ISO8601_timestamp"
        }
        ```
    *   **Логика обработки:** Создать новую транзакцию типа "purchase" со статусом "created".
*   `user.subscription.renewal_due.v1`:
    *   **Источник:** Account Service или Subscription Service
    *   **Назначение:** Уведомление о необходимости произвести периодический платеж за подписку.
    *   **Структура Payload (пример):**
        ```json
        {
          "subscription_id": "uuid_subscription",
          "user_id": "uuid_user",
          "plan_id": "uuid_subscription_plan",
          "renewal_amount": 499.00,
          "currency": "RUB",
          "next_billing_date": "ISO8601_date"
        }
        ```
    *   **Логика обработки:** Инициировать рекуррентный платеж, используя сохраненный платежный метод пользователя.
*   *(Примечание: Событие `catalog.price.changed` также может потребляться для аннулирования/обновления "корзин" или отложенных платежей, если такая функциональность предусмотрена. Событие `user.account.deleted` от Account Service также должно обрабатываться для анонимизации или удаления платежных данных пользователя).*

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Account Service**: Получение данных о пользователях.
*   **Catalog Service**: Получение данных о ценах игр и скидках.
*   **Library Service**: Уведомление о покупках/возвратах для управления библиотекой.
*   **Developer Service**: Информация о финансовых операциях и балансах разработчиков.
*   **Admin Service**: Административный доступ к финансовым операциям.
*   **Analytics Service**: Предоставление данных о транзакциях.
*   **Notification Service**: Инициирование уведомлений о финансовых операциях.
*   (Детали см. в Спецификации Payment Service, разделы 1.3 и 6).

### 6.2. Внешние Системы
*   **Платежные системы (СБП, МИР, ЮMoney)**: Для обработки платежей.
*   **Операторы фискальных данных (ОФД)**: Для передачи фискальных чеков.
*   (Детали см. в Спецификации Payment Service, раздел 6.7).

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `PAYMENT_SERVICE_PORT`: Порт сервиса.
*   `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DBNAME`, `POSTGRES_SSLMODE`
*   `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB_PAYMENT`
*   `KAFKA_BROKERS` (comma-separated)
*   `KAFKA_TOPIC_PAYMENT_EVENTS` (e.g., `payment.transactions`)
*   `SBP_GATEWAY_URL`, `SBP_GATEWAY_API_KEY_SECRET` (имя секрета в Kubernetes)
*   `MIR_GATEWAY_URL`, `MIR_GATEWAY_MERCHANT_ID_SECRET`, `MIR_GATEWAY_API_KEY_SECRET`
*   `YOOMONEY_SHOP_ID_SECRET`, `YOOMONEY_SECRET_KEY_SECRET`
*   `OFD_API_URL`, `OFD_API_KEY_SECRET`, `OFD_INN`, `OFD_KKT_REG_NUMBER`
*   `PLATFORM_COMMISSION_PERCENT` (e.g., `30.0`)
*   `MIN_PAYOUT_AMOUNT_RUB` (e.g., `1000.00`)
*   `DEFAULT_CURRENCY` (e.g., `RUB`)
*   `LOG_LEVEL` (e.g., `info`, `debug`)
*   `AUTH_SERVICE_GRPC_ADDR` (e.g., `auth-service:9090`)
*   `ACCOUNT_SERVICE_GRPC_ADDR` (e.g., `account-service:9090`)
*   `CATALOG_SERVICE_GRPC_ADDR` (e.g., `catalog-service:9090`)
*   `NOTIFICATION_SERVICE_KAFKA_TOPIC` (e.g., `notification.send.request`)
*   `OTEL_EXPORTER_JAEGER_ENDPOINT` (e.g., `http://jaeger-collector:14268/api/traces`)
*   `HTTP_SERVER_PORT` (e.g., `8080`)
*   `GRPC_SERVER_PORT` (e.g., `9090`)
*   `TRANSACTION_LOCK_TIMEOUT_SECONDS` (e.g., `60`)
*   `PAYMENT_SESSION_TTL_SECONDS` (e.g., `900` for 15 minutes)

### 7.2. Файлы Конфигурации (если применимо)
*   Конфигурация сервиса осуществляется преимущественно через переменные окружения. Если потребуются файлы конфигурации для сложных настроек (например, для правил маршрутизации между разными платежными шлюзами или детализированных параметров фискализации для разных типов товаров/услуг), их структура будет определена здесь.

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   Корректная обработка и логирование ошибок от платежных систем и ОФД.
*   Механизмы повторных попыток для временных сбоев.
*   Четкие сообщения об ошибках для пользователей и администраторов.

### 8.2. Распространенные Коды Ошибок
*   `INSUFFICIENT_FUNDS`
*   `PAYMENT_GATEWAY_ERROR`
*   `FISCALIZATION_ERROR`
*   `TRANSACTION_NOT_FOUND`
*   `REFUND_POLICY_VIOLATION`
*   `PROMO_CODE_INVALID_OR_EXPIRED`

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   JWT (через Auth Service) для всех API.
*   mTLS для межсервисных вызовов к критичным частям.

### 9.2. Авторизация
*   RBAC для доступа к операциям (например, только админ может делать выплаты или корректировки).

### 9.3. Защита Данных
*   **PCI DSS Compliance:** Соответствие требованиям (если планируется прямая обработка карточных данных, иначе через токенизацию на стороне шлюза).
*   Шифрование чувствительных данных (токены платежных методов, реквизиты для выплат) при хранении и передаче.
*   Защита от мошенничества.

### 9.4. Управление Секретами
*   API ключи платежных систем, ключи шифрования в Kubernetes Secrets или Vault.
*   **Аудит**: Детальное логирование всех финансовых операций.
*   (Детали см. в Спецификации Payment Service, раздел 7.1).

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile.**
*   **Kubernetes манифесты/Helm-чарты.**

### 10.2. Зависимости при Развертывании
*   PostgreSQL, Redis, Kafka.
*   Account Service, Catalog Service, Auth Service, Notification Service.

### 10.3. CI/CD
*   Автоматическая сборка, тестирование (с особым вниманием к финансовым расчетам), развертывание.
*   (Детали см. в Спецификации Payment Service, раздел 8).

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   Структурированные JSON логи.
*   Централизованный сбор (ELK).
*   Логирование всех этапов транзакций, выплат, фискализации.

### 11.2. Мониторинг
*   Prometheus, Grafana.
*   Метрики: количество и суммы транзакций (успешных/неуспешных), время обработки платежей, ошибки интеграции с платежными системами/ОФД, размеры очередей.

### 11.3. Трассировка
*   Интеграция с системой распределенной трассировки (например, Jaeger/OpenTelemetry) будет реализована для отслеживания полного жизненного цикла финансовых операций, включая вызовы к внешним платежным системам и ОФД. Контекст трассировки будет логироваться и передаваться.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:** API P95 < 500мс, обработка >= 100 транзакций/сек.
*   **Надежность:** Доступность >= 99.9%.
*   **Безопасность:** Шифрование, PCI DSS (если применимо), ФЗ-54, ФЗ-152, ФЗ-115.
*   **Масштабируемость:** Горизонтальная.
*   (Детали см. в Спецификации Payment Service, раздел 2.4).

## 13. Приложения (Appendices) (Опционально)
*   Детальные схемы DDL для PostgreSQL, полные примеры REST API запросов/ответов (включая обработку вебхуков и ошибок), форматы событий Kafka, а также примеры фискальных чеков (в соответствии с 54-ФЗ) будут добавлены по мере финализации дизайна и реализации интеграций с платежными системами.

---
*Этот шаблон является отправной точкой и может быть адаптирован под конкретные нужды проекта и сервиса.*
