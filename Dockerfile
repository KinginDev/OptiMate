# Start with a base image containing Go runtime
FROM golang:1.23 as userServiceBuilder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files from the user-service directory
COPY user-service/go.mod user-service/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY user-service/ .

WORKDIR /app/cmd/api

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/user-service

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /root/

# Install ca-certificates to allow SSL-based applications
RUN apk --no-cache add ca-certificates

# Copy the Pre-built binary file from the previous stage
COPY --from=userServiceBuilder /bin/user-service .

# Command to run the executable
CMD ["./user-service"]
