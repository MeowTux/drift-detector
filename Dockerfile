# Multi-stage build for Drift Detector
# Optimized for size and security

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=docker -X main.BuildDate=$(date -u '+%Y-%m-%d')" \
    -o drift-detector .

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 drift && \
    adduser -D -u 1000 -G drift drift

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/drift-detector /app/drift-detector

# Create directories
RUN mkdir -p /app/config && \
    chown -R drift:drift /app

# Switch to non-root user
USER drift

# Set environment variables
ENV PATH="/app:${PATH}"
ENV LOG_LEVEL="info"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/drift-detector", "--version"]

ENTRYPOINT ["/app/drift-detector"]
CMD ["detect", "--help"]
