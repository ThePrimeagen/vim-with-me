package tcp

import (
	"fmt"
	"net"
)

func (t *TCP) Close() {
	t.listener.Close()
}

func NewTCPServer(port uint16) (*TCP, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Error creating TCP server: %w", err)
	}

	tcps := &TCP{
		FromSockets: make(chan *TCPCommand, 10),
		listener:    listener,
		sockets:     make([]*Connection, 0, 10),
		welcomes:    make([]*TCPCommand, 0, 10),
	}

	go func() { tcps.listen() }()

	return tcps, nil

}
