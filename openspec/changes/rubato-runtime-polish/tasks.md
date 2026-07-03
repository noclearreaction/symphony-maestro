## 1. Operational Logging Refinement

- [x] 1.1 Implement structured logs for decision path events: anchor detection, declared plugins, execution outcomes.
- [x] 1.2 Ensure non-debug logs avoid unnecessary prompt/body content leakage.
- [x] 1.3 Validate error logs preserve actionable failure context.

## 2. Runtime Configuration Hardening

- [x] 2.1 Add and validate timeout bounds for plugin execution.
- [x] 2.2 Document default values and failure behavior for invalid config.
- [x] 2.3 Add or update `.opencode` configuration to route model traffic through Rubato in devcontainer workflows.
