## 1. Binary Implementation

- [ ] 1.1 Create `cmd/rplugin/main.go` with `flag` parsing: positional plugin name, `--working-dir` string flag, `--args` JSON string flag.
- [ ] 1.2 Print usage to stderr and exit 1 when no plugin name is provided.
- [ ] 1.3 Instantiate plugin registry with `NewGitStatus()` and `NewGoTest()`.
- [ ] 1.4 Merge `--args` JSON object with `--working-dir` into plugin args map; `--working-dir` takes precedence.
- [ ] 1.5 Execute the named plugin; write output to stdout on success.
- [ ] 1.6 On unknown plugin name, execution error, or invalid `--args` JSON: write error to stderr and exit 1.

## 2. Documentation

- [ ] 2.1 Add `rplugin` usage section to `cmd/rubato/README.md` with example invocations for `git_status` and `go_test`.
