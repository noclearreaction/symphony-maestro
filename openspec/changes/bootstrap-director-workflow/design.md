## Context

Currently, the Director repository is completely strategic (containing goals, roadmap, and decisions), but lacks the practical execution plumbing required to build and refine itself safely. To break this bootstrapping paradox, we are designing our very first OpenSpec change to dogfood the development workflow. This document details the technical implementation of our local-first, git-backed, and review-gated developer experience.

## Goals / Non-Goals

**Goals:**
- **Dynamic Builder Agent**: Specify and register a custom `builder` agent under `.opencode/agents/builder.md` with strict prompt-level boundaries.
- **Isolate Agent Workspace**: Set up a structured, untracked directory tree under `.symphony/` to isolate agent logs, plans, reviews, and test evidence.
- **Git Branch Automation**: Develop a local script (`bin/director-start`) to automate starting an OpenSpec change and switching to a dedicated `change/<name>` branch.
- **GitHub PR Integration**: Develop a local script (`bin/director-submit`) to automate syncing delta specs, pushing the branch, and creating a GitHub Pull Request using `gh`.

**Non-Goals:**
- **No External Orchestrator**: We are not running background daemons, queue workers, or listening to incoming webhooks in this bootstrap phase.
- **No Direct Merge**: The automated scripts are blocked from merging directly into `main`. The human must review and merge the PR on GitHub.
- **No Database**: All project state, task lists, and specifications live as plain, versioned files in git.

## Decisions

### 1. Custom `builder` Agent via Dynamic Self-Loading Markdown
- **Choice**: Create `.opencode/agents/builder.md` instead of enabling the default `build` agent in `opencode.json`.
- **Rationale**: OpenCode dynamically loads files in `.opencode/agents/` as callable custom agents. This allows us to keep `opencode.json` minimal and clean.
- **Alternatives Considered**: Modifying the global permissions schema in `opencode.json` for the generic `build` agent. This was rejected because the generic `build` agent lacks domain-specific instructions about keeping strategic files pristine.

### 2. Segmented `.symphony/` Folder Tree
- **Choice**: Introduce a single `.symphony/` directory structure containing `scratchpad/`, `plans/`, `reviews/`, and `evidence/`, and append `.symphony/` to `.gitignore`.
- **Rationale**: Agents need workspace space to log raw text, prepare step-by-step task plans, and record test evidence without polluting the main git tree.
- **Alternatives Considered**: Using a system-wide `/tmp` or placing untracked files directly in the root directory. This was rejected as it makes auditing harder and results in a messy git workspace.

### 3. Bash Shell Wrapper Scripts for Git & GitHub Automation
- **Choice**: Implement two lightweight bash scripts under `bin/`:
  1. `bin/director-start <change-name>`: Validates the name, creates/checks out the git branch `change/<change-name>`, runs `openspec new change "<change-name>"`, and prints out instructions to the user.
  2. `bin/director-submit`: Ensures the user is on a feature branch, runs `openspec sync` to promote delta specs to the main `specs/` directory, commits the synced changes, pushes the branch, and executes `gh pr create` to open a pull request.
- **Rationale**: Shell scripts are highly portable, fast, easily parsed by the human, and have zero running overhead.
- **Alternatives Considered**: A full NodeJS CLI or Python orchestrator. This was rejected as over-engineered for the bootstrap phase. We want to start simple and scale up if needed.

## Risks / Trade-offs

- **[Risk] GitHub CLI (`gh`) is not installed or authenticated**
  - *Mitigation*: The `director-submit` script will verify that `gh` is installed and the user is authenticated (`gh auth status`) before attempting to push. If verification fails, it will gracefully fallback and print the manual git push and PR commands for the user.
- **[Risk] Builder Agent bypasses Soft Boundaries**
  - *Mitigation*: The builder agent has no capabilities to push commits directly to the remote repository. The human reviews the entire git diff of the PR before merging.
- **[Risk] Diverged Main Specs**
  - *Mitigation*: The `director-submit` script automatically runs `openspec sync` prior to commit, ensuring that any spec changes are fully integrated into the main `specs/` folder and committed as a single unit on the feature branch.
