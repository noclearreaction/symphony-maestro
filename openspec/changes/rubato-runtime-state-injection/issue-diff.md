## Issue Comparison: rubato-runtime-state-injection vs GitHub issues

Date: 2026-06-19

## Sources Reviewed

- #61 rubato: marker-based runtime state injection
- #60 Introduce rubato: docker-compose scaffolding and Taskfile wrapper
- #62 rubato: structured log levels and single-file logging
- #64 rubato: subcommand validation rejects flags placed before subcommand
- #55 Spike SF-4a: Transparent Go proxy opencode -> OpenRouter
- #56 Spike SF-4b: Proxy traffic inspection and structured logging
- #24 Inject runtime environment context into agent prompts
- #43 Spike: empirically verify opencode cache behavior

---

## High-Impact Discrepancies

### 1) Unknown/invalid plugin handling: warn-and-continue in #61 vs fail-fast in this change

- Issue #61 scope says unknown inject keys should produce WARN and empty values.
- This change requires request failure for unknown plugin keys, execution failures, timeouts, and invalid plugin output.

Why this matters:
- Behavior changes from permissive degradation to strict correctness guarantees.
- Existing expectations from #61 would be incompatible unless issue text is updated.

### 2) System prompt mutation policy: #61 says never modify messages[0], this change allows idempotent guidance injection into messages[0]

- Issue #61 explicitly keeps `messages[0]` untouched to preserve cache warmth.
- This change deliberately augments `messages[0]` with deterministic plugin guidance (idempotent, byte-stable, no volatile fields).

Why this matters:
- This is a core contract difference, not an implementation detail.
- Caching strategy in #61 and this change are based on different assumptions.

### 3) Runtime architecture dependency: #61 depends on #60 compose stack, this change is topology-agnostic

- #61 depends on #60 (docker-compose + Taskfile + mounted `/workspace` workflow).
- This change explicitly avoids coupling MVP behavior to a fixed multi-container topology.

Why this matters:
- This change can be implemented without adopting #60 first.
- If #60 remains open/abandoned, this change can still proceed.

### 4) Plugin scope contract: #61 hardcodes key->command map and includes unittests; this change ships git_status MVP with extensible plugin contract and arg maps

- #61 design uses hardcoded keys (`git_status`, `unittests`) and command bindings.
- This change specifies `git_status` MVP and plugin extensibility/arguments as first-class contract requirements.

Why this matters:
- `unittests` is not part of MVP in this change.
- Implementers following #61 literally may overbuild beyond this approved scope.

---

## Medium Discrepancies

### 5) Logging approach differs from #56 and partially overlaps #62

- #56 expects full request/response payload inspection logs (spike observability).
- This change requires decision-path logging without unnecessary prompt leakage at non-debug levels.
- #62 proposes structured single-file logging and level controls; this change does not mandate file topology/schema, only logging outcomes.

Net:
- #56 behavior is intentionally superseded for production behavior.
- #62 is directionally compatible, but not required verbatim by this change.

### 6) Injection trigger format diverges from #61 marker text

- #61 uses a specific marker/comment style with Finding-prefixed block text.
- This change requires a valid rubato anchor block with plugin declarations and static per-plugin arg maps, but does not lock to #61 marker literal format.

Net:
- Existing marker syntax from #61 may need migration or parser compatibility handling.

### 7) Configuration surface: #24 focuses on .opencode command-template shell injection, this change shifts injection into rubato runtime path

- #24 centers OpenCode command markdown/template-level shell output injection.
- This change moves context injection responsibility to the proxy/runtime layer.

Net:
- The problem statement aligns (reduce redundant discovery/tool calls), but mechanism is fundamentally different.

---

## Low/Adjacent Gaps

### 8) #64 CLI subcommand validation bug is not covered by this change scope

- #64 targets entrypoint/argument parsing around `bin/rubato` subcommand detection.
- This change focuses request mutation/plugin semantics, not CLI parser behavior.

Net:
- No direct conflict, but #64 remains a separate implementation concern.

### 9) #43 spike-oriented measurement deliverables are only indirectly reflected

- #43 includes broad empirical cache-behavior investigation across many scenarios.
- This change bakes in deterministic guidance and per-request plugin refresh assumptions, but does not itself deliver #43 experiment matrix outputs.

Net:
- This change uses conclusions from spike work but does not replace all #43 deliverables.

---

## Alignment Summary (Where They Match)

- Strong alignment with #61 on anchor/marker-gated runtime injection into `messages[-1]`.
- Alignment with #55 on positioning proxy/rubato in `POST /v1/chat/completions` request path.
- Alignment with #24/#43 motivation: reduce redundant environment-discovery turns while preserving cache behavior.
- Partial alignment with #62 on safer logging defaults and structured observability intent.

---

## Recommended Follow-up (Issue Hygiene)

1. Update #61 to reflect fail-fast behavior and `messages[0]` deterministic guidance augmentation.
2. Mark #56 expectations as spike-only and superseded by production logging constraints.
3. Decide whether #60 is a hard prerequisite or optional deployment track for this change.
4. Open a follow-on issue if `unittests` plugin remains desired beyond `git_status` MVP.
5. Keep #64 separate; do not fold CLI parsing into this change unless scope is intentionally expanded.
