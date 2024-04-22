package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
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
	Background: window.NewColor(255, 0, 0, false),
	Foreground: window.NewColor(0, 255, 0, true),
	Value:      byte('X'),
}

func (x *X) Render() (window.Location, [][]window.Cell) {
	cells := [][]window.Cell{
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
		{cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell, cell},
	}
	return window.NewLocation(x.row-1, x.col-1), cells
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

	x := &X{row: 10, col: 10}
	renderer.Add(x)

	defer server.Close()
	go server.Start()

	ticker := time.NewTicker(time.Millisecond * 2500)

	count := 0
	for {
		<-ticker.C

		localCount := count % 20
		if localCount >= 10 {
			x.col--
			x.row--
		} else {
			x.col++
			x.row++
		}

		count++

		cells := renderer.Render()
		fmt.Printf("cells = %+v\n", cells)
		fmt.Printf("x = %+v\n", x)

		server.Send(commands.PartialRender(cells))

	}
}
