package commands

import (
	"maps"

	"github.com/theprimeagen/vim-with-me/pkg/tcp2"
)

type Change struct {
    Row byte
    Col byte
    Value byte
}

var CHANGE_LENGTH = 3
const (
    COMMANDS = iota
    RENDER
    PARTIAL_RENDER
    CLOSE
    ERROR
    OPEN_WINDOW
    // TODO: what should i do on the server for missing commands?
    // I am thinking about closing the connection
    MISSING
    EXT_START
)

var commandMap = map[string]byte {
    "render": RENDER,
    "partial": PARTIAL_RENDER,
    "close": CLOSE,
    "error": ERROR,
    "openWindow": OPEN_WINDOW,
    "commands": COMMANDS,
}

type Commander struct {
    extensions map[string]byte
    size byte
}

func NewCommander() Commander {
    return Commander {
        extensions: maps.Clone(commandMap),
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

func (c *Commander) ToCommands() *tcp2.TCPCommand {
    b := []byte{}
    for name, k := range c.extensions {
        b = append(b, []byte(name)...)
        b = append(b, '\n')
        b = append(b, k)
    }

    return &tcp2.TCPCommand{
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

func PartialRender(changes Changes) *tcp2.TCPCommand {
    bytes := make([]byte, 0, len(changes) * CHANGE_LENGTH)
    for _, change := range changes {
        bytes = append(bytes, change.Bytes()...)
    }

    return &tcp2.TCPCommand{
        Command: PARTIAL_RENDER,
        Data: bytes,
    }
}

func Render(data []byte) *tcp2.TCPCommand {
    return &tcp2.TCPCommand{
        Command: RENDER,
        Data: data,
    }
}

func Close(msg []byte) *tcp2.TCPCommand {
    return &tcp2.TCPCommand{
        Command: CLOSE,
        Data: msg,
    }
}

func Error(msg []byte) *tcp2.TCPCommand {
    return &tcp2.TCPCommand{
        Command: ERROR,
        Data: msg,
    }
}
