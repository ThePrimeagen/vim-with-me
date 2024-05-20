package net

const VERSION = byte(1)

type BaseFrameType byte

const (
	OPEN BaseFrameType = iota
	BRIGHTNESS_TO_ASCII
	FRAME
)

type Encodeable interface {
	Type() byte
	Into(into []byte, offset int) error
}

type Frameable struct {
	Item Encodeable
}

func (f *Frameable) Into(into []byte, offset int) error {
	into[offset] = VERSION
	into[offset+1] = f.Item.Type()
	return f.Item.Into(into, offset+2)
}
