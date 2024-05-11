package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	"github.com/theprimeagen/vim-with-me/pkg/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser/display"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	//"github.com/theprimeagen/vim-with-me/pkg/v2/encoding"
	"github.com/theprimeagen/vim-with-me/pkg/v2/program"
)

func main() {
    flag.Parse()
    args := flag.Args()
    name := args[0]

    d := doom.NewDoom()
    prog := program.
        NewProgram(name).
        WithArgs(args[1:]).
        WithWriter(d)

    ctx := context.Background()
    go func() {
        err := prog.Run(ctx)
        assert.NoError(err, "prog.Run(ctx)")
    }()

    <-d.Ready()

    frames := d.Frames()
    _ = ascii_buffer.NewFreqency()
    _ = ascii_buffer.NewFreqency()

    outFile, err := os.CreateTemp("/tmp", "doom")
    if err != nil {
        log.Fatal("couldn't create tmp")
    }
    d.Framer.DebugToFile(outFile)

    for range 100 {
        frame := <-frames

        fmt.Print("\033[H\033[2J")
        fmt.Println(display.Display(&frame, d.Rows, d.Cols))
    }

}


