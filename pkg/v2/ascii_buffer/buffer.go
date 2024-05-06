package ascii_buffer

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

type AsciiFrame struct {
	Buffer   []byte
}

func NewAsciiFrame(row, col int) *AsciiFrame {
	length := row * col
	return &AsciiFrame{
		Buffer:   make([]byte, length, length),
	}
}

func (a *AsciiFrame) Length() int {
    return len(a.Buffer)
}

func (a *AsciiFrame) PushFrame(data []byte) *AsciiFrame {
    assert.Assert(len(data) == len(a.Buffer), fmt.Sprintf("the frame MUST be the same size as the AsciiBuffer: Expected: %d Received: %d", len(data), len(a.Buffer)))
	copy(a.Buffer, data)

    return a
}

func (framer *AsciiFrame) Diff(a *AsciiFrame, b *AsciiFrame) *AsciiFrame {
    // TODO: Obvi perf win, just don't know how in go

    for i, aByte := range a.Buffer {
        framer.Buffer[i] = aByte ^ b.Buffer[i]
    }

    return framer
}
