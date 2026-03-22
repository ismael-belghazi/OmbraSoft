
# LANCER PROJECT

docker-compose up --build

docker buildx build --platform linux/amd64,linux/arm64 -t ombrasoft-backend ./backend --push

# aller sur le site:
http://localhost:5173