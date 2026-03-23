# Ombrasoft

Ombrasoft est une application web full-stack permettant la gestion de **bookmarks** et de **séries**, avec un backend en **Go** et un frontend en **React**. Elle supporte la gestion des utilisateurs, la visualisation de séries, et l’ajout/suppression de favoris.

## Fonctionnalités

- Gestion des utilisateurs  
- Ajout et suppression de **bookmarks** et **séries**  
- Affichage dynamique des séries avec leurs images correspondantes  
- Support multi-plateforme grâce à Docker  
- Frontend réactif en **React** avec gestion des états via hooks  
- Backend robuste en **Go**, utilisant **UUID** pour les identifiants  

## Prérequis

Avant de lancer le projet, assurez-vous d’avoir installé :  
- [Docker](https://www.docker.com/) et [Docker Compose](https://docs.docker.com/compose/)  
- [Go](https://golang.org/dl/) pour le backend  
- [Node.js](https://nodejs.org/) et [npm](https://www.npmjs.com/) pour le frontend  

## Installation et lancement

### Lancer le projet avec Docker
docker-compose up --build

Construire et pousser l'image backend multi-architecture
docker buildx build --platform linux/amd64,linux/arm64 -t ombrasoft-backend ./backend --push
Accéder à l'application
Frontend : http://localhost:5173
Backend : http://localhost:8080

Utilisation
Ajouter un bookmark : via le formulaire sur le frontend
Supprimer un bookmark : cliquer sur l’icône de suppression
Visualiser les séries : la page principale affiche toutes les séries liées aux bookmarks
Chaque série possède sa propre image et détails
Développement
Backend

API
GET /bookmarks : récupérer tous les bookmarks
POST /bookmarks : créer un bookmark
DELETE /bookmarks/:id : supprimer un bookmark et ses séries associées
GET /series : récupérer toutes les séries
POST /series : créer une série

Docker
Le projet est dockerisé pour un déploiement simple
Docker Compose gère le backend et le frontend simultanément
Support multi-architecture pour le backend


MIT License

