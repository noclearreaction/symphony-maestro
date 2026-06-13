## Why

The `reinstall-node-tools` change needs to add a `node-builder` service to the devcontainer
build so that image is available by name on the host Docker daemon for DooD tasks. Docker
Compose is the supported mechanism for this within devcontainers, but the current setup
uses `devcontainer.json` with a plain `build.dockerfile` entry.

This change migrates the devcontainer build to Docker Compose first, so `reinstall-node-tools`
can add the `node-builder` service on a clean foundation. Moving the ARG version defaults
into compose at the same time eliminates the current duplication (same values in two files)
and makes compose the single source of truth.

## What Changes

- Add `.devcontainer/docker-compose.yml` with YAML anchors for versions and shared build
  config; initial `devcontainer` service only (`--target final`)
- All tool version ARG defaults move from the Dockerfile into the compose file; Dockerfile
  ARGs become required (no defaults)
- Replace `build` block in `devcontainer.json` with `dockerComposeFile` and `service`
- All existing `mounts`, `features`, `postStartCommand`, `remoteUser`, and `customizations`
  are preserved unchanged

## Capabilities

### Modified Capabilities

- `devcontainer-environment`: devcontainer build mechanism changes from `build.dockerfile`
  to `dockerComposeFile`; compose is now the single source of truth for tool version defaults

## Impact

- `.devcontainer/docker-compose.yml`: new file
- `.devcontainer/Dockerfile`: ARG defaults removed (versions now supplied exclusively by compose)
- `.devcontainer/devcontainer.json`: `build` block replaced with `dockerComposeFile` + `service`
