## ADDED Requirements

### Requirement: Renovate config exists and is schema-valid
The repository SHALL contain a `renovate.json` at the repo root that is parseable JSON and conforms to the Renovate config schema.

#### Scenario: Valid config present
- **WHEN** Renovate processes the repository
- **THEN** no config parse errors are reported

### Requirement: renovate.json must be git-tracked
`renovate.json` SHALL be committed or staged in git. The local platform reads files via git and will silently skip untracked files.

#### Scenario: renovate.json is untracked
- **WHEN** `task devcontainer:check-versions` runs with an untracked `renovate.json`
- **THEN** the regex manager reports 0 files and 0 deps

#### Scenario: renovate.json is staged or committed
- **WHEN** `task devcontainer:check-versions` runs
- **THEN** the regex manager reports `fileCount: 1, depCount: 8`

### Requirement: Devcontainer tool versions are tracked
The Renovate config SHALL configure the `dockerfile` manager to track all version ARGs in `.devcontainer/Dockerfile`.

Tracked ARGs:
- `GO_VERSION` — golang releases
- `DENO_VERSION` — github-releases `denoland/deno`
- `TASK_VERSION` — github-releases `go-task/task`
- `NODE_VERSION` — node releases
- `PNPM_VERSION` — npm `pnpm`
- `OPENSPEC_VERSION` — npm `@fission-ai/openspec`
- `OPENCODE_VERSION` — npm `opencode-ai`

#### Scenario: Lookup finds a newer version
- **WHEN** `task devcontainer:check-versions` is run
- **THEN** the output reports the current pinned version and the latest available version for each tracked ARG

### Requirement: Version check is non-blocking
The `devcontainer:check-versions` task SHALL exit 0 whether or not updates are available. It is advisory only.

#### Scenario: Updates are available
- **WHEN** one or more tools have a newer version available
- **THEN** the task prints the update report and exits 0

#### Scenario: Network is unavailable
- **WHEN** the registry cannot be reached
- **THEN** the task prints an error message and exits 0 (post-start script uses `|| true`)

### Requirement: Version check runs on container start
The `.devcontainer/post-start` script SHALL call `task devcontainer:check-versions` after `task devcontainer:doctor`.

#### Scenario: Container starts
- **WHEN** the devcontainer starts
- **THEN** `devcontainer:check-versions` runs and its output is visible in the terminal

### Requirement: Pre-release versions are ignored
The Renovate config SHALL set `ignoreUnstable: true` so that pre-release versions are not reported as available updates.

#### Scenario: A pre-release npm version is published
- **WHEN** `opencode-ai@2.0.0-beta.1` is published
- **THEN** the check-versions task does NOT report it as an available update
