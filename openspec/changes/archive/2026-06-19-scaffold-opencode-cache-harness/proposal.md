## Why

Verifying opencode cache behavior empirically requires a controlled, reproducible environment. Without an isolated test harness, experiments are inconsistent across machines and sessions, making it impossible to trust token and cost measurements.

## What Changes

- New `harness/` directory in the repository with a minimal Dockerfile installing opencode
- Minimal project fixture: a single agent configuration with a ~200-token system prompt and no real application code
- `README.md` inside the harness describing how to build the image and start experimenting
- The specific mechanisms for triggering turns and observing cache/cost behavior are themselves goals of the spike, not preconditions

## Capabilities

### New Capabilities

- `opencode-cache-harness`: A self-contained Docker environment for running controlled opencode experiments; the mechanisms for triggering turns and reading cache/cost data will be discovered as part of the spike

### Modified Capabilities

## Impact

- Adds a new `docker/` (or `harness/`) directory to this repository
- Depends on Docker being available in the execution environment
- Unblocks SF-2 through SF-8 of the parent spike (#43)
- No changes to existing specs, governance, or agent instruction files
