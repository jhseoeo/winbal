import { webRtc } from './webrtc.js';

const video = document.getElementById('video');
const screenWidth = 1920;
const screenHeight = 1080;

class InputObserver {
    constructor() {
        this.mousePos = { x: 0, y: 0 };
        this.prevMousePos = { ...this.mousePos };
        window.addEventListener('keydown', this.keydown.bind(this));
        window.addEventListener('keyup', this.keyup.bind(this));
        video.addEventListener("mousedown", this.mousedown.bind(this));
        video.addEventListener("mouseup", this.mouseup.bind(this));
        video.addEventListener("mousemove", this.mousemove.bind(this));
        video.addEventListener("contextmenu", this.noEvent.bind(this));

        this.mousePosloop = setInterval(this.mouseLoopFunc.bind(this), 1000/30);
    }

    keydown(e) {
        console.log(e.keyCode);
        webRtc.sendKeyboardEvent(e.keyCode, "down");
    }

    keyup(e) {
        e.preventDefault();
        webRtc.sendKeyboardEvent(e.keyCode, "up");
    }

    mousedown(e) {
        e.preventDefault();
        webRtc.sendMouseEvent(this.mousePos.x, this.mousePos.y, e.button, "down");
    }

    mouseup(e) {
        e.preventDefault();
        webRtc.sendMouseEvent(this.mousePos.x, this.mousePos.y, e.button, "up");
    }

    mousemove(e) {
        e.preventDefault();
        const videoWidth = video.width;
        const videoHeight = video.height;
        const x = Math.round(e.offsetX * screenWidth / videoWidth);
        const y = Math.round(e.offsetY * screenHeight / videoHeight);
        this.mousePos = { x, y };
    }

    mouseLoopFunc() {
        if (this.mousePos.x === this.prevMousePos.x && this.mousePos.y === this.prevMousePos.y) return;
        if (this.mousePos.x < 0 || this.mousePos.y < 0 || this.mousePos.x > screenWidth || this.mousePos.y > screenHeight) return;
        webRtc.sendMouseEvent(this.mousePos.x, this.mousePos.y, 0, "move");
        this.prevMousePos = { ...this.mousePos };
    }

    noEvent(e) {
        e.preventDefault();
        return false;
    }
}

const keybinder = new InputObserver();