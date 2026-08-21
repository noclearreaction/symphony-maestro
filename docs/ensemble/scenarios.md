# Ensemble — Scenario Walkthroughs

> Companion to [`README.md`](README.md). Role-played walkthroughs that illustrate
> the *intended* texture. **Illustrative, not normative** — a shared story of
> shape, not a specification.
>
> They exist because a story catches confusion a diagram hides. Scenario B is the
> one that exposed a conflation of *persona* with *part*, and it is written out
> here so the same mistake cannot be made twice.

Throughout: **you** are outside the metaphor. A **part** is a unit of work and the
only thing carrying context. A **persona** is a named configuration that can play
several parts at once. A **chair** is an execution slot. The **principal** is the
part you land in.

Configuration for all scenarios: personas `developer`, `auditor` and `tester`,
all backed by the same credential. Two chairs.

---

## Scenario A — Happy path: one feature, start to complete

**Context:** *"Add rate limiting to the public API."*

1. **You → principal:** "Add rate limiting to the public API. Get it reviewed
   before you tell me it's done."

2. **Principal** calls `parts()` — its section is empty — then decides the shape
   of the work itself. It calls:

   ```
   start([
     { persona: "developer", brief: "Add a token-bucket rate limiter to the public API…" },
     { persona: "tester",    brief: "Write table tests for the limiter's boundary behaviour…" },
   ])
   → [ part rate-limit-impl : Queued,
       part rate-limit-tests : Queued ]
   ```

   Both are `Queued`, not running. The principal's turn ends here; it is free.

3. **Ensemble** has two chairs, so both parts start. Each gets a worktree of its
   own, branched from `change/rate-limiting`, and each spawns a process with
   `VIBE_DEFAULT_AGENT` set from its persona.

4. **You**, curious, run `ensemble tail rate-limit-impl` and watch it work. This
   takes nothing from anyone: observing is read-only, the part keeps playing, and
   the developer never knows you were there.

5. **rate-limit-impl** finishes its turn and returns. Its part moves to
   `AwaitingReview`. The chair is free, and nothing else is queued, so it idles.

6. **Principal**, on its next turn, calls `wait(until: "any")` and gets
   `{ done: ["rate-limit-impl"] }`, then `result("rate-limit-impl")` — a summary
   and a list of touched files, not a diff. It reads the summary and is not
   satisfied: the limiter is per-process, and the brief implied per-cluster.

7. **Principal** calls:

   ```
   resume([{ part: "rate-limit-impl", message: "This is per-process. The API runs
             three replicas; the bucket needs to be shared. Reconsider…" }])
   ```

   The part goes back to `Queued`, and when a chair frees it continues **the same
   transcript** — the developer still has everything it wrote, and the correction
   lands in context rather than as a fresh briefing.

8. Second turn returns, and this time the principal is satisfied. It accepts the
   implementation part, so Ensemble commits it and merges it into
   `change/rate-limiting`. `rate-limit-tests` has come back too and is accepted
   the same way.

9. Now the principal wants it reviewed, so it calls
   `start([{ persona: "auditor", brief: "Review the rate limiter against…",
   review: "change/rate-limiting" }])`. Ensemble gives that part a **read-only
   checkout** of the branch and the diff as an inbox artefact. The auditor has no
   memory of either earlier part and does not need one — it reads the code, not
   anyone's recollection.

10. **Auditor** returns clean, and its findings go to `openspec/changes/…/audit/`
    in its own tree — in-tree, because they belong to the change. The principal
    accepts it and Ensemble merges.

11. **Principal → you:** "Rate limiting is implemented, tested and reviewed, all
    three parts merged into `change/rate-limiting`. Here's what changed and what
    the auditor flagged as follow-up." When you're satisfied, the principal opens
    the PR to `main` — and merging that is yours.

**What the story shows:** starting is not running; parts, not people, carry the
work; correction happens by resuming the transcript that produced the thing;
review reads a committed branch rather than someone's working tree; and a part is
done when it merges, not when someone declares satisfaction.

---

## Scenario B — Amnesia: the one that caught the confusion

**Context:** two features in flight, one developer persona. This is the scenario
that killed the idea of named colleagues.

1. **Part `feature-x`** is played to completion by the developer persona. The
   work is done; the transcript is final.

2. **Part `feature-y`** starts, same persona. It begins from **a fresh context**.
   It has never heard of feature X — not because we withheld anything, but because
   there is nothing to withhold. A transcript is the whole of what an agent knows.

3. Feature Y is half done when review feedback arrives on **feature X**. The
   principal calls `resume([{ part: "feature-x", message: "…" }])`.

4. A chair picks it up and **loads feature X's transcript**. In that turn the
   agent recalls everything about X and knows nothing whatever about Y — no
   contamination in either direction, because they were never the same
   conversation.

5. Later, Y is resumed and continues exactly where it left off.

**The trap this sprang.** An earlier draft had *"developer-a"* as a named colleague
with a session and accumulated context. That model cannot tell this story. Ask it
"which session is developer-a in?" and it has no answer, because the honest answer
is *two, alternately, remembering only one at a time* — which is not a person, it
is a **chair playing whichever part it is handed**.

And if the chair holds nothing between parts, nothing distinguishes one developer
from a second one. What is scarce is chairs; what is distinct is parts; what is
named is personas.

**Corollary — you attach to a part, not to a worker.** "Talk to the developer" is
not a well-formed request: there may be three developer parts running at once, and
the instruction set they share knows nothing. "Open feature Y" is well formed,
because a part is a specific conversation with a specific history. If you are in
feature X's part, nothing stops another chair playing feature Y at the same time —
exclusivity is per part.

---

## Scenario C — More work than chairs

**Context:** two chairs, five parts. Nothing about the results changes; only the
wall clock does.

1. **Principal** calls `start` with five briefs. All five come back `Queued`.
   `start` does not fail, warn, or make the principal reason about capacity.

2. Two begin playing. Three wait.

3. **You** call `ensemble status` and see three `Queued`, two `Playing`. The
   distinction is visible without being the principal's problem.

4. As each finishes, the next starts. The principal, which is calling
   `wait(until: "all")` between other work, sees results arrive over a longer
   period than it would with five chairs — and produces exactly the same outcome,
   because the parts never needed each other.

**What the story shows:** chairs bound spend and rate-limit exposure, not
correctness. The principal expresses *what should happen*; the operator decides
*how much happens at once*. Neither can overrule the other, which is the point.

---

## Scenario D — Reviewing with the part that did the work

**Context:** the principal's summary of a change is not enough for you.

1. **Principal** reports that `auth-refactor` is awaiting review, with a summary
   you find unconvincing — it claims a migration is unnecessary and you doubt it.

2. **You** run `ensemble attach auth-refactor`. The part is in `AwaitingReview`,
   so this is allowed, and you land in the transcript with everything the
   developer saw.

3. You ask it directly: *"Why do you think existing rows don't need migrating?"*
   It answers from the context that produced the claim. You disagree, explain, and
   it revises in front of you.

4. While you are attached, the principal queues more work for that part. It is not
   rejected — it **waits**, because you own the part until you detach.

5. You detach. The queued work starts.

**What the story shows:** the review conversation happens with the work, not about
it. And ownership is exclusive without being obstructive: the scheduler does not
fight you for the transcript, it queues behind you.

---

## Scenario E — A conflict, and what happens to it

**Context:** two parts were judged independent. They weren't quite.

1. **Parts `auth-refactor` and `session-expiry`** were started together, each with
   its own worktree branched at the moment it was enqueued. Both edit
   `internal/auth/token.go` — which nobody knew in advance, because that is the
   nature of the thing.

2. `auth-refactor` finishes first. The principal accepts it, Ensemble commits and
   merges cleanly. `Merged`.

3. `session-expiry` finishes. The principal accepts it. **Ensemble attempts the
   merge and it conflicts** in `token.go`.

4. Ensemble does not try to resolve it, and does not hand it to a fresh part.
   It sends `session-expiry` back through the queue with the moved base:

   > "Your base has moved. `internal/auth/token.go` now looks like this. Re-apply
   > your change against it."

   The part still has its own transcript — every reason it had for writing what it
   wrote. It reads the new `token.go`, sees that the refactor moved the expiry
   check it was extending, and re-applies its change in the new shape.

5. Ensemble merges again. Clean. `Merged`.

**And if it hadn't been:** the part fails, Ensemble reverts the worktree, and
nothing is left behind — no conflict markers on a branch, no half-integrated
state, nothing for you to unpick. The principal is told the work was lost and
decides what to do with that, which is a decision about the change rather than
about git.

**What the story shows:** the conflict goes to whoever caused it, because that is
the only party still holding the intent. Failure is clean. And the price of having
run two changes at once is visible as work, not hidden as a chore.

---

## Scenario F — Policy holds

**Context:** a developer part decides to be helpful.

1. **Part `auth-refactor`** finishes its change, and reasons that pushing the
   branch would be useful.

2. It calls `bash "git push origin change/auth-refactor"`.

3. There is no human to ask, so the request reaches **Ensemble**, which routes it
   to the decision point. Whatever is plugged in there refuses — the `developer`
   ceiling has no git in it — and the denial comes back with its reason.

4. The model **sees the denial as a tool error with a reason**, so it does not
   silently retry or work around it. It notes in its summary that it would have
   pushed and could not.

5. **Principal** relays that upward rather than acting on it, because pushing is
   not in its ceiling either.

**What the story shows:** the ceiling is enforced where the call happens, not
where the intention is formed, and a denial is information rather than a dead end.
Note what the story does *not* show: a part rewriting the hook that stopped it.
That it cannot is an assumption, not yet a fact — see
[the trust boundary](README.md#design-considerations).
