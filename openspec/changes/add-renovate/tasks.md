## 1. Annotate Dockerfile ARGs

- [x] 1.1 Add `# renovate: datasource=golang-version depName=go` comment on the `ARG GO_VERSION` line
- [x] 1.2 Add `# renovate: datasource=github-releases depName=denoland/deno` comment on the `ARG DENO_VERSION` line
- [x] 1.3 Add `# renovate: datasource=npm depName=pnpm` comment on the `ARG PNPM_VERSION` line
- [x] 1.4 Add `# renovate: datasource=npm depName=@fission-ai/openspec` comment on the `ARG OPENSPEC_VERSION` line
- [x] 1.5 Add `# renovate: datasource=npm depName=opencode-ai` comment on the `ARG OPENCODE_VERSION` line

## 2. Create renovate.json

- [x] 2.1 Create `renovate.json` at repo root with `$schema`, `extends: ["config:recommended"]`, `automerge: false`, `ignoreUnstable: true`
- [x] 2.2 Add `packageRules` entry: `matchFileNames: [".devcontainer/Dockerfile"]`, `groupName: "devcontainer tools"`, `automerge: false`

## 3. Add check-versions task

- [x] 3.1 Add `RENOVATE_VERSION=43.220.0` ARG to Dockerfile with `# renovate: datasource=npm depName=renovate` annotation
- [x] 3.2 Add `renovate@${RENOVATE_VERSION}` to the `pnpm add` install block in the `node-runtime` stage
- [x] 3.3 Add `/usr/local/bin/renovate` wrapper script in the `base` stage (same pattern as openspec/opencode)
- [x] 3.4 Add `check-versions` task to `.devcontainer/Taskfile.yaml`: runs `renovate --platform=local --dry-run=lookup`; exits 0 regardless of result
- [x] 3.5 Update `.devcontainer/post-start` to call `task devcontainer:check-versions || true` after the `doctor` call

## 4. Verify

- [x] 4.1 Validate `renovate.json` parses as valid JSON (`python3 -m json.tool renovate.json`)
- [x] 4.2 Run `task devcontainer:check-versions` manually; confirm `.devcontainer/Dockerfile` is detected and all 8 ARGs appear in the output (`regex: fileCount: 1, depCount: 8`)
- [x] 4.3 Confirm no unknown datasource warnings in the output (GitHub token WARN for github-releases is expected without credentials)
- [x] 4.4 Confirm task exits 0
