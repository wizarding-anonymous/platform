# Интеграции Проекта

## 1. Введение

Данный документ описывает основные интеграции между микросервисами платформы "Российский Аналог Платформы Steam", а также интеграции с внешними системами и особенности взаимодействия с Flutter-клиентом. Цель документа — предоставить общее представление о потоках данных и зависимостях между компонентами системы.

## 2. Внутренние Интеграции (Между Микросервисами)

Внутренние интеграции обеспечивают взаимодействие между различными микросервисами платформы для выполнения бизнес-функций.

### 2.1. Общая Карта Интеграций

```mermaid
graph TD
    FlutterClient --> APIGateway[API Gateway]

    APIGateway --> AuthService[Auth Service]
    APIGateway --> AccountService[Account Service]
    APIGateway --> CatalogService[Catalog Service]
    APIGateway --> LibraryService[Library Service]
    APIGateway --> DownloadService[Download Service]
    APIGateway --> PaymentService[Payment Service]
    APIGateway --> SocialService[Social Service]
    APIGateway --> DeveloperService[Developer Service]
    APIGateway --> AdminService[Admin Service]
    APIGateway --> AnalyticsService[Analytics Service]
    APIGateway --> NotificationService[Notification Service]

    AuthService <--> AccountService
    AuthService <--> PaymentService
    AuthService <--> DeveloperService
    AuthService <--> AdminService

    AccountService <--> SocialService
    AccountService <--> LibraryService
    AccountService <--> PaymentService
    AccountService <--> DeveloperService
    AccountService <--> AdminService
    AccountService --> AnalyticsService
    AccountService --> NotificationService

    CatalogService <--> LibraryService
    CatalogService <--> DownloadService
    CatalogService <--> PaymentService
    CatalogService <--> DeveloperService
    CatalogService <--> AdminService
    CatalogService --> AnalyticsService

    LibraryService <--> DownloadService
    LibraryService <--> SocialService
    LibraryService <--> AdminService
    LibraryService --> AnalyticsService
    LibraryService --> NotificationService

    DownloadService <--> DeveloperService
    DownloadService <--> AdminService
    DownloadService --> AnalyticsService
    DownloadService --> NotificationService

    PaymentService <--> DeveloperService
    PaymentService <--> AdminService
    PaymentService --> AnalyticsService
    PaymentService --> NotificationService

    SocialService <--> AdminService
    SocialService --> AnalyticsService
    SocialService --> NotificationService

    DeveloperService <--> AdminService
    DeveloperService --> AnalyticsService
    DeveloperService --> NotificationService

    AdminService <--> AnalyticsService
    AdminService <--> NotificationService

    AnalyticsService --> NotificationService # Для алертов и отчетов, возможно
```

### 2.2. Матрица Интеграций (Типы взаимодействия)

| Микросервис        | API Gateway | Auth         | Account      | Catalog      | Library      | Download     | Payment      | Social       | Developer    | Admin        | Analytics    | Notification |
|--------------------|-------------|--------------|--------------|--------------|--------------|--------------|--------------|--------------|--------------|--------------|--------------|--------------|
| **API Gateway**    | -           | REST         | REST         | REST         | REST         | REST         | REST         | REST         | REST         | REST         | REST         | WebSocket    |
| **Auth Service**   | REST        | -            | REST, Events | -            | -            | -            | REST         | -            | REST         | REST         | -            | -            |
| **Account Service**| REST        | REST, Events | -            | -            | REST         | -            | REST         | REST         | REST         | REST         | Events       | Events       |
| **Catalog Service**| REST        | -            | -            | -            | REST         | REST         | REST         | -            | REST         | REST         | Events       | -            |
| **Library Service**| REST        | -            | REST         | REST         | -            | REST         | -            | REST         | -            | REST         | Events       | Events       |
| **Download Service**| REST       | -            | -            | REST         | REST         | -            | -            | -            | REST         | REST         | Events       | Events       |
| **Payment Service**| REST        | REST         | REST         | REST         | -            | -            | -            | -            | REST         | REST         | Events       | Events       |
| **Social Service** | REST        | -            | REST         | -            | REST         | -            | -            | -            | -            | REST         | Events       | Events       |
| **Developer Service**| REST      | REST         | REST         | REST         | -            | REST         | REST         | -            | -            | REST         | Events       | Events       |
| **Admin Service**  | REST        | REST         | REST         | REST         | REST         | REST         | REST         | REST         | REST         | -            | REST         | REST         |
| **Analytics Service**| REST      | -            | Events       | Events       | Events       | Events       | Events       | Events       | Events       | REST         | -            | -            |
| **Notification Service**| WebSocket| -          | Events       | -            | Events       | Events       | Events       | Events       | Events       | REST         | -            | -            |

*(Источник: Аудит интеграций между микросервисами.txt)*

### 2.3. Примеры Детальных Интеграций (Резюме)

*   **API Gateway → Auth Service:** Проксирование запросов аутентификации, валидация JWT, обновление токенов.
*   **Auth Service → Account Service:** Создание профиля пользователя при регистрации (через REST API или событие `user.registered`).
*   **Catalog Service → Payment Service:** Предоставление актуальных цен и скидок для формирования заказов.
*   **Library Service → Download Service:** Проверка прав доступа к игре для начала загрузки.
*   **Notification Service ← (Другие сервисы):** Получение событий от многих сервисов (Account, Auth, Library, Payment, Social, etc.) для отправки уведомлений пользователям.

*(Более детальный анализ каждой пары интеграций должен быть представлен в спецификациях соответствующих микросервисов. См. также "Аудит интеграций между микросервисами.txt" для анализа несоответствий и рекомендаций на момент аудита).*

## 3. Интеграция с Flutter-клиентом

Flutter-клиент взаимодействует с бэкендом исключительно через API Gateway.

### 3.1. Типы Взаимодействий
*   **REST API:** Основной метод для большинства операций (CRUD, авторизация, получение данных).
*   **WebSocket:** Для функций реального времени (чат, уведомления, статусы).
*   **gRPC (опционально):** Может рассматриваться для критичных к производительности сценариев.

### 3.2. Форматы Данных
*   **JSON:** Для REST API и сообщений WebSocket.
*   **Protocol Buffers:** Для gRPC (если используется).

### 3.3. Аутентификация и Авторизация
*   **JWT-токены:** Access и Refresh токены, безопасное хранение (`flutter_secure_storage`).
*   **Telegram авторизация:** Интеграция с Telegram Login Widget.
*   **Биометрическая аутентификация:** Локальная, с использованием `local_auth`.

### 3.4. Управление Сессиями
*   Автоматическое обновление токенов.
*   Многоуровневая авторизация (API Gateway + микросервисы).

### 3.5. Оптимизация для Мобильных Устройств
*   Пагинация, частичная загрузка данных.
*   Ленивая загрузка изображений.
*   Локальное кэширование (Hive, SQLite).
*   Сжатие данных (gzip).
*   Поддержка условных запросов (ETag, If-Modified-Since).

### 3.6. Обработка Ошибок и Отказоустойчивость
*   Стандартизированные ошибки API.
*   Локализованные сообщения.
*   Стратегии повторных попыток, поддержка оффлайн-режима с очередью операций.
*   Мониторинг клиентских ошибок.

*(Источник: Аудит интеграций между микросервисами.txt, раздел "Интеграция с Flutter-клиентом")*

## 4. Внешние Интеграции

Платформа интегрируется с рядом внешних систем для обеспечения своей функциональности.

### 4.1. Платежные Системы и Финансовые Сервисы
*   **Система быстрых платежей (СБП):** Обработка платежей. (Payment Service)
*   **Платежная система МИР:** Обработка платежей по картам МИР. (Payment Service)
*   **ЮMoney:** Обработка платежей через электронные кошельки. (Payment Service)
*   **Оператор фискальных данных (ОФД):** Формирование и регистрация фискальных чеков (54-ФЗ). (Payment Service)
*   **Банковские API (Сбербанк, Тинькофф, Альфа-Банк):** Проведение платежей. (Payment Service)

### 4.2. Системы Аутентификации и Авторизации
*   **ВКонтакте OAuth:** Аутентификация через ВКонтакте. (Auth Service)
*   **Telegram Login:** Аутентификация через Telegram. (Auth Service)
*   **Одноклассники OAuth:** Аутентификация через Одноклассники. (Auth Service)

### 4.3. Системы Уведомлений и Коммуникаций
*   **Email-провайдеры (SendPulse, Unisender, MailRu Cloud):** Отправка email-уведомлений. (Notification Service)
*   **SMS-шлюзы (SMSC, SMS.ru, МТС Коммуникатор):** Отправка SMS. (Notification Service, Auth Service для 2FA)
*   **Push-уведомления (Firebase Cloud Messaging - FCM):** Push на Android. (Notification Service)
*   **Push-уведомления (Apple Push Notification Service - APNS):** Push на iOS/macOS. (Notification Service)
*   **Telegram Bot API:** Отправка уведомлений через Telegram. (Notification Service)

### 4.4. Облачные Хранилища и CDN
*   **S3-совместимое хранилище (Yandex Object Storage, VK Cloud, SberCloud):** Хранение игровых файлов, медиа-контента. (Download Service, Catalog Service, Developer Service)
*   **CDN (Content Delivery Network):** Ускорение доставки контента. (Download Service, Catalog Service)

### 4.5. Аналитические и Мониторинговые Системы
*   **Яндекс.Метрика:** Сбор и анализ поведения пользователей. (Analytics Service)
*   **Sentry / Rollbar:** Мониторинг ошибок приложений. (Все микросервисы)

### 4.6. Системы Защиты и Безопасности
*   **Captcha (reCAPTCHA, hCaptcha, Yandex SmartCaptcha):** Защита от ботов. (Auth Service, Social Service)
*   **Антивирусные API (Kaspersky, Dr.Web):** Проверка загружаемых файлов. (Developer Service, Download Service)

### 4.7. Геолокационные и IP-сервисы
*   **GeoIP-базы (MaxMind, IP-API):** Определение местоположения по IP. (Auth Service, Analytics Service)

### 4.8. Интеграции с Социальными Сетями (для постинга и др.)
*   **ВКонтакте API:** Публикация контента, приглашения. (Social Service)
*   **Одноклассники API:** Публикация контента, приглашения. (Social Service)

### 4.9. Системы Проверки Возраста и Идентификации
*   **Системы верификации возраста:** Проверка возраста для доступа к контенту. (Auth Service, Catalog Service)

### 4.10. Системы Логирования и Мониторинга (Внешние аспекты)
*   **Elasticsearch / Logstash / Kibana (ELK Stack):** Централизованный сбор логов. (Все микросервисы)
*   **Prometheus / Grafana:** Сбор метрик и мониторинг. (Все микросервисы)

*(Источник: Список внешних интеграций российского аналога Steam.txt)*

## 5. Стандартизированные Контракты Интеграций

Для обеспечения согласованности и упрощения взаимодействия между микросервисами, а также с Flutter-клиентом, должны быть разработаны и приняты стандартизированные контракты:

*   **REST API Контракты:**
    *   Единый формат успешных ответов (включая пагинацию).
    *   Единый формат ошибок.
    *   (См. "Стандарты API, форматов данных, событий и конфигурационных файлов.txt" и "Аудит интеграций между микросервисами.txt")
*   **WebSocket Контракты:**
    *   Стандартизированный формат сообщений (тип, полезная нагрузка).
    *   Формат подтверждений.
    *   (См. "Аудит интеграций между микросервисами.txt")
*   **События (Kafka/CloudEvents):**
    *   Единый формат событий (например, CloudEvents).
    *   Четко определенные схемы (`data` payload) для каждого типа события.
    *   Централизованный реестр событий.
    *   (См. "Стандарты API, форматов данных, событий и конфигурационных файлов.txt")
*   **gRPC API Контракты:**
    *   Согласованное именование сервисов, методов и сообщений.
    *   Стандартные типы данных для общих сущностей.
    *   Единые подходы к обработке ошибок и передаче метаданных.

Детальные стандарты для каждого типа контракта описаны ниже. Эти стандарты основаны на документах "Стандарты API, форматов данных, событий и конфигурационных файлов.txt" и "Аудит интеграций между микросервисами.txt".

### 5.1. REST API Контракты

Все REST API должны соответствовать следующим принципам:

*   **Версионирование:** Версия API указывается в URL, например: `/api/v1/resource`. Мажорная версия (v1, v2) изменяется при несовместимых изменениях API. Минорные изменения (например, добавление новых полей в ответ) не требуют изменения версии.
*   **Формат URL:**
    *   Использовать существительные во множественном числе для идентификации ресурсов: `/api/v1/users`, `/api/v1/games`.
    *   Использовать вложенные ресурсы для выражения логических связей: `/api/v1/games/{game_id}/reviews`.
    *   Использовать kebab-case для имен ресурсов, состоящих из нескольких слов: `/api/v1/payment-methods`.
    *   Избегать использования глаголов в URL, за исключением случаев, когда действие не вписывается в стандартные CRUD-операции над ресурсом. Такие специальные действия оформляются через суффикс `/action`, например: `/api/v1/games/{game_id}/publish`.
*   **HTTP-методы:**
    *   `GET`: Получение ресурса или коллекции ресурсов.
    *   `POST`: Создание нового ресурса.
    *   `PUT`: Полное обновление существующего ресурса.
    *   `PATCH`: Частичное обновление существующего ресурса.
    *   `DELETE`: Удаление ресурса.
*   **Коды ответов HTTP:** Использовать стандартные коды состояния HTTP для индикации результата операции (200 OK, 201 Created, 204 No Content, 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 409 Conflict, 422 Unprocessable Entity, 429 Too Many Requests, 500 Internal Server Error).
*   **Пагинация:**
    *   Для коллекций ресурсов использовать параметры запроса `page` (номер страницы, начиная с 1) и `per_page` (количество элементов на странице, макс. 100).
    *   Ответ должен содержать объект `meta` с информацией о пагинации (`page`, `per_page`, `total_pages`, `total_items`) и объект `links` со ссылками на текущую, первую, предыдущую, следующую и последнюю страницы.
*   **Фильтрация:** Фильтрация данных осуществляется через параметры запроса (например, `/api/v1/games?genre=strategy&price_min=500`). Для сложных сценариев фильтрации может использоваться параметр `filter` с JSON-структурой.
*   **Сортировка:** Параметр `sort` используется для сортировки (например, `sort=price` для сортировки по возрастанию цены, `sort=-price` для сортировки по убыванию). Множественная сортировка поддерживается через запятую: `sort=genre,-price`.
*   **Выборка полей (Field Selection):** Клиенты могут запрашивать только необходимые поля с помощью параметра `fields` (например, `fields=id,title,developer{id,name}`).
*   **Формат ответа:**
    *   Успешные ответы для одиночного ресурса: `{ "data": { "id": "...", "type": "...", "attributes": { ... }, "relationships": { ... } } }`.
    *   Успешные ответы для коллекций: `{ "data": [ ... ], "meta": { ... }, "links": { ... } }`.
    *   Имена полей в JSON должны использовать `camelCase`.
*   **Формат ошибок:** Ответы об ошибках должны иметь стандартизированную структуру:
    ```json
    {
      "errors": [
        {
          "code": "ERROR_CODE_NAME",
          "title": "Человекочитаемый заголовок ошибки",
          "detail": "Детальное описание проблемы.",
          "source": { "pointer": "/data/attributes/field_name" } // Опционально, указывает на источник ошибки
        }
      ]
    }
    ```
*   **Общие заголовки:**
    *   `Content-Type: application/json` для тела запроса и ответа.
    *   `Accept: application/json` для указания предпочитаемого формата ответа.
    *   `Authorization: Bearer <jwt_token>` для аутентификации.
    *   `X-Request-ID: <uuid>` для трассировки запросов.
    *   При взаимодействии микросервисов через API Gateway, последний добавляет заголовки: `X-User-Id`, `X-User-Roles`, `X-Original-IP`.
*   **Документация:** Каждый REST API должен быть документирован с использованием спецификации OpenAPI 3.0 (Swagger). Документация должна быть доступна по стандартному пути, например, `/api/v1/docs`.

### 5.2. WebSocket Контракты

*   **Подключение:** URL для подключения имеет вид `/api/v1/ws/{service_name}`. Аутентификация при установлении соединения осуществляется через параметр запроса `token` или стандартный заголовок `Authorization: Bearer <jwt_token>`. Должна быть поддержка Ping/Pong фреймов для поддержания соединения и проверки его активности.
*   **Формат сообщений:** Все сообщения передаются в формате JSON и имеют следующую базовую структуру:
    ```json
    {
      "type": "unique_message_type_string", // Тип сообщения (например, "chat.message.send", "notification.new")
      "id": "client_generated_uuid",       // Уникальный ID сообщения (для отслеживания и подтверждений)
      "payload": { ... }                     // Полезная нагрузка, специфичная для типа сообщения
    }
    ```
*   **Обработка ошибок:** Ошибки на уровне WebSocket передаются специальным типом сообщения:
    ```json
    {
      "type": "error",
      "id": "original_message_id_if_applicable", // ID сообщения, вызвавшего ошибку
      "payload": {
        "code": "ERROR_CODE_NAME",
        "message": "Описание ошибки"
      }
    }
    ```
*   **Подтверждения (Acknowledgements):** Для критически важных сообщений, где требуется гарантия доставки, может использоваться механизм подтверждений. Клиент или сервер, получив сообщение, может отправить в ответ:
    ```json
    {
      "type": "ack",
      "id": "original_message_id",         // ID подтверждаемого сообщения
      "payload": {
        "status": "delivered" // или "received", "processed"
      }
    }
    ```
*   **Оптимизация для мобильных клиентов:** Необходимо учитывать особенности мобильных платформ, такие как управление энергопотреблением и нестабильность сети, реализуя механизмы переподключения и сжатия сообщений.

### 5.3. События (Kafka/CloudEvents)

Для асинхронного взаимодействия между сервисами используются события, передаваемые через Apache Kafka. Формат событий должен стремиться к совместимости со спецификацией CloudEvents.

*   **Формат события:**
    ```json
    {
      "id": "unique_event_uuid_v4",              // Уникальный идентификатор события
      "type": "com.projectname.domain.resource.action.v1", // Тип события (например, "com.gameplatform.user.registered.v1")
      "source": "/service_name/resource_path",    // Источник события (имя сервиса, опционально путь к ресурсу)
      "specversion": "1.0",                       // Версия спецификации CloudEvents (если применимо)
      "time": "ISO8601_timestamp_utc",            // Время генерации события в UTC
      "datacontenttype": "application/json",      // Тип контента в поле data
      "subject": "entity_id_or_relevant_identifier", // Идентификатор субъекта события (например, user_id, game_id)
      "correlationid": "trace_or_request_uuid",   // ID для корреляции событий в рамках одной операции/запроса
      "data": { ... }                             // Полезная нагрузка события, специфичная для типа
    }
    ```
*   **Именование типов событий (`type`):** Используется обратный DNS-нотации стиль, включающий домен, имя ресурса, выполненное действие и версию. Например: `com.gameplatform.user.registered.v1`, `com.gameplatform.payment.completed.v1`. Действия именуются в прошедшем времени.
*   **Версионирование событий:** Версия схемы полезной нагрузки (`data`) включается в поле `type`. При несовместимых изменениях схемы `data` создается новый тип события с инкрементированной версией.
*   **Топики Kafka:**
    *   Именование топиков: `{service_name}.{resource_name}.{action_name}` (например, `auth.user.registered`).
    *   Партиционирование: По ключу, релевантному для сохранения порядка обработки (часто `subject` или `user_id`, `game_id`).
    *   Репликация и Retention: Настраиваются согласно требованиям к надежности и хранению данных.
*   **Обработка событий:** Обработчики должны быть идемпотентными. Порядок событий важен для некоторых сценариев и должен обеспечиваться на уровне Kafka (партиционированием по ключу) и логикой подписчиков.

### 5.4. gRPC API Контракты

Для высокопроизводительного межсервисного взаимодействия используется gRPC.

*   **Protobuf Определения:** Контракты сервисов определяются с использованием Protocol Buffers (версия 3). Файлы `.proto` являются источником истины для структуры сообщений и сигнатур сервисов.
*   **Именование:**
    *   **Пакеты:** Используется для версионирования, например, `platform.service_name.v1` (e.g., `platform.user.v1`).
    *   **Сервисы:** `PascalCaseService` (e.g., `UserService`).
    *   **Методы RPC:** `PascalCaseAction` (e.g., `GetUser`, `CreateGame`).
    *   **Сообщения:** `PascalCaseNoun` (e.g., `UserResponse`, `CreateGameRequest`). Поля сообщений именуются в `snake_case`.
    *   **Перечисления (Enums):** `PascalCaseEnum` для типа, `ENUM_NAME_UPPER_SNAKE_CASE` для значений (e.g., `UserStatus`, `USER_STATUS_ACTIVE`).
*   **Стандартные паттерны сообщений:** Для RPC-методов используются стандартные суффиксы `Request` и `Response` для сообщений запроса и ответа соответственно.
*   **Обработка ошибок:** Используются стандартные коды состояния gRPC (например, `NOT_FOUND`, `INVALID_ARGUMENT`, `PERMISSION_DENIED`). Для передачи дополнительной информации об ошибке используется `google.rpc.Status` и поле `details`.
*   **Общие типы данных Protobuf:** Рекомендуется использовать стандартные типы из `google/protobuf/` для общих нужд, такие как `google.protobuf.Timestamp` для дат и времени, `google.protobuf.Empty` для пустых запросов/ответов, `google.protobuf.Wrappers` для опциональных скалярных типов.
*   **Безопасность:** Соединения должны быть защищены с использованием TLS. Аутентификация сервисов может осуществляться через mTLS или передачу токенов в метаданных запроса.
*   **Документация:** Комментарии в `.proto` файлах должны использоваться для документирования сервисов, методов, сообщений и полей. Эти комментарии могут быть использованы для автоматической генерации документации.
