package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/chat"
	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/tower_defense"
)

func main() {
	fake := false
	flag.BoolVar(&fake, "fake", false, "use fake chat instead of real chat")

	testies.SetupLogger()
	server, err := testies.CreateServerFromArgs()

	if err != nil {
		slog.Error("received error on startup", "err", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel

	fmt.Printf("starting twitch chat\n")
	chat, err := chat.NewTwitchChat(ctx)
	if err != nil {
		slog.Error("chat.Start()", "err", err)
		return
	}

	params := tower_defense.TDParams{}

	fmt.Printf("new tower defense\n")
	td := tower_defense.NewTD(params)
	server.WelcomeMessage(func() *tcp.TCPCommand { return td.Commander.ToCommands() })
	server.WelcomeMessage(func() *tcp.TCPCommand { return commands.OpenCommand(td.Renderer) })

	defer server.Close()

	go server.Start()
	go tower_defense.LinkChatToTowerDefense(&td, chat, ctx)

	ticker := time.NewTicker(time.Second * 10)
	fmt.Printf("about to start\n")
	for range ticker.C {

		td.Tick(time.Second * 10)

		cells := td.Renderer.Render()
		cmds := commands.PartialRender(cells)
		fmt.Printf("cells: %d --- %+v\n", len(cells), cells)

		if len(cells) == 0 {
			continue
		}

		server.Send(cmds)
	}

	ticker.Stop()
}
