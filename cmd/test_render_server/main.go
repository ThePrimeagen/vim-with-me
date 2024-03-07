package main

import (
	"flag"
	"fmt"
	"os"

	"chat.theprimeagen.com/pkg/commands"
	"chat.theprimeagen.com/pkg/tcp"
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

    count := 0
    for {
        cmd := <-server.FromSockets
        if cmd.Command == "c" {
            break
        }

        str := ""
        for i := 0; i < 1920; i++ {
            if (i + count) % 4 == 0 {
                str += "X"
            } else {
                str += " "
            }
        }

        count++
        server.Send(commands.Render(str))
    }
}

