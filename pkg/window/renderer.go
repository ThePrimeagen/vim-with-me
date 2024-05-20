package window

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

const CELL_ENCODING_LENGTH = COLOR_ENCODING_LENGTH*2 + 1
const CELL_AND_LOC_ENCODING_LENGTH = CELL_ENCODING_LENGTH + LOCATION_ENCODING_LENGTH

const CELL_VALUE_NO_PLACE = 0
const CELL_VALUE_BACKGROUND_COLOR_ONLY = 1

type Cell struct {
	Foreground Color `json:"foreground"`
	Background Color `json:"background"`
	Value      byte  `json:"value"`
}

func ForegroundCell(value byte, foreground Color) Cell {
	return Cell{
		Value:      value,
		Foreground: foreground,
		Background: DEFAULT_BACKGROUND,
	}
}

func BackgroundCellOnly(background Color) Cell {
	return Cell{
		Value:      CELL_VALUE_BACKGROUND_COLOR_ONLY,
		Foreground: DEFAULT_FOREGROUND,
		Background: background,
	}
}

func BackgroundCell(value byte, background Color) Cell {
	return Cell{
		Value:      value,
		Foreground: DEFAULT_FOREGROUND,
		Background: background,
	}
}

func DefaultCell(value byte) Cell {
	return Cell{
		Value:      value,
		Foreground: DEFAULT_FOREGROUND,
		Background: DEFAULT_BACKGROUND,
	}
}

func EmptyCell() Cell {
	return Cell{
		Value:      CELL_VALUE_NO_PLACE,
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

func (c *Cell) Merge(cell *Cell) {
	if cell.Value == CELL_VALUE_NO_PLACE {
		return
	}

	if cell.Value == CELL_VALUE_BACKGROUND_COLOR_ONLY {
		c.Background = cell.Background
		return
	}

	c.Value = cell.Value
	c.Background = cell.Background
	c.Foreground = cell.Foreground
}

func DebugCells(cells [][]Cell) {
	for _, cell_row := range cells {
		for _, cell := range cell_row {
			fmt.Printf("|%s|", string(cell.Value))
		}
		fmt.Println()
	}
}

type CellWithLocation struct {
	Cell     `json:"cell"`
	Location `json:"location"`
}

func NewCellWithLocation(cell Cell, row, col int) *CellWithLocation {
	return &CellWithLocation{
		Cell:     cell,
		Location: NewLocation(row, col),
	}
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
	assert.Assert(len(data) >= CELL_ENCODING_LENGTH, fmt.Sprintf("Cell#UnmarshalBinary not enough data to UnmarshalBinary: got %d -- expected %d", len(data), COLOR_ENCODING_LENGTH))

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

func (c *Cell) EqualWithLocation(other *CellWithLocation) bool {
	return c.Equal(&other.Cell)
}

func (c *Cell) IsEmpty() bool {
	return c.Value == CELL_VALUE_NO_PLACE
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

func NewRender(rows, cols int) *Renderer {
	length := cols * rows
	buffer := make([]Cell, 0, length)
	clean := make([]Cell, 0, length)

	for i := 0; i < int(length); i++ {
		buffer = append(buffer, Cell{
			Foreground: DEFAULT_FOREGROUND,
			Background: DEFAULT_BACKGROUND,
			Value:      byte(' '),
		})
		clean = append(clean, EmptyCell())
	}

	previous := make([]Cell, length, length)
	copy(previous, buffer)

	slog.Debug("new renderer", "rows", rows, "cols", cols)
	return &Renderer{
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

func translateFromIdx(idx, rowSize, colSize int) (bool, int, int) {
	row := idx / colSize
	if row >= rowSize {
		slog.Debug("translateFromIdx: exceeds", "idx", idx, "row", row, "rowSize", rowSize, "colSize", colSize)
		return true, 0, 0
	}

	return false, row, idx % colSize
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

func (r *Renderer) FromRemoteRenderer(cells []*CellWithLocation) {
	for _, c := range cells {
		exceeds, idx := translate(&c.Location, 0, 0, r.rows, r.cols)
		assert.Assert(exceeds == false, "you should never render from a canvas that is too big")

		r.buffer[idx].Merge(&c.Cell)
	}
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
	if cells == nil {
		return
	}

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

			cell := cells[row][col]
			if cell.Value == CELL_VALUE_NO_PLACE {
				continue
			} else if cell.Value == CELL_VALUE_BACKGROUND_COLOR_ONLY {
				r.buffer[idx].Background = cell.Background
			} else {
				r.buffer[idx] = cell
			}
		}
	}
}

func (r *Renderer) Clear() {
	r.renderables = make([]Render, 0)
}

func (r *Renderer) Render() []*CellWithLocation {
	for i := 0; i < len(r.renderables); i++ {
		slog.Debug("Render#renderable", "id", r.renderables[i].Id(), "z", r.renderables[i].Z())
		r.place(r.renderables[i])
	}

	out := make([]*CellWithLocation, 0)
	for i, cell := range r.buffer {
		other := r.previous[i]
		if !cell.Equal(&other) {
			exceeds, row, col := translateFromIdx(i, r.rows, r.cols)
			assert.Assert(exceeds == false, "exceeded bounds when partial rendering")

			// I probably care about this...
			// TODO(v1): LogValuer interface (LogAttr maybe?)
			slog.Debug("partial render cell with location", "row", row, "col", col, "cell", cell.String())

			if cell.IsEmpty() {
				continue
			}

			out = append(out, &CellWithLocation{
				Cell:     cell,
				Location: NewLocation(row, col),
			})

			r.previous[i].Merge(&cell)
		}
	}

	r.previousPartialRender = out
	copy(r.buffer, r.clean)
	return out
}

var id int = 0

func GetNextId() int {
	id++
	return id
}

func (r *Renderer) FullRender() []*CellWithLocation {
	cells := make([]*CellWithLocation, 0)
	for idx, cell := range r.previous {
		exceeds, row, col := translateFromIdx(idx, r.rows, r.cols)
		assert.Assert(exceeds == false, "somehow i have translated from our buffer and exceeded our buffer")

		cells = append(cells, NewCellWithLocation(cell, row, col))
	}
	return cells
}

func printBuff(buffer []Cell, rows, cols int) string {
	out := make([]string, 0)
	for row := 0; row < rows; row++ {
		strRow := ""
		for col := 0; col < cols; col++ {
			i := row*cols + col
			strRow += fmt.Sprintf("|%s%s", buffer[i].Background.ColorCode(), string(buffer[i].Value))
		}
		strRow += "|"
		out = append(out, strRow)
	}

	return strings.Join(out, "\n")
}

func (r *Renderer) Debug() string {
	return printBuff(r.previous, r.rows, r.cols)
}
