# Multi-stage Dockerfile for Go API with SQLite (CGO enabled)
# Force cache invalidation - update this comment to force rebuild
# Build version: 1759252953

FROM golang:1.22-bookworm AS builder

# Install build dependencies
RUN apt-get update -y && apt-get install -y --no-install-recommends \
    build-essential pkg-config libsqlite3-dev && \
    rm -rf /var/lib/apt/lists/*

# Create app user for build stage
RUN groupadd -g 1000 appuser && \
    useradd -u 1000 -g appuser -s /bin/bash -m appuser

WORKDIR /app

# Copy go mod files for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main .

# Final stage - Use Debian for better SQLite CGO compatibility
FROM debian:bookworm-slim

# Install runtime dependencies INCLUDING curl for health checks
RUN apt-get update -y && apt-get install -y --no-install-recommends \
    ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -g 1000 appuser && \
    useradd -u 1000 -g appuser -s /bin/bash -m appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Create data directory with correct permissions
RUN mkdir -p /app/data && chown -R appuser:appuser /app/data

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check using curl (NOT wget)
HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
