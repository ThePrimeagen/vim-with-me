package main

import (
	"log/slog"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

var empty window.Cell = window.Cell{
	Background: window.DEFAULT_BACKGROUND,
	Foreground: window.DEFAULT_FOREGROUND,
	Value:      byte('X'),
}

var x window.Cell = window.Cell{
	Background: window.NewColor(255, 0, 0, false),
	Foreground: window.NewColor(0, 255, 0, true),
	Value:      byte('X'),
}

func cell(c window.Cell, row, col int) *window.CellWithLocation {
	return &window.CellWithLocation{
		Location: window.Location{Row: row, Col: col},
		Cell:     c,
	}
}

var messageOne = commands.PartialRender([]*window.CellWithLocation{
	cell(x, 0, 0), cell(empty, 0, 1), cell(x, 0, 2),
	cell(empty, 1, 0), cell(x, 1, 1), cell(empty, 1, 2),
	cell(x, 2, 0), cell(empty, 2, 1), cell(x, 2, 2),
})

var messageTwo = commands.PartialRender([]*window.CellWithLocation{
	cell(empty, 0, 0), cell(x, 0, 1), cell(empty, 0, 2),
	cell(x, 1, 0), cell(empty, 1, 1), cell(x, 1, 2),
	cell(empty, 2, 0), cell(x, 2, 1), cell(empty, 2, 2),
})

func main() {
	testies.SetupLogger()
	server, err := testies.CreateServerFromArgs()
	defer server.Close()

	if err != nil {
		slog.Error("Error creating server: %s", err)
		return
	}

	commander := commands.NewCommander()
	renderer := window.NewRender(3, 3)

	server.WelcomeMessage(tcp.MakeWelcome(commander.ToCommands()))
	server.WelcomeMessage(tcp.MakeWelcome(commands.OpenCommand(renderer)))

	ticker := time.NewTicker(time.Second)

	go server.Start()

	for {
		<-ticker.C

		count := 0
		if count%2 == 0 {
			server.Send(messageOne)
		} else {
			server.Send(messageTwo)
		}
		count++
	}
}
