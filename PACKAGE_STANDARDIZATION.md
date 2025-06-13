# Package Standardization Guidelines

This document outlines the recommended standard packages and libraries for common tasks across microservices. Adhering to these standards helps in maintaining consistency, reducing boilerplate, improving security, and simplifying dependency management.

## 1. Guiding Principles

*   **Consistency:** Use the same library for the same task across different microservices where feasible.
*   **Community Support & Maturity:** Prefer well-maintained, widely adopted libraries with active communities.
*   **Performance:** Consider the performance implications of a library, especially for critical path operations.
*   **Security:** Choose libraries with a good security track record. Regularly update dependencies to patch vulnerabilities.
*   **Licensing:** Ensure library licenses are compatible with the project's licensing requirements.
*   **Minimalism:** Avoid adding unnecessary dependencies. If a feature can be easily achieved with native language capabilities or an existing library, prefer that.

## 2. Recommended Standard Packages

This section will be populated based on the primary programming language(s) used for microservices. Assuming a primary language like **Go** or **Python** or **Node.js (JavaScript/TypeScript)**, examples are provided. Specify the language and then list packages.

**If the primary language is not yet defined, this section should state that and recommend defining it first.**

---
**(Example for Node.js/TypeScript - Replace/Adapt as needed)**

### 2.1. Node.js / TypeScript
*   **TODO: This section is an example and should be removed or replaced if Node.js is not a primary backend language. The project currently focuses on Go, with potential Java/Kotlin and Python for specific services.**
*   **HTTP Client:**
    *   **`axios`**: Promise-based HTTP client for the browser and node.js. Widely used, feature-rich, and well-documented.
    *   *Alternative:* `node-fetch` (for a more lightweight, standard Fetch API experience).
*   **Logging:**
    *   **`pino`**: Extremely fast, JSON-based logger. Good for structured logging.
    *   *Alternative:* `winston` (more flexible, supports multiple transports).
*   **Web Framework (for REST APIs):**
    *   **`Express.js`**: Minimal and flexible Node.js web application framework. Vast ecosystem.
    *   *Alternative:* `Fastify` (focus on speed and low overhead). Consider if performance is paramount.
*   **gRPC Implementation:**
    *   **`@grpc/grpc-js`** and **`@grpc/proto-loader`**: Official gRPC libraries for Node.js.
*   **Environment Variables Management:**
    *   **`dotenv`**: Loads environment variables from a `.env` file.
*   **Date/Time Manipulation:**
    *   **`date-fns`** or **`luxon`**: Modern libraries for date/time manipulation, offering immutability and I18N support. Avoid Moment.js for new projects due to its mutability and size.
*   **Validation:**
    *   **`joi`** or **`zod`** (especially for TypeScript): Powerful schema description language and data validator.
*   **UUID Generation:**
    *   **`uuid`**: For creating RFC4122 UUIDs.
*   **Testing:**
    *   **`Jest`**: A delightful JavaScript Testing Framework with a focus on simplicity.
    *   *Alternatives:* `Mocha` (flexible, often paired with `Chai` for assertions).
*   **ORM/Database Interaction (PostgreSQL Example):**
    *   **`pg` (node-postgres)**: Non-blocking PostgreSQL client for Node.js.
    *   **`Sequelize`** or **`TypeORM`** (for TypeScript): Promise-based Node.js ORM for Postgres, MySQL, MariaDB, SQLite and Microsoft SQL Server. Choose one and stick to it.
*   **Caching Client (Redis Example):**
    *   **`ioredis`**: A robust, performance-focused and full-featured Redis client for Node.js.
*   **Message Queue Client (Kafka Example):**
    *   **`kafkajs`**: A modern Apache Kafka client for Node.js.

---
**(End of Node.js/TypeScript Example)**

### 2.2. Go
*   **HTTP Framework (REST APIs):**
    *   **`Echo` (`github.com/labstack/echo/v4`)**: High performance, extensible, minimalist Go web framework. (Used consistently across multiple service docs).
    *   *Alternative:* `Gin` (`github.com/gin-gonic/gin`) (Also mentioned, similar performance and features).
*   **gRPC Implementation:**
    *   **`google.golang.org/grpc`**: Official gRPC library for Go.
    *   **`google.golang.org/protobuf/cmd/protoc-gen-go`**: Protobuf compiler plugin for Go.
*   **Logging:**
    *   **`Zap` (`go.uber.org/zap`)**: Blazing fast, structured, leveled logging in Go. (Used consistently).
    *   *Alternative:* `Logrus` (`github.com/sirupsen/logrus`) (Popular, but Zap preferred for performance in many new projects).
*   **Configuration Management:**
    *   **`Viper` (`github.com/spf13/viper`)**: Go configuration with fangs. Supports multiple formats, environment variables, remote config. (Mentioned in service docs).
*   **Database Interaction (PostgreSQL):**
    *   **`GORM` (`gorm.io/gorm`)** with `gorm.io/driver/postgres`: Developer-friendly ORM. (Mentioned as primary or option in several services).
    *   **`pgx` (`github.com/jackc/pgx/v5`)**: Low-level PostgreSQL driver and toolkit. Often used with `sqlx` or directly for performance. (Mentioned as alternative/option).
    *   **`sqlx` (`github.com/jmoiron/sqlx`)**: General purpose SQL extension package.
    *   **Миграции:** `golang-migrate/migrate` (`github.com/golang-migrate/migrate/v4`).
*   **Caching Client (Redis):**
    *   **`go-redis` (`github.com/redis/go-redis/v9`)**: High-performance Redis client for Go.
*   **Message Queue Client (Kafka):**
    *   **`confluent-kafka-go` (`github.com/confluentinc/confluent-kafka-go`)**: Confluent's Go client for Apache Kafka. (Mentioned in service docs).
    *   *Alternative:* `Shopify/sarama` (`github.com/Shopify/sarama`) (Another popular Go client for Kafka).
*   **UUID Generation:**
    *   **`gofrs/uuid` (`github.com/gofrs/uuid`)** or **`google/uuid` (`github.com/google/uuid`)**: For creating RFC 4122 UUIDs.
*   **Validation:**
    *   **`go-playground/validator/v10` (`github.com/go-playground/validator/v10`)**: Package validator implements value validations for structs and individual fields based on tags.
*   **Testing:**
    *   Standard `testing` package.
    *   **`testify` (`github.com/stretchr/testify`)**: Toolkit with `assert`, `require`, and `mock` packages.
    *   **`testcontainers-go` (`github.com/testcontainers/testcontainers-go`)**: For integration testing with real dependencies in Docker containers.
*   **OpenTelemetry (Observability):**
    *   `go.opentelemetry.io/otel` (API & SDK)
    *   `go.opentelemetry.io/otel/exporters/jaeger` or `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc`
    *   `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc`
    *   `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`

### 2.3. Java/Kotlin (Spring Boot ecosystem - for Payment Service, potentially Notification)
*   **HTTP Client:**
    *   **`Spring RestTemplate`** or **`WebClient`** (non-blocking).
*   **Logging:**
    *   **SLF4J with Logback/Log4j2** (standard in Spring Boot).
*   **Database Interaction (PostgreSQL):**
    *   **Spring Data JPA** with Hibernate.
*   **Caching Client (Redis):**
    *   **Spring Data Redis** with Lettuce or Jedis.
*   **Message Queue Client (Kafka):**
    *   **Spring Kafka**.
*   **Testing:**
    *   **JUnit 5, Mockito, Spring Boot Test**.
*   **OpenTelemetry (Observability):**
    *   OpenTelemetry Java SDK and auto-instrumentation agent.

---

### 2.4. Frontend (Flutter/Dart)
*   **Note:** The choice of packages should be consistent with `project_technology_stack.md`. This list aims to standardize packages for common frontend development needs. One primary solution from the alternatives should be chosen for consistency across the project.
*   **State Management:**
    *   **`flutter_bloc` / `bloc`**: For BLoC pattern.
    *   *Alternative:* `provider`, `riverpod`.
*   **Routing / Navigation:**
    *   **`go_router`**: For declarative routing.
    *   *Alternative:* `auto_route`.
*   **HTTP Client:**
    *   **`dio`**: Powerful HTTP client for Dart, supports interceptors, FormData, request cancellation, etc.
    *   *Alternative:* `http` (standard Dart package, for simpler needs).
*   **Local Storage:**
    *   **`hive` / `hive_flutter`**: Lightweight and fast NoSQL database.
    *   **`shared_preferences`**: For simple key-value data.
    *   **`flutter_secure_storage`**: For securely storing sensitive data.
    *   *Alternative for SQL:* `sqflite`.
*   **JSON Serialization/Deserialization:**
    *   **`json_serializable`** (build_runner based) / **`freezed`**: For generating boilerplate code for models.
*   **Equality & Immutability for Models:**
    *   **`equatable`**: For value equality.
    *   **`freezed`**: (Also covers this, often used with `json_serializable`).
*   **Dependency Injection:**
    *   **`get_it`**: Simple service locator.
    *   *Alternative:* `injectable` (code generator for `get_it`).
*   **Testing:**
    *   **`flutter_test`**: (SDK testing framework for unit and widget tests).
    *   **`bloc_test`**: For testing BLoCs/Cubits.
    *   **`mockito` / `mocktail`**: For creating mock objects.
    *   **`integration_test`**: (SDK testing framework for integration tests).
*   **Linting:**
    *   **`flutter_lints`** or **`lints`**: Official lint rules.
    *   **`dart_code_metrics`**: Additional static analysis metrics.
*   **Localization / Internationalization (i18n):**
    *   **`intl`** package with Flutter's `flutter_localizations` delegate.
*   **Utility:**
    *   **`dartz`**: Functional programming utilities (e.g., `Either`, `Option`).
    *   **`cached_network_image`**: For displaying images from the internet and keeping them in the cache.

---

## 3. Process for Adding New Standard Packages

1.  **Proposal:** If a new common task requires a library, or a better alternative to an existing standard is found, a developer can propose it. The proposal should include justification (why this library, alternatives considered).
2.  **Review:** The proposal should be reviewed by the Architecture Team or designated senior developers. Factors like those in Section 1 (Guiding Principles) will be considered.
3.  **Decision & Documentation:** If approved, the library is added to this document.
4.  **Communication:** The decision should be communicated to all development teams.

## 4. Versioning and Updates

*   It is recommended to use a dependency management tool (e.g., `npm`, `yarn`, `go mod`, `pip`) and define specific versions for packages.
*   Regularly review and update dependencies to their latest stable versions to incorporate bug fixes, performance improvements, and security patches. Tools like `npm audit` or GitHub's Dependabot can help identify vulnerabilities.
*   Major version upgrades of standard packages should be discussed and planned, as they might introduce breaking changes.

This document serves as a living guideline and should be updated as the project evolves and new needs arise.
