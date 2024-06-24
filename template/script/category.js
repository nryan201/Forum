document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('createCategoryForm').addEventListener('submit', function(event) {
        event.preventDefault();

        const categoryName = document.querySelector('.createCategory').value;

        if (categoryName.trim() !== '') {
            const newCategory = document.createElement('li');
            newCategory.textContent = categoryName;

            document.getElementById('categoriesList').appendChild(newCategory);

            document.querySelector('.createCategory').value = ''; // Clear the input field

            // Optionally, you can use AJAX to send the form data without reloading the page
             let formData = new FormData(event.target);
             fetch('/CreateCategory', { method: 'POST', body: formData })
                .then(response => response.json())
                .then(data => console.log(data));
        }
    });
});
