## 1. Anchor Parameters Parsing

- [ ] 1.1 Add `Parameters []map[string]any` field to `anchor.Block`; parse from top-level `parameters` array in anchor JSON.
- [ ] 1.2 Add `MaxAge() int` method on `Block`: returns `parameters[0]["max_age"]` as int, defaulting to 100; returns 0 when explicitly set to 0.
- [ ] 1.3 Test: anchor with `parameters:[{"max_age":50}]` returns MaxAge 50.
- [ ] 1.4 Test: anchor without `parameters` returns MaxAge 100.
- [ ] 1.5 Test: anchor with `max_age:0` returns MaxAge 0.
- [ ] 1.6 Test: unknown parameter keys are preserved and do not cause parse errors.

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
