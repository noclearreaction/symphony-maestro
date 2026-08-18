# Fragment 9b — Additional working directories

| | |
| --- | --- |
| **Position** | Between [project context](09-project-context.md) and [AGENTS.md](10-agents-md.md) |
| **Gate** | `include_project_context` **and** at least one project root other than the cwd |
| **Source** | Built inline — [`system_prompt.py` L403–L415](https://github.com/mistralai/mistral-vibe/blob/v2.24.1/vibe/core/system_prompt.py#L403-L415) |
| **Measured** | Absent in this repo — no extra roots configured |

Emitted only when `harness_files.project_roots` contains a root that is not the
resolved working directory; a single-root session skips it entirely. Extra roots
come from `--add-dir`, so this fragment is the prompt-side counterpart of that
flag.

## Source

```python
cwd_resolved = cwd.resolve()
extra_roots = [
    root
    for root in harness_files.project_roots
    if root.resolve() != cwd_resolved
]
if extra_roots:
    dirs_lines = "\n".join(f" - {d}" for d in extra_roots)
    sections.append(
        "Additional working directories (treated with the same "
        "file-access permissions as the primary working directory):\n"
        + dirs_lines
    )
```

No template file and no substitution — the string is built inline.

## Content

```markdown
Additional working directories (treated with the same file-access permissions as the primary working directory):
 - $root
 - $root
```

## Example

Running `vibe --add-dir /srv/shared --add-dir /opt/reference` from `/workspace`:

```markdown
Additional working directories (treated with the same file-access permissions as the primary working directory):
 - /srv/shared
 - /opt/reference
```

## Permission note

The wording is not decorative. Extra roots are described to the model as carrying
**the same file-access permissions as the primary working directory**, so
`--add-dir` widens the boundary that the file-tool permission chain treats as
"inside the workdir" — step 5 of the resolution order documented in
[scratchpad](08-scratchpad.md#permission-interaction).

Adding a directory therefore has a permission consequence, not just a discovery
one. A path under an added root no longer requires the outside-workdir approval
it would otherwise trigger.
