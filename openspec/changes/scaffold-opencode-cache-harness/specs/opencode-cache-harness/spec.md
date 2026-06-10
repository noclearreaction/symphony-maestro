## ADDED Requirements

### Requirement: Docker image builds successfully
The harness SHALL provide a Dockerfile that builds a reproducible Docker image containing opencode and all required dependencies. The image SHALL install a pinned version of opencode via npm at build time.

#### Scenario: Clean build succeeds
- **WHEN** a user runs `docker build` against the harness Dockerfile on a machine with Docker installed
- **THEN** the build completes without errors and produces a tagged image

#### Scenario: Pinned version is used
- **WHEN** the image is built
- **THEN** the exact opencode version specified in the Dockerfile is installed, not the latest

### Requirement: Minimal project fixture is present
The harness SHALL include a minimal project fixture baked into the image. The fixture SHALL consist of a single agent configuration and a system prompt of approximately 200 tokens. The fixture SHALL contain no application code.

#### Scenario: Fixture is available inside container
- **WHEN** a container is started from the harness image
- **THEN** the agent configuration and system prompt files are present at documented paths inside the container

### Requirement: Container is usable for interactive experiments
The harness SHALL provide a documented way to start an interactive shell session inside the container so that an experimenter can invoke opencode manually, observe its behavior, and explore how to trigger turns and read output.

#### Scenario: Interactive session starts
- **WHEN** a user runs the container with an interactive shell command
- **THEN** they land in a shell with opencode available on the PATH and the fixture project present at a documented working directory

### Requirement: Harness is documented
The harness SHALL include a README that describes prerequisites, how to build the image, and how to start an experiment session inside the container.

#### Scenario: README covers getting started
- **WHEN** a user reads the README
- **THEN** they can find instructions for `docker build` and how to enter the container without consulting any other document
