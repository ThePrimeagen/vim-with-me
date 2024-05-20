package relay

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

type Conn struct {
	Conn *websocket.Conn
	id   int

	msgs  chan []byte
	relay *Relay

	authorized bool
}

func (c *Conn) read() {
	c.authorized = false
	for {
		mt, message, err := c.Conn.ReadMessage()
		if mt != websocket.BinaryMessage {
			break
		}

		if err != nil {
			break
		}

		slog.Debug("connection message received", "msg", message, "authorized", c.authorized, "id", c.id)
		if !c.authorized {
			if c.relay.uuid == string(message) {
				c.authorized = true
				continue
			} else {
				slog.Error("unauthorized message :: destroying connection", "id", c.id)
				break
			}
		}

		// relay this message to everyone else
		c.relay.relay(message)
	}

	slog.Warn("closing down connection", "id", c.id)
	c.relay.remove(c.id)
	c.Conn.Close()
}

func (c *Conn) write() {
	for {
		msg := <-c.msgs

		if c.authorized {
			continue
		}

		slog.Debug("writing message to client", "id", c.id)
		err := c.Conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			break
		}
	}

	c.relay.remove(c.id)
	c.Conn.Close()
}

func (c *Conn) msg(msg []byte) {
	select {
	case c.msgs <- msg:
	default:
		break
	}
}
