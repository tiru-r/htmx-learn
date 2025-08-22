# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git curl

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Download Tailwind standalone CLI
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 && \
    chmod +x tailwindcss-linux-x64 && \
    mv tailwindcss-linux-x64 tailwindcss

# Copy source code
COPY . .

# Generate templates and build CSS
RUN templ generate
RUN ./tailwindcss -i static/css/input.css -o static/css/output.css

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/htmx-learn

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/db ./db

# Create data directory for SQLite database
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Set environment variable for database path
ENV DATABASE_PATH=/data/app.db

# Run the binary
CMD ["./main"]