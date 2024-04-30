package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/memesweeper/pkg/memesweeper"
	"github.com/theprimeagen/vim-with-me/pkg/chat"
	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
)

func run() {
}

func main() {
    testies.SetupLogger()
    server, err := testies.CreateServerFromArgs()

    ctx, cancel := context.WithCancel(context.Background())
    _ = cancel

    ch, err := chat.NewTwitchChat(ctx)
    if err != nil {
        slog.Error("chat.Start()", "err", err)
        return
    }

    state := memesweeper.NewMemeSweeperState(10, 5).WithDims(5, 10)
    ms := memesweeper.NewMemeSweeper(state)

    commander := commands.NewCommander()
    server.WelcomeMessage(commander.ToCommands())
    server.WelcomeMessage(commands.OpenCommand(&ms))

    go server.Start()
    defer server.Close()

    go func() {
        ticker := time.NewTicker(time.Millisecond * 16)
        outer:
        for {
            start := time.Now().UnixMilli()
            select {
            case <-ctx.Done():
                break outer
            case msg := <-ch:
                slog.Debug("main: msg received", "msg", msg.Msg, "name", msg.Name)
                ms.Chat(&msg)
            case <-ticker.C:
                cells := ms.Render(time.Now().UnixMilli() - start)
                cmds := commands.PartialRender(cells)
                server.Send(cmds)
            }
        }
    }()

    for {
        ms.StartRound()
        <-time.After(time.Second * 10)
        ms.EndRound()
        <-time.After(time.Second * 5)
    }
}


