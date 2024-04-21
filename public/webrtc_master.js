const video = document.getElementById('video')

class WebRTC {
    constructor() {
        this.signaling = new Signaling("master");
        this.pc = new RTCPeerConnection({
            iceServers: [
                {
                    urls: "stun:stun.l.google.com:19302"
                }
            ]
        });
        this.pc.onicecandidate = (event) => {
            if (event.candidate) {
                console.log(event.candidate.toJSON());
                this.signaling.sendIceCandidate(event.candidate);
            }
        }
        this.pc.oniceconnectionstatechange = (event) => {
            console.log(this.pc.iceConnectionState)
        }
        this.signaling.onSdpOffer(async (answer) => {
            await this.pc.setRemoteDescription(answer);
            const sdp = await this.pc.createAnswer();
            await this.pc.setLocalDescription(sdp);
            this.signaling.sendSdpAnswer(sdp);
        });
        this.signaling.onIceCandidate(async (candidate) => {
            await this.pc.addIceCandidate(candidate);
        });

        navigator.mediaDevices.getDisplayMedia({ video: true, audio: false }).then((stream) => {
            stream.getTracks().forEach((track) => {
                this.pc.addTrack(track, stream);
            });
            video.srcObject = stream;
        });
    }

    async start() {
        await this.signaling.waitForConnection();
    }
}

const webRtc = new WebRTC();
webRtc.start();