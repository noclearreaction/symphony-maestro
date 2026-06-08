## Why

The custom wrapper scripts `bin/director-start.ts` and `bin/director-submit.ts` introduce unnecessary maintenance overhead, obscure standard OpenSpec CLI behaviors, and diverge from community best practices. Removing these wrappers and instructing agents to use standard native OpenSpec CLI commands directly simplifies the codebase and aligns perfectly with Constitutional Principle 5 (Tooling Discipline).

## What Changes

- **Script Removal**: Delete `bin/director-start.ts` and `bin/director-submit.ts` from the codebase completely.
- **Specification Deprecation**: Update the `director-workflow` specification to remove the custom automation requirements and specify standard native OpenSpec CLI and `git-operator` subagent workflows instead.
- **Agent Instruction Cleanup**: Remove all references to the custom scripts from the Orchestrator, Designer, Builder, and Issue agent instruction files.

## Capabilities

### New Capabilities
<!-- None, this is a modification of existing specifications -->

### Modified Capabilities
- `director-workflow`: Deprecates the custom automation scripts (`director-start` and `director-submit`) and replaces them with requirements for standard OpenSpec CLI toolchain workflows.

## Impact

- **Codebase Simplification**: Deletes 2 Deno TypeScript scripts under `bin/`.
- **Engineering Workflows**: Humans and agents will use the native `openspec` toolchain directly.
