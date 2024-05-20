package tcp

import (
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

var id int = 0

const MAX_PACKET_LENGTH = 10_000

var MAX_PACKET_ERROR = errors.New("maximum packet size exceeded")

type Connection struct {
	Reader FrameReader
	Writer FrameWriter
	Id     int
	conn   net.Conn

	previous []byte
	scratch  [1024]byte
}

// TODO: How do i close?
func NewConnection(conn net.Conn, id int) Connection {
	return Connection{
		Reader: NewFrameReader(conn),
		Writer: NewFrameWriter(conn),
		Id:     id,
		conn:   conn,
	}
}

func (c *Connection) Close() {
	c.conn.Close()
}

func (c *Connection) Next() (*TCPCommand, error) {
	cmdBytes, err := c.Reader.Read()
	if err != nil {
		return nil, err
	}

	var cmd TCPCommand
	err = cmd.UnmarshalBinary(cmdBytes)

	if err != nil {
		return nil, err
	}

	return &cmd, nil
}

type FrameReader struct {
	Reader   io.Reader
	previous []byte
	scratch  []byte
}

func NewFrameReader(Reader io.Reader) FrameReader {
	return FrameReader{
		Reader:   Reader,
		previous: []byte{},
		scratch:  make([]byte, 1024),
	}
}

func (f *FrameReader) canParse(data []byte) bool {
	if len(data) < HEADER_SIZE {
		return false
	}

	length := int(binary.BigEndian.Uint16(data[2:]))
	return len(data) >= HEADER_SIZE+length
}

func (f *FrameReader) packetLen(data []byte) int {
	if len(data) < HEADER_SIZE {
		return -1
	}

	return int(binary.BigEndian.Uint16(data[2:])) + HEADER_SIZE
}

func (f *FrameReader) Read() ([]byte, error) {
	for {
		n := f.packetLen(f.previous)
		if n > MAX_PACKET_LENGTH {
			return nil, fmt.Errorf("FrameReader#Read %d %w", n, MAX_PACKET_ERROR)
		}

		if f.canParse(f.previous) {
			out := f.previous[:n]
			remaining := len(f.previous) - n
			new := make([]byte, remaining, remaining)
			copy(new, f.previous[n:])
			f.previous = new
			return out, nil
		}

		n, err := f.Reader.Read(f.scratch)

		if err != nil {
			return nil, err
		}

		f.previous = append(f.previous, f.scratch[:n]...)
	}
}

type Bytes interface {
	Bytes() []byte
}

type FrameWriter struct {
	Writer io.Writer
}

func NewFrameWriter(Writer io.Writer) FrameWriter {
	return FrameWriter{
		Writer: Writer,
	}
}

func (w *FrameWriter) Write(bytes encoding.BinaryMarshaler) error {
	data, err := bytes.MarshalBinary()
	if err != nil {
		return err
	}

	length := len(data)

	for length > 0 {
		n, err := w.Writer.Write(data)

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
