package memesweeper

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
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
	Rand    *rand.Rand
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

func (ms MemeSweeperState) WithSeed(seed int64) MemeSweeperState {
	ms.Rand = rand.New(rand.NewSource(seed))
	return ms
}

type MemeSweeper struct {
	State    MemeSweeperState
	board    *Board
	Renderer *window.Renderer
	grid     grid
	chat     *ChatAggregator

	texts      []*components.Text
	clock      *components.Text
	smiley     *components.Text
	control    *components.Text
	timePassed int64

	currentPick *components.HighlightPoint
}

func getDimensions(params MemeSweeperState) (int, int) {
	return params.Height + 1 + 1 + 1 + 1, params.Width + 1 + 15
}

func getTime(ms int64) string {
	return fmt.Sprintf("Time: I Suck")
}

func NewMemeSweeper(state MemeSweeperState) MemeSweeper {
	assert.Assert(state.Height > 0, "please set the height of the minesweeper game")
	assert.Assert(state.Width > 0, "please set the width of the minesweeper game")

	clock := components.NewText(3, 14, getTime(0))
	smiley := components.NewText(0, 12, ":)")
	control := components.NewText(0, 0, "Waiting...")
	texts := []*components.Text{
		components.NewText(4, 13, fmt.Sprintf("BombCount: %d", state.Bombs)),
		clock,
		smiley,
	}

	rows, cols := getDimensions(state)
	render := window.NewRender(rows, cols)
	params := BoardParams{
		cols: state.Width,
		rows: state.Height,

		row:       2,
		col:       1,
		bombCount: state.Bombs,
		random:    state.Rand,
	}

	chatAgg := NewChatAggregator()

	pickPos := components.NewCompositePosition(chatAgg, window.NewLocation(params.row, params.col))
	pick := components.NewHighlightPoint(pickPos, 100, components.BACKGROUND_RED)
	board := NewBoard(params)
	grid := newGrid(1, 0, state.Width, state.Height)

	for _, t := range texts {
		render.Add(t)
	}
	render.Add(board)
	render.Add(pick)
	render.Add(grid)

	return MemeSweeper{
		State:       state,
		board:       board,
		texts:       texts,
		clock:       clock,
		control:     control,
		smiley:      smiley,
		chat:        chatAgg,
		Renderer:    render,
		timePassed:  0,
		currentPick: pick,
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
	rows, cols := getDimensions(m.State)
	return byte(rows), byte(cols)
}

func (m *MemeSweeper) Chat(msg *chat.ChatMsg) error {
	row, col, err := ParseChatMessage(msg.Msg)
	if err != nil {
		return err
	}

	// row is 1 based
	row -= 1

	if row >= m.State.Height {
		return errors.New("row is too big")
	}

	c := int(col[0] - 'A')
	if c >= m.State.Width {
		return errors.New("col is too big")
	}

	slog.Debug("MemeSweeper#Chat", "row", row, "col", col)
	m.chat.Add(row, c)

	return nil
}

func (m *MemeSweeper) StartRound() {
	m.control.SetText("Pick Pos!")
	m.currentPick.SetActiveState(true)
	m.chat.SetActiveState(true)
}

func (m *MemeSweeper) EndRound() {
	m.control.SetText("Waiting...")
	m.currentPick.SetActiveState(false)
	m.chat.SetActiveState(false)

	point := m.chat.Reset()
	m.board.PickSpot(point.row, point.col)

	if m.board.state == LOSE {
		m.smiley.SetText(";(")
		m.control.SetText("L Take")
	} else if m.board.state == WIN {
		m.smiley.SetText("8)")
		m.control.SetText("W Take")
	}
}

func (m *MemeSweeper) RevealBombs() {
	m.board.RevealBombs()
}

func (m *MemeSweeper) Reset() {
	m.chat.Reset()
	m.chat.SetActiveState(false)
	m.board.Reset()
	m.currentPick.SetActiveState(false)

	m.smiley.SetText(":)")
	m.control.SetText("Waiting...")
	m.clock.SetText(getTime(0))
	m.timePassed = 0
}

func (m *MemeSweeper) GameOver() bool {
	return m.board.state != PLAYING && m.board.state != INIT
}

func (m *MemeSweeper) Render(timePassedMS int64) []*window.CellWithLocation {
	if !m.GameOver() {
		m.timePassed += timePassedMS
		m.clock.SetText(getTime(m.timePassed))
	}

	return m.Renderer.Render()
}
