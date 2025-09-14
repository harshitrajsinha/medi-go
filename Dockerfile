# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main ./cmd/medigo/main.go

# Final minimal image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose (GraphQL/REST) server port, if present
EXPOSE 8080

# Run the app
CMD ["./main"]