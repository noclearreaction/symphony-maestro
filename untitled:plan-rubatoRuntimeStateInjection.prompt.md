## Rubato Hygiene Plan

Goal: align issues and change docs so they say the same thing before implementation work proceeds.

### Fixed Decisions
1. Work is split into three stages.
2. Error handling is fail-fast.
3. The committed-count definition must be explicit.
4. Rubato is stateless.
5. Complete Task 7 after Stage B and before Stage C.
6. Update issues to match the approved plan.

### Stages
1. Stage A: minimal non-mutating Rubato behavior.
2. Stage B: MVP injection using one plugin.
3. Stage C: polish and cleanup.

### Task 7 Sequence
1. Finish Stage A.
2. Finish Stage B.
3. Complete Task 7 and lock contract details.
4. Start Stage C.
5. Do not change contract details during Stage C unless explicitly reopened.

### Execution Steps
1. Build issue disposition matrix.
Columns: issue, mismatch, decision, doc updates, close condition.
2. Reconcile active change docs.
Fix naming drift, malformed-anchor behavior, committed-count meaning, and stateless wording.
3. Fill contract gaps.
Define error response shape, plugin argument validation, deterministic rendering rules, and request edge-case behavior.
4. Update issue statuses.
Close, reject, update, or split only where real residual scope exists.
5. Check sibling changes.
Only update or close sibling changes when they directly affect Rubato assumptions.
6. Run final consistency sweep.
Confirm proposal, design, tasks, spec, review, and issue diffs all align.

### Done Criteria
1. Every known discrepancy has one documented resolution.
2. Every affected issue has one clear terminal action.
3. No MVP-critical terms are ambiguous.
4. No unresolved contradictions remain across change artifacts.

### Primary Files
- /workspace/openspec/changes/rubato-runtime-state-injection/proposal.md
- /workspace/openspec/changes/rubato-runtime-state-injection/design.md
- /workspace/openspec/changes/rubato-runtime-state-injection/tasks.md
- /workspace/openspec/changes/rubato-runtime-state-injection/review.md
- /workspace/openspec/changes/rubato-runtime-state-injection/issue-diff.md
- /workspace/openspec/changes/rubato-runtime-state-injection/change-diff.md
- /workspace/openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-state-injection/spec.md