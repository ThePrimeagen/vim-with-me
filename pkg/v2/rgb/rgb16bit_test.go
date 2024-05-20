package rgb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRGB16Bit(t *testing.T) {
	rgb := newRGB16BitReader()
	buffer := []byte{
		0x0, 0x45,
		0x45, 0x45,
		0x45, 0x0,
	}

	val, nextIdx := rgb.read(buffer, 0)
	require.Equal(t, val, 0x45)
	require.Equal(t, nextIdx, 2)
	val, nextIdx = rgb.read(buffer, nextIdx)
	require.Equal(t, val, 0x4545)
	require.Equal(t, nextIdx, 4)
	val, nextIdx = rgb.read(buffer, nextIdx)
	require.Equal(t, val, 0x4500)
	require.Equal(t, nextIdx, 6)

}
