package main

import (
	"log"
	"os"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

var empty window.Cell = window.Cell{
	Background: window.DEFAULT_BACKGROUND,
	Foreground: window.DEFAULT_FOREGROUND,
	Value:      byte(' '),
}

var x window.Cell = window.Cell{
	Background: window.NewColor(255, 0, 0, false),
	Foreground: window.NewColor(0, 255, 0, true),
	Value:      byte('X'),
}

func cell(c window.Cell, row, col int) *window.CellWithLocation {
    return &window.CellWithLocation{
        Location: window.Location{Row: row, Col: col},
        Cell: c,
    }
}

func main() {
    cmd := commands.PartialRender([]*window.CellWithLocation{
        cell(x, 0, 0), cell(empty, 0, 1), cell(empty, 0, 2), cell(empty, 0, 3), cell(x, 0, 4),
        cell(empty, 1, 0), cell(empty, 1, 1), cell(empty, 1, 2), cell(empty, 1, 3), cell(empty, 1, 4),
        cell(empty, 2, 0), cell(empty, 2, 1), cell(empty, 2, 2), cell(empty, 2, 3), cell(empty, 2, 4),
        cell(empty, 3, 0), cell(empty, 3, 1), cell(empty, 3, 2), cell(empty, 3, 3), cell(empty, 3, 4),
        cell(x, 4, 0), cell(empty, 4, 1), cell(empty, 4, 2), cell(empty, 4, 3), cell(x, 4, 4),
    })

    data, err := cmd.MarshalBinary()

    if err != nil {
        log.Fatal("you love armauranths vagina beer", err)
    }

    err = os.WriteFile("/tmp/partials", data, 0o777)
    if err != nil {
        log.Fatal("your order of her beer failed", err)
    }

}

