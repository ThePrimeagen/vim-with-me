package display

import (
	"fmt"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/v2/encoding"
	colors "gitlab.com/ethanbakerdev/colors"
	"github.com/leaanthony/go-ansi-parser"
)

type Frame struct {
    Idx int
	Color []byte
	Chars []byte
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
            idx := r * cols + c
            color := frame.Color[idx]
            char := frame.Chars[idx]

            if color != prev && len(str) > 0 {
                row += colors.AnsiRGB(encoding.RGBByteToAnsiRGB(prev)) + str
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
    color := encoding.RGBTo8BitColor(&style.FgCol.Rgb)
    colorStr := colors.AnsiRGB(encoding.RGBByteToAnsiRGB(byte(color)))

    fmt.Printf("%+v %+v vs %+v \n", style, style.FgCol, colorStr + style.Label)

}
