# # Stage 1: Builder
# FROM golang:1.23 AS builder

# WORKDIR /app

# # Copy go.mod, go.sum first, then download dependencies
# COPY ./backend go.mod go.sum ./

# # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# # Copy the rest of the source code
# COPY . .

# # Build the Go binary
# RUN go build -o /baggsy/backend main.go

# # Stage 2: Runner
# FROM alpine:3.18
# WORKDIR /app

# # Copy compiled binary from builder
# COPY --from=builder /baggsy/backend /app/

# # Expose port 8080
# EXPOSE 8080

# CMD ["/app/baggsy/backend"]


# Baggsy/backend/Dockerfile

# Stage 1: Builder
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go binary
RUN go build -o /baggsy main.go

# Stage 2: Minimal runtime image
FROM alpine:3.18
WORKDIR /app

COPY --from=builder /baggsy /app/baggsy

# Expose the port your Go app listens on
EXPOSE 8080

CMD ["/app/baggsy"]
