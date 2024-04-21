package signaling

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type SignalingClient struct {
	dialer *websocket.Dialer
	conn   *websocket.Conn

	onIceCandidate func(candidate webrtc.ICECandidateInit)
	onSdpOffer     func(offer webrtc.SessionDescription)
	onClose        func()
	onError        func(err error)
}

func NewSignalingClient() *SignalingClient {
	dialer := &websocket.Dialer{}
	return &SignalingClient{
		dialer: dialer,
	}
}

func getValueFromMap[T any](m map[string]interface{}, key string) (T, error) {
	var ZERO T
	v, ok := m[key]
	if !ok {
		return ZERO, fmt.Errorf("key %s not found", key)
	}
	vv, ok := v.(T)
	if !ok {
		return ZERO, fmt.Errorf("value is not matched")
	}
	return vv, nil
}

func parseIceCandidateInit(data interface{}) (webrtc.ICECandidateInit, error) {
	res := webrtc.ICECandidateInit{}
	dmap, ok := data.(map[string]interface{})

	if !ok {
		return res, errors.New("data is not a type of map[string]interface{}")
	}

	var err error
	res.Candidate, err = getValueFromMap[string](dmap, "candidate")
	if err != nil {
		return res, err
	}
	sdpMLineIndex, err := getValueFromMap[int](dmap, "sdpMLineIndex")
	if err != nil {
		return res, err
	}
	sdpMLineIndex16 := uint16(sdpMLineIndex)
	res.SDPMLineIndex = &sdpMLineIndex16
	sdpMid, err := getValueFromMap[string](dmap, "sdpMid")
	if err != nil {
		return res, err
	}
	res.SDPMid = &sdpMid
	res.UsernameFragment, err = getValueFromMap[*string](dmap, "usernameFragment")
	if err != nil {
		return res, err
	}
	return res, nil
}

func parseSessionDescription(data interface{}) (webrtc.SessionDescription, error) {
	res := webrtc.SessionDescription{}
	dmap, ok := data.(map[string]interface{})

	if !ok {
		return res, errors.New("data is not a type of map[string]interface{}")
	}

	var err error
	res.SDP, err = getValueFromMap[string](dmap, "sdp")
	if err != nil {
		return res, err
	}
	typestr, err := getValueFromMap[string](dmap, "type")
	if err != nil {
		return res, err
	}
	res.Type = webrtc.NewSDPType(typestr)

	return res, nil
}

func (s *SignalingClient) Connect(host, port string) error {
	if s.conn != nil {
		return nil
	}
	if s.onIceCandidate == nil {
		return errors.New("onIceCandidate is not set")
	}
	if s.onSdpOffer == nil {
		return errors.New("onSdpOffer is not set")
	}

	conn, _, err := s.dialer.Dial("ws://"+host+":"+port+"/ws", nil)
	if err != nil {
		return err
	}
	s.conn = conn

	s.conn.WriteJSON(message{
		Type: "join",
		Data: "master",
	})

	// read loop
	go func() {
		defer s.conn.Close()

		for {
			mt, payload, err := s.conn.ReadMessage()
			if mt == websocket.CloseMessage || mt == websocket.CloseGoingAway {
				s.onClose()
				return
			} else if err != nil {
				s.onError(err)
				return
			} else if mt != websocket.TextMessage {
				log.Println("Received non-text message")
				continue
			}

			var msg message
			err = json.Unmarshal(payload, &msg)
			if err != nil {
				s.onError(err)
				return
			}

			switch msg.Type {
			case "candidate":
				ic, err := parseIceCandidateInit(msg.Data)
				if err != nil {
					s.onError(err)
				}
				s.onIceCandidate(ic)
			case "offer":
				sdp, err := parseSessionDescription(msg.Data)
				if err != nil {
					s.onError(err)
				}
				s.onSdpOffer(sdp)
			default:
				log.Println("Unknown message type: ", msg.Type)
			}
		}
	}()

	return nil
}

func (s *SignalingClient) Close() error {
	return s.conn.Close()
}

func (s *SignalingClient) SendIceCandidate(candidate webrtc.ICECandidateInit) error {
	return s.conn.WriteJSON(message{
		Type: "candidate",
		Data: candidate,
	})
}

func (s *SignalingClient) SendSdpAnswer(answer webrtc.SessionDescription) error {
	return s.conn.WriteJSON(message{
		Type: "answer",
		Data: answer,
	})
}

func (s *SignalingClient) OnIceCandidate(f func(webrtc.ICECandidateInit)) {
	s.onIceCandidate = f
}

func (s *SignalingClient) OnSdpOffer(f func(webrtc.SessionDescription)) {
	s.onSdpOffer = f
}

func (s *SignalingClient) OnError(f func(error)) {
	s.onError = f
}
