## Why

Runtime-state injection should be introduced only after Rubato proxy foundations are stable. This change isolates the first behavior layer: anchor-driven plugin execution with one MVP plugin (`git_status`).

## What Changes

- Add strict anchor parsing and eligibility checks for runtime injection.
- Add internal plugin contract and declared-plugin registry resolution.
- Implement fail-fast plugin execution semantics.
- Implement runtime-state block injection into `messages[-1]`.
- Implement deterministic, idempotent guidance augmentation in `messages[0]`.
- Implement MVP `git_status` plugin with explicit hygiene metrics and recognizable detached-HEAD/bare-repo states.

## Capabilities

### New Capabilities

- `rubato-plugin-git-status`: anchor parsing, plugin execution framework, runtime mutation semantics, MVP git plugin behavior.

### Modified Capabilities

- `rubato-proxy-foundation`: extends pass-through proxy with controlled mutation path for eligible requests.

## Impact

- Affected systems: Rubato request pre-processing, plugin execution path, request mutation path.
- Affected behavior: eligible requests can now fail fast or forward with injected runtime state.
- Dependencies: `rubato-proxy-foundation` change must be complete first.