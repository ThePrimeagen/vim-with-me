package tcp

import (
	"encoding"
	"encoding/binary"
	"net"
)

var id int = 0

// wraps underline net.Conn
type Connection struct {
	net.Conn
	Id       int
	scratch  []byte
	previous []byte
}

// TODO: How do i close?
func NewConnection(conn net.Conn) Connection {
	id++
	return Connection{
		Conn:     conn,
		Id:       id,
		scratch:  make([]byte, 1024),
		previous: make([]byte, 0),
	}
}

// refering to the underline connection
func (c *Connection) Close() error {
	return c.Conn.Close()
}

func (c *Connection) Next() (*TCPCommand, error) {

	for {
		n, err := c.Read(c.scratch)
		if err != nil {
			return nil, err
		}

		c.previous = append(c.previous, c.scratch[:n]...)

		packetN := c.parsePacket()
		if packetN == 0 {
			continue
		}

		var cmd TCPCommand
		err = cmd.UnmarshalBinary(c.previous[:packetN])
		if err != nil {
			return nil, err
		}

		c.previous = c.previous[packetN:]

		return &cmd, nil
	}
}

func (c *Connection) parsePacket() int {
	if len(c.previous) < HEADER_SIZE {
		return 0
	}

	length := int(binary.BigEndian.Uint16(c.previous[2:]))
	if len(c.previous) < HEADER_SIZE+length {
		return 0
	}

	return HEADER_SIZE + length
}

func (c *Connection) Read(b []byte) (int, error) {
	return c.Conn.Read(b)
}

type Bytes interface {
	Bytes() []byte
}

func (c *Connection) Write(bytes encoding.BinaryMarshaler) error {
	data, err := bytes.MarshalBinary()
	if err != nil {
		return err
	}

	length := len(data)

	for length > 0 {
		n, err := c.Conn.Write(data)

		// TODO: what's the error that means i need to wait for a moment to
		// write more?
		if err != nil {
			return err
		}

		data = data[n:]
		length -= n
	}
	return nil
}
