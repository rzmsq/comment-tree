# Stage 1: Build the application
FROM golang:latest AS builder

WORKDIR /app

# Copy go.mod and go.sum to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# The main package is in the comment_tree directory
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./comment_tree/main.go

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

COPY --from=builder /app/web ./web

# Copy the configuration file
COPY config.yaml .

EXPOSE 8080

# Command to run the application
CMD ["./main"]
