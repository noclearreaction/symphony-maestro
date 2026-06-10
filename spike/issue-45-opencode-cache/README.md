# opencode cache harness

A minimal, reproducible Docker environment for exploring opencode cache behavior as part of spike [#43](https://github.com/noclearreaction/symphony-director/issues/43).

This harness is the baseline for SF-1 ([#45](https://github.com/noclearreaction/symphony-director/issues/45)). It provides a clean, isolated container with opencode installed and a minimal project fixture ready to use.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed and running

## Pinned versions

| Software | Version |
|---|---|
| Node.js (base image) | `node:20-slim` |
| opencode | `1.16.2` |

To update the opencode version, edit the `RUN npm install -g opencode-ai@...` line in `Dockerfile` and rebuild.

## Fixture layout

The minimal project fixture is baked into the image at `/app/fixture/`:

```
/app/fixture/
├── Dockerfile        # image definition for opencode-cache-harness
├── opencode.json     # opencode config: model, provider, instructions
├── INSTRUCTIONS.md   # stable system instructions (~1024+ tokens, cache target)
└── prompt.md         # agent identity prompt
```

The fixture defines a single agent (`build`) with a stable, low-variability system prompt sized to exceed the Gemini 2.5 Flash implicit cache threshold (1024 tokens). It has no application code and all tools are denied.

## Build the images

```bash
# Fixture image
docker build -t opencode-cache-harness spike/issue-45-opencode-cache/fixture/

# Proxy image
docker build -t openrouter-proxy spike/issue-45-opencode-cache/proxy/
```

Run from the repository root.

## Start an experiment session

```bash
docker run --rm -it opencode-cache-harness
```

This drops you into a bash shell inside the container with:
- `opencode` available on the PATH
- Working directory set to `/app/fixture/` (the minimal project)

From there you can invoke opencode directly and explore its behavior:

```bash
# Check opencode is available
opencode --version

# Explore available commands
opencode --help

# Check what database tooling is available
opencode db --help

# View session stats
opencode stats
```

## Extending the fixture

To iterate on the system prompt without rebuilding:

```bash
docker run --rm -it \
  -v "$(pwd)/spike/issue-45-opencode-cache/fixture:/app/fixture" \
  opencode-cache-harness
```

This volume-mounts your local fixture over the baked-in one, so edits are reflected immediately without a rebuild. Use this for prompt iteration; the baked-in fixture is the reproducible baseline.

## Multi-turn experiments (cache testing)

`docker run --rm` destroys the DB after each run. To observe cache hits across turns, keep the container alive and reuse the session:

```bash
# Start a persistent container
docker run -d --name cache-exp opencode-cache-harness sleep infinity

# Turn 1 — establishes the session and primes the cache
docker exec cache-exp opencode run "What is 2+2? Reply with only the number."

# Capture the session ID
SESSION=$(docker exec cache-exp opencode db "SELECT id FROM session ORDER BY time_created DESC LIMIT 1" --format json | grep -o 'ses_[^"]*')

# Turn 2 — continues the session; system prompt served from cache
docker exec cache-exp opencode run --session "$SESSION" "What is 3+3? Reply with only the number."

# Read cache metrics
docker exec cache-exp opencode db \
  "SELECT tokens_input, tokens_cache_read, tokens_cache_write FROM session WHERE id=\"$SESSION\"" \
  --format json

# Clean up
docker stop cache-exp && docker rm cache-exp
```

Expected result: `tokens_cache_read` is ~0 after turn 1, then ~512 after turn 2.

## Proxy-routed experiments (SF-4+)

For experiments that require intercepting or mutating the request (SF-4 and later), run the `openrouter-proxy` container alongside the fixture container on a shared Docker network.

```bash
# Create shared network (once)
docker network create spike-net

# Start proxy (OPENROUTER_API_KEY required — store in .env at repo root)
docker run -d --name openrouter-proxy \
  --network spike-net \
  --env-file .env \
  openrouter-proxy

# Start fixture container on same network
docker run -d --name sf-experiments \
  --network spike-net \
  opencode-cache-harness sleep infinity

# Run a turn (opencode will route through the proxy)
docker exec sf-experiments opencode run "Say: acknowledged"

# Read cache metrics as usual
SESSION=$(docker exec sf-experiments opencode db "SELECT id FROM session ORDER BY time_created DESC LIMIT 1" --format json | grep -o 'ses_[^"]*')
docker exec sf-experiments opencode db \
  "SELECT tokens_input, tokens_cache_read FROM session WHERE id=\"$SESSION\"" \
  --format json

# Clean up
docker stop sf-experiments openrouter-proxy && docker rm sf-experiments openrouter-proxy
```

## Notes

- The opencode version is pinned in `fixture/Dockerfile`. Do not change it mid-spike without documenting the change as a variable.
- The current model is `google/gemini-2.5-flash` routed via `openrouter-proxy`. It uses implicit caching — no `cache_control` markers required. Cache activates at 1024+ input tokens.
- `INSTRUCTIONS.md` is the stable cache target (~800 words, ~1040 tokens). `prompt.md` adds ~160 more. Combined they comfortably exceed the 1024 token threshold.
- Cache metrics live in `prompt_tokens_details.cached_tokens` on the raw OpenRouter response, and in `cache.read` within the `data` JSON column of the opencode `message` table.
- SF-3 baseline was established on `opencode/deepseek-v4-flash-free` (opencode's own infrastructure). See `findings/sf-2-observability.md`.
