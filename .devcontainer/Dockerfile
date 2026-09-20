FROM mcr.microsoft.com/devcontainers/go:1.26-bookworm

USER root

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt/lists,sharing=locked \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        gcc \
        python3 \
        ffmpeg \
        nodejs

RUN mkdir -p /home/vscode/.cache && chown vscode:vscode /home/vscode/.cache

USER vscode

RUN curl -fsSL https://deno.land/install.sh | sh -s v2.1.4

ENV PATH="/home/vscode/.deno/bin:/go/bin:${PATH}"

RUN --mount=type=cache,target=/go/pkg,sharing=locked,uid=1000,gid=1000 \
    --mount=type=cache,target=/home/vscode/.cache/go-build,sharing=locked,uid=1000,gid=1000 \
    go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /workspaces/tunectl

COPY --chown=vscode:vscode go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg,sharing=locked,uid=1000,gid=1000 \
    go mod download