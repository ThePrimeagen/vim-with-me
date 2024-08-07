package td

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

type Tower struct {
    Row int `json:"row"`
    Col int `json:"col"`
    Ammo int `json:"ammo"`
    Level int `json:"level"`
}

type Range struct {
    StartRow uint `json:"startRow"`
    EndRow uint `json:"endRow"`
}

type GameState struct {
    Rows uint `json:"rows"`
    Cols uint `json:"cols"`
    AllowedTowers int `json:"allowedTowers"`
    YourCreepDamage uint `json:"yourCreepDamage"`
    EnemyCreepDamage uint `json:"enemyCreepDamage"`
    YourTowers []Tower `json:"yourTowers"`
    EnemyTowers []Tower `json:"enemyTowers"`
    TowerPlacementRange Range `json:"towerPlacementRange"`
    CreepSpawnRange Range `json:"creepSpawnRange"`
    Round uint `json:"round"`
    Finished bool `json:"finished"`
    Playing bool `json:"playing"`
    Winner uint `json:"winner"`
}

func (gs *GameState) String() string {
    b, err := json.Marshal(gs)
    assert.NoError(err, "unable to create gamestate string")
    return string(b)
}

type FetchPositions struct {
    ai AIFetcher
    maxTries uint
    debug *testies.DebugFile
}

type AIFetcher interface {
    ReadWithTimeout(prompt string, t time.Duration) (string, error)
}

func NewFetchPosition(ai AIFetcher, debug *testies.DebugFile) FetchPositions {
    return FetchPositions {
        ai: ai,
        debug: debug,
        maxTries: 3,
    };
}

type Stats struct {
    BadParses int
    RandomGuesses int
    TotalTowers int
}

func (s *Stats) Add(stats Stats) {
    s.BadParses += stats.BadParses
    s.RandomGuesses += stats.RandomGuesses
    s.TotalTowers += stats.TotalTowers
}

func (s *Stats) String() string {
    return fmt.Sprintf("%d,%d,%d", s.TotalTowers, s.RandomGuesses, s.BadParses)
}

type Position struct {
    Row uint
    Col uint
}

type Positions []Position
func (p Positions) String() string {
    out := []string{}
    for _, pos := range p {
        out = append(out, pos.String())
    }
    return strings.Join(out, ",")
}

func OutOfBoundPosition() Position {
    return Position{Row: 999, Col: 999}
}

func (p *Position) OutOfBounds() bool {
    return p.Row == 999 && p.Col == 999
}

func (p *Position) String() string {
    return fmt.Sprintf("%d,%d", p.Row, p.Col)
}

// if the ai provides some line that is unreadable i will just provide a value too large and thus random placement
func getPosFromAIResponse(line string) Position {
    parts := strings.Split(line, ",")
    for i, p := range parts {
        parts[i] = strings.TrimSpace(p)
    }

    if len(parts) != 2 {
        return OutOfBoundPosition()
    }

    row, err := strconv.Atoi(parts[0])
    if err != nil {
        return OutOfBoundPosition()
    }

    col, err := strconv.Atoi(parts[1])
    if err != nil {
        return OutOfBoundPosition()
    }
    return Position{
        Row: uint(row),
        Col: uint(col),
    }
}

func (f *FetchPositions) Fetch(gs *GameState) ([]Position, Stats) {
    count := gs.AllowedTowers
    responses := []Position{}
    var tries uint = 0
    guesses := 0
    badParses := 0

    for len(responses) < count && tries < f.maxTries {
        resp, err := f.ai.ReadWithTimeout(gs.String(), time.Second * 5)
        tries++

        if err != nil {
            f.debug.WriteStrLine(fmt.Sprintf("ai-placement response: \"%s\" err: \"%s\"", resp, err.Error()))
            parts := strings.Split(err.Error(), "try again in ")
            if len(parts) == 2 {
                secsStr := strings.Split(parts[1], " ")[0]
                secs, err := strconv.ParseFloat(secsStr[0:len(secsStr) - 2], 64)
                if err == nil {
                    dur := time.Duration(float64(time.Second) * secs)
                    f.debug.WriteStrLine(fmt.Sprintf("ai-placement wait time required: %d", dur))
                    <-time.NewTimer(dur).C
                }
            }
            continue
        }

        if resp == "" {
            <-time.NewTimer(time.Second).C
        }

        for _, line := range strings.Split(resp, "\n") {
            line = strings.TrimSpace(line)
            if line == "" || len(responses) == count {
                break;
            }

            parsedLine := getPosFromAIResponse(line)
            if parsedLine.OutOfBounds() {
                badParses++
                continue
            }

            responses = append(responses, parsedLine)
            f.debug.WriteStrLine(fmt.Sprintf("ai-placement: %s - %s", line, parsedLine.String()))
        }
    }

    for range gs.AllowedTowers - len(responses) {
        guesses++
        responses = append(responses, OutOfBoundPosition())
    }

    return responses, Stats{
        BadParses: badParses,
        RandomGuesses: guesses,
        TotalTowers: gs.AllowedTowers,
    }
}


