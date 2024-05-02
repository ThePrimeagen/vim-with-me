package ansiparser

import (
	//"github.com/leaanthony/go-ansi-parser"
	"bytes"
	"strings"

	"github.com/leaanthony/go-ansi-parser"
	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

type Ansi8BitFramer struct {
	rows int
	cols int

	ch         chan []byte
	currentIdx int
	currentCol int
	currentRow int
	buffer     []byte
	scratch    []byte
}

func nextAnsiChunk(data []byte) (bool, int, *ansi.StyledText, error) {
    assert.Assert(data[0] == '', "the ansi chunks should always start on an escape")
    nextEsc := bytes.Index(data[1:], escape) + 1

    var styles []*ansi.StyledText = nil
    var err error = nil
    var out int = 0
    var complete = true
    if nextEsc == 0 {
        styles, err = ansi.Parse(string(data))
        out = len(data)
        complete = false
    } else {
        out = nextEsc
        styles, err = ansi.Parse(string(data[:nextEsc]))
    }

    if styles != nil && len(styles) != 0 {
        assert.Assert(len(styles) == 1, "there must only be one style at a time parsed")
        return complete, out, styles[0], err
    }
    return complete, out, nil, err
}

// TODO: I could also use a ctx to close out everything
func New8BitFramer(rows, cols int) *Ansi8BitFramer {

	// 1 byte color, 1 byte ascii
	return &Ansi8BitFramer{
		rows:       rows,
		cols:       cols,
		ch:         make(chan []byte, 10),
		currentIdx: 0,
		currentCol: 0,
		currentRow: 0, // makes life easier
		buffer:     make([]byte, rows*cols*2, rows*cols*2),
		scratch:    make([]byte, 0),
	}
}

func RGBTo8BitColor(hex ansi.Rgb) uint {
	red := uint(hex.R*8) / 256
	green := uint(hex.G*8) / 256
	blue := uint(hex.B*4) / 256

	return (red << 5) | (green << 2) | blue
}

func remainingIsRegisteredNurse(data []byte) bool {
    if len(data) != 3 {
        return false
    }

    return data[1] == '\r' && data[2] == '\n'
}

func (framer *Ansi8BitFramer) place(color, char byte) {
    framer.buffer[framer.currentIdx] = color
    framer.buffer[framer.currentIdx+1] = char
    framer.currentIdx += 2
    framer.currentCol++
}

func (framer *Ansi8BitFramer) fillRemainingRow() {
    for framer.currentCol < framer.cols {
        framer.place(0, ' ')
    }
}

func (framer *Ansi8BitFramer) Write(data []byte) (int, error) {
	idx := 0
    if len(framer.scratch) != 0 {
        // this is terrible for perf
        data = append(framer.scratch, data...)
        framer.scratch = make([]byte, 0)
    }

	for idx < len(data) {
        completed, nextEsc, style, err := nextAnsiChunk(data[idx:])

        if !completed && framer.currentRow + 1 != framer.rows {
			framer.scratch = data[idx:]
            break
        }

        idx += nextEsc

        // errors happen when parsing non color commands
        // or there is just nothing that had any data when parsing
		if err != nil || style == nil {
			continue
		}

        color := RGBTo8BitColor(style.FgCol.Rgb)
        label := style.Label

        if strings.Contains(label, "\r\n") {
            strings.Replace(label, "\r\n", "\n", 1)
        }

        for _, char := range label {
            framer.produceFrame()

            c := byte(char)
            if c == '\n' {
                framer.fillRemainingRow()
                framer.currentCol = 0
                framer.currentRow++
                continue
            }

            if framer.currentCol >= framer.cols {
                continue
            }

            framer.place(byte(color), c)
        }
        framer.produceFrame()
	}

	return len(data), nil
}

func (a *Ansi8BitFramer) produceFrame() {
	if a.currentIdx == len(a.buffer) {
		out := a.buffer
		a.buffer = make([]byte, a.rows*a.cols*2, a.rows*a.cols*2)

		a.ch <- out
		a.currentIdx = 0
        a.currentCol = 0
        a.currentRow = 0
	}
}

func (a *Ansi8BitFramer) Frames() chan []byte {
	return a.ch
}

/*
func (a* AnsiFramer) Write(p []byte) (n int, err error) {
}
*/
