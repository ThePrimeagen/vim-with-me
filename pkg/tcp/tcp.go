package tcp

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var VERSION = 1

type TCPStream struct {
	outs []chan TCPCommand
	lock sync.RWMutex
    welcomes []*TCPCommand
}

func (t *TCPStream) Welcome(cmd *TCPCommand) {
    t.welcomes = append(t.welcomes, cmd)
}

func (t *TCPStream) Spread(command *TCPCommand) {
	t.lock.RLock()
	for i, listener := range t.outs {
        fmt.Printf("sending command to listener: %d\n", i)
		listener <- *command
        fmt.Printf("   sent command: %d\n", i)
	}
	t.lock.RUnlock()
}

func (t *TCPStream) listen() <-chan TCPCommand {
    fmt.Println("adding listener")
	t.lock.Lock()
	listener := make(chan TCPCommand, 10)
	t.outs = append(t.outs, listener)
	t.lock.Unlock()
	return listener
}

func (t *TCPStream) removeListen(rm <-chan TCPCommand) {
    fmt.Println("removing listener")
	t.lock.Lock()
	for i, listener := range t.outs {
		if listener == rm {
			t.outs = append(t.outs[:i], t.outs[i+1:]...)
			break
		}
	}
	t.lock.Unlock()
}

func createTCPCommandSpread() TCPStream {
	return TCPStream{
		outs: make([]chan TCPCommand, 0),
		lock: sync.RWMutex{},
	}
}

type TCPCommand struct {
	Command string
	Data    string
}

var malformedTCPCommand = TCPCommand{
	Command: "e",
	Data:    "Malformed TCP Command",
}

func versionMismatch(v1, v2 int) *TCPCommand {
	return &TCPCommand{
		Command: "e",
		Data:    fmt.Sprintf("Version Mismatch %d %d", v1, v2),
	}
}

var tcpClosedCommand = TCPCommand{
	Command: "c",
	Data:    "Connection Closed",
}

func (t *TCPCommand) Bytes() []byte {
	str := fmt.Sprintf("%s:%s", t.Command, t.Data)
	str = fmt.Sprintf("%d:%d:%s", VERSION, len(str), str)
	return []byte(str)
}

func CommandFromBytes(b string) (string, *TCPCommand) {
	parts := strings.SplitN(b, ":", 3)
	if len(parts) != 3 {
		return b, nil
	}

	versionStr := parts[0]
	lengthStr := parts[1]
	dataStr := parts[2]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return b, &malformedTCPCommand
	}

	if version != VERSION {
		return b, versionMismatch(version, VERSION)
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return b, &malformedTCPCommand
	}
	if len(dataStr) < length {
		return b, nil
	}

	remaining := dataStr[length:]
	commandStr := dataStr[:length]
	commandParts := strings.SplitN(commandStr, ":", 2)

	if len(commandParts) != 2 {
		return b, &malformedTCPCommand
	}

	cmd := &TCPCommand{
		Command: commandParts[0],
		Data:    commandParts[1],
	}

	return remaining, cmd
}

type TCP struct {
	FromSockets chan TCPCommand
	ToSockets   TCPStream
	listener    net.Listener
}

func (t *TCP) Send(cmd *TCPCommand) {
	t.ToSockets.Spread(cmd)
}

func CommandParser(r io.Reader) chan TCPCommand {
	out := make(chan TCPCommand)

	go func() {
		defer close(out)

		buffer := make([]byte, 1024)
		previous := ""
		for {
			n, err := r.Read(buffer)
			if err != nil {
                fmt.Printf("Error reading from connection: %s\n", err)
				out <- tcpClosedCommand
				return
			}
			current := previous + string(buffer[:n])

			for remaining, cmd := CommandFromBytes(current); cmd != nil; remaining, cmd = CommandFromBytes(current) {
				current = remaining
				out <- *cmd
			}

			previous = current
		}
	}()

	return out
}

func (t *TCP) listen() {
	for {

		conn, err := t.listener.Accept()
		if err != nil {
			// true and factual
			log.Fatal("You like amouranth", err)
		}

        // TODO: Think about this a bit more... i worry
        for _, cmd := range t.ToSockets.welcomes {
            _, err := conn.Write(cmd.Bytes())
            if err != nil {
                conn.Close()
                continue
            }
        }

		go func(c net.Conn) {
            defer c.Close()

            toTcp := t.ToSockets.listen()
            defer t.ToSockets.removeListen(toTcp)

            fromTcp := CommandParser(conn)
            defer func() {
                fmt.Println("connection has closed")
            }()

            timer := time.NewTicker(1 * time.Second)

		OuterLoop:
			for {
				select {
				case cmd := <-toTcp:
					_, err := c.Write(cmd.Bytes())
					if err != nil {
						fmt.Printf("Error writing to client: %s\n", err)
						break OuterLoop
					}

				case cmd := <-fromTcp:
					// NOTE: i am sure there is a better way to do this
					// TODO: Figure out that better way
					if cmd.Command == "c" {
                        fmt.Println("closing connection")
						break OuterLoop
					}

					t.FromSockets <- cmd

					if cmd.Command == "e" {
						break OuterLoop
					}

                case <-timer.C:
                    fmt.Println("tick")
				}
			}
            fmt.Println("I AM DONE WITH THIS SHIT")

		}(conn)
	}
}

func (t *TCP) Close() {
	t.listener.Close()
}

func NewTCPServer(port uint16) (*TCP, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Error creating TCP server: %w", err)
	}

	tcps := &TCP{
		FromSockets: make(chan TCPCommand, 10),
		ToSockets:   createTCPCommandSpread(),
		listener:    listener,
	}

	go func() { tcps.listen() }()

	return tcps, nil

}
