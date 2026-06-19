## Why

Tool versions are pinned manually in `.devcontainer/Dockerfile` ARG blocks. Without automated tracking, upgrades are discovered by chance. On every container start, a version check should run and report any stale tools — giving the developer immediate visibility with no external services, scheduling, or GitHub integration required.

## What Changes

- Add `renovate.json` configuration at the repo root to teach Renovate which registries and packages to check
- Annotate `ARG` lines in `.devcontainer/Dockerfile` with Renovate datasource hints for non-standard names
- Add `task devcontainer:check-versions` that runs Renovate via Docker in `--dry-run=lookup --platform=local` mode — queries registries, prints available updates, opens no PRs
- Wire `devcontainer:check-versions` into the `post-start` script so it runs on every container start

## Capabilities

### New Capabilities

- `renovate-config`: Renovate configuration and local version-check task for the devcontainer

### Modified Capabilities

<!-- none -->

## Impact

- Adds `renovate.json` to repo root
- Adds `devcontainer:check-versions` task to `.devcontainer/Taskfile.yaml`
- Updates `.devcontainer/post-start` to call the new task
- Annotates ARG lines in `.devcontainer/Dockerfile` (comments only, no functional change)
- Requires Docker socket access inside the container (already confirmed working via DooD)
