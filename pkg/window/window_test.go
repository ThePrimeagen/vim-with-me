package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/theprimeagen/vim-with-me/pkg/commands"
)

func TestWindowChangeSetOnly(t *testing.T) {
    w := NewWindow(10, 10)

    w.Set(0, 0, ' ')
    flushed := w.Flush()

    if len(flushed) != 0 {
        t.Errorf("Expected no changes, but got %d", len(flushed))
    }

    w.Set(0, 0, 'a')
    flushed = w.Flush()

    if len(flushed) != 1 {
        t.Errorf("Expected 1 change, but got %d", len(flushed))
    }

    assert.Equal(t, []commands.Change{{
        Row: 0,
        Col: 0,
        Value: 'a',
    }}, flushed)
}

