# Utiliser une image Node.js pour le frontend
FROM node:20-alpine

WORKDIR /app

# Copier package.json et package-lock.json
COPY package*.json ./

# Installer les dépendances
RUN npm install

# Copier tout le code source
COPY . .

# Exposer le port utilisé par Vite (généralement 5173)
EXPOSE 5173

# Lancer l'application avec npm run dev (mode développement)
CMD ["npm", "run", "dev"]
