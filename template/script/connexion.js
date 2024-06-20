document.addEventListener('DOMContentLoaded', function() {
    var loginForm = document.getElementById('loginForm');
    var registerForm = document.getElementById('registerForm');
    var forgotPassword = document.getElementById('forgotPassword');

    // Fonction pour afficher le formulaire d'enregistrement
    function showRegisterForm(e) {
        e.preventDefault();
        loginForm.style.display = 'none';
        forgotPassword.style.display = 'none';
        registerForm.style.display = 'block';
    }

    // Fonction pour afficher le formulaire de connexion
    function showLoginForm(e) {
        e.preventDefault();
        registerForm.style.display = 'none';
        forgotPassword.style.display = 'none';
        loginForm.style.display = 'block';
    }

    // Fonction pour afficher le formulaire de mot de passe oublié
    function showForgotPasswordForm(e) {
        e.preventDefault();
        loginForm.style.display = 'none';
        registerForm.style.display = 'none';
        forgotPassword.style.display = 'block';
    }

    // Attacher les écouteurs d'événements
    function attachEventListeners() {
        var registerLink = document.querySelector('a[href="#registerForm"]');
        var loginLink = document.querySelector('a[href="#loginForm"]');
        var forgotPasswordLink = document.querySelector('a[href="#ForgotPassword"]');

        if (registerLink) {
            registerLink.addEventListener('click', showRegisterForm);
        }
        if (loginLink) {
            loginLink.addEventListener('click', showLoginForm);
        }
        if (forgotPasswordLink) {
            forgotPasswordLink.addEventListener('click', showForgotPasswordForm);
        }
        if (loginLink) {
            loginLink.addEventListener('click', showLoginForm);
        }
    }

    // Initialiser les écouteurs d'événements
    attachEventListeners();
});