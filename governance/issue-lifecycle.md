# Issue Lifecycle and Governance

This document establishes the official development lifecycle, state machine, and labeling guidelines for the Symphony Director and Symphony repositories. Adherence to these guidelines ensures system state transparency and supports Constitutional Principle 3 (Intent Traceability).

---

## 1. Issue Categorization Schema

To preserve workspace hygiene, every issue MUST be categorized under exactly one **Type**, one **Status**, and one **Priority** label.

### 1.1 Work Type (`type:*`)
- `type:feature` (Green, `#0E8A16`): Introduction of new functionality, user capabilities, or features.
- `type:bug` (Orange-Red, `#D93F0B`): Repair of unexpected behavior, functional regression, or system failure.
- `type:chore` (Grey, `#EDEDED`): Non-functional tasks including CI/CD configuration, formatting, dependency updates, and maintenance.
- `type:spike` (Yellow, `#FBCA04`): Time-boxed investigation, architecture exploration, prototyping, or strategic planning (preferred over unstructured research).

### 1.2 Development Status (`status:*`)
- `status:backlog` (Grey, `#EDEDED`): Default state for newly created, un-triaged, or unapproved work.
- `status:accepted` (Light Yellow, `#FEF2C0`): Work has been triaged, approved, and prioritized. Ready for implementation.
- `status:in-progress` (Light Yellow, `#FEF2C0`): Active execution phase (e.g., branch checked out, development agent active).
- `status:completed` (Green, `#0E8A16`): The implementation satisfies all specifications, passed local validation, and has been successfully submitted or merged.
- `status:blocked` (Black, `#000000`): Active work is halted due to an external dependency, blocker, or unresolved question.

### 1.3 Priority (`priority:*`)
- `priority:high` (Dark Red, `#B60205`): Critical or blocking issues essential to immediate milestones.
- `priority:medium` (Yellow, `#FBCA04`): Standard prioritized backlog work.
- `priority:low` (Grey, `#EDEDED`): Elective improvements, polish, or minor optimizations.

---

## 2. State Transition Machine

All issues follow a predictable state machine from initial intake to final resolution.

```
+------------------+
|  status:backlog  | <--- Default on issue creation
+------------------+
         |
         v (Triage / Planning Approval)
+------------------+
| status:accepted  |
+------------------+
         |
         +----------------------------------+
         |                                  |
         v (Execution Start)                v (Halt on blocker)
+--------------------+              +------------------+
| status:in-progress | <=========>  |  status:blocked  |
+--------------------+ (Resolve)    +------------------+
         |
         v (Specs Verified & PR Submitted/Merged)
+--------------------+
|  status:completed  |
+--------------------+
```

### 2.1 Backlog State (`status:backlog`)
- **Entrance Criteria**: Automatically applied to all newly created issues.
- **Operator Action**: The `@issue` subagent applies `status:backlog` during drafting and creation.

### 2.2 Accepted State (`status:accepted`)
- **Entrance Criteria**: The issue is approved by the human operator, planned, or added to a milestone.
- **Operator Action**: Transitioned manually by the human or an automated orchestration agent during workspace planning.

### 2.3 In-Progress State (`status:in-progress`)
- **Entrance Criteria**: Implementation has begun. Specifically, when a local feature branch is instantiated (`change/*`) and a developer (human or `builder` agent) begins code modifications.
- **Operator Action**: Transitioned by the builder/orchestrator at the start of implementation.

### 2.4 Blocked State (`status:blocked`)
- **Entrance Criteria**: The implementation cannot proceed due to external factors, missing dependencies, or outstanding design questions.
- **Operator Action**: Transitioned by the operator. The operator MUST add a detailed comment to the issue documenting the blocking condition. Once resolved, the issue returns to `status:in-progress`.

### 2.5 Completed State (`status:completed`)
- **Entrance Criteria**: The change passes all tests and local verification (e.g., `openspec validate`), and a Pull Request is successfully submitted and merged into `main`.
- **Operator Action**: Transitioned automatically or manually as part of the PR merge/archive routine.

---

## 3. Tracing and Auditability Rules

To maintain absolute visibility and auditability:
1. **Comment on Transition**: Every transition into `status:blocked` or out of it MUST be accompanied by a clear commentary explaining the rationale.
2. **Commit Association**: All code commits and PR descriptions MUST reference the associated issue number (e.g., `Resolves #6` or `Addresses #6`).
3. **Specification Sync**: No issue can be transitioned to `status:completed` without all relevant specifications under `openspec/specs/` being synchronized to `main`.
