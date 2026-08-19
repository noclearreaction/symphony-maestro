# Rubato

Rubato is an OpenAI-compatible HTTP proxy for chat completions. It sits between a coding agent and the upstream model API, enabling request inspection and mutation without modifying the agent's configuration.

The upstream is a single environment variable, so serving two upstreams means running two instances rather than adding a second route:

| Instance | Upstream | Client | How it runs |
|---|---|---|---|
| `mistral-rubato` :8081 | `api.mistral.ai` | Vibe | `mistral-rubato` compose service |
| :8080 | `openrouter.ai/api` | opencode | **not currently running** — see below |

Only the Mistral instance is containerised. The `:8080` references throughout this document and in `opencode.json` are therefore dangling: nothing listens there, and opencode currently fails to reach a proxy.

This is deliberate and is not waiting on a decision. Opencode is slated for replacement by Vibe under Ensemble, so its proxy path is left unserved rather than maintained. If opencode is needed again before that lands, the fix is a second containerised instance following the `mistral-rubato` pattern; otherwise the service and these `:8080` references retire with it.

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

## How Vibe routes through Rubato

Vibe routes by replacing the whole `mistral` provider entry in `~/.vibe/config.toml`. A provider override is not merged, so every field is restated:

```toml
[[providers]]
name = "mistral"
api_base = "http://mistral-rubato:8081/v1"
api_key_env_var = "MISTRAL_API_KEY"
backend = "mistral"
```

The host is the compose **service name**, resolved by Docker's embedded DNS on the project network. The port is not published to the host, so this address only works from a container attached to that network.

This configuration currently lives in `~/.vibe/` on the shared `mistral-vibe` volume rather than in the repository, so a fresh clone gets the service but nothing routing through it.

## Configuration

Rubato reads configuration from environment variables:

| Variable | Default | Description |
|---|---|---|
| `RUBATO_UPSTREAM_URL` | `https://openrouter.ai/api` | Upstream API base (path is forwarded from the incoming request) |
| `RUBATO_LISTEN_ADDR` | `:8080` | Address to listen on |
| `OPENROUTER_API_KEY` | _(none)_ | Bearer token forwarded to upstream |

`OPENROUTER_API_KEY` is only a fallback. Rubato forwards the client's own `Authorization` header when there is one and consults this variable only when there is not, which is why the `mistral-rubato` service configures no key: Vibe's key passes straight through. The variable name is upstream-specific and predates the second instance.

## Prerequisites

Create `/workspace/.env` with your OpenRouter key:

```
OPENROUTER_API_KEY=sk-or-...
```

Nothing loads this automatically any more — the `dotenv` integration went with the local process lifecycle. Source it yourself when running Rubato or the smoke test by hand:

```bash
set -a && source /workspace/.env && set +a
```

## Lifecycle

Rubato runs as the `mistral-rubato` compose service, started by `runServices` in `.devcontainer/devcontainer.json` and supervised by compose via `restart: unless-stopped`. It is deliberately **not** started by `.devcontainer/bin/post-start`: the previous backgrounded process did not survive post-start's teardown, reporting success while leaving nothing listening.

```bash
task rubato:status   # state and healthcheck result
task rubato:logs     # last 100 lines; -- --follow to tail
task rubato:start    # idempotent
task rubato:stop     # stops the container, keeps it
task rubato:restart  # reuses the container — see below
task rubato:doctor   # reachability check, also run by devcontainer:doctor
```

None of these tasks ever recreates the container, so none of them picks up an edited `docker-compose.yml` or a rebuilt image — rebuild or reopen the devcontainer for that, then rebuild the image on the host with:

```bash
cd .devcontainer && docker buildx bake --allow=fs.read=.. -f versions.hcl -f docker-bake.hcl
```

That restraint is deliberate. VS Code creates these containers from the base compose file plus two generated override files, so a plain `docker compose up -d` against the base file alone sees different configuration and recreates the container. The devcontainer owns container creation; the tasks only start, stop and inspect.

## rplugin — run a plugin from the terminal

`rplugin` is a standalone binary that executes a single named rubato plugin and prints its output to stdout. Use it to verify plugin behavior, inspect ambient state, and test args without starting the proxy.

```bash
go run ./cmd/rplugin <plugin-name> [--working-dir <path>] [--args <json>]
```

### Examples

```bash
# git_status in the current workspace
go run ./cmd/rplugin git_status --working-dir /workspace

# go_test with a custom timeout
go run ./cmd/rplugin go_test --working-dir /workspace --args '{"timeout_seconds":30}'
```

Errors (unknown plugin, execution failure, invalid JSON) are written to stderr and the process exits with code 1.

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

The smoke test uses port `18080` so it does not conflict with a Rubato started
by hand on `:8080`. The containerised instance is unaffected either way: its
`:8081` is inside the compose network and never bound on the devcontainer. To
run the same scenario by hand, use two terminals.

**Terminal 1 — start Rubato**

```bash
cd /workspace
RUBATO_LISTEN_ADDR=:18080 \
RUBATO_UPSTREAM_URL=https://openrouter.ai/api \
go run ./cmd/rubato/
```

**Terminal 2 — run opencode**

The fixture `testdata/smoke/opencode.json` overrides the `openrouter` provider
`baseURL` to `http://127.0.0.1:18080/v1`. The smoke agent declares the
`git_status` anchor.

```bash
set -a && source /workspace/.env && set +a
cd /tmp
OPENCODE_CONFIG=/workspace/cmd/rubato/testdata/smoke/opencode.json \
OPENCODE_CONFIG_DIR=/workspace/cmd/rubato/testdata/smoke \
opencode run --pure \
  --model openrouter/openai/gpt-4o-mini \
  --agent smoke \
  "What branch am I on right now?"
```

Switch back to Terminal 1 to see the injection in the log.
