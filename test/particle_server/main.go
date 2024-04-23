package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/particles"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

func main() {
    testies.SetupLogger()
	server, err := testies.CreateServerFromArgs()
	if err != nil {
		slog.Error("Error creating server: %s", err)
	}

	commander := commands.NewCommander()
	renderer := window.NewRender(8, 61)

	server.WelcomeMessage(commander.ToCommands())
	server.WelcomeMessage(commands.OpenCommand(&renderer))

    defer server.Close()
    go server.Start()

    coffee := particles.NewCoffee(61, 8, 9.0)
    coffee.Start()
    renderer.Add(&coffee)

    timer := time.NewTicker(100 * time.Millisecond)
    for _ = range timer.C {
        coffee.Update()
        partials := renderer.Render()
        fmt.Printf("partials: %+v\n", partials)
        server.Send(commands.PartialRender(partials))
    }
}
