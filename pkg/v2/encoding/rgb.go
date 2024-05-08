package encoding

import (
	"fmt"

	"github.com/leaanthony/go-ansi-parser"
	color "gitlab.com/ethanbakerdev/colors"
)

const redMask = 0b111_00_000
const greenMask = 0b000_11_000
const blueMask = 0b000_00_111

func RGBTo8BitColor(hex ansi.Rgb) uint {
	red := uint(hex.R) * 8 / 256
	green := uint(hex.G) * 4 / 256
	blue := uint(hex.B) * 8 / 256

	return (red << 5) | (green << 3) | blue
}

func RGBBrightness(rgb byte) float64 {
	red := float64((rgb & redMask) >> 5)
	green := float64((rgb & greenMask) >> 3)
	blue := float64(rgb & blueMask)

	return red/7.0 + green/3.0 + blue/7.0
}

func RGBToString(rgb byte) string {
	red := (rgb & redMask) >> 5
	green := (rgb & greenMask) >> 3
	blue := rgb & blueMask

	return fmt.Sprintf("%02x%02x%02x", red, green, blue)
}

func RGBByteToAnsiRGB(rgb byte) color.RGB {
	red := float64((rgb & redMask) >> 5)
	green := float64((rgb & greenMask) >> 3)
	blue := float64(rgb & blueMask)

	return color.RGB{
        R: int(red / 7.0 * 255),
        G: int(green / 3.0 * 255),
        B: int(blue / 7.0 * 255),
    }
}

