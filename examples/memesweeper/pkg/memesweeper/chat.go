package memesweeper

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"
)

type Point struct {
	row   int
	col   int
	count int
}

type ChatAggregator struct {
	points []*Point
    max Point
}

func NewChatAggregator() ChatAggregator {
	return ChatAggregator{
		points: make([]*Point, 0, 100),
        max: Point{row: 0, col: 0, count: 0},
	}
}

func (c *ChatAggregator) Add(row, col int) {
    for _, p := range c.points {
        if p.row == row && p.col == col {
            p.count++
            if c.max.count < p.count {
                c.max = *p
            }
            return
        }
    }

    c.points = append(c.points, &Point{row: row, col: col, count: 1})
}

func (c *ChatAggregator) Count() (int, int) {
    count := 0
    for _, p := range c.points {
        count += p.count
    }
    return len(c.points), count
}

func (c *ChatAggregator) Reset() Point {
    slog.Debug("ChatAggregator#Reset", "points", len(c.points), "max", c.max)
	c.points = make([]*Point, 0, 100)

    out := c.max
    c.max = Point{row: 0, col: 0, count: 0}

    return out
}

func isCol(msg string) bool {
    return len(msg) == 1 && msg[0] >= 'A' && msg[0] <= 'Z'
}

func ParseChatMessage(msg string) (int, string, error) {

	parts := strings.SplitN(msg, ":", 2)
	if len(parts) != 2 {
		return 0, "", errors.New("malformed message")
	}

    idx := 0
    if isCol(parts[0]) {
        idx = 1
    }

    row, err := strconv.Atoi(parts[idx])
    if err != nil {
        return 0, "", err
    }

    return row, parts[0], nil
}

