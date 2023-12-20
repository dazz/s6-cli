# Stage 1: Build the Go binary
FROM golang:1.17 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o s6-cli .

# Stage 2: Create a minimal image for running the binary
FROM scratch

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the final image
COPY --from=builder /app/s6-cli .

# Command to run the binary
CMD ["./s6-cli"]
