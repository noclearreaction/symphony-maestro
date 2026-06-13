# =============================================================================
# VERSION CONFIGURATION — managed via .devcontainer/docker-compose.yml
# Renovate manages docker-compose.yml via the dockerfile datasource.
# =============================================================================
# renovate: datasource=golang-version depName=go
ARG GO_VERSION
# renovate: datasource=github-releases depName=denoland/deno
ARG DENO_VERSION
# renovate: datasource=github-releases depName=go-task/task
ARG TASK_VERSION

# =============================================================================
# Stage: download-base
# Purpose: Common utilities for binary downloads
# Output: /srv/files/ working directory with curl, ca-certificates, tar, xz-utils
# =============================================================================
FROM debian:bookworm-slim AS download-base
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        curl \
        ca-certificates \
        tar \
        xz-utils \
 && rm -rf /var/lib/apt/lists/*
WORKDIR /srv/files

# =============================================================================
# Stage: go-runtime
# Purpose: Complete Go toolchain from the official Go image
# Output: /srv/files/go<GO_VERSION>/ — full Go tree ready to COPY into base
# =============================================================================
FROM golang:${GO_VERSION}-alpine AS go-runtime
ARG GO_VERSION
RUN mkdir -p /srv/files/go${GO_VERSION} \
 && cp -r /usr/local/go/. /srv/files/go${GO_VERSION}/

# =============================================================================
# Stage: deno-runtime
# Purpose: Download the Deno binary from the official GitHub release
# Output: /srv/files/deno-<DENO_VERSION>/bin/deno
# Note: Uses debian:bookworm-slim directly (needs unzip, not in download-base)
# =============================================================================
FROM debian:bookworm-slim AS deno-runtime
ARG DENO_VERSION
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        curl \
        ca-certificates \
        unzip \
 && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /srv/files/deno-${DENO_VERSION}/bin \
 && curl -fL \
        -o /tmp/deno.zip \
        -- "https://github.com/denoland/deno/releases/download/v${DENO_VERSION}/deno-x86_64-unknown-linux-gnu.zip" \
 && cd /tmp \
 && unzip deno.zip \
 && mv /tmp/deno /srv/files/deno-${DENO_VERSION}/bin/deno \
 && chmod +x /srv/files/deno-${DENO_VERSION}/bin/deno

# =============================================================================
# Stage: task-binary
# Purpose: Download the go-task binary from the official GitHub release
# Output: /srv/files/task
# =============================================================================
FROM download-base AS task-binary
ARG TASK_VERSION
RUN curl -fL \
        -o /tmp/task.tar.gz \
        -- "https://github.com/go-task/task/releases/download/v${TASK_VERSION}/task_linux_amd64.tar.gz" \
 && tar -xzf /tmp/task.tar.gz -C /srv/files task \
 && rm /tmp/task.tar.gz \
 && chmod +x /srv/files/task

# =============================================================================
# Stage: base
# Purpose: Assembly — all runtimes combined into the devcontainers base image
# =============================================================================
FROM mcr.microsoft.com/devcontainers/base:ubuntu-22.04 AS base
ARG GO_VERSION
ARG DENO_VERSION

# --- Go ---
COPY --from=go-runtime /srv/files/go${GO_VERSION} /opt/go${GO_VERSION}
RUN update-alternatives --install /usr/local/bin/go go /opt/go${GO_VERSION}/bin/go 100 \
        --slave /usr/local/bin/gofmt gofmt /opt/go${GO_VERSION}/bin/gofmt
ENV GOROOT=/opt/go${GO_VERSION} \
    GOPATH=/home/vscode/go

# --- Deno ---
COPY --from=deno-runtime /srv/files/deno-${DENO_VERSION} /opt/deno-${DENO_VERSION}
RUN update-alternatives --install /usr/local/bin/deno deno /opt/deno-${DENO_VERSION}/bin/deno 100

# --- Task ---
COPY --from=task-binary /srv/files/task /usr/local/bin/task

# =============================================================================
# Stage: ci
# Purpose: CI-appropriate environment layered on base
# =============================================================================
FROM base AS ci
ENV CI=true \
    DEBIAN_FRONTEND=noninteractive

# =============================================================================
# Stage: final
# Purpose: Developer environment layered on ci — SSH agent, XDG, build tracking
# =============================================================================
FROM ci AS final

# SSH agent profile script
COPY ssh-agent.sh /etc/profile.d/ssh-agent.sh
RUN echo "source /etc/profile.d/ssh-agent.sh" >> /home/vscode/.bashrc \
 && echo "source /etc/profile.d/ssh-agent.sh" >> /home/vscode/.zshrc

# VS Code server directories
RUN mkdir -p \
        /home/vscode/.vscode-server/extensions \
        /home/vscode/.vscode-server/data/User \
 && chown -R vscode:vscode /home/vscode/.vscode-server

# XDG environment variables and pager
ENV XDG_CONFIG_HOME=/home/vscode/.config \
    XDG_DATA_HOME=/home/vscode/.local/share \
    XDG_CACHE_HOME=/home/vscode/.cache \
    XDG_STATE_HOME=/home/vscode/.local/state \
    PAGER=cat

# Build tracking
RUN date +%s > /container_build_id
