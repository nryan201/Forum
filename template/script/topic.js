document.addEventListener('DOMContentLoaded', function() {
    const topicContainer = document.getElementById('topic-container');

    function createTopicBox(topic) {
        const topicBox = document.createElement('div');
        topicBox.className = 'post-topic';

        const postFeed = document.createElement('div');
        postFeed.className = 'post-feed';

        const title = document.createElement('h3');
        title.textContent = topic.title;
        title.className = 'post-title';

        const description = document.createElement('p');
        description.textContent = topic.description;
        description.className = 'post-description';

        postFeed.appendChild(title);
        postFeed.appendChild(description);
        topicBox.appendChild(postFeed);

        // Adding event listener to the entire topic box
        topicBox.addEventListener('click', function() {
            window.location.href = '/post?id=' + topic.id; // Redirect to the detailed post view
        });

        return topicBox;
    }

    function fetchTopics() {
        fetch('/topics')
            .then(response => response.json())
            .then(topics => {
                topics.forEach(topic => {
                    const box = createTopicBox(topic);
                    topicContainer.appendChild(box);
                });
            })
            .catch(error => console.error('Failed to fetch topics:', error));
    }

    fetchTopics();
});
