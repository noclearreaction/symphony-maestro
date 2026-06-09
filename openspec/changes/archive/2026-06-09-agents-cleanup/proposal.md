## Why

`AGENTS.md` currently contains a mixture of immediate operational context and speculative, future-oriented system specifications (such as Symphony v2 conceptual requirements). This mixing of active system status with future speculative targets violates the Principle of Progress, the Symphony Constitution (Principle 3: Intent Traceability), and workspace hygiene guidelines. To ensure agents are guided by clear, unambiguous, and active context, `AGENTS.md` must be trimmed of future speculative targets, and those targets must be tracked systematically as individual GitHub Issues.

## What Changes

- **Section-by-Section Review of `AGENTS.md`**: Inspect all sections to isolate speculative, future-looking targets from active status and operational context.
- **Speculative Target Extraction**: Extract future-oriented spec/feature targets (specifically around the constraints and architecture of Symphony v2) into individual GitHub Issues to track them systematically.
- **`AGENTS.md` Trimming**: Trim `AGENTS.md` to serve strictly as an active, high-fidelity marker of system status, operational principles, and active system layers, aligning it fully with the Symphony Constitution and Vision.
- **Reference Updates**: Ensure that any other parts of the workspace pointing to `AGENTS.md` remain valid and aligned with the updated file structure.

## Capabilities

### New Capabilities
- None

### Modified Capabilities
- controlled-git-workflows: Ensure the AGENTS.md instruction file serves strictly as a high-fidelity operational context marker without speculative system requirements.

## Impact

- **`AGENTS.md`**: Streamlined and trimmed down by removing speculative v2 elements, serving strictly as a marker of active system status and operational context.
- **GitHub Issues**: New GitHub issues generated to track extracted speculative targets under the correct categorization schema (`type:feature`, `status:backlog`, `priority:medium/low`).
- **Orchestration / Developer Agents**: Will load a more focused, high-hygiene context file without being confused by future speculative architecture.
