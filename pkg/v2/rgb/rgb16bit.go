package rgb

import (
	"fmt"

	"github.com/leaanthony/go-ansi-parser"
	color "gitlab.com/ethanbakerdev/colors"
)

const redMask16 = 0xFC_00
const greenMask16 = 0x03_E0
const blueMask16 = 0x00_1F
const _16 = 256 * 256

func RGBTo16BitColor(hex *ansi.Rgb) uint {
	// TODO: Gray coloring needs to be considered

	red := uint(hex.R) * 64 / _16
	green := uint(hex.G) * 32 / _16
	blue := uint(hex.B) * 32 / _16

	return (red << 10) | (green << 5) | blue
}

func RGB16BitBrightness(rgb byte) float64 {
	red := float64((rgb & redMask) >> 10)
	green := float64((rgb & greenMask) >> 5)
	blue := float64(rgb & blueMask)

	return red/63.0 + green/31.0 + blue/31.0
}

func RGB16BitToString(rgb int) string {
	red := (rgb & redMask) >> 10
	green := (rgb & greenMask) >> 5
	blue := rgb & blueMask

	return fmt.Sprintf("%02x%02x%02x", red, green, blue)
}

func RGB16BitToAnsiRGB(rgb int) color.RGB {
	red := float64((rgb & redMask) >> 10)
	green := float64((rgb & greenMask) >> 5)
	blue := float64(rgb & blueMask)

	return color.RGB{
		R: int(red / 63.0 * _16),
		G: int(green / 31.0 * _16),
		B: int(blue / 31.0 * _16),
	}
}

type rgb16BitReader struct{}

func newRGB16BitReader() *rgb16BitReader {
	return &rgb16BitReader{}
}

// Note i am not using binary as it requires me to make slices
// I already know i make TONS of unneeded garbage, this one is simple to avoid
func (r *rgb16BitReader) read(buffer []byte, offset int) (int, int) {
	hi := buffer[offset]
	lo := buffer[offset+1]
	return int(hi)<<8 | int(lo), offset + 2
}

func (r *rgb16BitReader) write(buffer []byte, offset int, color *ansi.Rgb) int {
	val := RGBTo16BitColor(color)
	hi := byte((val & 0xFF_00) >> 8)
	lo := byte(val & 0xFF)
	buffer[offset] = hi
	buffer[offset+1] = lo

	return offset + 1
}
