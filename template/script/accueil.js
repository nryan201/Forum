document.getElementById("commentForm").onsubmit = function(event) {
    event.preventDefault();
    let formData = new FormData(this);
    let imageFile = document.getElementById("imageUpload").files[0];

    // Check if the image file is too large (more than 20 MB)
    if (imageFile && imageFile.size > 20971520) {
        alert("Image trop lourde");
        return;
    }

    fetch("api/commentaire", {
        method: "POST",
        body: formData,
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Erreur lors du téléchargement");
            }
            return response.json();
        })
        .then(data => {
            console.log(data);
            alert("Commentaire ajouté");
            document.getElementById("commentForm").reset(); // Reset the form after successful submission
        })
        .catch(error => {
            console.error("Erreur lors de l'ajout du commentaire", error);
            alert("Erreur lors de l'ajout du commentaire");
        });
};
