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

RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /workspaces/tunectl

COPY --chown=vscode:vscode go.mod go.sum ./

RUN go mod download