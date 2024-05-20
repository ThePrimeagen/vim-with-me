package window

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

// TODO(v1): Simple Ascii window needs to use Cells insteasd of bytes
type SimpleAsciiWindow struct {
	Rows    int
	Cols    int
	cache   [][]Cell
	changes []*CellWithLocation
}

func NewSimpleWindow(rows, cols int) SimpleAsciiWindow {
	cache := make([][]Cell, rows)
	for i := range cache {
		cache[i] = make([]Cell, cols)
		for j := range cache[i] {
			cache[i][j] = DefaultCell(' ')
		}
	}

	return SimpleAsciiWindow{
		Rows:    rows,
		Cols:    cols,
		cache:   cache,
		changes: make([]*CellWithLocation, 0),
	}
}

func (w *SimpleAsciiWindow) Set(row, col int, value byte) error {
	if row < 0 || row >= w.Rows {
		return fmt.Errorf("Row out of bounds: %d", row)
	}

	if col < 0 || col >= w.Cols {
		return fmt.Errorf("Col out of bounds: %d", col)
	}

	if w.cache[row][col].Value != value {
		w.cache[row][col].Value = value
		w.changes = append(w.changes, &CellWithLocation{
			Cell:     w.cache[row][col],
			Location: NewLocation(row, col),
		})
	}

	return nil
}

func (w *SimpleAsciiWindow) SetString(row int, value string) {
	assert.Assert(len(value) > int(w.Cols), fmt.Sprintf("String provided to Window is longer than columns: %d > %d", len(value), w.Cols))

	for i, v := range []byte(value) {
		if w.cache[row][i].Value != v {
			w.cache[row][i].Value = v
			w.changes = append(w.changes, &CellWithLocation{
				Cell:     w.cache[row][i],
				Location: NewLocation(row, i),
			})
		}
	}
}

func (w *SimpleAsciiWindow) SetWindow(value string) error {
	if len(value) != int(w.Rows)*int(w.Cols) {
		return fmt.Errorf("String provided to Window is not the correct length: %d != %d", len(value), w.Rows*w.Cols)
	}

	for i, v := range []byte(value) {
		row := i / int(w.Cols)
		col := i % int(w.Cols)

		if w.cache[row][col].Value != v {
			w.cache[row][col].Value = v
			w.changes = append(w.changes, &CellWithLocation{
				Cell:     w.cache[row][col],
				Location: NewLocation(row, col),
			})
		}
	}

	return nil
}

func (r *SimpleAsciiWindow) Dimensions() (byte, byte) {
	return byte(r.Rows), byte(r.Cols)
}

func (w *SimpleAsciiWindow) Render() []*CellWithLocation {
	w.changes = make([]*CellWithLocation, 0)

	out := make([]*CellWithLocation, 0)

	for r, row := range w.cache {
		for c, cell := range row {
			out = append(out, &CellWithLocation{
				Cell:     cell,
				Location: NewLocation(r, c),
			})
		}
	}

	return out
}

func (w *SimpleAsciiWindow) PartialRender() []*CellWithLocation {
	changes := w.changes
	w.changes = make([]*CellWithLocation, 0)

	return changes
}
