package tcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
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

func CommandFromBytes(b []byte) (TCPCommand, error) {
    if len(b) < MINIMUM {
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
	sockets     []*Connection
	welcomes    []*TCPCommand
	listener    net.Listener
}

func (t *TCP) Length() int {
    return len(t.sockets)
}

func (t *TCP) Debug() string {
    out := ""
    for _, k := range t.sockets {
        out += fmt.Sprintf("id=%d closed=%v ::", k.id, k.closed)
    }
    return out
}

func (t *TCP) Welcome(cmd *TCPCommand) {
	t.welcomes = append(t.welcomes, cmd)
}

// TODO: Think about project level logging and the ability to enable debug
// logging
func (t *TCP) Send(cmd *TCPCommand) {
	log.Printf("send %+v\n", cmd)
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

func (t *TCP) runSocket(conn *Connection) {
    go func() {
        defer conn.Close()

        fromTcp := CommandParser(conn.conn)

        defer func() {
            log.Printf("connection has closed: %d\n", id)
        }()

        timer := time.NewTicker(100 * time.Millisecond)

        for {
            select {
            case commandWrapper := <-fromTcp:

                if errors.Is(commandWrapper.Error, io.EOF) {
                    log.Printf("EOF received, closing connection", commandWrapper.Error)
                    conn.closed = true
                    return
                }

                if commandWrapper.Error != nil {
                    log.Printf("error from command parsing: %v\n", commandWrapper.Error)
                    return
                }

                cmd := commandWrapper.Command

                t.FromSockets <- cmd

                if cmd.Command == 'c' || cmd.Command == 'e' {
                    conn.Close()
                    return
                }

            case <-timer.C:
                // what to do here?
            }
        }
    }()
}

func (t *TCP) listen() {
	for {
        log.Printf("waiting for server incoming connection")
		conn, err := t.listener.Accept()

		if err != nil {
            // TODO: Logging?
            log.Printf("server shutting down: %+v", err)
            break
		}

		newConn := NewConnection(conn)
        slog.Info("new connection received", "id", newConn.id)
		err = send_cmds(&newConn, t.welcomes)
		if err != nil {
            log.Printf("this doesnt't happen: %+v\n", err)
			// TODO: Should i do something with the conn?
			continue
		}
		t.sockets = append(t.sockets, &newConn)
        t.runSocket(&newConn)
	}
    log.Println("server finished running for loop")
}
