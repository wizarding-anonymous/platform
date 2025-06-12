# Account Service

## Overview

The Account Service is a foundational microservice within the Russian Steam analog platform. Its primary purpose is to manage user accounts, user profiles, and associated data such as settings, verification details, and contact information.

**For detailed specification, please see: [./docs/README.md](./docs/README.md)**

## Core Functionality (Summary)

*   Account Management
*   User Profile Management
*   Contact Information & Verification
*   Settings Management
*   Event Generation (Kafka)

## Technologies (Summary)

*   Backend: Go
*   Database: PostgreSQL
*   Cache: Redis
*   Messaging: Kafka
*   APIs: REST, gRPC

## Integrations (Summary)

Interacts with Auth Service, Notification Service, Social Service, Payment Service, Admin Service, API Gateway, and Kafka.
