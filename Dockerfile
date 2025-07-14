# Multi-stage build for Resume API

# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go.mod and go.sum files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o resume-api ./cmd/api

# Final stage
FROM alpine:3.18

# Add non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/resume-api /app/
COPY --from=builder /app/migrations /app/migrations

# Set ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the API port
EXPOSE 8080

# Set environment variables with defaults
ENV RESUME_API_ENVIRONMENT=production \
    RESUME_API_SERVER_HOST=0.0.0.0 \
    RESUME_API_SERVER_PORT=8080 \
    RESUME_API_LOGGING_LEVEL=info \
    RESUME_API_LOGGING_FORMAT=json

# Command to run the application
CMD ["/app/resume-api"]