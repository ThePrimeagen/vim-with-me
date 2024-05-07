package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	"github.com/theprimeagen/vim-with-me/pkg/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
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
    chars := ascii_buffer.NewFreqency()
    colors := ascii_buffer.NewFreqency()

    outFile, err := os.CreateTemp("/tmp", "doom")
    if err != nil {
        log.Fatal("couldn't create tmp")
    }

    d.Framer.DebugToFile(outFile)

    for i := range 100 {
        fmt.Printf("count: %d\n", i)
        frame := <-frames
        chars.Freq(frame.Chars)
        colors.Freq(frame.Color)
    }
}


