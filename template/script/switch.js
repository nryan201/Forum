document.addEventListener("DOMContentLoaded", function () {
    fetch("/api/check-auth")
        .then(response => response.json())
        .then(data => {
            const authButton = document.getElementById('auth-button');
            const createPostButton = document.getElementById('creation-button');
            const messageButton = document.getElementById('Messages');

            if (data.authenticated) {
                // Mise à jour du bouton d'authentification
                if (authButton) {
                    authButton.textContent = "Profile";
                    authButton.href = "/profil";
                }
                // Afficher le bouton de création de post
                if (createPostButton) {
                    createPostButton.style.display = 'inline-block';
                }
                // Afficher le bouton de messagerie
                if (messageButton) {
                    messageButton.style.display = 'inline-block';
                }
            } else {
                // Mise à jour du bouton d'authentification
                if (authButton) {
                    authButton.textContent = "Connexion";
                    authButton.href = "/connexion";
                }
                // Masquer le bouton de création de post
                if (createPostButton) {
                    createPostButton.style.display = 'none';
                }
                // Masquer le bouton de messages
                if (messageButton) {
                    messageButton.style.display = 'none';
                }
            }
        })
        .catch(error => console.error("Erreur lors de la vérification d'authentification:", error));
});