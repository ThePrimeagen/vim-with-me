package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser/display"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/huffman"
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

	frames := d.Frames()
    colors := ascii_buffer.NewFreqency()
	_ = ascii_buffer.NewFreqency()

    outer:
	for range 1000 {
        select{
        case frame := <-frames:
            colors.Freq(frame.Color8BitIterator())
            huff := huffman.CalculateHuffman(colors)
            huffBuff := make([]byte, len(frame.Color), len(frame.Color))

            bitLen, err := huff.Encode(frame.Color8BitIterator(), huffBuff)
            fmt.Println(display.Display(&frame, d.Rows, d.Cols))
            fmt.Fprintf(os.Stderr, "huff: %d bitLen: %d -- err: %v\n", len(huff.DecodingTree), bitLen / 8 + 1, err)

            //out := make([]byte, len(frame.Color), len(frame.Color))
            //writer := byteutils.U8Writer{}
            //writer.Set(out)

            //huff.Decode(huffBuff, bitLen, &writer)
            //frame.Color = out
            //fmt.Println(display.Display(&frame, d.Rows, d.Cols))

            break outer;
        case <-finish:
            break outer;
        }
	}

}
