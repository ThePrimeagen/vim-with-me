package color

import (
	"fmt"
	"math"
	"slices"
	"strings"

	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
)

func round(f float64) float64 {
    return math.Round(f * 100) / 100
}

func RGBBrightness(rgb byte) float64 {
	red := float64((rgb & 0b111_000_00) >> 5)
	green := float64((rgb & 0b000_111_00) >> 2)
	blue := float64(rgb & 0b11)

    return round(red / 7.0) + round(green / 7.0) + round(blue / 3.0)
}

func RGBToString(rgb byte) string {
	red := (rgb & 0b111_000_00) >> 5
	green := (rgb & 0b000_111_00) >> 2
	blue := rgb & 0b11

	return fmt.Sprintf("%02x%02x%02x", red, green, blue)
}


type Lumin struct {
    buffer []byte
    idx int
}

type charable struct {
    brightness []float64
    rgb byte
    count int
    char string
}

func (c *charable) String() string {
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

    return fmt.Sprintf("%s: min=%f avg=%f max=%f same=%v diff=%f", c.char, min, sum / float64(len(c.brightness)), max, min == max, max - min)
}


type ColorToChar struct {
    CharOrder []*charable
}

func NewColorToChar() *ColorToChar {
    return &ColorToChar{
        CharOrder: make([]*charable, 0),
    }
}

func (c *ColorToChar) Reset() {
    c.CharOrder = make([]*charable, 0)
}

func (c *ColorToChar) Map(frame ansiparser.Frame) {
    for i, ch := range frame.Chars {

        char := string(ch)
        bright := RGBBrightness(frame.Color[i])

        var color *colorable = nil
        var charPtr *charable = nil
        for _, b := range c.BrightnessOrder {
            if b.brightness == bright {
                color = b
                break
            }
        }

        for _, c := range c.CharOrder {
            if c.char == char {
                charPtr = c
                break
            }
        }

        if color == nil {
            color = &colorable{
                brightness: bright,
                chars: []string{char},
                rgb: frame.Color[i],
            }
            c.BrightnessOrder = append(c.BrightnessOrder, color)
        }

        if charPtr == nil {
            charPtr = &charable{
                char: char,
                brightness: make([]float64, 0),
                count: 0,
                rgb: frame.Color[i],
            }

            c.CharOrder = append(c.CharOrder, charPtr)
        }

        if !slices.Contains(color.chars, char) {
            color.chars = append(color.chars, char)
        }

        color.count++
        charPtr.brightness = append(charPtr.brightness, bright)
        charPtr.count++
    }
}

func (c *ColorToChar) String() string {
    slices.SortFunc(c.BrightnessOrder, func(a, b *colorable) int {
        return int(b.brightness) - int(a.brightness)
    })

    out := ""
    for _, color := range c.BrightnessOrder {
        out += fmt.Sprintf("%s\n", color)
    }

    out += "-------chars----------\n"
    for _, char := range c.CharOrder {
        out += fmt.Sprintf("%s\n", char)
    }
    return out
}

func NewLumin(count int) *Lumin {
    return &Lumin{
        buffer: make([]byte, count, count),
        idx: 0,
    }
}


