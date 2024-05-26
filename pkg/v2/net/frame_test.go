package net_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/v2/net"
)

func quickRead(framer *net.ByteFramer) *net.Frame {
	select {
	case f := <-framer.Frames():
		return f
	case <-time.After(time.Millisecond * 10):
	}
	return nil
}

func TestFramer(t *testing.T) {
	framer := net.NewByteFramer()
	ch := make(chan []byte, 11)
	go framer.FrameChan(ch)

	ch <- []byte{net.VERSION}
	require.Nil(t, quickRead(framer))
	ch <- []byte{3}
	require.Nil(t, quickRead(framer))
	ch <- []byte{0b00001010}
	require.Nil(t, quickRead(framer))
	ch <- []byte{0} // flags
	require.Nil(t, quickRead(framer))
	ch <- []byte{0x00, 0x03} // length 3
	require.Nil(t, quickRead(framer))
	ch <- []byte{0x01}
	require.Nil(t, quickRead(framer))
	ch <- []byte{0x02}
	require.Nil(t, quickRead(framer))
	ch <- []byte{0x03}
	require.Equal(t, &net.Frame{
		Seq:     0b1010,
		Flags:   0b0000,
		Data:    []byte{0x01, 0x02, 0x03},
		CmdType: 3,
	}, quickRead(framer))
}

func TestFramerEncode(t *testing.T) {
	cmd := byte(3)
	seq := byte(0b1010)
	flags := byte(0b0101)
	frame := &net.Frame{
		Seq:     seq,
		Flags:   flags,
		Data:    []byte{0x01, 0x02, 0x03},
		CmdType: cmd,
	}

	offset := 2
	out := make([]byte, net.HEADER_SIZE+3+offset)
	n, err := frame.Into(out, offset)

	require.Equal(t, net.HEADER_SIZE+3, n)
	require.NoError(t, err)
	require.Equal(t, []byte{
		0, 0, // offset
		net.VERSION,
		cmd,
		seq,
        0,
		0, 0x03,
		0x01, 0x02, 0x03,
	}, out)
}
