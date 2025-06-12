# Admin Service

## Overview

The Admin Service is a crucial component of the Russian Steam analog platform, providing the necessary tools and interfaces for platform administrators and support staff. Its primary purpose is to manage and oversee various aspects of the platform, including content moderation, user management, technical support operations, security monitoring, system-wide configurations, administrative analytics, and platform-level marketing campaigns. It acts as the central hub for all administrative control and operational management.

## Core Functionality

*   **Content Moderation:** Manages workflows for moderating user-generated and developer-submitted content, including games, reviews, comments, messages, and user profiles. Supports both automated and manual moderation processes.
*   **User Management (Platform Users):** Allows administrators to search, view, and manage platform users, including updating their status (active, blocked, suspended), managing roles, and handling verification issues.
*   **Administrative User Management:** Manages accounts and permissions for administrative staff.
*   **Technical Support:** Provides a system for managing user support tickets, including categorization, assignment, response management, and a knowledge base for common issues.
*   **Security Monitoring:** Tools for identifying suspicious activities, analyzing potentially fraudulent payments, managing IP address blocklists, and auditing administrative actions.
*   **Platform Settings Management:** Enables configuration of global and regional platform parameters, operational limits, system notifications, and promotional banners.
*   **Administrative Analytics:** Offers dashboards and reporting capabilities on platform metrics such as sales, user activity, content engagement, moderation effectiveness, and support ticket statistics.
*   **Marketing Campaign Management:** Tools for creating and managing platform-wide marketing initiatives, including sales events, promotional codes, and special offers.

## Technologies

The Admin Service is designed with flexibility in its technology choices, to be aligned with the overall platform standards.
*   **Backend Language:** Go (with Gin), Python (with Django/Flask), Java (with Spring Boot), or Node.js (with Express). The final choice will be harmonized with the platform's dominant backend technologies.
*   **Database:**
    *   PostgreSQL: For structured relational data (e.g., admin users, moderation queues, support tickets, system settings).
    *   MongoDB: For complex data structures, logs, and audit trails.
*   **Search & Indexing:** Elasticsearch for efficient searching through logs, tickets, and user data.
*   **Caching:** Redis for session management and caching frequently accessed data.
*   **Messaging Queue:** Kafka or RabbitMQ for asynchronous task processing (e.g., report generation, bulk moderation tasks).
*   **Frontend (Admin Panel):** Likely a Single Page Application (SPA) using a modern JavaScript framework (e.g., React, Vue, Angular) – though this is external to the Admin Service backend itself.
*   **Infrastructure:** Docker, Kubernetes.

## API Summary

The Admin Service exposes a RESTful API for administrative operations, typically accessed via an API Gateway. All endpoints require administrator privileges, validated through JWT authentication.

**Base URL:** `/api/v1/admin`

Key endpoint categories:

*   **Administrative Users:**
    *   `GET /users`: List admin users.
    *   `POST /users`: Create an admin user.
    *   `GET /users/{admin_id}`: Get admin user details.
    *   `PUT /users/{admin_id}`: Update admin user.
    *   `PUT /users/{admin_id}/permissions`: Update admin user permissions.
*   **Content Moderation:**
    *   `GET /moderation/queues`: List moderation queues.
    *   `GET /moderation/queues/{queue_id}/items`: Get items in a specific queue.
    *   `PUT /moderation/items/{item_id}/decision`: Make a moderation decision (approve, reject).
*   **Platform User Management:**
    *   `GET /platform-users`: Search platform users.
    *   `PUT /platform-users/{user_id}/status`: Change user status (block, unblock).
    *   `PUT /platform-users/{user_id}/role`: Change user role.
*   **Support System:**
    *   `GET /support/tickets`: List support tickets.
    *   `POST /support/tickets/{ticket_id}/messages`: Add a message to a ticket.
    *   `GET /support/knowledge`: List knowledge base articles.
    *   `POST /support/knowledge`: Create a knowledge base article.
*   **Security Operations:**
    *   `GET /security/audit-log`: View admin action audit logs.
    *   `POST /security/blocked-ips`: Block an IP address.
*   **Platform Settings:**
    *   `GET /settings`: Retrieve platform settings.
    *   `PUT /settings/{setting_id}`: Update a specific setting.
*   **Marketing:**
    *   `GET /marketing/campaigns`: List marketing campaigns.
    *   `POST /marketing/campaigns`: Create a new campaign.

## Data Model

Key data entities managed by the Admin Service include:

*   **AdminUser:** Represents an administrative user with specific roles and permissions.
*   **ModerationQueue:** Defines queues for different types of content requiring moderation.
*   **ModerationItem:** An individual piece of content (e.g., a game submission, a user review) awaiting moderation.
*   **ModerationDecision:** Records the outcome of a moderation action.
*   **SupportTicket:** A user's request for technical support.
*   **KnowledgeBaseArticle:** Articles for self-help and agent assistance.
*   **SystemSetting:** Configurable parameters for the platform.
*   **SecurityIncident:** Records of security-related events or breaches.
*   **BlockedIP:** IP addresses blocked from accessing the platform.
*   **MarketingCampaign:** Details of platform-wide promotional campaigns.
*   **AdminActionLog:** A log of all significant actions performed by administrators.

## Integrations

The Admin Service is highly interconnected with most other platform microservices:

*   **Account Service:** To manage user account details, statuses, and roles.
*   **Auth Service:** For authentication and authorization of administrative staff.
*   **Catalog Service:** To moderate game content, manage visibility, and potentially adjust metadata.
*   **Developer Service:** To receive game submissions for moderation and send back moderation results.
*   **Social Service:** To moderate user-generated social content (reviews, comments).
*   **Payment Service:** To handle refund requests, investigate suspicious transactions, and manage platform-wide promotions.
*   **Analytics Service:** To access platform data for generating administrative reports and dashboards.
*   **Notification Service:** To send out system-wide announcements, user-specific administrative messages, or notifications related to moderation/support.
*   **Download Service:** To manage content availability.
*   **Library Service:** To manage user libraries in exceptional cases (e.g., revoking access due to fraud).

## Error Handling

Error responses from the REST API adhere to the platform's standard JSON structure:
```json
{
  "status": "error",
  "error": {
    "code": "ERROR_CODE", // e.g., "VALIDATION_ERROR", "UNAUTHORIZED", "RESOURCE_NOT_FOUND"
    "message": "Human-readable error description.",
    "details": { /* Optional, e.g., field-specific validation errors */ }
  }
}
```
Common HTTP status codes (400, 401, 403, 404, 500) are used to indicate the nature of errors. Detailed logging with trace IDs is implemented for all errors to facilitate debugging.
