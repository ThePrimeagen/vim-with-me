package color

import (
	"fmt"
	"slices"
	"strings"

	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
)

type colorable struct {
    brightness float64
    rgb byte
    count int
    chars []string
}

type charable struct {
    brightness []float64
    rgb byte
    count int
    char string
}


func (c *colorable) String() string {
    return fmt.Sprintf("0x%s: %f %s %d", ansiparser.RGBToString(c.rgb), c.brightness, strings.Join(c.chars, ", "), c.count)
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
    BrightnessOrder []*colorable
    CharOrder []*charable
}

func NewColorToChar() *ColorToChar {
    return &ColorToChar{
        BrightnessOrder: make([]*colorable, 0),
    }
}

func (c *ColorToChar) Reset() {
    c.BrightnessOrder = make([]*colorable, 0)
}

func (c *ColorToChar) Map(frame ansiparser.Frame) {
    for i, ch := range frame.Chars {

        char := string(ch)
        bright := ansiparser.RGBBrightness(frame.Color[i])

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

