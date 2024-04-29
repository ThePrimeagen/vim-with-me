package memesweeper

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/theprimeagen/vim-with-me/pkg/window"
)

type SweeperCell struct {
	bomb     bool
	revealed bool

	adjacentCount     int
	displayedAdjCount int
}

var colors = []window.Color{
	window.NewColor(255, 0, 0, true),
	window.NewColor(0, 255, 0, true),
	window.NewColor(0, 0, 255, true),

	window.NewColor(0, 255, 255, true),
	window.NewColor(255, 0, 255, true),
	window.NewColor(255, 255, 0, true),

	window.NewColor(128, 128, 128, true),
	window.NewColor(255, 255, 255, true),
}

func toWindowCell(sCell *SweeperCell) window.Cell {
	cell := window.DefaultCell(' ')
	if !sCell.revealed {
		return cell
	}

	if sCell.bomb {
		cell.Value = '*'
		cell.Foreground = colors[0]
	} else {
		str := strconv.Itoa(sCell.displayedAdjCount)
		cell.Value = str[0]
		cell.Foreground = colors[int(math.Max(0, float64(sCell.displayedAdjCount)-1))]
	}
	return cell
}

type BoardState int

const (
	INIT BoardState = iota
	PLAYING
	LOSE
	WIN
)

type Board struct {
	params BoardParams

	revealCount, id int
	countCount      int

	state BoardState
	cells [][]*SweeperCell
}

type BoardParams struct {
	width, height int
	row, col      int
	bombCount     int
}

func NewBoard(params BoardParams) *Board {
	cells := make([][]*SweeperCell, 0)
	for range params.height {
		cell_row := make([]*SweeperCell, 0, params.height)
		for range params.width {
			cell_row = append(cell_row, &SweeperCell{
				bomb:     false,
				revealed: false,

				adjacentCount:     0,
				displayedAdjCount: 0,
			})
		}
		cells = append(cells, cell_row)
	}

	return &Board{
		params: params,

		id:          window.GetNextId(),
		countCount:  0,
		state:       INIT,
		revealCount: 0,
	}
}

func (r *Board) PickSpot(row, col int) {
	if r.state == INIT {
		r.init(row, col)
	}

	cell := r.cells[row][col]
	r.revealSpot(row, col)

	if cell.bomb {
		r.state = LOSE
		return
	}

	if r.revealCount == r.params.height*r.params.width-r.params.bombCount {
		r.state = WIN
	}
}

var dirs = [][]int{
	{-1, 0},
	{1, 0},

	{-1, 1},
	{0, 1},
	{1, 1},

	{-1, -1},
	{0, -1},
	{1, -1},
}

func (r *Board) count(row, col int) int {
	c := 0
	for _, d := range dirs {
		rowNext := row + d[0]
		colNext := col + d[1]

		if rowNext >= r.params.height || colNext >= r.params.width || rowNext < 0 || colNext < 0 {
			continue
		}

		if r.cells[rowNext][colNext].bomb {
			c++
            r.countCount++
		}
	}
	return c
}

func (r *Board) revealSpot(row, col int) {
	cell := r.cells[row][col]
	if cell.revealed {
		return
	}

	r.revealCount++
	cell.revealed = true

	if cell.adjacentCount != 0 || cell.bomb {
		return
	}

	for _, d := range dirs {
		rowNext := row + d[0]
		colNext := col + d[1]

		if rowNext >= r.params.height || colNext >= r.params.width || rowNext < 0 || colNext < 0 {
			continue
		}

		r.revealSpot(rowNext, colNext)
	}
}

func (r *Board) init(row, col int) {
	for range r.params.bombCount {
		randomC := 0
		randomR := 0
		for {

			randomR = rand.Intn(r.params.height)
			randomC = rand.Intn(r.params.width)

			if (randomC != col || randomR != row) && !r.cells[row][col].bomb {
				break
			}
		}

		r.cells[randomR][randomC].bomb = true
	}

	params := r.params
	for row := range params.height {
		for c := range params.width {
			r.cells[row][c].adjacentCount = r.count(row, c)
		}
	}

	r.state = PLAYING
}

func (r *Board) ReduceOne() {
    if r.countCount == 0 {
        return
    }

    r.countCount--

	for {

		row := rand.Intn(r.params.height)
		c := rand.Intn(r.params.width)
		cell := r.cells[row][c]

		if cell.bomb || cell.displayedAdjCount == 0 {
			continue
		}
		cell.displayedAdjCount--
	}
}

func (r *Board) Render() (window.Location, [][]window.Cell) {
    cells := make([][]window.Cell, 0)
	for row := range r.params.height {
        cell_row := make([]window.Cell, 0)
		for c := range r.params.width {
			cell_row = append(cell_row, toWindowCell(r.cells[row][c]))
		}
        cells = append(cells, cell_row)
	}

    return window.NewLocation(r.params.row, r.params.col), cells
}

func (r *Board) Z() int {
	return 2
}

func (r *Board) Id() int {
	return r.id
}
