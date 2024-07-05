document.addEventListener('DOMContentLoaded', function() {
    var loginForm = document.getElementById('loginForm');
    var registerForm = document.getElementById('registerForm');
    var forgotPassword = document.getElementById('forgotPassword');

    // Fonctions pour afficher les formulaires
    function showRegisterForm(e) {
        e.preventDefault();
        loginForm.style.display = 'none';
        forgotPassword.style.display = 'none';
        registerForm.style.display = 'block';
    }

    function showLoginForm(e) {
        e.preventDefault();
        registerForm.style.display = 'none';
        forgotPassword.style.display = 'none';
        loginForm.style.display = 'block';
    }

    function showForgotPasswordForm(e) {
        e.preventDefault();
        loginForm.style.display = 'none';
        registerForm.style.display = 'none';
        forgotPassword.style.display = 'block';
    }

    // Attachement des écouteurs d'événements
    function attachEventListeners() {
        var registerLink = document.querySelector('a[href="#registerForm"]');
        var loginLink = document.querySelector('a[href="#loginForm"]');
        var forgotPasswordLink = document.querySelector('a[href="#ForgotPassword"]');
        var loginLink2 = document.querySelector('a[href="#loginForm2"]');

        if (registerLink) {
            registerLink.addEventListener('click', showRegisterForm);
        }
        if (loginLink) {
            loginLink.addEventListener('click', showLoginForm);
        }
        if (forgotPasswordLink) {
            forgotPasswordLink.addEventListener('click', showForgotPasswordForm);
        }
        if (loginLink2) {
            loginLink2.addEventListener('click', showLoginForm);
        }
    }

    // Initialisation des écouteurs d'événements
    attachEventListeners();

    // Vérification de l'état d'authentification au chargement de la page
    checkAuthStatus();
});
