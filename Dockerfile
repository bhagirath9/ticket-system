# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies for CGO (GORM SQLite driver requires CGO)
RUN apk add --no-cache build-base

WORKDIR /app

# Copy dependency definition and download packages
COPY go.mod ./
RUN go mod download

# Copy the application source code
COPY . .

# Build the production-ready static binary with optimization flags
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o ticket-system cmd/main.go

# Run stage
FROM alpine:latest

# Install libc compatibility for Go's dynamic linking on Alpine
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/ticket-system .

# Copy .env.example as standard configuration template
COPY --from=builder /app/.env.example .env

# Expose default application port
EXPOSE 8080

# Command to execute the server binary
CMD ["./ticket-system"]
