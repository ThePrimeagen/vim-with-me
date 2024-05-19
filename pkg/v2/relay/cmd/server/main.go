package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/relay"
)

func main() {
    godotenv.Load()

    var port uint = 0
    flag.UintVar(&port, "port", 0, "the port to run on for the websocket")

    if port == 42069 {
        portStr := os.Getenv("PORT")
        portEnv, err := strconv.Atoi(portStr)
        if err == nil {
            port = uint(portEnv)
        }
    }

    assert.Assert(port != 0, "please provide a port for the relay server")

	uuid := os.Getenv("AUTH_ID")
	assert.Assert(len(uuid) > 0, "empty auth id, unable to to start relay")

    r := relay.NewRelay(uint16(port), uuid)
    r.Start()
}

