package net

import (
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

const VERSION = byte(1)

type BaseFrameType byte

const (
	OPEN BaseFrameType = iota
	BRIGHTNESS_TO_ASCII
	FRAME
)

func typeToString(t byte) string {
	switch t {
	case byte(OPEN):
		return "open"
	case byte(BRIGHTNESS_TO_ASCII):
		return "brightness_to_ascii"
	case byte(FRAME):
		return "frame"
	}
	return "unknown"
}

const HEADER_SIZE = 6

type Encodeable interface {
	Type() byte
	Into(into []byte, offset int) (int, error)
}

type Frameable struct {
	Item Encodeable
	seq  byte
}

func NewFrameable(item Encodeable) *Frameable {
    return &Frameable{
        Item: item,
        seq: byte(nextSeqId()),
    }

}

type Open struct {
	Rows int
	Cols int
}

func (o *Open) Into(into []byte, offset int) (int, error) {
	byteutils.Write16(into, offset, o.Rows)
	byteutils.Write16(into, offset+2, o.Cols)
	return 4, nil
}

func (o *Open) Type() byte {
	return byte(OPEN)
}

func CreateOpen(rows, cols int) *Frameable {
	return NewFrameable(&Open{Rows: rows, Cols: cols})
}

func (f *Frameable) Into(into []byte, offset int) (int, error) {
	frameHeader(into, offset, f.Item.Type(), f.seq)

	n, err := f.Item.Into(into, offset+HEADER_SIZE)
	if err != nil {
		return 0, err
	}

	byteutils.Write16(into, offset+4, n)

	// bytes + 5 for header
	return n + HEADER_SIZE, nil
}
