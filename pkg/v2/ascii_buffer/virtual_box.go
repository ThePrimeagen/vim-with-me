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

	X int
	Y int

	Stride int

    iterX int
    iterY int
    res byteutils.ByteIteratorResult
}

func (v *VirtualBox) String() string {
    return fmt.Sprintf("VirtualBox(%d, %d) rows=%d cols=%d totalRows=%d totalCols=%d stride=%d (iter: x=%d y=%d res=%s)",
        v.X, v.Y,
        v.Rows, v.Cols,
        v.TotalRows, v.TotalCols,
        v.Stride,
        v.iterX, v.iterY, &v.res,
    )
}

func fromVirtualBox(v *VirtualBox) *VirtualBox {
	return &VirtualBox{
		buffer:    v.buffer,
		X:         0,
		Y:         0,
		TotalRows: v.TotalRows,
		TotalCols: v.TotalCols,
		Stride:  v.Stride,
        iterX: 0,
        iterY: 0,
        res: byteutils.ByteIteratorResult{Done: false, Value: 0},
	}
}

func newVirtualBox(buf []byte, totalRows, totalCols int) *VirtualBox {
	return &VirtualBox{
		buffer:    buf,
		X:         0,
		Y:         0,
		TotalRows: totalRows,
		TotalCols: totalCols,
        Rows: totalRows,
        Cols: totalCols,
        iterX: 0,
        iterY: 0,
        res: byteutils.ByteIteratorResult{Done: false, Value: 0},
	}
}

func (v *VirtualBox) withSize(rows, cols int) *VirtualBox {
	v.Rows = cols
	v.Cols = rows
	return v
}

func (v *VirtualBox) at(x, y int) *VirtualBox {
	v.X = x
	v.Y = y
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

	centerX := v.X + v.Cols/2
	centerY := v.Y + v.Rows/2
	centerX -= centerX % v.Stride

	rows := v.Rows / 2
	cols := v.Cols / 2
	cols -= cols % v.Stride

	tl := fromVirtualBox(v).
		at(v.X, v.Y).
		withSize(rows, cols)

	tr := fromVirtualBox(v).
		at(centerX, v.Y).
		withSize(rows, v.Cols-cols)

	bl := fromVirtualBox(v).
		at(v.X, centerY).
		withSize(v.Rows-rows, cols)

	br := fromVirtualBox(v).
		at(centerX, centerY).
		withSize(v.Rows-rows, v.Cols-cols)

	return tl, tr, bl, br
}

func (v *VirtualBox) Next() byteutils.ByteIteratorResult {
    assert.Assert(!v.res.Done, "you cannot call next when iterator is done", "res", v.res)

    x := v.X + v.iterX
    y := v.Y + v.iterY
    idx := Translate(x, y, v.TotalCols)

    assert.Assert(idx >= 0 && idx < len(v.buffer), "idx cannot exceed the bounds, this means translate is broken", "idx", idx)
    value := int(v.buffer[idx])
    if v.Stride == 2 {
        value = byteutils.Read16(v.buffer, idx)
    }

    v.iterX += v.Stride
    assert.Assert(v.iterX <= v.Cols, "somehow iterX is greator than cols which means we have a mishaped virtual box", "box", v)

    if v.iterX == v.Cols {
        v.iterX = 0
        v.iterY++
    }

    v.res.Done = v.iterY == v.Rows
    v.res.Value = value

    return v.res
}

func (v *VirtualBox) Reset() {
    v.iterX = 0
    v.iterY = 0
    v.res.Done = false
}

func partition(data []byte, current *VirtualBox, boxes *[]*VirtualBox, depth int) {
    if depth == 0 {
        *boxes = append(*boxes, current)
        return
    }

    tl, tr, bl, br := current.quad()
    partition(data, tl, boxes, depth - 1)
    partition(data, tr, boxes, depth - 1)
    partition(data, bl, boxes, depth - 1)
    partition(data, br, boxes, depth - 1)


}

// TODO: I hate this interface...
func Partition(data []byte, rows, cols, depth, stride int) []*VirtualBox {
	boxes := &([]*VirtualBox{})
    v := newVirtualBox(data, rows, cols).
        withStride(stride)

	partition(data, v, boxes, depth)

	return *boxes
}
