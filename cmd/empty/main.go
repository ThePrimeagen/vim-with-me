package main

import (
    "github.com/joho/godotenv"
	"github.com/gorilla/websocket"
)

func main() {
    _ = websocket.BinaryMessage
    godotenv.Load()
}

