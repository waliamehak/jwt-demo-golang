# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git and ca-certificates so Go can download modules
RUN apk add --no-cache git ca-certificates

# Enable Go modules and use the Go proxy
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

# Copy Go modules manifests and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o jwt-server main.go

# Stage 2: Minimal image to run the binary
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests (if needed)
RUN apk add --no-cache ca-certificates

# Copy the compiled binary from builder stage
COPY --from=builder /app/jwt-server .

# Set default JWT secret (can be overrideen via environment variable)
ENV JWT_SECRET="AVeryStrongSecretIndeed!"

# Expose port 8080
EXPOSE 8080

# Run the server
CMD ["./jwt-server"]
