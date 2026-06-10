## Why

The Symphony Director utilizes multiple specialized AI agents with strict boundaries to perform engineering and planning workflows. Currently, modifying agent prompts, configurations, and global agent documentation (such as `AGENTS.md`) is done manually or using general agents, which can lead to accidental modifications of source code or policy inconsistencies. Introducing a dedicated `@policy` worker agent will isolate the responsibility of managing AI prompts, configurations, and policies under `.opencode/` and `AGENTS.md`, while enforcing a strict firewall that prevents this agent from editing application runtime code or executing unauthorized bash commands.

## What Changes

- Register a new specialized worker agent named `@policy` in `.opencode/agents/policy.md`.
- Configure the `@policy` agent with strict read, edit, and execution boundaries:
  - **Read**: Allowed to read all files in the repository (excluding potential sensitive credentials, with `.opencode/*` read enabled specifically for this agent).
  - **Edit**: Allowed to edit files under `.opencode/*` (prompts, configs, scripts) and `AGENTS.md`. Denied from editing any other source code, test suites, or configuration files.
  - **Bash/Command Execution**: Denied general execution, but allowed to run standard safe validations such as `openspec validate`. Direct git actions are delegated to `@git-operator`.
- Document the `@policy` agent's purpose, boundaries, and collaboration contracts with other agents in the repository's main `AGENTS.md` file.

## Capabilities

### New Capabilities
- `policy-management`: Defines the roles, permissions, and validation rules for the `@policy` agent dedicated to managing AI prompts, configurations, and global agent documentation.

### Modified Capabilities

## Impact

- **AI Agents Configuration**: A new agent definition file at `.opencode/agents/policy.md`.
- **System Documentation**: Modification to `AGENTS.md` to document the `@policy` agent and establish its boundaries with other agents.
- **Workflow Security**: Strengthened execution firewall by ensuring that prompt engineering and agent policy management can only be performed by an agent with zero application-edit permissions, preventing security/policy drift during automation tasks.
