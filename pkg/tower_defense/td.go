package tower_defense

import (
	"context"
	"log/slog"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/chat"
	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

type TD struct {
	Commander commands.Commander
	Renderer  *window.Renderer

	towers []*Tower
	params TDParams
	agg    ChatAggregator
	done   bool
}

type TDParams struct{}

func NewTD(params TDParams) TD {
	return TD{
		done:   false,
		params: params,
		agg:    NewChatAggregator(),
		towers: make([]*Tower, 0),

		Commander: commands.NewCommander(),
		Renderer:  window.NewRender(24, 80),
	}
}

func (t *TD) Done() bool {
	return t.done
}

func (t *TD) Tick(delta time.Duration) {
	row, col := t.agg.Reset()

	for _, t := range t.towers {
		if t.Row == row && t.Col == col {
			t.Count++
			return
		}
	}

	tower := NewTower(row, col)
	t.Renderer.Add(tower)
	t.towers = append(t.towers, tower)
}

func (t *TD) Render() []*window.CellWithLocation {
	return t.Renderer.Render()
}

func (t *TD) NewChatMsg(msg string) {
	row, col, err := ParseChatMessage(msg)
	if err != nil {
		return
	}

	if row > 24 || row < 0 || col > 80 || col < 0 {
		return
	}

	t.agg.Add(row, col)
}

func (t *TD) Start() {
	// TODO: Do i need to do anything here?
}

func LinkChatToTowerDefense(t *TD, c chan chat.ChatMsg, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c:
			if msg.Bits == 0 {
				slog.Debug("LinkChatToTowerDefense", "msg", msg.Msg)
				t.NewChatMsg(msg.Msg)
			}
		}
	}
}
