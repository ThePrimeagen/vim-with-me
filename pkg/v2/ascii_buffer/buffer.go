package ascii_buffer

import (
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

func assertBounds(a, b, out []byte) {
    assert.Assert(len(a) == len(b), "cannot diff slices of differing lengths")
    assert.Assert(len(a) <= len(out), "cannot diff into an out array smaller than in arrays")
}

func Xor(a, b, out []byte) {
    assertBounds(a, b, out)

    for i, aByte := range a {
        out[i] = aByte ^ b[i]
    }
}

func RemoveAsciiStyledPixels(data []byte) []byte {
    assert.Assert(len(data) & 1 == 0, "you cannot remove ascii styled pixels if the array is not even length")

    idx := 1
    doubleIdx := 2

    for ;doubleIdx < len(data); doubleIdx += 2 {
        data[idx] = data[doubleIdx]
        idx++
    }

    return data[:len(data) / 2]
}

