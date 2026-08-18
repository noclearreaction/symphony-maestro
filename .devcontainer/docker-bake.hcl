# =============================================================================
# Symphony Studio — Docker Bake file
#
# Each stage lives in its own Dockerfile. This file is the wiring diagram.
# The original monolithic Dockerfile is preserved but not referenced here.
#
# Versions are defined in versions.hcl. Load both files together:
#
#   docker buildx bake --allow=fs.read=.. -f versions.hcl -f docker-bake.hcl
#   docker buildx bake --allow=fs.read=.. -f versions.hcl -f docker-bake.hcl --print
#
# --allow=fs.read=.. is required by the rubato-runtime target, whose context is
# the repository root rather than .devcontainer/. Run from .devcontainer/.
# =============================================================================

group "default" {
  targets = [
    "devcontainer",
    "node-apps",
    "renovate",
    "rubato-runtime"
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

target "vibe-runtime" {
  context    = "."
  dockerfile = "docker/Dockerfile.vibe-runtime"
  platforms  = [PLATFORM]
  contexts = {
    download-base = "target:download-base"
  }
  args = {
    VIBE_VERSION = VERSIONS.vibe
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
# Rubato — standalone proxy, built from repo source
#
# Context is the repository root, unlike every other target, because the Go
# sources live outside .devcontainer/. Deliberately does not build on the
# devcontainer image: rubato is a separate service and should not inherit the
# full development toolchain.
# ---------------------------------------------------------------------------

target "rubato-runtime" {
  context    = ".."
  dockerfile = ".devcontainer/docker/Dockerfile.rubato-runtime"
  platforms  = [PLATFORM]
  contexts = {
    go-runtime = "target:go-runtime"
  }
  args = {
    GO_VERSION     = VERSIONS.go
    DEBIAN_VERSION = VERSIONS.debian
  }
  tags = ["rubato:${PROJECT}"]
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
    vibe-runtime = "target:vibe-runtime"
  }
  args = {
    UBUNTU_VERSION = VERSIONS.ubuntu
    GO_VERSION     = VERSIONS.go
    DENO_VERSION   = VERSIONS.deno
    VIBE_VERSION   = VERSIONS.vibe
  }
  tags = ["devcontainer:${PROJECT}"]
}
