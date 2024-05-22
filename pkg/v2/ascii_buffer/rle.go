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
		buffer: nil,
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

	if a.idx+1 >= len(a.buffer) {
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

func (a *AsciiRLE) Reset(buf []byte) {
	a.idx = 0
	a.buffer = buf
}

func (a *AsciiRLE) String() string {
    return fmt.Sprintf("AsciiRLE: curr=%d count=%d bufLen=%d buf=%+v", a.curr, a.count, a.idx, a.buffer[:a.idx])
}

func (a *AsciiRLE) Write(data []byte) {
	assert.Assert(a.buffer != nil, "AsciiRLE#Write needs a buffer to write into")
	assert.Assert(len(data) > 0, "AsciiRLE#Write received 0 len data array")

	i := 0

	if a.count == 0 {
		a.curr = data[0]
		a.count = 1
		i = 1
	}

	for ; i < len(data); i++ {
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
    a.count = 0
}

func (a *AsciiRLE) Bytes() []byte {
	return a.buffer[:a.idx]
}

func ExpandRLE(frame []byte, offset int, previous []byte, out []byte) {
    assert.Assert(len(previous) <= len(out), "cannot decode an rle frame into a smaller frame (out is smaller than previous)", "previous", len(previous), "out", len(out))
    idx := 0
    for i := offset; i < len(frame); i += 2 {
        repeat := frame[i]
        char := frame[i + 1]
        for count := byte(0); count < repeat; count++ {
            out[idx] = char ^ previous[idx]
            idx++
        }
    }

}
