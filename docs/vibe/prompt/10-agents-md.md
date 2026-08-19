# Fragment 10 — AGENTS.md

| | |
| --- | --- |
| **Position** | Last |
| **Gate** | `include_project_context` (default `true`) — shared with [project context](09-project-context.md) |
| **Source** | A fixed preamble, then the contents of every discovered `AGENTS.md` |
| **Measured** | 5,239 chars = preamble 526 + wrapper 227 + this repo's 4,486-char `AGENTS.md` |

The largest fragment, and the only one whose content is entirely yours. It is
also the only sanctioned way to inject your own material into the system prompt.

## Gating is shared

`include_project_context` removes this entire fragment along with the git
snapshot — 5,711 chars combined, of which this is 5,239.

There is **no separate key** for AGENTS.md. You cannot keep project instructions
while dropping the stale git snapshot, or vice versa. There is also no per-agent
filtering — every agent, including subagents, receives an identical copy.

## Discovery

- `~/.vibe/AGENTS.md` — user level.
- Every `AGENTS.md` from the project root up through the trust chain, read **once
  at session start**.
- `AGENTS.md` in subdirectories — injected **lazily mid-session**, when
  `read_file` reads a file below them. Any `AGENTS.md` between the file's parent
  and the project root is added to context at that point.

Precedence, as stated in the preamble: project beats user, and among project
files, closer to the working directory wins.

## Content

Rendered from the bundled `UtilityPrompt.AGENTS_DOC` template, whose single
`$sections` variable is filled with the concatenated document sections
([L417–L435](https://github.com/mistralai/mistral-vibe/blob/v2.24.1/vibe/core/system_prompt.py#L417-L435)):

```python
if user_doc.strip():
    doc_sections.append(f"## User instructions\n\nContents of {VIBE_HOME.path}/AGENTS.md ...")
if project_docs:
    doc_sections.append("## Project instructions (checked into the codebase)")
for doc_dir, doc_content in project_docs:
    doc_sections.append(f"Contents of {doc_dir}/AGENTS.md:\n\n{doc_content.strip()}")
```

A user-level `~/.vibe/AGENTS.md` appears as its own `## User instructions`
section ahead of the project ones. The example below shows the project case only.

A fixed preamble, then one titled section per discovered file:

```markdown
Codebase and user instructions are shown below. Be sure to adhere to these instructions. IMPORTANT: These instructions OVERRIDE any default behavior and you MUST follow them exactly as written. When both user-level and project-level instructions are present, project instructions take priority over user instructions. When multiple project-level AGENTS.md files are present, instructions closer to the working directory take priority. Each AGENTS.md applies to its own directory and all of its descendants within the project.

## Project instructions (checked into the codebase)

Contents of $path:

$file_contents

IMPORTANT: this context may or may not be relevant to your tasks. You should act on these guidelines if they are relevant to your task.
```

## Example

Abbreviated — the real fragment carries all 4,486 chars of the file:

```markdown
Codebase and user instructions are shown below. Be sure to adhere to these instructions. IMPORTANT: These instructions OVERRIDE any default behavior and you MUST follow them exactly as written. ...

## Project instructions (checked into the codebase)

Contents of /workspace/AGENTS.md:

# Symphony Director Project Context

This repository represents the Symphony Director project.
...

IMPORTANT: this context may or may not be relevant to your tasks. You should act on these guidelines if they are relevant to your task.
```

Note the preamble and the closing line pull in opposite directions: the opening
says the instructions **OVERRIDE any default behavior** and must be followed
**exactly as written**, while the closing line says they **may or may not be
relevant**. Both framings reach the model in every session.

## Using it for dynamic content

`AGENTS.md` is an ordinary file, so "dynamic" means generating it before Vibe
starts:

```bash
# a wrapper, a Taskfile target, or a devcontainer postCreate step
{
  cat AGENTS.base.md
  echo
  echo "## Current state"
  echo "- open change: $(ls openspec/changes 2>/dev/null | head -1)"
} > AGENTS.md
vibe "$@"
```

Constraints to design around:

- **Read once per session** for the root chain. Regenerating mid-session has no
  effect until restart.
- **No include or import mechanism.** Composition has to happen before launch, as
  above.
- **No per-agent scoping.** Content here reaches every role identically. Material
  that should differ by role belongs in `prompts/<role>.md`, which *is* per-agent.
- **Subdirectory files are the one genuinely mid-session lever** — placing
  `AGENTS.md` in a subtree gets it injected when that subtree is first read.

## Per-turn alternatives

For content that must change every turn, hooks can inject into other channels —
though not into the system prompt:

| Hook | Mechanism | Where it lands |
| --- | --- | --- |
| `post_tool` | `hook_specific_output.additional_context` | Appended to the tool result the model sees |
| `post_agent` | `decision: "deny"` with a `reason` | Injected as a retry user message, max 3 per hook per turn |

See [../../../.vibe/hooks.toml](../../../.vibe/hooks.toml) for the full wire
protocol.
