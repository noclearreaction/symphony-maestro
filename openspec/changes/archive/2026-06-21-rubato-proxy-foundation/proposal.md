## Why

`rubato-runtime-state-injection` is too broad for clean execution and review. We need a narrow first change that establishes Rubato as a reliable proxy service and Go codebase foundation before any prompt mutation behavior.

## What Changes

- Create a production-grade Go project for Rubato at repository root using a standard layout with `cmd/`, `internal/`, `pkg/`, and `test/`, with Rubato code under package namespace `internal/rubato/`.
- Implement a minimal HTTP proxy path for chat-completions forwarding with no runtime injection behavior.
- Add baseline request validation, upstream forwarding, and deterministic error responses for malformed HTTP requests and upstream failures.
- Add unit tests and focused component tests for proxy pass-through behavior.
- Establish build and test commands for Rubato so follow-on plugin work lands on stable foundations.

## Capabilities

### New Capabilities

- `rubato-proxy-foundation`: Rubato executable, proxy routing, baseline request validation, and testable pass-through behavior.

### Modified Capabilities

- None.

## Impact

- Affected systems: repository-root Go source layout, `internal/rubato/` package namespace, local build/test workflow, proxy runtime entrypoint.
- Affected behavior: requests are forwarded unchanged (no injection yet).
- Dependencies: existing upstream-compatible chat-completions API contract.