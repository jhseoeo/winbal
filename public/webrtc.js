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
        this.pc.onicecandidate = (event) => {
            if (event.candidate) {
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
            video.srcObject = event.streams[0];
        }
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
}

const webRtc = new WebRTC();
webRtc.createOffer();
