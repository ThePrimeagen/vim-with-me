package objects

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

const TowerColSize = 5
const TowerRowSize = 3

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

    OneCreepDamage uint `json:"oneCreepDamage"`
    OneTowers []Tower `json:"oneTowers"`
    OneTowerPlacementRange Range `json:"oneTowerPlacementRange"`
    OneCreepSpawnRange Range `json:"oneCreepSpawnRange"`
    OneTotalTowersBuild uint `json:"oneTotalTowersBuild"`
    OneTotalProjectiles uint `json:"oneTotalProjectiles"`
    OneTotalTowerUpgrades uint `json:"oneTotalTowerUpgrades"`
    OneTotalCreepDamage uint `json:"oneTotalCreepDamage"`
    OneTotalTowerDamage uint `json:"oneTotalTowerDamage"`
    OneTotalDamageFromCreeps uint `json:"oneTotalDamageFromCreeps"`

    TwoCreepDamage uint `json:"twoCreepDamage"`
    TwoTowers []Tower `json:"TwoTowers"`
    TwoTowerPlacementRange Range `json:"twoTowerPlacementRange"`
    TwoCreepSpawnRange Range `json:"twoCreepSpawnRange"`
    TwoTotalTowersBuild uint `json:"twoTotalTowersBuild"`
    TwoTotalProjectiles uint `json:"twoTotalProjectiles"`
    TwoTotalTowerUpgrades uint `json:"twoTotalTowerUpgrades"`
    TwoTotalCreepDamage uint `json:"twoTotalCreepDamage"`
    TwoTotalTowerDamage uint `json:"twoTotalTowerDamage"`
    TwoTotalDamageFromCreeps uint `json:"twoTotalDamageFromCreeps"`

    Round uint `json:"round"`
    Finished bool `json:"finished"`
    Playing bool `json:"playing"`
    Winner uint `json:"winner"`
}

type PromptState struct {
    Rows uint `json:"rows"`
    Cols uint `json:"cols"`
    AllowedTowers int `json:"allowedTowers"`
    YourCreepDamage uint `json:"yourCreepDamage"`
    EnemyCreepDamage uint `json:"enemyCreepDamage"`
    YourTowers []Tower `json:"yourTowers"`
    EnemyTowers []Tower `json:"TwoTowers"`
    YourTowerPlacementRange Range `json:"yourTowerPlacementRange"`
    EnemyTowerPlacementRange Range `json:"enemyTowerPlacementRange"`
    YourCreepSpawnRange Range `json:"yourCreepSpawnRange"`
    EnemyCreepSpawnRange Range `json:"enemyCreepSpawnRange"`
    Round uint `json:"round"`
}

func (p *GameState) towers(team uint8) []Tower {
    if team == '1' {
        return p.OneTowers
    }
    return p.TwoTowers
}

func (p *GameState) towerRange(team uint8) Range {
    if team == '1' {
        return p.OneTowerPlacementRange
    }
    return p.TwoTowerPlacementRange
}

func (p *GameState) creepRange(team uint8) Range {
    if team == '1' {
        return p.OneCreepSpawnRange
    }
    return p.TwoCreepSpawnRange
}

func (p *GameState) creepDamage(team uint8) uint {
    if team == '1' {
        return p.OneCreepDamage
    }
    return p.TwoCreepDamage
}

func (p *GameState) PromptState(team uint8) PromptState {
    enemyTeam := uint8('2')
    if team == '2' {
        enemyTeam = '1'
    }

    return PromptState{
        Rows: p.Rows,
        Cols: p.Cols,
        AllowedTowers: p.AllowedTowers,
        Round: p.Round,

        YourTowers: p.towers(team),
        EnemyTowers: p.towers(enemyTeam),

        YourCreepDamage: p.creepDamage(team),
        EnemyCreepDamage: p.creepDamage(enemyTeam),

        YourCreepSpawnRange: p.creepRange(team),
        EnemyCreepSpawnRange: p.creepRange(enemyTeam),

        YourTowerPlacementRange: p.towerRange(team),
        EnemyTowerPlacementRange: p.towerRange(enemyTeam),
    }
}

func (p *PromptState) String() string {
    b, err := json.Marshal(p)
    assert.NoError(err, "unable to create gamestate string")
    return string(b)
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

    parts[0] = strings.TrimSpace(parts[0])
    parts[1] = strings.TrimSpace(parts[1])

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

func (p *Position) OutOfBounds(gs *GameState) bool {
    return p.Row >= gs.Rows || p.Col >= gs.Cols
}

func (p *Position) String() string {
    return fmt.Sprintf("%d,%d", p.Row, p.Col)
}
