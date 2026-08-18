Agent profiles for Mistral Vibe.

An agent profile configures the model's behavior, tools and system prompt. This
directory holds **profile metadata only** — one `*.toml` per agent. The prompt
text lives in `../prompts/`, and the two are linked by `system_prompt_id`.

## Layout

```
.vibe/
├── agents/
│   └── reviewer.toml     # profile: metadata, model, tool controls
└── prompts/
    └── reviewer.md       # system prompt, resolved via system_prompt_id
```

The agent's name is the `.toml` filename stem — `reviewer.toml` is selected as
`--agent reviewer`. There is no requirement that a `.md` file sit beside the
`.toml`; a prompt file in `agents/` is not read at all.

Profiles and prompts resolve independently from both scopes, so a profile here
may reference a prompt in `~/.vibe/prompts/`, and a profile in `~/.vibe/agents/`
may reference a prompt in `../prompts/`.

## Minimal profile

```toml
agent_type = "agent"              # "agent" (user-facing) or "subagent" (delegation-only)
system_prompt_id = "reviewer"     # resolves to prompts/reviewer.md
active_model = "mistral-medium-3.5"
display_name = "Reviewer"         # shown in the UI
description = "Code review specialist agent"
```

Optional keys:

```toml
safety = "neutral"                # "safe" | "neutral" | "destructive" | "yolo"
enabled_tools = ["read_file", "grep"]
disabled_tools = ["edit", "write_file"]

[tools.bash]
permission = "never"              # "ask" | "always" | "never"
```

> **An invalid value silently drops the profile.** Vibe does not report a
> validation error — the agent simply stops existing, and selecting it fails with
> `Error: Agent '<name>' not found.` If a profile you just wrote is "not found",
> suspect a bad field value before suspecting scope or trust.

## Selecting an agent

- At launch: `vibe --agent reviewer`
- Mid-session: `Shift+Tab` cycles through available agents; `/config` also sets it.
  The active agent is shown in the status line.
- By default: `default_agent = "reviewer"` in `config.toml`.

There is no command that lists agents — no `vibe agents list`, no `/agents`.
`Shift+Tab` cycling is the only listing UI. To verify a profile loads without
starting a session:

```bash
vibe --workdir "$PWD" --agent reviewer -p "reply OK" --max-turns 1
```

Selection is discovered relative to the launch directory and requires a trusted
folder, so run from the repository root or pass `--workdir`.

## Built-in agents

| Agent | Behavior |
| --- | --- |
| `ask` | Requests approval for every tool execution |
| `plan` | Planning-focused; auto-approves read-only tools |
| `accept-edits` | Default; auto-approves file edits, asks for other tools |
| `auto-approve` | Auto-approves all tool calls |
| `lean` | Lean 4 proof assistant. Opt-in — `/leanstall`, or `installed_agents = ["lean"]` |

Subagents are model-facing only: the model spawns them to delegate subtasks, and
they cannot be selected by the user. `explore` is a subagent, so
`--agent explore` is rejected.

## Related config keys

```toml
agent_paths = ["/path/to/shared/agents"]   # additional search paths
enabled_agents = ["ask", "plan"]
disabled_agents = ["auto-approve"]
installed_agents = ["lean"]                # opt-in builtins
default_agent = "plan"
```

`agent_paths` is the cleanest way to mount shared or team-wide agents, since it
avoids writing into `$VIBE_HOME`.

## VS Code extension

The extension lists agents that have a `.toml` profile; a prompt file alone will
not register one. Its exact discovery rules are not verified here — the layout
above is what the CLI reads.

Reference: https://docs.mistral.ai/vibe/code/cli/agents
