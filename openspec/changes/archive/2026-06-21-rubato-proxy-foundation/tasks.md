## 1. Go Project Foundation

- [x] 1.1 Create repository-root Go project and baseline tooling files, with Rubato under package namespace `internal/rubato/`.
- [x] 1.2 Scaffold standard directories: `cmd/`, `internal/`, `pkg/`, `test/`.
- [x] 1.3 Implement `cmd/rubato` entrypoint and runtime config loading.

## 2. Minimal Proxy Path

- [x] 2.1 Implement HTTP server with explicit route/method handling for chat completions.
- [x] 2.2 Implement pass-through forwarding to configured upstream without mutating messages.
- [x] 2.3 Add deterministic request and upstream failure responses.

## 3. Tests And Verification

- [x] 3.1 Add unit tests for handler validation and routing behavior.
- [x] 3.2 Add component tests using upstream test doubles for pass-through relay behavior.
- [x] 3.3 Add baseline build/test commands and verify they run cleanly in devcontainer workflow.