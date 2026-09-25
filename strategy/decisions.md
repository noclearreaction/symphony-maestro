# Strategic Decisions (SDR)

This file tracks major strategic tradeoffs and architectural decisions made by the User, including *why* certain paths were chosen or rejected.

## SDR-001: Separation of Symphony Director and Symphony
* **Status**: Decided
* **Context**: We need a way to track strategic direction without cluttering our actual runtime code, and without coupling orchestration to application-level goals.
* **Decision**: We split these concerns into two conceptual layers. "Symphony Director" owns strategic continuity and has no runtime code. "Symphony" owns orchestration mechanics.
* **Consequences**:
  * We maintain a zero-dependency strategic sandbox (Symphony Director) that operates out-of-tree.
  * We avoid premature architecture in the runtime project repositories.

## SDR-002: Ensemble builds on Vibe over ACP, not on the Claude Agent SDK

* **Status**: Decided (2026-09-25), on a due-diligence survey run 2026-09-23. Evidence: [`docs/ensemble/prior-art.md`](../docs/ensemble/prior-art.md).
* **Context**: Before building Ensemble we asked whether anything already solves "operate an AI swarm while maintaining control over the entire instruction set", and whether Claude would be a better substrate than Mistral Vibe. On capability alone Claude is the stronger substrate: the Agent SDK replaces the system prompt outright, adjudicates tool permissions through `canUseTool` and `PreToolUse` hooks, sets `cwd` per session, and Claude Code ships git worktrees and an attach surface (`claude agents`). Those capabilities retire three of Ensemble's stated constraints.
* **Decision**: Build on **Vibe over ACP**, as originally planned. The deciding factor is commercial, not technical: Anthropic's terms bar third-party products from using claude.ai login or subscription rate limits, including agents built on the Agent SDK, absent prior approval. We hold a subscription and cannot use API keys. Claude Code's own parallelism (subagents, agent teams) runs on the subscription but is interactive-only and cannot be driven by an external orchestrator, which is what Ensemble is.
* **Consequences**:
  * The Agent SDK's capabilities are unavailable to us, so Ensemble keeps the parts of its design that exist to work around Vibe's limits.
  * Persona delivery is a **per-agent shim**, not a protocol feature. ACP carries no system prompt and the RFD proposing one is open with its reference PR closed unmerged. This is a build constraint, recorded in [`first-build.md`](../docs/ensemble/first-build.md).
  * Two design claims are corrected: §1's one-level delegation limit is Vibe-specific rather than a property of ACP, and §5's permission routing is observability rather than containment — the worktree is the boundary.
  * Merge-on-finish and re-apply-on-conflict are **not novel**; Amp ships both and hadron gates the merge. `open-questions.md` marks re-apply as untested; it is not untested in the field, only by us.
  * What remains genuinely ours to build is worktree lifecycle, the bounded queue, and merge on finish. No surveyed project supplies all three.
  * Read [`Sannrox/bugyo`](https://github.com/Sannrox/bugyo) and [`s0lda/hadron`](https://github.com/s0lda/hadron) before writing code; both implement this design and bugyo ships a Phase-0 ACP spike that covers what `spikes/acp` has not.
  * Revisit if we gain API-key access or Anthropic approval, which would reopen the Claude path and shrink the build substantially.