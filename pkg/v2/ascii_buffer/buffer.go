package ascii_buffer

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

type AsciiFrame struct {
	buffer   []byte
}

func NewAsciiFrame(row, col int) *AsciiFrame {
	length := row * col
	return &AsciiFrame{
		buffer:   make([]byte, length, length),
	}
}

func (a *AsciiFrame) Length() int {
    return len(a.buffer)
}

func (a *AsciiFrame) PushFrame(data []byte) *AsciiFrame {
    assert.Assert(len(data) == len(a.buffer), fmt.Sprintf("the frame MUST be the same size as the AsciiBuffer: Expected: %d Received: %d", len(data), len(a.buffer)))
	copy(a.buffer, data)

    return a
}

func (framer *AsciiFrame) Diff(a *AsciiFrame, b *AsciiFrame) *AsciiFrame {
    // TODO: Obvi perf win, just don't know how in go

    for i, aByte := range a.buffer {
        framer.buffer[i] = aByte ^ b.buffer[i]
    }

    return framer
}
