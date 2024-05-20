package memesweeper

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/window"
	"github.com/theprimeagen/vim-with-me/pkg/window/components"
)

type SweeperCell struct {
	bomb     bool
	revealed bool
	killed   bool

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

func toWindowCell(sCell *SweeperCell, row, col int) window.Cell {
	cell := window.BackgroundCell(' ', components.BACKGROUND_GRAY)
	if (row+col)%2 == 0 {
		cell.Background = window.DEFAULT_BACKGROUND
	}
	if !sCell.revealed {
		return cell
	}

	if sCell.bomb && sCell.killed {
		cell.Value = 'x'
		cell.Foreground = colors[0]
	} else if sCell.bomb {
		cell.Value = '*'
		cell.Foreground = colors[0]
	} else if sCell.displayedAdjCount == 0 {
		cell.Value = '.'
		cell.Foreground = window.DEFAULT_FOREGROUND
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
	random    *rand.Rand
	cols      int
	rows      int
	row, col  int
	bombCount int
}

func newSweeperCells(rows, cols int) [][]*SweeperCell {
	cells := make([][]*SweeperCell, 0)
	for range rows {
		cell_row := make([]*SweeperCell, 0, rows)
		for range cols {
			cell_row = append(cell_row, &SweeperCell{
				bomb:     false,
				revealed: false,
				killed:   false,

				adjacentCount:     0,
				displayedAdjCount: 0,
			})
		}
		cells = append(cells, cell_row)
	}

	return cells
}

func NewBoard(params BoardParams) *Board {
	if params.random == nil {
		seed := int64(time.Now().UnixMilli())
		random := rand.New(rand.NewSource(seed))
		params.random = random
	}

	return &Board{
		params: params,
		cells:  newSweeperCells(params.rows, params.cols),

		id:          window.GetNextId(),
		countCount:  0,
		state:       INIT,
		revealCount: 0,
	}
}

func (r *Board) RevealBombs() {
	for _, row := range r.cells {
		for _, cell := range row {
			if cell.bomb {
				cell.revealed = true
			}
		}
	}
}

func (r *Board) Reset() {
	r.countCount = 0
	r.state = INIT
	r.revealCount = 0
	r.cells = newSweeperCells(r.params.rows, r.params.cols)
}

func (r *Board) PickSpot(row, col int) {
	if r.state == INIT {
		r.init(row, col)
	}

	cell := r.cells[row][col]
	r.revealSpot(row, col)

	if cell.bomb {
		r.state = LOSE
		cell.killed = true
		return
	}

	if r.revealCount == r.params.rows*r.params.cols-r.params.bombCount {
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

		if rowNext >= r.params.rows || colNext >= r.params.cols || rowNext < 0 || colNext < 0 {
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

		if rowNext >= r.params.rows || colNext >= r.params.cols || rowNext < 0 || colNext < 0 {
			continue
		}

		r.revealSpot(rowNext, colNext)
	}
}

func (r *Board) debug() {
	for _, row := range r.cells {
		for _, cell := range row {
			if cell.bomb {
				fmt.Printf("%v* ", cell.revealed)
			} else {
				fmt.Printf("%v%d ", cell.revealed, cell.displayedAdjCount)
			}
		}
		fmt.Println()
	}
}

func (r *Board) init(row, col int) {
	for range r.params.bombCount {
		randomC := 0
		randomR := 0
		for {

			randomR = r.params.random.Intn(r.params.rows)
			randomC = r.params.random.Intn(r.params.cols)

			if (randomC != col || randomR != row) && !r.cells[randomR][randomC].bomb {
				break
			}
		}

		r.cells[randomR][randomC].bomb = true
	}

	params := r.params
	for row := range params.rows {
		for c := range params.cols {
			cell := r.cells[row][c]
			if cell.bomb {
				continue
			}
			cell.adjacentCount = r.count(row, c)
			cell.displayedAdjCount = cell.adjacentCount
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

		row := r.params.random.Intn(r.params.rows)
		c := r.params.random.Intn(r.params.cols)

		cell := r.cells[row][c]

		if cell.bomb || cell.displayedAdjCount == 0 {
			continue
		}
		cell.displayedAdjCount--
	}
}

func (r *Board) Render() (window.Location, [][]window.Cell) {
	cells := make([][]window.Cell, 0)
	for row := range r.params.rows {
		cell_row := make([]window.Cell, 0)
		for c := range r.params.cols {
			cell_row = append(cell_row, toWindowCell(r.cells[row][c], row, c))
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
