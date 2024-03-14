package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/tcp"
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
		fmt.Printf("Could not start the server: %v", err)
		os.Exit(1)
	}

    defer server.Close()

	fmt.Printf("server started and waiting for command\n")
	cmd := <-server.FromSockets
	fmt.Printf("Command: %v\n", cmd)

    server.Send(&tcp.TCPCommand{
        Command: cmd.Data,
        Data: cmd.Command,
    })

    time.Sleep(1 * time.Second)
}

