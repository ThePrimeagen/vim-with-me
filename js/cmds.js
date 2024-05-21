export const types = {
	open: 0,
	brightnessToAscii: 1,
	frame: 2,
}

export const encodings = {
	NONE: 0,
	XOR_RLE: 1,
    XOR_BIT_DIFF: 2, // NOT USED
	HUFFMAN: 3,
	XOR_HUFFMAN: 5, // NOT USED
	HUFFMAN_QUADTREE: 6, // Maybe implement?
	XOR_RLE_QUADTREE: 7, // Maybe implement?
}
