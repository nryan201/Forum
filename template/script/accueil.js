document.getElementById("commentForm").onsubmit = function(event){
    event.preventDefault();
    let formData = new FormData(this);
    let imageFile = document.getElementById("imageUpload").files[0];

    if (imageFile && imageFile.size > 20971520) { // 20 Mo = 20971520 octets
        alert("Image trop grande");
        return;
    }

    fetch('/api/comment', {
        method: 'POST',
        body: formData,
    })
    then(response => response.json())
    then(data => console.log(data))
    .catch(error => console.error(error));
};