## Context

In our initial design, we considered updating the project-wide Git skill and injecting rules globally. However, this approach creates cognitive noise (context bloat) for agents when they are focused on other unrelated tasks (like designing specifications or analyzing requirements). To implement a clean separation of concerns, we will isolate all CLI-specific Git operations into a single specialized subagent: `git-operator`. This subagent will be the sole executor of git commands, keeping git mechanics separate from high-level reasoning and design contexts.

Additionally, to prevent any risk of destructive operations (like force-pushing or bypassing branch rules), we enforce strict, fine-grained permission gating on the subagent. Only a safe, read-only or low-impact subset of Git and GitHub CLI commands are allowed automatically; sensitive commands require explicit human permission (`ask`), and dangerous commands are completely blocked (`deny`).

## Goals / Non-Goals

**Goals:**
- **Exclusive Permission Gating**: Ensure only the `git-operator` subagent is permitted to run `git` and `gh` bash commands.
- **Fine-Grained Sandbox Gating**: Restrict the `git-operator` to a safe subset of commands. Forbid direct execution of force-pushes, hard resets, or direct commits to main.
- **Single Source of Truth**: Delete the redundant `.opencode/skills/git` directory, consolidating all git execution logic in the subagent's system prompt.
- **Linter Gating**: Enforce that the `git-operator` runs `bin/commit-lint.ts` on all commit messages to guarantee strict conventional commit compliance.
- **Workflow-Focused Orchestration**: Update the General Orchestrator and other agents to delegate Git CLI commands to the `git-operator` rather than executing them directly.

**Non-Goals:**
- **Altering Core Git Config**: We avoid mutating local git configurations. We rely purely on pre-commit git hooks, validation scripts, and agent instruction alignment.

## Decisions

### Decision 1: Creating the git-operator Subagent with Safe Gating
- **Choice**: Create `.opencode/agents/git-operator.md` with clean, broad, and robust permissions:
  - **Allowed (`allow`)**: `git status`, `git diff`, `git log`, `git add`, `git commit *`, `git checkout *`, `git branch`, `gh pr status`, `gh pr list`, `gh issue list`, `gh issue view *`.
  - **Prompt-on-Execution (`ask`)**: `git push`, `git push *`, `git checkout main`, `git pull`, `git pull *`, `git reset`, `git reset *`, `git merge *`, `git cherry-pick *`, `gh pr *`, `gh issue create *`.
- **Rationale**: Rather than attempting to block dangerous command flags (like `-f` or `--force`) with fragile and easily-bypassed glob combinations, we assign all state-transitioning commands (like `push`, `reset`, `pull`, `merge`) to `ask`. Any variation of these commands (e.g., `git push -vf` or `git reset --hard`) will automatically trigger a human permission prompt. The user is presented with the exact execution string, ensuring complete sovereignty and eliminating bypass vectors. Placing the broad catch-all `"*": deny` at the top ensures specific overrides below it apply correctly.

### Decision 2: Revoking direct git permissions on other agents
- **Choice**: Explicitly deny `git` and `gh` command execution in all other agent frontmatters (`builder.md`, `designer.md`, `orchestrator.md`, `issue.md`, `advisor.md`) and globally in `opencode.json`.
- **Rationale**: Enforces a bulletproof "Execution Firewall" aligned with Symphony's Constitution.

### Decision 3: Consolidating Git Guidelines into the Subagent System Prompt
- **Choice**: Put all branch naming conventions (`change/<name>`), Conventional Commits format criteria, and atomic commit practices directly inside `.opencode/agents/git-operator.md`.
- **Rationale**: Avoids duplicate rules across skills and global system instructions.

## Risks / Trade-offs

- **[Risk] Slower local workflows due to multi-agent delegation**: Running a subagent to commit changes adds a small delay.
  - **Mitigation**: This delay is tiny compared to the gain in safety, traceability, and prevention of untracked files.
