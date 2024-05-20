package memesweeper

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

func getBombs(board *Board) []window.Location {
	bombs := make([]window.Location, 0)
	for r, row := range board.cells {
		for c, cell := range row {
			if cell.bomb {
				bombs = append(bombs, window.NewLocation(r, c))
			}
		}
	}
	return bombs
}

func hasBomb(row, col int, bombs []window.Location) bool {
	return slices.ContainsFunc(bombs, func(b window.Location) bool {
		return b.Col == col && b.Row == row
	})
}

func pickGoodSpot(board *Board, bombs []window.Location) (int, int) {
	for r := range board.params.rows {
		for c := range board.params.cols {
			if !hasBomb(r, c, bombs) && !board.cells[r][c].revealed {
				return r, c
			}
		}
	}
	return -1, -1
}

func revealed(t *testing.T, board *Board, bools [][]bool) {
	for r, row := range bools {
		for c, is := range row {
			require.Equal(t, is, board.cells[r][c].revealed)
		}
	}
}

func TestBoard(t *testing.T) {
	params := BoardParams{
		random:    rand.New(rand.NewSource(69420)),
		cols:      3,
		rows:      3,
		row:       1,
		col:       1,
		bombCount: 2,
	}

	board := NewBoard(params)
	require.Equal(t, 0, len(getBombs(board)))

	// true0 false* false0
	// false* false0 false0
	// false0 false0 false0

	board.PickSpot(0, 0)
	revealed(t, board, [][]bool{
		{true, false, false},
		{false, false, false},
		{false, false, false},
	})
	require.Equal(t, 0, board.cells[0][1].adjacentCount) // bomb
	require.Equal(t, 2, board.cells[1][1].adjacentCount)

	board.PickSpot(2, 2)
	revealed(t, board, [][]bool{
		{true, false, false},
		{false, true, true},
		{false, true, true},
	})

	require.Equal(t, PLAYING, board.state)
	board.PickSpot(0, 1)
	require.Equal(t, LOSE, board.state)
}

func TestBoardWin(t *testing.T) {
	params := BoardParams{
		random:    rand.New(rand.NewSource(69420)),
		cols:      3,
		rows:      3,
		row:       1,
		col:       1,
		bombCount: 2,
	}

	board := NewBoard(params)
	require.Equal(t, 0, len(getBombs(board)))

	// true0 false* false0
	// false* false0 false0
	// false0 false0 false0

	board.PickSpot(0, 0)
	require.Equal(t, PLAYING, board.state)
	board.PickSpot(2, 2)
	require.Equal(t, PLAYING, board.state)
	board.PickSpot(0, 2)
	require.Equal(t, PLAYING, board.state)
	board.PickSpot(2, 0)
	require.Equal(t, WIN, board.state)
}
