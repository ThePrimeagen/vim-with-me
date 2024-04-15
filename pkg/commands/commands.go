package commands

import (
	"maps"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

const (
    COMMANDS = iota
    RENDER
    PARTIAL_RENDER
    COLOR
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

func PartialRender(cells []*window.CellWithLocation) *tcp.TCPCommand {
    bytes := make([]byte, 0, len(cells) * window.CELL_AND_LOC_ENCODING_LENGTH)
    for _, cell := range cells {
        data, err := cell.MarshalBinary()
        assert.Assert(err == nil, "encoding a cell should never fail")
        bytes = append(bytes, data...)
    }

    return &tcp.TCPCommand{
        Command: PARTIAL_RENDER,
        Data: bytes,
    }
}

type Openable interface {
    Dimensions() (byte, byte)
}

func OpenCommand(window Openable) *tcp.TCPCommand {
    rows, cols := window.Dimensions()
	return &tcp.TCPCommand{
		Command: OPEN_WINDOW,
		Data:    []byte{rows, cols},
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
