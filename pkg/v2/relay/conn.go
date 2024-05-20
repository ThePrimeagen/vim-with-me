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

        slog.Warn("connection message received", "msg", string(message), "authorized", c.authorized, "id", c.id)
		if !c.authorized {
			if c.relay.uuid == string(message) {
				c.authorized = true
				continue
			} else {
                slog.Warn("unauthorized message :: destroying connection", "id", c.id)
				break
			}
		}

		// relay this message to everyone else
		c.relay.relay(message)
	}

    slog.Warn("closing down connection", "id", c.id)
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
