package tcp

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTCPClient(port uint16) (*net.TCPConn, error) {
	servAddr := fmt.Sprintf("127.0.0.1:%d", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func TestTCPServer(t *testing.T) {
	port := uint16(42069)
	server, err := NewTCPServer(port)
	if err != nil {
		t.Fatalf("Error creating TCP server: %v", err)
	}

	client, err := createTCPClient(uint16(42069))
	if err != nil {
		t.Fatalf("Error creating TCP client: %v", err)
	}
	client2, err := createTCPClient(uint16(42069))
	if err != nil {
		t.Fatalf("Error creating TCP client: %v", err)
	}

	cmd := TCPCommand{
		Command: byte('g'),
		Data:    []byte("Hello World"),
	}

    _, err = client.Write(cmd.Bytes())

    if err != nil {
        t.Fatalf("Error writing cmd to the client: %v", err)
    }

    c2 := <- server.FromSockets
    assert.Equal(t, c2, cmd)

    cmd2 := TCPCommand{
        Command: byte('t'),
        Data: []byte("69:420"),
    }

    server.Send(&cmd2)

    clientCmd := CommandParser(client)
    clientCmd2 := CommandParser(client2)

    out := <- clientCmd
    out2 := <- clientCmd2

    assert.Equal(t, out, cmd2)
    assert.Equal(t, out2, cmd2)

    client.Close()

    server.Send(&cmd)
    out2 = <- clientCmd2
    assert.Equal(t, out2, cmd)
}


/*
func TestCommandParser(t *testing.T) {
    cmd := TCPCommand{
        Command: "g",
        Data: "Hello World",
    }

    cmd2 := TCPCommand{
        Command: "t",
        Data: "Goodbye, cruel world",
    }

    b := cmd.Bytes()
    b2 := cmd2.Bytes()
    reader := bytes.NewReader(append(b, b2...))

    parsedCmd := CommandParser(reader)

    c := <- parsedCmd
    c2 := <- parsedCmd
    assert.Equal(t, c, cmd)
    assert.Equal(t, c2, cmd2)
}
*/
