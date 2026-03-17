# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o main ./cmd

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary
COPY --from=builder /app/main .

# Expose port (Render expects 8080)
EXPOSE 8080

# Run the binary
CMD ["./main"]