## ADDED Requirements

### Requirement: Isolated Git Execution Subagent
The system SHALL isolate all Git and GitHub command execution permissions strictly to a specialized `git-operator` subagent. No other AI agent in the workspace (Orchestrator, Designer, Advisor, Issue, Builder) SHALL have permission to run local `git` or `gh` commands directly.

#### Scenario: Agent attempting git command
- **WHEN** an agent other than `git-operator` needs to execute a Git operation
- **THEN** it SHALL delegate the task to the `git-operator` subagent instead of running bash commands.

### Requirement: Fine-Grained Git Permission Gating
The `git-operator` subagent SHALL be restricted to a safe, controlled subset of Git and GitHub CLI commands. Non-destructive commands (e.g., status, diff, add, log, local branches, checkouts) SHALL be allowed (`allow`), while state-changing, remote-syncing, or history-altering commands (e.g., checkout main, push, pull, reset, merge, PR operations) MUST require explicit user validation and approval (`ask`).

#### Scenario: Subagent running non-destructive command
- **WHEN** `git-operator` executes `git status` or `git diff`
- **THEN** the system SHALL allow the execution automatically.

#### Scenario: Subagent attempting sensitive state-changing command
- **WHEN** `git-operator` attempts to run `git push -vf` or `git reset --hard`
- **THEN** the system SHALL intercept the execution and prompt the user for validation and permission.

## MODIFIED Requirements

### Requirement: AI Agent Workflow Enforcement
All AI agents configured in the repository (Orchestrator, Designer, Builder, Issue) SHALL be instructed to respect and use the branch boundaries, Conventional Commits, and atomic commit practices. These rules MUST be consolidated in the system instructions of the `git-operator` subagent, keeping git-specific mechanics and prompt instructions isolated from non-git agents.

#### Scenario: Agent running git actions
- **WHEN** the `git-operator` agent is executing commands
- **THEN** the agent SHALL verify it is on a `change/*` branch and use Conventional Commits for all git commits.
