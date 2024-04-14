package window

import (
	"fmt"
	"slices"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

type Cell struct {
	Foreground Color
	Background Color
	Value      byte
}

func (c *Cell) Equal(other *Cell) bool {
    return c.Value == other.Value &&
        c.Foreground.Equal(&other.Foreground) &&
        c.Background.Equal(&other.Background)
}

type Render interface {
	Render() (Location, [][]Cell)
	// Expect Z to be a constant for now
	Z() int
	Id() int
}

type Renderer struct {
	cols int
	rows int
	len  int

	buffer      []Cell
	previous    []Cell
	renderables []Render

	previousPartialRender []*Cell
}

func NewRender(cols, rows int) Renderer {
	length := cols * rows
	buffer := make([]Cell, 0, length)

	for i := 0; i < int(length); i++ {
		buffer = append(buffer, Cell{
			Foreground: NewColor(255, 255, 255, true),
			Background: NewColor(255, 255, 255, false),
			Value:      byte(' '),
		})
	}

	previous := make([]Cell, length, length)
    copy(previous, buffer)

	return Renderer{
		buffer:                buffer,
		previous:              previous,
		renderables:           make([]Render, 0, 100),
		previousPartialRender: make([]*Cell, 0),

		cols: cols,
		rows: rows,
		len:  rows * cols,
	}
}

func translate(loc *Location, offsetR, offsetC, rowSize, colSize int) (bool, int) {
	out := int((loc.Row+offsetR)*colSize + loc.Col + offsetC)

	exceeds :=
		// Off the board on right or down
		loc.Row+offsetR >= rowSize || loc.Col+offsetC >= colSize ||

			// Off the board on left or top
			(loc.Row+offsetR) < 0 || (loc.Col+offsetC) < 0

		// Off the board on right or down
	return exceeds, out
}

func (r *Renderer) Add(renderable Render) {
	length := len(r.renderables)

	Z := renderable.Z()
	lo := 0
	hi := length
	idx := 0

	for lo < hi {
		mid := lo + (hi-lo)/2
		v := r.renderables[mid].Z()

		idx = mid
		if v == Z {
			idx = mid
			break
		} else if v < Z {
			idx = mid + 1
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	r.renderables = slices.Insert(r.renderables, idx, renderable)
}

func (r *Renderer) Remove(renderable Render) {
	// I can make this faster with a map... profile later on
	for i, v := range r.renderables {
		if v.Id() == renderable.Id() {
			r.renderables = slices.Delete(r.renderables, i, i+1)
			break
		}
	}
}

func (r *Renderer) place(renderable Render) {
	location, cells := renderable.Render()

	assert.Assert(len(cells) > 0, "must contain rendering data")
	assert.Assert(len(cells[0]) > 0, "must contain rendering data")

	rows := len(cells)
	cols := len(cells[0])

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			exceeds, idx := translate(&location, row, col, r.rows, r.cols)
			if exceeds {
				continue
			}

			r.buffer[idx] = cells[row][col]
		}
	}
}

func (r *Renderer) Render() []*Cell {
    for i := 0; i < len(r.renderables); i++ {
		r.place(r.renderables[i])
    }

    out := make([]*Cell, 0)
	for i, cell := range r.buffer {
        other := r.previous[i]
        if !cell.Equal(&other) {
            out = append(out, &cell)
        }
	}

    r.previousPartialRender = out
    return out
}

func printBuff(buffer []Cell, rows, cols int) {
    for row := 0; row < rows; row++ {
        toPrint := make([]int, 0)
        for col := 0; col < cols; col++ {
            i := row * cols + col
            toPrint = append(toPrint, int(buffer[i].Value))
        }
        fmt.Printf("%+v\n", toPrint)
    }

}

func (r *Renderer) debug() {
    fmt.Println("buffer")
    printBuff(r.buffer, r.rows, r.cols)

    fmt.Println("previous")
    printBuff(r.previous, r.rows, r.cols)
}
