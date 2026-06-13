## Context

The devcontainer currently uses `devcontainer.json` with a `build.dockerfile` entry
targeting the `final` stage. The `reinstall-node-tools` change will need to add a
`node-builder` service to the compose file so that image is available by name on the
host Docker daemon. That requires compose to already be the build driver.

This change migrates the devcontainer build from `build.dockerfile` to `dockerComposeFile`,
establishing the compose foundation that `reinstall-node-tools` extends. It also moves
tool version ARG defaults into compose so they are maintained in one place.

## Goals / Non-Goals

**Goals:**
- Migrate devcontainer build from `build.dockerfile` to `dockerComposeFile`
- All existing devcontainer behaviour preserved (mounts, features, postStartCommand, remoteUser)
- Compose is the single source of truth for all tool version defaults
- Compose file structured with YAML anchors so adding services later is clean

**Non-Goals:**
- Adding the `node-builder` service (that is in `reinstall-node-tools`)
- Adding any new Dockerfile stages
- Changing any task or post-start behaviour

## Decisions

### D1: docker-compose.yml with devcontainer service only

**Decision**: `.devcontainer/docker-compose.yml` defines a single `devcontainer` service
(`--target final`). `devcontainer.json` switches from `build.dockerfile` to
`dockerComposeFile` + `service: devcontainer`. The file is structured with `x-versions`
and `x-build` anchors so that `reinstall-node-tools` can add `node-builder` cleanly.

**Rationale**: Separating the compose migration from the node-builder addition keeps each
change minimal and independently verifiable.

### D2: Version ARG defaults move from Dockerfile to compose

**Decision**: All `ARG` defaults are removed from the Dockerfile. The compose file becomes
the single source of truth for tool version defaults, declared via a YAML anchor
(`x-versions`) and passed as `build.args` to every service. Renovate manages version
bumps directly in `docker-compose.yml` via the same Renovate regex manager that
currently targets the Dockerfile.

**Rationale**: With compose as the build driver, the Dockerfile ARG defaults are never
used — compose always supplies them explicitly. Keeping defaults in two places creates
drift. Making compose the single source of truth eliminates the duplication.

**Note**: `GO_VERSION`, `DENO_VERSION`, and `TASK_VERSION` are already in the Dockerfile
with Renovate annotations. These move to compose. `NODE_VERSION` and `PNPM_VERSION`
will be added to the `x-versions` anchor by `reinstall-node-tools`.

### D3: YAML anchors for shared build configuration

**Decision**: `docker-compose.yml` uses a top-level `x-versions` YAML anchor for version
strings and an `x-build` anchor for the shared `build` block (`context`, `dockerfile`,
`args`). Services merge `x-build`.

**Rationale**: Consistent with the project's existing compose sample style. When
`reinstall-node-tools` adds `node-builder`, it merges the same anchors with no duplication.

### D4: workspaceFolder and shutdownAction in devcontainer.json

**Decision**: When switching to `dockerComposeFile`, `workspaceMount` must be explicitly
specified in `devcontainer.json` because compose does not automatically bind-mount the
workspace. `shutdownAction` changes from `stopContainer` to `stopCompose`.

**Rationale**: Required by the devcontainer spec for compose-based configurations.
