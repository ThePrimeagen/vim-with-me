package window

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

const LOCATION_BYTE_LENGTH = 2
type Location struct {
    Row int
    Col int
}

func NewLocation(r, c int) Location {
    assert.Assert(r < 256, "cannot exceed 256 for rows")
    assert.Assert(c < 256, "cannot exceed 256 for cols")

    return Location{Row: int(r), Col: int(c)}
}

func (l *Location) MarshalBinary() (data []byte, err error) {
	b := make([]byte, 0, 2)
	b = append(b, byte(l.Row))
	return append(b, byte(l.Col)), nil
}

func (l *Location) UnmarshalBinary(bytes []byte) error {
	if len(bytes) < 2 {
        return fmt.Errorf("expects at least 2 bytes: %d", len(bytes))
	}

	l.Row = int(bytes[0])
    l.Col = int(bytes[1])

    return nil
}
