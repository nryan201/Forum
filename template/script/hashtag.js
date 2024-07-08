document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('createHashtagForm');
    const konamiCode = ['ArrowUp', 'ArrowUp', 'ArrowDown', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'ArrowLeft', 'ArrowRight', 'b', 'a', 'Enter'];
    let currentInput = [];

    fetch('/hashtags')
        .then(response => response.json())
        .then(data => {
            const hashtagsList = document.getElementById('hashtagsList');
            hashtagsList.innerHTML = ''; // Clear existing hashtags
            data.forEach(hashtag => {
                const li = document.createElement('li');
                li.textContent = hashtag;
                hashtagsList.appendChild(li);
            });
        })
        .catch(error => console.error('Error:', error));

    if (form) {
        form.addEventListener('submit', function(event) {
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
                let formData = new FormData();
                formData.append('name', hashtagName);
                fetch('/addHashtag', { method: 'POST', body: formData })
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
            window.location.href = 'https://www.youtube.com/watch?v=dQw4w9WgXcQ';
        }
    });

    function launchConfetti() {
        if (typeof confetti === 'function') {
            confetti({
                particleCount: 100,
                spread: 70,
                origin: { y: 0.6 }
            });
        } else {
            console.log('Confetti function is not defined.');
        }
    }

});
