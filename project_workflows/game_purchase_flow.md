# Game Purchase and Library Update Workflow

This diagram illustrates the sequence of interactions between services when a user purchases a game and it's added to their library. This example assumes a successful payment flow. A more complex Saga pattern might be used for robust error handling and rollbacks, as described in `project_database_structure.md`.

```mermaid
sequenceDiagram
    actor User
    participant ClientApp as Client Application
    participant APIGW as API Gateway
    participant CatalogSvc as Catalog Service
    participant OrderSvc as Order Service (Conceptual - could be part of Payment or a separate service)
    participant PaymentSvc as Payment Service
    participant KafkaBus as Kafka Message Bus
    participant LibrarySvc as Library Service
    participant NotificationSvc as Notification Service

    User->>ClientApp: Selects game and initiates purchase
    ClientApp->>APIGW: POST /api/v1/orders (or similar to initiate purchase) (game_id, payment_method_details)
    APIGW->>OrderSvc: Forward POST /orders (payload)

    OrderSvc->>CatalogSvc: gRPC GetProductPrice(game_id, user_region)
    CatalogSvc-->>OrderSvc: ProductPriceResponse (price, currency)
    OrderSvc->>OrderSvc: Create Order record (status: pending_payment)
    OrderSvc->>PaymentSvc: gRPC InitiatePayment(order_id, amount, currency, payment_method_details, user_id)
    PaymentSvc-->>OrderSvc: InitiatePaymentResponse (transaction_id, status: processing/redirect_url)

    alt Payment requires redirect (e.g., 3DS)
        OrderSvc-->>APIGW: HTTP 200 OK (redirect_url_for_payment)
        APIGW-->>ClientApp: HTTP 200 OK (redirect_url_for_payment)
        ClientApp->>User: Redirects to Payment Provider
        User->>User: Completes payment on Provider's site
        Note over User, PaymentSvc: Payment Provider sends Webhook to PaymentSvc
        PaymentSvc->>PaymentSvc: Handle Webhook: Update transaction status (e.g., to 'completed')
        PaymentSvc-->>KafkaBus: Publish event `payment.transaction.completed.v1` (transaction_id, order_id, user_id, game_id, amount)
    else Direct payment success/failure
        PaymentSvc->>PaymentSvc: Process payment
        PaymentSvc-->>KafkaBus: Publish event `payment.transaction.completed.v1` or `payment.transaction.failed.v1`
        OrderSvc-->>APIGW: HTTP 200 OK (transaction_status)
        APIGW-->>ClientApp: HTTP 200 OK (transaction_status)
    end

    subgraph Asynchronous Processing Post-Payment
        KafkaBus-->>OrderSvc: Consume `payment.transaction.completed.v1` (optional, if OrderSvc tracks final status)
        OrderSvc->>OrderSvc: Update Order status to 'completed'
        OrderSvc-->>KafkaBus: Publish `order.completed.v1` (if other services need this specific event)

        KafkaBus-->>LibrarySvc: Consume `payment.transaction.completed.v1` (or `order.completed.v1`)
        LibrarySvc->>LibrarySvc: Verify game_id and user_id
        LibrarySvc->>LibrarySvc: Add game to user's library (UserLibraryItem record)
        LibrarySvc-->>KafkaBus: Publish event `library.game.added.v1` (user_id, game_id, purchase_date)

        KafkaBus-->>NotificationSvc: Consume `payment.transaction.completed.v1` (or `order.completed.v1`)
        NotificationSvc->>NotificationSvc: Prepare purchase confirmation message
        NotificationSvc->>NotificationSvc: Send email/push notification to User
        NotificationSvc-->>KafkaBus: Publish event `notification.sent.v1`

        KafkaBus-->>PaymentSvc: Consume `payment.transaction.completed.v1` (for fiscalization, if not done synchronously)
        PaymentSvc->>PaymentSvc: Initiate fiscal receipt generation (integrates with OFD)
        PaymentSvc-->>KafkaBus: Publish `payment.fiscal.receipt.created.v1`

        KafkaBus-->>AnalyticsSvc: Consume `payment.transaction.completed.v1`
        AnalyticsSvc->>AnalyticsSvc: Log purchase event for analytics
    end

    ClientApp-->>User: Display purchase success/failure message
    User->>ClientApp: (Later) Checks library
    ClientApp->>APIGW: GET /api/v1/library
    APIGW->>LibrarySvc: Forward GET /library
    LibrarySvc-->>APIGW: Library contents (including newly purchased game)
    APIGW-->>ClientApp: Library contents
    ClientApp-->>User: Displays updated library
```

This diagram shows a common flow. Depending on the exact implementation of the "Order Service" (whether it's a distinct service or part of Payment/Catalog), the initial steps might vary slightly. The use of Kafka events for post-payment processing ensures loose coupling and resilience.
```
