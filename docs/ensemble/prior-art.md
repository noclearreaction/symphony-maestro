# Prior art and substrate

What exists, what the substrate supplies, and what is left to build. Companion to
[README](README.md). Surveyed 2026-09-23.

Claims here were checked against schemas, source and official documentation.
Anything marked *unverified* was not.

## The binding constraint

The Agent SDK is unavailable on a Claude subscription:

> Unless previously approved, Anthropic does not allow third party developers to
> offer claude.ai login or rate limits for their products, including agents built
> on the Claude Agent SDK. Use the API key authentication methods described in the
> Quickstart instead.
> — [Agent SDK overview](https://code.claude.com/docs/en/agent-sdk/overview)

Claude Code's own parallelism — subagents, and agent teams under
`CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` — runs on a subscription but is
interactive-only: *"In non-interactive mode with the `-p` flag, including Agent
SDK sessions, Claude doesn't spawn teammates."* Agent teams also have no worktree
isolation, no nested teams, and no resumption of in-process teammates.

An external orchestrator on a subscription therefore cannot be built on Claude.
Vibe stands as the first build's agent, for the reason
[first-build.md](first-build.md) already gives.

## Instruction control by agent

ACP carries no system prompt. `session/new` takes `cwd`,
`additionalDirectories`, `mcpServers`, `_meta`, and `grep -c systemPrompt` over
the v1 and v2 schemas returns 0. ACP v2 states that clients cannot set arbitrary
system prompts. RFD [#1237](https://github.com/agentclientprotocol/agent-client-protocol/issues/1237)
is open and its reference PR #1494 was closed unmerged, so this is a standing gap
rather than a pending fix.

Every agent solves it its own way, which makes persona delivery a per-agent shim:

| Agent | Full replacement by |
| --- | --- |
| **Vibe** | markdown in `~/.vibe/prompts/`, selected by `system_prompt_id` in `config.toml` |
| Claude Code CLI | `--system-prompt`, `--system-prompt-file` |
| Claude Agent SDK | a bare `systemPrompt` string — *"the SDK sends only what you provide"* |
| Claude ACP bridge | `_meta.systemPrompt` as a string, replacing the `claude_code` preset (`agentclientprotocol/claude-agent-acp`, `src/acp-agent.ts`) |
| Codex CLI | `model_instructions_file` |
| OpenCode | agent `system` field |
| Kilo Code | agent `prompt` field |
| Cline | `-s` / `--system` |
| Goose | *contested* — read both as `system_prompt_override` and as append-only under `# Additional Instructions:`. Unresolved. |

Vibe's mechanism is cleaner than the `VIBE_DEFAULT_AGENT` environment variable
[`spikes/acp`](../../spikes/acp/) used.

## Resumption

`session/load` is an optional v1 capability, gated on
`agentCapabilities.loadSession`. A client must not call it when the agent does not
advertise it.

Advertised true by Vibe (`vibe/acp/agent.py:504`), the Claude agent, Gemini CLI,
Goose and Kiro CLI. False by the community `codex-acp` bridge.

ACP v2 removes `session/load` and replaces it with `session/resume`, which becomes
required with no capability marker. Design against resume semantics; keep
`session/load` as the v1 path.

## Permission requests are not a boundary

The spec says an agent **may** request permission. An agent that never asks is
compliant. Claude Code resolves `bypassPermissions` before its permission
callback, so a client sees nothing for those calls.

Adjudication over ACP is therefore observability. The containment is the worktree,
which is what [open-questions.md](open-questions.md#where-the-trust-boundary-sits)
already concludes.

## What exists

Orchestrating parallel coding agents over worktrees is a populated category.

| Project | Holds |
| --- | --- |
| [`Sannrox/bugyo`](https://github.com/Sannrox/bugyo) | ACP client, worktree per session, queue, diff-review-and-merge, human approvals. `docs/acp-notes.md` is a Phase-0 ACP spike with resume and capability probes. |
| [`s0lda/hadron`](https://github.com/s0lda/hadron) | Worktree pool, automated rebase-and-test merge gate. |
| [Compozy](https://github.com/compozy/compozy) | ACP CLIs, SQLite sessions, worktrees, approvals as durable objects. Parallel mode falls back to sequential. |
| [Groundcrew](https://github.com/ClipboardHealth/groundcrew) | Worktree per task, tmux attach, bounded queue (`maximumInProgress`). Exits to a PR. |
| `Untrivial-ai/agent-orchestrator` | Detached per-session hosts: a daemon restart reconnects to an in-flight turn. |
| Amp | Ships to the base branch by default — rebases, tests, pushes, and has the agent resolve conflicts and push again. |

Merge on finish is not novel: hadron gates it, Amp does it. Amp also runs the
re-apply-on-conflict behaviour that
[open-questions.md](open-questions.md#whether-a-part-can-re-apply-its-own-change)
records as untested.

Bounded concurrency is usually money rather than slots. Jules publishes counts
(3/15/60 by tier); Devin, Cursor and Copilot bound by spend. Amp allows a burst of
20 then one new slot every five minutes, queueing rather than failing.

**Dead or moved**, as of this survey: Crystal (deprecated Feb 2026), Terragon
(shut down Feb 2026), Vibe Kanban (company shut down Apr 2026, community
maintained), Roo Code (archived; its `.roo/system-prompt-<mode>` override was
removed before the archive). Products in this category do not last; protocols and
SDKs do.

## Operational findings worth keeping

- A hard-killed session leaves a stale lock and `session/load` then fails with
  *"Session is active in another process (PID N)"*. Release with SIGTERM and wait;
  on resume, delete the lock only when its PID is dead. Consistent with deriving
  liveness rather than storing it.
- Agent-to-client requests may carry string ids where client-to-agent carry
  integers. A JSON-RPC reader must accept both.
- Mastra's `@mastra/acp` auto-selects the first permission option by default. Fail
  open is the wrong default for a decision point.

## What is left to build

Worktree lifecycle, the bounded queue, and merge on finish. No surveyed project
supplies all three together, and no substrate supplies any of them.

## What this changed upstream

- §1's *"delegation inside a session is one level deep"* is Vibe's limit. Claude
  Code nests subagents three levels by default. The consequence — orchestration as
  an external process — holds for Vibe.
- §5 routes permission requests. That is observability, not containment.
