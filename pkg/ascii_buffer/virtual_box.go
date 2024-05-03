package ascii_buffer

type AsciiVirtualBox struct {
	outerRows int
	outerCols int

	rows    int
	cols    int
	offsetR int
	offsetC int

    buffer []byte
}

func NewAsciiVirtualBox(buffer []byte) *AsciiVirtualBox {
    return &AsciiVirtualBox{buffer: buffer}
}

func (a *AsciiVirtualBox) WithSize(rows, cols int) *AsciiVirtualBox {
    a.rows = rows
    a.cols = cols
    return a
}

func (a *AsciiVirtualBox) WithOffset(offsetR, offsetC int) *AsciiVirtualBox {
    a.offsetR = offsetR
    a.offsetC = offsetC
    return a
}

func (a *AsciiVirtualBox) WithTotalSize(outerRows, outerCols int) *AsciiVirtualBox {
    a.outerRows = outerRows
    a.outerCols = outerCols
    return a
}

