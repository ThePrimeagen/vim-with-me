package net

import (
	"errors"
	"fmt"
	"io"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

var FramerVersionMismatch = errors.New("version mismatch")

type Frame struct {
	CmdType byte
	Data []byte
}

func (f *Frame) Into(data []byte, offset int) (int, error) {
    assert.Assert(len(data) > HEADER_SIZE + len(f.Data), "unable to encode frame into cache packet")
    data[0] = VERSION
    data[1] = f.CmdType
    byteutils.Write16(data, 2, len(data))
    copy(data[HEADER_SIZE:], f.Data)

    return HEADER_SIZE + len(f.Data), nil
}

type ByteFramer struct {
	curr []byte
	ch   chan *Frame
}

func NewByteFramer() *ByteFramer {
	return &ByteFramer{
		curr: make([]byte, 0),
		ch:   make(chan *Frame, 10),
	}
}

func (b *ByteFramer) frame() error {
	if b.curr[0] != VERSION {
		return errors.Join(
			FramerVersionMismatch,
			fmt.Errorf("expected %d received %d", VERSION, b.curr[0]),
		)
	}

    length := byteutils.Read16(b.curr, 2)
    if len(b.curr) < length + HEADER_SIZE {
        return nil
    }

    b.ch <- &Frame{
        Type: b.curr[1],
        Data: b.curr[HEADER_SIZE:HEADER_SIZE + length],
    }

    return nil
}

func (b *ByteFramer) Frame(reader io.Reader) error {
	data := make([]byte, 1024, 1024)
	for {
		if len(b.curr) > HEADER_SIZE {
			b.frame()
		}
		n, err := reader.Read(data)
		if err != nil {
			return err
		}

        b.curr = append(b.curr, data[:n]...)
	}
}

func (b *ByteFramer) FrameChan(reader chan []byte) error {
	for {
		if len(b.curr) > HEADER_SIZE {
			b.frame()
		}
        b.curr = append(b.curr, <-reader...)
	}
}

func (b *ByteFramer) Frames() chan *Frame {
    return b.ch
}
