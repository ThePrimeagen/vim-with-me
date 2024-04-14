package window

import "github.com/theprimeagen/vim-with-me/pkg/assert"

const FOREGROUND = 1
const COLOR = 3
const COLOR_ENCODING_LENGTH = FOREGROUND + COLOR

type Color struct {
	red        byte
	blue       byte
	green      byte
	foreground bool
}

var DEFAULT_BACKGROUND = Color{red: 0, blue: 0, green: 0}

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
