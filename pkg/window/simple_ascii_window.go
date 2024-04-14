package window

/*

// TODO(v1): Simple Ascii window needs to use Cells insteasd of bytes
// TODO(v1): Simple Ascii window needs to use Cells insteasd of bytes
// TODO(v1): Simple Ascii window needs to use Cells insteasd of bytes
// TODO(v1): Simple Ascii window needs to use Cells insteasd of bytes
// TODO(v1): Simple Ascii window needs to use Cells insteasd of bytes

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
)

// TODO(v1): Simple Ascii window needs to use Cells insteasd of bytes
type SimpleAsciiWindow struct {
	Rows       byte
	Cols       byte
	cache      [][]byte
	changes    []commands.Change
}

func NewSimpleWindow(rows, cols byte) *SimpleAsciiWindow {
	cache := make([][]byte, rows)
	for i := range cache {
		cache[i] = make([]byte, cols)
		for j := range cache[i] {
			cache[i][j] = byte(' ')
		}
	}

	return &SimpleAsciiWindow{
		Rows:  rows,
		Cols:  cols,
		cache: cache,
        changes: make([]commands.Change, 0),
	}
}

func (w *SimpleAsciiWindow) Set(row, col byte, value byte) error {
	if row < 0 || row >= w.Rows {
		return fmt.Errorf("Row out of bounds: %d", row)
	}

	if col < 0 || col >= w.Cols {
		return fmt.Errorf("Col out of bounds: %d", col)
	}

	if w.cache[row][col] != value {
		w.cache[row][col] = value
		w.changes = append(w.changes, commands.Change{
			Row:   row,
			Col:   col,
			Value: value,
		})
	}

	return nil
}

func (w *SimpleAsciiWindow) SetString(row byte, value string) {
	assert.Assert(len(value) > int(w.Cols), fmt.Sprintf("String provided to Window is longer than columns: %d > %d", len(value), w.Cols))

	for i, r := range []byte(value) {
		if w.cache[row][i] != r {
			w.cache[row][i] = r
			w.changes = append(w.changes, commands.Change{
				Row:   row,
				Col:   byte(i),
				Value: r,
			})
		}
	}
}

func (w *SimpleAsciiWindow) SetWindow(value string) error {
	if len(value) != int(w.Rows)*int(w.Cols) {
		return fmt.Errorf("String provided to Window is not the correct length: %d != %d", len(value), w.Rows*w.Cols)
	}

	for i, r := range []byte(value) {
		row := i / int(w.Cols)
		col := i % int(w.Cols)

		if w.cache[row][col] != r {
			w.cache[row][col] = r
			w.changes = append(w.changes, commands.Change{
				Row:   byte(row),
				Col:   byte(col),
				Value: r,
			})
		}
	}

	return nil
}

func (r *SimpleAsciiWindow) Dimensions() (byte, byte) {
    return r.Rows, r.Cols
}

func (w *SimpleAsciiWindow) Render() string {
	out := ""
	for i := 0; i < int(w.Rows); i++ {
		out += string(w.cache[i])
	}
	w.changes = make([]commands.Change, 0)
	return out
}

func (w *SimpleAsciiWindow) PartialRender() commands.Changes {
	changes := w.changes
	w.changes = make([]commands.Change, 0)
	return commands.Changes(changes)
}
*/
