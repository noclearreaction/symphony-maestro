## ADDED Requirements

### Requirement: Anchor block accepts optional options array
The `rubato:anchor` JSON SHALL accept an optional top-level `options` array of `{"name", "setting"}` objects. Parsers SHALL scan all entries for recognised `name` values and ignore entries with unknown names. Absence of `options` SHALL be treated as equivalent to an empty array.

#### Scenario: Anchor with options array
- **WHEN** the anchor JSON contains `"options":[{"name":"max_age","setting":50}]`
- **THEN** `Block.Options` SHALL contain one entry with name `max_age` and setting `50`

#### Scenario: Anchor without options
- **WHEN** the anchor JSON does not contain an `options` key
- **THEN** `Block.Options` SHALL be nil or empty with no error

#### Scenario: Unknown option names are ignored
- **WHEN** an options entry has an unrecognised `name` value
- **THEN** parsing SHALL succeed and the entry SHALL be preserved in `Block.Options`

### Requirement: MaxAge is derived from the options array with a default of 100
The anchor package SHALL expose a `MaxAge() int` method on `Block` that scans `Options` for an entry where `name == "max_age"` and returns its `setting` as int, defaulting to 100 when not found.

#### Scenario: max_age present
- **WHEN** `parameters[0]["max_age"]` is a positive integer
- **THEN** `MaxAge()` SHALL return that value

#### Scenario: max_age absent
- **WHEN** the `parameters` array is absent or empty
- **THEN** `MaxAge()` SHALL return 100

#### Scenario: max_age zero means always inject
- **WHEN** `parameters[0]["max_age"]` is 0
- **THEN** `MaxAge()` SHALL return 0, and the caller SHALL treat this as always-inject
