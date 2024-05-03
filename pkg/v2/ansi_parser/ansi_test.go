package ansiparser_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
)

func TestDoom8BitParserOneFrame(t *testing.T) {
    data, err := os.ReadFile("./doomtest")
    require.NoError(t, err)

    doomHeader := ansiparser.NewDoomAsciiHeaderParser()
    n, err := doomHeader.Write(data)
    require.NoError(t, err)

    data = data[n:]
    rows, cols := doomHeader.GetDims()

    doomAscii := ansiparser.NewDoomAnsiFramer(rows, cols)
    doomAscii.Write(data)

    require.Equal(t, 50, rows)
    require.Equal(t, 160, cols)

    frames := doomAscii.Frames()

    <-frames
}

func TestDoom8BitParserManyFrame(t *testing.T) {
    data, err := os.ReadFile("./doomtest_large")
    require.NoError(t, err)

    doomHeader := ansiparser.NewDoomAsciiHeaderParser()
    n, err := doomHeader.Write(data)
    require.NoError(t, err)

    data = data[n:]
    rows, cols := doomHeader.GetDims()
    require.Equal(t, rows, 66)
    require.Equal(t, cols, 212)

    doomAscii := ansiparser.NewDoomAnsiFramer(rows, cols)

    go func() {
        doomAscii.Write(data)
    }()

    frames := doomAscii.Frames()
    count := 0
    length := rows * cols

    outer:
    for {
        timer := time.NewTimer(time.Millisecond * 50)
        select {
        case frame := <-frames:
            require.Equal(t, len(frame.Chars), length)
            require.Equal(t, len(frame.Color), length)
            count++
        case <-timer.C:
            break outer
        }
    }

    require.Equal(t, 129, count)
}
