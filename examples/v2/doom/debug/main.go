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
            fmt.Fprintf(os.Stderr, "TOTAL: %d\n", len(huff.DecodingTree) + bitLen / 8 + 1)

            boxes := ascii_buffer.Partition(frame.Color, d.Rows, d.Cols, 4, 1)

            totalHuff := 0
            totalBits := 0
            for i, b := range boxes {
                freq := ascii_buffer.NewFreqency()
                freq.Freq(b)

                huff := huffman.CalculateHuffman(freq)
                huffBuff := make([]byte, len(frame.Color), len(frame.Color))

                b.Reset()

                bitLen, err := huff.Encode(b, huffBuff)
                fmt.Fprintf(os.Stderr, "BOX(%d, %d): %d huff: %d bitLen: %d -- err: %v\n", i, freq.Length(), b.Len(), len(huff.DecodingTree), bitLen / 8 + 1, err)

                totalBits += bitLen
                totalHuff += len(huff.DecodingTree)
            }

            fmt.Fprintf(os.Stderr, "Total(quad): %d\n", totalHuff + totalBits / 8 + 1)

            break outer;
        case <-finish:
            break outer;
        }
	}

}
