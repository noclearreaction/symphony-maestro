## Why

The SF-4b implementation writes one log file per request. A 3-turn session produces 4+ files — one title-generation request and one agent request per turn — making the log unreadable as a conversation. The useful unit of inspection is a session (all turns in sequence), not an individual HTTP request.

## What Changes

- Log files change from per-request to per-session: one NDJSON file per session, one line appended per request
- Session identity derived from `messages` array length and content hash of `messages[0]` — groups requests that share the same stable system prompt as one session
- Proxy startup timestamp used as the session key when no stable grouping signal is available (e.g. title-gen requests)
- File named `<session-key>.ndjson` rather than `<timestamp>.json`
- Each appended line is a self-contained JSON object: `{timestamp, turn, request, response}`

> **Superseded by `issue-56c-proxy-logging`**: Empirical testing revealed that multiple agent roles produce multiple NDJSON files, and session boundaries are not distinguishable within a file (two sessions with the same system prompt share a file). A single log file per proxy startup with structured log levels is a simpler and more operationally sound approach. The per-session-key grouping logic is replaced by a single append-only log with session key and turn number as record fields.

## Capabilities

### Modified Capabilities

- `proxy-inspection`: Change log granularity from per-request files to per-session NDJSON files; update file naming and append behavior

## Impact

- `spike/issue-45-opencode-cache/proxy/main.go` — replace `writeLog` single-file approach with session-keyed append
- `spike/issue-45-opencode-cache/proxy/AGENTS.md` — update log format description
- `spike/issue-45-opencode-cache/openspec/changes/issue-56-proxy-inspection/specs/proxy-inspection/spec.md` — corrected spec (MODIFIED)
