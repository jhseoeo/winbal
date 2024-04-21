package main

// /*
// #include <stdio.h>

// void hello() {
//     printf("Hello, World!\n");
// }
// */
// import "C"
import (
	"fmt"
	"log"

	"github.com/jhseoeo/winbal/signaling"
	"github.com/jhseoeo/winbal/x264"
	"github.com/pion/mediadevices"
	_ "github.com/pion/mediadevices/pkg/driver/screen"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"
)

func main() {
	x264Params, err := x264.NewParams()
	if err != nil {
		panic(err)
	}
	x264Params.BitRate = 512000
	x264Params.KeyFrameInterval = 60
	x264Params.Preset = x264.PresetUltrafast

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&x264Params),
	)
	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	pc, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	stream, err := mediadevices.GetDisplayMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
			c.FrameRate = prop.Float(30)
			c.DeviceID = prop.String("screen:0:0")
		},
		Codec: codecSelector,
	})
	if err != nil {
		panic(err)
	}

	for _, track := range stream.GetTracks() {
		track.OnEnded(func(err error) {
			fmt.Printf("Track (ID: %s) ended with error: %v\n", track.ID(), err)
		})

		if _, err := pc.AddTransceiverFromTrack(track, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionSendonly,
		}); err != nil {
			panic(err)
		}
	}

	sc := signaling.NewSignalingClient()

	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		sc.SendIceCandidate(c.ToJSON())
	})

	sc.OnSdpOffer(func(offer webrtc.SessionDescription) {
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

	sc.OnIceCandidate(func(candidate webrtc.ICECandidateInit) {
		if err := pc.AddICECandidate(candidate); err != nil {
			panic(err)
		}
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
