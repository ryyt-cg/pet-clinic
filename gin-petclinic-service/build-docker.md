# Stage 1: Build the application

FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o fiber-app .

# Stage 2: Create a minimal runtime image

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/fiber-app .

EXPOSE 3000

CMD ["./fiber-app"]
