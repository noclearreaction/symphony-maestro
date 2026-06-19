## Stage A - Minimal Non-Mutating Runtime Behavior

## 1. Request Path Foundation

- [ ] 1.1 Add a request pre-processing stage in rubato that inspects `messages[0].content` for the runtime injection anchor block
- [ ] 1.2 Implement strict anchor parsing with deterministic errors for malformed anchor content
- [ ] 1.2.1 Extend anchor parsing to support per-plugin static argument maps with schema validation
- [ ] 1.3 Add request eligibility checks so non-anchor requests bypass injection and forward unchanged

## Stage B - MVP Injection Using Plugin Contract

## 2. Plugin Contract And Execution

- [ ] 2.1 Define an internal plugin contract (key, execute, normalize output, error shape) that supports multiple plugins
- [ ] 2.2 Implement plugin registry resolution so only declared plugin keys are selected per request
- [ ] 2.2.1 Pass declared per-plugin args into plugin execution and reject invalid arg shapes with explicit failures
- [ ] 2.3 Implement fail-fast behavior for unknown plugin keys, command failures, and timeouts with clear user-facing error responses
- [ ] 2.4 Ensure plugin execution is refreshed per request with no per-session state reuse

## 3. Injection Mutation Semantics

- [ ] 3.1 Implement runtime-state block assembly from selected plugin outputs
- [ ] 3.2 Prepend the runtime-state block to `messages[-1]`
- [ ] 3.3 Inject plugin-presence usage guidance into `messages[0]` when absent, and keep this augmentation idempotent
- [ ] 3.3.1 Render `messages[0]` guidance using a canonical deterministic template (stable plugin ordering, stable arg ordering, no volatile fields)
- [ ] 3.3.2 Add a guidance format version token so future template changes are explicit and compatible with cache strategy
- [ ] 3.4 Guard mutation path for edge cases (missing messages array, missing user message) with explicit request failures

## 4. Git Status MVP Plugin

- [ ] 4.1 Implement `git_status` plugin to collect branch name, ahead/behind counts, commits-ahead count, staged count, unstaged tracked-modified count, and untracked count
- [ ] 4.2 Normalize git plugin output into bounded AI-consumable text without making exact phrasing part of behavior contract
- [ ] 4.3 Add git-specific failure mapping for non-repo and command execution errors

## 5. Configuration And Operational Integration

- [ ] 5.1 Add or update project `.opencode` configuration to route model traffic to rubato in devcontainer workflows
- [ ] 5.2 Add runtime configuration knobs for repo path and plugin command timeout bounds
- [ ] 5.3 Ensure logs capture injection decision path (anchor detected, plugins requested, plugin failures) without leaking unnecessary prompt content at non-debug levels

## 6. Verification

- [ ] 6.1 Add tests for anchor parsing and declared-plugin selection behavior
- [ ] 6.2 Add tests verifying `messages[-1]` receives injected state and `messages[0]` guidance injection is idempotent for eligible requests
- [ ] 6.2.1 Add tests asserting byte-identical `messages[0]` guidance across independent sessions for equivalent anchors
- [ ] 6.2.2 Add tests asserting guidance excludes volatile fields (timestamps, random IDs, host paths, counters)
- [ ] 6.3 Add tests verifying fail-fast responses for unknown plugin keys and plugin execution failures
- [ ] 6.4 Add tests verifying per-request refresh behavior by changing repository state between consecutive requests
- [ ] 6.5 Run end-to-end verification with opencode routed through rubato inside devcontainer

## Task 7 Sequence Rule

- [ ] Complete all Stage A and Stage B items before starting Stage C refinement tasks.
- [ ] Complete Task 7 contract freeze before Stage C refinement tasks.

## Stage C - Refinement And Polish

## 7. POC-Guided API Contract Freeze

- [ ] 7.1 Refine the POC request/response path to exercise uncertain API behaviors (anchor grammar, plugin error shape, and status mapping)
- [ ] 7.2 Capture concrete interaction traces from the POC and select the canonical wire contract for anchor parsing and failure responses
- [ ] 7.3 Update spec and design artifacts to freeze only the validated API interaction details and keep unresolved items implementation-defined
- [ ] 7.4 Re-run end-to-end verification against the frozen contract and confirm no behavioral regressions
