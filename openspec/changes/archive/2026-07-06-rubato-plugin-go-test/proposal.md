## Why

Agents running in the rubato proxy frequently need to know the current unit test status before reasoning about code correctness, suggesting commits, or planning fixes. Without ambient injection, they must issue an explicit tool call or ask the user to run tests — adding a turn and breaking the flow. The `git_status` plugin demonstrated the value of ambient state injection; `go_test` is the natural next plugin given Go's fast cached test execution.

## What Changes

- Add a new `go_test` plugin to the rubato plugin registry
- Plugin runs `go test ./...` in a configurable working directory
- On success: outputs a compact summary (pass count, cached count, elapsed)
- On failure: outputs per-test failure names and messages, skipping passing output
- Plugin is opt-in via `rubato:anchor` — no change to existing behavior

## Capabilities

### New Capabilities

- `go-test-plugin`: Ambient injection of Go unit test status — pass/fail summary with failure details — into every proxied request that declares the `go_test` plugin in its `rubato:anchor` block.

### Modified Capabilities

<!-- none -->

## Impact

- New file: `internal/rubato/plugin/gotest.go`
- New test file: `internal/rubato/plugin/gotest_test.go`
- `cmd/rubato/main.go`: register `go_test` plugin alongside `git_status`
- No changes to proxy, mutate, or anchor packages
- No breaking changes
