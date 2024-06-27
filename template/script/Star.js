document.addEventListener("DOMContentLoaded", () => {
    const stars = document.querySelectorAll('.star, .star1, .star2, .star3, .star4');
    const container = document.querySelector('.profile-section');

    stars.forEach(star => {
        let posX = Math.random() * container.clientWidth - 1;
        let posY = Math.random() * container.clientHeight - 1 ;
        let velocityX = 2 + Math.random() * 2; // Randomize velocity
        let velocityY = 2 + Math.random() * 2;
        const starSize = 20; // Taille de l'étoile (doit correspondre à la taille définie dans le CSS)

        function moveStar() {
            posX += velocityX;
            posY += velocityY;

            // Collision avec les bords de la div
            if (posX <= 0 || posX + starSize >= container.clientWidth) {
                velocityX *= -1;
            }
            if (posY <= 0 || posY + starSize >= container.clientHeight) {
                velocityY *= -1;
            }

            star.style.left = posX + 'px';
            star.style.top = posY + 'px';

            requestAnimationFrame(moveStar);
        }

        moveStar();
    });
});
