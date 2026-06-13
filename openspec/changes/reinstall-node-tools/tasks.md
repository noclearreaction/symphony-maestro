## 1. Committed pnpm project files

- [x] 1.1 Create `.devcontainer/node/package.json` with empty `dependencies` — no `packageManager` field; pnpm version is owned by `PNPM_VERSION` in the Dockerfile
- [x] 1.2 Create `.devcontainer/node/pnpm-workspace.yaml` with supply chain hardening settings (`minimumReleaseAge: 10080`, `blockExoticSubdeps: true`, `trustPolicy: no-downgrade`) — trust entries are managed by `task node:trust:add/rm`

## 2. Named Docker volumes and docker-compose.yml

- [ ] 2.1 Add named volume `node-modules` to `docker-compose.yml` and mount it at `/opt/node/node_modules` on the `symphony-studio` service
- [ ] 2.2 Add named volume `pnpm-store` to `docker-compose.yml` (used only by builder containers via DooD, not mounted on the devcontainer)
- [ ] 2.3 Add `node-builder` service to `docker-compose.yml` with `profiles: [node-tools]` so it builds with `docker-compose build` but is never started automatically; tag the image `symphony-studio-node-builder`
- [ ] 2.4 Verify the node-modules volume is shared between the devcontainer and containers spawned via DooD (key MVP assumption)

## 3. Dockerfile: add node-builder stage, remove node-runtime stage

- [x] 3.1 Add `ARG NODE_VERSION` and `ARG PNPM_VERSION` back to the Dockerfile ARG block with Renovate annotations
- [x] 3.2 Add a `node-builder` stage: `FROM node:${NODE_VERSION}-bookworm-slim AS node-builder` with `RUN corepack enable && corepack prepare pnpm@${PNPM_VERSION} --activate`
- [x] 3.3 Remove the old node-runtime stage from the Dockerfile entirely
- [x] 3.4 Remove all old npm package ARGs (`OPENSPEC_VERSION`, `OPENCODE_VERSION`, `RENOVATE_VERSION`, `RE2_VERSION`) — versions now owned by `package.json`
- [x] 3.5 Add `RUN mkdir -p /opt/node && ln -s /opt/node/node_modules/.bin /opt/node/bin` and `ENV PATH=/opt/node/bin:$PATH` to the base stage

## 4. node:build task

- [ ] 4.1 Add `task node:build` that runs the `symphony-studio-node-builder` image via DooD with `.devcontainer/node/` bind-mounted as source, `node-modules` volume mounted at `/dest/node_modules`, and `pnpm-store` volume mounted at `/root/.local/share/pnpm/store`; the task wipes `/dest/node_modules`, runs `pnpm install --frozen-lockfile`, `pnpm deploy /dest`, then `pnpm store prune`
- [ ] 4.2 Add `task node:build` call to `.devcontainer/post-start`

## 5. MVP verification — stop here and confirm before continuing

- [ ] 5.1 Build full image (`--target final`) with `--no-cache` and confirm success
- [ ] 5.2 Run `task node:build` and confirm the volume is populated
- [ ] 5.3 Confirm `openspec --version`, `opencode --version`, and `renovate --version` all succeed in the devcontainer
- [ ] 5.4 Confirm renovate runs without RE2 warning
- [ ] 5.5 Confirm `pnpm`, `npm`, `npx` are not on PATH in the final image
- [ ] 5.6 Run `task node:package:add -- some-test-pkg 1.0.0`, then `task node:build`, confirm package available without image rebuild

## 6. Taskfile and post-start restore

- [ ] 6.1 Restore `openspec --version`, `opencode --version`, and `renovate --version` checks in `task devcontainer:doctor`
- [ ] 6.2 Restore `task devcontainer:check-versions` task using renovate from the volume
- [ ] 6.3 Restore `task devcontainer:check-versions` call in `.devcontainer/post-start`

## 7. renovate.json supply chain settings

- [ ] 7.1 Add `minimumReleaseAge: "7 days"` to `renovate.json` at the top level

## 8. Sandboxed node package management tasks

- [ ] 8.1 Add `task node:package:add` that runs `pnpm add <package>@<version>` in a throwaway `symphony-studio-node-builder` container, bind-mounting only `.devcontainer/node/`
- [ ] 8.2 Add `task node:package:update` that runs `pnpm update <package>@<version>` in the same sandboxed container
- [ ] 8.3 Add `task node:package:rm` that runs `pnpm remove <package>` in the same sandboxed container and removes the package's `allowBuilds` entry from `pnpm-workspace.yaml` if present
- [ ] 8.4 Add `task node:package:list` that prints direct deps from `package.json` without starting a container
- [ ] 8.5 Add `task node:package:audit` that runs `pnpm audit` in a sandboxed container against the current lockfile; exits non-zero if vulnerabilities found
- [ ] 8.6 Add `task node:package:prune` that runs `pnpm store prune` in a `symphony-maestro-node-builder` container with the `${localWorkspaceFolderBasename}-pnpm-store` volume mounted; removes store entries no longer referenced by the current lockfile
- [ ] 8.6 Add `task node:trust:add` that sets `<package>: true` in `allowBuilds` in `pnpm-workspace.yaml` without touching `package.json` or the lockfile
- [ ] 8.7 Add `task node:trust:rm` that removes a package's entry from `allowBuilds` in `pnpm-workspace.yaml`
- [ ] 8.8 Add `task node:trust:list` that prints the current `allowBuilds` entries from `pnpm-workspace.yaml` without starting a container
- [ ] 8.9 Add `task node:trust:verify` that runs pnpm in a sandboxed container to identify transitive deps with build scripts not in `allowBuilds`; exits non-zero if any found
- [ ] 8.10 Add a comment header to `.devcontainer/node/package.json` documenting `task node:package:*` and `task node:trust:*` as the only supported interface

## 9. Verification of task interface

- [ ] 9.1 Verify `task node:package:add -- renovate 43.220.0` updates `package.json` and lockfile but does not touch `allowBuilds`
- [ ] 9.2 Verify `task node:trust:add -- re2` sets `re2: true` in `allowBuilds` without modifying `package.json`
- [ ] 9.3 Verify `task node:package:rm -- opencode-ai` removes from `package.json` and cleans up `allowBuilds`
- [ ] 9.4 Verify `task node:package:prune` removes unreferenced entries from the pnpm store volume
- [ ] 9.5 Run `task devcontainer:doctor` and confirm all tools pass after a `task node:build`
