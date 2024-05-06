package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	"github.com/theprimeagen/vim-with-me/pkg/assert"
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
    frame := <-frames


    count := 0
    for range d.Rows {
        fmt.Println(string(frame.Chars[count:count + d.Cols]))
        count += d.Cols
    }
}

