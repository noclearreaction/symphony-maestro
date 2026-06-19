## 1. Anchor Parsing And Eligibility

- [ ] 1.1 Implement request pre-processing that inspects `messages[0].content` for runtime anchor blocks.
- [ ] 1.2 Implement strict anchor parsing with deterministic malformed-anchor failures.
- [ ] 1.3 Implement per-plugin static args parsing with schema validation.
- [ ] 1.4 Ensure non-anchor requests bypass injection and forward unchanged.

## 2. Plugin Contract And Execution

- [ ] 2.1 Implement internal plugin contract and registry for multi-plugin extensibility.
- [ ] 2.2 Execute only declared plugin keys.
- [ ] 2.3 Fail fast for unknown plugin keys, plugin command errors, and timeouts.
- [ ] 2.4 Ensure plugin execution refreshes per request (no session reuse).

## 3. Mutation Semantics

- [ ] 3.1 Build runtime-state block from selected plugin outputs.
- [ ] 3.2 Prepend runtime-state block to `messages[-1]`.
- [ ] 3.3 Inject plugin guidance in `messages[0]` when absent and keep it idempotent.
- [ ] 3.4 Implement canonical deterministic guidance template with explicit version token.
- [ ] 3.5 Fail request explicitly for missing/invalid message structures needed for mutation.

## 4. Git Status MVP Plugin

- [ ] 4.1 Implement `git_status` plugin metrics: branch-or-head state, ahead/behind, commits-ahead, staged, unstaged tracked-modified, untracked.
- [ ] 4.2 Ensure detached-HEAD and bare-repo states are represented as visible plugin output states.
- [ ] 4.3 Map non-repo/unreadable/git-exec failures to explicit git_status failure responses.

## 5. Verification

- [ ] 5.1 Add tests for anchor parsing, arg validation, and declared-plugin selection.
- [ ] 5.2 Add tests for `messages[-1]` runtime-state injection and `messages[0]` idempotent guidance injection.
- [ ] 5.3 Add tests for byte-identical guidance across sessions for equivalent anchors.
- [ ] 5.4 Add tests for fail-fast error paths and per-request refresh behavior.