package players

import (
	"context"
	"strings"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/v2/td"
	"github.com/theprimeagen/vim-with-me/examples/v2/td/objects"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

type PositionChan chan<- []objects.Position;
type Done chan<- struct{}

type Player interface {
    StartRound()
    EndRound(gs *objects.GameState, cmdr td.TDCommander)
    StreamResults(team uint8, gs *objects.GameState, out PositionChan, done Done, ctx context.Context)
    Stats() objects.Stats
    Run(ctx context.Context)
    Name() string
}

type TeamPlayer struct {
    Player Player
    team uint8
    cmdr td.TDCommander
}

func NewTeamPlayer(player Player, team uint8, cmdr td.TDCommander) TeamPlayer {
    return TeamPlayer{Player: player, team: team, cmdr: cmdr}
}

func (t *TeamPlayer) StreamMoves(ctx context.Context, gs *objects.GameState) {
    out := make(chan []objects.Position, 10)
    done := make(chan struct{}, 1)
    t.Player.StreamResults(t.team, gs, out, done, ctx)

    outer:
    for {
        select {
        case <-done:
            break outer
        case <-ctx.Done():
            break outer
        case pos := <-out:
            t.cmdr.WritePositions(objects.Positions(pos), t.team)
        }
    }

    time.Sleep(time.Millisecond * 100)
    close(out)
    close(done)
}

func NewTeamPlayerFromString(arg string, debug *testies.DebugFile, ctx context.Context, team uint8, cmdr td.TDCommander) TeamPlayer {
    parts := strings.Split(arg, ":")

    var player Player

    switch (parts[0]) {
    case "ai":
        ai := AIPlayerFromString(arg, team, debug, ctx)
        player = &ai
    case "strat":
        player = StratPlayerFromString(arg)
    case "twitch":
        twitch := TwitchPlayerFromString(arg, team)
        player = &twitch
    default:
        assert.Never("unknown player type", "arg", arg)
    }

    return NewTeamPlayer(player, team, cmdr)
}
