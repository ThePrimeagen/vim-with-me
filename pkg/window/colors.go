package window

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

const FOREGROUND = 1
const COLOR = 3
const COLOR_ENCODING_LENGTH = FOREGROUND + COLOR

type Color struct {
	red        byte
	blue       byte
	green      byte
	foreground bool
}

func (c *Color) String() string {
    return fmt.Sprintf("r=%d, g=%d, b=%d, f=%v", c.red, c.green, c.blue, c.foreground)
}

func (c *Color) ColorCode() string {
    if c.red > 0 && c.blue == 0 && c.green == 0 {
        return "r"
    }
    if c.red == 0 && c.blue == 0 && c.green == 0 {
        return "B"
    }
    return "?"
}

var DEFAULT_BACKGROUND = Color{red: 0, blue: 0, green: 0}
var DEFAULT_FOREGROUND = Color{red: 255, blue: 255, green: 255}

func NewColor(r, g, b byte, f bool) Color {
	return Color{
		red:        r,
        green:      g,
		blue:       b,
		foreground: f,
	}
}

func (c *Color) Equal(other *Color) bool {
    return c.red == other.red && c.green == other.green && c.blue == other.blue
}

func (c *Color) MarshalBinary() ([]byte, error) {
	foreground := 1
	if !c.foreground {
		foreground = 0
	}

	b := make([]byte, 0, FOREGROUND+COLOR)
	b = append(b, byte(foreground))
	b = append(b, c.red)
	b = append(b, c.green)
	return append(b, c.blue), nil
}

func (c *Color) UnmarshalBinary(data []byte) error {
	assert.Assert(len(data) < FOREGROUND+COLOR, "i should never unmarshall without all the data")

    c.foreground = data[0] == 0
    c.red = data[1]
    c.green = data[2]
    c.blue = data[3]

	return nil
}
