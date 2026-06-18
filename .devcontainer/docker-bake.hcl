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

group "default" {
  targets = [
    "devcontainer", 
    "node-apps", 
    "renovate"
  ]
}

variable "PROJECT" {
  default = "symphony-maestro"
}

variable "PLATFORM" {
  default = "linux/amd64"
}

# ---------------------------------------------------------------------------
# Leaf targets — no inter-target dependencies, build in parallel
# ---------------------------------------------------------------------------

target "download-base" {
  context    = "."
  dockerfile = "docker/Dockerfile.download-base"
  platforms  = [PLATFORM]
  args = {
    DEBIAN_VERSION = VERSIONS.debian
  }
}

target "go-runtime" {
  context    = "."
  dockerfile = "docker/Dockerfile.go-runtime"
  platforms  = [PLATFORM]
  args = {
    GO_VERSION = VERSIONS.go
  }
}

target "deno-runtime" {
  context    = "."
  dockerfile = "docker/Dockerfile.deno-runtime"
  platforms  = [PLATFORM]
  args = {
    DEBIAN_VERSION = VERSIONS.debian
    DENO_VERSION   = VERSIONS.deno
  }
}

target "node-builder" {
  context    = "."
  dockerfile = "docker/Dockerfile.node-builder"
  platforms  = [PLATFORM]
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
  platforms  = [PLATFORM]
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
  platforms  = [PLATFORM]
  contexts = {
    node-builder = "target:node-builder"
  }
  tags = ["node-apps:${PROJECT}"]
}

# ---------------------------------------------------------------------------
# Renovate — pull and retag upstream image as local
# ---------------------------------------------------------------------------

target "renovate" {
  context    = "."
  dockerfile = "docker/Dockerfile.renovate"
  platforms  = [PLATFORM]
  args = {
    RENOVATE_VERSION = VERSIONS.renovate
  }
  tags = ["renovate:${PROJECT}"]
}

# ---------------------------------------------------------------------------
# Final assembly — devcontainer
# ---------------------------------------------------------------------------

target "devcontainer" {
  context    = "."
  dockerfile = "docker/Dockerfile.symphony-studio"
  platforms  = [PLATFORM]
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
  tags = ["devcontainer:${PROJECT}"]
}
