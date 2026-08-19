Project-local Mistral Vibe configuration for this repository.

Vibe only reads this directory when the folder is **trusted** (recorded in
`$VIBE_HOME/trusted_folders.toml`, default `~/.vibe/trusted_folders.toml`), and
it is discovered relative to the directory Vibe is launched from. Run Vibe from
the repository root, or pass `--workdir /workspace`, or these files are silently
ignored.

## Layout

| Path | Purpose |
| --- | --- |
| `config.toml` | Project settings (model, tools, permissions). Overrides `~/.vibe/config.toml`. |
| `hooks.toml` | Project hook definitions. Requires a trusted folder. |
| `agents/` | Agent profiles — one `*.toml` per agent. See [agents/README.md](agents/README.md). |
| `prompts/` | System prompts — the `*.md` files that agent profiles reference. |
| `skills/` | Custom skills, one subdirectory per skill. See [skills/README.md](skills/README.md). |
| `tools/` | Custom tools (`<name>.py`), plus `tools/prompts/<name>.md` to set or override a tool description. |

Note that `agents/` holds only profile metadata; the prompt text lives in
`prompts/`. The two are linked by the profile's `system_prompt_id`.

## Scope precedence

User-level configuration lives in `$VIBE_HOME` (default `~/.vibe/`) and applies
everywhere. Project-level configuration lives here and takes precedence for
`config.toml`. Agents, prompts and skills resolve from both scopes, so a profile
in `agents/` may reference a prompt in `~/.vibe/prompts/` and vice versa.

Vibe writes no logs into this directory. Logging is user-scoped: `vibe.log` and
per-session transcripts go to `$VIBE_HOME/logs/` (here, the `mistral-vibe` Docker
volume mounted at `~/.vibe`).

## Reference

- Agents: https://docs.mistral.ai/vibe/code/cli/agents
- Approvals and permissions: https://docs.mistral.ai/vibe/code/safety-approvals-permissions
- The running CLI also documents itself: ask Vibe about its own configuration and
  it loads a built-in self-awareness skill covering every key in `config.toml`.
