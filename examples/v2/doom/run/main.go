package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/encoder"

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

    enc := encoder.NewEncoder(d.Rows * (d.Cols / 2), ascii_buffer.QuadtreeParam{
        Depth: 2,
        Stride: 1,
        Rows: d.Rows,
        Cols: d.Cols / 2,
    })

    enc.AddEncoder(encoder.XorRLE)
    enc.AddEncoder(encoder.Huffman)

	frames := d.Frames()

    for i := range 1000 {
        select{
        case frame := <-frames:
            original := len(frame.Color)
            data := ansiparser.RemoveAsciiStyledPixels(frame.Color)
            encFrame := enc.PushFrame(data)

            if encFrame == nil {
                fmt.Printf("encoded failed to produce smaller frame: %d\n", len(data))
                break
            }

            fmt.Printf("(%d) frame: %d -- encoded(%d): %d\n", original, i, encFrame.Encoding, encFrame.Len)
        }
	}
}


