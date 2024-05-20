package window

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestRenderable struct {
	cells [][]Cell
	loc   Location
	id    int
	z     int
}

func (t *TestRenderable) Z() int {
	return t.z
}

func (t *TestRenderable) Id() int {
	return t.id
}

var noCell = DefaultCell(0)
var backgroundColor = NewColor(255, 0, 0, false)
var background = BackgroundCell(CELL_VALUE_BACKGROUND_COLOR_ONLY, backgroundColor)
var cell69 = Cell{Value: 69}
var cell70 = Cell{Value: 70}
var cell71 = Cell{Value: 71}
var cell72 = Cell{Value: 72}
var innerCells = [][]Cell{
	{cell69, cell70},
	{cell71, cell72},
}

var hiddenValueCells = [][]Cell{
	{noCell, noCell},
	{cell69, background},
}

var cell73 = Cell{Value: 73}
var cell74 = Cell{Value: 74}
var cell75 = Cell{Value: 75}
var cell76 = Cell{Value: 76}

var outerCells = [][]Cell{
	{cell73, cell74},
	{cell75, cell76},
}

var allCells = []Cell{
	cell69, cell70,
	cell71, cell72,
	cell73, cell74,
	cell75, cell76,
}

func newTestRenderable(cells [][]Cell, loc Location, z int) TestRenderable {
	id++
	return TestRenderable{
		loc:   loc,
		cells: cells,
		id:    id,
		z:     z,
	}
}

func (t *TestRenderable) Render() (Location, [][]Cell) {
	return t.loc, t.cells
}

func has(cell Cell, cells []*CellWithLocation) bool {
	for _, c := range cells {
		if cell.EqualWithLocation(c) {
			return true
		}
	}

	return false
}

func TestRender(t *testing.T) {
	render := NewRender(5, 5)
	renderers := []TestRenderable{
		newTestRenderable(outerCells, NewLocation(-1, -1), 1),
		newTestRenderable(outerCells, NewLocation(4, 4), 1),
		newTestRenderable(outerCells, NewLocation(-1, 4), 1),
		newTestRenderable(outerCells, NewLocation(4, -1), 1),
		newTestRenderable(innerCells, NewLocation(1, 1), 1),
	}

	for _, loc := range renderers {
		render.Add(&loc)
	}

	values := []byte{
		76, byte(' '), byte(' '), byte(' '), 75,
		byte(' '), 69, 70, byte(' '), byte(' '),
		byte(' '), 71, 72, byte(' '), byte(' '),
		byte(' '), byte(' '), byte(' '), byte(' '), byte(' '),
		74, byte(' '), byte(' '), byte(' '), 73,
	}

	cells := render.Render()
	render.Debug()

	for i, value := range values {
		require.Equal(t, render.previous[i].Value, value)
	}

	require.Equal(t, len(cells), 8)

	for _, cell := range allCells {
		require.True(t, has(cell, cells))
	}

	cells = render.Render()
	require.Equal(t, len(cells), 0)
}

func TestRenderHiddenValues(t *testing.T) {
	render := NewRender(5, 5)
	renderers := []TestRenderable{
		newTestRenderable(innerCells, NewLocation(1, 1), 1),
		newTestRenderable(hiddenValueCells, NewLocation(1, 1), 2),
	}

	for _, loc := range renderers {
		render.Add(&loc)
	}

	values := []byte{
		byte(' '), byte(' '), byte(' '), byte(' '), byte(' '),
		byte(' '), 69, 70, byte(' '), byte(' '),
		byte(' '), 69, 72, byte(' '), byte(' '),
		byte(' '), byte(' '), byte(' '), byte(' '), byte(' '),
		byte(' '), byte(' '), byte(' '), byte(' '), byte(' '),
	}

	backgrounds := []Color{
		DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND,
		DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND,
		DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, backgroundColor, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND,
		DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND,
		DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND, DEFAULT_BACKGROUND,
	}

	cells := render.Render()
	render.Debug()

	for i, value := range values {
		require.Equal(t, render.previous[i].Value, value)
		require.Equal(t, backgrounds[i], render.previous[i].Background)
	}

	require.Equal(t, len(cells), 4)

	cells = render.Render()
	require.Equal(t, len(cells), 0)
}
