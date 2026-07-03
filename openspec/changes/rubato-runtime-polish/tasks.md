## 1. Operational Logging Refinement

- [ ] 1.1 Implement structured logs for decision path events: anchor detection, declared plugins, execution outcomes.
- [ ] 1.2 Ensure non-debug logs avoid unnecessary prompt/body content leakage.
- [ ] 1.3 Validate error logs preserve actionable failure context.

## 2. Runtime Configuration Hardening

- [ ] 2.1 Add and validate timeout bounds for plugin execution.
- [ ] 2.2 Document default values and failure behavior for invalid config.
- [ ] 2.3 Add or update `.opencode` configuration to route model traffic through Rubato in devcontainer workflows.
