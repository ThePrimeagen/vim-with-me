package main

import (
	"log/slog"

	"github.com/theprimeagen/vim-with-me/pkg/chat"
	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/tower_defense"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

func main() {
    testies.SetupLogger()
    server, err := testies.CreateServerFromArgs()
    if err != nil {
        slog.Error("received error on startup", "err", err)
        return
    }

    c := chat.NewChat("./chat-fake.js")
    done, err := c.Start()

    if err != nil {
        slog.Error("unable to start chat", "err", err)
        return
    }

    td := tower_defense.NewTD(&c)
    server.WelcomeMessage(td.Commander.ToCommands())
    server.WelcomeMessage(commands.OpenCommand(&td.Renderer))

    defer server.Close()

    go server.Start()
}

