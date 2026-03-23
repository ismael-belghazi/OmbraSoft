# ===== BUILD =====
FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/main.go

# Installer playwright
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest

# ===== RUNTIME =====
FROM debian:bookworm-slim

WORKDIR /app

# Dépendances pour Playwright
RUN apt-get update && apt-get install -y \
    ca-certificates \
    wget \
    curl \
    fonts-liberation \
    libnss3 \
    libatk1.0-0 \
    libatk-bridge2.0-0 \
    libcups2 \
    libdrm2 \
    libxkbcommon0 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    libgbm1 \
    libasound2 \
    libpangocairo-1.0-0 \
    libpango-1.0-0 \
    libcairo2 \
    libatspi2.0-0 \
    libgtk-3-0 \
    && rm -rf /var/lib/apt/lists/*

# Copier playwright depuis le builder
COPY --from=builder /go/bin/playwright /usr/local/bin/playwright

RUN playwright install --with-deps

ENV PLAYWRIGHT_BROWSERS_PATH=/root/.cache/ms-playwright

# Copier l'application
COPY --from=builder /app/app .

# Créer le dossier covers si absent
RUN mkdir -p ./public/covers

# Exposer le port de ton backend
EXPOSE 8080

CMD ["./app"]