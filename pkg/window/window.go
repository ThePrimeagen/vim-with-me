package window

import (
	"fmt"

	"chat.theprimeagen.com/pkg/tcp"
)

type Window struct {
	Rows int
	Cols int
}

func NewWindow(rows, cols int) *Window {
    return &Window{
        Rows: rows,
        Cols: cols,
    }
}

func OpenCommand(window *Window) *tcp.TCPCommand {
    return &tcp.TCPCommand{
        Command: "open-window",
        Data: fmt.Sprintf("%d:%d", window.Rows, window.Cols),
    }
}
