# Fragment 8 — Scratchpad directory

| | |
| --- | --- |
| **Position** | Last of the `include_prompt_detail` group |
| **Gate** | `include_prompt_detail` (default `true`) |
| **Source** | Rendered from a per-session generated path |
| **Measured** | 329 chars |

Tells the model where to put temporary files, and that writes there skip the
permission prompt.

## Disabling

Only via the shared group gate:

```toml
include_prompt_detail = false     # also removes OS/shell, skills, subagents
```

Disabling this fragment does not remove the scratchpad or its permission
exemption — it only stops telling the model the path exists. The directory is
still created and still auto-approved. Suppressing it therefore makes the model
less likely to use a facility it still has, rather than tightening anything.

## Content

```markdown
# Scratchpad Directory

You have a scratchpad directory at: `$scratchpad_path`

Use this for temporary files: intermediate results, draft scripts, working files, outputs that don't belong in the project.
Files here are automatically allowed — no permission prompts.
Session-scoped. Shared with subagents.
```

The path is `/tmp/vibe-scratchpad-<session-prefix>-<random-suffix>`, generated
per session.

## Example

```markdown
# Scratchpad Directory

You have a scratchpad directory at: `/tmp/vibe-scratchpad-d3a8fc7a-wex3gf47`

Use this for temporary files: intermediate results, draft scripts, working files, outputs that don't belong in the project.
Files here are automatically allowed — no permission prompts.
Session-scoped. Shared with subagents.
```

The prefix matches the session ID, so the directory can be correlated with a
session log:

```bash
ls -dt ~/.vibe/logs/session/session_*d3a8fc7a*/
```

## Permission interaction

The scratchpad is first in the file-tool permission resolution order, ahead of
the denylist:

1. **Scratchpad path → always allowed**
2. denylist glob → always denied
3. allowlist glob → always allowed
4. `sensitive_patterns` → requires approval
5. outside workdir → requires approval
6. default → the tool's `permission` setting

Because it outranks the denylist, a denylist pattern cannot protect a path inside
the scratchpad. Since the directory is session-scoped and under `/tmp`, that is
usually the intent, but it does mean the scratchpad is not a place to rely on
deny rules.
