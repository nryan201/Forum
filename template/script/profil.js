
// Exemple de récupération de données depuis le backend Go (utilisation de fetch pour simplifier)
fetch('/api/user/profile') // Assurez-vous d'adapter l'URL à votre API backend
    .then(response => response.json())
    .then(data => {
        // Remplir les éléments HTML avec les données récupérées
        document.getElementById('profile-pic').innerText = data.profilePicUrl; // Exemple de champ profilePicUrl dans les données
        document.getElementById('username').innerText = `Username: ${data.username}`;
        document.getElementById('name').innerText = `Name: ${data.name}`;
        document.getElementById('firstname').innerText = `Firstname: ${data.firstname}`;
        document.getElementById('email').innerText = `Mail: ${data.email}`;
        document.getElementById('other-link').innerText = `Other link: ${data.otherLink}`;

        // Exemple pour les liens supplémentaires
        document.getElementById('discord-link').href = data.discordLink; // Assurez-vous que data.discordLink contient l'URL correcte
        document.getElementById('steam-link').href = data.steamLink;
        document.getElementById('riot-client-link').href = data.riotClientLink;
        document.getElementById('psn-link').href = data.psnLink;

        document.getElementById('birthdate-section').innerText = `Date de naissance: ${data.birthdate}`;
    })
    .catch(error => {
        console.error('Error fetching profile data:', error);
    });
