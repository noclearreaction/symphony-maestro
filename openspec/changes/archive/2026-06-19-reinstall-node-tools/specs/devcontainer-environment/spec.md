## MODIFIED Requirements

### Requirement: Multi-stage Dockerfile assembles all runtimes
The devcontainer SHALL be built from a single multi-stage `Dockerfile` with no external host tooling prerequisites beyond `docker build`. The Dockerfile SHALL contain a `node-builder` stage with `NODE_VERSION` pinned via ARG. The pnpm version SHALL NOT be in the Dockerfile — it is owned by the `packageManager` field in `.devcontainer/node/package.json`, read by corepack at task runtime. The node-runtime stage SHALL be removed; node packages are provided at container start via a named Docker volume populated by `task node:build`.

#### Scenario: Container builds from clean Docker cache
- **WHEN** a developer runs `docker build --target final .devcontainer/`
- **THEN** the build succeeds, producing an image with Go, Deno, Task, `openspec`, `opencode`, and `gh` available

#### Scenario: CI target builds independently
- **WHEN** a CI runner builds with `--target ci`
- **THEN** the image contains Go, Deno, Task, `openspec`, `opencode`, and `gh`, but excludes SSH profile scripts and vscode user directory setup

---

### Requirement: Node tools available via /opt/node_modules/.bin on PATH
`/opt/node_modules/.bin` SHALL be on `PATH` in the devcontainer image (`ENV PATH=/opt/node_modules/.bin:$PATH` baked at build time). No symlink is needed — `modulesDir: "/opt/node_modules"` in `pnpm-workspace.yaml` causes pnpm to write binaries directly to `/opt/node_modules/.bin` in the named volume. `pnpm`, `npm`, `npx`, and `corepack` SHALL NOT appear on `PATH`.

#### Scenario: pnpm not reachable by default
- **WHEN** a developer runs `pnpm` in a container shell
- **THEN** the shell returns "command not found"

#### Scenario: openspec reachable without pnpm on PATH
- **WHEN** a developer runs `openspec --version`
- **THEN** the command succeeds

#### Scenario: opencode reachable without pnpm on PATH
- **WHEN** a developer runs `opencode --version`
- **THEN** the command succeeds

#### Scenario: renovate reachable without pnpm on PATH
- **WHEN** a developer runs `renovate --version`
- **THEN** the command succeeds

---

### Requirement: A `doctor` task verifies tool wiring on demand and at container start
A Taskfile task `doctor` (namespaced as `devcontainer:doctor` from the root Taskfile) SHALL verify that all critical tools are reachable and correctly wired. It SHALL run automatically as part of `postStartCommand`. It SHALL exit non-zero if any tool check fails.

#### Scenario: doctor passes in a healthy container
- **WHEN** `task devcontainer:doctor` is run inside the devcontainer
- **THEN** it exits 0 and reports the version of each tool: Go, gofmt, Deno, Task, openspec, opencode, gh, and Docker

#### Scenario: doctor fails when a tool is missing
- **WHEN** a tool is absent or misconfigured
- **THEN** `task devcontainer:doctor` exits non-zero and the container start output surfaces the failure

---

### Requirement: check-versions task reports available updates via renovate
A Taskfile task `check-versions` (namespaced as `devcontainer:check-versions`) SHALL run renovate in `--dry-run=lookup` mode using `docker run renovate/renovate` via DooD. It SHALL run at container start and exit 0 regardless of findings.

#### Scenario: check-versions runs and reports findings
- **WHEN** `task devcontainer:check-versions` is run inside the devcontainer
- **THEN** renovate runs via Docker, scans the repo, and outputs version update information
- **AND** the task exits 0

#### Scenario: check-versions does not fail container start
- **WHEN** renovate finds outdated packages or cannot reach registries
- **THEN** `task devcontainer:check-versions` still exits 0
