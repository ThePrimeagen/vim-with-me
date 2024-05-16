package ascii_buffer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

func TestVirtualBox(t *testing.T) {
    data := []byte{
        1, 2, 3, 4,
        5, 6, 7, 8,
        9, 10, 11, 12,
        13, 14, 15, 16,
    }

    expectedDataOrder := []byte{
        1, 2, 5, 6,
        3, 4, 7, 8,
        9, 10, 13, 14,
        11, 12, 15, 16,
    }

    boxes := ascii_buffer.Partition(data, 4, 4, 2, 1)
    require.Equal(t, 16, len(boxes))

    for i, b := range boxes {
        idx := ascii_buffer.Translate(b.X, b.Y, b.TotalCols)
        require.Equal(t, i + 1, int(expectedDataOrder[idx]))
    }
}

func TestVirtualBoxIterators(t *testing.T) {
    data := []byte{
        1, 2, 3, 4,
        5, 6, 7, 8,
        9, 10, 11, 12,
        13, 14, 15, 16,
    }

    expectedDataOrder := [][]int{
        { 1, 2, 5, 6, },
        { 3, 4, 7, 8, },
        { 9, 10, 13, 14, },
        { 11, 12, 15, 16, },
    }

    boxes := ascii_buffer.Partition(data, 4, 4, 1, 1)
    require.Equal(t, 4, len(boxes))

    for i, b := range boxes {
        idx := ascii_buffer.Translate(b.X, b.Y, b.TotalCols)
        expected := expectedDataOrder[i]

        // TODO: look into this
        // for res := b.Next()
        var res byteutils.ByteIteratorResult
        for j := range 4 {
            res = b.Next()
            require.Equal(t, expected[j], res.Value)

            if res.Done {
                break
            }
            idx++
        }
        require.Equal(t, true, res.Done)
    }
}

func TestVirtualBoxStride(t *testing.T) {

    data := []byte{
        1, 2, 3, 4,
        5, 6, 7, 8,
        9, 10, 11, 12,
        13, 14, 15, 16,
    }

    expectedDataOrder := [][]int{
        { 0x0102, 0x506, },
        { 0x0304, 0x0708, },
        { 0x090a, 0x0d0e, },
        { 0x0b0c, 0x0f10, },
    }

    boxes := ascii_buffer.Partition(data, 4, 4, 1, 2)
    require.Equal(t, 4, len(boxes))

    for i, b := range boxes {
        idx := ascii_buffer.Translate(b.X, b.Y, b.TotalCols)
        expected := expectedDataOrder[i]

        // TODO: look into this
        // for res := b.Next()
        var res byteutils.ByteIteratorResult
        for j := range 2 {
            res = b.Next()
            require.Equal(t, expected[j], res.Value)

            if res.Done {
                break
            }
            idx++
        }
        require.Equal(t, true, res.Done)
    }
}
