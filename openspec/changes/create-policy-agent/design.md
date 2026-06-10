## Context

The Symphony Director relies on a constellation of specialized AI agents with strict sandbox boundaries to ensure high-fidelity coordination and protect the codebase from strategic drift or unauthorized modifications. Currently, there is no dedicated agent for managing AI-related configurations, skills, commands, system prompts, or global agent documentation (`AGENTS.md`). While other agents like `@builder` are explicitly blocked from modifying `.opencode/*`, they also cannot assist in prompt engineering or policy updates. Conversely, allowing more general agents to write to `.opencode/*` risks breaking the execution firewall. 

Introducing a dedicated `@policy` worker agent will isolate prompt and policy management to a specialized sandbox, preventing other implementation agents from tampering with system instructions while offering a structured way to maintain and refine agent-related files.

## Goals / Non-Goals

**Goals:**
- Define the identity, prompt structure, and sandbox permissions for the new `@policy` agent.
- Ensure `@policy` can edit `.opencode/*` and `AGENTS.md` but is strictly blocked from editing application source files or running arbitrary bash commands.
- Securely integrate `@policy` with the existing agent architecture and coordinate its git actions via `@git-operator`.
- Formulate a precise, sequential implementation checklist in `tasks.md`.

**Non-Goals:**
- Creating or editing any runtime application code or build scripts.
- Modifying the sandbox constraints of existing agents (e.g., we will not grant `@builder` or `@designer` edit access to `.opencode/*`).
- Automating the PR merge gate (PR merge remains a human-in-the-loop task).

## Decisions

### Decision 1: Agent Registration and File Location
- **Choice**: Store the agent definition in `.opencode/agents/policy.md` and register the agent as `@policy`.
- **Rationale**: Follows the existing dynamic loading pattern where OpenCode automatically instantiates agents from `.opencode/agents/*.md`.
- **Alternatives Considered**: Creating a central config JSON or a nested directory structure. These were rejected as they do not align with OpenCode's native agent loading mechanism.

### Decision 2: Permission Set & Security Firewall
- **Choice**: Apply strict, fine-grained tool and file permissions in the agent's frontmatter:
  ```yaml
  permission:
    edit:
      "*": deny
      ".opencode/*": allow
      "AGENTS.md": allow
    read:
      "*": allow
      ".opencode/*": allow
    bash:
      "*": deny
      "openspec validate *": allow
      "openspec status *": allow
  ```
- **Rationale**: This enforces the principle of least privilege. `@policy` has full reading context and is authorized to edit configuration and documentation files, but cannot touch application code (`bin/*.ts`, etc.) or execute raw bash scripts. It can only execute safe `openspec validate` commands to verify schema compliance.
- **Alternatives Considered**: 
  - *Allowing git execution directly*: Rejected. To comply with `controlled-git-workflows` (Requirement: Isolated Git Execution Subagent), `@policy` must delegate git actions to `@git-operator`.
  - *Restricting read access to `.opencode/*`*: Rejected. `@policy` needs to read other agents' prompts to ensure they are synchronized with global policies and `AGENTS.md`.

### Decision 3: Collaborative Boundary in `AGENTS.md`
- **Choice**: Update the repository's `AGENTS.md` file to explicitly define `@policy`'s role and distinguish its scope from `@advisor` and `@designer`.
- **Rationale**: Avoids overlap in duties:
  - `@advisor` owns high-level strategy, goals, and decision logs (lives under `strategy/`).
  - `@designer` owns feature planning, technical designs, and OpenSpec artifacts (lives under `openspec/`).
  - `@policy` owns prompt engineering, agent configuration, permissions, and coordination rules (lives under `.opencode/` and `AGENTS.md`).

### Decision 4: Execution & Bootstrapping Strategy
- **Choice**: Since `@builder` cannot edit files under `.opencode/*` (due to its own sandbox), the initial creation of `.opencode/agents/policy.md` must be performed by the human operator, or by `@policy` itself once the file is bootstrapped.
- **Rationale**: Preserves the security model of `@builder`. We must not temporarily elevate `@builder`'s permissions to write this file, as doing so violates Constitutional Principle 3 (Intent Traceability) and compromises the firewall.

## Risks / Trade-offs

- **[Risk] Self-Modification Vulnerability** → Since `@policy` is authorized to edit `.opencode/*`, it can theoretically modify its own system prompt or sandbox permissions.
  - *Mitigation*: The human operator remains the final authority for all branch merges. Furthermore, the `controlled-git-workflows` prevents automatic pushing or merging of branches, and any prompt modification will be highly visible in the git diff during code review.
- **[Risk] Overlap with @advisor** → `@policy` and `@advisor` both deal with "rules" and "policies".
  - *Mitigation*: Define a strict boundary: `@advisor` handles strategic decisions and long-term project direction (business/strategy level), while `@policy` handles prompt system instructions and tool permissions (technical execution level).
