# opencode cache harness

<<<<<<< HEAD
A minimal, reproducible Docker environment for exploring opencode cache behavior as part of spike [#43](https://github.com/noclearreaction/symphony-director/issues/43).

This harness is the baseline for SF-1 ([#45](https://github.com/noclearreaction/symphony-director/issues/45)). It provides a clean, isolated container with opencode installed and a minimal project fixture ready to use.
=======
A minimal, reproducible Docker environment for exploring opencode cache behavior as part of spike [#43](https://github.com/noclearreaction/symphony-maestro/issues/43).

This harness is the baseline for SF-1 ([#45](https://github.com/noclearreaction/symphony-maestro/issues/45)). It provides a clean, isolated container with opencode installed and a minimal project fixture ready to use.
>>>>>>> actual-troubleshooting

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
├── opencode.json   # opencode config: loads AGENTS.md, sets default agent
└── AGENTS.md       # ~160-token system prompt for the experiment agent
```

The fixture defines a single agent (`experiment`) with a short, low-variability system prompt. It has no application code and no tools configured.

## Build the image

```bash
docker build -t opencode-cache-harness harness/
```

Run from the repository root. The build installs opencode globally inside the image.

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
  -v "$(pwd)/harness/fixture:/app/fixture" \
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

## Notes

- The opencode version is pinned in the Dockerfile. Do not change it mid-spike without documenting the change as a variable.
- Cache behavior confirmed on `opencode/deepseek-v4-flash-free` (free model, no API key required). See `findings/sf-2-observability.md`.
- The `experiment` agent system prompt is approximately 160 tokens. This is an approximation — actual token count depends on the model's tokenizer.
- How to trigger turns, where token usage is recorded, and how to read cache metrics are open questions this spike is designed to answer.
