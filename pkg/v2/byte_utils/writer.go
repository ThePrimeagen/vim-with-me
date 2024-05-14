package byteutils

import "github.com/theprimeagen/vim-with-me/pkg/v2/assert"

type ByteWriter interface {
    write(buf []byte, offset int, value int)
}

func Write16(buf []byte, offset, value int) {
    assert.Assert(len(buf) > offset + 1, "you cannot write outside of the buffer")
    hi := (value & 0xFF00) >> 8
    lo := value & 0xFF
    buf[offset] = byte(hi)
    buf[offset + 1] = byte(lo)
}

type U8Writer struct {}
func (b *U8Writer) write(buf []byte, offset int, value int) {
    assert.Assert(len(buf) > offset, "you cannot write outside of the buffer")
    buf[offset] = byte(value)
}

type U16Writer struct {}
func (b *U16Writer) write(buf []byte, offset, value int) {
    Write16(buf, offset, value)
}
