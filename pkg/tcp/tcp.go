package tcp

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var VERSION byte = 1
var MINIMUM = 4

type TCPCommand struct {
	Command byte
	Data    []byte
}

func (t *TCPCommand) Bytes() []byte {
    length := uint16(len(t.Data))
    lengthData := make([]byte, 2)
    binary.BigEndian.PutUint16(lengthData, length)

    b := make([]byte, 0, 1 + 1 + 2 + length)
    b = append(b, VERSION)
    b = append(b, t.Command)
    b = append(b, lengthData...)
    return append(b, t.Data...)
}

func CommandFromBytes(b []byte) ([]byte, *TCPCommand, error) {
    if len(b) < MINIMUM {
        return b, nil, nil
    }

    if b[0] != VERSION {
        return b, nil, fmt.Errorf("version mismatch %d != %d", b[0], VERSION)
    }

    length := int(binary.BigEndian.Uint16(b[2:]))
    end := MINIMUM + length

    if len(b) < end {
        return b, nil, nil
    }

    command := b[1]
    data := b[MINIMUM:end]

	cmd := &TCPCommand{
		Command: command,
		Data:    data,
	}

    return b[end:], cmd, nil
}

type TCP struct {
	FromSockets chan *TCPCommand
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

type TCPCommandResult struct {
    Error error
    Command *TCPCommand
}

func CommandParser(r io.Reader) chan TCPCommandResult {
	out := make(chan TCPCommandResult)

	go func() {
		defer close(out)

		buffer := make([]byte, 1024)
		previous := make([]byte, 0)
		for {
			n, err := r.Read(buffer)
			if err != nil {
				out <- TCPCommandResult{
                    Error: fmt.Errorf("error calling r.Read: %w", err),
                    Command: nil,
                }
				return
			}

			current := append(previous, buffer[:n]...)

			for remaining, cmd, err := CommandFromBytes(current); cmd != nil; remaining, cmd, err = CommandFromBytes(current) {
                if err != nil {
                    out <- TCPCommandResult{
                        Error: fmt.Errorf("error from parsing tcp command: %w", err),
                        Command: nil,
                    }
                } else {
                    out <- TCPCommandResult{Command: cmd, Error: nil}
                }

				current = remaining
			}

			previous = current
		}
	}()

	return out
}

func (t *TCP) runSocket(conn net.Conn) {
    go func(c net.Conn) {
        defer c.Close()

        fromTcp := CommandParser(conn)
        defer func() {
            log.Println("connection has closed")
        }()

        timer := time.NewTicker(100 * time.Millisecond)

        for {
            select {
            case commandWrapper := <-fromTcp:

                if commandWrapper.Error != nil {
                    log.Printf("error from command parsing: %v\n", commandWrapper.Error)
                    return
                }

                cmd := commandWrapper.Command

                t.FromSockets <- cmd

                if cmd.Command == 'c' || cmd.Command == 'e' {
                    return
                }

            case <-timer.C:
                // what to do here?
            }
        }
    }(conn)
}

func (t *TCP) listen() {
	for {

		conn, err := t.listener.Accept()

		if err != nil {
            // TODO: Logging?
            log.Printf("server shutting down")
            break;
		}

		newConn := NewConnection(conn)
		err = send_cmds(newConn, t.welcomes)
		if err != nil {
			// TODO: Should i do something with the conn?
			continue
		}
		t.sockets = append(t.sockets, newConn)
        t.runSocket(conn)
	}
}
