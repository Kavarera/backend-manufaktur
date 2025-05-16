# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install git (needed for go modules)
RUN apk add --no-cache git

# Set working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies (cache them if no changes)
RUN go mod download

# Copy all source code
COPY . .

# Build the binary statically, optimized
RUN CGO_ENABLED=0 GOOS=linux go build -o manufacture_API main.go

# Stage 2: Run
FROM alpine:latest

# Install CA certificates for HTTPS requests (if needed)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/manufacture_API .

# Expose the port your app runs on
EXPOSE 8080

# Set default environment variables (optional)
ENV PORT=8080

# Command to run the executable
CMD ["./manufacture_API"]
