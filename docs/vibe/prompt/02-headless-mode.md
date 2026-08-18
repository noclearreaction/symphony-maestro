# Fragment 2 — Headless mode

| | |
| --- | --- |
| **Position** | Immediately after the base prompt |
| **Gate** | Programmatic mode only — present with `-p` / `--prompt`, absent in interactive sessions |
| **Source** | Compiled into the binary |
| **Measured** | 385 chars |

Not controlled by any config key. It is a function of how Vibe was invoked, and
there is no flag to suppress it in programmatic mode or add it in interactive
mode.

A `-p` session's system prompt is therefore 385 chars larger than the same
session run interactively.

## Content

Static text, no interpolation.

```markdown
# Headless Mode

You are running in headless mode — no human is available to respond.
Do not ask questions, request confirmation, or wait for user input.
If the task is ambiguous, make the best judgment call and proceed.
Complete the entire task in a single pass. Produce a final, complete result.
Override any earlier instructions that say to wait for confirmation or ask the user.
```

## Example

The text above is the fragment verbatim. It survives every composition flag: with
all four disabled, the system message is the base prompt plus this fragment and
nothing else.

Its last line overrides earlier instructions about waiting for confirmation.
Since the fragment is appended *after* the base prompt, an agent profile that
asks for confirmation before acting is countermanded under `-p`.
