
function uploadProfilePic(input) {
    if (input.files && input.files[0]) {
        var reader = new FileReader();
        reader.onload = function (e) {
            document.getElementById('profile-pic').style.backgroundImage = 'url(' + e.target.result + ')';
        };
        reader.readAsDataURL(input.files[0]);
    }
}
document.getElementById('profile-pic').addEventListener('mouseover', function() {
    document.getElementById('download-icon').style.display = 'block';
    this.querySelector('label').style.opacity = 1; // Make the label fully visible
});

document.getElementById('profile-pic').addEventListener('mouseout', function() {
    document.getElementById('download-icon').style.display = 'none';
    this.querySelector('label').style.opacity = 0; // Make the label transparent again
});
