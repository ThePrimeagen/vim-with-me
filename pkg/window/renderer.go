package window

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

const CELL_ENCODING_LENGTH = COLOR_ENCODING_LENGTH*2 + 1
const CELL_AND_LOC_ENCODING_LENGTH = CELL_ENCODING_LENGTH + LOCATION_ENCODING_LENGTH

type Cell struct {
	Foreground Color
	Background Color
	Value      byte
}

func ForegroundCell(value byte, foreground Color) Cell {
    return Cell{
        Value: value,
        Foreground: foreground,
        Background: DEFAULT_BACKGROUND,
    }
}

func DefaultCell(value byte) Cell {
    return Cell{
        Value: value,
        Foreground: DEFAULT_FOREGROUND,
        Background: DEFAULT_BACKGROUND,
    }
}

func (c *Cell) String() string {
	return fmt.Sprintf(
		"value=%s foreground=%s background=%s",
		[]byte{c.Value},
		c.Foreground.String(),
		c.Background.String(),
	)
}

type CellWithLocation struct {
	Cell
	Location
}

func (c *CellWithLocation) MarshalBinary() ([]byte, error) {
	loc, err := c.Location.MarshalBinary()
	if err != nil {
		return nil, err
	}

	cell, err := c.Cell.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return append(loc, cell...), nil
}

func (c *CellWithLocation) UnmarshalBinary(data []byte) error {
	assert.Assert(len(data) >= CELL_ENCODING_LENGTH+LOCATION_ENCODING_LENGTH, "not enough bytes for unmarshaling CellWithLocation")

	var loc Location
	err := loc.UnmarshalBinary(data)
	if err != nil {
		return err
	}
	c.Location = loc

	var cell Cell
	err = cell.UnmarshalBinary(data[LOCATION_ENCODING_LENGTH:])
	if err != nil {
		return err
	}
	c.Cell = cell

	return nil
}

func (c *Cell) MarshalBinary() ([]byte, error) {
	foreground, err := c.Foreground.MarshalBinary()
	if err != nil {
		return []byte{}, err
	}

	background, err := c.Background.MarshalBinary()
	if err != nil {
		return []byte{}, err
	}

	b := make([]byte, 0, len(foreground)+len(background)+1)
	b = append(b, c.Value)
	b = append(b, foreground...)
	return append(b, background...), nil
}

func (c *Cell) UnmarshalBinary(data []byte) error {
	assert.Assert(len(data) < 1+COLOR_ENCODING_LENGTH*2, "i should never unmarshall without all the data")

	c.Value = data[0]
	var foreground Color
	err := foreground.UnmarshalBinary(data[1:])
	if err != nil {
		return err
	}
	c.Foreground = foreground

	var background Color
	err = background.UnmarshalBinary(data[1+COLOR_ENCODING_LENGTH:])
	if err != nil {
		return err
	}

	c.Background = background
	return nil
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
	clean       []Cell
	renderables []Render

	previousPartialRender []*CellWithLocation
}

func NewRender(rows, cols int) Renderer {
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
	clean := make([]Cell, length, length)
	copy(previous, buffer)
	copy(clean, buffer)

	slog.Debug("new renderer", "rows", rows, "cols", cols)
	return Renderer{
		buffer:                buffer,
		previous:              previous,
		clean:                 clean,
		renderables:           make([]Render, 0, 100),
		previousPartialRender: make([]*CellWithLocation, 0),

		cols: cols,
		rows: rows,
		len:  rows * cols,
	}
}

func translate(loc *Location, offsetR, offsetC, rowSize, colSize int) (bool, int) {
	out := (loc.Row+offsetR)*colSize + loc.Col + offsetC

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

func (r *Renderer) Dimensions() (byte, byte) {
	return byte(r.rows), byte(r.cols)
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

func (r *Renderer) Render() []*CellWithLocation {
	for i := 0; i < len(r.renderables); i++ {
		r.place(r.renderables[i])
	}

	out := make([]*CellWithLocation, 0)
	for i, cell := range r.buffer {
		other := r.previous[i]
		if !cell.Equal(&other) {
			row := i / r.cols
			col := i % r.cols

			// I probably care about this...
			// TODO(v1): LogValuer interface (LogAttr maybe?)
			slog.Debug("partial render cell with location", "row", row, "col", col, "cell", cell.String())

			out = append(out, &CellWithLocation{
				Cell:     cell,
				Location: NewLocation(row, col),
			})
		}
	}

	r.previousPartialRender = out
	copy(r.previous, r.buffer)
	copy(r.buffer, r.clean)
	return out
}

var id int = 0
func GetNextId() int {
    id++
    return id
}

func (r *Renderer) FullRender() []*Cell {
	assert.Assert(false, "please implement me")
	return nil
}

func printBuff(buffer []Cell, rows, cols int) {
	for row := 0; row < rows; row++ {
		toPrint := make([]int, 0)
		for col := 0; col < cols; col++ {
			i := row*cols + col
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
