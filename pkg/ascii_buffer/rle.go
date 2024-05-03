package ascii_buffer

import (
	"errors"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

type AsciiRLE struct {
	buffer  []byte
	idx     int
	errored bool
}

func NewAsciiRLE() *AsciiRLE {
	return &AsciiRLE{
		buffer: make([]byte, 0),
		idx:    0,
        errored: false,
	}
}

func (a *AsciiRLE) Length() int {
	return a.idx
}

func (a *AsciiRLE) RLE(data []byte) (int, error) {
	assert.Assert(len(data)&1 == 0, "you must hand an even length buffer")

	if len(a.buffer) < len(data) {
		a.buffer = make([]byte, len(data), len(data))
	}

	a.idx = 0
    a.errored = false

	count := byte(0)
	idx := 0

	for i := 0; i < len(data)-3; i += 2 {
		if idx+2 >= len(a.buffer) {
            a.errored = true
			return 0, errors.New("RLE produced a buffer larger than original.  Failed")
		}
		curr := int(data[i])<<8 + int(data[i+1])
		next := int(data[i+2])<<8 + int(data[i+3])
		count++

		// counts the last encode + current
		if curr == next {
			if i+4 == len(data) {
				count++
			} else {
				continue
			}
		}

		a.buffer[idx] = count
		a.buffer[idx+1] = byte(curr >> 8)
		a.buffer[idx+2] = byte(curr & 0xFF)

		count = 0
		idx += 3
	}

	a.idx = idx

	return a.idx, nil
}

func (a *AsciiRLE) Bytes() []byte {
    if a.errored {
        return make([]byte, 0)
    }

	return a.buffer[:a.idx]
}
