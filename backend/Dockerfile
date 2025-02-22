FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o baggsy .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/baggsy .
CMD ["./baggsy"]