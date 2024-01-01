# Stage 1: Build the application
FROM golang:latest as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
WORKDIR /app/cmd/proxy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o proxy .

WORKDIR /app/cmd/shard
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shard .

# Stage 2: Run the application using scratch
FROM scratch

# Set the working directory
WORKDIR /app
USER 1001
# Copy the binary from the builder stage
COPY --from=builder /app/cmd/proxy/proxy .
COPY --from=builder /app/cmd/shard/shard .

# Command to run the executable
CMD ["./proxy"]
