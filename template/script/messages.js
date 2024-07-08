async function fetchUsers() {
    const response = await fetch('/get-users');
    const users = await response.json();
    const dataList = document.getElementById('usernames');
    dataList.innerHTML = ''; // Clear previous options
    users.forEach(user => {
        const option = document.createElement('option');
        option.value = user.Username;
        dataList.appendChild(option);
    });
}
fetchUsers();