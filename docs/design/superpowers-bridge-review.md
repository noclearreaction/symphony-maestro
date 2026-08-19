# Review: `superpowers-bridge` OpenSpec Schema & Superpowers TDD

> **Status:** PROVISIONAL review notes — a reading of external prior art, not a
> decision. Informs
> [`tdd-openspec-reconciliation.md`](tdd-openspec-reconciliation.md).
> **Verdict:** *Borrow the mechanisms, reject the execution philosophy.* The schema
> is strong prior art that **validates Option 2** (custom schemas work and are
> first-class), but its TDD execution model risks exactly the "mock TDD" failure
> this project is trying to avoid (see
> [`tdd-orchestration-intent.md`](tdd-orchestration-intent.md) §2, §7).

## What it is

- **`superpowers-bridge`** ([repo](https://github.com/JiangWay/openspec-schemas/tree/main/superpowers-bridge),
  [schema.yaml](https://raw.githubusercontent.com/JiangWay/openspec-schemas/main/superpowers-bridge/schema.yaml))
  is a community OpenSpec custom schema that wires OpenSpec artifact governance to
  [obra/superpowers](https://github.com/obra/superpowers) execution skills.
- **Artifact DAG:** `brainstorm → proposal → specs → tasks → plan → verify →
  retrospective`, with `apply.requires: [plan]` and `apply.tracks: tasks.md`.
- **Execution:** `apply` runs `superpowers:using-git-worktrees` →
  `superpowers:subagent-driven-development` (transitively activating
  `test-driven-development` and `requesting-code-review`) →
  `openspec-verify-change` → `retrospective` → archive →
  `finishing-a-development-branch`.

## Where it risks "mock TDD" (reject)

The Superpowers pipeline front-loads the **how** into `plan.md` and then runs the
red-green cycle over a pre-decided implementation.

1. **`plan.md` pre-writes the implementation.** Superpowers' `writing-plans` skill
   produces tasks with "exact file paths, **complete code**, verification steps"
   ([obra/superpowers README](https://github.com/obra/superpowers), *How it works*
   / *The Basic Workflow*). If the code is already in the plan, then "write a
   failing test → make it pass" degenerates into *paste the pre-authored code*.
   This is precisely the failure named in intent §2 and in the project's founding
   concern: *"decide how everything works, then do mock TDD to pretend we had an
   incremental process."* The red-green becomes theatre over a settled design.

2. **No Canon TDD test list.** There is no living behavioural backlog
   ([Kent Beck, *Canon TDD* step 1](https://newsletter.kentbeck.com/p/canon-tdd)).
   The flow jumps `specs → plan`, so the "which behaviour do we drive out next,
   and what did we just discover" thinking has nowhere to live. Our intent makes
   the `test-list` a first-class **living** artifact owned by Discovery; the
   schema has no equivalent.

3. **Autonomy over human gating.** Superpowers advertises the executor working
   "autonomously for a couple hours at a time without deviating from the plan"
   ([README](https://github.com/obra/superpowers), *How it works*). That directly
   contradicts our requirement that specify/implement request Human review **before
   every commit** and that the orchestrator present each step as options
   ([`governance/constitution.md`](../../governance/constitution.md) §1, §4;
   [`governance/vision.md`](../../governance/vision.md) §1–§2). It optimises
   throughput; we optimise understanding/ownership.

4. **No escape hatch to discovery.** The pipeline is plan-execution oriented
   ("without deviating from the plan"). When prior knowledge breaks mid-loop —
   the Human's second hard constraint (intent §2) — there is no cheap path back to
   re-discovery; the design is assumed settled at `plan.md`.

5. **No named orchestrator.** The transition decisions live implicitly inside the
   subagent executor, not in a lightweight, human-facing state machine. Our design
   deliberately *names* that node (intent §7) and keeps the Human in it.

**Root divergence:** Superpowers is **plan-centric** (brainstorm → detailed plan
*with code* → execute autonomously). Our model is **test-list-centric and
human-gated** (discovery → living test list + DoD → orchestrator presents next
step → Human approves → one red test *with WHY* → review → minimal green →
review). The difference is not cosmetic: keeping the "how" genuinely emergent
(code appears only at green time) is the entire point.

## What is genuinely worth borrowing (adopt)

The **schema machinery and several hygiene patterns are excellent** and directly
reusable regardless of the execution critique above.

| Mechanism | Why adopt | Source |
| --- | --- | --- |
| **Custom schema via `openspec schema fork`** | Validates Option 2; CLI-validated, listed by `openspec schemas`, per-change selectable. No skill files modified. | [customization.md](https://github.com/Fission-AI/OpenSpec/blob/main/docs/customization.md); [opsx.md](https://github.com/Fission-AI/OpenSpec/blob/main/docs/opsx.md) *Custom Schemas* |
| **`apply.requires` ≠ `apply.tracks`** | Lets a fine execution artifact drive apply while coarse `tasks.md` tracks progress — exactly the granularity bridge from reconciliation Option 1. | schema.yaml `apply:` block; README *Why `plan` is separate from `tasks`* |
| **Evidence-based PRECHECK** | Shell checks (`git log … | wc -l > 0`, `grep -c '^- \[x\]'`) at an instruction's start enforce ordering the DAG engine can't (it only checks file existence). | README *Six design touches* #5–#6 |
| **Evidence-first `retrospective`** | §0 quantitative front-matter + analysis, each claim citing a commit/file/fact. Maps onto our Auditor/`audit/` artifacts and hygiene (vision §4). | schema.yaml `retrospective`; README *Apply phase* #4 |
| **`verify` artifact with explicit checks** | Structural validate, task completion, delta-spec sync — good DoD gating. | schema.yaml `verify`; README *Apply phase* #3 |
| **Prompt-level integration principle** | Integration lives in `instruction:` fields, not by editing skills — upstream upgrades don't break the schema. | README *Six design touches* #2 |
| **Fresh subagent per task** | Context separation; matches our blind, history-free phase agents (AGENTS.md). | schema.yaml `apply` step 2 |
| **"Delete code written before its test"** | A genuinely Canon-TDD-enforcing rule worth adopting in our implement (GREEN) agent. | [obra/superpowers](https://github.com/obra/superpowers) *test-driven-development* |
| **Entry/exit gates (ceremony scales with risk)** | "Bug fix / typo → direct PR, not a change" aligns with vision §7. | README *Entry & exit gates* |

## Impact on the decision

- **Option 2 is confirmed viable.** Custom schemas are first-class
  (`openspec schema fork/init/validate`, stored in `openspec/schemas/`,
  version-controlled). The earlier "critical unknown" in
  [`tdd-openspec-reconciliation.md`](tdd-openspec-reconciliation.md) (open question
  #1) is now **resolved: YES**.
- **We should build our own `tdd-driven` schema**, informed by
  `superpowers-bridge`'s *mechanisms* but **not** adopting Superpowers' plan-centric
  execution. Concretely, our schema differs by:
  - adding a **living `test-list`** artifact (no Superpowers equivalent);
  - replacing `plan.md`-with-code with a **behaviour-only** inner loop where code
    emerges at GREEN time;
  - inserting the **orchestrator** + **Human-approval-per-step** gates in place of
    autonomous multi-hour execution;
  - making **escape-to-discovery** a first-class transition.
- **Combination 1 + 2 (the Human's preference) is fully supported:** a custom
  `tdd-driven` schema (Option 2) whose `tasks.md` stays coarse and whose
  `apply.tracks: tasks.md` while a finer inner loop drives execution (Option 1).

## References

- [`superpowers-bridge` README](https://github.com/JiangWay/openspec-schemas/tree/main/superpowers-bridge)
  · [schema.yaml](https://raw.githubusercontent.com/JiangWay/openspec-schemas/main/superpowers-bridge/schema.yaml)
- [obra/superpowers](https://github.com/obra/superpowers) — execution skills & methodology
- [OpenSpec customization docs](https://github.com/Fission-AI/OpenSpec/blob/main/docs/customization.md)
  · [OpenSpec OPSX workflow docs](https://github.com/Fission-AI/OpenSpec/blob/main/docs/opsx.md)
- [Kent Beck, *Canon TDD*](https://newsletter.kentbeck.com/p/canon-tdd)
- Local: [`tdd-orchestration-intent.md`](tdd-orchestration-intent.md),
  [`tdd-openspec-reconciliation.md`](tdd-openspec-reconciliation.md)
