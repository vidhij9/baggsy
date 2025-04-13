# Stage 1: Build the Go binary with a secure Alpine-based image
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o baggsy .

# Stage 2: Minimal, secure final image
FROM amazonlinux:2023
WORKDIR /app
RUN yum update -y && yum install -y ca-certificates && yum clean all
COPY --from=builder /app/baggsy .
EXPOSE 8080
CMD ["./baggsy"]
