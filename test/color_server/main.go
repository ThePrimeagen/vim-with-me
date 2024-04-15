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

    x := &X{row: 10, col: 10, foreground: window.NewColor(100, 169, 69, true)}
    renderer.Add(x)

	defer server.Close()
	go server.Start()

    ticker := time.NewTicker(time.Millisecond * 250)

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
