package huffman

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

type Huffman struct {
	DecodingTree  []byte
	EncodingTable map[int][]byte
}

func left(decoder []byte, idx int) int {
	assert.Assert(len(decoder) > idx+5, "decoder length + idx is shorter than huffmanNode decode length")
	return byteutils.Read16(decoder, idx+2)
}

func right(decoder []byte, idx int) int {
	assert.Assert(len(decoder) > idx+5, "decoder length + idx is shorter than huffmanNode decode length")
	return byteutils.Read16(decoder, idx+4)
}

func jump(decoder []byte, idx int, bit int) int {
	if bit == 1 {
		return right(decoder, idx)
	}
	return left(decoder, idx)
}

func value(decoder []byte, idx int) int {
	return byteutils.Read16(decoder, idx)
}

func isLeaf(decoder []byte, idx int) bool {
	assert.Assert(len(decoder) > idx+5, "decoder length + idx is shorter than huffmanNode decode length")
	return byteutils.Read16(decoder, idx+2) == 0 &&
		byteutils.Read16(decoder, idx+4) == 0
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

func (h *Huffman) DebugDecodeTree() string {
	out := make([]string, 0)
	var walk func(curr []byte, length int, idx int)
	walk = func(curr []byte, length int, idx int) {
		if isLeaf(h.DecodingTree, idx) {
			str := ""
			for _, b := range curr {
				str += strconv.Itoa(int(b))
			}
			out = append(out, str+" "+fmt.Sprintf("0x%2x", value(h.DecodingTree, idx)))
			return
		}

		rightIdx := right(h.DecodingTree, idx)

		curr = append(curr, 0)
		leftIdx := left(h.DecodingTree, idx)
		walk(curr, length+1, leftIdx)

		curr[length] = 1
		walk(curr, length+1, rightIdx)
	}

	walk([]byte{}, 0, 0)

	return strings.Join(out, "\n")
}

// Will i even use a decoder?  i should write this in typescript
func (h *Huffman) Decode(data []byte, bitLength int, writer byteutils.ByteWriter) error {
	assert.Assert(len(data) >= (bitLength+7)/8, "you did not provide enough data")

	idx := 0
	decodeIdx := 0
	decodeIdxCount := 0
	debugArr := make([]int, 20, 20)

outer:
	for {
		for bitIdx := 7; bitIdx >= 0; bitIdx-- {
			bit := int((data[idx] >> bitIdx) & 0x1)
			bitLength--

			decodeIdx = jump(h.DecodingTree, decodeIdx, bit)
			debugArr[decodeIdxCount] = bit
			decodeIdxCount++

			if isLeaf(h.DecodingTree, decodeIdx) {

				v := value(h.DecodingTree, decodeIdx)
				str := ""
				for _, b := range debugArr[:decodeIdxCount] {
					str += strconv.Itoa(int(b))
				}

				if err := writer.Write(v); err != nil {
					return errors.Join(
						HuffmanDecodingFailedToWrite,
						err,
						fmt.Errorf("failed to write at decodeIdx=%d with writer#Len=%d", decodeIdx, writer.Len()))
				}

				decodeIdx = 0
				decodeIdxCount = 0
			}

			if bitLength == 0 {
				assert.Assert(decodeIdx == 0, "finished decoding with hanging state", "decodeIdx", decodeIdx)
				break outer
			}
		}

		idx++
	}

	return nil
}

var HuffmanDoesntFitIntoOutArray = errors.New("unable to fit huffman into output array")

func IntoBytes(huff *Huffman, bitLen int, data []byte, offset int) int {
	assert.Assert(bitLen < int(math.Pow(2.0, 16)), "unable to encode huffman frame larger than 65535 bits")

	fit := 4+len(huff.DecodingTree) < len(data)-offset
	assert.Assert(fit, "huffman tree is unable to fit into provided data array")

	byteutils.Write16(data, offset, bitLen)
	byteutils.Write16(data, offset+2, len(huff.DecodingTree))

	slog.Warn("huffman IntoBytes", "bitLen", bitLen, "decodeTree", len(huff.DecodingTree))
	copy(data[offset+4:], huff.DecodingTree)

	return 4 + len(huff.DecodingTree)
}
