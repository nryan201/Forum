document.addEventListener('DOMContentLoaded', function() {
    const topicContainer = document.getElementById('topic-container');

    function createTopicBox(topic) {
        const topicBox = document.createElement('div');
        topicBox.className = 'topic';

        const title = document.createElement('h3');
        title.textContent = topic.title;
        title.style.fontSize = '24px'; // Example of making the text larger

        const description = document.createElement('p');
        description.textContent = topic.description;
        description.style.fontSize = '18px'; // Larger text for description

        topicBox.appendChild(title);
        topicBox.appendChild(description);

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
