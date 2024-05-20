package tcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"syscall"
)

var VERSION byte = 1
var HEADER_SIZE = 4

type TCPCommand struct {
	Command byte
	Data    []byte
}

type WelcomeCB func() *TCPCommand

func (t *TCPCommand) MarshalBinary() (data []byte, err error) {
	length := uint16(len(t.Data))
	lengthData := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthData, length)

	b := make([]byte, 0, 1+1+2+length)
	b = append(b, VERSION)
	b = append(b, t.Command)
	b = append(b, lengthData...)
	return append(b, t.Data...), nil
}

func (t *TCPCommand) UnmarshalBinary(bytes []byte) error {

	if bytes[0] != VERSION {
		return fmt.Errorf("version mismatch %d != %d", bytes[0], VERSION)
	}

	length := int(binary.BigEndian.Uint16(bytes[2:]))
	end := HEADER_SIZE + length

	if len(bytes) < end {
		return fmt.Errorf("not enough data to parse packet: got %d expected %d", len(bytes), HEADER_SIZE+length)
	}

	command := bytes[1]
	data := bytes[HEADER_SIZE:end]

	t.Command = command
	t.Data = data

	return nil
}

type TCP struct {
	welcomes    []WelcomeCB
	sockets     []Connection
	listener    net.Listener
	mutex       sync.RWMutex
	FromSockets chan TCPCommandWrapper
	NewSocket   chan *Connection
}

func (t *TCP) ConnectionCount() int {
	return len(t.sockets)
}

func (t *TCP) Send(command *TCPCommand) {
	t.mutex.RLock()
	removals := make([]int, 0)
	slog.Debug("sending message", "msg", command)
	for i, conn := range t.sockets {
		err := conn.Writer.Write(command)
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				slog.Debug("connection closed by client", "index", i)
			} else {
				slog.Error("removing due to error", "index", i, "error", err)
			}
			removals = append(removals, i)
		}
	}
	t.mutex.RUnlock()

	if len(removals) > 0 {
		t.mutex.Lock()
		for i := len(removals) - 1; i >= 0; i-- {
			idx := removals[i]
			t.sockets = append(t.sockets[:idx], t.sockets[idx+1:]...)
		}
		t.mutex.Unlock()
	}
}

func (t *TCP) welcome(conn *Connection) error {
	for _, w := range t.welcomes {
		cmd := w()
		err := conn.Writer.Write(cmd)

		if err != nil {
			// TODO: Do i need to close the connection?
			return err
		}
	}

	return nil
}

func (t *TCP) WelcomeMessage(cmd WelcomeCB) {
	t.welcomes = append(t.welcomes, cmd)
}

func (t *TCP) Close() {
	t.listener.Close()
}

type TCPCommandWrapper struct {
	Conn    *Connection
	Command *TCPCommand
}

func NewTCPServer(port uint16) (*TCP, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	// TODO: Done channel
	return &TCP{
		sockets:     make([]Connection, 0, 10),
		welcomes:    make([]WelcomeCB, 0, 10),
		listener:    listener,
		FromSockets: make(chan TCPCommandWrapper, 10),
		mutex:       sync.RWMutex{},
	}, nil
}

func readConnection(tcp *TCP, conn *Connection) {
	for {
		cmd, err := conn.Next()
		slog.Debug("new command", "id", conn.Id, "cmd", cmd)

		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Debug("socket received EOF", "id", conn.Id, "error", err)
			} else {
				slog.Error("received error while reading from socket", "id", conn.Id, "error", err)
			}
			break
		}

		tcp.FromSockets <- TCPCommandWrapper{Command: cmd, Conn: conn}
	}
}

func (t *TCP) Start() {
	id := 0
	for {
		conn, err := t.listener.Accept()
		id++

		if err != nil {
			slog.Error("server error:", "error", err)
		}

		newConn := NewConnection(conn, id)
		slog.Debug("new connection", "id", newConn.Id)
		err = t.welcome(&newConn)

		if err != nil {
			slog.Error("could not send out welcome messages", "error", err)
			continue
		}

		t.mutex.Lock()
		t.sockets = append(t.sockets, newConn)
		t.mutex.Unlock()

		go readConnection(t, &newConn)
	}
}

func MakeWelcome(cmd *TCPCommand) WelcomeCB {
	return func() *TCPCommand {
		return cmd
	}
}
