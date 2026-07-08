## 1. Plugin Implementation

- [x] 1.1 Create `internal/rubato/plugin/gotest.go` with `GoTest` struct implementing the `Plugin` interface
- [x] 1.2 Implement `Execute`: parse `working_dir` and `timeout_seconds` from args, default timeout 60s
- [x] 1.3 Run `go test -json ./...` via `exec.CommandContext`, capture combined output
- [x] 1.4 Parse JSON event stream: accumulate ran/cached/failed counts and per-test failure output
- [x] 1.5 Truncate individual failure output to 20 lines with a truncation note when exceeded
- [x] 1.6 Return compact pass summary (`tests: pass\nran: N\ncached: M`) on success
- [x] 1.7 Return fail summary with per-failure blocks on non-zero exit

## 2. Plugin Tests

- [x] 2.1 Create `internal/rubato/plugin/gotest_test.go`
- [x] 2.2 Test pass case: fixture module with one passing test returns `tests: pass` output
- [x] 2.3 Test fail case: fixture module with one failing test returns `tests: fail` with failure detail
- [x] 2.4 Test timeout: verify `Execute` returns error when timeout is exceeded
- [x] 2.5 Test truncation: failure output exceeding 20 lines is capped with truncation note
- [x] 2.6 Test non-module directory: `Execute` returns error for directory without `go.mod`
- [x] 2.7 Test default CWD: omitting `working_dir` runs tests in process working directory

## 3. Registration

- [x] 3.1 Register `plugin.NewGoTest()` in the registry in `cmd/rubato/main.go` alongside `GitStatus`

## 4. Smoke Test Coverage

- [x] 4.1 Update `cmd/rubato/testdata/smoke/agents/smoke.md` anchor to include `go_test` in the plugins list
- [x] 4.2 Add assertion in `cmd/rubato/smoke_test.go` that `tests:` appears in the rubato log
