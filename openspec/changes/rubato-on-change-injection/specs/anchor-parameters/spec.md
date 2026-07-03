## ADDED Requirements

### Requirement: Anchor block accepts optional parameters array
The `rubato:anchor` JSON SHALL accept an optional top-level `parameters` array. Each element is an object. Unknown keys within parameter objects SHALL be ignored. Absence of `parameters` SHALL be treated as equivalent to an empty array.

#### Scenario: Anchor with parameters array
- **WHEN** the anchor JSON contains `"parameters":[{"max_age":50}]`
- **THEN** `Block.Parameters` SHALL contain one entry with `max_age: 50`

#### Scenario: Anchor without parameters
- **WHEN** the anchor JSON does not contain a `parameters` key
- **THEN** `Block.Parameters` SHALL be nil or empty with no error

#### Scenario: Unknown parameter keys are ignored
- **WHEN** a parameter object contains an unrecognised key
- **THEN** parsing SHALL succeed and the unknown key SHALL be accessible via `Block.Parameters`

### Requirement: MaxAge is derived from parameters with a default of 100
The anchor package SHALL expose a `MaxAge() int` helper (or equivalent) on `Block` that returns the `max_age` value from the first parameters entry, defaulting to 100 when absent or when `parameters` is empty.

#### Scenario: max_age present
- **WHEN** `parameters[0]["max_age"]` is a positive integer
- **THEN** `MaxAge()` SHALL return that value

#### Scenario: max_age absent
- **WHEN** the `parameters` array is absent or empty
- **THEN** `MaxAge()` SHALL return 100

#### Scenario: max_age zero means always inject
- **WHEN** `parameters[0]["max_age"]` is 0
- **THEN** `MaxAge()` SHALL return 0, and the caller SHALL treat this as always-inject
