package encoder

import (
	"errors"

	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
	"github.com/theprimeagen/vim-with-me/pkg/v2/huffman"
)

var EncoderExceededBufferSize = errors.New("encoder exceeded buffer size of out buffer, useless encoding")

func NoneEncoding(frame *EncodingFrame) error {
	frame.Len = len(frame.Curr)
	frame.Encoding = NONE

	return nil
}

func RLEEncoding(frame *EncodingFrame) error {
	frame.RLE.Reset(frame.Out)
	frame.RLE.Write(frame.Curr)
	frame.RLE.Finish()
	frame.Len = frame.RLE.Length()
	frame.Encoding = RLE

	return nil
}

func XorRLE(frame *EncodingFrame) error {
	if frame.Prev == nil {
		frame.Len = len(frame.Out) + 1
		return nil
	}

	ascii_buffer.Xor(frame.Curr, frame.Prev, frame.Tmp)
	frame.RLE.Reset(frame.Out)
	frame.RLE.Write(frame.Tmp)
	frame.RLE.Finish()
	frame.Len = frame.RLE.Length()
	frame.Encoding = XOR_RLE

	return nil
}

func createIterator(frame *EncodingFrame) byteutils.ByteIterator {
	var iter byteutils.ByteIterator
	iter = byteutils.New8BitIterator(frame.Curr)
	if frame.Stride == 2 {
		iter = byteutils.New16BitIterator(frame.Curr)
	}

	return iter
}

func Huffman(frame *EncodingFrame) error {

	frame.Freq.Reset()

	frame.Freq.Freq(createIterator(frame))
	huff := huffman.CalculateHuffman(frame.Freq)

	huffLen := len(huff.DecodingTree)
	bitLen, err := huff.Encode(createIterator(frame), frame.Tmp)
	frame.TmpLen = (bitLen + 7) / 8

	if err != nil {
		return err
	}
	byteLen := bitLen/8 + 1

	if huffLen+byteLen >= len(frame.Curr) {
		frame.Len = len(frame.Curr) + 1
	}

	frame.Len = huffLen + byteLen
	frame.Encoding = HUFFMAN
	frame.HuffBitLen = bitLen
	frame.Huff = huff

	return nil
}

type encoderEncodingFn func(e *EncodingFrame, data []byte, offset int) (int, error)
type encoderEncodingMap map[EncoderType]encoderEncodingFn

var encodeInto encoderEncodingMap = encoderEncodingMap{
	HUFFMAN: func(e *EncodingFrame, data []byte, offset int) (int, error) {

		assert.Assert(e.Huff != nil, "the encoding type is huffman but the huff object in nil")

		decodeTreeLength := len(e.Huff.DecodingTree)
		decodeLen := 4 + decodeTreeLength + e.TmpLen

		assert.Assert(decodeLen < len(data), "unable to encode frame into provided buffer")

		byteutils.Write16(data, offset, e.HuffBitLen)
		byteutils.Write16(data, offset+2, decodeTreeLength)

		copy(data[offset+4:], e.Huff.DecodingTree)
		copy(data[offset+4+decodeTreeLength:], e.Tmp[:e.TmpLen])

		return decodeLen, nil
	},

	XOR_RLE: func(e *EncodingFrame, data []byte, offset int) (int, error) {
		assert.Assert(e.Len < len(data), "unable to encode frame into provided buffer")
		copy(data[offset:], e.Out[:e.Len])
		return e.Len, nil
	},

	RLE: func(e *EncodingFrame, data []byte, offset int) (int, error) {
		assert.Assert(e.Len < len(data), "unable to encode frame into provided buffer")
		copy(data[offset:], e.Out[:e.Len])
		return e.Len, nil
	},
}

func EncodingName(enc EncoderType) string {
    switch enc {
    case HUFFMAN:
        return "Huffman"
    case XOR_RLE:
        return "XorRLE"
    default:
        assert.Assert(false, "unable to determine encoding type", "enc", enc)
    }
    return "unreachable"
}
