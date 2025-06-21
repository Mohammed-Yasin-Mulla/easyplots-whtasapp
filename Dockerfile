# Build Stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/server/main.go

# Final Stage
FROM alpine:latest
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/server .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./server"]

