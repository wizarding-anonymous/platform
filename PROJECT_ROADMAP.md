# Project Roadmap: Russian Steam Analog

This document outlines the strategic development roadmap for the Russian Steam analog platform. It details the planned phases, from the Minimum Viable Product (MVP) through subsequent major and minor releases, outlining the features, components, and update mechanisms at each stage.

---

## Phase 1: Minimum Viable Product (MVP)

The MVP focuses on delivering the core functionality of the platform, enabling users to register, log in, browse, purchase, download, and play games.

### 1. Components/Microservices Implemented in MVP:

*   **Auth Service (New - Core):**
    *   Handles user registration (email/password, initial status `pending_verification`).
    *   Manages user authentication (login/logout via email/password).
    *   Issues and validates JWT (Access Token RS256, Refresh Token). Manages Refresh Token lifecycle (storage, revocation via JTI blacklist in Redis).
    *   Handles basic email verification process (generates code, verifies code).
    *   Utilizes Go, PostgreSQL (users, refresh tokens, verification codes), Redis (JTI blacklist, session info if any). Key integration with API Gateway (token validation) and Account Service (user creation event).
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

*   User registration using email and password; subsequent email verification.
*   User login with email/password and logout.
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
*   **Initial Deployment & CI/CD:** Setup of CI/CD pipelines for MVP services and basic infrastructure on a chosen cloud provider located in Russia (e.g., Yandex Cloud, VK Cloud, SberCloud).
*   **Security:** Foundational security measures: HTTPS for all endpoints, input validation, protection against common web vulnerabilities (OWASP Top 10 basics like XSS, SQLi prevention), secure JWT handling.
*   **Logging & Monitoring:** Basic centralized logging (e.g., ELK/Loki) and monitoring (e.g., Prometheus/Grafana) for all MVP services, focusing on uptime and error rates.
*   **Testing:** Unit tests for business logic, basic integration tests for service interactions (e.g., Payment -> Library event flow).

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
        *   Initial setup for email notifications via a selected provider (e.g., SendGrid or a Russian alternative).
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
        *   Moderation tools for group content and user-generated content within groups.
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

The frontend client, built with Flutter, supports multiple platforms (iOS, Android, Windows, Linux, Web). The update mechanism is tailored for each platform to ensure a smooth user experience and timely delivery of new features and fixes.

### A. Mobile Platforms (iOS & Android)

1.  **Primary Distribution:** Updates are delivered through the official Apple App Store (iOS) and Google Play Store (Android). This is the standard and recommended approach for these platforms.
2.  **Process:**
    *   CI/CD pipelines (as defined in `project_deployment_standards.md`) automatically build, test, sign, and submit new app versions to the respective stores (e.g., using Fastlane or platform-specific APIs like Google Play Developer API).
    *   Following the store's review and approval process, the update becomes available to users. Rollout may be staged (e.g., percentage-based rollout in Google Play).
3.  **In-App Update Prompts:**
    *   **Android:** The application will utilize the `in_app_update` package to check for updates and prompt users for immediate or flexible updates directly within the app, as per Google Play policies.
    *   **iOS:** The application will periodically check a version endpoint (see Common Considerations) and, if a new version is available, display a custom, non-intrusive UI element notifying users and directing them to the App Store page for the update. Apple's guidelines generally discourage custom in-app update prompts that mimic system behavior.
4.  **Mandatory Updates:**
    *   For critical updates (e.g., security patches, major API compatibility changes), the application will implement a mechanism to enforce updates.
    *   This involves the client checking its version against a minimum required version fetched from a server endpoint.
    *   If the client version is below the mandatory level, it will display a blocking screen after a grace period, instructing the user to update from the store to continue using the application.

### B. Desktop Platforms (Windows & Linux)

1.  **In-App Self-Updater (Primary for Windows, Secondary for Linux):**
    *   The desktop application will periodically (e.g., on startup, daily) check a dedicated update server (or an API endpoint, possibly versioned and managed by Download Service or a simple static JSON file on a CDN) for information about the latest available version for its platform.
    *   This version information endpoint will provide the version number, release notes URL, and download URL(s) for the installer/package.
    *   If an update is detected, the application will offer to download the new installer/package (e.g., .msi for Windows, .deb/.rpm/.AppImage for Linux) in the background.
    *   Upon successful and verified (e.g., SHA256 checksum) download, the user will be prompted to install the update. This may involve launching the installer and closing the current application.
    *   **Tooling Consideration:**
        *   **Windows:** Libraries like `flutter_distributor` (which can wrap tools like Inno Setup) or a custom solution leveraging `Squirrel.Windows` principles (background downloads, easy install).
        *   **Linux:** While direct self-update is possible, packaging for system package managers is often preferred. If self-update is primary, it might involve downloading an AppImage or a script to manage .deb/.rpm updates.
2.  **Manual Download:** Users can always download the latest version of the desktop application from the official platform website (e.g., from a "Downloads" page).
3.  **Linux Package Managers (Preferred for Linux):**
    *   For Linux, the primary update mechanism should ideally be through system package managers.
    *   The application will be packaged and distributed via:
        *   **Snapcraft Store (Snap):** Cross-distro, transactional updates.
        *   **Flathub (Flatpak):** Cross-distro, sandboxed applications.
        *   **APT/YUM Repositories:** Distribution via dedicated repositories for Debian/Ubuntu-based and Fedora/RHEL-based systems respectively.
    *   Updates are then handled by the system's package manager, providing a native experience. The in-app checker can notify users to run `sudo apt update` or check their software center.

### C. Web Platform

1.  **Automatic Updates:** The web application, being composed of static assets (HTML, JS, CSS, images), updates automatically when users reload the page or open it in a new browser session. The web server/CDN will serve the latest deployed version of these assets.
2.  **Cache Management:**
    *   Cache-busting techniques (e.g., versioned asset filenames using content hashes, managed by the Flutter build process for JS/CSS) will be employed to ensure users' browsers fetch the latest assets and not serve stale content from cache.
    *   HTTP cache headers (ETag, Cache-Control) will be configured appropriately on the server/CDN.
3.  **Service Workers:**
    *   (Optional, for future enhancement) Service workers may be implemented to:
        *   Provide improved background updates: download new assets in the background and prompt the user to refresh when ready.
        *   Enhance offline capabilities for parts of the application.
        *   Offer faster load times by serving assets from cache first.

### D. Common Considerations

1.  **Versioning:**
    *   A clear semantic versioning scheme (SemVer: MAJOR.MINOR.PATCH, e.g., `1.0.0`, `1.1.0`, `1.1.1`) will be used for all client releases across all platforms.
    *   Build numbers or platform-specific version codes (e.g., Android `versionCode`, iOS `CFBundleVersion`) will also be incremented with each release.
2.  **Release Notes & Update Information:**
    *   Users will be informed of changes, new features, and bug fixes in each update.
    *   This information will be available through:
        *   In-app notifications or dialogs when an update is detected or applied.
        *   Store listings (for iOS and Android).
        *   A dedicated "What's New" or "Release Notes" section on the platform website or within the application.
    *   The update server/endpoint for desktop clients will also serve a link to the release notes for the latest version.
3.  **CI/CD Automation:**
    *   The build, testing, signing, and deployment/submission process for all client platforms will be automated as part of the CI/CD pipeline (defined in `project_deployment_standards.md`).
    *   This includes generating installers/packages for desktop, bundles/IPAs for mobile, and deploying static assets for web.
    *   Automated submission to app stores will be implemented where feasible (e.g., using Fastlane tools).
4.  **Rollback Strategy:**
    *   In case of critical issues discovered in a new client release:
        *   **Mobile:** Submit a patched version or a previous stable version to the app stores as quickly as possible. Utilize staged rollouts to catch issues before they affect all users.
        *   **Desktop (Self-Update):** Revert the 'latest version' information on the update server to point to a previous stable version. Users who haven't updated yet will not receive the faulty version. For users who have, a new "update" to the stable version can be pushed.
        *   **Desktop (Package Managers):** Publish a new package that reverts the changes or fixes the issue.
        *   **Web:** Re-deploy the previous stable version of static assets to the web server/CDN.
    *   **Server-Side Feature Flags:** Critical client-side features that depend on backend stability can be designed to be remotely disabled via server-side feature flags, mitigating the impact of certain types of client bugs without requiring an immediate client update.
    *   **Forced Update Revision:** The "mandatory update" mechanism can be used to quickly force users off a critically flawed version by setting its version as below the new minimum required version.

---

## VI. Game and Game Component Update Mechanism

The platform provides a robust mechanism for updating games, DLCs, and other game components, aiming for efficiency and reliability. This process involves coordination between the Developer Service, Catalog Service, Download Service, Library Service, and the client application.

### A. Update Submission by Developers

1.  **New Version Upload:** Developers and publishers use the Developer Portal (interfacing with the Developer Service) to submit new versions of their games or game components (e.g., DLCs, major patches). This includes specifying the version number (e.g., 1.1.0), release notes, and any changes to metadata.
2.  **Build and Manifest Upload:**
    *   For each new version, developers upload game builds for all supported platforms (e.g., Windows, Linux) to the S3-compatible storage designated by the platform. This upload process is managed by the Developer Service, which may provide pre-signed URLs for direct uploads.
    *   Developer Service also facilitates the creation or update of a detailed file manifest for the new version. This manifest includes a list of all files in the build, their sizes, and their SHA256 checksums. For delta updates, information about changed files or binary diffs might also be included or generated at this stage.
3.  **Update Notification & Moderation:**
    *   Upon successful upload and validation of a new version's builds and manifest, the Developer Service flags this version as ready.
    *   An event (e.g., `com.platform.developer.product.version.submitted.v1`) is published to Kafka.
    *   This event is consumed by the Admin Service to trigger a moderation process if required for updates (e.g., for significant changes or new DLCs).
    *   Once approved (or if no moderation is needed for a patch), Admin Service or Developer Service publishes an event (e.g., `com.platform.admin.product.version.approved.v1`).
    *   Catalog Service consumes this approval event. It updates its records with the new version information (version number, release notes, manifest location) and then publishes an event like `com.platform.catalog.product.version.published.v1`.
    *   Download Service consumes the `com.platform.catalog.product.version.published.v1` event to become aware of the new version, its manifest, and the location of its builds, preparing it for distribution.

### B. Client-Side Update Detection and Process

1.  **Automatic Update Check:**
    *   The user's client application, upon startup and periodically in the background (e.g., every few hours), checks for updates for all installed games.
    *   This check involves the client querying the Library Service (e.g., `GET /api/v1/library/me/items/updates-check`) with a list of installed product IDs and their current versions.
    *   Library Service, in turn, may query Catalog Service or Download Service (or have this information pushed/cached) to determine the latest available official version for each product.
2.  **Manual Update Check:** Users can also manually trigger an update check for specific games or all games from their library interface within the client application.
3.  **Update Notification:**
    *   If an update is available for one or more games, the client application will notify the user.
    *   The notification will typically include the game name, current version, new version, estimated download size (especially differentiating between delta patch and full download), and a link to release notes.
    *   Users are provided options to: update immediately, add to download queue, schedule the update (if download scheduling is supported), or view more details.

### C. Downloading and Applying Updates

1.  **Download Initiation & Management:**
    *   When the user initiates an update, the client application requests the download from the Download Service (e.g., `POST /api/v1/download/tasks` with product ID and target version ID).
    *   Download Service authorizes the request (checking ownership via Library Service) and then manages the download process, providing download URLs (typically CDN links) for game files or patches.
2.  **Delta Updates (Patches):**
    *   The platform prioritizes delta updates to minimize download sizes and times.
    *   Download Service, using the file manifests of the user's current version and the new version, determines if a delta update is possible and prepares/provides the patch data.
    *   The client application downloads the patch files.
    *   The client application is then responsible for applying the downloaded patch to the existing game files using a patching algorithm (e.g., bsdiff/bspatch, or a custom solution). This may involve creating temporary copies of files being patched.
3.  **Full Downloads:**
    *   If a delta update is not feasible (e.g., the version difference is too large, local game files are corrupted, or the developer did not provide necessary data for patching), or if a patch application fails, the system will fall back to downloading the full new version of the game or specific changed/corrupted files.
4.  **File Verification:**
    *   After downloading (either a patch or full files), the client application verifies the integrity of all downloaded and patched files using checksums (e.g., SHA256) provided in the file manifest for the new version.
    *   If verification fails for any file, the client will attempt to re-download that specific file/chunk or report an error.
5.  **Installation & Finalization:**
    *   The client application manages the final installation steps, which may include:
        *   Replacing old game files with new ones.
        *   Applying patches to existing files.
        *   Running any necessary post-update scripts or setup routines if provided by the game developer (e.g., updating configuration files, registering components).
    *   Once the update is successfully applied and verified, the client updates its local record of the installed game version and notifies the Library Service (e.g., `PATCH /api/v1/library/me/items/{itemId}` with new version info).

### D. DLC and Game Component Updates

*   Updates for DLCs and other game components (e.g., optional high-resolution texture packs) are managed similarly to base game updates.
*   Each DLC or component can have its own versioning and manifest.
*   The client application checks for updates for installed DLCs alongside the base game.
*   Download and patching mechanisms are identical.

### E. Rollback to Previous Versions

*   **Standard Policy:** The default behavior is to always update users to the latest approved version of a game or component. Direct user-initiated rollback to a previous version is generally **not** supported as a standard feature due to complexity in managing save game compatibility, server-side compatibility (for online games), and support.
*   **Developer/Admin Initiated Rollback (Emergency):** In rare cases of a critically flawed update, developers (via Developer Service) or platform administrators (via Admin Service) can "roll back" an update by:
    1.  Unpublishing the faulty version from the Catalog Service.
    2.  Re-publishing a previous stable version as the current "latest" version.
    *   Clients would then see the re-published older version as an "update" if they had already installed the faulty one, or simply get the stable version if they hadn't updated yet. This is not a true client-side rollback but a change in the official latest version.
*   **Optional Client-Side Feature (Future Consideration):** As noted in the `Download Service` documentation, an *optional* feature for users to select and download specific older versions might be explored in the far future if there's strong demand and the technical complexities (especially around save game compatibility and developer opt-in) can be addressed. This would require Download Service to retain access to older game builds via CDN/S3.

### F. Client Application's Role in Updates

*   **Update Checking:** Periodically queries Library/Catalog/Download services for available updates for installed content.
*   **User Interaction:** Notifies users of available updates, presents release notes, and allows users to manage when and how updates are downloaded and installed.
*   **Download Coordination:** Interacts with the Download Service to fetch game files or patches, respecting user settings like bandwidth limits.
*   **Patch Application:** Implements the logic to apply binary patches to existing game files.
*   **File Verification:** Verifies checksums of all downloaded and patched files against the manifest.
*   **Installation Management:** Handles the file system operations for replacing/updating game files.
*   **State Reporting:** Reports the currently installed version of games and DLCs to the Library Service.
```
