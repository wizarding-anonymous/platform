# Analytics Service

## Overview

The Analytics Service is a central component of the Russian Steam analog platform, responsible for collecting, processing, analyzing, and providing insights from the vast amounts of data generated across the platform. Its primary goal is to empower data-driven decision-making for business strategy, operational improvements, user experience enhancements, and to provide developers with relevant statistics about their products.

## Core Functionality

*   **Data Collection:** Ingests events and data from all other microservices, including user actions (registrations, game views, purchases, playtime), business transactions, system logs, and marketing campaign interactions. Supports both real-time event streaming and batch imports.
*   **Data Processing:** Performs stream and batch processing of raw data. This includes transformation, aggregation, cleansing, and anonymization of personally identifiable information (PII) in compliance with privacy regulations (e.g., 152-FZ).
*   **Metrics & Reporting:** Calculates key performance indicators (KPIs) and various metrics (e.g., DAU, MAU, ARPU, conversion rates, retention). Generates standard and customizable reports for different stakeholders.
*   **User Behavior Analysis:** Provides tools and data for analyzing user journeys, cohort behavior, A/B testing results, and patterns in user reviews and interactions.
*   **Audience Segmentation:** Enables the creation of static and dynamic user segments based on various criteria for targeted marketing and personalized experiences.
*   **Predictive Analytics:** Develops and deploys machine learning models for tasks such as sales forecasting, churn prediction, game recommendations, and fraud detection.
*   **System Performance Monitoring:** Aggregates and analyzes metrics related to the health, performance, and availability of the platform's microservices.

## Technologies

The Analytics Service employs a specialized technology stack suited for big data processing and machine learning:

*   **Primary Data Processing Languages:** Scala, Python, Java
*   **API Layer Languages:** Go, Java
*   **Data Processing Frameworks:**
    *   Batch Processing: Apache Spark
    *   Stream Processing: Kafka Streams or Apache Flink
*   **Data Storage:**
    *   Analytical Database: ClickHouse
    *   Raw Data / Data Lake: S3-compatible object storage (e.g., MinIO)
    *   Metadata & Configuration: PostgreSQL
*   **Messaging/Event Streaming:** Apache Kafka
*   **Machine Learning:**
    *   Libraries: TensorFlow, PyTorch, Scikit-learn
    *   Management: MLflow
*   **API Types:**
    *   REST (using Spring Boot for Java, or standard Go libraries/frameworks like Gin/Echo)
    *   GraphQL (e.g., using Apollo Server or Hasura for specific data access patterns)
*   **Data Visualization:** Grafana, Apache Superset (for ad-hoc queries and dashboards)
*   **Infrastructure:** Docker, Kubernetes

## API Summary

The Analytics Service provides APIs for other services and authorized administrative tools to access analytical data and insights. Adherence to platform API standards ([Стандарты API, форматов данных, событий и конфигурационных файлов.txt](https://placeholder.com/link-to-api-standards)) is maintained.

**Base URL:** `/api/v1/analytics`

### REST API

*   **Metrics:**
    *   `GET /metrics`: List available metrics.
    *   `GET /metrics/{metric_name}`: Retrieve specific metric data (supports query parameters for time range, dimensions, filters, granularity).
    *   `GET /metrics/realtime`: Access real-time aggregated metrics.
*   **Reports:**
    *   `GET /reports`: List available pre-defined and custom reports.
    *   `POST /reports/{report_id}/generate`: Trigger report generation.
    *   `GET /reports/instances/{instance_id}/download`: Download a generated report.
*   **Segments:**
    *   `GET /segments`: List user segments.
    *   `POST /segments`: Create a new user segment based on defined criteria.
    *   `POST /segments/{segment_id}/export`: Export users within a segment.
*   **Predictive Analytics:**
    *   `GET /predictions/models`: List available ML models.
    *   `POST /predictions/{prediction_type}`: Request a specific type of prediction (e.g., churn probability for a user).

### GraphQL API

*   **Endpoint:** `/api/v1/analytics/graphql`
*   Provides a flexible query language for clients to request specific datasets, metrics, and report configurations. Allows fetching complex nested data in a single request.

### WebSocket API (Streaming)

*   **Endpoint:** `/api/v1/analytics/streaming`
*   Used for subscribing to real-time data streams, such as live metrics updates or system alerts.

## Data Model

The Analytics Service primarily deals with:

*   **Events:** Raw and processed records of user actions, system events, and business operations. Stored typically in ClickHouse and S3.
*   **Metrics:** Aggregated numerical data representing specific aspects of platform performance or user behavior (e.g., DAU, conversion_rate). Stored in ClickHouse.
*   **Reports & Report Instances:** Definitions of reports and their generated outputs. Metadata in PostgreSQL, report files potentially in S3.
*   **Segments:** Definitions of user groups based on specific criteria. Metadata in PostgreSQL.
*   **MLModels & Predictions:** Definitions of machine learning models, their training parameters, and the predictions they generate. Metadata and potentially model artifacts stored appropriately (PostgreSQL, S3/MLflow).

## Integrations

*   **Data Ingestion (Primarily Kafka):**
    *   Consumes events from nearly all other microservices (Account, Auth, Catalog, Library, Payment, Download, Social, Notification, Developer, Admin) to build a comprehensive view of platform activity.
*   **Data Provisioning:**
    *   **Developer Service:** Provides analytics and sales data for developers regarding their games.
    *   **Admin Service:** Offers platform-wide dashboards and reports for administrative oversight.
    *   **Catalog/Recommendation Service (Potentially):** Could feed processed data (e.g., popularity metrics, user preferences) to services responsible for recommendations or content ranking.
    *   **Marketing Tools:** Exported segments can be used by external or internal marketing automation tools.
*   **Auth Service:** For authenticating and authorizing API requests to the Analytics Service.

## Error Handling

Error responses from the REST API follow the platform's standard JSON structure:
```json
{
  "status": "error",
  "error": {
    "code": "ERROR_CODE", // e.g., "DATA_VALIDATION_ERROR", "QUERY_EXECUTION_ERROR"
    "message": "Human-readable error description.",
    "details": { /* Optional, e.g., specific query errors */ }
  }
}
```
GraphQL and WebSocket APIs will have their own standard ways of reporting errors. Common issues include data processing failures, query execution errors, or issues with ML model operations. All significant errors are logged with a trace ID.
