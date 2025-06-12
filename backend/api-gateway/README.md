# API Gateway

## Overview

The API Gateway is a critical infrastructure component that serves as the single entry point for all external client requests to the Russian Steam analog platform. Its primary purpose is to provide a unified and secure facade for the various backend microservices. It handles common cross-cutting concerns such as request routing, authentication and authorization, rate limiting, and SSL/TLS termination, thereby simplifying client-side development and enhancing platform security and manageability.

## Core Functionality

*   **Request Routing:** Dynamically routes incoming client requests (HTTP, HTTPS, WebSocket) to the appropriate backend microservices based on path, host, method, or headers. Supports path rewriting.
*   **Authentication & Authorization Offloading:** Integrates with the Auth Service to validate client credentials (primarily JWTs, also API keys). It can perform initial authorization checks and injects user context (ID, roles) into requests forwarded to backend services.
*   **Rate Limiting:** Protects backend services from overload by enforcing limits on request frequency from individual clients (based on IP, user ID, or API key).
*   **SSL/TLS Termination:** Handles HTTPS requests, decrypting them and forwarding them (typically over HTTP or mTLS) to internal services.
*   **CORS Management:** Centralizes Cross-Origin Resource Sharing (CORS) policy enforcement for web applications.
*   **Security Enforcement:** Provides a layer of defense against common web attacks and can be integrated with a Web Application Firewall (WAF). Implements standard security headers.
*   **Request/Response Transformation:** (Optional) Modifies headers or, sparingly, the body of requests and responses.
*   **Service Discovery:** Integrates with Kubernetes to dynamically discover and route to healthy service instances.
*   **Monitoring & Logging:** Collects metrics, logs, and traces for all incoming traffic, providing observability into API usage and performance.
*   **Circuit Breaking:** (Typically a feature of the chosen gateway solution) Can prevent cascading failures by temporarily stopping requests to unhealthy backend services.
*   **WebSocket Proxying:** Supports proxying WebSocket connections to services like Notification or Social Service.

## Technologies

*   **Recommended Gateway Solutions:**
    *   **Kong Gateway:** Built on Nginx and Lua/OpenResty, known for high performance and an extensive plugin ecosystem.
    *   **Tyk Gateway:** Written in Go, offering a comprehensive feature set.
*   **Orchestration:** Kubernetes (v1.21+)
*   **Configuration Model:**
    *   Kubernetes Gateway API CRDs (Custom Resource Definitions) are a modern approach.
    *   Alternatively, specific CRDs for Kong (KongIngress, KongPlugin) or Tyk Operator.
    *   Configuration managed via Git (GitOps).
*   **Service Discovery:** Kubernetes Services.
*   **Ingress Control:** Typically works in conjunction with an Ingress Controller like Nginx Ingress or Traefik.
*   **Monitoring:** Prometheus, Grafana.
*   **Logging:** Loki, ELK Stack.
*   **Tracing:** Jaeger, OpenTelemetry.

## Configuration

The API Gateway's behavior, including routing rules, security policies, and plugin configurations, is defined declaratively.
*   **Routes:** Define how client requests (based on host, path, method, headers) are mapped to backend services.
*   **Services:** Represent the upstream microservices to which traffic is proxied.
*   **Policies/Plugins:** Implement cross-cutting concerns:
    *   **Authentication:** JWT validation (interacting with Auth Service), API key checks.
    *   **Rate Limiting:** Configured with thresholds and identifiers (IP, user, key).
    *   **CORS:** Defining allowed origins, methods, headers.
    *   **Header Modification:** Adding/removing/modifying headers (e.g., `X-User-Id`, `X-User-Roles`).
*   **Management:** Configuration is typically managed "as code" in Git and applied to the Kubernetes cluster using CI/CD and GitOps principles. Kubernetes Gateway API CRDs (like `HTTPRoute`, `Gateway`, `GatewayClass`) or solution-specific CRDs are used to define the gateway's behavior.

## Integrations

The API Gateway is the primary interaction point for external clients and fronts all backend microservices:

*   **Auth Service:** For validating JWTs and API keys. The gateway typically caches public keys from the Auth Service for efficient JWT signature verification.
*   **All Backend Microservices:** (Account, Catalog, Library, Payment, Social, Download, Notification, Analytics, Developer, Admin services) The API Gateway routes requests to these services based on its configuration. It enriches requests with user context after successful authentication.
*   **Kubernetes API:** For service discovery, allowing the gateway to find and route to active instances of backend services.
*   **Monitoring Systems (Prometheus, Jaeger):** Exports metrics and traces.
*   **Logging Systems (Loki, ELK):** Exports request/response logs.

## Security

The API Gateway plays a vital role in the platform's security posture:

*   **Centralized Authentication:** Enforces authentication for most API endpoints by validating JWTs or API keys.
*   **SSL/TLS Termination:** Encrypts external traffic using HTTPS.
*   **Hides Internal Network Topology:** Clients interact only with the gateway, not directly with internal services.
*   **Rate Limiting & IP Controls:** Mitigates DoS/DDoS attacks and abuse.
*   **Security Headers:** Applies standard HTTP security headers (HSTS, X-Frame-Options, CSP, etc.) to responses.
*   **Input Validation (Basic):** Can perform basic validation on incoming requests before forwarding.
*   **WAF Integration:** Can be placed behind or integrate with a Web Application Firewall for enhanced threat protection.
*   **CORS Control:** Manages cross-origin requests.
*   **Audit Logging:** Logs access attempts and policy enforcement actions.

Error responses are standardized as per [Стандарты API, форматов данных, событий и конфигурационных файлов.txt](https://placeholder.com/link-to-api-standards), typically involving a JSON body with `error.code` and `error.message`.
