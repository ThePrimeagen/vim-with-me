package main

import (
	"fmt"
	"log"
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
		log.Fatalf("Error reading file: %s", err)
	}

	_ = win.SetWindow(str)
}

func partialRender(win *window.Window, row, col byte, text []byte) {

	for i := 0; i < len(text); i++ {
		err := win.Set(row, col+byte(i), text[i])
		if err != nil {
			log.Fatalf("Error setting partial render: %s", err)
		}
	}
}

func main() {
	server, err := testies.CreateServerFromArgs()
	if err != nil {
		log.Fatalf("Error creating server: %s", err)
	}

    comms := commands.NewCommander()

    comms.AddCommand("open")

    server.Welcome(comms.ToCommands())

	win := window.NewWindow(24, 80)

	fmt.Printf("test\n")
	for {
		fmt.Printf("waiting on from socket \n")
		cmd := <-server.FromSockets
		fmt.Printf("command received %+v\n", cmd)

		switch cmd.Command {
        // Think about how to do better custom commands and really routing in
        // general
        case comms.GetCommandByte("open"):
			out := window.OpenCommand(win)
			server.Send(out)
        case commands.RENDER:
			render(win)
			str := win.Render()
			out := commands.Render([]byte(str))
			server.Send(out)
        case commands.PARTIAL_RENDER:
            row := cmd.Data[0]
            col := cmd.Data[1]
			partialRender(win, row, col, []byte("theprimeagen"))
			renders := win.PartialRender()
			fmt.Printf("partial render %d\n", len(renders))
			server.Send(commands.PartialRender(renders))
		}
	}
}
