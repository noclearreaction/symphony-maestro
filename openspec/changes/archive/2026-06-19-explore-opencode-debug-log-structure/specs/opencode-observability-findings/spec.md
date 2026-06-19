## ADDED Requirements

### Requirement: Debug log structure is documented
The discovery session SHALL produce a findings note that documents the full structure of the opencode debug log captured from a single turn: what sections appear, what field names are present, and representative example values.

#### Scenario: Debug log fields are recorded
- **WHEN** a turn is run inside the harness container with debug logging enabled
- **THEN** the findings note contains the field names and structure observed in the debug log output

#### Scenario: Gaps are noted
- **WHEN** a field assumed in #43 (e.g. `tokens_input`, `tokens_cache_read`) does not appear in the debug log
- **THEN** the findings note explicitly records that the field was not found

### Requirement: Database schema is documented
The discovery session SHALL produce a findings note that documents what opencode exposes via its database interface (whether `opencode db`, direct SQLite access, or another mechanism): what tables and columns exist, and what per-turn data is available.

#### Scenario: Schema is recorded
- **WHEN** the experimenter accesses opencode's database after a turn
- **THEN** the findings note lists the tables and columns observed, with example values where available

#### Scenario: Missing command is noted
- **WHEN** `opencode db` does not exist as a subcommand
- **THEN** the findings note records this explicitly and documents the alternative access method used (e.g. direct SQLite query)

### Requirement: Measurement approach is stated
The findings note SHALL include an explicit statement of the recommended measurement approach for SF-3–SF-8, based on what was actually observed — not on prior assumptions.

#### Scenario: Approach is actionable
- **WHEN** a reader of the findings note is planning SF-3
- **THEN** they can identify which specific fields or commands to use for measuring token counts and cost without needing to re-run the discovery

#### Scenario: Assumptions that failed are called out
- **WHEN** an assumption from #43 does not hold (e.g. a field does not exist)
- **THEN** the findings note states the assumption, states what was found instead, and recommends an alternative or notes the gap

### Requirement: Findings are committed to the repository
The findings note SHALL be committed as a versioned file in the repository so it is referenceable by downstream planning artifacts.

#### Scenario: Findings file is present after the session
- **WHEN** the discovery session is complete
- **THEN** a markdown findings file exists in the repository at a documented path and is included in the change commit
