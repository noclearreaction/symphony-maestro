# Fragment 1 — Base prompt

| | |
| --- | --- |
| **Position** | First. Everything else is appended after it. |
| **Gate** | `system_prompt_id` — always present, but its content is fully replaceable |
| **Source** | A `.md` file resolved through the prompt search path |
| **Measured** | 538 chars (`reviewer.md`); 9,234 chars for the bundled default |

This is the only replaceable slot. Setting `system_prompt_id` to a custom prompt
removes the bundled default **entirely** — no part of it survives.

## Resolution

`system_prompt_id` names a prompt ID, resolved in order (first match wins):

1. `<project>/.vibe/prompts/<id>.md`
2. `~/.vibe/prompts/<id>.md`
3. bundled `_internal/vibe/core/prompts/<id>.md`

Set it in an agent profile, which is the per-role approach:

```toml
# .vibe/agents/reviewer.toml
system_prompt_id = "reviewer"     # -> .vibe/prompts/reviewer.md
```

Or globally in `config.toml`, which changes the default for every agent that does
not override it:

```toml
system_prompt_id = "cli"
```

The default value `cli` comes from the `vibe_cli_system_prompt` feature flag
(`defaultValue: "cli"`), visible under `experiments.features` in any session's
`meta.json`.

## Bundled base prompts

| ID | Chars | Notes |
| --- | --- | --- |
| `cli_2026-07_v2` | 10,371 | |
| `cli` | 9,234 | Current default |
| `lean` | 8,023 | Lean 4 proof assistant |
| `explore` | 1,472 | Used by the `explore` subagent |
| `minimal` | 419 | Terse; drops tool-usage and instruction-hierarchy guidance |
| `tests` | 52 | |

Switching `cli` → `minimal` cuts ~8.8KB, but it removes behavioural guidance
rather than just text. Expect different behaviour, not merely a smaller prompt.

## Content

The file is inserted with one template substitution applied —
[`_interpolate_prompt`](https://github.com/mistralai/mistral-vibe/blob/v2.24.1/vibe/core/system_prompt.py#L323-L324):

```python
def _interpolate_prompt(prompt: str) -> str:
    return Template(prompt).safe_substitute(current_date=_format_current_date())
```

So exactly one variable is available in your own prompt, and it is the only
user-reachable dynamic value anywhere in the system prompt:

| Variable | Expands to |
| --- | --- |
| `$current_date` | `2026-08-18 (Tuesday)` |

Anything else is left as literal text, because `safe_substitute` does not raise
on unknown placeholders. A prompt containing `$abs_path` reaches the model with
that literal string — the project-context variables are a different mapping at a
different call site.

```markdown
<the contents of the resolved .md file, with $current_date substituted>
```

## Example

With `system_prompt_id = "reviewer"` resolving to `.vibe/prompts/reviewer.md`:

```markdown
# Reviewer Agent

You are a code reviewer. Your role is to critically examine code changes and provide feedback.

## Instructions

- Focus on code quality, correctness, and maintainability
- Check for consistency with existing patterns in the codebase
- Verify tests cover the changes adequately
- Point out potential bugs, edge cases, and performance issues
- Suggest improvements to readability and structure
- Do not approve your own work

## Scope

Review diffs, pull requests, and code changes. Provide concise, actionable feedback.
```

The bundled `cli` prompt, by contrast, opens:

```markdown
You are Mistral Vibe, a CLI coding agent built by Mistral AI. You work on a local codebase using tools.
Today's date is 2026-08-18 (Tuesday).

## Instruction hierarchy

When instructions conflict, resolve in this order (lowest number wins):

1. Critical instructions (never overridable)
2. User messages (more recent messages override older ones)
3. Repo AGENTS.md files — all files on the path from the task files up to
the repo root are active; closer to the task wins on conflict
4. The user's AGENTS.md
```

With `system_prompt_id = "reviewer"`, none of that text is present.
