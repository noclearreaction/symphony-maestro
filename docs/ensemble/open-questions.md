# Open questions

Undecided, and what would settle each. Companion to [README](README.md).

## How a persona reaches an agent

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


## Where the trust boundary sits

Adjudicating over ACP puts the decision in a
process the part cannot reach. What it does not bound is what a part does with
permission once granted, or with the ambient filesystem.

| Position | Bounds | Cost |
| --- | --- | --- |
| Adjudicate over ACP | Permission. A granted call is granted, and `bash` reaches the machine. | None — Ensemble is already the client. |
| Container per turn | Reach. What a part can touch at all. | Worktrees become mounts and Ensemble supervises containers. A container lives for one turn. |

The second changes the deployment model rather than the configuration.

## How sections reach each other

Open. Either the repo is the medium and
sections communicate through refs and issues; or a human carries intent across by
switching sessions; or Ensemble relays with an explicit `message(section, text)`.
The first two need no mechanism. The third is buildable, but a section that can
message another section's principal can task it, and nothing in the permission
model speaks to that.

## What survives an Ensemble restart

A part is a transcript on disk and a record
of state, so durability is close to free. A turn in flight does not survive, since
the process running it dies with its parent. Either Ensemble supervises those or
they return `Faulted` with worktrees intact. The second is cheaper and should be a
decision rather than a consequence.

## Whether a turn is at-most-once

Work that edited files and lost its result
channel has still done the work. The worktree persists and the transcript is
resumable, so recovery is inspection rather than replay.


## Whether the queue stays FIFO

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

