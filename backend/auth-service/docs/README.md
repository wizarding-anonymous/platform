# Спецификация Микросервиса: Auth Service (Микросервис Аутентификации)

**Версия:** 2.0
**Дата последнего обновления:** 2023-10-27

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение и Роль:** Auth Service является центральным компонентом платформы, отвечающим за управление идентификацией пользователей, аутентификацию и авторизацию. Он обеспечивает безопасный, масштабируемый и надежный фундамент для модели безопасности всей платформы. Его основная цель - проверка подлинности пользователей, контроль доступа к ресурсам на основе ролей и разрешений, управление сессиями и токенами (JWT), а также предоставление доверенной информации об аутентифицированных пользователях другим микросервисам.
*   (Источник: `New Docs/specification/auth_service_overview.md`, разделы 1, 2)

### 1.2. Ключевые Функциональности
*   Регистрация пользователей (в координации с Account Service).
*   Аутентификация по логину/паролю.
*   Двухфакторная аутентификация (2FA): TOTP, SMS/Email коды, резервные коды.
*   Внешняя аутентификация: OAuth 2.0 / OIDC (Telegram, VK, Odnoklassniki).
*   Управление JSON Web Token (JWT): Генерация (Access Token RS256, Refresh Token), валидация, ротация.
*   Управление сессиями: Отслеживание и отзыв активных сессий.
*   Управление паролями: Сброс и изменение пароля.
*   Подтверждение Email.
*   Role-Based Access Control (RBAC): Управление ролями и разрешениями пользователей.
*   Управление API ключами.
*   Аудит событий безопасности.
*   Обнаружение подозрительной активности (базовое).
*   Административные функции (управление пользователями, ролями через Admin Service).
*   (Источник: `New Docs/specification/auth_service_overview.md`, раздел 3)

### 1.3. Основные Технологии
*   **Язык программирования:** Go.
*   **API Фреймворки:** Gin (REST), standard `net/http` и `grpc-go` (gRPC).
*   **Базы данных:** PostgreSQL (пользовательские данные, роли, токены), Redis (кэш сессий, временные токены, rate limiting).
*   **Сообщения:** Kafka (`confluent-kafka-go`).
*   **Безопасность:** `golang-jwt/jwt/v5` (JWT), `golang.org/x/crypto/argon2` (хеширование паролей).
*   **Логирование:** Zap.
*   **Мониторинг:** Prometheus (`prometheus-client_golang`).
*   **Трассировка:** OpenTelemetry.
*   (Источник: `New Docs/specification/auth_service_overview.md`, раздел 5; `Auth_Service_Detailed_Specification.md`, раздел 10)

### 1.4. Термины и Определения (Glossary)
*   Для используемых терминов (Access Token, Refresh Token, JWT, RBAC, 2FA, TOTP, Argon2id, OAuth, OIDC, JWKS и др.) см. "Единый глоссарий терминов и определений для российского аналога Steam.txt".
*   Дополнительные определения можно найти в `project_docs/old_docs/Спецификация микросервиса Auth Service.md`, раздел 2.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Auth Service разработан как stateless микросервис на языке Go.
*   Он взаимодействует с PostgreSQL для постоянного хранения данных, Redis для кэширования и временных данных, и Kafka для асинхронного обмена событиями.
*   Архитектура ориентирована на безопасность, масштабируемость и надежность.
*   (Источник: `New Docs/specification/auth_service_overview.md`, раздел 4)

Сервис Auth Service спроектирован как stateless (не хранящий состояние между запросами) микросервис, написанный на языке Go. В основе его архитектуры лежит многослойный подход, включающий:
*   **API Layer (Слой Представления):** Отвечает за обработку всех входящих запросов (REST через Gin, gRPC для межсервисного взаимодействия), валидацию данных и вызов соответствующей бизнес-логики.
*   **Application Layer (Прикладной Слой / Слой Бизнес-Логики):** Содержит основную логику сервиса, включая процессы регистрации, аутентификации (логин/пароль, JWT, 2FA, OAuth, API ключи), авторизации (RBAC), управления сессиями и токенами. Оркестрирует взаимодействие между доменными сущностями и инфраструктурным слоем.
*   **Infrastructure Layer (Инфраструктурный Слой / Слой Доступа к Данным):** Реализует взаимодействие с внешними системами и хранилищами: PostgreSQL для персистентного хранения данных (пользователи, роли, токены и т.д.), Redis для кэширования сессий и временных данных (например, кодов 2FA, счетчиков ограничения скорости), и Apache Kafka для асинхронной публикации доменных событий (например, `auth.user.registered`).

Такой подход обеспечивает четкое разделение ответственности (separation of concerns) между компонентами, повышает тестируемость каждого слоя в изоляции, облегчает сопровождение и дальнейшее развитие сервиса.

**Диаграмма верхнеуровневой архитектуры:**
```mermaid
graph TD
    subgraph Auth Service
        direction LR
        subgraph PresentationLayer [API Layer (Presentation)]
            direction TB
            API_REST[REST API (Gin)]
            API_GRPC[gRPC API]
        end

        subgraph ApplicationLayer [Business Logic Layer (Application)]
            direction TB
            RegSvc[Registration Service]
            LoginSvc[Login Service]
            TokenSvc[Token Management (JWT, Refresh)]
            TwoFASvc[2FA Service (TOTP, SMS/Email)]
            RBACSvc[RBAC Service (Roles, Permissions)]
            ApiKeySvc[API Key Management]
            ExternalAuthSvc[External Auth (OAuth, Telegram)]
            PasswordSvc[Password Management (Reset, Change)]
            EmailVerifySvc[Email Verification Service]
        end

        subgraph InfrastructureLayer [Data Access & Infrastructure Layer]
            direction TB
            RepoPostgreSQL[PostgreSQL Repositories (Users, Roles, Tokens, Sessions, etc.)]
            CacheRedis[Redis Client (Session Cache, 2FA codes, Rate Limits)]
            ProducerKafka[Kafka Producer (Domain Events)]
            CryptoUtils[Cryptography (Argon2id, JWT Signing)]
        end

        PresentationLayer --> ApplicationLayer
        ApplicationLayer --> InfrastructureLayer
    end

    Clients[Clients (Web, Mobile, Desktop)] --> APIGateway[API Gateway]
    APIGateway -- REST/gRPC Requests --> PresentationLayer

    InfrastructureLayer -- CRUD Operations --> DB[(Database: PostgreSQL)]
    InfrastructureLayer -- Cache Operations --> Cache[(Cache: Redis)]

    ApplicationLayer -- Publishes Domain Events --> KafkaBroker[Kafka Broker]
    KafkaBroker -- auth.user.registered etc. --> AccountSvc[Account Service]
    KafkaBroker -- auth.user.verification_code_sent etc. --> NotificationSvc[Notification Service]

    InternalMS[Other Microservices] -- gRPC: ValidateToken, CheckPermission --> API_GRPC
```

### 2.2. Слои Сервиса
(Предполагаемая структура на основе функционала и технологий, детализирующая диаграмму выше)

#### 2.2.1. Presentation Layer (Слой Представления)
*   Ответственность: Обработка входящих REST и gRPC запросов, валидация, вызов Application Layer.
*   Ключевые компоненты/модули:
    *   HTTP Handlers (Gin): Эндпоинты для регистрации, логина, управления токенами, 2FA, API ключами и т.д. (см. `auth_api_rest.md`).
    *   gRPC Service Implementations: Реализация методов `ValidateToken`, `CheckPermission`, `GetUserInfo`, `GetJWKS` (см. `auth_api_grpc.md`).
    *   DTOs: Для запросов/ответов API.

#### 2.2.2. Application Layer (Прикладной Слой / Слой Сценариев Использования)
*   Ответственность: Реализация сценариев использования, таких как регистрация, аутентификация, управление токенами, 2FA, внешняя аутентификация, управление API ключами, RBAC.
*   Ключевые компоненты/модули:
    *   Сервисы для каждого рабочего процесса, описанного в `auth_service_logic.md` (например, `UserRegistrationService`, `LoginService`, `TokenManagementService`, `TwoFactorAuthService`).
    *   Интерфейсы для репозиториев и внешних клиентов.

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Бизнес-сущности и правила, связанные с аутентификацией и авторизацией.
*   Ключевые компоненты/модули:
    *   Entities: `User` (учетные данные, статус), `Role`, `Permission`, `Session`, `RefreshToken`, `ExternalAccount`, `MFASecret`, `MFABackupCode`, `APIKey`, `VerificationCode`, `AuditLog`.
    *   Value Objects: Например, для представления токенов, хешей паролей.
    *   Domain Events: `auth.user.registered`, `auth.user.login_success` и т.д. (см. `auth_event_streaming.md`).
    *   Интерфейсы репозиториев.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Взаимодействие с PostgreSQL, Redis, Kafka. Реализация криптографических операций.
*   Ключевые компоненты/модули:
    *   PostgreSQL Repositories: Для `users`, `roles`, `permissions`, `sessions`, `refresh_tokens` и др. (см. `auth_data_model.md`).
    *   Redis Client: Для кэша сессий, счетчиков rate limiting, временных токенов.
    *   Kafka Producers/Consumers: Для публикации и подписки на события (см. `auth_event_streaming.md`).
    *   Password Hashing (Argon2id).
    *   JWT Generation/Validation (RS256).
    *   TOTP Generation/Validation.

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL:** `/api/v1/auth` (маршрутизируется через API Gateway).
*   **Формат данных:** JSON.
*   **Аутентификация:** Большинство эндпоинтов требуют `Bearer <access_token>`, публичные эндпоинты (регистрация, логин и т.д.) не требуют.
*   **Основные группы эндпоинтов:**
    *   Регистрация и Вход (`POST /register`, `POST /login`).
    *   Управление токенами (`POST /refresh-token`, `POST /logout`, `POST /logout-all`).
    *   Верификация Email (`POST /verify-email`, `POST /resend-verification`).
    *   Управление Паролем (`POST /forgot-password`, `POST /reset-password`, `PUT /me/password`).
    *   Управление Двухфакторной Аутентификацией (`POST /me/2fa/totp/enable`, `POST /me/2fa/totp/verify`, `POST /login/2fa/verify`, `POST /me/2fa/disable`).
    *   Управление Текущим Пользователем (`GET /me`, `GET /me/sessions`, `DELETE /me/sessions/{session_id}`).
    *   Внешняя Аутентификация (`GET /oauth/{provider}`, `GET /oauth/{provider}/callback`, `POST /telegram-login`).
    *   Управление API Ключами (`GET /me/api-keys`, `POST /me/api-keys`, `DELETE /me/api-keys/{key_id}`).
    *   Административные Эндпоинты (`GET /admin/users`, `POST /admin/users/{user_id}/block` и др.).
*   (Полный список и детали см. в `New Docs/specification/auth_api_rest.md`).

### 3.2. gRPC API
*   Предназначен для внутреннего взаимодействия между микросервисами.
*   Определение Protobuf: `auth.v1.proto`.
*   **Основные RPC Методы:**
    *   `ValidateToken`: Проверка JWT токена.
    *   `CheckPermission`: Проверка прав доступа пользователя.
    *   `GetUserInfo`: Получение информации о пользователе.
    *   `GetJWKS`: Получение публичных ключей для верификации JWT.
    *   `HealthCheck`: Стандартная проверка работоспособности.
*   (Полный список и детали см. в `New Docs/specification/auth_api_grpc.md`).

### 3.3. WebSocket API (если применимо)
*   Information not found in existing documentation. (Не используется Auth Service напрямую для предоставления API).

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности
*   `User`: Учетные данные пользователя (username, email, хеш пароля, соль, статус).
*   `Role`: Роли в системе (user, admin, developer).
*   `Permission`: Гранулярные разрешения.
*   `Session`: Активные сессии пользователей.
*   `RefreshToken`: Токены обновления.
*   `ExternalAccount`: Связь с внешними OAuth провайдерами.
*   `MFASecret`: Секреты для 2FA (TOTP).
*   `MFABackupCode`: Резервные коды 2FA.
*   `APIKey`: API ключи для разработчиков/сервисов.
*   `AuditLog`: Журнал событий безопасности.
*   `VerificationCode`: Временные коды для email верификации, сброса пароля.

### 4.2. Схема Базы Данных
*   **PostgreSQL:** Используется для хранения `users`, `roles`, `permissions`, `role_permissions`, `user_roles`, `sessions`, `refresh_tokens`, `external_accounts`, `mfa_secrets`, `mfa_backup_codes`, `api_keys`, `audit_logs`, `verification_codes`.
    ```sql
    -- Пример таблицы users (сокращенно)
    CREATE TABLE users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        username VARCHAR(255) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password_hash VARCHAR(255) NOT NULL,
        salt VARCHAR(128) NOT NULL,
        status VARCHAR(50) NOT NULL DEFAULT 'pending_verification',
        -- ...
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    ```
*   **Redis:** Используется для кэширования сессий (опционально), счетчиков rate limiting, временных токенов (2FA, password reset), списка отозванных токенов (JTI blacklist).
*   (Полную схему PostgreSQL и описание структур Redis см. в `New Docs/specification/auth_data_model.md`).

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Kafka.
*   **Формат событий:** CloudEvents v1.0 JSON.
*   **Топик:** `auth.events`.
*   **Основные публикуемые события:**
    *   `auth.user.registered.v1`
    *   `auth.user.email_verified.v1`
    *   `auth.user.password_reset_requested.v1`
    *   `auth.user.password_changed.v1`
    *   `auth.user.login_success.v1`
    *   `auth.user.login_failed.v1`
    *   `auth.user.account_locked.v1`
    *   `auth.user.roles_changed.v1`
    *   `auth.session.created.v1`
    *   `auth.session.revoked.v1`
    *   `auth.2fa.enabled.v1`
    *   `auth.2fa.disabled.v1`
*   (Детальную структуру Payload для каждого события см. в `New Docs/specification/auth_event_streaming.md`).

### 5.2. Потребляемые События (Consumed Events)
*   **Основные потребляемые события:**
    *   `account.user.profile_updated.v1` (от Account Service): Для обновления кэшированной информации или реакции на изменение статуса.
    *   `admin.user.force_logout.v1` (от Admin Service): Для принудительного завершения сессий пользователя.
    *   `admin.user.block.v1` (от Admin Service): Для блокировки пользователя и отзыва сессий.
    *   `admin.user.unblock.v1` (от Admin Service): Для разблокировки пользователя.
*   (Детальную структуру Payload и логику обработки см. в `New Docs/specification/auth_event_streaming.md`).

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **API Gateway**: Валидация токенов (gRPC `ValidateToken`), проксирование публичных REST эндпоинтов Auth Service.
*   **Account Service**: Публикация события `auth.user.registered` для создания профиля; получение статуса аккаунта от Account Service; потребление событий `account.user.profile_updated` для реакции на изменения.
*   **Notification Service**: Публикация событий для отправки email (верификация, сброс пароля) и SMS (2FA).
*   **Admin Service**: Предоставление REST/gRPC API для управления пользователями, ролями, просмотра аудита; потребление событий от Admin Service (force_logout, block/unblock).
*   **Другие микросервисы**: Предоставление gRPC API (`ValidateToken`, `CheckPermission`) для авторизации их запросов.
*   (Детали см. в `New Docs/specification/auth_integrations.md`).

### 6.2. Внешние Системы
*   **OAuth Провайдеры (Telegram, VK, Odnoklassniki)**: Интеграция для внешней аутентификации.
*   (Детали см. в `New Docs/specification/auth_service_logic.md`, раздел про External Authentication).

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `AUTH_DB_HOST, AUTH_DB_PORT, AUTH_DB_USER, AUTH_DB_PASSWORD, AUTH_DB_NAME`: Параметры подключения к PostgreSQL.
*   `AUTH_REDIS_ADDR, AUTH_REDIS_PASSWORD`: Параметры Redis.
*   `AUTH_KAFKA_BROKERS`: Адреса Kafka.
*   `AUTH_JWT_PRIVATE_KEY_PATH`, `AUTH_JWT_PUBLIC_KEY_PATH`: Пути к RSA ключам для JWT.
*   `AUTH_JWT_ACCESS_TOKEN_TTL`, `AUTH_JWT_REFRESH_TOKEN_TTL`: Сроки жизни токенов.
*   `AUTH_ARGON2ID_MEMORY`, `AUTH_ARGON2ID_ITERATIONS`, `AUTH_ARGON2ID_PARALLELISM`, `AUTH_ARGON2ID_KEY_LENGTH`, `AUTH_ARGON2ID_SALT_LENGTH`: Параметры Argon2id.
*   `LOG_LEVEL`.
*   OAuth provider client IDs and secrets.
*   (Более полный список см. в старой спецификации `project_docs/old_docs/Спецификация микросервиса Auth Service.md`, раздел 8.3, и в `Auth_Service_Detailed_Specification.md`, раздел 10).

### 7.2. Файлы Конфигурации (если применимо)
*   `config.yaml` может использоваться для задания структуры настроек, которые затем переопределяются переменными окружения.
*   Хранение JWKS ключей.
*   (Пример см. в старой спецификации, раздел 8.3 и `New Docs/specification/Auth_Service_Detailed_Specification.md`, раздел 10).

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   REST API: Стандартные коды HTTP, JSON-формат ошибки с полями `code`, `title`, `detail`, `source`, `meta` (см. `auth_api_rest.md`).
*   gRPC API: Стандартные коды gRPC, детализация через `google.rpc.Status` и `google.rpc.ErrorInfo`.
*   Логирование всех ошибок с `request_id`/`trace_id`.

### 8.2. Распространенные Коды Ошибок
*   `invalid_credentials` (HTTP 401)
*   `validation_error` (HTTP 400)
*   `username_already_exists` (HTTP 409)
*   `email_already_exists` (HTTP 409)
*   `token_expired` (HTTP 401)
*   `invalid_token` (HTTP 401)
*   `user_blocked` (HTTP 403)
*   `email_not_verified` (HTTP 403)
*   `invalid_2fa_code` (HTTP 401)
*   (Полный список см. в `auth_api_rest.md` и старой спецификации раздел 5.5).

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   **Механизмы:** Логин/пароль, JWT (Access + Refresh), OAuth 2.0/OIDC (Telegram, VK, OK), API ключи.
*   **JWT:** RS256, Access Token (15 мин), Refresh Token (30 дней, ротация, HttpOnly cookie для web).
*   **2FA:** TOTP (RFC 6238), SMS/Email коды, резервные коды.

### 9.2. Авторизация
*   **RBAC:** Роли и разрешения хранятся в БД. JWT содержит роли/разрешения.
*   Проверка на API Gateway (грубая) и в микросервисах (точная, через Auth Service gRPC `CheckPermission`).

### 9.3. Защита Данных
*   **Хеширование паролей:** Argon2id.
*   **Шифрование:** TOTP секреты шифруются на уровне приложения (AES-256-GCM). TLS для всех коммуникаций.
*   **Защита от атак:** Rate limiting, блокировка аккаунтов, CAPTCHA (план), защита от CSRF/XSS (для админ UI), предотвращение user enumeration.
*   **ФЗ-152:** Локализация данных, согласие, минимизация данных, права субъектов.

### 9.4. Управление Секретами
*   Приватные ключи JWT, секреты OAuth провайдеров, мастер-ключ шифрования хранятся в HashiCorp Vault или зашифрованных Kubernetes Secrets.
*   Ротация ключей.
*   (Детали см. в `New Docs/specification/auth_security_compliance.md`).

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Многоэтапный, на основе Alpine.
*   **Kubernetes манифесты/Helm-чарты:** Deployment, Service, ConfigMap, Secret, HPA.
*   (Детали см. в `Auth_Service_Detailed_Specification.md`, раздел 10 и старой спецификации, раздел 8).

### 10.2. Зависимости при Развертывании
*   PostgreSQL, Redis, Kafka.
*   API Gateway, Account Service, Notification Service, Admin Service.

### 10.3. CI/CD
*   Автоматизированные пайплайны для сборки, тестирования (unit, integration, security scans), развертывания (dev, staging, prod).
*   Стратегии RollingUpdate/Canary.
*   (Детали см. в `Auth_Service_Detailed_Specification.md`, раздел 10 и старой спецификации, раздел 8.4).

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат:** Структурированные логи JSON (Zap).
*   **События:** Запросы, ответы, ошибки, события безопасности, изменения состояния.
*   **Интеграция:** FluentBit/Loki или ELK.
*   (Детали см. в `Auth_Service_Detailed_Specification.md`, раздел 10 и старой спецификации, раздел 9.2).

### 11.2. Мониторинг
*   **Метрики (Prometheus):** `auth_requests_total`, `auth_request_duration_seconds`, `auth_token_validation_total`, `auth_active_sessions_count`, и др.
*   **Дашборды (Grafana):** Обзор производительности, безопасности.
*   **Алертинг (AlertManager):** HighErrorRate, SlowResponseTime, TokenValidationFailures.
*   (Детали см. в `Auth_Service_Detailed_Specification.md`, раздел 10 и старой спецификации, раздел 9.1).

### 11.3. Трассировка
*   **Интеграция:** OpenTelemetry, Jaeger/Tempo.
*   Трассировка входящих/исходящих запросов, взаимодействия с БД/Kafka.
*   (Детали см. в `Auth_Service_Detailed_Specification.md`, раздел 10 и старой спецификации, раздел 9.3).
*   **Аудит безопасности:** Детальное логирование событий в `audit_logs` таблицу (см. `auth_data_model.md` и `auth_security_compliance.md`).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:** Логин/проверка токена: P95 < 150 мс, P99 < 300 мс (1000 RPS). Генерация токена: P95 < 50 мс.
*   **Масштабируемость:** Горизонтальное масштабирование до 5000 RPS.
*   **Надежность:** Доступность 99.98%. RTO < 5 мин, RPO < 1 мин.
*   **Безопасность:** Argon2id, RS256 JWT, OWASP Top 10, ФЗ-152.
*   **Сопровождаемость:** Покрытие тестами > 85%.
*   (Детали см. в `Auth_Service_Detailed_Specification.md`, раздел 9).

## 13. Приложения (Appendices) (Опционально)
*   Детальные примеры API запросов/ответов, схемы Protobuf, JSON Schemas для событий и моделей данных находятся в соответствующих файлах в директории `New Docs/specification/`:
    *   `auth_api_rest.md`
    *   `auth_api_grpc.md`
    *   `auth_data_model.md`
    *   `auth_event_streaming.md`
*   Старая спецификация (`project_docs/old_docs/Спецификация микросервиса Auth Service.md`) также содержит примеры в разделе 11.

---
*Этот документ является обобщением и структурированием информации из детализированных спецификаций Auth Service v2.0. Для полной информации следует обращаться к соответствующим файлам в `New Docs/specification/`.*
