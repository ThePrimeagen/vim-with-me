package tcp

import (
	"io"
	"log"
)

type Connection struct {
	conn io.Writer
	id   int
}

var id int = 0

func NewConnection(writer io.Writer) Connection {
	id++
	return Connection{id: id, conn: writer}
}

// TODO: I need to handle n if n is less than bytes length
// This will be a serious issue potentially.  once we are out of sync the
// connection may crash
func send(conns []Connection, cmd *TCPCommand) []Connection {
    removals := make([]int, 0)
    removed := make([]Connection, 0)
    msg := cmd.Bytes()

    for i, conn := range conns {
        log.Printf("sending to %d\n", i)
        _, err := conn.conn.Write(msg)
        if err != nil {
            log.Printf("removing due to close: %d\n", i)
            removals = append(removals, i)
        }
    }

    // TODO: on airplane, can i reverse iterate?
    for i := len(removals) - 1; i >= 0; i-- {
        removed = append(removed, conns[i])
        conns = append(conns[:i], conns[i + 1:]...)
    }
    return removed
}

func send_cmds(conn Connection, cmds []*TCPCommand) error {
    for _, cmd := range cmds {
        _, err := conn.conn.Write(cmd.Bytes())
        if err != nil {
            // TODO: Do i need to close the connection?
            return err
        }
    }

    return nil
}
