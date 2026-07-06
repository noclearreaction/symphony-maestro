## Context

Rubato's `mutate.Apply` currently executes all declared plugins and prepends a full `rubato:state` block on every request. The full message history is available in each request body, so rubato can determine what has changed since the last injected state. Per-plugin output is treated as an atomic unit for diffing purposes.

The `anchor.Block` struct currently holds `Plugins []string` and `Args map[string]map[string]any`. This change redesigns the anchor format, the `Block` struct, and the `Plugin` interface to use a consistent `Option` abstraction throughout.

## Goals / Non-Goals

**Goals:**
- Redesign anchor format to plugin descriptors with co-located options
- Introduce `Option{Name, Setting}` as the shared abstraction for both plugin options and rubato-level config
- Change `Plugin.Execute` to accept `[]anchor.Option` instead of `map[string]any`
- Expose `MaxAge()` from top-level anchor options (default 100, 0 = always inject)
- Scan backward through prior messages (bounded by `MaxAge`) for last injected output per plugin
- Inject only plugins that are new, changed, or beyond the `MaxAge` window
- Omit the state block entirely when no plugins need injection
- Full test coverage for all injection decision paths

**Non-Goals:**
- Guidance text changes (guidance remains silent on absence semantics)
- New anchor semantics beyond `options`/`max_age`

## Decisions

### D-1) Unified anchor format — plugin descriptors with co-located options

```json
{
  "plugins": [
    {"plugin": "git_status"},
    {"plugin": "go_test", "options": [{"name": "timeout_seconds", "setting": 30}]}
  ],
  "options": [{"name": "max_age", "setting": 50}]
}
```

`plugins` becomes an array of descriptor objects, each with a `plugin` name and an optional `options` array of `{name, setting}` pairs — `setting` itself is optional, allowing flag-style options with no value. Top-level `options` carry rubato-level config. The existing `args` top-level key is removed.

Rationale: plugin name and its config are co-located; the format is consistent throughout; no separate args lookup. The `{name, setting}` shape is used uniformly for both plugin options and rubato options. This is a breaking change to the current string-array `plugins` format — acceptable since nothing is deployed.

### D-2) `max_age: 0` means always inject

Zero is the natural "disable" value for an age threshold. It preserves the current always-inject behavior as an explicit opt-in rather than a magic string.

### D-3) Inject only changed/stale plugins (Option A — partial block)

When `git_status` changes but `go_test` does not, the injected block contains only `[git_status]`. Absence signals "unchanged since you last saw it". Guidance text does NOT explain this — any explicit hint risks triggering the model to use tools instead of trusting ambient state.

### D-4) Scan is bounded by MaxAge, backward from messages[-2]

`messages[-1]` is the current user turn being mutated. The scan starts at `messages[-2]` and walks backward up to `MaxAge` positions. A plugin not found within this window is treated as stale and injected unconditionally.

### D-5) State block parsing uses fence markers only

Backward scan looks for `` ```rubato:state `` open fences and parses `[plugin_name]` section headers to extract per-plugin output. No JSON parsing required — the format is line-oriented.

### D-6) `Option` is defined in `anchor`; `plugin` imports `anchor`

The `Option{Name string; Setting any}` type is defined in the `anchor` package. The `Plugin` interface's `Execute` method is updated to `Execute(ctx context.Context, options []anchor.Option) (string, error)`. The `plugin` package imports `anchor`; `mutate` imports both.

Dependency graph:
```
anchor  →  (nothing)
plugin  →  anchor
mutate  →  anchor, plugin
```

Rationale: anchor owns the wire format; plugins are consumers of it. `setting` is `any` and typed by the plugin implementation. `float64` coercion (from JSON number decoding) is handled by a helper in the anchor package.

## Risks / Trade-offs

- **Partial context**: model must recall plugin state from earlier turns when output is stable. Mitigated by `max_age` forcing a refresh within a bounded window.
- **Scan cost**: O(MaxAge × message_size) per request. At default MaxAge 100 with typical message sizes this is negligible.
- **Parse fragility**: state block parsing relies on the fence/section-header format being stable. Mitigated by the format being owned by rubato itself.
