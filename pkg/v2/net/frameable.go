package net

import byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"

const VERSION = byte(1)

type BaseFrameType byte

const (
	OPEN BaseFrameType = iota
	BRIGHTNESS_TO_ASCII
	FRAME
)

const HEADER_SIZE = 4

type Encodeable interface {
	Type() byte
	Into(into []byte, offset int) (int, error)
}

type Frameable struct {
	Item Encodeable
}

func (f *Frameable) Into(into []byte, offset int) (int, error) {
	into[offset] = VERSION
	into[offset+1] = f.Item.Type()

	n, err := f.Item.Into(into, offset+4)
	if err != nil {
		return 0, err
	}

	byteutils.Write16(into, offset+2, n)

    // bytes + 4 for header
	return n + 4, nil
}
