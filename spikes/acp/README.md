# ACP spike

> **Status:** Throwaway. Stage 0 of the Ensemble build order
> ([`docs/ensemble/README.md`](../../docs/ensemble/README.md) §8). Kept for the
> findings below, not as code to build on.

Drives a Vibe player over the [Agent Client Protocol](https://agentclientprotocol.com)
and reports what the contract gives back.

```bash
cd spikes/acp
deno task spike                                    # default agent
deno task spike --agent reviewer --prompt "..."    # a specific profile
```

## What it answers

| Question | Answer |
| --- | --- |
| Does ACP run a real turn against the installed `vibe-acp`? | **Yes.** `initialize` → `session/new` → `session/prompt` → `stopReason: end_turn`. |
| Does traffic route through rubato? | **Yes.** The proxy logs one `proxying POST` per turn. |
| Is the headless fragment absent? | **Yes**, on both runs. Confirms on the installed binary what was previously only seen via the VS Code extension's copy. |
| Do custom instructions land? | **Yes**, via `VIBE_DEFAULT_AGENT` at spawn — see below. |

## Findings

### Instructions are selected per process, not per session

`vibe/acp/agent.py` builds `SessionOptions` without an `agent` parameter, so the
protocol offers no way to choose a profile. Configuration is the only lever, and
`VIBE_`-prefixed environment variables override it — so a player's instructions
are fixed when its process is spawned.

That suits Ensemble: one `vibe-acp` process per player, `VIBE_DEFAULT_AGENT` set
from the player's role. It also means a player cannot change role mid-session,
which is a constraint worth designing around rather than discovering later.

Measured, same trivial prompt both times:

| Profile | System prompt | Tokens | First line |
| --- | --- | --- | --- |
| `accept-edits` (default) | 19,457 chars | 7,689 | *"You are Mistral Vibe, a CLI coding agent…"* |
| `reviewer` | 10,840 chars | 4,443 | `# Reviewer Agent` |

The profile is doing more than swapping the base prompt — 8,617 chars and 3,246
tokens separate them on an identical one-word exchange. Worth understanding
before assuming a player's floor cost, since it is paid on every turn.

### `vibe-acp` does not shut down cleanly

Killing the child after a turn produces a Python traceback
(`ValueError: I/O operation on closed file`) and a Sentry flush on the way out.
Cosmetic here; less so for a supervisor that starts and stops players constantly.
Two things follow: Ensemble should close the stream and let the agent exit rather
than killing it, and **`vibe-acp` emits telemetry to Sentry**, which a fleet
would do at fleet volume.

### Deno needs no scaffolding

`npm:@agentclientprotocol/sdk@1.3.0` resolves straight from the import map — no
`package.json`, no `node_modules`, no build step. `deno check` type-checks the
npm package's own types.

## Not covered

- `session/load`, so resuming a session over ACP is still unverified. This is
  what tier 2 depends on.
- Permission requests. Nothing in these runs triggered one, so the auto-approve
  handler is written but unexercised.
- `vibe-app-server`, the third substrate candidate.
- Anything policy-related. Cedar and Biscuit remain unevaluated.
