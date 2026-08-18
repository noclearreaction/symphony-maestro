# Fragment 9 — Project context

| | |
| --- | --- |
| **Position** | After the scratchpad block |
| **Gate** | `include_project_context` (default `true`) |
| **Source** | Bundled template `_internal/vibe/core/prompts/project_context.md` |
| **Measured** | 470 chars |

Working directory plus a snapshot of git state. This is one of the few fragments
with a template that exists on disk as a `.md` file.

## Disabling

```toml
# config.toml
include_project_context = false
```

```bash
VIBE_INCLUDE_PROJECT_CONTEXT=false vibe ...
```

Disabling removes 5,711 chars in total, because this flag **also gates the entire
[AGENTS.md fragment](10-agents-md.md)** — 5,239 of that total. If the goal is a
leaner prompt, this is the highest-value flag; if the goal is keeping project
instructions while dropping git noise, this flag cannot do that, because the two
are not separable.

## Content

The bundled template, using `string.Template` placeholders:

```markdown
Absolute path: $abs_path

gitStatus: This is the git status at the start of the conversation. Note that this status is a snapshot in time, and will not update during the conversation.
$git_status
```

`$git_status` expands to current branch, main branch, working-tree status, and
the five most recent commits with their subjects truncated to the first word.

## Example

```markdown
Absolute path: /workspace

gitStatus: This is the git status at the start of the conversation. Note that this status is a snapshot in time, and will not update during the conversation.
Current branch: change/tdd-orchestration-intent
Main branch (you will usually use this for PRs): main
Status: (clean)
Recent commits:
132582f (HEAD -> change/tdd-orchestration-intent) feat
b723ae7 feat
eccbda2 docs
9847592 docs
57e0095 docs(design): add acceptance validation strategy
```

Note that commit subjects are truncated — `feat(vibe): add project-local Vibe
configuration` appears as just `feat`. The oldest entry retains its full subject,
so the truncation is not uniform.

## The snapshot is stale by design

The fragment states it plainly: the status is captured once at session start and
never refreshed. In the example above `Status: (clean)` was true at startup; any
commit or edit made during the session leaves this fragment unchanged.

The model therefore cannot rely on it. Anything acting on current git state should
run `git status` rather than trust this block — the same caveat applies to a model
reading it as to a human.

## Related bundled template

`dangerous_directory.md` (316 chars) replaces this fragment when project context
scanning is suppressed:

```markdown
directoryStructure: Project context scanning has been disabled because $reason. This prevents permission dialogs and potential system slowdowns. Use the LS tool and other file tools to explore the project structure as needed.

Absolute path: $abs_path

gitStatus: Use git tools to check repository status if needed.
```

This substitution is triggered by Vibe's own dangerous-directory detection, not by
a config key.

## Not user-overridable

Placing a file at `~/.vibe/prompts/project_context.md` or
`.vibe/prompts/project_context.md` has **no effect**, in either scope.

The reason is structural rather than a special case. `PROJECT_CONTEXT` is a
member of the `UtilityPrompt` enum, and reading an enum member resolves straight
to its bundled path. The search path lives in a different function,
`load_prompt(setting_name, …)`, which resolves *a config setting's value* — so a
prompt is overridable only if some setting names it. Exactly two do:
`system_prompt_id` and `compaction_prompt_id`. Nothing names `PROJECT_CONTEXT`.

Note this is not a system-vs-utility distinction: `COMPACT` is also a
`UtilityPrompt` and *is* overridable, because `compaction_prompt_id` names it.

The `$abs_path` / `$git_status` placeholders are interpolated internally and are
not a user-facing templating feature — a custom base prompt containing
`$abs_path` reaches the model with that literal text.

## How it is rendered

[`ProjectContextProvider.get_full_context`](https://github.com/mistralai/mistral-vibe/blob/v2.24.1/vibe/core/system_prompt.py#L155-L162)
reads the bundled `UtilityPrompt.PROJECT_CONTEXT` and substitutes two names:

```python
template = UtilityPrompt.PROJECT_CONTEXT.read()
return Template(template).safe_substitute(
    abs_path=..., git_status=...
)
```

The values come from `_fetch_git_status` → `_parse_git_log` →
`_format_git_status`, which shell out to git in a thread pool with timeouts and
cap the result at `MAX_GIT_STATUS_SIZE`.

This is one of four substitution sites in the assembly — the others being your
base prompt (`current_date`), `DANGEROUS_DIRECTORY` (`reason`, `abs_path`) and
`AGENTS_DOC` (`sections`). Each mapping is local to its call site rather than
coming from a registry, so an overridable version of this file still could not
reference a new variable: nothing would supply it, and `safe_substitute` renders
unknown placeholders as literal text.
