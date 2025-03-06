# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o baggsy .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/baggsy .
RUN chmod +x baggsy
CMD ["./baggsy"]