# Ensemble

> **Status: provisional.** Nothing is built. The substrate is not finally
> committed and the tool surface is unimplemented. Items marked *measured* were
> verified against a running system; everything else is a proposal.

Ensemble is the orchestration and policy layer that **maestro** directs. It
schedules work across a section of parts, collects results, and routes permission
decisions to whatever decides them.

[`scenarios.md`](scenarios.md) walks the same model as narrative.

## Vocabulary

| Term | Is |
| --- | --- |
| **Maestro** | The framework. Directs the ensemble. |
| **Ensemble** | This layer: sections, their parts, and the chairs that play them. |
| **Section** | A long-lived working group with a principal. |
| **Base** | The git ref a part's worktree is created from and merged into. What it signifies belongs to the project. |
| **Persona** | An instruction set and which ACP agent runs it. Not a tree, not a permission, not an individual. |
| **Part** | A unit of work, and the only carrier of context: a transcript, a persona, a tree or none, a lifecycle. |
| **Chair** | An execution slot. Plays one part at a time. |
| **Principal** | The part a human lands in by default. Orchestrator persona, working in the base. |

### Everything is a part

Agents have no memory outside a transcript, so context lives in the part. Two
parts running the same persona are indistinguishable, and there is no individual
to name or track.

A part carries two properties:

| Property | Values |
| --- | --- |
| **Persona** | Which instruction set, and which agent runs it |
| **Tree** | A worktree of its own, the base, a read-only checkout for review, or none |

Whether a human is connected is ephemeral state, not a property — the `Attached`
of [§3](#lifecycle).

The principal is a combination of the same properties:

| | Developer part | Principal |
| --- | --- | --- |
| Git | None | Yes |
| Tree | A worktree from the base | The base |
| Ends when | Its tree merges into the base | A human says so |

A principal nobody is attached to is running unattended. That is a state, not a
mode.

What a part may do follows from its persona and is decided outside Ensemble
([§6](#6-tool-use-with-no-human-present)). A principal names a persona when it
starts a part and does not supply permissions, so it cannot grant more than a
human has sanctioned.

## 1. Constraints

Ensemble drives work over **ACP**, the
[Agent Client Protocol](https://agentclientprotocol.com). Vibe, Claude Code, Codex
CLI, Copilot CLI, OpenCode, Gemini CLI, Cline and Goose implement it, so a part
can be run by any of them.

Five properties of that substrate shape the design. *Measured* —
[`spikes/acp`](../../spikes/acp/) and [`docs/vibe/prompt/`](../vibe/prompt/) hold
the evidence.

| Constraint | Consequence |
| --- | --- |
| A session is a durable transcript, addressable and resumable. | A part *is* a session. Continuing one is a resume. |
| A process serves one turn, and instructions are bound when it starts. | Idle parts cost nothing. A part's persona is fixed for a turn. |
| Delegation inside a session is one level deep. | Orchestration must be an external process. |
| An agent driven without a human cannot ask a question, but its tool calls raise permission requests, and those reach the client. | Ensemble answers them ([§6](#6-tool-use-with-no-human-present)). |
| `session/new` carries `cwd`, `additionalDirectories`, `mcpServers` and `_meta`. Nothing about instructions. | A part's tree and tool surface are protocol-level. Its instructions are not. |

An agent run headlessly is typically instructed to make a judgment call and
continue rather than stop, which turns ambiguity discovered mid-turn into a silent
guess. ACP sessions do not carry that instruction. That is the reason for ACP
rather than a CLI.

> **Cost differs by agent.** Work run by Vibe draws on a Mistral subscription;
> work run by Claude Code is directed to API keys by Anthropic's Agent SDK terms.

### Unknown: how a persona reaches an agent

The protocol carries no instructions, so delivering a persona is agent-specific.
The two mechanisms we know differ in kind:

| Agent | Instructions arrive by |
| --- | --- |
| Claude Code, via its bridge | `_meta.systemPrompt` on `session/new` |
| Vibe | an environment variable read at process spawn |

Whatever bridges this spans launching the process as well as composing messages.
For Vibe the work is entirely at launch.

This is a prerequisite: without it Ensemble cannot start a part, because nothing
in the protocol tells an agent what it is. Undecided: what the seam is called,
whether it also owns credentials and model selection, and how much of it a second
agent shares.

## 2. Architecture

Ensemble knows which parts exist, which session id each one is, and runs their
turns. A principal reaches it over MCP; a human through the CLI.

Attaching is a lookup: a session is resumable by the agent's own command, so
Ensemble supplies the session id and the CLI execs the agent against it.

See [Architecture](diagrams/architecture.md).

Ensemble owns the sections and their parts, the mapping from part to session id,
the queue, and the routing of permission requests. It does not interpret results,
judge quality, make plans, or decide policy.

### Shape: one core, several entry points

The core is file-backed state and the logic over it — part records, transcripts,
worktrees, the queue — all on disk. Entry points are spawned on demand and share
nothing but those files: an MCP server, a CLI, later a TUI and an extension.

The filesystem carries coordination. A lock bounds how many turns run at once, and
any invocation can start queued work when a chair frees. A turn is driven by a
process living exactly as long as the turn; that process is the ACP client, so it
is also what answers permission requests, and it writes its result where the next
reader finds it.

The MCP surface serves an agent making tool calls. A human wants different things
from the same state: which parts exist, what they are doing, attach, logs. The
entry points are therefore separate rather than one surface with two audiences.

State stays in the core, so entry points cannot disagree and any of them can crash
without loss.

Attach is the one capability entry points cannot implement uniformly:

| Entry point | Attach means |
| --- | --- |
| CLI | exec the agent against the session in the foreground |
| TUI | suspend itself and hand over the terminal, or embed a pty |
| VS Code | open the agent's own extension against the session |

The core decides whether a part can be attached — it is in an attachable state and
nobody else owns it — and supplies the session id. Getting a terminal onto that
session belongs to the entry point.

#### CLI first

The CLI validates the control API without UI work confounding the question: if the
CLI can drive the ensemble, the API is complete. It remains the scripting surface
afterwards.

A CLI polls and exits, and a poll-shaped CLI produces a poll-shaped API. A TUI
needs a live view, which is a subscription. The first CLI therefore includes the
`tail`-shaped observer, which forces an event stream into the API early; the TUI
then becomes a rendering problem.

The TUI is the eventual daily driver. VS Code comes last: where an agent ships its
own extension, Ensemble's part is a section panel that launches an attach.

## 3. Parts

### One owner at a time

A part has exactly one owner. Either Ensemble is driving it on scheduled work, or
a human is driving it. Never both.

Two writers on one transcript leave the scheduler holding a conversation
containing turns it never sent, and the human in a session doing things they did
not ask for. Neither can reason about a conversation it only partly authored.

**Owning attachment** gives the part to a human. A part is *live* when `Playing` or
`Attached`, and a live part does not take work off the queue. Queued work waits and
starts when the human detaches.

**Observing** is read-only. It streams the same `session/update` notifications the
driver sees, works in any state including mid-turn, and any number of observers can
watch at once. In the CLI it is a `tail`-shaped command; in the TUI it is what the
section view is made of.

Separating the two allows watching work without stopping it.

Later: show a part's queued work in the TUI and let a human pull it into the
session they are in, rather than detaching and letting the scheduler start it.
Exclusivity is unchanged; the handoff becomes deliberate rather than timed.

### Lifecycle

See [Part lifecycle](diagrams/part-lifecycle.md).

A transcript is never deleted. `Reverted` concerns the tree: the worktree is
discarded, and the transcript remains for inspection.

### Merged or reverted

A part's tree, if it has one, is merged into its base when the principal accepts
the part. Outbox contents are not merged.

The merge is attempted once, at accept, not per turn. Acceptance is a judgment
about quality and completeness; the merge is not. Partial work merges as readily
as complete work, and whether the original problem is solved is answered by
starting more work.

A conflict goes back to the part that caused it, which re-applies its own change
against the moved base:

> "Your base has moved. `internal/auth/token.go` now looks like this. Re-apply
> your change against it."

The part that formed the intent still holds its reasoning. A fresh part would have
to infer two intents from diffs.

Ensemble presents this as an artefact rather than a git operation: it resets the
worktree to the moved base and attaches the part's own previous diff. The part
edits files as it always does — no conflict markers, no git tooling it would not
otherwise have.

If it cannot re-apply, its worktree is reverted. Failure leaves nothing
half-integrated.

Ordering follows completion order: first to merge wins, second re-applies.

> **Unverified.** Whether an agent reliably re-applies intent against a moved base
> is empirical. If it does not, the remedy is fewer concurrent parts, paying in
> wall-clock rather than in integration.

`resume` covers more than rejection: a follow-up question, new information, a
conflict to re-apply, a change of direction. There is no separate `reject` verb.
Resuming continues the existing transcript, which is cheaper — cached input bills
at a fraction of fresh — and keeps the context that produced the work.

`Faulted` is not retried automatically. The worktree is left for a human to see.

### Starting is not playing

`start` puts a part on the queue. It begins playing when a chair is free and the
part is not live. `start` fails for neither reason.

The queue holds only ready work. Sequencing is the orchestrator's judgment: it
enqueues "write the code" once the test exists, rather than enqueueing both and
leaving a scheduler to order them. There is no dependency graph to evaluate.

Dependent work therefore needs no shared tree. A part enqueued after its
predecessor merged branches from a base that already contains that predecessor's
output.

A part's worktree is created at enqueue, not at selection. At enqueue the intended
base state is known. Deferring creation to selection branches from whatever HEAD
has since become.

### What bounds concurrency

Two configured constraints, both operator configuration and neither reachable from
a tool:

| Constraint | Because |
| --- | --- |
| **Chairs** | Total spend and rate-limit exposure |
| **Per-credential concurrency** | The provider limits this, per account |

Nothing is capped per persona. Two parts running one persona are identical
configuration, and the limit that exists is the credential's.

Chairs have no relationship to parts: ten parts and one chair is ordinary, and so
is the reverse. The bound exists for tokens, spend and rate limits rather than CPU,
since a chair spends most of its time waiting on a model API.

#### Later: the queue need not be FIFO

Work can carry a priority and an estimated cost, and the queue can schedule on
them. The contract should not foreclose it.

Priority pays only when partial results are actionable. Where a principal needs
all five outputs before proceeding, ordering gains nothing — completion is set by
the slowest item. It pays when finishing one item unblocks someone while the rest
run on.

The signal worth scheduling on is whether finishing an item unblocks anyone, which
the principal knows and the scheduler cannot infer.

Priority is safe to expose to a principal where chair count is not: it reorders
work without widening resource use. That is the test for anything added to the tool
surface — does it let a principal spend more, or only spend it in a different
order?

### Who owns git

| Level | Owner | Because |
| --- | --- | --- |
| Part worktrees, commits, merge into the base, revert | **Ensemble** | Mechanical and deterministic, triggered by `accept`, on trees Ensemble created |
| Anything above the base | **The principal** | Judgment and authored content, and it already works in the base |

A part can be denied git outright, since it never needs it. A principal cannot,
and should not receive git through tools wrapping it: git flows vary, and a fixed
set of verbs constrains it to whatever we anticipated.

The hard boundary belongs outside Ensemble. A protected ref holds regardless of
what any agent attempts.

Committing a part's work is Ensemble's, but the message is content: the part
proposes it in its result and Ensemble commits with it. Whether a message is
checked against a project's conventions belongs to that project, and a rejection
returns the same way a conflict does.

### Subagents belong to the agent

A part or a principal may spawn subagents, and Ensemble neither sees nor controls
them. They run in-process, inheriting the working directory, so a subagent is
always in the same tree as whatever spawned it.

That bounds them usefully: a part's subagents cannot leave its worktree. It also
gives narrow-instruction work — commit-message drafting, for instance — a home
that costs Ensemble nothing.

### Instruction snippets

Every constraint here is something an agent must be told, or it will attempt to
commit, guess rather than report, or work outside its tree. Ensemble ships
composable snippets, and a persona prompt is role content plus Ensemble's
contract.

Ensemble's contract:

- Git is not yours. You edit files in your worktree; Ensemble commits and merges.
  Propose a commit message with your result.
- Your worktree is the only tree you touch.
- Output belonging to the work goes in-tree and is reviewed. Evidence *about* the
  work — measurements, logs, a report — goes to your outbox and is referenced.
- End your turn with an outcome. Partial progress is a legitimate result; say what
  you did and what you did not.
- If handed a re-apply brief, your previous change and the new base are attached.
  Re-apply the intent; do not attempt to merge.

These are invariants, which is why they are composed rather than hand-written. A
persona missing the git snippet fails invisibly: it attempts a commit, is denied,
and spends a turn learning what it should have been told.

How much of an agent's own instructions survive a persona replacing them is
agent-specific. For Vibe, [`docs/vibe/prompt/`](../vibe/prompt/) records which
fragments are replaced and which are appended.

### Channels: who reads and writes where

A developer part must not see an auditor's reasoning, and personas are
interchangeable only because nothing passes between parts. Every space therefore
has a declared writer and declared readers.

| Space | Writer | Readers | Notes |
| --- | --- | --- | --- |
| Part working tree | the part | Ensemble, to commit and merge; a human attached or observing | No other part is given a route to it |
| A part's commit | Ensemble, on the part's behalf | Ensemble; a review part, via a checkout Ensemble prepares | Immutable, so it cannot shift under a reviewer |
| Part scratchpad | the part and its subagents | same | Per session, outside the worktree, never committed |
| Part inbox — brief, previous diff, conflict material | Ensemble | the part | One-way, in |
| Part outbox — summary, proposed commit message, files | the part | Ensemble, then the principal | One-way, out. Durable and outside the worktree |
| The base | Ensemble, by merging | any part branched from it; a human | The only inter-part channel |
| Whatever the principal writes in the base | the principal | parts, read-only | Parts do not coordinate through it |

Ensemble gives no space two writers, and hands each part what it needs rather than
a route to someone else's work. That enables isolation without enforcing it — see
[what this does not attempt](#what-this-does-not-attempt).

Review reads a commit, not a tree. A working tree changes under anyone reading it;
a commit does not. Ensemble gives a review part a read-only checkout of the commit
plus the diff as an inbox artefact, so the reviewer reads files as usual and needs
no git.

A human can read anything.

Output that belongs to the work goes in-tree and is reviewed; evidence about the
work goes to the outbox and is referenced from the result. Which applies is a
judgment a part is told how to make, and is one of the snippets above.

#### What this does not attempt

A part with `bash` can write anywhere its process can reach, the devcontainer
filesystem included. That is a property of running an agent on a machine, and
channel design does not close it.

What Ensemble owes is not creating leaks of its own: no space writable by two
parts, no artefact without a declared writer, no visibility that skips the merge.
Closing ambient filesystem access belongs to the
[trust boundary](#design-considerations). Channel hygiene is not a substitute for
it, and treating it as one would present a convention as a guarantee.

## 4. Functional contract

The MCP tool surface. Calls are batch-shaped: each tool boundary is a potential
turn, and every turn re-sends the whole context, so instructing five parts costs
one turn.

See [Contract types](diagrams/contract-types.md).

Persona is absent: it is a name resolved out of configuration, not something
Ensemble holds. A part carries the name it runs under.

### Tool surface

| Tool | Blocks | Returns | Semantics |
| --- | --- | --- | --- |
| `personas()` | no | `Persona[]` | What can be played: names and the agent behind each. Read-only. |
| `parts()` | no | `Part[]` | The calling principal's section: ids, personas, states, queued work. |
| `start(items[])` | no | `Part[]` | `{persona, brief}`. Creates parts and queues their first turn. Does not fail for capacity or liveness. |
| `resume(items[])` | no | `Part[]` | `{part, message}`. Queues another turn on an existing transcript. |
| `wait(parts?, until, timeout)` | yes | `WaitResult` | `until` is `any` or `all`. Server caps the timeout and returns `pending`. |
| `status(parts?)` | no | `Status[]` | Cheap poll. Omit `parts` for the whole section. |
| `result(part, full?)` | no | `Result` | Summary by default; `full` returns the diff. |
| `accept(parts[])` | no | `Part[]` | Done with the work. Ensemble merges the tree if there is one: clean gives `Merged`, a conflict returns the part to re-apply. |
| `cancel(parts[])` | no | `void` | Drops queued work, or stops a playing turn. The worktree is reverted; the transcript remains. |

### Invariants

- **No call hangs.** Blocking tools take a server-side cap and return
  `{status: "pending"}` so the model can call again. A hung call also freezes the
  principal against its human.
- **Results are summaries.** Full diffs only on request. Ensemble builds a
  `Result` from what the substrate returns: the final message is the summary,
  reasoning is dropped.
- **Context lives in the part.** There is no cross-part memory.
- **No call widens what a part may do.** Permissions follow from its persona and
  are decided outside Ensemble.
- **Parts are durable.** A part is a transcript on disk and a record of state, not
  a process.

### How results reach the principal

Through `wait`. While a human is attached to a principal it is a live session, and
a live session cannot be driven from outside, so there is no callback path and no
delivery queue. The principal asks; Ensemble answers when it has something.

Ensemble could type into a live session and does not, for the same reason a part
has one owner.

## 5. Flows

- [Dispatch](diagrams/dispatch.md) — starting work, collecting a result, resuming
  a part.
- [Observing and attaching](diagrams/observe-attach.md) — observing takes nothing
  from a running part; attaching takes ownership.

## 6. Tool use with no human present

A part with nobody attached cannot ask, but its tool calls still need answering.
ACP raises them as `session/request_permission`, which arrives at the client —
Ensemble. Nothing else is there to reply.

Ensemble is already in the protocol's decision point, in a process the part cannot
reach.

It routes each request to a **pluggable decision point** and applies the answer. It
does not decide and does not define policy. Two things it owes the plug:

- **Fail closed.** A decision point that breaks refuses rather than permits.
- **Pass a denial back with its reason**, so the model adapts or reports rather
  than working around it.

Coverage depends on how much the agent asks about, which is a dial on the persona:
instructions and tool settings that ask for everything produce a request per call.

> **Untested.** The spike never triggered a permission request. Agent-specific
> escapes also exist, such as Vibe auto-allowing its scratchpad before permission
> resolution.

### Design considerations

**Where the trust boundary sits.** Adjudicating over ACP puts the decision in a
process the part cannot reach. What it does not bound is what a part does with
permission once granted, or with the ambient filesystem.

| Position | Bounds | Cost |
| --- | --- | --- |
| Adjudicate over ACP | Permission. A granted call is granted, and `bash` reaches the machine. | None — Ensemble is already the client. |
| Container per turn | Reach. What a part can touch at all. | Worktrees become mounts and Ensemble supervises containers. A container lives for one turn. |

The second changes the deployment model rather than the configuration.

**How sections reach each other.** Open. Either the repo is the medium and
sections communicate through refs and issues; or a human carries intent across by
switching sessions; or Ensemble relays with an explicit `message(section, text)`.
The first two need no mechanism. The third is buildable, but a section that can
message another section's principal can task it, and nothing in the permission
model speaks to that.

**What survives an Ensemble restart.** A part is a transcript on disk and a record
of state, so durability is close to free. A turn in flight does not survive, since
the process running it dies with its parent. Either Ensemble supervises those or
they return `Faulted` with worktrees intact. The second is cheaper and should be a
decision rather than a consequence.

**Whether a turn is at-most-once.** Work that edited files and lost its result
channel has still done the work. The worktree persists and the transcript is
resumable, so recovery is inspection rather than replay.

## 7. Choices for the first build

The design above is unspecific about agents, plugs and entry points. The first
build has to pick. These are picks, not commitments — each is the cheapest thing
that exercises the design.

| Concern | First build | Because |
| --- | --- | --- |
| ACP agent | **Vibe**, via `vibe-acp` | Installed, routed through rubato, funded by a subscription we hold. Verified in [`spikes/acp`](../../spikes/acp/). |
| Policy plug | a static allow/deny list from configuration | Proves the decision point is wired without deciding anything about policy. |
| Entry point | **CLI**, including the `tail`-shaped observer | Validates the control API, and forces an event stream in early. |
| Language | **TypeScript on Deno** | The ACP SDK is first-class there, and it matches the repository. |
| Sections | **One** | A second section proves the boundary, which is not the first thing worth knowing. |

### Not in the first build

- **Containers per turn.** The trust boundary stays where ACP puts it, bounding
  permission but not reach.
- **Priority, budgets, section reservation.** The queue is FIFO over ready work
  with a chair count.
- **TUI and VS Code.** The CLI proves the API first.
- **A second agent vendor.** The adapter seam exists so adding one is not a
  rewrite.

## 8. Build order

Each stage is independently useful, and each is a place to stop.

| Stage | Delivers | Proves |
| --- | --- | --- |
| 0 ✅ | **ACP spike** — [`spikes/acp`](../../spikes/acp/). | Done. A real turn, routed through rubato, no headless instruction, the persona ours. |
| 1 | `session/load` on the spike. | Resume works over ACP — the assumption every review round depends on. |
| 2 | Core plus CLI: start a part, attach to it, `tail` it. | The entry point, and an event stream in the control API. |
| 3 | Personas from configuration; `personas`, `parts`, `start`, `status`. | Work can be started and survives a restart. |
| 4 | `wait`, `result`, the queue, chairs, a worktree per part. | Work returns, bounded by chairs. |
| 5 | `resume`, `accept`, merge and re-apply on conflict, attach in `AwaitingReview`. | The review loop, and that a part merges or reverts cleanly. |
| 6 | Route permission requests to a decision point, with one plug behind it. | A part can use tools with no human present. |
| 7 | A second section. | Sections are a first-class boundary. |
