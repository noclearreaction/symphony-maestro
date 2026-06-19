# Issue Disposition Matrix

Purpose: align active issues to the split Rubato plan.

## Decisions Applied

1. Rubato work is split into three changes:
   - `rubato-proxy-foundation`
   - `rubato-plugin-git-status`
   - `rubato-runtime-polish`
2. Fail-fast behavior is required.
3. Rubato is stateless.
4. Issues are updated to match the plan.

## Matrix

| Issue | Current State | Mismatch | Decision | Action Now |
|---|---|---|---|---|
| #24 | CLOSED | Broad prompt-injection feature with cache concerns overlapped Rubato runtime injection scope | Superseded by split Rubato changes and implementation issue #61 | Closed with superseded rationale and residual-scope guidance |
| #43 | OPEN | Spike parent issue, still useful for observability findings | Keep open; update with stage linkage and what Rubato will consume | Add comment with current Rubato alignment |
| #55 | CLOSED | Spike SF-4a proxy groundwork fed Rubato implementation work | Completed/superseded by split Rubato track | Closed with resolution comment |
| #56 | CLOSED | None | No change | None |
| #60 | OPEN | Assumes compose scaffolding is a strict prerequisite; staged plan is more flexible | Update wording to supporting infrastructure, not hard gate for all stages | Add alignment comment |
| #61 | OPEN | Historical monolithic references caused execution ambiguity | Keep as umbrella implementation issue pointing to split changes in sequence | Maintain split-change links and completion criteria |
| #62 | OPEN | Logging issue depends on older sequencing | Keep open; move into Stage C refinement | Add sequencing comment |
| #64 | OPEN | Independent bug in subcommand parsing | Keep open; out of Rubato staging scope | Add note: tracked separately |

## Execution Order

1. Keep #61 as the umbrella issue for split change sequencing.
2. Keep #60 and #62 aligned as supporting/stage-C issues.
3. Keep #43 as spike parent context and #64 as independent out-of-scope bug.
4. Keep #24 and #55 closed as superseded/completed.

## Done Criteria

1. Each issue above has one explicit status direction.
2. #61 text clearly reflects split changes and locked behavior decisions.
3. No issue claims contradict the current change artifacts.