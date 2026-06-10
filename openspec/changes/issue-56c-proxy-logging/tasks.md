## 1. Proxy Changes

- [ ] 1.1 Replace `sync.Map` file-keyed NDJSON logic with `log/slog` JSON handler writing to a single file opened at startup
- [ ] 1.2 Derive log file path from `LOG_FILE` env var, falling back to `<LOG_DIR>/proxy-<unix-seconds>.log`
- [ ] 1.3 Parse `LOG_LEVEL` env var (default `info`); set `slog` level accordingly
- [ ] 1.4 On each forwarded request, emit INFO record: `session_key`, `turn`, `model`, `status`, `prompt_tokens`, `completion_tokens`, `cached_tokens`
- [ ] 1.5 In DEBUG mode, add `original_last_user` field (pre-injection `messages[-1].content`) to the record
- [ ] 1.6 In DEBUG mode, add `injected_last_user` field when injection has modified `messages[-1].content` (omit when no injection)
- [ ] 1.7 In DEBUG mode, add `response_body` field (full raw SSE response)
- [ ] 1.8 On upstream failure, emit ERROR record: `session_key`, `turn`, `error`
- [ ] 1.9 On JSON parse failure, emit WARN record and forward request unchanged
- [ ] 1.10 Remove `appendLog`, `fileMutex`, and `sync.Map` from previous implementation

## 2. Rebuild and Verify

- [ ] 2.1 Rebuild proxy image
- [ ] 2.2 Run with `LOG_LEVEL=info`; confirm log file created in `LOG_DIR`, one line per request, no content fields
- [ ] 2.3 Run with `LOG_LEVEL=debug`; confirm `original_last_user` and `response_body` appear
- [ ] 2.4 Confirm `session_key` and `turn` fields present in both modes
- [ ] 2.5 Confirm token counts appear after a session with cache hits

## 3. Documentation

- [ ] 3.1 Update `proxy/AGENTS.md`: describe single log file, log levels, env vars, inspection commands (`jq` filters by session_key, level)

## 4. Close Out

- [ ] 4.1 Commit and push on `feature/issue-55-openrouter-proxy`
