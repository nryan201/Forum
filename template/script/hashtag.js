document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('createHashtagForm');
    const konamiCode = ['ArrowUp', 'ArrowUp', 'ArrowDown', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'ArrowLeft', 'ArrowRight', 'b', 'a', 'Enter'];
    let currentInput = [];

    if (form) {
        form.addEventListener('submit', function (event) {
            event.preventDefault();
            let hashtagName = document.querySelector('.createHashtags').value.trim();
            if (hashtagName !== '') {
                if (!hashtagName.startsWith('#')) {
                    hashtagName = '#' + hashtagName;
                }

                const newHashtag = document.createElement('li');
                newHashtag.textContent = hashtagName;
                document.getElementById('hashtagsList').appendChild(newHashtag);
                document.querySelector('.createHashtags').value = '';

                // Send form data using AJAX
                let formData = new FormData(event.target);
                fetch('/CreateCategory', { method: 'POST', body: formData })
                    .then(response => response.json())
                    .then(data => console.log(data))
                    .catch(error => console.error('Error:', error));
            }
        });
    }

    document.addEventListener('keydown', (event) => {
        currentInput.push(event.key);
        if (currentInput.length > konamiCode.length) {
            currentInput.shift();
        }

        if (currentInput.join('') === konamiCode.join('')) {
            launchConfetti();
            alert('Konami Code activated!');
        }
    });

    function launchConfetti() {
        confetti({
            particleCount: 100,
            spread: 70,
            origin: { y: 0.6 }
        });
    }
});
