# =============================================================================
# VERSION CONFIGURATION — managed via .devcontainer/docker-compose.yml
# Renovate manages docker-compose.yml via the dockerfile datasource.
# =============================================================================
# renovate: datasource=golang-version depName=go
ARG GO_VERSION="1.26.4"
# renovate: datasource=github-releases depName=denoland/deno
ARG DENO_VERSION="2.8.2"
# renovate: datasource=github-releases depName=go-task/task
ARG TASK_VERSION="3.51.1"
# renovate: datasource=node depName=node
ARG NODE_VERSION="22.16.0"

# =============================================================================
# Stage: download-base
# Purpose: Common utilities for binary downloads
# Output: /srv/files/ working directory with curl, ca-certificates, tar, xz-utils
# =============================================================================
FROM debian:bookworm-slim AS download-base
RUN --mount=type=cache,id=apt-cache-bookworm,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,id=apt-lib-bookworm,target=/var/lib/apt,sharing=locked \
    apt-get update \
 && apt-get install -y --no-install-recommends \
        curl \
        ca-certificates \
        tar \
        xz-utils
WORKDIR /srv/files

# =============================================================================
# Stage: node-builder
# Purpose: Prepare Node + pnpm
# =============================================================================
FROM node:${NODE_VERSION}-bookworm-slim AS node-builder

ARG PNPM_VERSION

ENV COREPACK_ENABLE_AUTO_PIN=1 \
    COREPACK_ENABLE_STRICT=1 \
    COREPACK_ENABLE_DOWNLOAD_PROMPT=0

RUN --mount=type=cache,id=apt-cache-bookworm,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,id=apt-lib-bookworm,target=/var/lib/apt,sharing=locked \
    apt-get update \
 && apt-get install -y --no-install-recommends jq
RUN corepack enable \
 && install -d -o 1000 -g 1000 /opt/node

USER 1000:1000

RUN corepack prepare pnpm@${PNPM_VERSION} --activate

WORKDIR /home/node

# =============================================================================
# Stage: node-apps
# Purpose: Install node packages into /opt/node
# =============================================================================
FROM node-builder AS node-apps

COPY node/pnpm-workspace.yaml /home/node/pnpm-workspace.yaml
COPY node/package.json /home/node/package.json
COPY node/pnpm-lock.yaml* /home/node/

RUN --mount=type=cache,id=symphony-node-cache,target=/home/node/.cache/,uid=1000,gid=1000,sharing=locked \
    --mount=type=cache,id=symphony-node-module-cache,target=/home/node/node_modules/,uid=1000,gid=1000,sharing=locked \
    [ -f pnpm-lock.yaml ] && export ARGS="$ARGS --frozen-lockfile" ; \
    pnpm install $ARGS; \
    pnpm --filter=. --prod deploy dist

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
RUN --mount=type=cache,id=apt-cache-bookworm,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,id=apt-lib-bookworm,target=/var/lib/apt,sharing=locked \
    apt-get update \
 && apt-get install -y --no-install-recommends \
        curl \
        ca-certificates \
        unzip
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

# --- Node tools mount point and PATH ---
ENV PATH=/opt/node/node_modules/.bin:$PATH
COPY --from=node-apps /home/node/dist/ /opt/node/



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
FROM base AS final

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
    PAGER=less \
    TASK_TEMP_DIR=/tmp/task

# task -g global Taskfile (symlinks — resolves at runtime when workspace is mounted)
# Both vscode and root get the symlink so 'sudo task -g' works too
RUN ln -s /workspace/.devcontainer/resources/Taskfile-global.yaml /home/vscode/Taskfile.yaml \
 && ln -s /workspace/.devcontainer/resources/Taskfile-global.yaml /root/Taskfile.yaml

# Build tracking
RUN date +%s > /container_build_id
