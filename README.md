# Forum Project

## Description

Ce projet est un forum web développé en Go (Golang) pour le backend, avec des interfaces utilisateur en HTML, CSS et JavaScript. Il permet aux utilisateurs de créer des comptes, de se connecter, de poster des sujets et des commentaires, et d'interagir via une messagerie. Le projet utilise Docker pour la conteneurisation, assurant une installation et un déploiement simplifiés.

## Fonctionnalités

- Création de compte utilisateur
- Connexion et déconnexion
- Connexion et inscription via GitHub, Google et Facebook
- Publication de sujets et de commentaires
- Gestion des hashtags
- Messagerie privée entre utilisateurs
- Administration et modération des contenus
- Support de l'authentification via des cookies

## Installation

### Prérequis

- [Docker](https://www.docker.com/products/docker-desktop) installé sur votre machine

### Étapes d'installation

1. Clonez le dépôt du projet :

    ```sh
    git clone https://github.com/nryan201/Forum.git
    cd Forum
    ```

2. Initialisez et démarrez les conteneurs Docker :

    ```sh
    docker-compose up --build
    ```

3. Initialisez le certificat SSL avec Let's Encrypt (à exécuter une seule fois) :

    ```sh
    docker-compose run --rm certbot certonly --webroot --webroot-path=/var/www/certbot --email your-email@domain.com --agree-tos --no-eff-email -d yourdomain.com
    ```

4. Redémarrez les conteneurs Docker pour appliquer le certificat :

    ```sh
    docker-compose up --build
    ```

## Utilisation

### Accéder au forum

Ouvrez votre navigateur et allez à `https://yourdomain.com`

### Connexion et inscription

Utilisez les pages de connexion et d'inscription pour créer un compte ou vous connecter à votre compte existant. Vous pouvez également vous inscrire et vous connecter via GitHub, Google et Facebook pour plus de commodité.

### Publier un sujet ou un commentaire

Une fois connecté, vous pouvez créer de nouveaux sujets ou commenter sur des sujets existants.

### Utiliser la messagerie

Accédez à la section messagerie pour envoyer et recevoir des messages privés.

## Structure du Projet

```plaintext
/Forum
|-- back/
|   |-- .go/
|-- data/
|   |-- Forum.sql/
|-- Dockerfile
|-- docker-compose.yml
|-- nginx/
|   |-- nginx.conf
 main.go
|-- template/
|   |-- ressource/
|       |-- images/
|   |-- css/
|   |-- html/
|   |-- script/
|-- db.sqlite

```

## Contributeurs

Nous sommes une équipe de cinq développeurs qui ont collaboré pour créer ce projet :

- Ryan NEDJARI
- Alexis REDAUD
- Paul AMALRIC
- Louis MARGRAS
- Kevin CARLES

Merci à tous pour leurs efforts toutes au long de se projets aussi éprouvant qu'enrichissant.
