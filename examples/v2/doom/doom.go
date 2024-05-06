package doom

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/program"
)

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
	return escapeIdx - 1, nil

}

type Doom struct {
	DoomAsciiHeaderParser
	*ansiparser.Ansi8BitFramer

	program program.Program
	ready   chan struct{}
}

func NewDoom(program program.Program) *Doom {
	doom := &Doom{
		DoomAsciiHeaderParser: newDoomAsciiHeaderParser(),
		Ansi8BitFramer:        nil,

		program: program,
		ready:   make(chan struct{}),
	}

	go func() {
		_, _ = io.Copy(doom, &program)
	}()

	return doom
}

func (d *Doom) Ready() <-chan struct{} {
	return d.ready
}

func (d *Doom) Write(data []byte) (int, error) {
	if d.Ansi8BitFramer == nil {
		headerBytes, err := d.DoomAsciiHeaderParser.Write(data)
		assert.Assert(err == nil, "doom ascii header should never fail")

		if headerBytes == len(data) {
			return headerBytes, nil
		}

		rows, cols := d.DoomAsciiHeaderParser.GetDims()
		d.Ansi8BitFramer = ansiparser.New8BitFramer().WithDim(rows, cols)
		data = data[headerBytes+1:]

		d.ready <- struct{}{}
        close(d.ready)
	}

	return d.Ansi8BitFramer.Write(data)
}
