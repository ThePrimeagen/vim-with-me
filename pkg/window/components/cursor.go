package components

import "github.com/theprimeagen/vim-with-me/pkg/window"

type Position interface {
	Position() window.Location
}
type HighlightPoint struct {
    window.RenderBase
	position Position
	color    window.Color
}

func NewHighlightPoint(pos Position, z int, color window.Color) HighlightPoint {
	return HighlightPoint{
        RenderBase: window.NewRenderBase(z),
		position: pos,
		color:    color,
	}
}

func (h *HighlightPoint) Render() (window.Location, [][]window.Cell) {
    return h.position.Position(), [][]window.Cell{
        {window.BackgroundCellOnly(h.color)},
    }
}
