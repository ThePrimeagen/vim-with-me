package main

import (
	"flag"

	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/pkg/v2/relay"
)

func main() {
    godotenv.Load()

    var port uint = 0
    flag.UintVar(&port, "port", 42069, "the port to run on for the websocket")

    r := relay.NewRelay(uint16(port))
    r.Start()
}

