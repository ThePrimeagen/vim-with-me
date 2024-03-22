package tcp

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var VERSION = 1

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
	sockets     []Connection
	welcomes    []*TCPCommand
	listener    net.Listener
}

func (t *TCP) Welcome(cmd *TCPCommand) {
	t.welcomes = append(t.welcomes, cmd)
}

// TODO: Think about project level logging and the ability to enable debug
// logging
func (t *TCP) Send(cmd *TCPCommand) {
	send(t.sockets, cmd)
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
            // TODO: Logging?
            log.Printf("server shutting down")
            break;
		}

		new_conn := NewConnection(conn)
		err = send_cmds(new_conn, t.welcomes)
		if err != nil {
			// TODO: Should i do something with the conn?
			continue
		}
		t.sockets = append(t.sockets, new_conn)

		go func(c net.Conn) {
			defer c.Close()

			fromTcp := CommandParser(conn)
			defer func() {
				fmt.Println("connection has closed")
			}()

			timer := time.NewTicker(100 * time.Millisecond)

		OuterLoop:
			for {
                select {
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
