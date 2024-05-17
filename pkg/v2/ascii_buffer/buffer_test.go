package ascii_buffer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
)

func TestAsciiPixelTest(t *testing.T) {
    data := []byte("llllll;;llllllllllllllllIIllll>>llllllllll::llllll;;;;IIII;;")
    data = ascii_buffer.RemoveAsciiStyledPixels(data)

    require.Equal(t,
        []byte("lll;llllllllIll>lllll:lll;;II;"),
        data)
}

