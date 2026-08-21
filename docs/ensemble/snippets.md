# Instruction snippets

What a persona must be told about the system it is inside. Companion to
[README](README.md).


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
- Partial progress is a legitimate result. Say what you did and what you did not,
  rather than guessing to fill the gap.
- If handed a re-apply brief, your previous change and the new base are attached.
  Re-apply the intent; do not attempt to merge.

These are invariants, which is why they are composed rather than hand-written. A
persona missing the git snippet fails invisibly: it attempts a commit, is denied,
and spends a turn learning what it should have been told.

How much of an agent's own instructions survive a persona replacing them is
agent-specific. For Vibe, [`docs/vibe/prompt/`](../vibe/prompt/) records which
fragments are replaced and which are appended.

