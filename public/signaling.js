class Signaling {
    constructor() {
        this.ws = new WebSocket('ws://localhost:8080/ws');
        this.onSdpAnswer_ = () => {}
        this.onIceCandidate_ = () => {}
        this.ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            console.log(event.data)
            if (message.type === 'answer') {
                this.onSdpAnswer_(message.data);
            }
            if (message.type === 'candidate') {
                this.onIceCandidate_(message.data);
            }
        }
    }

    waitForConnection() {
        return new Promise((resolve, reject) => {
            if (this.ws.readyState === WebSocket.OPEN) {
                console.log('Connected to signaling server');
                this.ws.send(JSON.stringify({
                    type: "join",
                    data: "viewer"
                }));
                resolve();
            } else {
                this.ws.onopen = () => {
                    resolve();
                    console.log('Connected to signaling server');
                    this.ws.send(JSON.stringify({
                        type: "join",
                        data: "viewer"
                    }));
                };
            }
        });
    }

    onSdpAnswer(callback) {
        this.onSdpAnswer_ = callback
    }

    onIceCandidate(callback) {
        this.onIceCandidate_ = callback
    }

    sendSdpOffer(sdp) {
        this.ws.send(JSON.stringify({
            type: 'offer',
            data: sdp
        }));
    }

    sendIceCandidate(candidate) {
        this.ws.send(JSON.stringify({
            type: 'candidate',
            data: candidate
        }));
    }
}