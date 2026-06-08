## 1. Commit Message Linter Implementation

- [x] 1.1 Implement a zero-dependency Deno TypeScript commit linter `bin/commit-lint.ts` to parse and validate commit messages against Conventional Commits.
- [x] 1.2 Validate `bin/commit-lint.ts` against both compliant and non-compliant test commit messages.

## 2. Submit Script Integration

- [x] 2.1 Modify `bin/director-submit.ts` to run `bin/commit-lint.ts` and validate that any automatic commit messages conform to the Conventional Commit standard.

## 3. Git Hook Configuration

- [x] 3.1 Create a `.git/hooks/commit-msg` git hook file that automatically invokes `bin/commit-lint.ts` when a commit is made.
- [x] 3.2 Make the git hook executable.

## 4. AI Agent Guidelines Integration

- [ ] 4.1 Update system instructions for the Orchestrator Agent (`.opencode/agents/orchestrator.md`) to require Conventional Commits, single-topic branches, and atomic unit-of-work commits.
- [ ] 4.2 Update system instructions for the Builder Agent (`.opencode/agents/builder.md`) to require Conventional Commits, single-topic branches, and atomic unit-of-work commits.
- [ ] 4.3 Update system instructions for the Designer Agent (`.opencode/agents/designer.md`) to require Conventional Commits, single-topic branches, and atomic unit-of-work commits.
- [ ] 4.4 Update system instructions for the Issue Agent (`.opencode/agents/issue.md`) to recognize development workflows and conventional commits.
