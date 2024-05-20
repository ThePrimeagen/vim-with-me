package doom_test

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser/display"
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

	timer := time.NewTimer(1000000000 * time.Millisecond)

	select {
	case f := <-frames:
		fmt.Println(display.Display(&f, d.Rows, d.Cols))
	case <-timer.C:
		panic("YOU SUCK")
	}
}

/*
func TestDoom8BitParserManyFrame(t *testing.T) {
    data, err := os.Open("./doomtest_large")
    require.NoError(t, err)

    d := doom.NewDoom()

    go func() {
        defer data.Close()
        io.Copy(d, data)
    }()

    <-d.Ready()
    fh, err := os.CreateTemp("/tmp", "doomtest")
    require.NoError(t, err)
    d.Framer.DebugToFile(fh)

    require.Equal(t, d.Rows, 66)
    require.Equal(t, d.Cols, 212)

    frames := d.Frames()
    count := 0
    length := d.Rows * d.Cols

    outer:
    for {
        timer := time.NewTimer(time.Millisecond * 50)
        select {
        case frame := <-frames:
            require.Equal(t, len(frame.Chars), length)
            require.Equal(t, len(frame.Color), length)
            count++
            fmt.Println(display.Clear())
            fmt.Println(display.Display(&frame, d.Rows, d.Cols))
        case <-timer.C:
            break outer
        }
    }

    require.Equal(t, 129, count)
}
*/
