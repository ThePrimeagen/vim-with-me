package ansiparser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

type DoomAnsiFramer struct {
	FoundRows int
	FoundCols int

	scratch     []byte
	ansiParsing bool
	reader      io.Reader
	framer      *Ansi8BitFramer
}

func NewDoomAnsiFramer() *DoomAnsiFramer {
	return &DoomAnsiFramer{
        FoundRows: -1,
        FoundCols: -1,
		ansiParsing: false,
		framer:      nil,
        scratch: make([]byte, 0),
	}
}

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

func (d *DoomAnsiFramer) Frames() chan []byte {
    return d.framer.Frames()
}

func (d *DoomAnsiFramer) Write(data []byte) (int, error) {
	if d.ansiParsing {
		return d.framer.Write(data)
	}

    idx := 0
    for {
        escapeIdx := bytes.Index(data[idx:], escape)
        if escapeIdx == 0 {
            d.ansiParsing = true
            break
        }

        nextIdx := bytes.Index(data[idx:], newLine)
        if nextIdx == -1 {
            d.scratch = append(d.scratch, data[idx:]...)
            break
        }

        buf := data[idx:idx + nextIdx]
        if len(d.scratch) > 0 {
            buf = append(d.scratch, buf...)
            d.scratch = make([]byte, 0)
        }

        colIdx := bytes.Index(buf, cols)
        rowIdx := bytes.Index(buf, rows)
        if colIdx != -1 {
            d.FoundCols = parseOutNumber(buf[colIdx + len(cols):]) * 2 // because
        }

        if rowIdx != -1 {
            d.FoundRows = parseOutNumber(buf[rowIdx + len(rows):])
        }

        idx += nextIdx + 1
    }

    if d.ansiParsing && d.framer == nil {
        d.framer = New8BitFramer(d.FoundRows, d.FoundCols)
    }

    if idx < len(data) {
        _, err := d.framer.Write(data[idx:])
        return len(data), err
    }

    return len(data), nil
}
