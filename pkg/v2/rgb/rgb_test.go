package rgb_test

import (
	"fmt"
	"testing"

	"github.com/leaanthony/go-ansi-parser"
	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/v2/rgb"
)

func TestAnsiStyleRGBConversion(t *testing.T) {
	ansiStr := "[38;2;255;255;001mQQQQQQQQQQQQQQ"
	style, err := ansi.Parse(ansiStr)

	require.NoError(t, err)

	_8bit := rgb.RGBTo8BitColor(&style[0].FgCol.Rgb)
	from8Bit := rgb.RGBByteToAnsiRGB(byte(_8bit))
	styleRGB := style[0].FgCol.Rgb

	require.Equal(t, uint8(from8Bit.R), styleRGB.R)
	require.Equal(t, uint8(from8Bit.G), styleRGB.G)

	// blue should be wiped out because its value is 1
	require.Equal(t, uint8(from8Bit.B), uint8(0))
}

func TestAnsiStyleRGBConversionLarger(t *testing.T) {
	ansiStr := "[38;2;155;092;020m111111[38;2;175;124;032mtt[38;2;159;131;100mnn[38;2;255;255;116mhh[38;2;175;124;032mtt[38;2;096;076;056m]][38;2;064;044;028mii[38;2;088;068;052m??[38;2;024;016;008m^^"
	style, err := ansi.Parse(ansiStr)

	require.NoError(t, err)

	fmt.Printf("styles: %+v\n", style)
}
