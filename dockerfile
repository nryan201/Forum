# Utiliser l'image officielle de Go comme base
FROM golang:alpine

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers du projet
COPY . .

# Construire l'application
RUN go build -o main .

# Exposer le port utilisé par l'application
EXPOSE 8080

# Commande pour lancer l'application
CMD ["./main"]
