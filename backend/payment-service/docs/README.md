# Спецификация Микросервиса: Payment Service (Сервис Платежей)

**Версия:** 1.1 (адаптировано из предыдущей версии)
**Дата последнего обновления:** 2024-03-15

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение документа:** Данный документ представляет собой полную спецификацию микросервиса Payment Service платформы "Российский Аналог Steam".
*   **Роль в общей архитектуре платформы:** Payment Service является критически важным компонентом, отвечающим за обработку всех видов платежей, управление транзакциями и связанными финансовыми операциями. Он обеспечивает безопасное проведение платежей пользователей за продукты и услуги, фискализацию этих операций в соответствии с законодательством РФ (54-ФЗ), обработку возвратов, управление балансами разработчиков, применение промокодов и подарочных сертификатов, а также организацию выплат разработчикам.
*   **Основные бизнес-задачи:**
    *   Обеспечение возможности для пользователей оплачивать покупки различными способами.
    *   Гарантия безопасности и надежности проведения финансовых операций.
    *   Соблюдение требований фискального законодательства.
    *   Управление процессом возврата средств пользователям.
    *   Расчет и управление доходами разработчиков и платформы.
    *   Поддержка маркетинговых активностей через промокоды и подарочные сертификаты.
    *   Осуществление регулярных выплат разработчикам.
*   Разработка сервиса должна вестись в соответствии с `CODING_STANDARDS.md`.

### 1.2. Ключевые Функциональности
*   **Обработка входящих платежей:** Интеграция с российскими платежными системами и шлюзами (например, СБП, МИР, ЮMoney, банковские карты через PSP). Инициирование платежа, обработка колбэков от платежных систем, проверка статуса платежа.
*   **Управление транзакциями:** Создание и отслеживание статуса транзакций (покупка, пополнение баланса, возврат, выплата). Хранение истории транзакций.
*   **Фискализация:** Формирование фискальных чеков в соответствии с 54-ФЗ для всех приходных и расходных операций. Интеграция с Оператором Фискальных Данных (ОФД). Хранение и предоставление доступа к фискальным чекам. Поддержка коррекции чеков.
*   **Обработка возвратов:** API для инициирования возвратов. Проверка возможности возврата (например, политики возвратов, время с момента покупки). Проведение возврата через платежную систему. Фискализация чека возврата.
*   **Управление балансами разработчиков:** Учет доходов от продаж продуктов разработчиков. Расчет комиссий платформы. Отслеживание текущего баланса разработчика. История операций по балансу. Возможность блокировки или корректировки баланса администраторами.
*   **Промокоды и подарочные сертификаты:** Создание и управление промокодами (процентные, на фиксированную сумму, на конкретный товар). Валидация и применение промокодов к заказам. Создание, активация и погашение подарочных сертификатов.
*   **Выплаты разработчикам:** Управление методами выплат для разработчиков. Планирование и инициирование выплат. Отслеживание статуса выплат. Формирование отчетов по выплатам.
*   **Управление сохраненными платежными методами:** Безопасное сохранение токенизированных платежных методов пользователей (если поддерживается PSP) для упрощения повторных покупок.

### 1.3. Основные Технологии
*   **Языки программирования:** Java/Kotlin (для основных компонентов и бизнес-логики), Go (для высоконагруженных прокси или специфических утилит, если потребуется).
*   **Фреймворки:** Spring Boot (для Java/Kotlin).
*   **Базы данных:** PostgreSQL (для основного хранилища транзакционных данных, балансов, метаданных).
*   **Кэширование/Сессии:** Redis (для кэширования временных данных платежных сессий, idempotency ключей, лимитов).
*   **Очереди сообщений:** Apache Kafka (для асинхронной обработки событий, связанных с платежами, фискализацией, уведомлениями).
*   **API:** RESTful API (для взаимодействия с фронтендом и некоторыми внутренними сервисами), gRPC (для критичных внутренних межсервисных взаимодействий).
*   **Инфраструктура:** Docker, Kubernetes.
*   **Мониторинг/Трассировка:** OpenTelemetry, Prometheus, Grafana, Jaeger/Tempo.
*   **Безопасность:** HashiCorp Vault или Kubernetes Secrets для управления секретами.
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `PACKAGE_STANDARDIZATION.md`.
*   *Примечание: Приведенные ниже примеры API и конфигураций основаны на предположении использования Java/Kotlin (Spring Boot) для основного API и бизнес-логики, PostgreSQL для хранения данных, Redis для кэширования и Kafka для обмена событиями.*

### 1.4. Термины и Определения (Glossary)
*   **Транзакция (Transaction):** Любая финансовая операция, регистрируемая в системе (покупка, возврат, выплата, пополнение баланса).
*   **Платежный Шлюз (Payment Gateway / PSP - Payment Service Provider):** Внешняя система, обеспечивающая проведение онлайн-платежей (например, Сбербанк Эквайринг, ЮKassa, Тинькофф Оплата).
*   **Фискализация (Fiscalization):** Процесс формирования и отправки фискальных чеков в налоговые органы через ОФД в соответствии с 54-ФЗ.
*   **ОФД (Оператор Фискальных Данных):** Организация, уполномоченная на обработку, хранение и передачу фискальных данных от ККТ в ФНС.
*   **Холдирование (Payment Hold/Authorization):** Временная блокировка средств на карте покупателя до подтверждения или отмены операции.
*   **Клиринг (Clearing):** Процесс взаиморасчетов между участниками платежной системы.
*   **Выплата (Payout):** Перечисление денежных средств разработчикам за проданные продукты.
*   Для других общих терминов см. `project_glossary.md`.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Payment Service построен на микросервисной архитектуре, где каждый компонент отвечает за определенную часть функциональности (например, обработка транзакций, взаимодействие с платежными шлюзами, фискализация, управление балансами). Компоненты взаимодействуют друг с другом через синхронные (gRPC/REST) и асинхронные (Kafka) интерфейсы.
*   Сервис обеспечивает изоляцию взаимодействия с внешними платежными системами и ОФД, предоставляя унифицированный интерфейс для других сервисов платформы.
*   Диаграмма взаимодействия компонентов (из предыдущей версии документа, актуальна концептуально):
    ```mermaid
    graph TD
        subgraph Payment Service
            API_GW[API Gateway/Controller Layer] --> TS(Transaction Service Logic)
            API_GW --> PPS(Payment Processing Logic)
            API_GW --> FS(Fiscalization Logic)
            TS --> PPS; TS --> FS; TS --> BS(Balance Management Logic); TS --> RS(Refund Logic); TS --> PromoLogic(Promo & GiftCard Logic); TS --> PayoutLogic(Payout Logic)
            PPS --> ExtPS[External Payment Systems Adapters]
            FS --> OFDAdapter[OFD Adapters]
            BS --> DB_Storage[(Data Storage - PostgreSQL, Redis)]
            RS --> PPS
            PromoLogic --> DB_Storage
            PayoutLogic --> PPS
            InternalKafka[Kafka Producers/Consumers] -- Events & Commands
        end
        OtherMicroservices -- API Calls / Kafka Events --> API_GW
        OtherMicroservices -- API Calls / Kafka Events --> InternalKafka
        Clients -- User Actions --> API_GW
        AdminInterfaces -- Admin Actions --> API_GW

        ExtPS --> ExternalPaymentGateways[Внешние Платежные Шлюзы]
        OFDAdapter --> ExternalOFD[Внешние ОФД]

        classDef component fill:#d4edda,stroke:#28a745,color:#000
        classDef logic fill:#e6f0ff,stroke:#007bff,color:#000
        classDef adapter fill:#fff3cd,stroke:#ffc107,color:#000
        classDef datastore fill:#f8d7da,stroke:#dc3545,color:#000
        classDef external fill:#e2e3e5,stroke:#6c757d,stroke-width:2px;

        class API_GW component;
        class TS,PPS,FS,BS,RS,PromoLogic,PayoutLogic logic;
        class ExtPS,OFDAdapter adapter;
        class DB_Storage datastore;
        class InternalKafka component_minor;
        class OtherMicroservices,Clients,AdminInterfaces,ExternalPaymentGateways,ExternalOFD external;
    ```

### 2.2. Слои Сервиса (и ключевые компоненты)

#### 2.2.1. Presentation Layer (Слой Представления)
*   **Ответственность:** Прием HTTP REST запросов от клиентских приложений и административных интерфейсов (через API Gateway). Прием Webhook-уведомлений от внешних платежных систем. Прием gRPC запросов от других внутренних микросервисов. Валидация входящих данных (DTO).
*   **Ключевые компоненты/модули:** REST контроллеры (Spring MVC/WebFlux), Webhook-хендлеры, gRPC серверные реализации.

#### 2.2.2. Application Layer (Прикладной Слой / Слой Сценариев Использования)
*   **Ответственность:** Оркестрация выполнения бизнес-операций. Координация взаимодействия между различными доменными сущностями и инфраструктурными компонентами. Реализация сценариев использования (use cases).
*   **Ключевые компоненты/модули:**
    *   `TransactionApplicationService`: Управление жизненным циклом транзакций.
    *   `PaymentProcessingApplicationService`: Взаимодействие с платежными шлюзами.
    *   `FiscalizationApplicationService`: Управление фискальными чеками.
    *   `RefundApplicationService`: Обработка возвратов.
    *   `BalanceApplicationService`: Управление балансами разработчиков.
    *   `PromoApplicationService`: Логика промокодов и подарочных карт.
    *   `PayoutApplicationService`: Управление процессом выплат.

#### 2.2.3. Domain Layer (Доменный Слой)
*   **Ответственность:** Содержит бизнес-сущности, агрегаты, доменные события и бизнес-правила, связанные с платежами, транзакциями, балансами и т.д.
*   **Ключевые компоненты/модули:**
    *   **Entities (Сущности):** `Transaction`, `TransactionItem`, `PaymentMethod` (сохраненный токен), `FiscalReceipt`, `DeveloperBalance`, `DeveloperPayout`, `PromoCode`, `GiftCard`.
    *   **Aggregates:** Например, `TransactionAggregate` (включая саму транзакцию и ее элементы).
    *   **Value Objects:** `Money` (сумма + валюта), `PaymentStatus`, `TransactionType`.
    *   **Domain Services:** Для сложной логики, не относящейся к одной сущности (например, расчет комиссии, проверка условий для возврата).
    *   **Domain Events:** `PaymentCompletedEvent`, `RefundProcessedEvent`, `PayoutInitiatedEvent`.
    *   **Repository Interfaces:** Контракты для сохранения и извлечения сущностей.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   **Ответственность:** Реализация интерфейсов репозиториев для работы с PostgreSQL и Redis. Взаимодействие с внешними платежными шлюзами и ОФД через их API. Публикация доменных событий в Kafka.
*   **Ключевые компоненты/модули:** Реализации репозиториев (например, Spring Data JPA), клиенты для платежных шлюзов (HTTP/SOAP клиенты), клиенты для ОФД, Kafka продюсеры, RedisTemplate или аналоги.

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL:** `/api/v1/payments` (маршрутизируется через API Gateway).
*   **Аутентификация:** JWT Bearer Token (проверяется API Gateway, `X-User-Id` передается в заголовках).
*   **Авторизация:** На основе ролей (`user`, `developer_finance`, `payment_admin`).
*   **Стандартный формат ответа об ошибке:**
    ```json
    {
      "errors": [
        {
          "status": "4XX/5XX",
          "code": "ERROR_CODE_UPPER_SNAKE_CASE",
          "title": "Краткое описание ошибки на русском",
          "detail": "Полное описание ошибки с контекстом."
        }
      ]
    }
    ```

#### 3.1.1. Транзакции (Transactions)
*   **`POST /transactions/initiate`**
    *   Описание: Инициирование новой транзакции на покупку (например, игры, пополнения баланса). Сервис создает транзакцию в статусе "pending" и возвращает URL для редиректа на страницу оплаты или данные для SDK платежной системы.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "transactionInitiation",
            "attributes": {
              "user_id": "user-uuid-123",
              "order_id": "order-uuid-abc", // ID заказа из Order Service
              "amount_minor_units": 199900, // 1999.00 RUB
              "currency_code": "RUB",
              "payment_method_type_hint": "sbp", // "card", "yoomoney", etc. (опционально)
              "description": "Покупка игры 'Супер Гонки'",
              "items": [ // Для фискализации
                { "name": "Игра 'Супер Гонки'", "quantity": 1, "price_minor_units": 199900, "vat_code": "VAT_20" }
              ],
              "success_url": "https://example.com/payment/success",
              "fail_url": "https://example.com/payment/fail"
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created):
        ```json
        {
          "data": {
            "type": "transaction",
            "id": "txn-uuid-xyz",
            "attributes": {
              "status": "pending_payment_provider_redirect",
              "redirect_url": "https://psp.example.com/pay/session-token", // или данные для SDK
              "payment_provider_data": { /* ... специфичные данные для SDK PSP ... */ }
            }
          }
        }
        ```
    *   Требуемые права доступа: `user_self_only`.
*   **`GET /transactions/{transaction_id}`**
    *   Описание: Получение статуса и деталей транзакции.
    *   Пример ответа (Успех 200 OK):
        ```json
        {
          "data": {
            "type": "transaction",
            "id": "txn-uuid-xyz",
            "attributes": {
              "user_id": "user-uuid-123",
              "order_id": "order-uuid-abc",
              "status": "completed", // "pending", "failed", "refunded"
              "amount_minor_units": 199900,
              "currency_code": "RUB",
              "payment_method_used": "sbp",
              "created_at": "2024-03-15T10:00:00Z",
              "updated_at": "2024-03-15T10:05:00Z",
              "fiscal_receipt_id": "receipt-uuid-123" // Опционально
            }
          }
        }
        ```
    *   Требуемые права доступа: `user_self_only` (если это его транзакция) или `payment_admin`.

#### 3.1.2. Платежные Методы (Payment Methods)
*   **`POST /users/{user_id}/payment-methods`**
    *   Описание: Добавление (сохранение токенизированного) платежного метода для пользователя. Фактические данные карты не передаются, только токен от PSP.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "paymentMethodRegistration",
            "attributes": {
              "payment_provider": "yookassa", // или другой PSP
              "provider_payment_method_token": "psp-token-for-card-xyz",
              "card_details_hint": { // Не для хранения, а для отображения пользователю
                "card_type": "Visa",
                "last_four_digits": "1234",
                "expiry_month": 12,
                "expiry_year": 2028
              },
              "is_default": true
            }
          }
        }
        ```
    *   Пример ответа (Успех 201 Created): (Возвращает созданный `paymentMethod` ресурс)
    *   Требуемые права доступа: `user_self_only`.

#### 3.1.3. Возвраты (Refunds)
*   **`POST /refunds`**
    *   Описание: Инициирование возврата по существующей транзакции.
    *   Тело запроса:
        ```json
        {
          "data": {
            "type": "refundRequest",
            "attributes": {
              "original_transaction_id": "txn-uuid-xyz",
              "amount_minor_units_to_refund": 199900, // Может быть частичный возврат
              "reason": "Пользователь запросил возврат по правилам платформы.",
              "initiator_id": "admin-user-uuid-789" // ID админа или системы
            }
          }
        }
        ```
    *   Пример ответа (Успех 202 Accepted): (Возвращает созданную транзакцию возврата)
    *   Требуемые права доступа: `payment_admin` или `customer_support_lead`.

### 3.2. Webhook API (от Платежных Провайдеров)
*   **`POST /webhooks/{provider_name}`** (например, `/webhooks/yookassa`, `/webhooks/sbp_bank101`)
    *   Описание: Эндпоинт для получения асинхронных уведомлений (колбэков) от внешних платежных систем о статусе платежа или других событиях.
    *   Тело запроса: Зависит от конкретного провайдера. Должно содержать идентификатор транзакции и новый статус.
        *   Пример (концептуальный, от YooMoney):
            ```json
            {
              "type": "notification",
              "event": "payment.succeeded", // "payment.waiting_for_capture", "refund.succeeded"
              "object": {
                "id": "psp-payment-uuid-abc", // ID платежа в системе YooMoney
                "status": "succeeded",
                "amount": { "value": "1999.00", "currency": "RUB" },
                "metadata": { "internal_transaction_id": "txn-uuid-xyz" },
                // ... другие поля от YooMoney
              }
            }
            ```
    *   Обработка:
        1.  Валидация источника запроса (например, проверка IP-адреса, подписи запроса, если предоставляется провайдером).
        2.  Проверка идемпотентности (например, по `psp-payment-uuid-abc` и `event`).
        3.  Извлечение `internal_transaction_id` из метаданных или поиск транзакции по `psp-payment-uuid-abc`.
        4.  Обновление статуса транзакции в БД Payment Service.
        5.  Если транзакция завершена успешно (`payment.succeeded`):
            *   Инициировать фискализацию чека.
            *   Опубликовать событие `payment.transaction.completed.v1` в Kafka.
        6.  Ответ провайдеру (обычно HTTP 200 OK, если обработка принята).
    *   Требуемые права доступа: Публичный, но с проверкой источника/подписи.

### 3.3. Интеграционные API (внутренние)
*   В основном, взаимодействие с другими сервисами происходит через Kafka события (см. раздел 5).
*   Если требуется синхронное взаимодействие для специфичных нужд (например, немедленная проверка возможности возврата перед отображением кнопки пользователю), могут быть использованы gRPC эндпоинты (см. 3.4).
*   `POST /integration/purchase-complete` и `POST /integration/refund-complete` из предыдущей версии документа, скорее всего, будут заменены на публикацию событий в Kafka (`payment.transaction.completed.v1`, `payment.refund.processed.v1`).

### 3.4. gRPC API (для межсервисного взаимодействия)
*   Пакет: `payment.v1`.
*   Определение Protobuf: `payment/v1/payment_internal_service.proto`.

#### 3.4.1. Сервис: `PaymentInternalService`
*   **`rpc GetTransactionStatus(GetTransactionStatusRequest) returns (GetTransactionStatusResponse)`**
    *   Описание: Получение текущего статуса и деталей транзакции по ее ID.
    *   `message GetTransactionStatusRequest { string transaction_id = 1; }`
    *   `message GetTransactionStatusResponse { string transaction_id = 1; string status = 2; int64 amount_minor_units = 3; string currency_code = 4; string payment_method_type_used = 5; google.protobuf.Timestamp created_at = 6; google.protobuf.Timestamp updated_at = 7; string error_message = 8; /* если статус failed */ }`
*   **`rpc HoldPayment(HoldPaymentRequest) returns (HoldPaymentResponse)`**
    *   Описание: Инициирование операции удержания (холдирования) средств для последующего списания.
    *   `message HoldPaymentRequest { string user_id = 1; string order_id = 2; int64 amount_minor_units = 3; string currency_code = 4; string payment_method_id_hint = 5; /* опционально, ID сохраненного метода */ string description = 6; google.protobuf.Duration hold_duration = 7; /* опционально, на какой срок холдировать */ }`
    *   `message HoldPaymentResponse { string transaction_id = 1; /* ID созданной транзакции холдирования */ string status = 2; /* pending_hold, hold_successful, hold_failed */ string payment_provider_reference_id = 3; string redirect_url_if_needed = 4; /* если требуется 3DS или другая аутентификация */ }`
*   **`rpc ConfirmPaymentHold(ConfirmPaymentHoldRequest) returns (ConfirmPaymentHoldResponse)`**
    *   Описание: Подтверждение списания ранее захолдированных средств.
    *   `message ConfirmPaymentHoldRequest { string hold_transaction_id = 1; int64 final_amount_minor_units = 2; /* может отличаться от суммы холда, если часть товаров отменили */ repeated FiscalItem items_for_receipt = 3; }`
    *   `message FiscalItem { string name = 1; int64 price_minor_units = 2; int32 quantity = 3; string vat_code = 4; }`
    *   `message ConfirmPaymentHoldResponse { string capture_transaction_id = 1; string status = 2; /* completed, failed */ }`
*   **`rpc CancelPaymentHold(CancelPaymentHoldRequest) returns (CancelPaymentHoldResponse)`**
    *   Описание: Отмена ранее созданного холда.
    *   `message CancelPaymentHoldRequest { string hold_transaction_id = 1; string reason = 2; }`
    *   `message CancelPaymentHoldResponse { string status = 1; /* cancelled, failed_to_cancel */ }`

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности

*   **`Transaction` (Транзакция)**
    *   `id` (UUID): Уникальный идентификатор транзакции. Обязательность: Required.
    *   `user_id` (UUID): ID пользователя, инициировавшего транзакцию. Обязательность: Required (nullable для системных транзакций).
    *   `order_id` (UUID): ID заказа из Order Service (если применимо). Обязательность: Optional.
    *   `transaction_type` (ENUM: `purchase`, `refund`, `payout`, `deposit`, `hold`, `capture`): Тип транзакции. Обязательность: Required.
    *   `status` (ENUM: `pending`, `processing`, `awaiting_psp_callback`, `awaiting_capture`, `completed`, `failed`, `cancelled`, `refunded`): Статус транзакции. Обязательность: Required.
    *   `amount_minor_units` (BIGINT): Сумма транзакции в минимальных единицах валюты. Пример: `199900` для 1999.00. Обязательность: Required.
    *   `currency_code` (VARCHAR(3)): Код валюты (RUB, USD, EUR). Обязательность: Required.
    *   `payment_provider` (VARCHAR(50)): Имя платежного провайдера (yookassa, sbp_bankX). Обязательность: Optional (после выбора метода).
    *   `payment_method_type` (VARCHAR(50)): Тип использованного платежного метода (card, sbp, yoomoney_wallet). Обязательность: Optional.
    *   `psp_transaction_id` (VARCHAR(255)): ID транзакции в системе платежного провайдера. Обязательность: Optional.
    *   `description` (TEXT): Описание транзакции. Обязательность: Optional.
    *   `error_code` (VARCHAR(100)): Код ошибки от PSP или внутренний. Обязательность: Optional.
    *   `error_message` (TEXT): Сообщение об ошибке. Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ), `updated_at` (TIMESTAMPTZ).

*   **`PaymentMethod` (Сохраненный Платежный Метод Пользователя)**
    *   `id` (UUID): Уникальный идентификатор. Обязательность: Required.
    *   `user_id` (UUID): ID пользователя. Обязательность: Required.
    *   `provider_name` (VARCHAR(50)): Имя платежного провайдера. Обязательность: Required.
    *   `provider_method_token` (TEXT): Токен платежного метода от провайдера. Валидация: not null. Обязательность: Required.
    *   `method_type` (ENUM: `card`, `e_wallet`): Тип метода. Обязательность: Required.
    *   `details_hint` (JSONB): Маскированные детали для отображения пользователю. Пример: `{"card_type": "Visa", "last_four": "1234", "expiry_year": 2028, "expiry_month": 12}`. Обязательность: Optional.
    *   `is_default` (BOOLEAN): Является ли методом по умолчанию. Обязательность: Required.
    *   `added_at` (TIMESTAMPTZ).

*   **`FiscalReceipt` (Фискальный Чек)**
    *   `id` (UUID): Уникальный идентификатор. Обязательность: Required.
    *   `transaction_id` (UUID, FK to Transaction): ID связанной транзакции. Обязательность: Required.
    *   `receipt_type` (ENUM: `income`, `income_refund`): Тип чека. Обязательность: Required.
    *   `status` (ENUM: `pending_fiscalization`, `fiscalized_success`, `fiscalization_failed`): Статус фискализации. Обязательность: Required.
    *   `ofd_provider_name` (VARCHAR(100)): Имя ОФД-провайдера. Обязательность: Optional.
    *   `fiscal_document_number` (VARCHAR(255)): Номер фискального документа. Обязательность: Optional.
    *   `fiscal_document_sign` (VARCHAR(255)): Фискальный признак документа (ФПД/ФП). Обязательность: Optional.
    *   `receipt_url_ofd` (TEXT): URL чека на сайте ОФД. Обязательность: Optional.
    *   `qr_code_data` (TEXT): Данные для QR-кода чека. Обязательность: Optional.
    *   `created_at` (TIMESTAMPTZ), `fiscalized_at` (TIMESTAMPTZ).

*   **`DeveloperBalance` (Баланс Разработчика)**
    *   `developer_id` (UUID, PK): ID разработчика. Обязательность: Required.
    *   `balance_minor_units` (BIGINT): Текущий баланс в минимальных единицах валюты. Обязательность: Required.
    *   `currency_code` (VARCHAR(3)): Код валюты. Обязательность: Required.
    *   `on_hold_minor_units` (BIGINT): Сумма, замороженная для выплат или других операций. Обязательность: Required.
    *   `last_updated_at` (TIMESTAMPTZ).

*   **`DeveloperPayout` (Выплата Разработчику)**
    *   `id` (UUID): Уникальный идентификатор выплаты.
    *   `developer_id` (UUID, FK to DeveloperBalance): ID разработчика. Обязательность: Required.
    *   `amount_minor_units` (BIGINT): Сумма выплаты. Обязательность: Required.
    *   `currency_code` (VARCHAR(3)): Валюта. Обязательность: Required.
    *   `status` (ENUM: `requested`, `processing`, `completed`, `failed`, `cancelled`). Обязательность: Required.
    *   `payout_method_snapshot` (JSONB): Снимок реквизитов на момент выплаты. Обязательность: Required.
    *   `requested_at` (TIMESTAMPTZ), `processed_at` (TIMESTAMPTZ).
    *   `psp_transaction_id` (VARCHAR(255)): ID транзакции в платежной системе. Обязательность: Optional.

### 4.2. Схема Базы Данных

#### 4.2.1. PostgreSQL

**ERD Диаграмма (ключевые таблицы):**
```mermaid
erDiagram
    TRANSACTIONS {
        UUID id PK
        UUID user_id "FK (User, nullable)"
        UUID order_id "FK (Order, nullable)"
        VARCHAR transaction_type
        VARCHAR status
        BIGINT amount_minor_units
        VARCHAR currency_code
        VARCHAR payment_provider
        VARCHAR payment_method_type
        VARCHAR psp_transaction_id
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    TRANSACTION_ITEMS {
        UUID id PK
        UUID transaction_id FK
        VARCHAR item_id "ID продукта/услуги"
        VARCHAR item_name
        BIGINT price_minor_units
        INT quantity
        VARCHAR vat_code
    }
    USER_PAYMENT_METHODS {
        UUID id PK
        UUID user_id "FK (User)"
        VARCHAR provider_name
        TEXT provider_method_token
        VARCHAR method_type
        JSONB details_hint
        BOOLEAN is_default
        TIMESTAMPTZ added_at
    }
    FISCAL_RECEIPTS {
        UUID id PK
        UUID transaction_id FK
        VARCHAR receipt_type
        VARCHAR status
        VARCHAR ofd_provider_name
        VARCHAR fiscal_document_number
        VARCHAR fiscal_document_sign
        TEXT receipt_url_ofd
        TIMESTAMPTZ fiscalized_at
    }
    DEVELOPER_BALANCES {
        UUID developer_id PK "FK (Developer)"
        BIGINT balance_minor_units
        VARCHAR currency_code
        BIGINT on_hold_minor_units
        TIMESTAMPTZ last_updated_at
    }
    DEVELOPER_PAYOUTS {
        UUID id PK
        UUID developer_id FK
        BIGINT amount_minor_units
        VARCHAR currency_code
        VARCHAR status
        JSONB payout_method_snapshot
        TIMESTAMPTZ requested_at
        TIMESTAMPTZ processed_at
        VARCHAR psp_transaction_id
    }
    PROMO_CODES {
        UUID id PK
        VARCHAR code UK
        VARCHAR discount_type -- percent, fixed_amount
        DECIMAL discount_value
        TIMESTAMPTZ valid_from
        TIMESTAMPTZ valid_to
        INT max_usages
        INT current_usages
        JSONB applicability_rules -- { "product_ids": ["id1"], "min_order_amount": 100000 }
    }

    TRANSACTIONS ||--|{ TRANSACTION_ITEMS : "contains"
    TRANSACTIONS ||--o{ FISCAL_RECEIPTS : "has"
    USERS ||--o{ USER_PAYMENT_METHODS : "has"  # USERS - предполагаемая таблица из Auth/Account
    USERS ||--o{ TRANSACTIONS : "initiates"
    DEVELOPERS ||--|| DEVELOPER_BALANCES : "has" # DEVELOPERS - предполагаемая таблица из Developer Service
    DEVELOPER_BALANCES ||--o{ DEVELOPER_PAYOUTS : "receives"
```

**DDL (PostgreSQL - примеры ключевых таблиц):**
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID, -- Может быть NULL для системных транзакций или выплат
    order_id UUID, -- ID заказа из другого сервиса
    transaction_type VARCHAR(50) NOT NULL CHECK (transaction_type IN ('purchase', 'refund', 'payout', 'deposit', 'hold', 'capture')),
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'awaiting_psp_callback', 'awaiting_capture', 'completed', 'failed', 'cancelled', 'refunded')),
    amount_minor_units BIGINT NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    payment_provider VARCHAR(50),
    payment_method_type VARCHAR(50), -- card, sbp, yoomoney_wallet, etc.
    psp_transaction_id VARCHAR(255) UNIQUE, -- ID транзакции в системе платежного провайдера
    description TEXT,
    error_code VARCHAR(100),
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_transactions_user_id ON transactions(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_transactions_order_id ON transactions(order_id) WHERE order_id IS NOT NULL;
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);

CREATE TABLE user_payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL, -- FK to users table
    provider_name VARCHAR(50) NOT NULL, -- e.g., yookassa, cloudpayments
    provider_method_token TEXT NOT NULL, -- Токен, возвращаемый PSP
    method_type VARCHAR(50) NOT NULL, -- e.g., bank_card, yoo_money
    details_hint JSONB, -- {"card_type": "Visa", "last_four": "1234", "expiry_year": 2028, "expiry_month": 12}
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    added_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, provider_name, provider_method_token) -- У одного юзера не может быть дважды один и тот же токен у одного провайдера
);
CREATE INDEX idx_user_payment_methods_user_id_default ON user_payment_methods(user_id, is_default);

CREATE TABLE fiscal_receipts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE RESTRICT,
    receipt_type VARCHAR(50) NOT NULL CHECK (receipt_type IN ('income', 'income_refund', 'expense', 'expense_refund')),
    status VARCHAR(50) NOT NULL DEFAULT 'pending_fiscalization' CHECK (status IN ('pending_fiscalization', 'fiscalization_requested', 'fiscalized_success', 'fiscalization_failed')),
    ofd_provider_name VARCHAR(100),
    fiscal_document_number VARCHAR(255), -- ФН
    fiscal_storage_number VARCHAR(255), -- ФД
    fiscal_document_sign VARCHAR(255), -- ФП/ФПД
    receipt_url_ofd TEXT,
    qr_code_data TEXT,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    fiscalized_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_fiscal_receipts_transaction_id ON fiscal_receipts(transaction_id);

-- TODO: Добавить DDL для transaction_items, developer_balances, developer_payouts, promo_codes, gift_cards.
```

#### 4.2.2. Redis
*   **Платежные сессии:**
    *   Ключ: `payment_session:<session_token_or_transaction_id>` (HASH). Поля: `user_id`, `order_id`, `amount`, `status`, `psp_data`. TTL: минуты/часы.
*   **Идемпотентность Webhook-обработчиков:**
    *   Ключ: `webhook_processed:<provider_name>:<psp_event_id>` (STRING). Значение: `processed`. TTL: дни.
*   **Rate Limiting для взаимодействия с PSP:**
    *   Ключ: `rl:psp:<provider_name>:<api_endpoint_hash>` (COUNTER).
*   **Кэширование конфигурации платежных провайдеров/ОФД (если не меняется часто):**
    *   Ключ: `psp_config:<provider_name>` (JSON/HASH).
*   **Временное хранение токенов для редиректов или мобильных SDK.**

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Apache Kafka.
*   **Формат событий:** CloudEvents v1.0 (JSON encoding), согласно `project_api_standards.md`.
*   **Основной топик:** `payment.events.v1`.

*   **`payment.transaction.status.updated.v1`** (Заменяет `payment.transaction.completed` и `payment.transaction.failed`)
    *   Описание: Статус транзакции изменился (например, создана, ожидает оплаты, успешно завершена, ошибка, отменена, возвращена).
    *   Пример Payload (`data` секция CloudEvent):
        ```json
        {
          "transaction_id": "txn-uuid-xyz",
          "user_id": "user-uuid-123",
          "order_id": "order-uuid-abc", // Опционально
          "new_status": "completed", // "pending", "failed", "refunded", etc.
          "previous_status": "processing",
          "amount_minor_units": 199900,
          "currency_code": "RUB",
          "payment_provider": "yookassa",
          "payment_method_type_used": "bank_card",
          "updated_at": "2024-03-15T10:05:00Z",
          "error_code": null // или код ошибки, если new_status = "failed"
        }
        ```
    *   Потребители: Library Service (для предоставления доступа к контенту), Notification Service, Analytics Service, Order Service (для обновления статуса заказа).
*   **`payment.payout.status.updated.v1`**
    *   Описание: Статус выплаты разработчику изменился.
    *   Пример Payload:
        ```json
        {
          "payout_id": "payout-uuid-789",
          "developer_id": "dev-uuid-abc",
          "new_status": "completed", // "processing", "failed"
          "processed_at": "2024-03-18T10:00:00Z",
          "amount_minor_units": 5000000,
          "currency_code": "RUB",
          "external_transaction_id": "psp_payout_txn_id", // ID в системе PSP
          "failure_reason": null
        }
        ```
    *   Потребители: Developer Service, Notification Service.
*   **`payment.fiscal.receipt.status.updated.v1`** (Заменяет `payment.fiscal.receipt.created`)
    *   Описание: Статус фискального чека изменился (например, успешно фискализирован, ошибка фискализации).
    *   Пример Payload:
        ```json
        {
          "receipt_id": "receipt-uuid-123",
          "transaction_id": "txn-uuid-xyz",
          "new_status": "fiscalized_success", // "fiscalization_failed"
          "fiscal_document_number": "1234567890123456",
          "fiscal_document_sign": "0987654321098765",
          "receipt_url_ofd": "https://ofd.example.com/check/...",
          "fiscalized_at": "2024-03-15T10:06:00Z",
          "error_message": null
        }
        ```
    *   Потребители: Notification Service (для отправки чека пользователю), Account Service (для истории покупок).

### 5.2. Потребляемые События (Consumed Events)

*   **`order.payment.initiation.requested.v1`** (Заменяет `order.created.v1` для Payment Service)
    *   Источник: Order Service (или другой сервис, инициирующий платеж).
    *   Описание: Запрос на инициирование процесса оплаты для созданного заказа.
    *   Ожидаемый Payload:
        ```json
        {
          "order_id": "order-uuid-abc",
          "user_id": "user-uuid-123",
          "total_amount_minor_units": 199900,
          "currency_code": "RUB",
          "items_for_receipt": [ // Данные для фискального чека
            { "name": "Игра 'Супер Гонки'", "quantity": 1, "price_minor_units": 199900, "vat_code": "VAT_20_CALCULATED", "item_type_code": "DIGITAL_GOODS" }
          ],
          "success_redirect_url": "https://client.example.com/payment/success",
          "fail_redirect_url": "https://client.example.com/payment/failure",
          "preferred_payment_method_id": null // Опционально, ID сохраненного метода
        }
        ```
    *   Логика обработки: Создать транзакцию (`Transaction`) в статусе `pending`. Вызвать `PaymentProcessingApplicationService` для взаимодействия с платежным шлюзом и получения URL для редиректа пользователя или данных для SDK.
*   **`developer.payout.creation.requested.v1`** (Заменяет `developer.payout.requested`)
    *   Источник: Developer Service.
    *   Описание: Запрос на создание и обработку выплаты разработчику.
    *   Ожидаемый Payload:
        ```json
        {
          "payout_request_id_dev_service": "dev-req-uuid-789", // ID запроса из Developer Service
          "developer_id": "dev-uuid-abc",
          "amount_minor_units": 5000000,
          "currency_code": "RUB",
          "payout_method_id_dev_service": "dev-payout-method-uuid-456" // ID метода выплаты из Developer Service
        }
        ```
    *   Логика обработки: Валидировать запрос. Проверить баланс разработчика (`DeveloperBalance`). Создать транзакцию выплаты (`DeveloperPayout`) в статусе `processing`. Инициировать процесс перевода средств через соответствующий платежный шлюз или банковский API. Опубликовать событие `payment.payout.status.updated.v1`.
*   **`admin.transaction.refund.requested.v1`** (Заменяет `admin.refund.request.manual`)
    *   Источник: Admin Service.
    *   Описание: Запрос от администратора на полный или частичный возврат по существующей транзакции.
    *   Ожидаемый Payload:
        ```json
        {
          "original_transaction_id": "txn-uuid-xyz",
          "amount_minor_units_to_refund": 100000, // 0 для полного возврата исходной суммы
          "reason_code": "customer_request", // "fraud", "technical_issue"
          "reason_description": "Пользователь попросил возврат в течение 14 дней.",
          "admin_user_id": "admin-user-uuid-superuser"
        }
        ```
    *   Логика обработки: Найти оригинальную транзакцию. Проверить возможность возврата. Создать транзакцию возврата (`Transaction` с типом `refund`). Инициировать возврат через платежный шлюз. Инициировать фискализацию чека возврата. Опубликовать `payment.transaction.status.updated.v1` для транзакции возврата и, возможно, для оригинальной транзакции.
*   **`user.subscription.renewal.payment.requested.v1`** (Заменяет `user.subscription.renewal_due.v1`)
    *   Источник: Subscription Service (или Account Service).
    *   Описание: Запрос на автоматическое списание средств за продление подписки.
    *   Ожидаемый Payload:
        ```json
        {
            "user_id": "user-uuid-123",
            "subscription_id": "sub-uuid-qwerty",
            "renewal_amount_minor_units": 49900,
            "currency_code": "RUB",
            "default_payment_method_id": "pm-uuid-default", // ID сохраненного платежного метода
            "items_for_receipt": [
                { "name": "Подписка 'Плюс' на 1 месяц", "quantity": 1, "price_minor_units": 49900, "vat_code": "VAT_20_CALCULATED", "item_type_code": "SERVICE" }
            ]
        }
        ```
    *   Логика обработки: Попытаться выполнить рекуррентный платеж с использованием сохраненного платежного метода пользователя. Создать транзакцию, фискализировать, опубликовать события.

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Account Service:** Получение данных о пользователях (например, для фискальных чеков, проверки лимитов).
*   **Catalog Service:** Получение информации о товарах/услугах для фискальных чеков (названия, ставки НДС, если не переданы в запросе).
*   **Library Service:** Уведомление о результате покупки для предоставления доступа к контенту (через Kafka).
*   **Developer Service:** Управление балансами и выплатами разработчиков (через Kafka и, возможно, gRPC).
*   **Admin Service:** Предоставление API для административных операций (просмотр транзакций, ручные возвраты).
*   **Analytics Service:** Отправка данных о всех транзакциях для финансовой аналитики (через Kafka).
*   **Notification Service:** Инициирование отправки уведомлений пользователям и разработчикам о статусе платежей, возвратов, выплат, фискальных чеков (через Kafka).
*   **API Gateway:** Прием запросов от клиентов и маршрутизация их в Payment Service.
*   **Auth Service:** Валидация JWT токенов (обычно на уровне API Gateway).

### 6.2. Внешние Системы
*   **Платежные Шлюзы (СБП, МИР, ЮMoney, эквайринг банков):** Интеграция через их API для инициирования платежей, обработки колбэков, проведения возвратов и выплат.
*   **Операторы Фискальных Данных (ОФД):** Интеграция через их API для регистрации касс, отправки данных фискальных чеков и получения статусов.

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `PAYMENT_SERVICE_HTTP_PORT`: Порт для REST API.
*   `PAYMENT_SERVICE_GRPC_PORT` (если используется).
*   `POSTGRES_DSN`: DSN для PostgreSQL.
*   `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB_PAYMENT`: Параметры Redis.
*   `KAFKA_BROKERS`: Список брокеров Kafka.
*   `KAFKA_TOPIC_PAYMENT_EVENTS`: Топик для публикуемых событий Payment Service.
*   `KAFKA_CONSUMER_GROUP_ID_PAYMENT`: ID группы консьюмеров для входящих событий.
*   Различные API ключи, ID мерчантов, секреты для каждого платежного шлюза (например, `SBP_MERCHANT_ID_SECRET`, `YOOKASSA_SHOP_ID_SECRET`, `YOOKASSA_API_KEY_SECRET`). **Эти значения должны извлекаться из системы управления секретами (например, Vault или Kubernetes Secrets).**
*   Параметры для ОФД (например, `OFD_API_URL`, `OFD_API_KEY_SECRET`, `OFD_INN`, `OFD_KKT_REG_NUMBER`).
*   `PLATFORM_COMMISSION_PERCENTAGE`: Процент комиссии платформы с продаж.
*   `MIN_DEVELOPER_PAYOUT_AMOUNT_RUB`: Минимальная сумма для выплаты разработчику.
*   `DEFAULT_PAYMENT_CURRENCY`: Валюта по умолчанию (например, `RUB`).
*   `LOG_LEVEL`.
*   `AUTH_SERVICE_GRPC_ADDR`, `ACCOUNT_SERVICE_GRPC_ADDR`, `CATALOG_SERVICE_GRPC_ADDR`.
*   `NOTIFICATION_SERVICE_KAFKA_TOPIC_SEND_REQUESTS`.
*   `OTEL_EXPORTER_JAEGER_ENDPOINT`.
*   `PAYMENT_SESSION_TTL_SECONDS`: Время жизни сессии платежа на стороне PSP (например, 15-30 минут).
*   `TRANSACTION_LOCK_TIMEOUT_SECONDS`: Таймаут для блокировки транзакции при обработке.

### 7.2. Файлы Конфигурации (если применимо)
*   **`configs/payment_config.yaml`**: Может использоваться для более сложных или менее часто изменяемых настроек.
    ```yaml
    server:
      http_port: ${PAYMENT_SERVICE_HTTP_PORT:-8086}
      # grpc_port: ${PAYMENT_SERVICE_GRPC_PORT:-9096}

    # Параметры для взаимодействия с платежными шлюзами (не секреты)
    payment_providers:
      yookassa:
        base_url: "https://api.yookassa.ru/v3"
        timeout_seconds: 30
        retry_policy: { "max_attempts": 3, "backoff_ms": 1000 }
      sbp:
        # ...

    fiscalization:
      default_ofd_provider: "ofd_platforma" # Имя конфигурации ОФД-провайдера
      ofd_platforma:
        api_url: ${OFD_API_URL}
        inn: ${OFD_INN}
        kkt_reg_number: ${OFD_KKT_REG_NUMBER}
        default_tax_system_code: 1 # ОСН
        default_vat_rate_code_items: "VAT_20_CALCULATED" # Ставка НДС для товаров
        default_payment_object_code_items: "COMMODITY" # Признак предмета расчета
      # Настройки таймаутов и ретраев для ОФД

    commission:
      platform_percentage: ${PLATFORM_COMMISSION_PERCENTAGE:-30.0}
      # Могут быть более сложные правила расчета комиссии

    payouts:
      min_amount_rub: ${MIN_DEVELOPER_PAYOUT_AMOUNT_RUB:-1000.0}
      # Правила для разных методов выплат, лимиты и т.д.

    idempotency:
      webhook_processing_ttl_seconds: 86400 # 24 часа для хранения ключей идемпотентности вебхуков
    ```

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   Корректная обработка и логирование ошибок от платежных систем и ОФД. Использование специфичных кодов ошибок от PSP, если это возможно и полезно.
*   Механизмы повторных попыток (retry) с экспоненциальной задержкой для временных сбоев при взаимодействии с внешними системами.
*   Использование ключей идемпотентности при работе с API платежных систем для предотвращения двойных списаний/возвратов.
*   Четкие сообщения об ошибках для пользователей и администраторов, не раскрывающие излишних технических деталей.
*   Для критичных операций (например, подтверждение платежа) – механизм отложенных задач и очередей для гарантии выполнения.

### 8.2. Распространенные Коды Ошибок (для REST API)
*   **`400 Bad Request` (`VALIDATION_ERROR`, `INVALID_PAYMENT_METHOD`, `INSUFFICIENT_ORDER_DATA`)**: Некорректные входные данные.
*   **`401 Unauthorized` (`UNAUTHENTICATED`)**: Ошибка аутентификации.
*   **`402 Payment Required` (`INSUFFICIENT_FUNDS`, `PAYMENT_DECLINED_BY_PSP`, `CARD_EXPIRED`)**: Платеж не прошел.
*   **`403 Forbidden` (`PERMISSION_DENIED`, `OPERATION_NOT_ALLOWED_FOR_TRANSACTION_STATUS`)**: Недостаточно прав или операция невозможна для текущего статуса.
*   **`404 Not Found` (`TRANSACTION_NOT_FOUND`, `PAYMENT_METHOD_NOT_FOUND`)**: Запрашиваемый ресурс не найден.
*   **`409 Conflict` (`IDEMPOTENCY_KEY_VIOLATION`, `PAYMENT_ALREADY_PROCESSED`)**: Конфликт состояния (например, повторный запрос с тем же ключом идемпотентности, но другими данными).
*   **`422 Unprocessable Entity` (`FISCALIZATION_ERROR`, `REFUND_POLICY_VIOLATION`)**: Ошибка бизнес-логики (например, не удалось фискализировать чек, нарушение правил возврата).
*   **`500 Internal Server Error` (`INTERNAL_ERROR`)**: Внутренняя ошибка сервера.
*   **`502 Bad Gateway` (`PSP_ERROR`, `OFD_ERROR`)**: Ошибка на стороне платежного шлюза или ОФД.
*   **`503 Service Unavailable` (`SERVICE_UNAVAILABLE`)**: Сервис временно недоступен или его зависимости.
*   **`504 Gateway Timeout` (`PSP_TIMEOUT`, `OFD_TIMEOUT`)**: Таймаут при ожидании ответа от внешнего сервиса.

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   JWT (через Auth Service) для всех пользовательских и административных API.
*   Защищенные API-ключи или mTLS для межсервисных вызовов.
*   Аутентификация запросов от внешних платежных систем (вебхуки) через проверку подписи или IP-адресов.

### 9.2. Авторизация
*   RBAC для доступа к операциям (например, только пользователь может инициировать платеж от своего имени, только администратор может выполнять определенные типы возвратов или корректировки балансов).
*   Проверка принадлежности транзакций и платежных методов конкретному пользователю.

### 9.3. Защита Данных
*   **PCI DSS Compliance:** Прямая обработка полных номеров банковских карт **не предполагается**. Сервис должен полагаться на токенизацию на стороне платежных шлюзов (PSP) и использование iFrame/редиректов на страницы оплаты PSP для минимизации области действия PCI DSS. Если в будущем потребуется работа с PAN, этот раздел должен быть кардинально пересмотрен с привлечением специалистов по PCI DSS.
*   Шифрование чувствительных данных при хранении (например, токенизированные платежные методы, если они кэшируются/хранятся локально; реквизиты для выплат разработчикам) и при передаче (TLS 1.2+ для всех коммуникаций).
*   Защита от мошенничества (fraud prevention): базовая проверка транзакций, возможно, интеграция с внешними антифрод-системами.
*   Соблюдение требований 54-ФЗ (фискализация), ФЗ-152 "О персональных данных", ФЗ-115 "О противодействии легализации (отмыванию) доходов..." (AML).

### 9.4. Управление Секретами
*   API ключи для взаимодействия с платежными системами, ОФД, ключи шифрования должны храниться в HashiCorp Vault или Kubernetes Secrets и безопасно внедряться в приложение.
*   Регулярная ротация секретов.
*   **Аудит**: Детальное логирование всех финансовых операций, изменений статусов, действий администраторов.

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Для сборки сервиса.
*   **Kubernetes манифесты/Helm-чарты:** Для развертывания и управления в Kubernetes.
*   (Ссылка на `project_deployment_standards.md`).

### 10.2. Зависимости при Развертывании
*   Кластер Kubernetes.
*   PostgreSQL, Redis, Kafka.
*   Доступность Auth Service, Account Service, Catalog Service, Notification Service, Developer Service, Admin Service.
*   Сетевой доступ к API внешних платежных шлюзов и ОФД.

### 10.3. CI/CD
*   Автоматизированная сборка, тестирование (модульное, интеграционное с моками внешних систем, компонентное).
*   Статический анализ кода. Сканирование на уязвимости.
*   Развертывание по окружениям (dev, staging, production) с использованием GitOps.
*   Процедуры миграции схемы БД.

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат:** Структурированные JSON логи.
*   **Ключевые события:** Все этапы обработки транзакций (создание, запрос к PSP, колбэк от PSP, фискализация, обновление статуса), операции с балансами, выплатами, промокодами, ошибки интеграции.
*   **Интеграция:** С централизованной системой логирования (ELK Stack, Loki/Grafana).
*   (Ссылка на `project_observability_standards.md`).

### 11.2. Мониторинг
*   **Метрики (Prometheus):**
    *   Количество и суммы транзакций (успешных, неуспешных, в обработке) по типам, платежным методам, провайдерам.
    *   Время обработки платежей (P50, P90, P95, P99).
    *   Количество ошибок при взаимодействии с платежными системами и ОФД.
    *   Размеры очередей Kafka, используемых сервисом.
    *   Производительность и ошибки при работе с PostgreSQL, Redis.
    *   Количество успешных и неуспешных фискализаций.
*   **Дашборды (Grafana):** Для визуализации финансовых потоков, состояния интеграций, производительности.
*   **Алертинг (AlertManager):** Для критических ошибок (например, высокий процент отказов платежей, сбои фискализации, недоступность PSP/ОФД, аномальное изменение объемов транзакций).
*   (Ссылка на `project_observability_standards.md`).

### 11.3. Трассировка
*   **Интеграция:** OpenTelemetry SDK. Экспорт трейсов в Jaeger или Tempo.
*   Трассировка полного жизненного цикла финансовой операции: от запроса на инициирование платежа до его завершения и фискализации, включая все вызовы к внутренним и внешним системам.
*   (Ссылка на `project_observability_standards.md`).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:**
    *   Время ответа API для инициирования платежа (до редиректа на PSP): P95 < 500 мс.
    *   Время обработки колбэка от PSP (до постановки в очередь на фискализацию и отправки события): P95 < 200 мс.
    *   Пропускная способность: поддержка не менее 100 TPS (транзакций в секунду) на создание, до 500 TPS на проверку статуса.
*   **Надежность:**
    *   Доступность сервиса: >= 99.95% (для критичных операций платежей).
    *   Отсутствие потерь финансовых транзакций (RPO = 0).
    *   Корректность финансовых расчетов и учета балансов.
    *   RTO < 1 часа для восстановления после сбоя.
*   **Безопасность:**
    *   Соответствие требованиям PCI DSS (минимизация области действия).
    *   Надежная защита от мошеннических операций (фрод-мониторинг, возможно, интеграция с антифрод-системами).
    *   Соблюдение 54-ФЗ и других релевантных законодательных актов.
*   **Масштабируемость:** Горизонтальное масштабирование для обработки пиковых нагрузок (например, во время распродаж).
*   **Сопровождаемость:** Четкое логирование, полные трассировки, подробные метрики.

## 13. Приложения (Appendices)
*   Детальные OpenAPI схемы для REST API и Protobuf определения для gRPC API будут храниться в соответствующих репозиториях или артефактах CI/CD.
*   Полные DDL схемы базы данных будут поддерживаться в актуальном состоянии в системе миграций.
*   Примеры фискальных чеков и форматы взаимодействия с конкретными ОФД.
*   TODO: Добавить ссылки на репозиторий с OpenAPI/Protobuf и на систему управления миграциями БД, когда они будут определены.

---
*Этот документ является основной спецификацией для Payment Service и должен поддерживаться в актуальном состоянии.*

## 14. Связанные Рабочие Процессы (Related Workflows)
*   [Game Purchase and Library Update](../../../project_workflows/game_purchase_flow.md)
*   [Developer Payout Flow](../../../project_workflows/developer_payout_flow.md) (TODO: Создать этот документ)
*   [Refund Processing Flow](../../../project_workflows/refund_processing_flow.md) (TODO: Создать этот документ)
