# Client Update Mechanism

**Версия:** 1.0
**Дата последнего обновления:** 2024-07-16

Данный документ описывает механизмы обновления для различных платформ, поддерживаемых клиентским приложением. Обзор поддерживаемых платформ см. в [CROSS_PLATFORM_SUPPORT.md](./CROSS_PLATFORM_SUPPORT.md).

The frontend client, built with Flutter, supports multiple platforms (iOS, Android, Windows, Linux, Web). The update mechanism is tailored for each platform to ensure a smooth user experience and timely delivery of new features and fixes.

## A. Mobile Platforms (iOS & Android)

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

## B. Desktop Platforms (Windows & Linux)

1.  **In-App Self-Updater (Primary for Windows, Secondary for Linux):**
    *   The desktop application will periodically (e.g., on startup, daily) check a dedicated update server (or an API endpoint, possibly versioned and managed by Download Service or a simple static JSON file on a CDN) for information about the latest available version for its platform.
    *   This version information endpoint will provide the version number, release notes URL, and download URL(s) for the installer/package.
    *   If an update is detected, the application will offer to download the new installer/package (e.g., .msi for Windows, .deb/.rpm/.AppImage for Linux) in the background.
    *   Upon successful and verified (e.g., SHA256 checksum) download, the user will be prompted to install the update. This may involve launching the installer and closing the current application.
    *   **Tooling Consideration:**
        *   **Windows:** Libraries like `flutter_distributor` (which can wrap tools like Inno Setup) or a custom solution leveraging `Squirrel.Windows` principles (background downloads, easy install).
        *   **Linux:** While direct self-update is possible, packaging for system package managers is often preferred. If self-update is primary, it might involve downloading an AppImage or a script to manage .deb/.rpm updates.
        *   **macOS:** macOS (x64, Нативная поддержка Apple Silicon (Universal Binaries) и Intel x64) - Similar to Windows, can use Sparkle framework or custom logic for updates. For App Store distribution, updates go through the App Store.
2.  **Manual Download:** Users can always download the latest version of the desktop application from the official platform website (e.g., from a "Downloads" page).
3.  **Linux Package Managers (Preferred for Linux):**
    *   For Linux, the primary update mechanism should ideally be through system package managers.
    *   The application will be packaged and distributed via:
        *   **Snapcraft Store (Snap):** Cross-distro, transactional updates.
        *   **Flathub (Flatpak):** Cross-distro, sandboxed applications.
        *   **APT/YUM Repositories:** Distribution via dedicated repositories for Debian/Ubuntu-based and Fedora/RHEL-based systems respectively.
    *   Updates are then handled by the system's package manager, providing a native experience. The in-app checker can notify users to run `sudo apt update` or check their software center.

## C. Web Platform

1.  **Automatic Updates:** The web application, being composed of static assets (HTML, JS, CSS, images), updates automatically when users reload the page or open it in a new browser session. The web server/CDN will serve the latest deployed version of these assets.
2.  **Cache Management:**
    *   Cache-busting techniques (e.g., versioned asset filenames using content hashes, managed by the Flutter build process for JS/CSS) will be employed to ensure users' browsers fetch the latest assets and not serve stale content from cache.
    *   HTTP cache headers (ETag, Cache-Control) will be configured appropriately on the server/CDN.
3.  **Service Workers:**
    *   (Optional, for future enhancement) Service workers may be implemented to:
        *   Provide improved background updates: download new assets in the background and prompt the user to refresh when ready.
        *   Enhance offline capabilities for parts of the application.
        *   Offer faster load times by serving assets from cache first.

## D. Common Considerations

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
