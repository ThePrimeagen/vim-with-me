package encoder_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
	"github.com/theprimeagen/vim-with-me/pkg/v2/encoder"
	"github.com/theprimeagen/vim-with-me/pkg/v2/net"
)

func getEncoder() *encoder.Encoder {
	return encoder.NewEncoder(39, ascii_buffer.QuadtreeParam{
		Depth:  1,
		Stride: 1,
		Rows:   13,
		Cols:   3,
	})
}

func getTestData() []byte {
    return []byte{
		'A', 'A', 'A',
		'A', 'A', 'A',
		'A', 'A', 'A',
		'A', 'A', 'A',
		'A', 'A', 'A',
		'A', 'A', 'A',
		'A', 'A', 'A',
		'A', 'A', 'A',
		'A', 'A', 'A',
		'A', 'A', 'A',
		'A', 'A', 'A',
		'A', 'A', 'A',
		'C', 'B', 'B',
	}
}

func TestEncodeFrameXOR_RLE(t *testing.T) {
    data := getTestData()
    data2 := getTestData()
	enc := getEncoder().AddEncoder(encoder.XorRLE)

    encFrame := enc.PushFrame(data)
    require.Nil(t, encFrame, "encframe was nil")
    encFrame = enc.PushFrame(data2)
    require.NotNil(t, encFrame, "encframe was nil")

    require.Equal(t, 2, encFrame.Len)
    require.Equal(t, []byte{
        39, 0,
    }, encFrame.Out[:encFrame.Len])

    // write that buf
    out := make([]byte, 7, 7)
    frameable := net.Frameable{Item: encFrame}
    n, err := frameable.Into(out, 0)
    require.NoError(t, err)
    require.Equal(t, 7, n)
    require.Equal(t, []byte{
        net.VERSION,
        byte(net.FRAME),

        // length
        0,
        3, // 1 encoding, 2 XorRLE

        // encoding
        byte(encoder.XOR_RLE),

        // Data
        39, 0,
    }, out)
}

func TestEncodeFrameHuffman(t *testing.T) {
    data := getTestData()
	enc := getEncoder().AddEncoder(encoder.Huffman)

    encFrame := enc.PushFrame(data)
    require.NotNil(t, encFrame, "encframe was nil")

    huffLen := len(encFrame.Huff.DecodingTree) + (encFrame.HuffBitLen + 7) / 8

    require.Equal(t, huffLen, encFrame.Len)

    expectedOut := []byte{byte(encoder.HUFFMAN), 0, 0, 0, 0}
    byteutils.Write16(expectedOut, 1, encFrame.HuffBitLen)
    byteutils.Write16(expectedOut, 3, len(encFrame.Huff.DecodingTree))

    expectedOut = append(expectedOut, encFrame.Huff.DecodingTree...)
    expectedOut = append(expectedOut, 0b1111_1111) // 8 As
    expectedOut = append(expectedOut, 0b1111_1111) // 8 As
    expectedOut = append(expectedOut, 0b1111_1111) // 8 As
    expectedOut = append(expectedOut, 0b1111_1111) // 8 As
    expectedOut = append(expectedOut, 0b1111_0001) // 4 As, 1 C, 1 B
    expectedOut = append(expectedOut, 0b0100_0000) // 1B... 0

    // 1 for encoding
    // 4 for lengths of bits + decoding tree length
    // 6 for the huffed data
    // len of tree for tree data
    totalLen := 1 + 4 + 6 + len(encFrame.Huff.DecodingTree)
    out := make([]byte, totalLen, totalLen)
    n, err := encFrame.Into(out, 0)

    require.NoError(t, err)
    require.Equal(t, totalLen, n)
    require.Equal(t, expectedOut, out)
}
