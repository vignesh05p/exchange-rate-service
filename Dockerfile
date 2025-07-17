# Use the official Golang image as builder
FROM golang:1.21 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app from cmd/main.go
RUN go build -o app ./cmd

# Use a minimal image for final container
FROM alpine:latest

# Create a user (optional for security)
RUN adduser -D appuser

# Set working directory
WORKDIR /home/appuser

# Copy built binary from builder
COPY --from=builder /app/app .

# Set executable permission (just in case)
RUN chmod +x ./app

# Run as non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Command to run the app
CMD ["./app"]
