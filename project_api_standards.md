# Стандарты API, форматов данных, событий и конфигурационных файлов

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16

## 1. REST API Стандарты

### 1.1. Версионирование
*   Версия API указывается в URL: `/api/v1/resource`
*   Мажорная версия (v1, v2) изменяется при несовместимых изменениях
*   Минорные изменения (добавление полей) не требуют новой версии
*   Старые версии поддерживаются минимум 6 месяцев после выпуска новой

### 1.2. Формат URL
*   Использовать существительные во множественном числе: `/users`, `/games`
*   Вложенные ресурсы: `/games/{game_id}/reviews`
*   Фильтрация через query параметры: `/games?genre=action&platform=pc`
*   Действия как подресурсы для не-CRUD операций: `/users/{id}/reset-password`

### 1.3. HTTP Методы
*   **GET** - получение ресурса или коллекции
*   **POST** - создание нового ресурса
*   **PUT** - полное обновление ресурса
*   **PATCH** - частичное обновление ресурса
*   **DELETE** - удаление ресурса

### 1.4. Коды Состояния
*   **200 OK** - успешный GET, PUT, PATCH
*   **201 Created** - успешный POST
*   **204 No Content** - успешный DELETE
*   **400 Bad Request** - некорректный запрос
*   **401 Unauthorized** - требуется аутентификация
*   **403 Forbidden** - недостаточно прав
*   **404 Not Found** - ресурс не найден
*   **409 Conflict** - конфликт состояния
*   **422 Unprocessable Entity** - ошибка валидации
*   **429 Too Many Requests** - превышен лимит запросов
*   **500 Internal Server Error** - ошибка сервера
*   **503 Service Unavailable** - сервис недоступен

### 1.5. Формат Ответа

#### Успешный ответ (одиночный ресурс):
```json
{
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "type": "game",
    "attributes": {
      "title": "Название игры",
      "releaseDate": "2024-01-15",
      "price": 1999
    },
    "relationships": {
      "developer": {
        "data": { "type": "developer", "id": "456" }
      }
    }
  },
  "meta": {
    "timestamp": "2024-07-16T10:30:00Z"
  }
}
```

#### Успешный ответ (коллекция):
```json
{
  "data": [
    {
      "id": "123",
      "type": "game",
      "attributes": { ... }
    }
  ],
  "meta": {
    "totalItems": 150,
    "totalPages": 15,
    "currentPage": 1,
    "perPage": 10,
    "timestamp": "2024-07-16T10:30:00Z"
  },
  "links": {
    "self": "/api/v1/games?page=1",
    "first": "/api/v1/games?page=1",
    "last": "/api/v1/games?page=15",
    "next": "/api/v1/games?page=2"
  }
}
```

#### Формат ошибки:
```json
{
  "errors": [
    {
      "code": "VALIDATION_ERROR",
      "title": "Ошибка валидации",
      "detail": "Поле 'email' имеет неверный формат",
      "source": { 
        "pointer": "/data/attributes/email" 
      },
      "meta": {
        "field": "email",
        "rejectedValue": "invalid-email"
      }
    }
  ],
  "meta": {
    "timestamp": "2024-07-16T10:30:00Z",
    "traceId": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### 1.6. Пагинация
*   Query параметры: `page` и `per_page` (или `limit` и `offset`)
*   Максимум элементов на странице: 100
*   Значения по умолчанию: page=1, per_page=20
*   В ответе включать meta информацию и links

### 1.7. Сортировка
*   Query параметр: `sort`
*   Формат: `sort=field` (по возрастанию) или `sort=-field` (по убыванию)
*   Множественная сортировка: `sort=genre,-price`

### 1.8. Фильтрация
*   Простые фильтры: `?status=active&type=game`
*   Операторы сравнения: `?price[gte]=1000&price[lte]=5000`
*   Поиск по тексту: `?search=название`
*   Фильтр по массиву: `?tags[]=action&tags[]=multiplayer`

### 1.9. Частичные Ответы
*   Query параметр `fields` для выбора полей
*   Формат: `?fields=id,title,price`
*   Вложенные поля: `?fields=id,title,developer{id,name}`

### 1.10. Заголовки

#### Обязательные заголовки запроса:
*   `Accept: application/json`
*   `Content-Type: application/json` (для POST, PUT, PATCH)
*   `Authorization: Bearer <token>` (для защищенных эндпоинтов)

#### Стандартные заголовки ответа:
*   `Content-Type: application/json`
*   `X-Request-ID: <uuid>`
*   `X-RateLimit-Limit: 1000`
*   `X-RateLimit-Remaining: 999`
*   `X-RateLimit-Reset: 1620000000`

## 2. gRPC API Стандарты

### 2.1. Именование
*   Сервисы: PascalCase с суффиксом Service (например, `UserService`)
*   Методы: PascalCase (например, `GetUser`, `CreateOrder`)
*   Сообщения: PascalCase с суффиксами Request/Response
*   Поля: snake_case

### 2.2. Структура Proto файлов
```protobuf
syntax = "proto3";

package platform.servicename.v1;

option go_package = "github.com/platform/api/servicename/v1;servicev1";

import "google/protobuf/timestamp.proto";
import "google/api/annotations.proto";

service GameService {
  rpc GetGame(GetGameRequest) returns (GetGameResponse) {
    option (google.api.http) = {
      get: "/api/v1/games/{game_id}"
    };
  }
}

message GetGameRequest {
  string game_id = 1;
}

message GetGameResponse {
  Game game = 1;
}

message Game {
  string id = 1;
  string title = 2;
  google.protobuf.Timestamp created_at = 3;
}
```

### 2.3. Обработка Ошибок
*   Использовать стандартные gRPC коды
*   Добавлять детали через `google.rpc.Status`
*   Включать trace_id в metadata

### 2.4. Версионирование
*   Версия в package name: `platform.servicename.v1`
*   Новые версии в отдельных файлах
*   Обратная совместимость обязательна

## 3. События (Event Streaming)

### 3.1. Формат CloudEvents
```json
{
  "specversion": "1.0",
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "source": "platform.auth-service",
  "type": "com.platform.auth.user.registered.v1",
  "time": "2024-07-16T10:30:00Z",
  "datacontenttype": "application/json",
  "data": {
    "userId": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "registeredAt": "2024-07-16T10:30:00Z"
  },
  "traceid": "550e8400-e29b-41d4-a716-446655440000",
  "correlationid": "660e8400-e29b-41d4-a716-446655440000"
}
```

### 3.2. Именование Событий
*   Формат: `com.platform.[service].[entity].[action].v[version]`
*   Примеры:
    *   `com.platform.auth.user.registered.v1`
    *   `com.platform.payment.transaction.completed.v1`
    *   `com.platform.game.achievement.unlocked.v1`

### 3.3. Kafka Topics
*   Naming: `platform.[domain].[entity].events`
*   Примеры:
    *   `platform.auth.user.events`
    *   `platform.payment.transaction.events`
*   Partitioning по entity ID для порядка событий

### 3.4. Гарантии Доставки
*   At-least-once delivery
*   Идемпотентность на стороне получателя
*   Дедупликация по event ID

## 4. WebSocket Протокол

### 4.1. Формат Сообщений
```json
{
  "type": "message.send",
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-07-16T10:30:00Z",
  "payload": {
    "chatId": "123",
    "text": "Привет!"
  }
}
```

### 4.2. Типы Сообщений
*   Клиент → Сервер:
    *   `auth.token` - аутентификация
    *   `subscribe.channel` - подписка на канал
    *   `unsubscribe.channel` - отписка от канала
    *   `message.send` - отправка сообщения
*   Сервер → Клиент:
    *   `auth.success` / `auth.failed`
    *   `message.new` - новое сообщение
    *   `user.status` - изменение статуса пользователя
    *   `error` - ошибка

### 4.3. Heartbeat
*   Ping/Pong каждые 30 секунд
*   Отключение при отсутствии pong в течение 60 секунд

## 5. Форматы Данных

### 5.1. Дата и Время
*   Формат: ISO 8601 с временной зоной
*   Пример: `2024-07-16T10:30:00Z`
*   Всегда в UTC

### 5.2. UUID
*   Формат: UUID v4
*   Пример: `550e8400-e29b-41d4-a716-446655440000`
*   Lowercase, с дефисами

### 5.3. Денежные Значения
*   Хранить в копейках (minor units)
*   Тип: integer
*   Пример: 199900 (1999 рублей)

### 5.4. Локализация
*   Поля с локализацией как JSON объект
*   Ключи - ISO 639-1 коды языков
```json
{
  "title": {
    "ru": "Название",
    "en": "Title"
  }
}
```

### 5.5. Изображения
*   Абсолютные URL для CDN
*   Разные размеры в объекте:
```json
{
  "images": {
    "thumbnail": "https://cdn.example.com/thumb.jpg",
    "medium": "https://cdn.example.com/medium.jpg",
    "large": "https://cdn.example.com/large.jpg",
    "original": "https://cdn.example.com/original.jpg"
  }
}
```

## 6. Аутентификация и Авторизация

### 6.1. JWT Токены
*   Алгоритм: RS256
*   Access Token TTL: 15 минут
*   Refresh Token TTL: 30 дней
*   Структура claims:
```json
{
  "iss": "https://auth.steamanalog.ru",
  "sub": "user-id",
  "aud": ["https://api.steamanalog.ru"],
  "exp": 1620000000,
  "iat": 1619999000,
  "jti": "token-id",
  "roles": ["user", "developer"],
  "permissions": ["read:games", "write:reviews"]
}
```

### 6.2. API Keys
*   Для сервис-сервис взаимодействия
*   Формат: `pak_live_xxxxxxxxxxxxxxxx`
*   Передача в заголовке: `X-API-Key`

### 6.3. OAuth 2.0
*   Поддержка Authorization Code flow
*   PKCE обязателен для публичных клиентов
*   Scopes согласно бизнес-требованиям

## 7. Конфигурационные Файлы

### 7.1. Формат
*   Основной формат: YAML
*   Альтернатива: переменные окружения
*   Именование: `config.yaml`, `config.dev.yaml`

### 7.2. Структура
```yaml
server:
  host: 0.0.0.0
  port: ${SERVER_PORT:8080}
  
database:
  host: ${DB_HOST:localhost}
  port: ${DB_PORT:5432}
  name: ${DB_NAME:platform}
  user: ${DB_USER:postgres}
  password: ${DB_PASSWORD}
  
redis:
  url: ${REDIS_URL:redis://localhost:6379}
  
kafka:
  brokers: ${KAFKA_BROKERS:localhost:9092}
  
logging:
  level: ${LOG_LEVEL:info}
  format: json
  
metrics:
  enabled: true
  port: ${METRICS_PORT:9090}
```

### 7.3. Управление Секретами
*   Никогда не коммитить секреты
*   Использовать переменные окружения
*   В production - Kubernetes Secrets или Vault
*   Placeholder в примерах: `YOUR_SECRET_HERE`

## 8. Документация API

### 8.1. OpenAPI (Swagger)
*   Версия: OpenAPI 3.0
*   Автогенерация где возможно
*   Включать примеры запросов/ответов
*   Доступность по `/api/v1/docs`

### 8.2. gRPC Documentation
*   Комментарии в proto файлах
*   Генерация через protoc-gen-doc
*   Reflection API в dev окружении

## 9. Rate Limiting

### 9.1. Лимиты по умолчанию
*   Анонимные: 100 req/hour
*   Аутентифицированные: 1000 req/hour
*   По IP: 10000 req/hour

### 9.2. Заголовки
*   `X-RateLimit-Limit` - лимит запросов
*   `X-RateLimit-Remaining` - осталось запросов
*   `X-RateLimit-Reset` - время сброса (Unix timestamp)

### 9.3. Ответ при превышении
*   Status: 429 Too Many Requests
*   Body: стандартный формат ошибки
*   Заголовок `Retry-After`

## 10. CORS

### 10.1. Разрешенные источники
*   Production: `https://steamanalog.ru`, `https://www.steamanalog.ru`
*   Staging: `https://stage.steamanalog.ru`
*   Development: `http://localhost:*`

### 10.2. Заголовки
*   `Access-Control-Allow-Methods`: GET, POST, PUT, PATCH, DELETE, OPTIONS
*   `Access-Control-Allow-Headers`: Content-Type, Authorization, X-Request-ID
*   `Access-Control-Max-Age`: 86400

---
*Этот документ является обязательным к исполнению для всех разработчиков платформы. Любые отклонения должны быть согласованы с архитектурной командой.*