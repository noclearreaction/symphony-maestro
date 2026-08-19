Custom skills for Mistral Vibe.

A skill is a specialized instruction set the model loads on demand. Each skill is
a **directory** containing a `SKILL.md` with YAML frontmatter, plus any supporting
files it references.

```
.vibe/skills/
└── my-skill/
    └── SKILL.md
```

```markdown
---
name: my-skill
description: What this skill does and when to use it.
user-invocable: true
allowed-tools: bash read
---

# Skill Instructions

Detailed instructions for the model...
```

## Invocation

- The model loads a skill on demand via the `skill` tool.
- Users invoke a `user-invocable` skill by typing `/my-skill`, optionally followed
  by extra instructions. Typing `/` at the start of the input opens the slash menu.

Skills with `user-invocable: false` are model-only: hidden from the slash menu,
and `/my-skill` will not resolve them — though the model can still load them via
the `skill` tool.

## Search order (first match wins)

1. `skill_paths` from `config.toml`
2. `.vibe/skills/` (this directory, trusted project only)
3. `.agents/skills/` (trusted project only)
4. `~/.vibe/skills/`
5. `~/.agents/skills/` (Agent Skills standard)

Project-level entries require the folder to be trusted and are resolved relative
to the launch directory, so run Vibe from the repository root or pass `--workdir`.
