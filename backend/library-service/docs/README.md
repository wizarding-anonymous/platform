# Спецификация Микросервиса: Library Service

**Версия:** 1.1 (на основе исходной)
**Дата последнего обновления:** 2025-05-22

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   **Назначение сервиса:** Library Service предназначен для управления личными библиотеками пользователей, включая доступ к приобретенным играм и приложениям, отслеживание игрового времени, управление достижениями и синхронизацию сохранений.
*   **Роль в общей архитектуре платформы:** Центральный компонент для пользовательского контента и игровой активности.
*   **Основные бизнес-задачи:** Обеспечение управления библиотеками, отслеживание игрового времени и достижений, синхронизация сохранений, предоставление информации о доступности игр.
*   (Источник: Спецификация Library Service v1.1, раздел 2.1)

### 1.2. Ключевые Функциональности
*   **Управление библиотекой пользователя:** Добавление/скрытие игр, фильтрация, сортировка, пользовательские категории, отслеживание статуса установки, семейный доступ.
*   **Отслеживание игрового времени:** Регистрация начала/окончания сессий, подсчет общего времени, статистика, последняя активность, синхронизация.
*   **Управление достижениями:** Регистрация полученных достижений, хранение прогресса, отображение списка, синхронизация, валидация.
*   **Управление списком желаемого:** Добавление/удаление игр, приоритизация, уведомления о скидках, автоудаление после покупки.
*   **Синхронизация сохранений:** Загрузка/скачивание сохранений в/из облака (S3), разрешение конфликтов, версионирование, автоматическая/ручная синхронизация, управление квотами.
*   **Настройки игр:** Хранение и синхронизация пользовательских настроек игр.
*   (Источник: Спецификация Library Service v1.1, раздел 2.3)

### 1.3. Основные Технологии
*   **Backend:** Go 1.21+, Echo (REST), gRPC.
*   **Базы данных:** PostgreSQL 15+ (основное хранилище), Redis 7+ (кэш, оперативные данные).
*   **Сообщения:** Apache Kafka.
*   **Облачное хранилище:** S3-совместимое.
*   **Инфраструктура:** Docker, Kubernetes, Helm, Prometheus, Grafana, ELK/Loki, Jaeger.
*   (Источник: Спецификация Library Service v1.1, раздел 9.1)

### 1.4. Термины и Определения (Glossary)
*   **Библиотека (Library):** Коллекция игр и продуктов пользователя.
*   **Игровое время (Playtime):** Время, проведенное в игре.
*   **Достижение (Achievement):** Награда за выполнение целей в игре.
*   **Список желаемого (Wishlist):** Список игр для будущей покупки.
*   **Сохранение (Savegame):** Файл с прогрессом пользователя в игре.
*   **Семейный доступ (Family Sharing):** Функция для обмена играми.
*   (Полный глоссарий см. в Спецификации Library Service v1.1, раздел 1)

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   Library Service построен на многослойной архитектуре, используя принципы чистой архитектуры (Clean Architecture).
*   Это обеспечивает разделение ответственности, масштабируемость и тестируемость.
*   (Источник: Спецификация Library Service v1.1, раздел 3.1)

### 2.2. Слои Сервиса
(На основе Спецификации Library Service v1.1, раздел 3.1.1)

#### 2.2.1. Presentation Layer (Слой Представления / Транспортный слой)
*   Ответственность: Обработка входящих REST (Echo), gRPC запросов и событий Kafka. WebSocket (через API Gateway) для real-time уведомлений.
*   Ключевые компоненты/модули: REST контроллеры, gRPC серверы, обработчики Kafka.

#### 2.2.2. Application Layer (Прикладной Слой / Сервисный слой / Use Case Layer)
*   Ответственность: Реализация бизнес-логики и сценариев использования. Оркестрация взаимодействия компонентов домена и репозиториев.
*   Ключевые компоненты/модули: Сервисы для управления библиотекой, игровым временем, достижениями, списком желаемого, сохранениями.

#### 2.2.3. Domain Layer (Доменный Слой)
*   Ответственность: Определение основных бизнес-сущностей (UserLibraryItem, PlaytimeSession, UserAchievement, WishlistItem, SavegameMetadata), их состояний и поведения.
*   Ключевые компоненты/модули: Entities, Value Objects, Domain Services, Repository Interfaces.

#### 2.2.4. Infrastructure Layer (Инфраструктурный Слой)
*   Ответственность: Реализация интерфейсов репозиториев (PostgreSQL, Redis), взаимодействие с облачным хранилищем (S3), Kafka, клиенты для других микросервисов. Логирование, мониторинг, трассировка.
*   Ключевые компоненты/модули: PostgreSQL repositories, Redis cache implementations, S3 client, Kafka producers/consumers, gRPC clients.

*   (Подробную структуру проекта см. в Спецификации Library Service v1.1, раздел 3.1.2)

## 3. API Endpoints

### 3.1. REST API
*   **Базовый URL:** `/api/v1/library`
*   **Аутентификация:** JWT Bearer Token.
*   **Спецификация:** `/api/openapi/v1/library.yaml` (OpenAPI 3.0).
*   **Основные эндпоинты:**
    *   Библиотека: `GET /`, `GET /{game_id}`, `PATCH /{game_id}`.
    *   Игровое время: `POST /playtime/start`, `POST /playtime/heartbeat`, `POST /playtime/end`, `GET /playtime/stats`.
    *   Достижения: `GET /achievements`, `POST /achievements/unlock`.
    *   Список желаемого: `GET /wishlist`, `POST /wishlist`, `DELETE /wishlist/{game_id}`.
    *   Сохранения: `GET /savegames`, `POST /savegames/upload-url`, `POST /savegames/confirm-upload`, `POST /savegames/download-url`.
    *   Настройки игр: `GET /settings/{game_id}`, `PUT /settings/{game_id}`.
*   (Детали и примеры см. в Спецификации Library Service v1.1, разделы 5.2 и "Детальные примеры API Library Service")

### 3.2. gRPC API
*   **Назначение:** Для синхронного межсервисного взаимодействия.
*   **Спецификация:** `/api/proto/v1/library.proto`.
*   **Сервисы:** `LibraryQueryService` (например, `GetLibrary`, `CheckGameAccess`), `LibraryCommandService` (если необходимо).
*   (Детали и примеры см. в Спецификации Library Service v1.1, раздел 5.3 и "Детальные примеры API Library Service")

### 3.3. WebSocket API
*   **Назначение:** Уведомление Flutter-клиента об изменениях в реальном времени (через API Gateway).
*   **События:** `library.updated`, `playtime.updated`, `achievement.unlocked`, `wishlist.updated`, `savegame.sync.status`.
*   (Детали см. в Спецификации Library Service v1.1, раздел 5.4 и "Детальные примеры API Library Service")

## 4. Модели Данных (Data Models)

### 4.1. Основные Сущности
*   **UserLibraryItem**: Игра/продукт в библиотеке пользователя.
*   **PlaytimeSession**: Сессия игрового времени.
*   **UserAchievement**: Достижение пользователя.
*   **WishlistItem**: Элемент списка желаемого.
*   **SavegameMetadata**: Метаданные сохранения игры.
*   **UserGameCategory**: Пользовательская категория для игр.
*   **FamilySharingLink**: Связь для семейного доступа.
*   (Структуры JSON/Go см. в Спецификации Library Service v1.1, раздел 5.1)

### 4.2. Схема Базы Данных
*   **PostgreSQL**: Хранит `user_libraries`, `playtime_sessions`, `playtime_stats_daily/monthly`, `user_achievements`, `user_achievement_progress`, `achievement_stats`, `wishlists`, `wishlist_price_history`, `savegame_metadata`, `savegame_sync_settings`, `savegame_versions`, `game_settings`, `game_settings_history`, `library_audit_log`, `achievement_audit_log`.
    ```sql
    -- Пример таблицы user_libraries (сокращенно)
    CREATE TABLE user_libraries (id UUID PRIMARY KEY, user_id UUID NOT NULL, game_id UUID NOT NULL, acquisition_date TIMESTAMPTZ ...);
    -- Пример таблицы playtime_sessions (сокращенно)
    CREATE TABLE playtime_sessions (id UUID PRIMARY KEY, user_id UUID NOT NULL, game_id UUID NOT NULL, start_time TIMESTAMPTZ ...);
    ```
*   **Redis**: Кэш библиотек, текущие игровые сессии, временные токены S3, счетчики.
*   **S3-совместимое хранилище**: Файлы сохранений.
*   (Полные DDL и детальное описание см. в Спецификации Library Service v1.1, раздел "Схемы таблиц базы данных Library Service" и 3.3)

## 5. Потоковая Обработка Событий (Event Streaming)

### 5.1. Публикуемые События (Produced Events)
*   **Система сообщений:** Apache Kafka.
*   **Формат событий:** CloudEvents JSON.
*   **Топики:** `library-events`, `achievement-events`, `playtime-events`, `wishlist-events`, `savegame-events`.
*   **Основные публикуемые события:**
    *   `library.game.added`
    *   `library.achievement.unlocked`
    *   `library.playtime.updated`
    *   `library.wishlist.item.added` / `removed`
    *   `library.savegame.synchronized`
*   (Детали и примеры Payload см. в Спецификации Library Service v1.1, раздел 11.3 и "Kafka События")

### 5.2. Потребляемые События (Consumed Events)
*   **Основные потребляемые события:**
    *   `payment.purchase.completed` (от Payment Service): Добавить игру в библиотеку.
    *   `catalog.game.updated` (от Catalog Service): Обновить кэш метаданных.
    *   `download.game.installed` / `uninstalled` (от Download Service): Обновить статус установки.
    *   `account.user.deleted` (от Account Service): Анонимизировать/удалить данные.
*   (Детали см. в Спецификации Library Service v1.1, раздел 11.4)

## 6. Интеграции (Integrations)

### 6.1. Внутренние Микросервисы
*   **Account Service**: Получение данных пользователя (gRPC).
*   **Catalog Service**: Получение метаданных игр и достижений (gRPC).
*   **Payment Service**: Получение событий о покупках (Kafka).
*   **Download Service**: Получение статуса установки игр (Kafka), проверка прав доступа (gRPC).
*   **Social Service**: Уведомления об игровых событиях (Kafka).
*   **Analytics Service**: Отправка данных для анализа (Kafka).
*   **Notification Service**: Отправка событий для уведомлений (Kafka).
*   **API Gateway**: Обработка клиентских запросов (REST, WebSocket).
*   (Детали см. в Спецификации Library Service v1.1, раздел 6.1)

### 6.2. Внешние Системы
*   **Облачное хранилище (S3-совместимое)**: Хранение файлов сохранений (AWS S3 SDK Go).
*   (Детали см. в Спецификации Library Service v1.1, раздел 6.2)

## 7. Конфигурация (Configuration)

### 7.1. Переменные Окружения
*   `LIBRARY_HTTP_PORT`, `LIBRARY_GRPC_PORT`.
*   `POSTGRES_DSN`.
*   `REDIS_ADDR`.
*   `KAFKA_BROKERS`.
*   `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`, `S3_BUCKET_SAVEGAMES`.
*   Адреса gRPC других сервисов.
*   `LOG_LEVEL`.
*   Квоты на облачное хранилище.
*   Таймауты, лимиты.
*   (Полный список см. в `configs/config.yaml` и разделе 9.2 исходной спецификации).

### 7.2. Файлы Конфигурации (если применимо)
*   `configs/config.yaml` (с возможностью переопределения через переменные окружения).
*   (Источник: Спецификация Library Service v1.1, раздел 3.1.2, 9.2)

## 8. Обработка Ошибок (Error Handling)

### 8.1. Общие Принципы
*   REST API: Стандартные коды HTTP, JSON-формат ошибки.
*   gRPC API: Стандартные коды gRPC.
*   Использование таймаутов, retry с экспоненциальной задержкой, Circuit Breaker.
*   Идемпотентность обработчиков событий Kafka, использование DLQ.
*   (Детали см. в Спецификации Library Service v1.1, раздел 6.3)

### 8.2. Распространенные Коды Ошибок
*   `GAME_NOT_IN_LIBRARY` (HTTP 404)
*   `SAVEGAME_CONFLICT` (HTTP 409)
*   `QUOTA_EXCEEDED` (HTTP 403/422)
*   `ACHIEVEMENT_ALREADY_UNLOCKED` (HTTP 409)
*   `INVALID_SESSION_OPERATION` (HTTP 400)

## 9. Безопасность (Security)

### 9.1. Аутентификация
*   JWT (через Auth Service) для всех API.

### 9.2. Авторизация
*   RBAC. Проверка владения ресурсами (библиотека, сохранения и т.д.).
*   (Матрица доступа см. в Спецификации Library Service v1.1, раздел 10.2).

### 9.3. Защита Данных
*   TLS 1.2+ для всех коммуникаций.
*   Шифрование секретов. Рассмотреть шифрование сохранений в S3.
*   Валидация входных данных. Защита от подделки достижений/игрового времени.

### 9.4. Управление Секретами
*   Kubernetes Secrets или HashiCorp Vault.
*   (Детали см. в Спецификации Library Service v1.1, раздел 7.1).

## 10. Развертывание (Deployment)

### 10.1. Инфраструктурные Файлы
*   **Dockerfile:** Многоэтапный (golang:alpine -> alpine).
*   **Kubernetes манифесты/Helm-чарты.**
*   (Детали см. в Спецификации Library Service v1.1, раздел 9.2).

### 10.2. Зависимости при Развертывании
*   PostgreSQL, Redis, Kafka, S3-хранилище.
*   Account, Catalog, Auth, Notification, Download, API Gateway.

### 10.3. CI/CD
*   GitLab CI/CD. Сборка, тесты (unit, integration), публикация образа, развертывание (dev, staging, prod).
*   (Детали см. в Спецификации Library Service v1.1, раздел 9.3).
*   **Миграции БД:** golang-migrate, применяются автоматически или через CI/CD. (см. раздел 9.4).
*   **Резервное копирование:** PostgreSQL (PITR, pg_dump), Redis (RDB/AOF), S3 (версионирование, репликация). (см. раздел 9.5).

## 11. Мониторинг и Логирование (Logging and Monitoring)

### 11.1. Логирование
*   **Формат:** Структурированные JSON логи (Zap).
*   **Интеграция:** Fluentd/Fluent Bit -> ELK/Loki.
*   (Детали см. в Спецификации Library Service v1.1, раздел 8.2).

### 11.2. Мониторинг
*   **Метрики (Prometheus):** Запросы (количество, задержка, ошибки), ресурсы (CPU, RAM), БД, Kafka, бизнес-метрики (активные сессии, синхронизации).
*   **Дашборды (Grafana):** Стандартные и кастомные.
*   **Алертинг (AlertManager):** По пороговым значениям метрик.
*   (Детали см. в Спецификации Library Service v1.1, раздел 8.1).

### 11.3. Трассировка
*   **Инструментация:** OpenTelemetry SDK.
*   **Экспорт:** Jaeger.
*   Сквозная трассировка запросов.
*   (Детали см. в Спецификации Library Service v1.1, раздел 8.3).

## 12. Нефункциональные Требования (NFRs)
*   **Производительность:** REST API P95 < 100мс, P99 < 200мс; gRPC P95 < 50мс. RPS >= 5000.
*   **Масштабируемость:** Горизонтальная. Поддержка >= 10 млн. пользователей, >= 100 млн. сессий/день.
*   **Надежность:** Доступность 99.99%. RTO < 15 мин, RPO < 5 мин.
*   **Безопасность:** Защита ПД, шифрование, валидация.
*   (Полный список см. в Спецификации Library Service v1.1, раздел 2.4).

## 13. Приложения (Appendices) (Опционально)
*   Детальные примеры API запросов/ответов, схемы Protobuf, DDL.
*   (Многие из этих деталей содержатся в Спецификации Library Service v1.1, включая "Детальные примеры API Library Service" и "Схемы таблиц базы данных Library Service").

---
*Этот шаблон является отправной точкой и может быть адаптирован под конкретные нужды проекта и сервиса.*
