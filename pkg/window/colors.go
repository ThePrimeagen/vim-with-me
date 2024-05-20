package window

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

const FOREGROUND = 1
const COLOR = 3
const COLOR_ENCODING_LENGTH = FOREGROUND + COLOR

type Color struct {
	Red        byte `json:"red"`
	Blue       byte `json:"blue"`
	Green      byte `json:"green"`
	Foreground bool `json:"foreground"`
}

func (c *Color) String() string {
	return fmt.Sprintf("r=%d, g=%d, b=%d, f=%v", c.Red, c.Green, c.Blue, c.Foreground)
}

func (c *Color) ColorCode() string {
	if c.Red > 0 && c.Blue == 0 && c.Green == 0 {
		return "r"
	}
	if c.Red == 0 && c.Blue == 0 && c.Green == 0 {
		return "B"
	}
	return "?"
}

var DEFAULT_BACKGROUND = Color{Red: 0, Blue: 0, Green: 0}
var DEFAULT_FOREGROUND = Color{Red: 255, Blue: 255, Green: 255}

func NewColor(r, g, b byte, f bool) Color {
	return Color{
		Red:        r,
		Green:      g,
		Blue:       b,
		Foreground: f,
	}
}

func (c *Color) Equal(other *Color) bool {
	return c.Red == other.Red && c.Green == other.Green && c.Blue == other.Blue
}

func (c *Color) MarshalBinary() ([]byte, error) {
	foreground := 1
	if !c.Foreground {
		foreground = 0
	}

	b := make([]byte, 0, FOREGROUND+COLOR)
	b = append(b, byte(foreground))
	b = append(b, c.Red)
	b = append(b, c.Green)
	return append(b, c.Blue), nil
}

func (c *Color) UnmarshalBinary(data []byte) error {
	assert.Assert(len(data) >= COLOR_ENCODING_LENGTH, "Color#UnmarshalBinary - Not enough data to unmarshal")

	c.Foreground = data[0] == 0
	c.Red = data[1]
	c.Green = data[2]
	c.Blue = data[3]

	return nil
}
