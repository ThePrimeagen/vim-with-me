package tower_defense

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"
)

type point struct {
	row   int
	col   int
	count int
}

type ChatAggregator struct {
	points []*point
	max    point
}

func NewChatAggregator() ChatAggregator {
	return ChatAggregator{
		points: make([]*point, 0, 100),
		max:    point{row: 0, col: 0, count: 0},
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

	c.points = append(c.points, &point{row: row, col: col, count: 1})
}

func (c *ChatAggregator) Count() (int, int) {
	count := 0
	for _, p := range c.points {
		count += p.count
	}
	return len(c.points), count
}

func (c *ChatAggregator) Reset() (int, int) {
	slog.Debug("ChatAggregator#Reset", "points", len(c.points), "max", c.max)
	c.points = make([]*point, 0, 100)

	out := c.max
	c.max = point{row: 0, col: 0, count: 0}

	return out.row, out.col
}

func ParseChatMessage(msg string) (int, int, error) {

	parts := strings.SplitN(msg, ":", 3)
	if len(parts) != 3 {
		return 0, 0, errors.New("malformed message")
	}

	if parts[0] != "t" {
		return 0, 0, errors.New("not a tower placement message")
	}

	row, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	col, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, err
	}

	return row, col, nil
}
