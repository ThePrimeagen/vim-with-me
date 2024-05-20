package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/memesweeper/pkg/memesweeper"
	"github.com/theprimeagen/vim-with-me/pkg/chat"
	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
)

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

	state := memesweeper.NewMemeSweeperState(15, 5).WithDims(5, 13)
	ms := memesweeper.NewMemeSweeper(state)

	commander := commands.NewCommander()
	server.WelcomeMessage(func() *tcp.TCPCommand { return commander.ToCommands() })
	server.WelcomeMessage(func() *tcp.TCPCommand { return commands.OpenCommand(&ms) })

	go server.Start()
	defer server.Close()

	start := time.Now().UnixMilli()
	go func() {
		ticker := time.NewTicker(time.Millisecond * 16)
	outer:
		for {
			select {
			case <-ctx.Done():
				break outer
			case msg := <-ch:
				if ms.GameOver() {
					break
				}
				slog.Debug("main: msg received", "msg", msg.Msg, "name", msg.Name)
				ms.Chat(&msg)
			case <-ticker.C:
				if ms.GameOver() {
					break
				}
				cells := ms.Render(time.Now().UnixMilli() - start)
				if len(cells) == 0 {
					break
				}
				cmds := commands.PartialRender(cells)
				server.Send(cmds)
			}
		}
	}()

	server.WelcomeMessage(func() *tcp.TCPCommand {
		cells := ms.Renderer.FullRender()
		return commands.PartialRender(cells)
	})

	for {
		for !ms.GameOver() {
			<-time.After(time.Second * 7)
			ms.StartRound()
			<-time.After(time.Second * 12)
			ms.EndRound()
		}

		ms.RevealBombs()
		cells := ms.Render(time.Now().UnixMilli() - start)
		cmds := commands.PartialRender(cells)
		server.Send(cmds)
		<-time.After(time.Second * 30)
		ms.Reset()
	}
}
