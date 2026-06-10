## Why

The issue-56b per-session-keyed NDJSON approach produces multiple files per proxy process — one per agent role plus one for title-gen. While this separates concerns, it is operationally awkward: you cannot `tail -f` a single file to observe all traffic, log rotation requires managing multiple files, and session boundaries within a file are invisible (two sessions with the same system prompt share a file with no separator).

A single log file per proxy startup with structured log levels is simpler to operate, ingest, and inspect. The session key and turn number become record fields rather than file names.

## What Changes

- Single append-only log file per proxy startup (named `proxy-<unix-timestamp-seconds>.log` in `LOG_DIR`, or path overridden by `LOG_FILE` env var)
- Structured log levels: `DEBUG`, `INFO`, `WARN`, `ERROR` — controlled by `LOG_LEVEL` env var (default: `info`)
- `INFO`: session key, turn number, model, upstream status code, token counts — no message content
- `DEBUG`: all INFO fields plus full original `messages[-1].content` and full injected `messages[-1].content` (when injection is active); full response body
- `WARN`/`ERROR`: error condition details only
- Remove per-session-keyed NDJSON file mechanism from issue-56b
- No checksums: original and injected content are logged directly in debug mode

## Capabilities

### Modified Capabilities

- `proxy-logging`: Replace per-session-key NDJSON files with a single structured log file per startup with log level control

## Impact

- `spike/issue-45-opencode-cache/proxy/main.go` — replace `sync.Map`-based per-file NDJSON append with `slog`-based single-file structured logging
- `spike/issue-45-opencode-cache/proxy/AGENTS.md` — update log format description, env vars, and inspection commands
