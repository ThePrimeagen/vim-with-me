package testies

import (
	"flag"
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/tcp2"
)

func CreateServerFromArgs() (*tcp2.TCP, error) {
    var port uint
    flag.UintVar(&port, "port", 0, "Port to listen on")
    flag.Parse()

	if port == 0 {
        return nil, fmt.Errorf("You need to provide a port")
	}

	fmt.Printf("Port: %d\n", port)
	fmt.Printf("starting server\n")

	server, err := tcp2.NewTCPServer(uint16(port))
    if err != nil {
        return nil, fmt.Errorf("Error creating server: %w", err)
    }

    return server, nil
}
