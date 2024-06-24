document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('createHashtagForm').addEventListener('submit', function(event) {
        event.preventDefault();

        let hashtagName = document.querySelector('.createHashtags').value.trim();

        if (hashtagName !== '') {
            if (!hashtagName.startsWith('#')) {
                hashtagName = '#' + hashtagName;
            }

            const newHashtag = document.createElement('li');
            newHashtag.textContent = hashtagName;

            document.getElementById('hashtagsList').appendChild(newHashtag);

            document.querySelector('.createHashtags').value = ''; // Clear the input field

            // Send form data using AJAX
            let formData = new FormData(event.target);
            fetch('/CreateCategory', { method: 'POST', body: formData })
                .then(response => response.json())
                .then(data => console.log(data));
        }
    });
});
