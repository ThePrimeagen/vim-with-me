package chat_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	"github.com/theprimeagen/vim-with-me/pkg/v2/chat"
)

type ChatTest struct {
	input  []string
	answer string
}

func (c ChatTest) test(ch *chat.ChatAggregator, t *testing.T) {
	for _, str := range c.input {
		ch.Add(str)
	}

	occ := ch.Reset()
	require.Equal(t, c.answer, occ.Msg)
}

func TestChat(t *testing.T) {
	c := chat.
		NewChatAggregator().
		WithMap(doom.DoomChatMapFn).
		WithFilter(doom.DoomFilterFn)

	tests := []ChatTest{
		{input: []string{"w", "w", "w", "a", "a", "a", "f"}, answer: "w"},
		{input: []string{"w", "w", "w", "a", "a", "a", "f", "w"}, answer: "w"},
		{input: []string{"w", "w", "w", "a", "a", "a", "f", "a"}, answer: "a"},
		{input: []string{}, answer: ""},
		{input: []string{"aw", "w", "w", "wa", "wa", "w"}, answer: "aw"},
		{input: []string{"fw", "w", "w", "wf", "wf", "w"}, answer: "fw"},
	}

	for _, test := range tests {
		test.test(&c, t)
	}
}

func TestOneMessageChat(t *testing.T) {
	c := chat.
		NewChatAggregator().
		WithMap(doom.DoomChatMapFn).
		WithFilter(doom.DoomFilterFn)

    c.Add("w")
    occ := c.Reset()

    require.Equal(t, "w", occ.Msg)
    require.Equal(t, 1, occ.Count)
}
