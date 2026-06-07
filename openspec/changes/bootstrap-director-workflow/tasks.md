## 1. Custom Agent Registration

- [x] 1.1 Create the custom builder agent system prompt file under `.opencode/agents/builder.md` with strict soft boundaries.
- [x] 1.2 Ensure the advisor and designer agent files remain unmodified, and the builder agent registers successfully.

## 2. Local Workspace Isolation

- [x] 2.1 Update the repository root `.gitignore` to ignore the `.symphony/` directory and its contents.
- [x] 2.2 Create the structured directories `.symphony/scratchpad/`, `.symphony/plans/`, `.symphony/reviews/`, and `.symphony/evidence/`.

## 3. Implement Start Script (bin/director-start.ts)

- [x] 3.1 Create the script `bin/director-start.ts` with Deno TypeScript.
- [x] 3.2 Implement argument validation, branch checkout logic (`change/<change-name>`), and `openspec new change "<change-name>"` invocation.
- [x] 3.3 Ensure the script is executable (`chmod +x bin/director-start.ts`).

## 4. Implement Submit Script (bin/director-submit.ts)

- [x] 4.1 Create the script `bin/director-submit.ts` with Deno TypeScript.
- [x] 4.2 Implement git branch pre-flight validation (ensuring it is run on a `change/*` branch).
- [x] 4.3 Add execution logic for `openspec sync` to promote delta specs to the main `specs/` directory, automatically stage/commit the synced changes, push to remote, and execute `gh pr create` (with defensive fallback if `gh` CLI is unauthenticated or remote is upstream).
- [x] 4.4 Ensure the script is executable (`chmod +x bin/director-submit.ts`).

## 5. Verification & Testing

- [x] 5.1 Execute a local dry-run or error-case test of `bin/director-start` and `bin/director-submit` to verify robust error-handling.
- [x] 5.2 Run `openspec validate` on the active changes to verify absolute schema compliance.
