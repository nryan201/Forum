document.addEventListener('DOMContentLoaded', function() {
const form = document.getElementById('createHashtagForm');
const konamiCode = ['ArrowUp', 'ArrowUp', 'ArrowDown', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'ArrowLeft', 'ArrowRight', 'b', 'a', 'Enter'];
let currentInput = [];

document.addEventListener('keydown', (event) => {
    currentInput.push(event.key);
    if (currentInput.length > konamiCode.length) {
        currentInput.shift();
    }

    if (currentInput.join('') === konamiCode.join('')) {
        launchConfetti();
        alert('Konami Code activated!');
        window.location.href = 'https://www.youtube.com/watch?v=dQw4w9WgXcQ';
    }
});

function launchConfetti() {
    if (typeof confetti === 'function') {
        confetti({
            particleCount: 100,
            spread: 70,
            origin: { y: 0.6 }
        });
    } else {
        console.log('Confetti function is not defined.');
    }
}
});