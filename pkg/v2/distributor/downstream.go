package distributor

import (
	"github.com/gorilla/websocket"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/metrics"
	"log/slog"
	"net/url"
	"fmt"
)

type Downstream struct {
	addr  string
	conn  *websocket.Conn
	stats *metrics.Metrics
}

func NewDownstream(addr string, stats *metrics.Metrics) *Downstream {
	ds := &Downstream{
		addr:  addr,
		stats: stats,
	}
	go ds.Run()
	return ds
}

func (ds *Downstream) Run() {
	var err error
	for {
		slog.Info("Connecting to downstream server", "addr", ds.addr)
		u := url.URL{Scheme: "ws", Host: ds.addr, Path: "/ws"}
		ds.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		assert.NoError(err, "unable to connect to downstream server "+ds.addr)

		for {
			// Discard any incoming messages
			mt, msg, err := ds.conn.ReadMessage()
			if err != nil || mt != websocket.BinaryMessage {
				// Reconnect
				slog.Warn("Downstream connection closed", "addr", ds.addr, "err", err)
				break
			}

			ds.stats.Add(ds.makeMetricName(downstreamBytesReceivedMetricName), len(msg))
		}

		slog.Error("Downstream connection closed, reconnecting", "addr", ds.addr)

		_ = ds.conn.Close()
		ds.conn = nil
	}
}

func (ds *Downstream) SendMessage(msgType int, msg []byte) {
	err := ds.conn.WriteMessage(msgType, msg)
	if err != nil {
		// Reconnect
		slog.Warn("Failed to send to downstream, closing", "addr", ds.addr, "err", err)
		_ = ds.conn.Close()
		ds.conn = nil
		return
	}

	ds.stats.Add(ds.makeMetricName(downstreamBytesSentMetricName), len(msg))
}

func (ds *Downstream) makeMetricName(name string) string {
	return fmt.Sprintf("%s{addr=\"%s\"}", name, ds.addr)
}
