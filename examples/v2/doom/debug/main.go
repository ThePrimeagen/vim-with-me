package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/encoder"
)

func main() {
	debug := ""
	flag.StringVar(&debug, "debug", "", "runs the file like the program instead of running doom")
	flag.Parse()

	//outFile, err := os.CreateTemp("/tmp", "doom-assert")
	//assert.NoError(err, "couldnt create temp file")
	//assert.ToWriter(outFile)

	d := doom.NewDoom()
	finish := make(chan struct{}, 1)
	go func() {
		f, err := os.Open(debug)
		assert.NoError(err, "unable to read debug file")

        io.Copy(d, f)
        <-time.After(time.Millisecond * 100)
		finish <- struct{}{}
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

    fmt.Printf("starting\n")
    outer:
    for i := range 1000 {
        fmt.Printf("loop %d\n", i)
        select{
        case frame := <-frames:
            data := ansiparser.RemoveAsciiStyledPixels(frame.Color)
            fmt.Printf("encoding frame %d\n", len(data))
            encFrame := enc.PushFrame(data)

            if encFrame == nil {
                fmt.Printf("encoded failed to produce smaller frame: %d\n", len(data))
                break
            }

            fmt.Printf("encoded(%d): %d\n", encFrame.Encoding, encFrame.Len)

        case <-finish:
            break outer;
        }

        fmt.Printf("done with select on %d\n", i)
	}

}
