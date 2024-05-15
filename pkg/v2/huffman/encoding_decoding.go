package huffman

import (
	"errors"
	"fmt"

	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

type Huffman struct {
	DecodingTree []byte
	EncodingTable   map[int][]byte
}

var HuffmanEncodingExceededSize = errors.New("huffman encoding has exceeded the data array size")
var HuffmanUnknownValue = errors.New("the iterator produced a value that is not contained in the encoding table")

func (h *Huffman) Encode(iterator byteutils.ByteIterator, out []byte) (int, error) {
    currentVal := byte(0)
    byteIdx := 7
    length := 0
    bitLength := 0
    dirty := false

    for {
        res := iterator.Next()

        encoding, ok := h.EncodingTable[res.Value]
        if !ok {
            return 0, HuffmanUnknownValue
        }

        for _, bit := range encoding {
            bitLength++
            dirty = true

            currentVal = currentVal | (bit << byte(byteIdx))
            byteIdx--

            if byteIdx < 0 {
                byteIdx = 7
                out[length] = currentVal

                dirty = false
                currentVal = byte(0)
                length++

                if length == len(out) {
                    return 0, HuffmanEncodingExceededSize
                }
            }
        }

        if res.Done {
            break
        }
    }

    if dirty {
        out[length] = currentVal
        length++
    }

    return bitLength, nil
}


