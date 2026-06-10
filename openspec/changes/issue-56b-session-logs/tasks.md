## 1. Proxy Changes

- [x] 1.1 Add `crypto/sha256` and `sync` imports to `main.go`
- [x] 1.2 Add `sessionKey(messages)` function: hash first 512 bytes of `messages[0].content`, return first 8 hex chars; return `"unknown"` if messages is empty
- [x] 1.3 Add `turnNumber(messages)` function: return `(len(messages) + 1) / 2`
- [x] 1.4 Add per-file mutex map (`sync.Map` keyed by session key) to guard concurrent appends
- [x] 1.5 Replace `logEntry` struct and `writeLog` function: new `appendLog` function appends one JSON line to `<LOG_DIR>/<session-key>.ndjson`
- [x] 1.6 Update `handleChatCompletion`: parse request body as JSON to extract `messages`, derive session key and turn number, call `appendLog` on success and error

## 2. Rebuild and Verify

- [x] 2.1 Rebuild proxy image: `docker build -t openrouter-proxy spike/issue-45-opencode-cache/proxy/`
- [x] 2.2 Run 3-turn session; confirm `/tmp/proxy-logs/` contains exactly 2 files: one for agent turns, one for title-gen
- [x] 2.3 Confirm agent session file has 3 lines (one per turn), each parseable as JSON
- [x] 2.4 Confirm turn numbers increment correctly across lines

## 3. Documentation

- [x] 3.1 Update `proxy/AGENTS.md`: describe NDJSON format, session key derivation, and file naming

## 4. Close Out

- [ ] 4.1 Commit and push on `feature/issue-55-openrouter-proxy`
