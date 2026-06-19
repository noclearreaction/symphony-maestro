## Review: rubato-runtime-state-injection

Date: 2026-06-19
Reviewer: GitHub Copilot (GPT-5.3-Codex)

## Findings

Status key: `resolved` means the discrepancy is addressed in current artifacts.

### 1. High - Malformed-anchor behavior was contradictory/underspecified (`resolved`)

The spec states that requests without a valid anchor are forwarded unchanged, while tasks and design require strict parsing with deterministic parse failures.

Current texts leave malformed-but-present anchor behavior ambiguous:
- pass-through as if no anchor is present
- fail-fast with explicit parse error

References:
- openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-state-injection/spec.md:14
- openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-state-injection/spec.md:16
- openspec/changes/rubato-runtime-state-injection/tasks.md:4
- openspec/changes/rubato-runtime-state-injection/design.md:180

Risk:
- Implementers and tests may encode conflicting behavior.

Resolution:
- Added malformed-anchor fail-fast scenario to spec and aligned task/design intent.

### 2. Medium - Statelessness intent conflicted with session-oriented wording (`resolved`)

The change claims stateless per-request behavior, but guidance injection is framed as once per conversation/session and sequence flow branches on session presence.

References:
- openspec/changes/rubato-runtime-state-injection/proposal.md:12
- openspec/changes/rubato-runtime-state-injection/design.md:13
- openspec/changes/rubato-runtime-state-injection/design.md:16
- openspec/changes/rubato-runtime-state-injection/design.md:132

Risk:
- Implementation may accidentally introduce hidden per-session state.

Resolution:
- Updated wording to request-content-driven idempotence and removed session-memory implication from core flow text.

### 3. Medium - git hygiene metric "committed count" was ambiguous (`resolved`)

The MVP language previously used "committed count" without exact derivation, which allowed multiple incompatible interpretations.

References:
- openspec/changes/rubato-runtime-state-injection/proposal.md:10
- openspec/changes/rubato-runtime-state-injection/specs/rubato-runtime-state-injection/spec.md:68
- openspec/changes/rubato-runtime-state-injection/tasks.md:27

Risk:
- Inconsistent plugin output and brittle tests.

Resolution:
- Replaced ambiguous wording with explicit metric names (commits-ahead, staged, unstaged tracked-modified, untracked).

## Open Questions

1. ~~For guidance idempotence, should equality be byte-exact block match or semantic equivalence of declared plugins and args?~~ **Resolved: byte-identical.** The cache-stability requirement drives this — semantically equivalent but non-identical text would invalidate prefix caches. Artifacts already reflect this.

2. ~~Should detached-HEAD and bare-repo behavior be included in the MVP contract now or deferred?~~ **Resolved: in MVP contract as observable output.** Detached-HEAD and bare-repo are not error conditions. The git_status plugin SHALL report them as visible state in the output so the model can reason about repository context. Spec updated with explicit scenarios.

## Overall Assessment

Direction is strong, especially on fail-fast behavior, staged sequencing, and deterministic guidance. Remaining questions are now narrower and implementation-focused.
