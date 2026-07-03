## 1. Operational Logging Refinement

- [ ] 1.1 Implement structured logs for decision path events: anchor detection, declared plugins, execution outcomes.
- [ ] 1.2 Ensure non-debug logs avoid unnecessary prompt/body content leakage.
- [ ] 1.3 Validate error logs preserve actionable failure context.

## 2. Runtime Configuration Hardening

- [ ] 2.1 Add and validate timeout bounds for plugin execution.
- [ ] 2.2 Document default values and failure behavior for invalid config.
- [ ] 2.3 Add or update `.opencode` configuration to route model traffic through Rubato in devcontainer workflows.

## 3. End-To-End Verification

- [ ] 3.1 Run devcontainer-routed end-to-end requests through Rubato and verify expected outcomes.
- [ ] 3.2 Run regression suite to ensure no Stage A/B behavior regressions.
- [ ] 3.3 Record final verification evidence and close remaining hygiene items.

## 4. On-Change State Injection

- [ ] 4.1 Add `Parameters []map[string]any` field to `anchor.Block`; parse from top-level `parameters` array in anchor JSON. Extract `max_age` (default 100, 0 = always inject) from first entry.
- [ ] 4.2 Update anchor tests to cover `repeat` field parsing.
- [ ] 4.3 In `mutate.Apply`, scan backward through up to `block.Repeat` prior messages for `rubato:state` blocks; extract last-known output per plugin.
- [ ] 4.4 Per-plugin: inject if not found in window or output changed; skip if found and unchanged.
- [ ] 4.5 Build state block containing only plugins to inject; skip prepend entirely when none.
- [ ] 4.6 Test: first turn injects all plugins.
- [ ] 4.7 Test: stable turn (all outputs match) injects nothing.
- [ ] 4.8 Test: one plugin changes, other unchanged — only changed plugin appears in state block.
- [ ] 4.9 Test: plugin beyond repeat window is re-injected regardless of output match.

## 5. Plugin Runner CLI

- [ ] 5.1 Create `cmd/rplugin/main.go` — binary that runs a single named plugin and writes output to stdout.
- [ ] 5.2 Accept plugin name as positional argument; exit 1 with usage message if absent.
- [ ] 5.3 Accept `--working-dir` flag (passed as plugin arg `working_dir`); accept `--args` flag for arbitrary JSON plugin args.
- [ ] 5.4 Wire the same plugin registry as `cmd/rubato`; write plugin output to stdout, errors to stderr, exit 1 on failure.
- [ ] 5.5 Document usage in `cmd/rubato/README.md`.
