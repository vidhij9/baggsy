# Use the official Golang image as a base
FROM golang:1.20

# Set working directory
WORKDIR /app

# Copy Go modules and source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
