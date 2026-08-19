# Fragment 6 — Available skills

| | |
| --- | --- |
| **Position** | After the OS and shell line |
| **Gate** | `include_prompt_detail` (default `true`) |
| **Source** | Rendered from the resolved skill registry |
| **Measured** | 1,937 chars with two skills registered |

Lists every enabled skill with its description so the model knows what it can
load on demand. The second-largest fragment, and it grows with each skill you
add — the whole `description` field of every skill lands in the system prompt of
every session.

## Disabling

Whole fragment, via the shared group gate:

```toml
include_prompt_detail = false     # also removes OS/shell, subagents, scratchpad
```

Or trim it selectively without losing the other three fragments, by reducing the
registry itself:

```toml
disabled_skills = ["experimental-*"]
enabled_skills = ["vibe"]         # allowlist; everything else drops out
```

This is the better lever: skill descriptions are the part that scales, and
`disabled_skills` shrinks the fragment without discarding the scratchpad and
subagent information.

## Content

A fixed preamble, then one `<skill>` element per registered skill. Descriptions
are XML-escaped, so apostrophes appear as `&#x27;` and quotes as `&quot;`.

```markdown
# Available Skills

You have access to the following skills. When a task matches a skill's description,
use the `skill` tool if available to load the full skill instructions, if it is not available, read the files manually if they exist.

When a user message is exactly `/skill-name` (optionally followed by extra
instructions), the user has explicitly invoked that skill. Its instructions are
loaded for you automatically: you will see a `skill` tool call and result
immediately after that message. Treat the loaded content as the active
instructions and act on it — you do not need to call the `skill` tool yourself.

<available_skills>
  <skill>
    <name>$skill_name</name>
    <description>$skill_description</description>
  </skill>
  ...
</available_skills>
```

## Example

With only the two built-in skills registered, abbreviated at the second
description:

```markdown
<available_skills>
  <skill>
    <name>skill-creator</name>
    <description>Load this skill before creating, updating, or deleting a Vibe skill, or whenever you plan to add or modify SKILL.md files under a skills directory. It explains Vibe&#x27;s SKILL.md frontmatter, discovery order, support files, and the permission flow.</description>
  </skill>
  <skill>
    <name>vibe</name>
    <description>Authoritative reference for Mistral Vibe — the CLI agent you (the model) are running inside.

LOAD when the user:
- asks anything about Vibe itself, even by indirect name (&quot;this CLI&quot;, &quot;this tool&quot;, &quot;you&quot;);
...
SCOPE: config under `~/.vibe/` and project-local `.vibe/`; ... logs; themes; voice.</description>
  </skill>
</available_skills>
```

The `vibe` skill is the CLI's own self-documentation, and it is the authoritative
source for the config keys referenced throughout this directory. Its full body
lives at `_internal/vibe/core/skills/builtins/vibe.py`.

Because descriptions are reproduced in full, a multi-paragraph description like
`vibe`'s costs over a kilobyte in every session. Keep skill descriptions short
enough to route on, and put the detail in `SKILL.md` where it is loaded on demand.
