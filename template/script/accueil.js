document.addEventListener('DOMContentLoaded', function() {
    const topicContainer = document.getElementById('topic-container');
    const searchInput = document.getElementById('searchInput');
    const searchForm = document.getElementById('searchForm');

    // Function to create HTML structure for each topic
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

        topicBox.addEventListener('click', () => {
            window.location.href = '/post?id=' + topic.id;
        });

        return topicBox;
    }

    // Function to fetch topics based on search input
    function fetchTopics(query) {
        const url = query ? `/search?search=${encodeURIComponent(query)}` : '/topics'; // Adjust the URL based on search
        fetch(url)
            .then(response => response.json())
            .then(topics => {
                topicContainer.innerHTML = ''; // Clear the container before adding new topics
                topics.forEach(topic => {
                    const box = createTopicBox(topic);
                    topicContainer.appendChild(box);
                });
            })
            .catch(error => {
                console.error('Failed to fetch topics:', error);
                topicContainer.innerHTML = '<p>Error loading topics.</p>';
            });
    }

    // Event listener for form submission
    searchForm.addEventListener('submit', function(event) {
        event.preventDefault(); // Prevent the default form submission
        fetchTopics(searchInput.value); // Fetch topics based on the search input
    });

    // Initial fetch of all topics or based on query if available
    fetchTopics();
});
