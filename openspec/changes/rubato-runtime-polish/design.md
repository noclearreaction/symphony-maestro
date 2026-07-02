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

## Open Questions

### OQ-1) State injection frequency — always vs. on-change

Currently rubato injects a fresh `rubato:state` block on every request turn.
Since the full message history is present in each request body, rubato could
instead scan backwards through prior messages to find the most recent output
per plugin and skip injection when output is unchanged.

Decision needed:
- **Always inject** (current): simple, correct, ~60-80 tokens per plugin per
  turn, negligible for 5-10 small plugins against typical 32k+ contexts.
- **Inject on change**: scan `messages[1..-2]` for previous `rubato:state`
  blocks, compare fresh output per plugin, only prepend if any differ. ~20-30
  lines in the `mutate` package. Quieter in stable long conversations.

Sub-question: if only *some* plugins change, inject only changed plugins or
always inject the full set?

This decision should be made before closing rubato-runtime-polish. If
on-change is chosen, the implementation belongs in `mutate.Apply` and requires
a new test covering the suppression case.

- Refinement-only guardrails can defer legitimate improvements; accepted to preserve staged integrity.
- More verification work increases short-term cycle time; accepted to reduce regressions.
