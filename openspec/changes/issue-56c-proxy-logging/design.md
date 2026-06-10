## Context

The proxy currently writes per-session NDJSON files keyed by SHA-256 prefix of `messages[0].content`. Empirical testing with this spike produced:

- `d46d701d.ndjson` — agent turns (3 lines for a 3-turn session)
- `12effdd3.ndjson` — title-gen (1 line)
- `25db663d.ndjson` — variant agent (separate session, different prompt.md)

Each new agent role produces a new file. Files accumulate. There is no single view of all proxy traffic. Session boundaries within a file are invisible.

The simpler model: one file, opened at startup, all records appended. Session key and turn number are record fields. Log rotation is external (logrotate, Docker log driver, etc.).

## Goals / Non-Goals

**Goals:**
- Single append-only log file per proxy process lifetime
- Structured log levels usable without a log aggregation stack (plain JSON lines)
- DEBUG mode captures enough to reconstruct exactly what was forwarded (original vs. injected message)
- No checksums — direct content in DEBUG records is sufficient and simpler to read

**Non-Goals:**
- Built-in log rotation (delegate to logrotate or Docker log driver)
- Log shipping or aggregation
- Structured tracing or span correlation

## Decisions

**Use `log/slog` (stdlib)**: Available since Go 1.21, zero dependencies, JSON handler is built-in. The proxy has no external dependencies and this preserves that.

**File opened at startup, path derived from startup time**: `proxy-<unix-seconds>.log` gives a stable, time-ordered filename without requiring a UUID library. The path can be overridden with `LOG_FILE` for testing.

**Log levels:**
- `ERROR`: upstream request failed, file open failed, response write failed
- `WARN`: request body could not be parsed as JSON (forwarded as-is); missing expected fields
- `INFO`: one record per forwarded request — `session_key`, `turn`, `model`, `status`, `prompt_tokens`, `completion_tokens`, `cached_tokens`
- `DEBUG`: all INFO fields plus `original_last_user` (raw content of `messages[-1]` before any injection), `injected_last_user` (content after injection, or omitted if no injection), full `response_body`

**No content in INFO**: INFO records must be safe to leave on in production. Message content may contain sensitive data; token counts and metadata are sufficient for cost/performance monitoring.

**Original vs. injected logged separately**: When SF-4c injection is active, the DEBUG record includes both the pre-injection and post-injection content of `messages[-1]`. No checksum needed — the full strings are there for direct comparison.

**Session key retained as a field**: The SHA-256 prefix of `messages[0].content` remains useful as a grouping key for filtering log records by agent role (e.g., `jq 'select(.session_key=="d46d701d")'`).

## Log Record Schema

### INFO
```json
{
  "time": "2026-06-10T14:23:01Z",
  "level": "INFO",
  "msg": "request",
  "session_key": "d46d701d",
  "turn": 2,
  "model": "google/gemini-2.5-flash",
  "status": 200,
  "prompt_tokens": 1289,
  "completion_tokens": 47,
  "cached_tokens": 1164
}
```

### DEBUG (additional fields)
```json
{
  ...INFO fields...,
  "original_last_user": "What is 2+2?",
  "injected_last_user": "[state: turn=2, cached=1164]\nWhat is 2+2?",
  "response_body": "data: {...}\ndata: {...}\n..."
}
```

### ERROR
```json
{
  "time": "...",
  "level": "ERROR",
  "msg": "upstream request failed",
  "session_key": "d46d701d",
  "turn": 2,
  "error": "connection refused"
}
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_DIR` | `/logs` | Directory for log file (used if `LOG_FILE` not set) |
| `LOG_FILE` | `<LOG_DIR>/proxy-<unix>.log` | Override full log file path |
| `LOG_LEVEL` | `info` | One of `debug`, `info`, `warn`, `error` |

## Risks / Trade-offs

**Single file grows unbounded**: Logrotate or Docker log driver must be configured externally. Acceptable at spike scale; noted for production.

**DEBUG logs contain full message content**: Appropriate for development/experiment environments. Should not be used in production with real user data.

**No correlation to opencode session ID**: Still absent from the wire format. The session key (system prompt hash) groups by agent role, not conversation. Injection selection is via the `## Runtime state injection` marker in `messages[0]` — see `sf-4c-injection-strategy.md`.
