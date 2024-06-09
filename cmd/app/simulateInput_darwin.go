package main

import "fmt"

func keyDown(keyCode int) {
	fmt.Println("keyDown: ", keyCode)
}

func keyUp(keyCode int) {
	fmt.Println("keyUp: ", keyCode)
}

func mouseMove(x, y int) {
	fmt.Println("mouseMove: ", x, y)
}

func mouseDown(btn MouseBtnType) {
	fmt.Println("mouseDown", btn)
}

func mouseUp(btn MouseBtnType) {
	fmt.Println("mouseUp", btn)
}
