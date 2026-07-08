## 1. docker-compose.yml

- [x] 1.1 Define an `x-versions` YAML anchor at the top of `.devcontainer/docker-compose.yml` with current tool version strings (`GO_VERSION`, `DENO_VERSION`, `TASK_VERSION`) matching current Dockerfile ARG defaults; add Renovate annotations
- [x] 1.2 Define an `x-build` YAML anchor with the shared build block (`context: ..`, `dockerfile: .devcontainer/Dockerfile`, `args` merged from `x-versions`)
- [x] 1.3 Add `devcontainer` service merging `x-build` with `build.target: final` and `image: symphony-maestro`; no profile

## 2. Dockerfile

- [x] 2.1 Remove default values from all tool version `ARG` declarations in `.devcontainer/Dockerfile` (values are now supplied exclusively by compose); preserve Renovate annotation comments

## 3. devcontainer.json

- [x] 3.1 Replace the `build` block with `dockerComposeFile: "docker-compose.yml"` and `service: "devcontainer"`
- [x] 3.2 Add explicit `workspaceMount` to bind-mount the workspace (required by devcontainer spec for compose configurations)
- [x] 3.3 Change `shutdownAction` from `stopContainer` to `stopCompose`

## 4. Verification

- [ ] 4.1 Rebuild the devcontainer and confirm it opens successfully
- [ ] 4.2 Confirm `task devcontainer:doctor` passes inside the container
- [ ] 4.3 Confirm all existing named volumes are still mounted (`vscode-extensions`, `vscode-user-data`)
