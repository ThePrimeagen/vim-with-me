package ascii_buffer

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

type AsciiRLE struct {
	buffer []byte
	idx    int
	curr   byte
	count  int
}

func NewAsciiRLE() *AsciiRLE {
	return &AsciiRLE{
		buffer: make([]byte, 0, 256),
		idx:    0,
		curr:   0,
		count:  0,
	}
}

func (a *AsciiRLE) Length() int {
	return a.idx
}

func (a *AsciiRLE) place() {
    if a.count == 0 {
        return
    }

    if a.idx + 1 >= len(a.buffer) {
        a.buffer = append(a.buffer, make([]byte, 128, 128)...)
    }

	a.buffer[a.idx] = byte(a.count)
	a.buffer[a.idx+1] = a.curr
	a.idx += 2
}

func (a *AsciiRLE) Debug() {
	fmt.Printf("rle: ")
	for i := 0; i < a.idx; i += 2 {
		fmt.Printf("%s(%d) ", string(a.buffer[i+1]), a.buffer[i])
	}
	fmt.Println()
}

func (a *AsciiRLE) Reset() {
	a.idx = 0
}

func (a *AsciiRLE) Write(data []byte) {
	assert.Assert(len(data) > 0, "AsciiRLE#Write received 0 len data array")
	if a.count == 0 {
		a.curr = data[0]
        a.count = 1
	}

	for i := 1; i < len(data); i++ {
		if a.count < 255 && a.curr == data[i] {
            a.count++
			continue
		}

		a.place()
        a.curr = data[i]
        a.count = 1
	}
}

func (a *AsciiRLE) Finish() {
    a.place()
}

func (a *AsciiRLE) Bytes() []byte {
	return a.buffer[:a.idx]
}
