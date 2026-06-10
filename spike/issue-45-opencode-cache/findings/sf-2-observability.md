# SF-2 Observability Findings

opencode 1.16.2 | free model: `opencode/deepseek-v4-flash-free` | harness: SF-1

---

## How to measure

```bash
# Run a turn
opencode run "<message>"

# Read metrics (use inside container, or via docker exec)
opencode db "SELECT tokens_input, tokens_output, tokens_reasoning, tokens_cache_read, tokens_cache_write, cost FROM session ORDER BY time_created DESC LIMIT 1" --format json

# DB path
/root/.local/share/opencode/opencode.db

# Full session export (messages + tokens)
opencode export <sessionID>
```

Token counts are session-level cumulative totals. **They are not in the debug log** — DB only.

---

## Debug log

The debug log captures what is sent to the model and how the session is structured. Useful for verifying what content is in context and for diagnosing why cache hits do or don't occur.

```bash
# Capture debug log from a run
opencode run --print-logs --log-level DEBUG "<message>" 2>debug.log

# Log format: one entry per line
# LEVEL  TIMESTAMP +ELAPSEDms service=<name> [key=value ...] <message>
```

Key log entries for cache/context experiments:

| service | What it shows |
|---|---|
| `session` | Session created with initial token snapshot (all zeros — pre-turn) |
| `llm` | LLM call dispatched: `providerID`, `modelID`, `agent`, `mode=primary` |
| `session.prompt` | Loop steps: `step=0` (start), `step=1` (after response), `exiting loop` |
| `session.tools` + `tool.registry` | Which tools were registered (and thus whose schemas were included) |

**What the debug log does NOT contain**: post-turn token counts, cache hit/miss signals, or response content. Those come from the DB and `opencode export` respectively.

To confirm what is in context for a given turn (system prompt, tool schemas, prior messages), check the `service=llm` entry and cross-reference the session's token count from the DB.

---

## session table columns (relevant to experiments)

`id`, `slug`, `agent`, `model` (JSON), `cost` (REAL), `tokens_input`, `tokens_output`, `tokens_reasoning`, `tokens_cache_read`, `tokens_cache_write`, `time_created` (Unix ms)

---

## Fixture configuration

`instructions: ["AGENTS.md"]` **appends** to the built-in agent prompt — it does not replace it.  
`agent.build.prompt: "{file:./AGENTS.md}"` **replaces** the system prompt.  
Enabled tools add ~1650 tokens per turn from their schema definitions.

Token baseline with corrected fixture (prompt override + all tools denied): **515 input tokens**.

---

## Cache behavior confirmed (free model)

Prompt caching works with `opencode/deepseek-v4-flash-free` within a session. Tested by running two turns on the same session using `opencode run --session <id>`:

| Turn | `tokens_input` | `tokens_cache_read` | `tokens_cache_write` |
|---|---|---|---|
| 1 | 515 | 0 | 0 |
| 2 | 536 | 512 | 0 |

Turn 2: 512 of 536 input tokens served from cache (99%). The system prompt is cached after the first turn and reused on subsequent turns in the same session.

**Implication for experiments**: No paid model or API key is required to observe cache hits. The free model caches within a session.

**Implication for harness design**: `docker run --rm` destroys the DB after every run. To test cache behavior across turns, the DB must persist — either via a named volume, a long-running container, or `docker exec` into a running container.

```bash
# Multi-turn pattern (DB persists across turns)
docker run -d --name cache-exp opencode-cache-harness sleep infinity
docker exec cache-exp opencode run "<message 1>"
SESSION=$(docker exec cache-exp opencode db "SELECT id FROM session ORDER BY time_created DESC LIMIT 1" --format json | grep -o 'ses_[^"]*')
docker exec cache-exp opencode run --session "$SESSION" "<message 2>"
docker exec cache-exp opencode db "SELECT tokens_input, tokens_cache_read FROM session WHERE id=\"$SESSION\"" --format json
docker stop cache-exp && docker rm cache-exp
```

---

## Gaps vs #43 assumptions

- `tokens_input`, `tokens_cache_read`, `cost` — all correct, exact column names in `session` table
- `opencode db "<SQL>" --format json` — works exactly as assumed
- **Token counts not in debug log** — DB only (assumption incorrect)
- **`opencode run` stdout is empty** — response appears in TUI/stderr only
- **`tokens_cache_read` is 0 for free models** — prompt caching requires Anthropic (or cache-capable) model

