package components

import "github.com/theprimeagen/vim-with-me/pkg/window"

type HighlightPoint struct {
	window.RenderBase
	position Position
	color    window.Color
	active   bool
}

func NewHighlightPoint(pos Position, z int, color window.Color) *HighlightPoint {
	return &HighlightPoint{
		RenderBase: window.NewRenderBase(z),
		position:   pos,
		color:      color,
		active:     false,
	}
}

func (h *HighlightPoint) SetActiveState(state bool) {
	h.active = state
}

var defaultLocation = window.NewLocation(0, 0)

func (h *HighlightPoint) Render() (window.Location, [][]window.Cell) {
	if !h.active {
		return defaultLocation, nil
	}

	return h.position.Position(), [][]window.Cell{
		{window.BackgroundCellOnly(h.color)},
	}
}
