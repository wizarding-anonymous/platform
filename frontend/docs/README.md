<!-- frontend\docs\README.md -->
# Спецификация Frontend Приложения: Клиентское приложение платформы 'Российский Аналог Steam'

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16
**Ответственная команда:** Frontend Team

## 1. Обзор Приложения (Overview)

### 1.1. Назначение и Роль
*   Основной клиентский интерфейс для взаимодействия пользователей с платформой 'Российский Аналог Steam', предоставляющий доступ ко всем функциям платформы.
*   Целевые платформы: 
    *   Web (все современные браузеры)
    *   Android (минимальный API уровень: 21, Android 5.0+)
    *   iOS (минимальная версия: 11.0+)
    *   Desktop:
        *   Windows x64 (Windows 10+)
        *   Linux x64 (Ubuntu 20.04+, Fedora 33+, Debian 10+)
        *   macOS x64 и Apple Silicon с нативной поддержкой (macOS 10.14+)
*   Ключевые бизнес-задачи: 
    *   Регистрация и аутентификация пользователей
    *   Просмотр каталога игр и другого контента
    *   Управление личной библиотекой (покупка, загрузка, установка, запуск игр)
    *   Взаимодействие с социальными функциями (профиль, друзья, чаты, отзывы)
    *   Управление настройками аккаунта и приложения

### 1.2. Ключевые Функциональности
*   **Phase 1 (MVP):**
    *   Регистрация/вход (email/пароль)
    *   Просмотр каталога игр (базовый поиск, фильтры по жанрам)
    *   Страница деталей игры
    *   Управление библиотекой (просмотр купленных игр, загрузка, установка, запуск)
    *   Управление профилем пользователя (просмотр, базовое редактирование)
    *   Процесс покупки игры
    *   Список желаемого
*   **Phase 2:**
    *   OAuth вход через VK/Telegram
    *   Расширенный поиск и фильтры
    *   Социальные функции (просмотр профилей, добавление в друзья, чат)
    *   Уведомления (In-App, Push)
    *   Отзывы и рейтинги
*   **Phase 3:**
    *   Форумы и сообщества
    *   Достижения
    *   Облачные сохранения
    *   Поддержка (FAQ, тикеты)

### 1.3. Основные Технологии
*   **Язык:** Dart (версия 3.0+)
*   **Framework:** Flutter (версия 3.10+)
*   **Ключевые библиотеки и пакеты (согласно `../../../../PACKAGE_STANDARDIZATION.md` и `../../../../project_technology_stack.md`):**
    *   **Управление состоянием:** `flutter_bloc` / `bloc` (основное решение)
    *   **Навигация:** `go_router` (основное решение)
    *   **HTTP Клиент:** `dio` (основное решение)
    *   **Локальное хранилище:** 
        *   `hive` / `hive_flutter` (для структурированных данных)
        *   `shared_preferences` (для простых настроек)
        *   `flutter_secure_storage` (для токенов и чувствительных данных)
    *   **Сериализация JSON:** `json_serializable` с `build_runner`, `freezed`
    *   **Равенство и неизменяемость моделей:** `equatable` или `freezed`
    *   **Внедрение зависимостей (DI):** `get_it` с `injectable` для генерации кода
    *   **Тестирование:** `flutter_test`, `bloc_test`, `mockito` или `mocktail`, `integration_test`
*   Ссылка на `CODING_STANDARDS.md` для стандартов кодирования

## 2. Архитектура Приложения (Application Architecture)

### 2.1. Общее Описание
*   Приложение построено на принципах Clean Architecture для обеспечения разделения ответственности, тестируемости и масштабируемости
*   Строгое разделение на слои: Presentation (UI), Domain (Business Logic), Data (Repository)
*   Использование интерфейсов для инверсии зависимостей
*   Диаграмма архитектуры:
    ```mermaid
    graph TD
        A[Presentation Layer<br/>UI, Widgets, BLoCs] --> B[Domain Layer<br/>Entities, Use Cases, Repository Interfaces]
        B --> C[Data Layer<br/>Repository Implementations, Data Sources]
        C --> D[External<br/>API, Local Storage, Platform Services]
    ```

### 2.2. Слои Приложения (Layers)

#### 2.2.1. Presentation Layer (Слой Представления)
*   **Ответственность:** Отображение UI, обработка пользовательского ввода, взаимодействие с Domain Layer через BLoCs
*   **Компоненты:**
    *   **Widgets (Виджеты):** 
        *   Атомарные компоненты (кнопки, поля ввода, карточки)
        *   Составные виджеты (формы, списки, диалоги)
        *   Экраны (полноценные страницы приложения)
    *   **BLoCs/Cubits:** 
        *   Управление состоянием UI
        *   Обработка событий от пользователя
        *   Вызов Use Cases из Domain Layer
    *   **Навигация:** Определение маршрутов через go_router
    *   **Темы:** Светлая и темная тема, адаптивные стили для разных платформ

#### 2.2.2. Domain Layer (Доменный Слой)
*   **Ответственность:** Бизнес-правила, сущности, сценарии использования. Независим от UI и деталей источников данных
*   **Компоненты:**
    *   **Entities:** User, Game, Order, Review, Achievement, Profile и др.
    *   **Use Cases:** 
        *   AuthenticateUser
        *   GetGameCatalog
        *   PurchaseGame
        *   DownloadGame
        *   UpdateProfile
        *   SendMessage
    *   **Repository Interfaces:** Абстракции для доступа к данным
    *   **Exceptions:** Доменные исключения (NetworkException, AuthException и др.)

#### 2.2.3. Data Layer (Слой Данных)
*   **Ответственность:** Получение данных из различных источников, кэширование, маппинг
*   **Компоненты:**
    *   **Repository Implementations:** Реализация интерфейсов из Domain Layer
    *   **Data Sources:**
        *   Remote: REST API через dio
        *   Local: Hive для кэша, SharedPreferences для настроек
    *   **Models/DTOs:** Модели для сериализации/десериализации JSON
    *   **Mappers:** Преобразование между DTO и Entity

## 3. Управление Состоянием (State Management)

*   **Паттерн:** BLoC (Business Logic Component) через flutter_bloc
*   **Принципы:**
    *   Все состояния иммутабельны (используем freezed или equatable)
    *   Единственный источник правды для каждого экрана
    *   Четкое разделение Events, States и логики в BLoC/Cubit
*   **Структура BLoC:**
    ```dart
    // Events
    abstract class GameCatalogEvent {}
    class LoadGamesEvent extends GameCatalogEvent {}
    class SearchGamesEvent extends GameCatalogEvent { final String query; }
    
    // States
    abstract class GameCatalogState {}
    class GameCatalogInitial extends GameCatalogState {}
    class GameCatalogLoading extends GameCatalogState {}
    class GameCatalogLoaded extends GameCatalogState { final List<Game> games; }
    class GameCatalogError extends GameCatalogState { final String message; }
    
    // BLoC
    class GameCatalogBloc extends Bloc<GameCatalogEvent, GameCatalogState> {
      // Implementation
    }
    ```

## 4. Навигация (Routing)

*   **Библиотека:** go_router для декларативной навигации
*   **Структура маршрутов:**
    ```
    /                           - Главная/Каталог
    /games/:id                  - Страница игры
    /games/:id/reviews          - Отзывы игры
    /library                    - Библиотека пользователя
    /library/downloads          - Загрузки
    /profile                    - Профиль текущего пользователя
    /profile/:userId            - Профиль другого пользователя
    /profile/edit               - Редактирование профиля
    /friends                    - Список друзей
    /chat                       - Список чатов
    /chat/:chatId               - Конкретный чат
    /wishlist                   - Список желаемого
    /cart                       - Корзина
    /checkout                   - Оформление заказа
    /settings                   - Настройки приложения
    /settings/account           - Настройки аккаунта
    /settings/notifications     - Настройки уведомлений
    /auth/login                 - Вход
    /auth/register              - Регистрация
    /auth/forgot-password       - Восстановление пароля
    /support                    - Поддержка
    /support/faq                - FAQ
    /support/ticket/new         - Создание тикета
    ```
*   **Route Guards:** Проверка аутентификации для защищенных маршрутов
*   **Deep Linking:** Полная поддержка для всех платформ

## 5. Взаимодействие с API (API Integrations)

### 5.1. HTTP Клиент
*   **Библиотека:** dio
*   **Конфигурация:**
    *   **Базовые URL:**
        *   Production: `https://api.steamanalog.ru`
        *   Staging: `https://api.stage.steamanalog.ru`
        *   Development: `https://api.dev.steamanalog.ru`
    *   **Таймауты:**
        *   Connect timeout: 30 секунд
        *   Receive timeout: 60 секунд
        *   Send timeout: 30 секунд
*   **Interceptors:**
    *   **AuthInterceptor:** Добавление JWT токена в заголовок Authorization
    *   **RefreshTokenInterceptor:** Автоматическое обновление Access Token при 401
    *   **LoggingInterceptor:** Логирование запросов/ответов в dev режиме
    *   **ErrorInterceptor:** Стандартизация и обработка ошибок
    *   **RetryInterceptor:** Повторная попытка при сетевых ошибках

### 5.2. Обработка Ошибок API
*   **Стандартный формат ошибок (согласно project_api_standards.md):**
    ```json
    {
      "errors": [
        {
          "code": "VALIDATION_ERROR",
          "title": "Ошибка валидации",
          "detail": "Поле email имеет неверный формат",
          "source": { "pointer": "/data/attributes/email" }
        }
      ]
    }
    ```
*   **Обработка:**
    *   Парсинг ошибок в доменные исключения
    *   Локализация сообщений об ошибках
    *   Отображение через SnackBar или диалоги
    *   Retry механизмы для временных сбоев

### 5.3. WebSocket
*   **Назначение:** 
    *   Real-time уведомления (новые сообщения, достижения)
    *   Обновление статусов друзей (онлайн/оффлайн)
    *   Синхронизация игровых сессий
*   **Библиотека:** web_socket_channel
*   **URL:** `wss://ws.steamanalog.ru`
*   **Протокол:**
    ```json
    {
      "type": "message.new",
      "id": "uuid",
      "payload": { ... }
    }
    ```
*   **Переподключение:** Автоматическое с экспоненциальной задержкой (1с, 2с, 4с, 8с, макс 30с)

## 6. Управление Ресурсами (Asset Management)

*   **Структура папки assets:**
    ```
    assets/
    ├── images/
    │   ├── 1x/           # Обычное разрешение
    │   ├── 2x/           # Retina
    │   └── 3x/           # Super Retina
    ├── icons/
    │   └── svg/          # Векторные иконки
    ├── fonts/
    │   ├── Roboto/       # Основной шрифт
    │   └── RobotoMono/   # Моноширинный шрифт
    ├── animations/
    │   └── lottie/       # Lottie анимации
    └── translations/     # Файлы локализации
        ├── ru.arb
        └── en.arb
    ```
*   **Оптимизация:**
    *   WebP для растровых изображений
    *   SVG для иконок и векторной графики
    *   Lazy loading для больших изображений
    *   Кэширование через cached_network_image

## 7. Локализация (Localization - l10n)

*   **Поддерживаемые языки:**
    *   Русский (ru) - основной язык
    *   Английский (en) - дополнительный
*   **Инструменты:** 
    *   intl package для интернационализации
    *   ARB файлы для хранения переводов
    *   flutter_localizations для системных виджетов
*   **Процесс:**
    1. Добавление строк в ARB файлы
    2. Генерация кода через `flutter gen-l10n`
    3. Использование через `AppLocalizations.of(context)`
*   **Fallback:** ru -> en -> ключ перевода

## 8. Сборка Приложения (Build Process)

### 8.1. Окружения (Environments)
*   **Development:**
    *   API URL: `https://api.dev.steamanalog.ru`
    *   Включен debug режим
    *   Расширенное логирование
    *   DevTools активны
    *   Mock данные для тестирования
*   **Staging:**
    *   API URL: `https://api.stage.steamanalog.ru`
    *   Профилирование производительности
    *   Ограниченное логирование
    *   Тестовые платежные системы
*   **Production:**
    *   API URL: `https://api.steamanalog.ru`
    *   Минификация и обфускация кода
    *   Отключено логирование
    *   Реальные платежные системы
    *   Включена аналитика

### 8.2. Сборка для Различных Платформ
*   **Web:**
    ```bash
    flutter build web --release --dart-define=ENV=prod --web-renderer canvaskit
    ```
    Особенности: CanvasKit для лучшей производительности, поддержка PWA
*   **Android:**
    ```bash
    flutter build appbundle --release --flavor prod --obfuscate --split-debug-info=build/symbols
    ```
    Подпись: Использование upload keystore для Google Play
*   **iOS:**
    ```bash
    flutter build ios --release --flavor prod --obfuscate --split-debug-info=build/symbols
    ```
    Подпись: Через Xcode с production сертификатами
*   **Desktop (Windows):**
    ```bash
    flutter build windows --release --dart-define=ENV=prod
    ```
*   **Desktop (Linux):**
    ```bash
    flutter build linux --release --dart-define=ENV=prod
    ```
*   **Desktop (macOS):**
    ```bash
    flutter build macos --release --dart-define=ENV=prod
    ```

## 9. Тестирование (Testing)

### 9.1. Типы Тестов
*   **Unit Tests:** 
    *   Покрытие Use Cases, Repositories, Mappers
    *   Минимум 80% покрытия для бизнес-логики
*   **Widget Tests:** 
    *   Тестирование отдельных виджетов
    *   Проверка отображения и взаимодействия
*   **Integration Tests:** 
    *   E2E тесты критических путей пользователя
    *   Запуск на реальных устройствах/эмуляторах
*   **Golden Tests:** 
    *   Скриншот тесты для предотвращения UI регрессий
    *   Отдельные эталоны для разных платформ

### 9.2. Инструменты и Библиотеки
*   `flutter_test` - основной фреймворк
*   `bloc_test` - тестирование BLoCs
*   `mockito` или `mocktail` - создание моков
*   `integration_test` - интеграционные тесты
*   `golden_toolkit` - улучшенные golden тесты

### 9.3. Покрытие и Метрики
*   Общее покрытие: минимум 70%
*   Критические пути: 100%
*   Новый код: минимум 80%
*   Отчеты через lcov и codecov

## 10. Производительность (Performance)

### 10.1. Целевые Метрики
*   **Web:**
    *   First Contentful Paint: < 1.5 сек
    *   Time to Interactive: < 3 сек
    *   Lighthouse Score: > 90
*   **Mobile/Desktop:**
    *   Cold Start: < 2 сек
    *   Warm Start: < 0.5 сек
    *   Frame Rate: 60 FPS (минимум 30 FPS)
    *   Memory Usage: < 150 MB в покое
    *   Jank: < 1%

### 10.2. Оптимизации
*   **Lazy Loading:** 
    *   Отложенная загрузка экранов
    *   Постраничная загрузка списков
*   **Кэширование:**
    *   HTTP ответы через dio_cache_interceptor
    *   Изображения через cached_network_image
    *   Локальный кэш данных в Hive
*   **Оптимизация списков:**
    *   ListView.builder для больших списков
    *   Переиспользование виджетов
    *   Const конструкторы где возможно
*   **Оптимизация изображений:**
    *   Правильные размеры для разных экранов
    *   WebP формат
    *   Progressive JPEG для превью

## 11. Безопасность (Security)

*   **Хранение данных:**
    *   JWT токены: flutter_secure_storage
    *   Пользовательские данные: шифрование в Hive
    *   Пароли: никогда не хранятся локально
*   **Сетевая безопасность:**
    *   HTTPS для всех соединений
    *   Certificate Pinning для критических API
    *   Защита от MITM атак
*   **Защита кода:**
    *   Обфускация в production сборках
    *   Защита от reverse engineering
    *   Удаление debug информации
*   **Аутентификация:**
    *   Биометрия как дополнительный фактор
    *   Автоматический logout при неактивности
    *   Безопасное хранение refresh token

## 12. Аналитика и Мониторинг

*   **Crash Reporting:** 
    *   Sentry (self-hosted версия)
    *   Автоматическая отправка отчетов
    *   Символизация стек-трейсов
*   **Аналитика поведения:**
    *   AppMetrica от Яндекс
    *   События: просмотры, клики, покупки
    *   Воронки конверсии
*   **Performance Monitoring:**
    *   Flutter DevTools в development
    *   Custom метрики через Prometheus
    *   Мониторинг API latency
*   **Логирование:**
    *   Структурированные логи
    *   Уровни: debug, info, warning, error
    *   Ротация логов на устройстве

## 13. CI/CD и Развертывание

### 13.1. CI Pipeline
1. **Проверка кода:**
   - dart analyze - статический анализ
   - dart format --set-exit-if-changed - форматирование
2. **Тестирование:**
   - Unit тесты
   - Widget тесты
   - Golden тесты
3. **Сборка:**
   - Сборка для всех платформ
   - Генерация debug symbols
4. **Интеграционные тесты:**
   - Запуск на device farm
   - Smoke тесты основных сценариев
5. **Артефакты:**
   - Загрузка в артефактори
   - Версионирование

### 13.2. CD Pipeline
*   **Beta распространение:**
    *   Android: Internal Testing Track
    *   iOS: TestFlight
    *   Desktop: Прямые ссылки для тестеров
*   **Production релиз:**
    *   Поэтапный rollout (5% -> 25% -> 50% -> 100%)
    *   Мониторинг crash rate
    *   Возможность отката

### 13.3. Дистрибуция
*   **Web:** 
    *   Хостинг на S3 + CloudFront CDN
    *   Автоматический деплой при merge в main
*   **Android:** 
    *   RuStore - основной магазин
    *   APK файлы на сайте
    *   F-Droid (опционально)
*   **iOS:** 
    *   App Store
    *   Enterprise распространение для корпоративных клиентов
*   **Desktop:** 
    *   Прямая загрузка с официального сайта
    *   Авто-обновление через встроенный механизм

## 14. Поддержка и Обновления

*   **Версионирование:** 
    *   Semantic Versioning (MAJOR.MINOR.PATCH)
    *   Build number увеличивается автоматически
*   **Обратная совместимость:**
    *   API версии поддерживаются минимум 6 месяцев
    *   Graceful degradation для старых версий
*   **Force Update:**
    *   При критических изменениях безопасности
    *   При breaking changes в API
    *   Уведомление за 2 недели
*   **Механизм обновлений:**
    *   In-app уведомления о новых версиях
    *   Автообновление для desktop версий
    *   Постепенная миграция пользователей

## 15. Пользовательские Сценарии (User Flows)

### 15.1. Регистрация и Вход
1. Пользователь открывает приложение
2. Нажимает "Создать аккаунт"
3. Вводит email и пароль
4. Получает письмо с подтверждением
5. Подтверждает email
6. Попадает в каталог игр

### 15.2. Покупка Игры
1. Пользователь находит игру в каталоге
2. Открывает страницу игры
3. Нажимает "Купить"
4. Выбирает способ оплаты
5. Подтверждает покупку
6. Игра появляется в библиотеке

### 15.3. Загрузка и Установка
1. Пользователь открывает библиотеку
2. Выбирает купленную игру
3. Нажимает "Установить"
4. Выбирает путь установки
5. Отслеживает прогресс загрузки
6. Запускает игру после установки

## 16. Известные Ограничения

*   **Web версия:**
    *   Нет доступа к файловой системе для установки игр
    *   Ограничения WebGL для 3D превью
*   **iOS версия:**
    *   Покупки только через In-App Purchase
    *   Ограничения на загрузку больших файлов
*   **Общие:**
    *   Максимальный размер кэша: 500 MB
    *   Одновременно не более 3 загрузок

---
*Документ поддерживается командой Frontend разработки и обновляется при каждом значимом изменении архитектуры или функциональности.*