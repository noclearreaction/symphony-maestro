## ADDED Requirements

### Requirement: rplugin accepts a plugin name and executes it
The `rplugin` binary SHALL accept a plugin name as the first positional argument, execute that plugin, and write its output to stdout.

#### Scenario: Known plugin runs successfully
- **WHEN** `rplugin git_status` is invoked in a git repository
- **THEN** the plugin output SHALL be written to stdout
- **THEN** the process SHALL exit with code 0

#### Scenario: Unknown plugin name
- **WHEN** `rplugin no-such-plugin` is invoked
- **THEN** an error message SHALL be written to stderr
- **THEN** the process SHALL exit with code 1

#### Scenario: No plugin name provided
- **WHEN** `rplugin` is invoked with no arguments
- **THEN** a usage message SHALL be written to stderr
- **THEN** the process SHALL exit with code 1

### Requirement: rplugin accepts --working-dir flag
The `rplugin` binary SHALL accept a `--working-dir` flag and pass its value as `working_dir` in the plugin args.

#### Scenario: --working-dir flag provided
- **WHEN** `rplugin git_status --working-dir /some/path` is invoked
- **THEN** the plugin SHALL be executed with `args["working_dir"] = "/some/path"`

### Requirement: rplugin accepts --args flag for additional plugin args
The `rplugin` binary SHALL accept an `--args` flag containing a JSON object. Keys from `--args` SHALL be merged into the plugin args, with `--working-dir` taking precedence over a `working_dir` key in `--args`.

#### Scenario: --args flag with valid JSON
- **WHEN** `rplugin go_test --args '{"timeout_seconds":30}'` is invoked
- **THEN** the plugin SHALL receive `args["timeout_seconds"] = 30`

#### Scenario: --args flag with invalid JSON
- **WHEN** `--args` is provided with malformed JSON
- **THEN** an error message SHALL be written to stderr
- **THEN** the process SHALL exit with code 1

### Requirement: Plugin errors are reported to stderr
When plugin execution returns an error, `rplugin` SHALL write the error message to stderr and exit with code 1.

#### Scenario: Plugin execution fails
- **WHEN** a plugin returns an error (e.g. git_status in a non-repo directory)
- **THEN** the error SHALL be written to stderr
- **THEN** the process SHALL exit with code 1
- **THEN** nothing SHALL be written to stdout
