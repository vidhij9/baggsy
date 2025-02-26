FROM golang:1.24 AS builder
WORKDIR /app/backend
COPY . .
RUN go mod tidy && go build -o baggsy ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backend/baggsy .
EXPOSE 8080
CMD ["./baggsy"]