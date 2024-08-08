package objects

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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

type Stats struct {
    BadParses int
    RandomGuesses int
    TotalTowers int
}

func (s *Stats) Add(stats Stats) *Stats {
    s.BadParses += stats.BadParses
    s.RandomGuesses += stats.RandomGuesses
    s.TotalTowers += stats.TotalTowers
    return s
}

func (s *Stats) String() string {
    return fmt.Sprintf("%d,%d,%d", s.TotalTowers, s.RandomGuesses, s.BadParses)
}

type Position struct {
    Row uint
    Col uint
}

func PositionFromString(str string) (Position, error) {
    if (len(str) < 3) {
        return Position{}, fmt.Errorf("str not long enough")
    }

    parts := strings.Split(str, ",")
    if len(parts) != 2 {
        return Position{}, fmt.Errorf("invalid position")
    }

    row, err := strconv.Atoi(parts[0])
    if err != nil {
        return Position{}, fmt.Errorf("invalid position")
    }

    col, err := strconv.Atoi(parts[1])
    if err != nil {
        return Position{}, fmt.Errorf("invalid position")
    }

    return Position{
        Row: uint(row),
        Col: uint(col),
    }, nil
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
