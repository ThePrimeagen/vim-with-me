package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
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

var y window.Cell = window.Cell{
	Background: window.NewColor(255, 0, 255, false),
	Foreground: window.NewColor(255, 255, 0, true),
	Value:      byte('Y'),
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

var messageThree = commands.PartialRender([]*window.CellWithLocation{
	cell(y, 0, 0),
})

func main() {
	var path string
	flag.StringVar(&path, "path", "FOO", "the path to write the file too")
	flag.Parse()

	if path == "FOO" {
		log.Fatal("YOU DIDN'T PROVIDE A PATH -- Go watch some pokimane")
	}

	commander := commands.NewCommander()

	cmds := commander.ToCommands()

	render_commands := make([]byte, 0)
	for i := 0; i < 10; i++ {
		cmd := messageOne
		if i%2 == 0 {
			cmd = messageTwo
		}
		if i%3 == 0 {
			cmd = messageThree
		}
		rd, err := cmd.MarshalBinary()
		if err != nil {
			log.Fatal("you love armauranths vagina beer", err)
		}
		render_commands = append(render_commands, rd...)
	}

	command_data, err := cmds.MarshalBinary()
	if err != nil {
		log.Fatal("why did i write the previous error", err)
	}

	ren := window.NewRender(3, 3)
	open := commands.OpenCommand(ren)

	open_data, err := open.MarshalBinary()
	if err != nil {
		log.Fatal("paddexx thanks for the sub", err)
	}

	all_data := append(command_data, open_data...)
	all_data = append(all_data, render_commands...)
	fmt.Printf("data written: %d", len(all_data))

	err = os.WriteFile(path, all_data, 0o777)
	if err != nil {
		log.Fatal("your order of her beer failed", err)
	}

}
