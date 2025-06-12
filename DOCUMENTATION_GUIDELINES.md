# Documentation Guidelines and Maintenance Process

This document outlines the guidelines and processes for maintaining the documentation of the "Российский Аналог Платформы Steam" project. Up-to-date and consistent documentation is crucial for onboarding new team members, ensuring architectural alignment, and facilitating effective development and operations.

## 1. Ownership

*   **Project-Wide Standards Documents:**
    *   Documents prefixed with `project_` (e.g., `project_api_standards.md`, `project_technology_stack.md`) and other general standards like `CODING_STANDARDS.md`, `README.md` (root) are owned by the **Architecture Team** or a designated **Lead Documentation Maintainer**.
    *   Changes to these documents should be proposed via Pull Requests and require review from the Architecture Team.
*   **Microservice-Specific Documentation:**
    *   The documentation for each microservice (typically located in `backend/<service-name>/docs/README.md`) is owned by the **technical lead or team responsible for that microservice**.
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
*   **Microservice Documentation:**
    *   All microservices **must** use the `standard_microservice_template.md` as the basis for their detailed documentation in `backend/<service-name>/docs/README.md`.
    *   The root `backend/<service-name>/README.md` should provide a brief overview and a link to the detailed `docs/README.md`.
*   **Diagrams:**
    *   Use Mermaid.js for sequence diagrams, ERDs, flowcharts, etc., embedded directly in Markdown.
    *   Ensure diagrams are kept up-to-date with architectural changes.
    *   Store complex or numerous workflow diagrams in the `project_workflows/` directory.
*   **Cross-Referencing:**
    *   Use relative links to refer to other documents within the repository.
    *   Refer to project-wide standards documents instead of duplicating information. For example, service API docs should refer to `project_api_standards.md` for common principles.
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
```
