## Why

The project needs a reliable way to inject current runtime repository state into AI turns without depending on host tooling layout or manual copy/paste. A dedicated rubato runtime-injection capability is needed now to restore business workflow velocity after the devcontainer isolation pivot.

This change is delivered in stages to reduce risk and keep issue/doc alignment clean:
- Stage A: minimal non-mutating runtime behavior.
- Stage B: MVP plugin-based injection behavior.
- Stage C: refinement and polish.

## What Changes

- Introduce a marker-controlled runtime injection capability in rubato that reads an anchor block from `messages[0]` and injects selected plugin outputs into `messages[-1]` before forwarding.
- Allow each declared plugin in the anchor to include static arguments so callers can tune plugin behavior per request/session without changing daemon configuration.
- Add cache-stable plugin-instruction augmentation in `messages[0]`: each declared plugin contributes a canonical usage/integration instruction block when missing from request content, and identical anchor/plugin inputs must produce byte-identical guidance across sessions.
- Start MVP with a `git_status` plugin focused on branch and hygiene signals (branch name, ahead/behind, commits ahead count, staged count, unstaged tracked-modified count, untracked count).
- Require fail-fast behavior: if a declared plugin cannot execute, the request fails with a clear user-facing error instead of silently degrading.
- Keep the runtime path stateless: plugin outputs are refreshed on every request and no per-session state is required.
- Define a plugin contract that anticipates multiple plugins while allowing MVP rollout with git only.
- Complete Task 7 contract-freeze work after Stage B and before Stage C.

## Capabilities

### New Capabilities
- `rubato-runtime-state-injection`: Marker/anchor driven runtime state injection for AI requests using plugin-selected context blocks.

### Modified Capabilities
- None.

## Impact

- Affected systems: rubato proxy request mutation path and request validation path.
- Affected configuration: opencode project config under `.opencode` to route provider traffic through rubato.
- Affected runtime behavior: request forwarding can now fail early when declared runtime plugins fail.
- Dependencies: existing OpenRouter-compatible request flow; no requirement to add non-stdlib Go dependencies for MVP.
