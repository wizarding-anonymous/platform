<!-- PROJECT_ROADMAP.md -->
# Project Roadmap: Russian Steam Analog

This document outlines the strategic development roadmap for the Russian Steam analog platform. It details the planned phases, from the Minimum Viable Product (MVP) through subsequent major and minor releases, outlining the features, components, and update mechanisms at each stage.

---

## Phase 1: Minimum Viable Product (MVP)

The MVP focuses on delivering the core functionality of the platform, enabling users to register, log in, browse, purchase, download, and play games.

### 1. Components/Microservices Implemented in MVP:

*   **Auth Service (New - Core):**
    *   Handles user registration (email/password, initial status `pending_verification`). (DONE - MVP)
    *   Manages user authentication (login/logout via email/password). (DONE - MVP)
    *   Issues and validates JWT (Access Token RS256, Refresh Token). Manages Refresh Token lifecycle (storage, revocation via JTI blacklist in Redis). (DONE - MVP)
    *   Handles basic email verification process (generates code, verifies code). (DONE - MVP)
    *   Utilizes Go, PostgreSQL (users, refresh tokens, verification codes), Redis (JTI blacklist, session info if any). Key integration with API Gateway (token validation) and Account Service (user creation event). (DONE - MVP)
*   **Account Service (New - Core):**
    *   Manages basic user account information (UserID from Auth Service, platform-specific AccountID, status like `active`, `inactive`).
    *   Stores minimal user profile data (e.g., nickname - must be unique, avatar URL - basic).
    *   Consumes `com.platform.auth.user.registered.v1` event to create account records.
    *   Utilizes Go, PostgreSQL (accounts, profiles), Redis (profile cache).
*   **API Gateway (New - Core Infrastructure):**
    *   Single entry point for all client requests (REST, potentially WebSocket later).
    *   Routes requests to appropriate backend services based on path, method, host.
    *   Handles initial JWT validation by calling Auth Service or using a shared JWKS endpoint. Injects `X-User-ID`, `X-User-Roles` headers.
    *   Manages CORS policies.
    *   Technology: Kong or Tyk, configured via Kubernetes CRDs (GitOps).
*   **Catalog Service (New - Basic):**
    *   Manages a list of available games with essential metadata: localized titles, descriptions, developer/publisher names (simplified text fields), main cover image URL, video trailer URL(s).
    *   Stores basic pricing information: base price in RUB. Regional pricing/currencies are **not** in MVP.
    *   Provides API for browsing (e.g., list of all games, paginated) and viewing game details.
    *   No advanced search/filtering (e.g., by tags, complex genre combinations) or extensive taxonomy for MVP. Elasticsearch integration is deferred.
    *   Utilizes Go, PostgreSQL (products, basic prices), Redis (cache for product details, lists).
*   **Payment Service (New - Basic):**
    *   Integrates with **one** selected Russian payment provider (e.g., YooKassa, Tinkoff Acquiring, or SBP via a specific bank gateway - decision pending).
    *   Handles the purchase transaction for a single game (no cart functionality for MVP).
    *   Processes payment provider callbacks to confirm payment status.
    *   Basic fiscalization according to 54-ФЗ: generation of a fiscal receipt for successful purchases via an integrated OFD.
    *   No complex features like refunds, subscriptions, multiple payment methods, or stored payment methods in MVP.
    *   Publishes `com.platform.payment.transaction.status.updated.v1` event upon successful payment.
    *   Utilizes Java/Kotlin (Spring Boot), PostgreSQL (transactions, basic fiscal receipt log).
*   **Library Service (New - Basic):**
    *   Consumes `com.platform.payment.transaction.status.updated.v1` (where status is 'completed' for a 'purchase' transaction type) to add a game to the user's library.
    *   Tracks games owned by the user (UserID, ProductID, acquisition_timestamp).
    *   Provides API for viewing the user's game library (list of ProductIDs).
    *   No playtime tracking, achievements, cloud saves, or user categories in MVP.
    *   Utilizes Go, PostgreSQL (user_library_items).
*   **Download Service (New - Basic):**
    *   Provides functionality to download purchased game builds.
    *   Requires `productId` and `versionId` (latest stable) and `userId` for authorization.
    *   Authorizes download requests by checking ownership via Library Service (internal gRPC call or event-based consistency).
    *   Generates secure, time-limited direct download links to game files stored on a CDN or S3-compatible storage.
    *   Provides a basic file manifest (list of files, sizes, hashes) for a given game version (obtained from Catalog Service or Developer Service via an event).
    *   No delta updates, no advanced download management (pause/resume client-side if possible via HTTP range requests, but service provides simple, full downloads).
    *   Utilizes Go, potentially S3-compatible storage for game builds (managed by Developer/Admin services, consumed by Download service).

### 2. User-Facing Functions Supported in MVP:

*   User registration using email and password; subsequent email verification. (DONE - MVP)
*   User login with email/password and logout. (DONE - MVP)
*   Viewing a list of available games in the catalog (basic listing, no advanced filtering).
*   Viewing detailed information for a specific game (description, price in RUB, main image/video).
*   Ability to initiate a purchase for a single game (no shopping cart).
*   Completing a game purchase using the single integrated Russian payment provider (e.g., YooKassa, Tinkoff, SBP).
*   Receiving an email confirmation with a fiscal receipt for the purchase.
*   Viewing purchased games in a personal library (simple list of owned games).
*   Downloading purchased games (basic download mechanism, full builds).
*   (Client-side responsibility) Launching downloaded games.
*   Viewing a basic user profile (nickname, avatar). Ability to change nickname and avatar.

### 3. Planned Updates/Integrations for MVP:

*   **Frontend Client (Flutter):** Basic client application, initially targeting Windows, to support the user-facing functions listed above.
*   **Payment Gateway Integration:** Final selection and technical integration of one Russian payment gateway (e.g., YooKassa API, Tinkoff Acquiring API, or SBP QR/redirect flow).
*   **OFD Integration:** Integration with one Russian OFD provider for fiscal receipt generation.
*   **CDN/S3 Storage:** Setup and integration for game build hosting and delivery.
*   **Initial Deployment & CI/CD:** Setup of CI/CD pipelines for MVP services and basic infrastructure. **Проект будет размещен на мощностях российского хостинг-провайдера Beget.** Облачные сервисы (S3-совместимые хранилища, базы данных, Kubernetes), если используются, будут от российских провайдеров (например, Yandex Cloud, VK Cloud, SberCloud), совместимых с инфраструктурой Beget или как часть гибридного решения.
*   **Security:** Foundational security measures: HTTPS for all endpoints, input validation, protection against common web vulnerabilities (OWASP Top 10 basics like XSS, SQLi prevention), secure JWT handling. (secure JWT handling DONE - MVP)
*   **Logging & Monitoring:** Basic centralized logging (e.g., ELK/Loki) and monitoring (e.g., Prometheus/Grafana) for all MVP services, focusing on uptime and error rates.
*   **Testing:** Unit tests for business logic (DONE - MVP), basic integration tests for service interactions (e.g., Payment -> Library event flow).

---

## Phase 2: Social Platform & Enhanced User Engagement (Major Release)

This major release focuses on building community features, improving user engagement with content, and introducing robust notification capabilities.

### Minor Release 2.1: Core Social Features & Profile Expansion

*   **Components/Microservices:**
    *   **Social Service (New - Core):**
        *   Manages user profiles (extending Account Service data): custom status messages, "About Me" section, profile background customization (predefined options or color).
        *   Basic friend system: send, accept, decline, remove friend requests; view friend list. Utilizes Neo4j for graph relationships.
        *   User presence: online, offline, in-game (basic status, updated via WebSocket or client polling). Stored in Redis.
        *   Technology: Go, PostgreSQL (profiles, pending requests), Neo4j (friend graph), Redis (presence).
    *   **Account Service (Enhancement):**
        *   API for Social Service to manage/display extended profile fields.
        *   Stores additional profile fields: custom status, "About Me" text, profile background selection.
    *   **API Gateway (Enhancement):** New routes for Social Service (e.g., `/api/v1/social/profiles/{userId}`, `/api/v1/social/friends`). WebSocket proxying for presence.
*   **User-Facing Functions Supported:**
    *   Viewing extended user profiles (own and others, respecting privacy settings).
    *   Setting a custom status message and "About Me" text.
    *   Customizing profile background (from predefined options).
    *   Sending, accepting, declining, and removing friend requests.
    *   Viewing friend lists and basic friend profiles with their online/in-game status.
*   **Planned Updates/Integrations:**
    *   Frontend client implementation for viewing extended profiles, managing friends, and updating new profile fields.
    *   Initial deployment of Social Service with PostgreSQL, Neo4j, and Redis.

### Minor Release 2.2: Wishlist & Basic Notifications

*   **Components/Microservices:**
    *   **Library Service (Enhancement):**
        *   Wishlist functionality: API to add/remove games from a personal wishlist, view wishlist. Data stored in PostgreSQL.
    *   **Notification Service (New - Basic):**
        *   Initial setup for email notifications via a selected **Russian provider** (e.g., SendPulse, Unisender, MailRu Cloud).
        *   Manages user notification preferences (opt-in/opt-out for basic categories like "Social Activity", "Wishlist Updates"). Stored in PostgreSQL.
        *   API for other services to request notification sending (or consumes Kafka events).
        *   Basic email templates for: new friend request, friend request accepted, game from wishlist is now on sale.
        *   Technology: Go, Kafka (for consuming requests), PostgreSQL (preferences, templates, basic logs).
    *   **Catalog Service (Enhancement):**
        *   Event `com.platform.catalog.price.updated.v1` to include discount information.
    *   **Social Service (Enhancement):**
        *   Publishes `com.platform.social.friend.request.sent.v1` and `com.platform.social.friend.request.accepted.v1` events.
*   **User-Facing Functions Supported:**
    *   Adding games to a personal wishlist.
    *   Removing games from the wishlist.
    *   Viewing the wishlist.
    *   Receiving email notifications for new friend requests, accepted requests, and when a wishlisted game goes on sale.
    *   Basic management of notification preferences (e.g., enable/disable "Wishlist Sale Alerts").
*   **Planned Updates/Integrations:**
    *   Frontend client implementation for wishlist management and notification preferences.
    *   Initial deployment of Notification Service.
    *   Integration: Library Service (wishlist) -> Kafka -> Notification Service (sale alerts). Social Service (friend requests) -> Kafka -> Notification Service.

### Minor Release 2.3: Basic Game Achievements & Activity Feed

*   **Components/Microservices:**
    *   **Library Service (Enhancement):**
        *   Basic achievement tracking: stores unlocked status and timestamp for predefined achievements per user per game.
        *   API for game clients (or via a trusted server call) to report achievement unlocking: `POST /me/achievements/unlock` (product_id, achievement_api_name).
        *   Publishes `com.platform.library.achievement.unlocked.v1` event.
    *   **Social Service (Enhancement):**
        *   Simple activity feed on user profiles: displays recent events like "User X unlocked achievement Y in Game Z", "User X became friends with User Y".
        *   Consumes `com.platform.library.achievement.unlocked.v1` and internal friend events.
        *   Stores feed items in Cassandra for scalability.
    *   **Catalog Service (Enhancement):**
        *   API for developers/admins to define a list of achievements for their games (achievement_api_name, localized name, description, icon_url_locked, icon_url_unlocked, is_hidden, sort_order). Stored in PostgreSQL.
*   **User-Facing Functions Supported:**
    *   Viewing unlocked achievements for owned games in the library and on user profiles.
    *   Viewing a simple activity feed on own and friends' profiles (respecting privacy settings).
*   **Developer-Facing Functions (via Admin Panel initially):**
    *   Admins can define achievements for games.
*   **Planned Updates/Integrations:**
    *   Frontend client display for achievements (in library and profiles) and activity feed.
    *   Mechanism for developers (via a temporary admin process or simplified Developer Portal interface if available by then) to submit basic achievement lists for their games.

---

## Phase 3: Developer Empowerment & Platform Operations (Major Release)

This major release focuses on providing tools for game developers to manage their content and for administrators to oversee the platform.

### Minor Release 3.1: Developer Portal Foundation & Basic Game Submission

*   **Components/Microservices:**
    *   **Developer Service (New - Core):**
        *   Developer/Publisher account registration and profile management (company name, contact info). Verification process handled by Admin Service initially.
        *   Product submission (game): input of core metadata (localized titles, descriptions, platforms, languages), genre/tag selection (from Catalog Service taxonomy), definition of system requirements, upload of main game images/videos (links to S3).
        *   No direct build uploads yet; placeholder for linking to externally hosted builds or specifying "coming soon".
        *   Publishes `com.platform.developer.product.submitted.v1` event.
        *   Technology: Go, PostgreSQL, S3 (for media assets).
    *   **Admin Service (New - Basic):**
        *   Interface to view developer account applications and submitted products.
        *   Workflow for developer account verification (manual process).
        *   Basic approval/rejection workflow for submitted game metadata (marks game as ready for Catalog Service to ingest if approved). Publishes `com.platform.admin.product.moderation.decision.v1`.
        *   Technology: Go, PostgreSQL (admin users, moderation tasks).
    *   **Catalog Service (Enhancement):**
        *   Consumes `com.platform.admin.product.moderation.decision.v1` (where decision is 'approved') to create or update product entries from Developer Service submissions.
    *   **Auth Service (Enhancement):**
        *   Support for developer account authentication (distinct roles like `developer_admin`, `developer_editor`).
*   **Developer-Facing Functions Supported:**
    *   Developers can register for a developer account.
    *   Developers can submit new games with core metadata, promotional media, and links to system requirements.
    *   Developers can view the status of their submissions (e.g., `draft`, `in_review`, `approved`, `rejected`).
*   **Admin-Facing Functions Supported:**
    *   Admins can review and approve/reject developer account applications.
    *   Admins can review submitted game metadata and approve/reject them.
*   **Planned Updates/Integrations:**
    *   Initial web-based frontend for Developer Portal (registration, company profile, product submission form).
    *   Initial web-based frontend for Admin Panel (developer verification queue, game moderation queue).
    *   Deployment of Developer Service and Admin Service.

### Minor Release 3.2: Game Build Management & Basic Analytics for Developers

*   **Components/Microservices:**
    *   **Developer Service (Enhancement):**
        *   Secure game build uploads (single file per platform initially) directly to platform's S3 storage via pre-signed URLs generated by Developer Service.
        *   Basic version control for game builds (e.g., list versions like "1.0.0_windows", "1.0.1_linux"; mark one as "active" for a given platform).
        *   Publishes `com.platform.developer.build.uploaded.v1` event.
    *   **Download Service (Enhancement):**
        *   Consumes `com.platform.developer.build.uploaded.v1` to register new builds and make them available for download (after Catalog Service marks product/version as live).
        *   Manages file manifests based on uploaded builds.
    *   **Analytics Service (New - Basic Data Collection & API):**
        *   Collects basic sales data per game (units sold, gross revenue) from `com.platform.payment.transaction.status.updated.v1` events.
        *   Stores aggregated data in ClickHouse.
        *   Provides a basic internal API for Developer Service to fetch sales data for a specific developer's products.
        *   Technology: Kafka consumers, ClickHouse, Go (for API).
    *   **Developer Service (Enhancement):**
        *   Displays basic sales analytics (units sold, gross revenue) for their own games in the Developer Portal (data fetched from Analytics Service).
*   **Developer-Facing Functions Supported:**
    *   Developers can upload game builds directly to the platform for specific versions and platforms.
    *   Developers can manage different versions of their game builds (e.g., set an active build).
    *   Developers can view basic sales figures (units sold, gross revenue) for their games.
*   **Planned Updates/Integrations:**
    *   Frontend enhancement for Developer Portal (build upload interface, version management, basic sales dashboard).
    *   Deployment of Analytics Service (basic data collection and API).
    *   Integration: Developer Service -> S3 & Download Service (builds). Payment Service -> Kafka -> Analytics Service (sales data). Analytics Service -> Developer Service (displaying analytics).

### Minor Release 3.3: User Management & Basic Support Tools for Admins

*   **Components/Microservices:**
    *   **Admin Service (Enhancement):**
        *   Platform User Management: Search users (by ID, email, nickname from Account Service), view basic user profile information, view user's game library (from Library Service), view user's transaction history (from Payment Service).
        *   Ability to issue warnings or temporarily/permanently ban users (updates status in Account Service, which Auth Service then respects). Publishes `com.platform.admin.user.status.updated.v1`.
        *   Basic platform statistics view (e.g., total users, total games sold, DAU - from Analytics Service).
    *   **Account Service (Enhancement):**
        *   API for Admin Service to update user status (e.g., `active`, `warned`, `banned_temp`, `banned_perm`). Publishes `com.platform.account.status.updated.v1`.
    *   **Auth Service (Enhancement):**
        *   Consumes `com.platform.account.status.updated.v1` to enforce bans (prevent login, revoke sessions).
    *   **Analytics Service (Enhancement):**
        *   Provide aggregated platform statistics (total users, DAU, total sales) to Admin Service.
    *   **Notification Service (Enhancement):**
        *   Consumes `com.platform.admin.user.status.updated.v1` (or a specific event from Admin Service) to notify users if they receive a warning or ban.
*   **Admin-Facing Functions Supported:**
    *   Admins can search for platform users and view their detailed profiles, library, and transaction history.
    *   Admins can issue warnings or ban users (temporarily or permanently) with a specified reason.
    *   Admins can view high-level platform operational statistics.
*   **Planned Updates/Integrations:**
    *   Frontend enhancement for Admin Panel (user management tools, user details view, statistics dashboard).
    *   Integration of Admin Service with Account Service (for user status updates) and Notification Service (for user notifications).

---

## Phase 4: Advanced Platform Capabilities & Monetization (Major Release)

This major release focuses on enhancing core platform services with advanced features, improving monetization options, and refining the user and developer experience.

### Minor Release 4.1: Advanced Catalog & Search, User Reviews, Basic Refunds

*   **Components/Microservices:**
    *   **Catalog Service (Enhancement):**
        *   Full Elasticsearch integration for advanced search (fuzzy search, faceted search, filtering by multiple criteria simultaneously, sorting options).
        *   Support for game bundles (fixed set of products sold together, possibly at a discount) and editions (e.g., Standard, Deluxe with extra DLCs).
    *   **Social Service (Enhancement):**
        *   User reviews and ratings for games: submit review (text + rating score like 1-5 stars or thumbs up/down), view reviews on product pages, sort/filter reviews. Stored in PostgreSQL within Social Service. Publishes `com.platform.social.review.submitted.v1`.
    *   **Payment Service (Enhancement):**
        *   Basic refund processing: API for Admin Service to initiate a full refund for a specific transaction. Payment Service attempts refund via the original payment provider. Publishes `com.platform.payment.transaction.status.updated.v1` (with status 'refunded' or 'refund_failed').
    *   **Admin Service (Enhancement):**
        *   Interface for initiating refunds for specific transactions, with reason input.
*   **User-Facing Functions Supported:**
    *   Users can use advanced search and filtering in the game catalog.
    *   Users can browse and purchase game bundles and different game editions.
    *   Users can submit textual reviews and ratings for games they own.
    *   Users can read, sort, and filter reviews on game pages.
*   **Admin-Facing Functions Supported:**
    *   Admins can initiate full refunds for purchases, providing a reason.
*   **Planned Updates/Integrations:**
    *   Frontend enhancements for advanced catalog browsing/search, bundle/edition purchasing, submitting and viewing reviews.
    *   Full Elasticsearch setup and data synchronization for Catalog Service.
    *   Admin panel interface for refund processing.
    *   Integration: Social Service (reviews) with Catalog Service (display on product pages).

### Minor Release 4.2: Enhanced Download Options & Developer Financials

*   **Components/Microservices:**
    *   **Download Service (Enhancement):**
        *   Delta updates for games: calculates and delivers only changed file parts if possible. Requires robust versioning and manifest comparison.
        *   Pre-loading support for upcoming releases (download encrypted builds before release date, decrypt on release).
        *   User-configurable bandwidth limiting options in the client, respected by Download Service.
    *   **Developer Service (Enhancement):**
        *   More detailed financial reporting: breakdown of sales by region (if regional pricing is introduced), net revenue after platform commission, estimated payout amounts.
        *   Tools for developers to suggest participation in platform-wide sales campaigns.
    *   **Payment Service (Enhancement):**
        *   Support for developer payouts: processing approved payout requests from Developer Service, integration with payment providers for mass payments or bank transfers to Russian entities.
        *   (Optional) Introduction of one additional payment method for users (e.g., another popular e-wallet or mobile payment if SBP/cards are already primary).
*   **User-Facing Functions Supported:**
    *   Faster game updates via delta patching where applicable.
    *   Ability to pre-load highly anticipated games before official release.
    *   User options in the client to control download bandwidth usage.
*   **Developer-Facing Functions Supported:**
    *   Developers can access more detailed financial reports on their sales and revenue.
    *   Developers can request to participate in upcoming sales events.
    *   Developers can manage their payout methods and see payout history.
*   **Planned Updates/Integrations:**
    *   Client-side logic for delta patching and pre-loading management.
    *   Developer portal enhancements for financial reporting and payout management.
    *   Payment Service integration with systems capable of payouts to Russian bank accounts/entities.

### Minor Release 4.3: Expanded Notifications, Community Features & Admin Analytics

*   **Components/Microservices:**
    *   **Notification Service (Enhancement):**
        *   In-app notifications (via WebSocket connection managed by API Gateway or a dedicated WebSocket service, Notification Service prepares content).
        *   User preferences for notification channels (email, push, in-app) and types (social, wishlist, game updates, marketing).
        *   Support for admin-driven marketing notifications (targeted to segments from Analytics Service).
    *   **Social Service (Enhancement):**
        *   Basic group functionality: create/join public groups, view group pages, post simple text announcements in groups. Stored in PostgreSQL.
        *   User mentions (`@nickname`) in comments and group posts, triggering notifications.
    *   **Analytics Service (Enhancement):**
        *   More comprehensive platform analytics for admins (e.g., user engagement metrics like session duration, feature usage, content interaction).
        *   Basic A/B testing framework support (data collection and segmentation side).
    *   **Admin Service (Enhancement):**
        *   Tools for creating and managing targeted marketing notification campaigns via Notification Service.
        *   Advanced analytics dashboards for platform performance and user behavior.
        *   Moderation tools for group names, descriptions, and admin-posted announcements within groups.
*   **User-Facing Functions Supported:**
    *   Users receive in-app notifications for important events.
    *   Users can finely tune their notification preferences across different channels and types.
    *   Users can create, join, and participate in basic public groups.
    *   Users can mention other users in comments and group posts.
*   **Admin-Facing Functions Supported:**
    *   Admins can send targeted marketing notifications to user segments.
    *   Admins have access to richer platform analytics and user behavior reports.
    *   Admins can moderate group names, descriptions, and announcements.
*   **Planned Updates/Integrations:**
    *   Client-side support for in-app notification center and preferences screen.
    *   Frontend for group functionality.
    *   Admin panel tools for marketing campaigns and advanced analytics views.
    *   Integration: Analytics Service (segments) -> Notification Service (targeting). Social Service (mentions) -> Notification Service.

---

## V. Client Update Mechanism

Детальное описание механизма обновления клиента см. в документе [CLIENT_UPDATE_MECHANISM.md](./CLIENT_UPDATE_MECHANISM.md).

---

## VI. Game and Game Component Update Mechanism

Детальное описание механизма обновления игр и их компонентов см. в документе [GAME_UPDATE_MECHANISM.md](./GAME_UPDATE_MECHANISM.md).
---
*Последнее обновление документа: 2024-07-16*
