## ADDED Requirements

### Requirement: Deterministic Operational Context
The project-wide `AGENTS.md` instruction file SHALL serve strictly as an active marker of system status and operational context, aligned with the Symphony Constitution and Vision, and SHALL NOT contain speculative or future-oriented system requirements.

#### Scenario: Verify AGENTS.md content
- **WHEN** the `AGENTS.md` file is loaded by any agent or validator
- **THEN** it SHALL be verified to contain only active operational rules, with all future-oriented speculative targets extracted to GitHub issues.
