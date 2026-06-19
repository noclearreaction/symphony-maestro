## 1. Go Project Foundation

- [ ] 1.1 Create `tools/rubato/` Go module and baseline tooling files.
- [ ] 1.2 Scaffold standard directories: `cmd/`, `internal/`, `pkg/`, `test/`.
- [ ] 1.3 Implement `cmd/rubato` entrypoint and runtime config loading.

## 2. Minimal Proxy Path

- [ ] 2.1 Implement HTTP server with explicit route/method handling for chat completions.
- [ ] 2.2 Implement pass-through forwarding to configured upstream without mutating messages.
- [ ] 2.3 Add deterministic request and upstream failure responses.

## 3. Tests And Verification

- [ ] 3.1 Add unit tests for handler validation and routing behavior.
- [ ] 3.2 Add component tests using upstream test doubles for pass-through relay behavior.
- [ ] 3.3 Add baseline build/test commands and verify they run cleanly in devcontainer workflow.