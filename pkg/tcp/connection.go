package tcp

import (
	"encoding"
	"encoding/binary"
	"io"
	"net"
)

var id int = 0
type Connection struct {
    FrameReader
    FrameWriter
    Id int
}

// TODO: How do i close?
func NewConnection(conn net.Conn) Connection {
    id++
	return Connection{
        FrameReader: NewFrameReader(conn),
        FrameWriter: NewFrameWriter(conn),
        Id: id,
	}
}

func (c *Connection) Next() (*TCPCommand, error) {
    _, err := c.Read(c.scratch)

    if err != nil {
        return nil, err
    }

    var cmd TCPCommand
    err = cmd.UnmarshalBinary(c.scratch)

    if err != nil {
        return nil, err
    }

    return &cmd, nil
}

type FrameReader struct {
	reader   io.Reader
	previous []byte
	scratch  []byte
}

func NewFrameReader(reader io.Reader) FrameReader {
	return FrameReader{
		reader:   reader,
		previous: []byte{},
		scratch:  make([]byte, 1024),
	}
}

func (f *FrameReader) parsePacket(data []byte) int {
	if len(data) < HEADER_SIZE {
		return 0
	}

	length := int(binary.BigEndian.Uint16(data[2:]))
	if len(data) < HEADER_SIZE+length {
		return 0
	}

	return HEADER_SIZE + length
}

func (f *FrameReader) Read(data []byte) (int, error) {
	for {
		n := f.parsePacket(f.previous)
		if n > 0 {
			copy(data, f.previous[:n])
			f.previous = f.previous[n:]
			return n, nil
		}

		n, err := f.reader.Read(f.scratch)
		if err != nil {
			return 0, err
		}

		f.previous = append(f.previous, f.scratch[:n]...)
	}
}

type Bytes interface {
	Bytes() []byte
}

type FrameWriter struct {
	writer io.Writer
}

func NewFrameWriter(writer io.Writer) FrameWriter {
	return FrameWriter{
		writer: writer,
	}
}

func (w *FrameWriter) Write(bytes encoding.BinaryMarshaler) error {
	data, err := bytes.MarshalBinary()
    if err != nil {
        return err
    }

	length := len(data)

	for length > 0 {
		n, err := w.writer.Write(data)

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

