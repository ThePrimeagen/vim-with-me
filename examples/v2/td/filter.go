package td

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

func inBounds(str string, bound int) bool {
    item, err := strconv.Atoi(str)
    if err != nil {
        return false
    }

    if item >= bound || item < 0 {
        return false;
    }

    return true
}

func TDFilter(rows, cols int) func(msg string) bool {
    return func(msg string) bool {
        parts := strings.Split(msg, ",")
        if len(parts) != 2 {
            return false;
        }

        return inBounds(parts[0], rows) && inBounds(parts[1], cols)
    }
}

func TDAfterMap(input string) string {
    parts := strings.Split(input, ",")
    assert.Assert(len(parts) == 2, "expected two parts", "parts", parts)
    row, _ := strconv.Atoi(parts[0])
    col, _ := strconv.Atoi(parts[1])

    return fmt.Sprintf("%d,%d", row, col)
}

