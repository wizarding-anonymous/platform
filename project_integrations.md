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

**TODO:** Разработать и задокументировать детальные стандарты для каждого типа контракта. Ссылки на эти стандарты должны быть включены во все спецификации микросервисов.
