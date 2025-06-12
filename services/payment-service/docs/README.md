# Спецификация Микросервиса: Payment Service

**Версия:** 1.0
**Дата последнего обновления:** {{YYYY-MM-DD}} <!-- TODO: Update date -->

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
*   TODO: Определить gRPC API для внутреннего высокопроизводительного взаимодействия, если требуется (например, для проверки статуса транзакции, получения цены). Исходная спецификация упоминает gRPC в технологическом стеке.

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
    *   `payment.developer.payout.completed`: Выплата разработчику произведена. -> Developer Service, Notification Service.
    *   `payment.fiscal.receipt.created`: Фискальный чек создан. -> Notification Service (для отправки пользователю).
*   TODO: Детализировать структуру Payload для каждого события.

### 5.2. Потребляемые События (Consumed Events)
*   `catalog.price.changed`: От Catalog Service, для обновления информации о ценах, если есть активные корзины.
*   `user.account.deleted`: От Account Service, для обработки данных пользователя.
*   `admin.refund.request.manual`: От Admin Service, для ручного инициирования возврата.
*   TODO: Детализировать другие возможные потребляемые события.

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
*   `POSTGRES_DSN`.
*   `REDIS_ADDR`.
*   `KAFKA_BROKERS`.
*   API ключи и эндпоинты для платежных систем и ОФД (хранятся в Secrets).
*   Параметры комиссии платформы.
*   `LOG_LEVEL`.
*   TODO: Сформировать полный список.

### 7.2. Файлы Конфигурации (если применимо)
*   Могут использоваться для настроек правил фискализации, лимитов операций.
*   TODO: Детализировать структуру, если используется.

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
*   TODO: Уточнить интеграцию с Jaeger/OpenTelemetry для отслеживания полных циклов финансовых операций.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:** API P95 < 500мс, обработка >= 100 транзакций/сек.
*   **Надежность:** Доступность >= 99.9%.
*   **Безопасность:** Шифрование, PCI DSS (если применимо), ФЗ-54, ФЗ-152, ФЗ-115.
*   **Масштабируемость:** Горизонтальная.
*   (Детали см. в Спецификации Payment Service, раздел 2.4).

## 13. Приложения (Appendices) (Опционально)
*   TODO: Детальные схемы DDL, примеры API запросов/ответов, форматы фискальных чеков.

---
*Этот шаблон является отправной точкой и может быть адаптирован под конкретные нужды проекта и сервиса.*
