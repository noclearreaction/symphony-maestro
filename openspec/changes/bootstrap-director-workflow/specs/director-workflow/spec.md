## ADDED Requirements

### Requirement: Builder Agent Soft Boundaries
The custom `builder` agent system prompt SHALL enforce strict soft boundaries that prevent any direct modifications to strategy documents or other agent prompts.

#### Scenario: Refusing strategy modification
- **WHEN** the builder agent is requested to edit files under `strategy/`
- **THEN** the builder agent SHALL refuse the change and prompt the user to consult the Advisor or Designer agent instead.

### Requirement: Structured Local Workspace Isolation
The repository SHALL ignore a designated `.symphony/` directory structure to allow agents to write scratch records, plans, and evidence without polluting git history.

#### Scenario: Isolating temporary logs
- **WHEN** an agent writes files under `.symphony/scratchpad/`
- **THEN** the file SHALL not appear as an untracked file in git status.

### Requirement: Branch Initiation Automation
The system SHALL provide a script `bin/director-start` to automate checking out a standardized feature branch and instantiating an OpenSpec change.

#### Scenario: Running start script with valid name
- **WHEN** the user runs `bin/director-start "improve-auth"`
- **THEN** the script SHALL check out branch `change/improve-auth` and run `openspec new change "improve-auth"`.

### Requirement: GitHub PR and Sync Automation
The system SHALL provide a script `bin/director-submit` to synchronize specifications, commit changes, push to the remote, and open a GitHub Pull Request.

#### Scenario: Submitting active change with GitHub CLI
- **WHEN** the user runs `bin/director-submit` on a feature branch
- **THEN** the script SHALL execute `openspec sync`, commit the synchronized specification updates, push the branch, and run `gh pr create` to open a pull request.
