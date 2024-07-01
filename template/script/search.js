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


function displayResults(results) {
    const resultsContainer = document.getElementById('searchResults');
    resultsContainer.innerHTML = ''; // Clear previous results
    results.forEach(result => {
        const div = document.createElement('div');
        div.textContent = result.title || result.username;
        resultsContainer.appendChild(div);
    });
}