## Context

Tool versions are pinned as `ARG X=y.z` defaults in `.devcontainer/Dockerfile`. Currently there is no automated signal when a newer version is available. Renovate is a dependency update bot that opens PRs with proposed version bumps; it does not apply them unless explicitly configured to automerge.

The repo is hosted on GitHub. Renovate is available as a hosted GitHub App, a self-hosted runner, and a Docker image that can run locally. This change verifies configuration using the local Docker dry-run (`--platform=local`). Production scheduling (opening PRs on GitHub) can be wired via GitHub Actions or the hosted app as a follow-on step.

## Goals / Non-Goals

**Goals:**
- Raise PR-based proposals for version updates to all pinned tools in `.devcontainer/Dockerfile`
- Group all devcontainer-related bumps into a single PR to reduce noise
- Keep `automerge: false` on every rule — no version is applied without human review and approval
- Validate `renovate.json` is parseable and schema-correct

**Non-Goals:**
- Installing the Renovate GitHub App (requires org/repo admin access — done manually post-change)
- Auto-merging any updates
- Managing versions outside `.devcontainer/Dockerfile` (no `package.json`, no other lockfiles in this repo)
- Opening PRs, pushing branches, or any GitHub interaction

## Decisions

### D1 — `dockerfile` manager with inline datasource comments

Renovate's built-in `dockerfile` manager understands `ARG NAME=value` lines. For standard names like `GO_VERSION` and `NODE_VERSION` it resolves datasources automatically. For non-standard names (`DENO_VERSION`, `TASK_VERSION`, `PNPM_VERSION`, `OPENSPEC_VERSION`, `OPENCODE_VERSION`), a `# renovate: datasource=... depName=...` comment on the ARG line is required.

**Alternative considered**: A `.env` file with a custom regex manager. Rejected — the devcontainer framework builds directly from the Dockerfile, making the Dockerfile the correct single source of truth (see prior exploration in `add-devcontainer`).

### D2 — `npm` datasource for `OPENSPEC_VERSION` and `OPENCODE_VERSION`

These two ARGs reference npm packages, not GitHub releases. The `dockerfile` manager resolves them with a `# renovate: datasource=npm depName=...` comment.

### D3 — `--dry-run=lookup` not `--dry-run=full`

`--dry-run=lookup` queries registries and reports current vs available versions without performing any git operations. `--dry-run=full` simulates the full PR workflow. Only `lookup` is appropriate for a passive on-start check.

### D4 — Renovate installed in the devcontainer image

Renovate is an npm package (`renovate`). It is installed via pnpm in the `node-runtime` stage alongside `openspec` and `opencode`, and exposed via the same wrapper script pattern (`/usr/local/bin/renovate` → `node .../renovate.js`). Its version is pinned as `RENOVATE_VERSION` in the ARG block and tracked by itself via `datasource=npm depName=renovate`.

This is simpler and more reliable than running it via `docker run` (no DooD dependency, no silent exit issues with `--platform=local` inside a container).

### D5 — Non-blocking on container start

The `check-versions` task prints a report but does not exit non-zero if updates are available. It is advisory only — stale versions are not an error. The `post-start` script calls it after `doctor` so a network failure cannot break container startup.

## Risks / Trade-offs

- **ARG comment annotations are fragile** → If an ARG line is reformatted, Renovate may silently drop the annotation. Mitigation: spec requires annotations on the same line as the ARG.
- **`TASK_VERSION` points to `go-task/task` releases, not a package registry** → Must use `datasource=github-releases`. Verified pattern; Renovate supports it.
- **Network unavailable on container start** → `renovate` will fail to reach registries. Mitigation: `post-start` calls check-versions with `|| true` so a network failure is logged but does not break startup.

## Open Questions

- Which Renovate major version tag to pin in the task — resolve during implementation by checking the current latest major.
