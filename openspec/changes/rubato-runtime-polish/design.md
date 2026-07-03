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

### 5) On-change injection — Option A with anchor-configured window

Rubato injects only the plugins whose output has changed since they were last seen in message history, scanned backward up to `repeat` messages (default 100, configurable in the `rubato:anchor` block).

Per-plugin logic:
- Found in window, output unchanged → skip
- Found in window, output changed → inject
- Not found in window (first turn, or beyond `repeat` limit) → inject (reminder)

When any plugins require injection, only those plugins appear in the `rubato:state` block. When no plugins changed and all are within the window, no state block is prepended.

The `max_age` field is read from a `parameters` array in the anchor JSON:
```json
{"plugins":["git_status","go_test"],"parameters":[{"max_age":50}]}
```
The `Block` struct gains a `Parameters` field parsed from the top-level `parameters` array. `MaxAge` is extracted from the first object in that array; absent defaults to 100. `max_age: 0` means always inject regardless of history. The `parameters` key is reserved for rubato-level config, distinct from per-plugin `args`.

Guidance text does NOT explain absence semantics — any explicit hint risks triggering the model to substitute tool calls for ambient state it already has.

Rationale:
- Keeps injected blocks minimal and signal-rich (presence = something changed)
- Bounded scan (O(repeat)) prevents O(n) cost on long conversations
- Forced reminder after `repeat` turns handles context compression in long sessions
- Per-plugin atomicity lets each plugin signal independently

## Risks / Trade-offs
