# Ensemble

> **Status: provisional.** Nothing is built. The substrate is not finally
> committed and the tool surface is unimplemented. Items marked *measured* were
> verified against a running system; everything else is a proposal.

Ensemble is the orchestration and policy layer that **maestro** directs. It
schedules work across a section of parts, collects results, and routes permission
decisions to whatever decides them.

This document is the model. Everything that changes on a different clock lives
beside it:

| Document | Holds |
| --- | --- |
| [`scenarios.md`](scenarios.md) | The same model as narrative |
| [`contract.md`](contract.md) | The tool surface, its invariants, and the spaces work moves through |
| [`snippets.md`](snippets.md) | What a persona is told about the system it is inside |
| [`first-build.md`](first-build.md) | What the first build picks, its entry points, and the order |
| [`open-questions.md`](open-questions.md) | Undecided, and what would settle each |
| [`diagrams/`](diagrams/) | Architecture, lifecycle, contract types, flows |

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
([§5](#5-tool-use-with-no-human-present)). A principal names a persona when it
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
| An agent driven without a human cannot ask a question, but its tool calls raise permission requests, and those reach the client. | Ensemble answers them ([§5](#5-tool-use-with-no-human-present)). |
| `session/new` carries `cwd`, `additionalDirectories`, `mcpServers` and `_meta`. Nothing about instructions. | A part's tree and tool surface are protocol-level. Its instructions are not. |

An agent run headlessly is typically instructed to make a judgment call and
continue rather than stop, which turns ambiguity discovered mid-turn into a silent
guess. ACP sessions do not carry that instruction. That is the reason for ACP
rather than a CLI.

> **Cost differs by agent.** Work run by Vibe draws on a Mistral subscription;
> work run by Claude Code is directed to API keys by Anthropic's Agent SDK terms.

## 2. Architecture

Ensemble knows which parts exist, which session id each one is, and runs their
turns. A principal reaches it over MCP; a human through the CLI.

Attaching is a lookup: a session is resumable by the agent's own command, so
Ensemble supplies the session id and the CLI execs the agent against it.

See [Architecture](diagrams/architecture.md).

Ensemble owns the sections and their parts, the mapping from part to session id,
the queue, and the routing of permission requests. It does not interpret results,
judge quality, make plans, or decide policy.

### One core, several entry points

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

The core decides whether a part can be attached — it is in an attachable state and
nobody else owns it — and supplies the session id. Getting a terminal onto that
session belongs to the entry point.

## 3. Parts

How a part does its work is its own business. Ensemble bounds where it works and
what it may do, and does not prescribe method. Anything a part spawns runs inside
those same bounds: subagents inherit its working directory, and Ensemble neither
sees nor manages them.

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

## 4. Flows

- [Dispatch](diagrams/dispatch.md) — starting work, collecting a result, resuming
  a part.
- [Observing and attaching](diagrams/observe-attach.md) — observing takes nothing
  from a running part; attaching takes ownership.

## 5. Tool use with no human present

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

