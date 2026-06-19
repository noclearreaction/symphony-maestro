## Context

This change corresponds to Stage C from the prior staged plan and intentionally excludes new product scope. Its purpose is hygiene, operability, and contract confidence.

## Goals / Non-Goals

**Goals**
- Improve production-safe logging and diagnostics.
- Harden runtime config boundaries and defaults.
- Freeze contract details from implementation evidence.
- Verify behavior end-to-end in devcontainer routing path.

**Non-Goals**
- No new plugins.
- No new anchor semantics.
- No broad architecture rewrite.

## Decisions

### 1) Stage C is refinement-only

Stage C cannot redefine Stage A/B behavior unless explicitly reopened by a new approved change.

Rationale:
- Prevents hidden scope drift while polishing.

### 2) Logging favors decision-path observability over payload dumping

Operational logs capture anchor detection, declared plugins, mutation decision, and failure causes while avoiding unnecessary prompt content leakage at non-debug levels.

Rationale:
- Balances debugging utility with safety and hygiene.

### 3) Contract freeze is evidence-backed

Task 7 outcomes must include concrete request/response traces and artifact updates that reflect validated behavior, not speculative wording.

Rationale:
- Keeps spec language grounded in implemented behavior.

### 4) End-to-end verification is mandatory

Polish is complete only when routed workflow tests pass in the actual devcontainer environment used by operators.

Rationale:
- Prevents local-only confidence that fails in real workflow topology.

## Risks / Trade-offs

- Refinement-only guardrails can defer legitimate improvements; accepted to preserve staged integrity.
- More verification work increases short-term cycle time; accepted to reduce regressions.
