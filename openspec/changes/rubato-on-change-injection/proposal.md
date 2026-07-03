## Why

Rubato currently injects a full `rubato:state` block on every proxied request turn. In stable sessions this is redundant noise — the model already knows the git branch and test status from the previous turn. Injecting only what changed reduces token use and gives the model a useful signal: a new state block means something is different this turn.

## What Changes

- Extend `rubato:anchor` with an optional `options` array for rubato-level config, starting with `max_age` (the maximum number of messages before a plugin's state is considered stale and must be re-injected regardless of change).
- On each request, scan backward through prior messages (bounded by `max_age`) to find the last injected output per plugin.
- Inject only the plugins whose output has changed or that have not been seen within the `max_age` window.
- When no plugins require injection, skip the state block entirely.
- Per-plugin atomicity: each plugin's output is an independent unit; one plugin changing does not force re-injection of unchanged plugins.

## Capabilities

### New Capabilities

- `anchor-options`: The `rubato:anchor` block accepts a top-level `options` array of `{key, value}` objects for rubato-level config. Initially supports `max_age` (default 100, `0` = always inject).
- `on-change-injection`: Rubato injects only changed or stale plugin outputs per turn, rather than the full state block every turn.

### Modified Capabilities

- `rubato-proxy-injection`: Injection behavior changes from always-inject to on-change-inject. The `rubato:state` block is now partial (only changed/stale plugins) or absent (all stable within window).

## Impact

- `internal/rubato/anchor/anchor.go` and `anchor_test.go`: parse `options` array, expose `MaxAge()`
- `internal/rubato/mutate/mutate.go` and `mutate_test.go`: backward scan, per-plugin diff, partial state block
- No changes to plugin implementations, proxy handler, or main.go
- No breaking changes to existing anchors — absent `parameters` defaults to `max_age: 100`
