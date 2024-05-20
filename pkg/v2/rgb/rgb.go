package rgb

import (
	"github.com/leaanthony/go-ansi-parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

type rgbReader interface {
	read(buf []byte, offset int) (int, int)
}

type rgbWriter interface {
	write(buffer []byte, offset int, color *ansi.Rgb) int
	byteLength() int
}

type RGBIterator struct {
	buffer []byte
	idx    int
	reader rgbReader
	ret    byteutils.ByteIteratorResult
}

var empty = make([]byte, 0)

func New8BitRGBIterator() *RGBIterator {
	return &RGBIterator{
		buffer: empty,
		idx:    0,
		ret:    byteutils.ByteIteratorResult{Done: true, Value: 0},
		reader: newRGB8Bit(),
	}
}

func New16BitRGBIterator() *RGBIterator {
	return &RGBIterator{
		buffer: empty,
		idx:    0,
		ret:    byteutils.ByteIteratorResult{Done: true, Value: 0},
		reader: newRGB16BitReader(),
	}
}

func (i *RGBIterator) Set(buffer []byte) *RGBIterator {
	i.buffer = buffer
	i.idx = 0
	i.ret.Done = true

	return i
}

func (i *RGBIterator) Next() byteutils.ByteIteratorResult {
	assert.Assert(i.ret.Done, "iterator is done, you cannot call done")
	value, offset := i.reader.read(i.buffer, i.idx)
	i.idx = offset

	if offset == len(i.buffer) {
		i.ret.Done = true
	}

	i.ret.Value = value

	return i.ret
}

type RGBWriter struct {
	buffer []byte
	idx    int
	writer rgbWriter
}

func New8BitRGBWriter() *RGBWriter {
	return &RGBWriter{
		buffer: empty,
		idx:    0,
		writer: newRGB8Bit(),
	}
}

func (w *RGBWriter) Set(buffer []byte) {
	w.buffer = buffer
	w.idx = 0
}

func (w *RGBWriter) ByteLength() int {
	return w.writer.byteLength()
}

func (w *RGBWriter) Write(color *ansi.Rgb) int {
	assert.Assert(w.idx+w.writer.byteLength()-1 < len(w.buffer), "unable to write byte into rgb writer, buffer full")
	w.idx = w.writer.write(w.buffer, w.idx, color)

	return w.idx
}

func (w *RGBWriter) Full() bool {
	return w.idx == len(w.buffer)
}
