package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

func render(win *window.SimpleAsciiWindow) {
	bytes, err := os.ReadFile("lua/vim-with-me/integration/theprimeagen")
	str := string(bytes)
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "*", " ")

	if err != nil {
		slog.Error("Error reading file", "err", err)
	}

	_ = win.SetWindow(str)
}

func partialRender(win *window.SimpleAsciiWindow, row, col int, text []byte) {

	for i := 0; i < len(text); i++ {
		err := win.Set(row, col+i, text[i])
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
	server.WelcomeMessage(func() *tcp.TCPCommand { return commander.ToCommands() })
	win := window.NewSimpleWindow(24, 80)

	defer server.Close()
	go server.Start()

	for {
		wrapper := <-server.FromSockets
		slog.Info("command received", "cmd", wrapper)

		switch wrapper.Command.Command {
		// Think about how to do better custom commands and really routing in
		// general
		case commander.GetCommandByte("open"):
			out := commands.OpenCommand(&win)
			server.Send(out)
		case commands.RENDER:
			render(&win)
			cells := win.Render()
			out := commands.PartialRender(cells)
			server.Send(out)
		case commands.PARTIAL_RENDER:
			row := int(wrapper.Command.Data[0])
			col := int(wrapper.Command.Data[1])
			partialRender(&win, row, col, []byte("theprimeagen"))
			renders := win.PartialRender()
			server.Send(commands.PartialRender(renders))
		}
	}
}
