package byteutils

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

type ByteIteratorResult struct {
	Done  bool
	Value int
}

func (b *ByteIteratorResult) String() string {
	return fmt.Sprintf("value=%d done=%v", b.Value, b.Done)
}

type ByteIterator interface {
	Next() ByteIteratorResult
}

type SixteenBitIterator struct {
	buffer []byte
	idx    int
	res    ByteIteratorResult
}

func New16BitIterator(buf []byte) *SixteenBitIterator {
	assert.Assert(len(buf)&0x1 == 0, "you must pass in an even number len buf for 16 bit iterators")
	return &SixteenBitIterator{
		buffer: buf,
		idx:    0,
		res:    ByteIteratorResult{Done: false, Value: 0},
	}
}

func Read16(buf []byte, offset int) int {
	assert.Assert(len(buf) > offset+1, "you cannot read outside of the buffer")

	hi := int(buf[offset])
	lo := int(buf[offset+1])
	return hi<<8 | lo
}

func (i *SixteenBitIterator) Next() ByteIteratorResult {
	assert.Assert(!i.res.Done, "SixteenBitIterator#Next was called on an exhausted iterator")

	value := Read16(i.buffer, i.idx)
	i.idx += 2

	i.res.Done = i.idx == len(i.buffer)
	i.res.Value = value

	return i.res
}

type EightBitIterator struct {
	buffer []byte
	idx    int
	res    ByteIteratorResult
}

func New8BitIterator(buf []byte) *EightBitIterator {
	return &EightBitIterator{
		buffer: buf,
		idx:    0,
		res:    ByteIteratorResult{Done: false, Value: 0},
	}
}

func (i *EightBitIterator) Next() ByteIteratorResult {
	assert.Assert(!i.res.Done, "EightBitIterator#Next was called on an exhausted iterator")
	val := i.buffer[i.idx]
	i.idx++

	i.res.Done = i.idx == len(i.buffer)
	i.res.Value = int(val)

	return i.res
}
