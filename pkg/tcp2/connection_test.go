package tcp2_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/theprimeagen/vim-with-me/pkg/tcp2"
)

func TestConnection(t *testing.T) {

    cmd := &tcp2.TCPCommand{
        Command: byte('t'),
        Data: []byte("69:420"),
    }

    bin, err := cmd.MarshalBinary()
    assert.NoError(t, err)

    r := bytes.NewReader(bin)
    w := bytes.NewBuffer(nil)

    conn := tcp2.Connection{
        Id: 0,
        FrameReader: tcp2.NewFrameReader(r),
        FrameWriter: tcp2.NewFrameWriter(w),
    }

    outCommand, err := conn.Next()
    assert.Equal(t, outCommand, cmd)

    err = conn.Write(outCommand)
    assert.NoError(t, err)
    assert.Equal(t, w.Bytes(), bin)
}

