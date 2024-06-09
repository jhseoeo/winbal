import { Signaling } from './signaling.js';

const video = document.getElementById('video')

class WebRTC {
    constructor() {
        this.signaling = new Signaling();
        this.pc = new RTCPeerConnection({
            iceServers: [
                {
                    urls: "stun:stun.l.google.com:19302"
                }
            ]
        });
        this.pc.onicecandidate = async (event) => {
            if (event.candidate) {
                console.log(event.candidate);
                this.signaling.sendIceCandidate(event.candidate);
            }
        }
        this.pc.oniceconnectionstatechange = (event) => {
            console.log(this.pc.iceConnectionState)
        }
        this.signaling.onSdpAnswer(async (answer) => {
            await this.pc.setRemoteDescription(answer);
        });
        this.signaling.onIceCandidate(async (candidate) => {
            await this.pc.addIceCandidate(candidate);
        });
        this.pc.ontrack = (event) => {
            console.log("ontrack");
            video.srcObject = new MediaStream([event.track]);
        }
        this.keyboardChannel = this.pc.createDataChannel("keyboard");
        this.mouseChannel = this.pc.createDataChannel("mouse");
        this.keyboardChannel.onopen = () => {
            console.log("keyboard datachannel open");
        }
        this.mouseChannel.onopen = () => {
            console.log("mouse datachannel open");
        }

        setTimeout(() => {
            console.log(
                this.pc.iceGatheringState,
                this.pc.getReceivers()[0].track.readyState,
                this.pc.getReceivers()[0].getParameters(),
                this.pc.getTransceivers()[0].receiver.transport.state
            );
        }, 4000);
    }

    async createOffer() {
        await this.signaling.waitForConnection();
        const offer = await this.pc.createOffer({
            offerToReceiveVideo: true,
            offerToReceiveAudio: false,
        });
        await this.pc.setLocalDescription(offer); 
        this.signaling.sendSdpOffer(offer);
    }

    async sendKeyboardEvent(key, type) {
        if (this.keyboardChannel.readyState !== "open") return;
        this.keyboardChannel.send(JSON.stringify({ key, type }));
    }

    async sendMouseEvent(x, y, btn, type) {
        if (this.mouseChannel.readyState !== "open") return;
        this.mouseChannel.send(JSON.stringify({ x, y, btn, type }));
    }
}

const webRtc = new WebRTC();
webRtc.createOffer();

export { webRtc };