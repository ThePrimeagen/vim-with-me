package commands

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/tcp"
)

type Change struct {
    Row int
    Col int
    Value rune
}

func (c *Change) Command() *tcp.TCPCommand {
    return PartialRender(c)
}

func PartialRender(change *Change) *tcp.TCPCommand {
    return &tcp.TCPCommand{
        Command: "p",
        Data: fmt.Sprintf("%d:%d:%c", change.Row, change.Col, change.Value),
    }
}

func Render(data string) *tcp.TCPCommand {
    return &tcp.TCPCommand{
        Command: "r",
        Data: data,
    }
}

func Close(msg string) *tcp.TCPCommand {
    return &tcp.TCPCommand{
        Command: "c",
        Data: msg,
    }
}

func Error(msg string) *tcp.TCPCommand {
    return &tcp.TCPCommand{
        Command: "e",
        Data: msg,
    }
}
