package huffman_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
	"github.com/theprimeagen/vim-with-me/pkg/v2/huffman"
)

func getFreq() ascii_buffer.Frequency {
	freq := ascii_buffer.NewFreqency()
	freq.Freq(byteutils.New8BitIterator([]byte{
		'A', 'A', 'A',
		'B', 'B',
		'C', 'D',
	}))
	return freq
}

func TestHuffman(t *testing.T) {
	freq := getFreq()
	encodeLen := byte(huffman.HUFFMAN_ENCODE_LENGTH)
	data := huffman.CalculateHuffman(freq)
	require.Equal(t, []byte{
		0, 0, 0, encodeLen, 0, encodeLen * 2,
		0, 'A', 0, 0, 0, 0, // 0
		0, 0, 0, encodeLen * 3, 0, encodeLen * 4,
		0, 'B', 0, 0, 0, 0, // 10
		0, 0, 0, encodeLen * 5, 0, encodeLen * 6,
		0, 'D', 0, 0, 0, 0, // 110
		0, 'C', 0, 0, 0, 0, // 111
	}, data.DecodingTree)
}

func TestHuffmanTable(t *testing.T) {
	freq := getFreq()

	data := huffman.CalculateHuffman(freq)
	require.Equal(t, []byte{0}, data.EncodingTable['A'])
	require.Equal(t, []byte{1, 0}, data.EncodingTable['B'])
	require.Equal(t, []byte{1, 1, 0}, data.EncodingTable['D'])
	require.Equal(t, []byte{1, 1, 1}, data.EncodingTable['C'])
}

func TestHuffmanEncodingByteBoundary(t *testing.T) {
	freq := getFreq()
	huff := huffman.CalculateHuffman(freq)

	data := make([]byte, 2, 2)

	bitLength, err := huff.Encode(byteutils.New8BitIterator([]byte{
		'A', 'A', 'A', // 3 bits of 0
		'A', 'A', 'A', // 3 bits of 0
		'C', // 3 bits of 111 SHOULD CROSS BYTE BOUNDARY
	}), data)
	require.NoError(t, err)
	require.Equal(t, 9, bitLength)
	require.Equal(t, data, []byte{0x03, 0x80})
}

func TestHuffmanEncodingErrTooLarge(t *testing.T) {
	freq := getFreq()
	huff := huffman.CalculateHuffman(freq)

	data := make([]byte, 1, 1)

	_, err := huff.Encode(byteutils.New8BitIterator([]byte{
		'A', 'A', 'A', // 3 bits of 0
		'A', 'A', 'A', // 3 bits of 0
		'C', // 3 bits of 111 SHOULD CROSS BYTE BOUNDARY
	}), data)
	require.ErrorIs(t, err, huffman.HuffmanEncodingExceededSize)
}

func TestHuffmanEncodingErrUnknownChar(t *testing.T) {
	freq := getFreq()
	huff := huffman.CalculateHuffman(freq)

	data := make([]byte, 10, 10)

	_, err := huff.Encode(byteutils.New8BitIterator([]byte{
		'E',
	}), data)

	require.ErrorIs(t, err, huffman.HuffmanUnknownValue)
}

// TODO: Do i want to merge encoding and decoding into a single test?
func TestHuffmanDecoding(t *testing.T) {
	freq := getFreq()
	huff := huffman.CalculateHuffman(freq)

	data := make([]byte, 2, 2)

	bitLength, err := huff.Encode(byteutils.New8BitIterator([]byte{
		'A', 'A', 'A', // 3 bits of 0
		'A', 'A', 'A', // 3 bits of 0
		'C', // 3 bits of 111 SHOULD CROSS BYTE BOUNDARY
	}), data)
	require.NoError(t, err)
	require.Equal(t, 9, bitLength)
	require.Equal(t, data, []byte{0x03, 0x80})

	out := make([]byte, 7, 7)
	writer := byteutils.U8Writer{}
	writer.Set(out)

	err = huff.Decode(data, bitLength, &writer)
	require.NoError(t, err)

	require.Equal(t, []byte{
		'A', 'A', 'A', // 3 bits of 0
		'A', 'A', 'A', // 3 bits of 0
		'C', // 3 bits of 111 SHOULD CROSS BYTE BOUNDARY
	}, out)
}

func TestHuffmanDecodingBufferTooSmall(t *testing.T) {
	freq := getFreq()
	huff := huffman.CalculateHuffman(freq)

	data := make([]byte, 2, 2)

	bitLength, err := huff.Encode(byteutils.New8BitIterator([]byte{
		'A', 'A', 'A', // 3 bits of 0
		'A', 'A', 'A', // 3 bits of 0
		'C', // 3 bits of 111 SHOULD CROSS BYTE BOUNDARY
	}), data)
	require.NoError(t, err)
	require.Equal(t, 9, bitLength)
	require.Equal(t, data, []byte{0x03, 0x80})

	out := make([]byte, 4, 4)
	writer := byteutils.U8Writer{}
	writer.Set(out)

	err = huff.Decode(data, bitLength, &writer)
	require.ErrorContains(t, err, huffman.HuffmanDecodingFailedToWrite.Error())
}
