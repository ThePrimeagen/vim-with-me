package net

import (
	"fmt"

	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

const VERSION = byte(1)

type BaseFrameType byte

const (
	OPEN BaseFrameType = iota
	BRIGHTNESS_TO_ASCII
	FRAME
)

const HEADER_SIZE = 5

type Encodeable interface {
	Type() byte
	Into(into []byte, offset int) (int, error)
}

type Frameable struct {
	Item Encodeable
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
	return &Frameable{
		Item: &Open{Rows: rows, Cols: cols},
	}
}

var count = 0

func (f *Frameable) Into(into []byte, offset int) (int, error) {
	into[offset] = VERSION
	into[offset+1] = f.Item.Type()
	into[offset+2] = byte(count % 16)
    fmt.Printf("seq: %d\n", byte(count % 16))

	count++

	n, err := f.Item.Into(into, offset+HEADER_SIZE)
	if err != nil {
		return 0, err
	}

	byteutils.Write16(into, offset+3, n)

	fmt.Printf("Frameable#Into: %d + 4 for encoding HEADER\n", n)

	// bytes + 5 for header
	return n + HEADER_SIZE, nil
}
