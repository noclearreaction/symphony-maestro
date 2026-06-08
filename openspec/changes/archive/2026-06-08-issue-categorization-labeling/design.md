## Context

To support systematic, multi-agent orchestrations, the repository requires a predictable, machine-readable Issue Categorization and Labeling schema. Currently, issues are created without standard labels, leading to administrative drift. This design outlines how we will provision the standardized labels, update the `@issue` subagent to automate labeling, and define state transition guidelines for development.

## Goals / Non-Goals

**Goals:**
- Provision exactly 12 standard labels in the repository (`type:*`, `status:*`, `priority:*`) with their exact colors and descriptions.
- Update `@issue` subagent instructions in `.opencode/agents/issue.md` to recommend type and priority labels during drafting and enforce `status:backlog` upon creation.
- Document the state machine and transition rules in `governance/issue-lifecycle.md`.

**Non-Goals:**
- Automating state transitions via active GitHub Actions or webhook integration in this change (transitions will be done manually or by orchestrators following the guidelines).
- Modifying other agents except the `@issue` agent.

## Decisions

### Decision 1: Label Provisioning Method
We will create a Deno TypeScript script `bin/provision-labels.ts` that utilizes the GitHub CLI (`gh`) to check, create, and update the 12 required labels.
- **Alternatives Considered:**
  1. Manual creation: High error potential and non-reproducible.
  2. Simple Shell Script: Platform-dependent, less robust error handling, and harder to test.
  3. Deno TypeScript Script: High platform portability, strong error-handling, standard JSON/object structure parsing, and a secure execution sandbox.
- **Rationale:** A Deno TypeScript script running with restricted sandbox permissions (specifically `--allow-run=gh` to execute the local GitHub CLI) provides a robust, cross-platform, type-safe, and highly secure implementation.

### Decision 2: Subagent Instruction Integration
We will update `.opencode/agents/issue.md` to define explicit label selection rules for the `@issue` subagent.
- **Rules to Add:**
  - Standard label dictionary defining all 12 labels.
  - Requirement to include recommended `type:*` and `priority:*` labels in the draft output.
  - Requirement to pass `--label "status:backlog"` plus the approved type/priority labels during the `gh issue create` call.

### Decision 3: Transition and Lifecycle Documentation
We will publish the issue state machine guidelines in `governance/issue-lifecycle.md`.
- **Rationale:** Keeping lifecycle and governance rules version-controlled under `governance/` ensures human and agent operators have a shared, permanent source of truth.

## Risks / Trade-offs

- **[Risk] GitHub API/CLI rate limits or permission errors when modifying labels** → *Mitigation:* Ensure script outputs clear messages on failure, and require human verification of label list.
- **[Risk] Overwriting existing customized labels in the repository** → *Mitigation:* The provisioning script will target only the 12 specified labels. It will safely skip or update them, avoiding blanket deletion of custom labels unless they conflict.
