## MODIFIED Requirements

### Requirement: Devcontainer built via docker-compose
The devcontainer SHALL be defined using `dockerComposeFile` in `devcontainer.json` pointing
to `.devcontainer/docker-compose.yml`. The compose file SHALL define a `devcontainer`
service built from the Dockerfile with `--target final`.

All tool version ARG defaults SHALL be defined in `docker-compose.yml` only. The
Dockerfile ARGs SHALL have no default values. The compose file is the single source of
truth for all version strings.

The compose file SHALL use YAML anchors (`x-versions`, `x-build`) so that additional
services can be added cleanly in future changes.

#### Scenario: Devcontainer environment unchanged after compose migration
- **WHEN** the devcontainer is rebuilt using `dockerComposeFile`
- **THEN** all existing tools (Go, Deno, Task, gh, Docker) remain available
- **AND** `task devcontainer:doctor` passes
- **AND** all existing named volumes (`vscode-extensions`, `vscode-user-data`) remain mounted

#### Scenario: All tool versions supplied by compose
- **WHEN** the devcontainer is built via compose
- **THEN** all tool version ARGs (`GO_VERSION`, `DENO_VERSION`, `TASK_VERSION`) are supplied by the compose file
- **AND** the Dockerfile contains no ARG default values for tool versions
