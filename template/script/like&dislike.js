// code pour le like
document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('.like-button').forEach(button => {
        button.addEventListener('click', function(event) {
            const topicID = this.dataset.topicId;
            const userID = this.dataset.userId;

            fetch('/addLike', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `topicID=${topicID}&userID=${userID}`
            })
            .then(response => response.json())
            .then(data => {
                if (data.status === "like added") {
                    // Update the like count in the UI
                    const likeCount = document.querySelector(`#like-count-${topicID}`);
                    likeCount.textContent = parseInt(likeCount.textContent) + 1;
                }
            })
            .catch(error => console.error('Error:', error));
        });
    });
});

// code pour le dislike
document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('.dislike-button').forEach(button => {
        button.addEventListener('click', function(event) {
            const topicID = this.dataset.topicId;
            const userID = this.dataset.userId;

            fetch('/addDislike', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `topicID=${topicID}&userID=${userID}`
            })
            .then(response => response.json())
            .then(data => {
                if (data.status === "dislike added") {
                    // Update the dislike count in the UI
                    const dislikeCount = document.querySelector(`#dislike-count-${topicID}`);
                    dislikeCount.textContent = parseInt(dislikeCount.textContent) + 1;
                }
            })
            .catch(error => console.error('Error:', error));
        });
    });
});
