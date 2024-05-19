package main

import (
	"fmt"
	"log"

	_ "github.com/jhseoeo/winbal/internal/driver/screen"
	"github.com/jhseoeo/winbal/internal/signaling"
	"github.com/jhseoeo/winbal/vpx"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"
)

func main() {
	fmt.Println("WebRTC Application started using codec vpx")

	vpxp, err := vpx.NewVP8Params()
	if err != nil {
		panic(err)
	}
	vpxp.BitRate = 5000000

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&vpxp),
	)
	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))

	// devices := mediadevices.EnumerateDevices()
	// anyScreenId := devices[len(devices)-1].DeviceID

	stream, err := mediadevices.GetDisplayMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.FrameRate = prop.Float(60)
			// c.DeviceID = prop.String(anyScreenId)
		},
		Codec: codecSelector,
	})
	if err != nil {
		panic(err)
	}

	sc := signaling.NewSignalingClient()

	sc.OnSdpOffer(func(offer webrtc.SessionDescription) {
		log.Println("received offer from signaling server")
		pc, err := newPeerConnection(api, stream, sc)
		if err != nil {
			panic(err)
		}
		sc.OnIceCandidate(func(candidate webrtc.ICECandidateInit) {
			if err := pc.AddICECandidate(candidate); err != nil {
				panic(err)
			}
		})
		if err := pc.SetRemoteDescription(offer); err != nil {
			panic(err)
		}
		answer, err := pc.CreateAnswer(nil)
		if err != nil {
			panic(err)
		}
		if err := pc.SetLocalDescription(answer); err != nil {
			panic(err)
		}
		sc.SendSdpAnswer(answer)
	})

	sc.OnError(func(err error) {
		log.Println("an error occured on signaling client:", err)
	})

	err = sc.Connect("localhost", "8080") // TODO: Do not used hardcoded address
	if err != nil {
		panic(err)
	}

	select {}
}
