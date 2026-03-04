# Build stage
FROM golang:1.24-alpine AS builder

# Install git for fetching dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary (no C dependencies)
# GOOS=linux for Linux containers
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /blog-service ./cmd/blog-service

# Final stage (minimal image)
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /blog-service .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./blog-service"]
