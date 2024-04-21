package main

import "github.com/gorilla/websocket"

type Session struct {
	id     string
	master *Peer
	viewer *Peer
	done   chan struct{}
}

type Peer struct {
	conn *websocket.Conn
	ch   chan []byte
}

var mutex = make(chan struct{}, 1)
var sessions map[string]*Session

func init() {
	sessions = make(map[string]*Session)
}

func NewSession(id string) *Session {
	mutex <- struct{}{}
	defer func() { <-mutex }()

	if s, ok := sessions[id]; ok {
		return s
	}
	s := &Session{
		id: id,
		master: &Peer{
			ch: make(chan []byte),
		},
		viewer: &Peer{
			ch: make(chan []byte),
		},
		done: make(chan struct{}),
	}
	sessions[id] = s

	go func() {
		for {
			select {
			case msg := <-s.master.ch:
				if s.master.conn == nil {
					continue
				}
				err := s.master.conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					s.Close()
					return
				}
			case msg := <-s.viewer.ch:
				if s.viewer.conn == nil {
					continue
				}
				err := s.viewer.conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					s.Close()
					return
				}
			case <-s.done:
				return
			}
		}
	}()

	return s
}

func (s *Session) SetMaster(c *websocket.Conn) {
	s.master.conn = c
}

func (s *Session) SetViewer(c *websocket.Conn) {
	s.viewer.conn = c
}

func (s *Session) SendToAnother(c *websocket.Conn, msg []byte) {
	if c == s.master.conn {
		s.viewer.ch <- msg
	} else {
		s.master.ch <- msg
	}
}

func (s *Session) Leave(c *websocket.Conn) {
	if c == s.master.conn {
		s.master.conn.Close()
		s.master.conn = nil
	} else if c == s.viewer.conn {
		s.viewer.conn.Close()
		s.viewer.conn = nil
	}

	if s.master.conn == nil && s.viewer.conn == nil {
		s.Close()
	}
}

func (s *Session) Close() {
	mutex <- struct{}{}
	defer func() { <-mutex }()

	s.master.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	s.viewer.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	s.master.conn.Close()
	s.viewer.conn.Close()

	delete(sessions, s.id)
	close(s.done)
}
