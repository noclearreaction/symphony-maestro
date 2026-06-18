## Why

Symphony-maestro has no devcontainer definition, meaning every developer and agent session runs against the host environment — different tool versions, missing tools, and no reproducibility. This change establishes a standardised, reproducible development environment using a multi-stage Dockerfile that can evolve into a shared base image across project repositories.

## What Changes

- Add `.devcontainer/Dockerfile` with a multi-stage build assembling Go, Deno, Task, and Node (constrained) runtimes
- Add `.devcontainer/devcontainer.json` wiring VS Code to the Dockerfile-built image with appropriate features and mounts
- Install binaries via `update-alternatives` for all versioned runtimes (Go with gofmt slaved, Deno) following Linux FHS — versioned installs under `/opt/`, alternatives managing `/usr/local/bin/`
- Go installed at `/opt/go<version>/` (not `/usr/local/go`) with `GOROOT` pointing at the versioned path
- Node/pnpm isolated: not on dev `PATH`, accessed only through `openspec` and `opencode` wrappers in `/usr/local/bin/`
- `npm`, `npx`, `pnpm` intentionally absent from `PATH` — no wrapper magic, no workarounds, topology enforces the constraint
- Add a `ci` stage between `base` and `final` — CI runs against `ci`, which excludes user-facing tooling (SSH profiles, vscode user dirs, git hook wiring)
- Task installed from standalone go-task GitHub release binary, not npm
- No Docker Compose; DooD (Docker-outside-of-Docker) via devcontainer feature for future sidecar needs
- No `docker buildx bake` — plain `docker build` with multi-stage is sufficient; bake's host dependency is an antipattern here

## Capabilities

### New Capabilities

- `devcontainer-environment`: Reproducible development container definition — Dockerfile, devcontainer.json, stage structure, tool installation conventions

### Modified Capabilities

<!-- None: no existing specs change requirements -->

## Impact

- `.devcontainer/` directory created (new)
- Developers opening the repo in VS Code will be prompted to reopen in container
- Host environment no longer required to have Go, Deno, Task, or Node installed
- `bin/commit-lint.ts` and `bin/provision-labels.ts` shebangs will need updating — they currently hardcode the host Deno path; this is tracked as a task within this change
- CI pipelines (future) can target the `ci` stage directly
