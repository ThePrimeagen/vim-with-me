package ascii_buffer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/ascii_buffer"
)

func TestBufferRLE(t *testing.T) {
    buf := ascii_buffer.NewAsciiRLE()
    data := []byte{
        0, '{',
        0, '{',
        0, '{',
        0, '}',
        0, '{',
        1, '{',
        1, '{',
        1, '{',
    }
    expected := []byte{
        3, 0, '{',
        1, 0, '}',
        1, 0, '{',
        3, 1, '{',
    }

    size, err := buf.RLE(data)
    require.NoError(t, err)
    require.Equal(t, 12, size)
    require.Equal(t, expected, buf.Bytes())
}

func TestBufferRLEToBig(t *testing.T) {
    buf := ascii_buffer.NewAsciiRLE()
    data := []byte{
        0, '{',
        0, '}',
        0, '{',
        1, '{',
    }


    size, err := buf.RLE(data)
    require.Error(t, err)
    require.Equal(t, 0, size)
    require.Equal(t, []byte(nil), buf.Bytes())
}
