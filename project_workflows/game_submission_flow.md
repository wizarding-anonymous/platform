# Developer Submits a New Game for Moderation Workflow

This diagram illustrates the sequence of interactions when a game developer submits a new game for moderation.

```mermaid
sequenceDiagram
    actor Developer
    participant DevClientApp as Developer Portal (Client App)
    participant APIGW as API Gateway
    participant DevSvc as Developer Service
    participant CatalogSvc as Catalog Service
    participant UploadSvc as Upload Service (Conceptual - could be part of DevSvc or a separate service for large files)
    participant KafkaBus as Kafka Message Bus
    participant AdminSvc as Admin Service
    participant NotificationSvc as Notification Service

    Developer->>DevClientApp: Fills game metadata, pricing, uploads game build & media assets
    DevClientApp->>APIGW: POST /api/v1/developer/games (metadata, pricing info)
    APIGW->>DevSvc: Forward POST /games (payload)

    DevSvc->>DevSvc: Validate developer account & permissions
    DevSvc->>DevSvc: Create/Update GameProject record (status: draft_submission)
    DevSvc-->>APIGW: HTTP 201 Created (game_project_id)
    APIGW-->>DevClientApp: HTTP 201 Created (game_project_id)

    DevClientApp->>APIGW: POST /api/v1/developer/games/{game_project_id}/builds (build file via multipart/form-data or request for presigned URL)
    APIGW->>UploadSvc: Forward POST /builds (or to DevSvc if it handles uploads)
    UploadSvc->>UploadSvc: Store build file (e.g., to S3 staging area)
    UploadSvc-->>APIGW: HTTP 200 OK (build_id, stored_path)
    APIGW-->>DevClientApp: HTTP 200 OK

    DevClientApp->>APIGW: POST /api/v1/developer/games/{game_project_id}/media (media assets)
    APIGW->>UploadSvc: Forward POST /media
    UploadSvc->>UploadSvc: Store media files
    UploadSvc-->>APIGW: HTTP 200 OK (media_asset_ids)
    APIGW-->>DevClientApp: HTTP 200 OK

    Developer->>DevClientApp: Reviews all information and clicks "Submit for Moderation"
    DevClientApp->>APIGW: POST /api/v1/developer/games/{game_project_id}/submit-for-moderation
    APIGW->>DevSvc: Forward POST /submit-for-moderation

    DevSvc->>DevSvc: Validate all required data is present (metadata, build, media)
    DevSvc->>DevSvc: Update GameProject status to 'pending_moderation'
    DevSvc-->>CatalogSvc: gRPC CreateOrUpdateProductDraft(game_project_id, metadata, build_info, media_info)
    CatalogSvc->>CatalogSvc: Create/update Product record in 'in_review' or similar status. Store metadata.
    CatalogSvc-->>DevSvc: CreateOrUpdateProductDraftResponse (product_id_in_catalog)

    DevSvc-->>KafkaBus: Publish event `developer.game.submitted.v1` (developer_id, game_project_id, product_id_in_catalog, submission_date)
    DevSvc-->>APIGW: HTTP 200 OK (Submission successful)
    APIGW-->>DevClientApp: HTTP 200 OK
    DevClientApp-->>Developer: Display success message ("Game submitted for moderation")

    subgraph Asynchronous Processing Post-Submission
        KafkaBus-->>AdminSvc: Consume event `developer.game.submitted.v1`
        AdminSvc->>AdminSvc: Create ModerationItem record for the new game submission
        AdminSvc->>AdminSvc: Assign to relevant moderation queue (e.g., "New Game Submissions")
        Note over AdminSvc: Moderators can now see this item in their queue.

        KafkaBus-->>NotificationSvc: Consume event `developer.game.submitted.v1`
        NotificationSvc->>NotificationSvc: Prepare "submission received" notification for Developer
        NotificationSvc->>NotificationSvc: Send notification (email/in-app via Developer Portal)
        NotificationSvc-->>KafkaBus: Publish event `notification.sent.v1`
    end
```

This diagram outlines the submission process. The "Upload Service" is shown conceptually; file uploads might be handled directly by the Developer Service or a dedicated microservice that integrates with S3 or another object storage. The interaction with Catalog Service ensures that a draft or "in-review" version of the product is created in the catalog system.
```
