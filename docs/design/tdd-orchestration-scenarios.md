# TDD Orchestration — Scenario Walkthroughs

> Companion to [`tdd-orchestration-intent.md`](tdd-orchestration-intent.md).
> Role-played, natural-language walkthroughs that illustrate the *intended*
> interaction texture. These are **illustrative, not normative** — a shared
> "contract" of shape, not a specification. `[rubato]` marks state pre-injected
> onto the prompt (no knowledge-gathering turn required).

## Scenario A — Happy path: a new small feature

**Context:** feature = "a config loader rejects an unknown top-level key."

1. **Human → Orchestrator:** "Start a change for: config loader rejects unknown
   top-level keys."
2. **Orchestrator → Discovery** (escalates — no test list exists yet).
3. **Discovery:** reads existing config code + domain notes, then produces:
   - `specs/config-loader.md` (DoD, Gherkin): *WHEN a config has an unknown
     top-level key THEN loading fails with an error naming the key.*
   - `test-list.md` (living): "known keys still load"; "one unknown key → error
     names it"; "multiple unknown keys → error names all"; "unknown nested key
     (out of scope?)".
   - `domain/config.md`: where keys are enumerated, current error style.
   - Presents all to **Human** for acknowledgement. **Human approves.**
4. **Orchestrator** `[rubato: git_status clean · go_test pass]`: test list
   non-empty, nothing red → picks the first item. Presents options to Human:
   *(a) specify "known keys still load"*, *(b) reorder*, *(c) amend list*.
   **Human picks (a).**
5. **Orchestrator → Specify.** Specify writes one test asserting a known-key
   config loads. **Motivates WHY:** "Starting from the passing baseline pins
   existing behaviour before we drive the new rejection — protects against
   regression (GOOS: start where you have footing)." Presents the test to
   **Human** for review **before commit**. Human approves → git-operator commits.
6. **Orchestrator** `[rubato: go_test pass]`: green already (behaviour exists) →
   strikes item, picks "one unknown key → error names it". Options to Human.
   **Human proceeds.**
7. **Specify (RED):** writes a test feeding `{ "unknown_key": 1 }`, expecting an
   error naming `unknown_key`. **WHY:** "This is the smallest increment that forces
   the rejection branch to exist; naming the key is asserted now so the message
   contract is driven, not retrofitted." **Human reviews → approves → commit.**
8. **Orchestrator** `[rubato: go_test fail — 1 failing, right reason]`: red for the
   intended reason → **Implement**.
9. **Implement (GREEN):** adds the minimal unknown-key check that errors naming the
   key. Nothing else. Presents the diff to **Human**. **Human approves → commit.**
10. **Orchestrator** `[rubato: go_test pass]`: green → offers optional refactor.
    **Human:** "skip." → picks "multiple unknown keys → error names all."
11. Loop repeats for the multi-key item (RED → GREEN).
12. **Orchestrator:** test list empty (nested-key item deferred as out-of-scope),
    all DoD scenarios pass → marks the feature a **done candidate**. Options to
    Human: *(a) audit now*, *(b) add nested-key behaviour*. **Human picks (a).**
13. **Orchestrator → Auditor.** Auditor fans out **security / quality / style**;
    returns `audit/config-loader.md`: "clean; minor — error message could share a
    helper (refactor opportunity, non-blocking)."
14. **Orchestrator** presents the audit to Human. **Human:** "apply the code
    refactor." → **Refactor (code)** extracts the helper;
    `[rubato: go_test pass]` unchanged. **Human approves → commit.**
15. **Orchestrator:** DoD met, audit clean, ledger (`tasks.md`) records the cycles.
    Change is ready to sync/archive. **Done.**

## Scenario B — Escape hatch: prior knowledge breaks mid-loop

**Context:** mid-way through Scenario A, at step 9's implementation.

1. **Implement (GREEN):** while making "unknown key → error" pass, discovers the
   loader **silently drops unknown keys by design elsewhere** — the DoD assumption
   ("unknown keys are always invalid") conflicts with existing behaviour relied on
   by another module.
2. **Implement → Orchestrator:** reports that it cannot make the test pass without
   contradicting existing, relied-upon behaviour. (It does *not* decide what to do —
   it reports.)
3. **Orchestrator** `[rubato: go_test fail]`: state + output don't add up to an
   obvious next step → **does not** push on. Escalates. Presents options to
   **Human**: *(a) escape to Discovery to revise the DoD*, *(b) mark the conflicting
   behaviour out-of-scope*, *(c) abandon this test, reorder.* **Human picks (a).**
4. **Orchestrator → Discovery** (escape). Discovery investigates, updates
   `domain/config.md` ("unknown keys are tolerated in module X for
   forward-compat"), and **amends the DoD** Gherkin: *WHEN an unknown key appears
   outside the compatibility allow-list THEN loading fails naming the key.* Updates
   `test-list.md` (adds "allow-listed unknown key still loads"). Presents to
   **Human**. **Human approves.**
5. **Orchestrator:** with a corrected DoD and test list, resumes the inner loop —
   picks the new allow-list test first. Loop continues as in Scenario A.

*Illustrates:* the phase agent stays blind (reports, doesn't decide); the
orchestrator recognises "this doesn't add up" and escalates rather than reasoning
through it; Discovery is the single place domain knowledge and the DoD are revised;
the Human authorises the pivot.

## Scenario C — Human short-circuit (deviation from autonomy)

**Context:** any point where the Orchestrator presents options.

1. At an orchestrator options prompt, the **Human** overrides: "Stop — I don't
   trust this test ordering; specify the multi-key case first."
2. **Orchestrator** complies, re-sequences, and re-presents the options with the
   multi-key item first. No phase agent has acted yet, so nothing is undone.
3. The loop continues under the Human's chosen order.

*Illustrates:* [`constitution.md`](../../governance/constitution.md) §1 (Human
sovereignty) and [`vision.md`](../../governance/vision.md) §1 (either partner may
pause). The orchestrator's options are proposals, never commitments.
