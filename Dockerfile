# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy .env file if it exists
COPY .env* ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o maintenance ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/maintenance .

# Copy .env file from builder stage
COPY --from=builder /app/.env* ./

# Make it executable
RUN chmod +x ./maintenance

# Expose port (matching your docker-compose.yml)
EXPOSE 9999

# Command to run
CMD ["./maintenance"]
