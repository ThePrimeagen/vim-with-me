package relay

import "github.com/gorilla/websocket"

type Conn struct {
	conn *websocket.Conn
	id   int

	msgs  chan []byte
	relay *Relay
}

func (c *Conn) read() {
	authorized := false
	for {
		mt, message, err := c.conn.ReadMessage()
		if mt != websocket.BinaryMessage {
			break
		}

		if err != nil {
			break
		}

		if !authorized {
			if c.relay.uuid == string(message) {
				authorized = true
				continue
			} else {
				break
			}
		}

		// relay this message to everyone else
		c.relay.relay(message)
	}

	c.relay.remove(c.id)
}

func (c *Conn) write() {
	for {
		msg := <-c.msgs
		err := c.conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			break
		}
	}

	c.relay.remove(c.id)
}

func (c *Conn) msg(msg []byte) {
	select {
	case c.msgs <- msg:
	default:
		break
	}
}
