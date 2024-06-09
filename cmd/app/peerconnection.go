package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jhseoeo/winbal/internal/signaling"
	"github.com/pion/mediadevices"
	"github.com/pion/webrtc/v3"
)

func newPeerConnection(api *webrtc.API, stream mediadevices.MediaStream, sc *signaling.SignalingClient) (*webrtc.PeerConnection, error) {
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
			log.Println("data channel opened:", dc.Label())
		})
		if dc.Label() == "keyboard" {
			dc.OnMessage(func(msg webrtc.DataChannelMessage) {
				log.Println("received keyboard data channel message:", string(msg.Data))
				var ki KeyboardInput
				err := json.Unmarshal(msg.Data, &ki)
				if err != nil {
					log.Println("failed to unmarshal keyboard input:", err)
					return
				}
				handleKeyboardInput(ki)
			})
		} else if dc.Label() == "mouse" {
			dc.OnMessage(func(msg webrtc.DataChannelMessage) {
				log.Println("received mouse data channel message:", string(msg.Data))
				var mi MouseInput
				err := json.Unmarshal(msg.Data, &mi)
				if err != nil {
					log.Println("failed to unmarshal mouse input:", err)
					return
				}
				handleMouseInput(mi)
			})
		} else {
			log.Println("received unknown data channel:", dc.Label())
		}
	})

	pc.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		log.Println("peer connection state changed:", pcs)
		if pcs == webrtc.PeerConnectionStateFailed {
			pc.Close()
		}
	})

	return pc, nil
}
