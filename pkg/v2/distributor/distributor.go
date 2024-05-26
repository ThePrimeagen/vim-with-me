package distributor

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	"log"
	"log/slog"
	"net/http"
	"sync"
)

type Distributor struct {
	sync.Mutex

	listenPort         int
	authId             string
	upstreamConnection *websocket.Conn
	downstreams        []string
	downstreamConns    []*Downstream
	msgChan            chan []byte
}

func NewDistributor(listenPort int, authId string, downstreams []string) *Distributor {
	return &Distributor{
		listenPort:  listenPort,
		authId:      authId,
		downstreams: downstreams,
	}
}

func (d *Distributor) Run() {
	for _, addr := range d.downstreams {
		d.downstreamConns = append(d.downstreamConns, NewDownstream(addr))
	}

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

		for _, downstream := range d.downstreamConns {
			downstream.SendMessage(mt, msg)
		}
	}

	_ = d.upstreamConnection.Close()
	d.upstreamConnection = nil
}
