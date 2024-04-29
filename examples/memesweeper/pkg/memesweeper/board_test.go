package memesweeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

func TestBoard(t *testing.T) {
    params := BoardParams{
        width: 3,
        height: 3,
        row: 1,
        col: 1,
        bombCount: 2,
    }

    board := NewBoard(params)
    cells := board.cells

    bombs := make([]window.Location, 0)
    for r, row := range cells {
        for c, cell := range row {
            if cell.bomb {
                bombs = append(bombs, window.NewLocation(r, c))
            }
        }
    }

    require.Equal(t, 2, len(bombs))
}

