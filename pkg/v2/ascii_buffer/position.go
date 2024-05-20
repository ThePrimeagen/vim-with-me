package ascii_buffer

func Translate(row, col, cols int) int {
	return row*cols + col
}
