# Open questions

Undecided, and what would settle each. Companion to [README](README.md).

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
they come back `Halted` with worktrees intact — the same shape as a turn that hit
a bound. The second is cheaper and should be a decision rather than a consequence.

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


## Whether finishing should merge at all

Ensemble owns both ends of a part's tree: it creates the worktree and it merges the
branch when the part finishes. The alternative is to own only the first end — make
the tree, hand back the branch, and let the caller decide what happens to it.

Claude Code's subagent isolation is that alternative, and it is *more general*: any
completion policy can be built on top, including this one. What it costs is that
integration becomes a property of how the caller was prompted rather than of the
system, so a merge happens when someone remembered to ask for one.

| | Merge on finish | Hand back the branch |
| --- | --- | --- |
| Integration is | a mechanism | the caller's to arrange |
| Suits | a base that is a change branch, parts that are small | flows we have not anticipated |
| Failure mode | forces a policy on work that wanted a different one | work that lands nowhere, silently |

The pick here is deliberate — a section drives a change to completion, and work
that merges nowhere is the thing worth designing out. But it is a pick, and the
general form is a superset. What would settle it is a flow that genuinely wants a
branch left alone; we do not have one yet.

## Whether a part can re-apply its own change

A conflict is handed back to the part that caused it, on the assumption it can
re-apply its intent against a moved base. That is empirical and untested. If it
cannot, the remedy is fewer concurrent parts, paying in wall-clock rather than in
integration.
