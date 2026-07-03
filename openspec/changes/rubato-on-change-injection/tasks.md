## 1. Anchor Format Redesign

- [ ] 1.1 Define `Option struct { Name string; Setting any }` and `PluginDescriptor struct { Plugin string; Options []Option }` in the anchor package.
- [ ] 1.2 Update `Block`: replace `Plugins []string` with `Plugins []PluginDescriptor`; replace `Args map[string]map[string]any` with `Options []Option` (top-level rubato options).
- [ ] 1.3 Update anchor parser to parse `plugins` as `[]PluginDescriptor` and top-level `options` as `[]Option`.
- [ ] 1.4 Add `MaxAge() int` method on `Block`: scans `Options` for `name == "max_age"`, returns `setting` as int; defaults to 100; returns 0 when explicitly set to 0.
- [ ] 1.5 Update all existing anchor tests to use the new plugin descriptor format.
- [ ] 1.6 Test: plugin descriptor with no options parses correctly.
- [ ] 1.7 Test: plugin descriptor with options returns correct per-plugin option values.
- [ ] 1.8 Test: top-level options absent returns MaxAge 100.
- [ ] 1.9 Test: top-level `max_age` 0 returns MaxAge 0.
- [ ] 1.10 Test: unknown option names are preserved without error.
- [ ] 1.11 Update `Plugin` interface: change `Execute(ctx, args map[string]any)` to `Execute(ctx, options []anchor.Option) (string, error)`.
- [ ] 1.12 Update `git_status` and `go_test` plugin implementations to accept `[]anchor.Option`; add helper to extract a named option's setting (with type coercion for JSON float64).
- [ ] 1.13 Update `Registry.Execute` signature to pass `[]anchor.Option` per plugin.
- [ ] 1.14 Update all plugin tests for the new Execute signature.
- [ ] 1.15 Update `mutate` package to extract per-plugin options from descriptors and pass `[]anchor.Option` to Execute.
- [ ] 1.16 Update all existing `rubato:anchor` blocks in smoke fixtures, specs, and README examples to the new format.

## 2. Backward State Scan

- [ ] 2.1 In `mutate` package, implement `scanPriorState(msgs []json.RawMessage, maxAge int) map[string]string` — scans backward through `msgs[0..len-2]` up to `maxAge` positions, parses `rubato:state` fences, returns last known output per plugin name.
- [ ] 2.2 Parser reads `` ```rubato:state `` open fence and `[plugin_name]` section headers; accumulates lines per plugin until next section or close fence.
- [ ] 2.3 Test: single prior state block returns correct per-plugin outputs.
- [ ] 2.4 Test: multiple prior state blocks — most recent wins.
- [ ] 2.5 Test: scan stops at maxAge boundary — older blocks beyond window are ignored.
- [ ] 2.6 Test: no prior state blocks returns empty map.

## 3. On-Change Injection Logic

- [ ] 3.1 In `mutate.Apply`, after executing plugins, call `scanPriorState` with `block.MaxAge()`.
- [ ] 3.2 Build inject list: plugins whose fresh output differs from scanned output, or not found in scan.
- [ ] 3.3 When `MaxAge() == 0`, skip scan and inject all plugins unconditionally.
- [ ] 3.4 When inject list is empty, skip state block prepend entirely.
- [ ] 3.5 Build state block from inject list only (not all declared plugins).
- [ ] 3.6 Test: first turn (no history) — all plugins injected.
- [ ] 3.7 Test: stable turn — inject list empty, no state block prepended.
- [ ] 3.8 Test: one plugin changes — only that plugin in state block.
- [ ] 3.9 Test: plugin beyond max_age window — re-injected regardless of output match.
- [ ] 3.10 Test: max_age 0 — all plugins injected unconditionally.
