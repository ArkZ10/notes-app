# Use official Go image for building
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go.mod and go.sum first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Move into cmd/ where main.go is and build
WORKDIR /app/cmd
RUN go build -o /app/main .

# Use a minimal image for running
FROM debian:bookworm-slim

WORKDIR /app

# Install ca-certificates (for HTTPS requests, if needed)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy uploads folder
RUN mkdir -p /app/uploads

# Expose the port
EXPOSE 8080

# Run the app
CMD ["./main"]
