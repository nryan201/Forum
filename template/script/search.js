document.getElementById('searchForm').addEventListener('submit', function (event) {
    event.preventDefault(); // Prevent form submission
    const searchValue = document.getElementById('searchInput').value;
    fetch(`/search?query=${encodeURIComponent(searchValue)}`)
        .then(response => response.json())
        .then(data => {
            displayResults(data);
        })
        .catch(error => console.error('Error fetching search results:', error));
});


// func for search topics in the search bar
function searchTopics() {
    const query = document.getElementById('searchInput').value;
    fetch(`/post?query=${encodeURIComponent(query)}`)
        .then(response => response.json())
        .then(data => {
            const resultsDiv = document.getElementById('results');
            resultsDiv.innerHTML = '';
            data.forEach(topic => {
                const topicDiv = document.createElement('div');
                topicDiv.innerHTML = `<h2>${topic.Title}</h2><p>${topic.Description}</p>`;
                resultsDiv.appendChild(topicDiv);
            });
        })
        .catch(error => console.error('Error fetching topics:', error));
}


function displayResults(results) {
    const resultsContainer = document.getElementById('searchResults');
    resultsContainer.innerHTML = ''; // Clear previous results
    results.forEach(result => {
        const div = document.createElement('div');
        div.textContent = result.title || result.username;
        resultsContainer.appendChild(div);
    });
}

