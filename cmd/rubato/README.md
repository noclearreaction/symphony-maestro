# Rubato

Rubato is an OpenAI-compatible HTTP proxy for chat completions. It sits between opencode and the upstream model API, enabling request inspection and mutation without modifying opencode's configuration.

## How opencode routes through Rubato

Routing is configured per-provider in `opencode.json` using the `provider.<id>.options.baseURL` override documented at https://opencode.ai/docs/providers/#config.

Example — redirect the `openrouter` provider to a local Rubato instance:

```json
{
  "provider": {
    "openrouter": {
      "options": {
        "baseURL": "http://127.0.0.1:8080/v1"
      }
    }
  }
}
```

The devcontainer also sets `OPENAI_BASE_URL=http://127.0.0.1:8080/v1` in `.devcontainer/devcontainer.json`, which overrides the base URL for the built-in **OpenAI provider only**. Models from other providers (e.g. `openrouter/*`) require the `opencode.json` override above.

## Configuration

Rubato reads configuration from environment variables:

| Variable | Default | Description |
|---|---|---|
| `RUBATO_UPSTREAM_URL` | `https://openrouter.ai/api` | Upstream API base (path is forwarded from the incoming request) |
| `RUBATO_LISTEN_ADDR` | `:8080` | Address to listen on |
| `OPENROUTER_API_KEY` | _(none)_ | Bearer token forwarded to upstream |

## Prerequisites

Create `/workspace/.env` with your OpenRouter key:

```
OPENROUTER_API_KEY=sk-or-...
```

The Taskfile `dotenv` integration loads this automatically when starting Rubato.

## Lifecycle

Rubato autostarts on devcontainer start via `.devcontainer/bin/post-start`.

Manual lifecycle commands:

```bash
task rubato:start    # build and start in background
task rubato:stop     # stop listener
task rubato:restart  # stop then start
task rubato:status   # show listener/PID
task rubato:logs     # tail /tmp/rubato.log
```

## Smoke test

`smoke_test.go` verifies the full round-trip: opencode → Rubato → upstream → opencode. It is tagged `//go:build smoke` and never runs under `go test ./...`.

The test:
- Re-invokes the compiled test binary as a Rubato subprocess on fixed port `18080` (`TestMain` detects `RUBATO_TEST_SUBPROCESS=1` and calls `main()` directly — no separate build step). The test fails immediately if port `18080` is already in use.
- Runs opencode from `os.TempDir()` with `OPENCODE_CONFIG` pointing to `testdata/smoke/opencode.json`, which overrides the `openrouter` provider `baseURL` to `http://127.0.0.1:18080/v1`
- Embeds a unique UUID probe token in the prompt and asserts it appears in Rubato's logged request body — proving the specific request transited Rubato, not just that Rubato started

### Run

```bash
set -a && source /workspace/.env && set +a
go test -tags smoke -v -run TestSmokeRoundTrip ./cmd/rubato/
```

The test skips (not fails) if `OPENROUTER_API_KEY` is unset or `opencode` is not in `PATH`.

### Manual verification

To verify Rubato routing without running the test:

1. Confirm Rubato is listening:

   ```bash
   lsof -iTCP:8080 -sTCP:LISTEN -n -P
   ```

2. Run a request through the devcontainer-routed instance:

   ```bash
   cd /tmp && \
   OPENCODE_CONFIG=/workspace/cmd/rubato/testdata/smoke/opencode.json \
   OPENCODE_CONFIG_DIR=/workspace/cmd/rubato/testdata/smoke \
   opencode run --pure --model openrouter/openai/gpt-4o-mini --agent smoke --format json --title smoke "hello"
   ```

3. Inspect the Rubato log:

   ```bash
   tail -n 50 /tmp/rubato.log
   ```
