# Ombrasoft

Ombrasoft est une application web qui permet de sauvegarder et organiser des liens vers des séries ainsi que de gérer ses favoris.  

Elle se compose d’un site web interactif utilisé par l’utilisateur, et d’un système en arrière-plan qui s’occupe de stocker et gérer les données de manière fiable. L’ensemble fonctionne ensemble pour offrir une expérience fluide, avec la possibilité de créer un compte, ajouter des favoris et consulter des séries avec leurs informations.

## Fonctionnalités

- Gestion des utilisateurs  
- Ajout et suppression de **bookmarks**  
- Gestion et affichage des **séries**  
- Visualisation dynamique des séries avec leurs images  
- Architecture **full-stack**  
- Frontend réactif en **React** (hooks pour la gestion d’état)  
- Backend robuste en **Go** avec utilisation d’**UUID** pour les identifiants  
- Conteneurisation avec **Docker** et **Docker Compose**  
- Support multi-plateforme (AMD64 / ARM64)  

## Prérequis

Avant de lancer le projet, assurez-vous d’avoir installé :

- Docker : https://www.docker.com/  
- Docker Compose : https://docs.docker.com/compose/  

Optionnel (pour développement local sans Docker) :

- Go : https://golang.org/dl/  
- Node.js : https://nodejs.org/  
- npm : https://www.npmjs.com/  

## Configuration

Le projet nécessite la configuration d’un fichier `.env`.

- Ne pas utiliser directement le fichier `.env.example`
- Ce fichier sert uniquement de modèle pour comprendre les variables attendues
- Vous devez créer votre propre fichier `.env` à la racine du projet

Sans cette configuration, l’application ne fonctionnera pas correctement.

## Installation et lancement

### Lancer le projet avec Docker

```bash
docker-compose up --build
```

### Build et push de l’image backend multi-architecture
```Bash
docker buildx build --platform linux/amd64,linux/arm64 -t ombrasoft-backend ./backend --push
```
### Accès à l'application
Frontend : http://localhost:5173
Backend : http://localhost:8080

### Utilisation
Ajouter un bookmark
    Via le formulaire disponible sur le frontend
Supprimer un bookmark
    Cliquer sur l’icône de suppression associée
Consulter les séries
    La page principale affiche toutes les séries liées aux bookmarks
    Chaque série contient une image et des informations détaillées


### important les Source des séries actuellement pris en charge 

https://raijin-scans.fr/

### Exemple d’utilisation :
Aller sur une série, par exemple :
    https://raijin-scans.fr/manga/nano-machine-1/
Copier l’URL de la série
    L’ajouter dans l’application via le formulaire de bookmark


### API Backend
Bookmarks
    GET /bookmarks : récupérer tous les bookmarks
    POST /bookmarks : créer un bookmark
    DELETE /bookmarks/:id : supprimer un bookmark (et ses séries associées)
Séries
    GET /series : récupérer toutes les séries
    POST /series : créer une série


# licence

Ce projet est distribué sous la licence **Apache License 2.0**.

Vous êtes libre d’utiliser, modifier et redistribuer ce projet conformément aux termes de cette licence.