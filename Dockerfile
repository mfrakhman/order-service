# --- Stage 1: Build Stage ---
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application's source code
# (config/, internal/, cmd/, etc.)
COPY . .

# Build the application.
# We point to './cmd' as that is where your main.go is located.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd

# --- Stage 2: Final Stage ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the compiled binary from the 'builder' stage
COPY --from=builder /app/main .

# Expose the port your app runs on
EXPOSE 3002

# Run the binary
CMD ["./main"]