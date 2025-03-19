# Build stage
FROM golang:1.22-alpine AS builder

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

# Build the application with verbose output to see any errors
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Install dependencies for production
RUN apk --no-cache add ca-certificates tzdata

# Create directory structure for config
RUN mkdir -p ./internal/core/config/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy config files
COPY --from=builder /app/internal/core/config/config.yaml ./internal/core/config/
COPY --from=builder /app/internal/core/config/config.development.yaml ./internal/core/config/
COPY --from=builder /app/internal/core/config/config.production.yaml ./internal/core/config/

# Set environment variables
ENV APP_ENV=production

# Make sure the binary is executable
RUN chmod +x ./main

# Expose the application port
EXPOSE 8080

# Print directory contents to debug
RUN ls -la

# Run the application
CMD ["./main"]