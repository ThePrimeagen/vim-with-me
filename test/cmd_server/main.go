package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

func render(win *window.Window) {
	bytes, err := os.ReadFile("lua/vim-with-me/integration/theprimeagen")
	str := string(bytes)
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "*", " ")

	if err != nil {
		slog.Error("Error reading file", "err", err)
	}

	_ = win.SetWindow(str)
}

func partialRender(win *window.Window, row, col byte, text []byte) {

	for i := 0; i < len(text); i++ {
		err := win.Set(row, col+byte(i), text[i])
		if err != nil {
			slog.Error("Error setting partial render", "err", err)
		}
	}
}

func main() {
    testies.SetupLogger()
	server, err := testies.CreateServerFromArgs()
	if err != nil {
		slog.Error("Error creating server: %s", err)
	}

	commander := commands.NewCommander()
	commander.AddCommand("open")
	server.WelcomeMessage(commander.ToCommands())
	win := window.NewWindow(24, 80)

    defer server.Close()
    go server.Start()

	for {
		wrapper := <-server.FromSockets
		slog.Info("command received", "cmd", wrapper)

		switch wrapper.Command.Command {
		// Think about how to do better custom commands and really routing in
		// general
		case commander.GetCommandByte("open"):
			out := window.OpenCommand(win)
			server.Send(out)
		case commands.RENDER:
			render(win)
			str := win.Render()
			out := commands.Render([]byte(str))
			server.Send(out)
		case commands.PARTIAL_RENDER:
			row := wrapper.Command.Data[0]
			col := wrapper.Command.Data[1]
			partialRender(win, row, col, []byte("theprimeagen"))
			renders := win.PartialRender()
			fmt.Printf("partial render %d\n", len(renders))
			server.Send(commands.PartialRender(renders))
		}
	}
}
