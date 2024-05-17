package encoder

import (
	"errors"

	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
	"github.com/theprimeagen/vim-with-me/pkg/v2/huffman"
)

var EncoderExceededBufferSize = errors.New("encoder exceeded buffer size of out buffer, useless encoding")

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
    bitLen, err := huff.Encode(createIterator(frame), frame.Out)
    if err != nil {
        return err
    }
    byteLen := bitLen / 8 + 1

    if huffLen + byteLen >= len(frame.Curr) {
        frame.Len = len(frame.Curr) + 1
    }

    frame.Len = huffLen + byteLen
    frame.Encoding = HUFFMAN

    // TODO: Bit length is important...?

    return nil
}
