package td

import (
	"fmt"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/v2/td/objects"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/cmd"
)

type TDCommander struct {
    Cmdr *cmd.Cmder
    Debug *testies.DebugFile
}

func (t *TDCommander) WritePositions(positions objects.Positions, team uint8) {
    posStr := fmt.Sprintf("%c%s", team, positions.String())
    t.Debug.WriteStrLine(fmt.Sprintf("tdcmdr#WritePositions: %s", posStr))
    err := t.Cmdr.WriteLine([]byte(posStr))
    assert.NoError(err, "unable to communicate with the game")
}

func (t *TDCommander) PlayRound() {
    err := t.Cmdr.WriteLine([]byte{'p'})
    t.Debug.WriteStrLine("td#cmdr#PlayRound")
    assert.NoError(err, "unable to communicate play round with game")
}

func (t *TDCommander) Countdown(tme time.Duration) {
    countdown := fmt.Sprintf("c%d", tme.Microseconds())
    err := t.Cmdr.WriteLine([]byte(countdown))
    t.Debug.WriteStrLine(fmt.Sprintf("tdcmdr#WritePositions: %s", countdown))
    assert.NoError(err, "unable to communicate countdown with game")
}


