# Documentation Guidelines and Maintenance Process

This document outlines the guidelines and processes for maintaining the documentation of the "Российский Аналог Платформы Steam" project. Up-to-date and consistent documentation is crucial for onboarding new team members, ensuring architectural alignment, and facilitating effective development and operations.

## 1. Ownership

*   **Project-Wide Standards Documents:**
    *   Documents prefixed with `project_` (e.g., `project_api_standards.md`, `project_technology_stack.md`) and other general standards like `CODING_STANDARDS.md`, `README.md` (root) are owned by the **Architecture Team** or a designated **Lead Documentation Maintainer**.
    *   Changes to these documents should be proposed via Pull Requests and require review from the Architecture Team.
*   **Microservice-Specific Documentation:**
    *   The documentation for each microservice (typically located in `backend/<service-name>/docs/README.md`) is owned by the **technical lead or team responsible for that microservice**.
*   **Frontend Application Documentation:**
    *   The documentation for the frontend application (typically located in `frontend/docs/README.md` and based on `standard_frontend_template.md`) is owned by the **Frontend Team Lead** or a designated **Lead Frontend Documentation Maintainer**.
*   **Workflow Diagrams (`project_workflows/`):**
    *   Owned by the **Architecture Team**, with input from the teams involved in those workflows.

## 2. Process for Updating Documentation

### 2.1. Changes Driven by Code or Design Modifications
*   **Requirement:** Documentation updates **must** accompany any significant change in code, API design, architecture, or functionality.
*   **Process:**
    1.  When a feature is developed or a change is made, the developer or team responsible must identify the affected documentation.
    2.  Updates to the documentation should be included in the **same Pull Request (PR)** as the code changes.
    3.  The PR review process must include a check for documentation accuracy and completeness related to the changes.
    4.  For significant architectural changes impacting multiple services or project-wide standards, the proposer should first discuss the change with the Architecture Team and update the relevant high-level documents before (or alongside) service-specific changes.

### 2.2.发现问题或提出改进建议 (Reporting Issues or Suggesting Improvements)
*   Team members who identify errors, outdated information, or areas for improvement in the documentation should:
    *   Create an **Issue** in the project's issue tracking system (e.g., Jira, GitHub Issues), tagging it with "documentation".
    *   Clearly describe the problem and the suggested change.
    *   If possible, and for minor changes, create a PR directly with the proposed fix.

## 3. Documentation Style and Structure Guidelines

*   **Language:**
    *   Primary language for documentation is **Russian**.
    *   Code comments and identifiers should follow `CODING_STANDARDS.md` (typically English for code elements, Russian for detailed comments).
*   **Tone:** Clear, concise, and professional. Avoid jargon where possible or explain it in the `project_glossary.md`.
*   **Formatting:** Use Markdown. Ensure readability with proper headings, lists, code blocks, and diagrams.
*   **Scope Clarity:**
    *   Project and service-level documentation should clearly define the scope of implemented and planned functionality.
    *   Explicitly mention major features or areas that are considered out-of-scope if they are commonly associated with similar platforms, to manage expectations. For example, as noted in the main project `README.md`.
*   **Microservice Documentation:**
    *   All microservices **must** use the `standard_microservice_template.md` as the basis for their detailed documentation in `backend/<service-name>/docs/README.md`.
    *   The root `backend/<service-name>/README.md` should provide a brief overview and a link to the detailed `docs/README.md`.
*   **Frontend Documentation:**
    *   The main frontend application **must** use the `standard_frontend_template.md` as the basis for its detailed documentation in `frontend/docs/README.md`. The root `frontend/README.md` should provide a brief overview and a link to the detailed `frontend/docs/README.md`.
*   **Diagrams:**
    *   Use Mermaid.js for sequence diagrams, ERDs, flowcharts, etc., embedded directly in Markdown.
    *   Ensure diagrams are kept up-to-date with architectural changes.
    *   Store complex or numerous workflow diagrams in the `project_workflows/` directory.
*   **Cross-Referencing:**
    *   Use relative links to refer to other documents within the repository.
    *   Refer to project-wide standards documents instead of duplicating information. For example, service API docs should refer to `project_api_standards.md` for common principles.
*   **Placeholders and TODOs:** All documentation should be actively maintained. Any placeholders (e.g., "TODO", "{{PLACEHOLDER}}", "[будет дополнено]") must be treated as temporary and should be replaced with actual information as soon as it becomes available. If information cannot be provided immediately, the placeholder should clearly state what is missing, why, and who is responsible for providing it, along with an estimated timeline if possible. Stale placeholders that are no longer relevant should be removed.
*   **Glossary:** All project-specific terms or acronyms should be defined in `project_glossary.md`.

## 4. Periodic Review Process

*   **Schedule:** All project-wide standards documents and microservice specifications should undergo a formal review at least **once every quarter** or **before a major platform release**.
*   **Responsibility:**
    *   The Architecture Team is responsible for initiating and overseeing the review of project-wide standards.
    *   Microservice teams are responsible for reviewing their respective service documentation.
*   **Process:**
    1.  Reviewers check for accuracy, completeness, clarity, and consistency with the current state of the project and other documentation.
    2.  Updates are made via PRs.
    3.  A brief notification can be sent out to the development team highlighting significant changes or updates after a review cycle.

## 5. "Docs Archive" Folder

*   The `Docs Archive/` folder is intended for outdated or superseded documents.
*   **No current development should rely on information solely from the Docs Archive.**
*   If information from an archived document is still relevant, it should be migrated to an active document and the archived version clearly marked as obsolete.

By following these guidelines, we aim to maintain a high-quality, reliable, and useful set of documentation for the project.

## 6. Configuration Management Strategy

Effective configuration management is essential for deploying and running microservices reliably across different environments. This section outlines the standard approach to managing configuration within the project.

### 6.1. Naming Conventions and Formats
*   **Preferred Format:** YAML (`.yaml` or `.yml`) is the preferred format for configuration files due to its widespread support and readability. Environment variables can also be used, especially for sensitive data or settings that vary frequently between deployment environments.
*   **Default File Name:** Each microservice should include a default configuration file named `config.default.yaml` or `config.example.yaml` in its root directory. This file should contain all possible configuration keys with non-sensitive default values suitable for local development or as a template.
*   **Environment-Specific Files:** For different environments (development, testing, staging, production), configurations can be managed by:
    *   Using environment variables to override defaults.
    *   Employing environment-specific configuration files (e.g., `config.dev.yaml`, `config.prod.yaml`) that are loaded based on an environment indicator (like `NODE_ENV` or a custom environment variable). **These environment-specific files (especially those with secrets) should NOT be committed to version control.** Add them to the service's `.gitignore` file.
*   **Hierarchy:** Applications should load the `config.default.yaml` first, then override values with an environment-specific file if it exists, and finally override with any relevant environment variables.

### 6.2. Handling Secrets
*   **Never commit secrets** (API keys, database passwords, etc.) directly into configuration files in the repository.
*   Use environment variables for secrets, injected at deployment time (e.g., via CI/CD pipeline, Kubernetes Secrets, Docker Swarm Secrets).
*   Alternatively, use a dedicated secrets management tool (e.g., HashiCorp Vault). If such a tool is adopted, its usage should be documented here.
*   The `config.default.yaml` or `config.example.yaml` should use placeholder values for secrets (e.g., `apiKey: YOUR_API_KEY_HERE`).

### 6.3. Configuration Schema and Validation
*   While not strictly enforced by a central tool yet, it is highly recommended that each microservice:
    *   Clearly documents its expected configuration schema in its `docs/README.md`.
    *   Implements validation logic at startup to ensure all required configurations are present and correctly formatted. The service should fail to start if critical configurations are missing or invalid.

### 6.4. Example Configuration Structure (Illustrative)
A `config.default.json` for a typical microservice might look like this. **Note: This example uses JSON for illustration. YAML is the preferred format as per `project_api_standards.md`, which also contains YAML examples.**

```json
{
  "serviceName": "my-awesome-service",
  "port": 3000,
  "logLevel": "info",
  "database": {
    "host": "localhost",
    "port": 5432,
    "username": "user",
    "password": "POSTGRES_PASSWORD_ENV_VAR", // Indicates this should come from an env var
    "dbName": "mydb"
  },
  "externalApiService": {
    "url": "https://api.externalservice.com/v1",
    "apiKey": "EXTERNAL_API_KEY_ENV_VAR" // Indicates this should come from an env var
  },
  "featureFlags": {
    "newSearchAlgorithmEnabled": false
  }
}
```

### 6.5. Microservice Configuration Templates
*   Each microservice, when created, should include a `config.default.yaml` (or `config.example.yaml`) in its root directory, reflecting the necessary configuration keys for that service.
*   A `.gitignore` file within each microservice should list environment-specific config files like `config.dev.yaml`, `config.prod.yaml`, `*.local.yaml`, `.env` to prevent accidental commits of sensitive or environment-specific data.

## 7. Documenting and Executing Tests

Comprehensive testing is vital for ensuring the quality, reliability, and maintainability of all microservices. This section outlines the standards for documenting and executing tests.

### 7.1. Test Location and Structure
*   **Dedicated Test Directory:** All tests for a microservice **must** be located in a `tests` directory within the root folder of that microservice (e.g., `backend/account-service/tests/`).
*   **Test Types:** It's recommended to structure tests by type (e.g., `unit/`, `integration/`, `e2e/`) within the `tests` directory, as appropriate for the service's testing strategy.
*   **Test File Naming:** Test files should clearly indicate what they are testing. Refer to `CODING_STANDARDS.md` for language-specific naming conventions for test files and functions.

### 7.2. Documenting Tests
*   **Clarity:** Tests should be written clearly and be easy to understand. Test names should describe the specific scenario or behavior being tested.
*   **Microservice Documentation:** The `docs/README.md` for each microservice (based on `standard_microservice_template.md`) should include a section on testing. This section should specify:
    *   Types of tests implemented (Unit, Integration, E2E).
    *   Instructions on how to set up the testing environment (e.g., specific database setup, environment variables, dependencies).
    *   The precise command(s) to run all tests locally.
    *   Any specific tools or frameworks used for testing that require special attention.

### 7.3. Running Tests Locally
*   **Developer Responsibility:** Developers are responsible for running all relevant tests locally after making changes and before submitting a Pull Request.
*   **Ease of Execution:** Running tests locally should be straightforward, typically a single command per microservice or test suite, as documented in the microservice's `docs/README.md`.
*   **Environment:** Ensure that local testing can be performed without extensive manual configuration, possibly using default configurations (e.g., `config.default.json` or environment variables for local development databases/services).

### 7.4. Automated Testing (CI/CD)
*   **Server-Side Execution:** All tests (unit, integration, and where feasible, E2E) should be automated and run as part of the Continuous Integration (CI) pipeline for every Pull Request and before deployment to any environment.
*   **CI Configuration:** The CI pipeline configuration (e.g., GitHub Actions workflows, Jenkinsfile) will define the steps to build, test, and deploy services. These configurations are part of the project's infrastructure-as-code.
*   **Test Reports:** Test results, including coverage reports, should be captured and made available through the CI system.

### 7.5. Testing Assistance
*   While the primary responsibility for writing and executing tests lies with the development team, assistance can be requested for:
    *   Guidance on complex testing scenarios.
    *   Troubleshooting CI/CD testing issues.
    *   Understanding cross-service integration testing strategies.
    *   This assistance is supplementary to the established automated and local testing practices.

Refer to `CODING_STANDARDS.md` for language-specific testing guidelines, recommended libraries, and coverage targets.
