# Ensemble

> **Status: PROVISIONAL — not ratified, not agreed, nothing built.** The substrate
> is not finally committed, the tool surface is unimplemented, and the model below
> has been rewritten more than once as its assumptions failed. Treat every
> specific — names, states, the policy choice — as a proposal to argue with.
>
> Firmer than the rest are the [ground rules](#1-ground-rules) and anything marked
> as measured, because those were verified against a running system rather than
> reasoned about.

Design plan for the orchestration and policy layer that **maestro** directs: an
**orchestration layer** that schedules work and collects results, and a **policy
layer** that hands out privileges the work cannot exceed.

The ground rules were established against Mistral Vibe, since that is what we
run. The design is not tied to it: work is driven over ACP, which several agents
implement — see [One protocol, many agents](#one-protocol-many-agents).

[`scenarios.md`](scenarios.md) walks the same model as narrative. Read it if the
vocabulary below feels abstract — and read Scenario B in particular, which is
where an earlier version of this document was caught conflating a *persona* with
a *part*.

## Vocabulary

Orchestral throughout, and worth stating once. You are not in this table; you sit
outside the metaphor and own it.

| Term | Is |
| --- | --- |
| **Maestro** | The framework. Directs the ensemble. |
| **Ensemble** | This layer: the sections, their parts, and the chairs that play them. |
| **Section** | A long-lived working group — a build section, a design section — with a principal. Works through changes over time. |
| **Change** | The unit that becomes a pull request: an OpenSpec change, its `change/<name>` branch, and the parts working it. |
| **Persona** | An instruction set, and which ACP agent runs it — `developer`, `auditor`. Nothing more: not a tree, not a permission, not an individual. |
| **Ceiling** | What a persona's parts may do. Reviewed configuration, bound to a persona by the operator rather than carried by it. |
| **Part** | A unit of work, and the only thing that carries context. Has a transcript, a persona, a tree or none, an interactive flag, and a lifecycle. |
| **Chair** | An execution slot. Plays one part at a time. |
| **Principal** | A part whose properties are: orchestrator persona, works in the change branch, interactive, wide ceiling. Not a separate kind of thing. |

### Everything is a part

The agents have amnesia: nothing survives a transcript, so context lives entirely
in the part. Two parts running the same persona are therefore indistinguishable,
and there is no individual to name or track.

One kind of thing, carrying four properties:

| Property | Values |
| --- | --- |
| **Persona** | Which instruction set, and which agent runs it |
| **Tree** | Its own branch, the change branch, a read-only checkout of a branch under review, or none |
| **Interactive** | Whether a human is driving it, rather than a brief |
| **Ceiling** | What it may do — from reviewed configuration, keyed by its persona |

The principal is a combination of those rather than an exception to them, which
keeps scheduling, enforcement and the lifecycle uniform and makes *"run a
principal unattended overnight"* a flag rather than a new concept:

| | Developer part | Principal |
| --- | --- | --- |
| Git | None | Branch and PR operations |
| Tree | Its own part branch | The change branch |
| Interactive | No | Yes |
| Ends when | Its branch merges into the change branch | Its change branch merges into `main` |

The last row is the same shape one level up, with the human's gate on it — where
[`governance/constitution.md`](../../governance/constitution.md) puts it.

**The ceiling is not part of the persona.** A persona is a prompt, a
ceiling is a policy, and the operator's configuration binds one to the other. That
separation is what stops a principal handing a part more authority than a human
sanctioned.

## 1. Ground rules

Properties of Vibe 2.24.1 that determine the shape of everything else. Each was
verified against the shipped binary; the published docs are silent on most.

| Property | Consequence |
| --- | --- |
| **Delegation is one level deep.** A subagent's `task` tool refuses every target ([`task.py:97`](https://github.com/mistralai/mistral-vibe/blob/v2.24.1/vibe/core/tools/builtins/task.py)). | Hierarchy cannot live inside a session. Orchestration must be an external process. |
| **ACP strips the ability to ask.** `vibe/acp/agent.py` disables `ask_user_question` and `exit_plan_mode` via `NON_INTERACTIVE_DISABLED_TOOLS`, the same pair `-p` removes. | A non-interactive part cannot ask a question, whatever its instructions say. Approval gates belong to interactive parts. |
| **Permission requests do reach the client.** `client.request_permission()` is called for tool approvals, which is a different channel from asking a question. | Ensemble can adjudicate tool calls centrally. Untested — the spike never triggered one. |
| **Sessions are durable, addressable and resumable.** `--resume <id>` re-enters with full context and appends. | A part *is* a session. Every follow-up is a resume. |
| **`-p` forces a 385-char headless fragment** ending *"Override any earlier instructions that say to wait for confirmation."* Hardcoded in `_run_programmatic_mode()`; no flag disables it. | `-p` turns discovered ambiguity into a silent guess. Ruled out as the substrate. |
| **ACP sessions are not headless.** `vibe/acp/agent.py` omits `headless=` from `SessionOptions`, so it defaults to `False` — verified live: no `# Headless Mode` section, 384 chars shorter than the same work under `-p`. | The verified way to run work programmatically without the fragment. |
| **ACP sessions are ordinary Vibe sessions.** One created over ACP appears in `~/.vibe/logs/session/` under its `sessionId` like any other. | One session store behind both front-ends, so a human can attach with `vibe --resume <id>` whatever drove the part. |
| **Instructions are bound when the process spawns.** ACP offers no way to pick a profile, so `VIBE_DEFAULT_AGENT` at spawn is the lever — measured in [`spikes/acp`](../../spikes/acp/). | A part's persona is fixed for a turn. A process lives for one turn, so this is not a lifetime constraint. |

## 2. Architecture

The human runs Ensemble, not Vibe. Ensemble hosts every session — the
principal's included — and brokers the terminal to whichever the human wants.
It is reached as an MCP server over HTTP for the tool surface, but it does not
sit *under* a session; it sits around all of them.

The principal is a part like any other, with the properties set out in
[Everything is a part](#everything-is-a-part): the orchestrator persona, the
change branch as its tree, interactive, and a wide ceiling. It is scheduled,
enforced and tracked by the same machinery as the rest — a wider ceiling is not
an exemption from having one.

```mermaid
flowchart TB
  H["human"]
  ENS["<b>Ensemble</b><br/>session host · terminal broker · decision router"]

  subgraph S1["Section · build"]
    P1["principal<br/>interactive"]
    A1["part: auth-refactor<br/>persona developer · worktree"]
    A2["part: rate-limit-tests<br/>persona tester · worktree"]
    A3["part: review-auth<br/>persona auditor"]
  end

  subgraph S2["Section · design"]
    P2["principal<br/>interactive"]
    D1["part: queue-model<br/>persona designer"]
  end

  CH["chairs<br/>N execution slots"]
  DEC["decision point<br/>pluggable"]
  POL[("policy — whatever plugs in")]
  RUB["rubato"]
  API[("model API")]

  H <-->|"terminal"| ENS
  ENS -.->|"attach · observe"| P1
  ENS -.->|"observe"| A1
  ENS --> S1
  ENS --> S2
  ENS --> CH
  CH -->|"plays"| A1
  CH -->|"plays"| A2
  P1 -->|"MCP / http"| ENS
  P1 --> DEC
  A1 --> DEC
  DEC --> POL
  A1 --> RUB
  P1 --> RUB
  RUB --> API
```

Ensemble owns the sections and their parts, the mapping from part to session id,
the terminal attachment, and routing permission requests to whatever decides
them. It does not interpret results, judge quality, make plans, or decide policy.

### Substrate

The constraint that ruled the field: **the headless fragment is not acceptable**,
because it converts ambiguity discovered mid-turn into a silent guess.

| Candidate | Clears it | Standing |
| --- | --- | --- |
| `vibe -p` | No — hardcoded in `_run_programmatic_mode()` | Ruled out |
| **ACP** — [Agent Client Protocol](https://agentclientprotocol.com), JSON-RPC 2.0 over stdio | Yes, verified | Front-runner |
| `vibe/app_server` directly | Unknown — the core both Vibe front-ends sit on | Unexplored |

[`spikes/acp`](../../spikes/acp/) drove ACP end to end against the installed
`vibe-acp`: a real turn, routed through rubato, no headless fragment, and a
system prompt of our choosing. Mistral documents none of it, but the protocol is
not Mistral's — it is a published standard with an agent registry, so the risk is
"Mistral's implementation drifts", not "the contract is someone's internals".

Two spike findings shape what follows:

- **A process lives for one turn.** A session is a transcript on disk and needs
  nothing running; a process is required only while a turn executes. Spawn, load
  or create, prompt, read to `stopReason`, exit. Idle parts cost nothing.
- **`vibe-acp` does not exit cleanly when the client disconnects**, and reports to
  Sentry. Ensemble should close the stream and let it exit rather than killing it.

Still unverified: `session/load`, which every resume depends on.

### One protocol, many agents

ACP is **not Vibe-specific**. Claude Code, Codex CLI, Copilot CLI, OpenCode,
Gemini CLI, Cline and Goose implement it, Claude Code through the
[`claude-agent-acp`](https://github.com/agentclientprotocol/claude-agent-acp)
bridge. Both the spec and that bridge live under the neutral
`agentclientprotocol` org, Apache-2.0 — it is not Zed's, and Mistral implements
it independently.

So a persona names not only *what* a part does but *who plays it*. A section can
mix vendors, and Ensemble stays a client of one protocol rather than a wrapper
around one CLI.

#### Mixed sections are not free

One difference between vendors is contractual rather than technical. Anthropic's
Agent SDK — which the Claude bridge wraps — states that *"unless previously
approved, Anthropic does not allow third party developers to offer claude.ai
login or rate limits for their products, including agents built on the Claude
Agent SDK"*, directing integrators to API keys.

| Persona backed by | Auth | Cost |
| --- | --- | --- |
| Vibe | Mistral account | Included allowance first, then pay-as-you-go |
| Claude | API key, per the terms above | Per token from the first call |

Start Vibe-only. A Claude persona is a deliberate spend, not a free upgrade, and
changing that is a conversation with Anthropic rather than an engineering task.

#### The protocol standardises the conversation, not the configuration

Agent-specific knobs ride in ACP's `_meta`, so how a part is configured differs
per agent even though how it is driven does not. Ensemble needs a small adapter
per agent kind.

| | Vibe | Claude Code (via bridge) |
| --- | --- | --- |
| System prompt | Base replaceable; ten appended fragments, five gated | `_meta.systemPrompt: "…"` replaces outright, no residue |
| Project instructions | `AGENTS.md`, gated with the git snapshot | `settingSources` controls CLAUDE.md loading |
| Tool restriction | `enabled_tools` / `disabled_tools` per profile | `tools`, `disallowedTools`, `disableBuiltInTools` |
| Policy hook | `.vibe/hooks.toml` on disk | `hooks.PreToolUse` passed by the client at session start |

That last row matters for [§6](#6-enforcement). For a Claude-backed part the hook
is supplied by Ensemble when the session opens, rather than read from a file
inside the workspace the part can edit — closing the enforcer-in-the-blast-radius
problem for that agent without containers. Vibe has no equivalent, so the trust
boundary is not uniform across a mixed section, and policy has to assume the
weakest member.

### Shape: one backend, several frontends

A backend daemon with thin, replaceable frontends: CLI, TUI, later a VS Code
extension. Four things make that more than a preference.

**There are two consumer classes.** Principals consume the orchestration tool
surface over MCP; humans consume section views, attach, approvals and logs. Those
want different APIs, and collapsing them is how a tool-shaped API ends up
rendering tables for people.

**The backend is already a server.** It must serve MCP over HTTP regardless, so a
control API for frontends is an increment rather than a new tier.

**State must not live in a frontend**, or the TUI and the CLI become two products
that disagree. A frontend should be able to crash without losing anything.

**Attach is where the split gets tested**, and the one place frontends cannot be
uniform:

| Frontend | Attach means |
| --- | --- |
| CLI | exec `vibe --resume <id>` in the foreground — trivial |
| TUI | suspend itself and hand over the terminal, or embed a pty — the hard case |
| VS Code | likely open Vibe's own extension against the session rather than reimplementing it |

So `attach` is a backend-blessed capability — the backend says which session and
why — with a per-frontend mechanism.

#### CLI first

Not because it is cheapest, though it is. Because it validates the control API
without UI work confounding the question: if the CLI can drive the ensemble, the
API is complete. It stays useful permanently as the scripting surface, and it
makes owning attachment trivial to prove.

The trap is that a CLI naturally polls and exits, and a poll-shaped CLI bakes a
poll-shaped API. The TUI wants a live view, which is a subscription. So the first
CLI must include the `tail`-shaped observer — the smallest thing forcing a real
event stream into the API, after which the TUI is a rendering problem rather than
an API redesign.

The TUI is the eventual daily driver, since watching and reviewing beats issuing
one-shot commands. VS Code comes last and may be smaller than it sounds: Vibe's
own extension already provides the session UI, so Ensemble's part could be a
section panel that launches an attach.

## 3. Parts

### One owner at a time

**A part has exactly one owner.** Either Ensemble is driving it on scheduled work,
or a human is driving it interactively. Never both.

The reason is meaning, not concurrency safety. Two writers on one transcript break
both: the scheduler gets back a conversation containing turns it never sent, and
the human finds their session doing things they did not ask for. Neither can
reason about a conversation it only partly authored.

**Owning attachment** is interactive — the human takes the part and drives it. A
part is *live* when `Playing` or `Attached`, and a live part does not take work
off the queue. Queued work is not rejected; it waits, and starts when the human
detaches.

**Observing** is read-only and takes nothing. It streams the same
`session/update` notifications the driver sees, works in any state including
mid-turn, and any number can watch at once. In the CLI it is a `tail`-shaped
command; in the TUI it is what the section view is made of.

Keeping those separate is what lets you watch work happen without stopping it.

Later, the wait becomes a choice: show a part's queued work in the TUI and let the
human **pull** it into the session they are already in, rather than detaching and
letting the scheduler start it. Same exclusivity, with the handoff deliberate
instead of merely timed.

### Lifecycle

```mermaid
stateDiagram-v2
  direction LR
  [*] --> Queued: start(persona, brief)
  Queued --> Playing: a chair frees up
  Queued --> Reverted: cancel()
  Playing --> AwaitingReview: turn returned
  Playing --> Faulted: error · budget exceeded
  Playing --> Reverted: cancel()
  AwaitingReview --> Queued: resume(message)
  AwaitingReview --> Merging: accept()
  AwaitingReview --> Attached: attach()
  Faulted --> Attached: attach()
  Faulted --> Queued: resume(message)
  Faulted --> Reverted: give up
  Attached --> AwaitingReview: detach()
  Merging --> Merged: clean
  Merging --> Queued: conflict — re-apply against the moved base
  Merging --> Reverted: cannot re-apply
  Merged --> [*]
  Reverted --> [*]

  note right of Attached
    Live: a human owns the part.
    Queued work waits, it is not refused.
  end note
```

### Merged or reverted

**A part is done when its branch is merged.** Those are the only two ends, and
both are clean.

Merge is attempted **once**, at the end of the work — not per turn. Turns run
until the orchestrator accepts the part, and that judgment is about quality and
completeness. Any change at all is enough to proceed to a merge, **including a
partial one**: a part that made a coherent dent and merged has succeeded, and
whether the original problem is fully solved is answered by enqueuing more work,
not at the integration boundary.

Merge only happens clean. A conflict is **a signal to the part that caused it** —
it goes back through the queue with the moved base and re-applies its own change:

> "Your base has moved. `internal/auth/token.go` now looks like this. Re-apply
> your change against it."

That is what a developer does when a rebase turns ugly, and it beats handing the
conflict to a fresh part: the part that formed the intent is still there with its
reasoning intact, where a resolver would have to infer two intents from diffs.

**Presented as an artefact, not as a git operation.** Ensemble resets the part's
worktree to the moved base and hands it its own previous diff. No conflict
markers, no `--theirs`, no git tooling added to the part's ceiling — it edits
files exactly as it always does, which is the thing it is good at rather than a
git operation it is bad at.

If it cannot re-apply, the part **fails and its worktree is reverted**. Failure is
therefore clean — nothing half-integrated, nothing for a human to unpick, no tree
left behind. The change either landed or never existed.

Ordering falls out of completion order: first to merge wins, second re-applies.

Three tiers, in cost order — avoid by sharing a worktree, re-apply, or escalate to
a human when a part cannot re-apply its own change, since a second attempt will
not fix that.

> **Unverified.** Whether an agent reliably re-applies its intent against a moved
> base is an empirical question. If it does not, the answer collapses to the first
> tier: share more worktrees, run less concurrently, and pay in wall-clock rather
> than in integration.

`resume` covers more than rejection. Rejected output is one reason; a follow-up
question, new information, a conflict to re-apply, or a change of direction are
others. There is no separate `reject` verb — rejection is a resume whose message
says what was wrong. Resuming continues the existing transcript rather than
starting a new one: cheaper, since cached input bills at a tenth of fresh, and
correct, since the context that produced the work is still there.

`Faulted` is not retried automatically. A part that ran out of budget mid-edit
has left a worktree a human should see.

### Starting is not playing

`start` puts a part on the queue. It begins playing when a chair is free and the
part is not live — neither `Playing` nor held `Attached`. Until then it waits, and
`start` fails for neither reason.

**The queue only ever holds ready work.** Sequencing is the orchestrator's
judgment, not a mechanism: it does not enqueue "write the code" and "write the
test" together and leave a scheduler to work out the order — it enqueues the code
once the test exists. There is no dependency graph to evaluate, and termination is
structural rather than hoped for, because the work set is finite and nothing enters
the queue before it can run.

**A part's worktree is created at enqueue, not at selection.** That is the moment
the intended base state is known — the test exists, at *this* commit. Defer
creation to selection and the part branches from whatever HEAD has since become,
which is not what was decided.

### What bounds concurrency

Three configured constraints, each traceable to a real cause, all operator
configuration and none reachable from a tool:

| Constraint | Because |
| --- | --- |
| **Chairs** | Total spend and rate-limit exposure |
| **Per-credential concurrency** | The provider actually limits this — per account, not per persona |
| **Section reservation, then share** | One section must not stall another; reserve a minimum, share the surplus so idle sections do not waste chairs |

Nothing is capped per persona: a headcount says nothing about why, where a
credential's rate limit is the constraint that actually exists.

Chairs have no relationship to parts. Ten parts and one chair is ordinary, and so
is the reverse. The bound exists for tokens, spend and provider rate limits — not
CPU, since a chair spends nearly all its time waiting on a model API. Sizing it
like a CPU-bound pool would be a category error.

Beyond configuration, work can also be excluded by contention: parts needing the
same **mutex** cannot play at once. Mutexes are acquired atomically at dispatch —
a part gets its whole set or waits holding nothing — so there is no hold-and-wait
and therefore no deadlock. A mutex is held for the life of the part, not the turn,
because a worktree that shifts between turns leaves a part re-entering a checkout
its own transcript no longer matches.

The commonest mutex is a worktree, and there it is a correctness requirement
rather than a preference: **dependent work shares a worktree because it needs its
predecessor's output.** "Write the code" must see the test. Serialisation is what
"this builds on that" means, not a scheduling choice — and the scheduling
consequence follows from the logical structure of the work rather than driving it.

Independent work gets separate worktrees, runs concurrently, and pays at merge if
it turns out to overlap.

### Who owns git

Not one answer — the two levels have different owners, because one is mechanical
and the other is judgment.

| Level | Owner | Why |
| --- | --- | --- |
| Part worktrees, commits, merge into the change branch, revert on failure | **Ensemble** | Mechanical, must be deterministic, triggered by `accept`. Ensemble created those trees. |
| The change branch, the PR, its title and body, when it is ready | **The principal** | Authored content and judgment — and the principal is already sitting in the parent worktree. |

So **a part never touches git**, and its persona's ceiling can deny it outright.
The principal's cannot: it needs branch and PR operations, and it should not get
them through a tool wrapping `gh pr create`. Real git flows vary, and a principal
with git and judgment adapts where a principal with `open_pr(title, body)` is
stuck with whatever we imagined.

The boundary that matters is `main`, and it is better protected by **remote branch
protection** than by withholding git. The constitution's gate is that the human
decides what reaches `main`; branch protection enforces that regardless of what
any agent attempts, which is a stronger guarantee than a ceiling we wrote.

Committing a part's work is Ensemble's, but the message is content and this
repository cares about it — see [`bin/commit-lint.ts`](../../bin/commit-lint.ts).
So the part proposes the message in its result, and Ensemble commits and
validates. An unlintable message comes back the same way a conflict does.

### Instruction snippets are part of the deliverable

Every constraint in this document is something an agent has to be **told**, or it
will do the obvious wrong thing: try to commit, guess rather than report, work
outside its tree. Ensemble therefore ships a library of composable instruction
snippets, and a persona prompt is *role-specific content plus Ensemble's
contract* rather than something hand-written each time.

The snippets divide into two kinds.

**Ensemble's contract** — what the agent must know about the system it is inside:

- Git is not yours. You edit files in your worktree; Ensemble commits and merges.
  Propose a commit message with your result.
- Your worktree is the only tree you touch.
- Output that belongs to the change goes in-tree and is reviewed. Output that is
  evidence *about* the work — measurements, logs, a large report — goes to your
  outbox and is referenced, not committed.
- End your turn with an outcome. Partial progress is a legitimate result; say what
  you did and what you did not.
- If handed a re-apply brief, your previous change and the new base are attached.
  Re-apply the intent, do not attempt to merge.

**What replacing the base prompt costs.** A persona sets `system_prompt_id`,
which replaces Vibe's bundled base outright, so anything from it that is actually
load-bearing has to be reintroduced deliberately — the scratchpad rules, OS and
shell, which subagents exist. See [`docs/vibe/prompt/`](../vibe/prompt/) for what
each fragment contains and what gating it costs.

Composing rather than hand-writing matters because these are invariants. A
persona that forgets the git snippet does not fail visibly; it just tries to
commit, gets denied, and burns a turn learning what it should have been told.

### Channels: who reads and writes where

Context isolation is a design property, not a side effect — a developer part must
not see an auditor's reasoning, and personas are only interchangeable because
nothing leaks between parts. Artefacts are how that gets broken, so every space
needs a declared writer and declared readers.

| Space | Writer | Readers | Notes |
| --- | --- | --- | --- |
| Part working tree | the part | Ensemble, to commit and merge; a human attached or observing | No other part is given a route to it |
| A part's commit | Ensemble, on the part's behalf | Ensemble; a review part, via a checkout Ensemble prepares | Immutable, so it cannot shift under a reviewer |
| Part scratchpad | the part and its subagents | same | `/tmp/vibe-scratchpad-<prefix>-<random>`, per session, outside the worktree so it is never committed |
| Part inbox — brief, previous diff, conflict material | Ensemble | the part | One-way, in |
| Part outbox — summary, proposed commit message, and files | the part | Ensemble, then the principal | One-way, out. Durable and outside the worktree, so out-of-tree output does not pollute the change |
| Change branch | Ensemble, by merging | any part of that change, as its base; a human | The only inter-part channel |
| Change artefacts — proposal, design, the `tasks.md` ledger | the principal, in the parent worktree | parts, read-only | Parts do not coordinate through the ledger |
| `main` | the human, by merging a PR | everyone | |

**Ensemble gives no space two writers, and hands each part what it needs rather
than a route to someone else's work.** That enables isolation; it does not enforce
it, and the difference matters — see [what this does not
attempt](#what-this-does-not-attempt).

The designed path for review is a **commit**, not a tree. A working tree is where
a part is mid-thought and can change under anyone reading it; a commit cannot.
Ensemble gives a review part a read-only checkout of that commit plus the diff as
an inbox artefact, so the reviewer reads files as usual, needs no git, and does not
need a route into anyone's tree.

A human can read anything.

**Not everything a part produces is a code change.** Measurements, logs, a large
output, a decision writeup — those go to the outbox, and Ensemble references them
from the result rather than committing them. Some artefacts genuinely belong to
the change and should be in-tree and reviewed: the `audit/*.md` of
[`tdd-orchestration-intent.md`](../design/tdd-orchestration-intent.md) §5 is one.
Which of the two applies is a judgment a part has to be told how to make, so it is
one more instruction snippet: *in-tree because it belongs to the change, outbox
because it is evidence about the work.*

#### What this does not attempt

A part with `bash` can write anywhere its process can reach — the devcontainer
filesystem included. That is a property of running an agent on a machine, not
something Ensemble introduced, and no amount of channel design closes it.

What Ensemble owes is **not creating leaks of its own**: no space writable by two
parts, no artefact arriving without a declared writer, no visibility that skips
the merge. Closing ambient filesystem access is the [trust
boundary's](#design-considerations) business — hook, hook denying its own paths,
or container per turn — and treating careful channel design as a substitute for it
would be the worse error, because it would look like a guarantee while being a
convention.

**One thing to verify.** Scratchpad writes skip the permission prompt, and
[`docs/vibe/prompt/08-scratchpad.md`](../vibe/prompt/08-scratchpad.md) records
that the exemption survives disabling the fragment — suppressing it only stops
telling the model about a facility it still has. Whether it also bypasses a
`pre_tool` hook is a different question and the answer matters: hooks are
documented as firing *before* the permission prompt, which suggests policy still
sees the call. If that turns out to be wrong, the scratchpad is a hole in every
ceiling and needs closing rather than documenting.

### Subagents are Vibe's, not Ensemble's

A part or a principal may spawn subagents, and **Ensemble neither sees nor
controls them**. Vibe spawns them in-process at depth 1, inheriting the working
directory, so a subagent is always in the same worktree as whatever spawned it.

That is a useful property rather than a gap. Isolation holds for free — a part's
subagents cannot escape its worktree — and it gives the `git-operator` of
[`tdd-orchestration-intent.md`](../design/tdd-orchestration-intent.md) §6 a better
home than it had. Not an Ensemble concept, just a subagent the principal spawns
when it wants commit-message discipline or PR-body drafting done by something with
narrow instructions. It costs Ensemble nothing, because Vibe handles all of it.

### Two levels of git

**One branch matters: the change branch.** It is the tree the principal works in,
what parts branch from, what their work merges back into, and what becomes a pull
request. It is not a new concept — it is the `change/<name>` branch this repository
already uses, for the OpenSpec change already named in `openspec/changes/<name>/`.

Below that there are worktrees and commits, not branches anybody addresses. A part
gets a worktree; it publishes by committing; Ensemble merges that commit into the
change branch. Whether Ensemble also points a branch at the commit is an
implementation convenience — useful when a human wants to poke at it — and nothing
here depends on it.

That leaves Ensemble with exactly two mechanical git operations: **create a
worktree, and merge the work back.** Everything else at the change level is the
principal's, done with git directly rather than through tools wrapping it. Neither
merges the PR: that is the human's, which is where
[`governance/constitution.md`](../../governance/constitution.md) puts the gate.

A section is long-lived and works through many changes over time. It **should**
work them sequentially, one change in flight at a time, because concurrent changes
in one section is how you get the conflicts of [§3](#merged-or-reverted) without
the isolation of separate sections to justify them. But that is a **default, not
an invariant** — expressed as a configured limit on changes in flight, defaulting
to one. Claiming to know exactly how a git flow always works is how a design
acquires a rule someone has to fight later.

## 4. Functional contract

The MCP tool surface. Every call is batch-shaped: each tool boundary is a
potential turn, and every turn re-sends the whole context, so instructing five
parts must cost one turn rather than five.

```mermaid
classDiagram
  direction LR
  class Persona {
    +string name
    +string acp_agent
    +string instructions
    +string credential
  }
  class Ceiling {
    +string persona
    +Scope scope
  }
  class Part {
    +string id
    +string persona
    +string change
    +string session_id
    +string tree
    +bool interactive
    +State state
  }
  class Work {
    +string part
    +string message
    +int priority
    +Budget budget
  }
  class Status {
    +State state
    +int turns
    +float cost_usd
    +time since
    +string last_event
  }
  class Result {
    +string part
    +string summary
    +Artifact[] artifacts
    +DiffStat diff
    +float cost_usd
  }
  class Section {
    +string name
    +string principal
    +int max_changes_in_flight
  }
  class Change {
    +string name
    +string branch
    +string pr_url
  }
  Section "1" --> "0..n" Change : normally one
  Change "1" --> "0..n" Part
  Persona "1" --> "0..n" Part : played as
  Part "1" --> "0..1" Work : queued
  Persona "1" --> "1" Ceiling : bound by configuration
  Part --> Status
  Part --> Result : last
```

### Tool surface

| Tool | Blocks | Returns | Semantics |
| --- | --- | --- | --- |
| `personas()` | no | `Persona[]` | What can be played: names, the agent behind each, and the ceiling bound to it. Read-only — personas and ceilings change by review, not by call. |
| `parts()` | no | `Part[]` | The calling principal's section: ids, personas, states, queued work. |
| `start(items[])` | no | `Part[]` | `{persona, brief, priority?}`. Creates parts and queues their first turn. Does not fail for capacity or liveness. |
| `start_sync(items[], timeout)` | yes | `Result[]` | Start and wait in a single turn — the reason it exists. |
| `resume(items[])` | no | `Part[]` | `{part, message}`. Queues another turn on an existing transcript. |
| `wait(parts?, until, timeout)` | yes | `WaitResult` | `until` is `any` or `all`. Server caps the timeout and returns `pending` rather than hanging. |
| `status(parts?)` | no | `Status[]` | Cheap poll. Omit `parts` for the whole section. |
| `result(part, full?)` | no | `Result` | Summary by default; `full` returns the diff. |
| `accept(parts[])` | no | `Part[]` | Done with the work. Ensemble attempts the merge: clean means `Merged`, a conflict sends the part back to re-apply. |
| `cancel(parts[])` | no | `void` | Drops queued work, or stops a playing turn and keeps the transcript for inspection. |

### Invariants

- **No call hangs.** Blocking tools take a server-side cap — 60s working default —
  and return `{status: "pending"}` so the model can call again. MCP requests time
  out at both ends, and a hung call also freezes the principal against its human.
- **Results are summaries.** Full diffs only on request. Five parts reporting must
  not paste five diffs into the principal's context. Ensemble builds a `Result`
  from what the substrate returns — the final message is the summary, reasoning
  is dropped.
- **Context lives in the part, and nowhere else.** There is no cross-part memory,
  because the agents have none.
- **A part cannot exceed its ceiling, and no call widens one.** Ceilings come from
  reviewed configuration bound to a persona, not from whoever started the part.
- **Parts are durable.** They survive an Ensemble restart, because a part is a
  transcript on disk and a row of state, not a process.

### How results reach the principal

Through `wait`, and only through `wait`. A principal is a live interactive
session, and a live session cannot be driven from outside — so there is no
callback path and no delivery queue. The principal asks; Ensemble answers when it
has something.

Ensemble owns the terminal, so it *could* type into a live session. It does not,
for the same reason a part has one owner.

## 5. Dispatch flows

```mermaid
sequenceDiagram
  autonumber
  participant O as Principal
  participant E as Ensemble
  participant C as chair
  O->>E: start([{developer, "auth refactor"}, {tester, "rate limit tests"}])
  E-->>O: [Queued, Queued]
  Note over O: turn ends — principal free
  E->>C: a chair frees up, plays auth-refactor
  C-->>E: returned
  O->>E: wait(until=any, timeout=60)
  E-->>O: {done: [auth-refactor]}
  O->>E: result(auth-refactor)
  E-->>O: summary + artifacts
  O->>E: resume([{auth-refactor, "missing table test"}])
  E->>C: session/load, prompt
```

```mermaid
sequenceDiagram
  autonumber
  participant H as human
  participant E as Ensemble
  participant P as part
  H->>E: observe(auth-refactor)
  activate E
  E-->>H: session/update stream
  Note over P: still playing — observing takes nothing
  deactivate E
  H->>E: attach(auth-refactor)
  Note over E: only from AwaitingReview or Faulted
  E-->>H: terminal handed over
  Note over E: queued work waits until detach
```

## 6. Tool use with no human present

**The problem is narrow: a non-interactive part cannot ask, but its tool calls
still need answering.** ACP raises them as `session/request_permission`, which
arrives at the **client** — Ensemble. Something has to reply, and no human is
there to.

Ensemble is already sitting in the protocol.s own decision point, in a process the
part cannot reach.

Ensemble routes each request to a **pluggable decision point** and applies the
answer. It does not decide and does not define policy — keeping that separate is
what stops this layer acquiring one. Two things it owes the plug:

- **Fail closed.** A decision point that breaks must refuse, not permit.
- **Pass a denial back with its reason**, so the model adapts or reports rather
  than silently working around it.

How complete the coverage is depends on how much the agent asks about, which is a
dial on the persona rather than a limit of the mechanism: instructions and tool
settings that ask for everything produce a request for every call.

Untested. The spike never triggered a permission request, so this path has been
read and not run — and agent-specific escapes exist, such as Vibe auto-allowing
its scratchpad before permission resolution.

### Design considerations

**Where the trust boundary sits.** Adjudicating over ACP puts the decision in a
process the part cannot reach, which is a real boundary rather than a convention —
a part cannot rewrite a decision point it has no access to.

What it does not bound is what the part does with permission once granted, or with
the ambient filesystem. Two positions remain:

| Position | Strength | Cost |
| --- | --- | --- |
| Adjudicate over ACP | The decision is out of reach. A granted call is still granted, and `bash` still reaches the machine. | None — Ensemble is already the client. |
| Container per turn | The boundary is the kernel's. What a part can reach at all is bounded, not just what it is allowed to ask for. | Worktrees become mounts and Ensemble supervises containers — but a container lives for one turn, so nothing idles. |

The second is the only one that bounds reach rather than permission, and it
changes the deployment model rather than the configuration — which is why it is
worth deciding rather than drifting into.

**How sections reach each other.** Open, and worth leaving open. Either the repo
is the medium and sections communicate through branches, changes and issues; or
the human is the medium, carrying intent across by switching sessions; or Ensemble
relays with an explicit `message(section, text)`. The first two need no mechanism
and should be tried first. The third is buildable, but a section that can message
another section's principal is a section that can task it, and persona ceilings
say nothing about that.

**What survives an Ensemble restart.** More than it first appears. A part is a
transcript on disk and a row of state, so durability is close to free. What cannot
survive is a turn in flight, since the process running it dies with its parent.
The choice is whether Ensemble supervises those or lets them come back `Faulted`
with worktrees intact. The second is cheaper, and should be a decision rather than a consequence.

**Whether a turn is at-most-once.** Work that edited files and then lost its
result channel has still done the work. Since the worktree persists and the
transcript is resumable, recovery is inspection rather than replay — an argument
for worktrees arriving early.

## 7. Build order

Each stage is independently useful, and each is a place to stop if the value is
not there.

| Stage | Delivers | Proves |
| --- | --- | --- |
| 0 ✅ | **ACP spike** — [`spikes/acp`](../../spikes/acp/). | Done. A real turn, routed through rubato, no headless fragment, instructions ours via `VIBE_DEFAULT_AGENT`. |
| 1 | `session/load` on the spike. | Resume works over ACP — the assumption every review round depends on. |
| 2 | Backend plus CLI: host a part, own-attach to it, and `tail` it. | The entrypoint, and an event stream in the control API before anything is built on it. |
| 3 | Personas from configuration, composed from role content plus Ensemble's instruction snippets; `personas`, `parts`, `start`, `status`. Durable across restart. | Work can be started, told what system it is in, and survives Ensemble dying. |
| 4 | `wait`, `result`, the queue, chairs, and a worktree per part. | Work comes back, bounded by chairs rather than by how much was started. |
| 5 | `resume`, `accept`, the merge and re-apply on conflict, plus owning attach in `AwaitingReview`. | The review loop, and that a part either merges or reverts cleanly. |
| 6 | Route permission requests to a decision point, with one plug behind it. | A part can use tools with no human present. |
| 7 | A second section. | Sections are a first-class boundary. |
