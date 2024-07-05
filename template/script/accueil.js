document.addEventListener('DOMContentLoaded', function () {
    var commentForm = document.getElementById("commentForm");
    if (commentForm) {
        commentForm.onsubmit = function (event) {
            event.preventDefault();
            let formData = new FormData(this);
            let imageFile = document.getElementById("imageUpload").files[0];

            if (imageFile && imageFile.size > 20971520) { // 20 Mo = 20971520 octets
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
                    commentForm.reset();
                })
                .catch(error => {
                    console.error("Erreur lors de l'ajout du commentaire", error);
                    alert("Erreur lors de l'ajout du commentaire");
                });
        };
    } else {
        console.error('Formulaire de commentaire non trouvé');
    }
});


