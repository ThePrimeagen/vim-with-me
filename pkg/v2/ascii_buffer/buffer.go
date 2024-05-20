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
