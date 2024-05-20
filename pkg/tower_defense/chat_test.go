package tower_defense

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func makeChat(agg *ChatAggregator, row, col, count int) {
	for i := 0; i < count; i++ {
		agg.Add(row, col)
	}
}

func TestChatAgg(t *testing.T) {
	agg := NewChatAggregator()
	makeChat(&agg, 69, 420, 4)
	makeChat(&agg, 1, 1, 5)
	makeChat(&agg, 6, 9, 6)

	r, c := agg.Reset()
	require.Equal(t, 6, r)
	require.Equal(t, 9, c)

	makeChat(&agg, 69, 420, 6)
	makeChat(&agg, 1, 1, 5)
	makeChat(&agg, 6, 9, 4)

	r, c = agg.Reset()
	require.Equal(t, 69, r)
	require.Equal(t, 420, c)
}
