package ascii_buffer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
)

func TestBufferRLE(t *testing.T) {
    frame := ascii_buffer.NewAsciiFrame(8, 2)
    data := []byte{
        '{', '{',
        '{', '{',
        '{', '{',
        'a', 'b',
        'c', 'd',
        '{', '{',
        '{', '{',
        '{', '{',
    }
    frame.PushFrame(data)

    rle := ascii_buffer.NewAsciiRLE(frame.Length())
    expected := []byte{
        6, '{',
        1, 'a',
        1, 'b',
        1, 'c',
        1, 'd',
        6, '{',
    }

    err := rle.RLE(frame)
    require.NoError(t, err)

    require.Equal(t, 12, rle.Length())
    require.Equal(t, expected, rle.Bytes())
}

func TestBufRLEWithLastDiff(t *testing.T) {
    frame := ascii_buffer.NewAsciiFrame(8, 2)
    data := []byte{
        '{', '{',
        '{', '{',
        '{', '{',
        'a', 'b',
        'c', 'd',
        '{', '{',
        '{', '{',
        '{', '}',
    }
    frame.PushFrame(data)

    rle := ascii_buffer.NewAsciiRLE(frame.Length())
    expected := []byte{
        6, '{',
        1, 'a',
        1, 'b',
        1, 'c',
        1, 'd',
        5, '{',
        1, '}',
    }

    err := rle.RLE(frame)
    require.NoError(t, err)

    require.Equal(t, 14, rle.Length())
    require.Equal(t, expected, rle.Bytes())
}


func TestBufRLEWithError(t *testing.T) {
    frame := ascii_buffer.NewAsciiFrame(2, 2)
    data := []byte{
        'a', 'b',
        'c', 'd',
    }
    frame.PushFrame(data)

    rle := ascii_buffer.NewAsciiRLE(frame.Length())

    err := rle.RLE(frame)
    require.Error(t, err)
}
