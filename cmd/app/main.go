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

	vpxp, err := vpx.NewVP9Params()
	if err != nil {
		panic(err)
	}
	vpxp.BitRate = 512000
	vpxp.KeyFrameInterval = 60
	vpxp.RTPCodec().ClockRate = 90000

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&vpxp),
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

	devices := mediadevices.EnumerateDevices()
	anyScreenId := devices[len(devices)-1].DeviceID

	stream, err := mediadevices.GetDisplayMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			// c.FrameFormat = prop.FrameFormat("RGBA")
			// c.Width = prop.Int(1600)
			// c.Height = prop.Int(900)
			c.FrameRate = prop.Float(30)
			c.DeviceID = prop.String(anyScreenId)
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

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		dc.OnOpen(func() {
			err := dc.SendText("HI~")
			if err != nil {
				log.Println("failed to send data channel message:", err)
			}
		})

		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.Println("received data channel message:", string(msg.Data))
			if string(msg.Data) == "ping" {
				err := dc.SendText("pong")
				if err != nil {
					log.Println("failed to send data channel message:", err)
				}
			}
		})
	})

	sc.OnSdpOffer(func(offer webrtc.SessionDescription) {
		log.Println("received offer from signaling server")
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
