## Why

Debugging what a rubato plugin will inject requires running the full proxy, sending a crafted request, and reading the log. A standalone `rplugin` binary lets developers run any registered plugin directly from the terminal and see its output immediately — useful for verifying plugin behavior, testing args, and confirming ambient state before a session.

## What Changes

- Add `cmd/rplugin/main.go`: a CLI binary that accepts a plugin name, optional flags, runs the plugin, and writes output to stdout.
- Document usage in the rubato README.

## Capabilities

### New Capabilities

- `rplugin-cli`: A `go run ./cmd/rplugin` binary that executes a single named rubato plugin and prints its output to stdout.

### Modified Capabilities

<!-- none -->

## Impact

- New file: `cmd/rplugin/main.go`
- Updated: `cmd/rubato/README.md`
- No changes to proxy, plugin, mutate, or anchor packages
- Shares the plugin registry with `cmd/rubato` via a shared wiring helper or direct instantiation
