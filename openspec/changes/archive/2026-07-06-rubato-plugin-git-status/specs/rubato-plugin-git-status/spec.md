## ADDED Requirements

### Requirement: Anchor-Controlled Injection Eligibility
The system SHALL perform runtime injection only for requests with a valid runtime anchor block in `messages[0].content`.

#### Scenario: Missing anchor bypasses injection
- **WHEN** a request has no valid runtime anchor block in `messages[0].content`
- **THEN** Rubato forwards the request unchanged

#### Scenario: Malformed anchor fails fast
- **WHEN** a request contains a malformed runtime anchor block in `messages[0].content`
- **THEN** Rubato returns a deterministic request failure and does not forward upstream

### Requirement: Declared Plugin Contract
The system SHALL execute only plugins declared by the runtime anchor and SHALL fail fast for unknown or failed declared plugins.

#### Scenario: Declared plugins only
- **WHEN** the runtime anchor declares a subset of available plugin keys
- **THEN** Rubato executes only that declared subset

#### Scenario: Unknown plugin key
- **WHEN** the runtime anchor declares an unregistered plugin key
- **THEN** Rubato fails the request with an explicit unknown-plugin error

### Requirement: Runtime Mutation Semantics
The system SHALL prepend runtime state to `messages[-1]` and SHALL inject deterministic idempotent guidance into `messages[0]`.

#### Scenario: Runtime-state injection
- **WHEN** declared plugins execute successfully
- **THEN** Rubato prepends runtime-state output to `messages[-1]`

#### Scenario: Deterministic guidance rendering
- **WHEN** two requests declare equivalent anchors
- **THEN** generated guidance content in `messages[0]` is byte-identical

#### Scenario: Guidance idempotence
- **WHEN** `messages[0]` already includes guidance for declared plugins
- **THEN** Rubato does not append duplicate guidance

### Requirement: Git Status MVP Plugin
The system SHALL provide a `git_status` plugin with explicit hygiene metrics and recognizable repository-state modes.

#### Scenario: Git status reports hygiene metrics
- **WHEN** the runtime anchor includes `git_status` and repository status is readable
- **THEN** runtime-state includes branch-or-head state, ahead/behind, commits-ahead, staged, unstaged tracked-modified, and untracked counts

#### Scenario: Detached HEAD is visible state
- **WHEN** the repository is in detached HEAD state
- **THEN** runtime-state reports detached HEAD as explicit plugin output state

#### Scenario: Bare repository is visible state
- **WHEN** the repository is bare
- **THEN** runtime-state reports bare-repository status as explicit plugin output state

#### Scenario: Git status execution failure
- **WHEN** repository status cannot be determined due to non-repo or execution failure
- **THEN** Rubato fails the request with a git_status-specific failure response
