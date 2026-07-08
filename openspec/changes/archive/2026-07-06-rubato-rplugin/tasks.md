## 1. Binary Implementation

- [x] 1.1 Create `cmd/rplugin/main.go` with `flag` parsing: positional plugin name, `--working-dir` string flag, `--args` JSON string flag.
- [x] 1.2 Print usage to stderr and exit 1 when no plugin name is provided.
- [x] 1.3 Instantiate plugin registry with `NewGitStatus()` and `NewGoTest()`.
- [x] 1.4 Merge `--args` JSON object with `--working-dir` into plugin args map; `--working-dir` takes precedence.
- [x] 1.5 Execute the named plugin; write output to stdout on success.
- [x] 1.6 On unknown plugin name, execution error, or invalid `--args` JSON: write error to stderr and exit 1.

## 2. Documentation

- [x] 2.1 Add `rplugin` usage section to `cmd/rubato/README.md` with example invocations for `git_status` and `go_test`.
