package chat

import (
	"testing"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/stretchr/testify/require"
)

func toPrivateMessage(name, msg string, bits int) twitch.PrivateMessage {
	return twitch.PrivateMessage{
		User:    twitch.User{DisplayName: name},
		Bits:    bits,
		Message: msg,
	}
}

func TestParseChat(t *testing.T) {
	name := "foo"
	msg := "t:0:0"

	chat := toChatMsg(toPrivateMessage(name, msg, 0))

	require.Equal(t, ChatMsg{
		Name: name,
		Msg:  msg,
		Bits: 0,
	}, chat)
}

func TestParseBit(t *testing.T) {
	name := "foo"
	msg := "i like armoranth"

	chat := toChatMsg(toPrivateMessage(name, msg, 69))
	require.Equal(t, ChatMsg{
		Name: name,
		Msg:  msg,
		Bits: 69,
	}, chat)
}
