## 1. Label Provisioning

- [x] 1.1 Create label provisioning Deno script `bin/provision-labels.ts` using Deno's Subprocess API to run `gh`
- [x] 1.2 Run `bin/provision-labels.ts` with restricted sandbox flags (`deno run --allow-run=gh bin/provision-labels.ts`) to create the 12 specified labels in the repository

## 2. Issue Subagent Integration

- [x] 2.1 Update system instructions in `.opencode/agents/issue.md` to define standard label taxonomy, require recommending type and priority labels during drafting, and automatically apply `status:backlog` during creation

## 3. Governance and Documentation

- [x] 3.1 Draft the development lifecycle transition guidelines in `governance/issue-lifecycle.md` explaining state machine transitions
