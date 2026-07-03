## Context

Rubato's plugin system supports ambient state injection via `rubato:anchor` blocks. The `git_status` plugin demonstrated the pattern: per-request execution of a cheap command, compact key-value output, injected as a `rubato:state` block before the model sees the request.

The `go_test` plugin follows the same contract. Go's built-in test cache (`go test -count=0` bypass aside) means repeated runs of `go test ./...` return instantly when no source files have changed, making per-request execution feasible.

This design covers an MVP: run all tests, summarise as pass/fail/skip counts, include failure details when tests fail. No coverage, no benchmarks, no build cache warming.

## Goals / Non-Goals

**Goals:**
- Implement `go_test` plugin conforming to the `Plugin` interface
- Run `go test ./...` in a configurable `working_dir` (defaults to process CWD)
- On all-pass: single compact summary line (`tests: pass, N ran, M cached`)
- On failure: summary line + per-failure block (test name, failure message)
- Apply a configurable timeout (default 60s) to bound worst-case execution
- Register the plugin in `cmd/rubato/main.go`
- Full unit test coverage with a real Go module fixture

**Non-Goals:**
- Coverage reporting
- Benchmark output
- Build error vs. test failure distinction (both surface as failure)
- Per-package granularity in the summary
- Integration with any CI system

## Decisions

### D-1) Use `go test -json ./...` for structured output

`go test -json` emits newline-delimited JSON events (Action: run/pass/fail/skip/output). Parsing this is more robust than scraping text output — failure messages are associated with their test name without regex fragility.

Alternative: parse `go test -v ./...` text output. Rejected — format is less stable and harder to extract failure details from without false positives.

### D-2) Timeout default of 60 seconds

`git_status` uses 5s. Tests may legitimately take longer on cache miss. 60s is long enough for a medium project's first run and short enough to not hang the proxy indefinitely. Configurable via `args["timeout_seconds"]`.

### D-3) Output format matches `git_status` key-value style

Pass case:
```
tests: pass
ran: 12
cached: 47
```

Failure case:
```
tests: fail
ran: 15
cached: 44
failed: 2
--- FAIL: TestFoo/bar
    foo_test.go:42: expected 1, got 2
--- FAIL: TestBaz
    baz_test.go:17: nil pointer dereference
```

This keeps the `rubato:state` block skimmable and consistent with the existing plugin output style.

### D-4) Failure output truncation

Individual failure output is capped at 20 lines per test to prevent a single verbose failure from flooding the injected state block. Truncation is noted inline.

## Risks / Trade-offs

- **Slow on cache miss**: First run after a `go clean -testcache` or large change can take seconds. Mitigated by the 60s timeout and the fact that on-change injection means this cost is only paid when tests actually changed.
- **Working directory assumption**: Plugin assumes the working dir is a Go module root. If `working_dir` points to a subdirectory without `go.mod`, `go test ./...` fails. The error propagates as a plugin execution error (proxy returns 500). Acceptable for MVP.
- **Build failures**: If a package doesn't compile, `go test -json` emits build error output without a clean FAIL event. The plugin will report `tests: fail` with the build output as failure detail — correct behavior, slightly different formatting.
