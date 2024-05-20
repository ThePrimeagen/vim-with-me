package window

import (
	"fmt"

	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

const LOCATION_ENCODING_LENGTH = 2

type Location struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func NewLocation(r, c int) Location {
	assert.Assert(r < 256, fmt.Sprintf("cannot exceed 256 for rows: %d", r))
	assert.Assert(c < 256, fmt.Sprintf("cannot exceed 256 for cols: %d", c))

	return Location{Row: int(r), Col: int(c)}
}

func (l *Location) ToRowCol() (int, int) {
	return l.Row, l.Col
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
