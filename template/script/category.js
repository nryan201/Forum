document.addEventListener('DOMContentLoaded', function() {
    fetch('/categories')
        .then(response => response.json())
        .then(data => {
            const categoriesList = document.getElementById('categoriesList');
            categoriesList.innerHTML = ''; // Clear existing categories
            data.forEach(category => {
                const li = document.createElement('li');
                li.textContent = category;
                categoriesList.appendChild(li);
            });
        })
        .catch(error => console.error('Error:', error));

    const categoryForm = document.getElementById('createCategoryForm');

    if (categoryForm) {
        categoryForm.addEventListener('submit', function(event) {
            event.preventDefault();

            const categoryName = document.querySelector('.createCategory').value.trim();

            if (categoryName !== '') {
                const newCategory = document.createElement('li');
                newCategory.textContent = categoryName;
                document.getElementById('categoriesList').appendChild(newCategory);
                document.querySelector('.createCategory').value = ''; // Clear the input field

                // Send form data using AJAX
                let formData = new FormData();
                formData.append('name', categoryName);
                fetch('/addCategory', { method: 'POST', body: formData })
                    .then(response => response.json())
                    .then(data => console.log(data))
                    .catch(error => console.error('Error:', error));
            }
        });
    }
});
