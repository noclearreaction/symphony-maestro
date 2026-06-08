# Capability: controlled-git-workflows

## Purpose
Establishes the rules, specifications, and protocols for managing git branches, commit conventions, and coordination loops.

## Requirements

### Requirement: Direct Commit Prohibition
Direct commits to the `main` branch SHALL be prohibited. All modifications, whether by human operators or AI agents, MUST be executed on dedicated feature branches.

#### Scenario: Attempting to commit to main
- **WHEN** a user or agent attempts to make a commit while checked out on the `main` branch
- **THEN** the system SHALL reject the commit or the agent SHALL refuse to proceed.

### Requirement: Single-Topic Branch Binding
Every active branch SHALL correspond to exactly one OpenSpec change. The branch name MUST match the kebab-case change identifier and be prefixed with `change/` (e.g. `change/controlled-git-workflows`).

#### Scenario: Switching branches
- **WHEN** starting a new change with Git checkout
- **THEN** the active branch SHALL be checked out as `change/<name>`.

### Requirement: Conventional Commit Standard
All commits made to the repository MUST adhere to the Conventional Commits specification. The allowed commit types are `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, and `revert`.

#### Scenario: Committing with standard message
- **WHEN** a commit is formatted as `feat(git): add commit linting script`
- **THEN** the commit SHALL be accepted by the linter.

### Requirement: Automated Commit Message Linting
The system SHALL provide a zero-dependency validation script `bin/commit-lint.ts` that parses and validates commit messages according to the Conventional Commit standard.

#### Scenario: Linting an invalid message
- **WHEN** a commit message of `fixed some things` is checked by `bin/commit-lint.ts`
- **THEN** the linter SHALL exit with a non-zero code.

### Requirement: Atomic Commit Progression
Commits SHALL be performed on each logical unit of work. Commits SHOULD be made incrementally as individual files or tasks in `tasks.md` are completed, rather than holding all changes in a single massive commit.

#### Scenario: Incremental task completion
- **WHEN** an AI builder agent completes a specific subtask in `tasks.md`
- **THEN** it SHALL commit that logical change before proceeding to the next subtask.

### Requirement: AI Agent Workflow Enforcement
All AI agents configured in the repository (Orchestrator, Designer, Builder, Issue) SHALL be instructed to respect and use the branch boundaries, Conventional Commits, and atomic commit practices. These rules MUST be consolidated in the system instructions of the `git-operator` subagent, keeping git-specific mechanics and prompt instructions isolated from non-git agents.

#### Scenario: Agent running git actions
- **WHEN** the `git-operator` agent is executing commands
- **THEN** the agent SHALL verify it is on a `change/*` branch and use Conventional Commits for all git commits.

### Requirement: Isolated Git Execution Subagent
The system SHALL isolate all Git and GitHub command execution permissions strictly to a specialized `git-operator` subagent. No other AI agent in the workspace (Orchestrator, Designer, Advisor, Issue, Builder) SHALL have permission to run local `git` or `gh` commands directly.

#### Scenario: Agent attempting git command
- **WHEN** an agent other than `git-operator` needs to execute a Git operation
- **THEN** it SHALL delegate the task to the `git-operator` subagent instead of running bash commands.

### Requirement: Fine-Grained Git Permission Gating
The `git-operator` subagent SHALL be restricted to a fine-grained, safe subset of Git and GitHub CLI commands. Non-destructive commands (e.g., status, diff, add, log, atomic branch checkouts) SHALL be allowed (`allow`), while sensitive commands (e.g., checkout main, push, pull, PR creation) MUST require user permission (`ask`), and dangerous commands (e.g., force push, hard reset) MUST be blocked (`deny`).

#### Scenario: Subagent running non-destructive command
- **WHEN** `git-operator` executes `git status` or `git diff`
- **THEN** the system SHALL allow the execution automatically.

#### Scenario: Subagent attempting dangerous command
- **WHEN** `git-operator` attempts to run `git push --force` or `git reset --hard`
- **THEN** the system SHALL block the execution and fail with a permission violation.

#### Scenario: Global injection verification
- **WHEN** any agent is initialized in the workspace
- **THEN** the agent SHALL automatically load the controlled git and conventional commits rules via the project-wide `AGENTS.md` instruction file.
