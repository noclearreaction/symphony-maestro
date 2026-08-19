# Fragment 7 — Available subagents

| | |
| --- | --- |
| **Position** | After the skills block |
| **Gate** | `include_prompt_detail` (default `true`) |
| **Source** | Rendered from the agent registry, filtered to `agent_type = "subagent"` |
| **Measured** | 140 chars with one subagent registered |

Tells the model which subagents it may spawn to delegate work. Only profiles
declaring `agent_type = "subagent"` appear here; user-facing agents are never
listed, because the model cannot switch to them.

## Disabling

Only via the shared group gate:

```toml
include_prompt_detail = false     # also removes OS/shell, skills, scratchpad
```

Alternatively, remove entries from the registry itself — a subagent that is not
registered cannot be listed or spawned.

## Content

```markdown
# Available Subagents

The following subagents can be spawned via the Task tool:
- **$name**: $description
```

One bullet per registered subagent, using the profile's `name` and `description`
fields.

## Example

With only the bundled `explore` subagent:

```markdown
# Available Subagents

The following subagents can be spawned via the Task tool:
- **explore**: Read-only subagent for codebase exploration
```

## Adding your own

A subagent is an ordinary profile in `agents/` with `agent_type` set to
`subagent`:

```toml
# .vibe/agents/auditor.toml
agent_type = "subagent"
system_prompt_id = "auditor"
display_name = "Auditor"
description = "Read-only subagent that verifies claims against evidence"
```

It then appears in this fragment and becomes spawnable. Two consequences worth
knowing:

- Subagents cannot be selected by the user. `--agent auditor` is rejected with
  *"is a subagent and cannot be used as the primary agent"*.
- Subagent invocations route through the `task` tool, so hooks match them with
  `match = "task"` and discriminate by reading `tool_input.agent`.

Subagents inherit the parent's hook configuration, and their events are logged to
the subagent's own session log rather than surfacing in the parent's UI.
