package display

import (
	"fmt"
	"strings"

	ansiparser "github.com/theprimeagen/vim-with-me/pkg/v2/ansi_parser"
	"github.com/theprimeagen/vim-with-me/pkg/v2/encoding"
	colors "gitlab.com/ethanbakerdev/colors"
)

func Display(frame *ansiparser.Frame, rows, cols int) string {
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
                row += colors.AnsiRGB(encoding.RGBByteToAnsiRGB(color)) + str
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

        fmt.Printf("cols: %d vs %d\n", len(row), cols)
        out = append(out, row)
    }


    return strings.Join(out, "\r\n")
}

