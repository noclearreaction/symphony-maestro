## Why

The system lacks a structured and standardized issue categorization and labeling schema, leading to potential administrative drift and ambiguity in tracking development status. Standardized classification of type, status, and priority is needed to coordinate multi-agent workflows, enable machine-readable triaging, and satisfy the Symphony Constitution's principles of Intent Traceability and Transparency.

## What Changes

- Provisioning of 12 standard labels (type, status, and priority) with exact hex colors and descriptions in the repository.
- Updates to the `@issue` subagent's system instructions to integrate the new taxonomy and automate default label application during issue creation.
- Guidelines defining lifecycle state transitions from creation (`status:backlog`) through planning (`status:accepted`), execution (`status:in-progress`), exceptions (`status:blocked`), and resolution (`status:completed`).

## Capabilities

### New Capabilities
- `issue-workflow`: Specifies the standard label schemas, the automated issue-creation behavior of the `@issue` subagent, and the state transition logic for issues.

### Modified Capabilities

## Impact

- Repository configuration (GitHub labels).
- Agent definition files (specifically `.opencode/agents/issue.md`).
- Project documentation and guidelines for state transition management.
