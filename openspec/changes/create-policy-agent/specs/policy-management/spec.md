# Capability: policy-management

## Purpose
Establishes the roles, permissions, constraints, and validation rules for the `@policy` agent dedicated to managing AI prompts, configurations, and global agent documentation.

## ADDED Requirements

### Requirement: Dedicated Policy Agent
The system SHALL define a specialized worker agent named `@policy` defined at `.opencode/agents/policy.md` whose sole purpose is to manage AI-related files, prompts, commands, configurations, and global agent coordination documentation.

#### Scenario: Verify policy agent file existence
- **WHEN** the agent registry or operator checks the configured agents list
- **THEN** the `@policy` agent definition file SHALL be present at `.opencode/agents/policy.md`.

### Requirement: Policy Agent Edit Sandboxing
The `@policy` agent SHALL have edit permissions strictly restricted to `.opencode/*` and the global agent coordination file `AGENTS.md`. It SHALL be denied from editing any application source code, configuration files, test suites, or other project metadata files.

#### Scenario: Policy agent attempts to edit agent prompts
- **WHEN** the `@policy` agent is requested to edit `.opencode/agents/designer.md` or `.opencode/agents/builder.md`
- **THEN** the system SHALL allow the edit to proceed.

#### Scenario: Policy agent attempts to edit application code
- **WHEN** the `@policy` agent attempts to edit a typescript or javascript file outside `.opencode/` (such as `bin/commit-lint.ts`)
- **THEN** the system SHALL block the edit and throw a permission violation.

### Requirement: Policy Agent Execution Sandboxing
The `@policy` agent SHALL be denied general bash execution privileges. It SHALL be permitted to run safe, non-destructive validation commands (specifically `openspec validate`), and SHALL delegate all Git and GitHub command execution to the specialized `git-operator` subagent.

#### Scenario: Policy agent attempts to run git commands
- **WHEN** the `@policy` agent needs to make a commit or check out a branch
- **THEN** the agent SHALL delegate the git command to the `git-operator` subagent instead of running git commands directly.

#### Scenario: Policy agent runs validation commands
- **WHEN** the `@policy` agent executes the command `openspec validate`
- **THEN** the system SHALL allow the execution automatically.

### Requirement: Global Agent Boundary Documentation
The global documentation `AGENTS.md` SHALL clearly document the role of the `@policy` agent and define its boundary interface with the other workspace agents (`@advisor`, `@designer`, `@builder`, `@issue`, `@git-operator`, `@orchestrator`).

#### Scenario: Verify AGENTS.md documentation
- **WHEN** an operator or validator reads the global `AGENTS.md` file
- **THEN** the file SHALL contain a clear definition of the `@policy` agent and its operational boundaries.
