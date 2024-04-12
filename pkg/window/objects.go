package window

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

type Cell struct {
	Count      int
	Foreground Color
	Background Color
	Value      byte
}

type Render interface {
	Render() (Location, [][]Cell)
}

type Renderer struct {
    cols     int
    rows     int
    len      int

	buffer   []Cell
	previous []Cell

	previousPartialRender []Cell
}

func NewRender(cols, rows int) Renderer {
    length := cols * rows
    buffer := make([]Cell, 0, length)
    previous := make([]Cell, 0, length)

    for i := 0; i < int(length); i++ {
        buffer = append(buffer, Cell{
            Count: 0,
            Foreground: NewColor(255, 255, 255, true),
            Background: NewColor(255, 255, 255, false),
            Value: byte(' '),
        })
    }

    copy(previous, buffer)

    return Renderer{
        buffer: buffer,
        previous: previous,
        previousPartialRender: make([]Cell, 0),

        cols: cols,
        rows: rows,
        len: rows * cols,
    }
}

func translate(loc *Location, offsetR, offsetC, rowSize, colSize int) (bool, int) {
    out := int((loc.Row + offsetR) * colSize + loc.Col + offsetC)

    exceeds :=
        // Off the board on right or down
        loc.Row + offsetR >= rowSize || loc.Col + offsetC >= colSize ||

        // Off the board on left or top
        (loc.Row + offsetR) < 0 || (loc.Col + offsetC) < 0

        // Off the board on right or down
    return exceeds, out
}

func (r *Renderer) Place(renderable Render) {
    location, cells := renderable.Render()

    assert.Assert(len(cells) > 0, "must contain rendering data")
    assert.Assert(len(cells[0]) > 0, "must contain rendering data")

    rows := len(cells)
    cols := len(cells[0])

    for row := 0; row < rows; row++ {
        for col := 0; col < cols; col++ {
            exceeds, idx := translate(&location, row, col, r.rows, r.cols)
            fmt.Printf("exceeds: %v idx: %d\n", exceeds, idx)
            if exceeds {
                continue
            }

            count := r.buffer[idx].Count + 1
            r.buffer[idx] = cells[row][col]
            r.buffer[idx].Count = count
        }
    }
}

func (r *Renderer) Render() []Cell {
    copy(r.buffer, r.previous)
    return nil
}
