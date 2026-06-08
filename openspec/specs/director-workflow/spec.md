# Capability: director-workflow

## Purpose
Structured Git, GitHub PR, and custom builder agent local workflow for Director development.

## Requirements

### Requirement: Builder Agent Soft Boundaries
The custom `builder` agent system prompt SHALL enforce strict soft boundaries that prevent any direct modifications to strategy documents or other agent prompts.

#### Scenario: Refusing strategy modification
- **WHEN** the builder agent is requested to edit files under `strategy/`
- **THEN** the builder agent SHALL refuse the change and prompt the user to consult the Advisor or Designer agent instead.

### Requirement: Structured Local Workspace Isolation
The repository SHALL ignore a designated `.symphony/` directory structure to allow agents to write scratch records, plans, and evidence without polluting git history.

#### Scenario: Isolating temporary logs
- **WHEN** an agent writes files under `.symphony/scratchpad/`
- **THEN** the file SHALL not appear as an untracked file in git status.

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
