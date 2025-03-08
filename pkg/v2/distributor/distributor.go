package distributor

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"github.com/theprimeagen/vim-with-me/pkg/v2/metrics"
	"log"
	"log/slog"
	"net/http"
	"sync"
)

const (
	upstreamConnectedMetricName       = "distributor_upstream_connected"
	upstreamBytesReceivedMetricName   = "distributor_upstream_bytes_received"
	upstreamBytesSentMetricName       = "distributor_upstream_bytes_sent"
)

type Distributor struct {
	sync.Mutex

	listenPort         int
	authId             string
	upstreamConnection *websocket.Conn
	downstreams        []string
	downstreamConns    []*Downstream
	msgChan            chan []byte
	stats              *metrics.Metrics
}

func NewDistributor(listenPort int, authId string, downstreams []string) *Distributor {
	stats := metrics.New()
	stats.Set(upstreamConnectedMetricName, 0)
	stats.Set(upstreamBytesReceivedMetricName, 0)
	stats.Set(upstreamBytesSentMetricName, 0)

	return &Distributor{
		listenPort:  listenPort,
		authId:      authId,
		downstreams: downstreams,
		stats:       stats,
	}
}

func (d *Distributor) Run() {
	for _, addr := range d.downstreams {
		d.downstreamConns = append(d.downstreamConns, NewDownstream(addr, d.stats))
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		d.stats.WritePrometheus(w)
	})

	http.HandleFunc("/ws", d.handleIncomingConnection)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Location", "https://vim-with.me/")
		w.WriteHeader(http.StatusTemporaryRedirect)
		_, _ = w.Write([]byte("Not here, <a href=\"https://vim-with.me/\">over here</a>!"))
	})

	addr := fmt.Sprintf("0.0.0.0:%d", d.listenPort)
	slog.Warn("listening and serving http", "http", addr)
	err := http.ListenAndServe(addr, nil)

	log.Fatal(err)
}

func (d *Distributor) handleIncomingConnection(w http.ResponseWriter, r *http.Request) {
	func() {
		d.Lock()
		defer d.Unlock()

		if d.upstreamConnection != nil {
			// One connection at a time!
			slog.Warn("Rejected connection, already have one",
				"remote", r.RemoteAddr)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		slog.Info("New upstream connection", "remote", r.RemoteAddr)

		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		assert.NoError(err, "unable to upgrade connection")

		d.upstreamConnection = c

		d.stats.Set(upstreamConnectedMetricName, 1)
	}()

	for {
		mt, msg, err := d.upstreamConnection.ReadMessage()
		if mt != websocket.BinaryMessage {
			slog.Error("Upstream sent non-binary message, disconnecting")
			break
		}

		if err != nil {
			slog.Error("Upstream error, disconnecting", "err", err)
			break
		}

		d.stats.Add(upstreamBytesReceivedMetricName, len(msg))
		for _, downstream := range d.downstreamConns {
			downstream.SendMessage(mt, msg)
		}
	}

	d.stats.Set(upstreamConnectedMetricName, 0)
	_ = d.upstreamConnection.Close()
	d.upstreamConnection = nil
}
