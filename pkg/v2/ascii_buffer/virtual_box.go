package ascii_buffer

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

type VirtualBox struct {
	buffer []byte

	TotalRows int
	TotalCols int
	Cols      int
	Rows      int

	Row int
	Col int

	Stride int

	iterRow int
	iterCol int
	res     byteutils.ByteIteratorResult
}

type Quadtree []*VirtualBox
type QuadtreeParam struct {
	Depth  int
	Rows   int
	Cols   int
	Stride int
}

func (q *Quadtree) UpdateBuffer(buffer []byte) {
	for _, v := range []*VirtualBox(*q) {
		v.buffer = buffer
	}
}

func (v *VirtualBox) String() string {
	return fmt.Sprintf("VirtualBox(%d, %d) rows=%d cols=%d totalRows=%d totalCols=%d stride=%d (iter: x=%d y=%d res=%s)",
		v.Row, v.Col,
		v.Rows, v.Cols,
		v.TotalRows, v.TotalCols,
		v.Stride,
		v.iterRow, v.iterCol, &v.res,
	)
}

func fromVirtualBox(v *VirtualBox) *VirtualBox {
	return &VirtualBox{
		buffer:    v.buffer,
		Row:       0,
		Col:       0,
		TotalRows: v.TotalRows,
		TotalCols: v.TotalCols,
		Stride:    v.Stride,
		iterRow:   0,
		iterCol:   0,
		res:       byteutils.ByteIteratorResult{Done: false, Value: 0},
	}
}

func newVirtualBox(buf []byte, totalRows, totalCols int) *VirtualBox {
	assert.Assert(len(buf) == totalRows*totalCols, "the buf and rows * cols are not the same size")
	return &VirtualBox{
		buffer: buf,
		Row:    0,
		Col:    0,

		TotalRows: totalRows,
		Rows:      totalRows,

		TotalCols: totalCols,
		Cols:      totalCols,

		iterRow: 0,
		iterCol: 0,
		res:     byteutils.ByteIteratorResult{Done: false, Value: 0},
	}
}

func (v *VirtualBox) withSize(rows, cols int) *VirtualBox {
	v.Rows = rows
	v.Cols = cols
	return v
}

func (v *VirtualBox) Len() int {
	return v.Rows * v.Cols
}

func (v *VirtualBox) at(row, col int) *VirtualBox {
	v.Row = row
	v.Col = col
	return v
}

func (v *VirtualBox) withStride(size int) *VirtualBox {
	assert.Assert(size == 1 || size == 2, "virtual box only supports 8 or 16 bit value types")

	v.Stride = size
	return v
}

func (v *VirtualBox) quad() (*VirtualBox, *VirtualBox, *VirtualBox, *VirtualBox) {
	assert.Assert(v.Stride > 0, "you must set data size before you quad tree")
	assert.Assert(v.Rows >= 2, "cannot make a virtual box with 1 or less rows", "from", v)
	assert.Assert(v.Cols >= 2, "cannot make a virtual box with 1 or less rows", "from", v)

	rows := v.Rows / 2
	cols := v.Cols / 2
	cols -= cols % v.Stride

	centerCol := v.Col + cols
	centerRow := v.Row + rows

	maxRow := centerRow + (v.Rows - rows) - 1
	maxCol := centerCol + (v.Cols - cols) - 1

	maxIdx := Translate(maxRow, maxCol, v.TotalCols)
	assert.Assert(maxIdx < len(v.buffer),
		"translation of location is off the map",
		"idx", maxIdx,
		"CRow", centerRow,
		"CCol", centerCol,

		"Rows", v.Rows,
		"Cols", v.Cols,

		"midRow", rows,
		"midCol", cols,
		"strides", v.Stride,
		"totalCols", v.TotalCols,
		"maxRow", maxRow,
		"maxCol", maxCol,
	)

	tl := fromVirtualBox(v).
		at(v.Row, v.Col).
		withSize(rows, cols)

	tr := fromVirtualBox(v).
		at(v.Row, centerCol).
		withSize(rows, v.Cols-cols)

	bl := fromVirtualBox(v).
		at(centerRow, v.Col).
		withSize(v.Rows-rows, cols)

	br := fromVirtualBox(v).
		at(centerRow, centerCol).
		withSize(v.Rows-rows, v.Cols-cols)

	return tl, tr, bl, br
}

func (v *VirtualBox) Next() byteutils.ByteIteratorResult {
	assert.Assert(!v.res.Done, "you cannot call next when iterator is done", "res", v.res)

	row := v.Row + v.iterRow
	col := v.Col + v.iterCol
	idx := Translate(row, col, v.TotalCols)

	assert.Assert(idx >= 0 && idx < len(v.buffer), "idx cannot exceed the bounds, this means translate is broken", "idx", idx, "x", row, "y", col)
	value := int(v.buffer[idx])
	if v.Stride == 2 {
		value = byteutils.Read16(v.buffer, idx)
	}

	v.iterCol += v.Stride
	assert.Assert(v.iterRow <= v.Cols, "somehow iterX is greator than cols which means we have a mishaped virtual box", "box", v)

	if v.iterCol == v.Cols {
		v.iterCol = 0
		v.iterRow++
	}

	v.res.Done = v.iterRow == v.Rows
	v.res.Value = value

	return v.res
}

func (v *VirtualBox) Reset() {
	v.iterRow = 0
	v.iterCol = 0
	v.res.Done = false
}

func partition(data []byte, current *VirtualBox, boxes *[]*VirtualBox, depth int) {
	if depth == 0 {
		*boxes = append(*boxes, current)
		return
	}

	tl, tr, bl, br := current.quad()
	partition(data, tl, boxes, depth-1)
	partition(data, tr, boxes, depth-1)
	partition(data, bl, boxes, depth-1)
	partition(data, br, boxes, depth-1)

}

func Partition(data []byte, params QuadtreeParam) Quadtree {
	boxes := &([]*VirtualBox{})
	v := newVirtualBox(data, params.Rows, params.Cols).
		withStride(params.Stride)

	partition(data, v, boxes, params.Depth)

	return *boxes
}
