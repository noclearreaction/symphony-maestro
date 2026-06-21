# Handover — 2026-06-12

## Strategic direction: rubato → devcontainer

The rubato Docker compose stack is being replaced. The core problem it was solving (running opencode with proxy interception) is better addressed by running everything natively inside a devcontainer.

**Key facts:**
- Proxy runs as a Go process inside the devcontainer (`go build`, then `./proxy`)
- opencode runs natively — tools, auth, credentials all present in devcontainer
- opencode config routes to `http://localhost:8080/v1` — no injection magic, just config
- Proxy logs to a host-bound path via devcontainer volume mount
- Per-session log separation is lost (proxy is a rolling append log) — accepted tradeoff
- Go build cache in a named Docker volume (`go-build-cache`) shared across devcontainers
- Proxy binary is gitignored, built on demand

**What survives from rubato:** `proxy/` source, Taskfile (simplified)

**What is deleted:** `opencode/` image, `runner/` image, `docker-compose.yml`, `bin/rubato`

**`.devcontainer/` is committed** — it's the authoritative environment definition, not per-developer.

## Injection spike — confirmed working (2026-06-12)

Static prefix injected into `messages[-1]` when `role == "user"`. Model responded with words not digits. Full details in issue #61 comment.

**Confirmed:** proxy can mutate outbound requests and it observably affects model behaviour.

**Finding:** SSE streaming prevents response-side log parsing — responses log as raw chunks. Request side is clean.

**Design (issue #61):** Gate injection on `messages[0]` containing `proxy:inject` marker. Agents opt in via system prompt. `messages[0]` is never modified (cache stays warm). Dynamic content (git status, test output) populated from commands run by the proxy process directly in the devcontainer environment.

## Open issues of note

- #61 — marker-based runtime injection (next implementation target)
- #62 — structured logging
- #64 — subcommand validation bug (rubato — may be moot after compose removal)
- #60 series — rubato hygiene (60e, 60f largely moot after devcontainer pivot)

## Basic usage (devcontainer)

- Rubato autostarts on devcontainer start from `.devcontainer/bin/post-start`.
- Default route is enabled via devcontainer `remoteEnv`: `OPENAI_BASE_URL=http://127.0.0.1:8080/v1`.
- Default model in `opencode.json` is `opencode/deepseek-v4-flash-free`.

### Setup

1. Put your OpenRouter key in `/workspace/.env`:

	`OPENROUTER_API_KEY=...`

2. Rebuild/reopen the devcontainer (or rerun `.devcontainer/bin/post-start`).

### Verify

1. Check Rubato is listening:

	`lsof -iTCP:8080 -sTCP:LISTEN -n -P`

2. Run a smoke request through opencode:

	`opencode run --format json "Reply with exactly OK"`

3. Confirm reply text in output and inspect logs if needed:

	`tail -n 50 /tmp/rubato.log`
