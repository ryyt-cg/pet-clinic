# Stage 1: Builder
FROM golang:1.25-bookworm AS builder

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
FROM scratch

# Uncomment the following lines if you want to use a Debian base image instead of scratch
#FROM debian:bookworm
#RUN apt-get update \
# && DEBIAN_FRONTEND=noninteractive \
#    apt-get install --no-install-recommends --assume-yes \
#      libsqlite3-0
#
#WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/gin-petclinic-service .
COPY --from=builder /app/config/*.yaml config/.
COPY --from=builder /app/cmd/certs/* cmd/certs/.
COPY --from=builder /app/pet-clinic.sqlite .

# Expose the port your gin application listens on
EXPOSE 8092

# Command to run the application
CMD ["./gin-petclinic-service"]