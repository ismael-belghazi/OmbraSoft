# Étape 1 : Build
FROM golang:1.25-alpine AS builder

# Répertoire de travail
WORKDIR /app

# Installer les certificats CA pour HTTPS
RUN apk add --no-cache ca-certificates git

# Copier les fichiers de modules et télécharger les dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier tout le code source
COPY . .

# Compiler l'application en statique
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd

# Étape 2 : Runtime léger
FROM alpine:latest

# Installer les certificats CA pour HTTPS
RUN apk add --no-cache ca-certificates bash

WORKDIR /app

# Copier l'exécutable depuis l'étape builder
COPY --from=builder /app/app .

# Rendre l'exécutable exécutable
RUN chmod +x /app/app

# Ajouter un healthcheck pour Render
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Exposer le port par défaut
EXPOSE 8080

# Commande de lancement
CMD ["/app/app"]
