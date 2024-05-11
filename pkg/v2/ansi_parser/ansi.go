package ansiparser

import (
	//"github.com/leaanthony/go-ansi-parser"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/leaanthony/go-ansi-parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser/display"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/encoding"
)

type Ansi8BitFramer struct {
	State *Ansi8BitFramerState

	frameStart []byte

	debug       io.Writer
	ch          chan display.Frame
	buffer      []byte
	colorOffset int
	scratch     []byte

	lastStyle *ansi.StyledText
}

type Ansi8BitFramerState struct {
	Rows   int
	Cols   int
	Length int
	Empty  int

	Count     int
	ReadCount int

	AttachedScratch bool

	CurrentStyle      *ansi.StyledText
	CurrentStyleCount int
	CurrentStyleIdx   int
	CurrentStyledLine []*ansi.StyledText

	CurrentRow           int
	CurrentCol           int
	CurrentIdx           int
	CurrentLine          []byte
	CurrentInputLine     string
	CurrentInputByteLine []byte
}

func (s *Ansi8BitFramerState) String() string {
	return fmt.Sprintf(`Ansi8BitFramerState(%d, %d): empty=%d count=%d readCount=%d
currentRow=%d currentCol=%d currentIdx=%d
attachedScratch=%v
currentStyledLine=%s
currentInputLine=%s
currentInputByteLine=%+v
currentLine=%s
currentStyle(%d/%d)=%s`, s.Rows, s.Cols,
		s.Empty, s.Count, s.ReadCount,
		s.CurrentRow, s.CurrentCol, s.CurrentIdx,
        s.AttachedScratch,
		s.CurrentStyledLine,
		s.CurrentInputLine,
		s.CurrentInputByteLine,
		string(s.CurrentLine),
		s.CurrentStyleIdx, s.CurrentStyleCount, s.CurrentStyle)
}

func (s *Ansi8BitFramerState) Reset() {
	s.CurrentIdx = 0
	s.CurrentCol = 0
	s.Empty = 0
	s.CurrentRow = 0
    s.AttachedScratch = false
	s.Count++
	s.CurrentInputLine = ""
	s.CurrentLine = make([]byte, s.Cols, s.Cols)
	s.CurrentStyleCount = 0
	s.CurrentStyleIdx = 0
}

// TODO: 2 errors, row 21 seems to start without an escape
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
	state := Ansi8BitFramerState{
		Rows:   0,
		Cols:   0,
		Length: 0,
		Empty:  0,

		Count:     0,
		ReadCount: 0,

		CurrentRow: 0,
		CurrentCol: 0,
		CurrentIdx: 0,

		CurrentLine: make([]byte, 0),
	}

	assert.AddAssertData("Ansi8BitFramer", &state)

	// 1 byte color, 1 byte ascii
	return &Ansi8BitFramer{
		State: &state,

		ch:         make(chan display.Frame, 10),
		buffer:     make([]byte, 0, 0),
		scratch:    make([]byte, 0),
		frameStart: nil,
	}
}

func (a *Ansi8BitFramer) WithFrameStart(start []byte) *Ansi8BitFramer {
	a.frameStart = start

	return a
}

func (a *Ansi8BitFramer) WithDim(rows, cols int) *Ansi8BitFramer {
	length := rows * cols

	a.State.Rows = rows
	a.State.Cols = cols
	a.State.Length = rows * cols
	a.State.Reset()

	a.colorOffset = length
	a.buffer = make([]byte, length*2, length*2)

	return a
}

func (framer *Ansi8BitFramer) place(color, char byte) {
	assert.Assert(framer.State.CurrentCol != framer.State.Cols, "current cols equals the maximum state cols")
	assert.Assert(framer.colorOffset+framer.State.CurrentIdx < len(framer.buffer), "place failed", "color", color, "byte", char, "State.CurrentIdx", framer.State.CurrentIdx, "data length", len(framer.buffer))

	framer.State.CurrentLine[framer.State.CurrentCol] = char
	framer.buffer[framer.State.CurrentIdx] = char
	framer.buffer[framer.colorOffset+framer.State.CurrentIdx] = color

	framer.State.CurrentIdx++
	framer.State.CurrentCol++

}

func (framer *Ansi8BitFramer) fillRemainingRow() {
	framer.State.Empty += framer.State.Cols - framer.State.CurrentCol
	for framer.State.CurrentCol < framer.State.Cols {
		framer.place(0, ' ')
	}
}

var newline = []byte{'\r', '\n'}

func (framer *Ansi8BitFramer) Write(data []byte) (int, error) {
	read := len(data)
	if framer.debug != nil {
		framer.debug.Write(data)
	}

	scratchLen := len(framer.scratch)

	if scratchLen != 0 {

        framer.State.AttachedScratch = true

		// this is terrible for perf
		data = append(framer.scratch, data...)
		framer.scratch = make([]byte, 0)
	} else {
        framer.State.AttachedScratch = false
    }

	for len(data) > 0 {
		nextLine := bytes.Index(data, newline)
		if nextLine == -1 {
            framer.scratch = make([]byte, len(data), len(data))
            copy(framer.scratch, data)
			break
		}

		line := data[:nextLine]
		data = data[nextLine+2:]

		// If we have a "frameStart" sequence we are looking for (ansi codes)
		// then we have to look for it on every new line
		// TODO: I may have to look into this for perf
		if framer.frameStart != nil &&
			bytes.Contains(line, framer.frameStart) &&
			framer.State.CurrentRow > 0 {

			framer.produceFrame()
		}

		lineString := string(line)

		styles := parseAnsiRow(lineString)
		assert.Assert(styles != nil, "i should never have a nil row")
		assert.Assert(framer.State.CurrentCol == 0, "we should always start every line with current col as 0")

		framer.State.CurrentStyleCount = len(styles)
		framer.State.CurrentInputLine = lineString
		framer.State.CurrentInputByteLine = line
		framer.State.CurrentStyledLine = styles

		for i, style := range styles {
			framer.State.CurrentStyleIdx = i
			framer.State.CurrentStyle = style

			color := uint(255)
			fg := style.FgCol
			if fg == nil {
				fg = framer.lastStyle.FgCol
			}

			color = encoding.RGBTo8BitColor(fg.Rgb)

			for _, char := range style.Label {
				assert.Assert(framer.State.CurrentCol <= framer.State.Cols, "exceeded the number cols per row", "cols", framer.State.CurrentCol)
				c := byte(char)
				framer.place(byte(color), c)
			}

			if style.FgCol != nil {
				framer.lastStyle = style
			}
		}

		framer.State.CurrentRow++
		framer.State.CurrentCol = 0
		if framer.frameStart == nil && framer.State.CurrentRow == framer.State.Rows {
			framer.produceFrame()
		}
	}

	return read, nil
}

func (a *Ansi8BitFramer) produceFrame() {
	assert.Assert(a.State.CurrentRow == a.State.Rows, "must produce a correct amount of rows", "rows", a.State.CurrentRow)
	assert.Assert(a.State.CurrentIdx == a.State.Length, "current idx != state.length")

	out := a.buffer

	a.ch <- display.Frame{
		Empty: a.State.Empty,
		Chars: out[:a.colorOffset],
		Color: out[a.colorOffset:],
	}

	a.buffer = make([]byte, a.State.Length*2, a.State.Length*2)
	a.State.Reset()
}

func (a *Ansi8BitFramer) DebugToFile(writer io.Writer) {
	a.debug = writer
}

func (a *Ansi8BitFramer) Frames() chan display.Frame {
	return a.ch
}
