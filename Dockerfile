# Multi-stage Dockerfile for Quiver
# Stage 1: Build the Go application
FROM golang:1.24.2-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o quiver ./cmd/quiver

# Stage 2: Create the final image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S quiver && \
    adduser -u 1001 -S quiver -G quiver

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/quiver .

# Copy configuration files
COPY --from=builder /app/template ./template
COPY --from=builder /app/.gitignore ./.gitignore

# Create necessary directories
RUN mkdir -p logs pkgs && \
    chown -R quiver:quiver /app

# Switch to non-root user
USER quiver

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set environment variables
ENV QUIVER_HOST=0.0.0.0
ENV QUIVER_PORT=8080
ENV QUIVER_LOG_LEVEL=info

# Run application
CMD ["./quiver"] 