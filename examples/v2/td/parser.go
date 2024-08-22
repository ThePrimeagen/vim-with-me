package td

import (
	"encoding/json"
	"fmt"

	"github.com/theprimeagen/vim-with-me/examples/v2/td/objects"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
)

type CmdErrParser struct {
    debug *testies.DebugFile
    Gs chan objects.GameState
    Done chan string
}

func NewCmdErrParser(debug *testies.DebugFile) CmdErrParser {
    return CmdErrParser{
        Gs: make(chan objects.GameState, 1),
        Done: make(chan string, 1),
        debug: debug,
    }
}

func (c *CmdErrParser) Parse(b []byte) (int, error) {
    var gs objects.GameState;
    err := json.Unmarshal(b, &gs)
    if err != nil {
        c.debug.WriteStrLine(fmt.Sprintf("td: %s\n", string(b)))
        return len(b), nil
    }
    c.Gs <- gs

    return len(b), nil
}


