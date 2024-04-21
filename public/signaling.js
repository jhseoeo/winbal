class Signaling {
    constructor(role = 'viewer') {
        this.ws = new WebSocket('ws://localhost:8080/ws');
        this.role = role
        this.onSdpOffer_ = () => {}
        this.onSdpAnswer_ = () => {}
        this.onIceCandidate_ = () => {}
        this.ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            console.log(event.data)
            if (message.type === 'offer') {
                this.onSdpOffer_(message.data);
            }
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
                    data: this.role
                }));
                resolve();
            } else {
                this.ws.onopen = () => {
                    resolve();
                    console.log('Connected to signaling server');
                    this.ws.send(JSON.stringify({
                        type: "join",
                        data: this.role
                    }));
                };
            }
        });
    }

    onSdpOffer(f) {
        this.onSdpOffer_ = f
    }

    onSdpAnswer(f) {
        this.onSdpAnswer_ = f
    }

    onIceCandidate(f) {
        this.onIceCandidate_ = f
    }

    sendSdpOffer(sdp) {
        this.ws.send(JSON.stringify({
            type: 'offer',
            data: sdp
        }));
    }

    sendSdpAnswer(sdp) {
        this.ws.send(JSON.stringify({
            type: 'answer',
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