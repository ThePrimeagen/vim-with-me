package encoder

type EncoderType byte

const (
	NONE EncoderType = iota
	XOR_RLE
	XOR_BIT_DIFF // Not implement, but i am horned up for it
	HUFFMAN
	XOR_HUFFMAN      // not implemented
	RLE
	HUFFMAN_QUADTREE // Maybe implement?
	XOR_RLE_QUADTREE // Maybe implement?
)
