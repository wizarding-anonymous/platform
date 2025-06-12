# Стандарты API, Форматов Данных, Событий и Конфигурационных Файлов

## 1. Введение

Данный документ определяет единые стандарты для API (REST, gRPC, WebSocket), форматов данных, событий и конфигурационных файлов, используемых во всех микросервисах российского аналога платформы Steam. Цель документа — обеспечить согласованность и совместимость между всеми компонентами системы, упростить интеграцию и поддержку.

## 2. Стандарты REST API

### 2.1. Общие принципы

1.  **Версионирование:**
    *   Версия API указывается в URL: `/api/v1/resource`.
    *   Мажорная версия (v1, v2) меняется при несовместимых изменениях.
    *   Минорные изменения (добавление новых полей) не требуют изменения версии.

2.  **Формат URL:**
    *   Использовать существительные во множественном числе для ресурсов: `/api/v1/games`.
    *   Вложенные ресурсы: `/api/v1/games/{game_id}/reviews`.
    *   kebab-case для составных слов: `/api/v1/payment-methods`.
    *   Специальные действия: `/api/v1/games/{game_id}/publish`.

3.  **HTTP-методы:**
    *   `GET`: получение ресурса/коллекции.
    *   `POST`: создание нового ресурса.
    *   `PUT`: полное обновление ресурса.
    *   `PATCH`: частичное обновление ресурса.
    *   `DELETE`: удаление ресурса.

4.  **Коды ответов:**
    *   `200 OK`: Успешный запрос.
    *   `201 Created`: Успешное создание.
    *   `204 No Content`: Успешный запрос без тела ответа.
    *   `400 Bad Request`: Ошибка в запросе.
    *   `401 Unauthorized`: Отсутствие аутентификации.
    *   `403 Forbidden`: Недостаточно прав.
    *   `404 Not Found`: Ресурс не найден.
    *   `409 Conflict`: Конфликт создания/обновления.
    *   `422 Unprocessable Entity`: Ошибка валидации.
    *   `429 Too Many Requests`: Превышение лимита.
    *   `500 Internal Server Error`: Внутренняя ошибка сервера.

5.  **Пагинация:**
    *   Параметры: `page` (номер страницы, с 1), `per_page` (количество на странице, max 100).
    *   Ответ должен содержать метаданные пагинации и ссылки (`self`, `first`, `prev`, `next`, `last`).
    ```json
    {
      "data": [...],
      "meta": { "page": 1, "per_page": 20, "total_pages": 5, "total_items": 97 },
      "links": { ... }
    }
    ```

6.  **Фильтрация:**
    *   Простые фильтры: `/api/v1/games?genre=strategy&price_min=100`.
    *   Сложные фильтры: `/api/v1/games?filter={"genre":["strategy","rpg"]}`.

7.  **Сортировка:**
    *   Параметр `sort`: `?sort=price` (возрастание), `?sort=-price` (убывание).
    *   Множественная: `?sort=genre,-price`.

8.  **Выборка полей:**
    *   Параметр `fields`: `?fields=id,title,price`.
    *   Вложенные поля: `?fields=id,title,developer{id,name}`.

9.  **Формат ответа (JSON API-подобный):**
    *   Одиночный ресурс:
        ```json
        {
          "data": {
            "id": "uuid",
            "type": "resource_type",
            "attributes": { ... },
            "relationships": { ... }
          }
        }
        ```
    *   Коллекция:
        ```json
        {
          "data": [ { "id": "uuid", "type": "resource_type", "attributes": { ... } }, ... ],
          "meta": { ... },
          "links": { ... }
        }
        ```

10. **Формат ошибок:**
    ```json
    {
      "errors": [
        {
          "code": "error_code_string",
          "title": "Human-readable title",
          "detail": "Detailed message.",
          "source": { "pointer": "/data/attributes/field" }
        }
      ]
    }
    ```

11. **Заголовки:**
    *   `Content-Type: application/json`
    *   `Accept: application/json`
    *   `Authorization: Bearer <token>`
    *   `X-Request-ID: <uuid>`
    *   `X-API-Key: <key>` (альтернатива JWT)

12. **Документация:**
    *   OpenAPI (Swagger) 3.0 для каждого REST API.
    *   Доступна по `/api/v1/docs`.

### 2.2. Специфические требования для микросервисов

*   **Auth Service:** Публичные эндпоинты для регистрации/логина. Эндпоинты для валидации/обновления токенов.
*   **API Gateway:** Добавляет заголовки `X-User-Id`, `X-User-Roles`, `X-Original-IP` для внутренних сервисов.
*   **Catalog Service:** Публичные эндпоинты для каталога, защищенные для управления.
*   **Payment Service:** Все эндпоинты через HTTPS. Webhook для уведомлений от платежных систем.

## 3. Стандарты gRPC API

### 3.1. Общие принципы

1.  **Версионирование:** В имени пакета: `package platform.v1.service;`.
2.  **Именование:**
    *   Сервисы: `PascalCaseService` (e.g., `UserService`).
    *   Методы: `PascalCaseAction` (e.g., `GetUser`).
    *   Сообщения: `PascalCaseNoun` (e.g., `User`, `CreateUserRequest`).
    *   Поля: `snake_case` (e.g., `user_id`).
    *   Перечисления: `PascalCaseEnum`, значения `ENUM_UPPER_SNAKE_CASE`.
3.  **Структура .proto файлов:** Каждый сервис в отдельном файле. Общие сообщения/enum в `common.proto`.
4.  **Формат сообщений:** Запросы `{Method}Request`, ответы `{Method}Response`. `google.protobuf.Timestamp`, `google.protobuf.Empty`.
5.  **Типы методов:** Унарные, серверные/клиентские/двунаправленные потоки.
6.  **Обработка ошибок:** Стандартные коды gRPC, метаданные для деталей.
7.  **Документация:** Комментарии в формате Protodoc.
8.  **Безопасность:** TLS, токены в метаданных (`authorization: Bearer <token>`).

### 3.2. Пример определения сервиса
```protobuf
syntax = "proto3";

package platform.v1.user;

import "google/protobuf/timestamp.proto";
// import "common.proto"; // Если есть общие типы

option go_package = "github.com/company/platform/api/grpc/v1/user";

// UserService предоставляет методы для управления пользователями.
service UserService {
  // GetUser возвращает информацию о пользователе по ID.
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  // ... другие методы
}

message GetUserRequest {
  string user_id = 1; // ID пользователя.
}

message GetUserResponse {
  User user = 1; // Информация о пользователе.
}

message User {
  string id = 1;
  string username = 2;
  string email = 3;
  // ... другие поля
  google.protobuf.Timestamp created_at = 5;
}
```

## 4. Стандарты WebSocket API

### 4.1. Общие принципы

1.  **Подключение:** URL `/api/v1/ws/{service}`. Аутентификация через query параметр `token` или заголовок `Authorization`. Ping/Pong.
2.  **Формат сообщений (JSON):**
    ```json
    {
      "type": "message_type_string", // Тип сообщения
      "id": "unique_message_id_uuid", // Уникальный ID сообщения (для ack)
      "payload": { ... } // Полезная нагрузка
    }
    ```
3.  **Обработка ошибок:** Сообщение типа `error` с `code` и `message` в `payload`.
4.  **Подтверждение доставки (Ack):** Для важных сообщений, тип `ack` с `original_message_id`.
5.  **Документация:** Описание всех типов сообщений и их структур.

### 4.2. Специфические требования
*   **Social Service (Chat):** Типы `chat_message`, `typing_status`, `read_receipt`.
*   **Notification Service:** Тип `notification`, `notification_read`.

## 5. Форматы Данных

### 5.1. Общие принципы

1.  **JSON:** Для REST API и WebSocket. camelCase для имен полей. UTF-8. Даты ISO 8601 (`YYYY-MM-DDTHH:mm:ss.sssZ`).
2.  **Protocol Buffers:** Для gRPC. snake_case для полей. `google.protobuf.Timestamp` для дат/времени.
3.  **Общие типы:**
    *   Идентификаторы: UUID v4 (строка).
    *   Денежные значения: Целое число (копейки/центы) для внутренних операций, строка с десятичной точкой для отображения API.
    *   Перечисления: Строковые константы (REST), числовые (gRPC).
4.  **Локализация:** JSON объект `{"ru": "Текст", "en": "Text"}`. Коды языков ISO 639-1 (`ru`, `en`).
5.  **Валидация:** Ограничения для типов данных должны быть документированы в OpenAPI/Protobuf.

### 5.2. Стандартные объекты (Примеры)
*   **User:** `id`, `username`, `email`, `status`, `createdAt`, `updatedAt`, `roles`.
*   **Game:** `id`, `title` (локализовано), `description` (локализовано), `price`, `developer`, `publisher`, `genres`, `tags`.
*   **Transaction:** `id`, `userId`, `type`, `status`, `amount`, `currency`, `items`.
*   **Error:** `errors: [ { code, title, detail, source: { pointer } } ]`.
*   (Более подробные примеры см. в исходном документе "Стандарты API...")

## 6. Стандарты Событий (Kafka)

### 6.1. Общие принципы

1.  **Формат события (CloudEvents-подобный):**
    ```json
    {
      "id": "uuid", // Уникальный ID события
      "type": "domain.resource.action.v1", // Тип события с версией
      "source": "service-name", // Источник события
      "time": "YYYY-MM-DDTHH:mm:ssZ", // Время события (UTC)
      "dataContentType": "application/json",
      "data": { ... }, // Полезная нагрузка события
      "subject": "resource_id", // ID основного ресурса
      "correlationId": "uuid" // ID для трассировки
    }
    ```
2.  **Именование типов событий:** `domain.resource.action` (e.g., `user.registered`, `game.published`). Глагол в прошедшем времени.
3.  **Версионирование событий:** В типе события (e.g., `user.registered.v1`).
4.  **Обработка событий:** Идемпотентность, порядок (по `time`), отказоустойчивость.
5.  **Топики Kafka:** Именование `{service}.{resource}.{action}` (e.g., `auth.user.registered`). Партиционирование по `subject`. Репликация >= 3. Retention >= 7 дней.

### 6.2. Стандартные события (Примеры)
*   **User Events:** `user.registered`, `user.verified`, `user.updated`, `user.deleted`, `user.logged_in`.
*   **Game Events:** `game.created`, `game.updated`, `game.published`, `game.price_changed`.
*   **Payment Events:** `payment.initiated`, `payment.completed`, `payment.failed`.
*   (Более подробный список см. в исходном документе "Стандарты API...")

## 7. Стандарты Конфигурационных Файлов

### 7.1. Общие принципы

1.  **Формат:** YAML. snake_case для параметров. Комментарии `#`.
2.  **Структура:** Группировка по секциям, вложенность. Без дублирования.
3.  **Переменные окружения:** Чувствительные данные из переменных окружения (`SERVICE_SECTION_PARAMETER`). Плейсхолдеры в YAML: `${ENV_VAR_NAME}`.
4.  **Профили окружений:** `config.yaml` (базовый), `config.{env}.yaml` (специфичный для окружения).
5.  **Валидация:** При запуске сервиса.

### 7.2. Пример конфигурационного файла (YAML)
```yaml
service:
  name: auth-service
  version: 1.0.0

http:
  port: 8080

database:
  driver: postgres
  host: postgres
  port: 5432
  password: "${AUTH_DB_PASSWORD}" # Загружается из переменной окружения

kafka:
  brokers:
    - "kafka-1:9092"
  producer:
    acks: "all"

logger:
  level: "info"
  format: "json"
```
*   (Более подробный пример см. в исходном документе "Стандарты API...")

## 8. Стандарты Инфраструктурных Файлов и CI/CD
Стандарты для инфраструктурных файлов (Dockerfile, Docker Compose, Kubernetes манифесты) и CI/CD детально описаны в документе `project_deployment_standards.md`.

---
*Этот документ должен регулярно обновляться и дополняться по мере развития платформы.*
