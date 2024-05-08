package ascii_buffer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHuffman(t *testing.T) {
    freq := NewFreqency()
    freq.Freq([]byte{
        'A', 'A', 'A',
        'B', 'B',
        'C', 'D',
    })

    data, err := CalculateHuffman(freq)
    require.NoError(t, err)
    require.Equal(t, data, []byte{
        0, 3, 6,
        'A', 0, 0, // 0
        0, 9, 12,
        'B', 0, 0, // 10
        0, 15, 18,
        'D', 0, 0, // 110
        'C', 0, 0, // 111
    })
}

