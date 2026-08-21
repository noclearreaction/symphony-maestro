# First build

What the first build picks, in what order, and the entry points it exposes.
Companion to [README](README.md).

## Choices

The design is unspecific about agents, plugs and entry points. The first
build has to pick. These are picks, not commitments — each is the cheapest thing
that exercises the design.

| Concern | First build | Because |
| --- | --- | --- |
| ACP agent | **Vibe**, via `vibe-acp` | Installed, routed through rubato, funded by a subscription we hold. Verified in [`spikes/acp`](../../spikes/acp/). |
| Policy plug | a static allow/deny list from configuration | Proves the decision point is wired without deciding anything about policy. |
| Entry point | **CLI**, including the `tail`-shaped observer | Validates the control API, and forces an event stream in early. |
| Language | **TypeScript on Deno** | The ACP SDK is first-class there, and it matches the repository. |
| Sections | **One** | A second section proves the boundary, which is not the first thing worth knowing. |

### Not in the first build

- **Containers per turn.** The trust boundary stays where ACP puts it, bounding
  permission but not reach.
- **Priority, budgets, section reservation.** The queue is FIFO over ready work
  with a chair count.
- **TUI and VS Code.** The CLI proves the API first.
- **A second agent vendor.** The adapter seam exists so adding one is not a
  rewrite.

## Build order

Each stage is independently useful, and each is a place to stop.

| Stage | Delivers | Proves |
| --- | --- | --- |
| 0 ✅ | **ACP spike** — [`spikes/acp`](../../spikes/acp/). | Done. A real turn, routed through rubato, no headless instruction, the persona ours. |
| 1 | `session/load` on the spike. | Resume works over ACP — the assumption every review round depends on. |
| 2 | Core plus CLI: start a part, attach to it, `tail` it. | The entry point, and an event stream in the control API. |
| 3 | Personas from configuration; `personas`, `parts`, `start`, `status`. | Work can be started and survives a restart. |
| 4 | `wait`, `result`, the queue, chairs, a worktree per part. | Work returns, bounded by chairs. |
| 5 | `resume`, `accept`, merge and re-apply on conflict, attach in `AwaitingReview`. | The review loop, and that a part merges or reverts cleanly. |
| 6 | Route permission requests to a decision point, with one plug behind it. | A part can use tools with no human present. |
| 7 | A second section. | Sections are a first-class boundary. |


Attach is the one capability entry points cannot implement uniformly:

| Entry point | Attach means |
| --- | --- |
| CLI | exec the agent against the session in the foreground |
| TUI | suspend itself and hand over the terminal, or embed a pty |
| VS Code | open the agent's own extension against the session |


### CLI first

The CLI validates the control API without UI work confounding the question: if the
CLI can drive the ensemble, the API is complete. It remains the scripting surface
afterwards.

A CLI polls and exits, and a poll-shaped CLI produces a poll-shaped API. A TUI
needs a live view, which is a subscription. The first CLI therefore includes the
`tail`-shaped observer, which forces an event stream into the API early; the TUI
then becomes a rendering problem.

The TUI is the eventual daily driver. VS Code comes last: where an agent ships its
own extension, Ensemble's part is a section panel that launches an attach.

