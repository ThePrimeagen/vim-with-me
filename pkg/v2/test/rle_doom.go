package main

import (
	"fmt"
	"os"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
)

type Encoder struct {
	buffs []*ascii_buffer.AsciiFrame
	idx   int

	diff *ascii_buffer.AsciiFrame
	rle  *ascii_buffer.AsciiRLE

	totalData        int
	deflatedData     int
	diffDeflatedData int
}

func NewEncoder(rows, cols int) *Encoder {
	buffs := []*ascii_buffer.AsciiFrame{
		ascii_buffer.NewAsciiFrame(rows, cols),
		ascii_buffer.NewAsciiFrame(rows, cols),
	}
	return &Encoder{
		buffs:        buffs,
		rle:          ascii_buffer.NewAsciiRLE(buffs[0].Length()),
		diff:         ascii_buffer.NewAsciiFrame(66, 212),
		idx:          0,
		totalData:    0,
		deflatedData: 0,
		diffDeflatedData: 0,
	}
}

func (e *Encoder) Push(data []byte) {
	curr := e.buffs[e.idx%2]
	prev := e.buffs[(e.idx+1)%2]

	curr.PushFrame(data)
	assert.Assert(e.rle.RLE(curr) == nil, "unable to push data")

	e.diff.Diff(curr, prev)

	fmt.Printf("curr: %d -- rle: %d\n", curr.Length(), e.rle.Length())
	e.deflatedData += e.rle.Length()

	assert.Assert(e.rle.RLE(e.diff) == nil, "unable to push data")

	fmt.Printf("diff: %d -- rle: %d\n", e.diff.Length(), e.rle.Length())
	e.totalData += len(data)
    e.diffDeflatedData += e.rle.Length()

	e.idx++
}

func main() {
	data, err := os.ReadFile("./pkg/v2/ansi_parser/doomtest_very_large")
	assert.Assert(nil == err, "couldn't read doom large")

	doomHeader := ansiparser.NewDoomAsciiHeaderParser()
	n, err := doomHeader.Write(data)
	assert.Assert(nil == err, "couldn't get out the dimensions")

	data = data[n:]
	rows, cols := doomHeader.GetDims()
	assert.Assert(rows == 66, "rows not equal to 66")
	assert.Assert(cols == 212, "cols not equal to 212")

	doomAscii := ansiparser.NewDoomAnsiFramer(rows, cols)
	assert.Assert(err == nil, "errored on reading file")

	go func() {
		doomAscii.Write(data)
	}()

	frames := doomAscii.Frames()
	chars := NewEncoder(rows, cols)
	colors := NewEncoder(rows, cols)

	for i := range 2372 {
		frame := <-frames
		fmt.Printf("Chars at %d\n", i)
		chars.Push(frame.Chars)
		fmt.Printf("colors at %d\n", i)
		colors.Push(frame.Color)
	}

	fmt.Printf("total: %d\n", chars.totalData)

	fmt.Printf("color: %d -- chars: %d\n", colors.deflatedData, chars.deflatedData)
	fmt.Printf("color: %f -- chars: %f\n",
		float64(colors.deflatedData)/float64(colors.totalData),
		float64(chars.deflatedData)/float64(chars.totalData))

	fmt.Printf("color: %d -- chars: %d\n", colors.diffDeflatedData, chars.diffDeflatedData)
	fmt.Printf("color: %f -- chars: %f\n",
		float64(colors.diffDeflatedData)/float64(colors.totalData),
		float64(chars.diffDeflatedData)/float64(chars.totalData))

}
