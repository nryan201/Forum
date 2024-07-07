document.addEventListener("DOMContentLoaded", function () {
    fetch("/api/check-auth")
        .then(function (response) {
            console.log("Received response from /api/check-auth");
            return response.json();
        })
        .then(function (data) {
            console.log("Data received from /api/check-auth:", data);
            const button = document.getElementById('auth-button');
            if (button) {
                console.log("Auth button found:", button);
                if (data.authenticated) {
                    console.log("User is authenticated");
                    button.textContent = "Profile";
                    button.href = "/profil";
                } else {
                    console.log("User is not authenticated");
                    button.textContent = "Connexion";
                    button.href = "/connexion";
                }
            } else {
                console.error("Auth button not found");
            }
        })
        .catch(function (error) {
            console.error("Error fetching /api/check-auth:", error);
        });
});