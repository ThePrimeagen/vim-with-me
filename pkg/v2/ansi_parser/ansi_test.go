package ansiparser_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
)

func TestAnsi(t *testing.T) {
    data, err := os.ReadFile("./doom_frame_start")
    require.NoError(t, err)

    parser := ansiparser.New8BitFramer().WithDim(2, 160)
    parser.Write(data)

    timer := time.NewTimer(time.Millisecond * 10)
    select {
    case <-timer.C:
        require.Fail(t, "failed to get frame after 10 ms")
    case f := <-parser.Frames():
        require.Equal(t, f.Empty, 0, "expected framer to have no missing data pieces")
    }

}
