## 1. Prepare the Session

- [x] 1.1 Confirm the `opencode-cache-harness` Docker image is available (`docker images | grep opencode-cache-harness`)
- [x] 1.2 Create `harness/findings/` directory to hold the output of this session
- [x] 1.3 Create `harness/findings/sf-2-observability.md` as an empty findings file with placeholder sections

## 2. Explore the opencode CLI

- [x] 2.1 Start an interactive shell in the container (`docker run -it opencode-cache-harness`)
- [x] 2.2 Run `opencode --help` and record all available subcommands
- [x] 2.3 Check whether `opencode db` exists as a subcommand; record the result
- [x] 2.4 Locate the SQLite database file path (check `~/.local/share/opencode/` or equivalent) and confirm it is accessible with `sqlite3`

## 3. Run a Single Instrumented Turn

- [x] 3.1 Run one turn using the experiment agent with debug logging: capture both stdout/stderr and any debug output
- [x] 3.2 Save the raw debug log output to a file inside the container (or copy it out)
- [x] 3.3 Note the exact command(s) used to trigger the turn

## 4. Document Debug Log Structure

- [x] 4.1 Identify the top-level sections/phases visible in the debug log
- [x] 4.2 Record all field names observed in the log, especially any related to tokens, cache, or cost
- [x] 4.3 Note which assumed fields from #43 (`tokens_input`, `tokens_cache_read`, `cost`) appear, appear under different names, or are absent
- [x] 4.4 Write these observations into `harness/findings/sf-2-observability.md` under a "Debug Log" section

## 5. Document the Database Schema

- [x] 5.1 After the turn, query the database: list all tables (`SELECT name FROM sqlite_master WHERE type='table'`)
- [x] 5.2 For each relevant table, record the column names and types
- [x] 5.3 Query the most recent session/turn row and record what fields are populated and with what values
- [x] 5.4 Write these observations into `harness/findings/sf-2-observability.md` under a "Database Schema" section

## 6. State the Measurement Approach

- [x] 6.1 Based on findings, write an explicit "Recommended Measurement Approach" section in the findings note: which fields/commands SF-3–SF-8 should use
- [x] 6.2 List any assumptions from #43 that need correction, with a brief note on what to use instead

## 7. Commit and Report

- [x] 7.1 Commit `harness/findings/sf-2-observability.md` with message referencing #46
- [ ] 7.2 Note on issue #46 that findings are committed (or paste a summary as a comment)
