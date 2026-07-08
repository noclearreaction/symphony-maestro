## ADDED Requirements

### Requirement: Plugin executes Go tests and returns a structured summary
The `go_test` plugin SHALL run `go test -json ./...` in the configured working directory and return a UTF-8 string summarising the outcome. The summary SHALL be suitable for injection into a `rubato:state` block.

The first line SHALL always be `status: <value>` where value is one of:
- `pass` — all tests passed
- `fail` — one or more tests failed or one or more packages failed to build
- `error` — infrastructure failure (module not found, setup broken); no tests could run

#### Scenario: All tests pass
- **WHEN** `go test ./...` exits with code 0
- **THEN** the output SHALL start with `status: pass`
- **THEN** the output SHALL contain `ran: N`, `passed: N`, and `failed: 0` counts

#### Scenario: One or more tests fail
- **WHEN** `go test ./...` exits with non-zero code due to test failures
- **THEN** the output SHALL start with `status: fail`
- **THEN** the output SHALL contain `ran: N`, `passed: N`, and `failed: N` counts
- **THEN** the output SHALL contain a `test failures:` section
- **THEN** each failing test SHALL be listed under its package path with its name and failure message
- **THEN** Go test framework header lines (`=== RUN`, `--- FAIL:`, `--- PASS:`, timing) SHALL be omitted

#### Scenario: One or more packages fail to build
- **WHEN** `go test ./...` exits with non-zero code due to build errors
- **THEN** the output SHALL start with `status: fail`
- **THEN** the output SHALL contain `ran: N`, `passed: N`, and `failed: N` counts reflecting any tests that did execute
- **THEN** the output SHALL contain a `build errors:` section
- **THEN** each package with build errors SHALL appear as a labelled group containing its compiler error lines (`file:line: message`)

#### Scenario: Mixed — some packages pass, some fail to build
- **WHEN** `go test ./...` exits with non-zero code and some packages built and ran while others failed to build
- **THEN** the output SHALL start with `status: fail`
- **THEN** `passed: N` SHALL reflect the count of tests that passed in the packages that did build
- **THEN** the `build errors:` section SHALL list only the packages that failed to build

#### Scenario: Infrastructure failure (module not found, setup broken)
- **WHEN** `go test ./...` cannot resolve the package pattern (e.g. no `go.mod` found)
- **THEN** the output SHALL start with `status: error`
- **THEN** the output SHALL contain the error message from the Go toolchain

#### Scenario: Working directory is specified via args
- **WHEN** `args["working_dir"]` is a non-empty string
- **THEN** the plugin SHALL run tests rooted at that directory

#### Scenario: Working directory defaults to process CWD
- **WHEN** `args["working_dir"]` is absent or empty
- **THEN** the plugin SHALL run tests rooted at the process working directory

### Requirement: Plugin applies a configurable execution timeout
The `go_test` plugin SHALL abort test execution if it exceeds the configured timeout and return a hard error. The default timeout SHALL be 60 seconds.

#### Scenario: Timeout exceeded
- **WHEN** `go test ./...` has not completed within the timeout
- **THEN** `Execute` SHALL return an error indicating the timeout was exceeded

#### Scenario: Custom timeout via args
- **WHEN** `args["timeout_seconds"]` is a positive number
- **THEN** the plugin SHALL use that value as the execution timeout in seconds

### Requirement: Plugin is registered in the rubato registry
The `go_test` plugin SHALL be registered in the plugin registry wired up in `cmd/rubato/main.go` so that it is available to any request that declares it in a `rubato:anchor` block.

#### Scenario: Plugin declared in anchor
- **WHEN** a request's system message contains a `rubato:anchor` block with `"plugins":["go_test"]`
- **THEN** the proxy SHALL execute the plugin and inject its output into the request before forwarding

#### Scenario: Plugin not declared
- **WHEN** a request's `rubato:anchor` block does not include `go_test`
- **THEN** the plugin SHALL NOT be executed and no test state SHALL be injected
