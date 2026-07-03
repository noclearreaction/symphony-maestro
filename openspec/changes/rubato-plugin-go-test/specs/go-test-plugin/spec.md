## ADDED Requirements

### Requirement: Plugin executes Go tests and returns a compact summary
The `go_test` plugin SHALL run `go test -json ./...` in the configured working directory and return a UTF-8 string summarising the outcome. The summary SHALL be suitable for injection into a `rubato:state` block.

#### Scenario: All tests pass
- **WHEN** `go test ./...` exits with code 0
- **THEN** the output SHALL contain `tests: pass`, the count of tests that ran, and the count of cached tests

#### Scenario: One or more tests fail
- **WHEN** `go test ./...` exits with non-zero code due to test failures
- **THEN** the output SHALL contain `tests: fail`, the total ran and cached counts, the number of failed tests, and the name and failure message for each failing test

#### Scenario: Working directory is specified via args
- **WHEN** `args["working_dir"]` is a non-empty string
- **THEN** the plugin SHALL run tests rooted at that directory

#### Scenario: Working directory defaults to process CWD
- **WHEN** `args["working_dir"]` is absent or empty
- **THEN** the plugin SHALL run tests rooted at the process working directory

### Requirement: Plugin applies a configurable execution timeout
The `go_test` plugin SHALL abort test execution if it exceeds the configured timeout and return an error. The default timeout SHALL be 60 seconds.

#### Scenario: Timeout exceeded
- **WHEN** `go test ./...` has not completed within the timeout
- **THEN** `Execute` SHALL return an error indicating the timeout was exceeded

#### Scenario: Custom timeout via args
- **WHEN** `args["timeout_seconds"]` is a positive number
- **THEN** the plugin SHALL use that value as the execution timeout in seconds

### Requirement: Plugin truncates verbose failure output
The `go_test` plugin SHALL cap individual test failure output at 20 lines. If output is truncated, the summary SHALL include a note indicating truncation.

#### Scenario: Failure output exceeds 20 lines
- **WHEN** a single failing test emits more than 20 lines of output
- **THEN** only the first 20 lines SHALL appear in the summary, followed by a truncation note

#### Scenario: Failure output within limit
- **WHEN** a single failing test emits 20 or fewer lines of output
- **THEN** all output lines SHALL appear in the summary unchanged

### Requirement: Plugin is registered in the rubato registry
The `go_test` plugin SHALL be registered in the plugin registry wired up in `cmd/rubato/main.go` so that it is available to any request that declares it in a `rubato:anchor` block.

#### Scenario: Plugin declared in anchor
- **WHEN** a request's system message contains a `rubato:anchor` block with `"plugins":["go_test"]`
- **THEN** the proxy SHALL execute the plugin and inject its output into the request before forwarding

#### Scenario: Plugin not declared
- **WHEN** a request's `rubato:anchor` block does not include `go_test`
- **THEN** the plugin SHALL NOT be executed and no test state SHALL be injected
