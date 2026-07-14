# Decision Point: Reconciling OpenSpec Skills with the TDD Workflow

> **Status:** OPEN — awaiting Human decision.
> **Type:** Architecture decision record (format after
> [MADR](https://adr.github.io/madr/) / [Nygard ADR](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions)).
> **Related:** [`tdd-orchestration-intent.md`](tdd-orchestration-intent.md) §4–§5, §8(6);
> [`tdd-orchestration-scenarios.md`](tdd-orchestration-scenarios.md).
> **Supersedes intent §5 wording:** the "`tasks.md` = execution ledger" framing is
> revised here toward "`tasks.md` = coarse outer-loop worklist" (see Option 1).

## Context

The five OpenSpec framework skills (`openspec-explore`, `openspec-propose`,
`openspec-apply-change`, `openspec-sync-specs`, `openspec-archive-change`) are all
**schema-aware**: they read `openspec status --json` / `openspec instructions` to
discover the active schema's artifacts, `applyRequires`, and `contextFiles`. The
active schema is declared in [`openspec/config.yaml`](../../openspec/config.yaml)
as `schema: spec-driven`.

Under `spec-driven`, `tasks.md` is an **implementation plan** and
`openspec-apply-change` executes it as a **single-agent linear loop** ("make the
code changes required, mark `- [ ]` → `- [x]`"). That model conflicts with our
double-loop TDD intent, where implementation is driven test-first by a multi-agent
inner loop with per-step Human gates (see intent §4, §7; issue
[#42](https://github.com/noclearreaction/symphony-maestro/issues/42)).

The question: **do we fork/replace the skills, define a custom schema, or is this a
granularity problem we can bridge?**

## Key insight: two different "plans"

TDD refuses one kind of plan and requires another:

| Plan kind | Example | TDD stance |
| --- | --- | --- |
| **Implementation** plan | "write `foo`, then wire `bar`" | **Refused** — the "how" is driven by tests, not decided upfront |
| **Behavioural / acceptance** plan | "drive out scenario X, then scenario Y" | **Required** — this *is* the double-loop outer loop |

A coarse `tasks.md` whose items are *"use TDD to satisfy spec X"* is therefore **not**
the mock-planning trap. It is the outer-loop worklist, derivable directly from the
`specs/` Gherkin (Freeman & Pryce, double-loop TDD). The fine-grained, "unplanned"
red→green→refactor work happens **inside** each task, recorded in the living
`test-list`.

This reframes the conflict from *"the skills are incompatible"* to *"the skills
operate at the outer-loop granularity; we run a finer inner loop inside each
task."*

## Skill compatibility (independent of option)

| Skill | Status | Notes |
| --- | --- | --- |
| `openspec-explore` | ✅ Compatible | Thinking/investigation only; maps to Discovery's exploratory mode. |
| `openspec-sync-specs` | ✅ Compatible | Merges delta `specs/` → main specs; our DoD lives in `specs/`. |
| `openspec-archive-change` | ✅ Compatible | Bookkeeping, orthogonal to how work was produced. |
| `openspec-propose` | ⚠️ Depends on option | Fine if it emits a **coarse** `tasks.md`; conflicts if it emits an implementation plan. |
| `openspec-apply-change` | ⚠️ Depends on option | Its linear single-agent loop must be **wrapped** (delegate each task to the orchestrator) or replaced. |

The decision below is really about `propose` and `apply` only.

## Options

### Option 1 — Granularity bridge: coarse `tasks.md` on `spec-driven` *(recommended)*

Keep the `spec-driven` schema. Define each `tasks.md` item at the **acceptance /
BDD granularity** — e.g. *"Use TDD to satisfy spec `config-loader`, scenario
'unknown key → error'."* The inner red→green→refactor loop runs inside each task
via the orchestrator and phase agents; `test-list` holds the living inner detail.

- **`propose`**: usable — instruct (via schema `rules` in `openspec/config.yaml`)
  that `tasks.md` items are coarse "TDD-implement spec X" entries, and add
  `test-list` guidance. No implementation-step tasks.
- **`apply`**: **wrapped** — each task's "make the code changes" step becomes
  "invoke the orchestrator loop until this spec's scenarios pass, then mark `[x]`."
- **`explore` / `sync` / `archive`**: unchanged.

**Pros:** stays on community tooling (vision §6); smallest change; no dependency on
custom-schema support; `tasks.md` checkmarks double as coarse progress. **Cons:**
`apply`'s single-agent assumption must be wrapped; the coupling of "task" to
"orchestrator invocation" is a convention, not enforced by the CLI.

### Option 2 — Custom `tdd-driven` schema

Define a new OpenSpec schema whose artifacts are the enrichment set (`test-list`,
`specs/` DoD, `domain/`, `audit/`, coarse `tasks.md`) and whose apply semantics are
the red-green-refactor loop.

- **`propose` / `sync` / `archive` / `explore`**: reused via schema-awareness.
- **`apply`**: still needs replacing/wrapping for the multi-agent inner loop.

**Pros:** cleanest separation; artifacts named for what they are. **Cons:**
**depends on whether OpenSpec supports user-defined schemas — currently UNKNOWN and
must be verified**; higher upfront effort.

### Option 3 — Fork the generative skills

Copy `propose` and `apply` into TDD-specific variants under
[`.opencode/skills/`](../../.opencode/skills/) (copies already exist there) and edit
them to emit coarse tasks + drive the orchestrator loop.

- **`explore` / `sync` / `archive`**: keep upstream.

**Pros:** full control; no dependency on schema extensibility. **Cons:** we own
maintenance and drift from upstream (tension with vision §6 community alignment).

## Decision matrix

| Criterion | Opt 1 (granularity) | Opt 2 (custom schema) | Opt 3 (fork) |
| --- | --- | --- | --- |
| Community alignment (vision §6) | High | High | Low |
| Upfront effort | Low | High | Medium |
| Depends on unverified capability | No | **Yes** (custom schemas) | No |
| `apply` still needs work | Wrap | Wrap/replace | Rewrite |
| Artifact naming clarity | Medium (`tasks.md` overloaded) | High | Medium |
| Maintenance burden | Low | Low–Medium | High |

## Recommendation

Adopt **Option 1** (granularity bridge). It confirms the Human's hypothesis:
`tasks.md` at BDD/SDD granularity ("use TDD to implement spec X") bridges the gap
with the skills as written, keeps us on community tooling, and reduces `apply` from
"replace" to "wrap." Keep **Option 2** as a later refinement *if* custom schemas
turn out to be supported and `tasks.md` overloading becomes painful.

**Residual gap to resolve regardless of option:** `openspec-apply-change` assumes a
single agent editing inline and marking done. Our inner loop is multi-agent with
per-red/green Human approval. The apply behaviour must be adapted so "execute a
task" means "hand off to the orchestrator," not "edit files directly."

## Open questions feeding this decision

1. **Does OpenSpec support user-defined schemas?** Determines whether Option 2 is
   even available. *Needs verification before Option 2 can be chosen.*
2. **Can `tasks.md` item granularity be enforced**, or only conventionally
   encouraged via schema `rules`?
3. **Where does the `test-list` live** relative to the change folder so `propose`
   can scaffold it and `apply` can read it?

## Consequences (if Option 1 is accepted)

- Update intent §5 to describe `tasks.md` as the coarse **outer-loop worklist**
  (not an "execution ledger"); the ledger role folds into task checkmarks +
  commit history + `audit/`.
- Add `test-list` (and `domain/`, `audit/`) to the schema via
  `openspec/config.yaml` `rules`/`context` so `propose` scaffolds them.
- Specify the `apply` wrapping behaviour in the orchestrator contract change.
