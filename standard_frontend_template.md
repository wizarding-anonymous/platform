<!-- standard_frontend_template.md -->
# Спецификация Frontend Приложения: [Название Приложения]

**Версия:** 1.0
**Дата последнего обновления:** [YYYY-MM-DD]
**Ответственная команда:** Frontend Team

## 1. Обзор Приложения (Overview)

### 1.1. Назначение и Роль
*   Основной клиентский интерфейс для взаимодействия пользователей с платформой "Российский Аналог Steam"
*   Обеспечивает доступ ко всем функциям платформы через удобный пользовательский интерфейс
*   Целевые платформы: Web, Android (API 21+), iOS (11.0+), Desktop (Windows x64, Linux x64, macOS x64 и Apple Silicon с нативной поддержкой)

### 1.2. Ключевые Функциональности
*   Аутентификация и регистрация пользователей
*   Просмотр и поиск игр в каталоге
*   Управление библиотекой игр
*   Социальные функции (профиль, друзья, чат, отзывы)
*   Процесс покупки и оплаты
*   Управление загрузками и установками
*   Настройки профиля и приложения
*   Уведомления (push, in-app)

### 1.3. Основные Технологии
*   **Язык:** Dart (версия 3.0+)
*   **Framework:** Flutter (версия 3.10+)
*   **Управление состоянием:** flutter_bloc / bloc
*   **Навигация:** go_router
*   **HTTP клиент:** dio
*   **Локальное хранилище:** hive/hive_flutter, shared_preferences, flutter_secure_storage
*   **DI:** get_it с injectable
*   **Сериализация:** json_serializable, freezed
*   **Тестирование:** flutter_test, bloc_test, mockito/mocktail

## 2. Архитектура Приложения (Application Architecture)

### 2.1. Общее Описание
*   Приложение следует принципам Clean Architecture для обеспечения разделения ответственности, тестируемости и масштабируемости
*   Основные слои: Presentation (UI), Domain (Business Logic), Data (Repository)
*   Строгое разделение между слоями через интерфейсы

### 2.2. Слои Приложения (Layers)

#### 2.2.1. Presentation Layer (Слой Представления)
*   **Ответственность:** Отображение UI, обработка пользовательского ввода, управление состоянием экранов
*   **Компоненты:**
    *   **Widgets:** Атомарные компоненты, составные виджеты, экраны
    *   **BLoCs/Cubits:** Управление состоянием, обработка UI событий
    *   **Навигация:** Определение маршрутов, guards, deep links
    *   **Темы:** Светлая/темная тема, адаптивные стили

#### 2.2.2. Domain Layer (Доменный Слой)
*   **Ответственность:** Бизнес-логика, не зависящая от UI и источников данных
*   **Компоненты:**
    *   **Entities:** Основные бизнес-модели (User, Game, Order и т.д.)
    *   **Use Cases:** Конкретные бизнес-операции
    *   **Repository Interfaces:** Абстракции для доступа к данным
    *   **Exceptions:** Доменные исключения

#### 2.2.3. Data Layer (Слой Данных)
*   **Ответственность:** Получение данных из различных источников
*   **Компоненты:**
    *   **Repository Implementations:** Реализация интерфейсов из Domain
    *   **Data Sources:** Remote (API), Local (БД, кэш)
    *   **Models/DTOs:** Модели для сериализации/десериализации
    *   **Mappers:** Преобразование между DTO и Entity

## 3. Управление Состоянием (State Management)

*   **Основное решение:** BLoC pattern через flutter_bloc
*   **Принципы:**
    *   Immutable состояния
    *   Единый поток данных
    *   Разделение UI и бизнес-логики
*   **Структура:**
    *   Events: Действия пользователя или системы
    *   States: Состояния UI
    *   BLoC/Cubit: Обработка событий и эмиссия состояний

## 4. Навигация (Routing)

*   **Решение:** go_router для декларативной навигации
*   **Структура маршрутов:**
    ```
    /                      - Главная/Каталог
    /games/:id            - Страница игры
    /library              - Библиотека
    /profile              - Профиль пользователя
    /profile/edit         - Редактирование профиля
    /cart                 - Корзина
    /checkout             - Оформление заказа
    /settings             - Настройки
    /auth/login           - Вход
    /auth/register        - Регистрация
    ```
*   **Route Guards:** Проверка аутентификации для защищенных маршрутов
*   **Deep Linking:** Поддержка для всех платформ

## 5. Взаимодействие с API (API Integrations)

### 5.1. HTTP Клиент
*   **Библиотека:** dio
*   **Базовый URL:** 
    *   Production: https://api.steamanalog.ru
    *   Staging: https://api.stage.steamanalog.ru
    *   Development: https://api.dev.steamanalog.ru
*   **Таймауты:**
    *   Connect: 30 секунд
    *   Receive: 60 секунд
*   **Interceptors:**
    *   AuthInterceptor: Добавление JWT токена
    *   RefreshTokenInterceptor: Обновление токенов
    *   LoggingInterceptor: Логирование в dev режиме
    *   ErrorInterceptor: Стандартизация ошибок

### 5.2. Обработка Ошибок API
*   **Стандартный формат ошибок:**
    ```json
    {
      "errors": [
        {
          "code": "ERROR_CODE",
          "title": "Заголовок ошибки",
          "detail": "Детальное описание",
          "source": { "pointer": "/data/attributes/field" }
        }
      ]
    }
    ```
*   **Обработка:** Парсинг ошибок, локализация сообщений, retry механизмы

### 5.3. WebSocket
*   **Назначение:** Real-time уведомления, чат, обновления статусов
*   **Библиотека:** web_socket_channel
*   **URL:** wss://ws.steamanalog.ru
*   **Переподключение:** Автоматическое с экспоненциальной задержкой

## 6. Управление Ресурсами (Asset Management)

*   **Структура:**
    ```
    assets/
    ├── images/
    │   ├── 1x/
    │   ├── 2x/
    │   └── 3x/
    ├── icons/
    ├── fonts/
    └── animations/
    ```
*   **Оптимизация:** WebP для изображений, SVG для иконок
*   **Lazy Loading:** Для тяжелых ресурсов

## 7. Локализация (Localization - l10n)

*   **Поддерживаемые языки:**
    *   Русский (ru) - основной
    *   Английский (en)
*   **Инструменты:** intl package с ARB файлами
*   **Процесс:** Автогенерация через build_runner
*   **Fallback:** ru -> en

## 8. Сборка Приложения (Build Process)

### 8.1. Окружения (Environments)
*   **Development:**
    *   API: https://api.dev.steamanalog.ru
    *   Включен debug режим
    *   Расширенное логирование
*   **Staging:**
    *   API: https://api.stage.steamanalog.ru
    *   Профилирование производительности
*   **Production:**
    *   API: https://api.steamanalog.ru
    *   Минификация и обфускация
    *   Отключено логирование

### 8.2. Сборка для Различных Платформ
*   **Web:**
    ```bash
    flutter build web --release --dart-define=ENV=prod --web-renderer canvaskit
    ```
*   **Android:**
    ```bash
    flutter build appbundle --release --flavor prod --obfuscate --split-debug-info=build/symbols
    ```
*   **iOS:**
    ```bash
    flutter build ios --release --flavor prod --obfuscate --split-debug-info=build/symbols
    ```
*   **Desktop:**
    ```bash
    flutter build [windows|linux|macos] --release --dart-define=ENV=prod
    ```

## 9. Тестирование (Testing)

### 9.1. Типы Тестов
*   **Unit Tests:** Покрытие бизнес-логики (Use Cases, Repositories)
*   **Widget Tests:** Тестирование UI компонентов
*   **Integration Tests:** E2E тесты критических сценариев
*   **Golden Tests:** Скриншот тесты для UI регрессии

### 9.2. Покрытие
*   Минимальное покрытие: 80%
*   Критические пути: 100%
*   UI компоненты: 70%

### 9.3. CI/CD
*   Автоматический запуск тестов при PR
*   Блокировка merge при падении тестов
*   Отчеты о покрытии в PR

## 10. Производительность (Performance)

### 10.1. Метрики
*   First Meaningful Paint: < 2 сек
*   Time to Interactive: < 3 сек
*   Frame Rate: 60 FPS (минимум 30 FPS)
*   Memory Usage: < 150 MB в покое

### 10.2. Оптимизации
*   Lazy loading для тяжелых экранов
*   Image caching с ограничением размера
*   Pagination для списков
*   Debouncing для поиска
*   Memoization для вычислений

## 11. Безопасность (Security)

*   **Хранение токенов:** flutter_secure_storage
*   **Обфускация кода:** В production сборках
*   **Certificate Pinning:** Для критических API
*   **Биометрическая аутентификация:** Опционально
*   **Шифрование локальных данных:** Для чувствительной информации

## 12. Аналитика и Мониторинг

*   **Crash Reporting:** Sentry (self-hosted)
*   **Аналитика:** AppMetrica
*   **Performance Monitoring:** Flutter DevTools в dev
*   **Логирование:** Структурированное с уровнями

## 13. CI/CD и Развертывание

### 13.1. CI Pipeline
1. Линтинг (dart analyze)
2. Форматирование (dart format --set-exit-if-changed)
3. Unit/Widget тесты
4. Integration тесты (на эмуляторах)
5. Сборка артефактов
6. Загрузка в тестовые треки

### 13.2. Дистрибуция
*   **Web:** CDN (CloudFlare)
*   **Android:** RuStore, APK прямая загрузка
*   **iOS:** App Store
*   **Desktop:** Прямая загрузка с сайта

## 14. Поддержка и Обновления

*   **Версионирование:** SemVer
*   **Минимальная поддерживаемая версия:** Определяется API
*   **Force Update:** При критических изменениях
*   **Hot Reload:** Для критических исправлений (через CodePush аналоги)

---
*Этот шаблон является основой для документации frontend приложений. Каждый раздел должен быть адаптирован под специфику конкретного приложения.*