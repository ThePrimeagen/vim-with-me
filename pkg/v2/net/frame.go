package net

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

var nextId = 0

func nextSeqId() int {
	out := nextId
	nextId++
	return out
}

func frameHeader(data []byte, offset int, t, sequence byte) {
	data[offset] = VERSION
	data[offset+1] = t
	data[offset+2] = sequence
	data[offset+3] = 0 // flags later on
}

type Frame struct {
	CmdType byte
	Seq     byte
	Flags   byte
	Unused  byte
	Data    []byte
}

func (f *Frame) Bytes() []byte {
    out := make([]byte, HEADER_SIZE + len(f.Data))
    frameHeader(out, 0, f.Type(), f.Seq)
    byteutils.Write16(out, 4, len(f.Data))

    copy(out[HEADER_SIZE:], f.Data)

    return out
}

func (f *Frame) Sequence() byte {
	return f.Seq
}

func (f *Frame) String() string {
	return fmt.Sprintf("frame(%s): seq=%d flags=%d data=%d", typeToString(f.CmdType), f.Seq, f.Flags, len(f.Data))
}

func (f *Frame) Type() byte {
	return f.CmdType
}

func (f *Frame) Into(data []byte, offset int) (int, error) {
	assert.Assert(len(data) > HEADER_SIZE+len(f.Data), "unable to encode frame into cache packet")
	frameHeader(data, offset, f.Type(), f.Sequence())
	byteutils.Write16(data, offset+4, len(f.Data))
	copy(data[offset+HEADER_SIZE:], f.Data)

	return HEADER_SIZE + len(f.Data), nil
}
