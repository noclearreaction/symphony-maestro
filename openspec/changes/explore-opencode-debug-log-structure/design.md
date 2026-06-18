## Context

SF-1 (#45) delivered a working Docker environment with opencode installed and a minimal fixture project. The container is confirmed functional.

The parent spike (#43) was written with assumptions about what opencode exposes: specific field names like `tokens_input`, `tokens_cache_read`, `cost`, an `opencode db` command, and debug log structure. None of these have been verified. SF-2 exists specifically to ground-truth those assumptions before building any measurement tooling.

This is a discovery task, not a build task. The output is a findings document, not code.

## Goals / Non-Goals

**Goals:**
- Run a turn inside the harness container with debug logging enabled and capture raw output
- Document the complete structure of the debug log: what sections appear, what fields are present, example values
- Discover what `opencode db` (or equivalent) exposes: whether the command exists, what tables/columns are present, what per-turn data is available
- Identify exactly where (if anywhere) cache token counts, input token counts, and cost appear
- Note any gaps between #43's assumptions and reality
- Produce a committed findings note that explicitly states what measurement approach SF-3–SF-8 should use

**Non-Goals:**
- Writing any measurement scripts (that is SF-3+)
- Automating the discovery process
- Exploring more than one model or provider during this session
- Modifying the harness Dockerfile (that is SF-1 scope)

## Decisions

### Output format: committed markdown findings note

The findings are captured as a markdown file committed to the repository, not as an issue comment or ephemeral notes.

**Rationale**: Issue comments are hard to reference from tasks, easy to lose in notification noise, and cannot be updated with structured diffs. A committed file is reviewable, versioned, and directly referenceable by SF-3–SF-8 planning artifacts.

**Alternative considered**: Posting a comment on #46 — acceptable for visibility but insufficient as a durable reference artifact.

### Scope: one provider, one model, one turn

The discovery session runs a single turn against the default model configured in the fixture, with debug logging enabled.

**Rationale**: The goal is to understand the log and db schema, not to compare across models. A single controlled turn is sufficient to observe the full structure. Expanding scope risks producing noisy findings that obscure the basic schema.

### Discovery method: manual exploration inside the container

The experimenter enters an interactive shell, runs opencode commands manually, and records observations.

**Rationale**: This is intentionally unscripted. We do not know yet what commands are valid or what output looks like. Scripting before discovery would embed the assumptions we are trying to test.

## Risks / Trade-offs

- **opencode behavior may differ from documentation** → The findings document is the authority; ignore prior assumptions if they conflict.
- **Debug log may be large or hard to parse** → Focus on identifying structure and representative field names; do not try to document every line.
- **`opencode db` may not exist or may behave differently than assumed** → Document what actually exists; explicitly call out the gap if the command is absent or different.
- **Free model may have rate limits that interrupt the session** → Keep turns minimal; this is a one-turn discovery, not a load test.

## Open Questions

- Does `opencode db` exist as a subcommand? If not, where is the SQLite database stored and can it be queried directly with `sqlite3`?
- Are token counts reported in the debug log, the database, or both?
- Is cost reported at all, or is it calculated client-side from token counts and a price table?
