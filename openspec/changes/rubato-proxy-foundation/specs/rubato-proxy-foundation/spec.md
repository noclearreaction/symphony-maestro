## ADDED Requirements

### Requirement: Rubato Go Service Foundation
The system SHALL provide Rubato as a first-class Go service in this repository with a standard package layout and executable entrypoint.

#### Scenario: Standard project layout exists
- **WHEN** the Rubato source tree is inspected
- **THEN** it includes `cmd/`, `internal/`, `pkg/`, and `test/` directories under module root `rubato/`

#### Scenario: Executable entrypoint exists
- **WHEN** Rubato is built from source
- **THEN** an executable entrypoint is produced from `cmd/rubato`

### Requirement: Pass-Through Proxy Baseline
The system SHALL forward eligible chat-completions requests to upstream without prompt mutation.

#### Scenario: Pass-through relay
- **WHEN** a valid chat-completions request is received
- **THEN** Rubato forwards the request upstream without changing `messages[0]` or `messages[-1]`
- **AND** Rubato relays the upstream response back to the client

### Requirement: Baseline Deterministic Errors
The system SHALL return deterministic errors for malformed requests and forwarding failures.

#### Scenario: Malformed request body
- **WHEN** request JSON cannot be parsed
- **THEN** Rubato returns a deterministic client-error response and does not forward upstream

#### Scenario: Upstream forwarding failure
- **WHEN** upstream request execution fails
- **THEN** Rubato returns a deterministic upstream-failure response

### Requirement: Foundation Test Coverage
The system SHALL include automated tests for routing, request validation, and pass-through behavior.

#### Scenario: Foundation tests run
- **WHEN** Rubato tests are executed
- **THEN** unit and component tests validate route handling, deterministic errors, and pass-through forwarding behavior
