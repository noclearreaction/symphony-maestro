## Context

The rubato plugin registry is wired in `cmd/rubato/main.go`. Plugins have no standalone entry point — the only way to invoke one is through the proxy. A lightweight CLI binary in `cmd/rplugin` can instantiate the same registry and run a single named plugin, making plugin behavior directly observable.

## Goals / Non-Goals

**Goals:**
- Minimal CLI binary: plugin name as positional arg, `--working-dir` flag, `--args` JSON flag
- Stdout = plugin output, stderr = errors, exit 0/1
- Uses the same plugin instances as the proxy (no divergence)
- README documents the invocation

**Non-Goals:**
- Shell completion, man pages, or rich CLI framework
- Running multiple plugins in one invocation
- Streaming output

## Decisions

### D-1) Plugin registry is instantiated directly in main.go

Rather than extracting a shared `newRegistry()` helper into a library package, `cmd/rplugin/main.go` instantiates the plugins directly (same two lines as `cmd/rubato/main.go`). This avoids introducing a shared package for what is currently two lines of wiring.

Revisit if a third consumer appears.

### D-2) `--working-dir` is a first-class flag; arbitrary args via `--args` JSON

`working_dir` is the most common plugin arg. A dedicated flag improves ergonomics:
```
go run ./cmd/rplugin git_status --working-dir /workspace
go run ./cmd/rplugin go_test --working-dir /workspace --args '{"timeout_seconds":30}'
```
`--args` is a JSON object merged with the flag-derived args; explicit flags take precedence.

### D-3) No test file for rplugin main

`cmd/rplugin/main.go` is a thin wiring layer (arg parse → plugin.Execute → print). The plugin implementations are already tested. A test that exercises the binary end-to-end would duplicate smoke test infrastructure without adding coverage value.

## Risks / Trade-offs

- **Registry divergence**: if a plugin is added to `cmd/rubato` but not `cmd/rplugin`, the binary silently omits it. Mitigated by the two-line wiring being easy to review. Revisit with a shared registry helper if the plugin count grows.
