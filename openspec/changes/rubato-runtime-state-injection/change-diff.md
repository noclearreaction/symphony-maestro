## Change Comparison: rubato-runtime-state-injection vs prior proxy/rubato efforts

Date: 2026-06-19

## Scope Compared

This diff compares the current change (`rubato-runtime-state-injection`) to earlier, mostly-spike artifacts in this repo that also touched proxy/routing behavior:

1. `spike/issue-45-opencode-cache/proxy/main.go`
2. `spike/issue-45-opencode-cache/proxy/AGENTS.md`
3. `spike/issue-45-opencode-cache/findings/sf-2-observability.md`
4. `openspec/changes/scaffold-opencode-cache-harness/design.md` (related predecessor context)

Note: Among OpenSpec changes, this current change is the only one that explicitly defines a rubato/proxy runtime-injection capability. Prior proxy behavior mostly exists as spike prototype artifacts.

---

## Key Discrepancies

### 1) Purpose and scope: observability proxy vs runtime context injector

- Prior prototype (`spike/.../proxy/main.go`) is a thin forwarder + logger to OpenRouter (`const upstream = https://openrouter.ai/api/v1/chat/completions`) with no prompt mutation or plugin model.
- Current change defines a request mutation system (anchor parsing, plugin execution, mutation of `messages[-1]`, optional guidance mutation in `messages[0]`).

Discrepancy:
- Old: transport/logging utility for cache experiments.
- New: productized request semantics layer controlling model context.

References:
- `spike/issue-45-opencode-cache/proxy/main.go:16`
- `spike/issue-45-opencode-cache/proxy/main.go:110-200`
- `openspec/changes/rubato-runtime-state-injection/design.md:26-44`
- `openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-state-injection/spec.md:18-27`

### 2) Activation model: always-on forwarding vs anchor-gated behavior

- Prior proxy handles every `POST /v1/chat/completions` similarly.
- Current change runs injection only when a valid anchor exists in `messages[0]`; otherwise pass-through.

Discrepancy:
- Old: unconditional request processing.
- New: explicit in-band opt-in via anchor contract.

References:
- `spike/issue-45-opencode-cache/proxy/main.go:44-46`
- `openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-state-injection/spec.md:3-16`
- `openspec/changes/rubato-runtime-state-injection/tasks.md:3-6`

### 3) Error semantics: best-effort forwarding vs fail-fast plugin contract

- Prior proxy only fails for transport/upstream errors; there is no declared plugin failure model.
- Current change requires fail-fast on unknown plugin keys, plugin execution errors, timeouts, and invalid plugin output.

Discrepancy:
- Old: network-forwarding reliability errors only.
- New: semantic validation and plugin execution failures are first-class request failures.

References:
- `spike/issue-45-opencode-cache/proxy/main.go:156-165`
- `openspec/changes/rubato-runtime-state-injection/design.md:62-70`
- `openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-state-injection/spec.md:56-65`

### 4) Data handling and privacy posture: full payload logging vs bounded decision-path logging

- Prior proxy logs near-full request and response bodies into NDJSON per session key (including timestamps and body payloads).
- Current tasks require logging injection decisions and failures without leaking unnecessary prompt content at non-debug levels.

Discrepancy:
- Old: high-fidelity payload capture for experiment diagnostics.
- New: constrained operational logging with content-leak minimization.

References:
- `spike/issue-45-opencode-cache/proxy/AGENTS.md:22-27`
- `spike/issue-45-opencode-cache/proxy/main.go:184-189`
- `openspec/changes/rubato-runtime-state-injection/tasks.md:35`

### 5) Determinism target: timestamp/session-keyed logs vs canonical byte-stable guidance

- Prior proxy intentionally writes volatile fields (`timestamp`) and derives session keys from first 512 bytes of `messages[0]` hash for file grouping.
- Current change introduces deterministic canonical rendering for `messages[0]` guidance, with stable ordering and explicit prohibition on volatile fields.

Discrepancy:
- Old: observability artifacts optimized for tracing runs.
- New: cache-stability contract optimized for repeatable prompt bytes.

References:
- `spike/issue-45-opencode-cache/proxy/AGENTS.md:22-27`
- `spike/issue-45-opencode-cache/proxy/main.go:58-69`
- `spike/issue-45-opencode-cache/proxy/main.go:159`
- `openspec/changes/rubato-runtime-state-injection/design.md:49-60`
- `openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-state-injection/spec.md:33-39`

### 6) Architecture direction: single hardcoded upstream vs plugin-extensible execution

- Prior proxy hardcodes one upstream (`openrouter.ai`) and lacks extension points beyond forwarding.
- Current change defines plugin registry/contract semantics and explicitly plans beyond `git_status` MVP.

Discrepancy:
- Old: fixed transport endpoint with no plugin abstraction.
- New: extensible plugin architecture behind stable anchor semantics.

References:
- `spike/issue-45-opencode-cache/proxy/main.go:16`
- `openspec/changes/rubato-runtime-state-injection/design.md:82-90`
- `openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-state-injection/spec.md:78-83`
- `openspec/changes/rubato-runtime-state-injection/tasks.md:10-13`

### 7) Repository-state source: DB/session analytics focus vs git hygiene runtime signals

- Prior spike findings focus on opencode DB/session metrics (`tokens_*`, `cost`) and cache measurement methodology.
- Current change shifts runtime signal source to git-status hygiene (branch, ahead/behind, committed/staged/untracked) injected into model context.

Discrepancy:
- Old: session analytics and cache-observability experiment focus.
- New: repository-working-state awareness for request-time guidance.

References:
- `spike/issue-45-opencode-cache/findings/sf-2-observability.md:14-23`
- `spike/issue-45-opencode-cache/findings/sf-2-observability.md:54-57`
- `openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-state-injection/spec.md:67-76`
- `openspec/changes/rubato-runtime-state-injection/tasks.md:27-29`

### 8) Lifecycle intent: temporary spike harness vs durable runtime capability

- Prior harness design explicitly states temporary spike-local role and non-product intent.
- Current change explicitly defines a durable architecture independent of prototype code.

Discrepancy:
- Old: exploratory scaffold with temporary placement.
- New: durable capability intended for ongoing workflows.

References:
- `openspec/changes/scaffold-opencode-cache-harness/design.md:3-5`
- `openspec/changes/rubato-runtime-state-injection/design.md:3-5`

---

## Practical Interpretation

The current change is not a small iteration on the old proxy; it is a directional shift:

- from experiment instrumentation
- to deterministic, contract-driven runtime context injection

This means some spike-era behaviors should be treated as intentionally superseded (full payload logging, always-on forwarding semantics, and session-key file partitioning), while transport compatibility and minimal dependency posture can still be carried forward.
