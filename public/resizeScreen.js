const video = document.getElementById('video')
const videoRatio = 1920 / 1080;

function resizeScreen() {
    const windowWidth = window.innerWidth;
    const windowHeight = window.innerHeight;
    let width = windowWidth;
    let height = windowHeight;
    if (windowWidth / windowHeight > videoRatio) {
        width = windowHeight * videoRatio;
    } else {
        height = windowWidth / videoRatio;
    }
    video.width = width;
    video.height = height;
}
resizeScreen()

window.addEventListener('resize', resizeScreen)