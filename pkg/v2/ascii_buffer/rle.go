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

func RLE(data []byte) ([]byte, int, error) {
	assert.Assert(len(data)&1 == 0, "you must hand an even length buffer")

	count := byte(0)
	idx := 0
    buffer := make([]byte, 0, len(data))

	for i := 0; i < len(data)-3; i += 2 {
		if idx+2 >= len(buffer) {
			return nil, 0, errors.New("RLE produced a buffer larger than original.  Failed")
		}
		curr := int(data[i])<<8 + int(data[i+1])
		next := int(data[i+2])<<8 + int(data[i+3])
		count++

		if curr == next {
            // oddity: i have to count curr + next if we are at the last encode
			if i+4 == len(data) {
				count++
			} else {
				continue
			}
		}

		buffer[idx] = count
		buffer[idx+1] = byte(curr >> 8)
		buffer[idx+2] = byte(curr & 0xFF)

		count = 0
		idx += 3
	}

	return buffer, idx, nil
}

func (a *AsciiRLE) Bytes() []byte {
    if a.errored {
        return nil
    }

	return a.buffer[:a.idx]
}
