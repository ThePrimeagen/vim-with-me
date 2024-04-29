package memesweeper

import "github.com/theprimeagen/vim-with-me/pkg/window"

type Text struct {
	row, col, id int

	cells [][]window.Cell
}

func NewText(row, col int, txt string) *Text {
    t := &Text{
        row: row,
        col: col,
        id: window.GetNextId(),
    }

    t.SetText(txt)
    return t
}

func (t *Text) SetText(txt string) {
    text_cells := make([]window.Cell, 0, len(txt))

    for i, rune := range txt {
        text_cells = append(text_cells, window.ForegroundCell(byte(rune), colors[i % len(colors)]))
    }

    t.cells = [][]window.Cell{
        text_cells,
    }
}

func (t *Text) Render() (window.Location, [][]window.Cell) {
    return window.NewLocation(t.row, t.col), t.cells
}

func (t *Text) Z() int {
    return 1
}

func (t *Text) Id() int {
    return t.id
}
