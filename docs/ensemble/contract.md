# Contract

The MCP tool surface Ensemble offers a principal, and the spaces work moves
through. Companion to [README](README.md).

Calls are batch-shaped: each tool boundary is a potential turn, and every turn
re-sends the whole context, so instructing five parts costs one turn.

See [Contract types](diagrams/contract-types.md).

Persona is absent: it is a name resolved out of configuration, not something
Ensemble holds. A part carries the name it runs under.

## Tools

| Tool | Blocks | Returns | Semantics |
| --- | --- | --- | --- |
| `personas()` | no | `Persona[]` | What can be played: names and the agent behind each. Read-only. |
| `parts()` | no | `Part[]` | The calling principal's section: ids, personas, states, queued work. |
| `start(items[])` | no | `Part[]` | `{persona, brief}`. Creates parts and queues their first turn. Does not fail for capacity or liveness. |
| `resume(items[])` | no | `Part[]` | `{part, message}`. Queues another turn on an existing transcript. |
| `wait(parts?, until, timeout)` | yes | `WaitResult` | `until` is `any` or `all`. Server caps the timeout and returns `pending`. |
| `status(parts?)` | no | `Status[]` | Cheap poll. Omit `parts` for the whole section. |
| `result(part, full?)` | no | `Result` | Summary by default; `full` returns the diff. Carries why the part stopped, which is what a halted part is judged on. |
| `cancel(parts[])` | no | `void` | Drops queued work, stops a playing part, or abandons a halted one. The worktree is reverted; the transcript remains. |

## Invariants

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
- **Merging is not a tool call.** A part that stops with the work done is committed
  and merged. One that stopped on a bound is not, and waits for its caller to
  `resume` or `cancel` it.

## How results reach the principal

Through `wait`. While a human is attached to a principal it is a live session, and
a live session cannot be driven from outside, so there is no callback path and no
delivery queue. The principal asks; Ensemble answers when it has something.

Ensemble could type into a live session and does not, for the same reason a part
has one owner.


## Channels: who reads and writes where

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
judgment a part is told how to make — one of the [snippets](snippets.md).

### What this does not attempt

A part with `bash` can write anywhere its process can reach, the devcontainer
filesystem included. That is a property of running an agent on a machine, and
channel design does not close it.

What Ensemble owes is not creating leaks of its own: no space writable by two
parts, no artefact without a declared writer, no visibility that skips the merge.
Closing ambient filesystem access belongs to the
[trust boundary](open-questions.md#where-the-trust-boundary-sits). Channel hygiene is not a substitute for
it, and treating it as one would present a convention as a guarantee.

