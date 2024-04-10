package tcp2

import (
	"encoding/binary"
	"fmt"
)

var VERSION byte = 1
var HEADER_SIZE = 4

type TCPCommand struct {
	Command byte
	Data    []byte
}

func (t *TCPCommand) MarshalBinary() (data []byte, err error) {
    length := uint16(len(t.Data))
    lengthData := make([]byte, 2)
    binary.BigEndian.PutUint16(lengthData, length)

    b := make([]byte, 0, 1 + 1 + 2 + length)
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
        return fmt.Errorf("not enough data to parse packet: got %d expected %d", len(bytes), HEADER_SIZE + length)
    }

    command := bytes[1]
    data := bytes[HEADER_SIZE:end]

    t.Command = command
    t.Data = data

    return nil
}



