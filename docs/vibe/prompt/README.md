# How Vibe builds its system prompt

> **Status:** Observational, not provisional — these are measurements of a
> shipped binary rather than proposals. But they are **pinned to Vibe 2.24.1**
> and will drift: every fragment, gate and character count below is what that
> version does, and an upgrade can change any of it without notice. Re-measure
> before relying on a number.

Vibe sends the model a **single system message**, built by concatenating one
replaceable base prompt with a series of appended fragments.

Everything here describes Vibe 2.24.1. The assembly is one function,
[`get_universal_system_prompt`](https://github.com/mistralai/mistral-vibe/blob/v2.24.1/vibe/core/system_prompt.py#L352-L437),
with one `sections.append` per fragment and a final `"\n\n".join(sections)`.
Prompt resolution lives in
[`vibe/core/prompts/__init__.py`](https://github.com/mistralai/mistral-vibe/blob/v2.24.1/vibe/core/prompts/__init__.py).

## Assembly order

Sizes are from a `-p` run with `--agent reviewer` on this repo, 9,533 chars
total. Most fragments scale with the project, so treat them as proportions.

| # | Fragment | Gate | Chars |
| --- | --- | --- | --- |
| 1 | [Base prompt](01-base-prompt.md) | `system_prompt_id` | 538 |
| 2 | [Headless mode](02-headless-mode.md) | programmatic (`-p`) mode only | 385 |
| 3 | [Commit signature](03-commit-signature.md) | `include_commit_signature` | 393 |
| 4 | [Model info](04-model-info.md) | `include_model_info` | 41 |
| 5 | [OS and shell](05-os-shell.md) | `include_prompt_detail` | 52 |
| 6 | [Available skills](06-skills.md) | `include_prompt_detail` | 1,937 |
| 7 | [Available subagents](07-subagents.md) | `include_prompt_detail` | 140 |
| 8 | [Scratchpad](08-scratchpad.md) | `include_prompt_detail` | 329 |
| 9 | [Project context](09-project-context.md) | `include_project_context` | 470 |
| 9b | [Additional working directories](09b-additional-working-dirs.md) | `include_project_context`, and extra roots exist | absent here |
| 10 | [AGENTS.md](10-agents-md.md) | `include_project_context` | 5,239 |

## Controlling what is included

Five keys in `config.toml`, all enabled by default:

```toml
system_prompt_id = "cli"          # base prompt: "cli", "lean", or a custom .md name
include_commit_signature = true   # fragment 3
include_model_info = true         # fragment 4
include_prompt_detail = true      # fragments 5-8
include_project_context = true    # fragments 9-10
```

Each accepts a `VIBE_`-prefixed environment override:

```bash
VIBE_INCLUDE_PROJECT_CONTEXT=false vibe --agent reviewer -p "..."
```

Effect of disabling each:

| Flag set to `false` | Chars removed | Removes |
| --- | --- | --- |
| `include_project_context` | 5,711 | project context **and all AGENTS.md** |
| `include_prompt_detail` | 2,454 | OS/shell, skills, subagents, scratchpad |
| `include_commit_signature` | 394 | commit signature |
| `include_model_info` | 42 | model name |

`include_project_context` gates AGENTS.md as well as the git snapshot. There is
no separate key for the two.

With all four off, the system message is the base prompt alone — 921 chars here,
being 538 of prompt plus the 385-char headless fragment that `-p` adds.

## Base prompt replacement

Setting `system_prompt_id` replaces the bundled default outright; none of it
survives. A prompt ID resolves in order:

1. `<project>/.vibe/prompts/<id>.md`
2. `~/.vibe/prompts/<id>.md`
3. bundled `_internal/vibe/core/prompts/<id>.md`

Only two settings resolve this way: `system_prompt_id` and
`compaction_prompt_id`. Other bundled templates — `PROJECT_CONTEXT`,
`AGENTS_DOC`, `DANGEROUS_DIRECTORY` and the rest of the `UtilityPrompt` enum —
are read directly from the bundle and cannot be overridden, because no setting
names them. This is not a system-vs-utility split: `COMPACT` is a `UtilityPrompt`
and is overridable, because `compaction_prompt_id` names it.

## Dynamic content

Your base prompt is run through template substitution with a single variable:

| Variable | Expands to |
| --- | --- |
| `$current_date` | `2026-08-18 (Tuesday)` |

```markdown
Today is $current_date.        ->  Today is 2026-08-18 (Tuesday).
Path is $abs_path.             ->  Path is $abs_path.
```

Unknown placeholders are left as literal text rather than raising.

The assembly has four substitution sites, each with a mapping local to its call
site. There is no variable registry, so no new variable can be introduced:

| Site | Template | Variables |
| --- | --- | --- |
| [L159](https://github.com/mistralai/mistral-vibe/blob/v2.24.1/vibe/core/system_prompt.py#L159) | `PROJECT_CONTEXT` | `abs_path`, `git_status` |
| [L324](https://github.com/mistralai/mistral-vibe/blob/v2.24.1/vibe/core/system_prompt.py#L324) | your base prompt | `current_date` |
| [L393](https://github.com/mistralai/mistral-vibe/blob/v2.24.1/vibe/core/system_prompt.py#L393) | `DANGEROUS_DIRECTORY` | `reason`, `abs_path` |
| [L434](https://github.com/mistralai/mistral-vibe/blob/v2.24.1/vibe/core/system_prompt.py#L434) | `AGENTS_DOC` | `sections` |

For anything richer, compute it outside Vibe and write it where Vibe reads:

| Want | Use | Refresh point |
| --- | --- | --- |
| Current date in your prompt | `$current_date` | Every session |
| Project-wide instructions | `AGENTS.md` | Read at startup |
| Directory-scoped instructions | `AGENTS.md` in a subdirectory | Injected when a file below it is read |
| Per-turn generated content | Hooks | Every turn, but into tool results or user messages, not the system prompt |

`AGENTS.md` is an ordinary file, so generating it from a wrapper script,
Taskfile target or devcontainer step is the practical route. The root chain is
read once per session; regenerating mid-session has no effect until restart.

Hook channels, from [`.vibe/hooks.toml`](../../../.vibe/hooks.toml):

- `post_tool` → `hook_specific_output.additional_context` appends to the tool
  result the model sees.
- `post_agent` → `decision: "deny"` injects `reason` as a retry user message,
  capped at 3 per hook per user turn.

## Inspecting a session

The resolved prompt is persisted per session:

```bash
D=$(ls -dt ~/.vibe/logs/session/session_*/ | head -1)

jq -r '.config.system_prompt_id'      "$D/meta.json"
jq -r '.system_prompt.content'        "$D/meta.json"
jq -r '.system_prompt.content|length' "$D/meta.json"
```

`meta.json` also records `tools_available`, the resolved `agent_profile`, and the
session's active feature flags.
