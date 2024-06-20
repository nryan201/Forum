document.addEventListener('DOMContentLoaded', function() {
    var loginForm = document.getElementById('loginForm');
    var registerForm = document.getElementById('registerForm');

    function showRegisterForm(e) {
        e.preventDefault();
        loginForm.style.display = 'none';
        registerForm.style.display = 'block';
        attachLoginLinkListener();
    }

    function showLoginForm(e) {
        e.preventDefault();
        loginForm.style.display = 'block';
        registerForm.style.display = 'none';
        attachRegisterLinkListener();
    }

    function attachRegisterLinkListener() {
        var registerLink = document.querySelector('.register-link a[href="#registerForm"]');
        if (registerLink) {
            registerLink.addEventListener('click', showRegisterForm);
            console.log('Register link listener attached');
        } else {
            console.log('Register link not found.');
        }
    }

    function attachLoginLinkListener() {
        var loginLink = document.querySelector('.login-link a[href="#loginForm"]');
        if (loginLink) {
            loginLink.addEventListener('click', showLoginForm);
            console.log('Login link listener attached');
        } else {
            console.log('Login link not found.');
        }
    }

    // Attach initial listeners
    attachRegisterLinkListener();
    attachLoginLinkListener();
});
