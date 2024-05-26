package doom_test

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
	"github.com/theprimeagen/vim-with-me/pkg/v2/encoder"
	"github.com/theprimeagen/vim-with-me/pkg/v2/net"
)

func TestDoom8BitParserOneFrame(t *testing.T) {
	data, err := os.Open("./doomtest")
	require.NoError(t, err)

	d := doom.NewDoom()

	go func() {
		defer data.Close()
		io.Copy(d, data)
	}()

	<-d.Ready()

	require.Equal(t, 50, d.Rows)
	require.Equal(t, 160, d.Cols)

	frames := d.Frames()

	enc := encoder.NewEncoder(d.Rows*(d.Cols/2), ascii_buffer.QuadtreeParam{
		Depth:  2,
		Stride: 1,
		Rows:   d.Rows,
		Cols:   d.Cols / 2,
	})

	enc.AddEncoder(encoder.XorRLE)
	enc.AddEncoder(encoder.Huffman)

	timer := time.NewTimer(2000 * time.Millisecond)

	select {
	case frame := <-frames:
		data := ansiparser.RemoveAsciiStyledPixels(frame.Color)
		encFrame := enc.PushFrame(data)

		require.NotNil(t, encFrame)
		require.Equal(t, encoder.HUFFMAN, encFrame.Encoding)

		require.Equal(t, 3253, encFrame.Len)
		require.Equal(t, 834, len(encFrame.Huff.DecodingTree))

		// Into frame
		frameData := make([]byte, 4096, 4096)
		frameable := net.Frameable{Item: encFrame}
		n, err := frameable.Into(frameData, 0)

		require.NoError(t, err)
		require.Equal(t, n, 3264)

		// Decode huffman

		writer := byteutils.U8Writer{}
		writer.Set(frameData)

		fmt.Printf("first byte: 0x%2x\n", data[0])
		fmt.Printf("enc first byte: 0x%2x\n", encFrame.Curr[0])
		//encFrame.Huff.Decode(encFrame.Curr[:1], 4, &writer)

	case <-timer.C:
		assert.Assert(false, "YOU SUCK")
	}
}
