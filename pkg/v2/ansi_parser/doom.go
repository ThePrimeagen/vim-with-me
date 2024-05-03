package ansiparser

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

type DoomAsciiHeaderParser struct {
	scratch []byte
}

func NewDoomAsciiHeaderParser() *DoomAsciiHeaderParser {
	return &DoomAsciiHeaderParser{
		scratch: make([]byte, 0),
	}
}

var comma = []byte{','}
var newLine = []byte{'\n'}
var rows = []byte(string("y_res: "))
var cols = []byte(string("x_res: "))
var escape = []byte{27}

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

	d.scratch = data[:escapeIdx]
	return escapeIdx, nil

}

type DoomAnsiFramer struct {
	*Ansi8BitFramer
}

func NewDoomAnsiFramer(rows, cols int) *DoomAnsiFramer {
	return &DoomAnsiFramer{
		Ansi8BitFramer: New8BitFramer().WithDim(rows, cols),
	}
}

func parseOutNumber(data []byte) int {
	num := data[:bytes.Index(data, comma)]

	cols, err := strconv.Atoi(string(num))
	assert.Assert(err == nil, fmt.Sprintf("unable to parse out cols: %d -- %s", cols, string(data)))

	return cols
}
