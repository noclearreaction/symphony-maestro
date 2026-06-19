## Why

The parent spike (#43) makes assumptions about what data opencode exposes — specific field names, log sections, and database columns — but none of those assumptions have been verified against a running instance. Before writing any measurement scripts (SF-3 through SF-8), we need to know what data is actually available and where.

## What Changes

- A structured discovery session run inside the Docker harness (from SF-1 / #45)
- A findings note documenting the real debug log structure: sections, field names, example values
- A findings note documenting the `opencode db` schema: tables, columns, per-turn data available
- An explicit mapping (or gap report) between what #43 assumed and what actually exists
- Updates to SF-3–SF-8 issue descriptions if the measurement approach needs to change

## Capabilities

### New Capabilities

- `opencode-observability-findings`: A documented findings artifact capturing what opencode actually exposes via debug logs and its database, serving as the ground-truth reference for all downstream SF work

### Modified Capabilities

## Impact

- No code changes — this is a discovery and documentation task
- Output is a findings note (committed to the repo or posted as a comment on the issue)
- Directly unblocks SF-3 through SF-8 by providing verified field names and data availability
- May require updates to the designs of SF-3–SF-8 if assumptions in #43 are incorrect
