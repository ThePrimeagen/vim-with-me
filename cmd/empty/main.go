package main

import (
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	_ = websocket.BinaryMessage
	godotenv.Load()
}
