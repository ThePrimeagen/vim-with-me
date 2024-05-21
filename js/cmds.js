export const types = {
	open: 0,
	brightnessToAscii: 1,
	frame: 2,
}

export const encodings = {
	NONE: 1,
	XOR_RLE: 2,
	HUFFMAN: 3,
	XOR_BIT_DIFF: 4, // Not implement, but i am horned up for it
	XOR_HUFFMAN: 5,
	HUFFMAN_QUADTREE: 6, // Maybe implement?
	XOR_RLE_QUADTREE: 7, // Maybe implement?
}
