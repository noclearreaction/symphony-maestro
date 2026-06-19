## Why

The devcontainer previously installed openspec, opencode, and renovate via a pnpm
project with `node-linker=hoisted` — a workaround that broke native module
resolution across multi-stage Docker builds and left the install unreprodible,
unsecured, and unmaintainable. All three tools must be reinstalled using correct
patterns that work reliably in Docker and in CI.

## What Changes

- Add a `node-builder` stage to the Dockerfile with a `NODE_VERSION` ARG (managed by Renovate) — pnpm version is owned by the `packageManager` field in `package.json`, read by corepack at task runtime; no runtime Docker Hub pulls needed
- Add committed `.devcontainer/node/package.json` (with `packageManager` pinned), `pnpm-lock.yaml`, and `pnpm-workspace.yaml` — versions and approvals are now reviewed artifacts
- Add a named Docker volume (`node-modules`) declared in `docker-compose.yml`, mounted at `/opt/node_modules` in the devcontainer
- `pnpm-workspace.yaml` sets `modulesDir: "/opt/node_modules"` so `pnpm install` writes packages directly into the volume — no `pnpm deploy` needed
- At container start, the `symphony-studio-node-builder` container runs `pnpm install --frozen-lockfile` via DooD with `.devcontainer/node/` bind-mounted RW and the volume mounted at `/opt/node_modules`
- The Dockerfile bakes `ENV PATH=/opt/node_modules/.bin:$PATH` into the devcontainer image; no symlink or wrapper scripts needed
- Add `task node:build` — runs the builder container to populate or update the volume; called at container start and after any package change; no image rebuild needed
- Add `task node:package:prune` — no-op; pnpm store is ephemeral per-build
- Add `minimumReleaseAge: 10080` (7 days), `blockExoticSubdeps: true`, `trustPolicy: no-downgrade` to `pnpm-workspace.yaml`
- Add `minimumReleaseAge: "7 days"` to `renovate.json` for all version updates
- Add `task node:package:add`, `task node:package:rm`, `task node:package:update`, `task node:package:list`, `task node:package:audit`, `task node:package:prune` — sandboxed npm package management; pnpm never required in the devcontainer
- Add `task node:trust:add`, `task node:trust:rm`, `task node:trust:list`, `task node:trust:verify` — explicit build script approval management, independent of package installation
- Restore `task devcontainer:check-versions` using renovate from the volume
- Restore openspec, opencode, and renovate checks in `task devcontainer:doctor`

## Capabilities

### New Capabilities

- `devcontainer-node-install`: Reproducible, supply-chain-hardened npm package installation in Docker using committed pnpm project files, frozen lockfile, and `modulesDir` volume targeting

### Modified Capabilities

- `devcontainer-environment`: node-runtime stage removed; volume-based node_modules with `/opt/node/bin` on PATH; openspec, opencode, and renovate available via symlink into volume

## Impact

- `.devcontainer/Dockerfile`: node-runtime stage removed; `ENV PATH=/opt/node_modules/.bin:$PATH` added to base stage; `node-builder` stage has only `corepack enable` (no `corepack prepare`; pnpm version owned by `packageManager` in `package.json`)
- `devcontainer.json`: named volumes `${localWorkspaceFolderBasename}-node-modules` (mounted at `/opt/node/node_modules`) and `${localWorkspaceFolderBasename}-pnpm-store` (builder-only) added
- `.devcontainer/node/package.json`: new committed file
- `.devcontainer/node/pnpm-lock.yaml`: new committed file
- `.devcontainer/node/pnpm-workspace.yaml`: new committed file with allowBuilds + hardening
- `.devcontainer/Taskfile.yaml`: `node:build` task added; check-versions and doctor restored
- `Taskfile.yaml` (root or new `node` namespace): `node:package:*` and `node:trust:*` tasks added
- `renovate.json`: minimumReleaseAge added
