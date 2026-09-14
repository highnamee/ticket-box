# ==========================================
# Stage 1: Build binary
# ==========================================
FROM golang:1.25-alpine AS builder

# Install build dependencies and SSL certificates
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binaries for API server and migration tool
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/bin/api ./cmd/api/main.go

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/bin/migrate ./cmd/migrate/main.go

# ==========================================
# Stage 2: Lightweight Production Runtime
# ==========================================
FROM alpine:3.21 AS runner

# Install essential runtime dependencies (CA certificates for HTTPS, timezone data)
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for enhanced container security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder --chown=appuser:appgroup /app/bin/api /app/api
COPY --from=builder --chown=appuser:appgroup /app/bin/migrate /app/migrate

# Switch to non-root user
USER appuser

# Expose backend default port
EXPOSE 8080

# Default command to run the backend API server
ENTRYPOINT ["/app/api"]
