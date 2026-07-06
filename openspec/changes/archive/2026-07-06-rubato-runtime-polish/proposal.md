## Why

After proxy foundation and plugin MVP are delivered, we need a bounded polish phase to improve operability, enforce final contract details, and verify behavior stability without expanding product scope.

## What Changes

- Add structured, privacy-aware operational logging aligned to staged behavior.
- Add runtime configuration hardening (timeouts, repo path bounds, logging controls).
- Add/verify `.opencode` routing configuration so model traffic is consistently routed through Rubato in devcontainer workflows.
- Execute Task 7 contract-freeze activities with captured traces and final artifact alignment.
- Add end-to-end verification for routed opencode-to-rubato workflows in devcontainer.

## Capabilities

### New Capabilities

- `rubato-runtime-polish`: operational hardening, contract freeze evidence, and final verification.

### Modified Capabilities

- `rubato-plugin-git-status`: adds operational constraints and verification coverage without changing core behavior.

## Impact

- Affected systems: logging, runtime configuration surface, verification workflow.
- Affected behavior: no functional scope expansion; behavior should be clarified and hardened only.
- Dependencies: `rubato-proxy-foundation` and `rubato-plugin-git-status` must be complete first.