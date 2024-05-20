package tower_defense

import (
	"strconv"

	"github.com/theprimeagen/vim-with-me/pkg/window"
)

type Tower struct {
	Row int
	Col int

	Count int

	id int
}

func NewTower(row, col int) *Tower {
	return &Tower{
		Row:   row,
		Col:   col,
		Count: 1,

		id: window.GetNextId(),
	}
}

func (t *Tower) Render() (window.Location, [][]window.Cell) {
	return window.NewLocation(t.Row, t.Col), [][]window.Cell{
		{window.DefaultCell(byte(strconv.Itoa(t.Count)[0]))},
	}
}

func (t *Tower) Z() int {
	return 1
}

func (t *Tower) Id() int {
	return t.id
}
