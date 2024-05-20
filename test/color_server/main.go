package main

import (
	"fmt"
	"log/slog"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

type X struct {
	row int
	col int
}

func (x *X) Z() int {
	return 1
}

func (x *X) Id() int {
	return 1
}

var cell window.Cell = window.Cell{
	Background: window.NewColor(255, 0, 69, false),
	Foreground: window.NewColor(69, 255, 42, true),
	Value:      byte('X'),
}

func (x *X) Render() (window.Location, [][]window.Cell) {
	cells := [][]window.Cell{
		{cell},
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

	x := &X{row: 6, col: 9}
	renderer.Add(x)
	cells := renderer.Render()

	server.WelcomeMessage(func() *tcp.TCPCommand { return commander.ToCommands() })
	server.WelcomeMessage(func() *tcp.TCPCommand { return commands.OpenCommand(renderer) })
	server.WelcomeMessage(func() *tcp.TCPCommand { return commands.PartialRender(cells) })

	fmt.Printf("Does this work?\n")

	defer server.Close()
	server.Start()
}
