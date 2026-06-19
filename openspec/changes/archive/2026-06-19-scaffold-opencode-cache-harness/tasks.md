## 1. Repository Structure

- [x] 1.1 Create `harness/` directory at the repository root
- [x] 1.2 Create `harness/fixture/` subdirectory for the minimal project files

## 2. Docker Image

- [x] 2.1 Write `harness/Dockerfile` using `node:20-slim` as the base image
- [x] 2.2 Pin the opencode version in the Dockerfile (`npm install -g opencode-ai@1.16.2`)
- [x] 2.3 Verify `docker build -t opencode-cache-harness harness/` succeeds without errors

## 3. Minimal Project Fixture

- [x] 3.1 Create `harness/fixture/opencode.json` (opencode config referencing AGENTS.md, sets `experiment` as default agent)
- [x] 3.2 Write the agent system prompt to `harness/fixture/AGENTS.md` (~160 tokens, no application code)
- [x] 3.3 Confirm fixture files are copied into the image at a documented path (`/app/fixture/`)

## 4. End-to-End Validation

- [x] 4.1 Start a container interactively and confirm opencode is on the PATH
- [x] 4.2 Confirm the fixture project is present at the documented working directory inside the container
- [x] 4.3 Confirm opencode can be invoked manually (e.g., `opencode --version` returns `1.16.2`)

## 5. Documentation

- [x] 5.1 Write `harness/README.md` covering prerequisites (Docker), `docker build` command, and how to start an interactive experiment session
- [x] 5.2 Document the pinned opencode version and the fixture layout in the README
