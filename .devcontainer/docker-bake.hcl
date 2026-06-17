# =============================================================================
# Symphony Studio — Docker Bake file
#
# Each stage lives in its own Dockerfile. This file is the wiring diagram.
# The original monolithic Dockerfile is preserved but not referenced here.
#
# Versions are defined in versions.hcl. Load both files together:
#
#   docker buildx bake -f .devcontainer/versions.hcl -f .devcontainer/docker-bake.hcl
#   docker buildx bake -f .devcontainer/versions.hcl -f .devcontainer/docker-bake.hcl --print
# =============================================================================

# ---------------------------------------------------------------------------
# Leaf targets — no inter-target dependencies, build in parallel
# ---------------------------------------------------------------------------

target "download-base" {
  context    = "."
  dockerfile = "docker/Dockerfile.download-base"
  args = {
    DEBIAN_VERSION = VERSIONS.debian
  }
}

target "go-runtime" {
  context    = "."
  dockerfile = "docker/Dockerfile.go-runtime"
  args = {
    GO_VERSION = VERSIONS.go
  }
}

target "deno-runtime" {
  context    = "."
  dockerfile = "docker/Dockerfile.deno-runtime"
  args = {
    DEBIAN_VERSION = VERSIONS.debian
    DENO_VERSION   = VERSIONS.deno
  }
}

target "node-builder" {
  context    = "."
  dockerfile = "docker/Dockerfile.node-builder"
  args = {
    NODE_VERSION = VERSIONS.node
  }
}

# ---------------------------------------------------------------------------
# Dependent targets — consume leaf targets via named contexts
# ---------------------------------------------------------------------------

target "task-binary" {
  context    = "."
  dockerfile = "docker/Dockerfile.task-binary"
  contexts = {
    download-base = "target:download-base"
  }
  args = {
    TASK_VERSION = VERSIONS.task
  }
}

target "node-apps" {
  context    = "."
  dockerfile = "docker/Dockerfile.node-apps"
  contexts = {
    node-builder = "target:node-builder"
  }
  tags = ["symphony-studio-node-apps:local"]
}

# ---------------------------------------------------------------------------
# Final assembly — symphony-studio
# ---------------------------------------------------------------------------

target "symphony-studio" {
  context    = "."
  dockerfile = "docker/Dockerfile.symphony-studio"
  contexts = {
    go-runtime   = "target:go-runtime"
    deno-runtime = "target:deno-runtime"
    task-binary  = "target:task-binary"
    node-apps    = "target:node-apps"
  }
  args = {
    UBUNTU_VERSION = VERSIONS.ubuntu
    GO_VERSION     = VERSIONS.go
    DENO_VERSION   = VERSIONS.deno
  }
  tags = ["symphony-studio:local"]
}

# Default builds the devcontainer image (node-apps is an intermediate, not a
# separate deliverable — it is reached transitively via symphony-studio).
group "default" {
  targets = ["symphony-studio"]
}
