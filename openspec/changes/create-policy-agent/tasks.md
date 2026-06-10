## 1. Agent Definition Bootstrap

- [ ] 1.1 Create the specialized policy agent file at `.opencode/agents/policy.md` (Note: This must be created by the human operator or bootstrapped directly, since `@builder` is blocked from editing under `.opencode/*`).
- [ ] 1.2 Implement the dynamic agent frontmatter in `.opencode/agents/policy.md` with:
  - `mode: primary`
  - `edit` permissions allowing `.opencode/*` and `AGENTS.md` and denying everything else (`"*": deny`).
  - `read` permissions allowing all files in the repository (`"*": allow`, `.opencode/*`: allow).
  - `bash` permissions denying general command execution, but allowing `openspec validate` and `openspec status` commands (`"*": deny`, `openspec validate *`: allow, `openspec status *`: allow).

## 2. Policy Agent System Instructions

- [ ] 2.1 Author the `@policy` system prompt body in `.opencode/agents/policy.md`. State its core purpose as managing prompt instructions, agent configs, tool permissions, and global coordination docs.
- [ ] 2.2 Configure prompt boundaries: instruct `@policy` to refuse editing strategy documents under `strategy/` (which belong to `@advisor`), refuse editing application runtime code or test suites (which belong to `@builder`), and delegate all git operations to `@git-operator`.

## 3. Global Documentation Update

- [ ] 3.1 Update the main `AGENTS.md` file to introduce the `@policy` agent.
- [ ] 3.2 Document `@policy`'s specific boundaries and collaboration contracts with existing agents (`@advisor`, `@designer`, `@builder`, `@issue`, `@git-operator`, `@orchestrator`) inside `AGENTS.md`.

## 4. Verification and Validation

- [ ] 4.1 Run `openspec validate` to ensure the new specs and changes are fully compliant with the repository schema.
- [ ] 4.2 Verify the `@policy` agent registration is active by querying its status.
