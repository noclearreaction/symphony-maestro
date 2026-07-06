## ADDED Requirements

### Requirement: plugins is an array of descriptor objects
The `rubato:anchor` JSON `plugins` field SHALL be an array of objects, each with a `plugin` string (the plugin name) and an optional `options` array of `{name, setting}` pairs. The `setting` field within an option is itself optional, allowing flag-style options with no value. A descriptor with no `options` field is valid.

#### Scenario: Plugin with no options
- **WHEN** the anchor JSON contains `"plugins":[{"plugin":"git_status"}]`
- **THEN** `Block.Plugins` SHALL contain one descriptor with name `git_status` and `Options` defaulting to an empty list

#### Scenario: Plugin with options
- **WHEN** the anchor JSON contains `"plugins":[{"plugin":"go_test","options":[{"name":"timeout_seconds","setting":30}]}]`
- **THEN** the descriptor for `go_test` SHALL have one option with name `timeout_seconds` and setting `30`

#### Scenario: Option without setting (flag-style)
- **WHEN** an options entry contains only `{"name":"verbose"}` with no `setting` field
- **THEN** the option SHALL be parsed with `setting` as nil/absent with no error

#### Scenario: Multiple plugins
- **WHEN** the anchor JSON contains two plugin descriptors
- **THEN** `Block.Plugins` SHALL preserve declaration order

### Requirement: Anchor accepts top-level options array for rubato-level config
The `rubato:anchor` JSON SHALL accept an optional top-level `options` array of `{name, setting}` objects. Parsers SHALL scan all entries for recognised names; unknown names are ignored. Absence of `options` is equivalent to an empty array.

#### Scenario: Top-level options present
- **WHEN** the anchor JSON contains `"options":[{"name":"max_age","setting":50}]`
- **THEN** `Block.Options` SHALL contain one entry with name `max_age` and setting `50`

#### Scenario: Top-level options absent
- **WHEN** the anchor JSON has no `options` field
- **THEN** `Block.Options` SHALL default to an empty list with no error

#### Scenario: Unknown option names are ignored
- **WHEN** an options entry has an unrecognised `name`
- **THEN** parsing SHALL succeed and the entry SHALL be preserved in `Block.Options`

### Requirement: MaxAge is derived from top-level options with a default of 100
The anchor package SHALL expose a `MaxAge() int` method on `Block` that scans `Options` for an entry where `name == "max_age"` and returns its `setting` as int, defaulting to 100 when not found; returning 0 when explicitly set to 0.

#### Scenario: max_age present
- **WHEN** `options` contains `{"name":"max_age","setting":50}`
- **THEN** `MaxAge()` SHALL return 50

#### Scenario: max_age absent
- **WHEN** `options` is absent or contains no `max_age` entry
- **THEN** `MaxAge()` SHALL return 100

#### Scenario: max_age zero means always inject
- **WHEN** `options` contains `{"name":"max_age","setting":0}`
- **THEN** `MaxAge()` SHALL return 0
