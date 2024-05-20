package ascii_buffer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

func TestBufferRLE(t *testing.T) {
	assert.AddAssertData("current-test", "TestBufferRLE")
	rle := ascii_buffer.NewAsciiRLE()
	rle.Reset(make([]byte, 12, 12))

	expected := []byte{
		6, '{',
		1, 'a',
		1, 'b',
		1, 'c',
		1, 'd',
		6, '{',
	}

	rle.Write([]byte{
		'{', '{',
		'{', '{',
		'{', '{',
	})

	rle.Write([]byte{
		'a', 'b',
		'c', 'd',
	})

	rle.Write([]byte{
		'{', '{',
		'{', '{',
		'{',
	})

	rle.Write([]byte{
		'{',
	})

	rle.Finish()

	require.Equal(t, 12, rle.Length())
	require.Equal(t, expected, rle.Bytes())
}

func TestBufRLEWithLastDiff(t *testing.T) {
	assert.AddAssertData("current-test", "TestBufRLEWithLastDiff")
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

	rle := ascii_buffer.NewAsciiRLE()
	rle.Reset(make([]byte, 14, 14))
	expected := []byte{
		6, '{',
		1, 'a',
		1, 'b',
		1, 'c',
		1, 'd',
		5, '{',
		1, '}',
	}

	rle.Write(data)
	rle.Finish()
	require.Equal(t, 14, rle.Length())
	require.Equal(t, expected, rle.Bytes())
}

func TestMaximumSize(t *testing.T) {
	assert.AddAssertData("current-test", "TestMaximumSize")
	data := []byte{}
	for i := 0; i < 256; i++ {
		data = append(data, '{')
	}

	rle := ascii_buffer.NewAsciiRLE()
	rle.Reset(make([]byte, 4, 4))
	expected := []byte{
		255, '{',
		1, '{',
	}

	rle.Write(data)
	rle.Finish()

	require.Equal(t, 4, rle.Length())
	require.Equal(t, expected, rle.Bytes())
}
