## Review: rubato-runtime-state-injection

Date: 2026-06-19
Reviewer: GitHub Copilot (GPT-5.3-Codex)

## Findings

### 1. High - Malformed-anchor behavior is contradictory/underspecified

The spec states that requests without a valid anchor are forwarded unchanged, while tasks and design require strict parsing with deterministic parse failures.

Current texts leave malformed-but-present anchor behavior ambiguous:
- pass-through as if no anchor is present
- fail-fast with explicit parse error

References:
- openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-injection/spec.md:14
- openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-injection/spec.md:16
- openspec/changes/rubato-runtime-state-injection/tasks.md:4
- openspec/changes/rubato-runtime-state-injection/design.md:180

Risk:
- Implementers and tests may encode conflicting behavior.

Recommendation:
- Add a dedicated scenario for malformed anchor input with explicit response semantics and status/error shape.

### 2. Medium - Statelessness intent conflicts with session-oriented wording

The change claims stateless per-request behavior, but guidance injection is framed as once per conversation/session and sequence flow branches on session presence.

References:
- openspec/changes/rubato-runtime-state-injection/proposal.md:12
- openspec/changes/rubato-runtime-state-injection/design.md:13
- openspec/changes/rubato-runtime-state-injection/design.md:16
- openspec/changes/rubato-runtime-state-injection/design.md:132

Risk:
- Implementation may accidentally introduce hidden per-session state.

Recommendation:
- Specify that idempotence is derived solely from request content (for example, marker detection in messages[0]) and requires no server-side session memory.

### 3. Medium - git hygiene metric "committed count" is ambiguous

The MVP requires a committed count but does not define exact derivation. This can be interpreted in multiple, incompatible ways.

References:
- openspec/changes/rubato-runtime-state-injection/proposal.md:10
- openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-injection/spec.md:68
- openspec/changes/rubato-runtime-state-injection/tasks.md:27

Risk:
- Inconsistent plugin output and brittle tests.

Recommendation:
- Replace with unambiguous metric names and explicit derivation rules.

## Open Questions

1. Should malformed anchors be hard failures or no-anchor pass-through?
2. For guidance idempotence, should equality be byte-exact block match or semantic equivalence of declared plugins and args?
3. Should detached-HEAD and bare-repo behavior be included in the MVP contract now or deferred?

## Overall Assessment

Direction is strong, especially on fail-fast plugin behavior and deterministic guidance. Clarifying the three items above before implementation will reduce drift and test churn.
