## ADDED Requirements

### Requirement: Operational Decision-Path Logging
The system SHALL provide structured logs that describe runtime-injection decisions and failures without unnecessary prompt-content leakage at non-debug levels.

#### Scenario: Decision-path visibility
- **WHEN** Rubato processes a request
- **THEN** logs include whether an anchor was detected, which plugins were declared, and whether execution succeeded or failed

#### Scenario: Non-debug content safety
- **WHEN** logging level is non-debug
- **THEN** logs avoid unnecessary request prompt/body content

### Requirement: Runtime Configuration Guardrails
The system SHALL enforce configuration bounds for plugin execution and runtime path settings.

#### Scenario: Timeout bound enforcement
- **WHEN** plugin timeout configuration is invalid or out of bounds
- **THEN** Rubato fails startup or request handling with explicit configuration errors

#### Scenario: Repo path constraint enforcement
- **WHEN** configured repository path is invalid or inaccessible
- **THEN** Rubato produces explicit configuration/runtime failures

### Requirement: Routed Provider Integration
The system SHALL provide configuration integration so opencode model traffic is routed through Rubato in the devcontainer workflow.

#### Scenario: .opencode routing configured
- **WHEN** devcontainer workflow configuration is applied
- **THEN** `.opencode` routes model traffic through Rubato for end-to-end runtime behavior verification

### Requirement: Contract Freeze Evidence
The system SHALL freeze and document canonical API behavior using implementation traces before final polish completion.

#### Scenario: Evidence-backed freeze
- **WHEN** Stage C freeze activity is completed
- **THEN** canonical request/response behaviors are captured and reflected in design/spec artifacts

### Requirement: End-To-End Regression Verification
The system SHALL pass routed end-to-end verification in the devcontainer workflow without regressing Stage A/B behavior.

#### Scenario: Devcontainer routed verification
- **WHEN** end-to-end verification is executed with opencode traffic routed through Rubato
- **THEN** expected injection and failure behaviors pass and no known Stage A/B regressions are introduced
