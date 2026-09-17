# Multi-stage Dockerfile for KasirPro Backend
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and certs
RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# Final lightweight stage
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Jakarta
ENV PORT=8080

COPY --from=builder /app/server /app/server

# Uploads directory with full permissions for container runners
RUN mkdir -p /app/uploads && chmod 777 /app/uploads

EXPOSE 8080

ENTRYPOINT ["/app/server"]
