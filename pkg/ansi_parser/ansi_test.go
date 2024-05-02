package ansiparser_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/ansi_parser"
)

func TestDoom8BitParserOneFrame(t *testing.T) {
    doomAscii := ansiparser.NewDoomAnsiFramer()
    data, err := os.ReadFile("./doomtest")
    require.NoError(t, err)

    doomAscii.Write(data)

    require.Equal(t, 50, doomAscii.FoundRows)
    require.Equal(t, 160, doomAscii.FoundCols)

    frames := doomAscii.Frames()

    <-frames
}

func TestDoom8BitParserManyFrame(t *testing.T) {
    doomAscii := ansiparser.NewDoomAnsiFramer()
    data, err := os.ReadFile("./doomtest_large")
    require.NoError(t, err)

    doomAscii.Write(data[:5000])

    require.Equal(t, 66, doomAscii.FoundRows)
    require.Equal(t, 212, doomAscii.FoundCols)

    go func() {
        doomAscii.Write(data[5000:])
    }()
    frames := doomAscii.Frames()

    count := 0
    outer:
    for {
        timer := time.NewTimer(time.Second * 1)
        select {
        case <-frames:
            count++
        case <-timer.C:
            break outer
        }
    }

    require.Equal(t, 129, count)
}
