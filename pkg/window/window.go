package window

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
)

type Window struct {
	Rows byte
	Cols byte
    cache [][]byte
    changes []commands.Change
}

func NewWindow(rows, cols byte) *Window {
    cache := make([][]byte, rows)
    for i := range cache {
        cache[i] = make([]byte, cols)
        for j := range cache[i] {
            cache[i][j] = byte(' ')
        }
    }
    return &Window{
        Rows: rows,
        Cols: cols,
        cache: cache,
    }
}

func (w *Window) Set(row, col byte, value byte) error {
    if row < 0 || row >= w.Rows {
        return fmt.Errorf("Row out of bounds: %d", row)
    }

    if col < 0 || col >= w.Cols {
        return fmt.Errorf("Col out of bounds: %d", col)
    }

    if w.cache[row][col] != value {
        w.cache[row][col] = value
        w.changes = append(w.changes, commands.Change{
            Row: row,
            Col: col,
            Value: value,
        })
    }

    return nil
}

func (w *Window) SetString(row byte, value string) error {
    if len(value) > int(w.Cols) {
        return fmt.Errorf("String provided to Window is longer than columns: %d > %d", len(value), w.Cols)
    }

    for i, r := range []byte(value) {
        if w.cache[row][i] != r {
            w.cache[row][i] = r
            w.changes = append(w.changes, commands.Change{
                Row: row,
                Col: byte(i),
                Value: r,
            })
        }
    }

    return nil;
}

func (w *Window) SetWindow(value string) error {
    if len(value) != int(w.Rows) * int(w.Cols) {
        return fmt.Errorf("String provided to Window is not the correct length: %d != %d", len(value), w.Rows * w.Cols)
    }

    for i, r := range []byte(value) {
        row := i / int(w.Cols)
        col := i % int(w.Cols)

        if w.cache[row][col] != r {
            w.cache[row][col] = r
            w.changes = append(w.changes, commands.Change{
                Row: byte(row),
                Col: byte(col),
                Value: r,
            })
        }
    }

    return nil
}

func (w *Window) Render() string {
    out := ""
    for i := 0; i < int(w.Rows); i++ {
        out += string(w.cache[i])
    }
    w.changes = make([]commands.Change, 0)
    return out
}

func (w *Window) PartialRender() commands.Changes {
    changes := w.changes
    w.changes = make([]commands.Change, 0)
    return commands.Changes(changes)
}

func OpenCommand(window *Window) *tcp.TCPCommand {
    return &tcp.TCPCommand{
        Command: commands.OPEN_WINDOW,
        Data: []byte{window.Rows, window.Cols},
    }
}
