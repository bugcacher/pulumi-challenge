# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o app .

# Run stage
FROM alpine:latest

# Add basic security-related packages
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/app .

# Expose the port your application uses
EXPOSE 8080

# Run the application
CMD ["./app"]