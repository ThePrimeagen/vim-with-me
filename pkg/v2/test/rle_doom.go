package main

import (
	"fmt"
	"os"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
)

func main() {
	doomAscii := ansiparser.NewDoomAnsiFramer()
	data, err := os.ReadFile("./pkg/v2/ansi_parser/doomtest_large")
	assert.Assert(err == nil, "errored on reading file")

	go func() {
		doomAscii.Write(data)
	}()

	frames := doomAscii.Frames()

	buffers := []*ascii_buffer.AsciiFrame{
		ascii_buffer.NewAsciiFrame(66, 212),
		ascii_buffer.NewAsciiFrame(66, 212),
	}
    rle := ascii_buffer.NewAsciiRLE(buffers[0].Length())

    // diff := ascii_buffer.NewAsciiFrame(66, 212)
	idx := 0
    errCount := 0

	for range 129 {
		frame := <-frames
        fmt.Printf("frame %d %d\n", len(frame.Chars), len(frame.Color))

		curr := buffers[idx%2]
		// prev := buffers[(idx+1)%2]

        curr.PushFrame(frame.Chars)
        if err := rle.RLE(curr); err != nil {
            errCount++
        }

        fmt.Printf("curr: %d -- rle: %d\n", curr.Length(), rle.Length())

        idx++
	}
}
