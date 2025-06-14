# Auth Service (Микросервис Аутентификации и Авторизации)

## Обзор (Overview)

Auth Service является центральным компонентом платформы "Российский Аналог Steam", отвечающим за управление идентификацией пользователей, аутентификацию и базовую авторизацию для MVP.

**Подробная спецификация находится здесь: [./docs/README.md](./docs/README.md)** (Эта спецификация описывает полный набор запланированных функций. Текущая MVP реализация включает подмножество этих функций, как описано ниже.)

## Ключевые Функциональности (MVP - Core Functionality)

*   Регистрация пользователей с использованием email и пароля. Статус пользователя `pending_verification` до подтверждения email.
*   Аутентификация пользователей (логин/пароль).
*   Управление сессиями через JWT (Access Token RS256, Refresh Token).
    *   Выдача и валидация токенов.
    *   Механизм обновления Access Token с использованием Refresh Token.
    *   Базовый отзыв токенов при выходе из системы (JTI blacklist для Access Token, удаление Refresh Token из БД).
*   Базовый процесс верификации email (генерация кода, отправка через событие Kafka, проверка кода).
*   Предоставление gRPC эндпоинта для валидации Access Token другими сервисами.

## Основные Технологии (Technologies - MVP)

*   **Язык программирования:** Go (1.21+)
*   **API:** REST (Gin), gRPC
*   **Базы данных:**
    *   PostgreSQL: хранение пользователей, refresh-токенов, кодов верификации.
    *   Redis: JTI blacklist для Access Token.
*   **Сообщения/События:** Apache Kafka для асинхронного взаимодействия (например, для отправки email верификации).
*   **Безопасность:** JWT (RS256), хеширование паролей (Argon2id).

## Getting Started

### Prerequisites
*   Go 1.21+
*   Docker (optional, for containerized build/run)
*   PostgreSQL instance
*   Redis instance
*   Kafka instance
*   OpenSSL (for generating RSA keys)
*   librdkafka: Required for Kafka functionality.
    *   On Alpine (used in Docker): `apk add librdkafka-dev` (for building) and `apk add librdkafka` (for running)
    *   On Debian/Ubuntu: `sudo apt-get install librdkafka-dev`
    *   On macOS: `brew install librdkafka`

### Configuration
1.  **RSA Keys**: Generate RSA private and public keys for signing JWTs:
    ```bash
    ./scripts/generate_keys.sh
    ```
    This will create `jwtRS256.key` and `jwtRS256.key.pub` in the `configs/keys/` directory. Ensure these paths are correctly referenced in your configuration under `rsa_keys`.
2.  **Configuration File**:
    *   The service uses `configs/config.yaml`. Copy `configs/config.yaml.example` to `configs/config.yaml` if you are setting up for the first time, or ensure your existing `config.yaml` is correctly configured.
    *   Pay special attention to database connection details (`postgres`, `redis`), Kafka brokers (`kafka`), and RSA key paths (`rsa_keys`).
    *   Passwords and sensitive information should ideally be set via environment variables (e.g., `AUTH_POSTGRES_PASSWORD`, `AUTH_REDIS_PASSWORD`). The configuration loader supports overriding YAML values with environment variables (e.g., `SERVER_HTTP_PORT` overrides `server.http_port`, `POSTGRES_PASSWORD` overrides `postgres.password`).

### Build
```bash
go build -o auth-service ./cmd/auth-service/main.go
```

### Run
```bash
./auth-service
```
The service will start HTTP and gRPC servers on ports defined in `configs/config.yaml`.

### Docker (Optional)
Build the Docker image:
```bash
docker build -t russian-steam/auth-service .
```
Run the Docker container (ensure `configs/config.yaml` is properly mounted or configuration is provided via environment variables):
```bash
# Example:
# Adjust volume mounts and environment variables as needed for your setup.
# Ensure your config.yaml points to correct DB/Redis/Kafka hostnames accessible from Docker.
docker run -p 8080:8080 -p 50051:50051 \
    -v $(pwd)/configs:/app/configs \
    -e POSTGRES_HOST=docker_postgres_host \
    -e REDIS_HOST=docker_redis_host \
    -e KAFKA_BROKERS=docker_kafka_broker:9092 \
    russian-steam/auth-service
```

## Ключевые Интеграции (Key Integrations - MVP)

*   **API Gateway:** Валидация токенов (через gRPC) и проксирование запросов к REST API.
*   **Account Service:** Потребляет событие `com.platform.auth.user.registered.v1` из Kafka для создания профиля пользователя.
*   **Notification Service:** Потребляет событие `com.platform.auth.email.verification_requested.v1` из Kafka для отправки email с кодом верификации.
*   **Другие микросервисы платформы:** Могут использовать gRPC API (`ValidateToken`) для проверки токенов.

## API Endpoints (MVP)

### REST API (base: /api/v1/auth)
*   `POST /register`: Регистрация нового пользователя.
*   `POST /login`: Аутентификация пользователя.
*   `POST /refresh-token`: Обновление Access Token.
*   `POST /logout`: Выход из системы (инвалидация токенов).
*   `POST /verify-email`: Подтверждение email с использованием кода.

### gRPC API
*   `auth.v1.AuthService/ValidateToken`: Валидация Access Token. (Proto: `proto/auth/v1/auth.proto`)

## Kafka Events Published (MVP)
*   `com.platform.auth.user.registered.v1`
*   `com.platform.auth.email.verification_requested.v1`
*   `com.platform.auth.user.email_verified.v1`
*   `com.platform.auth.user.login_success.v1`
*   `com.platform.auth.session.revoked.v1` (при logout)

```
