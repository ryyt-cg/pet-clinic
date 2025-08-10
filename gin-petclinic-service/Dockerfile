# Stage 1: Builder
FROM golang:1.24-bookworm AS builder

ENV CGO_ENABLED=1

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      build-essential \
      libsqlite3-dev


WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o gin-petclinic-service cmd/server.go

# Stage 2: Final image
FROM debian:bookworm
RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      libsqlite3-0

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/gin-petclinic-service .
COPY --from=builder /app/config/*.yaml config/.
COPY --from=builder /app/cmd/certs/* cmd/certs/.
COPY --from=builder /app/pet-clinic.sqlite .

# Expose the port your gin application listens on
EXPOSE 8443

# Command to run the application
CMD ["./gin-petclinic-service"]