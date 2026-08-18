# Fragment 4 — Model info

| | |
| --- | --- |
| **Position** | After the commit signature |
| **Gate** | `include_model_info` (default `true`) |
| **Source** | Rendered from the resolved active model |
| **Measured** | 41 chars |

A single line telling the model its own name. The smallest fragment; disabling it
is a rounding error on prompt size, but it is the only thing that makes the model
able to answer "which model am I?" accurately.

## Disabling

```toml
# config.toml
include_model_info = false
```

```bash
VIBE_INCLUDE_MODEL_INFO=false vibe ...
```

Disabling removes 42 chars.

## Content

```markdown
Your model name is: `$model_alias`
```

The value is the resolved model **alias** from configuration, not the upstream
model ID. Session `meta.json` records the mapping — for example alias
`mistral-medium-3.5` maps to `config.models."mistral-medium-3.5".name =
mistral-vibe-cli-latest`.

## Example

With the default alias:

```markdown
Your model name is: `mistral-medium-3.5`
```

An agent profile that pins its own model changes this line. `reviewer.toml` sets
`active_model = "mistral-medium-3.5"`, so a `reviewer` session reports that alias
even if `config.toml` leaves `active_model` unset.

Because the alias is reported rather than the upstream ID, the model's
self-description will not match what the API actually served if the alias is
pointed at a moving target such as `mistral-vibe-cli-latest`.
