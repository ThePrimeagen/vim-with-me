package components

import (
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

type Text struct {
	window.RenderBase
	row, col int

	cells [][]window.Cell
}

func NewText(row, col int, txt string) *Text {
	t := &Text{
		RenderBase: window.NewRenderBase(1),
		row:        row,
		col:        col,
	}

	t.SetText(txt)
	return t
}

func NewTextZ(row, col, z int, txt string) *Text {
	t := &Text{
		RenderBase: window.NewRenderBase(z),
		row:        row,
		col:        col,
	}

	t.SetText(txt)
	return t
}

func (t *Text) SetText(txt string) {
	text_cells := make([]window.Cell, 0, len(txt))

	for _, rune := range txt {
		text_cells = append(text_cells, window.DefaultCell(byte(rune)))
	}

	t.cells = [][]window.Cell{
		text_cells,
	}
}

func (t *Text) Render() (window.Location, [][]window.Cell) {
	return window.NewLocation(t.row, t.col), t.cells
}
