# build stage
FROM golang:1.20.12-alpine3.19 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o s6-cli -v ./cmd/s6cli


# Stage 2: Create a minimal image for running the binary
FROM scratch

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the final image
COPY --from=builder /app/s6-cli .

# Command to run the binary
CMD ["./s6-cli"]
