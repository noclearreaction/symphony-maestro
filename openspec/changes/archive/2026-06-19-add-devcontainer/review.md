# Review: add-devcontainer

**Reviewer**: GitHub Copilot  
**Date**: 2026-06-12 (initial) · 2026-06-13 (updated)  
**Artifacts reviewed**: `proposal.md`, `design.md`, `tasks.md`, `specs/devcontainer-environment/spec.md`

---

## Round 2 — 2026-06-13

### What was addressed

All significant findings from the initial review have been resolved:

| Finding | Resolution |
|---|---|
| Open question about `postStartCommand` / git hooks scope | Closed via D9: minimal scope — `task devcontainer:doctor` only; future hooks left as commented stubs; `env-lgc` tag pattern rejected (no use case, caused first-start failures) |
| Task 5.5 unbounded scope | Constrained: exact script content specified in the task |
| Task 5.7 "to be determined" | Explicitly deferred: empty array now, follow-on change when extension needs are observed |
| Spec requirement paragraphs restating build layout (paths, stage names) | Go and Deno requirements rewritten to behavioral; stage names removed from Dockerfile requirement |
| Trivial `devcontainer.json` scenario | Replaced with "all expected tools present after container opens" |
| Testing gap — no coverage of `post-start` script | Task 7.6 added |
| `gofmt` slave had no verification task | Task 7.4 updated, task 7.5 retained |
| No automated runtime check | D10 + Group 8 added: `doctor` Taskfile task runs on container start and exits non-zero on failure |

---

### Updated Verdict

**Ready to implement.** No blocking issues remain.

Two minor observations follow; neither should block work.

---

### Remaining Minor Issues

**1. ~~`gofmt` slave scenario references implementation mechanism in the WHEN clause~~ — withdrawn**

The scenario invokes `update-alternatives --set` and `update-alternatives --display` directly. On a Debian-based container, `update-alternatives` is the Debian standard for managing versioned binaries — not an incidental implementation choice. A plain `ln -sf` would satisfy the behavioral outcome but would bypass the standard, making the versioned layout fragile and the shared `/opt` volume path harder to exploit later. The spec is entitled to opine on standards compliance, not just on function. No change needed.

**2. ~~Task requirement paragraph retains an installation constraint~~ — withdrawn**

> Task SHALL NOT be installed via npm or pnpm.

On reflection, this is a security constraint: installing Task via npm would exercise pnpm as an install vector, expanding the Node surface area that the PATH isolation policy is designed to contain. It belongs in the spec alongside the PATH isolation requirement. No change needed.

---

## Round 1 — 2026-06-12 (retained for record)

### Verdict

Ready to implement with two issues that should be resolved before or during implementation:

1. Close the open question in `design.md` about `postStartCommand` / git hooks scope, or constrain task 5.5 to a minimal no-op stub.
2. Resolve task 5.7 (VS Code extensions list) or explicitly defer it with a note so it is not an unbounded implementation decision.

The spec should be refactored to separate behavioral requirements from structural/implementation constraints (see below).

### Completeness

Structurally complete. Proposal, design, tasks, and spec are present and internally consistent.

**Unresolved items:**

- **Open question (`design.md`)**: "Does `postStartCommand` need to wire git hooks... and if so, what does that script contain?" Task 5.5 then instructs creating that script. The question should be closed before implementation begins; otherwise 5.5 has an unbounded scope.
- **Task 5.7**: VS Code extensions are "to be determined". Either resolve the list now or mark the task deferred with a named follow-on change.

No tool versions are pinned in any artifact — task 1.2 defers that to implementation time, which is acceptable, provided the Dockerfile ARG comment block becomes the canonical record.

### Sufficiency

Sufficient. The 8 keyed design decisions cover all contested choices with rationale and rejected alternatives documented. Tasks are granular, dependency-ordered, and traceable to spec requirements.

### Testing Strategy

**Coverage**: adequate for a manual first pass. Group 7 (Verification) covers the important cases: tool presence/absence on PATH, DooD socket access, `update-alternatives` state, and shebang portability.

**Gaps:**

- All verification is manual and one-shot. There is no automated test harness and no CI job that runs `docker build --target ci` on each PR. This is consistent with the change's own framing (CI integration is future work), but it means the spec scenarios have no runner and verification is not a recurring signal.
- Task 5.5 (`post-start` script) has no corresponding entry in the verification group. Once the script exists, nothing checks its behavior.
- The `gofmt` slave scenario (spec §2) has no verification task in group 7.

### Spec: Behavior vs. Implementation Detail

This is the most significant issue. The spec mixes behavioral requirements (good) with structural implementation constraints (belongs in `design.md`).

**Implementation detail currently in spec requirements:**

> Stages SHALL be: `download-base`, `go-runtime`, `deno-runtime`, `task-binary`, `node-runtime`, `base`, `ci`, `final`

> Go SHALL be installed at `/opt/go<VERSION>/`. `update-alternatives` SHALL register `/usr/local/bin/go` pointing at `/opt/go<VERSION>/bin/go`...

These specify *how it is built*. The design document already owns this rationale (D1–D5). Restating it in the spec means changes to the build layout would require editing both documents and could cause divergence.

**Behavioral scenarios (the strong part):**

> WHEN a developer runs `pnpm` in a container shell, THEN the shell returns "command not found"

> WHEN `bin/commit-lint.ts` is executed inside the devcontainer, THEN it invokes Deno successfully without a "No such file or directory" error

These are testable, implementation-agnostic, and correctly capture intent. They are the spec's best content.

**Recommendation**: Trim requirement paragraphs to capability-level statements — what the system enables for a developer or CI runner — and let `design.md` own structural decisions. A spec whose requirement paragraphs describe outcomes (not paths, stage names, or command sequences) will remain valid if the build layout is revised.

### Minor Notes

- `devcontainer.json` scenario ("VS Code detects the devcontainer configuration") is trivially satisfied by the file's existence — it is not a runtime behavior. Consider replacing it with a scenario that verifies post-open tooling state.
- "Scenario: Dockerfile ARGs enumerate all tool versions" is a static analysis check, not a behavioral scenario. It is useful, but consider phrasing it as a linting or review gate rather than a Gherkin scenario.
