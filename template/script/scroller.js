document.addEventListener('DOMContentLoaded', function() {

    const mainSection = document.querySelector('.main-section');

    mainSection.addEventListener('scroll', function() {
        console.log('scrolling');
    });
});