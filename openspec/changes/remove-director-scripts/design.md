## Context

Previously, the custom Deno scripts `bin/director-start.ts` and `bin/director-submit.ts` were written to simplify starting and submitting OpenSpec changes. While convenient, they diverge from OpenSpec community standards, hide the native tool's behavior, and add extra maintenance overhead. We will remove these scripts and update all specifications and agent prompt contexts to use native `openspec` commands and native Git/PR flows (via the `git-operator` subagent).

## Goals / Non-Goals

**Goals:**
- **Script Removal**: Decommission and delete `bin/director-start.ts` and `bin/director-submit.ts`.
- **Specification Deprecation**: Remove the automated script requirements from `director-workflow` spec and specify native `openspec` CLI workflows.
- **Agent Prompts Alignment**: Clean up references to these custom scripts in all agent configurations.

**Non-Goals:**
- **Implementing a Custom Git Sync CLI**: We will not write another custom tool to sync specs. We rely purely on native OpenSpec workflows and manual/agent edits.

## Decisions

### Decision 1: Physical Deletion of the Custom Scripts
- **Choice**: Permanently remove `bin/director-start.ts` and `bin/director-submit.ts` from the `bin/` directory.
- **Rationale**: Keeps the codebase minimal, clean, and free of redundant, non-standard wrapper scripts, strictly supporting Constitutional Principle 5 (Tooling Discipline).

### Decision 2: Refactoring specs to use standard OpenSpec CLI
- **Choice**: Replace the automated start/submit requirements in the `director-workflow` spec with native OpenSpec toolchain requirements.
- **Rationale**: Ensures the specification represents stable, community-accepted OpenSpec best practices.

### Decision 3: Aligning Agent Contexts
- **Choice**: Remove any reference or guidance toward `bin/director-start` or `bin/director-submit` from `.opencode/agents/orchestrator.md`, `.opencode/agents/designer.md`, and others.
- **Rationale**: Prevents agents from attempting to run deleted scripts, ensuring they natively understand the standard `openspec` CLI commands.

## Risks / Trade-offs

- **[Risk] Slower manual setup**: Starting an OpenSpec change now requires two commands (`git checkout -b change/<name>` and `openspec new change "<name>"`) instead of a single `bin/director-start` script.
  - **Mitigation**: This minor increase in manual steps matches community standards exactly, increases developer/agent familiarity with native OpenSpec mechanics, and reduces the risk of magical wrapper failure.
