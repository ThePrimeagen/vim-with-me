package relay

import (
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

type RelayDriver struct {
	url  url.URL
	uuid string
	conn *websocket.Conn
}

func NewRelayDriver(host, path, uuid string) *RelayDriver {

	u := url.URL{Scheme: "ws", Host: host, Path: path}
	return &RelayDriver{
		url:  u,
		uuid: uuid,
	}
}

func (r *RelayDriver) Connect() error {
	assert.Assert(r.conn == nil, "attempting to connect while connected")

	c, _, err := websocket.DefaultDialer.Dial(r.url.String(), nil)
	if err != nil {
		return err
	}

	r.conn = c
	return c.WriteMessage(websocket.BinaryMessage, []byte(r.uuid))
}

func (r *RelayDriver) Relay(data []byte) error {
	return r.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (r *RelayDriver) Close() {
	assert.Assert(r.conn != nil, "attempting to close a nil connection")
	r.conn.Close()
}
