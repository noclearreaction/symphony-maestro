# Issue Disposition Matrix

Purpose: align active issues to the staged Rubato plan.

## Decisions Applied

1. Rubato work is staged:
   - Stage A: minimal non-mutating runtime behavior
   - Stage B: MVP plugin-based injection
   - Stage C: refinement and polish
2. Fail-fast behavior is required.
3. Rubato is stateless.
4. Issues are updated to match the plan.

## Matrix

| Issue | Current State | Mismatch | Decision | Action Now |
|---|---|---|---|---|
| #24 | OPEN | Broad prompt-injection feature with cache concerns; only partially overlaps Rubato | Update issue to reference Rubato staging and keep broader scope open | Add cross-link comment to staged Rubato change |
| #43 | OPEN | Spike parent issue, still useful for observability findings | Keep open; update with stage linkage and what Rubato will consume | Add comment with current Rubato alignment |
| #55 | OPEN | Spike SF-4a proxy groundwork appears complete and now feeds staged Rubato work | Close as completed/superseded by Rubato staged track | Close with resolution comment |
| #56 | CLOSED | None | No change | None |
| #60 | OPEN | Assumes compose scaffolding is a strict prerequisite; staged plan is more flexible | Update wording to supporting infrastructure, not hard gate for all stages | Add alignment comment |
| #61 | OPEN | Current text assumes single-shot marker mutation flow and old assumptions | Update as primary implementation issue for staged plan (A/B/C), fail-fast, stateless | Add authoritative alignment comment and checklist |
| #62 | OPEN | Logging issue depends on older sequencing | Keep open; move into Stage C refinement | Add sequencing comment |
| #64 | OPEN | Independent bug in subcommand parsing | Keep open; out of Rubato staging scope | Add note: tracked separately |

## Execution Order

1. Update #61 first as primary source issue for staged implementation.
2. Update #60 and #62 to align dependencies and sequencing.
3. Update #24 and #43 with cross-links and scope boundaries.
4. Close #55 with superseded/completed rationale.
5. Confirm #64 remains independent and unchanged in scope.

## Done Criteria

1. Each issue above has one explicit status direction.
2. #61 text clearly reflects staged A/B/C plan, fail-fast, and stateless behavior.
3. No issue claims contradict the current change artifacts.