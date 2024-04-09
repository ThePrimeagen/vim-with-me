package commands

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/tcp"
)

type Change struct {
    Row byte
    Col byte
    Value byte
}

var CHANGE_LENGTH = 3
const (
    RENDER = iota
    PARTIAL_RENDER
    CLOSE
    ERROR
    OPEN_WINDOW
    COMMANDS
    MISSING
    EXT_START
)

type Commander struct {
    extensions map[string]byte
    size byte
}

func NewCommander() Commander {
    return Commander {
        extensions: map[string]byte{},
        size: byte(EXT_START),
    }
}

func (c *Commander) AddCommand(name string) {
    if _, ok := c.extensions[name]; ok {
        return
    }
    c.extensions[name] = c.size;
    c.size += 1
}

func (c *Commander) GetCommandByte(name string) byte {
    if b, ok := c.extensions[name]; ok {
        return b
    }
    return MISSING
}

func (c *Commander) ToCommands() *tcp.TCPCommand {
    b := []byte{}
    for name, k := range c.extensions {
        b = append(b, []byte(name)...)
        b = append(b, '\n')
        b = append(b, k)
    }

    return &tcp.TCPCommand{
        Command: COMMANDS,
        Data: b,
    }
}

func (c *Commander) ToString(b byte) string {
    // TODO: Probably improve performan}ce... maybe.. if there is more than 30?
    for name, k := range c.extensions {
        if k == b {
            return name
        }
    }
    return ""
}

func (c *Change) Bytes() []byte {
    return []byte{
        c.Row,
        c.Col,
        c.Value,
    }
}

type Changes []Change

func PartialRender(changes Changes) *tcp.TCPCommand {
    bytes := make([]byte, 0, len(changes) * CHANGE_LENGTH)
    for _, change := range changes {
        bytes = append(bytes, change.Bytes()...)
    }

    return &tcp.TCPCommand{
        Command: PARTIAL_RENDER,
        Data: bytes,
    }
}

func Render(data []byte) *tcp.TCPCommand {
    return &tcp.TCPCommand{
        Command: RENDER,
        Data: data,
    }
}

func Close(msg []byte) *tcp.TCPCommand {
    return &tcp.TCPCommand{
        Command: CLOSE,
        Data: msg,
    }
}

func Error(msg []byte) *tcp.TCPCommand {
    return &tcp.TCPCommand{
        Command: ERROR,
        Data: msg,
    }
}
