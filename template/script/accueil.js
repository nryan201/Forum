document.getElementById("commentForm").onsubmit = function(event){
    event.preventDefault(); // 
    let FormData = new FormData(this);
    let imageFile = document.getElementById("imageUpload").files[0];

    if (imageFile && imageFile.size > 20971520) { // 20 Mo = 20971520 octets
        alert("Image trop lourde");
        return;
    }

    fetch("api/commentaire", {
        method: "POST",
        body: FormData,
    })
    ,then(response => {
        if (!response.ok) {
            throw new Error("Erreur lors du telechargement");
        }
        return response.json();
    })
    .then(data => {
        console.log(data);
        alert("Commentaire ajouté");
        document.getElementById("commentForm").reset();
    })
    .catch(error => {
        console.error("Erreur lors de l'ajout du commentaire", error);
        alert("Erreur lors de l'ajout du commentaire");
    });
};


// possbility to scroll 

