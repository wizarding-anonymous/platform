# Game and Game Component Update Mechanism

The platform provides a robust mechanism for updating games, DLCs, and other game components, aiming for efficiency and reliability. This process involves coordination between the Developer Service, Catalog Service, Download Service, Library Service, and the client application.

## A. Update Submission by Developers

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

## B. Client-Side Update Detection and Process

1.  **Automatic Update Check:**
    *   The user's client application, upon startup and periodically in the background (e.g., every few hours), checks for updates for all installed games.
    *   This check involves the client querying the Library Service (e.g., `GET /api/v1/library/me/items/updates-check`) with a list of installed product IDs and their current versions.
    *   Library Service, in turn, may query Catalog Service or Download Service (or have this information pushed/cached) to determine the latest available official version for each product.
2.  **Manual Update Check:** Users can also manually trigger an update check for specific games or all games from their library interface within the client application.
3.  **Update Notification:**
    *   If an update is available for one or more games, the client application will notify the user.
    *   The notification will typically include the game name, current version, new version, estimated download size (especially differentiating between delta patch and full download), and a link to release notes.
    *   Users are provided options to: update immediately, add to download queue, schedule the update (if download scheduling is supported), or view more details.

## C. Downloading and Applying Updates

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

## D. DLC and Game Component Updates

*   Updates for DLCs and other game components (e.g., optional high-resolution texture packs) are managed similarly to base game updates.
*   Each DLC or component can have its own versioning and manifest.
*   The client application checks for updates for installed DLCs alongside the base game.
*   Download and patching mechanisms are identical.

## E. Rollback to Previous Versions

*   **Standard Policy:** The default behavior is to always update users to the latest approved version of a game or component. Direct user-initiated rollback to a previous version is generally **not** supported as a standard feature due to complexity in managing save game compatibility, server-side compatibility (for online games), and support.
*   **Developer/Admin Initiated Rollback (Emergency):** In rare cases of a critically flawed update, developers (via Developer Service) or platform administrators (via Admin Service) can "roll back" an update by:
    1.  Unpublishing the faulty version from the Catalog Service.
    2.  Re-publishing a previous stable version as the current "latest" version.
    *   Clients would then see the re-published older version as an "update" if they had already installed the faulty one, or simply get the stable version if they hadn't updated yet. This is not a true client-side rollback but a change in the official latest version.
*   **Optional Client-Side Feature (Future Consideration):** As noted in the `Download Service` documentation, an *optional* feature for users to select and download specific older versions might be explored in the far future if there's strong demand and the technical complexities (especially around save game compatibility and developer opt-in) can be addressed. This would require Download Service to retain access to older game builds via CDN/S3.

## F. Client Application's Role in Updates

*   **Update Checking:** Periodically queries Library/Catalog/Download services for available updates for installed content.
*   **User Interaction:** Notifies users of available updates, presents release notes, and allows users to manage when and how updates are downloaded and installed.
*   **Download Coordination:** Interacts with the Download Service to fetch game files or patches, respecting user settings like bandwidth limits.
*   **Patch Application:** Implements the logic to apply binary patches to existing game files.
*   **File Verification:** Verifies checksums of all downloaded and patched files against the manifest.
*   **Installation Management:** Handles the file system operations for replacing/updating game files.
*   **State Reporting:** Reports the currently installed version of games and DLCs to the Library Service.
