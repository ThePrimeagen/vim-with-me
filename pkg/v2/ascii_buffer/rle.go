package ascii_buffer

import (
	"errors"
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

type AsciiRLE struct {
	buffer  []byte
	idx     int
	errored bool
}

func NewAsciiRLE(maxSize int) *AsciiRLE {
	return &AsciiRLE{
		buffer: make([]byte, maxSize, maxSize),
		idx:    0,
        errored: false,
	}
}

func (a *AsciiRLE) Length() int {
	return a.idx
}

func (a *AsciiRLE) place(idx, count int, value byte) error {
    if idx+1 >= len(a.buffer) {
        a.errored = true
        return errors.New("RLE produced a buffer larger than original.  Failed")
    }

    a.buffer[idx] = byte(count)
    a.buffer[idx+1] = value

    return nil
}

func (a *AsciiRLE) Debug() {
    fmt.Printf("rle: ")
    for i := 0; i < a.idx; i += 2 {
        fmt.Printf("%s(%d) ", string(a.buffer[i + 1]), a.buffer[i])
    }
    fmt.Println()
}

func (a *AsciiRLE) RLE(frame *AsciiFrame) error {
    asciiBuffer := frame.buffer
    length := len(asciiBuffer)
	assert.Assert(length&1 == 0, "you must hand an even length buffer")

	count := 0
	idx := 0
    a.errored = false

	for i := 0; i < len(asciiBuffer)-1; i++ {
		count++

		if count < 255 && asciiBuffer[i] == asciiBuffer[i + 1] {
            continue
		}

        if err := a.place(idx, count, asciiBuffer[i]); err != nil {
            return err
        }

		count = 0
		idx += 2
	}

    if asciiBuffer[length - 2] == asciiBuffer[length - 1] {
        count++
    } else {
        count = 1
    }
    a.place(idx, count, asciiBuffer[length - 1])
    a.idx = idx + 2

	return nil
}

func (a *AsciiRLE) Bytes() []byte {
    if a.errored {
        return nil
    }

	return a.buffer[:a.idx]
}
