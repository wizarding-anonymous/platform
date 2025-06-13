# Game Purchase and Library Update Workflow

This diagram illustrates the sequence of interactions between services when a user purchases a game and it's added to their library. This example assumes a successful payment flow, **including interaction with typical Russian payment gateways (like ЮKassa, Tinkoff, SberPay) where `PaymentSvc` might redirect the user or handle API calls, and subsequently process webhooks/callbacks from these providers. The flow also covers the necessary fiscalization steps as per 54-ФЗ.** A more complex Saga pattern might be used for robust error handling and rollbacks, as described in `project_database_structure.md`.

```mermaid
sequenceDiagram
    actor User
    participant ClientApp as Client Application
    participant APIGW as API Gateway
    participant CatalogSvc as Catalog Service
    participant OrderSvc as Order Service (Conceptual - could be part of Payment or a separate service for managing orders, carts)
    participant PaymentSvc as Payment Service
    participant KafkaBus as Kafka Message Bus
    participant LibrarySvc as Library Service
    participant NotificationSvc as Notification Service
    participant AnalyticsSvc as Analytics Service

    User->>ClientApp: Selects game and initiates purchase (e.g., adds to cart, proceeds to checkout)
    ClientApp->>APIGW: POST /api/v1/orders/initiate (payload: product_id, quantity, selected_payment_method_hint)
    APIGW->>OrderSvc: Forward POST /orders/initiate (payload)

    OrderSvc->>CatalogSvc: gRPC GetProductPrice(product_id, user_region, currency)
    CatalogSvc-->>OrderSvc: ProductPriceResponse (price, currency, promotions_applied)
    OrderSvc->>OrderSvc: Create Order record in DB (status: pending_payment, order_id, user_id, items, total_amount)
    OrderSvc-->>KafkaBus: Publish event `order.payment.initiation.requested.v1` (order_id, user_id, total_amount, currency, items_for_receipt, preferred_payment_method_id)
    OrderSvc-->>APIGW: HTTP 201 Created (order_id, status: "pending_payment")
    APIGW-->>ClientApp: HTTP 201 Created (order_id)

    ClientApp->>APIGW: POST /api/v1/payments/transactions/initiate (payload: order_id, payment_method_type_hint, success_url, fail_url)
    APIGW->>PaymentSvc: Forward POST /transactions/initiate (payload)

    participant PSP as Payment Provider (e.g., ЮKassa, Tinkoff)
    PaymentSvc->>PaymentSvc: Create Transaction record (status: pending_psp_redirect or processing)
    PaymentSvc->>PSP: Interact with Payment Provider (e.g., ЮKassa, Tinkoff, SberPay or other PSP)
    PSP-->>PaymentSvc: Respond with redirect_url or payment_status
    PaymentSvc-->>APIGW: HTTP 201 Created (transaction_id, status, redirect_url_if_any, psp_sdk_data_if_any)
    APIGW-->>ClientApp: HTTP 201 Created (transaction_id, redirect_url or psp_sdk_data)

    alt Payment requires redirect (e.g., 3DS, SBP QR)
        ClientApp->>User: Redirects to Payment Provider's page or displays QR code
        User->>User: Completes payment on Provider's site/app
        Note over User, PaymentSvc: Payment Provider sends Webhook to PaymentSvc (`POST /webhooks/{provider}`)
        PaymentSvc->>PaymentSvc: Handle Webhook: Validate signature, check idempotency, update transaction status to 'processing' or 'completed'
        PaymentSvc-->>KafkaBus: Publish event `payment.transaction.status.updated.v1` (transaction_id, order_id, user_id, new_status: "completed", amount)
    else Direct API-based payment (e.g., saved card, internal balance - not shown in detail)
        PaymentSvc->>PaymentSvc: Process payment directly
        PaymentSvc-->>KafkaBus: Publish event `payment.transaction.status.updated.v1` (transaction_id, order_id, user_id, new_status: "completed" or "failed", amount)
    end

    subgraph Asynchronous Processing Post-Successful-Payment
        KafkaBus-->>OrderSvc: Consume `payment.transaction.status.updated.v1` (if status is "completed")
        OrderSvc->>OrderSvc: Update Order status to 'completed' in DB
        OrderSvc-->>KafkaBus: Publish `order.processing.completed.v1` (order_id, user_id, items)  # Renamed for clarity

        KafkaBus-->>LibrarySvc: Consume `order.processing.completed.v1`
        LibrarySvc->>LibrarySvc: For each item in order: Add product to user's library (UserLibraryItem record)
        LibrarySvc-->>KafkaBus: Publish event `library.game.added.v1` (user_id, product_id, acquisition_type: "purchase", purchase_date) for each item

        KafkaBus-->>NotificationSvc: Consume `order.processing.completed.v1`
        NotificationSvc->>NotificationSvc: Prepare purchase confirmation message (template: `order_confirmation`)
        NotificationSvc->>NotificationSvc: Send email/push notification to User
        NotificationSvc-->>KafkaBus: Publish event `notification.message.status.updated.v1` (status: sent/failed)

        KafkaBus-->>PaymentSvc: Consume `payment.transaction.status.updated.v1` (if status is "completed", for fiscalization)
        PaymentSvc->>PaymentSvc: Initiate fiscal receipt generation (integrates with OFD)
        Note right of PaymentSvc: Fiscalization via ОФД (54-ФЗ)
        PaymentSvc-->>KafkaBus: Publish `payment.fiscal.receipt.status.updated.v1` (status: fiscalized_success/failed)

        KafkaBus-->>NotificationSvc: Consume `payment.fiscal.receipt.status.updated.v1` (if fiscalized_success)
        NotificationSvc->>NotificationSvc: Prepare fiscal receipt notification (template: `fiscal_receipt_delivery`)
        NotificationSvc->>NotificationSvc: Send email/in-app to User with receipt details/link
        NotificationSvc-->>KafkaBus: Publish event `notification.message.status.updated.v1`

        KafkaBus-->>AnalyticsSvc: Consume `payment.transaction.status.updated.v1` (if status is "completed")
        AnalyticsSvc->>AnalyticsSvc: Log purchase event for analytics (sales, revenue, etc.)

        KafkaBus-->>DeveloperSvc: Consume `payment.transaction.status.updated.v1` (if status is "completed")
        DeveloperSvc->>DeveloperSvc: Update developer's balance (or queue for aggregation)
        DeveloperSvc-->>KafkaBus: Publish `developer.balance.updated.v1` (developer_id, new_balance, transaction_item_id)
    end

    ClientApp-->>User: Display purchase success/failure message (based on immediate PaymentSvc response or async update via WebSocket)
    User->>ClientApp: (Later) Checks library
    ClientApp->>APIGW: GET /api/v1/library/items
    APIGW->>LibrarySvc: Forward GET /items
    LibrarySvc-->>APIGW: Library contents (including newly purchased game)
    APIGW-->>ClientApp: Library contents
    ClientApp-->>User: Displays updated library
```

This diagram shows a common flow. The "Order Service" is conceptual and its responsibilities might be distributed. The use of Kafka events for post-payment processing ensures loose coupling and resilience. **Fiscalization:** After successful payment confirmation, `PaymentSvc` initiates fiscal receipt generation by integrating with an accredited Operator Fiscal Data (ОФД), fulfilling 54-ФЗ requirements. Notifications about the fiscal receipt are then sent to the user. Developer balance update is also included.
```
