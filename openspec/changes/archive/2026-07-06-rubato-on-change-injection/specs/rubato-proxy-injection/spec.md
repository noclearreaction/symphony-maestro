## MODIFIED Requirements

### Requirement: Injection pipeline applies on-change semantics
The proxy injection pipeline SHALL inject plugin state selectively per turn based on change detection and the `max_age` window, rather than always injecting the full state block.

#### Scenario: Partial state block on mixed-change turn
- **WHEN** one plugin's output has changed and another's has not changed within the window
- **THEN** the injected `rubato:state` block SHALL contain only the changed plugin's output
- **THEN** the unchanged plugin SHALL NOT appear in the state block for that turn

#### Scenario: No state block on fully stable turn
- **WHEN** all declared plugins are stable within the window
- **THEN** no `rubato:state` block SHALL be prepended to the last message
- **THEN** the request body SHALL be forwarded with only the guidance block (if first turn) and no state mutation
