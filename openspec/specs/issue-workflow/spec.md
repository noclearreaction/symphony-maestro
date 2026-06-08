# Capability: issue-workflow

## Purpose
Establishing a predictable, machine-readable Issue Categorization and Labeling schema, automating label application on issue creation, and defining standard guidelines for tracking issue state transitions.

## Requirements

### Requirement: Standardized Label Taxonomy
The system SHALL establish a standardized schema containing 12 specific labels with exact names, hex colors, and descriptions.

#### Scenario: Verify label definition list
- **WHEN** the label provisioning tool or operator lists repository labels
- **THEN** the labels SHALL exactly match the following:
  - type:feature (#0E8A16, "New functionality or intent.")
  - type:bug (#D93F0B, "Unexpected behavior or failure.")
  - type:chore (#EDEDED, "Internal maintenance, CI, or configuration.")
  - type:spike (#FBCA04, "Investigation, technical spikes, or documentation (preferred over research).")
  - status:backlog (#EDEDED, "Default for new/untriaged work.")
  - status:accepted (#FEF2C0, "Approved for implementation.")
  - status:in-progress (#FEF2C0, "Active execution phase.")
  - status:completed (#0E8A16, "Work finished, PR submitted.")
  - status:blocked (#000000, "Process halted by external factor.")
  - priority:high (#B60205, "Critical blocker for milestones.")
  - priority:medium (#FBCA04, "Standard prioritized work.")
  - priority:low (#EDEDED, "Elective polish or minor backlog items.")

### Requirement: Issue Subagent Draft Recommendations
The `@issue` subagent SHALL evaluate issue request context to recommend correct `type:*` and `priority:*` labels during issue drafting.

#### Scenario: Subagent recommends labels during drafting
- **WHEN** the user prompts the `@issue` subagent to draft an issue with a specific intent
- **THEN** the subagent SHALL evaluate the context and recommend corresponding `type:*` and `priority:*` labels in the draft.

### Requirement: Issue Subagent Automated Creation
The `@issue` subagent SHALL automatically apply the `status:backlog` label, along with the approved `type:*` and `priority:*` labels, when executing issue creation.

#### Scenario: Creating issue with default status
- **WHEN** the `@issue` subagent executes the issue creation workflow
- **THEN** the subagent SHALL invoke the creation CLI with the `status:backlog` label and any selected type and priority labels.

### Requirement: Development Lifecycle Transition Guidelines
The system SHALL maintain a version-controlled documentation file defining the rules and guidelines for human and machine operators to transition issues through the development lifecycle.

#### Scenario: Documenting issue state machine
- **WHEN** an operator consults the project governance
- **THEN** a document at `governance/issue-lifecycle.md` SHALL define rules for moving issues from backlog to accepted, in-progress, blocked, and completed states.
