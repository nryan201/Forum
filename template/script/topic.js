document.addEventListener('DOMContentLoaded', function() {
    const topicContainer = document.getElementById('topic-container');

    // function to create a topic box
    function createTopicBox(topic) {
        const postFeed = document.createElement('div');
        postFeed.className = 'post-feed';

        const postAuthor = document.createElement('div');
        postAuthor.className = 'post-author';
        postAuthor.innerHTML = `Auteur: <span class="author-name">${topic.author}</span>`;

        const postContent = document.createElement('div');
        postContent.className = 'post-content';
        postContent.innerHTML = `Post: <span class="content-post">${topic.content}</span>`;

        postFeed.appendChild(postAuthor);
        postFeed.appendChild(postContent);

        return postFeed;
    }

    // fetch topics data
    fetch('/bdd/topics') // Assurez-vous que l'URL est correcte
        .then(response => response.json())
        .then(topic => {
            topics.forEach(topic => {
                const topicBox = createTopicBox(topic);
                topicContainer.appendChild(topicBox);
            });
        })
        .catch(error => console.error('Erreur :', error));



    // topics data
    const topics = [
        { id: 1, author: "Utilisateur 1", content: "Contenu du topic 1" },
        { id: 2, author: "Utilisateur 2", content: "Contenu du topic 2" },
        { id: 3, author: "Utilisateur 3", content: "Contenu du topic 3" },
        { id: 4, author: "Utilisateur 4", content: "Contenu du topic 4" }
    ];

    // topic creation
    const routes = [
        "/topic"
    ];

    // Simulation de la récupération des topics
    if (topicContainer && createTopicBox) {
        topics.forEach(topic => {
            const topicBox = createTopicBox(topic);
            const topicWrapper = document.createElement('div');
            topicWrapper.className = 'post-topic';
            topicWrapper.appendChild(topicBox);
            topicContainer.appendChild(topicWrapper);
            topicBox.addEventListener('click', () => {
                window.location.href = routes;
            });
        });
    } else {
        console.error("topicContainer or createTopicBox is not defined");
    }
    // like, dislike and share buttons
    document.addEventListener('click', function (event) {
        if (event.target.classList.contains('like-button')) {
            const postFeed = event.target.closest('.post-feed');
            const likeCount = postFeed.querySelector('.like-count');
            likeCount.textContent = parseInt(likeCount.textContent) + 1;
        } else if (event.target.classList.contains('dislike-button')) {
            const postFeed = event.target.closest('.post-feed');
            const dislikeCount = postFeed.querySelector('.dislike-count');
            dislikeCount.textContent = parseInt(dislikeCount.textContent) - 1;
        } else if (event.target.classList.contains('share-button')) {
            const postFeed = event.target.closest('.post-feed');
            const shareCount = postFeed.querySelector('.share-count');
            shareCount.textContent = parseInt(shareCount.textContent) + 1;
        }
    });
});

//Fonction pour récupérer et afficher les données du post
function fetchPostData() {
    $.ajax({
        url: 'https://localhost:443/topic',
        type: 'GET',
        success: function(data) {
            $('.post-title').text(data.title);
            $('.postContent').text(data.content);
            $('.author-name').text(data.author);
        },
        error: function() {
            alert('Erreur lors de la récupération des données du post');
        }
    });
}

// Exécute la fonction de récupération des données au chargement du document
$(document).ready(function() {
    fetchPostData(); // Appel initial pour charger les données du post
});
