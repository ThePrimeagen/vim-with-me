package testies

import (
	"flag"
	"fmt"

	"chat.theprimeagen.com/pkg/tcp"
	"chat.theprimeagen.com/pkg/window"
)

func CreateServerFromArgs() (*tcp.TCP, *window.Window, error) {
    var port uint
    flag.UintVar(&port, "port", 0, "Port to listen on")
    flag.Parse()

	if port == 0 {
        return nil, nil, fmt.Errorf("You need to provide a port")
	}

	fmt.Printf("Port: %d\n", port)
	fmt.Printf("starting server\n")

	server, err := tcp.NewTCPServer(uint16(port))
    if err != nil {
        return nil, nil, fmt.Errorf("Error creating server: %w", err)
    }

    w := window.NewWindow(80, 24)
    server.ToSockets.Welcome(window.OpenCommand(w))

    return server, w, nil
}
