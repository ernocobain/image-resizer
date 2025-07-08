# Tahap 1: Build
FROM golang:1.24.4 AS builder

# Install build-essentials yang berisi C compiler
RUN apt-get update && apt-get install -y build-essential

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Hapus CGO_ENABLED=0 untuk mengizinkan CGo
RUN go build -o /app/server -ldflags "-w -s" .

FROM alpine:latest

RUN apk add --no-cache libwebp-dev

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8081
CMD ["/app/server"]