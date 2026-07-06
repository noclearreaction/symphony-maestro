## ADDED Requirements

### Requirement: Injector scans prior messages for last known plugin output
The `Injector.Apply` SHALL scan backward through prior messages (from `messages[-2]` up to `MaxAge` positions) to find the most recent `rubato:state` block containing each declared plugin's output.

#### Scenario: Prior state found within window
- **WHEN** a `rubato:state` block containing `[git_status]` output exists within the last `MaxAge` messages
- **THEN** the last known output for `git_status` SHALL be extracted for comparison

#### Scenario: No prior state in window
- **WHEN** no `rubato:state` block for a plugin exists within the last `MaxAge` messages
- **THEN** the plugin SHALL be treated as having no known prior state

### Requirement: Only changed or stale plugins are injected
The `Injector.Apply` SHALL inject a plugin's output only if its fresh output differs from the last known output, or if it was not found within the `MaxAge` window.

#### Scenario: Plugin output unchanged within window
- **WHEN** fresh plugin output matches the last known output and the last injection is within `MaxAge` messages
- **THEN** that plugin SHALL NOT appear in the injected state block

#### Scenario: Plugin output changed
- **WHEN** fresh plugin output differs from the last known output
- **THEN** that plugin SHALL appear in the injected state block with the new output

#### Scenario: Plugin beyond MaxAge window
- **WHEN** the last injection for a plugin is more than `MaxAge` messages ago
- **THEN** that plugin SHALL appear in the injected state block regardless of whether output changed

#### Scenario: First turn — no prior history
- **WHEN** there are no prior messages containing a rubato:state block
- **THEN** ALL declared plugins SHALL be injected

### Requirement: State block is omitted when no plugins require injection
When all declared plugins are within the `MaxAge` window and their outputs are unchanged, `Injector.Apply` SHALL NOT prepend a `rubato:state` block to the last message.

#### Scenario: All plugins stable
- **WHEN** all declared plugins have matching output within the window
- **THEN** the last message SHALL be returned unchanged (no state block prepended)

### Requirement: max_age zero disables change detection
When `MaxAge()` returns 0, the injector SHALL inject all declared plugins on every turn without scanning history.

#### Scenario: max_age is 0
- **WHEN** `parameters[0]["max_age"]` is 0
- **THEN** all declared plugins SHALL be injected unconditionally, equivalent to prior always-inject behavior
