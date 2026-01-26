FROM golang:1.25-alpine AS builder

WORKDIR /app

# Security updates
RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 pour une image complètement statique
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd

# Runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/app .

# Rendre l'exécutable executable
RUN chmod +x /app/app

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

EXPOSE 8080

CMD ["/app/app"]
