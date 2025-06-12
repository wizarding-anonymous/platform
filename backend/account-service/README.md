# Account Service

## Overview

The Account Service is a foundational microservice within the Russian Steam analog platform. Its primary purpose is to manage user accounts, user profiles, and associated data such as settings, verification details, and contact information. It serves as the authoritative source for user-related data, providing this information to other services and ensuring data integrity and security.

## Core Functionality

*   **Account Management:** Handles user registration (via various methods including email/password and social providers through Auth Service), stores and manages basic account information (ID, Username, status), and supports account status updates (activation, blocking, deletion).
*   **User Profile Management:** Manages user profiles, including creation, editing of details (nickname, bio, country, avatar), and visibility settings.
*   **Contact Information & Verification:** Manages user contact details (email, phone numbers) and handles the verification process for them.
*   **Settings Management:** Stores and allows updates to user-specific settings, categorized for areas like privacy, notifications, and interface preferences.
*   **Event Generation:** Publishes events related to account and profile changes (e.g., account creation, profile updates) to a Kafka message broker.

## Technologies

*   **Backend Language:** Go (version 1.21+)
*   **REST Framework:** Gin or Echo
*   **gRPC Framework:** google.golang.org/grpc
*   **Database:** PostgreSQL (version 15+)
*   **Cache/Session Management:** Redis
*   **Messaging:** Kafka (using confluent-kafka-go or sarama)
*   **Validation:** go-playground/validator
*   **ORM/DB Driver:** GORM or sqlx
*   **Frontend (Interacts With):** Flutter/Dart
*   **Infrastructure:** Docker, Kubernetes, Prometheus, Grafana, ELK Stack/Loki, Jaeger

## API Summary

The service exposes RESTful APIs for client interaction (typically via an API Gateway) and gRPC APIs for inter-service communication. API design adheres to the platform's [Стандарты API, форматов данных, событий и конфигурационных файлов.txt](https://placeholder.com/link-to-api-standards).

### REST API (Base URL: `/api/v1`)

Key endpoints (proxied via API Gateway, often under `/users` or `/accounts`):

*   **Accounts:**
    *   `POST /accounts`: Register a new user account. (Often initiated via Auth Service)
    *   `GET /accounts/me`: Get current user's account information.
    *   `GET /accounts/{id}`: Get specific user's account information.
    *   `PUT /accounts/{id}/status`: (Admin) Update user account status.
    *   `DELETE /accounts/{id}`: (Admin/User) Delete a user account (soft delete).
*   **Profiles:**
    *   `GET /accounts/me/profile`: Get current user's profile.
    *   `GET /accounts/{id}/profile`: Get specific user's profile.
    *   `PUT /accounts/{id}/profile`: Update user profile.
    *   `POST /accounts/{id}/avatar`: Upload user avatar.
*   **Contact Info & Verification:**
    *   `GET /accounts/{id}/contact-info`: Get user's contact methods.
    *   `POST /accounts/{id}/contact-info`: Add a new contact method.
    *   `POST /accounts/{id}/contact-info/{type}/verification-request`: Request a verification code.
    *   `POST /accounts/{id}/contact-info/{type}/verify`: Submit a verification code.
*   **Settings:**
    *   `GET /accounts/{id}/settings/{category}`: Get user settings for a category.
    *   `PUT /accounts/{id}/settings/{category}`: Update user settings.

### gRPC API (Package: `account.v1`)

Key service methods:

*   `rpc GetAccount(GetAccountRequest) returns (AccountResponse)`: Retrieve account details.
*   `rpc GetAccounts(GetAccountsRequest) returns (GetAccountsResponse)`: Retrieve multiple account details.
*   `rpc GetProfile(GetProfileRequest) returns (ProfileResponse)`: Retrieve user profile.
*   `rpc CheckUsernameExists(CheckUsernameExistsRequest) returns (CheckExistsResponse)`
*   `rpc CheckEmailExists(CheckEmailExistsRequest) returns (CheckExistsResponse)`

## Data Model

The service manages several core entities:

*   **Account:** Represents the user's primary account with ID, username, status (e.g., `pending`, `active`, `blocked`), and timestamps.
*   **AuthMethod:** (Often managed in conjunction with Auth Service) Details about authentication methods linked to an account (e.g., password, Google ID).
*   **Profile:** Contains user-displayable information like nickname, bio, avatar URL, country, city, birth date, and visibility settings.
*   **ContactInfo:** Stores user's email addresses and phone numbers, their verification status, and primary contact flags.
*   **Setting:** Stores user preferences categorized (e.g., `privacy`, `notifications`) as JSONB.
*   **Avatar:** Manages uploaded avatar images and their URLs.
*   **ProfileHistory:** Logs changes to user profiles.

Data is stored in PostgreSQL, with Redis used for caching frequently accessed data.

## Integrations

Account Service interacts with several other microservices:

*   **Auth Service:** Delegates credential management and relies on it for authentication processes. Receives user identifiers post-authentication.
*   **Notification Service:** Publishes events or calls its API to send notifications for actions like email/phone verification or account status changes.
*   **Social Service:** Provides profile data to the Social Service.
*   **Payment Service:** Provides basic account information for transaction processing.
*   **Admin Service:** Allows administrators to manage user accounts and profiles.
*   **API Gateway:** Exposes REST endpoints to clients and handles initial request authentication/authorization.
*   **Kafka:** Publishes events like `account.created`, `profile.updated`, `account.contact.verified` for other services to consume.

## Error Handling

Error responses from the REST API follow a standard JSON structure:
```json
{
  "status": "error",
  "error": {
    "code": "ERROR_CODE", // e.g., "VALIDATION_ERROR", "RESOURCE_NOT_FOUND"
    "message": "Human-readable error description.",
    "details": { /* Optional, e.g., field-specific validation errors */ }
  }
}
```
gRPC APIs use standard gRPC error codes and statuses. All significant errors are logged with a trace ID for debugging.
Common error codes include `VALIDATION_ERROR`, `RESOURCE_NOT_FOUND`, `CONFLICT` (e.g., username/email already exists), and `INTERNAL_SERVER_ERROR`.
