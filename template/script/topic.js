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
        .then(topics => {
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
    topics.forEach(topic => {
        const topicBox = createTopicBox(topic);
        const topicWrapper = document.createElement('div');
        topicWrapper.className = 'post-topic'; // Assurez-vous que cette classe correspond dans le CSS
        topicWrapper.appendChild(topicBox);
        topicContainer.appendChild(topicWrapper);
    });

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