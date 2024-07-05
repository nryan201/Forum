document.addEventListener('DOMContentLoaded', function() {
    const topicContainer = document.getElementById('topic-container');
    let currentPage = 1;
    const pageSize = 5;

    function loadTopics(page) {
        fetch(`/api/topics?page=${page}&pageSize=${pageSize}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                if (!Array.isArray(data)) {
                    console.error('Received data is not an array:', data);
                    return;
                }
                topicContainer.innerHTML = ''; // Clear previous topics
                data.forEach(topic => {
                    const topicBox = createTopicBox(topic);
                    topicContainer.appendChild(topicBox);
                });
                updatePaginationControls(data.length);
            })
            .catch(error => {
                console.error('Error fetching topics:', error);
                topicContainer.innerHTML = '<p>Failed to load topics.</p>'; // Display error message
            });
    }

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

    function updatePaginationControls(numItems) {
        const nextButton = document.getElementById('nextButton');
        const prevButton = document.getElementById('prevButton');
        if (numItems < pageSize) {
            nextButton.style.display = 'none'; // Hide next button if there are no more items
        } else {
            nextButton.style.display = 'inline';
        }
        if (currentPage === 1) {
            prevButton.style.display = 'none';
        } else {
            prevButton.style.display = 'inline';
        }
    }

    document.getElementById('nextButton').addEventListener('click', () => {
        currentPage++;
        loadTopics(currentPage);
    });

    document.getElementById('prevButton').addEventListener('click', () => {
        currentPage--;
        loadTopics(currentPage);
    });

    loadTopics(currentPage); // Load initial set of topics
});
