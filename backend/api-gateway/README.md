# API Gateway

## Overview

The API Gateway is a critical infrastructure component that serves as the single entry point for all external client requests to the Russian Steam analog platform. It handles request routing, authentication, rate limiting, SSL termination, and other cross-cutting concerns.

**For detailed specification, please see: [./docs/README.md](./docs/README.md)**

## Core Functionality (Summary)

*   Request Routing to microservices
*   Authentication & Authorization (JWT validation via Auth Service)
*   Rate Limiting
*   SSL/TLS Termination
*   CORS Management
*   Service Discovery (Kubernetes)
*   Monitoring, Logging, Tracing

## Technologies (Summary)

*   Kong Gateway or Tyk Gateway
*   Kubernetes (Configuration via CRDs, GitOps)
*   Integration with Prometheus, Grafana, Loki, Jaeger

## Configuration

Configuration is managed declaratively via Kubernetes CRDs stored in Git, following GitOps principles. This includes defining routes, services, and plugins for authentication, rate limiting, etc.
