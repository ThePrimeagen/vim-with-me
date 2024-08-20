package players

import (
	"context"
	"strings"

	"github.com/theprimeagen/vim-with-me/examples/v2/td"
	"github.com/theprimeagen/vim-with-me/examples/v2/td/objects"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

type BoxPos struct {
    maxRows int
    position int
    sendResults bool
}

func NewBoxPos(maxRows int) *BoxPos {
    return &BoxPos{
        maxRows: maxRows,
        position: 0,
        sendResults: false,
    }
}

func (r *BoxPos) nextPos() objects.Position {
    col := 6
    if r.position & 0x1 == 0 {
        col = 12
    }

    row := 1 + ((r.position % 4) / 2) * 6
    r.position++

    return objects.Position{
        Row: uint(row),
        Col: uint(col),
    }
}

func (r *BoxPos) Stats() objects.Stats {
    return objects.Stats{};
}

// As of right now, out of bounds guesses will place a tower within your area
// randomly
func (r *BoxPos) Moves(gs *objects.GameState, ctx context.Context) []objects.Position {
    out := []objects.Position{}
    for range gs.AllowedTowers {
        out = append(out, r.nextPos())
    }
    return out
}

// As of right now, out of bounds guesses will place a tower within your area
// randomly

func (r *BoxPos) StreamResults(team uint8, gs *objects.GameState, out PositionChan, done Done, ctx context.Context) {
    if !r.sendResults {
        done <- struct{}{}
        return
    }

    r.sendResults = false;
    pos := []objects.Position{}
    for range gs.AllowedTowers {
        pos = append(pos, r.nextPos())
    }

    out <- pos
    done <- struct{}{}
}

func (r *BoxPos) StartRound() {
    r.sendResults = true
}

type RandomPos struct {
    outOfBounds objects.Position
    sendResults bool
}

func (r *RandomPos) StartRound() {
    r.sendResults = true
}

func NewRandomPos(maxRows int) *RandomPos {
    return &RandomPos{
        outOfBounds: objects.OutOfBoundPosition(),
        sendResults: false,
    }
}

func (r *RandomPos) StreamResults(team uint8, gs *objects.GameState, out PositionChan, done Done, ctx context.Context) {
    if !r.sendResults {
        done <- struct{}{}
        return
    }

    go func() {
        r.sendResults = false;
        out <- r.Moves(gs, ctx)
        done <- struct{}{}
    }()
}

func (r *RandomPos) Moves(gs *objects.GameState, ctx context.Context) []objects.Position {
    out := []objects.Position{}
    for range gs.AllowedTowers {
        out = append(out, r.outOfBounds)
    }
    return out
}

func (r *RandomPos) Stats() objects.Stats {
    return objects.Stats{}
}

func (f *RandomPos) Name() string { return "random-pos" }
func (f *RandomPos) Run(ctx context.Context) { }
func (r *RandomPos) EndRound(gs *objects.GameState, cmdr td.TDCommander) { }
func (f *BoxPos) Run(ctx context.Context) { }
func (f *BoxPos) Name() string { return "box-pos" }
func (r *BoxPos) EndRound(gs *objects.GameState, cmdr td.TDCommander) { }

func StratPlayerFromString(arg string) Player {
    assert.Assert(strings.HasPrefix(arg, "strat"), "invalid player string for strat client", "arg", arg)

    parts := strings.Split(arg, ":")
    assert.Assert(len(parts) == 2, "invalid strat player string colon count", "parts", parts)

    switch (parts[1]) {
    case "rand":
        return NewRandomPos(24)
    case "box":
        return NewBoxPos(24)
    }
    assert.Assert(false, "invalid strat provided", "arg", arg)
    return nil
}
