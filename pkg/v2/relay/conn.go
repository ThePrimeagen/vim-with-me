package relay

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

type Conn struct {
	conn *websocket.Conn
	id   int

	msgs  chan []byte
	relay *Relay

	authorized bool
}

func (c *Conn) read() {
	c.authorized = false
	for {
		mt, message, err := c.conn.ReadMessage()
		if mt != websocket.BinaryMessage {
			break
		}

		if err != nil {
			break
		}

		if !c.authorized {
			if c.relay.uuid == string(message) {
				c.authorized = true
				continue
			} else {
				break
			}
		}

		// relay this message to everyone else
		c.relay.relay(message)
	}

	c.relay.remove(c.id)
    c.conn.Close()
}

func (c *Conn) write() {
	for {
		msg := <-c.msgs

        if c.authorized {
            continue
        }

		slog.Warn("writing message to client", "id", c.id)
		err := c.conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			break
		}
	}

	c.relay.remove(c.id)
    c.conn.Close()
}

func (c *Conn) msg(msg []byte) {
	select {
	case c.msgs <- msg:
	default:
		break
	}
}
