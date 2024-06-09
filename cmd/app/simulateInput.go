package main

import "fmt"

type KeyboardInput struct {
	KeyCode int    `json:"key"`
	Type    string `json:"type"`
}

type MouseInput struct {
	X      int          `json:"x"`
	Y      int          `json:"y"`
	Button MouseBtnType `json:"btn"`
	Type   string       `json:"type"`
}

type MouseBtnType int

const (
	MoustBtnLeft MouseBtnType = iota
	MoustBtnMiddle
	MoustBtnRight
)

func handleKeyboardInput(ki KeyboardInput) {
	fmt.Println("handleKeyboardInput", ki)
	switch ki.Type {
	case "down":
		keyDown(ki.KeyCode)
	case "up":
		keyUp(ki.KeyCode)
	}
}

func handleMouseInput(mi MouseInput) {
	fmt.Println("handleMouseInput", mi)
	switch mi.Type {
	case "move":
		mouseMove(mi.X, mi.Y)
	case "down":
		mouseDown(mi.Button)
	case "up":
		mouseUp(mi.Button)
	}
}
