package tcp

import (
	"log/slog"
	"net"
)

type Connection struct {
	conn   net.Conn
	closed bool
	id     int
}

var id int = 0

func NewConnection(conn net.Conn) Connection {
	id++
    return Connection{closed: false, id: id, conn: conn}
}

func (c *Connection) Close() {
    if c.closed {
        return
    }

    c.closed = true
    _ = c.conn.Close()
}

// TODO: I need to handle n if n is less than bytes length
// This will be a serious issue potentially.  once we are out of sync the
// connection may crash
func send(conns []*Connection, cmd *TCPCommand) {
	removals := make([]int, 0)
	msg := cmd.Bytes()

	for i, conn := range conns {
		slog.Info("sending message", "index", i, "msg", msg)
		if conn.closed {
			slog.Error("removing due to close", "index", i)
			removals = append(removals, i)
            continue
		}

		_, err := conn.conn.Write(msg)
		if err != nil {
			slog.Error("removing due to error", "index", i)
			removals = append(removals, i)
		}
	}

	// TODO: on airplane, can i reverse iterate?
	for i := len(removals) - 1; i >= 0; i-- {
		idx := removals[i]
		conns = append(conns[:idx], conns[idx+1:]...)
	}
}

func send_cmds(conn *Connection, cmds []*TCPCommand) error {
	for _, cmd := range cmds {
		_, err := conn.conn.Write(cmd.Bytes())
		if err != nil {
			// TODO: Do i need to close the connection?
			return err
		}
	}

	return nil
}
