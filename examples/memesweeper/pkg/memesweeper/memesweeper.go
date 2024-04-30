package memesweeper

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
	"github.com/theprimeagen/vim-with-me/pkg/chat"
	"github.com/theprimeagen/vim-with-me/pkg/window"
	"github.com/theprimeagen/vim-with-me/pkg/window/components"
)

type MemeSweeperState struct {
	Skips     int
	StartTime time.Duration

	Reduces int
	Bombs   int
	Width   int
	Height  int
}

func NewMemeSweeperState(bombs, skips int) MemeSweeperState {
	return MemeSweeperState{
		Bombs:   bombs,
		Reduces: 0,
		Skips:   skips,
	}
}

func (ms MemeSweeperState) WithDims(height, width int) MemeSweeperState {
    ms.Height = height
    ms.Width = width
    return ms
}

type MemeSweeper struct {
	State    MemeSweeperState
	board    *Board
	Renderer window.Renderer
	grid     grid
	chat     ChatAggregator

	texts     []*components.Text
	clock     *components.Text
    next      *components.HighlightPoint
	startTime int64
}

func getTime(ms int64) string {
	return fmt.Sprintf("Time: %ds", ms/1000)
}

func NewMemeSweeper(state MemeSweeperState) MemeSweeper {
    assert.Assert(state.Height > 0, "please set the height of the minesweeper game")
    assert.Assert(state.Width > 0, "please set the width of the minesweeper game")

	clock := components.NewText(3, 14, getTime(0))
	texts := []*components.Text{
		components.NewText(0, 12, ":)"),
		components.NewText(2, 13, "Skips: 5"),
		components.NewText(3, 13, fmt.Sprintf("BombCount: %d", state.Bombs)),
		clock,
	}

	render := window.NewRender(state.Height, state.Width)
	params := BoardParams{
		width:  state.Width,
		height: state.Height,

		row:       1,
		col:       1,
		bombCount: state.Bombs,
	}

	board := NewBoard(params)
	grid := newGrid(0, 0, state.Width, state.Height)

	for _, t := range texts {
		render.Add(t)
	}
	render.Add(board)
	render.Add(grid)

	return MemeSweeper{
		State:    state,
		board:    board,
		texts:    texts,
		clock:    clock,
		chat:     NewChatAggregator(),
		Renderer: render,
	}
}

func (m *MemeSweeper) Pick(row, col int) {
    slog.Debug("MemeSweeper#Pick", "row", row, "col", col)
	m.board.PickSpot(row, col)
}

func (m *MemeSweeper) ReduceCount() {
	m.board.ReduceOne()
}

func (m *MemeSweeper) Dimensions() (byte, byte) {
    return byte(m.board.params.height + 1 + 1 + 1), byte(m.board.params.width + 1 + 15)
}

func (m *MemeSweeper) Chat(msg *chat.ChatMsg) {
    row, col, err := ParseChatMessage(msg.Msg)
    if err != nil {
        return
    }

    if row >= m.State.Height {
        return
    }

    c := int(col[0] - 'A')
    if c >= m.State.Width {
        return
    }

    slog.Debug("MemeSweeper#Chat", "row", row, "col", col)
    m.chat.Add(row, c)
}

func (m *MemeSweeper) PlayRound(deltaMS int64) {
    point := m.chat.Reset()
    m.board.PickSpot(point.row, point.col)
}

func (m *MemeSweeper) Render() []*window.CellWithLocation {
    if m.board.state == LOSE {
        txt := components.NewText(0, 12, ";(")
        m.Renderer.Clear()
        m.Renderer.Add(txt)
    } else if m.board.state == WIN {
        txt := components.NewText(0, 12, "8)")
        m.Renderer.Clear()
        m.Renderer.Add(txt)
    } else {
        now := time.Now().UnixMilli()
        m.clock.SetText(getTime(now - m.startTime))
        m.startTime = now

        current := m.chat.Current()
        if current.count != 0 {
        }
    }

	return m.Renderer.Render()
}
