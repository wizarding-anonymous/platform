# Analytics Service

## Overview

The Analytics Service is a central component of the Russian Steam analog platform, responsible for collecting, processing, analyzing, and providing insights from the vast amounts of data generated across the platform.

**For detailed specification, please see: [./docs/README.md](./docs/README.md)**

## Core Functionality (Summary)

*   Data Collection from all microservices.
*   Data Processing (stream and batch).
*   Metrics & Reporting (KPIs, DAU, MAU, ARPU, etc.).
*   User Behavior Analysis.
*   Audience Segmentation.
*   Predictive Analytics (sales forecasting, churn prediction, recommendations).
*   System Performance Monitoring aggregation.

## Technologies (Summary)

*   Data Processing: Scala, Python, Java; Apache Spark, Kafka Streams/Apache Flink.
*   Data Storage: ClickHouse (DWH), S3-compatible (Data Lake), PostgreSQL (Metadata).
*   Messaging: Apache Kafka.
*   ML: TensorFlow, PyTorch, Scikit-learn, MLflow.
*   APIs: REST, GraphQL, WebSocket.

## Integrations (Summary)

Consumes events from all other microservices. Provides data to Developer Service, Admin Service, and potentially Catalog/Recommendation services.
