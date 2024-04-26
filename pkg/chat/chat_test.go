package chat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseChat(t *testing.T) {
    name := "foo"
    msgType := "message"
    msg := "t:0:0"

    chat, err := parse(fmt.Sprintf("%s:%s:%s", msgType, name, msg))
    require.NoError(t, err)
    require.Equal(t, &ChatMsg{
        Name: name,
        Msg: msg,
        Bits: 0,
    }, chat)
}

func TestParseBit(t *testing.T) {
    name := "foo"
    msgType := "bits"
    bits := "69"
    msg := "i like armoranth"

    chat, err := parse(fmt.Sprintf("%s:%s:%s:%s", msgType, name, bits, msg))
    require.NoError(t, err)
    require.Equal(t, &ChatMsg{
        Name: name,
        Msg: msg,
        Bits: 69,
    }, chat)
}

func TestBadMessage(t *testing.T) {
    name := "foo"
    msgType := "aoeu"
    msg := "i like piq more"

    chat, err := parse(fmt.Sprintf("%s:%s:%s", name, msgType, msg))
    require.Error(t, err)
    require.Nil(t, chat)
}

