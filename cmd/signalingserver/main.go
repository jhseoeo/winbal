package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func staticFiles() func(http.ResponseWriter, *http.Request) {
	h := http.FileServer(http.Dir("./public"))
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

func handleWs() func(http.ResponseWriter, *http.Request) {
	var upgrader = websocket.Upgrader{}

	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer c.Close()

		var msg Message
		err = c.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}
		if msg.Type != "join" {
			log.Println("Invalid message type")
		}
		connectionType := msg.Data

		const sessionId = "session"
		s := NewSession(sessionId)
		if connectionType == "master" {
			s.SetMaster(c)
		} else if connectionType == "viewer" {
			s.SetViewer(c)
		} else {
			log.Println("Invalid connection type")
		}

		for {
			code, payload, err := c.ReadMessage()
			if code == websocket.CloseGoingAway {
				s.Leave(c)
			} else if err != nil {
				log.Println(err)
				break
			}
			s.SendToAnother(c, payload)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWs())
	http.HandleFunc("/", staticFiles())
	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
