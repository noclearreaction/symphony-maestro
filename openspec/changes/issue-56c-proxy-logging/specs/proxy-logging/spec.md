## Requirement: Single log file per proxy startup

The proxy SHALL open one append-only log file when it starts. All subsequent log records SHALL be written to that file for the lifetime of the process.

The log file path SHALL be determined as follows:
1. If `LOG_FILE` env var is set, use that path
2. Otherwise, use `<LOG_DIR>/proxy-<unix-timestamp-seconds>.log`

#### Scenario: All requests land in one file
- **WHEN** the proxy handles multiple requests from different agent roles
- **THEN** all log records are in a single file, filterable by `session_key`

#### Scenario: Each proxy startup produces a new file
- **WHEN** the proxy container is restarted
- **THEN** a new file is created with the new startup timestamp; the old file is not modified

---

## Requirement: Structured log levels

The proxy SHALL support four log levels: `DEBUG`, `INFO`, `WARN`, `ERROR`. The active level SHALL be controlled by the `LOG_LEVEL` env var (default: `info`). Each log record SHALL be a valid JSON object on a single line (NDJSON).

#### Scenario: INFO level omits message content
- **GIVEN** `LOG_LEVEL=info`
- **WHEN** a request is forwarded
- **THEN** the log record contains `session_key`, `turn`, `model`, `status`, `prompt_tokens`, `completion_tokens`, `cached_tokens` — and no `content` fields

#### Scenario: DEBUG level includes message content
- **GIVEN** `LOG_LEVEL=debug`
- **WHEN** a request is forwarded
- **THEN** the log record includes all INFO fields plus `original_last_user` (pre-injection content of `messages[-1]`) and `response_body`

#### Scenario: DEBUG level logs injection diff
- **GIVEN** `LOG_LEVEL=debug` and injection is active
- **WHEN** a request is forwarded with injected content in `messages[-1]`
- **THEN** the log record includes both `original_last_user` and `injected_last_user` as separate fields

#### Scenario: ERROR level logs upstream failure
- **WHEN** the upstream request to OpenRouter fails
- **THEN** a log record at `ERROR` level is written with `session_key`, `turn`, and `error` fields

---

## Requirement: Token counts extracted from streaming response

The proxy SHALL parse the streaming SSE response to extract token usage from the final `data:` chunk. The extracted values SHALL be included in the INFO log record as `prompt_tokens`, `completion_tokens`, and `cached_tokens` (0 if absent).

#### Scenario: Cache hit reflected in log record
- **WHEN** OpenRouter returns `prompt_tokens_details.cached_tokens > 0`
- **THEN** the INFO log record has `cached_tokens` set to that value

---

## Requirement: Unparseable requests logged at WARN

- **WHEN** the request body cannot be parsed as JSON, or `messages` field is absent
- **THEN** the proxy forwards the request unchanged and logs a `WARN` record with the parse error; no `session_key` or `turn` fields are set
