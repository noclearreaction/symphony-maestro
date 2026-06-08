## 1. Custom Scripts Deletion

- [x] 1.1 Delete the custom start script `bin/director-start.ts`.
- [x] 1.2 Delete the custom submit script `bin/director-submit.ts`.

## 2. AI Agent Instructions Realignment

- [x] 2.1 Update the Orchestrator Agent (`.opencode/agents/orchestrator.md`) system prompt to remove any reference or instructions related to `bin/director-start.ts` or `bin/director-submit.ts`.
- [x] 2.2 Update the Designer Agent (`.opencode/agents/designer.md`) system prompt to align instructions and transition protocol with standard native `openspec` commands.
- [x] 2.3 Update the Builder Agent (`.opencode/agents/builder.md`) system prompt to remove any reference to custom scripts.
- [x] 2.4 Verify all other agents and instructions do not contain residual references to the custom scripts.
