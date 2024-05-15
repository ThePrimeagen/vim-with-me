package huffman

import (
	"errors"
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

type Huffman struct {
	DecodingTree  []byte
	EncodingTable map[int][]byte
}

func left(decoder []byte, idx int) int {
    assert.Assert(len(decoder) > idx + 5, "decoder length + idx is shorter than huffmanNode decode length")
    return byteutils.Read16(decoder, idx + 2);
}

func right(decoder []byte, idx int) int {
    assert.Assert(len(decoder) > idx + 5, "decoder length + idx is shorter than huffmanNode decode length")
    return byteutils.Read16(decoder, idx + 4);
}

func jump(decoder []byte, idx int, bit int) int {
    if bit == 1 {
        return right(decoder, idx)
    }
    return left(decoder, idx)
}

func value(decoder []byte, idx int) int {
    return byteutils.Read16(decoder, idx);
}

func isLeaf(decoder []byte, idx int) bool {
    assert.Assert(len(decoder) > idx + 5, "decoder length + idx is shorter than huffmanNode decode length")
    return byteutils.Read16(decoder, idx + 2) == 0 &&
        byteutils.Read16(decoder, idx + 4) == 0
}

var HuffmanEncodingExceededSize = errors.New("huffman encoding has exceeded the data array size")
var HuffmanDecodingExceededSize = errors.New("huffman decoding has exceeded the data array size")
var HuffmanUnknownValue = errors.New("the iterator produced a value that is not contained in the encoding table")
var HuffmanDecodingFailedToWrite = errors.New("failed to decode data")

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

// Will i even use a decoder?  i should write this in typescript
func (h *Huffman) Decode(data []byte, bitLength int, writer byteutils.ByteWriter) error {
    fmt.Printf("data=%d with bitLength=%d with calculatedBytes=%d\n", len(data), bitLength, bitLength / 8 + 1)
	assert.Assert(len(data) >= bitLength / 8 + 1, "you did not provide enough data")

	idx := 0
	decodeIdx := 0

    outer:
	for {
		for bitIdx := 7; bitIdx >= 0; bitIdx-- {
			bit := int((data[idx] >> bitIdx) & 0x1)
			bitLength--

            nextDecode := jump(h.DecodingTree, decodeIdx, bit)

            fmt.Printf("decode(%d, %d): writer=%d bit=%d decodeIdx=%d next=%d isLeaf=%v value=%d\n",
                bitLength,
                bitIdx,
                writer.Len(),
                bit,
                decodeIdx,
                nextDecode,
                isLeaf(h.DecodingTree, nextDecode),
                value(h.DecodingTree, nextDecode))

            decodeIdx = nextDecode

            if isLeaf(h.DecodingTree, decodeIdx) {

                if err := writer.Write(value(h.DecodingTree, decodeIdx)); err != nil {
                    return errors.Join(
                        HuffmanDecodingFailedToWrite,
                        err,
                        fmt.Errorf("failed to write at decodeIdx=%d with writer#Len=%d", decodeIdx, writer.Len()))
                }

                decodeIdx = 0
            }

            if bitLength == 0 {
                break outer
            }
		}

        idx++
	}

	return nil
}
