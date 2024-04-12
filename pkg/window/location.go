package window

import "fmt"

const LOCATION_BYTE_LENGTH = 2
type Location struct {
    Row byte
    Col byte
}

func NewLocation(r, c byte) Location {
    return Location{Row: r, Col: c}
}

func (l *Location) MarshalBinary() (data []byte, err error) {
	b := make([]byte, 0, 2)
	b = append(b, l.Row)
	return append(b, l.Col), nil
}

func (l *Location) UnmarshalBinary(bytes []byte) error {
	if len(bytes) < 2 {
        return fmt.Errorf("expects at least 2 bytes: %d", len(bytes))
	}

	l.Row = bytes[0]
    l.Col = bytes[1]

    return nil
}
