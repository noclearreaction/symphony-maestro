## Context

Rubato's `mutate.Apply` currently executes all declared plugins and prepends a full `rubato:state` block on every request. The full message history is available in each request body, so rubato can determine what has changed since the last injected state. Per-plugin output is treated as an atomic unit for diffing purposes.

The `anchor.Block` struct currently holds `Plugins []string` and `Args map[string]map[string]any`. A `Parameters` field is added to carry rubato-level config that is not per-plugin.

## Goals / Non-Goals

**Goals:**
- Extend anchor parsing to support a top-level `parameters` array
- Expose `MaxAge int` derived from `parameters[0]["max_age"]`; default 100; 0 = always inject
- Scan backward through prior messages (bounded by `MaxAge`) for last injected output per plugin
- Inject only plugins that are new, changed, or beyond the `MaxAge` window
- Omit the state block entirely when no plugins need injection
- Full test coverage for all injection decision paths

**Non-Goals:**
- Changes to plugin implementations
- Guidance text changes (guidance remains silent on absence semantics)
- New anchor semantics beyond `parameters`/`max_age`

## Decisions

### D-1) `options` is an array of explicit `{key, value}` objects

```json
{"plugins":["git_status"],"options":[{"name":"max_age","setting":50}]}
```

Rationale: each option is a self-describing name/setting pair. Parsers scan the array for entries where `name` matches a known option and read `setting`; unknown names are ignored. This avoids the implicit single-key-per-object convention and makes the structure unambiguous regardless of what future options are added.

### D-2) `max_age: 0` means always inject

Zero is the natural "disable" value for an age threshold. It preserves the current always-inject behavior as an explicit opt-in rather than a magic string.

### D-3) Inject only changed/stale plugins (Option A — partial block)

When `git_status` changes but `go_test` does not, the injected block contains only `[git_status]`. Absence signals "unchanged since you last saw it". Guidance text does NOT explain this — any explicit hint risks triggering the model to use tools instead of trusting ambient state.

### D-4) Scan is bounded by MaxAge, backward from messages[-2]

`messages[-1]` is the current user turn being mutated. The scan starts at `messages[-2]` and walks backward up to `MaxAge` positions. A plugin not found within this window is treated as stale and injected unconditionally.

### D-5) State block parsing uses fence markers only

Backward scan looks for `` ```rubato:state `` open fences and parses `[plugin_name]` section headers to extract per-plugin output. No JSON parsing required — the format is line-oriented.

## Risks / Trade-offs

- **Partial context**: model must recall plugin state from earlier turns when output is stable. Mitigated by `max_age` forcing a refresh within a bounded window.
- **Scan cost**: O(MaxAge × message_size) per request. At default MaxAge 100 with typical message sizes this is negligible.
- **Parse fragility**: state block parsing relies on the fence/section-header format being stable. Mitigated by the format being owned by rubato itself.
