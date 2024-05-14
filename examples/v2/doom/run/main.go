package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser/display"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	//"github.com/theprimeagen/vim-with-me/pkg/v2/encoding"
	"github.com/theprimeagen/vim-with-me/pkg/v2/program"
)

func main() {
	debug := ""
	flag.StringVar(&debug, "debug", "", "runs the file like the program instead of running doom")

	assertF := ""
	flag.StringVar(&assertF, "assert", "", "add an assert file")

	rounds := 1000
	flag.IntVar(&rounds, "rounds", 1000, "the rounds of doom to play")


    flag.Parse()
    args := flag.Args()
    name := args[0]

    fmt.Printf("assert file attached \"%s\"\n", assertF)
    fmt.Printf("debug file attached \"%s\"\n", debug)
    fmt.Printf("args file attached \"%v\"\n", args)

    d := doom.NewDoom()

    prog := program.
        NewProgram(name).
        WithArgs(args[1:]).
        WithWriter(d)

    if debug != "" {
        debugFile, err := os.Create(debug)
        assert.NoError(err, "unable to open debug file")
        prog = prog.WithWriter(debugFile)
    }

    if assertF != "" {
        assertFile, err := os.Create(assertF)
        assert.NoError(err, "unable to open assert file")
        assert.ToWriter(assertFile)
    }

    ctx := context.Background()
    go func() {
        err := prog.Run(ctx)
        assert.NoError(err, "prog.Run(ctx)")
    }()

    <-d.Ready()

    frames := d.Frames()
    _ = ascii_buffer.NewFreqency()
    _ = ascii_buffer.NewFreqency()

    for range rounds {
        frame := <-frames

        fmt.Print("\033[H\033[2J")
        fmt.Println(display.Display(&frame, d.Rows, d.Cols))
    }

}


