document.getElementById('profile-upload').addEventListener('change', function() {
    const fileInput = this;
    if (fileInput.files && fileInput.files[0]) {
        const formData = new FormData();
        formData.append('profile', fileInput.files[0]);

        fetch('/uploadProfilePic', {
            method: 'POST',
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    document.getElementById('profile-image').src = data.filePath;
                } else {
                    alert('Failed to upload image.');
                }
            })
            .catch(error => console.error('Error uploading image:', error));
    }
});

document.getElementById('profile-pic').addEventListener('mouseover', function() {
    document.getElementById('download-icon').style.display = 'block';
});

document.getElementById('profile-pic').addEventListener('mouseout', function() {
    document.getElementById('download-icon').style.display = 'none';
});
