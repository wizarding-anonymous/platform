# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -v -o auth-service ./cmd/auth-service

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/auth-service .
# Copy config files
COPY configs/ ./configs/

# Expose port 8080 and 50051 (for HTTP and gRPC)
EXPOSE 8080
EXPOSE 50051

# Command to run the executable
CMD ["./auth-service"]
