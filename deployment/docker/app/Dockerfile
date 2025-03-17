# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Install dependencies for production
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy config files
COPY --from=builder /app/internal/core/config/config.yaml ./internal/core/config/

# Set environment variables
ENV APP_ENV=production

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]