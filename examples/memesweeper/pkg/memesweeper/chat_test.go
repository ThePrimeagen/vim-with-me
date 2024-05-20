package memesweeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
)

func TestParseChatMsg(t *testing.T) {
	row, col, err := ParseChatMessage("B2")
	require.NoError(t, err)

	require.Equal(t, 2, row)
	require.Equal(t, "B", col)

	row, col, err = ParseChatMessage("2B")
	require.NoError(t, err)

	require.Equal(t, 2, row)
	require.Equal(t, "B", col)
}

func TestChat(t *testing.T) {
	testies.SetupLogger()
	chat := NewChatAggregator()
	require.Equal(t, Point{row: 0, col: 0, count: 0}, chat.Current())

	chat.SetActiveState(false)
	chat.Add(1, 1)
	require.Equal(t, Point{row: 0, col: 0, count: 0}, chat.Current())

	chat.SetActiveState(true)
	chat.Add(1, 1)
	require.Equal(t, Point{row: 1, col: 1, count: 1}, chat.Current())
	chat.Add(1, 1)
	require.Equal(t, Point{row: 1, col: 1, count: 2}, chat.Current())

	point := chat.Reset()
	require.Equal(t, Point{row: 1, col: 1, count: 2}, point)

	chat.Add(1, 1)
	require.Equal(t, Point{row: 1, col: 1, count: 1}, chat.Current())
	chat.Add(1, 1)
	require.Equal(t, Point{row: 1, col: 1, count: 2}, chat.Current())
	chat.Add(2, 2)
	require.Equal(t, Point{row: 1, col: 1, count: 2}, chat.Current())
	chat.Add(2, 2)
	require.Equal(t, Point{row: 1, col: 1, count: 2}, chat.Current())
	chat.Add(2, 2)
	require.Equal(t, Point{row: 2, col: 2, count: 3}, chat.Current())

	point = chat.Reset()
	require.Equal(t, Point{row: 2, col: 2, count: 3}, point)

}
