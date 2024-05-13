package encoding

import "github.com/theprimeagen/vim-with-me/pkg/v2/assert"

type rgbReader interface {
	read(buf []byte, offset int) (uint, int)
}

type rgbWriter interface {
	write(buf []byte, value uint) int
}

type IteratorResult struct {
	done  bool
	value uint
}

type RGBIterator struct {
	buffer []byte
	idx    int
	reader rgbReader
	ret    IteratorResult
}

var empty = make([]byte, 0)

func New8BitRGBIterator() *RGBIterator {
	return &RGBIterator{
		buffer: empty,
		idx:    0,
        ret: IteratorResult{done: true, value: 0},
		reader: newRGB8BitReader(),
	}
}

func New16BitRGBIterator() *RGBIterator {
	return &RGBIterator{
		buffer: empty,
		idx:    0,
        ret: IteratorResult{done: true, value: 0},
		reader: newRGB16BitReader(),
	}
}

func (i *RGBIterator) Set(buffer []byte) *RGBIterator {
	i.buffer = buffer
	i.idx = 0
    i.ret.done = true

	return i
}

func (i *RGBIterator) Next() IteratorResult {
	assert.Assert(i.ret.done, "iterator is done, you cannot call done")
	value, offset := i.reader.read(i.buffer, i.idx)
	i.idx = offset

    if offset == len(i.buffer) {
        i.ret.done = true
    }

    i.ret.value = value

    return i.ret
}
