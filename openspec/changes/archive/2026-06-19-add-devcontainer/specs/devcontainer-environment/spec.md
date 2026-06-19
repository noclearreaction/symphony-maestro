## ADDED Requirements

### Requirement: Multi-stage Dockerfile assembles all runtimes
The devcontainer SHALL be built from a single multi-stage `Dockerfile` with no external host tooling prerequisites beyond `docker build`.

#### Scenario: Container builds from clean Docker cache
- **WHEN** a developer runs `docker build --target final .devcontainer/`
- **THEN** the build succeeds, producing an image with Go, Deno, Task, `openspec`, `opencode`, and `gh` available

#### Scenario: CI target builds independently
- **WHEN** a CI runner builds with `--target ci`
- **THEN** the image contains Go, Deno, Task, `openspec`, `opencode`, and `gh`, but excludes SSH profile scripts and vscode user directory setup

---

### Requirement: Go and gofmt available on PATH with atomic version switching
Go SHALL be available on `PATH` at the pinned version. `go` and `gofmt` SHALL be managed as an atomic pair — switching the active Go version SHALL automatically switch `gofmt` to the corresponding version. `GOROOT` SHALL point at the concrete versioned installation directory, not a symlink.

#### Scenario: Go binary reachable on PATH
- **WHEN** a shell session starts in the container
- **THEN** `go version` outputs the expected version and `gofmt --help` succeeds

#### Scenario: gofmt slave follows go alternative
- **WHEN** `update-alternatives --set go /opt/go<VERSION>/bin/go` is called
- **THEN** `update-alternatives --display gofmt` shows the corresponding `/opt/go<VERSION>/bin/gofmt` as the active slave

---

### Requirement: Deno available on PATH at pinned version
Deno SHALL be available on `PATH` at the pinned version. The installation SHALL be at a versioned path to enable future side-by-side version management.

#### Scenario: Deno binary reachable on PATH
- **WHEN** a shell session starts in the container
- **THEN** `deno --version` outputs the expected version

---

### Requirement: Task installed from standalone binary
Task SHALL be installed directly to `/usr/local/bin/task` from the official go-task GitHub release archive. It SHALL NOT be installed via npm or pnpm.

#### Scenario: Task runs without Node
- **WHEN** `task --version` is run in the container
- **THEN** the command succeeds and the version matches the pinned `TASK_VERSION` ARG

---

### Requirement: Node and pnpm isolated from dev PATH
Node SHALL be installed to `/opt/node/`. Only the `node` binary SHALL be on `PATH`. `pnpm`, `npm`, `npx`, and `corepack` SHALL NOT appear on `PATH`. `openspec` and `opencode` CLI wrappers SHALL be the only pnpm-installed tools exposed in `/usr/local/bin/`.

#### Scenario: pnpm not reachable by default
- **WHEN** a developer runs `pnpm` in a container shell
- **THEN** the shell returns "command not found"

#### Scenario: openspec reachable without pnpm on PATH
- **WHEN** a developer runs `openspec --version`
- **THEN** the command succeeds

#### Scenario: opencode reachable without pnpm on PATH
- **WHEN** a developer runs `opencode --version`
- **THEN** the command succeeds

---

### Requirement: bin/ scripts use portable Deno shebang
All scripts in `bin/` that invoke Deno SHALL use `#!/usr/bin/env deno` (or `#!/usr/bin/env -S deno run ...`) rather than hardcoded host paths (e.g. `/home/tunnel49/.deno/bin/deno`).

#### Scenario: commit-lint runs inside container
- **WHEN** `bin/commit-lint.ts` is executed inside the devcontainer
- **THEN** it invokes Deno successfully without a "No such file or directory" error

#### Scenario: provision-labels runs inside container
- **WHEN** `bin/provision-labels.ts` is executed inside the devcontainer
- **THEN** it invokes Deno successfully without a "No such file or directory" error

---

### Requirement: devcontainer.json wires VS Code to the Dockerfile build
A `devcontainer.json` SHALL reference the Dockerfile via `build.dockerfile` targeting the `final` stage. It SHALL include the `docker-outside-of-docker` devcontainer feature. It SHALL mount named volumes for VS Code server extensions and user data, consistent with the existing sibling-repo pattern.

#### Scenario: All expected tools present after container opens
- **WHEN** a developer reopens the repository in the devcontainer
- **THEN** `go version`, `deno --version`, `task --version`, `openspec --version`, `opencode --version`, and `gh --version` all succeed in the integrated terminal

#### Scenario: Docker available inside container
- **WHEN** a developer runs `docker ps` inside the container
- **THEN** the command succeeds, communicating with the host Docker daemon via DooD

---

### Requirement: A `doctor` task verifies tool wiring on demand and at container start
A Taskfile task `doctor` (namespaced as `devcontainer:doctor` from the root Taskfile) SHALL verify that all critical tools are reachable and correctly wired. It SHALL run automatically as part of `postStartCommand`. It SHALL exit non-zero if any tool check fails.

#### Scenario: doctor passes in a healthy container
- **WHEN** `task devcontainer:doctor` is run inside the devcontainer
- **THEN** it exits 0 and reports the version of each tool: Go, gofmt, Deno, Task, openspec, opencode, gh, and Docker

#### Scenario: doctor fails when a tool is missing
- **WHEN** a tool is absent or misconfigured (e.g. a broken update-alternatives link)
- **THEN** `task devcontainer:doctor` exits non-zero and the container start output surfaces the failure

#### Scenario: doctor runs automatically on container start
- **WHEN** VS Code starts the devcontainer
- **THEN** `postStartCommand` invokes `task devcontainer:doctor` and its output is visible in the terminal

---

### Requirement: All tool versions explicitly pinned in Dockerfile
Every tool installed in the Dockerfile SHALL have a corresponding `ARG <TOOL>_VERSION` declared at the top of the file with a concrete pinned value. No tool SHALL be installed using `latest`, `@latest`, or any floating version reference at build time.

#### Scenario: Dockerfile ARGs enumerate all tool versions
- **WHEN** a developer reads the top of the Dockerfile
- **THEN** they find explicit ARG declarations for GO_VERSION, DENO_VERSION, TASK_VERSION, NODE_VERSION, PNPM_VERSION, OPENSPEC_VERSION, and OPENCODE_VERSION, each with a concrete pinned value

#### Scenario: No floating version references in RUN instructions
- **WHEN** the Dockerfile is scanned for `@latest` or unversioned install commands
- **THEN** none are found in any `RUN` instruction that installs a tool
