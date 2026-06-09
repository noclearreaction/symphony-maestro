# Symphony Director Project Context

This repository represents the Symphony Director project.

The Symphony Director is an out-of-tree strategic contact point for the User. It is used to keep track of goals, priorities, decisions, project direction, open questions, and cross-project concerns.

The Symphony Director is not the runtime orchestration system. It should not be required for Symphony or any project repository to function.

## Vision

The project vision lives in `governance/vision.md`. 
Use it as an evaluative benchmark for governance and agent-instruction work, not as a hard rule set.

## Conceptual Layers

### Symphony Director

The Symphony Director owns the User's strategic continuity.

It tracks:

* the User's goals
* the User's priorities
* active decisions
* project direction
* cross-project concerns
* lessons learned while building Symphony
* risks of drift between intent and implementation

The Symphony Director may help decide what work matters and why.

The Symphony Director should not contain product/runtime code.

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

### Project Repositories

Project repositories own project truth.

They contain:

* project goals
* OpenSpec artifacts
* Gherkin scenarios
* tests and contract tests
* application code
* issues
* branches
* PRs
* review evidence
* CI results

A project repository should remain meaningful without the Symphony Director.

## Core Principles

### Human Ownership

The User remains the owner of intent, tradeoffs, approval, and final decisions.

AI agents may propose, draft, review, summarize, and execute bounded work, but they do not own the project.

### Artifact Boundaries

Agents exchange artifacts and evidence, not hidden reasoning or role-played context.

Preferred boundary-crossing artifacts include:

* OpenSpec changes
* Gherkin scenarios
* task descriptions
* diffs
* test output
* contract test results
* review comments
* decision records
* status summaries

An agent should not review its own work from the same context.

### Context Separation

Different roles should receive different context.

Implementation agents need task-local context.

Review agents need artifacts, specs, diffs, and evidence.

Symphony Director-level agents need project status, decisions, goals, risks, and trajectory.

Context should not be mixed merely because the same model is capable of playing multiple roles.

### Explore Before Building

Explore before specifying.

Specify before designing.

Design before implementing.

Do not convert uncertainty into architecture.

Do not produce solution-shaped artifacts before the problem space is understood.

## Collaboration Guardrails & Failure Modes to Avoid

To ensure high-quality collaboration between agents and the User, the system adheres to the following principles:

* **No Premature Architecture**: Avoid inventing canonical schemas or producing adapter designs before roles and problems are fully understood.
* **Grounded Design**: Agents must read and verify existing workspace files before proposing changes or drafting task lists, rather than design in a vacuum.
* **Clear Artifact Promotion**: Drafts, recommendations, and temporary agent files are not treated as project truth until reviewed, approved, and integrated.
* **Independent Contexts**: Agent contexts are kept strictly separate. Do not use one agent's reasoning context as another agent's review context.
* **Task Boundaries**: Ensure agents stay within their defined boundaries (e.g., design agents do not implement code, and implementation agents do not redefine specs).

## Current Bootstrap Strategy

The current goal is to establish a stable advisory and specification workflow before building the full orchestration system.

The early workflow uses OpenCode as an interactive shell. OpenCode is a bootstrap interface, not the source of project truth.

Durable project truth lives in versioned artifacts, specifications, tests, reviews, and decisions.

The Advisory and Specification workflow (using the Advisor and Designer roles) is designed to preserve clarity and prevent premature implementation while the User explores the shape of Symphony Director, Symphony, and future project workflows.