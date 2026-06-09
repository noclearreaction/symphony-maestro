# Strategic Goals & Vision

This document tracks the overarching goals, boundaries, and vision for Symphony Director and Symphony. It is the north star for both the Advisor and Designer agents.

## 1. Vision & Core Philosophy
* **Human Ownership**: The human remains the sole owner of intent, tradeoffs, approval, and final decisions. AI agents propose and execute; they do not govern.
* **Artifact-Driven Coordination**: Agents communicate across boundaries using versioned, readable artifacts (OpenSpec, Gherkin, Markdown, test evidence) rather than hidden instructions or chat histories.

## 2. Active Strategic Goals (Bootstrap Phase)
* **Goal 1: Establish Advisory Continuity**
  * *What*: Ensure the Advisor and Designer have instant access to active goals, roadmaps, and decisions so they don't design in a vacuum.
  * *Success Criteria*: A `strategy/` directory exists with living files, and both agents are configured to read it.
* **Goal 2: Define Symphony v1 Architecture & Lifecycle**
  * *What*: Map out Symphony's core event-driven work intake and state transitions.
  * *Success Criteria*: An OpenSpec proposal for Symphony v1's core execution loop.

## 3. Explicit Non-Goals
* **No Runtime in Symphony Director**: We will not write application code, databases, or runtime systems in the Symphony Director repository.
* **No Premature Automation**: We will not build complex Symphony orchestration agents until we have manually run the workflow via OpenCode and proven its value.