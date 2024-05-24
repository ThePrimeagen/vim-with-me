package doom

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser/display"
)

var comma = []byte{','}
var newLine = []byte{'\n'}
var rows = []byte(string("y_res: "))
var cols = []byte(string("x_res: "))
var escape = []byte{27}

func parseOutNumber(data []byte) int {
	num := data[:bytes.Index(data, comma)]

	cols, err := strconv.Atoi(string(num))
	assert.Assert(err == nil, fmt.Sprintf("unable to parse out cols: %d -- %s", cols, string(data)))

	return cols
}

type DoomAsciiHeaderParser struct {
	scratch []byte
}

func newDoomAsciiHeaderParser() DoomAsciiHeaderParser {
	return DoomAsciiHeaderParser{
		scratch: make([]byte, 0),
	}
}

func (d *DoomAsciiHeaderParser) GetDims() (int, int) {
	colIdx := bytes.Index(d.scratch, cols)
	rowIdx := bytes.Index(d.scratch, rows)
	assert.Assert(colIdx != -1, "unable to parse out col numbers")
	assert.Assert(rowIdx != -1, "unable to parse out row numbers")

	rows := parseOutNumber(d.scratch[rowIdx+len(rows):])
	cols := parseOutNumber(d.scratch[colIdx+len(cols):]) * 2 // because

	return rows, cols
}

func (d *DoomAsciiHeaderParser) Write(data []byte) (int, error) {

	// Done parsing header
	escapeIdx := bytes.Index(data, escape)
	if escapeIdx == -1 {
		d.scratch = append(d.scratch, data...)
		return len(data), nil
	}

	d.scratch = append(d.scratch, data[:escapeIdx]...)
	return escapeIdx - 1, nil

}

type Doom struct {
	Framer *ansiparser.AnsiFramer

	header DoomAsciiHeaderParser
	ready  chan struct{}

	Rows int
	Cols int
}

func NewDoom() *Doom {
	doom := &Doom{
		header: newDoomAsciiHeaderParser(),
		Framer: nil,

		ready: make(chan struct{}, 0),

		Rows: 0,
		Cols: 0,
	}

	return doom
}

func (d *Doom) Ready() <-chan struct{} {
	return d.ready
}

func (d *Doom) Frames() chan display.Frame {
	return d.Framer.Frames()
}

func (d *Doom) Write(data []byte) (int, error) {
	consumed := 0
	if d.Framer == nil {
		headerBytes, err := d.header.Write(data)
		assert.Assert(err == nil, "doom ascii header should never fail")

		if headerBytes == len(data) {
			return headerBytes, nil
		}

		consumed += headerBytes + 1

		rows, cols := d.header.GetDims()

		d.Rows = rows
		d.Cols = cols
		d.Framer = ansiparser.
			NewFramer().
			WithDim(rows, cols).
			WithFrameStart([]byte("[;H")).
			WithInputStart([]byte{'7'})

		data = data[headerBytes+1:]

		d.ready <- struct{}{}
		close(d.ready)
	}

	n, err := d.Framer.Write(data)
	return n + consumed, err
}
