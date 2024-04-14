package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestRenderable struct {
	cells [][]Cell
	loc   Location
	id    int
	z     int
}

func (t *TestRenderable) Z() int {
    return t.z
}

func (t *TestRenderable) Id() int {
    return t.id
}

var cells = [][]Cell{
	{Cell{Value: 69}, Cell{Value: 70}},
	{Cell{Value: 71}, Cell{Value: 72}},
}

var id int = 0

func newTestRenderable(loc Location, z int) TestRenderable {
    id++
	return TestRenderable{
		loc:   loc,
		cells: cells,
        id: id,
        z: z,
	}
}

func (t *TestRenderable) Render() (Location, [][]Cell) {
	return t.loc, t.cells
}

func TestRender(t *testing.T) {
	render := NewRender(5, 5)
	renderers := []TestRenderable{
		newTestRenderable(NewLocation(-1, -1), 1),
		newTestRenderable(NewLocation(4, 4), 1),
		newTestRenderable(NewLocation(-1, 4), 1),
		newTestRenderable(NewLocation(4, -1), 1),
		newTestRenderable(NewLocation(1, 1), 1),
	}

	for _, loc := range renderers {
		render.Add(&loc)
	}

	values := []byte{
		72, byte(' '), byte(' '), byte(' '), 71,
		byte(' '), 69, 70, byte(' '), byte(' '),
		byte(' '), 71, 72, byte(' '), byte(' '),
		byte(' '), byte(' '), byte(' '), byte(' '), byte(' '),
		70, byte(' '), byte(' '), byte(' '), 69,
	}

    _ = render.Render()
    render.debug()

	for i, value := range values {
		assert.Equal(t, render.buffer[i].Value, value)
	}

}
