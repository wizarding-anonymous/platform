# Спецификация Компонента: API Gateway

**Версия:** 1.0
**Дата последнего обновления:** {{YYYY-MM-DD}}

## 1. Обзор Сервиса (Overview)

### 1.1. Назначение и Роль
*   API Gateway является критически важным инфраструктурным компонентом, служащим единой точкой входа для всех внешних клиентских запросов к платформе "Российский Аналог Steam".
*   Его основная роль - предоставить унифицированный и безопасный фасад для различных бэкенд-микросервисов. Он обрабатывает общие сквозные задачи, такие как маршрутизация запросов, аутентификация и авторизация, ограничение скорости (rate limiting) и терминирование SSL/TLS.
*   Основные бизнес-задачи: Упрощение клиентской разработки, повышение безопасности и управляемости платформы, централизация управления доступом к API, снижение нагрузки на бэкенд-сервисы.
*   Разработка конфигурации и управление API Gateway должны вестись в соответствии с `../../../../CODING_STANDARDS.md` (в части управления конфигурационными файлами) и `../../../../project_deployment_standards.md` (GitOps принципы).

### 1.2. Ключевые Функциональности
*   **Маршрутизация запросов:** Динамическая маршрутизация к бэкенд-микросервисам.
*   **Аутентификация и Авторизация:** Интеграция с Auth Service (JWT, API ключи), внедрение контекста пользователя.
*   **Ограничение скорости (Rate Limiting):** Защита от перегрузки.
*   **Терминирование SSL/TLS:** Обработка HTTPS, управление сертификатами.
*   **Управление CORS:** Централизованное применение политик.
*   **Обеспечение безопасности:** Слой защиты (WAF интеграция), заголовки безопасности.
*   **Трансформация запросов/ответов (опционально):** Модификация заголовков/тела.
*   **Обнаружение сервисов (Service Discovery):** Интеграция с Kubernetes.
*   **Мониторинг, Логирование и Трассировка:** Сбор данных для всего трафика.
*   **Проксирование WebSocket:** Поддержка WebSocket соединений.
*   **Кэширование ответов (опционально):** Для снижения нагрузки.

### 1.3. Основные Технологии
*   **Gateway Software:** Kong Gateway или Tyk Gateway (согласно `../../../../project_technology_stack.md`).
*   **Оркестрация и Управление Конфигурацией:** Kubernetes (1.25+), Kubernetes Gateway API CRDs (предпочтительно) или специфичные CRD вендора, GitOps.
*   **Обнаружение сервисов:** Kubernetes Services DNS.
*   **Ingress Control:** Связка с Ingress Controller (Nginx Ingress, Traefik Ingress).
*   **Мониторинг и Логирование:** Интеграция с Prometheus, Grafana, Loki/ELK, Jaeger/OpenTelemetry.
*   Выбор сторонних библиотек и пакетов должен осуществляться согласно `../../../../PACKAGE_STANDARDIZATION.md` (применимо к любым кастомным плагинам или управляющим скриптам).
*   Ссылки на: `../../../../project_technology_stack.md`, `../../../../project_deployment_standards.md`, `../../../../project_observability_standards.md`, `../../../../project_glossary.md`.

### 1.4. Термины и Определения (Glossary)
*   См. `../../../../project_glossary.md`.
*   **CRD (Custom Resource Definition):** Расширение API Kubernetes.
*   **Ingress Controller:** Компонент Kubernetes для управления внешним доступом.
*   **Kubernetes Gateway API:** Новая спецификация Kubernetes для моделирования сетевых сервисов.
*   **GitOps:** Подход к управлению инфраструктурой через Git.

## 2. Внутренняя Архитектура (Internal Architecture)

### 2.1. Общее Описание
*   API Gateway функционирует как обратный прокси. Его архитектура определяется выбранным ПО (Kong, Tyk) и интеграцией с Kubernetes.
*   Логика работы Gateway: применение набора правил и политик (маршруты, плагины/фильтры) к трафику, декларативно конфигурируемых через Kubernetes CRDs и GitOps.
*   Диаграмма верхнеуровневой архитектуры:
    ```mermaid
    graph TD
        Clients[Клиентские приложения (Web, Mobile, Desktop)] --> Internet[Интернет]
        Internet --> Firewall[Межсетевой экран/WAF]
        Firewall --> ExternalLoadBalancer[Внешний Балансировщик Нагрузки (L4)]
        ExternalLoadBalancer --> K8sIngress[Kubernetes Ingress Controller (Nginx/Traefik)]
        K8sIngress --> APIGateway[API Gateway (Kong/Tyk) Data Plane]

        subgraph "Kubernetes Cluster"
            APIGateway
            APIGatewayControlPlane[API Gateway Control Plane] -.-> APIGateway
            K8sAPIServer[Kubernetes API Server (CRDs)] <--> APIGatewayControlPlane
            GitOpsController[GitOps Controller (ArgoCD/Flux)] --> K8sAPIServer

            subgraph "Backend Microservices"
                AuthSvc[Auth Service]
                AccountSvc[Account Service]
                CatalogSvc[Catalog Service]
                OtherSvcs[...]
            end
            APIGateway --> AuthSvc
            APIGateway --> AccountSvc
            APIGateway --> CatalogSvc
            APIGateway --> OtherSvcs
        end

        ObservabilityStack[Системы Мониторинга/Логирования (Prometheus, Loki, Jaeger)]
        APIGateway -.-> ObservabilityStack
        APIGatewayControlPlane -.-> ObservabilityStack

        GitRepo[Git Репозиторий (Конфигурация Gateway)] --> GitOpsController

    classDef internet fill:#f0f8ff,stroke:#87cefa,stroke-width:2px;
    classDef security fill:#fff0f5,stroke:#dda0dd,stroke-width:2px;
    classDef loadbalancer fill:#e0ffff,stroke:#afeeee,stroke-width:2px;
    classDef gateway fill:#e6e6fa,stroke:#9370db,stroke-width:2px;
    classDef k8s fill:#f5f5f5,stroke:#d3d3d3,stroke-width:2px;
    classDef services fill:#d4edda,stroke:#28a745,stroke-width:2px;

    class Clients,Internet internet;
    class Firewall security;
    class ExternalLoadBalancer,K8sIngress loadbalancer;
    class APIGateway,APIGatewayControlPlane gateway;
    class K8sAPIServer,GitOpsController,ObservabilityStack k8s;
    class AuthSvc,AccountSvc,CatalogSvc,OtherSvcs services;
    class GitRepo external;
    ```

### 2.2. Слои Сервиса (применительно к функционированию и конфигурации Gateway)

#### 2.2.1. Presentation Layer (Слой Конфигурации / Configuration Plane)
*   Ответственность: Определение способов конфигурации Gateway администраторами или CI/CD.
*   Ключевые компоненты/модули:
    *   Kubernetes API Server (для CRDs).
    *   CRDs: `GatewayClass`, `Gateway`, `HTTPRoute`, etc. (Kubernetes Gateway API) или специфичные CRD вендора.
    *   GitOps контроллер (ArgoCD, Flux).

#### 2.2.2. Application Layer (Слой Управления / Control Plane - ядро Gateway)
*   Ответственность: Трансляция декларативной конфигурации (CRD) во внутренние правила прокси. Динамическое обновление конфигурации. Обнаружение сервисов.
*   Ключевые компоненты/модули: Внутренние компоненты выбранного ПО API Gateway.

#### 2.2.3. Domain Layer (Не применимо в традиционном смысле)
*   API Gateway не содержит бизнес-логики или доменных сущностей в том смысле, как это делают микросервисы. Его "домен" - это управление API-трафиком и применение политик.

#### 2.2.4. Infrastructure Layer (Слой Данных / Data Plane и Слой Интеграций)
*   **Data Plane (Проксирующий движок):**
    *   Ответственность: Непосредственная обработка и проксирование запросов. Применение плагинов/фильтров.
    *   Ключевые компоненты: Высокопроизводительный прокси-сервер (Nginx в Kong, собственный движок в Tyk).
*   **Слой Интеграций:**
    *   Ответственность: Взаимодействие с внешними системами.
    *   Ключевые компоненты: Клиенты Kubernetes API, HTTP/gRPC клиенты к Auth Service, экспортеры метрик/логов/трейсов.

## 3. API Endpoints (Конфигурируемые маршруты)
API Gateway **маршрутизирует** запросы к эндпоинтам бэкенд-микросервисов.
*   См. документацию каждого микросервиса для деталей их API.
*   См. `../../../../project_api_standards.md` для общих стандартов API.
*   Конфигурация маршрутов осуществляется декларативно (пример `HTTPRoute` см. в существующем документе).
*   Полная конфигурация находится в GitOps репозитории (например, `deploy/platform-gitops/gateway-config/`). См. `../../../../project_deployment_standards.md`.
*   [NEEDS DEVELOPER INPUT: Add more diverse or specific examples of routing rules if the existing one isn't sufficiently illustrative, e.g., WebSocket route, route with header manipulation, or specific plugin configurations beyond basic auth.]

### 3.1. REST API
*   Не применимо (Gateway проксирует, а не определяет REST API).

### 3.2. gRPC API
*   Не применимо (Gateway проксирует, а не определяет gRPC API).
*   Ссылка на `.proto` файл: Не применимо.

### 3.3. WebSocket API (если применимо)
*   Поддерживается проксирование WebSocket соединений. Конфигурация маршрутов для WebSocket аналогична HTTPRoute, но может использовать `TCPRoute` или специфичные аннотации/CRD в зависимости от Gateway.

## 4. Модели Данных (Конфигурационные)
*   "Модели данных" API Gateway - это структуры его конфигурации, определяемые CRD (например, `Gateway`, `HTTPRoute`, `KongPlugin`, `TykAPIDefinition`).
*   Бизнес-данные не хранятся.

## 5. Потоковая Обработка Событий (Event Streaming)
*   Напрямую не участвует в обработке бизнес-событий Kafka.
*   Генерирует операционные события/логи (доступ, ошибки) для систем мониторинга.
*   Проксирует WebSocket для сервисов real-time коммуникаций.

## 6. Интеграции (Integrations)
См. `../../../../project_integrations.md`.
[NEEDS DEVELOPER INPUT: Review and confirm service boundaries and reliability strategies for each integration point for api-gateway service, especially Auth Service and Kubernetes API.]

### 6.1. Внутренние Микросервисы
*   **Auth Service:** HTTP/gRPC вызовы для валидации токенов/ключей. Надежность: `[NEEDS DEVELOPER INPUT: Strategy for Auth Service unavailability, e.g., fail-open/closed, caching]`
*   **Все бэкенд-микросервисы:** Проксирование запросов, добавление заголовков (`X-User-ID`, `X-User-Roles`).
*   **Kubernetes API:** Service Discovery, чтение CRD. Надежность: `[NEEDS DEVELOPER INPUT: Impact of K8s API unavailability]`

### 6.2. Внешние Системы
*   **Клиентские приложения.**
*   **Системы Мониторинга/Логирования/Трассировки.**
*   **WAF (Web Application Firewall):** (Опционально, может быть перед Gateway).

## 7. Конфигурация (Configuration)
Общие стандарты: `../../../../project_api_standards.md` и `../../../../DOCUMENTATION_GUIDELINES.md`.

### 7.1. Переменные Окружения (для ПО Gateway)
*   Переменные для настройки базовых параметров работы самого ПО Gateway (уровни логирования, подключение к БД конфигурации Gateway если используется, адреса админ-API).
*   Примеры: `KONG_LOG_LEVEL`, `KONG_DATABASE`, `TYK_LOGLEVEL`.
*   Управляются через Kubernetes Deployments/StatefulSets и Helm-чарты. Секреты через Kubernetes Secrets.

### 7.2. Файлы Конфигурации (Декларативная конфигурация маршрутов и политик)
*   Расположение: GitOps репозиторий (например, `deploy/gitops/api-gateway-config/`).
*   Структура: Набор YAML файлов, определяющих CRD (HTTPRoute, KongPlugin, etc.).
*   [NEEDS DEVELOPER INPUT: Provide a link to the actual GitOps configuration repository or a more detailed example structure if deemed necessary beyond the HTTPRoute example.]

## 8. Обработка Ошибок (Error Handling)
См. `../../../../project_api_standards.md`. API Gateway генерирует собственные ошибки или проксирует ошибки бэкендов.

### 8.1. Общие Принципы
*   Стандартный JSON формат ответа об ошибке.

### 8.2. Распространенные Коды Ошибок (генерируемые Gateway)
*   `400 INVALID_REQUEST_SYNTAX`
*   `401 AUTHENTICATION_FAILED`
*   `403 FORBIDDEN_BY_GATEWAY_POLICY`
*   `404 ROUTE_NOT_FOUND`
*   `429 RATE_LIMIT_EXCEEDED`
*   `502 UPSTREAM_SERVICE_ERROR`
*   `503 SERVICE_UNAVAILABLE` (Gateway или бэкенд перегружен/недоступен)
*   `504 UPSTREAM_TIMEOUT`
*   Примеры ответов см. в существующем документе.

## 9. Безопасность (Security)
См. `../../../../project_security_standards.md`.

### 9.1. Аутентификация
*   Валидация JWT/API ключей через Auth Service.

### 9.2. Авторизация
*   Базовая проверка ролей/разрешений.

### 9.3. Защита Данных
*   Терминирование SSL/TLS.
*   Применение политик WAF.
*   Сокрытие внутренней архитектуры.

### 9.4. Управление Секретами
*   Сертификаты SSL/TLS, ключи для плагинов и т.д. управляются через Kubernetes Secrets.

## 10. Развертывание (Deployment)
См. `../../../../project_deployment_standards.md`.

### 10.1. Инфраструктурные Файлы
*   Dockerfile (для кастомных образов Gateway если необходимо): `[NEEDS DEVELOPER INPUT: Path to Dockerfile if custom images are built, otherwise state "Uses standard vendor images"]`
*   Helm-чарты: `deploy/charts/api-gateway/` (или чарт вендора).

### 10.2. Зависимости при Развертывании
*   Kubernetes кластер.
*   Auth Service.
*   Доступность бэкенд-сервисов для маршрутизации.
*   Системы мониторинга/логирования.

### 10.3. CI/CD
*   Конфигурация Gateway управляется через GitOps.
*   CI/CD пайплайны для валидации конфигурации, тестирования (например, с помощью инструментов типа `kuttl` для CRD) и развертывания обновлений ПО Gateway.

## 11. Мониторинг и Логирование (Logging and Monitoring)
См. `../../../../project_observability_standards.md`.

### 11.1. Логирование
*   Формат: JSON.
*   Ключевые события: Запросы, ответы, ошибки, информация о маршрутизации, результаты аутентификации/авторизации.
*   Интеграция: Loki/ELK. Логи доступа с `trace_id`.

### 11.2. Мониторинг
*   Метрики (Prometheus): Количество запросов, задержка (общая, бэкендов, Gateway), активные соединения, состояние upstream-сервисов, ошибки.
*   Дашборды (Grafana): Обзор трафика, производительности, ошибок.
*   Алерты (AlertManager): Высокий % ошибок, рост задержки, недоступность upstream-сервисов.

### 11.3. Трассировка
*   Интеграция: OpenTelemetry. Экспорт: Jaeger.

## 12. Нефункциональные Требования (NFRs)
*   **Производительность**: Доп. задержка (P99) Gateway < 5-10мс; JWT валидация (P99) < 2мс. Пропускная способность >10-20k RPS.
*   **Масштабируемость**: Горизонтальное масштабирование (HPA). До 5000 маршрутов.
*   **Надежность**: Доступность > 99.99% (data plane). Время перезагрузки конфигурации P95 < 30с.
*   **Безопасность**: Защита OWASP Top 10. Безопасная обработка токенов.
*   [NEEDS DEVELOPER INPUT: Confirm or update specific NFR values for api-gateway service, especially for throughput and route capacity based on chosen vendor/solution.]

## 13. Резервное Копирование и Восстановление (Backup and Recovery)
Общие принципы см. в `../../../../project_database_structure.md`.

### 13.1. Конфигурация API Gateway (CRDs)
*   Процедура: GitOps (Git-репозиторий - источник правды). Дополнительно Velero для бэкапа CRD из кластера.
*   RTO: < 1ч. RPO: Близко к нулю.

### 13.2. Собственная База Данных API Gateway (если используется)
*   Процедура: Стандартные бэкапы PostgreSQL (если Kong не DB-less).
*   RTO: < 2-4ч. RPO: < 5-15мин.
*   **Примечание:** Предпочтителен DB-less режим.

## 14. Приложения (Appendices) (Опционально)
*   Детальные примеры конфигурации CRD доступны в GitOps репозитории (например, `deploy/platform-gitops/gateway-config/`).
*   [NEEDS DEVELOPER INPUT: Add any other appendices if necessary for api-gateway service.]

## 15. Связанные Рабочие Процессы (Related Workflows)
*   [Процесс добавления нового микросервиса и его API на платформу] `[NEEDS DEVELOPER INPUT: Add link to a workflow document describing how a new service is exposed via the Gateway, e.g., project_workflows/new_service_onboarding_flow.md]`
*   [Процесс обновления SSL-сертификата для домена платформы] `[NEEDS DEVELOPER INPUT: Add link to SSL management workflow, e.g., project_workflows/ssl_certificate_management_flow.md]`

---
*Этот документ является отправной точкой и должен регулярно обновляться по мере развития конфигурации API Gateway.*
