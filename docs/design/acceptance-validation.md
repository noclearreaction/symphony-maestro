# Design Note: Acceptance (DoD) Validation Strategy

> **Status:** Design intent. Extends
> [`tdd-orchestration-intent.md`](tdd-orchestration-intent.md) §4 (outer loop /
> DoD) and §6.5 (Auditor). Captures how the Gherkin acceptance suite is *validated*
> across two horizons.

## The question

The outer-loop **Definition of Done** is the Gherkin acceptance suite in
`openspec/specs/*.md` (intent §4). How do we know an implementation actually
**satisfies** it — and that the suite itself is **sound**? Two mechanisms, on two
horizons.

## Long-term: godog (executable acceptance)

[**godog**](https://github.com/cucumber/godog) is the official Cucumber BDD
framework for Go. It runs Gherkin `.feature` files against Go **step definitions**
via `go test`, and — in Cucumber's own words — *"merges specification and test
documentation into one cohesive whole"* as **living documentation** (the specs are
always current because they are executed).

- **Role:** make the Gherkin DoD **executable**. A green godog run **mechanically
  proves fulfillment** of every scenario it covers.
- **Why:** this is the "**ubiquitous DSL**" goal — one shared Gherkin language
  between Human, agents, and the running system. A passing suite is *automation,
  not judgment*, which is exactly what [`governance/vision.md`](../../governance/vision.md)
  §3 (*Automation over Intelligence*) asks us to prefer.
- **Format bridge (open):** our current `specs/*.md` use OpenSpec's delta
  `### Requirement` / `#### Scenario` WHEN/THEN format, whereas godog consumes
  `.feature` files (`Given/When/Then`) plus registered step definitions
  ([Gherkin reference](https://cucumber.io/docs/gherkin/reference/)). A translation
  layer (generate `.feature` from the OpenSpec scenarios, or map step defs onto
  them) must be designed. See [§ Open questions](#open-questions).

## Interim (and later reinforcement): the AI Auditor

Until step definitions exist, the **Auditor agent** (intent §6.5) stands in. Its
job is **not to prove the work succeeds — it is to prove the work *fails*.** It
attacks two **distinct** properties of the cucumber spec, seeking a
counterexample in each:

- **(a) Reasonable** — the Auditor tries to show the Gherkin is *un*reasonable:
  ill-formed, internally contradictory, ambiguous, under-specified, or
  **untestable**. A single genuine contradiction or gap refutes "reasonable." This
  is a review of the **DoD itself**, independent of any code.
- **(b) Fulfilled** — the Auditor tries to show the implementation does *not*
  satisfy the spec: it hunts a scenario, edge case, or input for which the
  behaviour diverges, citing **evidence** (a failing/absent test, a code path, an
  observed behaviour). A single genuine divergence refutes "fulfilled." This is a
  review of the **implementation against the DoD**.

### The Auditor's stance: prove it *fails*

The Auditor operates by **refutation**, not confirmation. Its success condition is
**finding a failure** — a counterexample, a contradiction, an unmet scenario, a
coverage gap. This is a direct application of two long-standing results:

- **Popper (falsifiability):** there is a logical asymmetry between verification
  and refutation — no finite amount of confirming evidence verifies a universal
  claim, but *a single counter-instance refutes it*
  ([SEP: Karl Popper](https://plato.stanford.edu/entries/popper/); *Conjectures and
  Refutations*, 1963). "It works" is not establishable; "it fails here" is.
- **Dijkstra:** *"Program testing can be used to show the presence of bugs, but
  never to show their absence"* (EWD249, *Notes on Structured Programming*, 1970 —
  [E.W. Dijkstra Archive](https://www.cs.utexas.edu/users/EWD/ewd02xx/EWD249.PDF)).

Consequences for the verdict:

- A **PASS is never "proven correct."** It is *"sustained refutation attempts found
  no failure (yet)"* — provisional corroboration, not verification. It carries the
  record of *what was attacked and how hard*, so the pass is legible and
  re-testable.
- The Auditor's **evidence is its refutation attempts**, not confirmations: the
  scenarios probed, edge cases constructed, inputs fuzzed, contradictions searched
  for. An audit that lists no serious attack is not a pass — it is an empty test.
- This is the structural antidote to confirmation theatre: **an agent whose job is
  to break the work cannot rubber-stamp it.** A clean pass must be *earned* against
  a documented battery of break-attempts.

### Why the Auditor is *reinforcement*, not just a stopgap

The two mechanisms cover different things — which is why the Auditor survives even
after godog lands:

| Property | godog (executable) | AI Auditor |
| --- | --- | --- |
| **(a) Reasonable** | ✗ — will faithfully run an incoherent or wrong-intent spec | ✓ — judges spec soundness |
| **(b) Fulfilled** | ✓ — proven *mechanically* | ✓ — argued via *evidence* (weaker) |

- **Interim:** the Auditor owns **both** (a) and (b). Its (b) verdict is
  judgment-based, so it is explicitly a **stand-in** for godog.
- **Long-term:** godog owns **(b)** mechanically; the Auditor retains **(a)**
  reasonableness plus **coverage-gap detection** (does the godog suite actually
  exercise every scenario and the real intent?) and spot-check reinforcement of (b).

godog can prove a suite passes; it can never tell you the suite was *worth*
passing. That judgment stays with the Auditor (and, finally, the Human —
[`constitution.md`](../../governance/constitution.md) §1).

## Honesty & risk: avoiding "confirmation theatre"

An AI Auditor asserting "fulfilled" risks a failure analogous to the **mock-TDD**
concern (see [`superpowers-bridge-review.md`](superpowers-bridge-review.md)): an
eager auditor that rubber-stamps. The primary guard is the **refutation stance**
above — an auditor that succeeds only by *finding failure* cannot pass by
confirmation. Per vision §3, judgment is a *fallback* for where automation does not
yet exist; the interim (b) verdict is exactly that fallback, to be superseded by
godog. Reinforcing mitigations:

1. **Independent context.** The Auditor must not review work from the implementer's
   context ([`AGENTS.md`](../../AGENTS.md): *an agent should not review its own work
   from the same context*). It runs on artifacts + evidence only — otherwise it
   inherits the implementer's blind spots and cannot find their failures.
2. **Evidence-is-attempts.** Every verdict cites the **refutation attempts made**
   (scenarios probed, edge cases, inputs) and their outcomes — never a bare
   "it works." (Aligned with `superpowers-bridge`'s evidence-first retrospective.)
3. **Reachable failure verdict.** The Auditor must be able — and expected — to
   return **"NOT fulfilled"** or **"spec unreasonable"** and route back to Discovery
   via the escape hatch (intent §2, §6.1). A verdict with no reachable failure mode
   is theatre.
4. **Provisional PASS.** An interim "no failure found" is provisional (Popperian
   corroboration), flagged *judgment-based* until a godog scenario backs it.
5. **Human is final authority** (constitution §1, §4).

## Where it plugs into the loop

- The double loop's **outer-loop "test"** (intent §4) is:
  **godog run (long-term) OR Auditor fulfillment argument (interim).**
- In the interim, inner-loop *green* is Go **unit** tests (via rubato `go_test`),
  which sit *below* the acceptance DoD. There is a **gap** between unit-green and
  acceptance-fulfilled; the Auditor bridges it by arguing the unit-level
  implementation fulfills the acceptance Gherkin. **godog later closes that gap
  mechanically.**
- Artifact: Auditor findings live in `audit/*.md` (intent §5), including the
  explicit (a)/(b) verdicts and their evidence.

## Open questions

1. **OpenSpec-scenario → godog `.feature` bridge.** Generate `.feature` files from
   the OpenSpec delta scenarios, or map step definitions directly onto them? Affects
   whether the DoD is *one* artifact or a source + generated pair.
2. **Auditor independence mechanics.** How is the "separate context" enforced in the
   opencode/rubato runtime — fresh subagent, distinct system prompt, evidence-only
   inputs?
3. **Coverage-gap metric.** Once godog exists, how does the Auditor measure that the
   suite covers every scenario *and* the intent (not just that it is green)?
4. **Interim → long-term migration.** Do interim Auditor (b) verdicts get
   retro-actively replaced by godog runs, and is that tracked per scenario?

## References

- [godog — official Cucumber BDD framework for Go](https://github.com/cucumber/godog)
- [Cucumber — Gherkin Reference](https://cucumber.io/docs/gherkin/reference/)
- [Cucumber — BDD](https://cucumber.io/docs/bdd/)
- [Dan North — Introducing BDD](https://dannorth.net/introducing-bdd/)
- [Karl Popper — falsifiability (SEP)](https://plato.stanford.edu/entries/popper/); *Conjectures and Refutations* (1963)
- [Edsger W. Dijkstra — EWD249, *Notes on Structured Programming* (1970)](https://www.cs.utexas.edu/users/EWD/ewd02xx/EWD249.PDF) ("testing shows the presence, not the absence of bugs")
- Local: [`tdd-orchestration-intent.md`](tdd-orchestration-intent.md) §4, §6.5;
  [`superpowers-bridge-review.md`](superpowers-bridge-review.md) (evidence-first,
  mock-TDD risk); [`governance/vision.md`](../../governance/vision.md) §3–§4;
  [`AGENTS.md`](../../AGENTS.md) (independent review context).
