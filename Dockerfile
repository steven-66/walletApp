# Use the official Go image as the base image for building the app
FROM golang:latest AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o wallet-cli .

# Use Ubuntu 22.04 as the base image for the runtime environment
FROM ubuntu:22.04

# Set the working directory inside the container
WORKDIR /app

# Install necessary dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy the compiled binary from the build stage
COPY --from=build /app/wallet-cli .

# Run the application
CMD ["./wallet-cli"]