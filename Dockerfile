# Dockerfile for a Go application using Fiber framework
# This Dockerfile uses a multi-stage build to create a minimal image for the Go application
# --- Builder Stage ---
FROM golang:1.24 AS builder

# Setting the working directory inside the container
WORKDIR /app

# Copy Go module files to leverage Docker cache
COPY go.mod go.sum ./

# Downloading dependencies
RUN go mod download

# Copying the rest of the source code to the container (only the necessary files)
COPY cmd ./cmd
COPY internal ./internal

# Build the application
# CGO_ENABLED=0 creates a statically linked binary, good for minimal images
# -ldflags="-s -w" reduces binary size by removing debug information
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/main ./cmd

# --- Final Stage ---
# Using a minimal base image for the final stage
FROM alpine:3.20

# Set the working directory
WORKDIR /app

# Environment variable to indicate production mode
ENV GO_ENV=prod

# Copying the built executable from the builder stage
COPY --from=builder /app/main .

WORKDIR /app

# Command to run the executable in this case the main binary
CMD ["./main"]