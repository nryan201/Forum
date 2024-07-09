    document.addEventListener('DOMContentLoaded', function() {

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
 
});
