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

func (c *Change) String() string {
    return fmt.Sprintf("%d:%d:%c", c.Row, c.Col, c.Value)
}

type Changes []Change

func PartialRender(changes Changes) *tcp.TCPCommand {
    data := ""
    for _, change := range changes {
        data += change.String()
    }

    return &tcp.TCPCommand{
        Command: "p",
        Data: data,
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
