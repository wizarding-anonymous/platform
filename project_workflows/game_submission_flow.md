<!-- project_workflows\game_submission_flow.md -->
**Дата последнего обновления:** 2024-07-16

# Developer Submits a New Game for Moderation Workflow

This diagram illustrates the sequence of interactions when a game developer submits a new game for moderation.

```mermaid
sequenceDiagram
    actor Developer
    participant DevClientApp as Developer Portal (Client App)
    participant APIGW as API Gateway
    participant DevSvc as Developer Service
    participant CatalogSvc as Catalog Service
    participant KafkaBus as Kafka Message Bus
    participant AdminSvc as Admin Service
    participant NotificationSvc as Notification Service

    Developer->>DevClientApp: Fills game metadata, pricing, prepares game build & media assets
    Note over Developer, DevClientApp: Developer creates a new game entry in their portal.
    DevClientApp->>APIGW: POST /api/v1/developer/games (payload: initial title, product_type)
    APIGW->>DevSvc: Forward POST /games (payload)
    DevSvc->>DevSvc: Create Game record (status: 'draft')
    DevSvc-->>APIGW: HTTP 201 Created (game_id)
    APIGW-->>DevClientApp: HTTP 201 Created (game_id)

    Developer->>DevClientApp: Edits game metadata (descriptions, genres, tags, system reqs, etc.) for game_id
    DevClientApp->>APIGW: PUT /api/v1/developer/games/{game_id}/metadata (payload)
    APIGW->>DevSvc: Forward PUT /games/{game_id}/metadata (payload)
    DevSvc->>DevSvc: Update GameMetadata record
    DevSvc-->>APIGW: HTTP 200 OK
    APIGW-->>DevClientApp: HTTP 200 OK

    Developer->>DevClientApp: Manages pricing for game_id
    DevClientApp->>APIGW: PUT /api/v1/developer/games/{game_id}/pricing (payload)
    APIGW->>DevSvc: Forward PUT /games/{game_id}/pricing (payload)
    DevSvc->>DevSvc: Store/Update GamePricing information (can involve calls to CatalogSvc for validation or to reflect changes)
    Note over DevSvc,CatalogSvc: Pricing might be complex, potentially involving CatalogSvc for regional templates or validation.
    DevSvc-->>APIGW: HTTP 200 OK
    APIGW-->>DevClientApp: HTTP 200 OK

    Developer->>DevClientApp: Creates a new version (e.g., "1.0.0") for game_id
    DevClientApp->>APIGW: POST /api/v1/developer/games/{game_id}/versions (payload: version_name, changelog)
    APIGW->>DevSvc: Forward POST /games/{game_id}/versions (payload)
    DevSvc->>DevSvc: Create GameVersion record (status: 'uploading_build')
    DevSvc-->>APIGW: HTTP 201 Created (version_id)
    APIGW-->>DevClientApp: HTTP 201 Created (version_id)

    Developer->>DevClientApp: Uploads game build for game_id, version_id
    DevClientApp->>APIGW: POST /api/v1/developer/games/{game_id}/versions/{version_id}/builds/upload-url (payload: file_name, file_size, platform)
    APIGW->>DevSvc: Forward POST .../upload-url (payload)
    DevSvc->>DevSvc: Generate pre-signed S3 URL
    DevSvc-->>APIGW: HTTP 200 OK (upload_url, internal_build_id)
    APIGW-->>DevClientApp: HTTP 200 OK (upload_url, internal_build_id)

    ClientApp->>ClientApp: Uploads file directly to S3 using pre-signed URL
    Note over ClientApp, DevSvc: After S3 upload, client notifies Developer Service.
    DevClientApp->>APIGW: POST /api/v1/developer/games/{game_id}/versions/{version_id}/builds/upload-complete (payload: internal_build_id, s3_path, file_hash)
    APIGW->>DevSvc: Forward POST .../upload-complete (payload)
    DevSvc->>DevSvc: Validate upload (check hash, size), update GameVersion status to 'ready_for_review' (or 'processing_build' if server-side processing is needed)
    DevSvc-->>APIGW: HTTP 200 OK

    Developer->>DevClientApp: Uploads media assets (screenshots, trailers) for game_id
    DevClientApp->>APIGW: POST /api/v1/developer/games/{game_id}/media/upload-url (payload: file_name, media_type)
    APIGW->>DevSvc: Forward POST .../media/upload-url (payload)
    DevSvc-->>APIGW: HTTP 200 OK (upload_url, internal_media_id)
    APIGW-->>DevClientApp: HTTP 200 OK (upload_url, internal_media_id)
    ClientApp->>ClientApp: Uploads media to S3
    DevClientApp->>APIGW: POST /api/v1/developer/games/{game_id}/media/upload-complete (payload: internal_media_id, s3_path)
    APIGW->>DevSvc: Forward POST .../media/upload-complete (payload)
    DevSvc->>DevSvc: Link media to Game record.
    DevSvc-->>APIGW: HTTP 200 OK

    Developer->>DevClientApp: Reviews all information and clicks "Submit for Moderation" for version_id
    DevClientApp->>APIGW: POST /api/v1/developer/games/{game_id}/versions/{version_id}/submit-for-review
    APIGW->>DevSvc: Forward POST /submit-for-review

    DevSvc->>DevSvc: Validate all required data is present for the version (metadata, build, media, pricing)
    DevSvc->>DevSvc: Update GameVersion status to 'pending_moderation'
    DevSvc->>DevSvc: Update Game (product) status to 'in_review' if not already
    DevSvc-->>KafkaBus: Publish event `developer.game.submitted.v1` (developer_id, game_id, version_id, product_id_catalog (if known), submission_date)
    DevSvc-->>APIGW: HTTP 200 OK (Submission successful, status: 'pending_moderation')
    APIGW-->>DevClientApp: HTTP 200 OK
    DevClientApp-->>Developer: Display success message ("Игра отправлена на модерацию")

    subgraph Asynchronous Processing Post-Submission
        KafkaBus-->>AdminSvc: Consume event `developer.game.submitted.v1`
        AdminSvc->>AdminSvc: Create ModerationItem record for the game version
        AdminSvc->>AdminSvc: Assign to relevant moderation queue (e.g., "New Game Versions")
        Note over AdminSvc: Модераторы видят этот элемент в своей очереди. CatalogSvc может быть уведомлен AdminSvc или также слушать это событие для создания/обновления черновика в каталоге.

        KafkaBus-->>CatalogSvc: Consume event `developer.game.submitted.v1` (alternative or additional flow)
        CatalogSvc->>CatalogSvc: Create/Update Product record in 'in_review' or 'draft_for_review' status using submitted data.
        CatalogSvc-->>KafkaBus: Publish `catalog.product.draft.created.v1` or `catalog.product.updated.v1`

        KafkaBus-->>NotificationSvc: Consume event `developer.game.submitted.v1`
        NotificationSvc->>NotificationSvc: Prepare "submission received" notification (template `dev_game_submission_received`)
        NotificationSvc->>NotificationSvc: Send notification to Developer (email/in-app via Developer Portal)
        NotificationSvc-->>KafkaBus: Publish event `notification.message.status.updated.v1`
    end
```

This diagram outlines the submission process. File uploads (builds, media) are shown as a two-step process (request URL, then confirm upload), which is common for S3. The interaction with Catalog Service ensures that a draft or "in-review" version of the product is created or updated in the catalog system.
```
