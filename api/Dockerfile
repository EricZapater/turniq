# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for fetching dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Runtime Stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder
COPY --from=builder /app/main .

# Copy migrations so the app can run them on startup
COPY --from=builder /app/migrations ./migrations

# Expose the application port (adjusted to 8080 as per standard)
EXPOSE 8080

CMD ["./main"]
