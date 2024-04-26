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
	Renderer  window.Renderer

	params TDParams
	chat   *chat.Chat
	done   chan struct{}
}

type TDParams struct {
	ChatWaitTime time.Duration
}

func NewTD(c *chat.Chat, params TDParams) TD {
	return TD{
		chat:   c,
		done:   make(chan struct{}, 1),
		params: params,

		Commander: commands.NewCommander(),
		Renderer:  window.NewRender(24, 80),
	}
}

func (t *TD) run(ctx context.Context) {

    ticker := time.NewTicker(t.params.ChatWaitTime)
    defer ticker.Stop()

    outer:
    for {
        select {
        case <-ctx.Done():
            slog.Warn("tower defense game over: canceled by context")
            break outer
        case msg := <-t.chat.Messages:
        case <-ticker.C:
        }
    }

    t.done <- struct{}{}
}

func (t *TD) Start(ctx context.Context) chan struct{} {
	return t.done
}
