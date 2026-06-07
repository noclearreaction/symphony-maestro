## Why

We currently have a sophisticated conceptual model for our development workflow (Advisor, Designer, OpenSpec, agent boundaries), but because the local execution plumbing is missing, we cannot actually use the repository resources to build the project. If we try to design Symphony now, we do so in a vacuum. 

We need to ground ourselves. The absolute best way to learn how to use OpenSpec, test our agent boundaries, and establish a real workflow is to dogfood it. This proposal halts Symphony design and directs 100% of our energy to bootstrapping the Director's own local workflow, moving from an empty workspace to a fully functioning, git-backed, review-gated, and tool-integrated workflow.

## What Changes

- **Custom Builder Agent**: We introduce a custom, role-bounded `builder` agent by adding `.opencode/agents/builder.md`. This agent is dynamic, self-loading, and relies on strict prompt-level soft boundaries to execute tasks without risking strategic drift or direct main branch commits.
- **Git Flow for Artifacts**: We implement a standardized git feature-branch flow (`change/<name>`) for artifact creation, testing, and promotion.
- **Structured Workspace**: We establish a designated local scribble sanctuary under `.symphony/` which is untracked by git to isolate temporary files, logs, plans, reviews, and test evidence.
- **GitHub PR Integration**: We automate the progression from local feature branches to proper GitHub Pull Requests using the GitHub CLI (`gh`).
- **Automation Helper Scripts**: We write lightweight, local wrapper scripts under `bin/` to manage the lifecycle transitions cleanly.

## Capabilities

### New Capabilities
- `director-workflow`: Structured Git, GitHub PR, and custom builder agent local workflow for Director development.

### Modified Capabilities
<!-- No existing capabilities exist to modify -->

## Impact

- Modifies `.gitignore` to untrack `.symphony/` directories.
- Adds `.opencode/agents/builder.md` to dynamically register the builder agent.
- Introduces local bash scripts under `bin/`.
- Has zero impact on external application code or running runtime daemons.
