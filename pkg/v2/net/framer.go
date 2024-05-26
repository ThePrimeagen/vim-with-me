package net

import (
	"errors"
	"fmt"
	"io"

	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

var FramerVersionMismatch = errors.New("version mismatch")

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

	if len(b.curr) < HEADER_SIZE {
		return nil
	}

	length := byteutils.Read16(b.curr, 4)
	totalLength := length + HEADER_SIZE
	remaining := len(b.curr) - totalLength

	if len(b.curr) < totalLength {
		return nil
	}

	b.ch <- &Frame{
		CmdType: b.curr[1],
		Seq:     b.curr[2],
		Flags:   b.curr[3],
		Data:    b.curr[HEADER_SIZE:totalLength],
	}

	copy(b.curr, b.curr[totalLength:])
	b.curr = b.curr[:remaining]

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
