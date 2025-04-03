# Stage 1: Build the Go binary
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o baggsy .

# Stage 2: Create the final image using Amazon Linux 2
FROM amazonlinux:2
WORKDIR /app
# Install ca-certificates for HTTPS support
RUN yum update -y && yum install -y ca-certificates && yum clean all
COPY --from=builder /app/baggsy .
# Expose port 8080 (used by the backend)
EXPOSE 8080
CMD ["./baggsy"]