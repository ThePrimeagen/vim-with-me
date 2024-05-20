package components

import "github.com/theprimeagen/vim-with-me/pkg/window"

type Position interface {
	Position() window.Location
}

type CompositePosition struct {
	pos    Position
	offset window.Location
}

func NewCompositePosition(pos Position, offset window.Location) *CompositePosition {
	return &CompositePosition{pos: pos, offset: offset}
}

// TODO(v1): Rename location into vector and stop being a joke
func (c *CompositePosition) Position() window.Location {
	loc := c.offset
	pos := c.pos.Position()
	loc.Row += pos.Row
	loc.Col += pos.Col

	return loc
}
