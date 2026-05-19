# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
# Install git for go modules
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/

# Runtime stage
FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    ca-certificates \
    wget \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
RUN mkdir -p /app/recordings
EXPOSE 8080
CMD ["./server"]
