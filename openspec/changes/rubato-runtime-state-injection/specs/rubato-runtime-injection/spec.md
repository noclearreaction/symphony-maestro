## ADDED Requirements

### Requirement: Anchor-Controlled Runtime Injection
The system SHALL support runtime injection only when `messages[0].content` contains a valid rubato anchor block that declares plugin keys.

#### Scenario: Valid anchor declares plugin keys
- **WHEN** an inbound request contains a valid anchor block in `messages[0].content` declaring one or more plugin keys
- **THEN** the system parses the declared keys and marks only those plugins as active for that request

#### Scenario: Anchor declares plugin arguments
- **WHEN** an inbound request contains a valid anchor block where plugin declarations include static arguments
- **THEN** the system parses and validates those arguments and passes them only to the declared plugin instance

#### Scenario: Missing anchor bypasses injection
- **WHEN** an inbound request does not contain a valid rubato anchor block in `messages[0].content`
- **THEN** the system forwards the request without runtime injection

### Requirement: Injection Target And Controlled Prompt Mutation
The system SHALL prepend runtime state output to `messages[-1]` for eligible requests and SHALL augment `messages[0]` with plugin guidance in an idempotent manner.

#### Scenario: Eligible request receives injected runtime state
- **WHEN** a request is eligible for injection and declared plugins execute successfully
- **THEN** the system prepends a runtime state block to `messages[-1]` before forwarding upstream

#### Scenario: Plugin guidance injected once
- **WHEN** runtime injection is performed for a request whose `messages[0]` does not yet contain rubato plugin guidance
- **THEN** the system injects plugin-presence usage guidance for declared plugins into `messages[0]` before forwarding upstream

#### Scenario: Plugin guidance is not duplicated
- **WHEN** runtime injection is performed and `messages[0]` already contains rubato plugin guidance for declared plugins
- **THEN** the system does not append duplicate guidance text

#### Scenario: Guidance rendering is deterministic across sessions
- **WHEN** two separate sessions submit requests with equivalent anchor declarations (same plugin set and same plugin args)
- **THEN** the injected plugin guidance content in `messages[0]` is byte-identical across those sessions

#### Scenario: Guidance excludes volatile fields
- **WHEN** plugin guidance is generated for `messages[0]`
- **THEN** the guidance contains no volatile runtime fields such as timestamps, random IDs, host-specific paths, or request counters

### Requirement: Declared Plugin Execution Semantics
The system SHALL execute only plugin keys declared by the request anchor and SHALL refresh plugin output on every eligible request.

#### Scenario: Only declared plugins run
- **WHEN** the request anchor declares a subset of available plugin keys
- **THEN** the system executes only the declared subset and does not execute undeclared plugins

#### Scenario: Declared plugin receives only its own args
- **WHEN** multiple plugins are declared with distinct argument maps
- **THEN** each plugin execution receives only its declared argument map and not arguments of other plugins

#### Scenario: Per-request refresh
- **WHEN** two consecutive eligible requests declare the same plugin keys
- **THEN** the system re-executes plugin commands for each request rather than reusing prior request output

### Requirement: Fail-Fast Behavior For Declared Plugin Failures
The system SHALL fail the request when any declared plugin fails to execute, returns invalid output, or is unknown.

#### Scenario: Plugin execution failure
- **WHEN** a declared plugin command returns an execution error or timeout
- **THEN** the system returns a request failure response that identifies the plugin key and failure reason

#### Scenario: Unknown declared plugin key
- **WHEN** a request anchor declares a plugin key that is not registered
- **THEN** the system returns a request failure response indicating an unknown plugin key

### Requirement: Git Status MVP Plugin
The system SHALL provide a `git_status` plugin in MVP that reports current repository hygiene signals including branch name, ahead/behind counts, committed changes count, staged changes count, and untracked file count.

#### Scenario: Git status plugin success
- **WHEN** the request anchor includes `git_status` and repository state is readable
- **THEN** the injected runtime state includes git hygiene signals for branch, ahead/behind, committed, staged, and untracked counts

#### Scenario: Git status plugin repository error
- **WHEN** the request anchor includes `git_status` but repository status cannot be determined
- **THEN** the request fails with a git_status-specific error response

### Requirement: Plugin Extensibility Contract
The system SHALL expose a plugin contract that supports adding additional runtime injection plugins without changing anchor semantics.

#### Scenario: New plugin added
- **WHEN** a new plugin is registered according to the plugin contract
- **THEN** requests declaring that plugin key in the anchor can execute it without requiring anchor format changes
