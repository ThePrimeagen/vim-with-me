package ansiparser

import (
	//"github.com/leaanthony/go-ansi-parser"
	"bytes"
	"io"
	"strings"

	"github.com/leaanthony/go-ansi-parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser/display"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/encoding"
)

type Ansi8BitFramer struct {
	rows int
	cols int

	count     int
	readCount int

	debug       io.Writer
	ch          chan display.Frame
	currentIdx  int
	currentCol  int
	empty       int
	CurrentRow  int
	buffer      []byte
	colorOffset int
	scratch     []byte
}

// TODO: 2 errors, row 21 seems to start without an escape
// TODO: we need to parse out each of the ansi chunks and discard any errors
// TODO: perhaps i need to change the ansi parsing library?
func parseAnsiRow(data string) []*ansi.StyledText {
	out := make([]*ansi.StyledText, 0)
	for len(data) > 0 {
		nextEsc := strings.Index(data[1:], "") + 1

		if nextEsc == 0 {
			nextEsc = len(data)
		}

		styles, err := ansi.Parse(data[:nextEsc])
		if err == nil {
			out = append(out, styles...)
		} else {
			idx := nextEsc - 1
			char := data[idx]

			for {
				next := data[idx-1]
				if next != char || (idx-1) < 0 {
					break
				}

				idx--
			}

			length := ((nextEsc - idx) / 2) * 2
			if length > 0 {
				str := data[nextEsc-length : nextEsc]
				out = append(out, &ansi.StyledText{
					Label: str,
				})
			}
		}

		data = data[nextEsc:]
	}

	return out
}

// TODO: I could also use a ctx to close out everything
func New8BitFramer() *Ansi8BitFramer {

	// 1 byte color, 1 byte ascii
	return &Ansi8BitFramer{
		ch:         make(chan display.Frame, 10),
		currentIdx: 0,
		currentCol: 0,
		count:      0,
		empty:      0,
		CurrentRow: 0, // makes life easier
		buffer:     make([]byte, 0, 0),
		scratch:    make([]byte, 0),
	}
}

func (a *Ansi8BitFramer) WithDim(rows, cols int) *Ansi8BitFramer {
	length := rows * cols
	a.rows = rows
	a.cols = cols

	a.colorOffset = length
	a.buffer = make([]byte, length*2, length*2)

	return a
}

func remainingIsRegisteredNurse(data []byte) bool {
	if len(data) != 3 {
		return false
	}

	return data[1] == '\r' && data[2] == '\n'
}

func (framer *Ansi8BitFramer) place(color, char byte) {
    if framer.currentCol == framer.cols {
        return
    }
	assert.Assert(framer.colorOffset+framer.currentIdx < len(framer.buffer), "place failed", "color", color, "byte", char, "currentIdx", framer.currentIdx, "data length", len(framer.buffer))
	if framer.currentIdx == 0 {
		framer.buffer[framer.currentIdx] = byte(framer.count % 10)
	} else {
		framer.buffer[framer.currentIdx] = char
	}
	framer.buffer[framer.colorOffset+framer.currentIdx] = color
	framer.currentIdx++
	framer.currentCol++
}

func (framer *Ansi8BitFramer) fillRemainingRow() {
    framer.empty += framer.cols - framer.currentCol
	for framer.currentCol < framer.cols {
		framer.place(0, ' ')
	}
}

func (framer *Ansi8BitFramer) Write(data []byte) (int, error) {
	read := len(data)
	if framer.debug != nil {
		framer.debug.Write(data)
	}

	scratchLen := len(framer.scratch)

	if scratchLen != 0 {

		// this is terrible for perf
		data = append(framer.scratch, data...)
		framer.scratch = make([]byte, 0)
	}

	for len(data) > 0 {
		nextLine := bytes.Index(data, []byte{'\r', '\n'})
		if nextLine == -1 {
			framer.scratch = data
			break
		}

		line := data[:nextLine]
		data = data[nextLine+2:]

		styles := parseAnsiRow(string(line))

		assert.AddAssertData("line", line)
		assert.AddAssertData("styles", styles)
		assert.AddAssertData("currentIdx", framer.currentIdx)
		assert.AddAssertData("currentRow", framer.CurrentRow)
		assert.AddAssertData("currentCol", framer.currentCol)
		assert.Assert(styles != nil, "i should never have a nil row")
        assert.Assert(framer.currentCol == 0, "we should always start every line with current col as 0")

		for _, style := range styles {
			color := uint(255)
			if style.FgCol != nil {
				color = encoding.RGBTo8BitColor(style.FgCol.Rgb)
			}

			for _, char := range style.Label {
				c := byte(char)
				framer.place(byte(color), c)
			}
		}
		assert.Assert(framer.currentCol <= framer.cols, "exceeded the number cols per row", "cols", framer.currentCol)

		framer.fillRemainingRow()
		framer.CurrentRow++
		framer.currentCol = 0
		framer.produceFrame()
	}

	return read, nil
}

func (a *Ansi8BitFramer) produceFrame() {
	if a.currentIdx == a.colorOffset {
		assert.Assert(a.CurrentRow == a.rows, "i should only produce frames when i have complete rows", "current row", a.CurrentRow, "rows", a.rows)
		a.count++
		out := a.buffer

		a.ch <- display.Frame{
            Empty: a.empty,
			Chars: out[:a.colorOffset],
			Color: out[a.colorOffset:],
		}

		a.buffer = make([]byte, a.rows*a.cols*2, a.rows*a.cols*2)
		a.currentIdx = 0
		a.currentCol = 0
		a.empty = 0
		a.CurrentRow = 0
	}
}

func (a *Ansi8BitFramer) DebugToFile(writer io.Writer) {
	a.debug = writer
}

func (a *Ansi8BitFramer) Frames() chan display.Frame {
	return a.ch
}
