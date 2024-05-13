package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRGB8Bit(t *testing.T) {
    rgb := newRGB8BitReader()
    buffer := []byte{
        0x0, 0x45,
        0x45, 0x0,
    }

    val, nextIdx := rgb.read(buffer, 0)
    require.Equal(t, val, uint(0x0))
    require.Equal(t, nextIdx, 1)
    val, nextIdx = rgb.read(buffer, nextIdx)
    require.Equal(t, val, uint(0x45))
    require.Equal(t, nextIdx, 2)
    val, nextIdx = rgb.read(buffer, nextIdx)
    require.Equal(t, val, uint(0x45))
    require.Equal(t, nextIdx, 3)
    val, nextIdx = rgb.read(buffer, nextIdx)
    require.Equal(t, val, uint(0x0))
    require.Equal(t, nextIdx, 4)
}


