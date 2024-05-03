package ascii_buffer

import "github.com/theprimeagen/vim-with-me/pkg/assert"

type AsciiFrame struct {
	idx      int
	buffer   []byte
	previous []byte
	rle      *AsciiRLE
}

func NewAsciiFrame(row, col, renders int) *AsciiFrame {
	length := row * col
	return &AsciiFrame{
		buffer:   make([]byte, length, length),
		previous: make([]byte, length, length),
		rle:      NewAsciiRLE(),
	}
}

func (a *AsciiFrame) Write(data []byte) (int, error) {
	assert.Assert(a.idx+len(data) <= len(a.buffer), "attempting to encode too much")

	copy(a.buffer[a.idx:], data)

	return len(data), nil
}

func (a *AsciiFrame) Frame() []byte {
	return a.previous
}

func (a *AsciiFrame) Render() []byte {
	assert.Assert(a.idx == len(a.buffer), "you can only call render once you have created a full frame")

	a.idx = 0

	copy(a.previous, a.buffer)

	return a.previous
}
