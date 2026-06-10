# opencode cache spike — Docker usage

This spike uses two Docker images. Both must be built before running experiments.

## Images

### `opencode-cache-harness` (fixture)

Built from `fixture/Dockerfile`. Contains opencode pinned to a specific version and a minimal project fixture baked in at `/app/fixture/`.

```bash
docker build -t opencode-cache-harness spike/issue-45-opencode-cache/fixture/
```

### `openrouter-proxy` (proxy)

Built from `proxy/Dockerfile`. A static Go binary that forwards `/v1/chat/completions` to OpenRouter, injecting the API key.

```bash
docker build -t openrouter-proxy spike/issue-45-opencode-cache/proxy/
```

## Running experiments

Both containers must share the `spike-net` Docker network:

```bash
docker network create spike-net  # once only

# Start proxy (requires OPENROUTER_API_KEY in .env)
docker run -d --name openrouter-proxy --network spike-net --env-file .env openrouter-proxy

# Start fixture
docker run -d --name sf-experiments --network spike-net opencode-cache-harness sleep infinity

# Run a turn
docker exec sf-experiments opencode run "Say: acknowledged"

# Query token metrics
docker exec sf-experiments opencode db \
  "SELECT data FROM message ORDER BY rowid DESC LIMIT 4" --format json

# Clean up
docker stop sf-experiments openrouter-proxy && docker rm sf-experiments openrouter-proxy
```

## API key

Store your OpenRouter API key in `.env` at the repo root (gitignored):

```
OPENROUTER_API_KEY=sk-or-v1-...
```

See `.env.example` for the template.

## Current model

`google/gemini-2.5-flash` — supports implicit prompt caching at 1024+ token prompts (no `cache_control` markers needed). Cache metrics appear in `prompt_tokens_details.cached_tokens` on the OpenRouter response, and in `cache.read` in the opencode message store.

## Findings

See `findings/` for documented experiment results.
