<!-- project_workflows\project_workflows.md -->
# Пользовательские Сценарии и Рабочие Процессы

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16

## 1. Сценарии Пользователя Платформы

### 1.1. Регистрация и Первый Вход

```mermaid
sequenceDiagram
    actor User
    participant Client as Клиент
    participant Gateway as API Gateway
    participant Auth as Auth Service
    participant Account as Account Service
    participant Notification as Notification Service
    participant Kafka

    User->>Client: Нажимает "Регистрация"
    Client->>Gateway: POST /api/v1/auth/register
    Gateway->>Auth: Проксирование запроса
    
    Auth->>Auth: Валидация данных
    Auth->>Auth: Проверка уникальности email
    Auth->>Auth: Создание пользователя (статус: pending_verification)
    Auth->>Auth: Генерация кода верификации
    
    Auth-->>Kafka: Публикация user.registered.v1
    Auth-->>Gateway: 201 Created + tokens
    Gateway-->>Client: Ответ с токенами
    
    Kafka-->>Account: Событие user.registered.v1
    Account->>Account: Создание аккаунта и профиля
    
    Kafka-->>Notification: Событие user.registered.v1
    Notification->>User: Email с кодом верификации
    
    User->>Client: Вводит код верификации
    Client->>Gateway: POST /api/v1/auth/verify-email
    Gateway->>Auth: Проксирование
    Auth->>Auth: Проверка кода
    Auth->>Auth: Обновление статуса на active
    Auth-->>Kafka: Публикация user.email.verified.v1
    Auth-->>Gateway: 200 OK
    Gateway-->>Client: Успешная верификация
```

### 1.2. Покупка Игры

```mermaid
sequenceDiagram
    actor User
    participant Client
    participant Gateway as API Gateway
    participant Catalog as Catalog Service
    participant Payment as Payment Service
    participant Library as Library Service
    participant Download as Download Service
    participant Notification as Notification Service
    participant Kafka

    User->>Client: Просмотр страницы игры
    Client->>Gateway: GET /api/v1/catalog/games/{id}
    Gateway->>Catalog: Получение данных игры
    Catalog-->>Gateway: Данные игры с ценой
    Gateway-->>Client: Отображение информации
    
    User->>Client: Нажимает "Купить"
    Client->>Gateway: POST /api/v1/payment/transactions/initiate
    Gateway->>Payment: Инициализация платежа
    Payment->>Payment: Создание транзакции
    Payment->>Payment: Подготовка платежной сессии
    Payment-->>Gateway: URL платежной страницы
    Gateway-->>Client: Redirect на платежную систему
    
    User->>User: Оплата через платежную систему
    
    Note over Payment: Webhook от платежной системы
    Payment->>Payment: Обновление статуса транзакции
    Payment-->>Kafka: transaction.completed.v1
    
    Kafka-->>Library: Событие покупки
    Library->>Library: Добавление игры в библиотеку
    Library-->>Kafka: game.added.to.library.v1
    
    Kafka-->>Download: Подготовка к загрузке
    Download->>Download: Создание записей для загрузки
    
    Kafka-->>Notification: Уведомление о покупке
    Notification->>User: Email с чеком и подтверждением
    
    Client->>Gateway: GET /api/v1/library/games
    Gateway->>Library: Запрос библиотеки
    Library-->>Gateway: Список игр включая новую
    Gateway-->>Client: Обновленная библиотека
```

### 1.3. Загрузка и Установка Игры

```mermaid
sequenceDiagram
    actor User
    participant Client
    participant Gateway as API Gateway
    participant Library as Library Service
    participant Download as Download Service
    participant CDN
    participant WebSocket as WebSocket Gateway

    User->>Client: Открывает библиотеку
    Client->>Gateway: GET /api/v1/library/games
    Gateway->>Library: Запрос списка игр
    Library-->>Gateway: Список с статусами установки
    Gateway-->>Client: Отображение библиотеки
    
    User->>Client: Нажимает "Установить"
    Client->>Gateway: POST /api/v1/download/tasks
    Gateway->>Download: Создание задачи загрузки
    Download->>Download: Проверка прав доступа
    Download->>Download: Получение манифеста файлов
    Download->>Download: Создание задачи в очереди
    Download-->>Gateway: ID задачи и начальный статус
    Gateway-->>Client: Подтверждение начала
    
    Client->>WebSocket: Подписка на обновления загрузки
    
    loop Процесс загрузки
        Download->>CDN: Запрос части файла
        CDN-->>Download: Данные файла
        Download->>Download: Сохранение и проверка
        Download-->>WebSocket: Прогресс загрузки
        WebSocket-->>Client: Обновление прогресса
        Client->>User: Отображение прогресса
    end
    
    Download->>Download: Проверка целостности
    Download->>Library: Обновление статуса установки
    Download-->>WebSocket: Загрузка завершена
    WebSocket-->>Client: Уведомление о завершении
    Client->>User: "Играть" вместо "Установить"
```

### 1.4. Социальное Взаимодействие

```mermaid
sequenceDiagram
    actor User1 as Пользователь 1
    actor User2 as Пользователь 2
    participant Client1 as Клиент 1
    participant Client2 as Клиент 2
    participant Gateway as API Gateway
    participant Social as Social Service
    participant Notification as Notification Service
    participant WebSocket
    participant Kafka

    User1->>Client1: Поиск друга
    Client1->>Gateway: GET /api/v1/social/users/search?q=nickname
    Gateway->>Social: Поиск пользователей
    Social-->>Gateway: Результаты поиска
    Gateway-->>Client1: Список пользователей
    
    User1->>Client1: Отправить запрос в друзья
    Client1->>Gateway: POST /api/v1/social/friends/requests
    Gateway->>Social: Создание запроса
    Social->>Social: Сохранение запроса
    Social-->>Kafka: friend.request.sent.v1
    Social-->>Gateway: Подтверждение
    Gateway-->>Client1: Запрос отправлен
    
    Kafka-->>Notification: Событие запроса
    Notification->>WebSocket: Push уведомление
    WebSocket-->>Client2: Новый запрос в друзья
    Notification->>User2: Email уведомление
    
    User2->>Client2: Принять запрос
    Client2->>Gateway: PUT /api/v1/social/friends/requests/{id}/accept
    Gateway->>Social: Принятие запроса
    Social->>Social: Создание связи друзей
    Social-->>Kafka: friend.request.accepted.v1
    Social-->>Gateway: Подтверждение
    Gateway-->>Client2: Запрос принят
    
    Kafka-->>Notification: Уведомление о принятии
    Notification->>WebSocket: Push уведомление
    WebSocket-->>Client1: Запрос принят
    
    User1->>Client1: Отправить сообщение
    Client1->>WebSocket: Отправка сообщения
    WebSocket->>Social: Обработка сообщения
    Social->>Social: Проверка дружбы
    Social->>Social: Сохранение сообщения
    Social-->>WebSocket: Доставка получателю
    WebSocket-->>Client2: Новое сообщение
    Client2->>User2: Отображение сообщения
```

## 2. Сценарии Разработчика

### 2.1. Регистрация как Разработчик

```mermaid
sequenceDiagram
    actor Dev as Разработчик
    participant Client
    participant Gateway as API Gateway
    participant Auth as Auth Service
    participant Developer as Developer Service
    participant Admin as Admin Service
    participant Notification
    participant Kafka

    Note over Dev: Уже зарегистрирован как пользователь
    
    Dev->>Client: Переход в "Портал разработчика"
    Client->>Gateway: GET /api/v1/developer/account
    Gateway->>Developer: Проверка аккаунта
    Developer-->>Gateway: 404 - Аккаунт не найден
    Gateway-->>Client: Предложение создать аккаунт
    
    Dev->>Client: Заполнение формы разработчика
    Note over Client: ИНН, юр. название, адрес, банк. реквизиты
    
    Client->>Gateway: POST /api/v1/developer/account
    Gateway->>Developer: Создание аккаунта
    Developer->>Developer: Валидация данных
    Developer->>Developer: Создание аккаунта (pending_verification)
    Developer-->>Kafka: developer.account.created.v1
    Developer-->>Gateway: Аккаунт создан
    Gateway-->>Client: Ожидание верификации
    
    Kafka-->>Admin: Новый аккаунт для проверки
    Admin->>Admin: Добавление в очередь модерации
    
    Kafka-->>Notification: Уведомление разработчику
    Notification->>Dev: Email о начале проверки
    
    Note over Admin: Модератор проверяет документы
    Admin->>Admin: Проверка ИНН, документов
    Admin->>Developer: PUT /api/v1/developer/accounts/{id}/verify
    Developer->>Developer: Обновление статуса на verified
    Developer-->>Kafka: developer.account.verified.v1
    
    Kafka-->>Notification: Уведомление об одобрении
    Notification->>Dev: Email - аккаунт подтвержден
```

### 2.2. Публикация Новой Игры

```mermaid
sequenceDiagram
    actor Dev as Разработчик
    participant Portal as Dev Portal
    participant Gateway as API Gateway
    participant Developer as Developer Service
    participant S3
    participant Admin as Admin Service
    participant Catalog as Catalog Service
    participant Kafka

    Dev->>Portal: Создание нового продукта
    Portal->>Gateway: POST /api/v1/developer/products
    Gateway->>Developer: Создание черновика
    Developer->>Developer: Создание ProductSubmission
    Developer-->>Gateway: ID черновика
    Gateway-->>Portal: Форма заполнения
    
    Dev->>Portal: Заполнение метаданных
    Note over Portal: Название, описание, жанры, системные требования
    Portal->>Gateway: PATCH /api/v1/developer/products/{id}
    Gateway->>Developer: Обновление метаданных
    Developer-->>Gateway: Сохранено
    
    Dev->>Portal: Загрузка медиа
    Portal->>Gateway: POST /api/v1/developer/products/{id}/media/upload
    Gateway->>Developer: Генерация presigned URLs
    Developer->>S3: Резервирование места
    Developer-->>Gateway: URLs для загрузки
    Gateway-->>Portal: Presigned URLs
    
    Portal->>S3: Прямая загрузка файлов
    S3-->>Portal: Подтверждение загрузки
    
    Portal->>Gateway: POST /api/v1/developer/products/{id}/media/confirm
    Gateway->>Developer: Подтверждение загрузки
    Developer->>Developer: Проверка и обработка медиа
    
    Dev->>Portal: Загрузка билдов игры
    Note over Portal: Аналогичный процесс для билдов
    
    Dev->>Portal: Отправка на модерацию
    Portal->>Gateway: POST /api/v1/developer/products/{id}/submit
    Gateway->>Developer: Отправка на проверку
    Developer->>Developer: Валидация комплектности
    Developer->>Developer: Смена статуса на submitted
    Developer-->>Kafka: product.submitted.for.review.v1
    Developer-->>Gateway: Отправлено
    Gateway-->>Portal: В очереди на проверку
    
    Kafka-->>Admin: Новый продукт для модерации
    
    Note over Admin: Процесс модерации
    Admin->>Admin: Проверка контента
    Admin->>Developer: POST /api/v1/developer/products/{id}/approve
    Developer->>Developer: Смена статуса на approved
    Developer-->>Kafka: product.approved.v1
    
    Kafka-->>Catalog: Создание продукта в каталоге
    Catalog->>Catalog: Импорт данных из Developer Service
    Catalog-->>Kafka: product.published.v1
    
    Kafka-->>Developer: Обновление статуса
    Developer->>Developer: Статус = published
    Developer-->>Portal: Игра опубликована
```

### 2.3. Получение Выплат

```mermaid
sequenceDiagram
    actor Dev as Разработчик
    participant Portal as Dev Portal
    participant Gateway as API Gateway
    participant Developer as Developer Service
    participant Analytics as Analytics Service
    participant Payment as Payment Service
    participant Bank as Банковская система
    participant Kafka

    Note over Developer: Ежемесячный процесс (1-5 число)
    
    Developer->>Analytics: Запрос данных о продажах
    Analytics-->>Developer: Продажи за предыдущий месяц
    
    Developer->>Developer: Расчет выплат
    Note over Developer: Сумма продаж * (100% - комиссия платформы)
    
    Developer->>Developer: Создание PayoutRequest
    Developer-->>Kafka: payout.request.created.v1
    
    Dev->>Portal: Просмотр финансов
    Portal->>Gateway: GET /api/v1/developer/payouts
    Gateway->>Developer: Список выплат
    Developer-->>Gateway: Данные включая новый запрос
    Gateway-->>Portal: Отображение
    
    Note over Developer: Автоматическое одобрение если > порога
    Developer->>Developer: Проверка порога (1000 руб)
    Developer->>Developer: Смена статуса на approved
    Developer-->>Kafka: payout.request.approved.v1
    
    Kafka-->>Payment: Обработка выплаты
    Payment->>Payment: Создание исходящего платежа
    Payment->>Bank: API запрос на перевод
    Bank-->>Payment: Подтверждение обработки
    
    Payment-->>Kafka: payout.processing.v1
    
    Note over Bank: 1-3 рабочих дня
    Bank->>Payment: Webhook - платеж выполнен
    Payment->>Payment: Обновление статуса
    Payment-->>Kafka: payout.completed.v1
    
    Kafka-->>Developer: Обновление статуса выплаты
    Developer->>Developer: Статус = completed
    
    Kafka-->>Notification: Уведомление разработчику
    Notification->>Dev: Email - выплата получена
```

## 3. Сценарии Администратора

### 3.1. Модерация Контента

```mermaid
sequenceDiagram
    actor Admin as Администратор
    actor User as Пользователь
    participant AdminPanel as Админ-панель
    participant Gateway as API Gateway
    participant AdminSvc as Admin Service
    participant Social as Social Service
    participant Account as Account Service
    participant Notification
    participant Kafka

    User->>User: Публикует оскорбительный отзыв
    User->>User: Другой пользователь жалуется
    
    Note over AdminSvc: Жалоба попадает в очередь
    
    Admin->>AdminPanel: Открывает очередь модерации
    AdminPanel->>Gateway: GET /api/v1/admin/moderation/queue
    Gateway->>AdminSvc: Запрос очереди
    AdminSvc-->>Gateway: Список элементов
    Gateway-->>AdminPanel: Отображение
    
    Admin->>AdminPanel: Открывает жалобу
    AdminPanel->>Gateway: GET /api/v1/admin/moderation/items/{id}
    Gateway->>AdminSvc: Детали жалобы
    AdminSvc->>Social: GET отзыв по ID
    Social-->>AdminSvc: Содержимое отзыва
    AdminSvc-->>Gateway: Полная информация
    Gateway-->>AdminPanel: Отображение контекста
    
    Admin->>AdminPanel: Решение - удалить отзыв
    AdminPanel->>Gateway: POST /api/v1/admin/moderation/actions
    Gateway->>AdminSvc: Применение действия
    
    AdminSvc->>Social: DELETE отзыв
    Social->>Social: Мягкое удаление
    Social-->>Kafka: review.deleted.v1
    
    AdminSvc->>Account: Применить предупреждение
    Account->>Account: Добавление strike
    Account-->>Kafka: user.warned.v1
    
    AdminSvc->>AdminSvc: Логирование действия
    AdminSvc-->>Gateway: Действие выполнено
    Gateway-->>AdminPanel: Подтверждение
    
    Kafka-->>Notification: События модерации
    Notification->>User: Email о нарушении
```

### 3.2. Управление Платформой

```mermaid
sequenceDiagram
    actor SuperAdmin as Супер-администратор
    participant AdminPanel as Админ-панель
    participant Gateway as API Gateway
    participant AdminSvc as Admin Service
    participant ConfigSvc as Config Service
    participant AllServices as Все сервисы
    participant Kafka

    SuperAdmin->>AdminPanel: Настройки платформы
    AdminPanel->>Gateway: GET /api/v1/admin/platform/settings
    Gateway->>AdminSvc: Запрос настроек
    AdminSvc->>ConfigSvc: Получить конфигурацию
    ConfigSvc-->>AdminSvc: Текущие настройки
    AdminSvc-->>Gateway: Настройки
    Gateway-->>AdminPanel: Форма настроек
    
    SuperAdmin->>AdminPanel: Изменить комиссию с 30% на 25%
    AdminPanel->>Gateway: PUT /api/v1/admin/platform/settings
    Gateway->>AdminSvc: Обновление настроек
    
    AdminSvc->>AdminSvc: Валидация изменений
    AdminSvc->>AdminSvc: Аудит изменения
    AdminSvc->>ConfigSvc: Обновить конфигурацию
    ConfigSvc->>ConfigSvc: Сохранение
    ConfigSvc-->>Kafka: config.updated.v1
    
    Kafka-->>AllServices: Событие обновления
    AllServices->>AllServices: Перезагрузка конфигурации
    
    AdminSvc-->>Gateway: Изменения применены
    Gateway-->>AdminPanel: Подтверждение
    
    SuperAdmin->>AdminPanel: Включить режим техобслуживания
    AdminPanel->>Gateway: POST /api/v1/admin/platform/maintenance
    Gateway->>AdminSvc: Активация режима
    AdminSvc-->>Kafka: maintenance.mode.enabled.v1
    
    Kafka-->>Gateway: Обновление режима
    Gateway->>Gateway: Отдача 503 для пользователей
    
    Kafka-->>Notification: Массовое уведомление
    Notification->>Notification: Рассылка всем активным
```

## 4. Интеграционные Сценарии

### 4.1. Обработка Платежа с Фискализацией

```mermaid
sequenceDiagram
    participant Payment as Payment Service
    participant YooKassa as ЮKassa
    participant OFD as ОФД (АТОЛ/Яндекс.ОФД)
    participant FNS as ФНС
    participant Notification
    participant User

    Payment->>YooKassa: Создание платежа с чеком
    Note over Payment: Включены данные для чека
    YooKassa->>YooKassa: Обработка платежа
    YooKassa->>OFD: Формирование чека
    OFD->>OFD: Фискализация
    OFD->>FNS: Отправка данных
    FNS-->>OFD: Подтверждение
    OFD-->>YooKassa: Чек сформирован
    YooKassa-->>Payment: Webhook - платеж успешен
    
    Payment->>Payment: Обновление транзакции
    Payment->>Notification: Отправить чек пользователю
    Notification->>User: Email с чеком
    
    Note over OFD: Хранение чека 5 лет
```

### 4.2. Синхронизация Игровых Сохранений

```mermaid
sequenceDiagram
    participant Game as Игра
    participant Client as Клиент платформы
    participant Gateway as API Gateway
    participant Library as Library Service
    participant S3
    participant Sync as Sync Worker

    Game->>Client: Запрос на сохранение
    Client->>Gateway: PUT /api/v1/library/saves/{game_id}
    Gateway->>Library: Загрузка сохранения
    Library->>Library: Проверка квоты (100MB)
    Library->>S3: Сохранение файла
    S3-->>Library: URL сохранения
    Library->>Library: Версионирование
    Library-->>Gateway: Подтверждение
    Gateway-->>Client: Сохранено
    Client-->>Game: Успех
    
    Note over User: Запуск на другом устройстве
    
    Game->>Client: Запрос сохранений
    Client->>Gateway: GET /api/v1/library/saves/{game_id}
    Gateway->>Library: Получить последнее
    Library->>S3: Запрос файла
    S3-->>Library: Данные сохранения
    Library-->>Gateway: Сохранение
    Gateway-->>Client: Файл
    Client-->>Game: Восстановление
```

## 5. Аварийные Сценарии

### 5.1. Недоступность Payment Service

```mermaid
sequenceDiagram
    actor User
    participant Client
    participant Gateway as API Gateway
    participant Payment as Payment Service ❌
    participant Circuit as Circuit Breaker
    participant Cache as Redis Cache
    participant Queue as Kafka
    participant Notification

    User->>Client: Попытка покупки
    Client->>Gateway: POST /api/v1/payment/transactions
    Gateway->>Circuit: Проверка статуса Payment
    Circuit-->>Gateway: Сервис недоступен (Open)
    
    Gateway->>Gateway: Fallback стратегия
    Gateway->>Cache: Сохранить намерение покупки
    Cache-->>Gateway: Сохранено с TTL 24h
    
    Gateway->>Queue: Отложенная транзакция
    Queue-->>Gateway: В очереди
    
    Gateway-->>Client: 202 Accepted
    Client->>User: "Обработка платежа задерживается"
    
    Gateway->>Notification: Уведомить пользователя
    Notification->>User: Email о задержке
    
    Note over Payment: Сервис восстановлен
    
    Circuit->>Payment: Health check
    Payment-->>Circuit: 200 OK
    Circuit->>Circuit: Переход в Closed
    
    Queue->>Payment: Обработка очереди
    Payment->>Payment: Выполнение транзакций
    Payment-->>Notification: Результаты
    Notification->>User: Платеж обработан
```

### 5.2. DDoS Атака

```mermaid
sequenceDiagram
    actor Attacker
    participant CDN as CloudFlare
    participant WAF
    participant Gateway as API Gateway
    participant RateLimit as Rate Limiter
    participant Monitoring
    participant Admin

    Attacker->>CDN: Массовые запросы
    CDN->>CDN: Обнаружение аномалии
    CDN->>WAF: Фильтрация трафика
    WAF->>WAF: Блокировка по паттернам
    
    Note over WAF: Пропущено 10% атаки
    
    WAF->>Gateway: Остаточный трафик
    Gateway->>RateLimit: Проверка лимитов
    RateLimit->>RateLimit: Превышение порога
    RateLimit-->>Gateway: 429 Too Many Requests
    
    RateLimit->>Monitoring: Алерт о нагрузке
    Monitoring->>Admin: Уведомление в Telegram
    
    Admin->>CDN: Включение Under Attack Mode
    CDN->>CDN: JavaScript Challenge
    CDN->>Attacker: Проверка браузера
    Attacker-->>CDN: Fail
    CDN-->>Attacker: Block
    
    Admin->>WAF: Обновление правил
    WAF->>WAF: Блокировка подсетей
    
    Note over System: Нормализация трафика
```

---
*Этот документ содержит основные пользовательские сценарии платформы. Дополнительные сценарии добавляются по мере развития функциональности.*