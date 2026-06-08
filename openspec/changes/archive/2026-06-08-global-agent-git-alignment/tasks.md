## 1. Redundant Skill Cleanup

- [x] 1.1 Delete the redundant `.opencode/skills/git/` skill directory to remove duplicate logic.

## 2. Git Operator Subagent Implementation

- [x] 2.1 Create `.opencode/agents/git-operator.md` with exclusive, safe bash execution permissions (mapping safe commands to `allow`, and state-transitioning commands to `ask` to enforce human verification gates).
- [x] 2.2 Define the system prompt in `.opencode/agents/git-operator.md` to enforce the complete Git playbook (Conventional Commits, atomic commits, branch conventions, and pre-commit commit lint validation).

## 3. Permission Hardening

- [x] 3.1 Update `.opencode/agents/builder.md` to revoke direct `git` and `gh` command execution permissions (setting `bash` to deny them or ask).
- [x] 3.2 Update `.opencode/agents/orchestrator.md`, `.opencode/agents/designer.md`, and `.opencode/agents/issue.md` to ensure they have zero direct Git execution permissions.
- [x] 3.3 Ensure other files or guidelines do not reference direct git commands for non-git agents.
