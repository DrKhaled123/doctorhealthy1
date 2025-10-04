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

# Verify frontend directory exists (critical for serving UI)
RUN if [ ! -d "frontend" ]; then \
        echo "❌ ERROR: frontend directory not found!" && exit 1; \
    else \
        echo "✅ Frontend directory found with $(ls -1 frontend | wc -l) files"; \
    fi

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

# Expose port (matches internal application port)
EXPOSE 8081

# Enhanced health check with proper port and increased resilience
HEALTHCHECK --interval=15s --timeout=10s --start-period=30s --retries=5 \
    CMD curl -f http://localhost:8081/health || exit 1

# Run the application
CMD ["./main"]
