package encoding

import (
	"encoding/json"

	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
)

type charable struct {
	brightness []float64
	rgb        byte
	count      int
	char       string
}

func (c *charable) brightnessRange() BrightnessRange {
	min := 4.0
	max := 0.0
	sum := 0.0

	for _, bright := range c.brightness {
		if bright > max {
			max = bright
		}
		if bright < min {
			min = bright
		}

		sum += bright
	}

	return BrightnessRange{
		Min: min,
		Max: max,
		Avg: sum / float64(len(c.brightness)),
	}
}

type colorToChar struct {
	charOrdering []*charable
}

func NewColorToChar() *colorToChar {
	return &colorToChar{
		charOrdering: make([]*charable, 0),
	}
}

func (c *colorToChar) Map(frame ansiparser.Frame) {
	for i, ch := range frame.Chars {

		char := string(ch)
		bright := RGBBrightness(frame.Color[i])

		var charPtr *charable = nil

		for _, c := range c.charOrdering {
			if c.char == char {
				charPtr = c
				break
			}
		}

		if charPtr == nil {
			charPtr = &charable{
				char:       char,
				brightness: make([]float64, 0),
				count:      0,
				rgb:        frame.Color[i],
			}

			c.charOrdering = append(c.charOrdering, charPtr)
		}

		charPtr.brightness = append(charPtr.brightness, bright)
		charPtr.count++
	}
}

func (c *colorToChar) String() string {

	brightnesses := make(map[string]BrightnessRange, 0)
	for _, char := range c.charOrdering {
        brightnesses[char.char] = char.brightnessRange()
	}

    bytes, err := json.Marshal(brightnesses)
    if err != nil {
        return "errored"
    }


	return string(bytes)
}
