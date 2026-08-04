# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server/

# Runtime stage — needs real shell (bash) for PTY terminals
FROM debian:12-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    zsh \
    ca-certificates \
    curl \
    wget \
    procps \
    && rm -rf /var/lib/apt/lists/*

# Non-root user for running the server (PTY still spawns user shells)
RUN groupadd -r termhub && useradd -r -g termhub -m -s /bin/bash termhub

WORKDIR /app
COPY --from=builder /app/server .
RUN mkdir -p /app/recordings /app/uploads

# Note: NOT using "USER termhub" — PTY spawning requires fork/exec permissions
# which Docker's default seccomp profile blocks for non-root users.
# The container is isolated by Docker's namespace/cgroup boundaries.
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/health || exit 1
CMD ["./server"]
