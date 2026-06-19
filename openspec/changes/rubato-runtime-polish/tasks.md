## 1. Operational Logging Refinement

- [ ] 1.1 Implement structured logs for decision path events: anchor detection, declared plugins, execution outcomes.
- [ ] 1.2 Ensure non-debug logs avoid unnecessary prompt/body content leakage.
- [ ] 1.3 Validate error logs preserve actionable failure context.

## 2. Runtime Configuration Hardening

- [ ] 2.1 Add and validate timeout bounds for plugin execution.
- [ ] 2.2 Add and validate repo-path/runtime config constraints.
- [ ] 2.3 Document default values and failure behavior for invalid config.
- [ ] 2.4 Add or update `.opencode` configuration to route model traffic through Rubato in devcontainer workflows.

## 3. Task 7 Contract Freeze

- [ ] 3.1 Capture implementation traces for anchor parsing, plugin failures, and git status mappings.
- [ ] 3.2 Select canonical wire-level error and response shapes from traces.
- [ ] 3.3 Update design/spec/review artifacts to reflect only validated behaviors.

## 4. End-To-End Verification

- [ ] 4.1 Run devcontainer-routed end-to-end requests through Rubato and verify expected outcomes.
- [ ] 4.2 Run regression suite to ensure no Stage A/B behavior regressions.
- [ ] 4.3 Record final verification evidence and close remaining hygiene items.