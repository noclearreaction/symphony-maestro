# Fragment 5 — OS and shell

| | |
| --- | --- |
| **Position** | After the model info |
| **Gate** | `include_prompt_detail` (default `true`) |
| **Source** | Rendered from runtime environment detection |
| **Measured** | 52 chars |

A single line naming the operating system and the shell that the `bash` tool will
use. First of the four fragments sharing the `include_prompt_detail` gate —
see also [skills](06-skills.md), [subagents](07-subagents.md) and
[scratchpad](08-scratchpad.md).

## Disabling

There is no dedicated key. It goes only as part of the group:

```toml
# config.toml
include_prompt_detail = false     # also removes skills, subagents, scratchpad
```

```bash
VIBE_INCLUDE_PROMPT_DETAIL=false vibe ...
```

Disabling the group removes 2,454 chars in total, of which this fragment is 52.

## Content

```markdown
The operating system is Linux with shell `$shell`
```

The shell is taken from the environment, so it reflects the user's login shell
rather than a fixed value.

## Example

In this devcontainer:

```markdown
The operating system is Linux with shell `/bin/zsh`
```

Worth noting that the reported shell is `/bin/zsh` while Vibe's `bash` tool is
named for bash. The shell named here is the one commands are actually executed
with, so shell-specific syntax in generated commands follows this line, not the
tool name. Globbing differences between zsh and bash are the usual way this
surfaces — an unmatched glob is an error in zsh but passes through in bash.
