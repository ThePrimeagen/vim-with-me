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
	"github.com/theprimeagen/vim-with-me/pkg/v2/rgb"
)

type AnsiFramer struct {
	State *AnsiFramerState

	frameStart []byte
	inputStart []byte

	debug   io.Writer
	ch      chan display.Frame
	buffer  []byte
	scratch []byte
	writer  *rgb.RGBWriter

	lastStyle *ansi.StyledText
}

type AnsiFramerState struct {
	Rows   int
	Cols   int
	Length int

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

func (s *AnsiFramerState) String() string {
	return fmt.Sprintf(`AnsiFramerState(%d, %d): count=%d readCount=%d
currentRow=%d currentCol=%d currentIdx=%d
attachedScratch=%v
currentStyledLine=%s
currentInputLine=%s
currentInputByteLine=%+v
currentLine=%s
currentStyle(%d/%d)=%s`, s.Rows, s.Cols,
		s.Count, s.ReadCount,
		s.CurrentRow, s.CurrentCol, s.CurrentIdx,
		s.AttachedScratch,
		s.CurrentStyledLine,
		s.CurrentInputLine,
		s.CurrentInputByteLine,
		string(s.CurrentLine),
		s.CurrentStyleIdx, s.CurrentStyleCount, s.CurrentStyle)
}

func (s *AnsiFramerState) Reset() {
	s.CurrentIdx = 0
	s.CurrentCol = 0
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
func NewFramer() *AnsiFramer {
	state := AnsiFramerState{
		Rows:   0,
		Cols:   0,
		Length: 0,

		Count:     0,
		ReadCount: 0,

		CurrentRow: 0,
		CurrentCol: 0,
		CurrentIdx: 0,

		CurrentLine: make([]byte, 0),
	}

	assert.AddAssertData("AnsiFramer", &state)

	// 1 byte color, 1 byte ascii
	return &AnsiFramer{
		State: &state,

		ch:         make(chan display.Frame, 10),
		buffer:     make([]byte, 0, 0),
		scratch:    make([]byte, 0),
		frameStart: nil,
		writer:     rgb.New8BitRGBWriter(),
	}
}

func (a *AnsiFramer) WithInputStart(start []byte) *AnsiFramer {
	a.inputStart = start

	return a
}

func (a *AnsiFramer) WithFrameStart(start []byte) *AnsiFramer {
	a.frameStart = start

	return a
}

func (a *AnsiFramer) WithDim(rows, cols int) *AnsiFramer {
	a.State.Rows = rows
	a.State.Cols = cols
	a.State.Length = rows * cols

	a.reset()

	return a
}

func (a *AnsiFramer) WithColorWriter(writer *rgb.RGBWriter) {
	a.writer = writer
	a.reset()
}

func (a *AnsiFramer) reset() {
	a.State.Reset()

	colorLength := a.writer.ByteLength() * a.State.Length
	length := a.State.Length + colorLength
	a.buffer = make([]byte, length, length)
	a.writer.Set(a.buffer[a.State.Length:])
}

func (framer *AnsiFramer) place(color *ansi.Rgb, char byte) {
	assert.Assert(framer.State.CurrentCol != framer.State.Cols, "current cols equals the maximum state cols")
	assert.Assert(framer.State.CurrentIdx < framer.State.Length, "place failed", "color", color, "byte", char, "State.CurrentIdx", framer.State.CurrentIdx, "data length", framer.State.Length)

	framer.State.CurrentLine[framer.State.CurrentCol] = char
	framer.buffer[framer.State.CurrentIdx] = char
	framer.writer.Write(color)

	framer.State.CurrentIdx++
	framer.State.CurrentCol++

}

var black = ansi.Rgb{R: 0, G: 0, B: 0}

func (framer *AnsiFramer) fillRemainingRow() {
	for framer.State.CurrentCol < framer.State.Cols {
		framer.place(&black, ' ')
	}
}

var newline = []byte{'\r', '\n'}

func (framer *AnsiFramer) Write(data []byte) (int, error) {
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

			fg := style.FgCol
			if fg == nil {
				fg = framer.lastStyle.FgCol
			}

			for _, char := range style.Label {
				assert.Assert(framer.State.CurrentCol <= framer.State.Cols, "exceeded the number cols per row", "cols", framer.State.CurrentCol)
				c := byte(char)
				framer.place(&fg.Rgb, c)
			}

			if style.FgCol != nil {
				framer.lastStyle = style
			}
		}

		framer.fillRemainingRow()
		framer.State.CurrentRow++
		if framer.frameStart == nil && framer.State.CurrentRow == framer.State.Rows {
			framer.produceFrame()
		}
		framer.State.CurrentCol = 0
	}

	return read, nil
}

func (a *AnsiFramer) produceFrame() {
	assert.Assert(a.State.CurrentRow == a.State.Rows, "must produce a correct amount of rows", "rows", a.State.CurrentRow)
	assert.Assert(a.State.CurrentIdx == a.State.Length, "current idx != state.length")

	a.ch <- display.Frame{
		Idx:   0,
		Chars: a.buffer[:a.State.Length],
		Color: a.buffer[a.State.Length:],
	}

	a.reset()
}

func (a *AnsiFramer) DebugToFile(writer io.Writer) {
	a.debug = writer
}

func (a *AnsiFramer) Frames() chan display.Frame {
	return a.ch
}

func RemoveAsciiStyledPixels(data []byte) []byte {
	assert.Assert(len(data)&1 == 0, "you cannot remove ascii styled pixels if the array is not even length")

	idx := 1
	doubleIdx := 2

	for ; doubleIdx < len(data); doubleIdx += 2 {
		data[idx] = data[doubleIdx]
		idx++
	}

	return data[:len(data)/2]
}
