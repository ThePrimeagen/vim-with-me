package memesweeper

import "github.com/theprimeagen/vim-with-me/pkg/window"

type grid struct {
	row int
	col int

	cells [][]window.Cell
	id    int
}

func newGrid(row, col int, width, height int) *grid {
	cells := make([][]window.Cell, 0, height+1)
	for range height + 1 {
		cell_row := make([]window.Cell, 0, width+1)
		for range width + 1 {
			cell_row = append(cell_row, window.DefaultCell(' '))
		}
		cells = append(cells, cell_row)
	}

	for w := range width {
		value := byte('A' + w)
		cells[0][w+1].Value = string(value)[0]
	}

	for h := range height {
		value := byte('1' + h)
		cells[h+1][0].Value = string(value)[0]
	}

	window.DebugCells(cells)
	return &grid{row: row, col: col, cells: cells, id: window.GetNextId()}
}

func (g *grid) Render() (window.Location, [][]window.Cell) {
	return window.NewLocation(g.row, g.col), g.cells
}

func (g *grid) Z() int {
	return 1
}

func (g *grid) Id() int {
	return g.id
}
