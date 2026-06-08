## Why

To avoid context bloat and prevent rule drift across multiple agents, we must isolate Git CLI execution into a dedicated, sandboxed environment. Relying on a repository-wide Git skill causes redundant instruction loading for agents that have no business interacting with the terminal. By establishing a single specialized `git-operator` subagent, we cleanly enforce the execution firewall, lock down shell permissions, and maintain a singular, robust source of truth for commit linting and branch boundaries.

## What Changes

- **Skill Deprecation**: Remove the redundant OpenCode Git skill (`.opencode/skills/git/`) completely.
- **Subagent Creation**: Introduce `.opencode/agents/git-operator.md`, a highly constrained, single-purpose subagent. This is the **only** agent with permission to run local `git` and `gh` CLI commands.
- **Permission Hardening**: Revoke all `git *` and `gh *` bash permissions from other agents (Orchestrator, Designer, Advisor, Issue, and Builder) to establish a strict Execution Firewall.
- **Agent Awareness**: Train the `git-operator` system instructions to strictly enforce Conventional Commits, branch conventions, and atomic, logical commits.

## Capabilities

### New Capabilities
<!-- None, this is a modification of existing specifications -->

### Modified Capabilities
- `controlled-git-workflows`: Refactor requirements to enforce the "Subagent-Only" Git execution model.

## Impact

- **Execution Firewall**: Non-execution agents (Orchestrator, Designer, etc.) can no longer run `git` bash commands directly. They must delegate Git operations to the `git-operator` subagent.
- **Clean Contexts**: Agents only load the git-specific conventions and playbook when the `git-operator` context is spawned.
