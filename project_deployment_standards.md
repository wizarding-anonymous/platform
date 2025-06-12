# Стандарты Развертывания и Инфраструктурных Файлов

## 1. Введение

Данный документ определяет единые стандарты для инфраструктурных файлов (Dockerfile, Docker Compose, Kubernetes манифесты), процессов CI/CD и аспектов безопасности, связанных с развертыванием микросервисов российского аналога платформы Steam. Цель документа — обеспечить согласованность, автоматизацию и безопасность процессов развертывания и эксплуатации.

## 2. Стандарты Инфраструктурных Файлов

### 2.1. Dockerfile

1.  **Общие принципы:**
    *   Использовать многоэтапные сборки (multi-stage builds) для минимизации размера конечного образа.
    *   Использовать официальные и минимальные базовые образы (например, `alpine` для runtime, `golang:1.XX-alpine` для сборки Go-приложений).
    *   Не запускать процессы от `root` пользователя внутри контейнера. Создавать специального пользователя.
    *   Использовать тегированные версии базовых образов, избегать `latest`.
    *   Копировать только необходимые файлы в конечный образ.
    *   Определять `HEALTHCHECK` для контейнеров.
    *   Экспонировать только необходимые порты.

2.  **Пример Dockerfile для Go-сервиса:**
    ```dockerfile
    # Этап сборки
    FROM golang:1.21-alpine AS builder

    WORKDIR /app

    # Установка зависимостей
    COPY go.mod go.sum ./
    RUN go mod download && go mod verify

    # Копирование исходного кода
    COPY . .

    # Сборка приложения
    # RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server
    RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./cmd/your-service-main-package # Уточните путь к main

    # Финальный этап
    FROM alpine:latest # Или distroless/static для Go

    WORKDIR /app

    # Установка необходимых пакетов (например, ca-certificates, tzdata)
    RUN apk --no-cache add ca-certificates tzdata

    # Создание непривилегированного пользователя и группы
    RUN addgroup -S appgroup && adduser -S appuser -G appgroup

    # Копирование бинарного файла из этапа сборки
    COPY --from=builder /app/main .

    # Копирование конфигурационных файлов (если они не монтируются через ConfigMaps)
    # COPY --from=builder /app/configs ./configs
    # ENV CONFIG_PATH=/app/configs/config.yaml

    USER appuser

    # Определение порта
    EXPOSE 8080 # Замените на порт вашего сервиса

    # Запуск приложения
    ENTRYPOINT ["./main"]
    # CMD ["--config", "/app/configs/config.yaml"] # Если конфигурация через файл
    ```

### 2.2. Docker Compose (для локальной разработки)

1.  **Общие принципы:**
    *   Использовать актуальную версию формата Docker Compose (например, 3.8+).
    *   Группировать сервисы по их функциональности.
    *   Определять сети для изоляции и взаимодействия.
    *   Использовать именованные тома (volumes) для персистентных данных.
    *   Загружать конфигурацию и секреты через файлы или переменные окружения (из `.env` файла).
    *   Определять `healthcheck` для зависимых сервисов.

2.  **Пример `docker-compose.yml` (фрагмент):**
    ```yaml
    version: '3.8'

    services:
      api-gateway:
        build:
          context: ./api-gateway # Путь к Dockerfile API Gateway
          dockerfile: Dockerfile
        ports:
          - "8080:8080"
        environment:
          - SERVICE_ENV=development
        depends_on:
          auth-service:
            condition: service_healthy # Пример healthcheck-зависимости
        networks:
          - platform_network
        restart: unless-stopped
        healthcheck:
          test: ["CMD", "curl", "-f", "http://localhost:8080/health"] # Пример healthcheck
          # ... другие параметры healthcheck

      auth-service:
        build:
          context: ./auth-service
          dockerfile: Dockerfile
        environment:
          - SERVICE_ENV=development
          - AUTH_DB_PASSWORD=${AUTH_DB_PASSWORD} # Из .env файла
        depends_on:
          postgres_auth:
            condition: service_healthy
        networks:
          - platform_network
        restart: unless-stopped
        # ... healthcheck для auth-service

      postgres_auth: # Пример БД для Auth Service
        image: postgres:15-alpine
        environment:
          POSTGRES_USER: ${AUTH_DB_USER:-auth}
          POSTGRES_PASSWORD: ${AUTH_DB_PASSWORD:-password}
          POSTGRES_DB: ${AUTH_DB_NAME:-auth_db}
        volumes:
          - auth_postgres_data:/var/lib/postgresql/data
        networks:
          - platform_network
        restart: unless-stopped
        healthcheck:
          test: ["CMD-SHELL", "pg_isready -U ${AUTH_DB_USER:-auth}"]
          # ... другие параметры healthcheck

    networks:
      platform_network:
        driver: bridge

    volumes:
      auth_postgres_data:
    ```

### 2.3. Kubernetes Манифесты

1.  **Общие принципы:**
    *   Использовать Helm для управления манифестами и шаблонизации.
    *   Группировать ресурсы по микросервисам (например, в отдельные Helm charts).
    *   Использовать ConfigMaps для нечувствительной конфигурации и Secrets для чувствительных данных (пароли, API ключи).
    *   Определять `resources.requests` и `resources.limits` для CPU и памяти.
    *   Настраивать `livenessProbe` и `readinessProbe` для контейнеров.
    *   Использовать `RollingUpdate` стратегию для Deployments.
    *   Определять `NetworkPolicy` для контроля трафика между подами.
    *   Применять `PodSecurityPolicy` или Security Contexts для ограничения привилегий подов.

2.  **Примеры манифестов (Deployment, Service, ConfigMap, Secret):**
    *   (Подробные примеры YAML для Deployment, Service, ConfigMap, Secret см. в исходном документе "Стандарты API, форматов данных, событий и конфигурационных файлов.txt", раздел "Стандарты инфраструктурных файлов" -> "Kubernetes").
    *   Ключевые аспекты:
        *   **Deployment:** `replicas`, `selector`, `strategy`, `template` (с `image`, `ports`, `env`, `resources`, `probes`, `volumeMounts`).
        *   **Service:** `selector`, `ports`, `type` (обычно `ClusterIP` для внутренних сервисов).
        *   **ConfigMap:** Хранение конфигурационных файлов или ключ-значение пар.
        *   **Secret:** Хранение base64-кодированных секретов.

## 3. Стандарты CI/CD (Непрерывная Интеграция и Доставка)

### 3.1. Общие принципы
*   Автоматизация всех этапов: сборка, тестирование, анализ кода, сборка образов, развертывание.
*   Использование единой системы CI/CD для всех микросервисов (GitLab CI/CD или Jenkins, как указано в `project_technology_stack.md`).
*   Разделение на окружения: `development`, `testing`/`staging`, `production`.
*   Безопасное управление секретами и конфигурациями для разных окружений.
*   Возможность отката на предыдущую версию.
*   Уведомления о статусе сборки и развертывания.

### 3.2. Этапы пайплайна CI/CD
1.  **Сборка (Build):**
    *   Компиляция кода.
    *   Запуск статических анализаторов кода (`gosec`, `golangci-lint` для Go).
2.  **Тестирование (Test):**
    *   Запуск Unit-тестов.
    *   Запуск интеграционных тестов (с поднятием временных зависимостей, например, через Docker Compose или testcontainers).
    *   Расчет и проверка покрытия кода тестами.
3.  **Сканирование безопасности:**
    *   Сканирование зависимостей на известные уязвимости (`nancy` для Go, `trivy` для Docker-образов).
    *   Сканирование Docker-образов на уязвимости.
4.  **Сборка Образа (Package):**
    *   Сборка Docker-образа приложения.
    *   Тегирование образа версией коммита/тега Git.
    *   Загрузка образа в приватный реестр контейнеров.
5.  **Развертывание (Deploy):**
    *   **Dev/Testing:** Автоматическое развертывание в тестовое окружение после успешной сборки на ветке разработки.
    *   **Staging:** Развертывание в staging-окружение (может требовать ручного подтверждения) для E2E тестов и приемки.
    *   **Production:** Развертывание в production-окружение (требует ручного подтверждения или запускается по тегу Git) с использованием безопасных стратегий (Rolling Update, Canary, Blue/Green).
6.  **Пост-развертывание:**
    *   Запуск дымовых тестов (smoke tests).
    *   Мониторинг состояния развернутого приложения.

### 3.3. Примеры CI/CD конфигураций
*   (Примеры для GitHub Actions и GitLab CI/CD см. в исходном документе "Стандарты API, форматов данных, событий и конфигурационных файлов.txt", раздел "Стандарты инфраструктурных файлов" -> "CI/CD").
*   Ключевые элементы: определение `stages`, `jobs`, `scripts`, `rules`/`only` для управления потоком выполнения, использование переменных окружения и секретов CI/CD.

## 4. Безопасность Инфраструктуры и Развертывания

### 4.1. Статический анализ кода
*   **Инструменты для Go:** `gosec`, `golangci-lint`.
*   **Интеграция:** В CI/CD, блокировка сборки при критических уязвимостях.

### 4.2. Сканирование зависимостей
*   **Инструменты:** `nancy` (Go), `trivy` (Docker).
*   **Интеграция:** В CI/CD, блокировка сборки.

### 4.3. Безопасность контейнеров
*   **Базовые образы:** Минимальные (`scratch`, `alpine`).
*   **Пользователь:** Запуск от непривилегированного пользователя.
*   **Сканирование:** Перед деплоем.
*   **Иммутабельность:** Запрет на изменение запущенных контейнеров.

### 4.4. Безопасность Kubernetes
*   **RBAC:** Строгое разграничение прав доступа к API Kubernetes.
*   **Network Policies:** Изоляция сетевого трафика между подами.
*   **Pod Security Context / Policies:** Ограничение привилегий контейнеров (запрет `privileged`, `hostPID`, `hostNetwork`, настройка `runAsUser`, `fsGroup`).
*   **Secrets:** Шифрование секретов в `etcd` и использование механизмов вроде Vault для более безопасного управления.
*   Регулярное обновление Kubernetes до актуальных версий.
*   Использование Admission Controllers для применения политик безопасности.

---
*Разделы "Стандарты инфраструктурных файлов" и "CI/CD" извлечены из документа "(Дополнительный фаил) Стандарты API, форматов данных, событий и конфигурационных файлов.txt".*
*Разделы, касающиеся безопасности инфраструктуры, дополнены из документа "(Дополнительный фаил) Стандарты безопасности, мониторинга, логирования и трассировки.txt".*
*Этот документ должен регулярно обновляться и дополняться.*
