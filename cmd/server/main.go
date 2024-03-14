package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/tcp"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

func main() {
    var port uint
    flag.UintVar(&port, "port", 0, "Port to listen on")
    flag.Parse()

	if port == 0 {
		fmt.Printf("You need to provide a port")
		os.Exit(1)
	}

	fmt.Printf("Port: %d\n", port)
	fmt.Printf("starting server\n")

	server, err := tcp.NewTCPServer(uint16(port))
    if err != nil {
        fmt.Printf("Error starting server: %s", err)
        os.Exit(1)
    }

    w := window.NewWindow(80, 24)
    cmd := window.OpenCommand(w)
    server.ToSockets.Welcome(cmd)

    for {
        time.Sleep(1 * time.Second)
    }
}
