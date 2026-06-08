## ADDED Requirements

### Requirement: Native OpenSpec Change Initialization
The system SHALL support native OpenSpec change initialization. Starting a new change MUST be performed by checking out a dedicated local `change/<name>` branch via Git and running the standard `openspec new change "<name>"` command directly.

#### Scenario: Creating a change with native commands
- **WHEN** the user checks out a topic branch and runs `openspec new change "my-feature"`
- **THEN** the standard OpenSpec change directory SHALL be instantiated successfully.

### Requirement: Standard OpenSpec Verification and Sync
The system SHALL support standard OpenSpec verification and synchronization. Changes SHALL be validated using `openspec validate` and synchronized to main specifications using standard OpenSpec change-based workflow practices.

#### Scenario: Synchronizing specifications
- **WHEN** a change implementation is verified and complete
- **THEN** the specifications SHALL be synchronized to the main spec directory and the change archived to the planning archive directory.

## REMOVED Requirements

### Requirement: Branch Initiation Automation
**Reason**: Custom script `bin/director-start` is deprecated and removed to align with OpenSpec community standards and Constitutional Principle 5 (Tooling Discipline).
**Migration**: Perform branch checkout with standard Git commands and run standard `openspec new change "<name>"` directly.

### Requirement: GitHub PR and Sync Automation
**Reason**: Custom script `bin/director-submit` is deprecated and removed to align with OpenSpec community standards and Constitutional Principle 5 (Tooling Discipline).
**Migration**: Run standard `openspec validate` and manually/agent-edit main specifications to sync, then push and create a PR via GitHub CLI `gh pr create`.
