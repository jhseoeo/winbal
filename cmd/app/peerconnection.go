package main

import (
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

	pc.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		log.Println("peer connection state changed:", pcs)
		if pcs == webrtc.PeerConnectionStateFailed {
			pc.Close()
		}
	})

	return pc, nil
}
