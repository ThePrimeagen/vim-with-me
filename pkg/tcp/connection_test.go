package tcp_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
)

func TestConnection(t *testing.T) {

    cmd := &tcp.TCPCommand{
        Command: byte('t'),
        Data: []byte("69:420"),
    }

    bin, err := cmd.MarshalBinary()
    assert.NoError(t, err)

    r := bytes.NewReader(bin)
    w := bytes.NewBuffer(nil)

    conn := tcp.Connection{
        Id: 0,
        FrameReader: tcp.NewFrameReader(r),
        FrameWriter: tcp.NewFrameWriter(w),
    }

    outCommand, err := conn.Next()
    assert.Equal(t, outCommand, cmd)

    err = conn.Write(outCommand)
    assert.NoError(t, err)
    assert.Equal(t, w.Bytes(), bin)
}

