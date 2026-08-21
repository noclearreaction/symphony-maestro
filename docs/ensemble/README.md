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
| **Persona** | An instruction set and which ACP agent runs it. |
| **Part** | A unit of work and the carrier of its context: a transcript, a persona, a tree, a lifecycle. |
| **Chair** | An execution slot. Plays one part at a time. |
| **Principal** | The part a human lands in by default. Orchestrator persona, working in the base. |

### Everything is a part

Context lives in the part: an agent knows its transcript. A part carries two
properties:

| Property | Values |
| --- | --- |
| **Persona** | Which instruction set, and which agent runs it |
| **Tree** | A worktree of its own, the base, a read-only checkout for review, or none |

A human may be attached to a part. That is state, held in the lifecycle
([§3](#lifecycle)).

The principal is a part with the base as its tree.

A part's permissions follow from its persona and are resolved outside Ensemble
([§5](#5-tool-use-with-no-human-present)).

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
turns. A principal reaches it over MCP.

Attaching is a lookup: a session is resumable by the agent's own command, so a
frontend gets the session id and execs the agent against it.

See [Architecture](diagrams/architecture.md).

Ensemble owns the sections and their parts, the mapping from part to session id,
the queue, and the routing of permission requests. Results, quality, plans and
policy belong elsewhere.

### Backend and frontends

A backend/frontend split: the backend holds the state — part records, transcripts,
worktrees, the queue — as files, and runs on demand rather than resident.

A lock bounds concurrent turns, and any invocation starts queued work when a chair
frees. A turn's process is the ACP client, so it answers permission requests and
writes its result for the next reader.

The backend decides whether a part can be attached and supplies the session id.

## 3. Parts

Ensemble bounds where a part works and what it may do; method is the part's own.
Anything it spawns runs inside those bounds — subagents inherit its working
directory.

### One owner at a time

One party drives a part: Ensemble on scheduled work, or a human attached to it.
A transcript has a single author.

**Attaching** gives the part to a human. A live part — `Playing` or `Attached` —
leaves queued work on the queue until it is free again.

**Observing** streams the same `session/update` notifications the driver sees, to
any number of readers, in any state including mid-turn.

### Lifecycle

`Queued` while it waits for a chair, `Playing` while an agent works, then one of
two ends depending on why the agent stopped. See
[Part lifecycle](diagrams/part-lifecycle.md).

`Playing` spans however many turns the work takes. Ensemble sees the start and the
stop; what happens between is the agent's.

Transcripts persist. `Reverted` concerns the tree: the worktree is discarded, the
transcript remains.

`resume` is the only way work reaches an existing part, whatever the reason — a
follow-up, new information, a conflict to re-apply, more turns for a halted part.
It continues the transcript rather than starting one, which keeps the context that
produced the work.

### Merging

`stopReason` says which end a part reached.

**Finished** — the agent stopped because the work was done. Ensemble commits the
worktree and merges it into the base.

- if the merge conflicts, the part is handed the moved base and its own previous
  diff, and re-applies its change — no conflict markers, no git
- if it cannot re-apply, the worktree is reverted

**Halted** — the agent stopped because it hit a bound on turns or spend. The
worktree stays unmerged: a part cut off mid-edit may leave a tree nothing should
build on. Whoever started the part reads the result and either `resume`s it with
more, or abandons it and the tree is reverted.

That is the one judgment in the model, and it belongs to the caller.

The outbox is kept either way. Merge order follows completion order.

### Starting is not playing

`start` puts a part on the queue. It begins playing once a chair is free and the
part is idle.

The queue holds only ready work. Sequencing is the orchestrator's judgment: it
enqueues "write the code" once the test exists, rather than enqueueing both and
leaving a scheduler to order them. There is no dependency graph to evaluate.

Dependent work therefore needs no shared tree. A part enqueued after its
predecessor merged branches from a base that already contains that predecessor's
output.

A part's worktree is created at enqueue, when the intended base is known. Deferring creation to selection branches from whatever HEAD
has since become.

### What bounds concurrency

Two configured constraints, both operator configuration:

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
| Part worktrees, commits, merge into the base, revert | **Ensemble** | Mechanical and deterministic, on trees Ensemble created |
| Anything above the base | **The principal** | Judgment and authored content, and it already works in the base |

The principal gets git itself rather than tools wrapping it: git flows vary, and a
fixed set of verbs would constrain it to whatever we anticipated.

The hard boundary is a protected ref, outside Ensemble, holding regardless of what
any agent attempts.

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

