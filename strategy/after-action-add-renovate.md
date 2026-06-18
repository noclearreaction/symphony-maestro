# After-Action Report: add-renovate — 2026-06-13

## Objective

Install Renovate in the devcontainer image so that `task devcontainer:check-versions` runs on container start and reports available version updates for pinned tools in `.devcontainer/Dockerfile`. No GitHub, no scheduling, no PRs — local advisory output only.

## What Was Delivered

- `renovate.json` at repo root with a custom regex manager covering all `ARG` version pins in `.devcontainer/Dockerfile`
- 9 `ARG` lines with `# renovate:` annotations (go, deno, task, node, pnpm, openspec, opencode, renovate, re2)
- `RENOVATE_VERSION` and `RE2_VERSION` ARGs pinned and managed by renovate itself
- `task devcontainer:check-versions` task that runs `renovate --platform=local --dry-run=lookup` with correct env vars
- `task devcontainer:doctor` and `check-versions` called from `.devcontainer/post-start`

## What Went Wrong

### 1. Ignored the error message

The very first build failure said:

```
[ERR_PNPM_IGNORED_BUILDS] Ignored build scripts: ...re2...
Run "pnpm approve-builds" to pick which dependencies should be allowed to run scripts.
```

This was ignored. Instead, the following approaches were tried in sequence — none of which work in pnpm v11:

- `pnpm add --ignore-scripts` (blocks all scripts including re2 build)
- `onlyBuiltDependencies` in `package.json` (wrong config location in v11)
- `only-built-dependencies[]=re2` in `.npmrc` (not a valid pnpm v11 config key)
- Manually writing `pnpm-workspace.yaml` with `allowBuilds` (correct mechanism, but done manually instead of using the provided tool)

The user had to ask three separate times before `pnpm approve-builds` was actually tried.

### 2. Wrong order for approve-builds

When `pnpm approve-builds` was finally tried, it was run *before* `pnpm add`, so there were no placeholders in `pnpm-workspace.yaml` yet and it had nothing to approve. The correct sequence — which the docs describe — is:

1. `pnpm add ... || true` — fails but writes placeholder entries to `pnpm-workspace.yaml`
2. `pnpm approve-builds opencode-ai re2 '!...'` — flips approved packages to `true`
3. `pnpm install` — completes with approved build scripts

### 3. Hoisted linker does not copy native build artifacts

With `node-linker=hoisted`, pnpm writes package files into `node_modules/` but native build output (`.node` files) from install scripts stays in `.pnpm/<pkg>/node_modules/<pkg>/build/` — it is not copied to the hoisted location.

This meant `re2.node` was built but not findable by Node at runtime. The fix was to make `re2` an explicit direct dependency rather than relying on it as a transitive dep of renovate. When it is a direct dep, pnpm fully hoists it including build artifacts.

### 4. No fast feedback loop

Each iteration required a full Docker build (~5 minutes). There was no attempt to prototype the pnpm install sequence locally (in a throwaway container or the running devcontainer) before encoding it in the Dockerfile. This amplified every wrong guess into a 5-minute wait.

### 5. COPY path bug

`COPY ssh-agent.sh` failed because the build context is the workspace root, not `.devcontainer/`. This was a pre-existing bug exposed by the first full-image build. Fixed to `COPY .devcontainer/ssh-agent.sh`.

## What the Correct Approach Looks Like

```dockerfile
RUN pnpm add \
        "@fission-ai/openspec@${OPENSPEC_VERSION}" \
        "opencode-ai@${OPENCODE_VERSION}" \
        "renovate@${RENOVATE_VERSION}" \
        "re2@${RE2_VERSION}" || true \
 && pnpm approve-builds opencode-ai re2 '!@fission-ai/openspec' '!core-js-pure' '!dtrace-provider' '!protobufjs' \
 && pnpm install
```

The explicit allowlist is intentional. pnpm v11 blocks build scripts by default specifically because install scripts are a well-established npm supply chain attack vector (`ua-parser-js`, `node-ipc`, and others all used postinstall to execute malicious code). Using `--all` would defeat that protection. The allowlist must be reviewed and updated when new packages with build scripts are added. The `re2` explicit install ensures it is fully hoisted with its native build artifacts.

## Lessons

| # | Lesson |
|---|--------|
| 1 | Read the error message and act on it directly before trying alternatives |
| 2 | Test pnpm install sequences in a throwaway container before writing Dockerfile |
| 3 | `pnpm approve-builds` requires packages to already be partially resolved (run after a failed `pnpm add`) |
| 4 | `node-linker=hoisted` does not copy native build artifacts; make native-module packages explicit direct deps |
| 5 | `allowBuilds` in `pnpm-workspace.yaml` is the correct pnpm v11 config; `.npmrc` array syntax and `onlyBuiltDependencies` in `package.json` do not work |
| 6 | Docker `COPY` paths are relative to build context, not the Dockerfile location |

## Time Cost

Approximately 2 hours and many build iterations for what should have been a 15-minute task. The entire overage was caused by not acting on the error message in the first failure.
