package main

import (
	"log/slog"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

type X struct {
	row int
	col int

	foreground window.Color
}

func (x *X) Z() int {
    return 1
}

func (x *X) Id() int {
    return 1
}

func (x *X) Render() (window.Location, [][]window.Cell) {
    cells := [][]window.Cell{
        {{Background: window.DEFAULT_BACKGROUND, Foreground: x.foreground, Value: byte('X')}},
    }
    return window.NewLocation(x.row, x.col), cells
}

func main() {
	testies.SetupLogger()
	server, err := testies.CreateServerFromArgs()

	if err != nil {
		slog.Error("Error creating server: %s", err)
	}

	commander := commands.NewCommander()
	renderer := window.NewRender(24, 80)

	server.WelcomeMessage(commander.ToCommands())
	server.WelcomeMessage(commands.OpenCommand(&renderer))

    renderer.Add(&X{row: 8, col: 8, foreground: window.NewColor(100, 169, 69, true)})
    cells := renderer.Render()

	server.WelcomeMessage(commands.PartialRender(cells))

	defer server.Close()
	go server.Start()

	for {
		<-server.FromSockets
	}
}
