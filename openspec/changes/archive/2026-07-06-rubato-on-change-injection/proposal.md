## Why

Rubato currently injects a full `rubato:state` block on every proxied request turn. In stable sessions this is redundant noise — the model already knows the git branch and test status from the previous turn. Injecting only what changed reduces token use and gives the model a useful signal: a new state block means something is different this turn.

## What Changes

- Redesign `rubato:anchor` to use plugin descriptor objects (`{plugin, options}`) instead of a flat string array, co-locating each plugin's options with its declaration.
- Add a top-level `options` array for rubato-level config using the same `{name, setting}` format, starting with `max_age`.
- Remove the separate top-level `args` object.
- On each request, scan backward through prior messages (bounded by `max_age`) to find the last injected output per plugin.
- On each request, scan backward through prior messages (bounded by `max_age`) to find the last injected output per plugin.
- Inject only the plugins whose output has changed or that have not been seen within the `max_age` window.
- When no plugins require injection, skip the state block entirely.
- Per-plugin atomicity: each plugin's output is an independent unit; one plugin changing does not force re-injection of unchanged plugins.

## Capabilities

### New Capabilities

- `anchor-options`: The `rubato:anchor` block accepts plugin descriptors (`{plugin, options}`) and a top-level `options` array, both using `{name, setting}` pairs. Replaces the current string-array `plugins` and flat `args` object.
- `on-change-injection`: Rubato injects only changed or stale plugin outputs per turn, rather than the full state block every turn.

### Modified Capabilities

- `rubato-proxy-injection`: Injection behavior changes from always-inject to on-change-inject. The `rubato:state` block is now partial (only changed/stale plugins) or absent (all stable within window).

## Impact

- `internal/rubato/anchor/anchor.go` and `anchor_test.go`: new `PluginDescriptor` and `Option` types, updated `Block` struct (removes `Args`, adds `Options`), updated parser
- `internal/rubato/mutate/mutate.go` and `mutate_test.go`: extract per-plugin options from descriptors, backward scan, partial state block
- Update all existing `rubato:anchor` blocks in specs, smoke fixtures, and README examples
