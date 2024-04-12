package window

type Color struct {
	red        byte
	blue       byte
	green      byte
	foreground bool
}

func NewColor(r, b, g byte, f bool) Color {
    return Color{
        red: r,
        blue: b,
        green: g,
        foreground: f,
    }
}

/*
func (c *Color) MarshalBinary() (data []byte, err error) {
	b := make([]byte, 0, )
	b = append(b, VERSION)
	b = append(b, t.Command)
	b = append(b, lengthData...)
	return append(b, t.Data...), nil
}

func (t *TCPCommand) UnmarshalBinary(bytes []byte) error {
	if bytes[0] != VERSION {
		return fmt.Errorf("version mismatch %d != %d", bytes[0], VERSION)
	}

	length := int(binary.BigEndian.Uint16(bytes[2:]))
	end := HEADER_SIZE + length

	if len(bytes) < end {
		return fmt.Errorf("not enough data to parse packet: got %d expected %d", len(bytes), HEADER_SIZE+length)
	}

	command := bytes[1]
	data := bytes[HEADER_SIZE:end]

	t.Command = command
	t.Data = data

	return nil
}
*/
