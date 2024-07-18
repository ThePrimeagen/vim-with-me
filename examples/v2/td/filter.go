package td

import (
	"fmt"
	"strconv"
	"strings"
)

type Coord struct {
    Team int8
    Row int
    Col int
}

func (c *Coord) String() string {
    return fmt.Sprintf("%d%d,%d", c.Team, c.Row, c.Col)
}

func FromString(str string) (Coord, error) {
    if (len(str) < 4) {
        return Coord{}, fmt.Errorf("str not long enough")
    }

    team, err := strconv.Atoi(str[0:1])
    if err != nil {
        return Coord{}, fmt.Errorf("invalid team")
    }

    if team != 1 && team != 2 {
        return Coord{}, fmt.Errorf("invalid team")
    }

    parts := strings.Split(str[1:], ",")
    if len(parts) != 2 {
        return Coord{}, fmt.Errorf("invalid position")
    }

    row, err := strconv.Atoi(parts[0])
    if err != nil {
        return Coord{}, fmt.Errorf("invalid position")
    }

    col, err := strconv.Atoi(parts[1])
    if err != nil {
        return Coord{}, fmt.Errorf("invalid position")
    }

    return Coord{
        Row: row,
        Col: col,
        Team: int8(team),
    }, nil
}

func TDFilter(rows, cols int) func(msg string) bool {
    return func(msg string) bool {
        c, err := FromString(msg)
        if err != nil {
            return false
        }

        return c.Row <= rows && c.Row >= 2 &&
            c.Col <= cols && c.Col >= 2
    }
}
