## 1. Scaffold .devcontainer structure

- [x] 1.1 Create `.devcontainer/` directory with placeholder files: `Dockerfile`, `devcontainer.json`
- [x] 1.2 Determine and pin concrete versions for all ARGs: `GO_VERSION`, `DENO_VERSION`, `TASK_VERSION`, `NODE_VERSION`, `PNPM_VERSION`, `OPENSPEC_VERSION`, `OPENCODE_VERSION` — no `latest` or floating refs permitted; record resolved values in a comment block at the top of the Dockerfile

## 2. Dockerfile — download and runtime stages

- [x] 2.1 Write `download-base` stage: `debian:bookworm-slim` with `curl`, `ca-certificates`, `tar`, `xz-utils`; set `WORKDIR /srv/files`
- [x] 2.2 Write `go-runtime` stage: from `golang:<GO_VERSION>-alpine`; copy full Go tree to `/srv/files/go<GO_VERSION>/`; install any needed Go tools (e.g. `gopls`) to `/srv/files/go/bin/`
- [x] 2.3 Write `deno-runtime` stage: from `download-base`; download deno binary from official release URL to `/srv/files/deno-<DENO_VERSION>/bin/deno`; `chmod +x`
- [x] 2.4 Write `task-binary` stage: from `download-base`; download go-task release archive; extract `task` binary to `/srv/files/task`; `chmod +x`
- [x] 2.5 Write `node-runtime` stage: from `node:<NODE_VERSION>-bookworm-slim`; enable corepack; install pnpm at pinned version; install `openspec` and `opencode` globally via pnpm; stage output at `/srv/files/node/`

## 3. Dockerfile — base and ci stages

- [x] 3.1 Write `base` stage: from `mcr.microsoft.com/devcontainers/base:ubuntu-22.04`; `COPY --from=go-runtime` to `/opt/go<GO_VERSION>/`; register via `update-alternatives --install /usr/local/bin/go go /opt/go<GO_VERSION>/bin/go 100 --slave /usr/local/bin/gofmt gofmt /opt/go<GO_VERSION>/bin/gofmt`; set `GOROOT=/opt/go<GO_VERSION>` and `GOPATH=/home/vscode/go` in ENV
- [x] 3.2 `COPY --from=deno-runtime` to `/opt/deno-<DENO_VERSION>/`; register via `update-alternatives --install /usr/local/bin/deno deno /opt/deno-<DENO_VERSION>/bin/deno 100`
- [x] 3.3 `COPY --from=task-binary` to `/usr/local/bin/task`
- [x] 3.4 `COPY --from=node-runtime /srv/files/node /opt/node`; add only `node` binary to PATH; create `/usr/local/bin/openspec` and `/usr/local/bin/opencode` wrapper scripts pointing into `/opt/node/`; verify `pnpm`, `npm`, `npx` are NOT on PATH
- [x] 3.5 Write `ci` stage: `FROM base AS ci`; set CI-appropriate ENV (`CI=true`, `DEBIAN_FRONTEND=noninteractive`); no user-facing config

## 4. Dockerfile — final stage

- [x] 4.1 Write `final` stage: `FROM ci AS final`; copy SSH agent profile script to `/etc/profile.d/`; source it from `/home/vscode/.bashrc` and `/home/vscode/.zshrc`
- [x] 4.2 Create vscode user directories: `~/.vscode-server/extensions`, `~/.vscode-server/data/User`; set ownership
- [x] 4.3 Set XDG env vars (`XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_CACHE_HOME`, `XDG_STATE_HOME`) and `PAGER=cat`
- [x] 4.4 Add `RUN date +%s > /container_build_id` for build tracking

## 5. devcontainer.json

- [x] 5.1 Write `devcontainer.json`: `name`, `build.dockerfile`, `build.context`, `workspaceFolder`
- [x] 5.2 Add devcontainer features: `common-utils` (zsh as default shell), `github-cli`, `sshd`, `docker-outside-of-docker` (moby: false, installDockerBuildx: false)
- [x] 5.3 Add `mounts`: named volume for `vscode-server/extensions`, named volume for `vscode-server/data/User`
- [x] 5.4 Add `postCreateCommand` to fix ownership of vscode user dirs
- [x] 5.5 Add `postStartCommand` referencing `.devcontainer/post-start` script; create that script with: `task devcontainer:doctor`; add commented stubs for future hooks and MCP sidecar sections
- [x] 5.6 Set `remoteUser: vscode`, `shutdownAction: stopContainer`
- [x] 5.7 Add `customizations.vscode.extensions` — deferred: resolve in a follow-on change once container is in use and extension needs are observed; add an empty array for now

## 6. Fix bin/ script shebangs

- [x] 6.1 Update `bin/commit-lint.ts` shebang from hardcoded host path to `#!/usr/bin/env -S deno run --allow-read --allow-env`
- [x] 6.2 Update `bin/provision-labels.ts` shebang from hardcoded host path to `#!/usr/bin/env -S deno run --allow-run=gh`
- [ ] 6.3 Verify both scripts run correctly inside the container

## 7. Verification

- [x] 7.1 Build `--target ci` image; confirm Go, Deno, Task, openspec, opencode, gh present; confirm pnpm/npm/npx absent from PATH
- [ ] 7.2 Build `--target final` image; open in VS Code Dev Containers; confirm `task devcontainer:doctor` passes
- [ ] 7.3 Confirm `docker ps` succeeds inside the container (DooD working)
- [x] 7.4 Confirm `update-alternatives --display go` shows correct versioned path and lists `gofmt` as a slave; confirm `update-alternatives --display deno` shows correct versioned path
- [x] 7.5 Confirm `gofmt --help` works after switching go alternative
- [ ] 7.6 Confirm `post-start` script runs without error and `task devcontainer:doctor` exits 0

## 8. Taskfile

- [x] 8.1 Create `.devcontainer/Taskfile.yaml` with a `doctor` task that checks: `go version`, `gofmt -h`, `deno --version`, `task --version`, `openspec --version`, `opencode --version`, `gh --version`, `docker info` (socket access); exit non-zero on any failure
- [x] 8.2 Create root `Taskfile.yaml` including `.devcontainer/Taskfile.yaml` under the `devcontainer:` namespace; add a top-level `health` alias that calls `devcontainer:doctor`
- [ ] 8.3 Verify `task devcontainer:doctor` runs cleanly inside the built container
- [ ] 8.4 Verify `task devcontainer:doctor` exits non-zero when a tool is removed (manual destructive test, restore after)
