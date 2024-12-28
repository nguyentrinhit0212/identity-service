FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache gcc musl-dev

# Copy go.mod first
COPY go.mod ./

# Download dependencies and create go.sum if it doesn't exist
RUN go mod download && \
    go mod tidy

# Copy the rest of the code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/main.go

# Install migrate CLI
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    postgresql-client \
    curl

# Copy the binary and required files
COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations ./db/migrations
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY scripts/startup.sh .

# Make startup script executable
RUN chmod +x startup.sh

# Create a non-root user
RUN adduser -D appuser && \
    chown -R appuser:appuser /app
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:4000/health || exit 1

# Expose port
EXPOSE 4000

# Run migrations and start the service
ENTRYPOINT ["./startup.sh"]
