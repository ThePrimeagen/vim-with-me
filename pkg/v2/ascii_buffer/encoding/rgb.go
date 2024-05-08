package encoding

import "fmt"

func RGBBrightness(rgb byte) float64 {
	red := float64((rgb & 0b111_000_00) >> 5)
	green := float64((rgb & 0b000_111_00) >> 2)
	blue := float64(rgb & 0b11)

	return red/7.0 + green/7.0 + blue/3.0
}

func RGBToString(rgb byte) string {
	red := (rgb & 0b111_000_00) >> 5
	green := (rgb & 0b000_111_00) >> 2
	blue := rgb & 0b11

	return fmt.Sprintf("%02x%02x%02x", red, green, blue)
}

