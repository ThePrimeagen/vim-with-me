package encoder

import "github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"

type EncodingFrame struct {
	Prev []byte
	Curr []byte

	CurrQT ascii_buffer.Quadtree
	PrevQT ascii_buffer.Quadtree

	Out      []byte
	Len      int
	Encoding byte
}

func (e *EncodingFrame) pushFrame(frame []byte) {
	e.Prev = e.Curr
	e.Curr = frame

	e.CurrQT.UpdateBuffer(e.Curr)
	if e.Prev != nil {
		e.PrevQT.UpdateBuffer(e.Prev)
	}
}

func newEncodingFrame(size int, params ascii_buffer.QuadtreeParam) *EncodingFrame {
	out := make([]byte, size, size)
	prevQt := ascii_buffer.Partition(out, params)
	currQt := ascii_buffer.Partition(out, params)
	return &EncodingFrame{
		Prev:   nil,
		PrevQT: prevQt,

		Curr:   nil,
		CurrQT: currQt,

		Out: out,
		Len: 0,
	}
}

type EncodingCall func(frame *EncodingFrame) error

type Encoder struct {
	encodings []EncodingCall
	frames    []*EncodingFrame
	size      int
	params    ascii_buffer.QuadtreeParam
}

func NewEncoder(size int, treeParams ascii_buffer.QuadtreeParam) *Encoder {
	return &Encoder{
		encodings: make([]EncodingCall, 0),
		frames:    make([]*EncodingFrame, 0),
		size:      size,
		params:    treeParams,
	}
}

func (e *Encoder) AddEncoder(encoder EncodingCall) {
	e.encodings = append(e.encodings, encoder)
	e.frames = append(e.frames, newEncodingFrame(e.size, e.params))
}

func (e *Encoder) PushFrame(data []byte) (int, []byte) {
	min := len(data)
	minBytes := data

	for i, encoder := range e.encodings {
		frame := e.frames[i]
		frame.pushFrame(data)

		err := encoder(frame)
		if err != nil {
			continue
		}

		if min > frame.Len {
			min = frame.Len
			minBytes = frame.Out
		}
	}

	return min, minBytes
}
