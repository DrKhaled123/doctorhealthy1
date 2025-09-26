# Multi-stage Dockerfile for Go (Echo) API with SQLite (CGO enabled)

FROM golang:1.22-bookworm AS builder

RUN apt-get update -y && apt-get install -y --no-install-recommends \
    build-essential pkg-config && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

ENV CGO_ENABLED=1

# Build
RUN go build -ldflags="-s -w" -o server .


FROM debian:bookworm-slim AS runtime

RUN apt-get update -y && apt-get install -y --no-install-recommends \
    ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Non-root user
RUN useradd -m -u 10001 appuser

COPY --from=builder /app/server /app/server

ENV PORT=8081 \
    GIN_MODE=release

EXPOSE 8081

USER appuser

CMD ["/app/server"]


