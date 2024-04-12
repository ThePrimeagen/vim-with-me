package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestRenderable struct {
    cells [][]Cell
    loc Location
}

func (t *TestRenderable) Render() (Location, [][]Cell) {
    return t.loc, t.cells
}

var cells = [][]Cell{
    {Cell{Value: 69}, Cell{Value: 70}},
    {Cell{Value: 71}, Cell{Value: 72}},
}
func testLocationAndValue(t *testing.T, loc Location, value byte, index int, render Renderer) {
    render.Place(&TestRenderable{
        cells: cells,
        loc: loc,
    })

    assert.Equal(t, render.buffer[index].Value, value)
    assert.Equal(t, render.buffer[index].Count, 1)
}

func TestTranslate(t *testing.T) {
    render := NewRender(5, 5)
    locations := []Location{
        NewLocation(-1, -1),
        NewLocation(4, 4),
        NewLocation(-1, 4),
        NewLocation(4, -1),
    }
    values := []byte{
        byte(72),
        byte(69),
        byte(71),
        byte(70),
    }

    indexes := []int{
        0,
        len(render.buffer) - 1,
        4,
        20,
    }

    for i, loc := range locations {
        value := values[i]
        index := indexes[i]

        testLocationAndValue(t, loc, value, index, render)
    }

    render.Place(&TestRenderable{
        cells: cells,
        loc: NewLocation(1, 1),
    })

    assert.Equal(t, render.buffer[6].Value, byte(69))
    assert.Equal(t, render.buffer[7].Value, byte(70))
    assert.Equal(t, render.buffer[11].Value, byte(71))
    assert.Equal(t, render.buffer[12].Value, byte(72))

    render.Place(&TestRenderable{
        cells: cells,
        loc: NewLocation(1, 1),
    })

    assert.Equal(t, render.buffer[6].Count, 2)
    assert.Equal(t, render.buffer[7].Count, 2)
    assert.Equal(t, render.buffer[11].Count, 2)
    assert.Equal(t, render.buffer[12].Count, 2)
}

