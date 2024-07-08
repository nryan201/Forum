# Utiliser l'image officielle de Go comme base
FROM golang:alpine

# Installer git pour télécharger les dépendances
RUN apk add --no-cache git

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers go.mod et go.sum pour le téléchargement des dépendances
COPY src/go.mod src/go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier les fichiers du projet
COPY src/ .

# Construire l'application
RUN go build -o main .

# Exposer le port utilisé par l'application
EXPOSE 8080

# Commande pour lancer l'application
CMD ["./main"]
