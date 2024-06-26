document.addEventListener("DOMContentLoaded", () => {
    const topicContainer = document.getElementById('topic-container');

    // Fonction pour créer une box de topic
    function createTopicBox(topic) {
        const postFeed = document.createElement('div');
        postFeed.className = 'post-feed';

        const postAuthor = document.createElement('div');
        postAuthor.className = 'post-author';
        postAuthor.innerHTML = `Auteur: <span class="author-name">${topic.author}</span>`;

        const postContent = document.createElement('div');
        postContent.className = 'postContent';
        postContent.innerHTML = `Post: <span class="content-Post">${topic.content}</span>`;

        postFeed.appendChild(postAuthor);
        postFeed.appendChild(postContent);

        return postFeed;
    }

    // Données fictives pour les topics
    const topics = [
        { id: 1, author: "Utilisateur 1", content: "Contenu du topic 1" },
        { id: 2, author: "Utilisateur 2", content: "Contenu du topic 2" },
        { id: 3, author: "Utilisateur 3", content: "Contenu du topic 3" }
    ];

    // Simulation de la récupération des topics
    if (topicContainer && createTopicBox) {
        topics.forEach(topic => {
            const topicBox = createTopicBox(topic);
            const topicWrapper = document.createElement('div');
            topicWrapper.className = 'post-topic';
            topicWrapper.appendChild(topicBox);
            topicContainer.appendChild(topicWrapper);
        });
    } else {
        console.error("topicContainer or createTopicBox is not defined");
    }
});
/*// Récupération des topics depuis l'API
    fetch('/api/topics')
        .then(response => response.json())
        .then(topics => {
            topics.forEach(topic => {
                const topicBox = createTopicBox(topic);
                topicContainer.appendChild(topicBox);
            });
        })
        .catch(error => console.error('Erreur:', error));*/