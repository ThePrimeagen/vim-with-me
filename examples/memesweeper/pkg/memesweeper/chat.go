package memesweeper

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/window"
)

type Point struct {
	row   int
	col   int
	count int
}

type ChatAggregator struct {
	points []*Point
	max    *Point
	active bool
}

func NewChatAggregator() *ChatAggregator {
	max := &Point{row: 0, col: 0, count: 0}
	return &ChatAggregator{
		points: make([]*Point, 0, 100),
		max:    max,
		active: false,
	}
}

func (c *ChatAggregator) Add(row, col int) {
	if !c.active {
		slog.Debug("ChatAggregator#Add inactive", "row", row, "col", col)
		return
	}

	slog.Debug("ChatAggregator#Add placed", "row", row, "col", col, "max", c.max)
	for _, p := range c.points {
		if p.row == row && p.col == col {
			p.count++
			slog.Debug("ChatAggregator#Add incrementing", "count", p.count)
			if c.max.count < p.count {
				c.max = p
				slog.Debug("ChatAggregator#Add new max", "max", c.max)
			}
			return
		}
	}

	point := &Point{row: row, col: col, count: 1}
	c.points = append(c.points, point)
	if c.max.count < point.count {
		c.max = point
		slog.Debug("ChatAggregator#Add new max", "max", c.max)
	}
}

func (c *ChatAggregator) SetActiveState(state bool) {
	c.active = state
}

func (c *ChatAggregator) Current() Point {
	return *c.max
}

func (c *ChatAggregator) Position() window.Location {
	return window.NewLocation(c.max.row, c.max.col)
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
	for _, p := range c.points {
		slog.Debug("    ChatAggregator#Reset#points", "point", p)
	}
	c.points = make([]*Point, 0, 100)

	out := c.max
	c.max = &Point{row: 0, col: 0, count: 0}

	return *out
}

func isCol(msg byte) bool {
	return msg >= 'A' && msg <= 'Z'
}

func ParseChatMessage(msg string) (int, string, error) {
	if len(msg) != 2 {
		return 0, "", errors.New("malformed chat msg")
	}

	a := strings.ToUpper(string(msg[0]))
	b := strings.ToUpper(string(msg[1]))

	row := a
	col := b
	if isCol(a[0]) {
		row = b
		col = a
	}

	rowNum, err := strconv.Atoi(row)
	if err != nil {
		return 0, "", err
	}

	return rowNum, col, nil
}
