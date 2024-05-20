package display

import (
	"fmt"
	"strings"

	"github.com/leaanthony/go-ansi-parser"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
	"github.com/theprimeagen/vim-with-me/pkg/v2/rgb"
	colors "gitlab.com/ethanbakerdev/colors"
)

type Frame struct {
	Idx   int
	Color []byte
	Chars []byte
}

func (f *Frame) Color16BitIterator() byteutils.ByteIterator {
	return byteutils.New16BitIterator(f.Color)
}

func (f *Frame) Color8BitIterator() byteutils.ByteIterator {
	return byteutils.New8BitIterator(f.Color)
}

func Display(frame *Frame, rows, cols int) string {
	// TODO: I know this could be better... but i don't know if this will ever
	// be a perf area.  This is meant for testing and the mirror clients only

	out := []string{}
	prev := byte(255)
	str := ""

	for r := range rows {
		row := ""
		for c := range cols {
			idx := r*cols + c
			color := frame.Color[idx]
			char := frame.Chars[idx]

			if color != prev && len(str) > 0 {
				row += colors.AnsiRGB(rgb.RGBByteToAnsiRGB(prev)) + str
				//row += str
				prev = color
				str = ""
			}

			str += string(char)
		}

		if len(str) > 0 {
			row += str
			str = ""
		}

		out = append(out, row)
	}

	return strings.Join(out, "\r\n")
}

func Clear() string {
	return "\033[H\033[2J"
}

func DebugStyle(style *ansi.StyledText) {
	color := rgb.RGBTo8BitColor(&style.FgCol.Rgb)
	colorStr := colors.AnsiRGB(rgb.RGBByteToAnsiRGB(byte(color)))

	fmt.Printf("%+v %+v vs %+v \n", style, style.FgCol, colorStr+style.Label)

}
