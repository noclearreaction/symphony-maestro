## Context

`AGENTS.md` currently serves as both an immediate operational context file for agents and a loose placeholder for future system aspirations (specifically the evolutionary path from Symphony v1 bootstrap autonomy to Symphony v2 constrained artifact-driven orchestration). This mixing of active system state with speculative targets violates Principle 3 of the Symphony Constitution (Intent Traceability) and workspace hygiene guidelines. To ensure agents operate under precise, unambiguous, and active context, the speculative targets must be removed from `AGENTS.md` and tracked formally as structured GitHub Issues.

## Goals / Non-Goals

**Goals:**
- Trim `AGENTS.md` of future-oriented speculative targets (specifically, the Symphony v1 and v2 bootstrap guidelines under the `### Symphony` section).
- Extract and preserve these targets as individual, structured GitHub Issues categorized according to the `governance/issue-lifecycle.md` guidelines.
- Ensure `AGENTS.md` remains high-fidelity, compliant, and syntax-correct, serving strictly as a marker of active system status and operational context.

**Non-Goals:**
- Modifying the core principles, constitutional rules, or collaboration guardrails inside `AGENTS.md`.
- Modifying any of the three active specification workflows under `openspec/specs/` (e.g., `controlled-git-workflows`, `issue-workflow`, `director-workflow`).
- Initiating any code implementation or functional changes to Symphony or Symphony Director.

## Decisions

### Decision 1: Precise modifications to `AGENTS.md`

We will modify `AGENTS.md` by conducting a section-by-section review and trimming down the future-oriented Symphony v1/v2 speculative bootstrap targets.

- **Current Section under `### Symphony` (Lines 34–52):**
  ```markdown
  ### Symphony

  Symphony is the orchestration system being built.

  Symphony owns orchestration mechanics:

  * event-driven work intake
  * work item lifecycle
  * agent execution contracts
  * human-in-the-loop gates
  * review flows
  * PR and issue interaction
  * test and evidence handling
  * project workflow automation

  Symphony v1 may be bootstrapped with more AI autonomy to discover the shape of the problem.

  Symphony v2 should be more constrained, specified, test-driven, and artifact-driven.
  ```

- **Target Modifications:**
  We will remove the following lines (Lines 49–52) entirely:
  ```markdown
  Symphony v1 may be bootstrapped with more AI autonomy to discover the shape of the problem.

  Symphony v2 should be more constrained, specified, test-driven, and artifact-driven.
  ```

- **Alternative Considered:**
  Keep the lines but comment them out or label them as speculative.
  - *Rationale for Rejection*: Retaining dead or commented-out speculative content in active agent context files leads to context pollution and could still confuse or bias LLM execution paths. Moving them completely out of `AGENTS.md` is cleaner.

### Decision 2: GitHub Issues tracking the extracted targets

The speculative targets will be extracted and tracked systematically as two distinct GitHub Issues to ensure that this intent is not lost and is fully traceable under Constitutional Principle 3.

#### Issue 1: Symphony v1 bootstrapping with AI autonomy
- **Title**: `Symphony v1 bootstrapping with AI autonomy`
- **Labels**: `type:spike`, `status:backlog`, `priority:medium`
- **Description**:
  ```markdown
  Track the bootstrapping of Symphony v1. The initial implementation of Symphony v1 may be bootstrapped with higher AI autonomy to discover the shape of the problem space, explore interactive developer-agent workflows, and gather telemetry on agent interaction.
  ```

#### Issue 2: Symphony v2 transition to constrained artifact-driven orchestration
- **Title**: `Symphony v2 transition to constrained artifact-driven orchestration`
- **Labels**: `type:spike`, `status:backlog`, `priority:medium`
- **Description**:
  ```markdown
  Track the strategic design and implementation of Symphony v2. Transition the orchestration engine from the bootstrap v1 model to a highly constrained, specified, test-driven, and artifact-driven system to guarantee deterministic execution, predictability, and complete auditability.
  ```

### Decision 3: Verification Strategy

To ensure `AGENTS.md` remains high-fidelity, compliant, and syntax-correct:

1. **Section Readability Verification**: Verify that the resulting `### Symphony` section transitions cleanly into the subsequent `### Project Repositories` section without syntax errors or grammatical awkwardness.
2. **Markdown Lint Check**: Verify that the document retains valid markdown structure, proper spacing, and no trailing whitespaces or broken lists.
3. **Reference Verification**: Verify that other workspace files and specifications (such as `openspec/specs/controlled-git-workflows/spec.md`) that reference `AGENTS.md` for operational context are not broken or misaligned by this change.

## Risks / Trade-offs

- **[Risk] Lost Intent / Context Drift** → *Mitigation*: Ensure the GitHub Issues are created successfully with clear labels and detailed descriptions before `AGENTS.md` is modified, so no structural design context is lost.
- **[Risk] Broken Reference in Specs** → *Mitigation*: Verify that existing references to `AGENTS.md` are still valid. Since `AGENTS.md` will still exist and retain all active system rules (such as controlled git and conventional commits), reference paths remain unchanged and fully functional.
