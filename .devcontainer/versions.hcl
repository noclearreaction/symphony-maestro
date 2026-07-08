# .devcontainer/versions.hcl
# Single source of truth for all pinned tool versions.
# Consumed by docker-bake.hcl via: docker buildx bake -f versions.hcl -f docker-bake.hcl
# Renovate tracks these via the regex custom manager in renovate.json.

variable "VERSIONS" {
  default = {
    # renovate: datasource=docker depName=debian
    debian = "12.14-slim"
    # renovate: datasource=endoflife-date depName=ubuntu
    ubuntu = "24.04"
    # renovate: datasource=golang-version depName=go
    go = "1.26.4"
    # renovate: datasource=github-releases depName=denoland/deno
    deno = "2.8.2"
    # renovate: datasource=github-releases depName=go-task/task
    task = "3.51.1"
    # renovate: datasource=node depName=node
    node = "24.16.0"
    # renovate: datasource=github-releases depName=renovatebot/renovate
    renovate = "43.220.0"
  }
}
