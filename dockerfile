# Tahap 1: Build
FROM golang:1.24.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server -ldflags "-w -s" .

# Tahap 2: Final Image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8081
CMD ["/app/server"]