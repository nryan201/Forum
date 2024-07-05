document.addEventListener('DOMContentLoaded', function() {
    const topicContainer = document.getElementById('topic-container');

    function createTopicBox(topic) {
        // Create outer div for the topic with class 'post-topic'
        const topicBox = document.createElement('div');
        topicBox.className = 'post-topic';

        // Create inner div for the feed with class 'post-feed'
        const postFeed = document.createElement('div');
        postFeed.className = 'post-feed';

        // Create title element
        const title = document.createElement('h3');
        title.textContent = topic.title;
        title.style.fontSize = '24px'; // Example of making the text larger

        // Create description paragraph
        const description = document.createElement('p');
        description.textContent = topic.description;
        description.style.fontSize = '18px'; // Larger text for description

        // Append title and description to the 'post-feed' div
        postFeed.appendChild(title);
        postFeed.appendChild(description);

        // Append 'post-feed' to 'post-topic'
        topicBox.appendChild(postFeed);

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
