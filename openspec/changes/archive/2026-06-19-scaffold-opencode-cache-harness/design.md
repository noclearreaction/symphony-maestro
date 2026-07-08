## Context

Issue #45 is the first sub-feature (SF-1) of spike #43: *Empirically Verify Opencode Cache Behavior Across Prompt and State Boundaries*. All downstream sub-features (SF-2 through SF-8) depend on a stable, reproducible experiment baseline. Without it, token and cost measurements cannot be trusted across machines or sessions.

The harness lives in this repository temporarily so it can be co-located with planning artifacts during the spike. It has no runtime role in Symphony or any product system.

## Goals / Non-Goals

**Goals:**
- Docker image that installs opencode and all dependencies deterministically
- Minimal project fixture (single agent, ~200-token system prompt, no application code)
- A shell entry point that opens an interactive session inside the container so experiments can be run manually
- README explaining how to build the image and start an experiment session
- Reproducible: two people running `docker build` on the same inputs get the same environment

**Non-Goals:**
- Multi-turn automation or loop orchestration (belongs in SF-2+)
- CI integration for cache experiments (out of scope for this spike)
- Performance benchmarking beyond the fields listed in the issue
- Replacing or modifying any existing Symphony or Director tooling

## Decisions

### Base image: Debian slim + Node.js via official image

Use `node:20-slim` (Debian-based) as the base image rather than Alpine.

**Rationale**: opencode is a Node.js CLI. Alpine's musl libc can cause subtle compatibility issues with native addons. The slim Debian image is well-supported, small enough for experiment use, and avoids hidden debugging costs.

**Alternative considered**: Alpine (`node:20-alpine`) — rejected due to musl compatibility risk with opencode's SQLite bindings.

### opencode installation: global npm install at build time

Install opencode via `npm install -g opencode-ai` during `docker build`.

**Rationale**: Pinning a version in the Dockerfile makes the image reproducible. Installing at build time rather than runtime avoids network dependency during experiments.

**Alternative considered**: Copying a pre-built binary — rejected because it ties the image to a specific platform architecture unnecessarily.

### Project fixture: static files baked into the image

The minimal project (agent config + system prompt) is baked into the image as static files rather than mounted at runtime.

**Rationale**: Baking in the fixture ensures the baseline is identical across runs. Mounts introduce variability (file permissions, line endings, accidental edits).

**Alternative considered**: Volume-mount the fixture at `docker run` time — acceptable for iterating on prompts, but the baseline fixture should be locked.

## Risks / Trade-offs

- **opencode version drift** → Pin the exact version in the Dockerfile and document the pin. Update intentionally.
- **Docker not available in all environments** → README must clearly state Docker as a prerequisite. No fallback is provided (non-goal).
- **Fixture prompt token count varies by tokenizer** → The ~200-token target is approximate. Actual token count depends on the model's tokenizer. Document this as an approximation, not a hard constraint.

## Open Questions

- How does opencode actually accept a prompt or trigger a turn from the command line? This is a primary thing to discover in the spike.
- Where and how does opencode record token usage and cost? Is there a database, a log, an API? This is unknown until the environment is explored.
- Should the harness directory be `harness/` or `docker/`? `harness/` is more descriptive of the purpose.
